package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

const lines = 100
const n = 9
const c = 3

type Coord [2]int
type Set uint

type Game struct {
	board     [][]Set
	toProcess []Coord
}

func main() {
	file, err := os.Open("sudoku.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	scanner := bufio.NewScanner(file)
	i := 0
	scanner.Scan()
	for scanner.Scan() {
		if i > lines {
			return
		}
		s := strings.Split(scanner.Text(), ",")
		fmt.Println(s[0])
		game, err := ParseBoard(s[0])
		if err != nil {
			panic(err)
		}
		game.PrintBoard()
		valid, err := game.PropagateAll()
		if !valid {
			log.Fatal("expected board to be valid")
			return
		} else if err != nil {
			log.Fatalf("error: %e", err)
		}
		game.PrintBoard()
		i += 1
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func (game *Game) PrintBoard() {
	for _, line := range game.board {
		for _, set := range line {

			if Count(set) == 1 {
				fmt.Printf(" %d |", ToInt(set))
			} else {
				fmt.Print("   |")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (game *Game) PropagateAll() (bool, error) {
	for len(game.toProcess) > 0 {
		coord := game.toProcess[len(game.toProcess)-1]
		game.toProcess = game.toProcess[:len(game.toProcess)-1]
		valid, err := game.Propagate(coord[0], coord[1])
		if !valid || err != nil {
			return valid, err
		}
	}
	return true, nil
}

func (game *Game) Propagate(y int, x int) (bool, error) {
	count := Count(game.board[y][x])
	if count != 1 {
		return false, errors.New(fmt.Sprintf("%d, %d not determined, cannot propagate", y, x))
	}
	valid := true
	value := game.board[y][x]
	cellCornerY, cellCornerX := y/c, x/c
	for i := 0; i < n; i++ {
		if i != y {
			valid = valid && game.Remove(i, x, value)
		}
		if i != x {
			valid = valid && game.Remove(y, i, value)
		}
		cY, cX := c*cellCornerY+(i/3), c*cellCornerX+(i%3)
		if !(cY == y && cX == x) {
			valid = valid && game.Remove(cY, cX, value)
		}
	}
	return valid, nil
}

func (game *Game) Remove(y int, x int, value Set) bool {
	// fmt.Printf("Removing %d from %d,%d\n", ToInt(value), y, x)
	// fmt.Printf("Present count %d, value %d\n", Count(game.board[y][x]), ToInt(game.board[y][x]))
	hadValue := game.board[y][x]&value != 0
	// fmt.Printf("hadValue %t\n", hadValue)
	// fmt.Printf("bit fields:%b %b %b\n", game.board[y][x], (^value), value)
	game.board[y][x] = game.board[y][x] & (^value)
	// fmt.Printf("new count %d\n", Count(game.board[y][x]))
	if hadValue && Count(game.board[y][x]) == 1 {
		game.toProcess = append(game.toProcess, Coord{y, x})
	}
	return game.board[y][x] != 0
}

func ParseBoard(line string) (*Game, error) {
	if len(line) != n*n {
		return nil, errors.New(fmt.Sprintf("unexpected length of line: %s", line))
	}

	toProcess := make([]Coord, n*n)
	sets := make([]Set, n*n)
	for i, char := range line {
		sets[i] = ToSet(char)
		if char != '0' {
			toProcess = append(toProcess, Coord{i / n, i % n})
		}
	}

	board := make([][]Set, n)
	for i := 0; i < n; i++ {
		board[i] = sets[i*n : i*n+n]
	}

	return &Game{
		board:     board,
		toProcess: toProcess,
	}, nil
}

func ToSet(char rune) Set {
	if char == '0' {
		return Set(0b111111111)
	}
	val := Set(1 << int(char-'0'-1))
	// fmt.Printf("int %d, bin %b\n" ,ToInt(val),val)
	return val
}

func ToInt(set Set) int {
	if int(set) == 0 {
		return 0
	}
	i := 1
	for {
		if set%2 > 0 {
			return i
		}
		set = set >> 1
		i += 1
	}
}

//http://graphics.stanford.edu/~seander/bithacks.html#CountBitsSetParallel
func FastCount(set Set) uint {
	v := uint(set)
	v = v - ((v >> 1) & 0x55555555)                 // reuse input as temporary
	v = (v & 0x33333333) + ((v >> 2) & 0x33333333)  // temp
	c := ((v + (v>>4)&0xF0F0F0F) * 0x1010101) >> 24 // count
	return uint(c)
}

func Count(set Set) int {
	count := 0
	b := int(set)
	for b != 0 {
		count += b & 1
		b >>= 1
	}
	return count
}
