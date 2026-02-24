// fpart.go - draw ascii art from a byte slice
//
// Borrowed from: https://github.com/moul/drunken-bishop
// Licensed under same terms as the original above (Apache 2.0)
package fp

import (
	"fmt"
	"strings"
)

const (
	CellWidth  = 16
	CellHeight = 8
	StartY     = CellHeight / 2
	StartX     = CellWidth / 2
	StartCode  = 1000
	EndCode    = 2000

	StartChar   = 'S'
	EndChar     = 'E'
	Alphabet    = " .+oVXLZAHPDRB"
	UnknownChar = '!'

	NW = "00"
	NE = "01"
	SW = "10"
	SE = "11"
)

type FPArt struct {
	opt
}

type opt struct {
	alphabet string
	h, w     int
}

type Option func(o *opt)

func WithAlphabet(s string) Option {
	return func(o *opt) {
		o.alphabet = s
	}
}

func WithHeight(x int) Option {
	return func(o *opt) {
		o.h = x
	}
}

func WithWidth(x int) Option {
	return func(o *opt) {
		o.w = x
	}
}

func ToString(b []byte, opts ...Option) string {
	o := &opt{
		alphabet: Alphabet,
		h:        CellHeight,
		w:        CellWidth,
	}

	for _, fp := range opts {
		fp(o)
	}

	var w strings.Builder

	room := o.bytes2room(b)
	buf := make([]byte, o.w)
	alen := len(o.alphabet)

	fmt.Fprintf(&w, "+%s+\n", strings.Repeat("-", o.w))
	for _, row := range room {
		line := buf[:0]
		for _, col := range row {
			var char byte
			switch {
			case col == StartCode:
				char = StartChar
			case col == EndCode:
				char = EndChar
			default:
				char = o.alphabet[col%alen]
			}
			line = append(line, char)
		}
		fmt.Fprintf(&w, "|%s|\n", string(line))
	}
	fmt.Fprintf(&w, "+%s+\n", strings.Repeat("-", o.w))
	return w.String()
}

func (o *opt) makeroom() [][]int {
	room := make([][]int, o.h)
	for i := range o.h {
		room[i] = make([]int, o.w)
	}
	return room
}

func (o *opt) bytes2room(b []byte) [][]int {
	room := o.makeroom()
	sx := o.w / 2
	sy := o.h / 2
	pos := pos{sy, sx}
	for _, bitpair := range bytesToBitpairs(b) {
		switch bitpair {
		case NW:
			pos.Y--
			pos.X--
		case NE:
			pos.Y--
			pos.X++
		case SW:
			pos.Y++
			pos.X--
		case SE:
			pos.Y++
			pos.X++
		}
		switch {
		case pos.Y < 0:
			pos.Y = 0
		case pos.Y >= o.h:
			pos.Y = o.h - 1
		}
		switch {
		case pos.X < 0:
			pos.X = 0
		case pos.X >= o.w:
			pos.X = o.w - 1
		}
		room[pos.Y][pos.X]++
	}
	room[sy][sx] = StartCode
	room[pos.Y][pos.X] = EndCode
	return room
}

type pos struct{ Y, X int }

func bytesToBitpairs(input []byte) []string {
	bitpairs := []string{}
	for _, byte := range input {
		bin := fmt.Sprintf("%08b", byte)
		bitpairs = append(
			bitpairs,
			bin[6:8],
			bin[4:6],
			bin[2:4],
			bin[0:2],
		)
	}
	return bitpairs
}
