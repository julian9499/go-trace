package primitives

import (
	"math"
	"math/rand"
)

type Glass struct {
	C        Vector
	RefIndex float64
}

func (m Glass) Color() Vector {
	return m.C
}

func (d Glass) Bounce(input Ray, hit Hit) (bool, Ray) {
	var outwardNormal Vector
	var niOverNt, cosine float64

	if input.Direction.Dot(hit.Normal) > 0 {
		outwardNormal = hit.Normal.MultiplyScalar(-1)
		niOverNt = d.RefIndex

		a := input.Direction.Dot(hit.Normal) * d.RefIndex
		b := input.Direction.Length()

		cosine = a / b
	} else {
		outwardNormal = hit.Normal
		niOverNt = 1.0 / d.RefIndex

		a := input.Direction.Dot(hit.Normal) * d.RefIndex
		b := input.Direction.Length()

		cosine = -a / b
	}

	var success bool
	var refracted Vector
	var reflectProbability float64

	if success, refracted = input.Direction.Refract(outwardNormal, niOverNt); success {
		reflectProbability = d.schlick(cosine)
	} else {
		reflectProbability = 1.0
	}

	if rand.Float64() < reflectProbability {
		reflected := input.Direction.Reflect(hit.Normal)
		return true, Ray{hit.P, reflected}
	}

	return true, Ray{hit.P, refracted}
}

func (d Glass) schlick(cosine float64) float64 {
	r0 := (1 - d.RefIndex) / (1 + d.RefIndex)
	r0 = r0 * r0
	return r0 + (1-r0)*math.Pow((1-cosine), 5)
}
