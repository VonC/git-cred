package credhelper

import (
	"fmt"

	"github.com/VonC/gitcred/internal/syscall"
)

func (ch *credHelper) Erase(username, host string) error {
	fmt.Println("Erase")
	cmd := fmt.Sprintf("printf \"host=%s\\nprotocol=https\\nusername=%s\"|\"%s\" erase", host, username, ch.exe)
	//fmt.Println(cmd)
	_, _, err := syscall.ExecCmd(cmd)

	if err != nil {
		return fmt.Errorf("unable to erase credential.helper value for Host '%s' and username '%s':\n%w", host, username, err)
	}
	return nil
}
