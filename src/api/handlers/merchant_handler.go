package handlers

import (
	"belimang/src/api/presenter"

	"belimang/src/pkg/entities"
	"belimang/src/pkg/merchant"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// validator instance
var validateMerchant = validator.New()
var validateItems = validator.New()

// CreateMerchant is handler/controller which creates merchant
// @Summary      Create a new merchant
// @Description  Add a new merchant to the collection
// @Tags         Merchants
// @Accept       json
// @Produce      json
// @Param        merchant  body      entities.RequestMerchant  true  "Merchant object"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/admin/merchants [post]
func CreateMerchant(service merchant.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.RequestMerchant

		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body: " + err.Error()))
		}

		if errVal := validateMerchant.Struct(requestBody); errVal != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errVal.Error()))
		}

		// Validate the request body
		new_data := entities.Merchant{
			Name:             requestBody.Name,
			ImageUrl:         requestBody.ImageUrl,
			Lat:              requestBody.Location.Lat,
			Long:             requestBody.Location.Long,
			MerchantCategory: requestBody.MerchantCategory,
		}

		result, err := service.InsertMerchant(&new_data)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err.Error()))
		}
		return c.JSON(
			map[string]interface{}{
				"merchantId": result.ID,
			},
		)
	}
}

// CreateMerchantItems is handler/controller which creates an item for a merchant
// @Summary      Create a new item for a merchant
// @Description  Add a new item under a specific merchant
// @Tags         Merchants
// @Accept       json
// @Produce      json
// @Param        merchantId  path      string                true  "Merchant ID (UUID)"
// @Param        item        body      entities.RequestItems true  "New Item"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /api/v1/admin/merchants/{merchantId}/items [post]
func CreateMerchantItems(service merchant.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var requestBody entities.RequestItems
		merchantId := c.Params("merchantId")
		if err := c.BodyParser(&requestBody); err != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse("invalid request body: " + err.Error()))
		}

		if errVal := validateItems.Struct(requestBody); errVal != nil {
			return c.Status(http.StatusBadRequest).
				JSON(presenter.ErrorResponse(errVal.Error()))
		}

		merchantIdUUID, err := uuid.Parse(merchantId)
		if err != nil {
			// handle the error
			return c.Status(http.StatusBadRequest).JSON(presenter.ErrorResponse(err.Error()))
		}

		new_data := entities.Items{
			Name:            requestBody.Name,
			ImageUrl:        requestBody.ImageUrl,
			Price:           requestBody.Price,
			ProductCategory: requestBody.ProductCategory,
			MerchantID:      merchantIdUUID,
		}

		result, err := service.CreateItems(&new_data, merchantIdUUID)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return c.JSON(presenter.ErrorResponse(err.Error()))
		}
		return c.JSON(
			map[string]interface{}{
				"itemsId": result.ID,
			},
		)
	}
}
