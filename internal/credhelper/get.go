package credhelper

import (
	"fmt"
	"strings"

	"github.com/VonC/gitcred/internal/syscall"
)

func (ch *credHelper) Get(username string) (string, error) {
	fmt.Printf("Get")
	res := ""
	for _, cred := range ch.creds {
		res = res + "\n" + username + "@" + cred.host + ":\n"
		u := ""
		if username != "" {
			u = "\\nusername=" + username
		}
		cmd := fmt.Sprintf("printf \"host=%s\\nprotocol=%s%s\"|\"%s\" get", cred.host, ch.protocol, u, ch.exe)
		//fmt.Println(cmd)
		_, stdout, err := syscall.ExecCmd(cmd)

		if err != nil {
			return "", fmt.Errorf("unable to get credential.helper value for Host '%s@%s':\n%w", username, cred.host, err)
		}
		res = res + strings.TrimSpace(stdout.String())
	}
	return res, nil
}
