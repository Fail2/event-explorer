package controllers

import (
	"errors"
	"net/http/httptest"
	"testing"

	"event-explorer/models"
	"event-explorer/services"

	beegoContext "github.com/beego/beego/v2/server/web/context"
)

func newEventTestController(method, target string) *EventController {
	req := httptest.NewRequest(method, target, nil)
	recorder := httptest.NewRecorder()

	controller := &EventController{}
	controller.Ctx = beegoContext.NewContext()
	controller.Ctx.Reset(recorder, req)
	controller.Data = make(map[interface{}]interface{})

	return controller
}

func clearTestCache() {
	services.GetCacheInstance().ClearData()
}

func TestEventController_GetEvents(t *testing.T) {
	originalGetConcurrentEvents := getConcurrentEvents
	defer func() {
		getConcurrentEvents = originalGetConcurrentEvents
		clearTestCache()
	}()

	t.Run("invalid query", func(t *testing.T) {
		clearTestCache()

		controller := newEventTestController(
			"GET",
			"/events",
		)

		controller.GetEvents()

		if controller.Ctx.ResponseWriter.Status != 302 {
			t.Errorf(
				"status = %d, want 302",
				controller.Ctx.ResponseWriter.Status,
			)
		}
	})

	t.Run("both cache hits", func(t *testing.T) {
		clearTestCache()

		cache := services.GetCacheInstance()

		musicEvents := []models.UIEvent{
			{ID: "music-1", Name: "Music Event"},
		}

		sportsEvents := []models.UIEvent{
			{ID: "sports-1", Name: "Sports Event"},
		}

		cache.Set("London", "GB", "Music", musicEvents)
		cache.Set("London", "GB", "Sports", sportsEvents)

		getConcurrentEvents = func(city, countryCode string) (
			map[string][]models.UIEvent,
			map[string]error,
		) {
			t.Fatal("getConcurrentEvents should not be called")
			return nil, nil
		}

		controller := newEventTestController(
			"GET",
			"/events?city=London&countryCode=GB",
		)

		controller.GetEvents()

		if controller.Data["City"] != "London" {
			t.Errorf("City = %v, want London", controller.Data["City"])
		}

		if controller.Data["CountryCode"] != "GB" {
			t.Errorf("CountryCode = %v, want GB", controller.Data["CountryCode"])
		}

		if controller.Data["MusicCacheStatus"] != "Cache Hit" {
			t.Errorf("MusicCacheStatus = %v, want Cache Hit", controller.Data["MusicCacheStatus"])
		}

		if controller.Data["SportsCacheStatus"] != "Cache Hit" {
			t.Errorf("SportsCacheStatus = %v, want Cache Hit", controller.Data["SportsCacheStatus"])
		}

		if controller.TplName != "listing.tpl" {
			t.Errorf("TplName = %v, want listing.tpl", controller.TplName)
		}
	})

	t.Run("both fresh data", func(t *testing.T) {
		clearTestCache()

		getConcurrentEvents = func(city, countryCode string) (
			map[string][]models.UIEvent,
			map[string]error,
		) {
			return map[string][]models.UIEvent{
				"music": {
					{ID: "music-1", Name: "Music Event"},
				},
				"sports": {
					{ID: "sports-1", Name: "Sports Event"},
				},
			}, map[string]error{}
		}

		controller := newEventTestController(
			"GET",
			"/events?city=London&countryCode=GB",
		)

		controller.GetEvents()

		if controller.Data["MusicCacheStatus"] != "Fresh Data" {
			t.Errorf("MusicCacheStatus = %v, want Fresh Data", controller.Data["MusicCacheStatus"])
		}

		if controller.Data["SportsCacheStatus"] != "Fresh Data" {
			t.Errorf("SportsCacheStatus = %v, want Fresh Data", controller.Data["SportsCacheStatus"])
		}

		if controller.TplName != "listing.tpl" {
			t.Errorf("TplName = %v, want listing.tpl", controller.TplName)
		}
	})

	t.Run("music error and sports success", func(t *testing.T) {
		clearTestCache()

		getConcurrentEvents = func(city, countryCode string) (
			map[string][]models.UIEvent,
			map[string]error,
		) {
			return map[string][]models.UIEvent{
					"sports": {
						{ID: "sports-1", Name: "Sports Event"},
					},
				},
				map[string]error{
					"music": errors.New("music service error"),
				}
		}

		controller := newEventTestController(
			"GET",
			"/events?city=London&countryCode=GB",
		)

		controller.GetEvents()

		if controller.Data["MusicCacheStatus"] != "Expired Data" {
			t.Errorf("MusicCacheStatus = %v, want Expired Data", controller.Data["MusicCacheStatus"])
		}

		if controller.Data["MusicError"] == nil {
			t.Fatal("expected MusicError")
		}

		if controller.Data["SportsCacheStatus"] != "Fresh Data" {
			t.Errorf("SportsCacheStatus = %v, want Fresh Data", controller.Data["SportsCacheStatus"])
		}
	})

	t.Run("music success and sports error", func(t *testing.T) {
		clearTestCache()

		getConcurrentEvents = func(city, countryCode string) (
			map[string][]models.UIEvent,
			map[string]error,
		) {
			return map[string][]models.UIEvent{
					"music": {
						{ID: "music-1", Name: "Music Event"},
					},
				},
				map[string]error{
					"sports": errors.New("sports service error"),
				}
		}

		controller := newEventTestController(
			"GET",
			"/events?city=London&countryCode=GB",
		)

		controller.GetEvents()

		if controller.Data["MusicCacheStatus"] != "Fresh Data" {
			t.Errorf("MusicCacheStatus = %v, want Fresh Data", controller.Data["MusicCacheStatus"])
		}

		if controller.Data["SportsCacheStatus"] != "Expired Data" {
			t.Errorf("SportsCacheStatus = %v, want Expired Data", controller.Data["SportsCacheStatus"])
		}

		if controller.Data["SportsError"] == nil {
			t.Fatal("expected SportsError")
		}
	})

	t.Run("both errors", func(t *testing.T) {
		clearTestCache()

		getConcurrentEvents = func(city, countryCode string) (
			map[string][]models.UIEvent,
			map[string]error,
		) {
			return map[string][]models.UIEvent{},
				map[string]error{
					"music":  errors.New("music error"),
					"sports": errors.New("sports error"),
				}
		}

		controller := newEventTestController(
			"GET",
			"/events?city=London&countryCode=GB",
		)

		controller.GetEvents()

		if controller.Data["MusicError"] == nil {
			t.Fatal("expected MusicError")
		}

		if controller.Data["SportsError"] == nil {
			t.Fatal("expected SportsError")
		}

		if controller.Data["MusicCacheStatus"] != "Expired Data" {
			t.Errorf("MusicCacheStatus = %v, want Expired Data", controller.Data["MusicCacheStatus"])
		}

		if controller.Data["SportsCacheStatus"] != "Expired Data" {
			t.Errorf("SportsCacheStatus = %v, want Expired Data", controller.Data["SportsCacheStatus"])
		}
	})

	t.Run("music cache hit and sports fresh", func(t *testing.T) {
		clearTestCache()

		cache := services.GetCacheInstance()

		cache.Set(
			"London",
			"GB",
			"Music",
			[]models.UIEvent{
				{ID: "music-cached"},
			},
		)

		getConcurrentEvents = func(city, countryCode string) (
			map[string][]models.UIEvent,
			map[string]error,
		) {
			return map[string][]models.UIEvent{
				"sports": {
					{ID: "sports-fresh"},
				},
			}, map[string]error{}
		}

		controller := newEventTestController(
			"GET",
			"/events?city=London&countryCode=GB",
		)

		controller.GetEvents()

		if controller.Data["MusicCacheStatus"] != "Cache Hit" {
			t.Errorf("MusicCacheStatus = %v, want Cache Hit", controller.Data["MusicCacheStatus"])
		}

		if controller.Data["SportsCacheStatus"] != "Fresh Data" {
			t.Errorf("SportsCacheStatus = %v, want Fresh Data", controller.Data["SportsCacheStatus"])
		}
	})

	t.Run("sports cache hit and music fresh", func(t *testing.T) {
		clearTestCache()

		cache := services.GetCacheInstance()

		cache.Set(
			"London",
			"GB",
			"Sports",
			[]models.UIEvent{
				{ID: "sports-cached"},
			},
		)

		getConcurrentEvents = func(city, countryCode string) (
			map[string][]models.UIEvent,
			map[string]error,
		) {
			return map[string][]models.UIEvent{
				"music": {
					{ID: "music-fresh"},
				},
			}, map[string]error{}
		}

		controller := newEventTestController(
			"GET",
			"/events?city=London&countryCode=GB",
		)

		controller.GetEvents()

		if controller.Data["SportsCacheStatus"] != "Cache Hit" {
			t.Errorf("SportsCacheStatus = %v, want Cache Hit", controller.Data["SportsCacheStatus"])
		}

		if controller.Data["MusicCacheStatus"] != "Fresh Data" {
			t.Errorf("MusicCacheStatus = %v, want Fresh Data", controller.Data["MusicCacheStatus"])
		}
	})
}

func TestEventController_GetEventDetails(t *testing.T) {
	originalFetch := fetchTicketmasterEventDetails
	defer func() {
		fetchTicketmasterEventDetails = originalFetch
	}()

	t.Run("missing event id", func(t *testing.T) {
		controller := newEventTestController(
			"GET",
			"/events",
		)

		controller.GetEventDetails()

		if controller.Ctx.ResponseWriter.Status != 302 {
			t.Errorf(
				"status = %d, want 302",
				controller.Ctx.ResponseWriter.Status,
			)
		}
	})

	t.Run("service error", func(t *testing.T) {
		fetchTicketmasterEventDetails = func(eventID string) (models.UIEvent, error) {
			return models.UIEvent{}, errors.New("ticketmaster error")
		}

		controller := newEventTestController(
			"GET",
			"/events/123",
		)

		controller.Ctx.Input.SetParam(":id", "123")

		controller.GetEventDetails()

		if controller.Data["Error"] == nil {
			t.Fatal("expected Error data")
		}

		if controller.TplName != "error.tpl" {
			t.Errorf("TplName = %v, want error.tpl", controller.TplName)
		}
	})

	t.Run("success", func(t *testing.T) {
		fetchTicketmasterEventDetails = func(eventID string) (models.UIEvent, error) {
			return models.UIEvent{
				ID:   "event-123",
				Name: "Test Event",
			}, nil
		}

		controller := newEventTestController(
			"GET",
			"/events/123",
		)

		controller.Ctx.Input.SetParam(":id", "123")

		controller.GetEventDetails()

		if controller.Data["Event"] == nil {
			t.Fatal("expected Event data")
		}

		if controller.TplName != "details.tpl" {
			t.Errorf("TplName = %v, want details.tpl", controller.TplName)
		}
	})
}
