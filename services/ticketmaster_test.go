package services

import (
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"event-explorer/models"
)

type mockRoundTripper struct {
	response *http.Response
	err      error
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}

	return m.response, nil
}

func mockHTTPResponse(statusCode int, body string) *http.Response {
	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     make(http.Header),
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "valid time",
			input: "19:30:00",
			want:  "7:30 PM",
		},
		{
			name:  "invalid time",
			input: "invalid",
			want:  "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatTime(tt.input)

			if got != tt.want {
				t.Errorf("formatTime() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFetchTicketmasterEvents(t *testing.T) {
	oldAPIKey := getenv("TICKETMASTER_API_KEY")
	defer setenv("TICKETMASTER_API_KEY", oldAPIKey)

	t.Run("missing api key", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "")

		_, err := FetchTicketmasterEvents(
			"London",
			"GB",
			"Music",
			"6",
		)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("request creation error", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				err: errors.New("request failed"),
			},
		}

		_, err := FetchTicketmasterEvents(
			"London",
			"GB",
			"Music",
			"6",
		)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("http client error", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				err: errors.New("network error"),
			},
		}

		_, err := FetchTicketmasterEvents(
			"London",
			"GB",
			"Music",
			"6",
		)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("read body error", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(errorReader{}),
					Header:     make(http.Header),
				},
			},
		}

		_, err := FetchTicketmasterEvents(
			"London",
			"GB",
			"Music",
			"6",
		)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("non 200 response", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				response: mockHTTPResponse(
					http.StatusBadRequest,
					`{"error":"bad request"}`,
				),
			},
		}

		_, err := FetchTicketmasterEvents(
			"London",
			"GB",
			"Music",
			"6",
		)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				response: mockHTTPResponse(
					http.StatusOK,
					`invalid json`,
				),
			},
		}

		_, err := FetchTicketmasterEvents(
			"London",
			"GB",
			"Music",
			"6",
		)

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("success", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		body := `{
			"_embedded": {
				"events": []
			}
		}`

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				response: mockHTTPResponse(
					http.StatusOK,
					body,
				),
			},
		}

		result, err := FetchTicketmasterEvents(
			"London",
			"GB",
			"Music",
			"6",
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Embedded.Events == nil {
			t.Fatal("expected embedded events")
		}
	})
}

func TestFetchTicketmasterEventDetails(t *testing.T) {
	oldAPIKey := getenv("TICKETMASTER_API_KEY")
	defer setenv("TICKETMASTER_API_KEY", oldAPIKey)

	t.Run("missing api key", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "")

		_, err := FetchTicketmasterEventDetails("event-123")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("http client error", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				err: errors.New("network error"),
			},
		}

		_, err := FetchTicketmasterEventDetails("event-123")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("non 200 response", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				response: mockHTTPResponse(
					http.StatusNotFound,
					`{"error":"not found"}`,
				),
			},
		}

		_, err := FetchTicketmasterEventDetails("event-123")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("read body error", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				response: &http.Response{
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(errorReader{}),
					Header:     make(http.Header),
				},
			},
		}

		_, err := FetchTicketmasterEventDetails("event-123")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				response: mockHTTPResponse(
					http.StatusOK,
					`invalid json`,
				),
			},
		}

		_, err := FetchTicketmasterEventDetails("event-123")

		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	t.Run("success with image and venue", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		body := `{
			"id": "event-123",
			"name": "Test Concert",
			"url": "https://example.com/ticket",
			"description": "Test description",
			"dates": {
				"start": {
					"localDate": "2026-10-01",
					"localTime": "19:30:00",
					"timezone": "Europe/London"
				}
			},
			"classifications": [
				{
					"category": {
						"name": "Music"
					},
					"genre": {
						"name": "Rock"
					}
				}
			],
			"images": [
				{
					"url": "https://example.com/image.jpg"
				}
			],
			"_embedded": {
				"venues": [
					{
						"name": "Test Arena",
						"city": {
							"name": "London"
						},
						"state": {
							"name": "London"
						},
						"country": {
							"name": "United Kingdom",
							"countryCode": "GB"
						},
						"address": {
							"line1": "123 Test Street"
						},
						"generalInfo": {
							"generalRule": "No outside food",
							"childRule": "Children allowed"
						}
					}
				]
			}
		}`

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				response: mockHTTPResponse(
					http.StatusOK,
					body,
				),
			},
		}

		result, err := FetchTicketmasterEventDetails("event-123")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertUIEvent(t, result)
	})

	t.Run("success without image and venue", func(t *testing.T) {
		setenv("TICKETMASTER_API_KEY", "test-key")

		oldClient := httpClient
		defer func() {
			httpClient = oldClient
		}()

		body := `{
			"id": "event-456",
			"name": "Simple Event",
			"url": "https://example.com/ticket",
			"description": "Simple description",
			"dates": {
				"start": {
					"localDate": "2026-10-02",
					"localTime": "20:00:00",
					"timezone": "Europe/London"
				}
			},
			"classifications": [
				{
					"category": {
						"name": "Sports"
					},
					"genre": {
						"name": "Football"
					}
				}
			]
		}`

		httpClient = &http.Client{
			Transport: &mockRoundTripper{
				response: mockHTTPResponse(
					http.StatusOK,
					body,
				),
			},
		}

		result, err := FetchTicketmasterEventDetails("event-456")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID != "event-456" {
			t.Errorf("ID = %v, want event-456", result.ID)
		}

		if result.Image != "" {
			t.Errorf("Image = %v, want empty", result.Image)
		}

		if result.Venue != "" {
			t.Errorf("Venue = %v, want empty", result.Venue)
		}
	})
}

func assertUIEvent(t *testing.T, event models.UIEvent) {
	t.Helper()

	if event.ID != "event-123" {
		t.Errorf("ID = %v, want event-123", event.ID)
	}

	if event.Name != "Test Concert" {
		t.Errorf("Name = %v, want Test Concert", event.Name)
	}

	if event.TicketURL != "https://example.com/ticket" {
		t.Errorf("TicketURL = %v, want ticket URL", event.TicketURL)
	}

	if event.Description != "Test description" {
		t.Errorf("Description = %v, want Test description", event.Description)
	}

	if event.Date != "2026-10-01" {
		t.Errorf("Date = %v, want 2026-10-01", event.Date)
	}

	if event.Time != "7:30 PM" {
		t.Errorf("Time = %v, want 7:30 PM", event.Time)
	}

	if event.TimeZone != "Europe/London" {
		t.Errorf("TimeZone = %v, want Europe/London", event.TimeZone)
	}

	if event.Category != "Music" {
		t.Errorf("Category = %v, want Music", event.Category)
	}

	if event.Genre != "Rock" {
		t.Errorf("Genre = %v, want Rock", event.Genre)
	}

	if event.Image != "https://example.com/image.jpg" {
		t.Errorf("Image = %v, want image URL", event.Image)
	}

	if event.Venue != "Test Arena" {
		t.Errorf("Venue = %v, want Test Arena", event.Venue)
	}

	if event.City != "London" {
		t.Errorf("City = %v, want London", event.City)
	}

	if event.State != "London" {
		t.Errorf("State = %v, want London", event.State)
	}

	if event.Country != "United Kingdom" {
		t.Errorf("Country = %v, want United Kingdom", event.Country)
	}

	if event.CountryCode != "GB" {
		t.Errorf("CountryCode = %v, want GB", event.CountryCode)
	}

	if event.Address != "123 Test Street" {
		t.Errorf("Address = %v, want 123 Test Street", event.Address)
	}

	if event.GeneralRule != "No outside food" {
		t.Errorf("GeneralRule = %v, want No outside food", event.GeneralRule)
	}

	if event.ChildRule != "Children allowed" {
		t.Errorf("ChildRule = %v, want Children allowed", event.ChildRule)
	}
}

type errorReader struct{}

func (errorReader) Read([]byte) (int, error) {
	return 0, errors.New("read error")
}

func (errorReader) Close() error {
	return nil
}

func getenv(key string) string {
	return strings.TrimSpace(getenvRaw(key))
}

func getenvRaw(key string) string {
	return os.Getenv(key)
}

func setenv(key, value string) {
	_ = os.Setenv(key, value)
}
