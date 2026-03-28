# 🔬 Bộ Prompt Phân Tích Video — Universal Edition

> [!IMPORTANT]
> Bộ prompt này có thể phân tích **BẤT KỲ** thể loại video nào — từ documentary, vlog, music video, animation, đến gaming, cooking, tech review, v.v.
> Mỗi prompt dưới đây phục vụ một mục đích cụ thể cho template. Hãy dùng chúng **theo thứ tự** — kết quả prompt trước sẽ bổ trợ cho prompt sau. Cung cấp video/screenshots kèm prompt cho AI phân tích.

> [!TIP]
> **Cách dùng tốt nhất:** Điền `[TÊN KÊNH / VIDEO]` bằng tên kênh hoặc video bạn muốn phân tích. Chọn 3-5 video đại diện (khác chủ đề nếu cùng kênh). Upload video hoặc chụp 20-30 screenshots đại diện từ mỗi video, rồi gửi kèm prompt bên dưới.

> [!NOTE]
> **Hướng dẫn thay thế placeholder:** Trước khi dùng, thay tất cả `[TÊN KÊNH / VIDEO]` bằng tên kênh/video cụ thể. Các mục trong `{lựa chọn A | lựa chọn B | ...}` cho bạn tham khảo — AI sẽ tự xác định đáp án phù hợp.

---

## Prompt 1: 🎨 PHÂN TÍCH DNA HÌNH ẢNH (Visual DNA)

> **Mục đích:** Thu thập dữ liệu chính xác cho `style_prompt` & tất cả `image_*` prompts
> **Input:** Upload 20-30 screenshots đại diện từ 2-3 video (bao gồm: đa dạng cỡ cảnh, ánh sáng, nội dung)

```
Bạn là chuyên gia phân tích ngôn ngữ hình ảnh (visual language analyst). Tôi cung cấp các frame/screenshots từ kênh/video "[TÊN KÊNH / VIDEO]." Nhiệm vụ của bạn là phân tích CỰC KỲ CHI TIẾT về visual DNA để tôi có thể tái tạo chính xác phong cách này bằng AI image generation.

Hãy phân tích theo đúng cấu trúc sau, mỗi mục phải CỤ THỂ, ĐỊNH LƯỢNG, có ví dụ frame cụ thể. Nếu một mục không áp dụng cho thể loại video này, hãy ghi "N/A" và giải thích tại sao.

## 1. KỸ THUẬT TẠO HÌNH ẢNH (Image/Video Production Method)
- Hình ảnh/video được tạo bằng phương pháp gì? Phân loại chính xác:
  { Quay thật (live-action footage) | AI-generated stills | Digital painting/illustration | 2D animation | 3D render/CGI | Motion graphics | Screen recording | Composite/mixed-media | Stock footage | Khác }
- Nếu quay thật: camera loại gì? (DSLR/mirrorless? smartphone? action cam? drone? webcam?) Sensor size ước tính?
- Nếu AI-generated: model nào có khả năng cao nhất? (Midjourney, DALL-E, Stable Diffusion, Flux, Sora, Kling?)
- Nếu animation: phong cách gì? (anime, cartoon, motion graphics, whiteboard, stop-motion?)
- Mức độ photorealistic từ 1-10? (1 = highly stylized/cartoon, 10 = indistinguishable from real footage)
- Có dấu hiệu nào của AI artifacts không? (ngón tay lạ, text lỗi, texture lặp?)
- Hình ảnh có tĩnh 100% hay có animation/movement? (parallax, zoom chậm, particle effects, real motion?)
- Có face-cam / talking head không? Chiếm bao nhiêu % screentime?

## 2. MÀU SẮC CHÍNH XÁC (Exact Color Science)
Cho mỗi nhóm bối cảnh THỰC TẾ XUẤT HIỆN trong video, hãy xác định MÃ MÀU GẦN NHẤT (hex):

### 2a. Bảng màu Shadow/Dark zones:
- Màu shadow chính (hex): ?
- Màu shadow phụ (hex): ?
- Blacks có bị lift không? Nếu có, lift bao nhiêu (ước tính giá trị IRE hoặc 0-255)?

### 2b. Bảng màu Highlight/Bright zones:
- Màu highlight chính (hex): ?
- Highlight có glow/bloom không? Bán kính bloom ước tính?
- Highlight clipping hay roll-off mềm?

### 2c. Bảng màu Midtone:
- Skin tones (nếu có người): warm (hex)? cool (hex)?
- Dominant midtone surfaces (hex): ? (tùy thể loại: nội thất, thiên nhiên, UI, sản phẩm, v.v.)

### 2d. Bảng màu Accent / Brand Colors:
- Màu accent chính (hex): ?
- Màu accent phụ (hex): ?
- Có consistent brand color palette không? Liệt kê.

### 2e. Tỷ lệ màu trong frame:
- % diện tích tối (shadow): ?
- % diện tích trung tính (midtone): ?
- % diện tích sáng (highlight): ?
- Tổng thể frame thiên tối (low-key), trung tính, hay sáng (high-key)?
- Color grading tổng thể: warm? cool? neutral? desaturated? vibrant? split-toned?

## 3. ÁNH SÁNG CHÍNH XÁC (Exact Lighting)
Phân tích theo các loại cảnh THỰC TẾ XUẤT HIỆN trong video:

### 3a. Setup ánh sáng chủ đạo:
- Ánh sáng tự nhiên hay nhân tạo? Hay kết hợp?
- Key light: hướng? Cường độ tương đối? Hard/soft?
- Fill light: có không? Tỷ lệ key:fill ratio ước tính? (2:1? 4:1? 8:1?)
- Rim/back light: có không? Mạnh/yếu? Màu gì?
- Có dùng colored lighting không? (neon, RGB, gels?)

### 3b. Biến thể ánh sáng theo context (liệt kê các tình huống thực tế):
(VD: outdoor/indoor, ngày/đêm, studio/location, screen-lit, v.v.)
- Mỗi context: nguồn sáng, hướng, mood, contrast level

### 3c. Hiệu ứng ánh sáng đặc biệt:
- Volumetric light / god rays? Mức độ?
- Lens flare? Có/không? Style?
- Atmospheric perspective (xa bị mờ)? Mức nào?
- Bóng đổ: hard edge / soft edge? Direction consistency?

## 4. TEXTURE & CHẤT LIỆU (Visual Surface Quality)
Mô tả CHÍNH XÁC chất lượng bề mặt trong video, tập trung vào các elements THỰC TẾ XUẤT HIỆN:

### 4a. Tổng thể chất lượng hình ảnh:
- Resolution ước tính (720p/1080p/4K)?
- Sharpness: razor sharp / clean / slightly soft / intentionally soft?
- Compression artifacts visible? Mức độ?

### 4b. Các bề mặt/vật liệu chính xuất hiện trong video (liệt kê 5-8 loại):
Cho mỗi loại, mô tả:
- Level of detail (1-10)?
- Texture rendering: photorealistic / stylized / flat / painted?
- Weathering/aging: có/không? Mức độ?

### 4c. Con người (nếu xuất hiện):
- Skin rendering: pore-level detail / smooth / stylized?
- Có subsurface scattering (da hơi trong) không?
- Hair rendering: individual strands / clumps / stylized?
- Chi tiết bẩn/mồ hôi/makeup: mức nào?

## 5. FILM GRAIN & POST-PROCESSING
- Film grain có không? Nếu có: mịn/thô? Mức intensity ước tính (1-10)?
- Grain đồng đều hay nhiều hơn ở shadow?
- Chromatic aberration có không?
- Vignette: có không? Mức độ?
- Lens distortion nhận thấy không? (barrel/pincushion/fisheye?)
- Bloom/glow quanh nguồn sáng: có không? Mức độ? Màu?
- Sharpening: quá sharp / vừa / hơi soft?
- Noise pattern: digital noise hay film grain?
- Depth of field: deep (all sharp) / shallow (bokeh)?
- Có letterboxing (black bars) không? Tỷ lệ aspect ratio?

## 6. SO SÁNH VỚI CÁC STYLE QUEN THUỘC
- Giống phim/show nào nhất? Liệt kê 3-5 references gần nhất.
- Giống kênh YouTube nào nhất về mặt hình ảnh? Liệt kê 3-5 kênh.
- Giống phong cách nhiếp ảnh / nghệ thuật nào? (nếu áp dụng)
- Nếu là AI art: giống AI art style nào nhất?

## 7. SAMPLE PROMPT RECONSTRUCTION
Chọn 3 frame ĐẠI DIỆN NHẤT (khác nhau về nội dung/mood) và viết prompt tái tạo CHÍNH XÁC cho từng frame:

Mỗi prompt phải đủ chi tiết để AI image generator tái tạo được frame gần giống nhất. Format:
{
  "frame_description": "...(mô tả ngắn frame)...",
  "full_prompt": "...(200+ words)...",
  "negative_prompt": "...",
  "recommended_model": "...",
  "recommended_settings": {"cfg_scale": ?, "steps": ?, "sampler": "?"}
}

Trả lời bằng TIẾNG VIỆT nhưng các prompt tái tạo viết bằng TIẾNG ANH.
```

---

## Prompt 2: 📐 PHÂN TÍCH BỐ CỤC & CAMERA (Composition & Cinematography)

> **Mục đích:** Thu thập dữ liệu cho `storyboard_breakdown` & `image_*` prompts
> **Input:** Upload 30-40 screenshots theo thứ tự thời gian từ 1 video (mỗi 10-15 giây chụp 1 frame)

```
Bạn là đạo diễn hình ảnh (Director of Photography - DP) chuyên nghiệp. Tôi cung cấp các frame liên tục từ một video của kênh/video "[TÊN KÊNH / VIDEO]." Hãy phân tích CHÍNH XÁC ngôn ngữ điện ảnh/hình ảnh được sử dụng.

Lưu ý: Hãy chỉ phân tích CÁC YẾU TỐ THỰC SỰ XUẤT HIỆN trong video. Nếu video không có một yếu tố nào (VD: không có camera movement vì là slideshow, hoặc chỉ có một góc quay cố định vì là talking-head), hãy ghi rõ và giải thích thay vì bỏ trống.

## 1. THỐNG KÊ CỠ CẢNH (Shot Size Distribution)
Phân loại TỪNG frame vào một cỡ cảnh, sau đó tổng hợp:

| Cỡ cảnh | Số lần | % tổng | Dùng khi nào (context) |
|----------|--------|--------|------------------------|
| Extreme Wide Shot (EWS) | ? | ? | ? |
| Wide Shot (WS) | ? | ? | ? |
| Medium Wide (MWS) | ? | ? | ? |
| Medium Shot (MS) | ? | ? | ? |
| Medium Close-Up (MCU) | ? | ? | ? |
| Close-Up (CU) | ? | ? | ? |
| Extreme Close-Up (ECU) | ? | ? | ? |
| Insert/Detail Shot | ? | ? | ? |
| Screen Recording / UI | ? | ? | ? |
| Text/Title Card | ? | ? | ? |
| B-Roll / Cutaway | ? | ? | ? |

(Bỏ qua các dòng có số lần = 0)

## 2. THỐNG KÊ GÓC MÁY (Camera Angle)
| Góc máy | Số lần | % | Dùng khi nào |
|---------|--------|---|-------------|
| Eye-level | ? | ? | ? |
| Low angle (nhìn lên) | ? | ? | ? |
| High angle (nhìn xuống) | ? | ? | ? |
| Bird's eye / Overhead / Top-down | ? | ? | ? |
| Dutch angle (nghiêng) | ? | ? | ? |
| Worm's eye (sát đất nhìn lên) | ? | ? | ? |
| POV (first person) | ? | ? | ? |
| Selfie angle | ? | ? | ? |

(Bỏ qua các dòng có số lần = 0)

## 3. CHUYỂN ĐỘNG CAMERA (Camera Movement)
Xác định phương pháp camera chính: {Quay thật cầm tay | Quay thật trên gimbal/tripod | Ảnh tĩnh + hiệu ứng camera giả (post-production) | Drone | Animation camera | Screencast | Kết hợp}

| Loại chuyển động | Tần suất | Tốc độ | Mô tả chi tiết |
|-------------------|----------|--------|-----------------|
| Static (hoàn toàn tĩnh) | ? | - | ? |
| Slow zoom in (Ken Burns) | ? | ước tính ? | ? |
| Slow zoom out | ? | ước tính ? | ? |
| Pan left/right | ? | ? | ? |
| Tilt up/down | ? | ? | ? |
| Tracking/dolly | ? | ? | ? |
| Handheld (organic movement) | ? | ? | ? |
| Parallax (multi-layer) | ? | ? | ? |
| Drone / aerial movement | ? | ? | ? |
| Digital zoom (post-production) | ? | ? | ? |

- Thời gian trung bình mỗi shot (giây)?
- Có dùng transition effects không? (dissolve, cut, fade to black, wipe, swipe, zoom transition?)
- Transitions đặc trưng nhất của kênh/video?

## 4. QUY LUẬT BỐ CỤC (Composition Rules)
Cho từng frame, xác định:
- Vị trí chủ thể (subject placement): center? rule-of-thirds? golden ratio? full-frame?
- Leading lines: có không? Loại gì?
- Foreground elements: có không? Loại gì?
- Negative space: nhiều/ít? Hướng nào?
- Depth layers: có bao nhiêu layer depth rõ ràng? (FG, MG, BG?)
- Framing devices: chủ thể có bị "đóng khung" bởi gì không?
- Text/graphic overlays: có không? Vị trí? Style?

Tổng hợp thành quy luật bố cục chung (top 5 patterns hay dùng nhất).

## 5. NHỊP ĐỘ HÌNH ẢNH (Visual Pacing)
- Thời lượng trung bình mỗi shot (ước tính từ các frame)?
- Shot ngắn nhất & dài nhất ước tính?
- Nhịp thay đổi shot: đều đặn hay thay đổi theo nội dung/cảm xúc?
- Khi nội dung intense/nhanh → shot nhanh hay chậm?
- Khi nội dung chậm/reflection → shot được giữ bao lâu?
- Có hiệu ứng đặc biệt nào khi transition giữa các "phần" không?

## 6. MẪU CẮT CẢNH (Editing Patterns)
Xác định patterns lặp lại:
- Có tuân theo pattern establishing → detail → reaction không?
- Shot/reverse shot? Có dùng cho đoạn nào?
- Montage sequences? Có đoạn nào cắt nhanh liên tiếp không?
- Match cuts? Smash cuts?
- J-cut / L-cut (audio leads/trails the visual)?
- Các pattern đặc trưng khác của kênh/video?

Trả lời bằng TIẾNG VIỆT.
```

---

## Prompt 3: 🗣️ PHÂN TÍCH GIỌNG KỂ & CẤU TRÚC (Narration/Audio Style & Script Structure)

> **Mục đích:** Thu thập dữ liệu cho `script_outline`, `script_episode`
> **Input:** Transcript/phụ đề từ 2-3 video khác chủ đề (có thể dùng YouTube auto-subtitle hoặc transcript tool)

```
Bạn là nhà phân tích ngôn ngữ học (linguistic analyst) chuyên phân tích phong cách kể chuyện và trình bày nội dung. Tôi cung cấp transcript từ kênh/video "[TÊN KÊNH / VIDEO]." Hãy phân tích CHÍNH XÁC phong cách narration/presenting để tôi có thể tái tạo bằng AI writing.

Lưu ý: Hãy phân tích dựa trên THỰC TẾ nội dung video. Không giả định rằng video là documentary — nó có thể là bất kỳ thể loại nào (tutorial, review, vlog, essay, entertainment, v.v.)

## 1. CẤU TRÚC VĨ MÔ (Macro Structure)
Phân tích cấu trúc tổng thể của MỖI video:
- Video thuộc thể loại nội dung gì? { Documentary/essay | Tutorial/how-to | Review/opinion | Vlog/personal | Entertainment/sketch | News/commentary | Educational/lecture | Story-driven narrative | List/ranking | Interview | Khác }
- Video mở đầu bằng gì? (hook question? dramatic scene? context? personal anecdote? cold open? title card? meme?)
- Có cold open không? Bao lâu?
- Các "phần" (sections/chapters) trong video được chia như nào?
- Có bao nhiêu "tầng" narrative? (main content, tangents, recurring segments?)
- Kết thúc bằng gì? (summary? CTA? cliffhanger? joke? reflection? open question?)
- Có intro/outro music/animation cố định không?
- Có sponsor segment không? Đặt ở đâu? Transition vào/ra sponsor như nào?
- Tổng thời lượng narration/talking vs nhạc/silence/other?

Vẽ sơ đồ cấu trúc cho mỗi video:
```
00:00-??:?? → [Mô tả phần]
??:??-??:?? → [Mô tả phần]
...
```

## 2. PHONG CÁCH NGÔN NGỮ (Linguistic Style)
### 2a. Ngôi kể / Phong cách nói (Point of View & Voice):
- Ngôi thứ mấy? (first person "tôi/I"? third person? second person "bạn/you"? chuyển đổi?)
- Single narrator hay nhiều voices?
- Narrator/presenter có "personality" rõ ràng không? (neutral / opinionated / humorous / dramatic / casual / academic / energetic?)
- Có catchphrases hoặc verbal tics không? VD?

### 2b. Vocabulary & Sentence Pattern:
- Độ dài câu trung bình (ước tính số từ)?
- Câu ngắn nhất & dài nhất?
- Có dùng câu đơn giật gọn cho dramatic effect không? VD?
- Từ vựng level: { casual/slang | phổ thông | semi-formal | academic/literary | chuyên ngành }
- Có dùng thuật ngữ chuyên biệt không? (tech jargon? academic terms? industry lingo?) VD?
- Metaphor/simile: tần suất? Kiểu ẩn dụ gì? VD cụ thể?
- Có humor không? Kiểu gì? (wordplay? sarcasm? absurdist? self-deprecating? reference-based?)

### 2c. Rhetorical Devices (Thủ pháp tu từ):
- Câu hỏi tu từ (rhetorical questions)? Tần suất? VD?
- Dramatic irony? VD?
- Foreshadowing (báo trước)? VD?
- Repetition/Anaphora (lặp cấu trúc)? VD?
- Juxtaposition (đối lập)? VD?
- Direct address (nói thẳng với khán giả)? VD?
- Callback (tham chiếu lại điều đã nói trước đó)? VD?
- Các thủ pháp đặc trưng khác?

### 2d. Transition Phrases (Cụm từ chuyển ý):
Liệt kê TẤT CẢ các cụm từ chuyển ý đặc trưng mà narrator hay dùng.
Liệt kê ít nhất 15-20 cụm từ, phân loại theo chức năng:
- Chuyển sang phần mới / chủ đề mới: ?
- Chuyển sang hồi hộp/căng thẳng: ?
- Chuyển sang giải thích/context: ?
- Chuyển sang twist/bất ngờ: ?
- Chuyển sang kết luận/reflection: ?
- Chuyển thời gian/không gian: ?
- Chuyển từ tangent trở về main topic: ?

## 3. NHỊP ĐỘ KỂ (Narration/Presentation Pacing)
- Tốc độ nói trung bình (từ/phút)?
- Có đoạn nói nhanh không? Khi nào?
- Có đoạn nói CHẬM/nghỉ dài không? Khi nào?
- Silence/music-only/SFX-only moments: có không? Bao lâu? Dùng khi nào?
- Tỷ lệ narration vs silence vs music trong 1 phút?
- Có sự thay đổi energy level rõ ràng không? Pattern?

## 4. CẢM XÚC & TONE (Emotional Arc)
Cho MỖI video, vẽ emotional arc:
- Timestamp → Emotion/Energy → Intensity (1-5) → Mô tả
- VD: "02:30 → Curiosity → 3 → Đặt câu hỏi gây tò mò"

Xác định:
- Tone chủ đạo là gì? (serious? playful? dark? inspirational? neutral? passionate? contemplative?)
- Có humor/sarcasm không? Tần suất?
- Narrator có bias/opinion rõ ràng không?
- Cách xử lý nội dung nhạy cảm (nếu có): explicit hay implicit?
- Mức độ personal/emotional: detached / balanced / deeply personal?

## 5. SO SÁNH VỚI CÁC KÊNH/CREATOR QUEN BIẾT
So sánh phong cách với 5-6 kênh/creator CÓ PHONG CÁCH GẦN NHẤT (cùng thể loại nội dung):

| Khía cạnh | [TÊN KÊNH / VIDEO] | Kênh tương đồng nhất |
|-----------|---------------------|---------------------|
| Tone | ? | ? |
| Pacing | ? | ? |
| Detail level | ? | ? |
| Emotional involvement | ? | ? |
| Production quality | ? | ? |

## 6. SAMPLE SCRIPT RECREATION
Viết lại 1 đoạn script (200-300 từ) bắt chước CHÍNH XÁC phong cách này, về một chủ đề KHÁC (nhưng cùng thể loại nội dung). Đoạn này phải:
- Dùng đúng ngôi kể / voice
- Dùng đúng cấu trúc câu
- Dùng đúng loại transition phrases
- Dùng đúng level vocabulary & humor (nếu có)
- Có đúng nhịp độ cảm xúc / energy

Trả lời bằng TIẾNG VIỆT (trừ sample script viết bằng NGÔN NGỮ GỐC CỦA VIDEO).
```

---

## Prompt 4: 🏛️ PHÂN TÍCH THIẾT KẾ HÌNH ẢNH & SẢN XUẤT (Visual Design & Production)

> **Mục đích:** Thu thập dữ liệu cho `character_extraction`, `scene_extraction`, `prop_extraction`
> **Input:** 15-20 screenshots tập trung vào con người, vật thể, bối cảnh, đồ họa

```
Bạn là chuyên gia thiết kế sản xuất (Production Designer) và art director. Tôi cung cấp screenshots từ kênh/video "[TÊN KÊNH / VIDEO]." Hãy phân tích chi tiết về thiết kế hình ảnh để tôi tái tạo bằng AI.

Lưu ý: Hãy chỉ phân tích CÁC YẾU TỐ THỰC SỰ XUẤT HIỆN trong video. Bỏ qua các section không áp dụng.

## 1. CON NGƯỜI / NHÂN VẬT (People & Characters)

### 1a. Phong cách thể hiện con người:
- Có người thật (live-action) không? Hay AI-generated? Hay illustration/animation?
- Photorealistic 100% hay có stylization? Mức độ?
- Có face-cam / presenter on-screen không?

### 1b. Trang phục & Ngoại hình (theo nhóm người xuất hiện):
Mô tả CHI TIẾT cho mỗi nhóm THỰC TẾ XUẤT HIỆN:

**Nhóm 1: [Tên nhóm - VD: Presenter, Nhân vật chính, Background people...]**
- Trang phục: chất liệu, màu sắc, style?
- Mức độ detail / weathering?
- Phụ kiện?

**Nhóm 2: [Tên nhóm]**
- (tương tự)

(Thêm nhóm nếu cần)

### 1c. Biểu cảm & Body Language:
- Biểu cảm phổ biến nhất?
- Ánh mắt: nhìn vào camera? nhìn xa? nhìn sản phẩm/đối tượng?
- Gestures: có/không? Kiểu gì?
- Posing: natural/candid? staged/posed? action/dynamic?

## 2. BỐI CẢNH & MÔI TRƯỜNG (Settings & Environments)
Phân loại và mô tả CHI TIẾT từng loại bối cảnh THỰC TẾ XUẤT HIỆN:

### 2a. Bối cảnh chính (Primary Setting):
- Mô tả tổng quát: loại không gian? (studio, phòng, ngoài trời, virtual, composite?)
- Kiến trúc / layout: phong cách? Kích thước? Proportions?
- Ánh sáng: nguồn gì? Hướng? Cường độ?
- Đồ nội thất / props chính?
- Background: rõ ràng / blurred / virtual / green-screen?
- Color palette của bối cảnh?

### 2b. Bối cảnh phụ / B-roll (nếu có):
- Liệt kê các loại bối cảnh phụ
- Mỗi loại: mô tả ngắn + mục đích sử dụng

### 2c. Tính nhất quán môi trường:
- Bối cảnh có nhất quán xuyên suốt video/kênh không?
- Có thay đổi theo chủ đề/season không?

## 3. ĐỒ VẬT & ĐẠO CỤ (Objects & Props)
Liệt kê các đồ vật/đạo cụ QUAN TRỌNG xuất hiện, cho mỗi cái:
- Tên + loại
- Mức độ chi tiết (1-10)
- Chất liệu / trạng thái
- Rendering style (photorealistic? stylized? illustrated?)
- Vai trò trong video (functional? decorative? subject matter?)
- Có xuất hiện là cận cảnh riêng hay chỉ trong bối cảnh?

## 4. TYPOGRAPHY & ĐỒ HỌA ON-SCREEN
### 4a. Text elements:
- Có text on-screen không? Liệt kê các loại (title, subtitle, label, caption, annotation, lower third?)
- Font chữ: kiểu gì? (serif/sans-serif/display/handwritten/monospace?)
- Màu text? Có shadow/outline/background box?
- Vị trí text thường nằm ở đâu?
- Animation: text có animated không? Kiểu gì? (fade in, type-on, slide, bounce, glitch?)

### 4b. Graphic elements:
- Có bản đồ, diagram, chart, infographic không? Style?
- Có emoji, sticker, icon overlay không?
- Có meme, reaction image, screenshot insert không?
- Lower thirds / name plates style?
- Thumbnail style: mô tả pattern

### 4c. Brand elements:
- Logo / watermark? Vị trí? Style?
- Consistent graphic identity (màu, font, layout)?
- Intro/outro animation style?

## 5. OUTPUT: IMAGE GENERATION REFERENCE SHEET
Tổng hợp thành "cheat sheet" cho AI image generation:

```json
{
  "content_type": "...(thể loại video)...",
  "style_keywords": ["keyword1", "keyword2", "..."],
  "negative_keywords": ["avoid1", "avoid2", "..."],
  "color_palette": {
    "shadow_primary": "#??????",
    "shadow_secondary": "#??????",
    "highlight_warm": "#??????",
    "highlight_cool": "#??????",
    "accent_primary": "#??????",
    "accent_secondary": "#??????",
    "brand_colors": ["#??????", "#??????"],
    "dominant_midtones": ["#??????", "#??????"]
  },
  "lighting_setup": {
    "primary_method": "...(natural/studio/mixed/virtual)...",
    "key_fill_ratio": "?:?",
    "light_sources": ["?", "?"],
    "mood": "?"
  },
  "texture_detail_level": "?/10",
  "film_grain_intensity": "?/10",
  "realism_level": "?/10",
  "production_value": "?/10"
}
```

Trả lời bằng TIẾNG VIỆT, riêng cheat sheet JSON giữ tiếng Anh.
```

---

## Prompt 5: 🎬 PHÂN TÍCH CHUYỂN ĐỘNG & VIDEO (Motion & Animation)

> **Mục đích:** Thu thập dữ liệu cho `video_constraint` & pacing trong `storyboard_breakdown`
> **Input:** Upload trực tiếp 1-2 video clips (30-60 giây mỗi clip, chọn đoạn có nhiều chuyển cảnh)

```
Bạn là chuyên gia motion graphics và post-production / video editor. Tôi cung cấp video clips từ kênh/video "[TÊN KÊNH / VIDEO]." Hãy phân tích CHÍNH XÁC cách video được sản xuất về mặt chuyển động và kỹ thuật.

Lưu ý: Video có thể là bất kỳ thể loại nào. Hãy phân tích DỰA TRÊN THỰC TẾ, không giả định phương pháp sản xuất.

## 1. PHƯƠNG PHÁP SẢN XUẤT VIDEO (Video Production Method)
Xác định phương pháp chính:
- { Live-action footage (quay thật) | Ảnh tĩnh được animate (still images + post-production) | AI-generated video | Screen recording + voiceover | Full animation (2D/3D) | Motion graphics | Mixed media | Khác }
- Nếu kết hợp nhiều phương pháp: tỷ lệ mỗi phương pháp?
- Phần mềm/tools ước đoán được dùng?

## 2. CÁC KIỂU CHUYỂN ĐỘNG (Motion Types)
Phân tích TẤT CẢ các loại chuyển động THỰC TẾ XUẤT HIỆN:

### 2a. Camera Movement (thật hoặc giả):
- Tốc độ: rất chậm / chậm / trung bình / nhanh?
- Zoom: in/out, từ bao nhiêu % đến bao nhiêu %?
- Pan/tilt: range? speed?
- Tracking/dolly/orbit?
- Easing: linear / ease-in / ease-out / ease-in-out?
- Thời gian trung bình mỗi camera move (giây)?

### 2b. Subject Movement (nếu có):
- Người/nhân vật: static / subtle movement / full motion?
- Vật thể: static / animated / physics-based?
- Text/graphics: static / animated? Kiểu animation?

### 2c. Parallax / Depth Effects (nếu có):
- Có dùng parallax (tách layer foreground/background) không?
- Bao nhiêu layers?
- Hiệu ứng depth-of-field giả?

### 2d. Particle & Atmospheric Effects:
(Liệt kê chỉ những gì THỰC SỰ XUẤT HIỆN)
- Bụi bay (dust motes)? Tàn lửa (embers)? Tuyết/mưa?
- Khói/sương (fog/mist)? Animated hay static?
- Ánh sáng nhấp nháy (flickering)?
- Bokeh / light particles?
- Glitch effects?
- Các hiệu ứng đặc biệt khác?

### 2e. Speed Effects:
- Slow motion? Khi nào? Mức chậm?
- Time-lapse? Khi nào?
- Speed ramp (thay đổi tốc độ liên tục)?
- Freeze frame?

## 3. TRANSITIONS (Chuyển cảnh)
Cho MỖI kiểu chuyển cảnh xuất hiện trong video:

| Loại | Tần suất | Duration (ms) | Dùng khi nào |
|------|----------|---------------|-------------|
| Hard cut | ? | - | ? |
| Cross dissolve | ? | ? | ? |
| Fade to black | ? | ? | ? |
| Fade from black | ? | ? | ? |
| Swipe / Wipe | ? | ? | ? |
| Zoom transition | ? | ? | ? |
| Whip pan | ? | ? | ? |
| Match cut | ? | ? | ? |
| Morph / smooth transition | ? | ? | ? |
| Custom / unique (mô tả) | ? | ? | ? |

(Bỏ qua các loại không xuất hiện)

## 4. NHẠC & ÂM THANH (Audio Design)
### 4a. Nhạc nền:
- Thể loại: { orchestral | ambient | electronic | lo-fi | hip-hop | pop | acoustic | cinematic | royalty-free generic | custom | không có nhạc }
- Instruments chính?
- Mood nhạc: consistent hay thay đổi theo nội dung?
- Nhạc có "drop" / build-up / climax align với visual không?

### 4b. Sound Effects:
- Có SFX không? Loại gì? (whoosh, click, ding, environmental, foley?)
- SFX liên tục hay chỉ ở điểm nhấn?
- SFX có sync với visual transitions không?

### 4c. Voice:
- Voice quality: professional VO / casual / raw / processed (reverb, EQ)?
- Tỷ lệ voice : music : SFX?
- Có voice modulation (thay đổi tone/pitch cho effect) không?

## 5. TIMELINE BREAKDOWN
Phân tích 30 giây video thành từng shot:

| Thời gian | Duration | Nội dung | Cỡ cảnh | Camera move | Transition vào | Transition ra | FX/Overlay |
|-----------|----------|----------|---------|-------------|----------------|---------------|------------|
| 00:00 | ?s | ? | ? | ? | ? | ? | ? |
| 00:?? | ?s | ? | ? | ? | ? | ? | ? |
...

## 6. VIDEO GENERATION / RECREATION PARAMETERS
Nếu muốn tái tạo hiệu ứng video này:
- Phần mềm (After Effects, Premiere, DaVinci Resolve, CapCut): settings/workflow chính?
- AI video (Runway Gen-3, Kling, Sora, Pika): prompt style?
- Cho 3 ví dụ video prompt dựa trên các shot trong video:

```json
{
  "shot_description": "?",
  "video_prompt": "?(100+ words)?",
  "duration": "?s",
  "camera_motion": "?",
  "style_keywords": ["?", "?"]
}
```

Trả lời bằng TIẾNG VIỆT, video prompts viết bằng TIẾNG ANH.
```

---

## Prompt 6: 🔄 META-PROMPT — TỔNG HỢP TẤT CẢ (Final Synthesis)

> **Mục đích:** Sau khi có kết quả từ Prompt 1-5, dùng prompt này để đúc kết thành template hoàn chỉnh
> **Input:** Kết quả phân tích từ Prompt 1-5 + link video gốc

```
Bạn đã phân tích kênh/video "[TÊN KÊNH / VIDEO]" với các kết quả chi tiết sau:

[DÁN KẾT QUẢ PROMPT 1 VÀO ĐÂY]
[DÁN KẾT QUẢ PROMPT 2 VÀO ĐÂY]
[DÁN KẾT QUẢ PROMPT 3 VÀO ĐÂY]
[DÁN KẾT QUẢ PROMPT 4 VÀO ĐÂY]
[DÁN KẾT QUẢ PROMPT 5 VÀO ĐÂY]

Bây giờ, hãy tổng hợp TẤT CẢ thông tin trên thành một bảng tham chiếu tối ưu (Master Reference Sheet) với cấu trúc sau. Nếu một section không áp dụng cho thể loại video này, hãy ghi rõ "N/A — [lý do]" thay vì bỏ trống.

## A. STYLE DNA SUMMARY (1 đoạn văn 100-150 từ, tiếng Anh)
Tóm tắt visual + content identity của kênh/video trong 1 đoạn duy nhất, đủ ngắn để dùng làm prefix cho mọi AI prompt. Bao gồm:
- Thể loại nội dung
- Phong cách hình ảnh
- Tone & mood
- Đặc điểm nổi bật nhất

## B. IMAGE GENERATION MASTER PROMPT (tiếng Anh)
Viết 1 prompt master (300-500 từ) mô tả phong cách hình ảnh, có thể dùng trực tiếp kèm bất kỳ scene description nào để tạo ảnh đúng phong cách.

## C. NEGATIVE PROMPT (tiếng Anh)
Liệt kê tất cả thứ KHÔNG ĐƯỢC xuất hiện để giữ đúng phong cách.

## D. EXACT COLOR PALETTE
```json
{
  "shadows": ["#hex1", "#hex2"],
  "midtones": ["#hex1", "#hex2"],
  "highlights": ["#hex1", "#hex2"],
  "accents": ["#hex1", "#hex2"],
  "brand_colors": ["#hex1", "#hex2"],
  "skin_warm": "#hex (or N/A)",
  "skin_cool": "#hex (or N/A)"
}
```

## E. NARRATION / PRESENTATION STYLE GUIDE (tiếng Anh)
Tóm tắt rules cho AI viết narration/script đúng phong cách:
- Content type: ?
- Point of view: ?
- Tone & personality: ?
- Sentence patterns: ?
- Vocabulary level: ?
- Humor style (if applicable): ?
- 10-15 transition phrases hay dùng nhất
- Emotional arc pattern: ?
- Pacing notes: ?

## F. SHOT LANGUAGE RULES (tiếng Anh)
- Shot size distribution (% mỗi loại, chỉ liệt kê loại thực tế dùng)
- Camera angle preferences
- Average shot duration
- Transition types & frequency
- 5 "signature shots/frames" (những cảnh/frame đặc trưng nhất)

## G. MOTION/VIDEO RULES (tiếng Anh)
- Production method
- Animation/motion techniques used
- Speed & easing preferences
- Particle/atmospheric effects
- Transition timing
- Camera movement settings

## H. VISUAL DESIGN REFERENCE (tiếng Anh)
- People/character design rules (if applicable)
- Costume/appearance rules (if applicable)
- Environment/setting rules
- Prop/object styling rules
- Typography & graphic overlay rules
- Brand identity rules

## I. AUDIO DESIGN REFERENCE (tiếng Anh)
- Music genre & mood
- SFX usage patterns
- Voice characteristics
- Audio mixing ratios

Trả lời bằng TIẾNG ANH (toàn bộ — vì output sẽ dùng trực tiếp cho AI prompts).
```

---

## 📋 Checklist Sử Dụng

| # | Prompt | Input cần chuẩn bị | Thời gian ước tính | Dữ liệu thu được |
|---|--------|--------------------|--------------------|-------------------|
| 1 | Visual DNA | 20-30 screenshots đa dạng | 15-20 phút | Color, light, texture specs |
| 2 | Composition & Camera | 30-40 frames liên tục từ 1 video | 15-20 phút | Shot sizes, angles, pacing |
| 3 | Narration/Audio Style | Transcript 2-3 video | 15-20 phút | Script structure, tone, phrases |
| 4 | Visual Design | 15-20 screenshots con người/vật/đồ họa | 15-20 phút | People, scene, prop, graphic specs |
| 5 | Motion & Video | 1-2 video clips (30-60s) | 10-15 phút | Animation, transitions, audio |
| 6 | **Final Synthesis** | Kết quả Prompt 1-5 | 10-15 phút | **Master Reference Sheet** |

> [!CAUTION]
> Prompt 6 (Final Synthesis) là prompt QUAN TRỌNG NHẤT — nó tổng hợp tất cả kết quả thành format có thể dùng trực tiếp cho template. Đừng bỏ qua bước này!

---

## 🎯 Thể Loại Video Đã Được Hỗ Trợ

Bộ prompt này đã được thiết kế để phân tích tất cả các thể loại sau:

| Thể loại | Ví dụ kênh | Những gì prompt sẽ tập trung |
|-----------|------------|------------------------------|
| 📚 Documentary / Essay | Kurzgesagt, Veritasium, Tim – Reborn History | Visual style, narration, research structure |
| 🎮 Gaming | Markiplier, Dream, ibxtoycat | Gameplay footage, commentary, facecam |
| 🍳 Cooking / Food | Binging with Babish, Joshua Weissman | Overhead shots, close-ups, ingredient styling |
| 💻 Tech Review | MKBHD, Linus Tech Tips | Product shots, studio lighting, B-roll |
| 🎵 Music Video | Any artist | Color grading, choreography, editing rhythm |
| 📱 Vlog | Casey Neistat, Emma Chamberlain | Handheld, jump cuts, personal tone |
| 📖 Tutorial / How-to | 3Blue1Brown, Fireship | Screen recording, diagrams, pacing |
| 🎭 Entertainment / Sketch | Dude Perfect, MrBeast | High energy, multiple cameras, graphics |
| 📰 News / Commentary | Philip DeFranco, Last Week Tonight | Talking head, B-roll inserts, graphics |
| 🎨 Art / Animation | Corridor Crew, TheOdd1sOut | Process shots, animation style, tools |
