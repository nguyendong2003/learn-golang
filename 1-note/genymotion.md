### Cài đặt Genymotion để làm Android Emulator

1. Bước 1: Cài đặt lại Android SDK (chỉ lấy Platform Tools) qua Terminal

- Mở terminal của Ubuntu lên và chạy lệnh sau để tải thẳng bộ công cụ Android SDK chính thức từ kho ứng dụng của Ubuntu:

```
sudo apt update
sudo apt install android-sdk
```

- Sau khi cài xong, hệ thống sẽ tự động tạo cho bạn một thư mục SDK mặc định tại đường dẫn:

`/usr/lib/android-sdk`

2. Bước 2: Trỏ Genymotion về đường dẫn mới

- Bây giờ thư mục SDK đã có lại, bạn hãy cấu hình cho Genymotion:

- Mở Genymotion > Vào Settings > Chọn tab ADB.

- Chọn Use custom Android SDK tools.

- Tại ô chọn đường dẫn, bạn nhập hoặc bấm Browse tìm đến đúng thư mục:
  `/usr/lib/android-sdk`

- Cửa sổ Genymotion sẽ hiện dấu tích xanh thông báo thư mục hợp lệ (This folder is valid).
