// Package gui provides Fyne-based graphical user interface
package gui

import (
	"fmt"
	"log"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/mushroom-classifier/mushroom-classifier-go/base64"
	"github.com/mushroom-classifier/mushroom-classifier-go/config"
	"github.com/mushroom-classifier/mushroom-classifier-go/openai"
)

// App contains all GUI widgets and application state
type App struct {
	// Fyne application instance
	FyneApp fyne.App

	// Main application window
	Window fyne.Window

	// Image display widget for showing the selected mushroom photo
	ImageView *canvas.Image

	// Button to trigger file selection dialog
	UploadButton *widget.Button

	// Button to start classification process
	ClassifyButton *widget.Button

	// Text widget for displaying classification results
	ResultView *widget.Entry

	// Label showing current status/progress
	StatusLabel *widget.Label

	// Path to the currently loaded image file
	ImagePath string

	// Base64 encoded image data
	Base64Image string

	// Application configuration (API keys, etc.)
	Config *config.Config
}

// NewApp creates a new App instance with initialized Fyne widgets
func NewApp(cfg *config.Config) (*App, error) {
	// Create Fyne application
	fyneApp := app.New()

	// Create main window
	window := fyneApp.NewWindow("Mushroom Classifier")
	window.Resize(fyne.NewSize(800, 600))

	app := &App{
		FyneApp: fyneApp,
		Window:  window,
		Config:  cfg,
	}

	// Create UI components
	app.createUI()

	return app, nil
}

// createUI builds the user interface
func (app *App) createUI() {
	// Create header
	headerLabel := widget.NewLabel("Mushroom Classifier")
	headerLabel.Alignment = fyne.TextAlignCenter
	headerLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Create image view
	app.ImageView = &canvas.Image{
		FillMode: canvas.ImageFillContain,
	}
	app.ImageView.SetMinSize(fyne.NewSize(400, 300))
	
	// Wrap image in a bordered container
	imageContainer := container.NewBorder(
		nil, nil, nil, nil,
		container.NewCenter(app.ImageView),
	)

	// Create buttons
	app.UploadButton = widget.NewButton("Select Image", app.onUploadClicked)
	app.ClassifyButton = widget.NewButton("Classify Mushroom", app.onClassifyClicked)
	app.ClassifyButton.Disable()

	buttonContainer := container.New(layout.NewHBoxLayout(),
		app.UploadButton,
		app.ClassifyButton,
	)

	// Create status label
	app.StatusLabel = widget.NewLabel("Select an image to begin")

	// Create results section
	resultsLabel := widget.NewLabel("Results:")
	resultsLabel.TextStyle = fyne.TextStyle{Bold: true}
	
	app.ResultView = widget.NewMultiLineEntry()
	app.ResultView.Wrapping = fyne.TextWrapWord
	app.ResultView.Disable()
	
	resultScroll := container.NewScroll(app.ResultView)
	resultScroll.SetMinSize(fyne.NewSize(0, 200))

	// Create main layout
	content := container.NewVBox(
		headerLabel,
		widget.NewSeparator(),
		imageContainer,
		buttonContainer,
		app.StatusLabel,
		widget.NewSeparator(),
		resultsLabel,
		resultScroll,
	)

	// Wrap in padded container
	paddedContent := container.NewPadded(content)
	
	app.Window.SetContent(paddedContent)
	app.Window.CenterOnScreen()
}

// Run starts the Fyne application
func (app *App) Run() {
	app.Window.ShowAndRun()
}

// onUploadClicked handles the upload button click event
func (app *App) onUploadClicked() {
	// Create file open dialog
	fileDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if err != nil {
			app.showError("Failed to open file dialog", err)
			return
		}
		if reader == nil {
			return
		}
		defer reader.Close()

		// Get file path
		filename := reader.URI().Path()
		
		// Load and display image
		if err := app.loadImage(filename); err != nil {
			app.showError("Failed to load image", err)
			return
		}

		app.ImagePath = filename
		app.StatusLabel.SetText(fmt.Sprintf("Loaded: %s", filepath.Base(filename)))
		app.ClassifyButton.Enable()
	}, app.Window)

	// Set file filter for images
	fileDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".jpeg", ".png", ".JPG", ".JPEG", ".PNG"}))
	fileDialog.Show()
}

// onClassifyClicked handles the classify button click event
func (app *App) onClassifyClicked() {
	if app.Base64Image == "" {
		app.showError("No image loaded", nil)
		return
	}

	// Disable buttons during processing
	app.UploadButton.Disable()
	app.ClassifyButton.Disable()
	app.StatusLabel.SetText("Analyzing image...")
	app.ResultView.SetText("Processing...")

	// Create OpenAI request
	req := &openai.Request{
		APIKey:      app.Config.OpenAIAPIKey,
		APIURL:      app.Config.OpenAIAPIURL,
		Model:       "gpt-4o",
		Prompt:      getMushroomPrompt(),
		Base64Image: app.Base64Image,
		MaxTokens:   1000,
	}

	// Process in background
	go func() {
		// Analyze image
		resp, err := openai.AnalyzeImage(req)

		// Update UI (Fyne is thread-safe)
		if err != nil {
			app.showError("Analysis failed", err)
			app.StatusLabel.SetText("Analysis failed")
			app.ResultView.SetText("")
		} else if !resp.Success {
			app.showError("Analysis failed", fmt.Errorf(resp.ErrorMessage))
			app.StatusLabel.SetText("Analysis failed")
			app.ResultView.SetText("")
		} else {
			app.ResultView.SetText(resp.Content)
			app.StatusLabel.SetText("Analysis complete")
		}

		// Re-enable buttons
		app.UploadButton.Enable()
		app.ClassifyButton.Enable()
	}()
}

// loadImage loads and displays an image file
func (app *App) loadImage(filename string) error {
	// Read image to base64
	base64Image, err := base64.ReadImageToBase64(filename)
	if err != nil {
		return err
	}
	app.Base64Image = base64Image

	// Load image for display
	app.ImageView.File = filename
	app.ImageView.Refresh()

	return nil
}

// showError displays an error message dialog
func (app *App) showError(message string, err error) {
	errorMsg := message
	if err != nil {
		errorMsg = fmt.Sprintf("%s: %v", message, err)
	}
	log.Printf("Error: %s", errorMsg)
	
	dialog.ShowError(fmt.Errorf(errorMsg), app.Window)
}

// getMushroomPrompt returns the prompt for mushroom analysis
func getMushroomPrompt() string {
	return `You are an expert mycologist. Analyze this image of a mushroom and provide:

1. **Species Identification**: Common name and scientific name
2. **Confidence Level**: How certain you are of the identification (High/Medium/Low)
3. **Key Identifying Features**: What visual characteristics led to this identification
4. **Edibility**: Whether this mushroom is edible, poisonous, or unknown
5. **Safety Warning**: Any important safety information
6. **Similar Species**: Other mushrooms it might be confused with

IMPORTANT: Always err on the side of caution. If uncertain, clearly state so. Never encourage consumption of wild mushrooms without expert verification.`
}