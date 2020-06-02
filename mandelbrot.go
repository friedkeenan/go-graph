package main

import "math/cmplx"

const (
    MaxIterations = 200
)

func Mandelbrot(c *Coord) bool {
    seed := complex(c.X, c.Y)
    z := complex(0, 0)

    for i := 0; i < MaxIterations; i++ {
        z = cmplx.Pow(z, complex(2, 0)) + seed

        if cmplx.Abs(z) >= 2 {
            return false;
        }
    }

    return true
}