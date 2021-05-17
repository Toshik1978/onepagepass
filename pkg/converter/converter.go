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
	pdfPath       string
	convertedPath string
	dpi           float64
}

// New creates new instance of Converter.
func New(pdfPath string, dpi float64) *Converter {
	if dpi == 0 {
		dpi = defaultDPI
	}

	return &Converter{
		pdfPath:       pdfPath,
		convertedPath: strings.Replace(pdfPath, ".pdf", ".converted.pdf", 1),
		dpi:           dpi,
	}
}

// Run runs converting.
func (c Converter) Run() error {
	src := pdf.New()
	if err := src.Open(c.pdfPath); err != nil {
		return fmt.Errorf("failed to open PDF file: %w", err)
	}
	defer func() { _ = src.Close() }()

	dst := pdf.New()
	if err := dst.Create(); err != nil {
		return fmt.Errorf("failed to create PDF file: %w", err)
	}

	if err := c.processPDF(src, dst); err != nil {
		return fmt.Errorf("failed to process PDF file: %w", err)
	}

	if err := dst.Save(c.convertedPath); err != nil {
		return fmt.Errorf("failed to save PDF file: %w", err)
	}
	return nil
}

func (c Converter) processPDF(src, dst *pdf.File) error {
	for i := 0; i < src.NumPages(); i++ {
		curr, err := src.Page(i, c.dpi)
		if err != nil {
			return fmt.Errorf("failed to get PDF page: %w", err)
		}

		if i != src.NumPages()-1 {
			next, err := src.Page(i+1, c.dpi)
			if err != nil {
				return fmt.Errorf("failed to get PDF page: %w", err)
			}

			next = imaging.Crop(next, image.Rect(next.Bounds().Min.X, next.Bounds().Min.Y, next.Bounds().Max.X, next.Bounds().Max.Y/2))
			curr = imaging.Paste(curr, next, image.Point{X: 0, Y: curr.Bounds().Max.Y / 2})
			i++
		}

		if err := dst.AddPage(curr); err != nil {
			return fmt.Errorf("failed to add PDF page: %w", err)
		}
	}
	return nil
}
