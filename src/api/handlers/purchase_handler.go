package handlers

import (
	"belimang/src/pkg/entities"
	"belimang/src/pkg/purchase"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// FindNearbyMerchant godoc
// @Summary      Get nearby merchants
// @Description  Get nearest merchants based on user's lat & lon
// @Tags         Purchase
// @Accept       json
// @Produce      json
// @Param        lat   path   number  true  "Latitude (-90 to 90)"
// @Param        lon   path   number  true  "Longitude (-180 to 180)"
// @Param        merchantId  query  string  false "Merchant ID"
// @Param        name  query  string  false "Merchant name"
// @Param        merchantCategory  query  string  false "Merchant category"
// @Param        limit   query  int    false "Limit results (default: 5)"
// @Param        offset  query  int    false "Pagination offset (default: 0)"
// @Security     BearerAuth
// @Success      200  {array}   map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/merchants/nearby/{lat}/{lon} [get]
func FindNearbyMerchant(service purchase.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userID := c.Locals("user_id")
		fmt.Println(userID)
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		latStr := c.Params("lat")
		lonStr := c.Params("lon")
		lat, err1 := strconv.ParseFloat(latStr, 64)
		lon, err2 := strconv.ParseFloat(lonStr, 64)
		if err1 != nil || err2 != nil || lat < -90 || lat > 90 || lon < -180 || lon > 180 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "lat/long is not valid",
			})
		}

		// Query params: limit & offset
		limitParam := c.Query("limit", "5")
		offsetParam := c.Query("offset", "0")
		merchantId := c.Query("merchantId")
		name := c.Query("name")
		merchantCategory := c.Query("merchantCategory")

		limit, _ := strconv.Atoi(limitParam)
		offset, _ := strconv.Atoi(offsetParam)

		if limit <= 0 {
			limit = 5
		}

		data_merchants, errn := service.NearbyMerchant(lat, lon, map[string]interface{}{
			"limit":            limit,
			"offset":           offset,
			"merchantId":       merchantId,
			"name":             name,
			"merchantCategory": merchantCategory,
		})

		if errn != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch nearby merchants",
			})
		}

		return c.Status(fiber.StatusOK).JSON(data_merchants)
	}
}

// Estimate godoc
// @Summary      calculate estimate time
// @Description  calculate estimate time
// @Tags         Purchase
// @Accept       json
// @Produce      json
// @Param        request body entities.EstimateRequest true "Estimate request body"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/users/estimate [post]
func Estimate(service purchase.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		var req entities.EstimateRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		// === Validasi ===
		// validasi koordinat
		if req.UserLocation.Lat < -90 || req.UserLocation.Lat > 90 || req.UserLocation.Long < -180 || req.UserLocation.Long > 180 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "lat/long is not valid",
			})
		}

		countStarting := 0
		for _, o := range req.Orders {

			if o.MerchantID == "" {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "merchantId is required",
				})
			}
			// Pastikan hanya ada satu starting point
			if o.IsStartingPoint {
				countStarting++
			}

			for _, i := range o.Items {
				if i.ItemID == "" {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "itemId is required",
					})
				}
				if i.Quantity < 1 {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error": "quantity must be greater than 0",
					})
				}
			}
		}
		if countStarting != 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "there must be exactly one order with isStartingPoint = true",
			})
		}

		user := uuid.MustParse(userID.(string))
		est, errEs := service.Estimate(req, user)
		if errEs != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": errEs.Error(),
			})
		}
		// fmt.Println(est)
		// TODO: panggil service untuk hitung jarak, TSP, estimasi waktu, dll.
		result := map[string]interface{}{
			"totalPrice":                   est.TotalPrice,
			"estimatedDeliveryTimeMinutes": fmt.Sprintf("%.2f", est.EstimatedDelivery),
			"calculatedEstimateId":         est.ID,
		}

		return c.Status(fiber.StatusOK).JSON(result)
	}
}

// Order godoc
// @Summary      order
// @Description  order based on estimate time id
// @Tags         Purchase
// @Accept       json
// @Produce      json
// @Param        request body entities.OrderRequest true "Estimate request body"
// @Security     BearerAuth
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/users/orders [post]
func Order(service purchase.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {

		userID := c.Locals("user_id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		var req entities.OrderRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		// TODO: panggil service untuk hitung jarak, TSP, estimasi waktu, dll.
		result, err := service.Order(req, uuid.MustParse(userID.(string)))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		data := map[string]interface{}{
			"orderId": result,
		}
		return c.Status(fiber.StatusOK).JSON(data)
	}
}

// GetOrder godoc
// @Summary      Get orders
// @Description  Get orders
// @Tags         Purchase
// @Accept       json
// @Produce      json
// @Param        merchantId  query  string  false "Merchant ID"
// @Param        name  query  string  false "Merchant name"
// @Param        merchantCategory  query  string  false "Merchant category"
// @Param        limit   query  int    false "Limit results (default: 5)"
// @Param        offset  query  int    false "Pagination offset (default: 0)"
// @Security     BearerAuth
// @Success      200  {array}   map[string]interface{}
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /api/v1/users/orders [get]
func GetOrder(service purchase.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		if userID == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "unauthorized",
			})
		}

		// Query params: limit & offset
		limitParam := c.Query("limit", "5")
		offsetParam := c.Query("offset", "0")
		merchantId := c.Query("merchantId")
		name := c.Query("name")
		merchantCategory := c.Query("merchantCategory")

		limit, _ := strconv.Atoi(limitParam)
		offset, _ := strconv.Atoi(offsetParam)

		if limit <= 0 {
			limit = 5
		}

		data_merchants, errn := service.GetOrderData(map[string]interface{}{
			"limit":            limit,
			"offset":           offset,
			"merchantId":       merchantId,
			"name":             name,
			"merchantCategory": merchantCategory,
			"userId":           uuid.MustParse(userID.(string)),
		})

		if errn != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to fetch order data",
			})
		}

		return c.Status(fiber.StatusOK).JSON(data_merchants)
	}
}
