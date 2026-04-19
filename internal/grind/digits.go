package grind

import "strings"

// 5x7 pixel block glyphs for digits 0-9 and uppercase letters A-Z.
// 'X' = filled, ' ' = empty. Rendered via renderTextScaled.
var digitGlyphs = map[rune][]string{
	' ': {"     ", "     ", "     ", "     ", "     ", "     ", "     "},
	'A': {" XXX ", "X   X", "X   X", "XXXXX", "X   X", "X   X", "X   X"},
	'B': {"XXXX ", "X   X", "X   X", "XXXX ", "X   X", "X   X", "XXXX "},
	'C': {" XXXX", "X    ", "X    ", "X    ", "X    ", "X    ", " XXXX"},
	'D': {"XXXX ", "X   X", "X   X", "X   X", "X   X", "X   X", "XXXX "},
	'E': {"XXXXX", "X    ", "X    ", "XXX  ", "X    ", "X    ", "XXXXX"},
	'F': {"XXXXX", "X    ", "X    ", "XXX  ", "X    ", "X    ", "X    "},
	'G': {" XXXX", "X    ", "X    ", "X  XX", "X   X", "X   X", " XXXX"},
	'H': {"X   X", "X   X", "X   X", "XXXXX", "X   X", "X   X", "X   X"},
	'I': {"XXXXX", "  X  ", "  X  ", "  X  ", "  X  ", "  X  ", "XXXXX"},
	'J': {"XXXXX", "   X ", "   X ", "   X ", "   X ", "X  X ", " XX  "},
	'K': {"X   X", "X  X ", "X X  ", "XX   ", "X X  ", "X  X ", "X   X"},
	'L': {"X    ", "X    ", "X    ", "X    ", "X    ", "X    ", "XXXXX"},
	'M': {"X   X", "XX XX", "X X X", "X X X", "X   X", "X   X", "X   X"},
	'N': {"X   X", "XX  X", "X X X", "X X X", "X X X", "X  XX", "X   X"},
	'O': {" XXX ", "X   X", "X   X", "X   X", "X   X", "X   X", " XXX "},
	'P': {"XXXX ", "X   X", "X   X", "XXXX ", "X    ", "X    ", "X    "},
	'Q': {" XXX ", "X   X", "X   X", "X   X", "X X X", "X  X ", " XX X"},
	'R': {"XXXX ", "X   X", "X   X", "XXXX ", "X X  ", "X  X ", "X   X"},
	'S': {" XXXX", "X    ", "X    ", " XXX ", "    X", "    X", "XXXX "},
	'T': {"XXXXX", "  X  ", "  X  ", "  X  ", "  X  ", "  X  ", "  X  "},
	'U': {"X   X", "X   X", "X   X", "X   X", "X   X", "X   X", " XXX "},
	'V': {"X   X", "X   X", "X   X", "X   X", "X   X", " X X ", "  X  "},
	'W': {"X   X", "X   X", "X   X", "X X X", "X X X", "X X X", " X X "},
	'X': {"X   X", "X   X", " X X ", "  X  ", " X X ", "X   X", "X   X"},
	'Y': {"X   X", "X   X", " X X ", "  X  ", "  X  ", "  X  ", "  X  "},
	'Z': {"XXXXX", "    X", "   X ", "  X  ", " X   ", "X    ", "XXXXX"},
	'/': {"    X", "    X", "   X ", "  X  ", " X   ", "X    ", "X    "},
	'0': {"XXXXX", "X   X", "X   X", "X   X", "X   X", "X   X", "XXXXX"},
	'1': {"  X  ", " XX  ", "  X  ", "  X  ", "  X  ", "  X  ", " XXX "},
	'2': {"XXXXX", "    X", "    X", "XXXXX", "X    ", "X    ", "XXXXX"},
	'3': {"XXXXX", "    X", "    X", "XXXXX", "    X", "    X", "XXXXX"},
	'4': {"X   X", "X   X", "X   X", "XXXXX", "    X", "    X", "    X"},
	'5': {"XXXXX", "X    ", "X    ", "XXXXX", "    X", "    X", "XXXXX"},
	'6': {"XXXXX", "X    ", "X    ", "XXXXX", "X   X", "X   X", "XXXXX"},
	'7': {"XXXXX", "    X", "    X", "    X", "    X", "    X", "    X"},
	'8': {"XXXXX", "X   X", "X   X", "XXXXX", "X   X", "X   X", "XXXXX"},
	'9': {"XXXXX", "X   X", "X   X", "XXXXX", "    X", "    X", "XXXXX"},
}

// 3-wide x 7-tall colon glyph.
var colonGlyph = []string{
	"   ",
	" X ",
	"   ",
	"   ",
	"   ",
	" X ",
	"   ",
}

// renderTextScaled expands each pixel of a 5x7 glyph to `scaleW x scaleH`
// characters of `\u2588`, so at scale (2,1) each pixel becomes a roughly
// square block on terminals whose cell aspect is ~0.5x1.
func renderTextScaled(s string, scaleW, scaleH int) string {
	const glyphHeight = 7
	rows := make([]strings.Builder, glyphHeight*scaleH)

	for _, ch := range s {
		glyph, ok := digitGlyphs[ch]
		if !ok && ch == ':' {
			glyph = colonGlyph
		}
		if glyph == nil {
			continue
		}
		for py := 0; py < glyphHeight; py++ {
			for repeat := 0; repeat < scaleH; repeat++ {
				rowIdx := py*scaleH + repeat
				for _, pix := range glyph[py] {
					if pix == 'X' {
						rows[rowIdx].WriteString(strings.Repeat("\u2588", scaleW))
					} else {
						rows[rowIdx].WriteString(strings.Repeat(" ", scaleW))
					}
				}
			}
		}
	}

	lines := make([]string, len(rows))
	for i := range rows {
		lines[i] = rows[i].String()
	}
	return strings.Join(lines, "\n")
}
