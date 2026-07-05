// Package converter re-lays a rasterized PDF into a two-pages-per-A4 format.
package converter

import (
	"fmt"
	"image"
	"image/draw"
	"strings"

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

			curr = mergePages(curr, next)
			i++
		}

		if err := c.dst.AddPage(curr); err != nil {
			return fmt.Errorf("failed to add PDF page: %w", err)
		}
	}

	return nil
}

// mergePages packs two document scans onto a single sheet: curr keeps its upper
// half, and its lower half is replaced by the upper half of next.
func mergePages(curr, next image.Image) image.Image {
	b := curr.Bounds()
	merged := image.NewRGBA(b)
	draw.Draw(merged, b, curr, b.Min, draw.Src)

	// Copy the top of next into the lower half of curr.
	midY := b.Min.Y + b.Dy()/2
	dstRect := image.Rect(b.Min.X, midY, b.Max.X, b.Max.Y)
	draw.Draw(merged, dstRect, next, next.Bounds().Min, draw.Src)

	return merged
}
