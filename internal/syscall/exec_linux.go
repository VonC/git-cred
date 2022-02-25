// +build linux

package syscall

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)

// ExecCmd starts a sh -c 'scmd' session.
// If scmd ends with &, don't wait for result (background process)
func ExecCmd(scmd string) (berr *bytes.Buffer, bout *bytes.Buffer, err error) {
	berr = &bytes.Buffer{}
	bout = &bytes.Buffer{}
	cmd := exec.Command("sh", "-c", scmd)
	//cmd.SysProcAttr = &syscall.SysProcAttr{}
	//cmd.SysProcAttr.CmdLine = "cmd /C " + scmd
	//fmt.Printf("Execute '%s'\n%+v\n", scmd, cmd.SysProcAttr.CmdLine)
	//fmt.Printf("Execute '%s'\n", scmd)
	cmd.Stderr = berr
	cmd.Stdout = bout
	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("stdin error %s [%s]", err, berr.String())
		return berr, bout, err
	}
	stdin.Close()
	err = cmd.Start()
	if err != nil {
		log.Printf("start error %s [%s]", err, berr.String())
		return berr, bout, err
	}
	if strings.HasSuffix(scmd, "&") {
		return berr, bout, nil
	}
	err = cmd.Wait()
	if err != nil {
		log.Printf("exit error %s [%s]", err, berr.String())
	}
	// log.Printf("sout #%s# [%s]", bout.String(), berr.String())
	return berr, bout, err
}
