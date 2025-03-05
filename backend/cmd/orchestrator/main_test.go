package main_test

import (
	"bytes"
	"calc-website/config"
	"calc-website/internal/models"
	"calc-website/internal/orchestrator"
	"calc-website/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func startTestServer() *httptest.Server {
	service := orchestrator.NewAPIService(&config.Config{
		TimeAdditionMs:        100,
		TimeSubtractionMs:     100,
		TimeMultiplicationsMs: 100,
		TimeDivisionsMs:       100,
		ComputingPower:        10,
	})
	handler := orchestrator.NewAPIHandler(service)

	server := httptest.NewServer(handler.Router())

	return server
}

func checkStatusCode(t *testing.T, resp *http.Response, expected int) {
	t.Helper()
	if resp.StatusCode != expected {
		t.Errorf("Ожидался статус-код %d, но получен %d", expected, resp.StatusCode)
	}
}

func TestCalculateExpression(t *testing.T) {
	server := startTestServer()
	defer server.Close()

	requestBody, _ := json.Marshal(models.ExpressionRequest{Expression: "2 + 2 * 2"})
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/calculate", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer utils.CloseResponseBody(resp.Body)

	checkStatusCode(t, resp, http.StatusCreated)

	var responseMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		log.Printf("decode task error: %v", err.Error())
	}

	if _, ok := responseMap["expression"]; !ok {
		t.Error("Ответ не содержит id выражения")
	}
}

func TestGetExpressions(t *testing.T) {
	server := startTestServer()
	defer server.Close()

	req, _ := http.NewRequest("GET", server.URL+"/api/v1/expressions", nil)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer utils.CloseResponseBody(resp.Body)

	checkStatusCode(t, resp, http.StatusOK)

	var responseMap map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		log.Printf("decode task error: %v", err.Error())
	}

	if _, ok := responseMap["expressions"]; !ok {
		t.Error("Ответ не содержит списка выражений")
	}
}

func TestGetExpressionByID(t *testing.T) {
	server := startTestServer()
	defer server.Close()

	client := &http.Client{}

	requestBody, _ := json.Marshal(map[string]string{"expression": "2 + 2 * 2"})
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/calculate", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer utils.CloseResponseBody(resp.Body)

	checkStatusCode(t, resp, http.StatusCreated)

	var responseMap map[string]any
	err = json.NewDecoder(resp.Body).Decode(&responseMap)
	if err != nil {
		t.Fatal("Ошибка декодирования JSON:", err)
	}

	exprMap, ok := responseMap["expression"].(map[string]any)
	if !ok {
		t.Fatal("Ошибка: responseMap[\"expression\"] не является map[string]any")
	}

	id, ok := exprMap["id"].(string)
	if !ok {
		t.Fatal("Ошибка: exprMap[\"id\"] не является строкой")
	}

	req, _ = http.NewRequest("GET", server.URL+"/api/v1/expressions/"+id, nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer utils.CloseResponseBody(resp.Body)

	checkStatusCode(t, resp, http.StatusOK)

	var expression map[string]any
	err = json.NewDecoder(resp.Body).Decode(&expression)
	if err != nil {
		t.Fatal("Ошибка декодирования JSON:", err)
	}

	if _, ok := expression["expression"]; !ok {
		t.Error("Ответ не содержит выражение")
	}
}

// 📌 Тест на получение задачи для вычисления (GET /internal/task)
func TestGetTask(t *testing.T) {
	server := startTestServer()
	defer server.Close()

	requestBody, _ := json.Marshal(map[string]string{"expression": "2 + 2 * 2"})
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/calculate", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer utils.CloseResponseBody(resp.Body)

	checkStatusCode(t, resp, http.StatusCreated)

	req, _ = http.NewRequest("GET", server.URL+"/internal/task", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer utils.CloseResponseBody(resp.Body)

	checkStatusCode(t, resp, http.StatusOK)

	var task map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		log.Printf("decode task error: %v", err.Error())
	}

	if _, ok := task["id"]; !ok {
		t.Error("Ответ не содержит id задачи")
	}
}

func TestPostTaskResult(t *testing.T) {
	server := startTestServer()
	defer server.Close()

	requestBody, _ := json.Marshal(map[string]string{"expression": "2 + 2 * 2"})
	req, _ := http.NewRequest("POST", server.URL+"/api/v1/calculate", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer utils.CloseResponseBody(resp.Body)

	checkStatusCode(t, resp, http.StatusCreated)

	req, _ = http.NewRequest("GET", server.URL+"/internal/task", nil)
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer utils.CloseResponseBody(resp.Body)

	checkStatusCode(t, resp, http.StatusOK)

	var task map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&task)
	if err != nil {
		log.Printf("decode task error: %v", err.Error())
	}

	resultData := map[string]interface{}{
		"id":     task["id"],
		"result": 6.0,
	}
	resultBody, _ := json.Marshal(resultData)

	req, _ = http.NewRequest("POST", server.URL+"/internal/task", bytes.NewBuffer(resultBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	checkStatusCode(t, resp, http.StatusOK)
}
