package main

import (
    "os"
    "log"
    "math"
    "math/cmplx"
    "image/color"
)

func main() {
    g := NewGraph(NewArea(-5, 5, 5, -5), 1024, 1024)

    g.DrawGrid()

    g.DrawBoolExpression(Mandelbrot)

    g.DrawDiffExpressionWithColor(func (c *Coord) float64 {
        return (math.Pow(c.X, 2) - math.Pow(c.Y, 2)) - 1
    }, color.RGBA{0x80, 0xFF, 0x00, 0xFF})

    g.DrawDiffExpressionWithColor(func (c *Coord) float64 {
        return c.Y - math.Cosh(c.X)
    }, color.RGBA{0x80, 0x00, 0x80, 0xFF})
    g.DrawDiffExpressionWithColor(func (c *Coord) float64 {
        return c.Y - math.Sinh(c.X)
    }, color.RGBA{0x00, 0x80, 0x00, 0xFF})
    g.DrawDiffExpressionWithColor(func (c *Coord) float64 {
        return c.Y - math.Tanh(c.X)
    }, color.RGBA{0xFF, 0x80, 0x00, 0xFF})

    g.DrawDiffExpressionWithColor(func (c *Coord) float64 {
        return c.Y - math.Pow(math.E, c.X)
    }, color.RGBA{0x00, 0x80, 0xFF, 0xFF})
    g.DrawDiffExpressionWithColor(func (c *Coord) float64 {
        return c.Y - math.Pow(math.E, -c.X)
    }, color.RGBA{0x00, 0xFF, 0x80, 0xFF})

    g.ApplyComplexExpression(func (z complex128) complex128 {
        return cmplx.Pow(z, complex(2, 0)) / 5
    })

    f, err := os.Create("out.png")
    if err != nil {
        log.Fatal(err)
    }

    defer f.Close()

    err = g.SavePNG(f)
    if err != nil {
        log.Fatal(err)
    }
}