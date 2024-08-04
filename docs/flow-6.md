# 解決方案 - 控制流程練習
讓我們看看每個練習的可能解決方案。

---
## 撰寫 FizzBuzz 程式
使用 switch 陳述式的練習解決方案可能如下所示：

```Go
package main

import (
    "fmt"
    "strconv"
)

func fizzbuzz(num int) string {
    switch {
    case num%15 == 0:
        return "FizzBuzz"
    case num%3 == 0:
        return "Fizz"
    case num%5 == 0:
        return "Buzz"
    }
    return strconv.Itoa(num)
}

func main() {
    for num := 1; num <= 100; num++ {
        fmt.Println(fizzbuzz(num))
    }
}
```
在 FizzBuzz 案例中，3 會乘以 5，因為結果是可由 3 和 5 整除的數字。 您也可以包含 AND 條件，以檢查數字可否由 3 和 5 整除。

---
## 尋找質數
若要尋找小於 20 的質數，練習的解決方案可能如下所示：

```Go
package main

import "fmt"

func findprimes(number int) bool {	
	for i := 2; i < number; i++ {
        if number % i == 0 {
			return false
        }
    }

	if number > 1 {
		return true
	} else {
	    return false
	}	
}

func main() {
    fmt.Println("Prime numbers less than 20:")
	
    for number := 1; number <= 20; number++ {
        if findprimes(number) {
            fmt.Printf("%v ", number)
        }
    }
}
```
在 main 函式中，我們會執行迴圈從 1 到 20，並呼叫 findprimes 函式來檢查目前的數字。 在 findprimes 函式中，我們會在 2 開始 for 迴圈，並重複直到計數器大於 number 值為止。 如果 number 可由計數器整除，則 number 不是質數。 如果我們在未結束的情況下完成迴圈，則數字會為 1 或為其質數。

輸出如下：
```
Prime numbers less than 20:
2 3 5 7 11 13 17 19 
```

---
## 要求一個數字，如為負數，則出現緊急狀況
嘗試 panic 呼叫的練習解決方案可能如下所示：

```Go
package main

import "fmt"

func main() {
    val := 0

    for {
        fmt.Print("Enter number: ")
        fmt.Scanf("%d", &val)

        switch {
        case val < 0:
            panic("You entered a negative number!")
        case val == 0:
            fmt.Println("0 is neither negative nor positive")
        default:
            fmt.Println("You entered:", val)
        }
    }
}
```

請記住，此練習旨在熟悉無限迴圈和 switch 陳述式。