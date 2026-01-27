package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strings"
)

/*
| Interface                              | Ai gọi              | Khi nào                |
| -------------------------------------- | ------------------- | ---------------------- |
| `sql.Scanner` → `Scan()`               | GORM                | **SELECT từ DB**       |
| `driver.Valuer` → `Value()`            | GORM                | **INSERT / UPDATE DB** |
| `json.Marshaler` → `MarshalJSON()`     | Gin / encoding/json | **Trả response**       |
| `json.Unmarshaler` → `UnmarshalJSON()` | Gin                 | **Nhận request JSON**  |

*/

/* Luồng chạy
REQUEST JSON
   ↓
UnmarshalJSON()   ← Gin gọi
   ↓
ItemStatus (0,1,2)
   ↓
Value()           ← GORM gọi
   ↓
DB ("Doing")

----------------------------

DB ("Doing")
   ↓
Scan()            ← GORM gọi
   ↓
ItemStatus (0,1,2)
   ↓
MarshalJSON()     ← Gin gọi
   ↓
RESPONSE JSON
*/

type ItemStatus int

/*
iota là một hằng số đặc biệt trong Go, dùng để tự động tăng số khi khai báo const. Nó chỉ dùng được trong const block.
Trong const block, nếu bỏ trống giá trị, Go sẽ tự dùng lại biểu thức ở dòng trên, nhưng iota thì tự tăng.
*/
const (
	ItemStatusDoing   ItemStatus = iota // 0
	ItemStatusDone                      // iota -> 1
	ItemStatusDeleted                   // iota -> 2
)

var allItemStatuses = [3]string{"Doing", "Done", "Deleted"}

// String() – dùng khi in / trả JSON (map từ số 0,1,2 sang "Doing", "Done", "Deleted")
func (item *ItemStatus) String() string {
	return allItemStatuses[*item]
}

func parseStrToItemStatus(s string) (ItemStatus, error) {
	for index := range allItemStatuses {
		if allItemStatuses[index] == s {
			return ItemStatus(index), nil
		}
	}

	return ItemStatus(0), errors.New("Invalid status string")
}

// convert string → enum (GORM tự động gọi khi: SELECT từ DB)
func (item *ItemStatus) Scan(value interface{}) error {
	bytes, ok := value.([]byte)

	if !ok {
		return fmt.Errorf("Fail to scan data from sql: %s", value)
	}

	v, err := parseStrToItemStatus(string(bytes))

	if err != nil {
		return errors.New(err.Error())
	}

	*item = v // *item = 0,1,2

	return nil
}

func (item *ItemStatus) Value() (driver.Value, error) {
	if item == nil {
		return nil, nil
	}

	return item.String(), nil
}

// Nếu không có method này thì khi gọi api get item thì status sẽ trả về số 0,1,2 thay vì "Doing", "Done", "Deleted" (GORM gọi hàm này)
func (item *ItemStatus) MarshalJSON() ([]byte, error) {
	if item == nil {
		return nil, nil
	}

	return []byte(fmt.Sprintf("\"%s\"", item.String())), nil
}

func (item *ItemStatus) UnmarshalJSON(data []byte) error {
	str := strings.ReplaceAll(string(data), "\"", "")

	itemValue, err := parseStrToItemStatus(str)

	if err != nil {
		return err
	}

	*item = itemValue

	return nil
}
