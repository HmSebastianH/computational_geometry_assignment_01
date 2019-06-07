package geometry

import "fmt"

type Point struct {
	X float64
	Y float64
}

func (a *Point) equals( b Point) bool{
	return (a.X == b.X) && (a.Y == b.Y)
}

func (a *Point) String() string {
	return fmt.Sprintf("(%.2f, %.2f)", a.X, a.Y)
}