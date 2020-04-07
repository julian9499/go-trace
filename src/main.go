package main

import (
	"fmt"
	"image"
	color2 "image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	p "./primitives"
)

const (
	nx = 800 // size of x
	ny = 400 // size of y
	ns = 50   // number of AA sampling
	c  = 255.99
)

var (
	white    = p.Vector{1.0, 1.0, 1.0}
	blue     = p.Vector{0.5, 0.7, 1.0}
	upLeft   = image.Point{0, 0}
	lowRight = image.Point{nx, ny}
)

func check(e error, s string) {
	if e != nil {
		fmt.Fprintf(os.Stderr, s, e)
		os.Exit(1)
	}
}

func color(r p.Ray, world p.Hitable, depth int) p.Vector {
	hit, record := world.Hit(r, 0.001, math.MaxFloat64)

	if hit {
		if depth < 50 {
			bounced, bouncedRay := record.Bounce(r, record)
			if bounced {
				newColor := color(bouncedRay, world, depth+1)
				return record.Material.Color().Multiply(newColor)
			}
		}
		return p.Vector{}
	}

	return gradient(r)
}

func gradient(r p.Ray) p.Vector {
	// make unit vector so y is between -1.0 and 1.0
	v := r.Direction.Normalize()

	// scale t to be between 0.0 and 1.0
	t := 0.5 * (v.Y + 1.0)

	// linear blend: blended_value = (1 - t) * white + t * blue
	return white.MultiplyScalar(1.0 - t).Add(blue.MultiplyScalar(t))
}

func createFile() *os.File {
	f, err := os.Create("outNew.png")
	check(err, "Error opening file: %v\n")

	// http://netpbm.sourceforge.net/doc/ppm.html
	return f
}

func setImg(img *image.RGBA, rgb p.Vector, x, y int) {
	// get intensity of colors
	ir := uint8(c * math.Sqrt(rgb.X))
	ig := uint8(c * math.Sqrt(rgb.Y))
	ib := uint8(c * math.Sqrt(rgb.Z))

	color := color2.RGBA{ir, ig, ib, 0xff}

	img.Set(x, y, color)
}

// samples rays for anti-aliasing
func sample(world *p.World, camera *p.Camera, i, j int, rand *rand.Rand) p.Vector {
	rgb := p.Vector{}

	for s := 0; s < ns; s++ {
		u := (float64(i) + rand.Float64()) / float64(nx)
		v := (float64(j) + rand.Float64()) / float64(ny)

		ray := camera.RayAt(u, v, rand)
		rgb = rgb.Add(color(ray, world, 0))
	}

	// average
	return rgb.DivideScalar(float64(ns))
}

func render(world *p.World, camera *p.Camera) {
	ticker := time.NewTicker(time.Millisecond * 500)

	go func() {
		for {
			<-ticker.C
			fmt.Print(".")
		}
	}()

	f := createFile()
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})
	defer f.Close()

	start := time.Now()

	//type empty struct {}
	var wg sync.WaitGroup

	res := make([][]p.Vector, ny)

	for j := 0; j < ny; j++ {
		res[j] = make([]p.Vector, nx)
		wg.Add(1)
		go func(res [][]p.Vector, world *p.World, camera *p.Camera, j int, wg *sync.WaitGroup) {
			defer wg.Done()
			rand := rand.New(rand.NewSource(time.Now().UnixNano()))
			for i := 0; i < nx; i++ {
				res[j][i] = sample(world, camera, i, j, rand)
			}
		}(res, world, camera, j, &wg)
	}
	wg.Wait()

	for j := ny - 1; j >= 0; j-- {
		for i := 0; i < nx; i++ {
			setImg(img, res[j][i], i, ny-j)
		}
	}
	png.Encode(f, img)

	ticker.Stop()
	fmt.Printf("\nDone.\nElapsed: %v\n", time.Since(start))
}

func main() {
	camera := p.NewCamera(p.Vector{0, 0, -1}, p.Vector{0, 0, 1}, 75.0, 2.0, 0.01)

	world := p.World{}

	sphere := p.Sphere{p.Vector{0, 0, -1}, 0.5, p.Lambertian{p.Vector{0.8, 0.3, 0.3}}}
	floor := p.Sphere{p.Vector{0, -100000.5, -1}, 100000, p.Lambertian{p.Vector{0.4, 0.4, 0.4}}}
	left := p.Sphere{p.Vector{-1, 0, -1}, 0.5, p.Glass{p.Vector{0.8, 0.8, 0.8}, 1.5}}
	right := p.Sphere{p.Vector{1, 0, -1}, 0.5, p.Metal{p.Vector{0.8, 0.6, 0.2}, 0.3}}

	triangle := p.Triangle{p.Vector{0,0,3} , p.Vector{0,0.5,3}, p.Vector{0.5,0.5,3}, p.Lambertian{p.Vector{0.8, 0.3, 0.3}}}

	rand := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 800; i++ {
		//tempSphere := p.Sphere{Center: p.Vector{rand.Float64()*80 - 60, -0.25, rand.Float64()*80 - 60}, Radius: 0.25, Material: p.Lambertian{p.Vector{rand.Float64(), rand.Float64(), rand.Float64()}}}
		tempSphere := p.Sphere{Center: p.Vector{rand.Float64()*2.0, rand.Float64()*2.0 -1.0, rand.Float64()*2.0 -1.0}, Radius: 0.000005, Material: p.Lambertian{p.Vector{rand.Float64(), rand.Float64(), rand.Float64()}}}
		world.Add(&tempSphere)
	}

	world.Add(&sphere)
	world.Add(&floor)
	world.Add(&left)
	world.Add(&right)
	world.Add(&triangle)

	render(&world, &camera)
}
