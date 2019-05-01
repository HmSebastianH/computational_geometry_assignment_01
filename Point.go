package main


type Point struct {
	x float64
	y float64
}

func (a *Point) equals( b Point) bool{
	return (a.x == b.x) && (a.y == b.y)
}
