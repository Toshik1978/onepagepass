package converter

import (
	"fmt"
	"image"
	"strings"

	"github.com/disintegration/imaging"

	"github.com/Toshik1978/onepagepass/pkg/pdf"
)

const (
	defaultDPI = 300.0
)

// Converter declare implementation of the main module.
type Converter struct {
	src     PdfFile
	srcPath string

	dst     PdfFile
	dstPath string

	dpi float64
}

// New creates new instance of Converter.
func New(pdfPath string, dpi float64) *Converter {
	if dpi == 0 {
		dpi = defaultDPI
	}
	dstPath := strings.Replace(pdfPath, ".pdf", ".converted.pdf", 1)

	return &Converter{
		src:     pdf.New(),
		srcPath: pdfPath,
		dst:     pdf.New(),
		dstPath: dstPath,
		dpi:     dpi,
	}
}

// Run runs converting.
func (c Converter) Run() error {
	if err := c.src.Open(c.srcPath); err != nil {
		return fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer func() { _ = c.src.Close() }()

	if err := c.dst.Create(); err != nil {
		return fmt.Errorf("failed to create PDF file: %w", err)
	}

	if err := c.processPdf(); err != nil {
		return fmt.Errorf("failed to process PDF file: %w", err)
	}

	if err := c.dst.Save(c.dstPath); err != nil {
		return fmt.Errorf("failed to save PDF file: %w", err)
	}
	return nil
}

func (c Converter) processPdf() error {
	for i := 0; i < c.src.NumPages(); i++ {
		curr, err := c.src.Page(i, c.dpi)
		if err != nil {
			return fmt.Errorf("failed to get PDF page: %w", err)
		}

		if i != c.src.NumPages()-1 {
			next, err := c.src.Page(i+1, c.dpi)
			if err != nil {
				return fmt.Errorf("failed to get PDF page: %w", err)
			}

			next = imaging.Crop(next, image.Rect(next.Bounds().Min.X, next.Bounds().Min.Y, next.Bounds().Max.X, next.Bounds().Max.Y/2))
			curr = imaging.Paste(curr, next, image.Point{X: 0, Y: curr.Bounds().Max.Y / 2})
			i++
		}

		if err := c.dst.AddPage(curr); err != nil {
			return fmt.Errorf("failed to add PDF page: %w", err)
		}
	}
	return nil
}
