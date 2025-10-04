package purchase

import (
	"belimang/src/pkg/dtos"
	"belimang/src/pkg/entities"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
)

// Service is an interface from which our api module can access our repository of all our models
type Service interface {
	NearbyMerchant(lat, long float64, params map[string]interface{}) (*dtos.NearbyMerchantResponse, error)
	Estimate(req entities.EstimateRequest, userID uuid.UUID) (*entities.DeliveryEstimate, error)
	Order(req entities.OrderRequest, userID uuid.UUID) (string, error)
	GetOrderData(req map[string]interface{}) ([]map[string]interface{}, error)
}

type service struct {
	repository Repository
}

// NewService is used to create a single instance of the service
func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

func (s *service) NearbyMerchant(lat, long float64, params map[string]interface{}) (*dtos.NearbyMerchantResponse, error) {
	// Ambil data merchant dari repository
	merchants, err := s.repository.NearbyMerchant(lat, long, params)
	if err != nil {
		return nil, err
	}

	limit := params["limit"].(int)
	offset := params["offset"].(int)
	if v, ok := params["limit"]; ok {
		if l, ok := v.(int); ok && l > 0 {
			limit = l
		}
	}
	if v, ok := params["offset"]; ok {
		if o, ok := v.(int); ok && o >= 0 {
			offset = o
		}
	}

	limit_merchants := limit
	if len(merchants) < limit {
		limit_merchants = len(merchants)
	}
	tspMerchants, _ := NearestNeighborTSP(lat, long, merchants[offset:limit_merchants])

	// Apply offset
	if offset >= len(tspMerchants) {
		return &dtos.NearbyMerchantResponse{
			Data: []dtos.MerchantWithItems{},
			Meta: dtos.MetaResponse{
				Limit:  limit,
				Offset: offset,
				Total:  10,
			},
		}, nil
	}

	if offset > 0 {
		tspMerchants = tspMerchants[offset:]
	}
	// Get merchant IDs
	merchantIDs := make([]uuid.UUID, len(tspMerchants))
	for i, m := range tspMerchants {
		merchantIDs[i] = m.ID
	}
	// Get items for these merchants
	itemsByMerchant, err := s.repository.GetItemsByMerchantIDs(merchantIDs)
	if err != nil {
		return nil, err
	}
	// Build response
	data := make([]dtos.MerchantWithItems, len(tspMerchants))
	for i, merchant := range tspMerchants {
		// Convert merchant to response format
		merchantResp := dtos.MerchantResponse{
			MerchantID:       uuid.MustParse(merchant.ID.String()),
			Name:             merchant.Name,
			MerchantCategory: string(merchant.MerchantCategory),
			ImageURL:         merchant.ImageUrl,
			Location: dtos.LocationResponse{
				Lat:  merchant.Lat,
				Long: merchant.Long,
			},
			CreatedAt: formatNanosToISO8601(merchant.CreatedAt),
		}

		// Convert items to response format
		items := itemsByMerchant[merchant.ID]
		itemsResp := make([]dtos.ItemResponse, len(items))
		for j, item := range items {
			itemsResp[j] = dtos.ItemResponse{
				ItemID:          uuid.MustParse(item.ID.String()),
				Name:            item.Name,
				ProductCategory: string(item.ProductCategory),
				Price:           item.Price,
				ImageURL:        item.ImageUrl,
				CreatedAt:       formatNanosToISO8601(item.CreatedAt),
			}
		}

		data[i] = dtos.MerchantWithItems{
			Merchant: merchantResp,
			Items:    itemsResp,
		}
	}
	return &dtos.NearbyMerchantResponse{
		Data: data,
		Meta: dtos.MetaResponse{
			Limit:  limit,
			Offset: offset,
			Total:  len(merchants),
		},
	}, nil
}

func Haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Radius bumi dalam kilometer

	// Konversi derajat ke radian
	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	// Perbedaan koordinat
	dLat := lat2Rad - lat1Rad
	dLon := lon2Rad - lon1Rad

	// Formula Haversine
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c

}

func NearestNeighborTSP(userLat, userLong float64, merchants []entities.Merchant) ([]entities.Merchant, float64) {
	if len(merchants) == 0 {
		return []entities.Merchant{}, 0
	}

	visited := make(map[int]bool)
	route := make([]entities.Merchant, 0, len(merchants))
	totalDistance := 0.0

	currentLat := userLat
	currentLong := userLong

	for len(visited) < len(merchants) {
		nearestIdx := -1
		nearestDist := math.MaxFloat64

		for i, merchant := range merchants {
			if !visited[i] {
				dist := Haversine(currentLat, currentLong, merchant.Lat, merchant.Long)
				if dist < nearestDist {
					nearestDist = dist
					nearestIdx = i
				}
			}
		}

		if nearestIdx != -1 {
			visited[nearestIdx] = true
			route = append(route, merchants[nearestIdx])
			totalDistance += nearestDist

			currentLat = merchants[nearestIdx].Lat
			currentLong = merchants[nearestIdx].Long
		}
	}

	return route, totalDistance
}

func in_array(arr []string, target string) bool {
	for _, v := range arr {
		if v == target {
			return true
		}
	}
	return false
}

func (s *service) Estimate(req entities.EstimateRequest, userID uuid.UUID) (*entities.DeliveryEstimate, error) {
	merchantIds := make([]uuid.UUID, 0, len(req.Orders)) // kapasitas sesuai jumlah order
	for _, order := range req.Orders {
		id, err := uuid.Parse(order.MerchantID)
		if err != nil {
			return nil, fmt.Errorf("invalid merchantId: %s", order.MerchantID)
		}
		merchantIds = append(merchantIds, id)
	}
	// cek merchantId ada di DB
	merchants, err := s.repository.FindMerchantById(merchantIds)
	if err != nil {
		return nil, err
	}
	//get merchant id from db
	merchantIdsDB := make([]string, 0, len(merchants))
	for _, merchant := range merchants {
		merchantIdsDB = append(merchantIdsDB, merchant.ID.String())
	}

	for _, merchant := range merchantIds {
		if !in_array(merchantIdsDB, merchant.String()) {
			return nil, fmt.Errorf("invalid merchantId: %s", merchant.String())
		}
	}

	//hitung jarak antara user dan merchant starting point
	//jika > 3 km, maka tolak
	for _, mec := range req.Orders {
		if mec.IsStartingPoint {
			for _, toko := range merchants {
				if toko.ID.String() == mec.MerchantID {
					jrk := Haversine(req.UserLocation.Lat, req.UserLocation.Long, toko.Lat, toko.Long)

					if jrk > 3 {
						return nil, fmt.Errorf("starting point too far")
					}
				}
			}
		}
	}

	//htung total harga
	//ambil id items dari order
	ord_items := make(map[string]int)
	var itemsId []uuid.UUID
	for _, order := range req.Orders {
		for _, item := range order.Items {
			ord_items[item.ItemID] = item.Quantity
			itemsId = append(itemsId, uuid.MustParse(item.ItemID))
		}
	}

	//ambil items
	items, err := s.repository.FindItemsById(itemsId)
	if err != nil {
		return nil, err
	}
	totalHarga := 0.0
	for _, item := range items {
		qty := ord_items[item.ID.String()]
		totalHarga += item.Price * float64(qty)
	}

	//jika merchant berjumlah 1, tidak perlu TSP
	//hitung jarak antara user dan merchant
	//dapatkn waktu tempuh
	if len(merchantIds) == 1 {
		for _, merchant := range merchants {
			if merchant.ID.String() == merchantIds[0].String() {
				jrk := Haversine(req.UserLocation.Lat, req.UserLocation.Long, merchant.Lat, merchant.Long)
				waktu_menit := (jrk / 40.0) * 60.0

				ordersJSON, _ := json.Marshal(req.Orders)

				simpan_data := entities.DeliveryEstimate{
					UserID:            uuid.MustParse("550e8400-e29b-41d4-a716-446655440000"),
					Orders:            ordersJSON,
					EstimatedDelivery: waktu_menit,
					TotalPrice:        totalHarga,
				}

				hasilEstimasi, err := s.repository.simpanEstimate(simpan_data)
				if err != nil {
					return nil, err
				}
				return hasilEstimasi, nil
			}
		}
	}

	user_lat := req.UserLocation.Lat
	user_long := req.UserLocation.Long

	nearby_merchants, errNearby := s.repository.FindMerchantById(merchantIds)
	if errNearby != nil {
		return nil, errNearby
	}

	_, distance := NearestNeighborTSP(user_lat, user_long, nearby_merchants)

	waktu_menit := (distance / 40.0) * 60.0
	ordersJSON, _ := json.Marshal(req.Orders)
	simpan_data := entities.DeliveryEstimate{
		UserID:            userID,
		Orders:            ordersJSON,
		EstimatedDelivery: waktu_menit,
		TotalPrice:        totalHarga,
	}

	hasilEstimasi, err := s.repository.simpanEstimate(simpan_data)
	if err != nil {
		return nil, err
	}
	return hasilEstimasi, nil

}

func (s *service) Order(req entities.OrderRequest, userID uuid.UUID) (string, error) {

	estID := req.CalculatedEstimateId
	est, err := s.repository.FindEstimateById(uuid.MustParse(estID))
	if err != nil {
		return "", err
	}

	var wrappers []entities.OrderWrapper
	if err := json.Unmarshal(est.Orders, &wrappers); err != nil {
		return "", err
	}

	//hitung total harga
	ordItems := make(map[uuid.UUID]int)
	var itemsId []uuid.UUID
	for _, wp := range wrappers {
		for _, item := range wp.Items {
			ordItems[item.ItemID] = item.Quantity
			itemsId = append(itemsId, item.ItemID)
		}
	}

	//ambil items
	items, err := s.repository.FindItemsById(itemsId)
	if err != nil {
		return "", err
	}
	totalHarga := 0.0
	for _, item := range items {
		qty := ordItems[item.ID]
		totalHarga += item.Price * float64(qty)
	}

	OrderData, errSimpanOrder := s.repository.SimpanOrders(*est, userID, totalHarga)
	if errSimpanOrder != nil {
		return "", errSimpanOrder
	}
	return OrderData, nil
}
func formatNanosToISO8601(nanos int64) string {
	sec := nanos / 1e9
	nsec := nanos % 1e9
	t := time.Unix(sec, nsec)
	// Format manual agar tetap ada 9 digit nanodetik
	return t.Format("2006-01-02T15:04:05.000000000Z07:00")
}
func (s *service) GetOrderData(params map[string]interface{}) ([]map[string]interface{}, error) {
	ordersDB, err := s.repository.FindOrders(params)
	if err != nil {
		return nil, err
	}

	orderMap := map[string]map[string]interface{}{}
	for _, row := range ordersDB {
		orderID := row.OrderID.String()
		if _, exists := orderMap[orderID]; !exists {
			orderMap[orderID] = map[string]interface{}{
				"orderId": orderID,
				"orders":  []map[string]interface{}{},
			}
		}

		orders := orderMap[orderID]["orders"].([]map[string]interface{})
		// Cari merchant yang sudah ada di orders
		var existingOrder *map[string]interface{}
		for i := range orders {
			merchant := orders[i]["merchant"].(map[string]interface{})
			if merchant["merchantId"] == row.MerchantID {
				existingOrder = &orders[i]
				break
			}
		}

		item := map[string]interface{}{
			"itemId":          row.ItemID,
			"name":            row.ItemName,
			"productCategory": row.ProductCategory,
			"price":           row.Price,
			"quantity":        row.Quantity,
			"imageUrl":        row.ItemImageURL,
			"createdAt":       formatNanosToISO8601(row.ItemCreatedAt),
		}

		if existingOrder != nil {
			// Tambahkan item ke merchant yang sudah ada
			existingItems := (*existingOrder)["items"].([]map[string]interface{})
			(*existingOrder)["items"] = append(existingItems, item)
		} else {
			// Buat merchant baru
			merchant := map[string]interface{}{
				"merchantId":       row.MerchantID,
				"name":             row.MerchantName,
				"merchantCategory": row.MerchantCategory,
				"imageUrl":         row.MerchantImageURL,
				"location": map[string]interface{}{
					"lat":  row.MerchantLat,
					"long": row.MerchantLong,
				},
				"createdAt": formatNanosToISO8601(row.MerchantCreatedAt),
			}

			newOrder := map[string]interface{}{
				"merchant": merchant,
				"items":    []map[string]interface{}{item},
			}

			orderMap[orderID]["orders"] = append(orders, newOrder)
		}
	}

	orderData := make([]map[string]interface{}, 0, len(orderMap))
	for _, order := range orderMap {
		orderData = append(orderData, order)
	}

	return orderData, nil

}
