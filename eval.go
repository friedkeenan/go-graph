package main

import (
    "strings"
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
    return "No equality in relation"
}

func IsFunction(dep, indep string, vars0, vars1 []string, tokens0, tokens1 []govaluate.ExpressionToken) (bool, bool) {
    if len(tokens0) == 1 && len(vars0) == 1 && vars0[0] == dep {
        for _, v := range vars1 {
            if _, ok := Constants[v]; v != indep && !ok {
                return false, false
            }
        }

        return true, true
    }

    if len(tokens1) == 1 && len(vars1) == 1 && vars1[0] == dep {
        for _, v := range vars0 {
            if _, ok := Constants[v]; v != indep && !ok {
                return false, false
            }
        }

        return true, false
    }

    return false, false
}

func Eval(expr string) (interface{}, error) {
    if strings.Contains(expr, "==") {
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

        vars0, tokens0 := e0.Vars(), e0.Tokens()
        vars1, tokens1 := e1.Vars(), e1.Tokens()

        if is_func, dep_first := IsFunction("y", "x", vars0, vars1, tokens0, tokens1); is_func {
            e := e1
            if !dep_first {
                e = e0
            }

            return Function(func (x float64) float64 {
                params := map[string]interface{} {
                    "x": x,
                }

                for k, v := range Constants {
                    params[k] = v
                }

                result, err := e.Evaluate(params)
                if err != nil {
                    return math.NaN()
                }

                return result.(float64)
            }), nil
        }

        if is_func, dep_first := IsFunction("r", "theta", vars0, vars1, tokens0, tokens1); is_func {
            e := e1
            if !dep_first {
                e = e0
            }

            return PolarFunction(func (theta float64) float64 {
                params := map[string]interface{} {
                    "theta": theta,
                }

                for k, v := range Constants {
                    params[k] = v
                }

                result, err := e.Evaluate(params)
                if err != nil {
                    return math.NaN()
                }

                return result.(float64)
            }), nil
        }

        return Relation(func (c *Coord) interface{} {
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
                return err
            }

            result1, err := e1.Evaluate(params)
            if err != nil {
                return err
            }

            return result0.(float64) - result1.(float64)
        }), nil
    }

    e, err := govaluate.NewEvaluableExpressionWithFunctions(expr, Functions)
    if err != nil {
        return nil, err
    }

    return Relation(func (c *Coord) interface{} {
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
            return err
        }

        return result
    }), nil
}