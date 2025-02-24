/*
Intersects Method:
The Intersects method is used to determine if two Rect objects overlap or intersect. It takes another Rect (other) as a parameter and returns true if there is any overlap. The method checks if:

The right edge of the first rectangle (r.MaxX()) is to the right of the left edge of the second rectangle (other.X).
The left edge of the first rectangle (r.X) is to the left of the right edge of the second rectangle (other.MaxX()).
The bottom edge of the first rectangle (r.MaxY()) is below the top edge of the second rectangle (other.Y).
The top edge of the first rectangle (r.Y) is above the bottom edge of the second rectangle (other.MaxY()).
If all these conditions are true, the rectangles intersect.

	|    |
	|    |
	| |~~|	| ----> intersect happen
	  |		|
	  | 	|
*/
package utils

type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

func NewRect(x, y, width, height float64) Rect {
	return Rect{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

func (r Rect) MaxX() float64 {
	return r.X + r.Width
}

func (r Rect) MaxY() float64 {
	return r.Y + r.Height
}

// Intersects checks if the current rectangle intersects with another rectangle.
// It returns true if the rectangles overlap, otherwise false.
//
// Parameters:
// - other: The rectangle to check for intersection with.
//
// Returns:
// - bool: True if the rectangles intersect, false otherwise.
func (r Rect) Intersects(other Rect) bool {
	return r.X <= other.MaxX() &&
		other.X <= r.MaxX() &&
		r.Y <= other.MaxY() &&
		other.Y <= r.MaxY()
}
