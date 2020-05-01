// +build linux

package command

import (
	"bytes"
	"os/exec"
)

// RunSystemLoad Возвращает строку с информацией о загрузке системы
// выполняемая команда: top -b -n1 | grep load
func RunSystemLoad() (int, string, string) {
	var output, stderr bytes.Buffer

	c1 := exec.Command("top", "-b", "-n1") // top -b -n1
	c2 := exec.Command("grep", "load")
	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = &output
	c2.Stderr = &stderr

	if err := c2.Start(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), output.String(), stderr.String()
		}
	}

	if err := c1.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), output.String(), stderr.String()
		}
	}

	if err := c2.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), output.String(), stderr.String()
		}
	}

	return 0, output.String(), stderr.String()
}

// RunLoadCPU Возвращает строку с информацией о загрузке системы
// выполняемая команда: top -b -n1 | grep Cpu
func RunLoadCPU() (int, string, string) {
	var output, stderr bytes.Buffer

	c1 := exec.Command("top", "-b", "-n1") // top -b -n1
	c2 := exec.Command("grep", "Cpu")
	c2.Stdin, _ = c1.StdoutPipe()
	c2.Stdout = &output
	c2.Stderr = &stderr

	if err := c2.Start(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), output.String(), stderr.String()
		}
	}

	if err := c1.Run(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), output.String(), stderr.String()
		}
	}

	if err := c2.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			return exitError.ExitCode(), output.String(), stderr.String()
		}
	}

	return 0, output.String(), stderr.String()
}
