// Package base64 provides Base64 encoding utilities for image processing
package base64

import (
	"encoding/base64"
	"fmt"
	"os"
)

// EncodeData encodes binary data to Base64 string
//
// Converts binary data to Base64 encoded string following RFC 4648.
// The output is padded with '=' characters as needed.
func EncodeData(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// ReadImageToBase64 reads an image file and encodes it as Base64
//
// Opens the specified image file in binary mode, reads its entire contents,
// and returns a Base64 encoded representation. This is commonly used for
// embedding images in JSON requests to vision APIs.
func ReadImageToBase64(filename string) (string, error) {
	// Read the entire file
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filename, err)
	}

	// Check if file is empty
	if len(data) == 0 {
		return "", fmt.Errorf("file %s is empty", filename)
	}

	// Encode to base64
	encoded := EncodeData(data)
	return encoded, nil
}