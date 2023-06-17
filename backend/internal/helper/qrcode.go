package helper

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func GenerateQRCodeMarkdown(str string) (string, error) {
	// Generate the QR code as an image
	qrCode, err := qr.Encode(str, qr.H, qr.Auto)
	if err != nil {
		return "", err
	}

	// Scale the QR code image to a reasonable size
	qrCode, err = barcode.Scale(qrCode, 400, 400)
	if err != nil {
		return "", err
	}

	// Create a buffer to hold the PNG image
	buffer := new(bytes.Buffer)

	// Encode the QR code image as PNG and write it to the buffer
	err = png.Encode(buffer, qrCode)
	if err != nil {
		return "", err
	}

	// Convert the buffer to base64 string
	base64Str := base64.StdEncoding.EncodeToString(buffer.Bytes())

	// Generate the Markdown image syntax with the base64-encoded image data
	markdownImage := fmt.Sprintf("![QR Code](data:image/png;base64,%s)", base64Str)

	return markdownImage, nil
}
