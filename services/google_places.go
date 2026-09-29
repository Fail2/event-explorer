package services

import (
	"bytes"
	"encoding/json"
	"event-explorer/models"
	"fmt"
	"io"
	"net/http"
	"os"
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

	fmt.Println("---------------------------------")
	fmt.Println("Request Method:", req.Method)
	fmt.Println("Request URL:", req.URL.String())
	fmt.Println("Request Body:", string(jsonBody))
	fmt.Println("---------------------------------")

	resp, err := client.Do(req)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, err
	}

	fmt.Println("---------------------------------")
	fmt.Println("Google API Status:", resp.StatusCode)
	fmt.Println("Google Response:", string(respBody))
	fmt.Println("---------------------------------")

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

	fmt.Println(output)
	return output, nil
}
