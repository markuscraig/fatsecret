package fatsecret

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	FatSecretBarcodeLength = 13
)

type Food struct {
	ID          string `json:"food_id,omitempty"`
	Name        string `json:"food_name,omitempty"`
	Type        string `json:"food_type,omitempty"`
	BrandName   string `json:"brand_name,omitempty"`
	URL         string `json:"food_url,omitempty"`
	Description string `json:"food_description,omitempty"`
}

type FoodSearchResponseFoods struct {
	PageNumber   int    `json:"page_number,string"`
	PageSize     int    `json:"max_results,string"`
	TotalResults int    `json:"total_results,string"`
	Food         []Food `json:"food"`
}

type FoodSearchResponse struct {
	Foods *FoodSearchResponseFoods `json:"foods,omitempty"`
	Error *ErrorResponse           `json:"error,omitempty"`
}

type FoodID struct {
	Value string `json:"value"`
}

type FoodIDResponse struct {
	ID    *FoodID        `json:"food_id,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
}

// FoodSearch invokes the FatSecret 'foods.search' API call and
// returns the response as a slice of Food structs
func (c *Client) FoodSearch(query string) ([]Food, error) {
	// invoke the api call
	body, err := c.InvokeAPI(
		"foods.search",
		map[string]string{
			"search_expression": query,
		},
	)
	if err != nil {
		return nil, err
	}

	// parse the api response
	foodResp := FoodSearchResponse{}
	if err := json.Unmarshal(body, &foodResp); err != nil {
		return nil, err
	}

	// if an error response was returned
	if foodResp.Error != nil {
		// return the response error message
		return nil, errors.New(foodResp.Error.Message)
	}

	// return the slice of food items
	return foodResp.Foods.Food, nil
}

// FoodIDForBarcode invokes the FatSecret 'food.find_id_for_barcode' API call and
// returns the response as a slice of Food structs
func (c *Client) FoodIDForBarcode(barcode string) (string, error) {
	// if the barcode is invalid
	barcodeLen := len(barcode)
	if barcodeLen == 0 || barcodeLen > FatSecretBarcodeLength {
		return "", fmt.Errorf("Invalid barcode length '%d' given", barcodeLen)
	}

	// pad the barcode to GTIN-13 format
	barcode = PadLeft(barcode, FatSecretBarcodeLength, "0")

	// invoke the api call
	body, err := c.InvokeAPI(
		"food.find_id_for_barcode",
		map[string]string{
			"barcode": barcode,
		},
	)
	if err != nil {
		return "", err
	}

	// parse the api response
	foodIDResp := FoodIDResponse{}
	if err := json.Unmarshal(body, &foodIDResp); err != nil {
		return "", err
	}

	// if an error response was returned
	if foodIDResp.Error != nil {
		// return the response error message
		return "", errors.New(foodIDResp.Error.Message)
	}

	// return the slice of food items
	return foodIDResp.ID.Value, nil
}
