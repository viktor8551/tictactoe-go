package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"golang.org/x/sys/windows"
)

const (
	rows = 3
	cols = 3
)

const (
	colorReset = "\033[0m"

	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
)

var board [rows][cols]int
var clear_terminal func()

func init_linux() {
	if runtime.GOOS != "linux" {
		return
	}

	// Set clear terminal function for linux
	clear_terminal = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func init_windows() {
	if runtime.GOOS != "windows" {
		return
	}

	// Set clear terminal function for windows
	clear_terminal = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}

	// Set console mode to be able to use ANSI color codes in windows cmd
	stdout := windows.Handle(os.Stdout.Fd())
	var originalMode uint32

	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
}

func colorize(text string, color string) string {
	return string(color) + text + string(colorReset)
}

func getBoardChar(number int) string {
	switch number {
	case 1:
		return colorize("X", colorBlue)
	case 2:
		return colorize("O", colorPurple)
	default:
		return "-"
	}
}

func printBoard() {
	for i := 0; i < rows; i++ {
		for k := 0; k < cols; k++ {
			fmt.Print("+---")
		}
		fmt.Println("+")

		fmt.Print("| ")
		for j := 0; j < cols; j++ {
			fmt.Print(getBoardChar(board[i][j]), " | ")
		}
		fmt.Println()
	}

	for k := 0; k < cols; k++ {
		fmt.Print("+---")
	}
	fmt.Println("+")
}

func checkWin(player int) bool {
	// Check rows
	for i := 0; i < rows; i++ {
		win := true

		for j := 0; j < cols; j++ {
			if board[i][j] != player {
				win = false
				break
			}
		}

		if win {
			return true
		}
	}

	// Check cols
	for i := 0; i < cols; i++ {
		win := true

		for j := 0; j < rows; j++ {
			if board[j][i] != player {
				win = false
				break
			}
		}

		if win {
			return true
		}
	}

	// Check diagonals
	win := true
	for i := 0; i < rows; i++ {
		if board[i][i] != player {
			win = false
			break
		}
	}
	if win {
		return true
	}
	win = true
	for i := 0; i < rows; i++ {
		if board[i][rows-i-1] != player {
			win = false
			break
		}
	}

	return win
}

func start_game() {
	player, selectedSlotsCount := 1, 0
	clear_terminal()

	// Game loop
	for {
		row, col := 0, 0
		printBoard()

		// Wait for user inputs
		fmt.Println(colorize(fmt.Sprintf("Player %d's turn", player), colorBlue))
		fmt.Printf("Enter row (1-%d): ", rows)
		fmt.Scanf("%d\n", &row)
		row = row - 1
		fmt.Printf("Enter column (1-%d): ", cols)
		fmt.Scanf("%d\n", &col)
		col = col - 1
		fmt.Println()

		// Input error handling
		if row >= rows || row < 0 || col >= cols || col < 0 {
			clear_terminal()
			fmt.Println(colorize("[ERROR]: Input value is invalid!", colorRed))
			continue
		}

		if board[row][col] != 0 {
			clear_terminal()
			fmt.Println(colorize("[ERROR]: Board slot is already taken!", colorRed))
			continue
		}

		// Set board slot value to player value
		board[row][col] = player
		selectedSlotsCount++

		// Check if player wins
		if selectedSlotsCount >= rows+cols-1 {
			winner := checkWin(player)
			if winner {
				clear_terminal()
				printBoard()
				fmt.Println(colorize(fmt.Sprintf("WINNER: Player %d", player), colorGreen))
				break
			}
		}

		// If all slots are selected
		if selectedSlotsCount >= rows*cols {
			clear_terminal()
			printBoard()
			fmt.Println(colorize("TIE!", colorBlue))
			break
		}

		// Switch player
		if player == 1 {
			player = 2
		} else {
			player = 1
		}

		clear_terminal()
	}

	// Loop for replay
	for {
		var restart string
		fmt.Println("Play again? (yes/no)")
		fmt.Scanf("%s\n", &restart)

		if strings.ToLower(restart) == "yes" {
			board = [rows][cols]int{}
			start_game()
			break
		} else if restart == "no" {
			os.Exit(3)
		} else {
			fmt.Println(colorize("[ERROR]: Given an invalid input (yes or no)", colorRed))
		}
	}
}

func main() {
	init_linux()
	init_windows()

	start_game()
}
