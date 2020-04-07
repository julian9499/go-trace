package primitives

import "math"

const kEpsilon = 0.00001

type Triangle struct{
	V0,V1,V2 Vector
	Material
}

func (tr *Triangle) Hit(r Ray, tMin float64, tMax float64) (bool, Hit) {
	rec := Hit{}

	v0v1 := tr.V1.Subtract(tr.V0)
	v0v2 := tr.V2.Subtract(tr.V0)

	N := v0v1.Cross(v0v2).Normalize()

	NdotRayDirection := N.Dot(r.Direction)

	if math.Abs(NdotRayDirection) < kEpsilon {
		return false ,rec
	}

	d := N.Dot(tr.V0)
	t := (N.Dot(r.Origin) + d) / NdotRayDirection

	if t < 0 || t < tMin || t > tMax {
		return false, rec
	}

	P := r.Origin.Add(r.Direction.MultiplyScalar(t))
	vp0 := P.Subtract(tr.V0)
	C := v0v1.Cross(vp0)

	if N.Dot(C) < 0 {
		return false,rec
	}

	v2v1 := tr.V2.Subtract(tr.V1)
	vp1 := P.Subtract(tr.V1)
	C = v2v1.Cross(vp1)

	if N.Dot(C) < 0 {
		return false,rec
	}

	v2v0 := tr.V0.Subtract(tr.V2)
	vp2 := P.Subtract(tr.V2)
	C = v2v0.Cross(vp2)

	if N.Dot(C) < 0 {
		return false,rec
	}

	rec.P = P
	rec.Normal = N
	rec.T = t
	rec.Material = tr.Material

	return true, rec
}
