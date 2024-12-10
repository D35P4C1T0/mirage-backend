package utils

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"image"
	"io"
	"mime/multipart"
	"net/http"
)

// RetrieveImageFromHTTPForm retrieves and validates an image file uploaded through an HTTP form.
//
// Parameters:
//   - c: The Gin context containing the HTTP request and response objects.
//   - paramName: The name of the form parameter that holds the uploaded file.
//
// Returns:
//   - A byte slice containing the file's data if the operation is successful.
//   - A boolean indicating whether an error occurred (true if there was an error, false otherwise).
//
// This function performs the following steps:
//  1. Retrieves the uploaded file from the HTTP form using the specified parameter name.
//  2. Reads the file's contents into memory as a byte slice.
//  3. Validates the file's format by decoding it as an image.
//  4. If any step fails, it sends an appropriate error response via the Gin context and returns `nil` along with `true`.
//
// Notes:
//   - The function ensures the file is closed after reading its contents.
//   - If the image format is invalid, an error response is returned with the "Invalid image format" message.
func RetrieveImageFromHTTPForm(c *gin.Context, paramName string) ([]byte, bool) {
	// Retrieve the uploaded file
	file, _, fileErr := c.Request.FormFile(paramName)
	if fileErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file", "details": fileErr.Error()})
		return nil, true
	}
	defer func(file multipart.File) {
		err := file.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to close file", "details": err.Error()})
		}
	}(file)

	//log.Printf("Received file: %s, size: %d bytes", header.Filename, header.Size)

	// Read the file into memory
	fileBytes, readErr := io.ReadAll(file)
	if readErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file", "details": readErr.Error()})
		return nil, true
	}

	// Decode and validate the image
	_, _, decodeErr := image.Decode(bytes.NewReader(fileBytes))
	if decodeErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format", "details": decodeErr.Error()})
		return nil, true
	}

	return fileBytes, false
}
