# Map 映射

Go 中的對應基本上是雜湊表，其為機碼和值組的集合。 對應中的所有機碼都必須是相同類型，而值也是如此。 但是，您可以使用不同類型的機碼和值。 例如，機碼可以是數字，而值可以是字串。 若要存取對應中的特定項目，您可以參考其機碼。

---

## 針對Map進行宣告及初始化
若要宣告Map，您必須使用 map 關鍵字。 然後，您會定義機碼和值類型，如下所示：map[T]T。 例如，若您想要建立包含學生年齡的 map，可以使用下列程式碼：

```Go
package main

import "fmt"

func main() {
    studentsAge := map[string]int{
        "john": 32,
        "bob":  31,
    }
    fmt.Println(studentsAge)
}
```
當您執行上述程式碼時，會看到下列輸出：
```
map[bob:31 john:32]
```
您可以使用內建 make() 函式來在上一個區段中建立map。 您可以使用下列程式碼來建立空map：

```Go
studentsAge := make(map[string]int)
```
map是動態的。您可以在建立 map 之後新增、存取或移除項目。讓我們探索那些動作。

---
## 新增項目
若要新增項目，您不需像 slice 那樣使用內建函式。map 會更簡單明瞭。您只需定義key和 value。 如果配對不存在，則會將項目新增至 map

讓我們使用 make 函式來重寫先前用來建立 map 的程式碼。接著我們會將項目新增至map。您可以使用以下程式碼：

```Go
package main

import "fmt"

func main() {
    studentsAge := make(map[string]int)
    studentsAge["john"] = 32
    studentsAge["bob"] = 31
    fmt.Println(studentsAge)
}
```
當您執行程式碼時，會得到與先前相同的輸出：
```
map[bob:31 john:32]
```
請注意，我們已將項目新增至已初始化的 map。 但是，若您嘗試使用 nil 對應進行相同的作業，就會收到錯誤訊息。 例如，下列程式碼將無法運作：

```Go
package main

import "fmt"

func main() {
    var studentsAge map[string]int
    studentsAge["john"] = 32
    studentsAge["bob"] = 31
    fmt.Println(studentsAge)
}
```
當您執行上述程式碼時，會發生下列錯誤：

```
panic: assignment to entry in nil map

goroutine 1 [running]:
main.main()
        demos/helloworld/main.go:7 +0x4f
exit status 2
```
若要避免在將項目新增至對應時遇到問題，請確定您是使用 make 函式 (如同我們在先前程式碼片段中所示) 來建立空白的對應 (而不是 nil 對應)。 只有當您新增項目時，才適用此規則。 如果您在 nil 對應中執行查閱、刪除或迴圈作業，Go 並不會發生緊急錯誤。 我們很快就會確認該行為。

---
## 存取項目
若要存取 map 中的項目，您會和處理array或slice一樣，使用一般的下標標記法 m[key]。 以下是如何存取項目的簡單範例：

```Go
package main

import "fmt"

func main() {
    studentsAge := make(map[string]int)
    studentsAge["john"] = 32
    studentsAge["bob"] = 31
    fmt.Println("Bob's age is", studentsAge["bob"])
}
```
當您在 map 中使用下標標記法時，即使key不存在 map 上，也一律會得到回應。當您存取不存在的項目時，Go 並不會發生緊急錯誤。相反地，您會取得預設值。您可以使用下列程式碼來確認該行為：

```Go
package main

import "fmt"

func main() {
    studentsAge := make(map[string]int)
    studentsAge["john"] = 32
    studentsAge["bob"] = 31
    fmt.Println("Christy's age is", studentsAge["christy"])
}
```
當您執行上述程式碼時，會看到下列輸出：
```
Christy's age is 0
```
在許多情況下，當您存取不存在於 map 中的項目時，Go 不會傳回錯誤的行為是合情合理的。 但有時候您會需要知道某個項目是否存在。 在 Go 中，map 的下標標記法可能會產生兩個值。 第一個是項目的value。 第二個是布林值旗標，表示key是否存在。

若要修正上一個程式碼片段的問題，您可以使用下列程式碼：

```Go
package main

import "fmt"

func main() {
    studentsAge := make(map[string]int)
    studentsAge["john"] = 32
    studentsAge["bob"] = 31

    age, exist := studentsAge["christy"]
    if exist {
        fmt.Println("Christy's age is", age)
    } else {
        fmt.Println("Christy's age couldn't be found")
    }
}
```
當您執行上述程式碼時，會看到下列輸出：

```
Christy's age couldn't be found
```
您可以使用第二個程式碼片段來檢查對應中的機碼是否存在，然後再加以存取。

---
## 移除項目
若要從 map 中移除項目，請使用內建 delete() 函式。 以下是如何從 map 中移除項目的範例：

```Go
package main

import "fmt"

func main() {
    studentsAge := make(map[string]int)
    studentsAge["john"] = 32
    studentsAge["bob"] = 31
    delete(studentsAge, "john")
    fmt.Println(studentsAge)
}
```

當您執行程式碼時，會取得下列輸出：
```
map[bob:31]
```
如先前所述，如果您嘗試刪除不存在的項目，Go 並不會發生緊急錯誤。 以下是該行為的範例：

```Go
package main

import "fmt"

func main() {
    studentsAge := make(map[string]int)
    studentsAge["john"] = 32
    studentsAge["bob"] = 31
    delete(studentsAge, "christy")
    fmt.Println(studentsAge)
}
```
當您執行程式碼時，您不會收到錯誤，而且會看到下列輸出：

```
map[bob:31 john:32]
```

---

## map中的迴圈
最後，讓我們看看如何在 map 中執行迴圈，以程式設計方式存取其所有項目。 若要這樣做，您可以使用範圍型迴圈，如下列範例所示：

```Go
package main

import (
    "fmt"
)

func main() {
    studentsAge := make(map[string]int)
    studentsAge["john"] = 32
    studentsAge["bob"] = 31
    for name, age := range studentsAge {
        fmt.Printf("%s\t%d\n", name, age)
    }
}
```

當您執行上述程式碼時，會看到下列輸出：
```

john    32
bob     31
```
請注意，您可以將key和value資訊存放在不同變數中。 在此案例中，我們將key保留於 name 變數中，並將值保留於 age 變數中。 因此，range 會先產生項目的key，然後再產生該項目的value。您可以使用 `_` 變數來忽略這兩者中的任何一個，如下列範例所示：

```Go
package main

import (
    "fmt"
)

func main() {
    studentsAge := make(map[string]int)
    studentsAge["john"] = 32
    studentsAge["bob"] = 31

    for _, age := range studentsAge {
        fmt.Printf("Ages %d\n", age)
    }
}
```
儘管在此案例中，以該方式列印年齡並無任何意義，但在某些情況下，您並不需要知道項目的key。或者，您也可以只使用項目的key，如下列範例所示：

```Go
package main

import (
    "fmt"
)

func main() {
    studentsAge := make(map[string]int)
    studentsAge["john"] = 32
    studentsAge["bob"] = 31

    for name := range studentsAge {
        fmt.Printf("Names %s\n", name)
    }
}
```

---
[下一個單元](./struct.md): 使用結構