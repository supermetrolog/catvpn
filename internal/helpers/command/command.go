package command

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"
)

func RunCommand(cmd string) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer

	c := exec.Command("bash", "-c", cmd)
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()

	logrus.Debugf("Runned command: %s", cmd)

	if err != nil {
		return stderr.String(), fmt.Errorf("command error: %v", err)
	}

	return stdout.String(), err
}
