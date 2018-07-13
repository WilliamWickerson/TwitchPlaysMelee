package main

import (
	"vJoy"
)

type GamecubeController struct {
	A, B, X, Y, Z, L, R, START bool;
	DLEFT, DUP, DRIGHT, DDOWN bool;
	STICKX, STICKY, CX, CY, LANA, RANA float64;
}

//Floats default to 0 so having the sticks default to center is desired
func NewGamecubeController() GamecubeController {
	return GamecubeController{
		STICKX : .5,
		STICKY : .5,
		CX : .5,
		CY : .5,
	}
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
)

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
	c.SetAxis(int(gc.STICKX*32684), STICKX);
	c.SetAxis(int(gc.STICKY*32684), STICKY);
	c.SetAxis(int(gc.CX*32684), CX);
	c.SetAxis(int(gc.CY*32684), CY);
	c.SetAxis(int(gc.LANA*32684), LANA);
	c.SetAxis(int(gc.RANA*32684), RANA);
}