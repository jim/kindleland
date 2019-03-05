package kindleland

import (
	"strings"
)

type TextBuffer struct {
	Paragraphs     []string
	paragraphIndex int
	byteIndex      int
	currentWord    string
	currentSpace   string
	eof            bool
}

func NewTextBuffer(text string) *TextBuffer {
	var paragraphs []string
	for _, p := range strings.Split(text, "\n") {
		if len(p) > 0 {
			paragraphs = append(paragraphs, p)
		}
	}
	return &TextBuffer{
		Paragraphs: paragraphs,
	}
}

func (tb *TextBuffer) NextWord() (string, string, bool) {
	if tb.eof {
		return "", "", false
	}
	tb.advance()
	return tb.currentWord, tb.currentSpace, len(tb.currentWord) > 0
}

func (tb *TextBuffer) NextParagraph() bool {
	if tb.eof {
		return false
	}

	if tb.paragraphIndex >= len(tb.Paragraphs) {
		tb.eof = true
		return false
	}

	tb.byteIndex = 0
	tb.paragraphIndex++
	return true
}

func (tb *TextBuffer) advance() {
	var word, space []rune
	var bi int

	// fmt.Printf("%d %d %d '%s' '%s'\n", tb.paragraphIndex, tb.byteIndex, len(tb.Paragraphs[tb.paragraphIndex]), tb.currentWord, tb.currentSpace)

	tb.currentWord = ""
	tb.currentSpace = ""

	currentParagraph := tb.Paragraphs[tb.paragraphIndex]
	if tb.byteIndex >= len(currentParagraph) {
		tb.byteIndex = 0
		tb.paragraphIndex++
	}

	if tb.paragraphIndex >= len(tb.Paragraphs) {
		tb.eof = true
		return
	}

	// fmt.Println(tb.byteIndex, tb.Paragraphs[tb.paragraphIndex][tb.byteIndex:])

	for b, r := range currentParagraph[tb.byteIndex:] {
		// fmt.Println(b, string(r))
		s := string(r)
		if s == " " {
			space = append(space, r)
		} else {
			if len(space) > 0 {
				break
			}
			word = append(word, r)
		}
		bi = b
	}

	tb.byteIndex += bi + 1
	tb.currentWord = string(word)
	tb.currentSpace = string(space)
	// fmt.Printf("%d, '%s', '%s'\n", bi, string(word), string(space))

	return
}
