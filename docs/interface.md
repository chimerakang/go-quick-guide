# 介面interface

by [@chimerakang](https://github.com/chimerakang)

---

## 介紹
編寫靈活、可重複使用且模組化的程式碼對於開發多功能程式至關重要。以這種方式工作可以避免在多個地方進行相同的更改，從而確保程式碼更易於維護。完成此操作的方式因語言而異。例如，繼承是 Java、C++、C# 等語言中使用的常見方法。

開發人員還可以透過組合來實現相同的設計目標。組合是將物件或資料類型組合成更複雜的物件或資料類型的一種方法。這是 Go 用來促進程式碼重複使用、模組化和靈活性的方法。 Go 中的介面提供了一種組織複雜組合的方法，學習如何使用它們將使您能夠建立通用的、可重複使用的程式碼。

在本文中，我們將學習如何編寫具有共同行為的自訂類型，這將使我們能夠重複使用程式碼。我們還將學習如何為我們自己的自訂類型實現接口，以滿足從另一個包定義的接口。

## 定義行為
組合的核心實作之一是介面的使用。介面定義了類型的行為。 Go標準庫中最常用的介面之一是以下fmt.Stringer介面：
```go
type Stringer interface {
    String() string
}
```
第一行程式碼定義了一個type被呼叫的Stringer.然後它聲明它是一個interface.就像定義結構一樣，Go 使用大括號 ( {}) 包圍介面的定義。與定義結構體相比，我們只定義介面的行為；即「這個類型能做什麼」。

對介面來說Stringer，唯一的行為就是String()方法。該方法不帶參數並傳回一個字串。

接下來，讓我們來看看一些具有以下fmt.Stringer行為的程式碼：

```go
package main

import "fmt"

type Article struct {
	Title string
	Author string
}

func (a Article) String() string {
	return fmt.Sprintf("The %q article was written by %s.", a.Title, a.Author)
}

func main() {
	a := Article{
		Title: "Understanding Interfaces in Go",
		Author: "Sammy Shark",
	}
	fmt.Println(a.String())
}
```
我們要做的第一件事是建立一個名為 的新類型Article。該類型有一個Title和一個Author字段，並且都是字串資料類型：
```go
type Article struct {
	Title string
	Author string
}
```
接下來，我們定義一個method被呼叫String的類型Article。該String方法將傳回一個表示Article類型的字串：

```go
func (a Article) String() string {
	return fmt.Sprintf("The %q article was written by %s.", a.Title, a.Author)
}
```
然後，在我們的main 函數中，我們建立該Article類型的實例並將其指派給名為 的變數a。我們"Understanding Interfaces in Go"為該欄位提供 的值Title，並"Sammy Shark"為該Author欄位提供：

```go
a := Article{
	Title: "Understanding Interfaces in Go",
	Author: "Sammy Shark",
}
```
String然後，我們透過呼叫fmt.Println並傳入方法呼叫的結果來列印出方法的結果a.String()：

```go
fmt.Println(a.String())
```
運行程式後您將看到以下輸出：

```
The "Understanding Interfaces in Go" article was written by Sammy Shark.
```
到目前為止，我們還沒有使用接口，但我們確實創建了一個具有行為的類型。該行為與fmt.Stringer介面相符。接下來，讓我們看看如何使用該行為來使我們的程式碼更具可重複使用性。

---
## 定義介面
現在我們已經用所需的行為定義了類型，我們可以看看如何使用該行為。

然而，在此之前，讓我們看看如果我們想從函數中的類型呼叫該String方法，我們需要做什麼：Article

```go
package main

import "fmt"

type Article struct {
	Title string
	Author string
}

func (a Article) String() string {
	return fmt.Sprintf("The %q article was written by %s.", a.Title, a.Author)
}

func main() {
	a := Article{
		Title: "Understanding Interfaces in Go",
		Author: "Sammy Shark",
	}
	Print(a)
}

func Print(a Article) {
	fmt.Println(a.String())
}
```
在此程式碼中，我們新增一個名為的新函數Print，該函數接受 anArticle作為參數。請注意，該函數所做的唯一事情Print就是呼叫該String方法。因此，我們可以定義一個介面來傳遞給函數：

```go
package main

import "fmt"

type Article struct {
	Title string
	Author string
}

func (a Article) String() string {
	return fmt.Sprintf("The %q article was written by %s.", a.Title, a.Author)
}

type Stringer interface {
	String() string
}

func main() {
	a := Article{
		Title: "Understanding Interfaces in Go",
		Author: "Sammy Shark",
	}
	Print(a)
}

func Print(s Stringer) {
	fmt.Println(s.String())
}
```
這裡我們創建了一個名為的介面Stringer：

```go
type Stringer interface {
	String() string
}
```
此Stringer介面只有一個方法，稱為String()返回string.方法是一種特殊函數，其作用域為 Go 中的特定類型。與函數不同，方法只能從定義它的類型的實例中呼叫。

然後我們更新方法的簽章Print以採用 a Stringer，而不是具體型別Article。因為編譯器知道Stringer介面定義了該String方法，所以它只接受也具有該String方法的類型。

現在我們可以將該Print方法與任何滿足Stringer介面的東西一起使用。讓我們創建另一種類型來演示這一點：

```go
package main

import "fmt"

type Article struct {
	Title  string
	Author string
}

func (a Article) String() string {
	return fmt.Sprintf("The %q article was written by %s.", a.Title, a.Author)
}

type Book struct {
	Title  string
	Author string
	Pages  int
}

func (b Book) String() string {
	return fmt.Sprintf("The %q book was written by %s.", b.Title, b.Author)
}

type Stringer interface {
	String() string
}

func main() {
	a := Article{
		Title:  "Understanding Interfaces in Go",
		Author: "Sammy Shark",
	}
	Print(a)

	b := Book{
		Title:  "All About Go",
		Author: "Jenny Dolphin",
		Pages:  25,
	}
	Print(b)
}

func Print(s Stringer) {
	fmt.Println(s.String())
}
```
我們現在新增第二種類型，稱為Book.它也String定義了方法。這意味著它也滿足Stringer接口。因此，我們也可以將其發送到我們的Print函數：

```
The "Understanding Interfaces in Go" article was written by Sammy Shark.
The "All About Go" book was written by Jenny Dolphin. It has 25 pages.
```
到目前為止，我們已經演示瞭如何使用單一介面。然而，一個介面可以定義多個行為。接下來，我們將了解如何透過聲明更多方法來使我們的介面更加通用。

---
## 介面中的多種行為
編寫 Go 程式碼的核心租用戶之一是編寫小而簡潔的類型，並將它們組合成更大、更複雜的類型。編寫接口時也是如此。要了解如何建立接口，我們首先從僅定義一個接口開始。我們將定義兩個形狀，aCircle和Square，它們都將定義一個名為 的方法Area。此方法將返回它們各自形狀的幾何面積：

```go
package main

import (
	"fmt"
	"math"
)

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

type Square struct {
	Width  float64
	Height float64
}

func (s Square) Area() float64 {
	return s.Width * s.Height
}

type Sizer interface {
	Area() float64
}

func main() {
	c := Circle{Radius: 10}
	s := Square{Height: 10, Width: 5}

	l := Less(c, s)
	fmt.Printf("%+v is the smallest\n", l)
}

func Less(s1, s2 Sizer) Sizer {
	if s1.Area() < s2.Area() {
		return s1
	}
	return s2
}
```
因為每種類型都聲明了Area方法，所以我們可以建立一個定義該行為的介面。我們建立以下Sizer介面：

```go
type Sizer interface {
	Area() float64
}
```
然後我們定義一個名為 的函數Less，它接受兩個值Sizer並傳回最小的一個：

```go
func Less(s1, s2 Sizer) Sizer {
	if s1.Area() < s2.Area() {
		return s1
	}
	return s2
}
```

請注意，我們不僅接受兩個參數作為 type Sizer，而且還會傳回結果 a Sizer。這意味著我們不再返回 aSquare或 a Circle，而是返回 的介面Sizer。

最後，我們印出面積最小的東西：

```
{Width:5 Height:10} is the smallest
```
接下來，讓我們為每種類型新增另一個行為。這次我們將添加String()返回字串的方法。這將滿足fmt.Stringer接口：

```go
package main

import (
	"fmt"
	"math"
)

type Circle struct {
	Radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * math.Pow(c.Radius, 2)
}

func (c Circle) String() string {
	return fmt.Sprintf("Circle {Radius: %.2f}", c.Radius)
}

type Square struct {
	Width  float64
	Height float64
}

func (s Square) Area() float64 {
	return s.Width * s.Height
}

func (s Square) String() string {
	return fmt.Sprintf("Square {Width: %.2f, Height: %.2f}", s.Width, s.Height)
}

type Sizer interface {
	Area() float64
}

type Shaper interface {
	Sizer
	fmt.Stringer
}

func main() {
	c := Circle{Radius: 10}
	PrintArea(c)

	s := Square{Height: 10, Width: 5}
	PrintArea(s)

	l := Less(c, s)
	fmt.Printf("%v is the smallest\n", l)

}

func Less(s1, s2 Sizer) Sizer {
	if s1.Area() < s2.Area() {
		return s1
	}
	return s2
}

func PrintArea(s Shaper) {
	fmt.Printf("area of %s is %.2f\n", s.String(), s.Area())
}
```
因為Circle和Square類型都實作了Area和String方法，所以我們現在可以建立另一個介面來描述更廣泛的行為集。為此，我們將建立一個名為 的介面Shaper。我們將把這個Sizer介面和fmt.Stringer介面組合起來：

```go

type Shaper interface {
	Sizer
	fmt.Stringer
}
```
注意： 嘗試以 結尾來命名您的介面被認為是慣用的做法，er例如fmt.Stringer,等io.Writer。ShaperShape

現在我們可以建立一個名為 的函數PrintArea，它以 aShaper作為參數。這意味著我們可以對Area和方法傳入的值呼叫這兩個方法String：

```go

func PrintArea(s Shaper) {
	fmt.Printf("area of %s is %.2f\n", s.String(), s.Area())
}
```
如果我們運行該程序，我們將收到以下輸出：

```
$ go run interface/interface.go 
area of Circle {Radius: 10.00} is 314.16
area of Square {Width: 5.00, Height: 10.00} is 50.00
Square {Width: 5.00, Height: 10.00} is the smallest 
```
我們現在已經了解瞭如何創建較小的介面並根據需要將它們建構成更大的介面。雖然我們可以從較大的介面開始並將其傳遞給所有函數，但僅將最小的介面發送給所需的函數被認為是最佳實踐。這通常會產生更清晰的程式碼，因為任何接受特定較小介面的東西都只會與該定義的行為一起工作。

例如，如果我們傳遞Shaper給Less函數，我們可以假設它將呼叫Area和String方法。但是，由於我們只想呼叫該Area方法，因此它使該Less函數變得清晰，因為我們知道我們只能呼叫Area傳遞給它的任何參數的方法。

---
## 結論

我們已經看到如何創建較小的接口並將其構建為較大的接口，使我們能夠僅共享函數或方法所需的內容。我們還了解到，我們可以從其他接口組合我們的接口，包括從其他包定義的接口，而不僅僅是我們的package