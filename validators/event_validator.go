package validators

import (
	"fmt"
	"net/url"
)

func EventsValidator(query url.Values) error {

	if city, exists := query["city"]; !exists {
		return fmt.Errorf("city parameter missing")
	} else if city[0] == "" {
		return fmt.Errorf("city value can't be empty")
	}

	if countryCode, exists := query["countryCode"]; !exists {
		return fmt.Errorf("countryCode parameter missing")
	} else if countryCode[0] == "" {
		return fmt.Errorf("countryCode value can't be empty")
	}
	return nil
}

func CachePurgeValidator(query url.Values) error {
	city, hasCity := query["city"]
	countryCode, hasCountryCode := query["countryCode"]
	category, hasCategory := query["category"]

	if !hasCity && !hasCountryCode && !hasCategory {
		return nil
	}

	if !hasCity || !hasCountryCode || !hasCategory {
		return fmt.Errorf(
			"city, countryCode and category must either all be provided or all be omitted",
		)
	}

	if len(city) == 0 || city[0] == "" {
		return fmt.Errorf("city value can't be empty")
	}

	if len(countryCode) == 0 || countryCode[0] == "" {
		return fmt.Errorf("countryCode value can't be empty")
	}

	if len(category) == 0 || category[0] == "" {
		return fmt.Errorf("category value can't be empty")
	}

	if category[0] != "Music" && category[0] != "Sports" {
		return fmt.Errorf(
			"invalid category: %s. Allowed values are Music and Sports",
			category[0],
		)
	}

	return nil
}