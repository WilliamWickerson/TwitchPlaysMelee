package controller

import (
	"controller/vJoy"
	"math"
	"fmt"
)

type GamecubeController struct {
	A, B, X, Y, Z, L, R, START bool;
	DLEFT, DUP, DRIGHT, DDOWN bool;
	STICKX, STICKY, CX, CY, LANA, RANA float64;
}

const (
	//Constants to map the 12 gamecube buttons
	A = 1;
	B = 2;
	X = 3;
	Y = 4;
	Z = 5;
	L = 6;
	R = 7;
	START = 8;
	DLEFT = 9;
	DUP = 10;
	DRIGHT = 11;
	DDOWN = 12;
	//Constants to map the 6 axes
	STICKX = vJoy.X;
	STICKY = vJoy.Y;
	CX = vJoy.Rx;
	CY = vJoy.Ry;
	LANA = vJoy.Sl0;
	RANA = vJoy.Sl1;
	//Constants to set vJoy correctly
	center = 16384; //about .00025 off towards negative at 16384 mirrored at 16385
	stickMultiplier = (16384*.300)/.475;
	stickNegMultiplier = (16384*.600)/.9625;
	sliderMultipler = (16384*.300)/.550;
)

func round(x float64) int {
	return int(math.Round(x));
}

func stickValue(x float64) int {
	if x >= 0 {
		x = x * stickMultiplier;
	} else {
		x = x * stickNegMultiplier;
	}
	return round(x) + center;
}

func SetvJoy(c vJoy.VJoyController, gc GamecubeController) {
	//Set the 12 vJoy buttons based on the Gamecube Controller State
	c.SetButton(gc.A, A);
	c.SetButton(gc.B, B);
	c.SetButton(gc.X, X);
	c.SetButton(gc.Y, Y);
	c.SetButton(gc.Z, Z);
	c.SetButton(gc.L, L);
	c.SetButton(gc.R, R);
	c.SetButton(gc.START, START);
	c.SetButton(gc.DLEFT, DLEFT);
	c.SetButton(gc.DUP, DUP);
	c.SetButton(gc.DRIGHT, DRIGHT);
	c.SetButton(gc.DDOWN, DDOWN);
	//Set the 6 axes based on the Gamecube Controller State
	c.SetAxis(stickValue(gc.STICKX), STICKX);
	c.SetAxis(stickValue(gc.STICKY), STICKY);
	c.SetAxis(stickValue(gc.CX), CX);
	c.SetAxis(stickValue(gc.CY), CY);
	c.SetAxis(int(round(gc.LANA*sliderMultipler) + center), LANA);
	c.SetAxis(int(round(gc.RANA*sliderMultipler) + center), RANA);
}

func PrintGC(gc GamecubeController) {
	//IDE gives tab a length of 4 characters while cmd gives it a length of 8
	//since the latter is preferred, the spaces are used to force that behavior
	fmt.Printf("A    \tB    \tX    \tY    \tZ    \tSTART\tDLEFT\tDRIGHT\tDUP  \tDDOWN\n");
	fmt.Printf("%t\t%t\t%t\t%t\t%t\t%t\t%t\t%t\t%t\t%t\n", gc.A, gc.B, gc.X, gc.Y, gc.Z, gc.START, gc.DLEFT, gc.DRIGHT, gc.DUP, gc.DDOWN);
	fmt.Printf("L: %t %f\n", gc.L, gc.LANA);
	fmt.Printf("R: %t %f\n", gc.R, gc.RANA);
	fmt.Printf("Stick: (%f, %f)\n", gc.STICKX, gc.STICKY);
	fmt.Printf("CStick: (%f, %f)\n", gc.CX, gc.CY);
}