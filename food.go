package fatsecret

import (
	"encoding/json"
	"errors"
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
