/*
Normalize Method:
The Normalize method is a receiver method for the Vector type. This method normalizes the vector, meaning it converts the vector into a unit vector (a vector of length 1) that points in the same direction as the original vector.
The normalization process involves:
Calculating the Magnitude: The magnitude (or length) of the vector is calculated using the Pythagorean theorem. For a vector (X, Y), the magnitude is sqrt(X^2 + Y^2).
Dividing Each Component by the Magnitude: The X and Y components of the vector are each divided by the magnitude. This scales the vector so that its length is 1.
The result is a new Vector where X and Y are the normalized components.
*/
package game

import "math"

type Vector struct {
	X float64
	Y float64
}

func (v Vector) Normalize() Vector {
	magnitude := math.Sqrt(v.X*v.X + v.Y*v.Y)
	return Vector{v.X / magnitude, v.Y / magnitude}
}
