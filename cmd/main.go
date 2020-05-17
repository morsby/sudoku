package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/morsby/sudoku/pkg/sudoku"
)

func main() {
	example := flag.Bool("example", false, "to demo an example")
	speed := flag.Int("speed", 200, "ns between moves")
	flag.Parse()
	input := []string{}
	if *example == true {
		input = []string{
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
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for len(input) < 9 {
			fmt.Printf("Please enter row %d:\t", len(input)+1)
			scanner.Scan()
			input = append(input, scanner.Text())
		}
	}
	board := sudoku.New()
	err := board.Parse(input)
	if err != nil {
		fmt.Println(err)
		return
	}

	origin, moves := board.Solve()
	sudoku.Animate(os.Stdout, *speed, origin, moves)

}
