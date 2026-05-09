---
name: cg5-poppy-playtime
description: Generate production-ready MV Maker script input for CG5-style Poppy Playtime music videos
---

# CG5 Poppy Playtime — Script Generator Skill

Tạo **script text sẵn sàng** (MV Maker input format) cho video nhạc Poppy Playtime theo phong cách CG5.

## Khi nào dùng

User yêu cầu tạo script input cho một bài hát CG5 Poppy Playtime. Đầu vào thường là:
- Timestamp + lyrics của bài hát (copy từ YouTube / lyrics site)
- Hoặc tên bài hát CG5 (Agent tự tìm lyrics)

## Đầu ra

Script text ở format **MV Maker** (xem `docs/features/music_video_input_guide.md`), sẵn sàng paste vào hệ thống.

---

## Instrumentation

```bash
./scripts/log-skill.sh "cg5-poppy-playtime" "manual" "$$"
```

---

## Quy trình

### Bước 1: Thu thập đầu vào

User cung cấp lyrics có timestamp. Format nhận dạng:
```
(M:SS - M:SS) lyrics line
```
Nếu user chỉ cung cấp tên bài hát, tìm lyrics + timestamp trước.

### Bước 2: Xác định Chapter & nhân vật

Dựa vào lyrics/title, xác định:
- **Chapter nào** của Poppy Playtime (xem bảng CG5 Songs bên dưới)
- **Villain chính** (POV character — ai đang "hát")
- **Nhân vật phụ** xuất hiện trong MV
- **Locations** phù hợp với chapter đó

### Bước 3: Generate Script

Áp dụng các quy tắc bên dưới để tạo output.

---

## Output Format — QUY TẮC BẮT BUỘC

### Cấu trúc file output

```text
[NOTE] Global notes — thiết lập thế giới, nhân vật, phong cách
[NOTE] Global notes — mô tả nhân vật villain chính
[NOTE] Global notes — mô tả nhân vật phụ (nếu có)

[SECTION TAG]
(timestamp) lyrics hoặc [INSTRUMENTAL]
[NOTE] Block note — chỉ đạo camera/hành động/ánh sáng cho đoạn này
```

### Quy tắc cốt lõi

1. **Dòng có Timestamp = Lyrics.** Format `(M:SS - M:SS)` — đây là lời hát hiển thị trên video. KHÔNG chèn ghi chú vào đây.
2. **`[NOTE]` = Chỉ đạo nghệ thuật.** Mô tả camera, hành động nhân vật, ánh sáng, hiệu ứng. Sẽ bị cắt bỏ khỏi lyrics cuối cùng.
3. **Section tags** = `[INTRO]`, `[VERSE 1]`, `[VERSE 2]`, `[PRE-CHORUS]`, `[CHORUS]`, `[BRIDGE]`, `[DROP]`, `[BUILDUP]`, `[OUTRO]` — giúp AI hiểu nhịp độ cắt cảnh.
4. **`[INSTRUMENTAL]`** dùng cho đoạn không có lời.

### Global Notes — BẮT BUỘC ở đầu file

Global notes thiết lập thế giới Poppy Playtime cho toàn bộ video. **PHẢI** bao gồm:

#### Note 1: Bối cảnh & phong cách
```text
[NOTE] Bối cảnh: nhà máy đồ chơi Playtime Co. bỏ hoang, {location cụ thể theo chapter}.
Phong cách: cinematic 3D CGI, mascot horror. Ánh sáng cực kỳ u tối (under-lighting từ dưới lên),
sương mù dày đặc màu xanh-tím đậm. Rim light màu {purple/cyan/red}.
Mọi bề mặt rỉ sét, bụi bẩn, hư hỏng (nhà máy bỏ hoang từ 1995).
```

#### Note 2: Villain chính (POV — người đang hát)
```text
[NOTE] {Tên villain} là {mô tả ngắn gọn từ Character Reference}.
{Chi tiết ngoại hình: vật liệu, kích thước, đặc điểm nhận dạng, vũ khí/khả năng}.
{Cách di chuyển đặc trưng}.
```

#### Note 3+: Nhân vật phụ (nếu xuất hiện trong MV)
```text
[NOTE] {Tên nhân vật phụ}: {mô tả ngắn gọn, vai trò trong MV này}.
```

### Block Notes — Quy tắc cho từng đoạn

Mỗi `[NOTE]` sau timestamp phải tuân thủ **CG5 Visual DNA**:

| Yếu tố | Quy tắc | Ví dụ |
|---|---|---|
| **Góc máy** | 70% low-angle (nhìn lên), villain trông đồ sộ | "Low angle nhìn lên Huggy Wuggy sừng sững" |
| **Shot type** | 40% MS, 20% MCU, 15% CU, 20% text card, 5% POV | "Medium shot Prototype đứng giữa hành lang" |
| **Camera** | Slow push-in (40%), handheld shake ở chorus | "Camera từ từ tiến về phía mặt CatNap" |
| **Ánh sáng** | LUÔN u tối, under-lighting, rim light màu | "Under-lit harsh, purple rim light từ trái" |
| **Sương mù** | LUÔN có volumetric fog | "Sương mù dày đặc che khuất hành lang phía sau" |
| **3D Text** | Mô tả text nổi trong không gian 3D | "Chữ 'SLEEP WELL' phát sáng vàng neon, nổi trong sương mù" |
| **Glitch** | Tăng dần từ verse → bridge | "Chromatic aberration mạnh, screen shake theo nhịp trống" |
| **Nhịp cắt** | Verse: 2-4s/shot. Chorus: 1-2s. Bridge: <1s | "CẮT CẢNH NHANH THEO NHỊP TRỐNG" |

### Transition Notes — Quy tắc chuyển cảnh

- 80% Hard cut ON BEAT (cắt theo nhịp)
- 20% Glitch transition (chromatic aberration burst, VHS tracking, frame tear)
- **KHÔNG dùng**: fade, dissolve, cross-fade (trừ fade from black ở intro)

### Section-specific Pacing

| Section | Nhịp | Ghi chú cho [NOTE] |
|---|---|---|
| `[INTRO]` | 5-10s/shot | Atmospheric, fog, establishing factory |
| `[VERSE]` | 2-4s/shot | Villain giới thiệu, deceptive calm, slow push-in |
| `[PRE-CHORUS]` | 2-3s/shot | Tension build, mask slips, lighting harsher |
| `[CHORUS]` | 1-2s/shot | RAPID CUTS, Dutch angle, max menace, handheld shake |
| `[BRIDGE]` | <1s/shot | JUMP CUTS, glitch heavy, breakdown, chaotic |
| `[BUILDUP]` | 4-6s hold | Dramatic pause → godhood declaration |
| `[OUTRO]` | 5-10s | Fade to darkness/CRT/glowing eye/poppy gel |

---

## Character Reference (Canonical)

Dùng chính xác các mô tả này trong Global Notes.

### Villains (POV Characters)

**The Prototype (Experiment 1006)**
> Thực thể khổng lồ, cơ thể lắp ghép từ titanium, bánh răng, dây điện, xương và mô hữu cơ.
> Mặt sứ trắng hình hề (mũ ba đỉnh có chuông, trang phục đỏ/vàng/xanh).
> Răng nhọn dài, mắt xoay như camera. Cánh tay gầy guộc với móng vuốt sắc.
> Ghép các bộ phận đồ chơi bị đánh bại lên cơ thể mình.
> Di chuyển: teleport-glitch, xuất hiện/biến mất bất ngờ.

**Huggy Wuggy (Experiment 1170)**
> Cao 5.5m, lông XANH DƯƠNG rối bẩn. Tay/chân vàng có miếng Velcro, ngón dính liền.
> Mắt đen to trợn, môi đỏ lớn che hàng răng kim nhọn. Hàm kép (kiểu lươn moray).
> Nơ xanh nhỏ. Cơ bắp/xương dưới lớp lông.
> Di chuyển: rình mò săn đuổi, kéo dãn chui qua ống thông gió.

**Mommy Long Legs (Experiment 1222)**
> Cao, mảnh, giống nhện. Chủ đạo màu HỒNG.
> Đầu hình elip, mắt xanh lá to tròn (3 lông mi), miệng đen cười rộng với son hồng.
> Tóc hồng đậm buộc đuôi ngựa dài. Chi cực dài ĐÀN HỒI co giãn.
> Di chuyển: kéo giãn chi, quấn/vươn không giới hạn.

**CatNap (Experiment 1188)**
> Mèo lông TÍM, tai tam giác, đuôi rất dài.
> Bigger Body: mèo khổng lồ đáng sợ, mắt trắng lạnh, nụ cười cố định.
> Phát ra KHÓI ĐỎ gây ảo giác/chết. Tín đồ cuồng tín của Prototype.
> Di chuyển: lơ lửng buồn ngủ → tấn công chớp nhoáng.

**Lily Lovebraids (Ch.5)**
> Búp bê thời trang 90s. Bím tóc điều khiển được (tóc ngựa + xương ngón ghép) di chuyển như đuôi bọ cạp.
> Sinh lý giống nhện. Mang đầu Candy Cat vỡ dùng nói bụng (giọng ngọt qua búp bê, hung hăng qua chính mình).
> Di chuyển: bò kiểu nhện, bím tóc quất.

**Miss Delight (Ch.3)**
> Búp bê giáo viên khớp cầu. Váy chấm bi đỏ thẫm, yếm xanh (biểu tượng táo), áo vàng, cà vạt đỏ.
> Mặt vỡ lộ thịt/răng bên trong. Vũ khí: "Barb" — chùy sao từ bút chì mài nhọn.
> Di chuyển: đứng yên khi bị nhìn, chạy khi không bị quan sát ("Weeping Angel").

**Yarnaby (Experiment 1166, Ch.4 — ĐÃ CHẾT)**
> Sư tử hoạt hình, cơ thể LEN SỢI nhiều màu. Mặt hổ phách, đường khâu đen dọc.
> Mắt trắng to, đồng tử đen. Miệng cười rộng, 3 răng nanh tam giác lệch.
> Di chuyển: nhảy múa vui tươi → tấn công hung hãn.

### Allies (không phải POV, nhưng có thể xuất hiện)

**Poppy (Experiment 1007)**
> Búp bê sứ ~40cm. Da trắng phấn, tóc ĐỎ xoăn buộc hai bím (nơ xanh).
> Mắt xanh to (đỏ ngầu theo thời gian), tàn nhang, má hồng. Váy Victorian xanh.
> Dây kéo sau lưng. Mặt nứt (Ch.5+), băng bó, váy rách.

**Kissy Missy (Experiment 1172)**
> Cấu trúc giống Huggy nhưng lông HỒNG, nơ xanh nhạt, lông mi dài. Đồng minh.

**Giblet (Ch.5)**
> Plush chihuahua/cáo lộn trái. Cam-nâu, lông rối. Mắt trái: đồng tử magenta/mống hoa.
> Mắt phải: mất (miếng da đen che). Mũi bu-lông kim loại. Áo frock, nút Huggy, găng vá, gậy-taser.

**Chum Chompkins (Ch.5)**
> To tròn, lông ĐỎ dày, bụng/tay/chân vàng. MIỆNG TRÊN BỤNG. Mắt lệch ngộ nghĩnh. Câm.

---

## Location Reference (theo Chapter)

| Chapter | Locations chính |
|---|---|
| Ch.1 | Sàn nhà máy chính, kho Make-A-Friend, hệ thống thông gió |
| Ch.2 | Game Station (Statues arena, Musical Memory), hang Mommy, tơ nhện |
| Ch.3 | Playcare (nhà trẻ ngầm), School, Home Sweet Home, hang CatNap (khói tím đỏ) |
| Ch.4 | Safe Haven, lãnh thổ Yarnaby, lò nung |
| Ch.5 | Phòng thí nghiệm (sâu nhất), Boiler Room, Biodiversity Labs, Lily's Dollhouse/Sweet Street, Outimal Tunnels (sợ ánh sáng), Reanimation Chamber, hệ thống tàu |

---

## CG5 Songs → Chapter Mapping

| Bài hát CG5 | Chapter | Villain POV | Nhân vật phụ gợi ý |
|---|---|---|---|
| **Poison Blooms** | Ch.1 | Huggy Wuggy | Poppy, Player (GrabPack POV) |
| **Mommy's Here** | Ch.2 | Mommy Long Legs | Poppy, Player |
| **Sleep Well** | Ch.3 | CatNap | DogDay, Smiling Critters, Miss Delight |
| **HELL LIKE THIS** | Ch.4 | Yarnaby | Poppy, Player |
| **Wrong Side Out** | Ch.5 | The Prototype | Lily, Giblet, Chum, Huggy, Kissy, Poppy |

---

## Ví dụ Output Hoàn Chỉnh

### Input của user:
> Tạo script cho "Wrong Side Out" (CG5 - Poppy Playtime Ch.5)
> ```
> (0:00 - 0:12) [INSTRUMENTAL]
> (0:12 - 0:17) I can see you from the inside
> (0:17 - 0:22) Every stitch and every seam
> (0:22 - 0:27) I have made a new design
> (0:27 - 0:32) You will be my greatest dream
> (0:32 - 0:37) Don't you worry, don't you cry
> (0:37 - 0:42) I will turn you wrong side out
> (0:42 - 0:52) I CAN MAKE YOU BETTER!
> (0:52 - 0:57) Rip the fabric, pull the thread
> (0:57 - 1:02) Everything you were is dead
> (1:02 - 1:12) I CAN MAKE YOU BETTER!
> ```

### Output (sẵn sàng paste vào hệ thống):

```text
[NOTE] Bối cảnh: tầng sâu nhất của nhà máy Playtime Co. — Phòng thí nghiệm (The Laboratories).
Hành lang xám xịt, thiết bị phẫu thuật, bể poppy gel phát sáng mờ, xác thí nghiệm Bigger Bodies trong ống.
Phong cách: cinematic 3D CGI, mascot horror. Ánh sáng cực kỳ u tối (under-lighting từ dưới lên),
sương mù dày đặc màu xanh-tím đậm (#0A0B1A). Rim light tím (#7D12FF) và đỏ (#FF0033).
Mọi bề mặt rỉ sét, bụi bẩn, hư hỏng (nhà máy bỏ hoang từ 1995).

[NOTE] The Prototype (Experiment 1006) — villain chính, đang hát POV.
Thực thể khổng lồ, cơ thể lắp ghép từ titanium, bánh răng, dây điện, xương và mô hữu cơ.
Mặt sứ trắng hình hề (mũ 3 đỉnh có chuông, trang phục đỏ/vàng/xanh), răng nhọn dài, mắt xoay như camera.
Cánh tay gầy guộc với móng vuốt sắc. Ghép bộ phận đồ chơi bại trận lên cơ thể.
Di chuyển: teleport-glitch, xuất hiện/biến mất đột ngột.

[NOTE] Giblet: plush chihuahua lộn trái cam-nâu, miếng bịt mắt phải, gậy-taser. Đồng minh, bị Prototype đe dọa.
[NOTE] Outimals (Rag Bags/Gutter Plushes): sinh vật lộn trái, mắt nhựa đảo ngược, sợ ánh sáng. Bầy đàn do Prototype tạo ra.

[INTRO]
(0:00 - 0:12) [INSTRUMENTAL]
[NOTE] Fade from black. Camera từ từ tiến qua hành lang phòng thí nghiệm tối đen, sương mù xanh-tím lùa qua sàn.
Đèn huỳnh quang nhấp nháy. Bể poppy gel phát sáng mờ ở xa. Bóng đen khổng lồ của Prototype hiện ra cuối hành lang.
Chỉ có đôi mắt camera xoay phát sáng vàng trong bóng tối. Chữ "WRONG SIDE OUT" hiện lên bằng neon cam, nhỏ giọt như poppy gel.

[VERSE 1: Sự thao túng]
(0:12 - 0:17) I can see you from the inside
[NOTE] Low angle nhìn lên. Prototype bước ra từ sương mù, under-lit harsh, mặt sứ hề phát sáng nhẹ.
Mắt camera xoay zoom vào. Chữ "FROM THE INSIDE" nổi mờ trong sương mù phía sau.
(0:17 - 0:22) Every stitch and every seam
[NOTE] MCU mặt Prototype. Camera slow push-in. Các bộ phận ghép (tay plush, mắt búp bê) co giật trên cơ thể.
Rim light tím (#7D12FF) viền bên trái. Sương mù cuộn quanh vai.
(0:22 - 0:27) I have made a new design
[NOTE] MS Prototype dang rộng cánh tay. Phía sau: hàng Outimals lộn trái xếp hàng trong bóng tối, mắt nhựa đảo phát sáng đỏ.
Chữ "NEW DESIGN" neon vàng (#FFCC00) nổi 3D quanh đầu Prototype như hào quang.
(0:27 - 0:32) You will be my greatest dream
[NOTE] Camera quay sang Giblet bị nhốt trong lồng sắt, một mắt magenta nhìn ra sợ hãi. Gậy-taser nằm ngoài tầm với.
Under-lit từ dưới sàn lồng, bóng sọc qua mặt Giblet. Prototype's shadow đổ lên lồng.

[PRE-CHORUS]
(0:32 - 0:37) Don't you worry, don't you cry
[NOTE] Camera quay lại Prototype. Giọng chuyển từ dịu dàng → đe dọa. Nụ cười sứ mở rộng lộ răng nhọn.
Ánh sáng chuyển harsher. Sương mù đỏ (CatNap legacy) bắt đầu trộn vào sương mù tím.
Chromatic aberration nhẹ ở rìa khung hình.
(0:37 - 0:42) I will turn you wrong side out
[NOTE] CU tay vuốt Prototype vươn ra — từ từ, đe dọa. Các Outimals bắt đầu bò ra từ bóng tối phía sau.
Chữ "WRONG SIDE OUT" bùng nổ neon cam (#FF9900), chữ bị lộn trái/méo dạng, nổi 3D.
Camera shake nhẹ theo nhịp bass tăng dần.

[CHORUS]
[NOTE] ĐOẠN NÀY CẮT CẢNH NHANH THEO NHỊP TRỐNG! Dutch angle, handheld shake, max chromatic aberration.
(0:42 - 0:52) I CAN MAKE YOU BETTER!
[NOTE] Smash cut: Low angle cực thấp (worm's eye) nhìn lên Prototype dang tay, mặt sứ hề chiếm hết khung hình.
Mắt camera phát sáng vàng tối đa bloom. Outimals ào qua foreground như bầy đàn.
Chữ "I CAN MAKE YOU BETTER" KHỔNG LỒ neon vàng (#FFCC00), nổi 3D lấp đầy khung hình, rung theo nhịp.
GLITCH: screen shake mạnh, VHS tracking noise, chromatic aberration burst.
Rapid cuts xen kẽ: CU răng Prototype → CU mắt → CU vuốt → Outimals lao tới camera.

[VERSE 2: Hủy diệt]
(0:52 - 0:57) Rip the fabric, pull the thread
[NOTE] MS Prototype xé một con plush toy ra từng mảnh. Bông nhồi bay tung. Under-lit harsh.
Rim light đỏ (#FF0033). Chữ "RIP" và "PULL" xuất hiện rồi vỡ tan như kính.
(0:57 - 1:02) Everything you were is dead
[NOTE] Camera quay qua hàng xác Outimals — sinh vật lộn trái nằm la liệt.
Poppy gel đỏ phát sáng chảy trên sàn. Prototype đứng giữa, backlit silhouette.
Chữ "IS DEAD" chìm từ từ vào sương mù, ánh sáng tắt dần.

[CHORUS 2]
(1:02 - 1:12) I CAN MAKE YOU BETTER!
[NOTE] FINAL CHORUS — cường độ tối đa. Prototype ở trung tâm, cánh tay dang rộng, tất cả bộ phận ghép co giật.
Backlit bằng neon tím + đỏ. Sương mù cuộn cuồng. Outimals quỳ xung quanh như tín đồ.
Chữ "I CAN MAKE YOU BETTER" quay quanh Prototype như quỹ đạo, phát sáng rực.
GLITCH tối đa: frame tear, pixel corruption, VHS noise.
Cut cuối: màn hình CRT hiện mắt Prototype — tín hiệu nhiễu — fade to black.
```

---

## Checklist trước khi giao output

- [ ] Có ít nhất 2-3 `[NOTE]` global ở đầu (bối cảnh + villain + phụ)
- [ ] Mỗi dòng timestamp giữ đúng lyrics gốc, KHÔNG thêm bớt
- [ ] Mỗi đoạn lyrics có `[NOTE]` block kèm theo (camera + hành động + ánh sáng)
- [ ] Có section tags (`[VERSE]`, `[CHORUS]`, etc.)
- [ ] Villain mô tả đúng canonical reference (không bịa đặt ngoại hình)
- [ ] Location phù hợp chapter
- [ ] Có mô tả 3D text / kinetic typography cho các câu hook
- [ ] Nhịp cắt tăng dần: verse (chậm) → chorus (nhanh) → bridge (cực nhanh)
- [ ] LUÔN có: under-lighting, volumetric fog, rim light, glitch effects
- [ ] KHÔNG có: daylight, clean surfaces, cute elements, 2D, text overlay
- [ ] Transition: 80% hard cut on beat, 20% glitch transition

---

## References

- Template chi tiết: `docs/cg5_poppy_playtime_template.md`
- Template gốc (multi-game): `docs/cg5_template.md`
- Format input: `docs/features/music_video_input_guide.md`
- Pipeline: `ai/indexes/pipeline-map.md` (Pipeline 1: Script → Storyboard)
