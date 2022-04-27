package credhelper

import (
	"fmt"
	"strings"

	"github.com/VonC/gitcred/internal/syscall"
)

func (ch *credHelper) Get(username, servername string) (string, error) {
	fmt.Printf("Get user='%s', server='%s', %d creds\n\n", username, servername, len(ch.creds))
	res := ""
	if servername != "" {
		return ch.getus(username, servername)
	}
	aUsername := ""
	for _, cred := range ch.creds {
		//spew.Dump(cred)
		aUsername = username
		if username == "" {
			aUsername = cred.username
		}
		res = res + "\n\n" + aUsername + "@" + cred.servername + ":\n"
		ares, err := ch.getus(aUsername, cred.servername)
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
	// fmt.Println(cmd)
	_, stdout, err := syscall.ExecCmd(cmd)

	if err != nil {
		if strings.Contains(err.Error(), "Cannot prompt because user interactivity has been disabled") {
			return res + "\n" + "<no credentials registered>", nil
		}
		return "", fmt.Errorf("unable to get credential.helper value for Host '%s@%s':\n%w", username, servername, err)
	}
	return res + "\n" + strings.TrimSpace(stdout.String()), nil
}
