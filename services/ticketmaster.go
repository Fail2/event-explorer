package services

import (
	"encoding/json"
	"event-explorer/models"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

var httpClient = &http.Client{}

func formatTime(timeStr string) string {
	t, err := time.Parse("15:04:05", timeStr)
	if err != nil {
		return timeStr
	}

	return t.Format("3:04 PM")
}

func FetchTicketmasterEvents(city, countryCode, classification, size string) (models.TicketmasterResponse, error) {

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

	resp, err := httpClient.Do(req)
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

func FetchTicketmasterEventDetails(eventID string) (models.UIEvent, error) {

	var output models.UIEvent

	apiKey := os.Getenv("TICKETMASTER_API_KEY")
	if apiKey == "" {
		return output, fmt.Errorf("TICKETMASTER_API_KEY is not configured")
	}

	baseURL := fmt.Sprintf(
		"https://app.ticketmaster.com/discovery/v2/events/%s.json",
		eventID,
	)

	req, err := http.NewRequest(http.MethodGet, baseURL, nil)
	if err != nil {
		return output, err
	}

	q := req.URL.Query()
	q.Add("apikey", apiKey)
	req.URL.RawQuery = q.Encode()

	resp, err := httpClient.Do(req)
	if err != nil {
		return output, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return output, fmt.Errorf(
			"ticketmaster api details returned status %d",
			resp.StatusCode,
		)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return output, err
	}

	var tmEvent models.TMEvent

	if err := json.Unmarshal(body, &tmEvent); err != nil {
		return output, err
	}

	output = models.UIEvent{
		ID:          tmEvent.ID,
		Name:        tmEvent.Name,
		TicketURL:   tmEvent.URL,
		Description: tmEvent.Description,
		Date:        tmEvent.Dates.Start.LocalDate,
		Time:        formatTime(tmEvent.Dates.Start.LocalTime),
		TimeZone:    tmEvent.Dates.TimeZone,
		Category:    tmEvent.Classifications[0].Category.Name,
		Genre:       tmEvent.Classifications[0].Genre.Name,
	}

	if len(tmEvent.Images) > 0 {
		output.Image = tmEvent.Images[0].URL
	}

	if len(tmEvent.Embedded.Venues) > 0 {
		v := tmEvent.Embedded.Venues[0]

		output.Venue = v.Name
		output.City = v.City.Name
		output.State = v.State.Name
		output.Country = v.Country.Name
		output.CountryCode = v.Country.CountryCode
		output.Address = v.Address.Adress
		output.GeneralRule = v.GeneralInfo.GeneralRule
		output.ChildRule = v.GeneralInfo.ChildRule
	}

	return output, nil
}