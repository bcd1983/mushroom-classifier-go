// Package gui provides GTK+ graphical user interface
package gui

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/mushroom-classifier/mushroom-classifier-go/base64"
	"github.com/mushroom-classifier/mushroom-classifier-go/config"
	"github.com/mushroom-classifier/mushroom-classifier-go/openai"
)

// App contains all GUI widgets and application state
type App struct {
	// Main application window
	Window *gtk.Window

	// Image display widget for showing the selected mushroom photo
	ImageView *gtk.Image

	// Button to trigger file selection dialog
	UploadButton *gtk.Button

	// Button to start classification process
	ClassifyButton *gtk.Button

	// Text view for displaying classification results
	ResultView *gtk.TextView

	// Label showing current status/progress
	StatusLabel *gtk.Label

	// Path to the currently loaded image file
	ImagePath string

	// Base64 encoded image data
	Base64Image string

	// Application configuration (API keys, etc.)
	Config *config.Config
}

// NewApp creates a new App instance with initialized GTK widgets
func NewApp(cfg *config.Config) (*App, error) {
	app := &App{
		Config: cfg,
	}

	// Create main window
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		return nil, fmt.Errorf("unable to create window: %w", err)
	}
	app.Window = win

	// Set window properties
	win.SetTitle("Mushroom Classifier")
	win.SetDefaultSize(800, 600)
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create main container
	vbox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	if err != nil {
		return nil, fmt.Errorf("unable to create box: %w", err)
	}
	vbox.SetMarginTop(10)
	vbox.SetMarginBottom(10)
	vbox.SetMarginStart(10)
	vbox.SetMarginEnd(10)
	win.Add(vbox)

	// Create header label
	headerLabel, err := gtk.LabelNew("")
	if err != nil {
		return nil, fmt.Errorf("unable to create header label: %w", err)
	}
	headerLabel.SetMarkup("<span size='x-large' weight='bold'>Mushroom Classifier</span>")
	vbox.PackStart(headerLabel, false, false, 0)

	// Create image view
	imageView, err := gtk.ImageNew()
	if err != nil {
		return nil, fmt.Errorf("unable to create image: %w", err)
	}
	imageView.SetSizeRequest(400, 300)
	app.ImageView = imageView

	// Create scrolled window for image
	scrolledImage, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create scrolled window: %w", err)
	}
	scrolledImage.Add(imageView)
	scrolledImage.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	vbox.PackStart(scrolledImage, true, true, 0)

	// Create button box
	buttonBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	if err != nil {
		return nil, fmt.Errorf("unable to create button box: %w", err)
	}
	vbox.PackStart(buttonBox, false, false, 0)

	// Create upload button
	uploadButton, err := gtk.ButtonNewWithLabel("Select Image")
	if err != nil {
		return nil, fmt.Errorf("unable to create upload button: %w", err)
	}
	uploadButton.Connect("clicked", app.onUploadClicked)
	app.UploadButton = uploadButton
	buttonBox.PackStart(uploadButton, true, true, 0)

	// Create classify button
	classifyButton, err := gtk.ButtonNewWithLabel("Classify Mushroom")
	if err != nil {
		return nil, fmt.Errorf("unable to create classify button: %w", err)
	}
	classifyButton.SetSensitive(false)
	classifyButton.Connect("clicked", app.onClassifyClicked)
	app.ClassifyButton = classifyButton
	buttonBox.PackStart(classifyButton, true, true, 0)

	// Create status label
	statusLabel, err := gtk.LabelNew("Select an image to begin")
	if err != nil {
		return nil, fmt.Errorf("unable to create status label: %w", err)
	}
	app.StatusLabel = statusLabel
	vbox.PackStart(statusLabel, false, false, 0)

	// Create results label
	resultsLabel, err := gtk.LabelNew("")
	if err != nil {
		return nil, fmt.Errorf("unable to create results label: %w", err)
	}
	resultsLabel.SetMarkup("<b>Results:</b>")
	vbox.PackStart(resultsLabel, false, false, 0)

	// Create text view for results
	textView, err := gtk.TextViewNew()
	if err != nil {
		return nil, fmt.Errorf("unable to create text view: %w", err)
	}
	textView.SetWrapMode(gtk.WRAP_WORD)
	textView.SetEditable(false)
	app.ResultView = textView

	// Create scrolled window for text view
	scrolledText, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("unable to create scrolled window: %w", err)
	}
	scrolledText.Add(textView)
	scrolledText.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledText.SetSizeRequest(-1, 200)
	vbox.PackStart(scrolledText, true, true, 0)

	return app, nil
}

// Run starts the GTK main event loop
func (app *App) Run() {
	app.Window.ShowAll()
	gtk.Main()
}

// onUploadClicked handles the upload button click event
func (app *App) onUploadClicked() {
	// Create file chooser dialog
	dialog, err := gtk.FileChooserDialogNewWith2Buttons(
		"Select Image",
		app.Window,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		"Cancel", gtk.RESPONSE_CANCEL,
		"Open", gtk.RESPONSE_ACCEPT)
	if err != nil {
		app.showError("Failed to create file dialog: %v", err)
		return
	}
	defer dialog.Destroy()

	// Add file filters
	filter, err := gtk.FileFilterNew()
	if err == nil {
		filter.SetName("Image files")
		filter.AddPattern("*.jpg")
		filter.AddPattern("*.jpeg")
		filter.AddPattern("*.png")
		filter.AddPattern("*.JPG")
		filter.AddPattern("*.JPEG")
		filter.AddPattern("*.PNG")
		dialog.AddFilter(filter)
	}

	// Run dialog
	response := dialog.Run()
	if response != gtk.RESPONSE_ACCEPT {
		return
	}

	// Get selected file
	filename := dialog.GetFilename()
	if filename == "" {
		return
	}

	// Load and display image
	if err := app.loadImage(filename); err != nil {
		app.showError("Failed to load image: %v", err)
		return
	}

	app.ImagePath = filename
	app.StatusLabel.SetText(fmt.Sprintf("Loaded: %s", filepath.Base(filename)))
	app.ClassifyButton.SetSensitive(true)
}

// onClassifyClicked handles the classify button click event
func (app *App) onClassifyClicked() {
	if app.Base64Image == "" {
		app.showError("No image loaded")
		return
	}

	// Disable buttons during processing
	app.UploadButton.SetSensitive(false)
	app.ClassifyButton.SetSensitive(false)
	app.StatusLabel.SetText("Analyzing image...")

	// Clear previous results
	buffer, err := app.ResultView.GetBuffer()
	if err == nil {
		buffer.SetText("")
	}

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
		
		// Update UI in main thread using glib.IdleAdd
		glib.IdleAdd(func() {
			// Re-enable buttons
			app.UploadButton.SetSensitive(true)
			app.ClassifyButton.SetSensitive(true)

			if err != nil {
				app.showError("Analysis failed: %v", err)
				app.StatusLabel.SetText("Analysis failed")
				return
			}

			if !resp.Success {
				app.showError("Analysis failed: %s", resp.ErrorMessage)
				app.StatusLabel.SetText("Analysis failed")
				return
			}

			// Display results
			buffer, err := app.ResultView.GetBuffer()
			if err == nil {
				buffer.SetText(resp.Content)
			}
			app.StatusLabel.SetText("Analysis complete")
		})
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

	// Load pixbuf for display
	pixbuf, err := gdk.PixbufNewFromFile(filename)
	if err != nil {
		return err
	}

	// Scale image to fit
	width := pixbuf.GetWidth()
	height := pixbuf.GetHeight()
	maxSize := 400

	if width > maxSize || height > maxSize {
		scale := float64(maxSize) / float64(max(width, height))
		newWidth := int(float64(width) * scale)
		newHeight := int(float64(height) * scale)
		
		scaled, err := pixbuf.ScaleSimple(newWidth, newHeight, gdk.INTERP_BILINEAR)
		if err == nil {
			pixbuf = scaled
		}
	}

	// Set image
	app.ImageView.SetFromPixbuf(pixbuf)
	return nil
}

// showError displays an error message dialog
func (app *App) showError(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	log.Printf("Error: %s", msg)
	
	dialog := gtk.MessageDialogNew(
		app.Window,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_ERROR,
		gtk.BUTTONS_OK,
		msg)
	dialog.Run()
	dialog.Destroy()
}

// max returns the maximum of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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