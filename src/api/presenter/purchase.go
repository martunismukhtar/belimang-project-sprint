package presenter

import (
	"belimang/src/pkg/dtos"
	"belimang/src/pkg/entities"
	"fmt"
)

func BuildNearbyMerchantResponse(merchants []entities.Merchant, items map[string][]entities.Items, limit, offset, total int) dtos.NearbyMerchantResponse {
	var data []dtos.MerchantWithItems

	for _, m := range merchants {
		merchantRes := dtos.MerchantResponse{
			MerchantID:       m.ID,
			Name:             m.Name,
			MerchantCategory: string(m.MerchantCategory),
			ImageURL:         m.ImageUrl,
			Location: dtos.LocationResponse{
				Lat:  m.Lat,
				Long: m.Long,
			},
			// CreatedAt: m.CreatedAt,
		}

		var itemRes []dtos.ItemResponse
		for _, i := range items[fmt.Sprintf("%s", m.ID)] {
			itemRes = append(itemRes, dtos.ItemResponse{
				ItemID:          i.ID,
				Name:            i.Name,
				ProductCategory: string(i.ProductCategory),
				Price:           i.Price,
				ImageURL:        i.ImageUrl,
				// CreatedAt:       i.CreatedAt,
			})
		}

		data = append(data, dtos.MerchantWithItems{
			Merchant: merchantRes,
			Items:    itemRes,
		})
	}

	return dtos.NearbyMerchantResponse{
		Data: data,
		Meta: dtos.MetaResponse{
			Limit:  limit,
			Offset: offset,
			Total:  total,
		},
	}
}
