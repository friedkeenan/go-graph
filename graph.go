package main

import (
    "math"
    "image"
    "image/color"
    "image/png"
    "io"
)

const (
    /*
        The maximum chunk size that goroutines
        will use for drawing relations and functions
    */
    ChunkSize = 64

    /*
        The maximum angle size that goroutines
        will use for drawing polar functions
    */
    AngleSize = math.Pi / 4

    /*
        The angle step that goroutines will increment
        the angle by when drawing polar functions
    */
    AngleStep = AngleSize / 100
)

var (
    /* The default graph background color */
    DefaultBackgroundColor = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}

    /* The default graph axis color */
    DefaultAxisColor       = color.RGBA{0xFF, 0x00, 0x00, 0xFF}

    /* The default graph grid color */
    DefaultGridColor       = color.RGBA{0xE0, 0xE0, 0xE0, 0xFF}

    /* The default graph relation color */
    DefaultRelationColor   = color.RGBA{0x00, 0x00, 0x00, 0xFF}
)

/* A coordinate on a graph */
type Coord struct {
    X, Y float64
}

/*
    A function that takes in a coordinate
    and outputs one of the following types:
    
    A bool to indicate that the coordinate should filled in.
    This is what you want for relations that are of a similar
    form as "f(x, y) <= g(x, y)".

    A float64 that indicates the return of a function that
    needs to equal zero to be drawn. This is what you want
    for relations that are of a similar form as "f(x, y) == g(x, y)"
    but slightly modified to look like "f(x, y) - g(x, y) == 0".
    This is because very rarely will any coordinate the relation
    is called with be exactly what you need to make it true, and
    so the DrawRelation function compares the signs of the returned
    values of the surrounding coordinates and if there is a difference
    in the signs of said values, and the relation is continuous, then
    we know that the returned value would equal zero somewhere
    between those coordinates.

    An error that indicates the current chunk should stop
    being drawn immediately.
*/
type Relation func (c *Coord) interface{}

/*
    A function that takes in an x value
    and returns a y value. This is what
    you want for relations of the form
    "y == f(x)" as it will be drawn faster
    and more accurately than using the
    Relation version.
*/
type Function func (x float64) float64

/*
    A function that takes in a theta value
    and returns a radius value. This is what
    you want for relations of the form
    "r == f(theta)" as it will be drawn faster
    and more accurately than using the
    Relation version.
*/
type PolarFunction func (theta float64) float64

/*
    A function that takes is a complex number
    and returns a complex number. This is used
    to treat a graph as the complex plane and
    then move the coordinates that correspond
    to the input to the return value of the ComplexRelation.
*/
type ComplexRelation func (z complex128) complex128

/*
    A function that takes in a coordinate and
    returns the derivate of a function at that point
*/
type DifferentialFunction func (c *Coord) float64

/*
    An area in coordinate space. To work nicely,
    Pos0 must be above and to the left (lower x
    value and higher y value) of Pos1.
*/
type Area struct {
    Pos0, Pos1 *Coord
}

type Graph struct {
    Bounds *Area
    Image *image.RGBA

    BackgroundColor, RelationColor, AxisColor, GridColor color.Color
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

func (c *Coord) IsValid() bool {
    return !math.IsInf(c.X, 1) && !math.IsInf(c.X, -1) && !math.IsNaN(c.X) &&
           !math.IsInf(c.Y, 1) && !math.IsInf(c.Y, -1) && !math.IsNaN(c.Y)
}

func (f Function) ToRelation() Relation {
    return func (c *Coord) interface{} {
        return c.Y - f(c.X)
    }
}

func (f PolarFunction) ToRelation() Relation {
    return func (c *Coord) interface{} {
        r, theta := c.Polar()

        return r - f(theta)
    }
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

func NewGraphWithColors(bounds *Area, scale float64, bg_col, rel_col, axis_col, grid_col color.Color) *Graph {
    g := &Graph{}

    g.Bounds = bounds
    g.Image = image.NewRGBA(image.Rect(0, 0, int(bounds.Width() * scale), int(bounds.Height() * scale)))

    g.BackgroundColor = bg_col
    g.RelationColor = rel_col
    g.AxisColor = axis_col
    g.GridColor = grid_col

    for x := 0; x < g.ImageWidth(); x++ {
        for y := 0; y < g.ImageHeight(); y++ {
            g.SetPixel(image.Pt(x, y), bg_col)
        }
    }

    return g
}

func NewGraph(bounds *Area, scale float64) *Graph {
    return NewGraphWithColors(bounds, scale, DefaultBackgroundColor, DefaultRelationColor, DefaultAxisColor, DefaultGridColor)
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
    if !c0.IsValid() || !c1.IsValid() {
        return
    }

    var p0, p1 image.Point

    if (c0.X <= c1.X) {
        p0 = g.CoordToPixel(c0)
        p1 = g.CoordToPixel(c1)
    } else {
        p0 = g.CoordToPixel(c1)
        p1 = g.CoordToPixel(c0)
    }

    delta := p1.Sub(p0)

    if delta.X == 0 { // Vertical line
        diff := 1
        if delta.Y < 0 {
            diff = -1
        }

        for ; p0.Y != p1.Y; p0.Y += diff {
            g.SetPixel(p0, col)
        }
    } else if delta.Y == 0 { // Horizontal line
        for ; p0.X < p1.X; p0.X++ {
            g.SetPixel(p0, col)
        }
    } else {
        y_diff := -1
        if delta.Y > 0 {
            delta.Y = -delta.Y
            y_diff = 1
        }

        err := delta.X + delta.Y

        for {
            g.SetPixel(p0, col)

            if p0.Eq(p1) {
                break
            }

            tmp_err := 2 * err

            if tmp_err >= delta.Y {
                err += delta.Y
                p0.X++
            }

            if tmp_err <= delta.X {
                err += delta.X
                p0.Y += y_diff
            }
        }
    }
}

func (g *Graph) DrawAxes() {
    g.DrawLine(NewCoord(0, g.Bounds.Pos0.Y), NewCoord(0, g.Bounds.Pos1.Y), g.AxisColor)
    g.DrawLine(NewCoord(g.Bounds.Pos0.X, 0), NewCoord(g.Bounds.Pos1.X, 0), g.AxisColor)
}

func (g *Graph) DrawGrid() {
    for x := 1.0; x < g.Bounds.Pos1.X; x++ {
        g.DrawLine(NewCoord(x, g.Bounds.Pos0.Y), NewCoord(x, g.Bounds.Pos1.Y), g.GridColor)
        g.DrawLine(NewCoord(-x, g.Bounds.Pos0.Y), NewCoord(-x, g.Bounds.Pos1.Y), g.GridColor)

        for y := 1.0; y < g.Bounds.Pos0.Y; y++ {
            g.DrawLine(NewCoord(g.Bounds.Pos0.X, y), NewCoord(g.Bounds.Pos1.X, y), g.GridColor)
            g.DrawLine(NewCoord(g.Bounds.Pos0.X, -y), NewCoord(g.Bounds.Pos1.X, -y), g.GridColor)
        }
    }

    g.DrawAxes()
}

func (g *Graph) DrawRelationInChunk(rel Relation, r *image.Rectangle, col color.Color, ch chan struct{}) {
    for x := r.Min.X; x < r.Max.X; x++ {
        for y := r.Min.Y; y < r.Max.Y; y++ {
            pt := image.Pt(x, y)
            c := g.PixelToCoord(pt)

            switch ret := rel(c); ret.(type) {
                case bool:
                    if ret.(bool) {
                        g.SetPixel(pt, col)
                    }

                case float64:
                    if math.IsNaN(ret.(float64)) || math.IsInf(ret.(float64), 1) || math.IsInf(ret.(float64), -1) {
                        break
                    }

                    coords := [3]*Coord {
                        g.PixelToCoord(pt.Add(image.Pt(1, 0))),
                        g.PixelToCoord(pt.Add(image.Pt(0, 1))),
                        g.PixelToCoord(pt.Add(image.Pt(1, 1))),
                    }

                    diff := [3]float64 {
                        rel(coords[0]).(float64),
                        rel(coords[1]).(float64),
                        rel(coords[2]).(float64),
                    }

                    if ret.(float64) == 0 {
                        g.SetPixel(pt, col)
                        continue
                    }

                    for _, d := range diff {
                        if (ret.(float64) > 0 && d < 0) || (ret.(float64) < 0 && d > 0) {
                            g.SetPixel(pt, col)
                            break
                        }
                    }

                case error:
                    ch <- struct{}{}
                    return
            }
        }
    }

    ch <- struct{}{}
}

func (g *Graph) DrawRelationWithColor(rel Relation, col color.Color) {
    var channels []chan struct{}

    for x := 0; x < g.ImageWidth(); x += ChunkSize {
        for y := 0; y < g.ImageHeight(); y += ChunkSize {
            ch := make(chan struct {})
            channels = append(channels, ch)

            r := image.Rect(x, y, MinInt(x + ChunkSize, g.ImageWidth()), MinInt(y + ChunkSize, g.ImageHeight()))
            go g.DrawRelationInChunk(rel, &r, col, ch)
        }
    }

    for _, ch := range channels {
        <-ch
    }
}

func (g *Graph) DrawRelation(rel Relation) {
    g.DrawRelationWithColor(rel, g.RelationColor)
}

func (g *Graph) ApplyComplexRelationInChunk(rel ComplexRelation, dst *image.RGBA, r *image.Rectangle, ch chan struct{}) {
    for x := r.Min.X; x < r.Max.X; x++ {
        for y := r.Min.Y; y < r.Max.Y; y++ {
            pt := image.Pt(x, y)
            c := g.PixelToCoord(pt)
            z := complex(c.X, c.Y)

            new_z := rel(z)
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

func (g *Graph) ApplyComplexRelation(rel ComplexRelation) {
    img := image.NewRGBA(g.Image.Bounds())
    for x := 0; x < g.ImageWidth(); x++ {
        for y := 0; y < g.ImageHeight(); y++ {
            g.SetPixel(image.Pt(x, y), g.BackgroundColor)
        }
    }

    var channels []chan struct{}

    for x := 0; x < g.ImageWidth(); x += ChunkSize {
        for y := 0; y < g.ImageHeight(); y += ChunkSize {
            ch := make(chan struct {})
            channels = append(channels, ch)

            r := image.Rect(x, y, MinInt(x + ChunkSize, g.ImageWidth()), MinInt(y + ChunkSize, g.ImageHeight()))
            go g.ApplyComplexRelationInChunk(rel, img, &r, ch)
        }
    }

    for _, ch := range channels {
        <-ch
    }

    g.Image = img
}

func (g *Graph) DrawDifferentialFunctionInDirection(d DifferentialFunction, start *Coord, dx float64, col color.Color, ch chan struct{}) {
    old := start

    for i := 0; i < g.ImageWidth(); i++ {
        start = start.Add(NewCoord(dx, d(start) * dx))
        g.DrawLine(start, old, col)
        old = start

        if !g.Bounds.Contains(start) {
            break
        }
    }

    ch <- struct{}{}
}

func (g *Graph) DrawDifferentialFunctionWithColor(d DifferentialFunction, start *Coord, col color.Color) {
    channels := [2]chan struct{} {
        make(chan struct{}),
        make(chan struct{}),
    }

    dx := g.Bounds.Width() / float64(g.ImageWidth())

    go g.DrawDifferentialFunctionInDirection(d, start, dx, col, channels[0])
    go g.DrawDifferentialFunctionInDirection(d, start, -dx, col, channels[1])

    for _, ch := range channels {
        <-ch
    }
}

func (g *Graph) DrawDifferentialFunction(d DifferentialFunction, start *Coord) {
    g.DrawDifferentialFunctionWithColor(d, start, g.RelationColor)
}

func (g *Graph) DrawFunctionInRange(f Function, start, end int, col color.Color, ch chan struct{}) {
    real_start := g.PixelToCoord(image.Pt(start, 0)).X
    old := NewCoord(real_start, f(real_start))

    for x := start + 1; x <= end; x++ {
        real_x := g.PixelToCoord(image.Pt(x, 0)).X

        new := NewCoord(real_x, f(real_x))
        g.DrawLine(new, old, col)
        old = new
    }

    ch <- struct{}{}
}

func (g *Graph) DrawFunctionWithColor(f Function, col color.Color) {
    var channels []chan struct{}

    for x := 0; x < g.ImageWidth(); x += ChunkSize {
        ch := make(chan struct{})
        channels = append(channels, ch)

        go g.DrawFunctionInRange(f, x, MinInt(x + ChunkSize, g.ImageWidth()), col, ch)
    }

    for _, ch := range channels {
        <-ch
    }
}

func (g *Graph) DrawFunction(f Function) {
    g.DrawFunctionWithColor(f, g.RelationColor)
}

func (g *Graph) DrawPolarFunctionInRange(f PolarFunction, start, end float64, col color.Color, ch chan struct {}) {
    old := NewCoordFromPolar(f(start), start)

    for theta := start + AngleStep; theta <= end + AngleStep; theta += AngleStep {
        new := NewCoordFromPolar(f(theta), theta)
        g.DrawLine(new, old, col)
        old = new
    }

    ch <- struct{}{}
}

func (g *Graph) DrawPolarFunctionWithColor(f PolarFunction, col color.Color) {
    var channels []chan struct{}

    for theta := 0.0; theta < 2 * math.Pi; theta += AngleSize {
        ch := make(chan struct{})
        channels = append(channels, ch)

        go g.DrawPolarFunctionInRange(f, theta, math.Min(theta + AngleSize, 2 * math.Pi), col, ch)
    }

    for _, ch := range channels {
        <-ch
    }
}

func (g *Graph) DrawPolarFunction(f PolarFunction) {
    g.DrawPolarFunctionWithColor(f, g.RelationColor)
}