package naglevelparse

import (
	"errors"
	kb "github.com/jamesandariese/koalaburrito"
	"math"
	"strconv"
)

type Comparator struct {
	min float64
	max float64
	at  bool
}

func parseNagiosFloat(s string, def float64) float64 {
	if s == "~" {
		return math.Inf(-1)
	}
	if s == "" {
		return def
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		// this should be impossible.
		panic(err)
	}
	return f
}

func (c *Comparator) Compare(value float64) bool {
	rv := false
	if value < c.min || value > c.max {
		rv = true
	}
	if c.at {
		rv = !rv
	}
	return rv
}

var EmptyInputError = errors.New("empty string is not a valid range")
var ParseError = errors.New("pattern has unrecognized components")

func Compile(f string) (*Comparator, error) {
	c := &Comparator{}

	t := kb.MakeTokenizer()
	t.AddPattern(`([-+]?([0-9]+\.[0-9]*|[0-9]*\.[0-9]+|[0-9]+))`, "NUMBER")
	t.AddPattern("~", "TILDE")
	t.AddPattern("@", "AT")
	t.AddPattern(":", "COLON")

	offset := 0

	tokens, err := t.MatchAll(f)
	if err != nil {
		return nil, err
	}
	if len(tokens) < 1 {
		return nil, EmptyInputError
	}

	c.at = false
	if tokens[0].Value() == "@" {
		c.at = true
		offset += 1
	}
	if len(tokens)-offset == 3 {
		if tokens[offset+1].Name() == "COLON" {
			if tokens[offset+0].Name() == "NUMBER" {
				c.min, _ = strconv.ParseFloat(tokens[offset+0].Value(), 64)
				if tokens[offset+2].Name() == "NUMBER" {
					c.max, _ = strconv.ParseFloat(tokens[offset+2].Value(), 64)
					return c, nil
				} else if tokens[offset+2].Name() == "TILDE" {
					c.max = math.Inf(-1)
					return c, nil
				}
			} else if tokens[offset+0].Name() == "TILDE" {
				c.min = math.Inf(-1)
				if tokens[offset+2].Name() == "NUMBER" {
					c.max, _ = strconv.ParseFloat(tokens[offset+2].Value(), 64)
					return c, nil
				} else if tokens[offset+2].Name() == "TILDE" {
					c.max = math.Inf(-1)
					return c, nil
				}
			}
		}
	} else if len(tokens)-offset == 2 {
		if tokens[offset+0].Name() == "COLON" {
			c.min = 0
			if tokens[offset+1].Name() == "NUMBER" {
				c.max, _ = strconv.ParseFloat(tokens[offset+1].Value(), 64)
				return c, nil
			} else if tokens[offset+1].Name() == "TILDE" {
				c.max = math.Inf(-1)
				return c, nil
			}
		} else if tokens[offset+1].Name() == "COLON" {
			c.max = math.Inf(1)
			if tokens[offset+0].Name() == "NUMBER" {
				c.min, _ = strconv.ParseFloat(tokens[offset+0].Value(), 64)
				return c, nil
			} else if tokens[offset+0].Name() == "TILDE" {
				c.min = math.Inf(-1)
				return c, nil
			}
		}
	} else if len(tokens)-offset == 1 {
		c.min = 0
		if tokens[offset+0].Name() == "NUMBER" {
			c.max, _ = strconv.ParseFloat(tokens[offset+0].Value(), 64)
			return c, nil
		} else if tokens[offset+0].Name() == "TILDE" {
			c.max = math.Inf(-1)
			return c, nil
		}
	}
	return nil, ParseError
}
