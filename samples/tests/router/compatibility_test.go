// Copyright (c) 2026 Hygon Information Technology Co., Ltd.
// SPDX-License-Identifier: Apache-2.0

package router_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/dcgm"
	"github.com/HYGON-AI/hcu-dcgm/v3/pkg/service/router"
)

func TestCompatibleReturnsWarningForVersionMismatch(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := gin.New()
	handler.GET("/compatible", router.Compatible)

	request := httptest.NewRequest(http.MethodGet, "/compatible?cardModel=BW1000&driverVersion=6.3.27-V1.2.5&dtkVersion=25.05.0", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	var body struct {
		Data dcgm.CompatibilityResult `json:"data"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Status != dcgm.CompatibilityWarn {
		t.Fatalf("status = %q, want %q", body.Data.Status, dcgm.CompatibilityWarn)
	}
	if body.Data.RecommendedDriver != ">= 6.3.8" || body.Data.RecommendedDTK != "25.04.*" {
		t.Fatalf("recommendations = (%q, %q)", body.Data.RecommendedDriver, body.Data.RecommendedDTK)
	}
}

func TestCompatibleReturnsErrorForUnknownModel(t *testing.T) {
	gin.SetMode(gin.TestMode)
	handler := gin.New()
	handler.GET("/compatible", router.Compatible)

	request := httptest.NewRequest(http.MethodGet, "/compatible?cardModel=UNKNOWN&driverVersion=6.3.30&dtkVersion=26.04", nil)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
}
