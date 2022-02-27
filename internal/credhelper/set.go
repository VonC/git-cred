package credhelper

import (
	"fmt"
	"gitcred/internal/syscall"
)

func (ch *credHelper) Set(user, password string) error {
	fmt.Println("Set")
	ch.username = user
	ch.password = password
	cmd := fmt.Sprintf("printf \"host=%s\\nprotocol=https\\nusername=%s\\npassword=%s\"|\"%s\" store", ch.host, ch.username, ch.password, ch.exe)
	fmt.Println(cmd)
	_, _, err := syscall.ExecCmd(cmd)

	if err != nil {
		return fmt.Errorf("unable to set credential.helper value for Host '%s' and username '%s':\n%w", ch.host, ch.username, err)
	}
	return nil
}
