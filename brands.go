package fatsecret

import (
	"encoding/json"
	"errors"
)

type FoodBrands struct {
	Brands []string `json:"food_brand"`
}

type FoodBrandsResponse struct {
	Brands *FoodBrands    `json:"food_brands,omitempty"`
	Error  *ErrorResponse `json:"error,omitempty"`
}

type BrandType int

const (
	// food brand types
	BrandTypeManufacturer BrandType = iota
	BrandTypeRestaurant
	BrandTypeSupermarket
)

// FoodBrandsByType invokes the FatSecret 'food_brands.get' API call using
// the 'brand_type' parameter and returns the response as a slice of brand strings
func (c *Client) FoodBrandsByType(brandType BrandType) ([]string, error) {
	// determine the brand type name
	var name string
	switch brandType {
	case BrandTypeManufacturer:
		name = "manufacturer"
	case BrandTypeRestaurant:
		name = "restaurant"
	case BrandTypeSupermarket:
		name = "supermarket"
	default:
		name = "manufacturer"
	}

	// invoke the api call
	body, err := c.InvokeAPI(
		"food_brands.get",
		map[string]string{
			"brand_type": name,
		},
	)
	if err != nil {
		return nil, err
	}

	// parse the api response
	brandsResp := FoodBrandsResponse{}
	if err := json.Unmarshal(body, &brandsResp); err != nil {
		return nil, err
	}

	// if an error response was returned
	if brandsResp.Error != nil {
		// return the response error message
		return nil, errors.New(brandsResp.Error.Message)
	}

	return brandsResp.Brands.Brands, nil
}

// FoodBrandsStartingWith invokes the FatSecret 'food_brands.get' API call using
// the 'starts_with' parameter and returns the response as a slice of brand strings
func (c *Client) FoodBrandsStartingWith(startsWith string) ([]string, error) {
	// invoke the api call
	body, err := c.InvokeAPI(
		"food_brands.get",
		map[string]string{
			"starts_with": startsWith,
		},
	)
	if err != nil {
		return nil, err
	}

	// parse the api response
	brandsResp := FoodBrandsResponse{}
	if err := json.Unmarshal(body, &brandsResp); err != nil {
		return nil, err
	}

	// if an error response was returned
	if brandsResp.Error != nil {
		// return the response error message
		return nil, errors.New(brandsResp.Error.Message)
	}

	return brandsResp.Brands.Brands, nil
}
