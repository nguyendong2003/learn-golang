package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
)

// 1. Adaptee – thư viện analytics (JSON only)
type JsonAnalyticsLibrary struct{}

func (j *JsonAnalyticsLibrary) AnalyzeJSON(jsonData string) {
	fmt.Println("Analyzing JSON data:", jsonData)
}

// 2. Target interface – app đang dùng
type XmlAnalytics interface {
	AnalyzeXML(xmlData string)
}

// 3. Adapter – chuyển XML → JSON
// Struct trung gian để parse XML
type Stock struct {
	XMLName xml.Name `xml:"stock"`
	Name    string   `xml:"name" json:"name"`
	Price   float64  `xml:"price" json:"price"`
}

type XmlToJsonAnalyticsAdapter struct {
	jsonAnalytics *JsonAnalyticsLibrary
}

func NewXmlToJsonAnalyticsAdapter(lib *JsonAnalyticsLibrary) *XmlToJsonAnalyticsAdapter {
	return &XmlToJsonAnalyticsAdapter{
		jsonAnalytics: lib,
	}
}

func (a *XmlToJsonAnalyticsAdapter) AnalyzeXML(xmlData string) {
	// 1. Parse XML
	var stock Stock
	xml.Unmarshal([]byte(xmlData), &stock)

	// 2. Convert sang JSON
	jsonBytes, _ := json.Marshal(stock)

	// 3. Gọi thư viện JSON
	a.jsonAnalytics.AnalyzeJSON(string(jsonBytes))
}

// 4. Client – Stock Market App
func main() {
	xmlData := `
	<stock>
		<name>AAPL</name>
		<price>150.5</price>
	</stock>`

	jsonLib := &JsonAnalyticsLibrary{}
	adapter := NewXmlToJsonAnalyticsAdapter(jsonLib)

	// App chỉ biết XML
	adapter.AnalyzeXML(xmlData)
}
