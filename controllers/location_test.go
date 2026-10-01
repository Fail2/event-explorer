package controllers

import (
	"errors"
	"net/http/httptest"
	"testing"

	"event-explorer/models"

	beegoContext "github.com/beego/beego/v2/server/web/context"
)

func newTestController(method, target string) *LocationController {
	req := httptest.NewRequest(method, target, nil)
	recorder := httptest.NewRecorder()

	controller := &LocationController{}
	controller.Ctx = beegoContext.NewContext()
	controller.Ctx.Reset(recorder, req)
	controller.Data = make(map[interface{}]interface{})

	return controller
}

func TestLocationController_Autocomplete(t *testing.T) {
	original := fetchAutocompleteData
	defer func() {
		fetchAutocompleteData = original
	}()

	t.Run("success", func(t *testing.T) {
		fetchAutocompleteData = func(input, sessionToken string) (models.UIAutocompleteResponse, error) {
			return models.UIAutocompleteResponse{}, nil
		}

		controller := newTestController(
			"GET",
			"/?input=London&sessionToken=test-token",
		)

		controller.Autocomplete()

		if controller.Data["json"] == nil {
			t.Fatal("expected json data")
		}
	})

	t.Run("error", func(t *testing.T) {
		fetchAutocompleteData = func(input, sessionToken string) (models.UIAutocompleteResponse, error) {
			return models.UIAutocompleteResponse{}, errors.New("service error")
		}

		controller := newTestController(
			"GET",
			"/?input=London&sessionToken=test-token",
		)

		controller.Autocomplete()

		result, ok := controller.Data["json"].(map[string]string)
		if !ok {
			t.Fatal("expected error response")
		}

		if result["error"] != "Failed to fetch autocomplete predictions" {
			t.Errorf(
				"error = %q, want %q",
				result["error"],
				"Failed to fetch autocomplete predictions",
			)
		}
	})
}

func TestLocationController_GetPlaceDetails(t *testing.T) {
	original := fetchPlaceDetailsData
	defer func() {
		fetchPlaceDetailsData = original
	}()

	t.Run("success", func(t *testing.T) {
		fetchPlaceDetailsData = func(placeID, sessionToken string) (models.CleanLocationResult, error) {
			return models.CleanLocationResult{
				City:        "London",
				CountryCode: "GB",
			}, nil
		}

		controller := newTestController(
			"GET",
			"/?sessionToken=test-token",
		)

		controller.Ctx.Input.SetParam(":placeId", "place-123")

		controller.GetPlaceDetails()

		if controller.Data["json"] == nil {
			t.Fatal("expected json data")
		}
	})

	t.Run("error", func(t *testing.T) {
		fetchPlaceDetailsData = func(placeID, sessionToken string) (models.CleanLocationResult, error) {
			return models.CleanLocationResult{}, errors.New("place service error")
		}

		controller := newTestController(
			"GET",
			"/?sessionToken=test-token",
		)

		controller.Ctx.Input.SetParam(":placeId", "place-123")

		controller.GetPlaceDetails()

		result, ok := controller.Data["json"].(map[string]string)
		if !ok {
			t.Fatal("expected error response")
		}

		if result["error"] != "place service error" {
			t.Errorf(
				"error = %q, want %q",
				result["error"],
				"place service error",
			)
		}
	})
}
