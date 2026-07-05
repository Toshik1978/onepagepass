package converter

import (
	"errors"
	"image"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	convertermock "github.com/Toshik1978/onepagepass/pkg/converter/mocks"
)

func TestConverter(t *testing.T) {
	suite.Run(t, new(converterTestSuite))
}

type converterTestSuite struct {
	suite.Suite

	img     image.Image
	srcPath string
	dstPath string
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

func (s *converterTestSuite) TestRunOpenFailed() {
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

	s.Require().ErrorIs(err, errFailed)
}

func (s *converterTestSuite) TestRunCreateFailed() {
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

	s.Require().ErrorIs(err, errFailed)
}

func (s *converterTestSuite) TestRunSaveFailed() {
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

	s.Require().ErrorIs(err, errFailed)
}

func (s *converterTestSuite) TestRunSucceeded() {
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

func (s *converterTestSuite) TestNewDefaults() {
	c := New("document.pdf", 0)

	s.Equal("document.pdf", c.srcPath)
	s.Equal("document.converted.pdf", c.dstPath)
	s.InEpsilon(defaultDPI, c.dpi, 0.0001)
	s.NotNil(c.src)
	s.NotNil(c.dst)
}

func (s *converterTestSuite) TestNewCustomDPI() {
	c := New("scan.pdf", 150)

	s.InEpsilon(150.0, c.dpi, 0.0001)
	s.Equal("scan.converted.pdf", c.dstPath)
}

func (s *converterTestSuite) TestRunFirstPageFailed() {
	errFailed := errors.New("failed")

	srcMock := convertermock.PdfFile{}
	srcMock.On("Open", s.srcPath).Return(nil).Once()
	srcMock.On("Close").Return(nil).Once()
	srcMock.On("NumPages").Return(3).Once()
	srcMock.On("Page", 0, defaultDPI).Return(nil, errFailed).Once()

	dstMock := convertermock.PdfFile{}
	dstMock.On("Create").Return(nil).Once()

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

	s.Require().ErrorIs(err, errFailed)
}

func (s *converterTestSuite) TestRunNextPageFailed() {
	errFailed := errors.New("failed")

	srcMock := convertermock.PdfFile{}
	srcMock.On("Open", s.srcPath).Return(nil).Once()
	srcMock.On("Close").Return(nil).Once()
	srcMock.On("NumPages").Return(3).Times(2)
	srcMock.On("Page", 0, defaultDPI).Return(s.img, nil).Once()
	srcMock.On("Page", 1, defaultDPI).Return(nil, errFailed).Once()

	dstMock := convertermock.PdfFile{}
	dstMock.On("Create").Return(nil).Once()

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

	s.Require().ErrorIs(err, errFailed)
}

func (s *converterTestSuite) TestRunAddPageFailed() {
	errFailed := errors.New("failed")

	srcMock := convertermock.PdfFile{}
	srcMock.On("Open", s.srcPath).Return(nil).Once()
	srcMock.On("Close").Return(nil).Once()
	srcMock.On("NumPages").Return(3).Times(2)
	srcMock.On("Page", mock.Anything, defaultDPI).Return(s.img, nil).Times(2)

	dstMock := convertermock.PdfFile{}
	dstMock.On("Create").Return(nil).Once()
	dstMock.On("AddPage", mock.Anything).Return(errFailed).Once()

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

	s.Require().ErrorIs(err, errFailed)
}
