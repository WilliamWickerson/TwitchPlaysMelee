package scripting

import (
	"testing"
	"controller"
	"scripting/AST"
)

func TestButtonCommand(t *testing.T) {
	parser := NewParser(NewScanner("press a"));
	gcArray := make([]controller.GamecubeController, AST.MAXFRAMEWINDOW);
	for _,command := range parser.Parse() {
		command.Execute(gcArray);
	}
	if !gcArray[0].A {
		t.Errorf("Error, expected A button to be: true, but got: %t", gcArray[0].A);
	}
}
