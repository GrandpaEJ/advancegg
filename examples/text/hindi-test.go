package main

import (
	"log"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	const W = 1000
	const H = 900
	dc := advancegg.NewContext(W, H)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	err := dc.LoadFontFace("assets/fonts/NotoSans-Regular.ttf", 28)
	if err != nil {
		log.Fatalf("Failed to load Latin font: %v", err)
	}

	err = dc.LoadScriptFont(advancegg.ScriptDevanagari, "/usr/share/fonts/truetype/noto/NotoSansDevanagari-Regular.ttf", 32)
	if err != nil {
		log.Fatalf("Failed to load Devanagari font: %v", err)
	}
	log.Println("✓ Loaded Noto Sans Devanagari")

	dc.SetRGB(0, 0, 0)
	dc.DrawString("Hindi (Devanagari) GPOS Test — Noto Sans Devanagari", 30, 40)
	dc.DrawString("Matras should attach to base consonants, conjuncts should form.", 30, 65)

	// ── All matras with क ──
	y := 100.0
	dc.DrawString("All matras with क (क + vowel sign):", 30, y)
	matras := []struct{ txt, label string }{
		{"का", "ा aa (post-base)"},
		{"कि", "ि i (above-base)"},
		{"की", "ी ii (above-base)"},
		{"कु", "ु u (below-base)"},
		{"कू", "ू uu (below-base)"},
		{"कृ", "ृ ri (below-base)"},
		{"के", "े e (above-base)"},
		{"कै", "ै ai (above-base)"},
		{"को", "ो o (above-base)"},
		{"कौ", "ौ au (above-base)"},
		{"कं", "ं anusvara (above)"},
		{"कः", "ः visarga (post)"},
		{"कँ", "ँ candrabindu (above)"},
	}

	for i, m := range matras {
		ry := y + float64(i)*26 + 30
		dc.DrawStringAnchored(m.txt, 60, ry, 0.5, 0.5)
		dc.DrawString(m.label, 90, ry-6)
	}

	// ── All matras with त ──
	dc.DrawString("Matras with त:", 300, y)
	dc.DrawString("त ता ति ती तु तू तृ ते तै तो तौ तं तः", 300, y+30)

	// ── Matras with प, ब, म, र, ल, व, स, ह ──
	otherBases := []struct{ col, row float64; txt string }{
		{300, y + 52, "प पा पि पी पु पू पृ पे पै पो पौ"},
		{300, y + 72, "ब बा बि बी बु बू बृ बे बै बो बौ"},
		{300, y + 92, "म मा मि मी मु मू मृ मे मै मो मौ"},
		{620, y + 30, "र रा रि री रु रू रृ रे रै रो रौ"},
		{620, y + 52, "ल ला लि ली लु लू लृ ले लै लो लौ"},
		{620, y + 72, "व वा वि वी वु वू वृ वे वै वो वौ"},
		{620, y + 92, "स सा सि सी सु सू सृ से सै सो सौ"},
		{620, y + 112, "ह हा हि ही हु हू हृ है है हो हौ"},
	}
	for _, b := range otherBases {
		dc.DrawString(b.txt, b.col, b.row)
	}

	// ── Conjuncts ──
	y2 := y + 155
	dc.DrawString("Conjuncts (संयुक्त व्यंजन):", 30, y2)
	conjuncts := []string{
		"क्क क्त क्म क्य क्र क्ल क्व क्ष",
		"ग्ग ग्ध ग्न ग्य ग्र ग्ल ग्व घ्न",
		"च्च च्छ च्ञ ज्ज ज्ञ ज्य ज्र",
		"ट्ट त्त त्न त्म त्य त्र त्व",
		"द्द द्ध द्न द्म द्य द्र द्व",
		"न्न न्त न्थ न्द न्ध न्म न्य",
		"प्त प्न प्य प्र प्ल प्स फ्र",
		"ब्द ब्ध ब्न ब्य ब्र ब्ल भ्र",
		"म्न म्प म्ब म्म म्य म्र म्ल",
		"य्य र्ग र्घ र्च र्ज र्ण र्त",
		"र्त र्द र्ध र्न र्प र्म र्स",
		"ल्क ल्ग ल्प ल्ब ल्ल ल्य श्च",
		"ष्क ष्ट ष्ठ ष्ण ष्प स्क स्त",
		"स्थ स्न स्प स्म स्य स्व ह्न",
		"ह्म ह्य ह्र श्र त्त्व त्स्न",
	}
	for i, c := range conjuncts {
		ry := y2 + float64(i)*22 + 28
		dc.DrawString(c, 30, ry)
	}

	// ── Words & sentences ──
	y3 := y2 + float64(len(conjuncts))*22 + 28 + 20
	dc.DrawString("Hindi words and sentences:", 30, y3)

	words := []string{
		"नमस्ते हिन्दी भाषा स्वागत है",
		"शिक्षा विद्यालय पुस्तकालय अध्यापक",
		"ग्राम नगर राष्ट्र विश्व संसार",
		"सूर्य चन्द्र तारा आकाश पृथ्वी",
		"ज्ञान विज्ञान प्रौद्योगिकी",
		"हृदय प्रेम मैत्री शान्ति सत्य",
		"स्वतन्त्रता समानता बन्धुत्व",
		"कृषि उद्योग वाणिज्य अर्थव्यवस्था",
		"संविधान लोकतन्त्र गणतन्त्र",
		"हिन्दी भारत की राजभाषा है।",
		"सभी मनुष्य जन्म से स्वतन्त्र हैं।",
		"शिक्षा ही सबसे बड़ा धन है।",
	}

	for i, w := range words {
		rx := 30 + float64(i%2)*460
		ry := y3 + float64(i/2)*26 + 28
		dc.DrawString(w, rx, ry)
	}

	outputPath := "images/text/hindi_test.png"
	err = dc.SavePNG(outputPath)
	if err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}
	log.Printf("✓ Created Hindi test image: %s", outputPath)
}
