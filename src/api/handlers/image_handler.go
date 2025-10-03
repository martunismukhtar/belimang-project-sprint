package handlers

import (
	"belimang/src/pkg/image"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// UploadImage handles image upload
// @Summary Upload an image
// @Description Upload an image file (jpg/jpeg only, 10KB-2MB)
// @Tags Image
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file"
// @Success 200 {object} image.UploadResponse
// @Failure 400 {object} image.ErrorResponse
// @Failure 401 {object} image.ErrorResponse
// @Security BearerAuth
// @Router /image [post]
func UploadImage(imageService image.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the file from the form
		fileHeader, err := c.FormFile("file")
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(image.ErrorResponse{
				Status:  false,
				Message: "image is wrong (not *jpg | *jpeg, more than 2MB or less than 10KB)",
				Error:   "file not found in request",
			})
		}

		// Upload the image
		imageURL, err := imageService.UploadImage(fileHeader)
		if err != nil {
			return c.Status(http.StatusBadRequest).JSON(image.ErrorResponse{
				Status:  false,
				Message: "image is wrong (not *jpg | *jpeg, more than 2MB or less than 10KB)",
				Error:   err.Error(),
			})
		}

		// Return success response
		return c.Status(http.StatusOK).JSON(image.UploadResponse{
			Message: "File uploaded successfully",
			Data: image.UploadData{
				ImageURL: imageURL,
			},
		})
	}
}