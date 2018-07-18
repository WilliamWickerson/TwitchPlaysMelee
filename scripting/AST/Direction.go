package AST

import (
	"scripting/token"
	"math"
)

type Direction interface {
	Components() (float64, float64);
}

type direction struct {
	angle float64;
	radius float64;
}

func NewDirectionPair(x float64, y float64, tilt bool) Direction {
	return nil;
}

func NewDirection(directions []token.Type, tilt bool) Direction {
	return nil;
}

func (d direction) Components() (float64, float64) {
	return d.radius*math.Cos(d.angle), d.radius*math.Sin(d.angle);
}
