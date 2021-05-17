package pdf

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"

	"github.com/gen2brain/go-fitz"
	"github.com/signintech/gopdf"
)

var (
	ErrAlreadyCreated = errors.New("PDF file already created")
	ErrNotCreated     = errors.New("PDF file not created")
	ErrAlreadyOpened  = errors.New("PDF file already opened")
	ErrNotOpened      = errors.New("PDF file not opened")
)

type File struct {
	r *fitz.Document
	w *gopdf.GoPdf
}

// New initializes new PDF file helper.
func New() *File {
	return &File{}
}

// Create initializes new empty PDF file.
func (p *File) Create() error {
	if p.r != nil {
		return ErrAlreadyOpened
	}

	p.w = &gopdf.GoPdf{}
	p.w.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})
	return nil
}

// Save saved PDF file in given destination.
func (p File) Save(path string) error {
	if p.w != nil {
		if err := p.w.WritePdf(path); err != nil {
			return fmt.Errorf("failed to save PDF file: %w", err)
		}
		return nil
	}
	return ErrNotCreated
}

// Open opens existing PDF file.
func (p *File) Open(path string) error {
	if p.w != nil {
		return ErrAlreadyCreated
	}

	var err error
	p.r, err = fitz.New(path)
	if err != nil {
		return fmt.Errorf("failed to open PDF file: %w", err)
	}
	return nil
}

// Close closes PDF file.
func (p File) Close() error {
	if p.r != nil {
		if err := p.r.Close(); err != nil {
			return fmt.Errorf("failed to close PDF file: %w", err)
		}
	}
	return nil
}

// NumPages return number of pages in PDF file.
func (p File) NumPages() int {
	if p.r != nil {
		return p.r.NumPage()
	}
	if p.w != nil {
		return p.w.GetNumberOfPages()
	}
	return 0
}

// Page return given page of PDF file.
func (p File) Page(page int, dpi float64) (image.Image, error) {
	if p.r != nil {
		img, err := p.r.ImageDPI(page, dpi)
		if err != nil {
			return nil, fmt.Errorf("failed to get image from PDF: %w", err)
		}
		return img, nil
	}
	return nil, ErrNotOpened
}

// AddPage adds page to PDF file.
func (p File) AddPage(img image.Image) error {
	if p.w != nil {
		p.w.AddPage()

		buffer := new(bytes.Buffer)
		err := jpeg.Encode(buffer, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
		if err != nil {
			return fmt.Errorf("failed to encode jpeg: %w", err)
		}

		pdfImage, err := gopdf.ImageHolderByReader(buffer)
		if err != nil {
			return fmt.Errorf("failed to create PDF image: %w", err)
		}
		if err := p.w.ImageByHolder(pdfImage, 0, 0, gopdf.PageSizeA4); err != nil {
			return fmt.Errorf("failed to add PDF image: %w", err)
		}
		return nil
	}
	return ErrNotCreated
}
