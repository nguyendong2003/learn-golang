# Hướng dẫn thiết lập React Native CLI với Genymotion trên Ubuntu (Không cần Android Studio)

Tài liệu này hướng dẫn chi tiết các bước thiết lập môi trường phát triển ứng dụng di động bằng React Native CLI trên hệ điều hành Ubuntu, sử dụng máy ảo Genymotion làm thiết bị kiểm thử và quản lý SDK thông qua Command Line Tools chính thức từ Google.

---

## Bước 1: Gỡ bỏ hoàn toàn cấu hình cũ (Nếu có)

Nếu hệ thống của bạn từng cài đặt gói Android SDK qua trình quản lý gói `apt` của Ubuntu, hãy dọn dẹp sạch sẽ để tránh xung đột đường dẫn:

```bash
sudo apt purge android-sdk android-sdk-platform-tools-common -y
sudo apt autoremove -y
sudo rm -rf /usr/lib/android-sdk
```

## Bước 2: Chuẩn bị môi trường Java

React Native CLI yêu cầu môi trường Java mã nguồn mở để biên dịch mã nguồn Android. Cài đặt OpenJDK 17 bằng lệnh:

```bash
sudo apt update
sudo apt install openjdk-17-jdk -y
```

- Kiểm tra phiên bản Java vừa cài:

```bash
java -version
```

- Nếu hiển thị OpenJDK 17 thì ok, nếu khác version thì thực hiện các bước sau:

```bash
sudo update-alternatives --config java
```

- Sau đó nhập số để chọn OpenJDK 17

## Bước 3: Cài đặt Android SDK

- Vào trang https://developer.android.com/studio/index.html#command-line-tools-only và tải file zip cho linux (commandlinetools-linux-14742923_latest.zip)

- Tạo thư mục để chứa Android SDK:

```bash
mkdir -p ~/Android/Sdk/cmdline-tools
```

- Giải nén file zip vào thư mục vừa tạo trong folder `cmdline-tools` và đổi tên folder vừa giải nén thành `latest`

## Bước 4: Cấu hình biến môi trường

- Mở file ~/.bashrc để khai báo đường dẫn SDK:

```bash
nano ~/.bashrc
```

- Dán đoạn code sau vào cuối file:

```bash
export ANDROID_HOME=$HOME/Android/Sdk
export PATH=$PATH:$ANDROID_HOME/cmdline-tools/latest/bin
export PATH=$PATH:$ANDROID_HOME/platform-tools
```

- Lưu lại (Ctrl+O, Enter, Ctrl+X) và cập nhật Terminal:

```bash
source ~/.bashrc
```

## Bước 5: Tải các thành phần SDK cần thiết cho React Native

- Bây giờ bạn đã có lệnh `sdkmanager`. Chúng ta sẽ tiến hành tải Build-Tools và Platform-Tools (thứ mà React Native cần để build app).

1. Chấp nhận bản quyền của Google (Bắt buộc):

```bash
sdkmanager --licenses
```

- (Bấm y và Enter liên tục cho đến khi xong).

2. Tải các gói SDK (React Native hiện tại cần API 34):

```bash
sdkmanager "platform-tools" "platforms;android-34" "build-tools;34.0.0"
```

## Bước 6: Cấu hình Genymotion

1. Nếu đang mở Genymotion thì tắt nó đi và khởi động lại
2. Vào Settings -> Chọn tab ADB.
3. Chọn Use custom Android SDK tools.
4. Trỏ đường dẫn đến thư mục: /home/<ubuntu-username>/Android/Sdk (thay <ubuntu-username> bằng username Ubuntu của bạn nếu khác).
5. Bật một máy ảo Android trên Genymotion lên.
6. Kiểm tra xem máy ảo Genymotion đã kết nối với máy tính chưa:

```bash
adb devices
```

- Nếu hiện như này là ok

```bash
List of devices attached
127.0.0.1:6555	device
```

## Bước 7: Tạo Project React Native và Chạy (https://reactnative.dev/docs/getting-started-without-a-framework)

1. Tạo project react native

```bash
npx @react-native-community/cli init MyAwesomeApp
cd MyAwesomeApp
```

2. Chạy ứng dụng:

- Mở terminal 1 tại thư mục project: `npm start`
- Mở terminal 2 tại thư mục project: `npm run android`
- Lúc này dự án của bạn sẽ tự động build thông qua các công cụ dòng lệnh vừa tải và nạp thẳng vào máy ảo Genymotion một cách mượt mà!

## Mẹo xử lý khi App bị "đơ" hoặc không tự reload
1. Khi bạn chạy dự án React Native, cơ chế hoạt động của nó cực kỳ tiện lợi nhờ vào tính năng gọi là Fast Refresh (Tự động tải lại khi sửa code). Bạn không cần phải build lại từ đầu hay chạy lại lệnh `npm run android` mỗi khi thay đổi code.

2. Đôi khi bạn sửa các file cấu hình cấu trúc lớn hoặc file native, tính năng Fast Refresh có thể không tự nhận diện được. Lúc này bạn xử lý như sau:
- Reload lại giao diện JavaScript: Click chuột vào màn hình máy ảo Genymotion và ấn phím R hai lần liên tiếp (hoặc bấm Ctrl + M để mở Developer Menu -> chọn Reload).
- Build lại từ đầu: Nếu bạn cài thêm một thư viện mới có can thiệp vào code hệ thống (ví dụ thư viện camera, bluetooth, định vị...), bạn bắt buộc phải tắt Terminal 1 đi, bật lại, rồi chạy lại lệnh `npm run android` ở Terminal 2 để nó nạp code native mới vào máy ảo.
