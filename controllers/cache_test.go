package controllers

import (
	"errors"
	"io"
	"net/http/httptest"
	"testing"

	beegoContext "github.com/beego/beego/v2/server/web/context"
)

type errorReadCloser struct{}

func (e errorReadCloser) Read([]byte) (int, error) {
	return 0, errors.New("request body read error")
}

func (e errorReadCloser) Close() error {
	return nil
}

func newCacheTestController(method, target string) *CacheController {
	req := httptest.NewRequest(method, target, nil)
	recorder := httptest.NewRecorder()

	controller := &CacheController{}
	controller.Ctx = beegoContext.NewContext()
	controller.Ctx.Reset(recorder, req)
	controller.Data = make(map[interface{}]interface{})

	return controller
}

func TestCacheController_Purge(t *testing.T) {
	t.Run("parse form error", func(t *testing.T) {
		controller := newCacheTestController(
			"POST",
			"/cache/purge",
		)

		controller.Ctx.Request.Body = io.NopCloser(errorReadCloser{})
		controller.Ctx.Request.ContentLength = 1
		controller.Ctx.Request.Header.Set(
			"Content-Type",
			"application/x-www-form-urlencoded",
		)

		controller.Purge()

		response, ok := controller.Data["json"].(map[string]interface{})
		if !ok {
			t.Fatal("expected json response")
		}

		if response["success"] != false {
			t.Errorf("success = %v, want false", response["success"])
		}

		if response["message"] != "failed to parse request" {
			t.Errorf(
				"message = %v, want failed to parse request",
				response["message"],
			)
		}

		if response["count"] != 0 {
			t.Errorf("count = %v, want 0", response["count"])
		}
	})

	t.Run("validation error", func(t *testing.T) {
		controller := newCacheTestController(
			"POST",
			"/cache/purge?city=London",
		)

		controller.Purge()

		response, ok := controller.Data["json"].(map[string]interface{})
		if !ok {
			t.Fatal("expected json response")
		}

		if response["success"] != false {
			t.Errorf("success = %v, want false", response["success"])
		}

		if response["count"] != 0 {
			t.Errorf("count = %v, want 0", response["count"])
		}

		if response["message"] == nil {
			t.Fatal("expected validation error message")
		}
	})

	t.Run("clear all cache", func(t *testing.T) {
		controller := newCacheTestController(
			"POST",
			"/cache/purge",
		)

		controller.Purge()

		response, ok := controller.Data["json"].(map[string]interface{})
		if !ok {
			t.Fatal("expected json response")
		}

		if response["success"] != true {
			t.Errorf("success = %v, want true", response["success"])
		}

		if response["count"] == nil {
			t.Fatal("expected count")
		}

		if response["message"] == nil {
			t.Fatal("expected message")
		}
	})

	t.Run("delete specific cache", func(t *testing.T) {
		controller := newCacheTestController(
			"POST",
			"/cache/purge?city=London&countryCode=GB&category=Music",
		)

		controller.Purge()

		response, ok := controller.Data["json"].(map[string]interface{})
		if !ok {
			t.Fatal("expected json response")
		}

		if response["success"] != true {
			t.Errorf("success = %v, want true", response["success"])
		}

		if response["count"] == nil {
			t.Fatal("expected count")
		}

		if response["message"] == nil {
			t.Fatal("expected message")
		}
	})
}
