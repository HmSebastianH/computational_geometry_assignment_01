package geometry


type Point struct {
	X float64
	Y float64
}

func (a *Point) equals( b Point) bool{
	return (a.X == b.X) && (a.Y == b.Y)
}
