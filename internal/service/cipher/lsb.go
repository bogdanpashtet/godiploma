package cipher

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"

	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
	"github.com/samber/lo"
	"golang.org/x/image/bmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var losslessFormats = []filed.Type{filed.TypeBMP, filed.TypePNG}

//nolint:mnd,gocyclo,lll
func cipherLSB(ctx context.Context, plaintext string, inputFile filed.File) (filed.File, error) {
	if !lo.Contains(losslessFormats, inputFile.Metadata.Type) {
		return filed.File{}, status.Error(codes.InvalidArgument, fmt.Sprintf("LSB is not suitable for lossy format '%s'. Returning original file.", inputFile.Metadata.Type))
	}

	img, _, err := image.Decode(bytes.NewReader(inputFile.File))
	if err != nil {
		return filed.File{}, fmt.Errorf("fail to decode file: %w", err)
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	plaintextBytes := []byte(plaintext)
	dataLen := uint64(len(plaintextBytes))

	dataToEmbed := make([]byte, 8+len(plaintextBytes))
	binary.BigEndian.PutUint64(dataToEmbed[:8], dataLen)
	copy(dataToEmbed[8:], plaintextBytes)
	totalDataBytesToEmbed := len(dataToEmbed)

	maxEmbeddableBytes := (width * height * 3) / 8
	if totalDataBytesToEmbed > maxEmbeddableBytes {
		return filed.File{}, fmt.Errorf("image too small (%d bytes capacity) to hold data (%d bytes including length)", maxEmbeddableBytes, totalDataBytesToEmbed)
	}

	imgRGBA := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			imgRGBA.Set(x, y, img.At(x, y))
		}
	}

	dataIndex := 0
	bitIndex := 0

pixelLoop:
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if dataIndex >= totalDataBytesToEmbed {
				break pixelLoop
			}

			c := imgRGBA.RGBAAt(x, y)
			channels := []*uint8{&c.R, &c.G, &c.B}

			for i := 0; i < 3; i++ {
				if dataIndex >= totalDataBytesToEmbed {
					break pixelLoop
				}

				secretBit := (dataToEmbed[dataIndex] >> (7 - bitIndex)) & 1

				*channels[i] = (*channels[i] & 0xFE) | secretBit

				bitIndex++
				if bitIndex > 7 {
					bitIndex = 0
					dataIndex++
				}
			}

			imgRGBA.SetRGBA(x, y, c)

			if (x*y)%100000 == 0 {
				select {
				case <-ctx.Done():
					return filed.File{}, context.Canceled
				default:
				}
			}
		}
	}

	var buf bytes.Buffer
	var encodeErr error
	outputFileType := inputFile.Metadata.Type

	switch outputFileType {
	case filed.TypePNG:
		encodeErr = png.Encode(&buf, imgRGBA)
	case filed.TypeBMP:
		encodeErr = bmp.Encode(&buf, imgRGBA)
	default:
		encodeErr = fmt.Errorf("unknown output format '%s'", outputFileType)
	}
	if encodeErr != nil {
		return filed.File{}, fmt.Errorf("fail to encode image: %w", encodeErr)
	}

	outputFile := filed.File{
		Metadata: filed.Metadata{
			Type: outputFileType,
		},
		File: buf.Bytes(),
	}

	return outputFile, nil
}

//nolint:mnd,gocyclo,lll,funlen
func extractLSB(ctx context.Context, stegoImageBytes []byte) (string, error) {
	img, format, err := image.Decode(bytes.NewReader(stegoImageBytes))
	if err != nil {
		return "", fmt.Errorf("fail to extract: cannot decode image (format: %s): %w", format, err)
	}
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()

	totalBitsCapacity := int64(width) * int64(height) * 3
	if totalBitsCapacity < 64 {
		return "", fmt.Errorf("fail to extract: image too small for length prefix (capacity %d bits)", totalBitsCapacity)
	}

	lengthBytes := make([]byte, 8)
	neededBits := uint64(64)
	bitsExtracted := uint64(0)
	byteIndex := 0
	currentByte := byte(0)
	bitsInCurrentByte := 0

lengthLoopExtract:
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if x == bounds.Min.X && y%20 == 0 {
				select {
				case <-ctx.Done():
					return "", ctx.Err()
				default:
				}
			}

			r, g, b, _ := img.At(x, y).RGBA()
			lsbChannels := []byte{byte(r), byte(g), byte(b)}

			for i := 0; i < 3; i++ {
				if bitsExtracted >= neededBits {
					break lengthLoopExtract
				}
				extractedBit := lsbChannels[i] & 1
				currentByte = (currentByte << 1) | extractedBit
				bitsInCurrentByte++

				if bitsInCurrentByte == 8 {
					if byteIndex < len(lengthBytes) {
						lengthBytes[byteIndex] = currentByte
					} else {
						return "", fmt.Errorf("fail to extract length: buffer overflow (internal error)")
					}
					byteIndex++
					bitsInCurrentByte = 0
					currentByte = 0
				}
				bitsExtracted++
			}
		}
	}

	if bitsExtracted < neededBits {
		return "", fmt.Errorf("fail to extract: could not extract full 8-byte length prefix (got %d bits)", bitsExtracted)
	}

	expectedDataLength := binary.BigEndian.Uint64(lengthBytes)

	totalBytesToExtract := uint64(8) + expectedDataLength
	totalBitsToExtract := totalBytesToExtract * 8
	if totalBitsToExtract > uint64(totalBitsCapacity) {
		return "", fmt.Errorf("fail to extract: image capacity (%d bits) < required bits (%d) for declared data length %d", totalBitsCapacity, totalBitsToExtract, expectedDataLength)
	}

	if expectedDataLength == 0 {
		return "", nil
	}

	extractedData := make([]byte, expectedDataLength)
	byteIndex = 0
	currentByte = 0
	bitsInCurrentByte = 0
	bitsExtracted = 0
	neededBits = expectedDataLength * 8
	overallBitsRead := uint64(0)

dataLoopExtract:
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if x == bounds.Min.X && y%20 == 0 {
				select {
				case <-ctx.Done():
					return "", ctx.Err()
				default:
				}
			}

			r, g, b, _ := img.At(x, y).RGBA()
			lsbChannels := []byte{byte(r), byte(g), byte(b)}

			for i := 0; i < 3; i++ {
				if overallBitsRead < 64 {
					overallBitsRead++
					continue
				}

				if bitsExtracted >= neededBits {
					break dataLoopExtract
				}

				extractedBit := lsbChannels[i] & 1
				currentByte = (currentByte << 1) | extractedBit
				bitsInCurrentByte++

				if bitsInCurrentByte == 8 {
					if byteIndex < len(extractedData) {
						extractedData[byteIndex] = currentByte
					} else {
						return "", fmt.Errorf("fail to extract data: buffer overflow (internal error)")
					}
					byteIndex++
					bitsInCurrentByte = 0
					currentByte = 0
				}
				bitsExtracted++
				overallBitsRead++
			}
		}
	}

	if bitsExtracted < neededBits {
		return "", fmt.Errorf("fail to extract: could not extract full data payload (got %d data bits, expected %d)", bitsExtracted, neededBits)
	}

	return string(extractedData), nil
}
