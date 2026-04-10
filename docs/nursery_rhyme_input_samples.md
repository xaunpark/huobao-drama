# Nursery Rhyme — Mẫu Input cho Tool

> **Mục đích**: Đây là mẫu input hoàn chỉnh để paste vào ô Script Content trong Episode. 
> Dùng file này làm reference khi yêu cầu AI khác tách lời bài hát thành format chuẩn.

---

## Quy tắc Format

```
[SECTION_TYPE N: Subject]
(M:SS – M:SS) [INSTRUMENTAL] Mô tả nếu là đoạn nhạc không lời
(M:SS – M:SS) Lời bài hát dòng 1
(M:SS – M:SS) Lời bài hát dòng 2
```

### Giải thích:
- `[VERSE N: Subject]` — Đánh dấu đoạn verse mới. `N` = số thứ tự, `Subject` = chủ đề/đối tượng chính
- `[CHORUS]` — Đánh dấu điệp khúc (nếu có)
- `[INTRO]` / `[OUTRO]` — Đánh dấu mở đầu/kết thúc
- `(M:SS – M:SS)` — Timestamp bắt đầu và kết thúc (phút:giây). Dùng dấu `–` (en-dash) hoặc `-` (hyphen)
- `[INSTRUMENTAL]` — Tag cho đoạn nhạc không có lời, theo sau là mô tả hình ảnh
- Dòng trống giữa các section header OK, sẽ bị bỏ qua
- Mỗi dòng lời = 1 shot tiềm năng (AI sẽ quyết định merge/split)

---

## Mẫu 1: Narrative Structure — "Wheels on the Bus"

> **Narrative** = Mỗi verse giới thiệu 1 hành động/đối tượng MỚI, không cộng dồn từ verse trước.

```
[INTRO]
(0:00 – 0:04) [INSTRUMENTAL] Logo intro animation
(0:05 – 0:11) [INSTRUMENTAL] Bus driving down a sunny road, establishing the town

[VERSE 1: The Wheels]
(0:12 – 0:17) The wheels on the bus go round and round, round and round, round and round
(0:18 – 0:23) The wheels on the bus go round and round, all through the town

[VERSE 2: The Door]
(0:24 – 0:25) [INSTRUMENTAL] Bus arriving at a stop
(0:26 – 0:31) The door on the bus goes open and shut, open and shut, open and shut
(0:32 – 0:37) The door on the bus goes open and shut, all through the town

[VERSE 3: The Wipers]
(0:38 – 0:39) [INSTRUMENTAL] Rain starts falling
(0:40 – 0:45) The wipers on the bus go swish swish swish, swish swish swish, swish swish swish
(0:46 – 0:51) The wipers on the bus go swish swish swish, all through the town

[VERSE 4: The Horn]
(0:52 – 0:53) [INSTRUMENTAL] Bus approaching an intersection
(0:54 – 0:59) The horn on the bus goes beep beep beep, beep beep beep, beep beep beep
(1:00 – 1:05) The horn on the bus goes beep beep beep, all through the town

[VERSE 5: The People]
(1:06 – 1:07) [INSTRUMENTAL] Camera showing inside the bus
(1:08 – 1:13) The people on the bus go up and down, up and down, up and down
(1:14 – 1:19) The people on the bus go up and down, all through the town

[VERSE 6: The Babies]
(1:20 – 1:25) The babies on the bus go wah wah wah, wah wah wah, wah wah wah
(1:26 – 1:31) The babies on the bus go wah wah wah, all through the town

[VERSE 7: The Mommies]
(1:32 – 1:37) The mommies on the bus go shh shh shh, shh shh shh, shh shh shh
(1:38 – 1:43) The mommies on the bus go shh shh shh, all through the town

[VERSE 8: The Wheels Reprise]
(1:44 – 1:49) The wheels on the bus go round and round, round and round, round and round
(1:50 – 1:55) The wheels on the bus go round and round, all through the town

[OUTRO]
(1:56 – 2:02) [INSTRUMENTAL] Bus drives off into the sunset, logo outro
```

---

## Mẫu 2: Cumulative Structure — "Old MacDonald Had A Farm"

> **Cumulative** = Mỗi verse thêm 1 phần tử MỚI, và lặp lại TẤT CẢ phần tử cũ.

```
[INTRO]
(0:00 – 0:05) [INSTRUMENTAL] Logo intro
(0:06 – 0:14) [INSTRUMENTAL] Acoustic fiddle and banjo playing, establishing the farm scene

[VERSE 1: The Pig]
(0:15 – 0:20) Old MacDonald had a farm, E-I-E-I-O
(0:21 – 0:25) And on that farm he had a pig, E-I-E-I-O
(0:26 – 0:29) With an oink oink here and an oink oink there
(0:30 – 0:34) Here an oink, there an oink, everywhere an oink oink
(0:35 – 0:39) Old MacDonald had a farm, E-I-E-I-O

[VERSE 2: The Duck]
(0:40 – 0:44) Old MacDonald had a farm, E-I-E-I-O
(0:45 – 0:49) And on that farm he had a duck, E-I-E-I-O
(0:50 – 0:53) With a quack quack here and a quack quack there
(0:54 – 0:58) Here a quack, there a quack, everywhere a quack quack
(0:59 – 1:02) Oink oink here and an oink oink there
(1:03 – 1:07) Old MacDonald had a farm, E-I-E-I-O

[VERSE 3: The Horse]
(1:08 – 1:12) Old MacDonald had a farm, E-I-E-I-O
(1:13 – 1:17) And on that farm he had a horse, E-I-E-I-O
(1:18 – 1:21) With a neigh neigh here and a neigh neigh there
(1:22 – 1:26) Here a neigh, there a neigh, everywhere a neigh neigh
(1:27 – 1:30) Quack quack here and a quack quack there
(1:31 – 1:34) Oink oink here and an oink oink there
(1:35 – 1:39) Old MacDonald had a farm, E-I-E-I-O

[VERSE 4: The Cow]
(1:40 – 1:44) Old MacDonald had a farm, E-I-E-I-O
(1:45 – 1:49) And on that farm he had a cow, E-I-E-I-O
(1:50 – 1:53) With a moo moo here and a moo moo there
(1:54 – 1:58) Here a moo, there a moo, everywhere a moo moo
(1:59 – 2:02) Neigh neigh here and a neigh neigh there
(2:03 – 2:06) Quack quack here and a quack quack there
(2:07 – 2:10) Oink oink here and an oink oink there
(2:11 – 2:15) Old MacDonald had a farm, E-I-E-I-O

[VERSE 5: The Sheep]
(2:16 – 2:20) Old MacDonald had a farm, E-I-E-I-O
(2:21 – 2:25) And on that farm he had a sheep, E-I-E-I-O
(2:26 – 2:29) With a baa baa here and a baa baa there
(2:30 – 2:34) Here a baa, there a baa, everywhere a baa baa
(2:35 – 2:38) Moo moo here and a moo moo there
(2:39 – 2:42) Neigh neigh here and a neigh neigh there
(2:43 – 2:46) Quack quack here and a quack quack there
(2:47 – 2:50) Oink oink here and an oink oink there
(2:51 – 2:55) Old MacDonald had a farm, E-I-E-I-O

[OUTRO]
(2:56 – 3:05) [INSTRUMENTAL] All animals gather together, finale celebration
(3:06 – 3:15) [INSTRUMENTAL] Logo outro with all characters waving
```

---

## Prompt để yêu cầu AI khác tạo input theo format này

Dùng prompt sau để gửi cho AI có khả năng phân tích video/audio, kèm theo link video hoặc file audio:

```
Hãy tách lời bài hát nursery rhyme từ video/audio này thành format timestamp theo mẫu dưới đây. 

QUY TẮC BẮT BUỘC:
1. Mỗi dòng lời phải có timestamp dạng (M:SS – M:SS) ở đầu
2. Gom các dòng lời thành các section, đánh dấu bằng header [VERSE N: Chủ đề] 
   - N = số thứ tự verse
   - Chủ đề = đối tượng/hành động chính của verse đó (ví dụ: "The Wheels", "The Pig")
3. Đoạn nhạc không có lời: ghi [INSTRUMENTAL] sau timestamp, theo sau là mô tả ngắn
4. Thêm [INTRO] cho phần mở đầu và [OUTRO] cho phần kết thúc
5. Timestamp phải CHÍNH XÁC theo video — bắt đầu và kết thúc của từng dòng lời
6. Mỗi dòng lời nên có thời lượng 3-8 giây
7. Nếu 1 dòng quá dài (>8 giây), tách thành 2 dòng
8. Dòng trống giữa các section là OK

MẪU OUTPUT:

[INTRO]
(0:00 – 0:05) [INSTRUMENTAL] Logo intro
(0:06 – 0:12) [INSTRUMENTAL] Music intro, establishing scene

[VERSE 1: The Wheels]
(0:13 – 0:18) The wheels on the bus go round and round, round and round, round and round
(0:19 – 0:24) The wheels on the bus go round and round, all through the town

[VERSE 2: The Door]
(0:25 – 0:26) [INSTRUMENTAL] Transition
(0:27 – 0:32) The door on the bus goes open and shut, open and shut, open and shut
(0:33 – 0:38) The door on the bus goes open and shut, all through the town

[OUTRO]
(1:50 – 1:58) [INSTRUMENTAL] Ending music, logo

Hãy xử lý video/audio sau và trả về kết quả theo ĐÚNG format trên. 
Không thêm giải thích, chỉ trả về nội dung đã format.
```

---

## Checklist trước khi paste vào Tool

- [ ] Mỗi dòng có timestamp `(M:SS – M:SS)` hợp lệ?
- [ ] Timestamp tuần tự, không chồng chéo?
- [ ] Có ít nhất 1 `[VERSE N: Subject]` header?
- [ ] Có `[INTRO]` và `[OUTRO]` (tùy chọn)?
- [ ] Đoạn nhạc không lời có tag `[INSTRUMENTAL]`?
- [ ] Mỗi dòng lời 3-8 giây?
- [ ] Tổng thời lượng khớp với video gốc?
