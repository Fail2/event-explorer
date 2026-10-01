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

type googlePlacesRoundTripFunc func(*http.Request) (*http.Response, error)

func (f googlePlacesRoundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func mockGooglePlacesHTTPTransport(
	statusCode int,
	body string,
	err error,
	check func(*http.Request),
) func() {
	oldTransport := http.DefaultTransport

	http.DefaultTransport = googlePlacesRoundTripFunc(func(req *http.Request) (*http.Response, error) {
		if check != nil {
			check(req)
		}

		if err != nil {
			return nil, err
		}

		return &http.Response{
			StatusCode: statusCode,
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})

	return func() {
		http.DefaultTransport = oldTransport
	}
}

func TestFetchAutocompleteData(t *testing.T) {
	t.Run("missing api key", func(t *testing.T) {
		os.Unsetenv("GOOGLE_PLACES_API_KEY")

		got, err := FetchAutocompleteData("Dhaka", "session-token")

		if err == nil {
			t.Fatal("expected error")
		}

		if got.Suggestions != nil {
			t.Fatalf("expected empty suggestions, got %+v", got.Suggestions)
		}
	})

	t.Run("http error", func(t *testing.T) {
		os.Setenv("GOOGLE_PLACES_API_KEY", "test-api-key")
		defer os.Unsetenv("GOOGLE_PLACES_API_KEY")

		restore := mockGooglePlacesHTTPTransport(
			0,
			"",
			errors.New("network error"),
			nil,
		)
		defer restore()

		_, err := FetchAutocompleteData("Dhaka", "session-token")

		if err == nil {
			t.Fatal("expected network error")
		}
	})

	t.Run("non 200 response", func(t *testing.T) {
		os.Setenv("GOOGLE_PLACES_API_KEY", "test-api-key")
		defer os.Unsetenv("GOOGLE_PLACES_API_KEY")

		restore := mockGooglePlacesHTTPTransport(
			http.StatusBadRequest,
			`{"error":"bad request"}`,
			nil,
			nil,
		)
		defer restore()

		_, err := FetchAutocompleteData("Dhaka", "session-token")

		if err == nil {
			t.Fatal("expected status error")
		}
	})

	t.Run("invalid json response", func(t *testing.T) {
		os.Setenv("GOOGLE_PLACES_API_KEY", "test-api-key")
		defer os.Unsetenv("GOOGLE_PLACES_API_KEY")

		restore := mockGooglePlacesHTTPTransport(
			http.StatusOK,
			`invalid-json`,
			nil,
			nil,
		)
		defer restore()

		_, err := FetchAutocompleteData("Dhaka", "session-token")

		if err == nil {
			t.Fatal("expected json error")
		}
	})

	t.Run("success", func(t *testing.T) {
		os.Setenv("GOOGLE_PLACES_API_KEY", "test-api-key")
		defer os.Unsetenv("GOOGLE_PLACES_API_KEY")

		body := `{
			"suggestions": [
				{
					"placePrediction": {
						"text": {
							"text": "Dhaka, Bangladesh"
						},
						"placeId": "place-123"
					}
				},
				{
					"placePrediction": {
						"text": {
							"text": "Chittagong, Bangladesh"
						},
						"placeId": "place-456"
					}
				}
			]
		}`

		restore := mockGooglePlacesHTTPTransport(
			http.StatusOK,
			body,
			nil,
			func(req *http.Request) {
				if req.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", req.Method)
				}

				if req.Header.Get("X-Goog-Api-Key") != "test-api-key" {
					t.Errorf("unexpected API key")
				}

				if req.Header.Get("Content-Type") != "application/json" {
					t.Errorf("unexpected content type")
				}

				expectedFieldMask :=
					"suggestions.placePrediction.text,suggestions.placePrediction.placeId"

				if req.Header.Get("X-Goog-FieldMask") != expectedFieldMask {
					t.Errorf("unexpected field mask")
				}
			},
		)
		defer restore()

		got, err := FetchAutocompleteData("Dhaka", "session-token")

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(got.Suggestions) != 2 {
			t.Fatalf("expected 2 suggestions, got %d", len(got.Suggestions))
		}

		if got.Suggestions[0].Text != "Dhaka, Bangladesh" {
			t.Errorf(
				"unexpected first text: %s",
				got.Suggestions[0].Text,
			)
		}

		if got.Suggestions[0].PlaceID != "place-123" {
			t.Errorf(
				"unexpected first place ID: %s",
				got.Suggestions[0].PlaceID,
			)
		}

		if got.Suggestions[1].Text != "Chittagong, Bangladesh" {
			t.Errorf(
				"unexpected second text: %s",
				got.Suggestions[1].Text,
			)
		}

		if got.Suggestions[1].PlaceID != "place-456" {
			t.Errorf(
				"unexpected second place ID: %s",
				got.Suggestions[1].PlaceID,
			)
		}
	})
}

func TestFetchPlaceDetailsData(t *testing.T) {
	t.Run("missing api key", func(t *testing.T) {
		os.Unsetenv("GOOGLE_PLACES_API_KEY")

		got, err := FetchPlaceDetailsData(
			"place-123",
			"session-token",
		)

		if err == nil {
			t.Fatal("expected error")
		}

		if got != (models.CleanLocationResult{}) {
			t.Fatalf(
				"expected empty result, got %+v",
				got,
			)
		}
	})

	t.Run("http error", func(t *testing.T) {
		os.Setenv("GOOGLE_PLACES_API_KEY", "test-api-key")
		defer os.Unsetenv("GOOGLE_PLACES_API_KEY")

		restore := mockGooglePlacesHTTPTransport(
			0,
			"",
			errors.New("network error"),
			nil,
		)
		defer restore()

		_, err := FetchPlaceDetailsData(
			"place-123",
			"session-token-123",
		)

		if err == nil {
			t.Fatal("expected network error")
		}
	})

	t.Run("non 200 response", func(t *testing.T) {
		os.Setenv("GOOGLE_PLACES_API_KEY", "test-api-key")
		defer os.Unsetenv("GOOGLE_PLACES_API_KEY")

		restore := mockGooglePlacesHTTPTransport(
			http.StatusNotFound,
			`{"error":"not found"}`,
			nil,
			nil,
		)
		defer restore()

		_, err := FetchPlaceDetailsData(
			"place-123",
			"session-token-123",
		)

		if err == nil {
			t.Fatal("expected status error")
		}
	})

	t.Run("read body error", func(t *testing.T) {
		os.Setenv("GOOGLE_PLACES_API_KEY", "test-api-key")
		defer os.Unsetenv("GOOGLE_PLACES_API_KEY")

		oldTransport := http.DefaultTransport

		http.DefaultTransport = googlePlacesRoundTripFunc(
			func(req *http.Request) (*http.Response, error) {
				return &http.Response{
					StatusCode: http.StatusOK,
					Body:       errorReader{},
					Header:     make(http.Header),
				}, nil
			},
		)

		defer func() {
			http.DefaultTransport = oldTransport
		}()

		_, err := FetchPlaceDetailsData(
			"place-123",
			"session-token-123",
		)

		if err == nil {
			t.Fatal("expected read error")
		}
	})

	t.Run("invalid json response", func(t *testing.T) {
		os.Setenv("GOOGLE_PLACES_API_KEY", "test-api-key")
		defer os.Unsetenv("GOOGLE_PLACES_API_KEY")

		restore := mockGooglePlacesHTTPTransport(
			http.StatusOK,
			`invalid-json`,
			nil,
			nil,
		)
		defer restore()

		_, err := FetchPlaceDetailsData(
			"place-123",
			"session-token-123",
		)

		if err == nil {
			t.Fatal("expected json error")
		}
	})

	t.Run("success", func(t *testing.T) {
		os.Setenv("GOOGLE_PLACES_API_KEY", "test-api-key")
		defer os.Unsetenv("GOOGLE_PLACES_API_KEY")

		body := `{
			"addressComponents": [
				{
					"longText": "Dhaka",
					"shortText": "DHK",
					"types": ["locality"]
				},
				{
					"longText": "Bangladesh",
					"shortText": "BD",
					"types": ["country"]
				}
			]
		}`

		restore := mockGooglePlacesHTTPTransport(
			http.StatusOK,
			body,
			nil,
			func(req *http.Request) {
				if req.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", req.Method)
				}

				if req.URL.Query().Get("sessionToken") != "sessiontoken123" {
					t.Errorf(
						"unexpected session token: %s",
						req.URL.Query().Get("sessionToken"),
					)
				}

				if req.Header.Get("X-Goog-Api-Key") != "test-api-key" {
					t.Errorf("unexpected API key")
				}

				if req.Header.Get("X-Goog-FieldMask") != "addressComponents" {
					t.Errorf("unexpected field mask")
				}
			},
		)
		defer restore()

		got, err := FetchPlaceDetailsData(
			"place-123",
			"session-token-123",
		)

		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if got.City != "Dhaka" {
			t.Errorf("expected Dhaka, got %s", got.City)
		}

		if got.CountryCode != "BD" {
			t.Errorf(
				"expected BD, got %s",
				got.CountryCode,
			)
		}
	})
}
