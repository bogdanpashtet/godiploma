package cipher

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"testing"
	"time"

	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestPNG(t *testing.T, width, height int) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, color.RGBA{R: uint8(x % 255), G: uint8(y % 255), B: uint8((x + y) % 255), A: 255})
		}
	}
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	require.NoError(t, err, "Failed to encode test PNG")
	return buf.Bytes()
}

func extractLSBData(t *testing.T, fileBytes []byte, totalBytesToExtract uint64) ([]byte, error) {
	t.Helper()
	img, _, err := image.Decode(bytes.NewReader(fileBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to decode image for extraction: %w", err)
	}
	bounds := img.Bounds()
	extractedData := make([]byte, totalBytesToExtract)
	dataIndex := uint64(0)
	bitIndex := 0
	currentByte := byte(0)

pixelLoopExtract:
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if dataIndex >= totalBytesToExtract {
				break pixelLoopExtract
			}
			r, g, b, _ := img.At(x, y).RGBA()
			lsbChannels := []byte{byte(r), byte(g), byte(b)}
			for i := 0; i < 3; i++ {
				if dataIndex >= totalBytesToExtract {
					break pixelLoopExtract
				}
				extractedBit := lsbChannels[i] & 1
				currentByte = (currentByte << 1) | extractedBit
				bitIndex++
				if bitIndex > 7 {
					extractedData[dataIndex] = currentByte
					bitIndex = 0
					dataIndex++
					currentByte = 0
				}
			}
		}
	}
	if dataIndex < totalBytesToExtract {
		return nil, fmt.Errorf("extracted only %d bytes, expected %d", dataIndex, totalBytesToExtract)
	}
	return extractedData, nil
}

func TestCipherLSBTableDriven(t *testing.T) {
	validPNG100x50Bytes := createTestPNG(t, 100, 50)
	validPNGSmallBytes := createTestPNG(t, 5, 5)
	corruptedPNGBytes := validPNG100x50Bytes[:len(validPNG100x50Bytes)/2]

	testCases := []struct {
		name               string
		inputFile          filed.File
		plaintext          string
		ctx                context.Context
		wantErr            bool
		wantErrMsgContains string
		checkOutput        func(t *testing.T, output filed.File, input filed.File)
	}{
		{
			name: "Success Case PNG",
			inputFile: filed.File{
				Metadata: filed.Metadata{Type: filed.TypePNG},
				File:     validPNG100x50Bytes,
			},
			plaintext: "Hello LSB!",
			ctx:       t.Context(),
			wantErr:   false,
			checkOutput: func(t *testing.T, output filed.File, input filed.File) {
				require.NotEmpty(t, output.File)
				assert.Equal(t, filed.TypePNG, output.Metadata.Type)
				assert.NotEqual(t, input.File, output.File)

				expectedLen := uint64(len("Hello LSB!"))
				extractedFullData, errExtract := extractLSBData(t, output.File, 8+expectedLen)
				require.NoError(t, errExtract)
				extractedLen := binary.BigEndian.Uint64(extractedFullData[:8])
				require.Equal(t, expectedLen, extractedLen)
				assert.Equal(t, "Hello LSB!", string(extractedFullData[8:]))
			},
		},
		{
			name: "Error Case Unsupported Format JPEG",
			inputFile: filed.File{
				Metadata: filed.Metadata{Type: filed.TypeUnspecified},
				File:     validPNGSmallBytes,
			},
			plaintext:          "test",
			ctx:                t.Context(),
			wantErr:            true,
			wantErrMsgContains: "LSB is not suitable",
		},
		{
			name: "Error Case Corrupted PNG Data",
			inputFile: filed.File{
				Metadata: filed.Metadata{Type: filed.TypePNG},
				File:     corruptedPNGBytes,
			},
			plaintext:          "test",
			ctx:                t.Context(),
			wantErr:            true,
			wantErrMsgContains: "fail to decode file",
		},
		{
			name: "Error Case Insufficient Capacity",
			inputFile: filed.File{
				Metadata: filed.Metadata{Type: filed.TypePNG},
				File:     validPNGSmallBytes,
			},
			plaintext:          "This message is way too long to fit",
			ctx:                t.Context(),
			wantErr:            true,
			wantErrMsgContains: "image too small",
		},
		{
			name: "Edge Case Empty Plaintext",
			inputFile: filed.File{
				Metadata: filed.Metadata{Type: filed.TypePNG},
				File:     createTestPNG(t, 10, 10),
			},
			plaintext: "",
			ctx:       t.Context(),
			wantErr:   false,
			checkOutput: func(t *testing.T, output filed.File, input filed.File) {
				require.NotEmpty(t, output.File)
				assert.NotEqual(t, input.File, output.File)

				prefixBytes, errExtract := extractLSBData(t, output.File, 8)
				require.NoError(t, errExtract)
				extractedLen := binary.BigEndian.Uint64(prefixBytes)
				assert.Equal(t, uint64(1), extractedLen)
			},
		},
		{
			name: "Error Case Context Cancelled",
			inputFile: filed.File{
				Metadata: filed.Metadata{Type: filed.TypePNG},
				File:     createTestPNG(t, 150, 150),
			},
			plaintext: "Some data",
			ctx: func() context.Context {
				c, cancel := context.WithCancel(t.Context())
				cancel()
				return c
			}(),
			wantErr:            true,
			wantErrMsgContains: context.Canceled.Error(),
		},
		{
			name: "Error Case Context Timeout",
			inputFile: filed.File{
				Metadata: filed.Metadata{Type: filed.TypePNG},
				File:     createTestPNG(t, 150, 150),
			},
			plaintext: "Some data",
			ctx: func() context.Context {
				c, cancel := context.WithTimeout(t.Context(), 1*time.Nanosecond)
				_ = cancel
				return c
			}(),
			wantErr:            true,
			wantErrMsgContains: context.Canceled.Error(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			outputFile, err := cipherLSB(tc.ctx, tc.plaintext, tc.inputFile)

			if tc.wantErr {
				require.Error(t, err, "Expected an error but got none")
				if tc.wantErrMsgContains != "" {
					assert.Contains(t, err.Error(), tc.wantErrMsgContains, "Error message mismatch")
				}
				if tc.name == "ErrorCaseContextCancelled" {
					assert.True(t, errors.Is(err, context.Canceled))
				}
				if tc.name == "ErrorCaseContextTimeout" {
					assert.True(t, errors.Is(err, context.DeadlineExceeded))
				}
				assert.Empty(t, outputFile.File)
				assert.Empty(t, outputFile.Metadata.Type)
			} else {
				require.NoError(t, err, "Expected no error but got one: %v", err)
				if tc.checkOutput != nil {
					tc.checkOutput(t, outputFile, tc.inputFile)
				} else {
					require.NotEmpty(t, outputFile.File, "Output file should not be empty on success")
				}
			}
		})
	}
}
