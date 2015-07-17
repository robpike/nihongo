// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package nihongo

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"testing"
)

type testPair struct {
	in  string
	out string
}

func testString(name string, t *testing.T, test testPair, tr func(string) string) {
	result := tr(test.in)
	if result != test.out {
		t.Errorf("string %s: expected %q got %q\n", name, test.out, result)
	}
}

func testBytes(name string, t *testing.T, test testPair, tr func([]byte) []byte) {
	result := string(tr([]byte(test.in)))
	if result != test.out {
		t.Errorf("byte %s: expected %q got %q\n", name, test.out, result)
	}
}

func testReader(name string, t *testing.T, test testPair, tr func(io.Reader) io.Reader) {
	data, err := ioutil.ReadAll(tr(strings.NewReader(test.in)))
	if err != nil {
		t.Errorf("reader %s: v\n", name, err)
		return
	}
	result := string(data)
	if result != test.out {
		t.Errorf("reader %s: expected %q got %q\n", name, test.out, result)
	}
}

var romajiTests = []testPair{
	// Unchanged.
	{"", ""},
	{"now is the time\n", "now is the time\n"},
	{"ひらがな", "hiragana"},
	{"カタカナ", "katakana"},
	// Leave non-kana alone.
	{"a日本語ひらがなカタカナb\n", "a日本語 hiraganakatakana b\n"},
}

func TestRomaji(t *testing.T) {
	for i, test := range romajiTests {
		name := fmt.Sprintf("#%d: romaji:", i)
		testString(name, t, test, RomajiString)
		testBytes(name, t, test, Romaji)
		testReader(name, t, test, RomajiReader)
	}
}

var hiraganaTests = []testPair{
	// Unchanged.
	{"", ""},
	{"xxqq", "xxqq"},
	{"hiragana", "ひらがな"},
	// Multi-kana inputs.
	{"chachodhowi", "ちゃちょぢょうぃ"},
	// Consonant modifier (tsu).
	{"tcho", "っちょ"},
	{"atcho", "あっちょ"},
	// Plosive
	{"shuppatu", "しゅっぱつ"},
	{"syuppatu", "しゅっぱつ"},
	{"syuppatsu", "しゅっぱつ"},
	// Leave existing kana alone.
	{"カタカナ日本語カタカナひらがなカタカナ\n", "カタカナ日本語カタカナひらがなカタカナ\n"},
}

func TestHiragana(t *testing.T) {
	for i, test := range hiraganaTests {
		name := fmt.Sprintf("#%d: hiragana:", i)
		testString(name, t, test, HiraganaString)
		testBytes(name, t, test, Hiragana)
		testReader(name, t, test, HiraganaReader)
	}
}

var katakanaTests = []testPair{
	// Unchanged.
	{"", ""},
	{"xxqq", "xxqq"},
	{"katakana", "カタカナ"},
	// Leave existing kana alone.
	{"カタカナ日本語カタカナひらがなカタカナ\n", "カタカナ日本語カタカナひらがなカタカナ\n"},
}

func TestKatakana(t *testing.T) {
	for i, test := range katakanaTests {
		name := fmt.Sprintf("#%d: katakana:", i)
		testString(name, t, test, KatakanaString)
		testBytes(name, t, test, Katakana)
		testReader(name, t, test, KatakanaReader)
	}
}
