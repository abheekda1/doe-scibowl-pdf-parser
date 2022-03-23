package parse

import (
	"bytes"

	"github.com/ledongthuc/pdf"
)

// Read PDF data from a bytes reader and return raw text
func ReadPdfToString(file *bytes.Reader) (string, error) {
	r, err := pdf.NewReader(file, file.Size())
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	pdfData, err := r.GetPlainText()
	if err != nil {
		return "", err
	}
	buf.ReadFrom(pdfData)
	return buf.String(), nil
}
