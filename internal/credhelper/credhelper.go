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
	host     string
	username string
}

func NewCredHelper(host, username string) (*credHelper, error) {
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

	credHelperFullName := filepath.Join(filepath.Dir(filepath.Dir(fname)), "mingw64/libexec/git-core", fmt.Sprintf("git-credential-%s", credHelperName))
	// fmt.Println(credHelperFullName)

	ch.exe = strings.TrimSpace(credHelperFullName)

	if host == "" {
		cred := &cred{}
		hosts, err := getLocalHosts()
		if err != nil {
			return nil, err
		}
		for _, host := range hosts {
			hh := strings.Split(host, "@")
			cred.host = hh[0]
			if len(hh) == 2 {
				cred.username = hh[0]
				cred.host = hh[1]
			}
			ch.creds = ch.creds.append(cred)
		}
	} else if username != "" {
		cred := &cred{
			host:     host,
			username: username,
		}
		ch.creds = ch.creds.append(cred)
	}

	return ch, nil
}

type orderedSet map[string]bool
type remotes orderedSet
type hosts orderedSet

func getLocalHosts() ([]string, error) {
	res := make([]string, 0)
	remotes, err := getUniqueRemoteNames()
	if err != nil {
		return res, err
	}
	hosts, err := remotes.getUniqueHosts()
	if err != nil {
		return res, err
	}
	return orderedSet(hosts).set(), nil
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

func (remotes remotes) getUniqueHosts() (hosts, error) {
	res := make(orderedSet)
	for remote := range remotes {
		hosts, err := getUniqueHostsFromRemote(remote)
		if err != nil {
			return nil, fmt.Errorf("unable to get host from remote name '%s'", remote)
		}
		for host := range hosts {
			res.add(host)
		}
	}
	return hosts(res), nil
}

func getUniqueHostsFromRemote(remote string) (hosts, error) {
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
		host := u.Hostname()
		if u.User.Username() != "" {
			host = u.User.Username() + "@" + host
		} else {
			pp := strings.Split(u.Path, "/")
			if len(pp) < 2 {
				continue // for instance U:\ or git@server:...
			}
			host = pp[1] + "@" + host
		}
		res.add(host)
	}
	return hosts(res), nil
}

func (ch *credHelper) Host() string {
	if len(ch.creds) == 1 {
		return ch.creds[0].host
	}
	return ""
}

func (ch *credHelper) User() string {
	if len(ch.creds) == 1 {
		return ch.creds[0].username
	}
	return ""
}

func (creds creds) append(cred *cred) creds {
	if cred.host == "" || cred.username == "" {
		return creds
	}
	for _, acred := range creds {
		if acred.host == cred.host && acred.username == cred.username {
			return creds
		}
	}
	return append(creds, cred)
}
