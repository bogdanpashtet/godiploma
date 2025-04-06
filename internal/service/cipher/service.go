package cipher

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"image"
	"image/png"

	"golang.org/x/image/bmp"

	cipherd "github.com/bogdanpashtet/godiploma/internal/domain/cipher"
	filed "github.com/bogdanpashtet/godiploma/internal/domain/file"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const kind = "CipherService"

type Service struct {
	l *zap.Logger
}

func New(
	l *zap.Logger,
) *Service {
	return &Service{
		l: l,
	}
}

func (s *Service) CreateStegoImage(ctx context.Context, req cipherd.CreateStegoImageRequest) ([]filed.File, error) {
	s.l = s.l.With(zap.String("kind", kind))

	resFiles := make([]filed.File, 0, len(req.Files))
	for _, file := range req.Files {
		switch req.Method {
		case cipherd.MethodLSB:
			res, err := cipherLSB(ctx, req.Plaintext, file)
			if err != nil {
				return nil, fmt.Errorf("fail to cipherLSB: %w", err)
			}
			resFiles = append(resFiles, res)
		default:
		}
	}

	return resFiles, nil
}

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
