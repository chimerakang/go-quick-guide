# Struct 結構

有時您需要以單一結構表示一組欄位。 例如，當您需要撰寫薪資程式時，便需要使用員工資料結構。 在 Go 中，您可以使用結構將可形成記錄的不同欄位組成群組。

Go 中的「結構」是「另一種資料類型」，其可包含零或多個任意類型的欄位，並將其表示為單一實體。

在此節中，我們將探索結構的必要性，以及其使用方式。

---

## 針對結構進行宣告及初始化
若要宣告結構，您必須使用 struct 關鍵字，以及您想要新資料類型具備的欄位及其類型清單。 例如，若要定義員工結構，您可以使用下列程式碼：

```Go
type Employee struct {
    ID        int
    FirstName string
    LastName  string
    Address   string
}
```
然後，您可以使用新類型來宣告變數，就像您通常會搭配其他類型使用的方式一樣，如下：

```Go
var john Employee
```
而且，如果您想要同時對變數進行宣告及初始化，您可以透過這種方式來達成：

```Go

employee := Employee{1001, "John", "Doe", "Doe's Street"}
```
請注意，您必須為結構中的每個欄位指定值。 但那有時可能會造成問題。 或者，您可以更具體地描述您想要在結構中將其初始化的欄位：

```Go
employee := Employee{LastName: "Doe", FirstName: "John"}
```
請注意，在上一個陳述式中，您將值指派給每個欄位的順序並不重要。 此外，如果您未針對任何其他欄位指定值，也不會有任何影響。 Go 將根據欄位資料類型來指派預設值。

若要存取結構的個別欄位，您可以使用點標記法 (.) 來達成，如下列範例所示：

```Go
employee.ID = 1001
fmt.Println(employee.FirstName)
```
最後，您可以使用 & 運算子來產生針對結構的指標，如下列程式碼所示：

```Go
package main

import "fmt"

type Employee struct {
    ID        int
    FirstName string
    LastName  string
    Address   string
}

func main() {
    employee := Employee{LastName: "Doe", FirstName: "John"}
    fmt.Println(employee)
    employeeCopy := &employee
    employeeCopy.FirstName = "David"
    fmt.Println(employee)
}
```
當您執行上述程式碼時，會看到下列輸出：
```
{0 John Doe }
{0 David Doe }
```
請注意，當您使用指標時，結構會變成可變的。

---
## 結構內嵌
Go 中的結構可讓您在結構中內嵌另一個結構。 有時您會想要減少重複次數，並重複使用通用結構。 例如，假設您想要將先前的程式碼重構，來使其中一個資料類型適用於員工，並使另一個資料類型適用於約聘員工。 您可以有一個 Person 結構來保留通用欄位，如下列範例所示：

```Go
type Person struct {
    ID        int
    FirstName string
    LastName  string
    Address   string
}
```
然後，您可以宣告內嵌 Person 類型的其他類型，例如 Employee 和 Contractor。 若要內嵌另一個結構，您必須建立新欄位，如下列範例所示：

```Go
type Employee struct {
    Information Person
    ManagerID   int
}
```
但是，若要參考 Person 結構中的欄位，您必須包括來自員工變數的 Information 欄位，如下列範例所示：

```Go
var employee Employee
employee.Information.FirstName = "John"
```
如果您像我們正在做的一樣對程式碼進行重構，將會使程式碼中斷。 或者，您可以包括與您要內嵌的結構具有相同名稱的新欄位，如下列範例所示：

```Go
type Employee struct {
    Person
    ManagerID int
}
```
作為示範，您可以使用以下程式碼：

```Go
package main

import "fmt"

type Person struct {
    ID        int
    FirstName string
    LastName  string
    Address   string
}

type Employee struct {
    Person
    ManagerID int
}

type Contractor struct {
    Person
    CompanyID int
}

func main() {
    employee := Employee{
        Person: Person{
            FirstName: "John",
        },
    }
    employee.LastName = "Doe"
    fmt.Println(employee.FirstName)
}
```
請注意，您可以從 Employee 結構存取 FirstName 欄位，而不需指定 Person 欄位，因為其會自動內嵌其所有欄位。 但是，當您將結構初始化時，便必須指定要指派值的欄位。

---

## 使用 JSON 對結構進行編碼和解碼
最後，您可以使用結構，以 JSON 對資料進行編碼和解碼。 Go 對 JSON 格式提供絕佳支援，而且已經隨附於標準程式庫套件中。

您也可以進行如重新命名結構中的欄位等動作。 例如，假設您不想要 JSON 輸出顯示 FirstName，而是只顯示 name，或忽略空白欄位。 您可以使用如下列範例所示的欄位標籤：

```Go
type Person struct {
    ID        int    
    FirstName string `json:"name"`
    LastName  string
    Address   string `json:"address,omitempty"`
}
```
然後，若要將結構編碼為 JSON，您可以使用 json.Marshal 函式。 此外，若要將 JSON 字串解碼為資料結構，您可以使用 json.Unmarshal 函式。 以下範例會將所有項目放在一起、將員工的陣列編碼為 JSON，並將輸出解碼為新變數：

```Go
package main

import (
    "encoding/json"
    "fmt"
)

type Person struct {
    ID        int
    FirstName string `json:"name"`
    LastName  string
    Address   string `json:"address,omitempty"`
}

type Employee struct {
    Person
    ManagerID int
}

type Contractor struct {
    Person
    CompanyID int
}

func main() {
    employees := []Employee{
        Employee{
            Person: Person{
                LastName: "Doe", FirstName: "John",
            },
        },
        Employee{
            Person: Person{
                LastName: "Campbell", FirstName: "David",
            },
        },
    }

    data, _ := json.Marshal(employees)
    fmt.Printf("%s\n", data)

    var decoded []Employee
    json.Unmarshal(data, &decoded)
    fmt.Printf("%v", decoded)
}
```
當您執行上述程式碼時，會看到下列輸出：
```
[{"ID":0,"name":"John","LastName":"Doe","ManagerID":0},{"ID":0,"name":"David","LastName":"Campbell","ManagerID":0}]
[{{0 John Doe } 0} {{0 David Campbell } 0}]
```

---

[下一個單元](./fibonacci.md): 資料結構作業