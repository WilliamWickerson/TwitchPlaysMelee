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
	angle := angle(x, y);
	var radius float64;
	if tilt {
		radius = .3;
	} else {
		radius = math.Min(math.Hypot(x, y), 1);
	}
	return direction{angle, radius};
}

func NewDirection(directions []token.Type, tilt bool) Direction {
	xVal := 0;
	yVal := 0;
	center := 1;
	for _,dir := range directions {
		switch(dir) {
			case token.KW_LEFT:
				xVal--;
				break;
			case token.KW_RIGHT:
				xVal++;
				break;
			case token.KW_UP:
				yVal++;
				break;
			case token.KW_DOWN:
				yVal--;
				break;
			case token.KW_CENTER:
				center++;
				break;
		}
	}
	angle := angle(float64(xVal), float64(yVal));
	var radius float64;
	if tilt {
		radius = .3;
	} else {
		radius = 1 / float64(center);
	}
	return direction{angle, radius};
}

func (d direction) Components() (float64, float64) {
	return d.radius*math.Cos(d.angle), d.radius*math.Sin(d.angle);
}

func angle(x, y float64) float64 {
	angle := math.Acos(x / math.Hypot(x,y));
	if y < 0 {
		angle = 2*math.Pi - angle;
	}
	return angle;
}
