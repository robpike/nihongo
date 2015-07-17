// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package nihongo implements simple transliteration between romaji
// and the two syllabic Japanese scripts, hiragana and katakana, encoded
// as UTF-8-encoded Unicode. Romaji output may include injected spaces
// to separate converted text from unconverted, and other markers.
// Invalid sequences, such as small kanas with no preceding kana,
// are passed unaltered. Hiragana and katakana may be inaccurate
// due to false matches. Katakana may be further inaccurate because
// of the inability to generate the tsu consonant-extending symbol.
package nihongo // import "robpike.io/nihongo"

import (
	"bufio"
	"bytes"
	"io"
	"unicode/utf8"
)

const eof = -1

// The getters return a function that gets the next rune from the various input sources.

func stringGetter(s string) func() rune {
	return func() rune {
		if len(s) == 0 {
			return eof
		}
		r, w := utf8.DecodeRuneInString(s)
		s = s[w:]
		return r
	}
}

func bytesGetter(b []byte) func() rune {
	return func() rune {
		if len(b) == 0 {
			return eof
		}
		r, w := utf8.DecodeRune(b)
		b = b[w:]
		return r
	}
}

func readerGetter(r io.Reader) func() rune {
	rr, ok := r.(io.RuneReader)
	if !ok {
		rr = bufio.NewReader(r)
	}
	return func() rune {
		c, _, err := rr.ReadRune()
		if err != nil {
			return eof
		}
		return c
	}
}

// The putters return a function that delivers to the various output sinks.

func bufPutter(buf *bytes.Buffer) func(byte) {
	return func(b byte) {
		buf.WriteByte(b)
	}
}

func chanPutter(ch chan byte) func(byte) {
	return func(b byte) {
		ch <- b
	}
}

// translator handles the io.
type translator struct {
	get   func() rune
	put   func(byte)
	ch    chan byte
	peekc rune
	// These are used only when next3 doing the input (hiragana, katakana).
	save    []byte
	runeBuf [utf8.UTFMax]byte
}

func newTranslator(get func() rune, put func(byte), ch chan byte) *translator {
	return &translator{
		get:   get,
		put:   put,
		ch:    ch,
		peekc: eof,
		save:  make([]byte, 0, 2*utf8.UTFMax),
	}
}

func (t *translator) next() rune {
	if t.peekc >= 0 {
		f := t.peekc
		t.peekc = eof
		return f
	}
	return t.get()
}

func (t *translator) peek() rune {
	r := t.next()
	if r != eof {
		t.pushback(r)
	}
	return r
}

func (t *translator) pushback(r rune) {
	t.peekc = r
}

func (t *translator) putRune(r rune) {
	var buf [utf8.UTFMax]byte
	n := utf8.EncodeRune(buf[:], r)
	for i := 0; i < n; i++ {
		t.put(buf[i])
	}
}

func (t *translator) putString(s string) {
	for i := 0; i < len(s); i++ {
		t.put(s[i])
	}
}

func (t *translator) next3() string {
	for len(t.save) < 3 {
		r := t.get()
		if r == eof {
			return string(t.save)
		}
		n := utf8.EncodeRune(t.runeBuf[:], r)
		t.save = append(t.save, t.runeBuf[:n]...)
	}
	return string(t.save[:3])
}

func (t *translator) advance(n int) {
	copy(t.save, t.save[n:])
	t.save = t.save[:len(t.save)-n]
}

func (t *translator) Read(p []byte) (int, error) {
	n := 0
	for n < len(p) {
		c, ok := <-t.ch
		if !ok {
			return n, io.EOF
		}
		p[n] = c
		n++
	}
	return n, nil
}
