package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
)

//测试用例  - 测试 路由ping 是否返回正确
func TestPingRoute(t *testing.T) {
	router := setupRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	if w.Code == 200 {
		t.Log("success")
	} else {
		t.Error(fmt.Sprintf("code is %d", w.Code))
	}

	if w.Body.String() == "pong" {
		t.Log("success")
	} else {
		t.Error(fmt.Sprintf("body is %s", w.Body.String()))
	}

}
