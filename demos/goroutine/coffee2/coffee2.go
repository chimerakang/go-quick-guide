package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/trace"
	"time"
)

type (
	Bean       int
	GroundBean int
	Water      int
	HotWater   int
	Coffee     int
)

const (
	GramBeans          Bean       = 1
	GramGroundBeans    GroundBean = 1
	MilliLiterWater    Water      = 1
	MilliLiterHotWater HotWater   = 1
	CupsCoffee         Coffee     = 1
)

func (w Water) String() string {
	return fmt.Sprintf("%d[ml] water", int(w))
}

func (hw HotWater) String() string {
	return fmt.Sprintf("%d[ml] hot water", int(hw))
}

func (b Bean) String() string {
	return fmt.Sprintf("%d[g] beans", int(b))
}

func (gb GroundBean) String() string {
	return fmt.Sprintf("%d[g] ground beans", int(gb))
}

func (cups Coffee) String() string {
	return fmt.Sprintf("%d cup(s) coffee", int(cups))
}

// 沖泡 1 杯咖啡所需的水量
func (cups Coffee) Water() Water {
	return Water(180*cups) / MilliLiterWater
}

// 沖泡 1 杯咖啡所需的熱水量
func (cups Coffee) HotWater() HotWater {
	return HotWater(180*cups) / MilliLiterHotWater
}

// 沖泡 1 杯咖啡所需的咖啡豆量
func (cups Coffee) Beans() Bean {
	return Bean(20*cups) / GramBeans
}

// 沖泡 1 杯咖啡所需咖啡粉量
func (cups Coffee) GroundBeans() GroundBean {
	return GroundBean(20*cups) / GramGroundBeans
}

func main() {
	f, err := os.Create("trace.out")
	if err != nil {
		log.Fatalln("Error:", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Fatalln("Error:", err)
		}
	}()

	if err := trace.Start(f); err != nil {
		log.Fatalln("Error:", err)
	}
	defer trace.Stop()

	RunMyProgram()
}

func RunMyProgram() {
	// 預計沖泡多少杯咖啡
	const amountCoffee = 20 * CupsCoffee

	ctx, task := trace.NewTask(context.Background(), "make coffee")
	defer task.End()

	// 材料
	water := amountCoffee.Water()
	beans := amountCoffee.Beans()

	fmt.Println(water)
	fmt.Println(beans)

	// 熱水
	var hotWater HotWater
	for water > 0 {
		water -= 600 * MilliLiterWater
		hotWater += boil(ctx, 600*MilliLiterWater)
	}
	fmt.Println(hotWater)

	// 咖啡粉
	var groundBeans GroundBean
	for beans > 0 {
		beans -= 20 * GramBeans
		groundBeans += grind(ctx, 20*GramBeans)
	}
	fmt.Println(groundBeans)

	// 沖泡咖啡
	var coffee Coffee
	cups := 4 * CupsCoffee
	for hotWater >= cups.HotWater() && groundBeans >= cups.GroundBeans() {
		hotWater -= cups.HotWater()
		groundBeans -= cups.GroundBeans()
		coffee += brew(ctx, cups.HotWater(), cups.GroundBeans())
	}

	fmt.Println(coffee)
}

// 燒開水
func boil(ctx context.Context, water Water) HotWater {
	defer trace.StartRegion(ctx, "boil").End()
	time.Sleep(400 * time.Millisecond)
	return HotWater(water)
}

// 研磨
func grind(ctx context.Context, beans Bean) GroundBean {
	defer trace.StartRegion(ctx, "grind").End()
	time.Sleep(200 * time.Millisecond)
	return GroundBean(beans)
}

// 沖泡
func brew(ctx context.Context, hotWater HotWater, groundBeans GroundBean) Coffee {
	defer trace.StartRegion(ctx, "brew").End()
	time.Sleep(1 * time.Second)
	// 少量者優先處理
	cups1 := Coffee(hotWater / (1 * CupsCoffee).HotWater())
	cups2 := Coffee(groundBeans / (1 * CupsCoffee).GroundBeans())
	if cups1 < cups2 {
		return cups1
	}
	return cups2
}
