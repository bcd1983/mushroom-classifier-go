# 🍄 Mushroom Classifier (Go Version)

[![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)](https://golang.org)
[![Fyne](https://img.shields.io/badge/Fyne-007ACC?style=for-the-badge&logo=go&logoColor=white)](https://fyne.io)
[![OpenAI](https://img.shields.io/badge/OpenAI-412991?style=for-the-badge&logo=openai&logoColor=white)](https://openai.com/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)

A desktop application for mushroom identification using OpenAI's vision models. Built with Go and Fyne for a modern, cross-platform GUI experience.

## ⚠️ Important Safety Notice

**This tool is for educational purposes only.** Never consume wild mushrooms based solely on this or any app's identification. Always consult with local mycology experts before consuming any wild mushrooms. Misidentification can lead to serious illness or death.

## ✨ Features

- 🖼️ **Image Analysis**: Upload photos of mushrooms for AI-powered identification
- 🧠 **Advanced AI**: Leverages OpenAI's GPT-4 Vision for accurate analysis
- 🎨 **Modern GUI**: Clean, intuitive interface built with Fyne - pure Go, no C dependencies
- 📊 **Detailed Results**: Get species names, confidence levels, and safety information
- 🔒 **Secure**: API credentials stored safely in environment files
- 🏗️ **Modern Go Architecture**: Well-structured codebase with idiomatic Go patterns
- 🚀 **Easy Deployment**: Single binary with no external dependencies

## 📋 Prerequisites

- Go 1.21 or higher
- OpenAI API key with GPT-4 Vision access

That's it! Fyne handles all GUI dependencies automatically.

## 🚀 Quick Start

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/mushroom-classifier-go.git
   cd mushroom-classifier-go
   ```

2. **Set up your environment**
   ```bash
   cp .env.example .env
   # Edit .env with your OpenAI API key
   ```

3. **Build the application**
   ```bash
   make build
   ```

4. **Run the application**
   ```bash
   make run
   ```

## 🛠️ Building from Source

### Standard Build
```bash
make clean
make build
```

### Build Specific Targets
```bash
make mushroom-classifier  # Build only the main application
make test-api            # Build only the API test utility
```

### Cross-Platform Builds
```bash
make build-linux    # Build for Linux
make build-windows  # Build for Windows
make build-darwin   # Build for macOS (Intel and Apple Silicon)
make build-all      # Build for all platforms
```

## 📁 Project Structure

```
mushroom-classifier-go/
├── main.go                 # Main application entry point
├── base64/                 # Base64 encoding utilities
│   └── base64.go
├── config/                 # Configuration management
│   └── config.go
├── httpclient/            # HTTP client utilities
│   └── httpclient.go
├── openai/                # OpenAI API integration
│   └── openai.go
├── gui/                   # GTK+ GUI implementation
│   └── gui.go
├── cmd/                   # Command line tools
│   └── test-api/         # API testing utility
│       └── main.go
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── Makefile              # Build configuration
├── .env.example          # Example environment file
├── LICENSE               # MIT License
└── README.md             # This file
```

## 🔧 Configuration

Create a `.env` file in the project root with your OpenAI credentials:

```env
OPENAI_API_KEY=your-api-key-here
OPENAI_API_URL=https://api.openai.com/v1/chat/completions
```

## 📖 Usage

1. **Launch the application**
   ```bash
   ./build/mushroom-classifier
   ```

2. **Select an image**
   - Click "Select Image" to choose a mushroom photo
   - Supported formats: JPEG, PNG

3. **Classify the mushroom**
   - Click "Classify Mushroom" to analyze the image
   - Wait for the AI to process and return results

4. **Review the results**
   - Species identification (common and scientific names)
   - Confidence level
   - Key identifying features
   - Edibility status
   - Safety warnings
   - Similar species to be aware of

## 🧪 Testing

### API Connection Test
```bash
make test
```

This will verify your OpenAI API credentials and connection.

## 🏗️ Architecture

The application follows Go best practices with a modular architecture:

- **Config Package**: Handles environment variables and application settings
- **Base64 Package**: Provides image encoding functionality
- **HTTPClient Package**: Manages API communications
- **OpenAI Package**: Interfaces with OpenAI's vision models
- **GUI Package**: Implements the Fyne-based user interface
- **Main Package**: Orchestrates the application lifecycle

### Key Design Principles

- **Package-Oriented Design**: Each package has a single, well-defined responsibility
- **Error Handling**: Comprehensive error checking with Go's error patterns
- **Concurrency**: Proper use of goroutines for non-blocking UI operations
- **Security**: API keys stored in environment files, never hardcoded
- **Documentation**: Extensive godoc comments throughout
- **Pure Go**: No CGO or C dependencies thanks to Fyne

## 🤝 Contributing

We welcome contributions! Please follow these guidelines:

### Development Setup

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Code Style

- Follow standard Go formatting (`go fmt`)
- Use meaningful variable and function names
- Add godoc comments for exported functions
- Ensure all tests pass
- Run `make lint` before submitting

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- OpenAI for providing the GPT-4 Vision API
- Fyne team for the excellent pure Go GUI framework
- The mycology community for mushroom identification resources

## 🐛 Known Issues

- Large images (>10MB) may take longer to process
- Some rare mushroom species may not be accurately identified
- Requires active internet connection for API calls

## 🚀 Future Enhancements

- [ ] Offline mode with local AI models
- [ ] Batch processing for multiple images
- [ ] Export results to PDF/CSV
- [ ] Integration with mushroom databases
- [ ] Mobile companion app
- [ ] Support for more image formats
- [ ] Caching for repeated classifications

---

Made with ❤️ by the Mushroom Classifier Contributors