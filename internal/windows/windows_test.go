package windows

import (
	"testing"
)

func TestListProcesses(t *testing.T) {
	ps, err := ListProcesses()
	if err != nil {
		t.Errorf("failed to list processes. processes: %+v err: %s\n", ps, err)
	}
}

func TestExec(t *testing.T) {
	c := "cmd /k echo test"
	out, err := Exec(c)
	if err != nil {
		t.Errorf("failed to execute command. command: %s, err: %s\n", c, err)
	}
	if out == "" {
		t.Errorf("execute command response null. command: %s, response: %s", c, out)
	}
}