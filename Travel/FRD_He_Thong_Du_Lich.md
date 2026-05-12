# TÀI LIỆU YÊU CẦU CHỨC NĂNG (FRD)
# HỆ THỐNG ĐẶT TOUR DU LỊCH & KHÁCH SẠN

| **Thông tin** | **Chi tiết** |
|---|---|
| Tên dự án | Hệ thống đặt Tour Du lịch & Khách sạn trực tuyến |
| Phạm vi | Việt Nam (nội địa) |
| Phiên bản | 1.0 |
| Ngày tạo | 12/05/2026 |
| Người soạn thảo | Business Analyst |
| Đối tượng đọc | Project Manager, Developer, Tester, Designer |
| Trạng thái | Approved |

---

## 📑 MỤC LỤC

1. [Tổng quan hệ thống](#1-tổng-quan-hệ-thống)
2. [Đối tượng sử dụng (Actors)](#2-đối-tượng-sử-dụng-actors)
3. [Phân quyền hệ thống](#3-phân-quyền-hệ-thống)
4. [Danh mục trạng thái](#4-danh-mục-trạng-thái-status-codes)
5. [Chức năng phía Admin](#5-chức-năng-phía-admin)
6. [Chức năng phía Đại lý](#6-chức-năng-phía-đại-lý)
7. [Chức năng phía Người dùng](#7-chức-năng-phía-người-dùng-khách-hàng)
8. [Chức năng nâng cao](#8-chức-năng-nâng-cao)
9. [Yêu cầu phi chức năng](#9-yêu-cầu-phi-chức-năng-tham-khảo)
10. [Phụ lục](#10-phụ-lục)

---

## 1. TỔNG QUAN HỆ THỐNG

### 1.1. Mục đích

Hệ thống là nền tảng trung gian (marketplace) kết nối **khách hàng** có nhu cầu du lịch trong nước với **các đại lý** (chủ khách sạn, công ty du lịch) cung cấp dịch vụ tại Việt Nam.

### 1.2. Phạm vi

| **Trong phạm vi** | **Ngoài phạm vi** |
|---|---|
| Tour, khách sạn, combo nội địa Việt Nam | Tour, KS quốc tế |
| Quản lý đặt chỗ, kiểm duyệt, đánh giá | Cổng thanh toán trực tuyến |
| Chatbot AI tư vấn | Hệ thống kế toán doanh nghiệp |
| Affiliate, Corporate Booking | HRM của đại lý |
| Google Maps tích hợp | Booking vé máy bay độc lập |
| VND, chặng bay nội địa | Đa tiền tệ, chặng bay quốc tế |

> **⚠️ Lưu ý quan trọng:** Hệ thống **KHÔNG** xử lý thanh toán trực tuyến. Việc thanh toán được thực hiện ngoài hệ thống (chuyển khoản trực tiếp cho đại lý, thanh toán tại quầy...). Hệ thống chỉ ghi nhận trạng thái booking và xác nhận từ đại lý.

### 1.3. Mô hình kinh doanh

```
┌─────────────┐      ┌──────────────────┐      ┌─────────────┐
│  Người dùng │ ───▶ │  Hệ thống (Web)  │ ◀─── │   Đại lý    │
│  (Khách)    │      │                  │      │ (KS / Tour) │
└─────────────┘      └────────┬─────────┘      └─────────────┘
                              │
                              ▼
                       ┌──────────────┐
                       │    Admin     │
                       │  (Vận hành)  │
                       └──────────────┘
```

### 1.4. Quy ước ký hiệu

| **Ký hiệu** | **Ý nghĩa** |
|---|---|
| FR-AD-xxx | Functional Requirement - Admin |
| FR-AG-xxx | Functional Requirement - Agent (Đại lý) |
| FR-US-xxx | Functional Requirement - User (Khách hàng) |
| FR-CO-xxx | Functional Requirement - Common (Chung) |
| BR-xx | Business Rule |
| ⭐ | Chức năng quan trọng / Phải có trong MVP |

---

## 2. ĐỐI TƯỢNG SỬ DỤNG (ACTORS)

| **Actor** | **Mô tả** | **Quyền chính** |
|---|---|---|
| **Khách vãng lai** | Người truy cập chưa đăng ký | Xem, tìm kiếm, đọc blog, chat AI |
| **Người dùng** | Khách hàng cá nhân đã đăng ký | Đặt tour/KS/combo, đánh giá, tích điểm |
| **Khách hàng doanh nghiệp** | Tổ chức đăng ký booking nhóm | Corporate booking, hợp đồng |
| **Cộng tác viên (Affiliate)** | Người chia sẻ link kiếm hoa hồng | Quản lý link, xem hoa hồng |
| **Đại lý - Chủ khách sạn** | Khách sạn / Homestay / Resort | Quản lý KS, phòng, giá |
| **Đại lý - Công ty du lịch** | Công ty tổ chức tour | Quản lý tour, lịch trình |
| **Admin** | Quản trị viên hệ thống (Super Admin) | Toàn quyền hệ thống |
| **Nhân viên hỗ trợ** | Staff vận hành (Admin tạo) | Quyền giới hạn theo phân công |

---

## 3. PHÂN QUYỀN HỆ THỐNG

Hệ thống áp dụng mô hình **RBAC (Role-Based Access Control)**.

| **Module** | **Khách vãng lai** | **User** | **Đại lý** | **Affiliate** | **Admin** |
|---|:---:|:---:|:---:|:---:|:---:|
| Xem tour/KS/Combo | ✅ | ✅ | ✅ | ✅ | ✅ |
| Đăng ký/Đăng nhập | ✅ | ✅ | ✅ | ✅ | ✅ |
| Đặt tour/KS/Combo | ❌ | ✅ | ✅ | ✅ | ✅ |
| Đánh giá | ❌ | ✅ | ❌ | ✅ | ✅ |
| Tạo tour/KS | ❌ | ❌ | ✅ | ❌ | ✅ |
| Quản lý sản phẩm | ❌ | ❌ | ✅ | ❌ | ✅ |
| Duyệt nội dung | ❌ | ❌ | ❌ | ❌ | ✅ |
| Quản lý user | ❌ | ❌ | ❌ | ❌ | ✅ |
| Cấu hình hệ thống | ❌ | ❌ | ❌ | ❌ | ✅ |
| Xem hoa hồng | ❌ | ❌ | ✅ | ✅ | ✅ |
| Tạo link affiliate | ❌ | ✅¹ | ✅¹ | ✅ | ✅ |
| Chat với AI | ✅ | ✅ | ✅ | ✅ | ✅ |

¹ *User và Đại lý có thể đăng ký thành Affiliate để có quyền tạo link.*

---

## 4. DANH MỤC TRẠNG THÁI (STATUS CODES)

### 4.1. Trạng thái Booking

| **Mã** | **Tên** | **Mô tả** |
|---|---|---|
| `PENDING_CONFIRMATION` | Chờ xác nhận | Khách vừa đặt, chờ đại lý xác nhận |
| `CONFIRMED` | Đã xác nhận | Đại lý đã xác nhận, chờ khách thanh toán |
| `WAITING_PAYMENT` | Chờ thanh toán | Khách đã được hướng dẫn thanh toán ngoài hệ thống |
| `PAID` | Đã thanh toán | Đại lý đã xác nhận nhận được tiền |
| `IN_PROGRESS` | Đang sử dụng | Khách đang sử dụng dịch vụ (đã check-in / đang trong tour) |
| `COMPLETED` | Hoàn tất | Khách đã sử dụng xong dịch vụ |
| `CANCELLED_BY_USER` | Khách hủy | User chủ động hủy |
| `CANCELLED_BY_AGENT` | Đại lý hủy | Đại lý từ chối / hủy |
| `CANCELLED_BY_SYSTEM` | Hệ thống hủy | Hết hạn xác nhận / hết hạn thanh toán |
| `REFUND_REQUESTED` | Yêu cầu hoàn tiền | Khách gửi yêu cầu hoàn tiền |
| `REFUNDED` | Đã hoàn tiền | Đại lý đã hoàn tiền cho khách |

### 4.2. Trạng thái Sản phẩm (Tour/KS/Combo)

| **Mã** | **Tên** | **Mô tả** |
|---|---|---|
| `DRAFT` | Bản nháp | Đại lý đang soạn thảo |
| `PENDING_REVIEW` | Chờ duyệt | Đã gửi admin duyệt |
| `APPROVED` | Đã duyệt | Admin duyệt nhưng đại lý chưa xuất bản |
| `PUBLISHED` | Đang hiển thị | Hiển thị công khai trên website |
| `REJECTED` | Bị từ chối | Admin từ chối, đại lý cần sửa |
| `PAUSED` | Tạm ngưng | Đại lý/Admin tạm ẩn |
| `ARCHIVED` | Đã lưu trữ | Không còn kinh doanh |

### 4.3. Trạng thái Đại lý

| **Mã** | **Tên** | **Mô tả** |
|---|---|---|
| `PENDING` | Chờ duyệt | Mới đăng ký, chờ admin xác minh |
| `ACTIVE` | Đang hoạt động | Đại lý đã được duyệt |
| `SUSPENDED` | Tạm khóa | Vi phạm, tạm khóa có thời hạn |
| `BANNED` | Bị cấm | Vi phạm nghiêm trọng, khóa vĩnh viễn |

---

## 5. CHỨC NĂNG PHÍA ADMIN

### 5.1. Module Quản lý tài khoản & Phân quyền

#### ⭐ FR-AD-001: Quản lý người dùng (Khách hàng)

**Mô tả:** Admin xem, tìm kiếm và quản lý toàn bộ tài khoản khách hàng.

**Actor:** Admin, Nhân viên hỗ trợ (chỉ quyền xem)

**Flow chính:**
1. Admin truy cập menu **"Quản lý người dùng"**.
2. Hệ thống hiển thị danh sách user dạng bảng, phân trang (mặc định 20 user/trang).
3. Admin có thể:
   - **Tìm kiếm** theo: họ tên, email, SĐT, ngày đăng ký.
   - **Lọc** theo: trạng thái (Active/Locked/Deleted), hạng thành viên, có booking hay không.
   - **Sắp xếp** theo: ngày đăng ký, tổng chi tiêu, số booking.
4. Click vào 1 user → xem trang chi tiết:
   - Thông tin cá nhân
   - Lịch sử booking (tour/KS/combo)
   - Lịch sử đánh giá
   - Điểm thưởng & hạng thành viên
   - Lịch sử voucher đã dùng
   - Log đăng nhập (IP, thiết bị)
5. Hành động: Khóa/Mở khóa, Đặt lại mật khẩu, Xóa mềm, Cộng/Trừ điểm thưởng, Thêm ghi chú nội bộ.

**Business Rule:**
- **BR-01:** Không cho phép xóa cứng tài khoản đã có booking. Chỉ xóa mềm (soft delete).
- **BR-02:** Khi khóa user, các booking đang xử lý vẫn được giữ nguyên.
- **BR-03:** Lý do khóa phải được ghi log và gửi email thông báo cho user.
- **BR-04:** User bị khóa không thể đăng nhập, hệ thống hiển thị thông báo lý do.

---

#### ⭐ FR-AD-002: Quản lý đại lý

**Mô tả:** Admin duyệt đăng ký, quản lý và xác minh các đại lý.

**Actor:** Admin

**Flow chính:**
1. Admin xem 2 tab:
   - **Chờ duyệt** (`PENDING`)
   - **Đã duyệt** (`ACTIVE`, `SUSPENDED`, `BANNED`)
2. Với đại lý chờ duyệt:
   - Xem hồ sơ: tên công ty, MST, GPKD (file), CCCD người đại diện, địa chỉ.
   - Có thể yêu cầu bổ sung tài liệu.
   - **Duyệt** → `ACTIVE`, gửi email hướng dẫn.
   - **Từ chối** (kèm lý do) → gửi email.
3. Với đại lý đã duyệt:
   - Xem sản phẩm, doanh thu, hoa hồng, đánh giá.
   - **Tạm khóa** khi vi phạm.
   - **Cấm vĩnh viễn** khi vi phạm nghiêm trọng.
   - **Điều chỉnh tỷ lệ hoa hồng** riêng cho đại lý.

**Business Rule:**
- **BR-05:** Đại lý chỉ có thể đăng tour/KS sau khi được duyệt.
- **BR-06:** Mỗi đại lý có 1 mã định danh duy nhất (format: `AG-YYYYMM-XXXXX`).
- **BR-07:** Đại lý bị `SUSPENDED`/`BANNED` → tất cả sản phẩm tự động `PAUSED`.
- **BR-08:** MST phải là MST Việt Nam hợp lệ (10-13 chữ số).
- **BR-09:** 1 MST chỉ đăng ký được 1 tài khoản đại lý.

---

#### FR-AD-003: Quản lý nhân viên admin & phân quyền (RBAC)

**Mô tả:** Tạo tài khoản nhân viên hỗ trợ với quyền giới hạn.

**Actor:** Super Admin

**Flow chính:**
1. Truy cập **"Quản lý nhân viên"** → **"Tạo mới"**.
2. Nhập: họ tên, email, SĐT, vai trò.
3. Chọn nhóm quyền:
   - `Content Moderator`: Duyệt nội dung, blog, review
   - `Customer Support`: Xử lý ticket, hỗ trợ khách
   - `Finance`: Báo cáo tài chính, đối soát
   - `Marketing`: Voucher, campaign
   - `Super Admin`: Toàn quyền
4. Gửi email kích hoạt (hiệu lực 24h).

**Quản lý Roles:**
- Tạo/sửa/xóa Role.
- Mỗi Role gồm tập Permissions (`user.view`, `user.edit`, `booking.cancel`, `agent.approve`...).
- Có thể tạo Role tùy chỉnh.

**Business Rule:**
- **BR-10:** Chỉ Super Admin được tạo/sửa Roles.
- **BR-11:** Mọi hành động của nhân viên đều ghi log audit: ai, làm gì, thời điểm, IP.
- **BR-12:** Không thể xóa Super Admin cuối cùng của hệ thống.

---

#### FR-AD-004: Quản lý cộng tác viên (Affiliate)

**Mô tả:** Quản lý CTV đăng ký chương trình affiliate.

**Actor:** Admin

**Flow chính:**
1. Xem danh sách CTV: tên, kênh quảng bá, click, conversion, hoa hồng.
2. Duyệt đơn đăng ký mới.
3. Cấu hình tỷ lệ hoa hồng affiliate (mặc định hoặc theo từng CTV).
4. Theo dõi hiệu suất.
5. Khóa CTV nếu phát hiện gian lận.

**Business Rule:**
- **BR-13:** User và Đại lý đều có thể đăng ký làm Affiliate.
- **BR-14:** Hoa hồng affiliate chỉ tính khi booking `COMPLETED`.
- **BR-15:** Phát hiện self-purchase → không tính hoa hồng.

---

#### FR-AD-005: Quản lý tài khoản doanh nghiệp (Corporate)

**Mô tả:** Quản lý khách hàng doanh nghiệp ký hợp đồng dài hạn.

**Actor:** Admin, Customer Support

**Flow chính:**
1. Xem danh sách tài khoản doanh nghiệp.
2. Duyệt đăng ký mới (kiểm tra MST, GPKD).
3. Cấu hình hợp đồng:
   - Mức chiết khấu riêng (% trên giá niêm yết)
   - Hạn mức tín dụng
   - Thời hạn thanh toán (NET 15/30/60 ngày)
   - Người đại diện duyệt booking
4. Xem lịch sử booking, công nợ.

**Business Rule:**
- **BR-16:** Booking doanh nghiệp có giá riêng, không cộng dồn voucher công khai.
- **BR-17:** Hợp đồng có thời hạn, hết hạn → cần gia hạn.

---

### 5.2. Module Quản lý nội dung & Danh mục

#### ⭐ FR-AD-006: Quản lý danh mục

**Mô tả:** Quản lý các danh mục dùng chung trên hệ thống.

**Actor:** Admin, Content Moderator

**Các loại danh mục:**

| **Danh mục** | **Ví dụ giá trị** |
|---|---|
| Loại tour | Tour ngày, Tour dài ngày, Tour mạo hiểm, Tour nghỉ dưỡng, Tour văn hóa, Tour ẩm thực, Tour biển đảo, Tour miền núi, Tour caravan |
| Loại khách sạn | Hotel, Resort, Homestay, Villa, Hostel, Căn hộ dịch vụ, Bungalow |
| Hạng sao | 1 sao, 2 sao, 3 sao, 4 sao, 5 sao |
| Loại phòng | Standard, Deluxe, Superior, Suite, Family, Connecting |
| Tiện nghi KS | Wifi, Bể bơi, Phòng gym, Spa, Nhà hàng, Bãi đỗ xe, Ăn sáng, Đưa đón sân bay, Cho phép thú cưng |
| Tiện nghi phòng | Smart TV, Minibar, Ban công, Bồn tắm, Ấm đun nước, Két sắt, Điều hòa, Máy sấy tóc |
| Phương tiện tour | Xe khách, Xe limousine, Tàu hỏa, Máy bay, Tàu thủy, Xe máy |
| Loại bữa ăn | Sáng, Trưa, Tối, Buffet, Set menu |

**Flow chính:**
1. Chọn loại danh mục.
2. Thêm/Sửa/Xóa/Sắp xếp thứ tự.
3. Mỗi mục có: tên, mô tả, icon/ảnh, trạng thái Active/Inactive.

**Business Rule:**
- **BR-18:** Không cho phép xóa danh mục đang được sử dụng. Phải chuyển sang Inactive.

---

#### ⭐ FR-AD-007: Quản lý điểm đến (Destination) - Việt Nam

**Mô tả:** Quản lý cây địa lý Việt Nam phục vụ tìm kiếm tour/KS.

**Actor:** Admin

**Cấu trúc 3 cấp:**
```
Tỉnh/Thành phố
    └── Phường/Xã
            └── Khu vực / Điểm đến cụ thể
```

**Ví dụ:**
- `Hà Nội → Phường Hoàn Kiếm → Phố cổ Hà Nội`
- `Lào Cai → Xã Sa Pa → Bản Cát Cát`
- `Khánh Hòa → Phường Lộc Thọ → Bãi biển Trần Phú`
- `Kiên Giang → Xã Dương Đông → Bãi Sao - Phú Quốc`

**Flow chính:**
1. Quản lý theo cây phân cấp (tree view).
2. Thêm/Sửa/Xóa từng cấp.
3. Mỗi điểm đến có:
   - Tên (chính thức theo hành chính)
   - Mã hành chính
   - Mô tả
   - Hình ảnh đại diện
   - Tọa độ GPS (lat, lng)
   - SEO meta
   - Trạng thái Active/Inactive
   - Đánh dấu **"Điểm đến nổi bật"** (hiển thị trang chủ)
4. Hỗ trợ import danh sách Tỉnh/Phường/Xã từ file Excel.

**Business Rule:**
- **BR-19:** Dữ liệu Tỉnh/TP, Phường/Xã đồng bộ với danh mục hành chính chuẩn của Việt Nam.
- **BR-20:** Khi đổi tên/sáp nhập đơn vị hành chính, hệ thống có chức năng cập nhật & redirect URL cũ → mới (SEO).
- **BR-21:** Không xóa điểm đến đang có tour/KS gắn vào.

---

#### FR-AD-008: Quản lý trang tĩnh

**Mô tả:** Quản lý các trang nội dung tĩnh trên website.

**Actor:** Admin, Content Moderator

**Danh sách trang tĩnh tối thiểu:**
- Giới thiệu / Về chúng tôi
- Điều khoản sử dụng
- Chính sách bảo mật
- Chính sách hoàn/hủy
- Hướng dẫn đặt tour
- Hướng dẫn đặt khách sạn
- Câu hỏi thường gặp (FAQ)
- Liên hệ
- Tuyển dụng
- Hợp tác đại lý (landing page)
- Chương trình affiliate

**Flow chính:**
1. Chọn trang cần chỉnh sửa.
2. Sử dụng WYSIWYG editor:
   - Định dạng văn bản (heading, bold, italic, list)
   - Chèn ảnh, video, link, bảng
   - Embed HTML
3. Cấu hình SEO: meta title, description, keywords, OG tags.
4. Lưu nháp / Xuất bản / Lên lịch.

**Business Rule:**
- **BR-22:** Lưu lịch sử thay đổi (version history) - có thể rollback.
- **BR-23:** URL slug có thể tùy chỉnh, không trùng nhau.

---

#### FR-AD-009: Quản lý banner & slider

**Mô tả:** Quản lý banner/slider hiển thị trên website.

**Actor:** Admin, Marketing

**Vị trí banner:**
- Hero slider trang chủ (1920x600 px)
- Sidebar trang tìm kiếm (300x250 px)
- Popup khuyến mãi (600x400 px)
- Banner trang tour/KS
- Footer

**Flow chính:**
1. Tạo banner:
   - Upload ảnh, tiêu đề, mô tả, CTA button
   - Link đích
   - Vị trí, thứ tự hiển thị
   - Thời gian hiệu lực
   - Đối tượng hiển thị (tất cả/đã đăng nhập/hạng thành viên)
2. Bật/Tắt banner.
3. Xem thống kê: impression, click, CTR.

**Business Rule:**
- **BR-24:** Banner hết hạn tự động ẩn.
- **BR-25:** Kiểm tra kích thước ảnh đúng chuẩn cho từng vị trí.

---

#### FR-AD-010: Quản lý blog/tin tức du lịch

**Mô tả:** Quản lý bài viết blog, cẩm nang du lịch.

**Actor:** Admin, Content Moderator

**Flow chính:**
1. Tạo bài viết:
   - Tiêu đề, slug
   - Tóm tắt, nội dung (rich text)
   - Ảnh đại diện
   - Danh mục, tags, tác giả
   - SEO meta
   - Gắn tour/KS liên quan
2. Lên lịch hoặc xuất bản ngay.
3. Quản lý bình luận: duyệt / xóa / báo cáo spam.

**Business Rule:**
- **BR-26:** Bài viết có thể gắn tối đa 5 tour/KS để hiển thị "Sản phẩm liên quan".

---

### 5.3. Module Kiểm duyệt

#### ⭐ FR-AD-011: Duyệt tour, khách sạn, combo

**Mô tả:** Admin kiểm duyệt sản phẩm do đại lý gửi trước khi công khai.

**Actor:** Admin, Content Moderator

**Flow chính:**
1. Truy cập **"Hàng chờ duyệt"** → 3 tab: Tour / Khách sạn / Combo.
2. Mỗi item: tên, đại lý, ngày gửi, độ ưu tiên.
3. Click vào item → xem chi tiết đầy đủ.
4. Chọn hành động:
   - **Duyệt** → `APPROVED`, đại lý có thể publish.
   - **Yêu cầu chỉnh sửa** → gửi feedback, sản phẩm về `DRAFT`.
   - **Từ chối** → `REJECTED`, kèm lý do.
5. Gửi email + in-app notification cho đại lý.

**Tiêu chí duyệt (Checklist):**
- ✅ Thông tin đầy đủ, không sai chính tả
- ✅ Hình ảnh chất lượng, không vi phạm bản quyền
- ✅ Giá hợp lý
- ✅ Chính sách hủy rõ ràng
- ✅ Nội dung không vi phạm pháp luật, thuần phong mỹ tục
- ✅ Tọa độ GPS chính xác

**Business Rule:**
- **BR-27:** SLA duyệt tối đa: 48 giờ làm việc.
- **BR-28:** Chỉnh sửa lớn (giá, lịch trình, chính sách) → duyệt lại. Chỉnh sửa nhỏ → không cần.
- **BR-29:** Bị từ chối 3 lần liên tiếp → cảnh báo đại lý.

---

#### FR-AD-012: Duyệt đánh giá (Review)

**Mô tả:** Kiểm duyệt review trước khi công khai.

**Actor:** Admin, Content Moderator

**Flow chính:**
1. Xem danh sách review chờ duyệt.
2. Kiểm tra: spam, ngôn từ tục, quảng cáo, link ngoài, SĐT/email, ảnh nhạy cảm.
3. **Duyệt** / **Từ chối** / **Ẩn**.

**Business Rule:**
- **BR-30:** Review chỉ viết được bởi user đã `COMPLETED` booking.
- **BR-31:** 1 user / 1 booking / 1 review. Có thể sửa trong 7 ngày.
- **BR-32:** Review chứa SĐT/email/link → tự đưa vào hàng chờ duyệt.
- **BR-33:** Có thể bật/tắt chế độ "Auto-approve" (kiểm duyệt hậu kỳ).

---

#### FR-AD-013: Duyệt khuyến mãi của đại lý

**Mô tả:** Duyệt khuyến mãi vượt ngưỡng.

**Business Rule:**
- **BR-34:** Discount ≤ 50% → tự động duyệt.
- **BR-35:** Discount > 50% → bắt buộc duyệt thủ công, đại lý phải nhập lý do.

---

### 5.4. Module Quản lý đặt chỗ

#### ⭐ FR-AD-014: Quản lý booking toàn hệ thống

**Mô tả:** Admin xem & can thiệp toàn bộ booking.

**Actor:** Admin, Customer Support

**Flow chính:**
1. Truy cập **"Quản lý booking"**.
2. Danh sách: Mã booking, Khách hàng, Sản phẩm, Đại lý, Ngày đặt, Ngày sử dụng, Tổng tiền, Trạng thái.
3. Bộ lọc: trạng thái, loại sản phẩm, đại lý, thời gian, khoảng giá.
4. Click chi tiết: thông tin khách, sản phẩm, hành khách, timeline trạng thái, ghi chú, lịch sử trao đổi.
5. Hành động:
   - Chuyển trạng thái thủ công (kèm lý do)
   - Thêm ghi chú nội bộ
   - Gửi email/SMS lại
   - Liên hệ đại lý

**Business Rule:**
- **BR-36:** Mỗi thay đổi trạng thái → log + thông báo các bên.
- **BR-37:** Admin can thiệp phải ghi rõ lý do (mandatory).

---

#### FR-AD-015: Xử lý tranh chấp & khiếu nại

**Mô tả:** Xử lý tranh chấp giữa user và đại lý.

**Actor:** Admin, Customer Support

**Flow chính:**
1. User/Đại lý gửi khiếu nại từ chi tiết booking.
2. Hệ thống tạo ticket, phân loại (thấp/trung bình/cao).
3. Admin xem nội dung, chứng cứ (ảnh, file, tin nhắn).
4. Trao đổi với cả 2 bên qua chat khiếu nại.
5. Quyết định:
   - Hoàn tiền toàn bộ / một phần
   - Đền bù voucher / điểm thưởng
   - Từ chối khiếu nại
   - Cảnh cáo / phạt đại lý
6. Đóng ticket khi 2 bên xác nhận.

**Business Rule:**
- **BR-38:** Mọi tranh chấp có ticket ID, lưu vĩnh viễn.
- **BR-39:** SLA xử lý: 72 giờ làm việc.
- **BR-40:** Khiếu nại an toàn/khẩn cấp → ưu tiên 24h.

---

### 5.5. Module Quản lý tài chính

> ⚠️ Hệ thống KHÔNG xử lý thanh toán trực tuyến. Chỉ **ghi nhận** và **đối soát**.

#### ⭐ FR-AD-016: Quản lý hoa hồng đại lý

**Mô tả:** Cấu hình và theo dõi hoa hồng phải thu từ đại lý.

**Actor:** Admin, Finance

**Flow chính:**
1. **Cấu hình tỷ lệ:**
   - Mặc định toàn hệ thống (VD: 10%)
   - Override theo loại sản phẩm (tour 12%, KS 10%, combo 15%)
   - Override theo từng đại lý
2. **Tính tự động:** Khi booking `COMPLETED` → hoa hồng = Tổng × % hoa hồng.
3. **Đối soát định kỳ:**
   - Chốt kỳ (tuần/tháng)
   - Xuất bảng đối soát chi tiết
   - Gửi đại lý xác nhận
4. **Ghi nhận thanh toán:** Đại lý chuyển khoản → admin đánh dấu "Đã thu" + ngày + số chứng từ.

**Business Rule:**
- **BR-41:** Hoa hồng tính tại thời điểm `COMPLETED`.
- **BR-42:** Booking bị hoàn tiền → rollback hoa hồng.
- **BR-43:** Quá hạn 30 ngày chưa thanh toán → cảnh báo, có thể tạm khóa đại lý.

---

#### FR-AD-017: Quản lý hoa hồng cộng tác viên (Affiliate)

**Mô tả:** Theo dõi và chi trả hoa hồng cho CTV.

**Actor:** Admin, Finance

**Flow chính:**
1. Track: link → click → booking → conversion.
2. Booking `COMPLETED` → hoa hồng CTV = Tổng × % hoa hồng affiliate.
3. CTV gửi yêu cầu rút (đạt ngưỡng tối thiểu, VD: 500.000đ).
4. Admin duyệt, chuyển khoản (ngoài hệ thống), ghi nhận "Đã chi".

**Business Rule:**
- **BR-44:** Hoa hồng affiliate tính sau `COMPLETED` + 7 ngày (qua thời hạn hoàn hủy).
- **BR-45:** Phát hiện gian lận → hủy hoa hồng + cảnh cáo/khóa CTV.

---

#### FR-AD-018: Báo cáo dòng tiền

**Output:** Bảng tổng hợp + biểu đồ + xuất Excel/PDF.

**Chỉ số chính:**
- Tổng GMV
- Hoa hồng dự kiến (chưa COMPLETED)
- Hoa hồng đã ghi nhận (COMPLETED, chưa thu)
- Hoa hồng đã thu
- Hoa hồng affiliate đã chi
- Doanh thu thuần

---

### 5.6. Module Khuyến mãi & Marketing

#### ⭐ FR-AD-019: Quản lý voucher / mã giảm giá

**Actor:** Admin, Marketing

**Thuộc tính voucher:**

| **Trường** | **Mô tả** |
|---|---|
| Mã code | Tự sinh hoặc nhập (unique) |
| Tên chương trình | VD: "Sale 30/4 - Giảm 20%" |
| Loại giảm | % / Số tiền cố định / Free dịch vụ |
| Giá trị giảm | VD: 20% hoặc 200.000đ |
| Giảm tối đa | VD: 500.000đ (cho loại %) |
| Giá trị tối thiểu | Booking ≥ X mới áp dụng |
| Loại sản phẩm | Tour / KS / Combo / Tất cả |
| Điểm đến áp dụng | Toàn quốc / Tỉnh cụ thể |
| Đối tượng | Tất cả / User mới / Hạng thành viên / Tài khoản cụ thể |
| Tổng số lượng | VD: 1000 voucher |
| Giới hạn / user | VD: 1 lần / user |
| Thời gian hiệu lực | Từ - đến |
| Trạng thái | Active / Inactive |

**Flow chính:**
1. Tạo voucher với thuộc tính trên.
2. Voucher hiển thị tại: trang Khuyến mãi, chi tiết tour/KS, email marketing.
3. User nhập code khi đặt → validate → áp dụng.

**Business Rule:**
- **BR-46:** 1 booking chỉ áp dụng 1 voucher hệ thống.
- **BR-47:** Voucher hệ thống có thể cộng dồn voucher đại lý nếu đại lý cho phép.
- **BR-48:** Voucher hết hạn/hết số lượng → tự ẩn.

---

#### FR-AD-020: Chương trình tích điểm (Loyalty)

**Cấu hình:**
- **Tỷ lệ tích điểm:** 10.000đ chi tiêu = 1 điểm
- **Tỷ lệ đổi điểm:** 100 điểm = 10.000đ giảm giá
- **Hạng thành viên:**

| **Hạng** | **Điều kiện (chi tiêu/năm)** | **Quyền lợi** |
|---|---|---|
| Đồng | < 5 triệu | Tích điểm cơ bản |
| Bạc | 5-15 triệu | Tích điểm x1.2, sale sớm |
| Vàng | 15-50 triệu | Tích điểm x1.5, hỗ trợ ưu tiên |
| Bạch kim | > 50 triệu | Tích điểm x2, hỗ trợ VIP, quà sinh nhật |

**Business Rule:**
- **BR-49:** Điểm tích sau khi booking `COMPLETED`.
- **BR-50:** Điểm có hạn 12 tháng từ ngày tích.
- **BR-51:** Hạng cập nhật đầu tháng dựa trên chi tiêu 12 tháng gần nhất.

---

### 5.7. Module Báo cáo & Dashboard

#### ⭐ FR-AD-021: Dashboard tổng quan

**Widget hiển thị:**
- Tổng booking (hôm nay / 7 ngày / 30 ngày) - kèm % so với kỳ trước
- Tổng GMV
- User mới, đại lý mới
- Top 5 tour, top 5 khách sạn bán chạy
- Top 10 đại lý theo doanh thu
- Đánh giá trung bình toàn hệ thống
- Tickets chưa xử lý
- Sản phẩm chờ duyệt
- Biểu đồ booking theo thời gian (line chart)
- Biểu đồ phân bổ booking theo điểm đến (heatmap VN)

---

#### FR-AD-022: Báo cáo chi tiết

| **Báo cáo** | **Mục đích** |
|---|---|
| Doanh thu theo thời gian | Xu hướng |
| Doanh thu theo đại lý | Đánh giá hiệu suất |
| Doanh thu theo điểm đến | Phân tích điểm đến hot |
| Doanh thu theo loại sản phẩm | Tour/KS/Combo |
| Booking analytics | Tỷ lệ hủy, confirm, lead time |
| User analytics | New vs Returning, cohort, retention |
| Affiliate report | Hiệu suất CTV |
| Voucher report | Hiệu quả voucher |
| Campaign report | ROI campaign |
| Review report | Tỷ lệ review, điểm trung bình |

**Tính năng:** Lọc đa chiều, xuất Excel/PDF/CSV, lên lịch gửi email định kỳ.

---

### 5.8. Module Cấu hình hệ thống

#### FR-AD-023: Cấu hình chung

- Thông tin website: tên, logo, favicon, hotline, email
- Email SMTP
- SMS gateway
- Mạng xã hội: Facebook, Zalo, Instagram, YouTube
- Tracking: Google Analytics, Facebook Pixel, GTM

---

#### ⭐ FR-AD-024: Cấu hình Google Maps

- Google Maps API key
- Default zoom level, default center (trung tâm Việt Nam)
- Style map (default / custom JSON)
- Bật/Tắt:
  - Bản đồ chi tiết KS
  - Bản đồ chi tiết tour
  - Bản đồ kết quả tìm kiếm
  - Chỉ đường (Directions API)
  - Street View
  - Place Autocomplete

---

#### ⭐ FR-AD-025: Cấu hình Chatbot AI

**Cấu hình:**
- AI provider & API key
- Model
- System prompt
- Knowledge base: FAQ, dữ liệu tour/KS, chính sách
- Giới hạn: số câu hỏi/user/ngày, token max/response
- Fallback:
  - AI không trả lời được → chuyển human agent
  - Câu trả lời mặc định khi lỗi
- Tone: thân thiện / chuyên nghiệp / hài hước

**Business Rule:**
- **BR-52:** Chatbot không tự ý cam kết giá hoặc booking. Phải redirect đến trang booking.
- **BR-53:** Lưu lịch sử hội thoại tối thiểu 30 ngày để phân tích & cải thiện.

---

#### FR-AD-026: Cấu hình hoa hồng & phí

- % hoa hồng mặc định (theo loại sản phẩm)
- % hoa hồng affiliate mặc định
- % VAT
- Phí dịch vụ
- Ngưỡng rút hoa hồng affiliate

---

#### FR-AD-027: Cấu hình email/SMS/notification template

**Template chính:**

| **Loại** | **Trigger** |
|---|---|
| Email/SMS xác nhận đăng ký | User đăng ký |
| Email/SMS OTP | Đăng nhập / xác thực |
| Email xác nhận booking | Booking thành công |
| Email confirm từ đại lý | Đại lý xác nhận |
| Email hướng dẫn thanh toán | Sau confirm |
| Email nhắc lịch | 1 ngày trước ngày sử dụng |
| Email xin review | Sau `COMPLETED` |
| Email marketing | Campaign |

**Tính năng:** WYSIWYG editor, variable placeholders (`{{user_name}}`, `{{booking_code}}`...), preview, test gửi.

---

### 5.9. Module Hỗ trợ khách hàng

#### FR-AD-028: Quản lý ticket hỗ trợ

**Actor:** Admin, Customer Support

**Flow chính:**
1. Ticket tạo từ: form liên hệ, chi tiết booking, live chat (AI→human), email.
2. Hệ thống tự động: phân loại theo từ khóa, gán nhân viên, đặt ưu tiên.
3. Nhân viên xử lý: trả lời, đóng/mở, chuyển, escalate.
4. User đánh giá chất lượng (CSAT 1-5 sao).

**Business Rule:**
- **BR-54:** SLA phản hồi: 1h giờ hành chính, 4h ngoài giờ.
- **BR-55:** Ticket Urgent (an toàn, khẩn cấp) → SLA 15 phút.

---

#### ⭐ FR-AD-029: Giám sát Chatbot AI

**Chức năng:**
- Dashboard chatbot: số phiên/ngày, tỷ lệ resolved, tỷ lệ chuyển human, điểm hài lòng TB.
- Xem lịch sử hội thoại (filter theo trạng thái, ngày).
- Phân tích câu hỏi không trả lời được → bổ sung knowledge base.
- Thiết lập điều kiện chuyển human:
  - User gõ "gặp nhân viên"
  - AI confidence < threshold
  - Câu hỏi về khiếu nại, hoàn tiền

---


## 6. CHỨC NĂNG PHÍA ĐẠI LÝ

### 6.1. Module Đăng ký & Hồ sơ

#### ⭐ FR-AG-001: Đăng ký tài khoản đại lý

**Mô tả:** Cá nhân/doanh nghiệp Việt Nam đăng ký trở thành đại lý.

**Actor:** Người muốn trở thành đại lý

**Flow chính:**
1. Truy cập **"Trở thành đại lý"** từ trang chủ.
2. Chọn loại đại lý:
   - **Chủ khách sạn** (KS / Homestay / Resort)
   - **Công ty du lịch** (tổ chức tour)
   - **Cả hai**
3. Nhập thông tin:
   - Tên công ty / tên hộ kinh doanh
   - Mã số thuế (validate định dạng MST Việt Nam)
   - Giấy phép kinh doanh (upload file PDF/ảnh, max 10MB)
   - Tên, CCCD, SĐT người đại diện
   - Email công ty (xác thực OTP)
   - Địa chỉ trụ sở (chọn Tỉnh/TP → Phường/Xã + nhập số nhà, đường)
   - Mô tả ngắn (200-500 ký tự)
   - Lĩnh vực hoạt động (chọn từ danh mục)
4. Đồng ý điều khoản hợp tác (bắt buộc tích).
5. Gửi đăng ký → trạng thái `PENDING`.
6. Admin duyệt → email thông báo + link đăng nhập trang quản trị.

**Business Rule:**
- **BR-56:** 1 MST = 1 tài khoản đại lý.
- **BR-57:** Bắt buộc xác thực email + SĐT trước khi gửi đăng ký.
- **BR-58:** Hồ sơ chưa duyệt có thể chỉnh sửa, bổ sung tài liệu.
- **BR-59:** Đại lý phải có MST Việt Nam (không nhận đại lý nước ngoài).

---

#### FR-AG-002: Quản lý hồ sơ đại lý

**Mô tả:** Đại lý cập nhật thông tin công ty.

**Trường thông tin:**
- Logo (vuông, min 500x500)
- Ảnh bìa (1920x400)
- Gallery (tối đa 10 ảnh)
- Giới thiệu công ty (rich text)
- Lĩnh vực hoạt động
- Năm thành lập
- Quy mô (số nhân viên)
- Website
- Mạng xã hội (Facebook, Zalo, Instagram)
- Hotline, email hỗ trợ
- Giờ làm việc

**Business Rule:**
- **BR-60:** Thay đổi thông tin pháp lý (MST, GPKD, người đại diện) phải admin duyệt lại.
- **BR-61:** Thay đổi thông tin giới thiệu (mô tả, logo, ảnh) tự động cập nhật, không cần duyệt.

---

### 6.2. Module Quản lý khách sạn

#### ⭐ FR-AG-003: Quản lý khách sạn

**Mô tả:** Tạo và quản lý thông tin khách sạn.

**Actor:** Đại lý - Chủ khách sạn

**Trường dữ liệu khách sạn:**

| **Trường** | **Bắt buộc** | **Mô tả** |
|---|:---:|---|
| Tên khách sạn | ✅ | Tên chính thức |
| Loại hình | ✅ | Hotel / Resort / Homestay / Villa / Hostel... |
| Hạng sao | ✅ | 1-5 sao (nếu là Hotel/Resort) |
| Địa chỉ | ✅ | Số nhà, đường + Tỉnh/TP → Phường/Xã (dropdown) |
| Tọa độ GPS | ✅ | Google Maps picker (kéo ghim trên bản đồ) |
| Mô tả tổng quan | ✅ | Min 200 ký tự, rich text |
| Hình ảnh | ✅ | Tối thiểu 5 ảnh, max 30 ảnh, mỗi ảnh max 5MB |
| Video giới thiệu | ❌ | URL YouTube |
| Tiện nghi | ✅ | Multi-select từ danh mục |
| Giờ check-in/check-out | ✅ | VD: 14:00 / 12:00 |
| Chính sách hủy | ✅ | Chọn từ template hoặc tự nhập |
| Chính sách trẻ em | ❌ | Quy định miễn phí/phụ thu theo độ tuổi |
| Chính sách thú cưng | ❌ | Có/Không cho phép |
| Ngôn ngữ phục vụ | ❌ | Tiếng Việt, Anh, Trung... |

**Flow chính:**
1. Đại lý → **"Khách sạn của tôi"** → **Thêm mới**.
2. Wizard 4 bước:
   - Bước 1: Thông tin cơ bản
   - Bước 2: Vị trí (Google Maps)
   - Bước 3: Tiện nghi & chính sách
   - Bước 4: Hình ảnh
3. **Lưu nháp** (`DRAFT`) hoặc **Gửi duyệt** (`PENDING_REVIEW`).
4. Admin duyệt → `APPROVED` → đại lý publish → `PUBLISHED`.

**Business Rule:**
- **BR-62:** Khách sạn phải có ít nhất 1 loại phòng được tạo trước khi publish.
- **BR-63:** Tọa độ GPS phải nằm trong lãnh thổ Việt Nam.

---

#### ⭐ FR-AG-004: Quản lý loại phòng

**Trường dữ liệu loại phòng:**

| **Trường** | **Bắt buộc** | **Mô tả** |
|---|:---:|---|
| Tên loại phòng | ✅ | VD: Deluxe Sea View |
| Mô tả | ✅ | Rich text |
| Số khách tối đa | ✅ | Người lớn + trẻ em |
| Diện tích phòng | ✅ | m² |
| Loại giường | ✅ | 1 giường đôi / 2 giường đơn / king / queen |
| Số giường | ✅ | Số lượng |
| View phòng | ❌ | City view, Sea view, Garden view, Mountain view |
| Tiện nghi phòng | ✅ | Multi-select (TV, minibar, ban công, bồn tắm...) |
| Hình ảnh | ✅ | Min 3 ảnh, max 10 |
| Số phòng tổng | ✅ | Tổng số phòng loại này KS có |

**Flow chính:**
1. Trong trang chi tiết KS → tab **"Loại phòng"** → **Thêm mới**.
2. Nhập thông tin theo bảng trên.
3. Lưu.

---

#### ⭐ FR-AG-005: Quản lý giá phòng

**Mô tả:** Thiết lập giá phòng linh hoạt.

**Các loại giá:**

| **Loại giá** | **Mô tả** |
|---|---|
| Giá cơ bản | Giá / đêm mặc định |
| Giá theo mùa | Mùa cao điểm / thấp điểm (khoảng ngày cụ thể) |
| Giá theo cuối tuần | Tăng/giảm theo Thứ 6, 7, CN |
| Giá theo ngày lễ | Tăng vào lễ Tết (cấu hình từng ngày) |
| Giá khuyến mãi | Giá đặc biệt trong khoảng ngày |

**Phụ thu:**
- Thêm người lớn: VD +200.000đ/người/đêm
- Thêm trẻ em: theo độ tuổi (free dưới 6, 50% từ 6-12)
- Thêm giường: VD +150.000đ
- Ăn sáng: VD +100.000đ/người

**Flow chính:**
1. Trong loại phòng → tab **"Giá"**.
2. Cấu hình giá cơ bản.
3. Thêm rule giá đặc biệt (theo mùa/cuối tuần/lễ).
4. Cấu hình phụ thu.
5. Preview giá theo ngày cụ thể.

**Business Rule:**
- **BR-64:** Khi đặt phòng, hệ thống tính tổng giá = sum(giá từng đêm) dựa trên rule áp dụng cho đêm đó (ưu tiên giá khuyến mãi > ngày lễ > cuối tuần > theo mùa > cơ bản).
- **BR-65:** Giá hiển thị cho khách = giá / đêm + thuế + phí dịch vụ (nếu có).

---

#### ⭐ FR-AG-006: Quản lý phòng trống (Inventory Calendar)

**Mô tả:** Lịch hiển thị số phòng trống mỗi ngày.

**Flow chính:**
1. Truy cập **"Lịch phòng"** → grid view:
   - Cột: ngày (hiển thị 30-60 ngày)
   - Hàng: loại phòng
   - Ô: số phòng trống / tổng + biểu tượng trạng thái
2. Click ô để:
   - Cập nhật số phòng trống thủ công
   - Khóa phòng (Stop Sale)
   - Mở khóa
3. Thao tác hàng loạt (Bulk update):
   - Chọn khoảng ngày
   - Áp dụng: số phòng cụ thể / khóa tất cả / mở khóa
4. Sync với booking:
   - Booking `CONFIRMED` → số phòng trống tự giảm
   - Booking hủy → số phòng trống tự tăng lại

**Business Rule:**
- **BR-66:** Không cho phép book quá số phòng trống (chống overbooking).
- **BR-67:** Số phòng trống mặc định = số phòng tổng của loại phòng (nếu chưa cấu hình).
- **BR-68:** Đại lý có thể khóa phòng để bảo trì → khách không đặt được nhưng vẫn hiển thị thông tin KS.

---

### 6.3. Module Quản lý tour

#### ⭐ FR-AG-007: Quản lý tour

**Mô tả:** Tạo và quản lý sản phẩm tour nội địa.

**Actor:** Đại lý - Công ty du lịch

**Trường dữ liệu tour:**

| **Trường** | **Bắt buộc** | **Mô tả** |
|---|:---:|---|
| Tên tour | ✅ | VD: "Hà Nội - Sapa - Fansipan 3N2Đ" |
| Mã tour | ✅ | Tự sinh hoặc nhập (unique) |
| Loại tour | ✅ | Multi-select từ danh mục |
| Điểm khởi hành | ✅ | Tỉnh/TP |
| Điểm đến | ✅ | 1 hoặc nhiều (Tỉnh/TP, có thể chi tiết Phường/Xã) |
| Thời lượng | ✅ | VD: 3 ngày 2 đêm |
| Phương tiện | ✅ | Multi-select |
| Số khách tối thiểu | ✅ | VD: 10 khách |
| Số khách tối đa | ✅ | VD: 40 khách |
| Mô tả tổng quan | ✅ | Min 300 ký tự |
| Điểm nổi bật | ❌ | List 3-5 điểm |
| Bao gồm | ✅ | List (xe, KS, ăn uống, vé tham quan, HDV...) |
| Không bao gồm | ✅ | List (phí cá nhân, đồ uống, tip...) |
| Lưu ý quan trọng | ❌ | Rich text |
| Chính sách hủy | ✅ | Chọn template hoặc tự nhập |
| Ảnh gallery | ✅ | Min 5, max 30 |
| Video | ❌ | URL YouTube |
| Brochure PDF | ❌ | Upload file PDF |

**Flow chính:**
1. **Tour của tôi** → **Thêm mới**.
2. Wizard 5 bước:
   - Bước 1: Thông tin cơ bản
   - Bước 2: Lịch trình chi tiết
   - Bước 3: Bao gồm / Không bao gồm / Chính sách
   - Bước 4: Hình ảnh, video, brochure
   - Bước 5: Đợt khởi hành (xem FR-AG-008)
3. Lưu nháp / Gửi duyệt.

---

#### ⭐ FR-AG-008: Quản lý đợt khởi hành (Tour Departure)

**Mô tả:** Mỗi tour có nhiều đợt khởi hành theo lịch.

**Trường dữ liệu:**

| **Trường** | **Bắt buộc** | **Mô tả** |
|---|:---:|---|
| Ngày khởi hành | ✅ | DD/MM/YYYY |
| Ngày kết thúc | ✅ | Tự tính theo thời lượng |
| Số chỗ tổng | ✅ | VD: 30 chỗ |
| Số chỗ đã đặt | Auto | Tự tính |
| Giá người lớn | ✅ | VND |
| Giá trẻ em (5-10 tuổi) | ✅ | VND |
| Giá trẻ em (2-4 tuổi) | ✅ | VND |
| Giá em bé (< 2 tuổi) | ✅ | Thường = 0 hoặc phí nhỏ |
| Phụ thu phòng đơn | ❌ | VND |
| Trạng thái | Auto | Mở / Đóng / Còn ít chỗ (< 5) / Hết chỗ |

**Tính năng tạo nhanh:**
- **Tạo lịch lặp:** VD "Mỗi thứ 7 hàng tuần từ 01/06 đến 31/12" → sinh nhanh nhiều đợt.
- **Copy đợt cũ:** Tạo đợt mới dựa trên đợt đã có.

**Business Rule:**
- **BR-69:** Khi đặt tour, số chỗ đã đặt tăng. Khi hủy, giảm lại.
- **BR-70:** Đợt đầy chỗ tự đổi trạng thái "Hết chỗ", không cho đặt thêm.
- **BR-71:** Đợt khởi hành quá khứ tự ẩn khỏi danh sách tìm kiếm.

---

#### ⭐ FR-AG-009: Quản lý lịch trình tour

**Mô tả:** Soạn lịch trình từng ngày của tour.

**Mỗi ngày có:**

| **Trường** | **Bắt buộc** | **Mô tả** |
|---|:---:|---|
| Tiêu đề | ✅ | VD: "Ngày 1: Hà Nội - Sapa" |
| Mô tả hoạt động | ✅ | Rich text, chi tiết từng giờ |
| Bữa sáng | ✅ | Có / Không |
| Bữa trưa | ✅ | Có / Không |
| Bữa tối | ✅ | Có / Không |
| Khách sạn nghỉ đêm | ❌ | Tên + hạng sao |
| Ảnh minh họa | ❌ | Max 5 ảnh / ngày |
| Điểm đến chính | ✅ | Tọa độ GPS cho Google Maps (xem FR-CO-001) |

**Flow chính:**
1. Trong wizard tạo tour → bước 2.
2. Thêm từng ngày → nhập đầy đủ.
3. Có thể kéo thả sắp xếp lại thứ tự ngày.
4. Bản đồ tổng quan tự render dựa trên tọa độ các ngày.

---

### 6.4. Module Quản lý Combo

#### ⭐ FR-AG-010: Quản lý Combo (Tour + Khách sạn + Vé máy bay)

**Mô tả:** Đại lý tạo combo trọn gói cho khách.

**Actor:** Đại lý (Chủ KS hoặc Công ty du lịch - tùy combo)

**Trường dữ liệu combo:**

| **Trường** | **Bắt buộc** | **Mô tả** |
|---|:---:|---|
| Tên combo | ✅ | VD: "Combo Phú Quốc 3N2Đ - Bay Vietjet + KS 4 sao" |
| Mô tả | ✅ | Rich text |
| Loại combo | ✅ | Tour+KS / KS+Vé máy bay / Tour+KS+Vé máy bay |
| Thành phần | ✅ | Gắn tour ID + KS ID + thông tin vé máy bay |
| Thời lượng | ✅ | Số ngày đêm |
| Số khách áp dụng | ✅ | Min - max |
| Giá gốc (tổng các thành phần) | Auto | Tự tính |
| Giá combo | ✅ | Giá khuyến mãi |
| Mức giảm | Auto | % so với giá gốc |
| Bao gồm | ✅ | List |
| Không bao gồm | ✅ | List |
| Ảnh combo | ✅ | Min 3 ảnh |
| Chính sách hủy | ✅ | Có thể khác với từng thành phần |
| Đợt khởi hành | ✅ | Tương tự tour |

**Thông tin vé máy bay nội địa:**

| **Trường** | **Mô tả** |
|---|---|
| Hãng bay | Vietnam Airlines, Vietjet, Bamboo, Vietravel Airlines... |
| Chặng bay đi | Sân bay khởi hành → Sân bay đến (VD: HAN → SGN) |
| Chặng bay về | Tương tự |
| Hạng vé | Economy / Business |
| Hành lý | Số kg ký gửi, xách tay |
| Ngày bay | Khớp với đợt khởi hành |
| Giờ bay (dự kiến) | Sáng/Trưa/Chiều/Tối hoặc giờ cụ thể |

**Flow chính:**
1. **Combo của tôi** → **Tạo combo**.
2. Chọn loại combo.
3. Chọn tour và/hoặc KS từ danh sách của mình.
4. Nhập thông tin vé máy bay (nếu có).
5. Cấu hình giá combo, đợt khởi hành.
6. Lưu nháp / Gửi duyệt.

**Business Rule:**
- **BR-72:** Chỉ áp dụng chặng bay **nội địa Việt Nam**.
- **BR-73:** Giá combo phải ≤ Tổng giá các thành phần (bắt buộc có discount).
- **BR-74:** Khi book combo, hệ thống tự trừ inventory: chỗ tour + phòng KS.
- **BR-75:** Vé máy bay trong combo là **thông tin tham khảo**, đại lý chịu trách nhiệm đặt vé thực tế cho khách (không tích hợp API hãng bay).

---

### 6.5. Module Quản lý đặt chỗ của đại lý

#### ⭐ FR-AG-011: Danh sách booking đến đại lý

**Flow chính:**
1. Truy cập **"Booking của tôi"**.
2. Danh sách hiển thị: Mã booking, Khách hàng, Sản phẩm, Ngày đặt, Ngày sử dụng, Tổng tiền, Trạng thái.
3. Bộ lọc: trạng thái, loại (tour/KS/combo), thời gian.
4. Click chi tiết → xem đầy đủ.

---

#### ⭐ FR-AG-012: Xác nhận / Từ chối booking

**Flow chính:**
1. Booking mới → `PENDING_CONFIRMATION` → đại lý nhận email + in-app notification.
2. Đại lý xem chi tiết:
   - Thông tin khách
   - Thông tin hành khách
   - Sản phẩm, ngày, số người
   - Yêu cầu đặc biệt
   - Tổng tiền, voucher áp dụng
3. Kiểm tra inventory thực tế.
4. **Xác nhận** → `CONFIRMED` → hệ thống gửi email cho khách kèm hướng dẫn thanh toán (số tài khoản đại lý, nội dung CK).
5. **Từ chối** (kèm lý do) → `CANCELLED_BY_AGENT` → đề xuất phương án thay thế (nếu có).

**Business Rule:**
- **BR-76:** Đại lý phải xác nhận trong 24h. Quá hạn → tự động `CANCELLED_BY_SYSTEM`.
- **BR-77:** Tỷ lệ từ chối > 10% / tháng → admin cảnh cáo.
- **BR-78:** Đại lý phải cập nhật số tài khoản nhận chuyển khoản trong hồ sơ trước khi nhận booking.

---

#### FR-AG-013: Xác nhận thanh toán

**Mô tả:** Khi khách chuyển khoản, đại lý xác nhận đã nhận tiền.

**Flow chính:**
1. Booking ở `WAITING_PAYMENT`.
2. Đại lý kiểm tra tài khoản ngân hàng.
3. Khi nhận được tiền → mở booking → **"Xác nhận đã thanh toán"**.
4. Nhập: ngày nhận, số tham chiếu giao dịch, số tiền.
5. Booking chuyển `PAID`.
6. Hệ thống gửi email xác nhận cho khách.

**Business Rule:**
- **BR-79:** Booking `CONFIRMED` quá 48h chưa thanh toán → cảnh báo cho khách. Quá 72h → tự `CANCELLED_BY_SYSTEM`.
- **BR-80:** Đại lý nhập số tham chiếu giao dịch để admin có thể đối chiếu khi đối soát.

---

#### FR-AG-014: Quản lý check-in/check-out (Khách sạn)

**Flow chính:**
1. Trong chi tiết booking KS → tab **"Tình trạng lưu trú"**.
2. Đến ngày check-in → đại lý đánh dấu **"Đã check-in"** → `IN_PROGRESS`.
3. Đến ngày check-out → đánh dấu **"Đã check-out"** → `COMPLETED`.
4. Có thể nhập ghi chú: ghi nhận sự cố, hư hỏng, đánh giá khách...

---

#### FR-AG-015: Quản lý danh sách khách tour

**Mô tả:** Xem danh sách hành khách theo đợt khởi hành.

**Chức năng:**
- Xem danh sách hành khách: họ tên, CCCD/Passport, SĐT, giới tính, ngày sinh.
- Xuất Excel/PDF.
- Gửi thông báo hàng loạt (email/SMS) cho cả đoàn (lịch trình, lưu ý...).
- In phiếu đoàn.

---

### 6.6. Module Khuyến mãi của đại lý

#### ⭐ FR-AG-016: Tạo voucher / khuyến mãi riêng

**Loại khuyến mãi:**

| **Loại** | **Mô tả** |
|---|---|
| Giảm % | VD: Giảm 15% giá phòng |
| Giảm số tiền | VD: Giảm 500.000đ |
| Early Bird | Đặt trước X ngày, giảm Y% |
| Last Minute | Đặt sát ngày (≤ 3 ngày), giảm Y% |
| Flash Sale | Giảm mạnh trong khoảng giờ ngắn (VD: 12h-14h hôm nay) |
| Combo nhiều đêm | Đặt 3 đêm tặng 1, đặt 5 đêm tặng 2 |
| Buy 1 Get 1 | Mua 1 phòng tặng 1 (theo điều kiện) |

**Trường dữ liệu:** Tương tự voucher hệ thống (xem FR-AD-019), nhưng:
- Phạm vi áp dụng: **chỉ tour/KS/combo của đại lý đó**.
- Không giới hạn đối tượng theo hạng thành viên (chỉ áp dụng cho tất cả).

**Business Rule:**
- **BR-81:** Discount ≤ 50% → tự động active.
- **BR-82:** Discount > 50% → cần admin duyệt (xem BR-35).
- **BR-83:** Voucher đại lý có thể cộng dồn với voucher hệ thống nếu đại lý bật cấu hình "Cho phép cộng dồn".

---

### 6.7. Module Tài chính đại lý

#### ⭐ FR-AG-017: Xem doanh thu & hoa hồng

**Thông tin hiển thị:**
- Tổng doanh thu (theo trạng thái booking)
- Hoa hồng phải trả cho hệ thống
- Số tiền đã đối soát / chưa đối soát
- Lịch sử đối soát

**Tính năng:**
- Lọc theo khoảng thời gian, loại sản phẩm.
- Xuất Excel/PDF.
- Biểu đồ doanh thu theo thời gian.

---

#### FR-AG-018: Đối soát

**Flow chính:**
1. Cuối kỳ (tuần/tháng) → admin gửi bảng đối soát.
2. Đại lý nhận thông báo + xem bảng chi tiết.
3. Đại lý:
   - **Đồng ý** → chuyển khoản hoa hồng → đánh dấu **"Đã chuyển khoản"** + upload chứng từ.
   - **Khiếu nại** → ghi rõ booking có sai số → tạo ticket.
4. Admin xác nhận nhận tiền → đóng kỳ đối soát.

---

### 6.8. Module Tương tác với khách

#### FR-AG-019: Trả lời đánh giá

**Flow chính:**
1. Đại lý xem danh sách review → có thể lọc theo: chưa trả lời / đã trả lời / sao thấp.
2. Click vào review → viết phản hồi.
3. Phản hồi hiển thị công khai dưới review.

**Business Rule:**
- **BR-84:** Mỗi review chỉ phản hồi 1 lần. Có thể sửa trong 24h.
- **BR-85:** Phản hồi không phù hợp (tục, công kích khách) → admin có quyền xóa.

---

#### FR-AG-020: Inbox Q&A với khách

**Mô tả:** Hộp thư trao đổi giữa đại lý và khách.

**Tính năng:**
- Inbox hội thoại theo từng khách.
- Đánh dấu đã đọc / chưa đọc.
- Trả lời nhanh (saved replies / templates).
- Tìm kiếm hội thoại.
- Nhận thông báo realtime.

**Business Rule:**
- **BR-86:** Cấm trao đổi thông tin liên hệ ngoài hệ thống. Hệ thống tự động lọc SĐT, email, link Zalo/Telegram trong tin nhắn → cảnh báo đại lý.
- **BR-87:** Đại lý vi phạm BR-86 nhiều lần → bị cảnh cáo / suspend.

---

### 6.9. Module Báo cáo đại lý

#### FR-AG-021: Dashboard đại lý

**Widget hiển thị:**
- Tổng booking (hôm nay / 7 ngày / 30 ngày)
- Doanh thu, hoa hồng phải trả
- Tỷ lệ lấp đầy phòng (Occupancy rate) - cho KS
- Tỷ lệ chốt tour (số chỗ đã đặt / tổng) - cho tour
- Đánh giá trung bình
- Top 5 tour/phòng bán chạy
- Lượt xem sản phẩm
- Booking pending cần xử lý
- Tin nhắn chưa đọc

---

#### FR-AG-022: Báo cáo chi tiết đại lý

**Các báo cáo:**
- Doanh thu theo thời gian / sản phẩm / kênh (web, affiliate)
- Booking analytics (tỷ lệ confirm, tỷ lệ hủy)
- Phân tích khách hàng (mới / quay lại, độ tuổi, vùng miền)
- Hiệu quả voucher đại lý
- So sánh hiệu suất với kỳ trước

**Tính năng:** Lọc đa chiều, xuất Excel/PDF, lên lịch email báo cáo định kỳ.

---

