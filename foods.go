package fatsecret

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	fatSecretBarcodeLength = 13
)

type FoodSearchItem struct {
	ID          string `json:"food_id,omitempty"`
	Name        string `json:"food_name,omitempty"`
	Type        string `json:"food_type,omitempty"`
	BrandName   string `json:"brand_name,omitempty"`
	URL         string `json:"food_url,omitempty"`
	Description string `json:"food_description,omitempty"`
}

type FoodSearchResponseFoods struct {
	PageNumber   int              `json:"page_number,string"`
	PageSize     int              `json:"max_results,string"`
	TotalResults int              `json:"total_results,string"`
	Food         []FoodSearchItem `json:"food"`
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

type FoodInfo struct {
	ID        string       `json:"food_id"`
	Name      string       `json:"food_name"`
	Type      string       `json:"food_type"`
	URL       string       `json:"food_url"`
	BrandName string       `json:"brand_name"`
	Servings  FoodServings `json:"servings"`
}

type FoodServings struct {
	Serving FoodServing `json:"serving"`
}

type FoodServing struct {
	// serving info
	ServingID              string `json:"serving_id"`              // serving_id – the unique serving identifier.
	ServingDescription     string `json:"serving_description"`     // serving_description – the full description of the serving size. E.G.: "1 cup" or "100 g".
	ServingURL             string `json:"serving_url"`             // serving_url – URL of the serving size for this food item on www.fatsecret.com.
	MetricServingAmount    string `json:"metric_serving_amount"`   // metric_serving_amount is a Decimal - the metric quantity combined with metric_serving_unit to derive the total standardized quantity of the serving (where available).
	MetricServingUnit      string `json:"metric_serving_unit"`     // metric_serving_unit – the metric unit of measure for the serving size – either "g" or "ml" or "oz" – combined with metric_serving_amount to derive the total standardized quantity of the serving (where available).
	NumberOfUnits          string `json:"number_of_units"`         // number_of_units is a Decimal - the number of units in this standard serving size. For instance, if the serving description is "2 tablespoons" the number of units is "2", while if the serving size is "1 cup" the number of units is "1".
	MeasurementDescription string `json:"measurement_description"` // measurement_description – a description of the unit of measure used in the serving description. For instance, if the description is "1/2 cup" the measurement description is "cup", while if the serving size is "100 g" the measurement description is "g".

	// nutrient info
	Calories           string `json:"calories"`            // calories is a Decimal – the energy content in kcal.
	Carbohydrate       string `json:"carbohydrate"`        // carbohydrate is a Decimal – the total carbohydrate content in grams.
	Protein            string `json:"protein"`             // protein is a Decimal – the protein content in grams.
	Fat                string `json:"fat"`                 // fat is a Decimal – the total fat content in grams.
	SaturatedFat       string `json:"saturated_fat"`       // saturated_fat is a Decimal – the saturated fat content in grams (where available).
	PolyunsaturatedFat string `json:"polyunsaturated_fat"` // polyunsaturated_fat is a Decimal – the polyunsaturated fat content in grams (where available).
	MonounsaturatedFat string `json:"monounsaturated_fat"` // monounsaturated_fat is a Decimal – the monounsaturated fat content in grams (where available).
	TransFat           string `json:"trans_fat"`           // trans_fat is a Decimal – the trans fat content in grams (where available).
	Cholesterol        string `json:"cholesterol"`         // cholesterol is a Decimal – the cholesterol content in milligrams (where available).
	Sodium             string `json:"sodium"`              // sodium is a Decimal – the sodium content in milligrams (where available).
	Potassium          string `json:"potassium"`           // potassium is a Decimal – the potassium content in milligrams (where available).
	Fiber              string `json:"fiber"`               // fiber is a Decimal – the fiber content in grams (where available).
	Sugar              string `json:"sugar"`               // sugar is a Decimal – the sugar content in grams (where available).
	VitaminA           string `json:"vitamin_a"`           // vitamin_a is a Decimal – the percentage of daily recommended Vitamin A, based on a 2000 calorie diet (where available).
	VitaminC           string `json:"vitamin_c"`           // vitamin_c is a Decimal – the percentage of daily recommended Vitamin C, based on a 2000 calorie diet (where available).
	Calcium            string `json:"calcium"`             // calcium is a Decimal – the percentage of daily recommended Calcium, based on a 2000 calorie diet (where available).
	Iron               string `json:"iron"`                // iron is a Decimal – the percentage of daily recommended Iron, based on a 2000 calorie diet (where available).
}

type FoodInfoResponse struct {
	Food  *FoodInfo      `json:"food,omitempty"`
	Error *ErrorResponse `json:"error,omitempty"`
}

// FoodSearch invokes the FatSecret 'foods.search' API call and
// returns the response as a slice of FoodSearchItem structs
func (c *Client) FoodSearch(query string) ([]FoodSearchItem, error) {
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
	if barcodeLen == 0 || barcodeLen > fatSecretBarcodeLength {
		return "", fmt.Errorf("Invalid barcode length '%d' given", barcodeLen)
	}

	// pad the barcode to GTIN-13 format
	barcode = padLeft(barcode, fatSecretBarcodeLength, "0")

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

// FoodByID invokes the FatSecret 'food.get' API call for the
// given food-id and returns the response
func (c *Client) FoodByID(id string) (*FoodInfo, error) {
	// if the food id is invalid
	if len(id) == 0 {
		return nil, fmt.Errorf("Invalid food id '%s' given", id)
	}

	// invoke the api call
	body, err := c.InvokeAPI(
		"food.get",
		map[string]string{
			"food_id": id,
		},
	)
	if err != nil {
		return nil, err
	}

	// parse the api response
	resp := FoodInfoResponse{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}

	// if an error response was returned
	if resp.Error != nil {
		// return the response error message
		return nil, errors.New(resp.Error.Message)
	}

	// return the food info
	return resp.Food, nil
}
