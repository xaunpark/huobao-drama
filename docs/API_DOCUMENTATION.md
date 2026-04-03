# Veo Flow API Documentation

## Tổng quan
API Service cho phép các hệ thống bên ngoài gửi yêu cầu tạo video/ảnh AI thông qua Google Veo/VideoFX.

**Base URL:** `http://localhost:8000`  
**Swagger Docs:** `http://localhost:8000/docs`

---

## Chế độ Thực thi (Execution Mode)

Server chỉ hỗ trợ **DIRECT_API** mode (sử dụng Direct API với captcha bypass).

Xem mode hiện tại:
```bash
GET /v1/settings/mode
```

---

## Captcha Provider

Khi dùng Direct API, server sử dụng dịch vụ bypass reCAPTCHA:

| Provider | Mô tả |
|----------|-------|
| `nanoai` (mặc định) | nanoai.pics - Ổn định |
| `omocaptcha` | OmoCaptcha - Backup |

**Lưu ý**: Chỉ requests **tạo video/image** mới cần captcha. Upload ảnh **KHÔNG** cần captcha.

```bash
# Xem provider
GET /v1/settings/captcha-provider

# Đổi provider
POST /v1/settings/captcha-provider?provider=nanoai
```

---

## Generation Modes

| Code | UI Name | Mô tả | Images | Max |
|------|---------|-------|--------|-----|
| `T2V` | Text to Video | Tạo video từ text | ❌ | - |
| `I2V_S` | Image to Video (Only Start) | Video từ 1 ảnh start | ✅ | 1 |
| `I2V_SE` | Image to Video (Start + End) | Video từ 2 ảnh | ✅ | 2 |
| `R2V` | Reference Images (R2V) | Video từ ảnh tham khảo | ✅ | 8 |
| `T2I` | Text to Image | Tạo ảnh từ text | ❌ | - |
| `I2I` | Image to Image | Tạo biến thể ảnh | ✅ | 8 |

---

## Sticky Batch Routing (Quan trọng)

Để hỗ trợ đa tài khoản (multi-profile), hệ thống sử dụng cơ chế **Sticky Routing**. Nếu bạn upload nhiều ảnh để tạo 1 video (như chế độ R2V), bạn **PHẢI** sử dụng chung một `batch_id` cho tất cả các yêu cầu để đảm bảo chúng được lưu trên cùng một tài khoản Google.

**Quy trình chuẩn:**
1. Tự sinh một mã `batch_id` duy nhất (ví dụ: UUID).
2. Gửi `batch_id` này kèm theo trong **tất cả** các request `POST /v1/upload`.
3. Gửi `batch_id` này kèm theo trong request `POST /v1/jobs`.

Server sẽ đảm bảo chọn duy nhất 1 tài khoản cho toàn bộ lô (batch) này, tránh lỗi 404 "Not Found" do ảnh nằm ở tài khoản khác.

---

## API Endpoints

### 1. Upload Image
**`POST /v1/upload`**

Upload ảnh để lấy media_id trước khi tạo video/image.

#### Request
```json
{
  "image_data": "data:image/png;base64,iVBORw0KGgo...",
  "mime_type": "image/png",
  "batch_id": "uuid-batch-123" // Optional nhưng khuyên dùng cho batch jobs
}
```

#### Response
```json
{
  "success": true,
  "media_id": "CAMaJDUxNjE5ZmIyLTQ5NjItN...",
  "error": null
}
```

---

### 2. Create Job
**`POST /v1/jobs`**

Tạo job video/image mới. Trả về `job_id` để polling status.

#### Request Body

```json
{
  "prompt": "Mô tả video/ảnh cần tạo",
  "mode": "T2V",
  "quality": "fast",
  "ratio": "landscape",
  "reference_image_ids": [],
  "start_image_id": null,
  "end_image_id": null,
  "webhook_url": "https://callback-domain.com/webhook",
  "settings": {},
  "wait_for_result": false,
  "batch_id": "uuid-batch-123" // Đảm bảo khớp với batch_id lúc upload ảnh
}
```

#### Parameters

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `prompt` | string | ✅ | Mô tả nội dung cần tạo |
| `mode` | string | ✅ | `T2V`, `I2V_S`, `I2V_SE`, `R2V`, `T2I`, `I2I` |
| `quality` | string | ❌ | `fast` (mặc định), `quality`, hoặc `fast_lp` (Fast Lower Priority) |
| `ratio` | string | ❌ | `landscape` (mặc định) hoặc `portrait` |
| `reference_image_ids` | array | ❌ | Media IDs cho R2V/I2I (max 8) |
| `start_image_id` | string | ❌ | Media ID cho I2V_S/I2V_SE |
| `end_image_id` | string | ❌ | Media ID cho I2V_SE |
| `webhook_url` | string | ❌ | Địa chỉ webhook để nhận thông báo khi job hoàn thành |
| `settings` | object | ❌ | Ghi đè các parameters tuỳ chỉnh nâng cao |
| `wait_for_result` | bool | ❌ | Chờ kết quả (default: false) |
| `batch_id` | string | ❌ | Mã định danh lô để đảm bảo "dính" vào cùng một profile |

#### Response
```json
{
  "job_id": "abc12345",
  "status": "queued",
  "progress": 0,
  "message": "Job queued",
  "result_urls": [],
  "error": null
}
```

---

### 3. Get Job Status
**`GET /v1/jobs/{job_id}`**

#### Response
```json
{
  "job_id": "abc12345",
  "status": "success",
  "progress": 100,
  "message": "",
  "result_urls": ["https://storage.googleapis.com/..."],
  "error": null
}
```

### 4. Upscale Video (1080p)
**`POST /v1/jobs/{job_id}/upscale`**

Kích hoạt tính năng Upscale 1080p (sử dụng model `veo_3_1_upsampler_1080p`) cho một video đã hoàn thành. Tự động nhận diện `aspect ratio` từ job gốc.

#### Behavior
- Chỉ áp dụng cho Job Video đã `COMPLETED` hoặc Job Upscale bị `FAILED` (để Retry).
- Cập nhật trạng thái Job hiện tại thành `QUEUED` và Mode thành `UPSCALE`.
- **Sticky Assignment:** Job upscale sẽ **ưu tiên sử dụng chung Profile/Account** đã tạo ra video gốc để đảm bảo tính xác thực `mediaId`.
- **KHÔNG** tạo Job ID mới. Kết quả video upscale sẽ là file mới tải về, có thể ghi đè/nằm cạnh file cũ.

#### Response
```json
{
  "success": true,
  "job_id": "abc12345",
  "message": "Upscale queued"
}
```

---

### 5. Upscale Image (2K/4K)
**`POST /v1/jobs/{job_id}/upscale-image`**

Kích hoạt Upscale Ảnh (sử dụng Google Flow `upsampleImage`) cho một ảnh đã được tạo thành công trên hệ thống. 

#### Request Body
Tuỳ chọn truyền thêm `target_resolution`, nếu không mặc định là `2K`.
```json
{
  "target_resolution": "4K" // Hoặc "2K"
}
```

#### Behavior
- Tương tự Upscale Video, **Sticky Assignment** sẽ khóa chặt Profile/Account đã tạo ảnh gốc.
- Job sẽ được cập nhật thành mode `UPSCALE_IMAGE`.

#### Response
```json
{
  "success": true,
  "job_id": "abc12345",
  "message": "Image Upscale (4K) queued"
}
```

---

## Queue Management

### 1. Queue Statistics
**`GET /v1/queue/stats`**

Lấy thông tin thống kê hàng đợi.

#### Response
```json
{
  "queue": {
    "total": 50,
    "queued": 5,
    "processing": 3,
    "completed": 40,
    "failed": 2,
    "cancelled": 0,
    "available_slots": 495
  },
  "dispatcher": {
    "is_running": true,
    "is_paused": false,
    "active_threads": 3
  }
}
```

### 2. List Jobs
**`GET /v1/queue/jobs`**

Lấy danh sách job trong hàng đợi.

**Params:**
- `limit`: Số lượng job (default: 50)
- `status`: Lọc theo status (optional)

### 3. Queue Control
- **`POST /v1/queue/pause`**: Tạm dừng dispatching job mới.
- **`POST /v1/queue/resume`**: Tiếp tục dispatching.
- **`POST /v1/queue/clear`**: Xóa tất cả job đang chờ (QUEUED).

#### Status Values
| Status | Meaning |
|--------|---------|
| `queued` | Đang chờ xử lý |
| `pending` | Đang chuẩn bị |
| `processing` | Đang tạo |
| `success` | Hoàn thành ✅ |
| `failed` | Thất bại ❌ |

---

## Ví dụ Usage

### Text to Video (T2V)
```bash
curl -X POST http://localhost:8000/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A cat playing piano in a jazz club",
    "mode": "T2V",
    "quality": "fast",
    "ratio": "landscape"
  }'
```

### Image to Video - Start Only (I2V_S)

**Step 1: Upload image**
```bash
curl -X POST http://localhost:8000/v1/upload \
  -H "Content-Type: application/json" \
  -d '{
    "image_data": "data:image/png;base64,iVBORw0KGgo...",
    "mime_type": "image/png"
  }'
# Returns: {"success": true, "media_id": "CAMaJDUx..."}
```

**Step 2: Create job**
```bash
curl -X POST http://localhost:8000/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Camera slowly zooms in on the character",
    "mode": "I2V_S",
    "quality": "fast",
    "ratio": "landscape",
    "start_image_id": "CAMaJDUx..."
  }'
```

### Image to Video - Start + End (I2V_SE)
```bash
curl -X POST http://localhost:8000/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Smooth morphing transition between frames",
    "mode": "I2V_SE",
    "quality": "quality",
    "ratio": "landscape",
    "start_image_id": "CAMaJDUx...",
    "end_image_id": "CAMbKE..."
  }'
```

### Reference Images to Video (R2V)
```bash
curl -X POST http://localhost:8000/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "The character is dancing in a forest",
    "mode": "R2V",
    "quality": "fast",
    "ratio": "landscape",
    "reference_image_ids": ["mediaId1", "mediaId2", "mediaId3"]
  }'
```

### Text to Image (T2I)
```bash
curl -X POST http://localhost:8000/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A fantasy castle on a mountain at sunset",
    "mode": "T2I",
    "ratio": "landscape"
  }'
```

### Image to Image (I2I)
```bash
curl -X POST http://localhost:8000/v1/jobs \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Same scene but in cyberpunk style",
    "mode": "I2I",
    "ratio": "landscape",
    "reference_image_ids": ["sourceImageMediaId"]
  }'
```

---

## Python Example (Full Workflow)

```python
import requests
import time
import base64

API_BASE = "http://localhost:8000"

def upload_image(image_path: str) -> str:
    """Upload image and return media_id"""
    with open(image_path, 'rb') as f:
        data = base64.b64encode(f.read()).decode()
    
    resp = requests.post(f"{API_BASE}/v1/upload", json={
        "image_data": f"data:image/png;base64,{data}",
        "mime_type": "image/png"
    })
    result = resp.json()
    if result.get("success"):
        return result["media_id"]
    raise Exception(result.get("error"))

def create_r2v_job(prompt: str, image_paths: list) -> str:
    """Create R2V job with reference images using Batch ID (Best Practice)"""
    import uuid
    batch_id = str(uuid.uuid4())
    
    # Upload all images with SAME batch_id
    media_ids = []
    for path in image_paths[:8]:
        with open(path, 'rb') as f:
            data = base64.b64encode(f.read()).decode()
            
        resp = requests.post(f"{API_BASE}/v1/upload", json={
            "image_data": f"data:image/png;base64,{data}",
            "mime_type": "image/png",
            "batch_id": batch_id
        })
        res = resp.json()
        if res.get("success"):
            media_ids.append(res["media_id"])

    # Create job with SAME batch_id
    resp = requests.post(f"{API_BASE}/v1/jobs", json={
        "prompt": prompt,
        "mode": "R2V",
        "quality": "fast",
        "ratio": "landscape",
        "reference_image_ids": media_ids,
        "batch_id": batch_id
    })
    return resp.json()["job_id"]

def poll_job(job_id: str, max_wait: int = 300) -> dict:
    """Poll job until complete"""
    start = time.time()
    while time.time() - start < max_wait:
        resp = requests.get(f"{API_BASE}/v1/jobs/{job_id}")
        status = resp.json()
        
        if status["status"] in ["success", "completed", "done"]:
            return {"success": True, "urls": status.get("result_urls", [])}
        elif status["status"] in ["failed", "error"]:
            return {"success": False, "error": status.get("error")}
        
        time.sleep(5)
    return {"success": False, "error": "Timeout"}

# Usage Example
job_id = create_r2v_job(
    "Character dancing in a forest",
    ["C:/images/char1.png", "C:/images/char2.png"]
)
print(f"Job created: {job_id}")

result = poll_job(job_id)
if result["success"]:
    print(f"Video ready: {result['urls']}")
else:
    print(f"Failed: {result['error']}")
```

---

## UI Prompt Syntax

Khi sử dụng UI (`python main.py`), có thể embed ảnh trực tiếp trong prompt:

| Tag | Mode | Mô tả |
|-----|------|-------|
| `[START_IMG: path]` | I2V_S, I2V_SE | Ảnh bắt đầu |
| `[END_IMG: path]` | I2V_SE | Ảnh kết thúc |
| `[REF_IMG: path]` | R2V, I2I | Ảnh tham khảo (max 8) |

**Ví dụ:**
```
[REF_IMG: C:/images/char.png] [REF_IMG: C:/images/bg.jpg] Character dancing in forest
```

---

## Technical Notes: Google Flow Image API

Khi tích hợp với API nội bộ của Google (`flowMedia:batchGenerateImages`) cho việc sinh ảnh (I2I/T2I), tính toán payload có một số khác biệt cốt lõi kể từ bản cập nhật gần đây:

**1. Schema Hỗ Trợ Nhiều Ảnh (Multi-Reference):**
Để truyền nhiều ảnh tham chiếu (tối đa 8 ảnh) thông qua array `imageInputs`, bắt buộc payload phải thay đổi từ trường `"prompt": "..."` truyền thống sang `"structuredPrompt": {"parts": [{"text": "..."}]}`. Nếu vẫn dùng `prompt` thông thường mà `imageInputs` > 1, Google API sẽ lập tức trả về lỗi `HTTP 400: INVALID_ARGUMENT`.

**2. Required Payload Keys:**
- `"useNewMedia": true` (Root level)
- `"mediaGenerationContext": {"batchId": "<uuid4>"}` (Khai báo UUID batch id ở root level payload)
- Từng ảnh truyền vào `imageInputs` chỉ mang kiểu `"IMAGE_INPUT_TYPE_REFERENCE"` và `"name"` sẽ là UUID (assetId) hoặc raw schema từ `uploadUserImage` (vd `cbab2...` từ `mediaGenerationId`).

(Hệ thống Flow Tool đằng sau `utils/direct_api.py` đã tự động handling và wrap schema mới này, client chỉ cần post vào `POST /v1/jobs` như thường lệ)

---

## Rate Limits & Notes

1. **Captcha Bypass**: Mỗi request video/image cần captcha bypass (~2-5 giây)
2. **Video Generation**: 1-5 phút tùy quality
3. **Image Generation**: 5-15 giây (synchronous)
4. **Max Reference Images**: 8 ảnh cho R2V và I2I
5. **File Paths**: Sử dụng đường dẫn tuyệt đối

---

## Quick Reference

| Action | Endpoint | Method |
|--------|----------|--------|
| Upload image | `/v1/upload` | POST |
| Create job | `/v1/jobs` | POST |
| Get status | `/v1/jobs/{job_id}` | GET |
| List profiles | `/v1/profiles` | GET |
| Get mode | `/v1/settings/mode` | GET |
| Get captcha | `/v1/settings/captcha-provider` | GET |
