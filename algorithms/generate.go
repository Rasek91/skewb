package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/Rasek91/skewb"
)

var (
	allPreMoves   = []string{""}
	allSolveMoves = map[int][]string{
		0: {""},
		1: {},
		2: {},
		3: {},
		4: {},
		5: {},
		6: {},
		7: {},
		8: {},
		9: {},
	}
	preMoves   = []string{"x", "x'", "x2", "y", "y'", "y2", "z", "z'", "z2"}
	solveMoves = []string{"F", "F'", "f", "f'", "R", "R'", "r", "r'", "b", "b'"}
)

func iteratorPreMoves(previousMoves string, currentIteration, maxIteration int) {
	for _, move := range preMoves {
		if currentIteration == 1 {
			if (previousMoves != "") && (len(strings.Split(previousMoves, " ")) == currentIteration) && (lastMoveIsDifferent(previousMoves, move)) {
				move = fmt.Sprintf("%v %v", previousMoves, move)
			}
		} else {
			if (previousMoves != "") && (len(strings.Split(previousMoves, " ")) == currentIteration-1) && (lastMoveIsDifferent(previousMoves, move)) {
				move = fmt.Sprintf("%v %v", previousMoves, move)
			}
		}

		if currentIteration != maxIteration {
			iteratorPreMoves(move, currentIteration+1, maxIteration)
		} else {
			if len(strings.Split(move, " ")) == maxIteration {
				if !slices.ContainsFunc(allPreMoves, func(newMove string) bool {
					s1 := skewb.New("#FFFFFFFF", "#00FF00FF", "#FF0000FF", "#0000FFFF", "#D67200FF", "#FBFF00FF")
					s2 := skewb.New("#FFFFFFFF", "#00FF00FF", "#FF0000FF", "#0000FFFF", "#D67200FF", "#FBFF00FF")
					s1.ApplyRubiskewbMoves(move)
					s2.ApplyRubiskewbMoves(newMove)

					return s1.ExactEqual(&s2)
				}) {
					allPreMoves = append(allPreMoves, move)
				}
			}
		}
	}
}

func iteratorSolveMoves(previousMoves string, currentIteration, maxIteration int) {
	for _, move := range solveMoves {
		if currentIteration == 1 {
			if (previousMoves != "") && (len(strings.Split(previousMoves, " ")) == currentIteration) && (lastMoveIsDifferent(previousMoves, move)) {
				move = fmt.Sprintf("%v %v", previousMoves, move)
			}
		} else {
			if (previousMoves != "") && (len(strings.Split(previousMoves, " ")) == currentIteration-1) && (lastMoveIsDifferent(previousMoves, move)) {
				move = fmt.Sprintf("%v %v", previousMoves, move)
			}
		}

		if currentIteration != maxIteration {
			iteratorSolveMoves(move, currentIteration+1, maxIteration)
		} else {
			moveNumber := len(strings.Split(move, " "))

			if moveNumber == maxIteration {
				isIn := false

				for i := 0; i <= moveNumber; i++ {
					if slices.ContainsFunc(allPreMoves, func(newMove string) bool {
						s1 := skewb.New("#FFFFFFFF", "#00FF00FF", "#FF0000FF", "#0000FFFF", "#D67200FF", "#FBFF00FF")
						s2 := skewb.New("#FFFFFFFF", "#00FF00FF", "#FF0000FF", "#0000FFFF", "#D67200FF", "#FBFF00FF")
						s1.ApplyRubiskewbMoves(move)
						s2.ApplyRubiskewbMoves(newMove)

						return s1.Equal(&s2)
					}) {
						isIn = true

						break
					}
				}

				if !isIn {
					allSolveMoves[moveNumber] = append(allSolveMoves[moveNumber], move)
				}
			}
		}
	}
}

func lastMoveIsDifferent(move1, move2 string) bool {
	if strings.Contains(move1, " ") {
		moveList := strings.Split(move1, " ")

		move1 = moveList[len(moveList)-1]
	}

	switch {
	case ((move1 == "U") || (move1 == "U'")) && ((move2 == "U") || (move2 == "U'")):
		return false
	case ((move1 == "R") || (move1 == "R'")) && ((move2 == "R") || (move2 == "R'")):
		return false
	case ((move1 == "r") || (move1 == "r'")) && ((move2 == "r") || (move2 == "r'")):
		return false
	case ((move1 == "B") || (move1 == "B'")) && ((move2 == "B") || (move2 == "B'")):
		return false
	case ((move1 == "b") || (move1 == "b'")) && ((move2 == "b") || (move2 == "b'")):
		return false
	case ((move1 == "L") || (move1 == "L'")) && ((move2 == "L") || (move2 == "L'")):
		return false
	case ((move1 == "l") || (move1 == "l'")) && ((move2 == "l") || (move2 == "l'")):
		return false
	case ((move1 == "F") || (move1 == "F'")) && ((move2 == "F") || (move2 == "F'")):
		return false
	case ((move1 == "f") || (move1 == "f'")) && ((move2 == "f") || (move2 == "f'")):
		return false
	case ((move1 == "x") || (move1 == "x'") || (move1 == "x2")) && ((move2 == "x") || (move2 == "x'") || (move2 == "x2")):
		return false
	case ((move1 == "y") || (move1 == "y'") || (move1 == "y2")) && ((move2 == "y") || (move2 == "y'") || (move2 == "y2")):
		return false
	case ((move1 == "z") || (move1 == "z'") || (move1 == "z2")) && ((move2 == "z") || (move2 == "z'") || (move2 == "z2")):
		return false
	default:
		return true
	}
}

func main() {
	previousTime := time.Now()
	iteratorPreMoves("", 1, 1)
	iteratorPreMoves("", 1, 2)
	iteratorPreMoves("", 1, 3)
	fmt.Printf("%s %v premoves finished\n", time.Since(previousTime), len(allPreMoves))

	for i := 1; i <= 8; i++ {
		previousTime := time.Now()

		for _, solveMove := range allSolveMoves[i-1] {
			iteratorSolveMoves(solveMove, i, i)
		}

		fmt.Printf("%s %v %v mover scrambles finished\n", time.Since(previousTime), len(allSolveMoves[i]), i)
	}

	archive, err := os.OpenFile(filepath.Join("algorithms", "moves.zip"), os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		panic(err)
	}

	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	defer zipWriter.Close()
	ioWriter, err := zipWriter.Create("allPreMoves.json")

	if err != nil {
		panic(err)
	}

	content, err := json.MarshalIndent(allPreMoves, "", "\t")

	if err != nil {
		panic(err)
	}

	reader := bytes.NewReader(content)
	if _, err := io.Copy(ioWriter, reader); err != nil {
		panic(err)
	}

	ioWriter, err = zipWriter.Create("allSolveMoves.json")

	if err != nil {
		panic(err)
	}

	content, err = json.MarshalIndent(allSolveMoves, "", "\t")

	if err != nil {
		panic(err)
	}

	reader = bytes.NewReader(content)
	if _, err := io.Copy(ioWriter, reader); err != nil {
		panic(err)
	}
}
