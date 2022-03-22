package credhelper

import (
	"fmt"

	"github.com/VonC/gitcred/internal/syscall"
)

func (ch *credHelper) Erase(username, servername string) error {
	fmt.Println("Erase")
	if username == "" {
		return fmt.Errorf("erase: --username is mandatory")
	}
	if servername == "" {
		return fmt.Errorf("erase: --servername is mandatory")
	}
	cmd := fmt.Sprintf("printf \"host=%s\\nprotocol=https\\nusername=%s\"|\"%s\" erase", servername, username, ch.exe)
	fmt.Println(cmd)
	_, _, err := syscall.ExecCmd(cmd)

	if err != nil {
		return fmt.Errorf("unable to erase credential.helper value for Host '%s' and username '%s':\n%w", servername, username, err)
	}
	return nil
}
