package credhelper

import (
	"fmt"

	"github.com/VonC/gitcred/internal/syscall"
)

func (ch *credHelper) Set(username, password, host string) error {
	fmt.Println("Set")
	if username == "" {
		return fmt.Errorf("set: --username is mandatory")
	}
	if host == "" {
		return fmt.Errorf("set: --servername is mandatory")
	}
	if password == "" {
		return fmt.Errorf("set: password is mandatory")
	}
	cmd := fmt.Sprintf("printf \"host=%s\\nprotocol=https\\nusername=%s\\npassword=%s\"|\"%s\" store", host, username, password, ch.exe)
	fmt.Println(cmd)
	_, _, err := syscall.ExecCmd(cmd)

	if err != nil {
		return fmt.Errorf("unable to set credential.helper value for Host '%s' and username '%s':\n%w", host, username, err)
	}
	return nil
}
