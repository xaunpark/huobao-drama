# BubiLoo Input Guide — Hướng Dẫn Nhập Liệu Kênh BubiLoo

Hướng dẫn này dành riêng cho kênh **BubiLoo** (Nursery Rhyme 3D CGI, phong cách Cocomelon-inspired). Kênh sử dụng **hệ thống nhân vật gốc** đã được thiết kế sẵn — bạn chỉ cần gọi tên nhân vật, hệ thống sẽ tự động nhận diện từ ảnh reference.

> [!IMPORTANT]
> **Nguyên tắc vàng:** KHÔNG mô tả ngoại hình nhân vật trong input. Chỉ gọi tên (Bubi, Mama, Papa...) và mô tả **hành động + cảm xúc**. Hệ thống đã có ảnh reference cho mọi nhân vật.

---

## Hệ Thống Nhân Vật BubiLoo (Tham Chiếu Nhanh)

| Tên | Vai trò | Khi nào dùng |
|---|---|---|
| **Bubi** | Toddler chính (~2 tuổi) | **MỌI** episode. Luôn là tâm điểm |
| **Mama** | Mẹ, narrator chính | ~80% episodes. Giọng hát chính người lớn |
| **Papa** | Bố, dạy qua hành động | ~50% episodes. Hoạt động ngoài trời, sửa chữa |
| **Luli** | Chị gái (~6 tuổi) | ~40% episodes. Dẫn dắt, dạy em |
| **Mochi** | Mèo trắng, comic relief | ~60% episodes. Gây cười, phản ứng ngộ nghĩnh |
| **Nana** | Bà ngoại | ~20% episodes. Nấu ăn, kể chuyện |
| **Popo** | Ông ngoại | ~20% episodes. Chơi nhạc, kể chuyện |
| **Ziggy** | Bạn trai (African) | ~25% episodes. Năng động, hào hứng |
| **Mei** | Bạn gái (Asian) | ~25% episodes. Thông minh, kiên nhẫn |
| **Rio** | Bạn trai (Latin) | ~25% episodes. Hài hước, comic relief |
| **Teacher Sunny** | Cô giáo mầm non | ~15% episodes. Chỉ dùng cho school-themed |

---

## Bộ Quy Tắc Cốt Lõi

1. **Dòng có Timestamp = Lời bài hát (Lyrics):** Dòng bắt đầu bằng `(M:SS - M:SS)` sẽ được hiển thị làm lời hát trên video.
2. **Thẻ Section = Cấu trúc bài hát:** `[VERSE]`, `[CHORUS]`, `[INTRO]`, `[OUTRO]` giúp AI hiểu nhịp độ và chuyển cảnh.
3. **Thẻ `[NOTE]` = Chỉ đạo nghệ thuật:** Mô tả hành động, góc máy, bối cảnh. Sẽ được **tự động cắt bỏ** khỏi phần Lyrics hiển thị trên video.

> [!CAUTION]
> **KHÔNG BAO GIỜ** viết trong `[NOTE]`:
> - Mô tả ngoại hình nhân vật (đầu trọc, mắt sao, yếm vàng...)
> - Tên "JJ", "Cocomelon", logo dưa hấu, hoặc bất kỳ branding nào
> - Text/chữ xuất hiện trên màn hình
>
> **CHỈ NÊN** viết trong `[NOTE]`:
> - Hành động: "Bubi nhảy lên vui sướng", "Mama ôm Bubi"
> - Bối cảnh: "Phòng tắm sáng sủa", "Công viên nắng đẹp"
> - Góc máy: "Cận cảnh mặt Bubi", "Wide shot cả gia đình"
> - Cảm xúc: "Bubi ngạc nhiên miệng tròn O"

---

## Quy Trình Chuyển Đổi Bài Hát Truyền Thống

Nhiều nursery rhyme nổi tiếng (Baa Baa Black Sheep, Mary Had a Little Lamb, Wheels on the Bus...) có nhân vật riêng không thuộc hệ thống BubiLoo. Quy trình chuyển đổi:

### Bước 1: Xác định nhân vật BubiLoo đóng vai

Gán nhân vật gốc của bài hát cho nhân vật BubiLoo phù hợp nhất:

| Nhân vật gốc trong bài hát | Gán cho nhân vật BubiLoo | Lý do |
|---|---|---|
| "Mary" (Mary Had a Little Lamb) | **Bubi** hoặc **Luli** | Trẻ em — nhân vật chính |
| "Little Bo Peep" | **Luli** | Cô bé chăn cừu → chị gái |
| "Old MacDonald" | **Papa** hoặc **Popo** | Người lớn nam giới |
| "Miss Polly" | **Mama** | Người phụ nữ chăm sóc |
| "The Doctor" | **Papa** (hoặc nhân vật mới) | Người lớn nam |
| Nhóm bạn / trẻ em khác | **Ziggy, Mei, Rio** | Nhóm bạn đa dạng |
| Thầy/cô giáo | **Teacher Sunny** | Cô giáo có sẵn |

### Bước 2: Động vật trong bài hát → Nhân vật mới

Các con vật trong bài hát (cừu, mèo, nhện...) **KHÔNG cần gán** cho nhân vật BubiLoo. Thay vào đó, mô tả chúng trong `[NOTE]` dưới dạng "thú cưng/vật nuôi phong cách 3D toy":

```
[NOTE] Con cừu: nhỏ, mềm mại, lông trắng bông, mắt to tròn, phong cách đồ chơi nhựa mềm.
```

> [!TIP]
> **Mochi (mèo nhà)** vẫn có thể xuất hiện cùng lúc với các con vật khác trong bài hát, đóng vai quan sát/phản ứng hài hước.

### Bước 3: Viết input với hệ thống BubiLoo

Sau khi xác định nhân vật, viết input theo format chuẩn. Xem ví dụ bên dưới.

---

## 1. Input Mẫu: Bài Hát Gốc BubiLoo

Bài hát được sáng tác riêng cho kênh, sử dụng đúng nhân vật hệ thống.

```text
[NOTE] Bubi học rửa tay. Bối cảnh: phòng tắm nhà Bubi, tươi sáng, sạch sẽ.
[NOTE] Mama hướng dẫn Bubi từng bước. Mochi ngồi trên mép bồn tắm quan sát.

[INTRO]
(0:00 - 0:07) [INSTRUMENTAL]
[NOTE] Wide shot phòng tắm sáng sủa. Bubi đứng trước bồn rửa, ngước lên nhìn vòi nước.

[VERSE 1: Rửa tay nào]
(0:08 - 0:12) Let's wash our hands, wash wash wash!
[NOTE] Mama đứng cạnh Bubi, bật vòi nước. Bubi giơ hai tay lên hào hứng.
(0:13 - 0:17) Scrub scrub scrub, squish squish squish!
[NOTE] Cận cảnh hai bàn tay bé Bubi đầy bọt xà phòng, bong bóng bay lên.
(0:18 - 0:22) Bubbles here, bubbles there!
[NOTE] Mochi vươn chân bắt bong bóng, mất thăng bằng suýt ngã. Bubi cười.

[CHORUS]
(0:23 - 0:30) Clean clean hands, everywhere!
[NOTE] Bubi giơ hai tay sạch bóng lên cao, cười rạng rỡ. Mama vỗ tay khen.
(0:31 - 0:37) We wash our hands because we care!
[NOTE] Wide shot: Bubi và Mama high-five. Mochi nhảy xuống bồn tắm chạy đi.

[VERSE 2: Lần nữa nào]
(0:38 - 0:42) One more time, wash wash wash!
[NOTE] Bubi tự bật vòi nước lần này, tự tin hơn. Mama gật đầu khích lệ.
(0:43 - 0:47) Under the water, splash splash splash!
[NOTE] Nước bắn tung tóe nhẹ, Bubi cười khanh khách.

[OUTRO]
(0:48 - 0:55) All done! Squeaky clean!
[NOTE] Bubi ôm Mama. Mochi quay lại với bọt xà phòng dính trên mũi. Cả nhà cười.
(0:56 - 1:00) [INSTRUMENTAL]
[NOTE] Fade out nhẹ nhàng. Star sparkle wipe ✨.
```

---

## 2. Input Mẫu: Bài Hát Truyền Thống — "Baa Baa Black Sheep"

Bài hát truyền thống được chuyển đổi sang thế giới BubiLoo. Con cừu là nhân vật mới, còn "master", "dame", "little boy" được gán cho nhân vật có sẵn.

```text
[NOTE] Bubi gặp một chú cừu đen ở nông trại. Cừu đen: nhỏ nhắn, lông xù mềm, mắt to hiền lành, phong cách đồ chơi nhựa mềm.
[NOTE] Bối cảnh: nông trại nắng đẹp, hàng rào gỗ tròn, cỏ xanh mướt, bầu trời trong xanh.
[NOTE] Papa = "master", Mama = "dame", Bubi = "little boy who lives down the lane".

[VERSE 1]
(0:00 - 0:04) Baa baa black sheep, have you any wool?
[NOTE] Wide shot nông trại. Bubi ngồi xổm trước mặt chú cừu đen, nghiêng đầu tò mò.
(0:05 - 0:09) Yes sir, yes sir, three bags full!
[NOTE] Cừu đen gật đầu vui vẻ. Ba túi len đầy ắp xuất hiện bên cạnh.

[VERSE 2]
(0:10 - 0:14) One for the master,
[NOTE] Bubi bê một túi len chạy đến Papa. Papa mỉm cười đón nhận.
(0:15 - 0:19) One for the dame,
[NOTE] Bubi bê túi thứ hai đến Mama. Mama ôm túi len và xoa đầu Bubi.
(0:20 - 0:24) And one for the little boy who lives down the lane!
[NOTE] Bubi ôm túi len cuối cùng vào ngực, cười tít mắt. Cừu đen đi theo Bubi.

[CHORUS]
(0:25 - 0:33) Baa baa black sheep, have you any wool?
[NOTE] Tất cả cùng hát. Mochi xuất hiện, nhảy lên đống len rồi bị lún xuống. Cả nhà cười.
(0:34 - 0:40) Yes sir, yes sir, three bags full!
[NOTE] Wide shot nông trại: Bubi, Mama, Papa, cừu đen và Mochi. Cầu vồng nhẹ phía sau.

[OUTRO]
(0:41 - 0:48) [INSTRUMENTAL]
[NOTE] Bubi vẫy tay chào cừu đen. Cừu kêu "baa" lần cuối. Cloud dissolve ☁️ kết thúc.
```

---

## 3. Input Mẫu: Bài Hát Truyền Thống — "Wheels on the Bus"

Bài hát dài, nhiều nhân vật, phù hợp để đưa nhóm bạn vào.

```text
[NOTE] Bubi và các bạn đi xe buýt đến trường. Xe buýt: tròn trĩnh, màu vàng tươi, phong cách đồ chơi Fisher-Price.
[NOTE] Bối cảnh: con đường làng xanh mát, trường mầm non phía xa.
[NOTE] Mama = người lái xe buýt (bus driver). Ziggy, Mei, Rio = babies on the bus.

[VERSE 1: Bánh xe]
(0:00 - 0:08) The wheels on the bus go round and round, round and round, round and round!
[NOTE] Wide shot xe buýt chạy trên con đường. Bánh xe quay tròn mượt mà.
(0:09 - 0:15) The wheels on the bus go round and round, all through the town!
[NOTE] Bên trong xe: Bubi ngồi hàng đầu, nhìn ra cửa sổ mắt sáng rỡ.

[VERSE 2: Cửa xe]
(0:16 - 0:24) The doors on the bus go open and shut!
[NOTE] Xe dừng lại. Cửa mở ra, Ziggy nhảy lên xe vẫy tay chào Bubi.
(0:25 - 0:31) Open and shut, open and shut, all through the town!
[NOTE] Mei và Rio cũng lên xe. Tất cả ngồi cạnh nhau, cười nói.

[VERSE 3: Em bé khóc]
(0:32 - 0:40) The babies on the bus go wah wah wah!
[NOTE] Ziggy, Mei, Rio giả vờ khóc "wah wah" rồi bật cười khanh khách. Bubi bụm miệng cười.
(0:41 - 0:47) Wah wah wah, all through the town!
[NOTE] Mochi nhảy từ ghế này sang ghế kia hoảng hốt vì tiếng khóc giả.

[VERSE 4: Mama lái xe]
(0:48 - 0:56) The driver on the bus says move on back!
[NOTE] Mama quay lại mỉm cười nhắc nhở. Bubi và các bạn ngồi ngay ngắn.
(0:57 - 1:03) Move on back, move on back, all through the town!
[NOTE] Xe đến trường. Teacher Sunny đứng trước cổng vẫy tay đón.

[OUTRO]
(1:04 - 1:10) [INSTRUMENTAL]
[NOTE] Tất cả xuống xe chạy vào trường. Bubi quay lại vẫy tay chào Mama trên xe. Star sparkle ✨.
```

---

## 4. Input Mẫu: Bài Hát Truyền Thống — "Mary Had a Little Lamb"

Ví dụ gán "Mary" cho Luli (chị gái lớn hơn, phù hợp với hình ảnh cô bé dắt cừu đi học).

```text
[NOTE] Luli có một chú cừu trắng bé nhỏ đi theo khắp nơi. Cừu trắng: nhỏ xíu, lông bông trắng mịn, đeo nơ hồng nhỏ trên cổ, mắt to ngây thơ.
[NOTE] Bối cảnh: nhà → con đường → trường mầm non.
[NOTE] Luli = "Mary". Bubi đi cùng Luli. Teacher Sunny = cô giáo ở trường.

[VERSE 1]
(0:00 - 0:07) Mary had a little lamb, little lamb, little lamb!
[NOTE] Luli dắt tay Bubi đi trên con đường, cừu trắng nhỏ lẽo đẽo phía sau.
(0:08 - 0:14) Mary had a little lamb, its fleece was white as snow!
[NOTE] Cận cảnh cừu trắng dụi mũi vào chân Luli. Bubi cúi xuống vuốt ve.

[VERSE 2]
(0:15 - 0:22) It followed her to school one day, school one day, school one day!
[NOTE] Trước cổng trường. Luli và Bubi đi vào, cừu lẻn theo phía sau.
(0:23 - 0:29) Which was against the rules!
[NOTE] Teacher Sunny nhìn thấy cừu, đưa tay lên miệng ngạc nhiên. Các bạn (Ziggy, Mei) quay lại nhìn.

[VERSE 3]
(0:30 - 0:37) It made the children laugh and play, laugh and play!
[NOTE] Cừu nhảy lên bàn, Ziggy và Mei cười vỗ tay. Rio đuổi theo cừu quanh lớp. Bubi ôm bụng cười.
(0:38 - 0:44) To see a lamb at school!
[NOTE] Teacher Sunny lắc đầu cười, nhẹ nhàng bế cừu ra ngoài. Mochi đang ngủ trên bậu cửa sổ, giật mình.

[OUTRO]
(0:45 - 0:53) [INSTRUMENTAL]
[NOTE] Sau giờ học, Luli và Bubi ra sân. Cừu trắng chạy đến mừng. Cả ba ôm nhau. Cloud dissolve ☁️.
```

---

## 💡 Mẹo Sử Dụng

### Thẻ `[NOTE]` — Cách viết đúng cho BubiLoo

- **Global Notes (trên cùng):** Dùng để thiết lập bối cảnh chung, giới thiệu nhân vật MỚI (không thuộc hệ thống), và gán vai.
  ```
  [NOTE] Bối cảnh: bãi biển nhiệt đới, sóng nhẹ, cát vàng, trời nắng.
  [NOTE] Con cua: nhỏ, đỏ cam, mắt trên cuống, dáng tròn trĩnh phong cách đồ chơi.
  [NOTE] Papa = "sailor" trong bài hát.
  ```

- **Block Notes (sau timestamp):** Mô tả hành động và cảm xúc cho từng đoạn.
  ```
  (0:10 - 0:15) Row row row your boat!
  [NOTE] Bubi và Papa ngồi trên thuyền nhỏ, chèo tay đồng bộ. Mochi nằm mũi thuyền.
  ```

- **Inline Notes (cùng dòng):** Viết nhanh gọn.
  ```
  (0:20 - 0:25) Merrily merrily merrily! [NOTE] Bubi cười tít mắt, Papa chèo nhanh hơn.
  ```

### Quy tắc gán nhân vật cho bài hát truyền thống

| Vai trò trong bài hát | Ưu tiên gán | Ghi chú |
|---|---|---|
| Trẻ em chính (chủ ngữ) | **Bubi** | Luôn là lựa chọn đầu tiên |
| Trẻ em lớn hơn / chị/anh | **Luli** | Khi cần nhân vật trẻ em "già hơn" |
| Mẹ / phụ nữ trưởng thành | **Mama** | Vai trò chăm sóc, dạy dỗ |
| Bố / đàn ông trưởng thành | **Papa** | Vai trò hành động, làm việc |
| Ông bà / người già | **Nana / Popo** | Vai trò kể chuyện, truyền thống |
| Nhóm bạn / trẻ em khác | **Ziggy, Mei, Rio** | Chia đều, tạo đa dạng |
| Cô giáo / người hướng dẫn | **Teacher Sunny** | Chỉ khi bối cảnh ở trường |
| Động vật | **Mô tả mới trong [NOTE]** | Giữ phong cách "3D toy" |

### Bubi luôn xuất hiện

Dù bài hát gốc không có nhân vật phù hợp với Bubi, **hãy luôn tìm cách đưa Bubi vào**. Bubi có thể:
- Là người quan sát/tương tác với nhân vật chính
- Là người nghe Mama kể/hát câu chuyện
- Là "little boy/girl" trong bài hát
- Đơn giản là xuất hiện ở bên cạnh, phản ứng vui vẻ

> [!TIP]
> **Mochi** cũng nên xuất hiện thường xuyên để tạo comic relief — nhảy lên đồ vật, bắt bướm, ngủ quên ở góc, hoặc phản ứng ngạc nhiên trước tình huống hài hước.

### Transition kết thúc

Luôn kết thúc bằng một trong hai branded transitions:
- `Star sparkle wipe ✨` — cho kết thúc vui vẻ, năng động
- `Cloud dissolve ☁️` — cho kết thúc nhẹ nhàng, ấm áp
