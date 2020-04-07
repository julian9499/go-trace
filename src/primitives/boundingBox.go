package primitives

type BoundingBox struct {
	MinV, MaxV Vector
}

func (b *BoundingBox) Hit(r *Ray, tMin, tMax float64) (bool, Hit) {

	rec := Hit{}

	txmin := (b.MinV.X - r.Origin.X) / r.Direction.X
	txmax := (b.MaxV.X - r.Origin.X) / r.Direction.X

	if txmin > txmax {
		temp := txmin
		txmin = txmax
		txmax = temp
	}

	tymin := (b.MinV.Y - r.Origin.Y) / r.Direction.Y
	tymax := (b.MaxV.Y - r.Origin.Y) / r.Direction.Y

	if tymin > tymax {
		temp := tymin
		tymin = tymax
		tymax = temp
	}

	if (txmin > tymax) || (tymin > txmax) {
		return false, rec
	}
	if tymin > txmin {
		txmin = tymin
	}

	if tymax < txmax {
		txmax = tymax
	}

	tzmin := (b.MinV.Z - r.Origin.Z) / r.Direction.Z
	tzmax := (b.MaxV.Z - r.Origin.Z) / r.Direction.Z

	if tzmin > tzmax {
		temp := tzmin
		tzmin = txmax
		txmax = temp
	}

	if (txmin > tzmax) || (tzmin > txmax) {
		return false, rec
	}

	if tzmin > txmin {
		txmin = tzmin
	}

	if tzmax < txmax {
		txmax = tzmax
	}

	return true, rec
}
