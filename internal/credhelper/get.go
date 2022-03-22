package credhelper

import (
	"fmt"
	"strings"

	"github.com/VonC/gitcred/internal/syscall"
)

func (ch *credHelper) Get(username, servername string) (string, error) {
	fmt.Printf("Get user='%s', server='%s', %d creds", username, servername, len(ch.creds))
	res := ""
	if servername != "" {
		return ch.getus(username, servername)
	}
	for _, cred := range ch.creds {
		if username == "" {
			username = cred.username
		}
		res = res + "\n\n" + username + "@" + cred.servername + ":\n"
		ares, err := ch.getus(username, cred.servername)
		if err != nil {
			return res, err
		}
		res = res + ares
	}
	return strings.TrimSpace(res), nil
}

func (ch *credHelper) getus(username, servername string) (string, error) {
	u := ""
	if username != "" {
		u = "\\nusername=" + username
	}
	cmd := fmt.Sprintf("printf \"host=%s\\nprotocol=%s%s\"|\"%s\" get", servername, ch.protocol, u, ch.exe)
	res := "\n" + cmd
	_, stdout, err := syscall.ExecCmd(cmd)

	if err != nil {
		return "", fmt.Errorf("unable to get credential.helper value for Host '%s@%s':\n%w", username, servername, err)
	}
	return res + "\n" + strings.TrimSpace(stdout.String()), nil
}
