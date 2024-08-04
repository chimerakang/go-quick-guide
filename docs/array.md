# Array 陣列
by [@chimerakang](https://github.com/chimerakang)

---

## 介紹

Go 的陣列是特定類型的固定長度資料結構。 其可以有零或多個元素，而您必須在加以宣告或初始化時定義大小。 此外，「您無法在加以建立之後調整其大小」。 基於這些理由，陣列並不常用於 Go 程式，但其為配量和對應的基礎。

## 宣告陣列

若要在 Go 中宣告陣列，您必須定義其元素的資料類型，以及陣列可保存的元素數目。 然後，您可以使用下標標記法來存取陣列中的每個元素，其中零是第一個元素，而最後一個元素則是陣列長度減一 (長度 - 1)。

讓我們使用下列程式碼作為範例：

```go
package main

import "fmt"

func main() {
    var a [3]int
    a[1] = 10
    fmt.Println(a[0])
    fmt.Println(a[1])
    fmt.Println(a[len(a)-1])
}
```

當您執行上述程式碼時，會得到如下的輸出：
```
0
10
0
```

即使您已宣告陣列，您在存取其元素時也不會收到錯誤。 依預設，Go 會利用預設資料類型來將每個元素初始化。 在此案例中，`int` 的預設值為零。 但是，您可以將值指派至特定位置，如同我們處理 `a[1] = 10` 一樣。 而且，您可以使用相同的標記法來存取該元素。 另外，請注意，為了參考第一個元素，我們使用了 `a[0]`。 為了參考最後一個元素，我們使用了 `a[len(a)-1]`。 `len` 函式是 Go 中的內建函式，可用來取得陣列、配量或對應中的元素數目。

## 將陣列初始化

您也可以在宣告陣列時，使用預設值以外的值來將陣列初始化。 例如，您可以使用下列程式碼來查看和測試語法：

```go
package main

import "fmt"

func main() {
    cities := [5]string{"New York", "Paris", "Berlin", "Madrid"}
    fmt.Println("Cities:", cities)
}
```

執行上述程式碼，您應該會看到下列輸出：

```
Cities: [New York Paris Berlin Madrid ]

```

雖然陣列應該要有五個元素，我們並不需要指派值給所有元素。 如先前已知的，最後的位置具有空字串，因為它是字串資料類型的預設值。

## 陣列中的省略符號

當您不知道將需要多少個位置，但知道您擁有多少組資料元素時，另一種對陣列進行宣告和初始化的方式是使用省略符號 (`...`)，如下列範例所示：

```go
q := [...]int{1, 2, 3}

```

讓我們修改上一節所使用的程式來使用省略符號。 該程式碼看起來應該像下列範例：

```go
package main

import "fmt"

func main() {
    cities := [...]string{"New York", "Paris", "Berlin", "Madrid"}
    fmt.Println("Cities:", cities)
}

```

執行上述程式碼，您應該會看到類似的輸出，如下列範例所示：

```
Cities: [New York Paris Berlin Madrid]

```

您可以看出差異嗎？ 結尾沒有空字串。 陣列長度是由您在加以初始化時所放置的字串所決定。 您不會保留您不知道最後是否會需要的記憶體。

將陣列初始化的另一種有趣方式是使用省略符號，並只指定最後位置的值。 例如，使用下列程式碼：

```go
package main

import "fmt"

func main() {
    numbers := [...]int{99: -1}
    fmt.Println("First Position:", numbers[0])
    fmt.Println("Last Position:", numbers[99])
    fmt.Println("Length:", len(numbers))
}
```

執行此程式碼，您將會得到此輸出：

```
First Position: 0
Last Position: -1
Length: 100

```

請注意陣列長度為 100，因為您為第 99 個位置指定值。 第一個位置會列印出預設值 (零)。

## 多維陣列

當您需要使用複雜的資料結構時，Go 支援多維度陣列。 讓我們建立一個程式，您會在其中對二維陣列進行宣告和初始化。 使用下列程式碼：

```go
package main

import "fmt"

func main() {
	var a[5] int
	fmt.Println("a:", a) // by default, an array is *zero-valued*

	a[4] = 100 // access by operator[]
	fmt.Println("after set, a:", a)
	fmt.Println("a[4]:", a[4])

	fmt.Println("len(a):", len(a)) // keyword len

	b := [5] int {1,2,3,4,5} // declare and initlize
	fmt.Println("b:", b)

	var tmp[2][3] int // two dimession
	fmt.Println("len(tmp):", len(tmp))
	fmt.Println("len(tmp[0]):", len(tmp[0]))
	fmt.Println("tmp:", tmp)
	for i := 0; i < 2; i++ {
		for j := 0; j < 3; j++ {
			tmp[i][j] = i + j
		}
	}
	fmt.Println("tmp:", tmp)
}
```

執行上述程式，您應該會看到類似下列範例的輸出：

```
$ go run arrays/arrays.go 
a: [0 0 0 0 0]
after set, a: [0 0 0 0 100]
a[4]: 100
len(a): 5
b: [1 2 3 4 5]
len(tmp): 2
len(tmp[0]): 3
tmp: [[0 0 0] [0 0 0]]
tmp: [[0 1 2] [1 2 3]]
```

您已宣告一個二維陣列，其會指定陣列在第二個維度中有多少個位置，例如 `var tmp [2][3]int`。 您可以將這個陣列視為具有欄和列的資料結構，例如試算表或矩陣。 到目前為止，所有位置的預設值均為零。 在 `for` 迴圈中，我們會針對每個列上具有不同值模式的每個位置進行初始化。 最後，將其所有值列印到終端。

如果您想要宣告三維陣列，該怎麼辦？ 您應該能猜到語法是什麼，對吧？ 您可以依下列範例所示的方式來執行：

```go
package main

import "fmt"

func main() {
    var threeD [3][5][2]int
    for i := 0; i < 3; i++ {
        for j := 0; j < 5; j++ {
            for k := 0; k < 2; k++ {
                threeD[i][j][k] = (i + 1) * (j + 1) * (k + 1)
            }
        }
    }
    fmt.Println("\nAll at once:", threeD)
}
```

執行上述程式碼，您應該會看到類似下列範例的輸出：

```go
All at once: [[[1 2] [2 4] [3 6] [4 8] [5 10]] [[2 4] [4 8] [6 12] [8 16] [10 20]] [[3 6] [6 12] [9 18] [12 24] [15 30]]]
```

如果我們將輸出格式化為更容易閱讀的格式，您便可能會有類似下列範例的內容：

```go
All at once: 
[
    [
        [1 2] [2 4] [3 6] [4 8] [5 10]
    ] 
    [
        [2 4] [4 8] [6 12] [8 16] [10 20]
    ] 
    [
        [3 6] [6 12] [9 18] [12 24] [15 30]
    ]
]
```

請注意結構與二維陣列之結構間的差異。 您可以視需要繼續增加更多維度，但我們的示範就到此為止，因為還有其他資料類型需要探索。

---
## Next
[Slice(切片)](./slice.md)
