//go:build windows
// +build windows

package syscall

import (
	"bytes"
	"fmt"
	"minimalcli/internal/logger"
	"os/exec"
	"strings"
	"syscall"
)

// ExecCmd starts a sh -c 'scmd' session.
// If scmd ends with &, don't wait for result (background process)
func ExecCmd(scmd string) (berr *bytes.Buffer, bout *bytes.Buffer, err error) {
	berr = &bytes.Buffer{}
	bout = &bytes.Buffer{}
	cmd := exec.Command("cmd", "/C", scmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.CmdLine = "cmd /C " + scmd
	//fmt.Printf("Execute '%s'\n%+v\n", scmd, cmd.SysProcAttr.CmdLine)
	logger.Debug("Execute '%s'\n", scmd)
	//fmt.Printf("Execute '%s'\n", cmd.SysProcAttr.CmdLine)
	cmd.Stderr = berr
	cmd.Stdout = bout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		//log.Printf("stdin error %s [%s]", err, berr.String())
		err = fmt.Errorf("stderr '%s':\n==> Lead to stdin:%w", berr.String(), err)
		return berr, bout, err
	}
	stdin.Close()
	err = cmd.Start()
	if err != nil {
		//log.Printf("start error %s [%s]", err, berr.String())
		err = fmt.Errorf("stderr '%s':\n==> Lead to start:%w", berr.String(), err)
		return berr, bout, err
	}
	if strings.HasPrefix(scmd, "start ") && strings.Contains(scmd, " /B ") {
		return berr, bout, nil
	}
	err = cmd.Wait()
	if err != nil {
		//log.Printf("Exit error %s [%s]", err, berr.String())
		err = fmt.Errorf("stderr '%s':\n==> Lead to:%w", berr.String(), err)
	}
	return berr, bout, err
}
