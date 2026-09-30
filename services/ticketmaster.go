package services

import (
	"encoding/json"
	"event-explorer/models"
	"fmt"
	"io"
	"net/http"
	"os"
)

func FetchTicketmasterEvents(
	city,
	countryCode,
	classification,
	size string,
) (models.TicketmasterResponse, error) {

	var output models.TicketmasterResponse

	apiKey := os.Getenv("TICKETMASTER_API_KEY")
	if apiKey == "" {
		return output, fmt.Errorf("TICKETMASTER_API_KEY is not configured")
	}

	baseURL := "https://app.ticketmaster.com/discovery/v2/events.json"

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return output, err
	}

	q := req.URL.Query()

	q.Add("city", city)
	q.Add("countryCode", countryCode)
	q.Add("classificationName", classification)
	q.Add("size", size)
	q.Add("apikey", apiKey)

	req.URL.RawQuery = q.Encode()

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, err
	}

	if resp.StatusCode != http.StatusOK {
		return output, fmt.Errorf(
			"ticketmaster api returned status %d: %s",
			resp.StatusCode,
			string(body),
		)
	}

	if err := json.Unmarshal(body, &output); err != nil {
		return output, err
	}

	return output, nil
}
