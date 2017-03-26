package fatsecret

import (
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	consumerKey  = os.Getenv("FATSECRET_CONSUMER_KEY")
	sharedSecret = os.Getenv("FATSECRET_CONSUMER_SECRET")
)

func TestBrandTypeName(t *testing.T) {
	// define the test-cases
	testCases := []struct {
		brandType BrandType
		name      string
	}{
		{BrandTypeManufacturer, "manufacturer"},
		{BrandTypeRestaurant, "restaurant"},
		{BrandTypeSupermarket, "supermarket"},
		{42, "manufacturer"}, // test default is 'manufacturer'
	}

	// iterate through each test-case
	for _, tc := range testCases {
		// run the next sub-test
		t.Run(fmt.Sprintf("Brand Type for '%s' (%v)", tc.name, tc.brandType), func(t *testing.T) {
			name := brandTypeName(tc.brandType)
			if name != tc.name {
				t.Errorf("got '%s'; want '%s'", name, tc.name)
			}
		})
	}
}

func TestFoodBrandsByType(t *testing.T) {
	// define the test-cases
	testCases := []struct {
		brandType BrandType
		name      string
	}{
		{BrandTypeManufacturer, "manufacturer"},
		{BrandTypeRestaurant, "restaurant"},
		{BrandTypeSupermarket, "supermarket"},
		{42, "manufacturer"}, // test default is 'manufacturer'
	}

	// create a fatsecret client
	c, err := NewClient(consumerKey, sharedSecret)
	if err != nil {
		t.Errorf("Could not create client: '%v'", err)
	}

	// iterate through each test-case
	for _, tc := range testCases {
		// run the next sub-test
		t.Run(fmt.Sprintf("Brand Type for '%s' (%v)", tc.name, tc.brandType), func(t *testing.T) {
			// invoke the api call
			brands, err := c.FoodBrandsByType(tc.brandType)
			if err != nil {
				t.Errorf("Could not fetch brands: '%v'", err)
			}

			// verify api results
			if len(brands) <= 0 {
				t.Errorf("No brands returned for brand type '%s'", tc.name)
			}
		})
	}
}

func TestFoodBrandsStartingWith(t *testing.T) {
	// define the test-cases
	testCases := []struct {
		brandType BrandType
		name      string
	}{
		{BrandTypeManufacturer, "manufacturer"},
		{BrandTypeRestaurant, "restaurant"},
		{BrandTypeSupermarket, "supermarket"},
		{42, "manufacturer"}, // test default is 'manufacturer'
	}

	// define all of the starting characters
	startingChars := "*abcxyz"

	// iterate through each test-case
	for _, tc := range testCases {
		// iterate through each starting character
		for _, c := range startingChars {
			// create a fatsecret client
			client, err := NewClient(consumerKey, sharedSecret)
			if err != nil {
				t.Errorf("Could not create client: '%v'", err)
			}

			// wait to avoid flooding the api
			time.Sleep(time.Millisecond * 200)

			// run the next sub-test
			startsWith := string(c)
			t.Run(fmt.Sprintf("Brands of type '%s' starting with '%s'", tc.name, startsWith), func(t *testing.T) {

				// invoke the api call
				brands, err := client.FoodBrandsStartingWith(startsWith)
				if err != nil {
					t.Errorf("Could not fetch brands: '%v'", err)
				}

				// verify api results
				if len(brands) <= 0 {
					t.Errorf("No brands returned for brand type '%s'", tc.name)
				}
			})
		}
	}
}
