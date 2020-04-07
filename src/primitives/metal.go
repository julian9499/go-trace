package primitives

type Metal struct {
	C    Vector
	Fuzz float64
}

func (m Metal) Bounce(input Ray, hit Hit) (bool, Ray) {
	direction := input.Direction.Reflect(hit.Normal)
	bouncedRay := Ray{hit.P, direction.Add(VectorInUnitSphere().MultiplyScalar(m.Fuzz))}
	bounced := direction.Dot(hit.Normal) > 0
	return bounced, bouncedRay
}

func (m Metal) Color() Vector {
	return m.C
}