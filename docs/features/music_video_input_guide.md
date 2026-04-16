# Music Video Input Guide (MV Maker & Nursery Rhyme)

Cả hai chế độ chia phân cảnh âm nhạc là **MV Maker** (Music Video) và **Nursery Rhyme** (Nhạc Thiếu Nhi) dùng chung một nhân lõi xử lý (Lyrics Parser). Do đó, cách thức nhập liệu và sử dụng thẻ `[NOTE]` của chúng là hoàn toàn giống nhau.

## Bộ Quy Tắc Cốt Lõi (Core Rules)
1. **Dòng có Timestamp = Lời bài hát (Lyrics):** Bất kỳ dòng nào bắt đầu bằng `(M:SS - M:SS)` hoặc `(MM:SS - MM:SS)` sẽ được hệ thống ngầm định là lời bài hát để hiển thị trên màn hình.
2. **Thẻ Section = Cấu trúc bài hát:** Các thẻ như `[VERSE]`, `[CHORUS]`, `[DROP]`, `[INTRO]` giúp AI hiểu nhịp độ. Ví dụ: khi gặp `[CHORUS]`, AI sẽ chuyển sang cắt cảnh nhanh (fast cuts).
3. **Thẻ [NOTE] = Chỉ đạo nghệ thuật (Creative Direction):** Bất kỳ dòng nào bắt đầu bằng `[NOTE]`, hoặc chữ nằm sau `[NOTE]` trên cùng một dòng với timestamp, sẽ được AI dùng để bố trí camera, góc máy, ánh sáng, hành động nhân vật... nhưng sẽ được **tự động cắt bỏ hoàn toàn** khỏi phần Lời bài hát để đảm bảo không bị lẫn lộn vào video cuối cùng.

---

## 1. Input Sample: MV Maker (Phong cách Gaming Horror)
Được sử dụng cho các MV có nhịp độ nhanh, kịch tính, phong cách kinh dị/hành động (VD: nhạc FNAF, Poppy Playtime, Sprunki...).

```text
[NOTE] Jester là con hề khổng lồ, răng nhọn mọc lỉa chỉa, cánh tay robot. Bối cảnh là nhà máy đồ chơi bỏ hoang rỉ sét.
[NOTE] Ánh sáng cực kỳ u tối, chủ yếu dùng đèn pin chiếu từ dưới lên mặt (under-lighting) kết hợp sương mù dày đặc.

[INTRO]
(0:00 - 0:10) [INSTRUMENTAL] 
[NOTE] Máy quay từ từ tiến qua làn sương mù đỏ sẫm. Một đôi mắt vàng phát sáng hiện ra trong bóng tối.

[VERSE 1: Sự thức tỉnh]
(0:11 - 0:15) Welcome to the hour of joy
[NOTE] Jester treo ngược từ trên trần nhà xuống như một con dơi, miệng cười nhếch mép.
(0:16 - 0:20) We've built a home for every girl and boy
[NOTE] Máy quay lướt qua sàn nhà đầy đồ chơi vỡ nát.

[CHORUS: Vụ nổ]
[NOTE] ĐOẠN NÀY CẮT CẢNH NHANH THEO NHỊP TRỐNG! Camera rung kịch liệt, góc quay nghiêng (Dutch angle).
(0:21 - 0:25) Breathe it in, the scarlet smoke!
[NOTE] Jester phun ra một lượng lớn khói đỏ từ miệng.
(0:26 - 0:30) You'll sleep forever until you choke!
[NOTE] Smash cut: Chuyển cảnh đột ngột vào cực gần (extreme close-up) hàm răng nhọn của Jester choán hết màn hình.
```

---

## 2. Input Sample: Nursery Rhyme (Nhạc Thiếu Nhi)
Được sử dụng cho các video ca nhạc với nhịp độ chậm, minh họa trực quan, nội dung an toàn và tươi sáng.

```text
[NOTE] Có ba nhân vật chính: Gấu Leo (dũng cảm, mặc áo choàng đỏ), Chuột Mimi (đeo kính). Phong cách hình ảnh 3D đất sét, tươi sáng, mượt mà và thân thiện với trẻ em.

[VERSE 1: Gặp gỡ Leo]
(0:00 - 0:04) Here comes Leo, the king of the land
[NOTE] Trời nắng đẹp, bãi cỏ xanh mướt. Gấu Leo đứng chống hông tạo dáng anh hùng.
(0:05 - 0:09) With a big red cape and a waving hand
[NOTE] Cận cảnh (Close-up) Gấu Leo vui vẻ vẫy tay. Góc máy cố định, an toàn.

[CHORUS: Tình bạn]
(0:10 - 0:15) We are friends, we like to play!
[NOTE] Leo và Mimi nhảy múa quanh một bông hoa hướng dương khổng lồ. Cầu vồng rực rỡ ở phía sau.
(0:16 - 0:20) Learning new things every day!

[INSTRUMENTAL]
(0:21 - 0:26) [INSTRUMENTAL]
[NOTE] Một chú chim xanh dương bay ngang qua màn hình, thả xuống các nốt nhạc vàng lung linh.
```

---

## 💡 Mẹo sử dụng thẻ `[NOTE]`
- **Global Notes (Ghi chú toàn cục):** Đặt thẻ `[NOTE]` ở **trên cùng** của kịch bản, trước bất kỳ thẻ `[VERSE]` nào. Mục đích để thiết lập thiết kế nhân vật, mô tả bối cảnh chung, thời tiết, và tone màu tổng thể. Mọi Extractors (Character/Scene/Prop) đều sẽ đọc chúng để tạo kho tài nguyên chuẩn xác.
- **Block Notes (Ghi chú cho phân cảnh):** Đặt thẻ `[NOTE]` **bên dưới** dòng thời gian `(M:SS - M:SS)` để yêu cầu AI tạo ra hiệu ứng đặc biệt, góc máy nghiêng hay hành động riêng biệt chỉ dành cho khoảng thời gian bài hát đó.
- **Inline Notes (Ghi chú trên cùng một dòng):** Nếu lười bạn hoàn toàn có thể viết liền một mạch:
  `(0:40 - 0:45) I will find you [NOTE] Jester áp sát đe dọa camera.`
  Hệ thống Code sẽ tự động "cắt bỏ" cụm `[NOTE] Jester...` và chỉ giữ lại đúng câu `"I will find you"` cho vào Lời bài hát để hiển thị trên video.
