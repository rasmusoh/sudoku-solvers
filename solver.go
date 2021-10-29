package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"math/bits"
	"os"
	"runtime"
	"strings"
	"time"
)

const max = 1000000
const n = 9
const c = 3

type Coord [2]int

type Game struct {
	board     [][]uint16
	toProcess []Coord
}

func main() {
	start := time.Now()
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
	scanner.Scan()

	ch := make(chan string, 100)
	for i := 0; i < runtime.NumCPU(); i++ {
		go solver(ch)
	}
	i := 1
	defer func() {
		t := time.Now()
		elapsed := t.Sub(start)
		fmt.Printf("solved %d puzzles in %v", i, elapsed)
	}()
	for scanner.Scan() {
		if i > max {
			return
		}
		ch <- scanner.Text()
		i++
		if i%1000 == 0 {
			fmt.Printf("solving at %d \n", i)
		}
	}
	close(ch)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func solver(ch chan string) {
	for line := range ch {
		s := strings.Split(line, ",")
		game, err := ParseBoard(s[0])
		if err != nil {
			panic(err)
		}
		err = game.Solve()
		if err != nil {
			log.Fatalf("error for line %s: %e", line, err)
		}
		if game.ToString() != s[1] {
			fmt.Printf("discrepancy for line %s\n", line)
		}
	}
}

func (game *Game) Solve() error {
	valid, err := game.PropagateAll()
	if err != nil {
		return err
	}
	if !valid {
		return errors.New("expected initial board to be valid")
	}
	if game.IsSolved() {
		return nil
	}
	hyps := game.GetNextHypotheticals()
	for len(hyps) > 0 {
		hyp := hyps[len(hyps)-1]
		hyps = hyps[:len(hyps)-1]
		valid, err = hyp.PropagateAll()
		if err != nil {
			return err
		}
		if !valid {
			continue
		} else if hyp.IsSolved() {
			game.board = hyp.board
			return nil
		} else {
			for _, next_hyp := range hyp.GetNextHypotheticals() {
				hyps = append(hyps, next_hyp)
			}
		}
	}
	return nil
}

func (game *Game) GetNextHypotheticals() []*Game {
	minAlts := 10
	minX, minY := -1, -1
	for y, line := range game.board {
		for x, set := range line {
			alts := bits.OnesCount16(set)
			if alts > 1 && alts < minAlts {
				minAlts = alts
				minX, minY = x, y
			}
		}
	}
	hyps := make([]*Game, 0, 2)
	set := game.board[minY][minX]
	for i := 0; i < n; i++ {
		if set%2 > 0 {
			hyp := game.Copy()
			hyp.board[minY][minX] = uint16(1 << i)
			hyp.toProcess = append(hyp.toProcess, Coord{minY, minX})
			hyps = append(hyps, hyp)
		}
		set = set >> 1
	}
	return hyps
}

func (game *Game) Copy() *Game {
	toProcess := make([]Coord, 0, n*n)
	sets := make([]uint16, n*n)
	for i := 0; i < n*n; i++ {
		sets[i] = game.board[i/n][i%n]
	}

	board := make([][]uint16, n)
	for i := 0; i < n; i++ {
		board[i] = sets[i*n : i*n+n]
	}

	return &Game{
		board:     board,
		toProcess: toProcess,
	}
}

func (game *Game) IsSolved() bool {
	solved := true
	for _, line := range game.board {
		for _, set := range line {
			solved = solved && (bits.OnesCount16(set) == 1)
		}
	}
	return solved
}

func (game *Game) PrintBoard() {
	for _, line := range game.board {
		for _, set := range line {

			if bits.OnesCount16(set) == 1 {
				fmt.Printf(" %d |", ToInt(set))
			} else {
				fmt.Print("   |")
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

func (game *Game) ToString() string {
	var b strings.Builder
	b.Grow(n * n)
	for _, line := range game.board {
		for _, set := range line {
			if bits.OnesCount16(set) == 1 {
				fmt.Fprintf(&b, "%d", ToInt(set))
			} else {
				fmt.Fprintf(&b, "%d", 0)
			}
		}
	}
	return b.String()
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
	count := bits.OnesCount16(game.board[y][x])
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

func (game *Game) Remove(y int, x int, value uint16) bool {
	hadValue := game.board[y][x]&value != 0
	game.board[y][x] = game.board[y][x] & (^value)
	if hadValue && bits.OnesCount16(game.board[y][x]) == 1 {
		game.toProcess = append(game.toProcess, Coord{y, x})
	}
	return game.board[y][x] != 0
}

func ParseBoard(line string) (*Game, error) {
	if len(line) != n*n {
		return nil, errors.New(fmt.Sprintf("unexpected length of line: %s", line))
	}

	toProcess := make([]Coord, 0, n*n)
	sets := make([]uint16, n*n)
	for i, char := range line {
		sets[i] = ToSet(char)
		if char != '0' {
			toProcess = append(toProcess, Coord{i / n, i % n})
		}
	}

	board := make([][]uint16, n)
	for i := 0; i < n; i++ {
		board[i] = sets[i*n : i*n+n]
	}

	return &Game{
		board:     board,
		toProcess: toProcess,
	}, nil
}

func ToSet(char rune) uint16 {
	if char == '0' {
		return uint16(0b111111111)
	}
	return uint16(1 << int(char-'0'-1))
}

func ToInt(set uint16) int {
	for i := 1; i <= n; i++ {
		if set%2 > 0 {
			return i
		}
		set = set >> 1
	}
	return 0
}
