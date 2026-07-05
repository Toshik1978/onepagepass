package pdf

import (
	"image"
	"image/color"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestPDF(t *testing.T) {
	suite.Run(t, new(pdfTestSuite))
}

type pdfTestSuite struct {
	suite.Suite

	img image.Image
}

func (s *pdfTestSuite) SetupSuite() {
	img := image.NewRGBA(image.Rect(0, 0, 200, 300))
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, color.White)
		}
	}
	for x := 20; x < 180; x++ {
		img.Set(x, 40, color.Black)
		img.Set(x, 260, color.Black)
	}
	s.img = img
}

// writeSamplePDF builds a real one-page PDF on disk and returns its path.
func (s *pdfTestSuite) writeSamplePDF() string {
	path := filepath.Join(s.T().TempDir(), "sample.pdf")

	w := New()
	s.Require().NoError(w.Create())
	s.Require().NoError(w.AddPage(s.img))
	s.Require().NoError(w.Save(path))

	return path
}

func (s *pdfTestSuite) TestWriteLifecycle() {
	path := filepath.Join(s.T().TempDir(), "out.pdf")

	w := New()
	s.Equal(0, w.NumPages())
	s.Require().NoError(w.Create())
	s.Require().NoError(w.AddPage(s.img))
	s.Require().NoError(w.AddPage(s.img))
	s.Equal(2, w.NumPages())
	s.Require().NoError(w.Save(path))
	s.FileExists(path)
}

func (s *pdfTestSuite) TestReadLifecycle() {
	r := New()
	s.Require().NoError(r.Open(s.writeSamplePDF()))
	s.Equal(1, r.NumPages())

	img, err := r.Page(0, 150)
	s.Require().NoError(err)
	s.Require().NotNil(img)
	s.Positive(img.Bounds().Dx())
	s.Positive(img.Bounds().Dy())

	s.Require().NoError(r.Close())
}

func (s *pdfTestSuite) TestCreateAfterOpen() {
	r := New()
	s.Require().NoError(r.Open(s.writeSamplePDF()))
	s.Require().ErrorIs(r.Create(), ErrAlreadyOpened)
	s.Require().NoError(r.Close())
}

func (s *pdfTestSuite) TestOpenAfterCreate() {
	w := New()
	s.Require().NoError(w.Create())
	s.Require().ErrorIs(w.Open("whatever.pdf"), ErrAlreadyCreated)
}

func (s *pdfTestSuite) TestOpenMissingFile() {
	path := filepath.Join(s.T().TempDir(), "missing.pdf")
	s.Require().Error(New().Open(path))
}

func (s *pdfTestSuite) TestSaveNotCreated() {
	s.Require().ErrorIs(New().Save("x.pdf"), ErrNotCreated)
}

func (s *pdfTestSuite) TestSaveWriteError() {
	w := New()
	s.Require().NoError(w.Create())
	s.Require().NoError(w.AddPage(s.img))

	// Parent directory does not exist, so the underlying write must fail.
	badPath := filepath.Join(s.T().TempDir(), "missing-dir", "out.pdf")
	s.Require().Error(w.Save(badPath))
}

func (s *pdfTestSuite) TestPageInvalidIndex() {
	r := New()
	s.Require().NoError(r.Open(s.writeSamplePDF()))
	defer func() { s.Require().NoError(r.Close()) }()

	_, err := r.Page(99, 150)
	s.Require().Error(err)
}

func (s *pdfTestSuite) TestPageNotOpened() {
	_, err := New().Page(0, 150)
	s.Require().ErrorIs(err, ErrNotOpened)
}

func (s *pdfTestSuite) TestAddPageNotCreated() {
	s.Require().ErrorIs(New().AddPage(s.img), ErrNotCreated)
}

func (s *pdfTestSuite) TestCloseWithoutOpen() {
	s.Require().NoError(New().Close())
}
