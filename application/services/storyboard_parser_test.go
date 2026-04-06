package services

import (
	"fmt"
	"testing"
)

func TestParseMarkedScript_FullProductionDoc(t *testing.T) {
	// Simulate the user's full production document pasted as-is
	script := `# SCRIPT: 3 BIG Questions — Why Do We Pray?

**Kênh:** Minno Laugh and Grow Bible
**Format:** [D] HYBRID — AI Director Marked Script
**Thời lượng mục tiêu:** 4 phút (~37 shots × 6.5s trung bình)
**Cú pháp:** Narrator (không tag) | ` + "`[Tên]`" + ` dialogue | ` + "`[CROWD]`" + ` | ` + "`[SFX]`" + `

---

## STRUCTURE MAP

| Phần | Shots | Timecode | Nội dung |
|------|-------|----------|----------|
| INTRO | 01–05 | 00:00–00:30 | Brand ident + câu hỏi trung tâm |

---

## FULL MARKED SCRIPT

---

### PHẦN 1 — INTRO (00:00 – 00:30)

---

**// SHOT 01 | 00:00–00:06 | 6s | BRAND IDENT**

[SFX] Bubbles popping — pop pop pop pop pop

---

**// SHOT 02 | 00:06–00:12 | 6s | NARRATED DIALOGUE**

[Hopeful World] Chào các bạn nhỏ! Mình là Hopeful World đây!

[SFX] Bell ding

---

**// SHOT 03 | 00:12–00:18 | 6s | NARRATED DIALOGUE**

Hôm nay chúng mình có một câu hỏi RẤT LỚN cần giải đáp!

[SFX] Da-dum musical sting

---

**// SHOT 07 | 00:36–00:42 | 6s | HYBRID EXPLAINER**

[Hopeful World] Các bạn sẵn sàng chưa?

[CROWD] Sẵn sàng ạ!!!

[SFX] Cheering — yayyy!

---

### PHẦN 5 — RECAP + SONG (03:10 – 04:00)

---

**// SHOT 32 | 03:24–03:32 | 8s | SONG — MUSICAL RECAP — OPEN**

[NARRATOR] Talk to God in every way — He hears your song every day!

[SFX] Upbeat music bursts in

[CROWD] Yeahhh!

---

## PRODUCTION NOTES

` + "```" + `
Shot duration:   5–8s (avg 6.5s)
Shots per minute: ~9
` + "```" + `
`

	segments, analysis := parseMarkedScript(script)

	if analysis == "" {
		t.Fatal("Expected analysis text, got empty (parser didn't detect any tags)")
	}

	// Count segment types
	counts := map[string]int{}
	for _, seg := range segments {
		counts[seg.Type]++
	}

	fmt.Println("=== PARSED SEGMENTS ===")
	for _, seg := range segments {
		fmt.Printf("[%s] (%s) %s\n", seg.Type, seg.Character, seg.Text)
	}

	fmt.Println("\n=== COUNTS ===")
	fmt.Printf("Narrator: %d\n", counts["narrator"])
	fmt.Printf("Dialogue: %d\n", counts["dialogue"])  
	fmt.Printf("Crowd:    %d\n", counts["crowd"])
	fmt.Printf("SFX:      %d\n", counts["sfx"])

	fmt.Println("\n=== PRE-ANALYSIS TEXT ===")
	fmt.Println(analysis)

	// Verify NO metadata leaked through
	for _, seg := range segments {
		if seg.Text == "HYBRID — AI Director Marked Script" {
			t.Errorf("LEAKED: [D] tag from metadata line was parsed as dialogue")
		}
		if seg.Type == "narrator" {
			// These should never appear as narrator text
			badPatterns := []string{
				"# SCRIPT",
				"STRUCTURE MAP",
				"PRODUCTION NOTES",
				"Shot duration",
				"| Phần |",
				"**//",
			}
			for _, bad := range badPatterns {
				if len(seg.Text) >= len(bad) && seg.Text[:len(bad)] == bad {
					t.Errorf("LEAKED metadata as narrator: %q", seg.Text)
				}
			}
		}
	}

	// Basic sanity checks
	if counts["sfx"] < 5 {
		t.Errorf("Expected at least 5 SFX segments, got %d", counts["sfx"])
	}
	if counts["dialogue"] < 2 {
		t.Errorf("Expected at least 2 dialogue segments, got %d", counts["dialogue"])
	}
	if counts["crowd"] < 2 {
		t.Errorf("Expected at least 2 crowd segments, got %d", counts["crowd"])
	}
}
