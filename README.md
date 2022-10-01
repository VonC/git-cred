# git-cred

Reads cached credentials using the configured Git credential helper.  
<sup>([MIT license](LICENSE.md))</sup>

    go install github.com/VonC/git-cred@latest

[![Open in VS Code](https://img.shields.io/static/v1?logo=visualstudiocode&label=&message=Open%20in%20Visual%20Studio%20Code&labelColor=2c2c32&color=007acc&logoColor=007acc)](https://vscode.dev/github/VonC/git-cred)

## Problem

A git credential helper is used to cache HTTPS remote Git hosting service credentials (username/password or token)

You can used that helper to query credentials manually:

```bash
credhelper=$(git config credential.helper)
printf "host=github.com\nprotocol=https" | git-credential-${credhelper} get
```

This works for any Mac/Linux/Windows cached credentials.

And if you want to set a new password/token, it is even more cumbersome:

```bash
credhelper=$(git config credential.helper)
printf "host=github.com\nprotocol=https\nusername=VonC\npassword=xxx" | git-credential-${credhelper} set
```

## Goal

- Replace the complex command line by a tool able to quickly read/set/erase cached credentials, no matter your credential helper.
- Cross-platform.

## Solution

- `git-cred` will read your current credential helper
- By default, in a repository, it will display cached credentials for the current folder/repository

Since the executable follows the naming convention `git-xxx` (here `git-cred` or `git-cred.exe`), that means you can also type:  
`git cred`.  
As if "`cred`" was a `git` command. It works if the executable `git-cred`(`.exe`) is in your `$PATH`/`%PATH%`.

## get

`get` is the default command for `git-cred`.  
You do not need to add `get`.

### get, from outside a repository

```bash
git cred -u VonC -s github.com
# or (same)
git cred -u VonC -s github.com get
```

### get, from inside a cloned repository folder

```bash
cd /path/to/local/github.com/cloned/repository
git cred
# or (same)
git cred get
```

## set

### set, from outside a repository

```bash
git cred -u VonC -s github.com set <password or token>
```

### set, from inside a cloned repository folder

```bash
cd /path/to/local/github.com/cloned/repository
git cred -u VonC set <password or token>
```

## erase

### erase, from outside a repository

```bash
git cred -u VonC -s github.com erase
```

### erase, from inside a cloned repository folder

```bash
cd /path/to/local/github.com/cloned/repository
git cred -u VonC erase
```
