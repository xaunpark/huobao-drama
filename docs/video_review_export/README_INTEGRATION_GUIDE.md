# AI Vision Video Review System - Integration Guide

Tài liệu này đóng gói toàn bộ quy trình chấm điểm video AI bằng Claude Vision, cho phép bạn mang theo và tích hợp (clone) vào bất kỳ dự án sinh video AI nào khác.

## 1. Thành phần lõi (Core Components)

Trong thư mục này, tôi đã copy sẵn 3 file quan trọng nhất để project khác của bạn có thể sử dụng lại:
- **`video_reviewer.py`**: Trái tim của hệ thống. Chứa logic băm mổ video bằng FFMPEG, ghép ảnh Contact Sheet, và chứa đoạn văn bản (System Prompt) thần thánh để "dạy" Claude nhận diện mọi loại lỗi Generate Video như chảy mặt, lỗi Reverse, lỗi mọc thêm chi.
- **`review_models.py`**: Chứa cấu trúc dữ liệu Pydantic (Model) giúp ép output của AI thành kiểu JSON tiêu chuẩn theo thang điểm từ `0.0` đến `10.0` qua 6 hạng mục thẩm định.
- **`fk-review-video.md`**: File hướng dẫn (Skill Workflow) dành cho lập trình viên/LLM về cách vận hành hệ thống này trong quy trình AI tự động lấy lỗi để Regen.

## 2. dependencies cần có ở Project Mới

Để code Python trong thư mục này chạy được trên Project mới của bạn, bạn cần thiết lập:

### System (Hệ điều hành)
- Phải cài đặt **FFMPEG**. Code dùng `subprocess.run` để gọi `ffmpeg` (để chiết xuất ảnh) và `ffprobe` (để lấy độ dài video).
- (Optional) Cài đặt **Claude CLI** toàn cầu (`npm install -g @anthropic-ai/claude-cli`) nếu bạn muốn dùng đường dẫn bypass chạy phân tích không cần cắm Key thẳng vào code.

### Python Packages (`pip install ...`)
- `aiohttp` / `asyncio` (Dùng gọi API song song).
- `anthropic` (SDK chính chủ của Claude nếu bạn chơi hệ nạp thẻ API Key thẳng vào `.env`).
- `pydantic` (Dùng ép khuôn kiểu dữ liệu JSON).

## 3. Kiến trúc chạy & Tích hợp (Integration Architecture)

Project mục tiêu của bạn cần sửa lại / điều biến làm sao để luồng chạy khớp với 3 bước sau:

### Bước 1: Thu thập bộ Data Đầu Vào (Inputs)
Để hệ thống chấm khách quan, bạn phải giữ lại được các thông số bạn thu thập lúc sinh video (vứt cho hàm Review):
1. **Video File**: `shot_01.mp4` thao tác nội bộ (tải về máy trước khi chấm).
2. **Text Prompts**: `image_prompt` (Mô tả bối cảnh) + `video_prompt` (Mô tả hành động của Camera/Nhân vật).
3. *(Chỉ với tính năng I2V nâng cao)* **Target References**: Đường link của ảnh Nhân Vật tham chiếu, để khi chẻ frame, Claude còn lấy cái đó làm điểm neo (Reference Anchor).

### Bước 2: Gọi Module `video_reviewer.py`
Hãy import hàm `review_scene_video` từ `video_reviewer.py`. Hàm này sẽ lo việc:
- **Fast Extract**: Tính độ dài video (Thường Fx cho ra khoảng 8 giây). Đem nhân với Frame rate mong muốn (Mặc định `REVIEW_FPS_LIGHT = 4fps`). Ra tổng cộng 32 frames. 
- Mọc ra `contact_sheet.jpg` (Grid 8 cột chứa mọi frame, độ phân giải nhẹ) hoặc tách riêng mớ base64 Jpegs nếu dùng SDK xịn.
- Nạp khối Prompt định nghĩa Rubric chấm điểm ở dòng ~205. Claude không chỉ chấm điểm vô tri mà nó bắt buộc nhận diện 14 Lỗi Kinh Địa (Drift, Breed Swap, Reverse Motion, Brand Logo...). 

### Bước 3: Đọc và Xử lý Lỗi (Retake / Regen Loop)
Khi Module trả về output JSON (chuẩn theo form của `review_models.py`), bạn hứng lấy kết quả tại Project mới của bạn.
- Đoạn quan trọng nhất không nằm ở tính năng In điểm, mà là **Sửa Sai Tự Động (Auto-Correct)**!
- Lấy `overall_score`. Nếu `< 7.5`:
    - Lặp qua dãy `errors` (Mảng json chứa các lỗi `CRITICAL` và `HIGH`).
    - Lấy mô tả lỗi của AI, gắn nối chuỗi đuôi vào câu Lệnh `video_prompt` ban đầu.
    - VD: Video Prompt = *"Sơn Tinh múa gậy."* => AI trả lỗi: Múa gậy trật tay đập vào mặt. => Gắn lại: *"Sơn Tinh múa gậy. DO NOT HIT THE FACE. KEEP MOTION SMOOTH"*.
    - Bắn lệnh lên AI Generator (Luma/Google/Runway) bắt đẻ video lại!

Bằng cách bứng cụm Code này đi theo các project sinh video của bạn, bạn đã có trong tay một **Gatekeeper (Người gác cổng)** chất lượng cao cho bất kỳ pipeline sản xuất hàng loạt nào. Chúc bạn Clone thành công!
