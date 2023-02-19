package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func main() {
	//input := "000007020 712030690 080600050 200000100 050000060 008000009 040002080 063080274 020900000"
	input := "607500020 020009305 000630000 400000010 076000450 080000002 000098000 804100090 060005701"
	//input := "802096000 005018030 106700024 078902105 000105603 001000098 984031700 250049080 600000000"
	lines := strings.Split(input, " ")
	var data [9][9]Point
	var xUsed [9][10]bool
	var yUsed [9][10]bool
	var zUsed [9][10]bool
	var usedCnt [10]int
	for i, line := range lines {
		for j, char := range line {
			data[i][j] = Point{
				X:     i,
				Y:     j,
				Value: int(char - '0'),
			}
			if data[i][j].Value != 0 {
				usedCnt[data[i][j].Value] += 1
				data[i][j].Acc = 1
				xUsed[i][data[i][j].Value] = true
				yUsed[j][data[i][j].Value] = true
				zUsed[getZ(i, j)][data[i][j].Value] = true
			}
		}
	}

	//debugInfo, _ := json.Marshal(map[string]interface{}{
	//	"data":  data,
	//	"xUsed": xUsed,
	//	"yUsed": yUsed,
	//	"zUsed": zUsed,
	//})
	//fmt.Println(string(debugInfo))

	// calc acc
	if !calcAllAcc(&xUsed, &yUsed, &zUsed, &data) {
		panic("impossible")
	}

	var attempts int
	if !loop(&xUsed, &yUsed, &zUsed, &usedCnt, &data, 0, &attempts, 1.0) {
		panic("no results")
	}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			fmt.Print(data[i][j].Value)
		}
		fmt.Println()
	}
}

func loop(xUsed, yUsed, zUsed *[9][10]bool, usedCnt *[10]int, data *[9][9]Point, deep int, attempts *int, acc float64) bool {
	*attempts += 1
	fmt.Printf("deep %d attempt %d acc %f\n", deep, *attempts, acc)

	// check is ok
	ok := true
	for i := 0; i < 9 && ok; i++ {
		for j := 0; j < 9 && ok; j++ {
			if data[i][j].Value == 0 {
				ok = false
			}
		}
	}
	if ok {
		return true
	}

	// calc steps
	var steps []Step
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if data[i][j].Value == 0 {
				for _, v := range data[i][j].CanUse {
					steps = append(steps, Step{
						X:     i,
						Y:     j,
						Value: v,
						Acc:   data[i][j].Acc,
					})
				}
			}
		}
	}
	canSteps := NewSteps(steps, usedCnt)
	sort.Sort(canSteps)

	if deep == 0 {
		fmt.Printf("deep %d total %d steps\n", deep, len(canSteps.steps))
	}
	failedMap := map[string]bool{}
mainLoop:
	for stepIndex, step := range canSteps.steps {
		stepIndex = stepIndex
		stepKey := getStepKey(step.X, step.Y, step.Value)
		if failedMap[stepKey] {
			continue
		}
		failedMap[stepKey] = true
		//if deep >= 0 {
		//	fmt.Printf("deep %d step %d\n", deep, stepIndex+1)
		//	for i := 0; i < 9; i++ {
		//		for j := 0; j < 9; j++ {
		//			fmt.Print(data[i][j].Value)
		//		}
		//		fmt.Println()
		//	}
		//	fmt.Println()
		//}
		dataNew := *data
		xUsedNew, yUsedNew, zUsedNew := *xUsed, *yUsed, *zUsed
		data[step.X][step.Y].Value = step.Value
		xUsed[step.X][step.Value] = true
		yUsed[step.Y][step.Value] = true
		zUsed[getZ(step.X, step.Y)][step.Value] = true
		for i := 0; i < 9; i++ {
			if data[i][step.Y].Value == 0 {
				data[i][step.Y].Acc = calcAcc(xUsed, yUsed, zUsed, data, i, step.Y)
				if data[i][step.Y].Acc < 0 {
					*xUsed, *yUsed, *zUsed, *data = xUsedNew, yUsedNew, zUsedNew, dataNew
					//fmt.Printf("deep %d step %d failed and rollback condition 1\n", deep, stepIndex+1)
					continue mainLoop
				}
			}
			if data[step.X][i].Value == 0 {
				data[step.X][i].Acc = calcAcc(xUsed, yUsed, zUsed, data, step.X, i)
				if data[step.X][i].Acc < 0 {
					*xUsed, *yUsed, *zUsed, *data = xUsedNew, yUsedNew, zUsedNew, dataNew
					//fmt.Printf("deep %d step %d failed and rollback contidion 2\n", deep, stepIndex+1)
					continue mainLoop
				}
			}
		}
		for i := step.X / 3 * 3; i < step.X/3*3+3; i++ {
			for j := step.Y / 3 * 3; j < step.Y/3*3+3; j++ {
				if data[i][j].Value == 0 {
					data[i][j].Acc = calcAcc(xUsed, yUsed, zUsed, data, i, j)
					if data[i][j].Acc < 0 {
						*xUsed, *yUsed, *zUsed, *data = xUsedNew, yUsedNew, zUsedNew, dataNew
						//fmt.Printf("deep %d step %d failed and rollback contidion 3\n", deep, stepIndex+1)
						continue mainLoop
					}
				}
			}
		}
		if !loop(xUsed, yUsed, zUsed, usedCnt, data, deep+1, attempts, acc*step.Acc) {
			*xUsed, *yUsed, *zUsed, *data = xUsedNew, yUsedNew, zUsedNew, dataNew
			//fmt.Printf("deep %d step %d failed and rollback\n", deep, stepIndex+1)
			continue
		}
		return true
	}

	return false
}

func getStepKey(i, j int, value int) string {
	return strconv.Itoa(i) + "_" + strconv.Itoa(j) + "_" + strconv.Itoa(value)
}

func getZ(i, j int) int {
	return (i/3)*3 + j/3
}

func calcAllAcc(xUsed, yUsed, zUsed *[9][10]bool, data *[9][9]Point) bool {
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if data[i][j].Value > 0 {
				xUsed[i][data[i][j].Value] = true
				yUsed[j][data[i][j].Value] = true
				zUsed[getZ(i, j)][data[i][j].Value] = true
			}
		}
	}
	for i := 0; i < 9; i++ {
		for j := 0; j < 9; j++ {
			if data[i][j].Value == 0 {
				data[i][j].Acc = calcAcc(xUsed, yUsed, zUsed, data, i, j)
				if data[i][j].Acc < 0 {
					return false
				}
			}
		}
	}
	return true
}

func calcAcc(xUsed, yUsed, zUsed *[9][10]bool, data *[9][9]Point, i, j int) float64 {
	var canUse = [10]bool{true, true, true, true, true, true, true, true, true, true}
	for value, used := range xUsed[i] {
		if value > 0 && used {
			canUse[value] = false
		}
	}
	for value, used := range yUsed[j] {
		if value > 0 && used {
			canUse[value] = false
		}
	}
	for value, used := range zUsed[getZ(i, j)] {
		if value > 0 && used {
			canUse[value] = false
		}
	}
	data[i][j].CanUse = nil
	for value, can := range canUse {
		if value > 0 && can {
			data[i][j].CanUse = append(data[i][j].CanUse, value)
		}
	}
	if len(data[i][j].CanUse) == 0 {
		return -1
	}
	return 1.0 / float64(len(data[i][j].CanUse))
}
