---
description: Tài liệu tham khảo API chuẩn dành cho các bên thứ ba thao tác Sinh Ảnh và Video
category: api-integration
---

# Flow Tool - Tài Liệu Tích Hợp API Sinh Media (VEO 3.1)

Tài liệu này là đặc tả kỹ thuật dành cho các công cụ (Tool) bên thứ 3, AI Agent hoặc Frontend Client muốn giao tiếp tự động với hệ thống Flow Tool Backend.

Hệ thống đã được trang bị **Kiến trúc Gom Cụm Đa Tài Khoản (Multi-Account Clustering)**. Nhờ đó, người dùng (Tool) **chỉ cần gửi DUY NHẤT một thông điệp sinh Media (Job) bao gồm cả file ảnh**. Hệ thống sẽ tự động tìm kiếm các tài khoản Google đang nhàn rỗi, tải (upload) ngầm tất cả tài liệu tham chiếu và gọi API khởi tạo hình ảnh/video trực tiếp mà không xảy ra bất kì rủi ro "chuyển mạch" hình ảnh nào.

---

## 🚀 Quy Trình Tích Hợp Giao Tiếp (Workflow)

Workflow tích hợp cực kì đơn giản với 2 bước:
1. **Submit Job**: Gửi một HTTP `POST /v1/jobs` chứa Lệnh Sinh (Prompt) và các Ảnh cục bộ/Ảnh Base64.
2. **Polling Result**: Dùng vòng lặp gọi HTTP `GET /v1/jobs/{job_id}` liên tục (khuyên dùng khoảng cách 3-5 giây/lần) cho đến khi trạng thái trả về là `SUCCESS` hoặc `FAILED`.

---

## 📍 1. API Gửi Yêu Cầu Tạo Media (Submit Job)

Hệ thống tự động thiết lập và upload cấu hình dựa trên biến `mode`.

**Endpoint:** `POST /v1/jobs`

### Cơ Bản Dành Cho Tạo Hình Ảnh (Text-To-Image / Image-To-Image)

#### 📝 T2I (Text to Image)
Dùng để tạo hình ảnh hoàn toàn từ văn bản. Không cần gửi ảnh tham chiếu.
```json
{
  "mode": "T2I",
  "prompt": "Một chú mèo máy đang uống trà sữa trên sao hỏa, phong cách cyberpunk",
  "quality": "quality", // hoặc "fast"
  "ratio": "portrait"  // "portrait", "landscape", "square"
}
```

#### 🖼️ I2I (Image to Image)
Chế độ này cho phép gửi tới đa 8 ảnh nhằm tạo ảnh mới tương đồng hình khối / giao diện. Hệ thống sẽ **tự động Upload tất cả đường dẫn ảnh bạn truyền vào `images`**.
```json
{
  "mode": "I2I",
  "prompt": "Anime style, 4k, masterpiece, biến đổi khung cảnh này",
  "images": [
    "C:\\Users\\dinht\\images\\ref1.png",
    "http://domain.com/url_ref2.jpg"
  ],
  "quality": "high",
  "ratio": "landscape"
}
```
### 🧩 Hỗ Trợ Đa Định Dạng Dữ Liệu (`images`)
Hệ thống giải mã của Flow Tool cho phép Client truyền dữ liệu linh hoạt nhất có thể. Tại tham số mảng `"images"`, mỗi phần tử có thể mang một trong **4 định dạng** sau:
1. **Chuỗi Base64 Data URI** (Tuyệt Đỉnh): `data:image/png;base64,iVBORw0KGgoa...` (Cho phép Frontend nhúng ảnh trực tiếp vào payload, không cần lưu file, backend tự bóc tách và upload).
2. **HTTP/HTTPS URL**: `https://cdn.domain.com/anh_1.jpg` (Backend tự động Download ảnh về máy, giải mã thành Base64 và dọn dẹp bộ nhớ đệm).
3. **Local File Path**: `C:\Users\Admin\images\anh.png` (Sử dụng rất tốt nếu Tools gọi đang nằm chung trên một Server Vật lý).
4. **Media ID (Có sẵn)**: `eef2c1-xxxx-xxxx...` (Sử dụng khi bạn đã biết chính xác file ảnh nằm trên cloud Google, tiết kiệm băng thông upload dư thừa).

---

### Cơ Bản Dành Cho Tạo Video (Text/Image/Ref-To-Video)

#### 🎬 T2V (Text To Video)
Tạo Video thần kì từ ký tự.
```json
{
  "mode": "T2V",
  "prompt": "Một đám mây hình chú cún bay lượn qua dãy núi tuyết hùng vĩ",
  "quality": "fast_lp",
  "ratio": "landscape"
}
```

#### 🎥 R2V (Reference To Video)
Giống như I2I, hệ thống cho phép tạo video được **mô phỏng từ nhiều bức ảnh tham chiếu phong cách/khuôn mặt**.
```json
{
  "mode": "R2V",
  "prompt": "Hoạt hình nhân vật này đang đi dọc bờ biển, hoàng hôn",
  "images": [
    "D:\\ProjectAssets\\chara_design_1.jpg",
    "D:\\ProjectAssets\\chara_design_2.jpg"
  ],
  "quality": "fast_lp",
  "ratio": "landscape"
}
```

#### 🌅 I2V_S (Image To Video - Standard Khởi Đầu)
Đây là chế độ **Tạo Video từ 1 Ảnh Khởi Đầu (Start Frame)**. Bạn truyền **duy nhất 1 ảnh** vào mảng `images`. Hệ thống sẽ tự động Map bức ảnh đầu tiên này thành `start_image`.
```json
{
  "mode": "I2V_S",
  "prompt": "Nhân vật trong tranh bắt đầu chớp mắt và quay đầu lại nhìn thẳng vào camera",
  "images": [
    "C:\\Data\\start_frame_01.jpg"
  ],
  "quality": "fast_lp",
  "ratio": "portrait"
}
```

#### 🔄 I2V_SE (Image To Video - Standard Extend Khởi Đầu & Kết Thúc)
Chế độ chuyển cảnh mượt mà giữa 2 bức ảnh. Bạn truyền **đúng 2 ảnh** vào mảng `images`. Hệ thống quy định:
- `images[0]`: Trở thành ảnh bắt đầu Video (Start Frame)
- `images[1]`: Trở thành ảnh kết thúc Video (End Frame)
```json
{
  "mode": "I2V_SE",
  "prompt": "Zoom xa camera từ ảnh 1 trải dài mở rộng phong cảnh biến thành ảnh 2",
  "images": [
    "C:\\Data\\start_scene.png",
    "C:\\Data\\end_scene.png"
  ],
  "quality": "fast_lp",
  "ratio": "landscape"
}
```

**Response Trả Về Mẫu:**
```json
{
  "success": true,
  "job_id": "92f7c13a-xxxx-xxxx",
  "message": "Job added to Producer-Consumer queue."
}
```

---

## ⚙️ Bảng Tra Cứu Tham Số Quality & Ánh Xạ Model

Một trong những thiết lập rất quan trọng khi gửi request là tham số `"quality"`. Phía Backend của Flow Tool chia cấu hình sinh Media thành 2 ngạch Tách Biệt: Hình Ảnh (Image) và Video.

### 1. Đối Với Cụm Tạo Video (T2V, R2V, I2V_S, I2V_SE)
Khác với ảnh, Video được Google phân bậc thời gian Render rất chặt chẽ thông qua các nhánh Model VEO 3.1 khác nhau. Có **3 giá trị** bạn có thể điền vào trường `"quality"`:

- `"quality"` (Độ Nét Cao Khung Hình Chuẩn): 
  - Sẽ ánh xạ tới các Model: `veo_3_1_t2v`, `veo_3_1_i2v_s`... 
  - (Khuyến cáo sử dụng cho Job quan trọng, nhưng tốn nhiều thời gian render).
- `"fast"` (Thời Gian Sinh Nhanh - Fast Ultra):
  - Ánh xạ tới các Model: `veo_3_1_t2v_fast_ultra`, `veo_3_1_i2v_s_fast_ultra`...
  - (Tốc độ sinh video lý tưởng, chất lượng không giảm đáng kể).
- `"fast_lp"` (**🏆 KHUYÊN DÙNG:** Fast Lower Priority - Rảnh Mới Nuôi):
  - Tham số đặc biệt. Sẽ ánh xạ tới model kết thúc bằng `_relaxed` (Ví dụ: `veo_3_1_t2v_fast_ultra_relaxed`).
  - **Khuyến cáo bắt buộc nên biến đây thành tham số mặc định cho toàn bộ Tools**. Chế độ này sinh Video với rủi ro bị Google chặn thấp nhất. Nó xếp hàng thông minh và nhường đường khi Server Google quá tải, giúp bảo vệ tính toàn vẹn của Token lâu dài hơn. Các ví dụ Payload bên trên đều đã được thiết lập sẵn là `fast_lp`.

*(Lưu ý: Riêng chế độ **R2V**, API của Google hoàn toàn **Không Hỗ Trợ** chế độ `"quality"` cao cấp. Nếu bạn bắn `mode: R2V` với `"quality": "quality"`, Job sẽ báo `FAILED`)*.

### 2. Đối Với Cụm Tạo Hình Ảnh (T2I, I2I)
Sinh hình ảnh hiện tại được định tuyến cứng vào một nhánh Workflow cực thấp (chạy siêu nhanh) của FlowMedia.
- Đối với Endpoint ảnh (T2I, I2I), giá trị `"quality"` (như `fast` hay `quality`) bạn truyền lên **Tuyệt đối không ảnh hưởng** tới cái Model được gọi. 
- Nó chỉ mang tính chất ghi nhận log, trong khi Backend mặc định ánh xạ mọi trường hợp sinh Ảnh vào kiến trúc Nano Banana Pro (Model `GEM_PIX_2`). Nên bạn cứ tự tin truyền `"fast"` làm giá trị chuẩn xác.

---

## 🔎 2. API Theo Dõi Kết Quả Job (Polling Status)

Bạn sử dụng `job_id` nhận được ở phía trên để theo dõi kết quả.

**Endpoint:** `GET /v1/jobs/{job_id}`

**Response Trả Về Mẫu (Khởi tạo):**
```json
{
  "job_id": "92f7c13a-xxxx-xxxx",
  "status": "QUEUED",
  "progress": 0,
  "message": "In Queue: QUEUED",
  "result_urls": [],
  "error": null
}
```

**Response Trả Về Mẫu (Đang Sinh - Xử lý trong hệ thống Load Balancer):**
```json
{
  "job_id": "92f7c13a-xxxx-xxxx",
  "status": "PROCESSING",
  "progress": 50,
  "message": "Processing on profile 5f7ddcc2",
  "result_urls": [],
  "error": null,
  "profile_id": "5f7ddcc2"
}
```

**Response Trả Về Mẫu (Nhận Kết Quả Hoàn Thiện):**
Khi `status` = `SUCCESS`, bạn sẽ nhận được mảng chứa đường dẫn File kết quả ở `result_urls`.

```json
{
  "job_id": "92f7c13a-xxxx-xxxx",
  "status": "SUCCESS",
  "progress": 100,
  "message": "Generation completed successfully",
  "result_urls": [
    "https://storage.googleapis.com/path_to_video_or_image_out.mp4"
  ],
  "error": null,
  "profile_id": "5f7ddcc2"
}
```

### 🛑 Các Trạng Thái (Status) Bạn Cần Lắng Nghe:
- **`QUEUED` / `PENDING`**: Request đang được đưa vào hàng chờ để tìm kiếm tải rảnh của các Tài khoản Google chưa đạt Limit (Mỗi Account gánh song song 4 tiến trình). Không việc gì phải sốt ruột.
- **`PROCESSING`**: Giai đoạn dài nhất. Ở chế độ này Backend đang upload ảnh tĩnh ngầm cho bạn hoặc đang khấn vái Video API Generator. Quá trình này mất khoảng từ 20 giây - 1 phút rưỡi. Mời tiếp tục Polling.
- **`SUCCESS`**: Tạo Media thành công. Nhanh tay chộp lấy link tại `result_urls`. Vòng lặp Polling của Client nên Stop ở đây.
- **`FAILED` / `CANCELLED`**: Sự cố xuất hiện (Lỗi token, lỗi logic tham số, hoặc lỗi Google Sandbox). Chuỗi lỗi giải trình sẽ có ở field `error`. Polling nên Stop ở đây.
