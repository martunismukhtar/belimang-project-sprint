package purchase

import (
	"belimang/src/pkg/dtos"
	"belimang/src/pkg/entities"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository interface allows us to access the CRUD Operations here.
type Repository interface {
	NearbyMerchant(lat, long float64, params map[string]interface{}) ([]entities.Merchant, error)
	GetItemsByMerchantIDs(merchantIDs []uuid.UUID) (map[uuid.UUID][]entities.Items, error)
	FindMerchantById(merchantIDs []uuid.UUID) ([]entities.Merchant, error)
	simpanEstimate(req entities.DeliveryEstimate) (*entities.DeliveryEstimate, error)
	FindItemsById(itemIDs []uuid.UUID) ([]entities.Items, error)
	FindEstimateById(estimateID uuid.UUID) (*entities.DeliveryEstimate, error)
	SimpanOrders(req entities.DeliveryEstimate, userID uuid.UUID, TotalHarga float64) (string, error)
	FindOrders(req map[string]interface{}) ([]dtos.OrderDetail, error)
}
type repository struct {
	DB *gorm.DB
}

// NewRepo is the single instance repo that is being created.
func NewRepo(db *gorm.DB) Repository {
	return &repository{
		DB: db,
	}
}

// CreateBook is a GORM repository that helps to create books

func (r *repository) NearbyMerchant(userLat, userLong float64, params map[string]interface{}) ([]entities.Merchant, error) {
	var results []entities.Merchant
	//menghitung Haversine
	query := r.DB.Model(&entities.Merchant{}).
		Table("merchants m").
		Select(`m.*,
			(6371 * 2 * ASIN(SQRT(
				POWER(SIN(RADIANS(lat - ?) / 2),2) + 
				COS(RADIANS(?)) * COS(RADIANS(lat)) * 
				POWER(SIN(RADIANS(long - ?) / 2), 2)
			))) AS distance_km
		`, userLat, userLat, userLong)

	// Apply filters dynamically from params
	if v, ok := params["merchantId"]; ok {
		if v != "" {
			switch val := v.(type) {
			case string:
				id, err := uuid.Parse(val)
				if err == nil {
					merchantId := []uuid.UUID{id}
					query = query.Where("m.id IN ?", merchantId)

				} else {
					return nil, fmt.Errorf("invalid merchantId: %s", val)
				}
			case uuid.UUID:
				merchantId := []uuid.UUID{val}
				query = query.Where("m.id IN ?", merchantId)

			case []uuid.UUID:
				query = query.Where("m.id IN ?", val)

			default:
				return nil, fmt.Errorf("invalid merchantId: %v", v)
			}
		}
	}

	if v, ok := params["name"]; ok {
		if v != "" {
			name := "%" + strings.ToLower(v.(string)) + "%"
			query = query.Where("(LOWER(m.name) LIKE ? OR EXISTS ("+
				"SELECT 1 FROM items i WHERE i.merchant_id = m.id AND LOWER(i.name) LIKE ?"+
				"))", name, name)
		}
	}
	if v, ok := params["merchantCategory"]; ok {

		if v != "" {
			fmt.Println("merchantCategory", v)
			query = query.Where("merchant_category = ?", v)
		}
	}

	// Order by nearest & limit
	err := query. //Offset(params["offset"].(int)).Limit(params["limit"].(int)).
			Order("distance_km ASC").Scan(&results).Error
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *repository) GetItemsByMerchantIDs(merchantIDs []uuid.UUID) (map[uuid.UUID][]entities.Items, error) {
	var items []entities.Items

	err := r.DB.Where("merchant_id IN ?", merchantIDs).Find(&items).Error
	if err != nil {
		return nil, err
	}

	// Group items by merchant_id
	itemsByMerchant := make(map[uuid.UUID][]entities.Items)
	for _, item := range items {
		itemsByMerchant[item.MerchantID] = append(itemsByMerchant[item.MerchantID], item)
	}

	return itemsByMerchant, nil
}

func (r *repository) FindMerchantById(merchantIDs []uuid.UUID) ([]entities.Merchant, error) {
	var merchant []entities.Merchant
	if err := r.DB.
		Table("merchants").
		Where("id IN ?", merchantIDs).Scan(&merchant).Error; err != nil {
		return nil, err
	}
	return merchant, nil
}
func (r *repository) FindItemsById(itemIDs []uuid.UUID) ([]entities.Items, error) {
	var items []entities.Items
	if err := r.DB.
		Table("items").
		Where("id IN ?", itemIDs).Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) FindEstimateById(estimateID uuid.UUID) (*entities.DeliveryEstimate, error) {
	var estimate entities.DeliveryEstimate

	if err := r.DB.
		Table("delivery_estimate").
		Where("id = ?", estimateID).
		First(&estimate).Error; err != nil {
		return nil, err
	}
	return &estimate, nil
}

func (r *repository) simpanEstimate(req entities.DeliveryEstimate) (*entities.DeliveryEstimate, error) {
	err := r.DB.Create(&req).Error
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func (r *repository) SimpanOrders(req entities.DeliveryEstimate, userID uuid.UUID, TotalHarga float64) (string, error) {
	OrderID := uuid.New()
	// UserID := userID
	var wrappers []entities.OrderWrapper
	if err := json.Unmarshal(req.Orders, &wrappers); err != nil {
		return "", err
	}

	err := r.DB.Transaction(func(tx *gorm.DB) error {
		var allOrders []entities.Order
		var allOrderItems []entities.OrderItem

		orderID := OrderID
		order := entities.Order{
			ID:         orderID,
			UserID:     userID,
			TotalPrice: TotalHarga, // bisa dihitung kalau ada harga item
			CreateAt:   time.Now().Unix(),
		}
		allOrders = append(allOrders, order)

		for _, w := range wrappers {
			// siapkan order_items untuk batch insert
			for _, item := range w.Items {
				orderItem := entities.OrderItem{
					ID:         uuid.New(),
					MerchantID: w.MerchantID,
					OrderID:    orderID,
					ItemID:     item.ItemID,
					Quantity:   item.Quantity,
				}
				allOrderItems = append(allOrderItems, orderItem)
			}
		}

		// batch insert orders
		if err := tx.Create(&allOrders).Error; err != nil {
			return err
		}

		// batch insert order_items
		if err := tx.Create(&allOrderItems).Error; err != nil {
			return err
		}

		return nil
	})
	return OrderID.String(), err
}

func (r *repository) FindOrders(params map[string]interface{}) ([]dtos.OrderDetail, error) {

	var orders []dtos.OrderDetail
	query := r.DB.
		Table("orders AS o").
		Select(`o.id AS order_id, 
	        oi.merchant_id, 
	        m.name AS merchant_name, 
	        m.merchant_category, 
	        m.image_url AS merchant_image_url, 
	        m.long AS merchant_long, 
	        m.lat AS merchant_lat, 
			m.created_at AS merchant_created_at,
	        oi.item_id, 
	        i.name AS item_name, 
	        i.product_category, 
	        i.price, 
	        oi.quantity, 
	        i.image_url AS item_image_url, 
			i.created_at AS item_created_at`).
		Joins("JOIN order_items oi ON o.id = oi.order_id").
		Joins("JOIN merchants m ON m.id = oi.merchant_id").
		Joins("JOIN items i ON oi.item_id = i.id")

	if v, ok := params["merchantId"]; ok {
		if v != "" {
			switch val := v.(type) {
			case string:
				id, err := uuid.Parse(val)
				if err == nil {
					merchantId := []uuid.UUID{id}
					query = query.Where("m.id IN ?", merchantId)

				} else {
					return nil, fmt.Errorf("invalid merchantId: %s", val)
				}
			case uuid.UUID:
				merchantId := []uuid.UUID{val}
				query = query.Where("m.id IN ?", merchantId)

			case []uuid.UUID:
				query = query.Where("m.id IN ?", val)

			default:
				return nil, fmt.Errorf("invalid merchantId: %v", v)
			}
		}
	}

	if v, ok := params["name"]; ok {
		if v != "" {
			name := "%" + strings.ToLower(v.(string)) + "%"
			query = query.Where("(LOWER(m.name) LIKE ? OR EXISTS ("+
				"SELECT 1 FROM items i WHERE i.merchant_id = m.id AND LOWER(i.name) LIKE ?"+
				"))", name, name)
		}
	}
	if v, ok := params["merchantCategory"]; ok {
		if v != "" {
			query = query.Where("m.merchant_category = ?", v)
		}
	}
	err := query.Where("o.user_id = ?", params["userId"]).
		Scan(&orders).Error

	if err != nil {
		return nil, err
	}

	return orders, nil
}
