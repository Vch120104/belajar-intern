package utils

import (
	"github.com/go-pdf/fpdf"
)

// CreatePdf creates a pdf file
func GeneratePDF(filename, text string) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, text)
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return err
	}
	return nil
}
