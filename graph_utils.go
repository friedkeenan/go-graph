package main

import (
    "math"
    "math/cmplx"
)

const (
    MaxIterations = 200
)

func OffsetExpression(expr Expression, off *Coord) Expression {
    return func (c *Coord) interface{} {
        c = c.Sub(off)

        return expr(c)
    }
}

func ScaleExpression(expr Expression, scale float64) Expression {
    return func (c *Coord) interface{} {
        c = c.Div(scale)

        return expr(c)
    }
}

func ScaleExpressionAround(expr Expression, scale float64, coord *Coord) Expression {
    return func (c *Coord) interface{} {
        c = c.Sub(coord).Div(scale).Add(coord)

        return expr(c)
    }
}

func ScaleExpressionPerAxis(expr Expression, scale_x, scale_y float64) Expression {
    return func (c *Coord) interface{} {
        c = NewCoord(c.X / scale_x, c.Y / scale_y)

        return expr(c)
    }
}

func ScaleExpressionPerAxisAround(expr Expression, scale_x, scale_y float64, coord *Coord) Expression {
    return func (c *Coord) interface{} {
        c = c.Sub(coord)
        c = NewCoord(c.X / scale_x, c.Y / scale_y)
        c = c.Add(coord)

        return expr(c)
    }
}

func RotateExpression(expr Expression, theta float64) Expression {
    return func (c *Coord) interface{} {
        c = c.Rotate(-theta)

        return expr(c)
    }
}

func RotateExpressionAround(expr Expression, theta float64, coord *Coord) Expression {
    return func (c *Coord) interface{} {
        c = c.RotateAround(-theta, coord)

        return expr(c)
    }
}

func InvertExpression(expr Expression) Expression {
    return func (c *Coord) interface{} {
        c = NewCoord(c.Y, c.X)

        return expr(c)
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

func Circle(r float64) Expression {
    return ScaleExpression(UnitCircle, r)
}

func CircleAt(r float64, center *Coord) Expression {
    return OffsetExpression(Circle(r), center)
}

func Ellipse(a, b float64) Expression {
    return ScaleExpressionPerAxis(UnitCircle, a, b)
}

func EllipseAt(a, b float64, center *Coord) Expression {
    return OffsetExpression(Ellipse(a, b), center)
}