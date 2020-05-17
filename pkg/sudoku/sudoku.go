// Package sudoku allows for solving of sudoku
// games through a back-tracking algorithm.
//
// Currently only supports a 9x9 grid, consisting of
// 9 boxes, each containing 9 cells.
package sudoku

import (
	"fmt"
	"strconv"
	"strings"
)

// Cell holds information on each cell;
// its value, whether is visible and/or locked.
type Cell struct {
	value   int
	visible bool
	locked  bool
}

// Move contains information on a move made
// during the solving algorithm. It stores the cell
// index and a value (0 for unsetting/backtracking).
type Move struct {
	cell  int
	value int
}

// Board is a sudoku board
type Board [][]Cell

// Sprint string-prints the current board prettily.
// Indents the board with pad
func (b Board) Sprint(pad int) string {

	paddingString := ""
	for i := 0; i < pad; i++ {
		paddingString = paddingString + " "
	}

	var board strings.Builder
	// header
	board.WriteString(paddingString + "|-----|-----|-----|\n")
	for x := 0; x < len(b); x++ {
		board.WriteString(paddingString + "|")
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
			board.WriteString(paddingString + "|-----|-----|-----|\n")
		}
	}
	return strings.TrimRight(board.String(), "\n")
}

// New ceates an empty Board
func New() Board {
	size := 9

	board := make(Board, size)
	for i := 0; i < size; i++ {
		board[i] = make([]Cell, size)
	}

	return board
}

// Parse parses a slice of strings and fills the board with it.
// Zeroes/empty values are used for unset cells.
func (b Board) Parse(cols []string) error {
	if len(cols) > 9 {
		return fmt.Errorf("too many cols; want 9 got %d", len(cols))
	}

	for x, col := range cols {
		if len(col) > 9 {
			return fmt.Errorf("row too long; want 9, length %d", len(col))
		}
		for y, valueRune := range col {
			value := int(valueRune - '0')
			if value > 0 {
				ok, err := b.Set(x, y, value, true)
				if !ok || err != nil {
					err = fmt.Errorf("could not set %d at %2d,%2d - err: %v", value, x, y, err)
					return err
				}
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

	if !b.isValid(x, y, value) {
		return false, nil
	}

	b[y][x] = Cell{value, true, lock}
	return true, nil
}

// Unset unsets the cell at a given position
func (b Board) Unset(x, y int) (ok bool, err error) {
	if x > len(b) || y > len(b) || x < 0 || y < 0 {
		return false, fmt.Errorf("invalid coordinate: %d,%d", x, y)
	}

	if b[y][x].locked {
		return false, fmt.Errorf("cell at %2d,%2d is locked", x, y)
	}

	b[y][x] = Cell{}
	return true, nil
}

// isValid checks whether it's valid to set value at position x,y.
func (b Board) isValid(x, y, value int) bool {
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

	// check box:
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
func (b Board) Solve() (origin Board, movesToSolve []Move) {
	// Make a copy of the board:
	for _, row := range b {
		row = append([]Cell{}, row...)
		origin = append(origin, row)
	}

	noOfSquares := len(b) * len(b)
	moves := []Move{}

	for sq := 0; sq < noOfSquares; sq++ {
		for value := 1; value < 10; value++ {
			// already has value
			if s, _ := b.getSquare(sq); s.value > 0 {
				break
			}
			// we could assign a value and break the inner loop
			if ok, _ := b.setSquare(sq, value); ok {
				moves = append(moves, Move{sq, value})
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
					moves = append(moves, Move{sq, 0})
				}
			}
		}
	}
	return origin, moves
}

/*
** HELPER FUNCTIONS
 */
// getCoord is a helper function that gets coordinates from a given cell.
func (b Board) getCoord(n int) (x, y int) {
	col := n / len(b)
	return col, n - (col)*len(b)
}

// setSquare is a helper function that sets a value by index in stead of
// coordinates.
// As this is a helper function and is only called doing solving, it does not lock
// cells.
func (b Board) setSquare(n, value int) (ok bool, err error) {
	x, y := b.getCoord(n)
	return b.Set(x, y, value, false)
}

// unsetSquare is a helper function that unsets a cell by index in stead of
// coordinates.
func (b Board) unsetSquare(n int) (ok bool, err error) {
	x, y := b.getCoord(n)
	return b.Unset(x, y)
}

// getSquare gets a Cell from its index.
func (b Board) getSquare(n int) (sq Cell, err error) {
	if n >= len(b)*len(b) {
		return Cell{}, fmt.Errorf("invalid Cell %d; larger than grid size", n)
	}
	x, y := b.getCoord(n)
	return b[y][x], nil
}
