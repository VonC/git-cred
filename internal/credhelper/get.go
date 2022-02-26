package credhelper

import (
	"fmt"
	"gitcred/internal/syscall"
)

func (ch *credHelper) Get() (string, error) {
	fmt.Println("Get")
	cmd := fmt.Sprintf("printf \"host=%s\\nprotocol=https\"|\"%s\" get", ch.host, ch.exe)
	//fmt.Println(cmd)
	_, stdout, err := syscall.ExecCmd(cmd)

	if err != nil {
		return "", fmt.Errorf("unable to get credential.helper value for Host '%s':\n%w", ch.host, err)
	}
	return stdout.String(), nil
}
