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

# Hướng dẫn build React Native CLI (https://reactnative.dev/docs/signed-apk-android)

## BƯỚC 1: Tạo khóa ký ứng dụng (Signing Key) ngay tại thư mục đích

1. Mở Terminal trên máy tính của bạn (đảm bảo đang ở thư mục gốc của dự án AwesomeProject).

2. Di chuyển thẳng vào thư mục android/app bằng lệnh:

```Bash
cd android/app
```

3. Chạy lệnh keytool sau để tạo file khóa (lúc này file sinh ra sẽ nằm cố định, đúng vị trí 100%) -> câu lệnh này tạo ra file `my-upload-key.keystore`:

```bash
keytool -genkeypair -v -storetype PKCS12 -keystore my-upload-key.keystore -alias my-key-alias -keyalg RSA -keysize 2048 -validity 10000
```

4. Lưu ý khi chạy lệnh:

- Hệ thống yêu cầu nhập mật khẩu: Bạn gõ mã bí mật của bạn (ví dụ: 123456). Lưu ý: Khi gõ mật khẩu trên terminal Linux/macOS sẽ không hiển thị ký tự dấu sao \*, bạn cứ gõ bình thường rồi nhấn Enter.
- Các câu hỏi tiếp theo (Họ tên, Tổ chức...): Nhấn Enter để bỏ qua hết.
- Câu hỏi xác nhận cuối cùng: Gõ yes (hoặc y) rồi nhấn Enter.

## BƯỚC 2: Cấu hình mật khẩu vào file gradle.properties

- Quay trở lại thư mục gốc của dự án hoặc mở VS Code lên, tìm đến file `📁 android/gradle.properties`
- Kéo xuống dưới cùng của file đó, xuống dòng và dán vào 4 dòng mã sau (thay `123456` bằng mật khẩu bạn vừa đặt ở Bước 1):

```bash
MYAPP_UPLOAD_STORE_FILE=my-upload-key.keystore
MYAPP_UPLOAD_KEY_ALIAS=my-key-alias
MYAPP_UPLOAD_STORE_PASSWORD=123456
MYAPP_UPLOAD_KEY_PASSWORD=123456
```

## BƯỚC 3: Cấu hình file build.gradle (File bên trong thư mục app)

1. Mở file `android/app/build.gradle` lên.
2. Tìm đến đoạn có chữ `signingConfigs` (thường ở khoảng dòng 60).
3. Thêm block release vào trong `signingConfigs`, và sửa dòng `signingConfig` trong `buildTypes -> release` giống hệt như thế này:

```gradle
android {
    ...
    defaultConfig { ... }

    signingConfigs {
        debug {
            storeFile file('debug.keystore')
            storePassword 'android'
            keyAlias 'androidexternal'
            keyPassword 'android'
        }

        // 1. DÁN THÊM BLOCK RELEASE NÀY VÀO ĐÂY:
        release {
            if (project.hasProperty('MYAPP_UPLOAD_STORE_FILE')) {
                storeFile file(MYAPP_UPLOAD_STORE_FILE)
                storePassword MYAPP_UPLOAD_STORE_PASSWORD
                keyAlias MYAPP_UPLOAD_KEY_ALIAS
                keyPassword MYAPP_UPLOAD_KEY_PASSWORD
            }
        }
    }

    buildTypes {
        debug {
            signingConfig signingConfigs.debug
        }
        release {
            // 2. SỬA DÒNG NÀY: Thay đổi từ signingConfigs.debug thành signingConfigs.release
            signingConfig signingConfigs.release

            minifyEnabled enableProguardInReleaseBuilds
            proguardFiles getDefaultProguardFile("proguard-android.txt"), "proguard-rules.pro"
        }
    }
}
```

## BƯỚC 4: Chạy lệnh dọn dẹp và đóng gói APK

1. Di chuyển đến thư mục `android/app`

```bash
cd android
```

2. Chạy lệnh xóa toàn bộ cache lỗi cũ:

```bash
./gradlew clean
```

3. Chạy lệnh build file APK Release hoàn chỉnh:

```bash
./gradlew assembleRelease
```

## BƯỚC 5: Lấy file chạy ứng dụng

- File cài đặt ứng dụng độc lập nằm tại đường dẫn: `📁 android/app/build/outputs/apk/release/app-release.apk`
- Bạn chỉ cần lấy đúng file `app-release.apk` này gửi lên Google Drive hoặc gửi qua Zalo/Telegram sang điện thoại thật là có thể bấm cài đặt và mở lên dùng bình thường.
- Kéo thả file `app-release.apk` đó vào cửa sổ máy ảo `genymotion` để cài đặt.

## Lưu ý bảo mật trong dự án với git

- File khóa ký ứng dụng (`my-upload-key.keystore`): Nếu kẻ xấu lấy được file này kèm mật khẩu, họ có thể giải nén app của bạn, chèn mã độc vào, ký lại bằng chính chữ ký của bạn rồi tải lên chợ ứng dụng dưới danh nghĩa là bạn.

- File chứa mật khẩu (`android/gradle.properties`): File này chứa mật khẩu dạng chữ rõ (Plain text) như `123456` mà bạn vừa điền. Kẻ xấu chỉ cần đọc file này là biết toàn bộ mật khẩu key của bạn.
- Cách khắc phục: Không đặt 4 dòng thông tin đó trong file `android/gradle.properties` và vẫn push lên git file đó như bình thường nhưng dùng 2 cách sau:

### Cách 1: Đặt mật khẩu ở file gradle.properties cục bộ (Khuyên dùng)

- Trên Windows: `C:\Users\Tên_Của_Bạn\.gradle\gradle.properties`
- Trên Linux/macOS: `~/.gradle/gradle.properties`
- Dán vào đó 4 dòng sau:

```bash
MYAPP_UPLOAD_STORE_FILE=/đường_dẫn_tuyệt_đối_đến_file/my-upload-key.keystore
MYAPP_UPLOAD_KEY_ALIAS=my-key-alias
MYAPP_UPLOAD_STORE_PASSWORD=123456
MYAPP_UPLOAD_KEY_PASSWORD=123456
```

- Khi bạn chạy lệnh `./gradlew assembleRelease`, Gradle sẽ tự động ưu tiên đọc file nằm trong máy bạn trước để lấy mật khẩu, còn file gradle.properties nằm trong dự án của bạn vẫn hoàn toàn "sạch sẽ" để push lên GitHub.

### Cách 2: Sử dụng biến môi trường (Environment Variables)

- Thay vì ghi chữ rõ mật khẩu, người ta cấu hình trong file `android/app/build.gradle` để nó đọc từ biến môi trường của hệ điều hành (thường áp dụng khi build tự động qua GitHub Actions, GitLab CI/CD).

# Hướng dẫn build React Native Expo

## BƯỚC 1: Khởi tạo dự án Expo từ đầu

```bash
npx create-expo-app@latest MyNewApp
cd MyNewApp
```

## BƯỚC 2: Cài đặt công cụ EAS CLI và Đăng ký tài khoản

- EAS CLI là công cụ giúp máy tính của bạn giao tiếp với máy chủ đám mây của Expo.

1.  Cài đặt EAS CLI toàn cục:

```Bash
npm install -g eas-cli
```

2. Đăng ký tài khoản: Hãy truy cập vào `expo.dev` và đăng ký một tài khoản miễn phí.

3. Đăng nhập trên máy tính: Quay lại Terminal và gõ lệnh sau để liên kết máy tính với tài khoản Expo:

```Bash
eas login
```

- Nhập chính xác tài khoản và mật khẩu Expo của bạn.

## BƯỚC 3: Cấu hình EAS cho dự án (Chỉ làm 1 lần đầu)

- Tại thư mục gốc của dự án, chạy lệnh khởi tạo:

```Bash
eas build:configure
```

- Hệ thống sẽ hỏi: `Which platforms would you like to configure for EAS Build?`

- Bạn dùng phím mũi tên chọn `All` (để cấu hình cho cả Android và iOS) rồi nhấn `Enter`.

- Lúc này, một file cấu hình tên là `eas.json` sẽ tự động sinh ra ở thư mục gốc dự án.

## BƯỚC 4: Chỉnh sửa cấu hình đầu ra (Tạo file .apk và .ipa)

- Mặc định, cấu hình của Expo sẽ build ứng dụng để đưa lên Store (Android ra `.aab`, iOS ra `.ipa` dạng App Store). Để build ra file cài trực tiếp lên máy thật để test, bạn cần chỉnh sửa file eas.json.

- Mở file `eas.json` lên và cập nhật lại block preview giống như sau:

```bash
{
  "cli": {
    "version": ">= 14.0.0"
  },
  "build": {
    "development": {
      "developmentClient": true,
      "distribution": "internal"
    },
    "preview": {
      "distribution": "internal",
      "android": {
        "buildType": "apk"
      },
      "ios": {
        "simulator": false
      }
    },
    "production": {}
  }
}
```

## BƯỚC 5: Tiến hành chạy lệnh Build

🤖 1. Hướng dẫn Build cho ANDROID (Ra file `.apk`)

- Chạy lệnh sau tại thư mục gốc:

```bash
eas build --platform android --profile preview
```

- Quá trình xử lý: Expo sẽ hỏi bạn có muốn họ tự động quản lý khóa ký (Keystore) không, hãy chọn Y (Yes). Dự án sẽ được nén và đẩy lên Cloud để build.

- Kết quả: Sau khi chạy xong (3-5 phút), Terminal sẽ hiển thị một Mã QR. Bạn chỉ cần lấy điện thoại Android quét mã này là tải được file .apk về cài đặt trực tiếp.

🍏 2. Hướng dẫn Build cho iOS (Ra file .ipa)

- Build iOS phức tạp hơn một chút vì Apple quản lý bảo mật rất nghiêm ngặt. Bạn bắt buộc phải có Tài khoản nhà phát triển Apple (Apple Developer Account) (loại có phí 99 USD/năm) thì mới build cài vào máy thật được.

- Chạy lệnh sau tại thư mục gốc:

```Bash
eas build --platform ios --profile preview
```

- Quá trình xử lý:

* EAS CLI sẽ yêu cầu bạn đăng nhập tài khoản Apple Developer của bạn ngay trên Terminal.

* Expo sẽ hỏi bạn có muốn họ tự động lo phần cấp chứng chỉ (Provisioning Profile và Signing Certificate) không? Bạn chọn Y (Yes).

* Để cài được vào máy thật, bạn cần đăng ký mã định danh của điện thoại (gọi là UDID) với tài khoản Apple. Expo cũng sẽ tự động hướng dẫn bạn quét một mã QR trên iPhone để tự động lấy và đăng ký mã UDID này lên hệ thống của Apple luôn! Mọi thứ diễn ra hoàn toàn tự động.

* Kết quả: Khi hoàn thành, bạn cũng sẽ nhận được một đường link/mã QR. Sử dụng chính chiếc iPhone đã đăng ký UDID quét mã để cài đặt app trực tiếp qua môi trường mạng (Ad-Hoc distribution).
