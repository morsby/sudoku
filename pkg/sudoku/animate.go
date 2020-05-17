package sudoku

import (
	"io"
	"time"

	"github.com/inancgumus/screen"
	"github.com/schollz/progressbar/v3"
)

// Animate shows each move made to solve a sudoku
func Animate(w io.Writer, speed int, b Board, moves []Move) {
	screen.Clear()
	width, height := screen.Size()

	// output = 13 lines
	// bar = 1 line
	padHeight := (height - 14) / 2
	padWidth := (width - 14) / 2

	var paddingString string
	for i := 0; i < padHeight; i++ {
		paddingString = paddingString + "\n"
	}

	w.Write([]byte(paddingString + b.Sprint(padWidth) + paddingString + "\n"))

	bar := progressbar.Default(int64(len(moves)))

	for _, move := range moves {
		time.Sleep(time.Duration(speed) * time.Millisecond)
		screen.MoveTopLeft()
		if move.value == 0 {
			b.unsetSquare(move.cell)
		} else {
			b.setSquare(move.cell, move.value)
		}
		w.Write([]byte(paddingString + b.Sprint(padWidth) + paddingString))
		bar.Add(1)
	}
}
