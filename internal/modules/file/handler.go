package file

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"path/filepath"
	"strings"

	"github.com/gen2brain/avif"
	"github.com/gofiber/fiber/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// UploadFilesHanlder godoc
// @Summary Upload multiple files
// @Description Uploads multiple files and saves metadata in the database
// @Tags files
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param files formData []file true "Files to upload" collectionFormat(multi)
// @Success 200 {object} map[string]interface{} "Returns file metadata"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /files/upload [post]
func UploadFilesHanlder(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse multipart form",
		})
	}
	files := form.File["files"]

	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot connect to MinIO",
		})
	}

	for _, fileHeader := range files {

		file, err := fileHeader.Open()
		if err != nil {
			continue
		}
		defer file.Close()

		var img image.Image
		ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
		switch ext {
		case ".png":
			img, err = png.Decode(file)
		case ".jpg", ".jpeg":
			img, err = jpeg.Decode(file)
		default:
			continue
		}
		if err != nil {
			continue
		}

		var buf bytes.Buffer
		if err := avif.Encode(&buf, img); err != nil {
			continue
		}

		objectName := fmt.Sprintf("%s.avif", strings.TrimSuffix(fileHeader.Filename, ext))

		_, err = minioClient.PutObject(
			c.Context(),
			"go-bucket",
			objectName,
			bytes.NewReader(buf.Bytes()),
			int64(buf.Len()),
			minio.PutObjectOptions{ContentType: "image/avif"},
		)
		if err != nil {
			continue
		}

	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Files uploaded successfully",
	})
}

// DeleteFilesHandler godoc
// @Summary Delete multiple files
// @Description Delete multiple files
// @Tags files
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param payload body DeleteFilesParams true "Keys"
// @Success 200 {object} map[string]interface{} "Returns file metadata"
// @Failure 400 {object} map[string]string "Bad request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /files/delete [delete]
func DeleteFilesHandler(c *fiber.Ctx) error {
	var req DeleteFilesParams

	if err := c.BodyParser(&req); err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	minioClient, err := minio.New("localhost:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	ctx := context.Background()
	objectCh := make(chan minio.ObjectInfo)
	go func() {
		defer close(objectCh)
		for _, key := range req.Keys {
			objectCh <- minio.ObjectInfo{Key: key}
		}
	}()

	minioClient.RemoveObjects(ctx, "go-bucket", objectCh, minio.RemoveObjectsOptions{})
	return c.SendStatus(fiber.StatusOK)
}
