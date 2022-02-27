package credhelper

import (
	"fmt"
	"gitcred/internal/syscall"
)

func (ch *credHelper) Erase(username string) error {
	fmt.Println("Erase")
	ch.username = username
	cmd := fmt.Sprintf("printf \"host=%s\\nprotocol=https\\nusername=%s\"|\"%s\" erase", ch.host, ch.username, ch.exe)
	//fmt.Println(cmd)
	_, _, err := syscall.ExecCmd(cmd)

	if err != nil {
		return fmt.Errorf("unable to erase credential.helper value for Host '%s' and username '%s':\n%w", ch.host, ch.username, err)
	}
	return nil
}
