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
