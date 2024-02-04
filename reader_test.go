package pgzip

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func decompress(data []byte) ([]byte, error) {
	r, err := NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()

	return io.ReadAll(r)
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func TestCompressingReader(t *testing.T) {
	type testCase struct {
		testName    string
		input       []byte
		expectedErr bool
	}

	testCases := []testCase{
		{
			testName: "SuccessfulCompression",
			input:    []byte("test data"),
		},
		{
			testName: "EmptyInput",
			input:    []byte{},
		},
		{
			testName: "NonUTF8Bytes",
			input:    []byte{0x80, 0xFF, 0x12, 0x34, 0xAA},
		},
		{
			testName:    "ReadError",
			expectedErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.testName, func(t *testing.T) {
			var r io.Reader
			if tt.expectedErr {
				r = &errorReader{}
			} else {
				r = bytes.NewReader(tt.input)
			}
			compressingReader := CompressingReader(r)

			outputData, err := io.ReadAll(compressingReader)

			if tt.expectedErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)

			decompressedData, err := decompress(outputData)
			assert.NoError(t, err)
			assert.Equal(t, tt.input, decompressedData)
		})
	}
}
