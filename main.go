package main

import (
	"fmt"
	"strings"
)

func main() {
	for i, input := range []string{
		"000007020 712030690 080600050 200000100 050000060 008000009 040002080 063080274 020900000",
		"607500020 020009305 000630000 400000010 076000450 080000002 000098000 804100090 060005701",
		"802096000 005018030 106700024 078902105 000105603 001000098 984031700 250049080 600000000",
		"503000070 010309800 000040000 001000009 000006000 090208300 007402080 020060000 000010400",
	} {
		fmt.Printf("=== case %d ===\n", i+1)
		data, ok := solve(input)
		if !ok {
			panic("impossible")
		}
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				fmt.Print(data[i][j])
			}
			fmt.Println()
		}
		fmt.Println()
	}
}

func solve(input string) ([9][9]int, bool) {
	lines := strings.Split(input, " ")
	var data [9][9]int
	var xUsed [9][10]bool
	var yUsed [9][10]bool
	var zUsed [9][10]bool
	for i, line := range lines {
		for j, char := range line {
			data[i][j] = int(char - '0')
			if data[i][j] != 0 {
				xUsed[i][data[i][j]] = true
				yUsed[j][data[i][j]] = true
				zUsed[getZ(i, j)][data[i][j]] = true
			}
		}
	}
	return data, loop(&xUsed, &yUsed, &zUsed, &data, 0, 0)
}

func loop(xUsed, yUsed, zUsed *[9][10]bool, data *[9][9]int, i, j int) bool {
	loopNext := func() bool {
		if i == 8 && j == 8 {
			return true
		}
		if j < 8 {
			return loop(xUsed, yUsed, zUsed, data, i, j+1)
		} else {
			return loop(xUsed, yUsed, zUsed, data, i+1, 0)
		}
	}

	if data[i][j] > 0 {
		return loopNext()
	}

	k := getZ(i, j)
	for value := 1; value <= 9; value++ {
		if xUsed[i][value] || yUsed[j][value] || zUsed[k][value] {
			continue
		}
		data[i][j] = value
		xUsed[i][value] = true
		yUsed[j][value] = true
		zUsed[k][value] = true
		if ok := loopNext(); ok {
			return ok
		}
		data[i][j] = 0
		xUsed[i][value] = false
		yUsed[j][value] = false
		zUsed[k][value] = false
	}

	return false
}

func getZ(i, j int) int {
	return (i/3)*3 + j/3
}
