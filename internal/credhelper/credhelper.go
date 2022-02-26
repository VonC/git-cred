package credhelper

import (
	"fmt"
	"gitcred/internal/syscall"
	"os/exec"
	"path/filepath"
	"strings"
)

type credHelper struct {
	host     string
	username string
	protocol string
	password string
	exe      string
}

func NewCredHelper(host string) (*credHelper, error) {
	ch := &credHelper{
		host: host,
	}

	stderr, stdout, err := syscall.ExecCmd("git config --global credential.helper")
	serr := stderr.String()
	if err != nil && serr != "" {
		return nil, fmt.Errorf("unable to get global credential.helper Git config (stderr '%s'): %w", serr, err)
	}
	if err != nil {
		stderr, stdout, err = syscall.ExecCmd("git config --system credential.helper")
		serr := stderr.String()
		if err != nil {
			return nil, fmt.Errorf("unable to get system credential.helper Git config (stderr '%s'): %w", serr, err)
		}
	}
	credHelperName := stdout.String()
	//fmt.Println(credHelperName)

	fname, err := exec.LookPath("git")
	if err == nil {
		fname, err = filepath.Abs(fname)
	}
	if err != nil {
		return nil, fmt.Errorf("unable to get Git path (stderr '%s'): %w", serr, err)
	}

	credHelperFullName := filepath.Join(filepath.Dir(filepath.Dir(fname)), "mingw64/libexec/git-core", fmt.Sprintf("git-credential-%s", credHelperName))
	// fmt.Println(credHelperFullName)

	ch.exe = strings.TrimSpace(credHelperFullName)

	return ch, nil
}
