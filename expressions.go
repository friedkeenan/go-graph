package main

import (
    "math"
    "math/cmplx"
)

const (
    MaxIterations = 200
)

func OffsetBoolExpression(expr BoolExpression, off *Coord) BoolExpression {
    return func (c *Coord) bool {
        c = c.Sub(off)

        return expr(c)
    }
}

func ScaleBoolExpression(expr BoolExpression, scale float64) BoolExpression {
    return func (c *Coord) bool {
        c = c.Div(scale)

        return expr(c)
    }
}

func ScaleBoolExpressionAround(expr BoolExpression, scale float64, coord *Coord) BoolExpression {
    return func (c *Coord) bool {
        c = c.Sub(coord).Div(scale).Add(coord)

        return expr(c)
    }
}

func ScaleBoolExpressionPerAxis(expr BoolExpression, scale_x, scale_y float64) BoolExpression {
    return func (c *Coord) bool {
        c = NewCoord(c.X / scale_x, c.Y / scale_y)

        return expr(c)
    }
}

func ScaleBoolExpressionPerAxisAround(expr BoolExpression, scale_x, scale_y float64, coord *Coord) BoolExpression {
    return func (c *Coord) bool {
        c = c.Sub(coord)
        c = NewCoord(c.X / scale_x, c.Y / scale_y)
        c = c.Add(coord)

        return expr(c)
    }
}

func RotateBoolExpression(expr BoolExpression, theta float64) BoolExpression {
    return func (c *Coord) bool {
        c = c.Rotate(-theta)

        return expr(c)
    }
}

func RotateBoolExpressionAround(expr BoolExpression, theta float64, coord *Coord) BoolExpression {
    return func (c *Coord) bool {
        c = c.RotateAround(-theta, coord)

        return expr(c)
    }
}

func OffsetDiffExpression(expr DiffExpression, off *Coord) DiffExpression {
    return func (c *Coord) float64 {
        c = c.Sub(off)

        return expr(c)
    }
}

func ScaleDiffExpression(expr DiffExpression, scale float64) DiffExpression {
    return func (c *Coord) float64 {
        c = c.Div(scale)

        return expr(c)
    }
}

func ScaleDiffExpressionAround(expr DiffExpression, scale float64, coord *Coord) DiffExpression {
    return func (c *Coord) float64 {
        c = c.Sub(coord).Div(scale).Add(coord)

        return expr(c)
    }
}

func ScaleDiffExpressionPerAxis(expr DiffExpression, scale_x, scale_y float64) DiffExpression {
    return func (c *Coord) float64 {
        c = NewCoord(c.X / scale_x, c.Y / scale_y)

        return expr(c)
    }
}

func ScaleDiffExpressionPerAxisAround(expr DiffExpression, scale_x, scale_y float64, coord *Coord) DiffExpression {
    return func (c *Coord) float64 {
        c = c.Sub(coord)
        c = NewCoord(c.X / scale_x, c.Y / scale_y)
        c = c.Add(coord)

        return expr(c)
    }
}

func RotateDiffExpression(expr DiffExpression, theta float64) DiffExpression {
    return func (c *Coord) float64 {
        c = c.Rotate(-theta)

        return expr(c)
    }
}

func RotateDiffExpressionAround(expr DiffExpression, theta float64, coord *Coord) DiffExpression {
    return func (c *Coord) float64 {
        c = c.RotateAround(-theta, coord)

        return expr(c)
    }
}

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

func UnitCircle(c *Coord) float64 {
    return math.Pow(c.X, 2) + math.Pow(c.Y, 2) - 1
}

func Circle(r float64) DiffExpression {
    return ScaleDiffExpression(UnitCircle, r)
}

func CircleAt(r float64, center *Coord) DiffExpression {
    return OffsetDiffExpression(Circle(r), center)
}

func Ellipse(a, b float64) DiffExpression {
    return ScaleDiffExpressionPerAxis(UnitCircle, a, b)
}

func EllipseAt(a, b float64, center *Coord) DiffExpression {
    return OffsetDiffExpression(Ellipse(a, b), center)
}