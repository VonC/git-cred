package credhelper

import (
	"fmt"
	"regexp"

	"github.com/VonC/git-cred/internal/syscall"
)

// https://regex101.com/r/VbVbfs/1
var rePassword = regexp.MustCompile(`(?m)password=[^"]*?"`)

func (ch *credHelper) Set(username, password, servername string) error {
	fmt.Println("Set")
	if username == "" {
		return fmt.Errorf("set: --username is mandatory")
	}
	if servername == "" {
		return fmt.Errorf("set: --servername is mandatory")
	}
	if password == "" {
		return fmt.Errorf("set: password is mandatory")
	}
	obfuscatedCmd := fmt.Sprintf("printf \"host=%s\\nprotocol=https\\nusername=%s\\npassword=xxxx\"|\"%s\" store", servername, username, ch.exe)
	cmd := rePassword.ReplaceAllString(obfuscatedCmd, fmt.Sprintf(`password=%s"`, password))
	fmt.Println(obfuscatedCmd)
	_, _, err := syscall.ExecCmd(cmd)

	if err != nil {
		return fmt.Errorf("unable to set credential.helper value for Host '%s' and username '%s':\n%w", servername, username, err)
	}
	return nil
}
