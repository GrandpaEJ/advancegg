package main

import (
	"log"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	const W = 1200
	const H = 1600
	dc := advancegg.NewContext(W, H)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	err := dc.LoadFontFace("assets/fonts/NotoSans-Regular.ttf", 28)
	if err != nil {
		log.Fatalf("Failed to load main font: %v", err)
	}

	err = dc.LoadScriptFont(advancegg.ScriptBengali, "assets/fonts/NotoSansBengali-Regular.ttf", 36)
	if err != nil {
		log.Fatalf("Failed to load Bengali font: %v", err)
	}
	log.Println("✓ Bengali script font loaded")

	dc.SetRGB(0, 0, 0)

	// ── Title ──
	dc.DrawString("Bengali GPOS Test — All Matras + 25 Conjuncts", 30, 40)
	dc.DrawString("(If matras float separately, GPOS is broken)", 30, 70)

	y := 120.0
	col1 := 50.0
	col2 := 430.0
	col3 := 810.0

	// ══════════════════════════════════════════════
	// SECTION 1: ALL 11 MATRAS (vowel signs)
	// ══════════════════════════════════════════════
	dc.DrawString("─ All 11 Matras (vowel signs) with ক ─", col1, y)

	matras := []struct {
		sign  string
		label string
	}{
		{"কা", "া a-kar (post-base)"},
		{"কি", "ি i-kar (above-left)"},
		{"কী", "ী dirgho i-kar (above)"},
		{"কু", "ু u-kar (below)"},
		{"কূ", "ূ dirgho u-kar (below)"},
		{"কৃ", "ৃ ri-kar (below)"},
		{"কে", "ে e-kar (pre-base)"},
		{"কৈ", "ৈ oi-kar (pre-base)"},
		{"কো", "ো o-kar (pre+post)"},
		{"কৌ", "ৌ ou-kar (pre+post)"},
		{"কং", "ং anusvara (above)"},
		{"কঃ", "ঃ visarga (post)"},
		{"কঁ", "ঁ candrabindu (above)"},
	}

	for i, m := range matras {
		ry := y + float64(i)*36
		dc.DrawStringAnchored(m.sign, col1+40, ry, 0.5, 0.5)
		dc.DrawString(m.label, col1+80, ry-8)
	}

	// ══════════════════════════════════════════════
	// SECTION 2: MATRAS with different consonant shapes
	// ══════════════════════════════════════════════
	y2 := y
	dc.DrawString("─ Matras across different bases ─", col2, y2)

	matraRows := []struct {
		base string
		row  float64
	}{
		{"ক কা কি কী কু কূ কৃ কে কৈ কো কৌ", 0},
		{"খ খা খি খী খু খূ খৃ খে খৈ খো খৌ", 1},
		{"গ গা গি গী গু গূ গৃ গে গৈ গো গৌ", 2},
		{"ঘ ঘা ঘি ঘী ঘু ঘূ ঘৃ ঘে ঘৈ ঘো ঘৌ", 3},
		{"চ চা চি চী চু চূ চৃ চে চৈ চো চৌ", 4},
		{"ছ ছা ছি ছী ছু ছূ ছৃ ছে ছৈ ছো ছৌ", 5},
		{"জ জা জি জী জু জূ জৃ জে জৈ জো জৌ", 6},
		{"ঝ ঝা ঝি ঝী ঝু ঝূ ঝৃ ঝে ঝৈ ঝো ঝৌ", 7},
		{"ট টা টি টী টু টূ টৃ টে টৈ টো টৌ", 8},
		{"ঠ ঠা ঠি ঠী ঠু ঠূ ঠৃ ঠে ঠৈ ঠো ঠৌ", 9},
		{"ড ডা ডি ডী ডু ডূ ডৃ ডে ডৈ ডো ডৌ", 10},
		{"ণ ণা ণি ণী ণু ণূ ণৃ ণে ণৈ ণো ণৌ", 11},
		{"ত তা তি তী তু তূ তৃ তে তৈ তো তৌ", 12},
		{"থ থা থি থী থু থূ থৃ থে থৈ থো থৌ", 13},
		{"দ দা দি দী দু দূ দৃ দে দৈ দো দৌ", 14},
		{"ধ ধা ধি ধী ধু ধূ ধৃ ধে ধৈ ধো ধৌ", 15},
		{"ন না নি নী নু নূ নৃ নে নৈ নো নৌ", 16},
		{"প পা পি পী পু পূ পৃ পে পৈ পো পৌ", 17},
		{"ফ ফা ফি ফী ফু ফূ ফৃ ফে ফৈ ফো ফৌ", 18},
		{"ব বা বি বী বু বূ বৃ বে বৈ বো বৌ", 19},
		{"ভ ভা ভি ভী ভু ভূ ভৃ ভে ভৈ ভো ভৌ", 20},
		{"ম মা মি মী মু মূ মৃ মে মৈ মো মৌ", 21},
		{"শ শা শি শী শু শূ শৃ শে শৈ শো শৌ", 22},
		{"ষ ষা ষি ষী ষু ষূ ষৃ ষে ষৈ ষো ষৌ", 23},
		{"স সা সি সী সু সূ সৃ সে সৈ সো সৌ", 24},
		{"হ হা হি হী হু হূ হৃ হে হৈ হো হৌ", 25},
		{"ড় ড়া ড়ি ড়ী ড়ু ড়ূ ড়ৃ ড়ে ড়ৈ ড়ো ড়ৌ", 26},
		{"ঢ় ঢ়া ঢ়ি ঢ়ী ঢ়ু ঢ়ূ ঢ়ৃ ঢ়ে ঢ়ৈ ঢ়ো ঢ়ৌ", 27},
	}

	for _, mr := range matraRows {
		ry := y2 + mr.row*18 + 36
		dc.DrawString(mr.base, col2, ry)
	}

	// ══════════════════════════════════════════════
	// SECTION 3: 25+ CONJUNCTS (joktoborno)
	// ══════════════════════════════════════════════
	y3 := y + 36
	dc.DrawString("─ 25 Conjuncts (যুক্তবর্ণ) in words ─", col3, y3)

	conjuncts := []struct {
		word  string
		gloss string
	}{
		{"ক্ত", "ক্ত — kta (যুক্ত ক+ত)"},
		{"ক্ষ", "ক্ষ — kṣa (ক+ষ)"},
		{"গ্ধ", "গ্ধ — gdha (গ+ধ)"},
		{"গ্ন", "গ্ন — gna (গ+ন)"},
		{"গ্ব", "গ্ব — gba (গ+ব)"},
		{"গ্ল", "গ্ল — gla (গ+ল)"},
		{"চ্চ", "চ্চ — cca (চ+চ)"},
		{"চ্ছ", "চ্ছ — ccha (চ+ছ)"},
		{"জ্জ", "জ্জ — jja (জ+জ)"},
		{"জ্ঞ", "জ্ঞ — jña (জ+ঞ)"},
		{"ঞ্চ", "ঞ্চ — ñca (ঞ+চ)"},
		{"ণ্ড", "ণ্ড — ṇḍa (ণ+ড)"},
		{"ত্ত", "ত্ত — tta (ত+ত)"},
		{"ত্ত্ব", "ত্ত্ব — ttva (ত+ত+ব)"},
		{"ত্র", "ত্র — tra (ত+র)"},
		{"দ্ধ", "দ্ধ — ddha (দ+ধ)"},
		{"দ্ধ্ব", "দ্ধ্ব — ddhva (দ+ধ+ব)"},
		{"দ্ধ্র", "দ্ধ্র — ddhra (দ+ধ+র)"},
		{"দ্ব", "দ্ব — dba (দ+ব)"},
		{"দ্ধ", "দ্ধ — ddha (দ+ধ)"},
		{"ন্ধ", "ন্ধ — ndha (ন+ধ)"},
		{"ন্ত", "ন্ত — nta (ন+ত)"},
		{"ন্দ", "ন্দ — nda (ন+দ)"},
		{"ন্ধ্র", "ন্ধ্র — ndhra (ন+ধ+র)"},
		{"প্ত", "প্ত — pta (প+ত)"},
		{"প্র", "প্র — pra (প+র)"},
		{"প্ল", "প্ল — pla (প+ল)"},
		{"ব্র", "ব্র — bra (ব+র)"},
		{"ভ্র", "ভ্র — bhra (ভ+র)"},
		{"ম্ব", "ম্ব — mba (ম+ব)"},
		{"ম্প", "ম্প — mpa (ম+প)"},
		{"ম্ভ", "ম্ভ — mbha (ম+ভ)"},
		{"ম্ম", "ম্ম — mma (ম+ম)"},
		{"ম্ল", "ম্ল — mla (ম+ল)"},
		{"য্য", "য্য — yya (য+য)"},
		{"র্ত", "র্ত — rta (র+ত)"},
		{"র্তৃ", "র্তৃ — rtṛ (র+ত+ৃ)"},
		{"র্থ", "র্থ — rtha (র+থ)"},
		{"র্দ", "র্দ — rda (র+দ)"},
		{"র্ধ", "র্ধ — rdha (র+ধ)"},
		{"র্ণ", "র্ণ — rṇa (র+ণ)"},
		{"র্ম", "র্ম — rma (র+ম)"},
		{"র্শ", "র্শ — rśa (র+শ)"},
		{"র্শ্ব", "র্শ্ব — rśba (র+শ+ব)"},
		{"র্ষ", "র্ষ — rṣa (র+ষ)"},
		{"র্হ", "র্হ — rha (র+হ)"},
		{"ল্প", "ল্প — lpa (ল+প)"},
		{"ল্প্র", "ল্প্র — lpra (ল+প+র)"},
		{"ল্ব", "ল্ব — lba (ল+ব)"},
		{"শ্চ", "শ্চ — śca (শ+চ)"},
		{"শ্ছ", "শ্ছ — ścha (শ+ছ)"},
		{"শ্ন", "শ্ন — śna (শ+ন)"},
		{"শ্ব", "শ্ব — śba (শ+ব)"},
		{"শ্ম", "শ্ম — śma (শ+ম)"},
		{"শ্র", "শ্র — śra (শ+র)"},
		{"ষ্ট", "ষ্ট — ṣṭa (ষ+ট)"},
		{"ষ্ঠ", "ষ্ঠ — ṣṭha (ষ+ঠ)"},
		{"ষ্ণ", "ষ্ণ — ṣṇa (ষ+ণ)"},
		{"ষ্প", "ষ্প — ṣpa (ষ+প)"},
		{"ষ্ফ", "ষ্ফ — ṣpha (ষ+ফ)"},
		{"ষ্ম", "ষ্ম — ṣma (ষ+ম)"},
		{"স্ক", "স্ক — ska (স+ক)"},
		{"স্ট", "স্ট — sṭa (স+ট)"},
		{"স্ত", "স্ত — sta (স+ত)"},
		{"স্ত্র", "স্ত্র — stra (স+ত+র)"},
		{"স্থ", "স্থ — stha (স+থ)"},
		{"স্ন", "স্ন — sna (স+ন)"},
		{"স্প", "স্প — spa (স+প)"},
		{"স্প্র", "স্প্র — spra (স+প+র)"},
		{"স্প্ল", "স্প্ল — spla (স+প+ল)"},
		{"স্ফ", "স্ফ — spha (স+ফ)"},
		{"স্ব", "স্ব — sba (স+ব)"},
		{"স্ম", "স্ম — sma (স+ম)"},
		{"স্র", "স্র — sra (স+র)"},
		{"স্ল", "স্ল — sla (স+ল)"},
		{"হ্ন", "হ্ন — hna (হ+ন)"},
		{"হ্ম", "হ্ম — hma (হ+ম)"},
		{"হ্র", "হ্র — hra (হ+র)"},
		{"হ্ল", "হ্ল — hla (হ+ল)"},
	}

	for i, c := range conjuncts {
		ry := y3 + float64(i)*18 + 72
		dc.DrawStringAnchored(c.word, col3+30, ry, 0.5, 0.5)
		dc.DrawString(c.gloss, col3+65, ry-5)
	}

	// ══════════════════════════════════════════════
	// SECTION 4: REAL BENGALI WORDS (exercise multiple features)
	// ══════════════════════════════════════════════
	y4 := y3 + float64(len(conjuncts))*18 + 72 + 40

	dc.DrawString("─ Real Bengali words (mixed matras + conjuncts) ─", col1, y4)

	words := []string{
		"বাংলাদেশ",      // Bangladesh
		"ভাষা",          // language
		"আমার",          // my
		"সুন্দর",        // beautiful
		"মাত্রা",        // matra
		"যুক্তবর্ণ",     // conjunct
		"বিদ্যালয়",     // school
		"সূর্য",         // sun
		"গ্রহ",          // planet
		"বিশ্ব",         // world
		"শিক্ষা",        // education
		"সৃষ্টি",        // creation
		"গ্রহণ",         // accept
		"মুক্তি",        // freedom
		"ব্যক্তি",       // person
		"সত্য",          // truth
		"গ্রন্থ",        // book
		"চন্দ্র",        // moon
		"ক্ষুধা",        // hunger
		"শ্রম",          // labor
		"হৃদয়",        // heart
		"কর্ণ",          // ear
		"মর্ম",          // essence
		"বর্ষা",         // rain
		"কষ্ট",          // pain
		"সপ্তাহ",       // week
		"প্রথম",        // first
		"বৃদ্ধি",        // growth
		"স্নান",         // bath
		"গ্রাম",         // village
	}

	for i, w := range words {
		rx := col1 + float64(i%5)*220
		ry := y4 + float64(i/5)*32 + 36
		dc.DrawStringAnchored(w, rx+40, ry, 0.5, 0.5)
	}

	// ══════════════════════════════════════════════
	// SECTION 5: Sentence-level test
	// ══════════════════════════════════════════════
	y5 := y4 + float64((len(words)+4)/5)*32 + 36 + 20

	sentences := []string{
		"বাংলা ভাষা একটি সুন্দর ভাষা।",
		"আমার বিদ্যালয়ে প্রতিদিন পাঠদান হয়।",
		"সূর্য পূর্ব দিকে উদয় হয় এবং পশ্চিমে অস্ত যায়।",
		"শিক্ষক শিক্ষার্থীদের পাঠদান করাচ্ছেন।",
		"গ্রন্থাগারে বহু মূল্যবান গ্রন্থ রয়েছে।",
		"বৃষ্টির পরে আকাশে রংধনু দেখা যায়।",
		"আমার মা রান্নাঘরে সুস্বাদু খাবার রান্না করছেন।",
		"পৃথিবী সূর্যের চারদিকে ঘুরছে।",
		"বিজ্ঞান ও প্রযুক্তির অগ্রগতি বিশ্বকে বদলে দিয়েছে।",
		"শ্রমিকরা সকাল থেকে সন্ধ্যা পর্যন্ত কাজ করে।",
		"চন্দ্র ও সূর্যের গ্রহণ একটি প্রাকৃতিক ঘটনা।",
		"হৃদয়ের কথা কেউ কি বুঝতে পারে?",
		"মুক্তিযুদ্ধে বাঙালি জাতির আত্মত্যাগ চিরস্মরণীয়।",
		"বাংলা সাহিত্যে রবীন্দ্রনাথ ঠাকুর একজন শ্রেষ্ঠ কবি।",
		"শিক্ষা ও সংস্কৃতি একটি জাতির প্রাণ।",
	}

	for i, s := range sentences {
		ry := y5 + float64(i)*30
		dc.DrawString(s, col1, ry)
	}

	outputPath := "images/text/bengali_gpos_test.png"
	err = dc.SavePNG(outputPath)
	if err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	log.Printf("✓ Created Bengali GPOS test image: %s", outputPath)
}
