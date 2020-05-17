package sudoku

import (
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	b := New()
	rows := len(b)
	cols := len(b[0])
	if rows != 9 || cols != 9 {
		t.Errorf("expected board to be 9x9, is %dx%d", cols, rows)
	}
}

func TestSprint(t *testing.T) {
	board := New()
	print := board.Sprint(0)
	expected := strings.TrimSpace(`
|-----|-----|-----|
|0 0 0|0 0 0|0 0 0|
|0 0 0|0 0 0|0 0 0|
|0 0 0|0 0 0|0 0 0|
|-----|-----|-----|
|0 0 0|0 0 0|0 0 0|
|0 0 0|0 0 0|0 0 0|
|0 0 0|0 0 0|0 0 0|
|-----|-----|-----|
|0 0 0|0 0 0|0 0 0|
|0 0 0|0 0 0|0 0 0|
|0 0 0|0 0 0|0 0 0|
|-----|-----|-----|
`)
	if print != expected {
		t.Errorf("Board looks wrong; got \n'%v'\n; wanted \n'%v'", print, expected)
	}
}

func TestSet(t *testing.T) {
	board := New()
	ok, err := board.Set(0, 2, 3, true)
	if !ok {
		t.Error("expected to be able to set value")
	}
	if err != nil {
		t.Errorf("Did not expect an err: %v", err)
	}
	got := board[2][0].value
	wanted := 3
	if got != wanted {
		t.Errorf("Value not set correctly; got %d; wanted %d", got, wanted)
	}

}

func TestIsValid(t *testing.T) {
	board := New()

	if board.isValid(0, 0, 1) != true {
		t.Error("Expected 1 to be valid to put; was not.")
	}

	board.Set(0, 0, 1, true)
	if board.isValid(0, 1, 1) != false {
		t.Error("Did not expect 1 to be valid to put; was.")
	}
	if board.isValid(0, 1, 2) != true {
		t.Error("Expected 2 to be valid to put; was not.")
	}

	if board.isValid(1, 1, 1) != false {
		t.Error("Expected 1 to be unvalid; was not.")
	}
}

func TestGetCoord(t *testing.T) {
	type expect struct {
		input  int
		output []int
	}

	ts := []expect{
		expect{0, []int{0, 0}},
		expect{1, []int{0, 1}},
		expect{8, []int{0, 8}},
		expect{9, []int{1, 0}},
		expect{10, []int{1, 1}},
		expect{80, []int{8, 8}},
	}
	board := New()
	for _, v := range ts {
		x, y := board.getCoord(v.input)
		got := []int{x, y}
		if !reflect.DeepEqual(got, v.output) {
			t.Errorf("Input: %v, got: %v; wanted: %v", v.input, got, v.output)
		}
	}
}

func TestParse(t *testing.T) {
	board := New()
	seed := []string{
		"000000000",
		"000001098",
		"005700030",
		"006050801",
		"000060000",
		"070900000",
		"068000407",
		"340000000",
		"090502600",
	}
	board.Parse(seed)
	expected := strings.TrimSpace(`
|-----|-----|-----|
|0 0 0|0 0 0|0 0 0|
|0 0 0|0 0 1|0 9 8|
|0 0 5|7 0 0|0 3 0|
|-----|-----|-----|
|0 0 6|0 5 0|8 0 1|
|0 0 0|0 6 0|0 0 0|
|0 7 0|9 0 0|0 0 0|
|-----|-----|-----|
|0 6 8|0 0 0|4 0 7|
|3 4 0|0 0 0|0 0 0|
|0 9 0|5 0 2|6 0 0|
|-----|-----|-----|
`)

	b := board.Sprint(0)
	if expected != b {
		t.Error("Board is not right.")
	}
}

func TestSolve(t *testing.T) {
	board := New()
	seed := []string{
		"000080070",
		"058030100",
		"000000000",
		"026000090",
		"400000006",
		"700029300",
		"007000900",
		"100203000",
		"060000054",
	}
	board.Parse(seed)
	beforeSolving := board.Sprint(0)
	origin, _ := board.Solve()
	solved := board.Sprint(0)
	expected := strings.TrimSpace(`
|-----|-----|-----|
|3 1 2|5 8 6|4 7 9|
|9 5 8|4 3 7|1 6 2|
|6 7 4|9 1 2|5 3 8|
|-----|-----|-----|
|5 2 6|3 7 4|8 9 1|
|4 3 9|1 5 8|7 2 6|
|7 8 1|6 2 9|3 4 5|
|-----|-----|-----|
|2 4 7|8 6 5|9 1 3|
|1 9 5|2 4 3|6 8 7|
|8 6 3|7 9 1|2 5 4|
|-----|-----|-----|`)

	if beforeSolving != origin.Sprint(0) {
		t.Error("Expected Solve origin to be unmodified by solving.")
	}

	if expected != solved {
		t.Errorf("Solution wrong; got:\n%v\n \t\t\t\twanted: \n%v", solved, expected)
	}

}
