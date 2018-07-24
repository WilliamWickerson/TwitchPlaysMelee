package scripting

import (
	"testing"
	"controller"
	"scripting/AST"
	"math"
)

func TestButtonPress(t *testing.T) {
	parser := NewParser(NewScanner(`press a; press b 2; press x 3; press y 4; press z 5;
		press start 6; press dleft 7; press dright 8; press ddown 9; press dup 10`));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	for _,command := range parser.Parse() {
		command.Execute(gcArray);
	}
	if !gcArray[0].A {
		t.Errorf("Error, expected A button to be: true, but got: %t", gcArray[0].A);
	}
	if !gcArray[1].B {
		t.Errorf("Error, expected B button to be: true, but got: %t", gcArray[1].B);
	}
	if !gcArray[2].X {
		t.Errorf("Error, expected X button to be: true, but got: %t", gcArray[2].X);
	}
	if !gcArray[3].Y {
		t.Errorf("Error, expected Y button to be: true, but got: %t", gcArray[3].Y);
	}
	if !gcArray[4].Z {
		t.Errorf("Error, expected Z button to be: true, but got: %t", gcArray[4].Z);
	}
	if !gcArray[5].START {
		t.Errorf("Error, expected START button to be: true, but got: %t", gcArray[5].START);
	}
	if !gcArray[6].DLEFT {
		t.Errorf("Error, expected DLEFT button to be: true, but got: %t", gcArray[6].DLEFT);
	}
	if !gcArray[7].DRIGHT {
		t.Errorf("Error, expected DRIGHT button to be: true, but got: %t", gcArray[7].DRIGHT);
	}
	if !gcArray[8].DDOWN {
		t.Errorf("Error, expected DDOWN button to be: true, but got: %t", gcArray[8].DDOWN);
	}
	if !gcArray[9].DUP {
		t.Errorf("Error, expected DUP button to be: true, but got: %t", gcArray[9].DUP);
	}
}

func TestButtonUnpress(t *testing.T) {
	parser := NewParser(NewScanner("press A 13"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	for _,command := range parser.Parse() {
		command.Execute(gcArray);
	}
	if !gcArray[12].A {
		t.Errorf("Error, expected A button to be: true, but got: %t", gcArray[0].A);
	}
	parser = NewParser(NewScanner("unpress A 13"));
	for _,command := range parser.Parse() {
		command.Execute(gcArray);
	}
	if gcArray[12].A {
		t.Errorf("Error, expected A button to be: false, but got: %t", gcArray[0].A);
	}
}

func TestSliderCommand(t *testing.T) {
	parser := NewParser(NewScanner("press L 7; press R .6 8"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	for _,command := range parser.Parse() {
		command.Execute(gcArray);
	}
	if !gcArray[6].L || gcArray[6].LANA != 1 {
		t.Errorf("Error, expected L to be fully pressed, but got: %f", gcArray[6].LANA);
	}
	if gcArray[7].R || gcArray[7].RANA != .6 {
		t.Errorf("Error, expected R to be pressed: .6, but got: %f", gcArray[7].RANA);
	}
	parser = NewParser(NewScanner("unpress L 7"));
	for _,command := range parser.Parse() {
		command.Execute(gcArray);
	}
	if gcArray[6].L || gcArray[6].LANA != 0 {
		t.Errorf("Error, expected L to be unpressed, but got: %f", gcArray[6].LANA);
	}
}

func TestStickPairCommand(t *testing.T) {
	parser := NewParser(NewScanner("cstick (.5,-.5) 2"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	for _,command := range parser.Parse() {
		command.Execute(gcArray);
	}
	eps := .0000001;
	if math.Abs(gcArray[1].CX - .5) > eps {
		t.Errorf("Error, expected stick to be: %f, but got: %f", .5, gcArray[1].CX);
	}
	if math.Abs(gcArray[1].CX - .5) > eps {
		t.Errorf("Error, expected stick to be: %f, but got: %f", -.5, gcArray[1].CY);
	}
}

func TestStickNamedDirections(t *testing.T) {
	parser := NewParser(NewScanner("stick left right up down left up 33"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	for _,command := range parser.Parse() {
		command.Execute(gcArray);
	}
	eps := .0000001;
	if math.Abs(gcArray[32].STICKX - -math.Sqrt2 / 2) > eps {
		t.Errorf("Error, expected stick to be: %f, but got: %f", math.Sqrt2 / 2, gcArray[32].STICKX);
	}
	if math.Abs(gcArray[32].STICKY - math.Sqrt2 / 2) > eps {
		t.Errorf("Error, expected stick to be: %f, but got: %f", math.Sqrt2 / 2, gcArray[32].STICKY);
	}
}

func TestDuration(t *testing.T) {
	parser := NewParser(NewScanner("press A 1-10,21-30,41-50"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	for _,command := range parser.Parse() {
		command.Execute(gcArray);
	}
	for i := 1; i <= 10; i++ {
		for j := i; j <= 50; j += 20 {
			if !gcArray[j-1].A {
				t.Errorf("Error, expected A button to be: true, on frame %d, but got: %t", i, gcArray[i-1].A);
			}
		}
	}
	for i := 11; i <= 20; i++ {
		for j := i; j <= 60; j += 20 {
			if gcArray[j-1].A {
				t.Errorf("Error, expected A button to be: false, on frame %d, but got: %t", i, gcArray[i-1].A);
			}
		}
	}
}

func TestMacroCommand(t *testing.T) {
	parser := NewParser(NewScanner("wavedash 3 left 2"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	for _,command := range parser.Parse() {
		command.Execute(gcArray);
	}
	if !gcArray[1].X {
		t.Errorf("Error, expected X button to be: true, but got: %t", gcArray[1].X);
	}
	if !gcArray[4].L {
		t.Errorf("Error, expected L button to be: true, but got: %t", gcArray[4].L);
	}
	eps := .0000001;
	if math.Abs(gcArray[4].STICKX - -math.Sqrt2 / 2) > eps {
		t.Errorf("Error, expected stick to be: %f, but got: %f", math.Sqrt2 / 2, gcArray[4].STICKX);
	}
	if math.Abs(gcArray[4].STICKY - -math.Sqrt2 / 2) > eps {
		t.Errorf("Error, expected stick to be: %f, but got: %f", math.Sqrt2 / 2, gcArray[4].STICKY);
	}
}

func TestMacroCommandDelay(t *testing.T) {
	parser := NewParser(NewScanner("stick up; grab 2"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	commands := parser.Parse();
	if len(commands) != 2 {
		t.Errorf("Error, expected there to be 2 commands, but got: %d", len(commands));
	}
	for _,command := range commands {
		command.Execute(gcArray);
	}
	eps := .0000001;
	if math.Abs(gcArray[0].STICKY - 1) > eps {
		t.Errorf("Error, expected stick to be: %f, but got: %f", 1, gcArray[0].STICKY);
	}
	if gcArray[0].Z {
		t.Errorf("Error, expected Z button to be: false, but got: %t", gcArray[0].Z);
	}
	if !gcArray[1].Z {
		t.Errorf("Error, expected Z button to be: true, but got: %t", gcArray[1].Z);
	}
}

func TestBadCommand(t *testing.T) {
	parser := NewParser(NewScanner("7 2 3 4; press A"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	commands := parser.Parse();
	if len(commands) != 1 {
		t.Errorf("Error, expected only 1 valid command, but got: %d", len(commands));
	}
	for _,command := range commands {
		command.Execute(gcArray);
	}
	if !gcArray[0].A {
		t.Errorf("Error, expected A button to be: true, but got: %t", gcArray[0].A);
	}
}

func TestBadButtonCommand(t *testing.T) {
	parser := NewParser(NewScanner("Press Stick; press A"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	commands := parser.Parse();
	if len(commands) != 1 {
		t.Errorf("Error, expected only 1 valid command, but got: %d", len(commands));
	}
	for _,command := range commands {
		command.Execute(gcArray);
	}
	if !gcArray[0].A {
		t.Errorf("Error, expected A button to be: true, but got: %t", gcArray[0].A);
	}
}

func TestBadStickCommand(t *testing.T) {
	parser := NewParser(NewScanner("stick 3; press A"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	commands := parser.Parse();
	if len(commands) != 1 {
		t.Errorf("Error, expected only 1 valid command, but got: %d", len(commands));
	}
	for _,command := range commands {
		command.Execute(gcArray);
	}
	if !gcArray[0].A {
		t.Errorf("Error, expected A button to be: true, but got: %t", gcArray[0].A);
	}
}

func TestBadDuration(t *testing.T) {
	parser := NewParser(NewScanner("press X 1.3; press A"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	commands := parser.Parse();
	if len(commands) != 1 {
		t.Errorf("Error, expected only 1 valid command, but got: %d", len(commands));
	}
	for _,command := range commands {
		command.Execute(gcArray);
	}
	if !gcArray[0].A {
		t.Errorf("Error, expected A button to be: true, but got: %t", gcArray[0].A);
	}
}

func TestBadNamedDirection(t *testing.T) {
	parser := NewParser(NewScanner("stick left-right; press A"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	commands := parser.Parse();
	if len(commands) != 1 {
		t.Errorf("Error, expected only 1 valid command, but got: %d", len(commands));
	}
	for _,command := range commands {
		command.Execute(gcArray);
	}
	if !gcArray[0].A {
		t.Errorf("Error, expected A button to be: true, but got: %t", gcArray[0].A);
	}
}

func TestBadPairDirection(t * testing.T) {
	parser := NewParser(NewScanner("stick (3 4); press A"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	commands := parser.Parse();
	if len(commands) != 1 {
		t.Errorf("Error, expected only 1 valid command, but got: %d", len(commands));
	}
	for _,command := range commands {
		command.Execute(gcArray);
	}
	if !gcArray[0].A {
		t.Errorf("Error, expected A button to be: true, but got: %t", gcArray[0].A);
	}
}
