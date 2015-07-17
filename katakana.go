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
	for {
		s := t.next3()
		if len(s) == 0 {
			break
		}
		if len(s) == 3 {
			cha, ok := threeK[s]
			if ok {
				t.putString(cha)
				t.advance(3)
				continue
			}
		}
		if len(s) >= 2 {
			ka, ok := twoK[s[:2]]
			if ok {
				t.putString(ka)
				t.advance(2)
				continue
			}
		}
		a, ok := oneK[s[:1]]
		if ok {
			t.putString(a)
			t.advance(1)
			continue
		}
		t.put(s[0])
		t.advance(1)
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
	"m": "ン",
}

var twoK = map[string]string{
	"ka": "カ",
	"ca": "カ",
	"ga": "ガ",
	"ki": "キ",
	"gi": "ギ",
	"ku": "ク",
	"cu": "ク",
	"gu": "グ",
	"ke": "ケ",
	"ge": "ゲ",
	"ko": "コ",
	"co": "コ",
	"go": "ゴ",
	"sa": "サ",
	"za": "ザ",
	"si": "シ",
	"zi": "ジ",
	"su": "ス",
	"zu": "ズ",
	"se": "セ",
	"ze": "ゼ",
	"so": "ソ",
	"zo": "ゾ",
	"ta": "タ",
	"da": "ダ",
	"ti": "チ",
	"di": "ヂ",
	"tu": "ツ",
	"du": "ヅ",
	"te": "テ",
	"de": "デ",
	"to": "ト",
	"do": "ド",
	"na": "ナ",
	"ni": "ニ",
	"nu": "ヌ",
	"ne": "ネ",
	"no": "ノ",
	"ha": "ハ",
	"va": "バ",
	"pa": "パ",
	"hi": "ヒ",
	"vi": "ビ",
	"pi": "ピ",
	"fe": "フェ",
	"fu": "フ",
	"bu": "ブ",
	"pu": "プ",
	"he": "ヘ",
	"ve": "ベ",
	"pe": "ペ",
	"ho": "ホ",
	"vo": "ボ",
	"po": "ポ",
	"ma": "マ",
	"mi": "ミ",
	"mu": "ム",
	"me": "メ",
	"mo": "モ",
	"ya": "ヤ",
	"yu": "ユ",
	"yo": "ヨ",
	"ra": "ラ",
	"ri": "リ",
	"ru": "ル",
	"re": "レ",
	"ro": "ロ",
	"wa": "ワ",
	"wi": "ヰ",
	"we": "ヱ",
	"wo": "ヲ",
	"vu": "ヴ",
}

var threeK = map[string]string{
	"shi": "シ",
	"chi": "チ",
	"tsu": "ツ",
	"cha": "チャ", // TODO IS THIS RIGHT
	"chu": "チュ", // TODO IS THIS RIGHT
	"cho": "チョ", // TODO IS THIS RIGHT
	// NEED MORE
}
