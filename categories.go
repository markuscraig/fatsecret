package fatsecret

import (
	"encoding/json"
	"errors"
)

type FoodCategory struct {
	ID          string `json:"food_category_id"`
	Name        string `json:"food_category_name"`
	Description string `json:"food_category_description"`
}

type FoodCategories struct {
	Categories []FoodCategory `json:"food_category"`
}

type FoodCategoriesResponse struct {
	Categories *FoodCategories `json:"food_categories,omitempty"`
	Error      *ErrorResponse  `json:"error,omitempty"`
}

// FoodCategories invokes the FatSecret 'food_categories.get' API call
// and returns the response as a slice of FoodCategory structs
func (c *Client) FoodCategories() ([]FoodCategory, error) {
	// invoke the api call
	body, err := c.InvokeAPI(
		"food_categories.get",
		map[string]string{},
	)
	if err != nil {
		return nil, err
	}

	// parse the api response
	resp := FoodCategoriesResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// if an error response was returned
	if resp.Error != nil {
		// return the response error message
		return nil, errors.New(resp.Error.Message)
	}

	// return the slice of food category entries
	return resp.Categories.Categories, nil
}
