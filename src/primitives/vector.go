package primitives

import (
	"math"
	"math/rand"
)

var (
	UnitVector = Vector{1,1,1}
)

type Vector struct {
	X, Y, Z float64
}

func VectorInUnitSphere() Vector {
	for {
		r := Vector{rand.Float64(), rand.Float64(), rand.Float64()}
		p := r.MultiplyScalar(2.0).Subtract(UnitVector)
		if p.SquaredLength() >= 1.0 {
			return p
		}
	}
}

func (v Vector) SquaredLength() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vector) Length() float64 {
	return math.Sqrt(v.Dot(v))
}

func (v Vector) Dot(o Vector) float64 {
	return v.X*o.X + v.Y*o.Y + v.Z*o.Z
}

func (v Vector) Normalize() Vector {
	length := v.Length()
	return Vector{v.X / length, v.Y / length, v.Z / length}
}

func (v Vector) Cross(o Vector) Vector {
	crossX := v.Y*o.Z - v.Z*v.Y
	crossY := -(v.X*o.Z - v.Z*o.X)
	crossZ := v.X*o.Y - o.X*v.Y
	return Vector{crossX, crossY, crossZ}
}

func (v Vector) MultiplyScalar(t float64) Vector {
	return Vector{v.X * t, v.Y * t, v.Z * t}
}

func (v Vector) Multiply(o Vector) Vector {
	return Vector{v.X * o.X, v.Y * o.Y, v.Z * o.Z}
}

func (v Vector) DivideScalar(t float64) Vector {
	return Vector{v.X / t, v.Y / t, v.Z / t}
}

func (v Vector) Add(o Vector) Vector {
	return Vector{v.X + o.X, v.Y + o.Y, v.Z + o.Z}
}

func (v Vector) AddScalar(t float64) Vector {
	return Vector{v.X + t, v.Y + t, v.Z + t}
}

func (v Vector) Subtract(o Vector) Vector {
	return Vector{v.X - o.X, v.Y - o.Y, v.Z - o.Z}
}

func (v Vector) Refract(o Vector, n float64) (bool, Vector) {
	uv := v.Normalize()
	uo := o.Normalize()
	dt := uv.Dot(uo)
	discriminant := 1.0 - (n * n * (1 - dt*dt))
	if discriminant > 0 {
		a := uv.Subtract(o.MultiplyScalar(dt)).MultiplyScalar(n)
		b := o.MultiplyScalar(math.Sqrt(discriminant))
		return true, a.Subtract(b)
	}
	return false, Vector{}
}

func(v Vector) Reflect(n Vector) Vector {
	b := 2 * v.Dot(n)
	return v.Subtract(n.MultiplyScalar(b))
}