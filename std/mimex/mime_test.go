package mimex

import (
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/fasionchan/goutils/baseutils"
	"gotest.tools/assert"
)

type PeekContentTypeTestCase struct {
	contentType        string
	contentDisposition string
	fileName           string
	expected           *MimeType
}

func (testCase *PeekContentTypeTestCase) Run(t *testing.T) {
	t.Helper()
	_, testFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	testDir := filepath.Dir(testFile)
	filePath := filepath.Join(testDir, "testdata", testCase.fileName)

	file, err := os.Open(filePath)
	if err != nil {
		t.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	mimeType, contentReader, err := PeekContentType(testCase.contentType, testCase.contentDisposition, file, 0)
	if err != nil {
		t.Fatalf("failed to peek content type: %v", err)
	}

	AssertMimeType(t, mimeType, testCase.expected)

	md51 := baseutils.GetMd5Hasher().SumRtos(contentReader)
	if err != nil {
		t.Fatalf("failed to get md5: %v", err)
	}

	file.Seek(0, io.SeekStart)
	md52 := baseutils.GetMd5Hasher().SumRtos(file)
	if err != nil {
		t.Fatalf("failed to get md5: %v", err)
	}

	if md51 != md52 {
		t.Fatalf("md5 mismatch: expected `%s` but got `%s`", md51, md52)
	}
}

func TestPeekContentType(t *testing.T) {
	for _, testCase := range []PeekContentTypeTestCase{
		// image
		{
			fileName: "file.png",
			expected: &MimeType{
				Type:    "image/png",
				TopType: "image",
				SubType: "png",
			},
		},

		// document
		{
			fileName: "file.pdf",
			expected: &MimeType{
				Type:    "application/pdf",
				TopType: "application",
				SubType: "pdf",
			},
		},
		{
			fileName: "file.docx",
			expected: &MimeType{
				Type:    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				TopType: "application",
				SubType: "vnd.openxmlformats-officedocument.wordprocessingml.document",
			},
		},
		{
			fileName: "file.xlsx",
			expected: &MimeType{
				Type:    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				TopType: "application",
				SubType: "vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			},
		},
		{
			fileName: "file.pptx",
			expected: &MimeType{
				Type:    "application/vnd.openxmlformats-officedocument.presentationml.presentation",
				TopType: "application",
				SubType: "vnd.openxmlformats-officedocument.presentationml.presentation",
			},
		},
		{
			fileName: "file.doc",
			expected: &MimeType{
				Type:    "application/msword",
				TopType: "application",
				SubType: "msword",
			},
		},
		{
			fileName: "file.xls",
			expected: &MimeType{
				Type:    "application/vnd.ms-excel",
				TopType: "application",
				SubType: "vnd.ms-excel",
			},
		},
		// {
		// 	fileName: "file.ppt",
		// 	expected: &MimeType{
		// 		Type:    "application/vnd.ms-powerpoint",
		// 		TopType: "application",
		// 		SubType: "vnd.ms-powerpoint",
		// 	},
		// },
	} {
		testCase.Run(t)
	}
}

func TestParseDispositionMimeType(t *testing.T) {
	for _, testCase := range []struct {
		contentDisposition string
		expected           *MimeType
	}{
		// image
		{
			contentDisposition: "attachment; filename=test.jpg",
			expected: &MimeType{
				Type:    "image/jpeg",
				TopType: "image",
				SubType: "jpeg",
			},
		},
		{
			contentDisposition: "attachment; filename=test.png",
			expected: &MimeType{
				Type:    "image/png",
				TopType: "image",
				SubType: "png",
			},
		},

		// document
		{
			contentDisposition: "attachment; filename=test.txt",
			expected: &MimeType{
				Type:    "text/plain",
				TopType: "text",
				SubType: "plain",
			},
		},
		{
			contentDisposition: "attachment; filename=test.pdf",
			expected: &MimeType{
				Type:    "application/pdf",
				TopType: "application",
				SubType: "pdf",
			},
		},
		{
			contentDisposition: "attachment; filename=test.docx",
			expected: &MimeType{
				Type:    "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
				TopType: "application",
				SubType: "vnd.openxmlformats-officedocument.wordprocessingml.document",
			},
		},
		{
			contentDisposition: "attachment; filename=test.xlsx",
			expected: &MimeType{
				Type:    "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
				TopType: "application",
				SubType: "vnd.openxmlformats-officedocument.spreadsheetml.sheet",
			},
		},
		{
			contentDisposition: "attachment; filename=test.pptx",
			expected: &MimeType{
				Type:    "application/vnd.openxmlformats-officedocument.presentationml.presentation",
				TopType: "application",
				SubType: "vnd.openxmlformats-officedocument.presentationml.presentation",
			},
		},
		{
			contentDisposition: "attachment; filename=test.doc",
			expected: &MimeType{
				Type:    "application/msword",
				TopType: "application",
				SubType: "msword",
			},
		},
		{
			contentDisposition: "attachment; filename=test.xls",
			expected: &MimeType{
				Type:    "application/vnd.ms-excel",
				TopType: "application",
				SubType: "vnd.ms-excel",
			},
		},
		{
			contentDisposition: "attachment; filename=test.ppt",
			expected: &MimeType{
				Type:    "application/vnd.ms-powerpoint",
				TopType: "application",
				SubType: "vnd.ms-powerpoint",
			},
		},

		// archive
		{
			contentDisposition: "attachment; filename=test.zip",
			expected: &MimeType{
				Type:    "application/zip",
				TopType: "application",
				SubType: "zip",
			},
		},
		// {
		// 	contentDisposition: "attachment; filename=test.rar",
		// 	expected: &MimeType{
		// 		Type: "application/vnd.rar",
		// 	},
		// },
	} {
		mt, err := ParseDispositionMimeType(testCase.contentDisposition)
		if err != nil {
			t.Fatalf("failed to parse disposition mime type: %v", err)
		}
		mt.Params = nil

		AssertMimeType(t, mt, testCase.expected)
	}
}

func TestParseMimeType(t *testing.T) {
	for _, testCase := range []struct {
		mimeType string
		expected *MimeType
	}{
		{
			mimeType: "text/plain",
			expected: &MimeType{
				Type:    "text/plain",
				TopType: "text",
				SubType: "plain",
				Params:  map[string]string{},
			},
		},
		{
			mimeType: "text/plain; charset=utf-8",
			expected: &MimeType{
				Type:    "text/plain",
				TopType: "text",
				SubType: "plain",
				Params:  map[string]string{"charset": "utf-8"},
			},
		},
		{
			mimeType: "text/plain; charset=utf-8; foo=bar",
			expected: &MimeType{
				Type:    "text/plain",
				TopType: "text",
				SubType: "plain",
				Params:  map[string]string{"charset": "utf-8", "foo": "bar"},
			},
		},
		{
			mimeType: "text/plain; charset=utf-8; foo=bar; baz=qux",
			expected: &MimeType{
				Type:    "text/plain",
				TopType: "text",
				SubType: "plain",
				Params:  map[string]string{"charset": "utf-8", "foo": "bar", "baz": "qux"},
			},
		},
	} {
		mt, err := ParseMimeType(testCase.mimeType)
		if err != nil {
			t.Fatalf("failed to parse mime type: %v", err)
		}

		AssertMimeType(t, mt, testCase.expected)
	}
}

func AssertMimeType(t *testing.T, actual *MimeType, expected *MimeType) {
	if actual.GetType() != expected.GetType() {
		t.Fatalf("expected mime type `%s` but got `%s`", expected.GetType(), actual.GetType())
	}

	if actual.GetTopType() != expected.GetTopType() {
		t.Fatalf("expected top type `%s` but got `%s`", expected.GetTopType(), actual.GetTopType())
	}

	if actual.GetSubType() != expected.GetSubType() {
		t.Fatalf("expected sub type `%s` but got `%s`", expected.GetSubType(), actual.GetSubType())
	}

	if len(actual.GetParams()) != len(expected.GetParams()) {
		t.Fatalf("expected params `%v` but got `%v`", expected.GetParams(), actual.GetParams())
	}

	for key, value := range expected.GetParams() {
		if actual.GetParams()[key] != value {
			t.Fatalf("expected param `%s` but got `%s`", value, actual.GetParams()[key])
		}
	}

	if actual.String() != expected.String() {
		t.Fatalf("expected string `%s` but got `%s`", expected.String(), actual.String())
	}
}

func TestMimeType(t *testing.T) {
	for _, ext := range []string{
		".zip",
		".rar",
		".7z",
		".tar",
		".gz",
		".bz2",
		".xz",
	} {
		fmt.Println(mime.TypeByExtension(ext))
	}
	assert.Assert(t, false)
}