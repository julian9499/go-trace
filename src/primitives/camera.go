package primitives

import (
	"math"
	"math/rand"
)

var vUp = Vector{X: 0, Y: 1, Z: 0}

type Camera struct {
	lowerLeft, horizontal, vertical, origin, w, u, v Vector
	lensRadius                                       float64
}

func NewCamera(LookFrom, LookAt Vector, vFov, aspect, aperture float64) Camera {
	c := Camera{}

	c.origin = LookFrom
	c.lensRadius = aperture / 2

	theta := vFov * math.Pi / 180
	halfHeight := math.Tan(theta / 2)
	halfWidth := aspect * halfHeight

	c.w = LookFrom.Subtract(LookAt).Normalize()
	c.u = vUp.Cross(c.w).Normalize()
	c.v = c.w.Cross(c.u)

	focusDist := LookFrom.Subtract(LookAt).Length()

	x := c.u.MultiplyScalar(halfWidth * focusDist)
	y := c.v.MultiplyScalar(halfHeight * focusDist)

	c.lowerLeft = c.origin.Subtract(x).Subtract(y).Subtract(c.w.MultiplyScalar(focusDist))
	c.horizontal = x.MultiplyScalar(2)
	c.vertical = y.MultiplyScalar(2)

	return c
}

func (c *Camera) RayAt(s, t float64, rnd *rand.Rand) Ray {
	var randomInUnitDisc Vector
	for {
		randomInUnitDisc = Vector{rnd.Float64(), rnd.Float64(), 0}.MultiplyScalar(2).Subtract(Vector{1, 1, 0})
		if randomInUnitDisc.Dot(randomInUnitDisc) < 1.0 {
			break
		}
	}

	rd := randomInUnitDisc.MultiplyScalar(c.lensRadius)
	offset := c.u.MultiplyScalar(rd.X).Add(c.v.MultiplyScalar(rd.Y))

	horizontal := c.horizontal.MultiplyScalar(s)
	vertical := c.vertical.MultiplyScalar(t)

	origin := c.origin.Add(offset)
	direction := c.lowerLeft.Add(horizontal).Add(vertical).Subtract(c.origin).Subtract(offset)
	return Ray{origin, direction}
}
