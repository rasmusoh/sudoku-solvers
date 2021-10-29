package main

import (
	"testing"
)

const testboard = "070000043040009610800634900094052000358460020000800530080070091902100005007040802"
const solution = "679518243543729618821634957794352186358461729216897534485276391962183475137945862"
const testboard_difficult = "068700900004000071030809050300080100040005007007304092602001005000020600059030028"
const solution_difficult = "568712943924653871731849256395287164246195387817364592682971435473528619159436728"


func TestSolve(t *testing.T) {
	game, err := ParseBoard(testboard)
	if err != nil {
		t.Errorf("error when parsing, %e", err)
	}
	err = game.Solve()
	if err != nil {
		t.Errorf("error when solving, %e", err)
	}
	if !game.IsSolved() {
		t.Error("game wasn't solved")
	}
	if game.ToString() != solution {
		t.Errorf("game wasn't solved correctly, expected %s got %s", game.ToString(), solution)
	}
}

func TestSolveDifficult(t *testing.T) {
	game, err := ParseBoard(testboard_difficult)
	if err != nil {
		t.Errorf("error when parsing, %e", err)
	}
	err = game.Solve()
	if err != nil {
		t.Errorf("error when solving, %e", err)
	}
	if !game.IsSolved() {
		t.Error("game wasn't solved")
	}
	if game.ToString() != solution_difficult {
		t.Errorf("game wasn't solved correctly, expected %s got %s", game.ToString(), solution_difficult)
	}
}

func TestToString(t *testing.T) {
	game, err := ParseBoard(testboard)
	if err != nil {
		t.Errorf("error when parsing, %e", err)
	}
	if game.ToString() != testboard {
		t.Errorf("game wasn't serialized correctly, expected %s got %s", game.ToString(), testboard)
    }
}

func TestCopy(t *testing.T) {
    const test = "070000043040009610800634900094052000358460020000800530080070091902100005007040802"
	game, err := ParseBoard(test)
	if err != nil {
		t.Errorf("error when parsing, %e", err)
	}
    newgame := game.Copy()
    newgame.board[0][0] = uint16(1)
    const expected = "170000043040009610800634900094052000358460020000800530080070091902100005007040802"
	if newgame.ToString() != expected {
		t.Errorf("newgame wasn't modified correctly, expected %s got %s", newgame.ToString(), expected)
    }
	if game.ToString() != test {
		t.Errorf("old game was also modified")
    }
}

func TestToInt(t *testing.T) {
	got := ToInt(0b1)
	if got != 1 {
		t.Errorf("ToInt(0b1) = %d; want 1", got)
	}
	got = ToInt(0b100)
	if got != 3 {
		t.Errorf("ToInt(0b100) = %d; want 3", got)
	}
	got = ToInt(0b1000)
	if got != 4 {
		t.Errorf("ToInt(0b1000) = %d; want 4", got)
	}
}
