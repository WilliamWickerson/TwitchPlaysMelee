package vJoy

// #cgo CFLAGS: -I./vJoySdk
// #cgo LDFLAGS: -L./vJoySdk -lvJoyInterface
// #include "windows.h"
// #include "public.h"
// #include "vJoyInterface.h"
import "C"

import (
	"unsafe"
)

type VJoyController interface {
	SetAxis(int,  Axis);
	SetButton(bool, int);
	Reset();
	Close();
}

type vJoyController struct {
	handle C.uint;
}

var (
	//Initialize so first controller gets value 1
	nextHandle C.uint = 1;
)

/*
Axes are implemented using #define in the header.
using an enumeration with those values works well
and avoids a switch statement later
*/
type Axis int;

const (
	X = C.HID_USAGE_X;
	Y = C.HID_USAGE_Y;
	Z = C.HID_USAGE_Z;
	Rx = C.HID_USAGE_RX;
	Ry = C.HID_USAGE_RY;
	Rz = C.HID_USAGE_RZ;
	Sl0 = C.HID_USAGE_SL0;
	Sl1 = C.HID_USAGE_SL1;
)

func NewVJoyController() VJoyController {
	//Acquire the next controller and create object
	C.AcquireVJD(C.uint(nextHandle));
	c := vJoyController{nextHandle};
	//Increment value for subsequent call
	nextHandle++;
	return c;
}

func (c vJoyController) SetAxis(value int, axis Axis) {
	//Need to convert to required types
	C.SetAxis(C.long(value), c.handle, C.uint(axis));
}

func (c vJoyController) SetButton(press bool, button int) {
	//Convert boolean to int since c doesn't have a built in bool type
	var intPress C.int;
	if press {
		intPress = 1;
	} else {
		intPress = 0;
	}
	//Convert button to character as required
	C.SetBtn(intPress, c.handle, C.uchar(button));
}

func (c vJoyController) Reset() {
	//For some reason resetVJD(handle) doesn't work so reset buttons then axes
	C.ResetButtons(c.handle);
	//Put all axes to center
	c.SetAxis(16384, X);
	c.SetAxis(16384, Y);
	c.SetAxis(16384, Z);
	c.SetAxis(16384, Rx);
	c.SetAxis(16384, Ry);
	c.SetAxis(16384, Rz);
	//For some reason <16384 reads negative values? so set them to 16384
	c.SetAxis(16384, Sl0);
	c.SetAxis(16384, Sl1);
}

func (c vJoyController) Close() {
	C.RelinquishVJD(c.handle);
}

func Enabled() bool {
	return C.vJoyEnabled() != 0;
}

//Converts wchar_t string to GoString
func convertwchar_t(wcharPtr *C.wchar_t) string {
	//Create space for char string plus '\0' and make sure it gets deleted
	charPtr := (*C.char)(unsafe.Pointer(C.malloc(C.wcslen(wcharPtr) + 1)));
	defer C.free(unsafe.Pointer(charPtr));
	//Convert the wchar_t string to char string
	C.wcstombs(charPtr, wcharPtr, C.wcslen(wcharPtr) + 1);
	//Convert the char string to a GoString
	return C.GoString(charPtr);
}

func GetManufacturerString() string {
	wcharPtr := (*C.wchar_t)(unsafe.Pointer(C.GetvJoyManufacturerString()));
	return convertwchar_t(wcharPtr);
}

func GetProductString() string {
	wcharPtr := (*C.wchar_t)(unsafe.Pointer(C.GetvJoyProductString()));
	return convertwchar_t(wcharPtr);
}

func GetSerialNumberString() string {
	wcharPtr := (*C.wchar_t)(unsafe.Pointer(C.GetvJoySerialNumberString()));
	return convertwchar_t(wcharPtr);
}

func DriverMatch() (bool,int,int) {
	//Get the versions and return whether they match along with the version numbers
	var verDrv, verDll C.WORD;
	success := C.DriverMatch(&verDrv, &verDll) != 0;
	return success, int(verDrv), int(verDll);
}