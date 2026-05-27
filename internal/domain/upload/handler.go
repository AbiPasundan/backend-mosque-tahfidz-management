package upload

import (
	"backend-mosque-tahfidz-management/pkg/utils"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog/log"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/webp": true,
	"image/avif": true,
	"image/gif":  true,
}

const maxFileSize = 5 * 1024 * 1024 // 5MB

type UploadHandler struct {
	cloudinary *CloudinaryService
}

func NewUploadHandler(cloudinary *CloudinaryService) *UploadHandler {
	return &UploadHandler{cloudinary: cloudinary}
}

func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusBadRequest, "file is required")
	}

	// Validate size
	if file.Size > maxFileSize {
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("file too large (%.1fMB), max %dMB", float64(file.Size)/(1024*1024), maxFileSize/(1024*1024)))
	}

	// Validate content type
	contentType := file.Header.Get("Content-Type")
	if !allowedMimeTypes[contentType] {
		allowed := make([]string, 0, len(allowedMimeTypes))
		for k := range allowedMimeTypes {
			allowed = append(allowed, k)
		}
		return utils.ErrorResponse(c, fiber.StatusBadRequest,
			fmt.Sprintf("unsupported file type: %s. Allowed: %s", contentType, strings.Join(allowed, ", ")))
	}

	// Open the file
	src, err := file.Open()
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, "failed to read file")
	}
	defer src.Close()

	// Compress the image (resize + re-encode at lower quality)
	compressed, err := CompressImage(src, contentType)
	if err != nil {
		log.Warn().Err(err).Str("filename", file.Filename).Msg("compression failed, uploading original")
		// Fallback: re-open and upload original if compression fails
		src2, _ := file.Open()
		defer src2.Close()

		result, uploadErr := h.cloudinary.Upload(src2, file.Filename)
		if uploadErr != nil {
			return utils.ErrorResponse(c, fiber.StatusInternalServerError, uploadErr.Error())
		}
		return utils.SuccessResponse(c, fiber.StatusOK, "file uploaded successfully (uncompressed)", fiber.Map{
			"url":       result.SecureURL,
			"public_id": result.PublicID,
			"format":    result.Format,
			"size":      result.Bytes,
			"width":     result.Width,
			"height":    result.Height,
		})
	}

	log.Info().
		Str("filename", file.Filename).
		Int64("original_size", file.Size).
		Int("compressed_size", compressed.Size).
		Str("output_type", compressed.ContentType).
		Msgf("compressed image: %.0f%% reduction", (1-float64(compressed.Size)/float64(file.Size))*100)

	// Upload compressed data to Cloudinary
	result, err := h.cloudinary.Upload(compressed.Data, file.Filename)
	if err != nil {
		return utils.ErrorResponse(c, fiber.StatusInternalServerError, err.Error())
	}

	return utils.SuccessResponse(c, fiber.StatusOK, "file uploaded successfully", fiber.Map{
		"url":             result.SecureURL,
		"public_id":       result.PublicID,
		"format":          result.Format,
		"size":            result.Bytes,
		"width":           result.Width,
		"height":          result.Height,
		"original_size":   file.Size,
		"compressed_size": compressed.Size,
	})
}
