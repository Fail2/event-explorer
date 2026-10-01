package services

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"event-explorer/models"
)

type eventServiceRoundTripFunc func(*http.Request) (*http.Response, error)

func (f eventServiceRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newEventServiceResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestMapTMToUIEvents(t *testing.T) {
	t.Run("nil events", func(t *testing.T) {
		resp := models.TicketmasterResponse{}

		got := mapTMToUIEvents(resp)

		if len(got) != 0 {
			t.Errorf("expected 0 events, got %d", len(got))
		}
	})

	t.Run("event without image and venue", func(t *testing.T) {
		body := `{
			"_embedded": {
				"events": [
					{
						"id": "event-1",
						"name": "Test Event",
						"url": "https://example.com/event",
						"pleaseNote": "Test description",
						"dates": {
							"start": {
								"localDate": "2025-01-01",
								"localTime": "18:00:00"
							},
							"timeZone": "Asia/Dhaka"
						}
					}
				]
			}
		}`

		var resp models.TicketmasterResponse

		if err := json.Unmarshal([]byte(body), &resp); err != nil {
			t.Fatalf("failed to unmarshal test JSON: %v", err)
		}

		got := mapTMToUIEvents(resp)

		if len(got) != 1 {
			t.Fatalf("expected 1 event, got %d", len(got))
		}

		if got[0].ID != "event-1" {
			t.Errorf("ID = %q, want %q", got[0].ID, "event-1")
		}

		if got[0].Name != "Test Event" {
			t.Errorf("Name = %q, want %q", got[0].Name, "Test Event")
		}

		if got[0].TicketURL != "https://example.com/event" {
			t.Errorf("TicketURL = %q, want %q", got[0].TicketURL, "https://example.com/event")
		}

		if got[0].Description != "Test description" {
			t.Errorf("Description = %q, want %q", got[0].Description, "Test description")
		}

		if got[0].Date != "2025-01-01" {
			t.Errorf("Date = %q, want %q", got[0].Date, "2025-01-01")
		}

		if got[0].Time != "18:00:00" {
			t.Errorf("Time = %q, want %q", got[0].Time, "18:00:00")
		}

		if got[0].TimeZone != "Asia/Dhaka" {
			t.Errorf("TimeZone = %q, want %q", got[0].TimeZone, "Asia/Dhaka")
		}

		if got[0].Image != "" {
			t.Errorf("Image = %q, want empty", got[0].Image)
		}

		if got[0].Venue != "" {
			t.Errorf("Venue = %q, want empty", got[0].Venue)
		}

		if got[0].City != "" {
			t.Errorf("City = %q, want empty", got[0].City)
		}

		if got[0].State != "" {
			t.Errorf("State = %q, want empty", got[0].State)
		}

		if got[0].Country != "" {
			t.Errorf("Country = %q, want empty", got[0].Country)
		}

		if got[0].Address != "" {
			t.Errorf("Address = %q, want empty", got[0].Address)
		}
	})

	t.Run("event with image and venue", func(t *testing.T) {
		body := `{
			"_embedded": {
				"events": [
					{
						"id": "event-2",
						"name": "Music Event",
						"url": "https://example.com/music",
						"pleaseNote": "Amazing music event",
						"dates": {
							"start": {
								"localDate": "2025-02-01",
								"localTime": "20:00:00"
							},
							"timeZone": "Asia/Dhaka"
						},
						"images": [
							{
								"url": "https://example.com/image.jpg"
							}
						],
						"_embedded": {
							"venues": [
								{
									"name": "Test Venue",
									"city": {
										"name": "Dhaka"
									},
									"state": {
										"name": "Dhaka"
									},
									"country": {
										"name": "Bangladesh",
										"countryCode": "BD"
									},
									"address": {
										"line1": "123 Test Street"
									}
								}
							]
						}
					}
				]
			}
		}`

		var resp models.TicketmasterResponse

		if err := json.Unmarshal([]byte(body), &resp); err != nil {
			t.Fatalf("failed to unmarshal test JSON: %v", err)
		}

		got := mapTMToUIEvents(resp)

		if len(got) != 1 {
			t.Fatalf("expected 1 event, got %d", len(got))
		}

		if got[0].ID != "event-2" {
			t.Errorf("ID = %q, want %q", got[0].ID, "event-2")
		}

		if got[0].Name != "Music Event" {
			t.Errorf("Name = %q, want %q", got[0].Name, "Music Event")
		}

		if got[0].TicketURL != "https://example.com/music" {
			t.Errorf("TicketURL = %q, want %q", got[0].TicketURL, "https://example.com/music")
		}

		if got[0].Description != "Amazing music event" {
			t.Errorf("Description = %q, want %q", got[0].Description, "Amazing music event")
		}

		if got[0].Date != "2025-02-01" {
			t.Errorf("Date = %q, want %q", got[0].Date, "2025-02-01")
		}

		if got[0].Time != "20:00:00" {
			t.Errorf("Time = %q, want %q", got[0].Time, "20:00:00")
		}

		if got[0].TimeZone != "Asia/Dhaka" {
			t.Errorf("TimeZone = %q, want %q", got[0].TimeZone, "Asia/Dhaka")
		}

		if got[0].Image != "https://example.com/image.jpg" {
			t.Errorf("Image = %q, want %q", got[0].Image, "https://example.com/image.jpg")
		}

		if got[0].Venue != "Test Venue" {
			t.Errorf("Venue = %q, want %q", got[0].Venue, "Test Venue")
		}

		if got[0].City != "Dhaka" {
			t.Errorf("City = %q, want %q", got[0].City, "Dhaka")
		}

		if got[0].State != "Dhaka" {
			t.Errorf("State = %q, want %q", got[0].State, "Dhaka")
		}

		if got[0].Country != "Bangladesh" {
			t.Errorf("Country = %q, want %q", got[0].Country, "Bangladesh")
		}

		if got[0].Address != "123 Test Street" {
			t.Errorf("Address = %q, want %q", got[0].Address, "123 Test Street")
		}
	})

	t.Run("multiple events", func(t *testing.T) {
		body := `{
			"_embedded": {
				"events": [
					{
						"id": "event-1",
						"name": "Event One",
						"url": "https://example.com/one",
						"pleaseNote": "First event",
						"dates": {
							"start": {
								"localDate": "2025-01-01",
								"localTime": "10:00:00"
							},
							"timeZone": "Asia/Dhaka"
						}
					},
					{
						"id": "event-2",
						"name": "Event Two",
						"url": "https://example.com/two",
						"pleaseNote": "Second event",
						"dates": {
							"start": {
								"localDate": "2025-01-02",
								"localTime": "12:00:00"
							},
							"timeZone": "Asia/Dhaka"
						}
					}
				]
			}
		}`

		var resp models.TicketmasterResponse

		if err := json.Unmarshal([]byte(body), &resp); err != nil {
			t.Fatalf("failed to unmarshal test JSON: %v", err)
		}

		got := mapTMToUIEvents(resp)

		if len(got) != 2 {
			t.Fatalf("expected 2 events, got %d", len(got))
		}

		if got[0].ID != "event-1" {
			t.Errorf("first event ID = %q, want %q", got[0].ID, "event-1")
		}

		if got[1].ID != "event-2" {
			t.Errorf("second event ID = %q, want %q", got[1].ID, "event-2")
		}
	})
}

func TestGetConcurrentEvents(t *testing.T) {
	oldAPIKey := os.Getenv("TICKETMASTER_API_KEY")
	if err := os.Setenv("TICKETMASTER_API_KEY", "test-api-key"); err != nil {
		t.Fatalf("failed to set API key: %v", err)
	}

	t.Cleanup(func() {
		if oldAPIKey == "" {
			_ = os.Unsetenv("TICKETMASTER_API_KEY")
		} else {
			_ = os.Setenv("TICKETMASTER_API_KEY", oldAPIKey)
		}
	})

	originalTransport := httpClient.Transport

	t.Cleanup(func() {
		httpClient.Transport = originalTransport
	})

	t.Run("both success", func(t *testing.T) {
		body := `{
			"_embedded": {
				"events": [
					{
						"id": "event-1",
						"name": "Test Event",
						"url": "https://example.com/event",
						"pleaseNote": "Test description",
						"dates": {
							"start": {
								"localDate": "2025-01-01",
								"localTime": "18:00:00"
							},
							"timeZone": "Asia/Dhaka"
						}
					}
				]
			}
		}`

		httpClient.Transport = eventServiceRoundTripFunc(
			func(req *http.Request) (*http.Response, error) {
				return newEventServiceResponse(http.StatusOK, body), nil
			},
		)

		results, errorsMap := GetConcurrentEvents("Dhaka", "BD")

		if len(errorsMap) != 0 {
			t.Fatalf("expected no errors, got %v", errorsMap)
		}

		if len(results) != 2 {
			t.Fatalf("expected 2 result categories, got %d", len(results))
		}

		if _, ok := results["music"]; !ok {
			t.Error("expected music results")
		}

		if _, ok := results["sports"]; !ok {
			t.Error("expected sports results")
		}
	})

	t.Run("both errors", func(t *testing.T) {
		httpClient.Transport = eventServiceRoundTripFunc(
			func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("ticketmaster unavailable")
			},
		)

		results, errorsMap := GetConcurrentEvents("Dhaka", "BD")

		if len(results) != 0 {
			t.Errorf("expected no results, got %v", results)
		}

		if len(errorsMap) != 2 {
			t.Fatalf("expected 2 errors, got %d", len(errorsMap))
		}

		if errorsMap["music"] == nil {
			t.Error("expected music error")
		}

		if errorsMap["sports"] == nil {
			t.Error("expected sports error")
		}
	})

	t.Run("music error sports success", func(t *testing.T) {
		body := `{
			"_embedded": {
				"events": [
					{
						"id": "sports-1",
						"name": "Sports Event",
						"url": "https://example.com/sports",
						"pleaseNote": "Sports description",
						"dates": {
							"start": {
								"localDate": "2025-03-01",
								"localTime": "15:00:00"
							},
							"timeZone": "Asia/Dhaka"
						}
					}
				]
			}
		}`

		httpClient.Transport = eventServiceRoundTripFunc(
			func(req *http.Request) (*http.Response, error) {
				if req.URL.Query().Get("classificationName") == "Music" {
					return nil, errors.New("music request failed")
				}

				return newEventServiceResponse(http.StatusOK, body), nil
			},
		)

		results, errorsMap := GetConcurrentEvents("Dhaka", "BD")

		if len(results) != 1 {
			t.Fatalf("expected 1 result category, got %d", len(results))
		}

		if len(errorsMap) != 1 {
			t.Fatalf("expected 1 error, got %d", len(errorsMap))
		}

		if errorsMap["music"] == nil {
			t.Error("expected music error")
		}

		if _, ok := results["sports"]; !ok {
			t.Error("expected sports results")
		}
	})

	t.Run("music success sports error", func(t *testing.T) {
		body := `{
			"_embedded": {
				"events": [
					{
						"id": "music-1",
						"name": "Music Event",
						"url": "https://example.com/music",
						"pleaseNote": "Music description",
						"dates": {
							"start": {
								"localDate": "2025-04-01",
								"localTime": "19:00:00"
							},
							"timeZone": "Asia/Dhaka"
						}
					}
				]
			}
		}`

		httpClient.Transport = eventServiceRoundTripFunc(
			func(req *http.Request) (*http.Response, error) {
				if req.URL.Query().Get("classificationName") == "Sports" {
					return nil, errors.New("sports request failed")
				}

				return newEventServiceResponse(http.StatusOK, body), nil
			},
		)

		results, errorsMap := GetConcurrentEvents("Dhaka", "BD")

		if len(results) != 1 {
			t.Fatalf("expected 1 result category, got %d", len(results))
		}

		if len(errorsMap) != 1 {
			t.Fatalf("expected 1 error, got %d", len(errorsMap))
		}

		if _, ok := results["music"]; !ok {
			t.Error("expected music results")
		}

		if errorsMap["sports"] == nil {
			t.Error("expected sports error")
		}
	})
}
