package main

import (
	"image/color"
	"log"
	"math"

	"github.com/fogleman/gg"
)

func drawStar(dc *gg.Context, x, y, radius, g float64) {
	// 绘制星星
	for i := 0; i < 5; i++ {
		angle := float64(i)*144.0*(math.Pi/180.0) + g
		x1 := x + radius*math.Cos(angle)
		y1 := y + radius*math.Sin(angle)
		if i == 0 {
			dc.MoveTo(x1, y1)
		} else {
			dc.LineTo(x1, y1)
		}
	}
	dc.ClosePath()
	dc.Fill()
}

func drawChinaFlag() {
	const W = 600
	const H = 400

	dc := gg.NewContext(W, H)

	// 绘制背景
	dc.SetColor(color.NRGBA{R: 255, G: 0, B: 0, A: 255}) // 红色
	dc.Clear()

	// 绘制大星星
	dc.SetColor(color.NRGBA{R: 255, G: 215, B: 0, A: 255}) // 黄色
	dc.Push()
	dc.Translate(100, 100) // 大星星位置
	drawStar(dc, 0, 0, 60, +270*(math.Pi/180))
	dc.Pop()

	// 绘制四颗小星星
	dc.SetColor(color.NRGBA{R: 255, G: 215, B: 0, A: 255}) // 黄色

	// 左1
	dc.Push()
	dc.Translate(200, 40) // 小星星位置
	drawStar(dc, 0, 0, 20, 294*(math.Pi/180))
	dc.Pop()

	// 右1
	dc.Push()
	dc.Translate(240, 80) // 小星星位置
	drawStar(dc, 0, 0, 20, 318*(math.Pi/180))
	dc.Pop()

	// 右2
	dc.Push()
	dc.Translate(240, 140) // 小星星位置
	drawStar(dc, 0, 0, 20, 342*(math.Pi/180))
	dc.Pop()

	// 左2
	dc.Push()
	dc.Translate(200, 180) // 小星星位置
	drawStar(dc, 0, 0, 20, 294*(math.Pi/180))
	dc.Pop()

	// 保存为PNG文件
	if err := dc.SavePNG("china_flag.png"); err != nil {
		log.Fatal(err)
	}
}

func main() {
	drawChinaFlag()
}
