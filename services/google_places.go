package services

import (
	"bytes"
	"encoding/json"
	"event-explorer/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func FetchAutocompleteData(input string, sessionToken string) (models.UIAutocompleteResponse, error) {
	var output models.UIAutocompleteResponse

	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		return output, fmt.Errorf("GOOGLE_PLACES_API_KEY is empty")
	}

	url := "https://places.googleapis.com/v1/places:autocomplete"

	bodyPayload := models.AutocompleteInput{
		Input:                input,
		IncludedPrimaryTypes: []string{"(cities)"},
		SessionToken:         sessionToken,
	}

	jsonBody, err := json.Marshal(bodyPayload)
	if err != nil {
		return output, err
	}

	req, err := http.NewRequest(
		http.MethodPost,
		url,
		bytes.NewBuffer(jsonBody),
	)
	if err != nil {
		return output, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Goog-Api-Key", apiKey)
	req.Header.Set(
		"X-Goog-FieldMask",
		"suggestions.placePrediction.text,suggestions.placePrediction.placeId",
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, err
	}
	if resp.StatusCode != http.StatusOK {
		return output, fmt.Errorf("google api returned status %d", resp.StatusCode)
	}

	var googleResp models.GoogleAutocompleteResponse

	if err := json.Unmarshal(respBody, &googleResp); err != nil {
		return output, err
	}

	for _, s := range googleResp.Suggestions {
		output.Suggestions = append(
			output.Suggestions,
			models.UIAutocompleteSuggestion{
				Text:    s.PlacePrediction.Text.Text,
				PlaceID: s.PlacePrediction.PlaceID,
			},
		)
	}

	return output, nil
}

func FetchPlaceDetailsData(placeID, sessionToken string) (models.CleanLocationResult, error) {
	var output models.CleanLocationResult
	var googleDetails models.GooglePlaceDetailsResponse

	apiKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if apiKey == "" {
		return output, fmt.Errorf("GOOGLE_PLACES_API_KEY is empty")
	}

	sanitizedToken := strings.ReplaceAll(sessionToken, "-", "")
	url := "https://places.googleapis.com/v1/places/" + placeID

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return output, nil
	}

	q := req.URL.Query()
	q.Add("sessionToken", sanitizedToken)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("X-Goog-Api-Key", apiKey)
	req.Header.Set("X-Goog-FieldMask", "addressComponents")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return output, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return output, fmt.Errorf("api status %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, err
	}

	if err := json.Unmarshal(respBody, &googleDetails); err != nil {
		return output, err
	}

	for _, comp := range googleDetails.AddressComponents {
		for _, t := range comp.Types {
			if t == "locality" {
				output.City = comp.LongText
			}

			if t == "country" {
				output.CountryCode = comp.ShortText
			}
		}
	}

	return output, nil

}
