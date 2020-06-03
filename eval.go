package main

import (
    "strings"
    "log"
    "math"
    "github.com/Knetic/govaluate"
)

var (
    Functions = map[string]govaluate.ExpressionFunction {
        "abs": func (args ...interface{}) (interface{}, error) {
            return math.Abs(args[0].(float64)), nil
        },

        "acos": func (args ...interface{}) (interface{}, error) {
            return math.Acos(args[0].(float64)), nil
        },

        "acosh": func (args ...interface{}) (interface{}, error) {
            return math.Acosh(args[0].(float64)), nil
        },

        "asin": func (args ...interface{}) (interface{}, error) {
            return math.Asin(args[0].(float64)), nil
        },

        "asinh": func (args ...interface{}) (interface{}, error) {
            return math.Asinh(args[0].(float64)), nil
        },

        "atan": func (args ...interface{}) (interface{}, error) {
            return math.Atan(args[0].(float64)), nil
        },

        "atan2": func (args ...interface{}) (interface{}, error) {
            return math.Atan2(args[0].(float64), args[1].(float64)), nil
        },

        "atanh": func (args ...interface{}) (interface{}, error) {
            return math.Atanh(args[0].(float64)), nil
        },

        "ceil": func (args ...interface{}) (interface{}, error) {
            return math.Ceil(args[0].(float64)), nil
        },

        "cos": func (args ...interface{}) (interface{}, error) {
            return math.Cos(args[0].(float64)), nil
        },

        "cosh": func (args ...interface{}) (interface{}, error) {
            return math.Cosh(args[0].(float64)), nil
        },

        "exp": func (args ...interface{}) (interface{}, error) {
            return math.Exp(args[0].(float64)), nil
        },

        "floor": func (args ...interface{}) (interface{}, error) {
            return math.Floor(args[0].(float64)), nil
        },

        "gamma": func (args ...interface{}) (interface{}, error) {
            return math.Gamma(args[0].(float64)), nil
        },

        "ln": func (args ...interface{}) (interface{}, error) {
            return math.Log(args[0].(float64)), nil
        },

        "log": func (args ...interface{}) (interface{}, error) {
            return math.Log10(args[0].(float64)), nil
        },

        "sin": func (args ...interface{}) (interface{}, error) {
            return math.Sin(args[0].(float64)), nil
        },

        "sinh": func (args ...interface{}) (interface{}, error) {
            return math.Sinh(args[0].(float64)), nil
        },

        "sqrt": func (args ...interface{}) (interface{}, error) {
            return math.Sqrt(args[0].(float64)), nil
        },

        "tan": func (args ...interface{}) (interface{}, error) {
            return math.Tan(args[0].(float64)), nil
        },

        "tanh": func (args ...interface{}) (interface{}, error) {
            return math.Tanh(args[0].(float64)), nil
        },

    }

    Constants = map[string]interface{} {
        "pi":  math.Pi,
        "tau": 2 * math.Pi,
        "e":   math.E,
        "phi": math.Phi,
    }
)

type NoEqualityError struct{}

func (e NoEqualityError) Error() string {
    return "No equality in expression"
}

func EvalBoolExpression(expr string) (BoolExpression, error) {
    e, err := govaluate.NewEvaluableExpressionWithFunctions(expr, Functions)
    if err != nil {
        return nil, err
    }

    return func (c *Coord) bool {
        params := map[string]interface{} {
            "x": c.X,
            "y": c.Y,
        }
        params["r"], params["theta"] = c.Polar()

        for k, v := range Constants {
            params[k] = v
        }

        result, err := e.Evaluate(params)
        if err != nil {
            log.Fatal(err)
        }

        return result.(bool)
    }, nil
}

func EvalDiffExpression(expr string) (DiffExpression, error) {
    sides := strings.Split(expr, "==")
    if len(sides) != 2 {
        return nil, NoEqualityError{}
    }

    e0, err := govaluate.NewEvaluableExpressionWithFunctions(sides[0], Functions)
    if err != nil {
        return nil, err
    }

    e1, err := govaluate.NewEvaluableExpressionWithFunctions(sides[1], Functions)
    if err != nil {
        return nil, err
    }

    return func (c *Coord) float64 {
        params := map[string]interface{} {
            "x": c.X,
            "y": c.Y,
        }
        params["r"], params["theta"] = c.Polar()

        for k, v := range Constants {
            params[k] = v
        }

        result0, err := e0.Evaluate(params)
        if err != nil {
            log.Fatal(err)
        }

        result1, err := e1.Evaluate(params)
        if err != nil {
            log.Fatal(err)
        }

        return result0.(float64) - result1.(float64)
    }, nil
}