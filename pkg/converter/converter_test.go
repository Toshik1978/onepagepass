package converter

import (
	"errors"
	"image"
	"os"
	"testing"

	convertermock "github.com/Toshik1978/onepagepass/pkg/converter/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestConverter(t *testing.T) {
	suite.Run(t, new(converterTestSuite))
}

type converterTestSuite struct {
	suite.Suite

	srcPath string
	dstPath string

	img image.Image
}

func (s *converterTestSuite) SetupSuite() {
	s.srcPath = "source/path"
	s.dstPath = "destination/path"

	reader, err := os.Open("testdata/sample.jpeg")
	if err != nil {
		panic(err)
	}
	img, _, err := image.Decode(reader)
	if err != nil {
		panic(err)
	}
	s.img = img
}

func (s converterTestSuite) TestRunOpenFailed() {
	errFailed := errors.New("failed")

	srcMock := convertermock.PdfFile{}
	srcMock.
		On("Open", s.srcPath).
		Return(errFailed).
		Once()

	dstMock := convertermock.PdfFile{}

	convert := Converter{
		src:     &srcMock,
		srcPath: s.srcPath,
		dst:     &dstMock,
		dstPath: s.dstPath,
		dpi:     defaultDPI,
	}
	err := convert.Run()

	srcMock.AssertExpectations(s.T())
	dstMock.AssertExpectations(s.T())

	s.Error(err)
	s.ErrorIs(err, errFailed)
}

func (s converterTestSuite) TestRunCreateFailed() {
	errFailed := errors.New("failed")

	srcMock := convertermock.PdfFile{}
	srcMock.
		On("Open", s.srcPath).
		Return(nil).
		Once()
	srcMock.
		On("Close").
		Return(nil).
		Once()

	dstMock := convertermock.PdfFile{}
	dstMock.
		On("Create").
		Return(errFailed).
		Once()

	convert := Converter{
		src:     &srcMock,
		srcPath: s.srcPath,
		dst:     &dstMock,
		dstPath: s.dstPath,
		dpi:     defaultDPI,
	}
	err := convert.Run()

	srcMock.AssertExpectations(s.T())
	dstMock.AssertExpectations(s.T())

	s.Error(err)
	s.ErrorIs(err, errFailed)
}

func (s converterTestSuite) TestRunSaveFailed() {
	errFailed := errors.New("failed")

	srcMock := convertermock.PdfFile{}
	srcMock.
		On("Open", s.srcPath).
		Return(nil).
		Once()
	srcMock.
		On("Close").
		Return(nil).
		Once()
	srcMock.
		On("NumPages").
		Return(3).
		Times(5)
	srcMock.
		On("Page", mock.Anything, defaultDPI).
		Return(s.img, nil).
		Times(3)

	dstMock := convertermock.PdfFile{}
	dstMock.
		On("Create").
		Return(nil).
		Once()
	dstMock.
		On("Save", s.dstPath).
		Return(errFailed).
		Once()
	dstMock.
		On("AddPage", mock.Anything).
		Return(nil).
		Times(2)

	convert := Converter{
		src:     &srcMock,
		srcPath: s.srcPath,
		dst:     &dstMock,
		dstPath: s.dstPath,
		dpi:     defaultDPI,
	}
	err := convert.Run()

	srcMock.AssertExpectations(s.T())
	dstMock.AssertExpectations(s.T())

	s.Error(err)
	s.ErrorIs(err, errFailed)
}

func (s converterTestSuite) TestRunSucceeded() {
	srcMock := convertermock.PdfFile{}
	srcMock.
		On("Open", s.srcPath).
		Return(nil).
		Once()
	srcMock.
		On("Close").
		Return(nil).
		Once()
	srcMock.
		On("NumPages").
		Return(3).
		Times(5)
	srcMock.
		On("Page", mock.Anything, defaultDPI).
		Return(s.img, nil).
		Times(3)

	dstMock := convertermock.PdfFile{}
	dstMock.
		On("Create").
		Return(nil).
		Once()
	dstMock.
		On("Save", s.dstPath).
		Return(nil).
		Once()
	dstMock.
		On("AddPage", mock.Anything).
		Return(nil).
		Times(2)

	convert := Converter{
		src:     &srcMock,
		srcPath: s.srcPath,
		dst:     &dstMock,
		dstPath: s.dstPath,
		dpi:     defaultDPI,
	}
	err := convert.Run()

	srcMock.AssertExpectations(s.T())
	dstMock.AssertExpectations(s.T())

	s.NoError(err)
}
