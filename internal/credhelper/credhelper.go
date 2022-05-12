package credhelper

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/VonC/gitcred/internal/syscall"
)

type credHelper struct {
	creds    creds
	protocol string
	exe      string
}

type creds []*cred

type cred struct {
	servername string
	username   string
}

func NewCredHelper(servername, username string) (*credHelper, error) {
	ch := &credHelper{
		protocol: "https",
		creds:    make(creds, 0),
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

	rootf := filepath.Dir(filepath.Dir(fname))
	if idx := strings.Index(rootf, "mingw64"); idx != -1 {
		rootf = rootf[:idx]
	}
	credHelperFullName := filepath.Join(rootf, "mingw64/libexec/git-core", fmt.Sprintf("git-credential-%s", credHelperName))
	fmt.Println("credHelperFullName = '" + credHelperFullName + "'")

	ch.exe = strings.TrimSpace(credHelperFullName)

	if servername == "" {
		userServerNames, err := getRemoteUserServernames()
		if err != nil {
			return nil, err
		}
		for _, userServername := range userServerNames {
			cred := &cred{}
			hh := strings.Split(userServername, "@")
			cred.servername = hh[0]
			if len(hh) == 2 {
				cred.username = hh[0]
				cred.servername = hh[1]
			}
			ch.creds = ch.creds.append(cred)
		}
	} else {
		cred := &cred{
			servername: servername,
			username:   username,
		}
		ch.creds = ch.creds.append(cred)
	}

	return ch, nil
}

type orderedSet map[string]bool
type remotes orderedSet
type userServernames orderedSet

func getRemoteUserServernames() ([]string, error) {
	res := make([]string, 0)
	remotes, err := getUniqueRemoteNames()
	if err != nil {
		return res, err
	}
	userServernames, err := remotes.getUniqueUserServernames()
	if err != nil {
		return res, err
	}
	return orderedSet(userServernames).set(), nil
}

func getUniqueRemoteNames() (remotes, error) {
	res := make(orderedSet)
	stderr, stdout, err := syscall.ExecCmd("git remote")
	serr := stderr.String()
	if err != nil && serr != "" {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		return nil, fmt.Errorf("unable to get remote list in current folder '%s' (stderr '%s'): %w", pwd, serr, err)
	}
	for _, line := range strings.Split(strings.TrimRight(stdout.String(), "\n"), "\n") {
		res.add(line)
	}
	return remotes(res), nil
}

func (os orderedSet) add(s string) {
	if !os[s] {
		os[s] = true
	}
}

func (os orderedSet) set() []string {
	keys := make([]string, 0, len(os))
	for k := range os {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func (remotes remotes) getUniqueUserServernames() (userServernames, error) {
	res := make(orderedSet)
	for remote := range remotes {
		userServernames, err := getUniqueUserServernameFromRemote(remote)
		if err != nil {
			return nil, fmt.Errorf("unable to get user@server name from remote name '%s'", remote)
		}
		for userServername := range userServernames {
			res.add(userServername)
		}
	}
	return userServernames(res), nil
}

func getUniqueUserServernameFromRemote(remote string) (userServernames, error) {
	res := make(orderedSet)
	stderr, stdout, err := syscall.ExecCmd("git remote get-url --all " + remote)
	serr := stderr.String()
	if err != nil && serr != "" {
		return nil, fmt.Errorf("unable to get remote URL list from remote '%s' (stderr '%s'): %w", remote, serr, err)
	}
	for _, line := range strings.Split(strings.TrimRight(stdout.String(), "\n"), "\n") {
		// https://stackoverflow.com/questions/45537134/parse-a-url-with-in-go
		u, err := url.Parse(line)
		if err != nil {
			continue
		}
		if u.Scheme != "https" {
			continue
		}
		// spew.Dump(u.Path)
		userServername := u.Hostname()
		if u.User.Username() != "" {
			userServername = u.User.Username() + "@" + userServername
		}
		res.add(userServername)
	}
	return userServernames(res), nil
}

func (creds creds) append(cred *cred) creds {
	if cred.servername == "" {
		return creds
	}
	for _, acred := range creds {
		if acred.servername == cred.servername && acred.username == cred.username {
			return creds
		}
	}
	return append(creds, cred)
}
