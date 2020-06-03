package main

import (
    "math"
    "image"
    "image/color"
    "image/png"
    "io"
)

const (
    ChunkSize = 64
)

var (
    AxisColor       = color.RGBA{0xFF, 0x00, 0x00, 0xFF}
    GridColor       = color.RGBA{0xE0, 0xE0, 0xE0, 0xFF}
    ExpressionColor = color.RGBA{0x00, 0x00, 0x00, 0xFF}
)

type Coord struct {
    X, Y float64
}

type Expression func (c *Coord) interface{}

type ComplexExpression func (z complex128) complex128

type Area struct {
    Pos0, Pos1 *Coord
}

type Graph struct {
    Bounds *Area
    Image *image.RGBA
}

func NewCoord(x, y float64) *Coord {
    return &Coord{x, y}
}

func NewCoordFromPolar(r, theta float64) *Coord {
    return NewCoord(r * math.Cos(theta), r * math.Sin(theta))
}

func (c *Coord) Equals(other *Coord) bool {
    return c.X == other.X && c.Y == other.Y
}

func (c *Coord) Add(other *Coord) *Coord {
    return NewCoord(c.X + other.X, c.Y + other.Y)
}

func (c *Coord) Sub(other *Coord) *Coord {
    return NewCoord(c.X - other.X, c.Y - other.Y)
}

func (c *Coord) Mult(mult float64) *Coord {
    return NewCoord(mult * c.X, mult * c.Y)
}

func (c *Coord) Div(div float64) *Coord {
    return NewCoord(c.X / div, c.Y / div)
}

func (c *Coord) Dist(other *Coord) float64 {
    return math.Sqrt(math.Pow(c.X - other.X, 2) + math.Pow(c.Y - other.Y, 2))
}

func (c *Coord) DistOrigin() float64 {
    return c.Dist(NewCoord(0, 0))
}

func (c *Coord) WithinDist(other *Coord, dist float64) bool {
    return c.Dist(other) <= dist
}

func (c *Coord) Polar() (r, theta float64) {
    r = c.DistOrigin()
    theta = math.Atan2(c.Y, c.X)

    if theta < 0 {
        theta += 2 * math.Pi
    }

    return
}

func (c *Coord) Rotate(theta float64) *Coord {
    r, t := c.Polar()

    return NewCoordFromPolar(r, t + theta)
}

func (c *Coord) RotateAround(theta float64, other *Coord) *Coord {
    return c.Sub(other).Rotate(theta).Add(other)
}

func NewArea(x0, y0, x1, y1 float64) *Area {
    return &Area{NewCoord(x0, y0), NewCoord(x1, y1)}
}

func (a *Area) Width() float64 {
    return math.Abs(a.Pos1.X - a.Pos0.X)
}

func (a *Area) Height() float64 {
    return math.Abs(a.Pos1.Y - a.Pos0.Y)
}

func (a *Area) Size() *Coord {
    return NewCoord(a.Width(), a.Height())
}

func (a *Area) CenterX() float64 {
    return (a.Pos0.X + a.Pos1.X) / 2
}

func (a *Area) CenterY() float64 {
    return (a.Pos0.Y + a.Pos1.Y) / 2
}

func (a *Area) Center() *Coord {
    return NewCoord(a.CenterX(), a.CenterY())
}

func (a *Area) Contains(c *Coord) bool {
    return a.Pos0.X <= c.X && c.X < a.Pos1.X && a.Pos0.Y >= c.Y && c.Y > a.Pos1.Y
}

func NewGraph(bounds *Area, scale float64) *Graph {
    g := Graph{}

    g.Bounds = bounds
    g.Image = image.NewRGBA(image.Rect(0, 0, int(bounds.Width() * scale), int(bounds.Height() * scale)))

    for i, _ := range g.Image.Pix {
        g.Image.Pix[i] = 0xFF
    }

    return &g
}

func (g *Graph) SavePNG(w io.Writer) error {
    return png.Encode(w, g.Image)
}

func (g *Graph) ImageWidth() int {
    return g.Image.Bounds().Dx()
}

func (g *Graph) ImageHeight() int {
    return g.Image.Bounds().Dy()
}

func (g *Graph) CoordToPixel(c *Coord) image.Point {
    tmp_c := c.Sub(g.Bounds.Pos0)
    tmp_c.X *= float64(g.ImageWidth()) / g.Bounds.Width()
    tmp_c.Y *= -float64(g.ImageHeight()) / g.Bounds.Height()

    return image.Pt(int(tmp_c.X), int(tmp_c.Y))
}

func (g *Graph) PixelToCoord(pt image.Point) *Coord {
    c := NewCoord(float64(pt.X), float64(pt.Y))
    c.X *= g.Bounds.Width() / float64(g.ImageWidth())
    c.Y *= -g.Bounds.Height() / float64(g.ImageHeight())
    c = c.Add(g.Bounds.Pos0)

    return c
}

func (g *Graph) SetPixel(pt image.Point, col color.Color) {
    g.Image.Set(pt.X, pt.Y, BlendColor(g.Image.At(pt.X, pt.Y), col))
}

func (g *Graph) SetCoord(c *Coord, col color.Color) {
    pt := g.CoordToPixel(c)
    g.SetPixel(pt, col)
}

func (g *Graph) AtPixel(pt image.Point) color.Color {
    return g.Image.At(pt.X, pt.Y)
}

func (g *Graph) AtCoord(c *Coord) color.Color {
    pt := g.CoordToPixel(c)
    return g.AtPixel(pt)
}

func (g *Graph) DrawLine(c0, c1 *Coord, col color.Color) {
    var p0, p1 image.Point

    if (c0.X <= c1.X) {
        p0 = g.CoordToPixel(c0)
        p1 = g.CoordToPixel(c1)
    } else {
        p0 = g.CoordToPixel(c1)
        p1 = g.CoordToPixel(c0)
    }

    delta := p1.Sub(p0)

    if (delta.X == 0) { // Vertical line
        var diff int
        if (p0.Y <= p1.Y) {
            diff = 1
        } else {
            diff = -1
        }

        for ; p0.Y < p1.Y; p0.Y += diff {
            g.Image.Set(p0.X, p0.Y, col)
        }
    } else {
        slope := float64(delta.Y) / float64(delta.X)
        yint := int(float64(p0.Y) - slope * float64(p0.X))

        for x := p0.X; x <= p1.X; x++ {
            y := int(slope * float64(x)) + yint

            g.SetPixel(image.Pt(x, y), col)
        }
    }
}

func (g *Graph) DrawAxes() {
    g.DrawLine(NewCoord(g.Bounds.CenterX(), g.Bounds.Pos0.Y), NewCoord(g.Bounds.CenterX(), g.Bounds.Pos1.Y), AxisColor)
    g.DrawLine(NewCoord(g.Bounds.Pos0.X, g.Bounds.CenterY()), NewCoord(g.Bounds.Pos1.X, g.Bounds.CenterY()), AxisColor)
}

func (g *Graph) DrawGrid() {
    for x := 1.0; x < g.Bounds.Pos1.X; x++ {
        g.DrawLine(NewCoord(x, g.Bounds.Pos0.Y), NewCoord(x, g.Bounds.Pos1.Y), GridColor)
        g.DrawLine(NewCoord(-x, g.Bounds.Pos0.Y), NewCoord(-x, g.Bounds.Pos1.Y), GridColor)

        for y := 1.0; y < g.Bounds.Pos0.Y; y++ {
            g.DrawLine(NewCoord(g.Bounds.Pos0.X, y), NewCoord(g.Bounds.Pos1.X, y), GridColor)
            g.DrawLine(NewCoord(g.Bounds.Pos0.X, -y), NewCoord(g.Bounds.Pos1.X, -y), GridColor)
        }
    }

    g.DrawAxes()
}

func (g *Graph) DrawExpressionInChunk(expr Expression, r *image.Rectangle, col color.Color, ch chan struct{}) {
    for x := r.Min.X; x < r.Max.X; x++ {
        for y := r.Min.Y; y < r.Max.Y; y++ {
            pt := image.Pt(x, y)
            c := g.PixelToCoord(pt)

            switch ret := expr(c); ret.(type) {
                case bool:
                    if ret.(bool) {
                        g.SetPixel(pt, col)
                    }

                case float64:
                    coords := [3]*Coord {
                        g.PixelToCoord(pt.Add(image.Pt(1, 0))),
                        g.PixelToCoord(pt.Add(image.Pt(0, 1))),
                        g.PixelToCoord(pt.Add(image.Pt(1, 1))),
                    }

                    diff := [3]float64 {
                        expr(coords[0]).(float64),
                        expr(coords[1]).(float64),
                        expr(coords[2]).(float64),
                    }

                    if ret.(float64) == 0 {
                        g.SetPixel(pt, col)
                        continue
                    }

                    for _, d := range diff {
                        if (ret.(float64) > 0 && d < 0) || (ret.(float64) < 0 && d > 0) {
                            g.SetPixel(pt, col)
                            continue
                        }
                    }
            }
        }
    }

    ch <- struct{}{}
}

func (g *Graph) DrawExpressionWithColor(expr Expression, col color.Color) {
    var channels []chan struct{}

    for x := 0; x < g.ImageWidth(); x += ChunkSize {
        for y := 0; y < g.ImageHeight(); y += ChunkSize {
            ch := make(chan struct {})
            channels = append(channels, ch)

            r := image.Rect(x, y, MinInt(x + ChunkSize, g.ImageWidth()), MinInt(y + ChunkSize, g.ImageHeight()))
            go g.DrawExpressionInChunk(expr, &r, col, ch)
        }
    }

    for _, ch := range channels {
        <-ch
    }
}

func (g *Graph) DrawExpression(expr Expression) {
    g.DrawExpressionWithColor(expr, ExpressionColor)
}

func (g *Graph) ApplyComplexExpressionInChunk(expr ComplexExpression, dst *image.RGBA, r *image.Rectangle, ch chan struct{}) {
    for x := r.Min.X; x < r.Max.X; x++ {
        for y := r.Min.Y; y < r.Max.Y; y++ {
            pt := image.Pt(x, y)
            c := g.PixelToCoord(pt)
            z := complex(c.X, c.Y)

            new_z := expr(z)
            new_c := NewCoord(real(new_z), imag(new_z))

            if g.Bounds.Contains(new_c) {
                new_pt := g.CoordToPixel(new_c)
                col := g.AtPixel(pt)
                dst.Set(new_pt.X, new_pt.Y, col)
            }
        }
    }

    ch <- struct{}{}
}

func (g *Graph) ApplyComplexExpression(expr ComplexExpression) {
    img := image.NewRGBA(g.Image.Bounds())
    for i, _ := range g.Image.Pix {
        img.Pix[i] = 0xFF
    }

    var channels []chan struct{}

    for x := 0; x < g.ImageWidth(); x += ChunkSize {
        for y := 0; y < g.ImageHeight(); y += ChunkSize {
            ch := make(chan struct {})
            channels = append(channels, ch)

            r := image.Rect(x, y, MinInt(x + ChunkSize, g.ImageWidth()), MinInt(y + ChunkSize, g.ImageHeight()))
            go g.ApplyComplexExpressionInChunk(expr, img, &r, ch)
        }
    }

    for _, ch := range channels {
        <-ch
    }

    g.Image = img
}