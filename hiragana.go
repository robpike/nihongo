// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nihongo

import (
	"bytes"
	"io"
)

// hiragana implements transliteration of romaji to hiragana.
type hiragana struct {
	t *translator
}

// Hiragana translates romaji into hiragana and returns the result.
func Hiragana(romaji []byte) []byte {
	var buf bytes.Buffer
	r := hiragana{
		t: newTranslator(bytesGetter(romaji), bufPutter(&buf), nil),
	}
	translateHiragana(r.t)
	return buf.Bytes()
}

// HiraganaString translates romaji into hiragana and returns the result.
func HiraganaString(romaji string) string {
	var buf bytes.Buffer
	h := hiragana{
		t: newTranslator(stringGetter(romaji), bufPutter(&buf), nil),
	}
	translateHiragana(h.t)
	return buf.String()
}

// HiraganaReader returns an io.Reader that will translate romaji in its input into hiragana.
func HiraganaReader(rd io.Reader) io.Reader {
	ch := make(chan byte, 100)
	h := &hiragana{
		t: newTranslator(readerGetter(rd), chanPutter(ch), ch),
	}
	go translateHiragana(h.t)
	return h
}

func (h *hiragana) Read(p []byte) (int, error) {
	return h.t.Read(p)
}

func translateHiragana(t *translator) {
	prevByte := -1
	mark := func() {
		if isConsonant[prevByte] {
			t.putString("っ")
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
			cha, ok := threeH[s]
			if ok {
				mark()
				prevByte = -1
				t.putString(cha)
				t.advance(3)
				continue
			}
		}
		if len(s) >= 2 {
			ka, ok := twoH[s[:2]]
			if ok {
				mark()
				t.putString(ka)
				t.advance(2)
				continue
			}
		}
		a, ok := oneH[s[:1]]
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

// Note the absence of n and m.
var isConsonant = map[int]bool{
	'b': true,
	'c': true,
	'd': true,
	'f': true,
	'g': true,
	'h': true,
	'j': true,
	'k': true,
	'l': true,
	'p': true,
	'q': true,
	'r': true,
	's': true,
	't': true,
	'v': true,
	'w': true,
	'x': true,
	'y': true,
	'z': true,
}

var oneH = map[string]string{
	"a": "あ",
	"i": "い",
	"u": "う",
	"e": "え",
	"o": "お",
	"n": "ん",
}

var twoH = map[string]string{
	"xa": "ぁ",
	"xi": "ぃ",
	"xu": "ぅ",
	"xe": "ぇ",
	"xo": "ぉ",

	"ka": "か",
	"ki": "き",
	"ku": "く",
	"ke": "け",
	"ko": "こ",

	"ga": "が",
	"gi": "ぎ",
	"gu": "ぐ",
	"ge": "げ",
	"go": "ご",

	"sa": "さ",
	"su": "す",
	"se": "せ",
	"so": "そ",

	"za": "ざ",
	"ji": "じ",
	"zu": "ず",
	"ze": "ぜ",
	"zo": "ぞ",

	"ja": "じゃ",
	"ju": "じゅ",
	"jo": "じょ",

	"ta": "た",
	"te": "て",
	"to": "と",

	"ti": "てぃ",
	"tu": "とぅ",

	"di": "でぃ",
	"du": "どぅ",

	"da": "だ",
	"de": "で",
	"do": "ど",

	"na": "な",
	"ni": "に",
	"nu": "ぬ",
	"ne": "ね",
	"no": "の",

	"ha": "は",
	"hi": "ひ",
	"fu": "ふ",
	"he": "へ",
	"ho": "ほ",

	"ba": "ば",
	"bi": "び",
	"bu": "ぶ",
	"be": "べ",
	"bo": "ぼ",

	"pa": "ぱ",
	"pi": "ぴ",
	"pu": "ぷ",
	"pe": "ぺ",
	"po": "ぽ",

	"fa": "ふぁ",
	"fi": "ふぃ",
	"fe": "ふぇ",
	"fo": "ふぉ",

	"ma": "ま",
	"mi": "み",
	"mu": "む",
	"me": "め",
	"mo": "も",

	"ya": "や",
	"yu": "ゆ",
	"ye": "いぇ",
	"yo": "よ",

	"ra": "ら",
	"ri": "り",
	"ru": "る",
	"re": "れ",
	"ro": "ろ",

	"wa": "わ",
	//"wi": "ゐ",
	//"we": "ゑ",
	"wo": "を",

	//"wa": "うぁ",
	"wi": "うぃ",
	"we": "うぇ",
	//"wo": "うぉ",

	"va": "ゔぁ",
	"vi": "ゔぃ",
	"vu": "ゔ",
	"ve": "ゔぇ",
	"vo": "ゔぉ",
}

var threeH = map[string]string{
	"kya": "きゃ",
	"kyu": "きゅ",
	"kyo": "きょ",

	"gya": "ぎゃ",
	"gyu": "ぎゅ",
	"gyo": "ぎょ",

	"shi": "し",

	"sha": "しゃ",
	"shu": "しゅ",
	"sho": "しょ",

	"chi": "ち",
	"tsu": "つ",

	"dhi": "ぢ",
	"dhu": "づ",

	"cha": "ちゃ",
	"chu": "ちゅ",
	"che": "ちぇ",
	"cho": "ちょ",

	"dha": "ぢゃ",
	//"dhu": "ぢゅ",
	"dhe": "ぢぇ",
	"dho": "ぢょ",

	"nya": "にゃ",
	"nyu": "にゅ",
	"nyo": "にょ",

	"hya": "ひゃ",
	"hyu": "ひゅ",
	"hyo": "ひょ",

	"bya": "びゃ",
	"byu": "びゅ",
	"byo": "びょ",

	"pya": "ぴゃ",
	"pyu": "ぴゅ",
	"pyo": "ぴょ",

	"mya": "みゃ",
	"myu": "みゅ",
	"myo": "みょ",

	"xya": "ゃ",
	"xyu": "ゅ",
	"xyo": "ょ",

	"rya": "りゃ",
	"ryu": "りゅ",
	"ryo": "りょ",
}
