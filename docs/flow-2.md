# 使用 switch 陳述式的控制流程
和其他程式設計語言一樣，Go 也支援 switch 陳述式。 使用 switch 陳述式可避免鏈結多個 if 陳述式。 若使用 switch 陳述式，則在維護和閱讀包含許多 if 陳述式的程式碼時可免去許多麻煩， 這些陳述式也可讓建構複雜條件變得更輕鬆。 下列各節將探討 switch 陳述式。

## 基本 switch 語法
和 if 陳述式一樣，switch 條件不需要括弧。 最簡單的 switch 陳述式看起來像這樣：

```Go

package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    sec := time.Now().Unix()
    rand.Seed(sec)
    i := rand.Int31n(10)

    switch i {
    case 0:
        fmt.Print("zero...")
    case 1:
        fmt.Print("one...")
    case 2:
        fmt.Print("two...")
    }

    fmt.Println("ok")
}
```
如果執行多次上述程式碼，每次都會看到不同的輸出。 (但如果在 Go Playground 中執行這段程式碼，每次都會得到相同的結果，這是服務的限制之一。)

Go 會比較 switch 陳述式的每個案例，直到找出符合條件的結果。 但請注意，之前的程式碼並未涵蓋 num 變數值的所有可能案例。 如果 num 的結果是 5，則程式輸出就是 ok。

或者，您也可以更明確地指定預設使用案例，並以下列方式包含此案例：

```Go

switch i {
case 0:
    fmt.Print("zero...")
case 1:
    fmt.Print("one...")
case 2:
    fmt.Print("two...")
default:
    fmt.Print("no match...")
}
```
請注意，您不需要為 default 案例撰寫驗證運算式， i 變數的值會根據 case 陳述式進行驗證，而案例會 default 處理任何未驗證的值。

---
## 使用多個運算式
有時候，一個 case 陳述式會有多個符合條件的運算式。 在 Go 中，如果您想要 case 陳述式包含多個運算式，請使用逗號 (,) 分隔運算式， 這個技巧可讓您避免重複的程式碼。

下列程式碼範例示範如何包含多個運算式。

```Go
package main

import "fmt"

func location(city string) (string, string) {
    var region string
    var continent string
    switch city {
    case "Delhi", "Hyderabad", "Mumbai", "Chennai", "Kochi":
        region, continent = "India", "Asia"
    case "Lafayette", "Louisville", "Boulder":
        region, continent = "Colorado", "USA"
    case "Irvine", "Los Angeles", "San Diego":
        region, continent = "California", "USA"
    default:
        region, continent = "Unknown", "Unknown"
    }
    return region, continent
}
func main() {
    region, continent := location("Irvine")
    fmt.Printf("John works in %s, %s\n", region, continent)
}
```
請注意，您在 case 陳述式的運算式中包含的值會對應至 switch 陳述式所驗證的變數資料類型。 如果您將整數值包含為新的 case 陳述式，就不會編譯程式。

## 呼叫函式
switch 也可以呼叫函式。 您可以在該函式中撰寫 case 陳述式以取得可能的傳回值。 例如，下列程式碼會呼叫 time.Now() 函式。 其列印的輸出取決於當時為星期幾。

```Go

package main

import (
    "fmt"
    "time"
)

func main() {
    switch time.Now().Weekday().String() {
    case "Monday", "Tuesday", "Wednesday", "Thursday", "Friday":
        fmt.Println("It's time to learn some Go.")
    default:
        fmt.Println("It's the weekend, time to rest!")
    }

    fmt.Println(time.Now().Weekday().String())
}
```
當您從 switch 陳述式呼叫函式時，您可在不變更運算式的情況下修改其邏輯，因為您一直在驗證函式傳回的內容。

您也可以從 case 陳述式呼叫函式。 例如，您可以使用這個技巧，使用規則運算式來比對特定的模式。 以下是範例：

```Go
package main

import "fmt"

import "regexp"

func main() {
    var email = regexp.MustCompile(`^[^@]+@[^@.]+\.[^@.]+`)
    var phone = regexp.MustCompile(`^[(]?[0-9][0-9][0-9][). \-]*[0-9][0-9][0-9][.\-]?[0-9][0-9][0-9][0-9]`)

    contact := "foo@bar.com"

    switch {
    case email.MatchString(contact):
        fmt.Println(contact, "is an email")
    case phone.MatchString(contact):
        fmt.Println(contact, "is a phone number")
    default:
        fmt.Println(contact, "is not recognized")
    }
}
```
請注意，switch 區塊中沒有驗證運算式。 我們會在下一節討論此概念。

## 省略條件
使用 Go 時，可以在 switch 陳述式中省略條件，就像在 if 陳述式中做的一樣。 這種模式就像是您在強制 switch 陳述式不停執行時，比較 true 值一樣。

以下範例示範如何撰寫不含條件的 switch 陳述式：

```Go
package main

import (
    "fmt"
    "math/rand"
    "time"
)

func main() {
    rand.Seed(time.Now().Unix())
    r := rand.Float64()
    switch {
    case r > 0.1:
        fmt.Println("Common case, 90% of the time")
    default:
        fmt.Println("10% of the time")
    }
}
```

程式一律會執行這種 switch 類型的陳述式，因為條件一律為 true。 條件式 switch 區塊可能比長串的 if 和 else if 陳述式更容易維護。

## 讓邏輯通過下一個案例
在某些程式設計語言中，您要在每個 case 陳述式的結尾寫入 break 關鍵字。 但在 Go 中，除非您明確叫停，否則只要邏輯進入一個案例，就會結束 switch 區塊。 為使邏輯通過下一個緊跟著的案例，請使用 fallthrough 關鍵字。

若要深入了解此模式，請參閱下列程式碼範例。

```Go
package main

import (
    "fmt"
)

func main() {
    switch num := 15; {
    case num < 50:
        fmt.Printf("%d is less than 50\n", num)
        fallthrough
    case num > 100:
        fmt.Printf("%d is greater than 100\n", num)
        fallthrough
    case num < 200:
        fmt.Printf("%d is less than 200", num)
    }
}
```
執行程式碼及分析輸出：

輸出
```
15 is less than 50
15 is greater than 100
15 is less than 200
```
您是否發現任何錯誤？

請注意，因為 num 是 15 (小於 50)，所以符合第一個案例。 但 num 不大於 100。 而且因為第一個 case 陳述式有 fallthrough 關鍵字，所以邏輯會在不驗證案例下立即進入下一個 case 陳述式。 因此，使用 fallthrough 關鍵字時請務必謹慎小心。 您可能不希望出現此程式碼建立的行為。

---
## Next
[下一個單元](./flow-3.md): 使用 for 運算式對資料進行迴圈