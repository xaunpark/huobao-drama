# Multi-Profile Anti-Detect Browser System

## Tổng quan

Hệ thống quản lý nhiều browser profile với tính năng Anti-Detect, giúp:
- Chạy nhiều tài khoản Google VideoFX song song
- Mỗi profile có fingerprint riêng biệt (User-Agent, Viewport, Timezone, Locale)
- Hỗ trợ Proxy per-profile
- Inject stealth scripts để bypass detection

## Cấu trúc

```
automation/
├── profile_manager.py    # Quản lý profiles
├── browser_manager.py    # Khởi tạo browser với anti-detect
└── engine.py             # Hỗ trợ chọn profile per-job

profiles/                  # Thư mục chứa profile data
├── profiles_config.json  # Cấu hình tất cả profiles
├── abc12345/             # Profile data folder
│   └── storage.json
└── def67890/
    └── storage.json

manage_profiles.py        # CLI tool quản lý profiles
```

## Sử dụng

### 1. Tạo Profile mới

```bash
python manage_profiles.py
# Chọn option 2: Create new profile
```

Hoặc trong code:

```python
from automation.profile_manager import get_profile_manager

pm = get_profile_manager()

# Tạo profile với fingerprint random
profile = pm.create_profile(
    name="Account 1",
    google_account="myemail@gmail.com",
    platform="windows_chrome"  # hoặc mac_chrome, linux_chrome
)

# Tạo profile với proxy
profile = pm.create_profile(
    name="Account 2",
    google_account="another@gmail.com",
    proxy={'server': 'http://proxy-ip:port'}
)

# Tạo profile với fingerprint tùy chỉnh
profile = pm.create_profile(
    name="Custom Profile",
    custom_user_agent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) ...",
    custom_viewport={'width': 1920, 'height': 1080},
    custom_timezone="America/New_York",
    custom_locale="en-US"
)
```

### 2. Import Chrome Profile có sẵn

```python
pm = get_profile_manager()

# Import folder Chrome profile đã login
profile = pm.import_from_folder(
    folder_path="C:/Users/xxx/AppData/Local/Google/Chrome/User Data/Profile 1",
    name="Existing Account",
    google_account="existing@gmail.com"
)
```

### 3. Sử dụng Profile trong API

```json
POST /v1/jobs
{
  "prompt": "...",
  "profile_id": "abc12345",
  "settings": {...}
}
```

### 4. Xem danh sách Profiles

```python
for p in profiles:
    print(f"{p.profile_id}: {p.name} - {p.google_account}")
```

### 5. Sticky Batch Routing (Quan trọng cho Đa Tài Khoản)

Hệ thống hỗ trợ cơ chế "Dính" tài khoản theo lô (Sticky Batch Routing). Vì tài sản (Media IDs) trong Google VideoFX là riêng tư cho từng tài khoản, việc gán nhầm tài khoản khi tạo Video từ ảnh đã upload sẽ gây lỗi 404.

**Cách dùng:**
1. Client gửi kèm `batch_id` (ví dụ UUID) trong tất cả các request Upload.
2. Server sẽ "khóa" (pin) `batch_id` đó vào 1 Profile duy nhất.
3. Khi nhận lệnh tạo Video với cùng `batch_id`, Server tự động điều phối về Profile đã chọn.

Điều này đảm bảo tính toàn vẹn của dữ liệu mà Client không cần phải tự quản lý việc chọn tài khoản.

## Anti-Detect Features

### Fingerprint được random:
- **User-Agent**: Chrome trên Windows/Mac/Linux
- **Viewport**: Nhiều resolution phổ biến (1920x1080, 1366x768, etc.)
- **Timezone**: Múi giờ phổ biến
- **Locale**: Ngôn ngữ (en-US, vi-VN, etc.)

### Stealth Scripts được inject:
- Override `navigator.webdriver`
- Spoof `navigator.plugins`
- Randomize Canvas fingerprint
- Override permissions query

### Proxy Support:
```python
profile = pm.create_profile(
    name="With Proxy",
    proxy={
        'server': 'http://ip:port',
        'username': 'user',  # optional
        'password': 'pass'   # optional
    }
)
```

## Best Practices

1. **Mỗi tài khoản = 1 profile**: Không dùng chung profile cho nhiều tài khoản
2. **Đa dạng fingerprint**: Để các profile có User-Agent khác nhau
3. **Sử dụng proxy**: Nếu chạy nhiều profile, nên dùng proxy khác nhau
4. **Rotate profiles**: Dùng `get_least_used_profile()` để cân bằng tải
5. **Monitor stats**: Xem `jobs_completed` và `jobs_failed` để phát hiện profile bị block

## Troubleshooting

### Profile bị block?
```python
# Deactivate profile
pm.deactivate_profile("abc12345")

# Profile sẽ không được chọn trong list_profiles(active_only=True)
```

### Cần đăng nhập lại?
1. Mở browser với mode không headless
2. Đăng nhập Google
3. Profile sẽ lưu cookies tự động
