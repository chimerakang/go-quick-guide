# 解決方案 - 資料結構作業

---
## 撰寫程式來計算 Fibonacci 數列


此挑戰的解決方法之一如下列範例所示：

```Go

package main

import "fmt"

func fibonacci(n int) []int {
	if n < 2 {
		return make([]int, 0)
	}

	nums := make([]int, n)
	nums[0], nums[1] = 1, 1

	for i := 2; i < n; i++ {
		nums[i] = nums[i-1] + nums[i-2]
	}

	return nums
}

func main() {
	var num int

	fmt.Print("What's the Fibonacci sequence you want? ")
	fmt.Scanln(&num)
	fmt.Println("The Fibonacci sequence is:", fibonacci(num))
}
```

---
## 建立羅馬數字轉譯程式
此挑戰的解決方法之一如下列範例所示：

```Go
package main

import (
	"fmt"
)

func romanToArabic(numeral string) int {
	romanMap := map[rune]int{
		'M': 1000,
		'D': 500,
		'C': 100,
		'L': 50,
		'X': 10,
		'V': 5,
		'I': 1,
	}

	arabicVals := make([]int, len(numeral)+1)

	for index, digit := range numeral {
		if val, present := romanMap[digit]; present {
			arabicVals[index] = val
		} else {
			fmt.Printf("Error: The roman numeral %s has a bad digit: %c\n", numeral, digit)
			return 0
		}
	}

	total := 0

	for index := 0; index < len(numeral); index++ {
		if arabicVals[index] < arabicVals[index+1] {
			arabicVals[index] = -arabicVals[index]
		}
		total += arabicVals[index]
	}

	return total
}

func main() {
	fmt.Println("MCLX is", romanToArabic("MCLX"))
	fmt.Println("MCMXCIX is ", romanToArabic("MCMXCIX"))
	fmt.Println("MCMZ is", romanToArabic("MCMZ"))
}
```