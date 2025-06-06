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
			img.SetRGBA(x, y, color.RGBA{R: uint8(x % 255), G: uint8(y % 255), B: uint8((x + y) % 255), A: 255}) //nolint:gosec
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

func Test_cipherLSB(t *testing.T) {
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
				t.Helper()
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
				t.Helper()
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

func Test_extractLSB(t *testing.T) {
	validPNG100x50Bytes := createTestPNG(t, 100, 50)
	corruptedPNGBytes := validPNG100x50Bytes[:len(validPNG100x50Bytes)/2]

	type extractLSBTestCase struct {
		name               string
		ctx                context.Context
		stegoImageProvider func(t *testing.T) []byte
		wantErr            bool
		wantErrMsgContains string
		wantPlaintext      string
	}

	testCases := []extractLSBTestCase{
		{
			name: "SuccessCaseValidData",
			ctx:  t.Context(),
			stegoImageProvider: func(t *testing.T) []byte {
				t.Helper()
				pt := "Valid hidden text"
				inputFile := filed.File{Metadata: filed.Metadata{Type: filed.TypePNG}, File: createTestPNG(t, 150, 100)}
				outputFile, err := cipherLSB(t.Context(), pt, inputFile)
				require.NoError(t, err, "Setup failed: cipherLSB errored")
				return outputFile.File
			},
			wantErr:       false,
			wantPlaintext: "Valid hidden text",
		},
		{
			name: "SuccessCaseEmptyPlaintext",
			ctx:  t.Context(),
			stegoImageProvider: func(t *testing.T) []byte {
				t.Helper()
				pt := ""
				inputFile := filed.File{Metadata: filed.Metadata{Type: filed.TypePNG}, File: createTestPNG(t, 30, 30)}
				outputFile, err := cipherLSB(t.Context(), pt, inputFile)
				require.NoError(t, err, "Setup failed: cipherLSB errored for empty plaintext")
				return outputFile.File
			},
			wantErr:       false,
			wantPlaintext: "E",
		},
		{
			name: "ErrorCaseCorruptedImageData",
			ctx:  t.Context(),
			stegoImageProvider: func(t *testing.T) []byte {
				t.Helper()
				return corruptedPNGBytes
			},
			wantErr:            true,
			wantErrMsgContains: "cannot decode image",
		},
		{
			name: "ErrorCaseDataTruncated",
			ctx:  t.Context(),
			stegoImageProvider: func(t *testing.T) []byte {
				t.Helper()
				pt := "Short msg"
				inputFile := filed.File{Metadata: filed.Metadata{Type: filed.TypePNG}, File: createTestPNG(t, 10, 6)}
				outputFile, err := cipherLSB(t.Context(), pt, inputFile)
				require.NoError(t, err, "Setup failed: cipherLSB errored")
				return outputFile.File[:60]
			},
			wantErr:            true,
			wantErrMsgContains: "fail to extract:",
		},
		{
			name: "ErrorCaseImageTooSmallForDeclaredLength",
			ctx:  t.Context(),
			stegoImageProvider: func(t *testing.T) []byte {
				t.Helper()
				img := image.NewRGBA(image.Rect(0, 0, 8, 8))
				lenBytes := make([]byte, 8)
				binary.BigEndian.PutUint64(lenBytes, uint64(20))
				dataIndex := 0
				bitIndex := 0
				pixelCount := 0
				for y := 0; y < 8 && pixelCount < 22; y++ {
					for x := 0; x < 8 && pixelCount < 22; x++ {
						c := color.RGBA{R: byte(x * 10), G: byte(y * 10), B: 100, A: 255}
						channels := []*uint8{&c.R, &c.G, &c.B}
						for i := 0; i < 3 && pixelCount*3+i < 64; i++ {
							secretBit := (lenBytes[dataIndex] >> (7 - bitIndex)) & 1
							*channels[i] = (*channels[i] & 0xFE) | secretBit
							bitIndex++
							if bitIndex > 7 {
								bitIndex = 0
								dataIndex++
							}
						}
						img.SetRGBA(x, y, c)
						pixelCount++
					}
				}
				var buf bytes.Buffer
				err := png.Encode(&buf, img)
				require.NoError(t, err)
				return buf.Bytes()
			},
			wantErr:            true,
			wantErrMsgContains: "image capacity",
		},
		{
			name: "ErrorCaseContextCancelled",
			ctx: func() context.Context {
				c, cancel := context.WithCancel(t.Context())
				cancel()
				return c
			}(),
			stegoImageProvider: func(t *testing.T) []byte {
				t.Helper()
				pt := "Some data for cancellation test..."
				inputFile := filed.File{Metadata: filed.Metadata{Type: filed.TypePNG}, File: createTestPNG(t, 250, 250)}
				outputFile, err := cipherLSB(t.Context(), pt, inputFile)
				require.NoError(t, err, "Setup failed: cipherLSB errored")
				return outputFile.File
			},
			wantErr:            true,
			wantErrMsgContains: context.Canceled.Error(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			stegoBytes := tc.stegoImageProvider(t)
			actualPlaintext, err := extractLSB(tc.ctx, stegoBytes)

			if tc.wantErr {
				require.Error(t, err, "Expected an error but got none")
				if tc.wantErrMsgContains != "" {
					assert.Contains(t, err.Error(), tc.wantErrMsgContains, "Error message mismatch")
				}
				if tc.name == "ErrorCaseContextCancelled" {
					assert.True(t, errors.Is(err, context.Canceled))
				}
				if errors.Is(tc.ctx.Err(), context.DeadlineExceeded) {
					assert.True(t, errors.Is(err, context.DeadlineExceeded))
				}
				assert.Equal(t, "", actualPlaintext, "Plaintext should be empty on error")
			} else {
				require.NoError(t, err, "Expected no error but got one: %v", err)
				assert.Equal(t, tc.wantPlaintext, actualPlaintext, "Extracted plaintext mismatch")
			}
		})
	}
}
