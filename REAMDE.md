Package nihongo implements simple transliteration between romaji
and the two syllabic Japanese scripts, hiragana and katakana, encoded
as UTF-8-encoded Unicode. Romaji output may include injected spaces
to separate converted text from unconverted, and other markers.
Invalid sequences, such as small kanas with no preceding kana,
are passed unaltered. Hiragana and katakana may be inaccurate
due to false matches. Katakana may be further inaccurate because
of the inability to generate the tsu consonant-extending symbol.

For full documentation, see [the godoc page](http://godoc.org/robpike.io/nihongo).
