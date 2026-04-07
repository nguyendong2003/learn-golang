package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// --- Cấu trúc Hoppscotch ---
type HoppHeader struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Active bool   `json:"active"`
}

type HoppParam struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Active bool   `json:"active"`
}

type HoppRequest struct {
	Name             string       `json:"name"`
	Method           string       `json:"method"`
	Endpoint         string       `json:"endpoint"`
	Headers          []HoppHeader `json:"headers"`
	Params           []HoppParam  `json:"params"`
	PreRequestScript string       `json:"preRequestScript"` // Script trước request
	TestScript       string       `json:"testScript"`       // Script sau request (Tests)
	Auth             struct {
		AuthType string `json:"authType"`
		Token    string `json:"token"`
	} `json:"auth"`
	Body struct {
		ContentType string      `json:"contentType"`
		Body        interface{} `json:"body"`
	} `json:"body"`
}

type HoppFolder struct {
	Name     string        `json:"name"`
	Folders  []HoppFolder  `json:"folders"`
	Requests []HoppRequest `json:"requests"`
}

type HoppCollection struct {
	Name    string       `json:"name"`
	Folders []HoppFolder `json:"folders"`
}

// --- Helpers ---

func fixVars(s string) string {
	s = strings.ReplaceAll(s, "<<", "{{")
	s = strings.ReplaceAll(s, ">>", "}}")
	return s
}

// Chuyển đổi cú pháp script từ Hoppscotch sang Postman
func fixScript(s string) string {
	if s == "" {
		return ""
	}
	s = strings.ReplaceAll(s, "hopp.response.body.asJSON()", "pm.response.json()")
	s = strings.ReplaceAll(s, "hopp.env.set(", "pm.environment.set(")
	s = strings.ReplaceAll(s, "hopp.expect", "pm.expect")
	return s
}

func transform(folders []HoppFolder, requests []HoppRequest) []interface{} {
	var items []interface{}

	for _, f := range folders {
		items = append(items, map[string]interface{}{
			"name": f.Name,
			"item": transform(f.Folders, f.Requests),
		})
	}

	for _, r := range requests {
		// 1. Khởi tạo Request Item
		reqItem := map[string]interface{}{
			"name": r.Name,
		}

		// 2. Cấu trúc Request (Method, URL, Headers, Body)
		var postmanHeaders []map[string]interface{}
		for _, h := range r.Headers {
			if h.Active {
				postmanHeaders = append(postmanHeaders, map[string]interface{}{
					"key":   h.Key,
					"value": fixVars(h.Value),
				})
			}
		}

		urlRaw := fixVars(r.Endpoint)
		urlMap := map[string]interface{}{
			"raw":   urlRaw,
			"query": extractQueryParams(r.Params),
		}
		// Tách Host cho đẹp
		if strings.HasPrefix(urlRaw, "{{") {
			parts := strings.SplitN(urlRaw, "/", 2)
			urlMap["host"] = []string{parts[0]}
			if len(parts) > 1 {
				urlMap["path"] = strings.Split(parts[1], "/")
			}
		}

		reqItem["request"] = map[string]interface{}{
			"method": r.Method,
			"header": postmanHeaders,
			"url":    urlMap,
		}

		// 3. Body JSON
		if r.Body.ContentType == "application/json" && r.Body.Body != nil {
			if b, ok := r.Body.Body.(string); ok {
				reqItem["request"].(map[string]interface{})["body"] = map[string]interface{}{
					"mode": "raw",
					"raw":  fixVars(b),
					"options": map[string]interface{}{
						"raw": map[string]interface{}{"language": "json"},
					},
				}
			}
		}

		// 4. Auth Bearer
		if r.Auth.AuthType == "bearer" {
			reqItem["request"].(map[string]interface{})["auth"] = map[string]interface{}{
				"type": "bearer",
				"bearer": []map[string]interface{}{
					{"key": "token", "value": fixVars(r.Auth.Token), "type": "string"},
				},
			}
		}

		// 5. Xử lý Scripts (Event trong Postman)
		var events []map[string]interface{}

		// Pre-request Script
		if r.PreRequestScript != "" {
			events = append(events, map[string]interface{}{
				"listen": "prerequest",
				"script": map[string]interface{}{
					"exec": strings.Split(fixScript(r.PreRequestScript), "\n"),
					"type": "text/javascript",
				},
			})
		}

		// Test Script (Post-response)
		if r.TestScript != "" {
			events = append(events, map[string]interface{}{
				"listen": "test",
				"script": map[string]interface{}{
					"exec": strings.Split(fixScript(r.TestScript), "\n"),
					"type": "text/javascript",
				},
			})
		}

		if len(events) > 0 {
			reqItem["event"] = events
		}

		items = append(items, reqItem)
	}
	return items
}

func extractQueryParams(params []HoppParam) []map[string]interface{} {
	var q []map[string]interface{}
	for _, p := range params {
		if p.Active {
			q = append(q, map[string]interface{}{
				"key":   p.Key,
				"value": fixVars(p.Value),
			})
		}
	}
	return q
}

func main() {
	content, err := os.ReadFile("elearning_hoppscotch.json")
	if err != nil {
		fmt.Println("Vui lòng để file export vào cùng thư mục và đặt tên là elearning_hoppscotch.json")
		return
	}

	var hopp HoppCollection
	json.Unmarshal(content, &hopp)

	postman := map[string]interface{}{
		"info": map[string]string{
			"name":   hopp.Name,
			"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		"item": transform(hopp.Folders, nil),
	}

	out, _ := json.MarshalIndent(postman, "", "  ")
	os.WriteFile("elearning_converted_postman.json", out, 0644)
	fmt.Println("🚀 Thành công! Đã chuyển đổi: Params, Auth, Headers và cả SCRIPTS.")
}
