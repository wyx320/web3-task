package exam5

import "math"

/*
题目五：
定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。
然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
考察点：
接口的定义与实现、面向对象编程风格。
*/

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectange struct {
	Length float64
	Width  float64
}
type Circle struct {
	Radius float64
}

func (rec *Rectange) Area() float64 {
	return rec.Length * rec.Width
}
func (rec *Rectange) Perimeter() float64 {
	return (rec.Length + rec.Width) * 2
}
func (cir *Circle) Area() float64 {
	return math.Pi * cir.Radius * cir.Radius
}
func (cir *Circle) Perimeter() float64 {
	return 2 * math.Pi * cir.Radius
}
