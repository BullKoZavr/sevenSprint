package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// var cafeList = map[string][]string{
// 	"moscow": []string{"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
// }

// func mainHandle(w http.ResponseWriter, req *http.Request) {
// 	countStr := req.URL.Query().Get("count")
// 	if countStr == "" {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("count missing"))
// 		return
// 	}

// 	count, err := strconv.Atoi(countStr)
// 	if err != nil {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("wrong count value"))
// 		return
// 	}

// 	city := req.URL.Query().Get("city")

// 	cafe, ok := cafeList[city]
// 	if !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte("wrong city value"))
// 		return
// 	}

// 	if count > len(cafe) {
// 		count = len(cafe)
// 	}

// 	answer := strings.Join(cafe[:count], ",")

// 	w.WriteHeader(http.StatusOK)
// 	w.Write([]byte(answer))
// }

func TestMainHandlerReturns200AndNonEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?count=2&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Проверяем статус ответа
	require.Equal(t, http.StatusOK, responseRecorder.Code)
	// Проверяем, что тело ответа не пустое
	assert.NotEmpty(t, responseRecorder.Body.String())
}

func TestMainHandlerReturns400AndWrongCityValue(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?count=2&city=undefined", nil)

	reqponseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(reqponseRecorder, req)

	require.Equal(t, http.StatusBadRequest, reqponseRecorder.Code)
	assert.Equal(t, "wrong city value", reqponseRecorder.Body.String())

}
func TestMainHandlerWhenCountMoreThanTotal(t *testing.T) {
	totalCount := 4

	req := httptest.NewRequest(http.MethodGet, "/?count=5&city=moscow", nil)

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	// Запрос сформирован корректно, сервис возвращает код ответа 200 и тело ответа не пустое.
	require.Equal(t, http.StatusOK, responseRecorder.Code)
	// Город, который передаётся в параметре city, не поддерживается. Сервис возвращает код ответа 400 и ошибку wrong city value в теле ответа.
	expected := "Мир кофе,Сладкоежка,Кофе и завтраки,Сытый студент"
	assert.Equal(t, expected, responseRecorder.Body.String())
	// Если в параметре count указано больше, чем есть всего, должны вернуться все доступные кафе.
	expectedSlice := strings.Split(responseRecorder.Body.String(), ",")
	assert.Len(t, expectedSlice, totalCount)

}
