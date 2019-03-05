package kindleland

import (
	"strings"
	"testing"
)

type nextWord struct {
	word  string
	space string
	ok    bool
}

func TestNextWord(t *testing.T) {
	text := `This is  a sentence.`

	nextWords := []nextWord{
		{"This", " ", true},
		{"is", "  ", true},
		{"a", " ", true},
		{"sentence.", "", true},
		{"", "", false},
	}

	tb := NewTextBuffer(text)

	for i, n := range nextWords {
		word, space, ok := tb.NextWord()
		if word != n.word {
			t.Errorf("wrong word on call %d; got '%s', expected '%s'", i, word, n.word)
		}
		if space != n.space {
			t.Errorf("wrong space on call %d; got '%s', expected '%s'", i, space, n.space)
		}
		if ok != n.ok {
			t.Errorf("wrong ok on call %d; got '%v', expected '%v'", i, ok, n.ok)
		}
	}

	if !tb.eof {
		t.Error("textbuffer was not eof")
	}
}
func TestNextParagraph(t *testing.T) {
	text := `This \nis  \na\n paragraph.`

	paragraphs := strings.Split(text, "\n")

	tb := NewTextBuffer(text)

	for i, p := range paragraphs {
		if tb.paragraphIndex != i {
			t.Errorf("wrong paragraph index on call %d; got '%d', expected '%d'", i, tb.paragraphIndex, i)
		}

		var paragraph string
		for {
			word, space, ok := tb.NextWord()
			if !ok {
				break
			}
			paragraph = paragraph + word + space
		}
		if paragraph != p {
			t.Errorf("wrong paragraph on call %d; got '%s', expected '%s'", i, paragraph, p)
		}
		tb.NextParagraph()
	}

	if !tb.eof {
		t.Error("textbuffer was not eof")
	}
}
