package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}

func checkRsSize(rs string) int {
	if rs != "" {
		return len(strings.Split(rs, ","))
	} else {
		return 0
	}

}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "tula"
	cafeListLen := len(cafeList[city])
	queryPath := "/cafe?city=" + city + "&count="

	requests := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, 100},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		target := queryPath + strconv.Itoa(v.count)
		request := httptest.NewRequest("GET", target, nil)

		handler.ServeHTTP(response, request)

		rs := strings.TrimSpace(response.Body.String())
		rsSize := checkRsSize(rs)

		require.Equal(t, http.StatusOK, response.Code)
		if v.count == 100 {
			assert.Equal(t, cafeListLen, rsSize)
		} else {
			assert.Equal(t, v.want, rsSize)
		}
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)
	city := "moscow"
	queryPath := "/cafe?city=" + city + "&search="

	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе в ответе
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		target := queryPath + v.search
		request := httptest.NewRequest("GET", target, nil)

		handler.ServeHTTP(response, request)

		responseLower := strings.ToLower(response.Body.String())

		require.Equal(t, http.StatusOK, response.Code)

		rsSize := checkRsSize(responseLower)

		if rsSize != 0 {
			assert.True(t, strings.Contains(responseLower, v.search))
		}
		assert.Equal(t, v.wantCount, rsSize)
	}
}
