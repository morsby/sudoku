package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/morsby/sudoku/pkg/sudoku"
)

func main() {
	example := flag.String("example", "", "to demo an example [easy or expert]")
	speed := flag.Int("speed", 0, "ms between moves")
	flag.Parse()
	input := []string{}
	if *example == "easy" {
		input = []string{
			"020500000",
			"649831752",
			"500600000",
			"003084690",
			"018390047",
			"006100208",
			"080000924",
			"074008105",
			"060000000",
		}
	} else if *example == "expert" {
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
