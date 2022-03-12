# gitcred

Reads cached credentials using the configured Git credential helper.  
<sup>([MIT license](LICENSE.md))</sup>

    go install github.com/VonC/gitcred@latest

## Problem

A git credential helper is used to cache HTTPS remote Git hosting service credentials (username/password or token)

You can used that helper to query credentials manually:

```bash
credhelper=$(git config credential.helper)
printf "host=github.com\nprotocol=https" | git-credential-${credhelper} get
```

This works for any Mac/Linux/Windows password.

And if you want to set a new password/token, it is even more cumbersome:

```bash
credhelper=$(git config credential.helper)
printf "host=github.com\nprotocol=https\nusername=VonC\npassword=xxx" | git-credential-${credhelper} set
```

## Goal

Replace the complex command line by a tool able to quickly read/set/erase cached credentials, no matter your credential helper.

Cross-platform.

## Solution

- `gitcred` will read your current credential helper
- By default, in a repository, it will display cached credentials for the current folder/repository

## get

`get` is the default command for `gitcred`.  
You do not need to add `get`.

### get, from outside a repository

```bash
gitcred -u VonC -s github.com
# or (same)
gitcred -u VonC -s github.com get
```

### get, from inside a cloned repository folder

```bash
cd /path/to/local/github.com/cloned/repository
gitcred
# or (same)
gitcred get
```

## set

### set, from outside a repository

```bash
gitcred -u VonC -s github.com set <password or token>
```

### set, from inside a cloned repository folder

```bash
cd /path/to/local/github.com/cloned/repository
gitcred -u VonC set <password or token>
```

## erase

### erase, from outside a repository

```bash
gitcred -u VonC -s github.com erase
```

### erase, from inside a cloned repository folder

```bash
cd /path/to/local/github.com/cloned/repository
gitcred -u VonC erase
```
