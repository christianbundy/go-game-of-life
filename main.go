package main

import (
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"time"
)

func main() {
	term_width := flag.Int("width", 80, "number of columns")
	term_height := flag.Int("height", 20, "number of rows")

	flag.Parse()

	var width = *term_width * 2
	var height = *term_height * 2
	var size = width * height

	var board = make([]bool, size)
	for index := range board {
		// TODO: Seed RNG.
		if rand.Intn(2) == 0 {
			board[index] = true
		}
	}

	create_display(height)

	start := time.Now()

	var target_fps = 30
	var target_pause = 1.0 / float64(target_fps)

	for {
		elapsed := time.Now().Sub(start)

		if elapsed.Seconds() > target_pause {
			start = time.Now()
			display(board, width, height)
			board = next(board, size, width)
		}
	}
}

func next(board []bool, size int, width int) []bool {
	var new_board = make([]bool, size)

	for index := range board {
		var neighbor_count = 0

		// ???: Why doesn't this have a var?
		// TODO: Un-Rustify this code. It's okay to have negative integers because
		// the type is int, not usize.
		neighbors := [8]int{
			(index + size - width) % size,     // up
			(index + size - width + 1) % size, // up right
			(index + size + 1) % size,         // right
			(index + size + width + 1) % size, // down right
			(index + size + width) % size,     // down down
			(index + size + width - 1) % size, // down left
			(index + size - 1) % size,         // left
			(index + size - width - 1) % size, // up left
		}

		for _, neighbor_index := range neighbors {
			if board[neighbor_index] {
				neighbor_count += 1
			}
		}

		if board[index] && neighbor_count == 2 || neighbor_count == 3 {
			new_board[index] = true
		} else if board[index] == false && neighbor_count == 3 {
			new_board[index] = true
		} else {
			new_board[index] = false
		}
	}

	return new_board
}

func create_display(height int) {
	var count = 0
	for count != height/2 {
		fmt.Print("\n")
		count += 1
	}
}

func display(board []bool, width int, height int) {
	var output bytes.Buffer
	var size = width * height

	for index := range board {
		if ((index % (width * 2)) >= width) && index%2 == 1 {
			// TODO: Fix bad variable names
			var up = board[(index+size-width)%size]   // up
			var me = board[index]                     // current
			var le = board[(index+size-1)%size]       // left
			var lu = board[(index+size-width-1)%size] // left up

			// I really wish we had pattern matching here...
			if up == false && me == false && le == false && lu == false {
				output.WriteString(" ")
			}
			if up == false && me == false && le == false && lu == true {
				output.WriteString("▘")
			}
			if up == false && me == false && le == true && lu == false {
				output.WriteString("▖")
			}
			if up == false && me == false && le == true && lu == true {
				output.WriteString("▌")
			}
			if up == false && me == true && le == false && lu == false {
				output.WriteString("▗")
			}
			if up == false && me == true && le == false && lu == true {
				output.WriteString("▚")
			}
			if up == false && me == true && le == true && lu == false {
				output.WriteString("▄")
			}
			if up == false && me == true && le == true && lu == true {
				output.WriteString("▙")
			}
			if up == true && me == false && le == false && lu == false {
				output.WriteString("▝")
			}
			if up == true && me == false && le == false && lu == true {
				output.WriteString("▀")
			}
			if up == true && me == false && le == true && lu == false {
				output.WriteString("▞")
			}
			if up == true && me == false && le == true && lu == true {
				output.WriteString("▛")
			}
			if up == true && me == true && le == false && lu == false {
				output.WriteString("▐")
			}
			if up == true && me == true && le == false && lu == true {
				output.WriteString("▜")
			}
			if up == true && me == true && le == true && lu == false {
				output.WriteString("▟")
			}
			if up == true && me == true && le == true && lu == true {
				output.WriteString("█")
			}

			if index%width == width-1 {
				output.WriteString("\n")
			}
		}
	}

	fmt.Printf("\u001b[%dA%s", height/2, output.String())
}
