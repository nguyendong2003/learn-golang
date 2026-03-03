# Setup ubuntu tool

## Bộ gõ tiếng việt ibus-bamboo
- Buoc 1:
```bash
sudo add-apt-repository ppa:bamboo-engine/ibus-bamboo
sudo apt update
sudo apt install ibus-bamboo
```

- Buoc 2:
    + Run bash:
    ```bash
    ibus restart
    ```
    
    + Mở Settings (Cài đặt): Vào mục Keyboard (Bàn phím)
    + Input Sources: Nhấn dấu +, chọn Vietnamese, sau đó tìm và chọn Vietnamese (Bamboo).

    + Lưu ý: Nếu không thấy "Vietnamese (Bamboo)", hãy thử đăng xuất (Log out) rồi đăng nhập lại.

## Cai dat vscode
- Tải từ trang chủ: https://code.visualstudio.com/
- Chuột phải vào file và chọn: Open With App Center -> Install
- Lưu ý: Tải vscode từ trang chủ mới nhận ibus-bamboo chứ tải từ App Center không nhận


## Phan mem Play video - VLC
```bash
sudo apt update && sudo apt install vlc
```

## Phần mềm văn phòng LibreOffice
```bash
sudo apt update
sudo apt install libreoffice
```

## Git
- Cài đặt
```bash
sudo apt update
sudo apt install git
```

- Kiểm tra
```bash
git --version
```

- Cấu hình
```bash
git config --global user.name "Nguyen Dong"
git config --global user.email "dongcoi14122003@gmail.com"
```

- Kiểm tra cấu hình
```bash
git config --list
```

