package vJoy

import (
	"testing"
)

func TestSetup(t *testing.T) {
	if !Enabled(){
		t.Errorf("vJoy is not enabled");
	}
	if GetManufacturerString() != "Shaul Eizikovich" {
		t.Errorf("Manufacturer incorrect, expected: Shaul Eizikovich, got: %s", GetManufacturerString());
	}
	if GetProductString() != "vJoy - Virtual Joystick" {
		t.Errorf("Product incorrect, expected: vJoy - Virtual Joystick, got: %s", GetProductString());
	}
	if GetSerialNumberString() != "2.1.8" {
		t.Errorf("Serial Code incorrect, expected: 2.1.8, got: %s", GetSerialNumberString());
	}
	if success,verDrv,verDll := DriverMatch(); !success {
		t.Errorf("vJoy driver version: %d, does not match dll version: %d", verDrv, verDll);
	}
}

func TestController(t *testing.T) {
	c := NewVJoyController();
	c.SetAxis(3000, Rx);
	c.SetButton(true, 8);
	c.Reset();
	c.Close();
}
