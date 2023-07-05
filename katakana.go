// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nihongo

import (
	"bytes"
	"io"
)

// katakana implements transliteration of romaji to katakana.
type katakana struct {
	t *translator
}

// Katakana translates romaji into katakana and returns the result.
func Katakana(romaji []byte) []byte {
	var buf bytes.Buffer
	r := katakana{
		t: newTranslator(bytesGetter(romaji), bufPutter(&buf), nil),
	}
	translateKatakana(r.t)
	return buf.Bytes()
}

// KatakanaString translates romaji into katakana and returns the result.
func KatakanaString(romaji string) string {
	var buf bytes.Buffer
	k := katakana{
		t: newTranslator(stringGetter(romaji), bufPutter(&buf), nil),
	}
	translateKatakana(k.t)
	return buf.String()
}

// KatakanaReader returns an io.Reader that will translate romaji in its input into katakana.
func KatakanaReader(rd io.Reader) io.Reader {
	ch := make(chan byte, 100)
	k := &katakana{
		t: newTranslator(readerGetter(rd), chanPutter(ch), ch),
	}
	go translateKatakana(k.t)
	return k
}

func (k *katakana) Read(p []byte) (int, error) {
	return k.t.Read(p)
}

func translateKatakana(t *translator) {
	prevByte := -1
	mark := func() {
		if isConsonant[prevByte] {
			t.putString("ッ")
		} else if prevByte >= 0 {
			t.put(byte(prevByte))
		}
		prevByte = -1
	}
	for {
		s := t.next3()
		if len(s) == 0 {
			break
		}
		if len(s) == 3 {
			cha, ok := threeK[s]
			if ok {
				mark()
				prevByte = -1
				t.putString(cha)
				t.advance(3)
				continue
			}
		}
		if len(s) >= 2 {
			ka, ok := twoK[s[:2]]
			if ok {
				mark()
				t.putString(ka)
				t.advance(2)
				continue
			}
		}
		a, ok := oneK[s[:1]]
		if ok {
			mark()
			t.putString(a)
			t.advance(1)
			continue
		}
		if prevByte >= 0 {
			t.put(byte(prevByte))
		}
		prevByte = int(s[0])
		t.advance(1)
	}
	if prevByte >= 0 {
		t.put(byte(prevByte))
	}
	if t.ch != nil {
		close(t.ch)
	}
}

var oneK = map[string]string{
	"a": "ア",
	"i": "イ",
	"u": "ウ",
	"e": "エ",
	"o": "オ",
	"n": "ン",
}

var twoK = map[string]string{
	"xa": "ァ",
	"xi": "ィ",
	"xu": "ゥ",
	"xe": "ェ",
	"xo": "ォ",
	"la": "ァ",
	"li": "ィ",
	"lu": "ゥ",
	"le": "ェ",
	"lo": "ォ",

	"ka": "カ",
	"ki": "キ",
	"ku": "ク",
	"ke": "ケ",
	"ko": "コ",

	"ga": "ガ",
	"gi": "ギ",
	"gu": "グ",
	"ge": "ゲ",
	"go": "ゴ",

	"sa": "サ",
	"si": "シ",
	"su": "ス",
	"se": "セ",
	"so": "ソ",

	"za": "ザ",
	"ji": "ジ",
	"zi": "ジ",
	"zu": "ズ",
	"ze": "ゼ",
	"zo": "ゾ",

	"ja": "ジャ",
	"ju": "ジュ",
	"jo": "ジョ",

	"ta": "タ",
	"te": "テ",
	"to": "ト",

	"ti": "チ",
	"tu": "ツ",

	"di": "ヂ",
	"du": "ヅ",

	"da": "ダ",
	"de": "デ",
	"do": "ド",

	"na": "ナ",
	"ni": "ニ",
	"nu": "ヌ",
	"ne": "ネ",
	"no": "ノ",

	"ha": "ハ",
	"hi": "ヒ",
	"fu": "フ",
	"hu": "フ",
	"he": "へ",
	"ho": "ホ",

	"ba": "バ",
	"bi": "ビ",
	"bu": "ブ",
	"be": "べ",
	"bo": "ボ",

	"pa": "パ",
	"pi": "ピ",
	"pu": "プ",
	"pe": "ペ",
	"po": "ポ",

	"fa": "ファ",
	"fi": "フィ",
	"fe": "フェ",
	"fo": "フォ",

	"ma": "マ",
	"mi": "ミ",
	"mu": "ム",
	"me": "メ",
	"mo": "モ",

	"ya": "ヤ",
	"yu": "ユ",
	"ye": "イェ",
	"yo": "ヨ",

	"ra": "ラ",
	"ri": "リ",
	"ru": "ル",
	"re": "レ",
	"ro": "ロ",

	"wa": "ワ",
	//"wi": "ヰ",
	//"we": "ヱ",
	"wo": "ヲ",

	//"wa": "ウァ",
	"wi": "ウィ",
	"we": "ウェ",
	//"wo": "ウォ",

	"va": "ヴァ",
	"vi": "ヴィ",
	"vu": "ヴ",
	"ve": "ヴェ",
	"vo": "ヴォ",
}

var threeK = map[string]string{
	"kya": "キャ",
	"kyu": "キュ",
	"kyo": "キョ",

	"gya": "ギャ",
	"gyu": "ギュ",
	"gyo": "ギョ",

	"shi": "シ",

	"sha": "シャ",
	"shu": "シュ",
	"sho": "ショ",
	"sya": "シャ",
	"syu": "シュ",
	"syo": "ショ",

	"chi": "チ",
	"tsu": "ツ",

	"dhi": "ディ",
	"dhu": "デュ",
	"dwu": "ドゥ",

	"cha": "チャ",
	"tya": "チャ",
	"chu": "チュ",
	"tyu": "チュ",
	"che": "チェ",
	"tye": "チェ",
	"cho": "チョ",
	"tyo": "チョ",

	"dha": "ヂァ",
	//"dhu": "ヂゥ",
	"dhe": "ヂェ",
	"dho": "ヂョ",

	"nya": "ニャ",
	"nyu": "ニュ",
	"nyo": "ニョ",

	"hya": "ヒャ",
	"hyu": "ヒュ",
	"hyo": "ヒョ",

	"bya": "ビャ",
	"byu": "ビュ",
	"byo": "ビョ",

	"pya": "ピャ",
	"pyu": "ピュ",
	"pyo": "ピョ",

	"mya": "ミャ",
	"myu": "ミュ",
	"myo": "ミョ",

	"xya": "ャ",
	"xyu": "ュ",
	"xyo": "ョ",
	"lya": "ャ",
	"lyu": "ュ",
	"lyo": "ョ",

	"rya": "リャ",
	"ryu": "リュ",
	"ryo": "リョ",
}
