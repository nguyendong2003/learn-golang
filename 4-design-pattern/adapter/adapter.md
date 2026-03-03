# Adapter Design Pattern 
- Adapter Design Pattern (mẫu thiết kế bộ chuyển đổi) là một structural pattern dùng để kết nối hai interface không tương thích với nhau, để các class/struct có thể làm việc chung mà không cần sửa code gốc.

- Nói đơn giản: Adapter đóng vai “ổ chuyển đổi” 🔌 – giống như khi bạn có phích cắm 3 chấu nhưng ổ điện chỉ nhận 2 chấu.

- Nói đời thường: Adapter giống như cái đầu chuyển sạc 🔌
– ổ điện không đổi
– sạc không đổi
– thêm adapter là dùng được

## Viết lại cho “sách giáo khoa” hơn một chút

Adapter Pattern cho phép các object có interface không tương thích làm việc với nhau bằng cách đưa vào một lớp trung gian (Adapter) để chuyển đổi interface của Adaptee thành interface mà Client mong đợi.

## Khi nào nên dùng Adapter?
- Trong Golang, Adapter rất hay dùng khi:
    + Dùng thư viện bên thứ 3 nhưng interface không khớp
    + Có legacy code không thể sửa
    + Muốn tuân theo interface mà client đã định nghĩa
    + Áp dụng Dependency Inversion (code phụ thuộc interface, không phụ thuộc implementation)

## Ý tưởng cốt lõi
- Bạn có Client cần dùng một interface A
- Nhưng object thực tế (Adaptee) lại cung cấp interface B (không khớp)
- Adapter đứng giữa, chuyển đổi A → B

## Cấu trúc
- Target: interface mà client mong đợi
- Adaptee: class có sẵn, nhưng interface không phù hợp
- Adapter: implements Target, và wrap Adaptee bên trong

## Ví dụ
- Hãy tưởng tượng bạn đang xây dựng một ứng dụng theo dõi thị trường chứng khoán. Ứng dụng tải dữ liệu chứng khoán từ nhiều nguồn khác nhau dưới định dạng XML, sau đó hiển thị cho người dùng các biểu đồ và sơ đồ trực quan, đẹp mắt.

- Đến một thời điểm, bạn quyết định cải tiến ứng dụng bằng cách tích hợp một thư viện phân tích thông minh của bên thứ ba. Tuy nhiên, có một vấn đề: thư viện phân tích này chỉ hoạt động với dữ liệu ở định dạng JSON.

- Ánh xạ sang các thành phần Adapter pattern

| Vai trò     | Trong bài toán                                  |
| ----------- | ----------------------------------------------- |
| **Client**  | Stock Market App                                |
| **Target**  | Interface mà app mong đợi (XML-based analytics) |
| **Adaptee** | 3rd-party Analytics Library (JSON-based)        |
| **Adapter** | XML → JSON Adapter                              |

- Code Java:
```java
// 1. Analytics library (Adaptee – không sửa được)
class JsonAnalyticsLibrary {
    public void analyzeJson(String jsonData) {
        System.out.println("Analyzing JSON data: " + jsonData);
    }
}

// 2. Target interface (app đang dùng)
interface XmlAnalytics {
    void analyzeXml(String xmlData);
}

// 3. Adapter (chìa khóa giải quyết)
class XmlToJsonAnalyticsAdapter implements XmlAnalytics {

    private JsonAnalyticsLibrary jsonAnalytics;

    public XmlToJsonAnalyticsAdapter(JsonAnalyticsLibrary jsonAnalytics) {
        this.jsonAnalytics = jsonAnalytics;
    }

    @Override
    public void analyzeXml(String xmlData) {
        // 1. Convert XML → JSON
        String jsonData = convertXmlToJson(xmlData);

        // 2. Delegate cho thư viện JSON
        jsonAnalytics.analyzeJson(jsonData);
    }

    private String convertXmlToJson(String xmlData) {
        // Giả lập chuyển đổi (thực tế dùng Jackson, Gson, org.json...)
        return "{ \"stock\": \"AAPL\", \"price\": 150 }";
    }
}

// 4. Client code (Stock App – không cần biết JSON tồn tại)
public class StockApp {
    public static void main(String[] args) {
        XmlAnalytics analytics =
            new XmlToJsonAnalyticsAdapter(new JsonAnalyticsLibrary());

        analytics.analyzeXml("<stock><name>AAPL</name><price>150</price></stock>");
    }
}
```