// Package sudoku allows for solving of sudoku
// games through a back-tracking algorithm.
//
// Currently only supports a 9x9 grid.
package sudoku

import (
	"fmt"
	"strconv"
	"strings"
)

type square struct {
	value   int
	visible bool
	locked  bool
}

// Board is a sudoku board
type Board [][]square

// Sprint string-prints the current board prettily.
func (b Board) Sprint() string {
	var board strings.Builder
	// header
	board.WriteString("|-----|-----|-----|\n")
	for x := 0; x < len(b); x++ {
		board.WriteString("|")
		for y := 0; y < len(b[0]); y++ {
			board.WriteString(strconv.Itoa(b[y][x].value))
			// spacing between cells
			if y%3 == 2 {
				board.WriteRune('|')
			} else {
				board.WriteRune(' ')
			}
		}
		board.WriteString("\n")
		if x%3 == 2 {
			board.WriteString("|-----|-----|-----|\n")
		}
	}
	return strings.TrimSpace(board.String())
}

// New ceates an empty Board
func New() Board {
	size := 9

	board := make(Board, size)
	for i := 0; i < size; i++ {
		board[i] = make([]square, size)
	}

	return board
}

// Parse parses a slice of strings and fills the board with it.
// Zeroes are used for unset cells.
func (b Board) Parse(cols []string) error {
	if len(cols) != 9 {
		return fmt.Errorf("invalid number of cols; want 9 got %d", len(cols))
	}

	for x, col := range cols {
		if len(col) != 9 {
			return fmt.Errorf("invalid number of rows; want 9 length %d", len(col))
		}
		for y, valueRune := range col {
			value := int(valueRune - '0')
			if value > 0 {
				b.Set(x, y, value, true)
			}
		}
	}

	return nil
}

// Set sets a value through a 0-based index (i.e. first position is [0,0])
func (b Board) Set(x, y, value int, lock bool) (ok bool, err error) {
	if value > 9 {
		return false, fmt.Errorf("invalid value: %d > 9", value)
	}

	if x > len(b) || y > len(b) || x < 0 || y < 0 {
		return false, fmt.Errorf("invalid coordinate: %d,%d", x, y)
	}

	if !b.isSafe(x, y, value) {
		return false, nil
	}

	b[y][x] = square{value, true, lock}
	return true, nil
}

// Unset unsets the value at a given position
func (b Board) Unset(x, y int) (ok bool, err error) {
	if x > len(b) || y > len(b) || x < 0 || y < 0 {
		return false, fmt.Errorf("invalid coordinate: %d,%d", x, y)
	}

	if b[y][x].locked {
		return false, fmt.Errorf("cell at %2d,%2d is locked", x, y)
	}

	b[y][x] = square{}
	return true, nil
}

// isSafe checks whether it's valid to set value at position x,y.
func (b Board) isSafe(x, y, value int) bool {
	for i := 0; i < len(b); i++ {
		// check col;
		if b[y][i].value == value {
			return false
		}

		// check row:
		if b[i][x].value == value {
			return false
		}
	}

	// check small square:
	rowBox := y / 3 * 3
	colBox := x / 3 * 3
	for i := rowBox; i < rowBox+3; i++ {
		for j := colBox; j < colBox+3; j++ {
			if b[i][j].value == value {
				return false
			}
		}
	}

	return true
}

// Solve solves the sudoku
func (b Board) Solve() {
	type move struct {
		square int
		value  int
	}
	noOfSquares := len(b) * len(b)

	for sq := 0; sq < noOfSquares; sq++ {
		for value := 1; value < 10; value++ {
			// already has value
			if s, _ := b.getSquare(sq); s.value > 0 {
				break
			}
			// we could assign a value and break the inner loop
			if ok, _ := b.setSquare(sq, value); ok {
				break
			}

			// we could find no allowed value,
			// backtracking until we can try again.
			for value == 9 {
				sq--
				lastSquare, _ := b.getSquare(sq)
				if !lastSquare.locked {
					// increments by one when the loop restarts
					value = lastSquare.value
					b.unsetSquare(sq)
				}
			}
		}
	}
}

/*
** HELPER FUNCTIONS
 */
// getCoord is a helper function that gets coordinates from a given square.
func (b Board) getCoord(n int) (x, y int) {
	col := n / len(b)
	return col, n - (col)*len(b)
}

// setSquare is a helper function that sets a value by square number in stead of
// coordinates.
// As this is a helper function and is only called doing solving, it does not lock
// cells.
func (b Board) setSquare(n, value int) (ok bool, err error) {
	x, y := b.getCoord(n)
	return b.Set(x, y, value, false)
}

// unsetSquare is a helper function that unsets a square by square number in stead of
// coordinates.
func (b Board) unsetSquare(n int) (ok bool, err error) {
	x, y := b.getCoord(n)
	return b.Unset(x, y)
}

// getSquare gets a square from its number.
func (b Board) getSquare(n int) (sq square, err error) {
	if n >= len(b)*len(b) {
		return square{}, fmt.Errorf("invalid square %d; larger than grid size", n)
	}
	x, y := b.getCoord(n)
	return b[y][x], nil
}
