package credhelper

import (
	"fmt"
	"gitcred/internal/syscall"
)

func (ch *credHelper) Get(username string) (string, error) {
	fmt.Println("Get")
	ch.username = username
	u := ""
	if ch.username != "" {
		u = "\\nusername=" + ch.username
	}
	cmd := fmt.Sprintf("printf \"host=%s\\nprotocol=https%s\"|\"%s\" get", ch.host, u, ch.exe)
	//fmt.Println(cmd)
	_, stdout, err := syscall.ExecCmd(cmd)

	if err != nil {
		return "", fmt.Errorf("unable to get credential.helper value for Host '%s':\n%w", ch.host, err)
	}
	return stdout.String(), nil
}
