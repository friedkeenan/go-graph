package main

import (
    "math"
    "math/cmplx"
)

const (
    MaxIterations = 200
    DifferentiateDx = 0.01
)

func OffsetRelation(rel Relation, off *Coord) Relation {
    return func (c *Coord) interface{} {
        c = c.Sub(off)

        return rel(c)
    }
}

func ScaleRelation(rel Relation, scale float64) Relation {
    return func (c *Coord) interface{} {
        c = c.Div(scale)

        return rel(c)
    }
}

func ScaleRelationAround(rel Relation, scale float64, coord *Coord) Relation {
    return func (c *Coord) interface{} {
        c = c.Sub(coord).Div(scale).Add(coord)

        return rel(c)
    }
}

func ScaleRelationPerAxis(rel Relation, scale_x, scale_y float64) Relation {
    return func (c *Coord) interface{} {
        c = NewCoord(c.X / scale_x, c.Y / scale_y)

        return rel(c)
    }
}

func ScaleRelationPerAxisAround(rel Relation, scale_x, scale_y float64, coord *Coord) Relation {
    return func (c *Coord) interface{} {
        c = c.Sub(coord)
        c = NewCoord(c.X / scale_x, c.Y / scale_y)
        c = c.Add(coord)

        return rel(c)
    }
}

func RotateRelation(rel Relation, theta float64) Relation {
    return func (c *Coord) interface{} {
        c = c.Rotate(-theta)

        return rel(c)
    }
}

func RotateRelationAround(rel Relation, theta float64, coord *Coord) Relation {
    return func (c *Coord) interface{} {
        c = c.RotateAround(-theta, coord)

        return rel(c)
    }
}

func InvertRelation(rel Relation) Relation {
    return func (c *Coord) interface{} {
        c = NewCoord(c.Y, c.X)

        return rel(c)
    }
}

func OffsetFunction(f Function, off *Coord) Function {
    return func (x float64) float64 {
        return f(x - off.X) + off.Y
    }
}

func ScaleFunction(f Function, scale float64) Function {
    return func (x float64) float64 {
        return scale * f(x / scale)
    }
}

func ScaleFunctionPerAxis(f Function, scale_x, scale_y float64) Function {
    return func (x float64) float64 {
        return scale_y * f(x / scale_x)
    }
}

func DifferentiateFunction(f Function) Function {
    return func(x float64) float64 {
        return (f(x + DifferentiateDx) - f(x)) / DifferentiateDx
    }
}

func IntegrateFunction(f Function, a, b float64) float64 {
    sum := 0.0

    if a < b {
        for x := a; x < b; x += DifferentiateDx {
            sum += f(x) * DifferentiateDx
        }
    } else {
        for x := a; x > b; x -= DifferentiateDx {
            sum -= f(x) * DifferentiateDx
        }
    }

    return sum
}

func AntiDifferentiateFunction(f Function, a float64) Function {
    return func (x float64) float64 {
        return IntegrateFunction(f, a, x)
    }
}

func Mandelbrot(c *Coord) interface{} {
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

func UnitCircle(c *Coord) interface{} {
    return math.Pow(c.X, 2) + math.Pow(c.Y, 2) - 1
}

func Circle(r float64) Relation {
    return ScaleRelation(UnitCircle, r)
}

func CircleAt(r float64, center *Coord) Relation {
    return OffsetRelation(Circle(r), center)
}

func Ellipse(a, b float64) Relation {
    return ScaleRelationPerAxis(UnitCircle, a, b)
}

func EllipseAt(a, b float64, center *Coord) Relation {
    return OffsetRelation(Ellipse(a, b), center)
}