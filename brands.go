package fatsecret

import (
	"encoding/json"
	"errors"
)

// FoodBrands is a component of the 'food_brands.get' API response data
type FoodBrands struct {
	Brands []string `json:"food_brand"`
}

// FoodBrandsResponse is the response format of the 'food_brands.get' API call
type FoodBrandsResponse struct {
	Brands *FoodBrands    `json:"food_brands,omitempty"`
	Error  *ErrorResponse `json:"error,omitempty"`
}

// BrandType is the enum type
type BrandType int

const (
	// BrandTypeManufacturer is for the 'manufacturer' API brand type
	BrandTypeManufacturer BrandType = iota
	// BrandTypeRestaurant is for the 'restaurant' API brand type
	BrandTypeRestaurant
	// BrandTypeSupermarket is for the 'supermarket' API brand type
	BrandTypeSupermarket
)

// FoodBrandsByType invokes the FatSecret 'food_brands.get' API call using
// the 'brand_type' parameter and returns the response as a slice of brand strings
func (c *Client) FoodBrandsByType(brandType BrandType) ([]string, error) {
	// get the brand type name string
	name := brandTypeName(brandType)

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

// brandTypeName converts the given enum type into the
// associated string name
func brandTypeName(brandType BrandType) string {
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
		// use 'manufacturer' by default, like the API does
		name = "manufacturer"
	}

	// return the brand type name
	return name
}
