# Hướng Dẫn Sử Dụng Hệ Thống AI Agents

> Dành cho người dùng (bạn) — không phải cho agent.

---

## Phần 1: Cách Dùng Hàng Ngày

### Bạn KHÔNG CẦN làm gì đặc biệt

AGENTS.md đã được tự động load. Bạn chỉ cần gửi yêu cầu bình thường:

```
❌ Sai: "Hãy đọc AGENTS.md trước, rồi tìm router phù hợp, rồi..."
✅ Đúng: "Thêm provider Kling vào video generation"
```

Agent đã biết cấu trúc project và sẽ tự tham khảo `ai/` khi cần.

### Khi nào nên THÊM ngữ cảnh vào prompt

| Tình huống | Prompt đơn giản đủ? | Nên thêm gì |
|-----------|---------------------|-------------|
| Thêm feature nhỏ | ✅ Đủ | — |
| Feature phức tạp, nhiều file | ⚠️ Nên thêm | `"Tham khảo ai/memory/conventions.md trước khi bắt đầu"` |
| Debug lỗi lạ | ⚠️ Nên thêm | `"Kiểm tra ai/memory/risks.md xem có known issue không"` |
| Thêm storyboard mode mới | ✅ Agent sẽ tự tìm | Nhưng có thể gợi ý: `"Theo quy trình trong ai/skills/new-storyboard-mode.md"` |
| Sửa prompt template | ✅ Đủ | — |
| Refactor code lớn | ⚠️ Nên thêm | `"Tuân thủ ai/rules/forbidden-patterns.md"` |

---

## Phần 2: Các Mẫu Prompt Hiệu Quả

### Mẫu 1: Yêu cầu trực tiếp (dùng phổ biến nhất)

```
Thêm endpoint API mới cho quản lý audio tracks.
```

Agent sẽ tự biết: cần handler + service + route + types.

### Mẫu 2: Yêu cầu + chỉ định tham khảo (khi cần chính xác hơn)

```
Thêm image provider mới tên "Midjourney". 
Tham khảo ai/memory/conventions.md mục "Adding a New AI Provider" để đúng quy trình.
```

### Mẫu 3: Debug với context (khi gặp lỗi lạ)

```
Video generation bị stuck ở status "processing" không bao giờ complete.
Kiểm tra ai/memory/risks.md xem có known issue liên quan không, 
rồi trace qua ai/systems/video-pipeline.md để tìm root cause.
```

### Mẫu 4: Yêu cầu kế hoạch trước khi code

```
Tôi muốn thêm chế độ "Podcast" cho storyboard generation.
Lên kế hoạch chi tiết theo ai/skills/new-storyboard-mode.md trước, 
đừng code ngay.
```

### Mẫu 5: Yêu cầu tuân thủ rules cụ thể

```
Refactor image_generation_service.go để tách batch logic ra file riêng.
Tuân thủ nghiêm ngặt ai/rules/forbidden-patterns.md và ai/rules/architecture-rules.md.
```

---

## Phần 3: Cập Nhật Hệ Thống AGENTS

### Nguyên tắc: Agent system phải SỐNG cùng code

Mỗi khi project thay đổi đáng kể, hệ thống `ai/` cũng cần cập nhật.

### 3A. Cập nhật tự động (thêm vào cuối mỗi task lớn)

Thêm dòng này vào cuối prompt của bạn:

```
Sau khi hoàn thành, cập nhật các file ai/ liên quan nếu có thay đổi đáng kể 
(thêm feature vào feature-map, thêm risk vào risks.md, cập nhật system docs, v.v.)
```

### 3B. Cập nhật theo sự kiện cụ thể

#### Khi thêm feature mới:
```
Cập nhật ai/indexes/feature-map.md — thêm feature [tên] vào danh sách 
với đúng file mapping backend + frontend.
```

#### Khi thêm AI provider mới:
```
Cập nhật ai/indexes/dependency-map.md — thêm provider [tên] vào bảng External AI Service Dependencies.
Cập nhật ai/systems/video-pipeline.md (hoặc image-pipeline.md) — thêm provider vào bảng providers.
```

#### Khi phát hiện bug pattern mới:
```
Thêm vào ai/memory/risks.md — mô tả risk mới phát hiện, 
mức độ nghiêm trọng, file liên quan, và workaround.
```

#### Khi quyết định architecture quan trọng:
```
Thêm decision mới vào ai/memory/decisions.md theo format D-0XX hiện có.
```

#### Khi thêm storyboard mode mới:
```
Cập nhật ai/systems/storyboard-system.md — thêm mode mới vào bảng dispatch 
và bảng prompt templates.
Cập nhật ai/skills/new-storyboard-mode.md — thêm mode vào bảng "Existing Modes for Reference".
```

#### Khi thay đổi conventions:
```
Cập nhật ai/memory/conventions.md và ai/memory/coding-style.md 
với convention mới.
```

### 3C. Bảo trì định kỳ (khuyến nghị mỗi 1-2 tuần)

Gửi prompt này:

```
Audit hệ thống ai/ cho tôi:
1. Kiểm tra ai/indexes/feature-map.md — có feature nào mới chưa được thêm không?
2. Kiểm tra ai/memory/risks.md — có risk nào đã fix hoặc risk mới chưa document?
3. Kiểm tra ai/indexes/dependency-map.md — có dependency mới chưa?
4. Kiểm tra ai/systems/*.md — có subsystem nào thay đổi đáng kể chưa cập nhật?
Liệt kê những gì cần update, rồi update.
```

### 3D. Khi thêm loại file hoàn toàn mới vào ai/

```
Tạo skill mới ai/skills/[tên-skill].md cho [mô tả capability].
Sau đó cập nhật AGENTS.md System Map nếu cần,
và cập nhật router liên quan để reference đến skill mới.
```

---

## Phần 4: Tóm Tắt Quy Trình

```
┌─────────────────────────────────────────────────┐
│                 VÒNG LẶP LÀM VIỆC               │
│                                                   │
│  1. Gửi yêu cầu bình thường                     │
│     (Agent tự tham khảo AGENTS.md)               │
│                                                   │
│  2. Nếu task phức tạp → gợi ý đọc file cụ thể   │
│     "Tham khảo ai/skills/xxx.md"                 │
│                                                   │
│  3. Sau khi hoàn thành task lớn → yêu cầu update │
│     "Cập nhật ai/ files liên quan"               │
│                                                   │
│  4. Mỗi 1-2 tuần → chạy audit                   │
│     "Audit hệ thống ai/ cho tôi"                │
│                                                   │
└─────────────────────────────────────────────────┘
```

---

## Phần 5: Cheat Sheet — Prompt Templates

### Phát triển feature
```
[Mô tả feature]. Tham khảo conventions trước khi code.
Sau khi xong, update ai/indexes/feature-map.md.
```

### Fix bug
```
[Mô tả lỗi]. Kiểm tra ai/memory/risks.md trước.
Nếu là bug pattern mới, thêm vào risks.md sau khi fix.
```

### Thêm provider
```
Thêm [provider] vào [image/video] generation.
Update dependency-map.md và [image/video]-pipeline.md sau khi xong.
```

### Thêm storyboard mode
```
Thêm chế độ [tên] cho storyboard. 
Theo quy trình ai/skills/new-storyboard-mode.md.
Update storyboard-system.md sau khi xong.
```

### Refactor
```
Refactor [mô tả]. Tuân thủ ai/rules/forbidden-patterns.md.
```

### Audit
```
Audit ai/ system — kiểm tra feature-map, risks, dependencies có up-to-date không.
```
