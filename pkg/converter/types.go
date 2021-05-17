package converter

import (
	"image"
)

//go:generate mockery --name PdfFile --filename converter.go --outpkg convertermock
// PdfFile declare interface for PDF file.
type PdfFile interface {
	Create() error
	Save(path string) error
	Open(path string) error
	Close() error

	NumPages() int
	Page(page int, dpi float64) (image.Image, error)
	AddPage(img image.Image) error
}
