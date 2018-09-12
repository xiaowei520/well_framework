package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"reflect"
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

func Test_setupRouter(t *testing.T) {
	tests := []struct {
		name string
		want int
	}{
		// TODO: Add test cases.
		{
			name: "111",
			want: 200,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := setupRouter()
			req, _ := http.NewRequest("GET", "/ping", nil)

			w := httptest.NewRecorder()
			got.ServeHTTP(w, req)

			if !reflect.DeepEqual(w.Code, tt.want) {

				t.Errorf("setupRouter() = %v, want %v", w.Code, tt.want)
			} else {
				t.Log("setupRouter() = %v, want %v", w.Code, tt.want)
			}
		})
	}
}
