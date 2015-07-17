// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nihongo

import (
	"bytes"
	"io"
)

// romaji implements transliteration to romaji.
type romaji struct {
	t *translator
}

// Romaji translates text into romaji and returns the result.
func Romaji(text []byte) []byte {
	var buf bytes.Buffer
	r := romaji{
		t: newTranslator(bytesGetter(text), bufPutter(&buf), nil),
	}
	translateRomaji(r.t)
	return buf.Bytes()
}

// RomajiString translates text into romaji and returns the result.
func RomajiString(text string) string {
	var buf bytes.Buffer
	r := romaji{
		t: newTranslator(stringGetter(text), bufPutter(&buf), nil),
	}
	translateRomaji(r.t)
	return buf.String()
}

// RomajiReader returns an io.Reader that will translate its input into romaji.
func RomajiReader(rd io.Reader) io.Reader {
	ch := make(chan byte, 100)
	r := &romaji{
		t: newTranslator(readerGetter(rd), chanPutter(ch), ch),
	}
	go translateRomaji(r.t)
	return r
}

func (r *romaji) Read(p []byte) (int, error) {
	return r.t.Read(p)
}

func translateRomaji(t *translator) {
	prevKana := false
	skip := false
	for first := true; ; first = false {
		r := t.next()
		if r == eof {
			break
		}
		k, ok := kana[r]
		if skip {
			skip = false
			continue
		}
		if !ok {
			if prevKana {
				t.put(' ')
			}
			t.putRune(r)
			prevKana = false
			continue
		}
		if !first && !prevKana {
			t.put(' ')
		}
		prevKana = true
		// Is there a modifier?
		if small[t.peek()] {
			skip = true
			r2 := t.next()
			k2, ok := mod[r2]
			if ok {
				t.putString(k[:len(k)-1])
				t.putString(k2[1:])
				continue
			}
			k2, ok = vowel[r2]
			if ok {
				t.putString(k)
				t.put('-')
				continue
			}
			// Otherwise it's just odd.
			t.put('<')
			t.putString(k)
			t.put('.')
			t.putString(odd[r2])
			t.put('>')
			continue
		}
		t.putString(k)
	}
	if t.ch != nil {
		close(t.ch)
	}
}

var kana = map[rune]string{
	'あ': "a",
	'い': "i",
	'う': "u",
	'え': "e",
	'お': "o",
	'か': "ka",
	'が': "ga",
	'き': "ki",
	'ぎ': "gi",
	'く': "ku",
	'ぐ': "gu",
	'け': "ke",
	'げ': "ge",
	'こ': "ko",
	'ご': "go",
	'さ': "sa",
	'ざ': "za",
	'し': "shi",
	'じ': "zi",
	'す': "su",
	'ず': "zu",
	'せ': "se",
	'ぜ': "ze",
	'そ': "so",
	'ぞ': "zo",
	'た': "ta",
	'だ': "da",
	'ち': "chi",
	'ぢ': "di",
	'つ': "tsu",
	'づ': "du",
	'て': "te",
	'で': "de",
	'と': "to",
	'ど': "do",
	'な': "na",
	'に': "ni",
	'ぬ': "nu",
	'ね': "ne",
	'の': "no",
	'は': "ha",
	'ば': "va",
	'ぱ': "pa",
	'ひ': "hi",
	'び': "vi",
	'ぴ': "pi",
	'ふ': "fu",
	'ぶ': "bu",
	'ぷ': "pu",
	'へ': "he",
	'べ': "ve",
	'ぺ': "pe",
	'ほ': "ho",
	'ぼ': "vo",
	'ぽ': "po",
	'ま': "ma",
	'み': "mi",
	'む': "mu",
	'め': "me",
	'も': "mo",
	'や': "ya",
	'ゆ': "yu",
	'よ': "yo",
	'ら': "ra",
	'り': "ri",
	'る': "ru",
	'れ': "re",
	'ろ': "ro",
	'わ': "wa",
	'ゐ': "wi",
	'ゑ': "we",
	'を': "wo",
	'ん': "n",
	'ゔ': "vu",
	'ア': "a",
	'イ': "i",
	'ウ': "u",
	'エ': "e",
	'オ': "o",
	'カ': "ka",
	'ガ': "ga",
	'キ': "ki",
	'ギ': "gi",
	'ク': "ku",
	'グ': "gu",
	'ケ': "ke",
	'ゲ': "ge",
	'コ': "ko",
	'ゴ': "go",
	'サ': "sa",
	'ザ': "za",
	'シ': "shi",
	'ジ': "zi",
	'ス': "su",
	'ズ': "zu",
	'セ': "se",
	'ゼ': "ze",
	'ソ': "so",
	'ゾ': "zo",
	'タ': "ta",
	'ダ': "da",
	'チ': "chi",
	'ヂ': "di",
	'ツ': "tsu",
	'ヅ': "du",
	'テ': "te",
	'デ': "de",
	'ト': "to",
	'ド': "do",
	'ナ': "na",
	'ニ': "ni",
	'ヌ': "nu",
	'ネ': "ne",
	'ノ': "no",
	'ハ': "ha",
	'バ': "va",
	'パ': "pa",
	'ヒ': "hi",
	'ビ': "vi",
	'ピ': "pi",
	'フ': "fu",
	'ブ': "bu",
	'プ': "pu",
	'ヘ': "he",
	'ベ': "ve",
	'ペ': "pe",
	'ホ': "ho",
	'ボ': "vo",
	'ポ': "po",
	'マ': "ma",
	'ミ': "mi",
	'ム': "mu",
	'メ': "me",
	'モ': "mo",
	'ヤ': "ya",
	'ユ': "yu",
	'ヨ': "yo",
	'ラ': "ra",
	'リ': "ri",
	'ル': "ru",
	'レ': "re",
	'ロ': "ro",
	'ワ': "wa",
	'ヰ': "wi",
	'ヱ': "we",
	'ヲ': "wo",
	'ン': "n",
	'ヴ': "vu",
}

var small = map[rune]bool{
	'ぁ': true,
	'ぃ': true,
	'ぅ': true,
	'ぇ': true,
	'ぉ': true,
	'っ': true,
	'ゃ': true,
	'ゅ': true,
	'ょ': true,
	'ゎ': true,
	'ゕ': true,
	'ゖ': true,

	'ァ': true,
	'ィ': true,
	'ゥ': true,
	'ェ': true,
	'ォ': true,
	'ッ': true,
	'ャ': true,
	'ュ': true,
	'ョ': true,
	'ヮ': true,
	'ヵ': true,
	'ヶ': true,
}

var vowel = map[rune]string{
	'ぁ': "a",
	'ぃ': "i",
	'ぅ': "u",
	'ぇ': "e",
	'ぉ': "o",

	'ァ': "a",
	'ィ': "i",
	'ゥ': "u",
	'ェ': "e",
	'ォ': "o",
}

var mod = map[rune]string{
	'ゃ': "ya",
	'ゅ': "yu",
	'ょ': "yo",

	'ャ': "ya",
	'ュ': "yu",
	'ョ': "yo",
}

var odd = map[rune]string{
	'っ': "hold",  // tsu == hold consonant
	'ゕ': "count", // ka == counting mark
	'ゖ': "count", // ke == counting mark
	'ッ': "hold",  // tsu == hold consonant
	'ヵ': "count", // ka == counting mark
	'ヶ': "count", // ke == counting mark
}
