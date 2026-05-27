package upload

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp" // register WebP decoder
)

const (
	maxWidth   = 1200 // Max pixel width — sufficient for profile/cover images
	maxHeight  = 1200
	jpegQuality = 80  // JPEG output quality (0-100)
)

// CompressResult contains the compressed image data and its new content type.
type CompressResult struct {
	Data        io.Reader
	ContentType string
	Size        int
}

// CompressImage decodes an image, resizes it if larger than maxWidth/maxHeight,
// and re-encodes as JPEG at reduced quality. This reduces file size significantly
// before uploading to Cloudinary.
//
// Supported input formats: JPEG, PNG, WebP (via x/image/webp decoder).
// AVIF and GIF are passed through without compression (no native Go decoder).
func CompressImage(src io.Reader, contentType string) (*CompressResult, error) {
	// AVIF and GIF: pass through as-is (no native Go encoder/decoder)
	if contentType == "image/avif" || contentType == "image/gif" {
		data, err := io.ReadAll(src)
		if err != nil {
			return nil, fmt.Errorf("failed to read file: %w", err)
		}
		return &CompressResult{
			Data:        bytes.NewReader(data),
			ContentType: contentType,
			Size:        len(data),
		}, nil
	}

	// Decode the image (jpeg, png, webp are auto-detected)
	img, _, err := image.Decode(src)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	// Resize if needed
	img = resizeIfNeeded(img)

	// Re-encode as JPEG for lossy compression (best size reduction)
	var buf bytes.Buffer

	if contentType == "image/png" {
		// Keep PNG for transparency support, use default compression
		err = png.Encode(&buf, img)
	} else {
		// Everything else → JPEG at controlled quality
		err = jpeg.Encode(&buf, img, &jpeg.Options{Quality: jpegQuality})
		contentType = "image/jpeg"
	}

	if err != nil {
		return nil, fmt.Errorf("failed to encode image: %w", err)
	}

	return &CompressResult{
		Data:        bytes.NewReader(buf.Bytes()),
		ContentType: contentType,
		Size:        buf.Len(),
	}, nil
}

// resizeIfNeeded scales the image down proportionally if either dimension
// exceeds the max. Uses high-quality CatmullRom interpolation.
func resizeIfNeeded(img image.Image) image.Image {
	bounds := img.Bounds()
	w := bounds.Dx()
	h := bounds.Dy()

	if w <= maxWidth && h <= maxHeight {
		return img // No resize needed
	}

	// Calculate new dimensions preserving aspect ratio
	ratio := float64(w) / float64(h)
	newW, newH := maxWidth, maxHeight

	if ratio > 1 {
		// Landscape: constrain by width
		newH = int(float64(newW) / ratio)
	} else {
		// Portrait: constrain by height
		newW = int(float64(newH) * ratio)
	}

	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, bounds, draw.Over, nil)
	return dst
}
