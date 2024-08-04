package main

import "fmt"

func add(x, y int) int {
	return x + y
}

func multiply(x, y int) int {
	return x * y
}

func power(x, y int) int {
	result := 1
	for i := 0; i < y; i++ {
		result *= x
	}
	return result
}

func main() {
	monsterID := 5
	monsterPower := getPower(monsterID)
	monsterPowerEx := power(monsterPower, 2)
	fmt.Printf("%d 號怪物的力量是 %d, 進化之後力量變成 %d!\n", monsterID, monsterPower, monsterPowerEx)
}
