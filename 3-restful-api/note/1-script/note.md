## Golang
1. Các tag trong struct:
- Name string `json:"name"` -> Bind dữ liệu từ JSON body
- Status string `form:"status"` -> Bind từ query string / form-data / x-www-form-urlencoded
- ID int `uri:"id"` -> Bind từ path param (GET /users/:id)
- Token string `header:"Authorization"` -> Lấy dữ liệu từ HTTP header
- Title string `xml:"title"` -> Dùng khi API chơi XML (ít gặp nhưng vẫn có)
- Port int `yaml:"port"` -> Dùng để đọc config YAML
- Tag cho validation (rất hay dùng):
    + `binding (Gin)` 
        + Một số rule phổ biến: `required`, `omitempty`, `email`, `min=3`, `max=100`, `oneof=active inactive`
        + Ví dụ: Email string `json:"email" binding:"required,email"`
    + ORM / Database tags
        + `GORM`: ID uint `gorm:"primaryKey;autoIncrement"`

2. Thứ tự code các tầng

- business -> storage -> transport 





