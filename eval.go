package main

import (
    "strings"
    "log"
    "math"
    "github.com/Knetic/govaluate"
)

var (
    Functions = map[string]govaluate.ExpressionFunction {
        "sin": func (args ...interface{}) (interface{}, error) {
            return math.Sin(args[0].(float64)), nil
        },

        "cos": func (args ...interface{}) (interface{}, error) {
            return math.Cos(args[0].(float64)), nil
        },

        "tan": func (args ...interface{}) (interface{}, error) {
            return math.Tan(args[0].(float64)), nil
        },

        "exp": func (args ...interface{}) (interface{}, error) {
            return math.Exp(args[0].(float64)), nil
        },

        "ln": func (args ...interface{}) (interface{}, error) {
            return math.Log(args[0].(float64)), nil
        },

        "log": func (args ...interface{}) (interface{}, error) {
            return math.Log10(args[0].(float64)), nil
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