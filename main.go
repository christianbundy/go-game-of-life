package main

import (
	"bytes"
	"flag"
	"fmt"
	"golang.org/x/term"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	term_width, term_height, err := term.GetSize(0)
	if err != nil {
		log.Fatal(err)
	}

	flag.Parse()

	width := term_width * 2
	height := term_height * 2
	size := width * height

	board := make([]bool, size)
	for index := range board {
		// TODO: Seed RNG.
		if rand.Intn(2) == 0 {
			board[index] = true
		}
	}

	fmt.Println("\033[?47h")

	start := time.Now()

	var target_fps = 30
	var target_pause = 1.0 / float64(target_fps)

	for {
		select {
		case <-sigs:
			fmt.Println("\033[?47l")
			os.Exit(0)
		default:

			elapsed := time.Now().Sub(start)

			if elapsed.Seconds() > target_pause {
				start = time.Now()
				display(board, width, height)
				board = next(board, size, width)
			}
		}
	}
}

func next(board []bool, size int, width int) []bool {
	var new_board = make([]bool, size)

	for index := range board {
		neighbor_count := 0

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

func display(board []bool, width int, height int) {
	var output bytes.Buffer
	size := width * height

	for index := range board {
		if ((index % (width * 2)) >= width) && index%2 == 1 {
			// TODO: Fix bad variable names
			up := board[(index+size-width)%size]   // up
			me := board[index]                     // current
			le := board[(index+size-1)%size]       // left
			lu := board[(index+size-width-1)%size] // left up

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
