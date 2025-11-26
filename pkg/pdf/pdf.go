package pdf

import (
	"fmt"
	"github.com/phpdave11/gofpdf"
	"links-cheker-2511/internal/storage"
	"time"
)

// TODO: Подумать над зависимостью от storage.Links
func GeneratePDF(links []storage.Links) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 18)
	pdf.Cell(40, 10, "Link Status Report")
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 14)

	headers := []string{"URL", "Status"}
	colWidths := []float64{80, 50}

	for i, header := range headers {
		pdf.CellFormat(colWidths[i], 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	for _, link := range links {
		pdf.CellFormat(colWidths[0], 6, link.URL, "1", 0, "", false, 0, "")
		pdf.CellFormat(colWidths[1], 6, link.Status, "1", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}

	filename := fmt.Sprintf("report_%d.pdf", time.Now().Unix())
	err := pdf.OutputFileAndClose(filename)
	if err != nil {
		return "", err
	}
	//defer os.Remove(filename)

	return filename, nil
}
