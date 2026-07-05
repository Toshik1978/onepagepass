package main

import (
	"context"
	"image"
	"image/color"
	"io"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/Toshik1978/onepagepass/pkg/pdf"
)

func TestCommand(t *testing.T) {
	suite.Run(t, new(commandTestSuite))
}

type commandTestSuite struct {
	suite.Suite

	img image.Image
}

func (s *commandTestSuite) SetupSuite() {
	img := image.NewRGBA(image.Rect(0, 0, 200, 300))
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, color.White)
		}
	}
	for x := 20; x < 180; x++ {
		img.Set(x, 40, color.Black)
	}
	s.img = img
}

// run executes the CLI with the given args, discarding its output.
func (s *commandTestSuite) run(args ...string) error {
	cmd := newRootCommand()
	cmd.Writer = io.Discard
	cmd.ErrWriter = io.Discard

	return cmd.Run(context.Background(), append([]string{"onepagepass"}, args...))
}

// writeSourcePDF builds a real multi-page PDF on disk and returns its path.
func (s *commandTestSuite) writeSourcePDF(pages int) string {
	path := filepath.Join(s.T().TempDir(), "src.pdf")

	w := pdf.New()
	s.Require().NoError(w.Create())
	for range pages {
		s.Require().NoError(w.AddPage(s.img))
	}
	s.Require().NoError(w.Save(path))

	return path
}

func (s *commandTestSuite) TestConvertSucceeds() {
	src := s.writeSourcePDF(3)

	s.Require().NoError(s.run("convert", "--pdf", src, "--dpi", "150"))

	converted := strings.Replace(src, ".pdf", ".converted.pdf", 1)
	s.FileExists(converted)

	r := pdf.New()
	s.Require().NoError(r.Open(converted))
	s.Equal(2, r.NumPages())
	s.Require().NoError(r.Close())
}

func (s *commandTestSuite) TestConvertAliasSucceeds() {
	src := s.writeSourcePDF(2)

	s.Require().NoError(s.run("c", "--pdf", src))

	s.FileExists(strings.Replace(src, ".pdf", ".converted.pdf", 1))
}

func (s *commandTestSuite) TestConvertMissingRequiredFlag() {
	s.Require().Error(s.run("convert"))
}

func (s *commandTestSuite) TestConvertOpenError() {
	missing := filepath.Join(s.T().TempDir(), "missing.pdf")
	s.Require().Error(s.run("convert", "--pdf", missing))
}

func (s *commandTestSuite) TestRootShowsHelp() {
	s.Require().NoError(s.run())
}
