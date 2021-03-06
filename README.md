# git-nfs

Make git repo as a nfs server storage backend

## Features

* NFS v3
* In memory filesystem to cache files
* Sync files to git repository every 5s (option to set)

## Getting Started

* build and run

```shell
git clone github.com/iineva/git-nfs
cd git-nfs/cmd/gitnfs/
# change port number whatever you want
go run gitnfs.go -d -a ":1234" https://github.com/iineva/git-nfs
```

* nfs client mount option

```shell
mkdir /tmp/test_gitnfs
mount -o "port=1234,mountport=1234,intr,noresvport,nolock,noacl" -t nfs localhost:/ /tmp/test_gitnfs
ls /tmp/test_gitnfs
# umount /tmp/test_gitnfs
```

## Usage

```
Usage: gitnfs [options] <YOUR_GIT_REPO_URL>

Options:
  -K string
        private key password
  -a string
        nfs listen addr (default ":0")
  -d    enable debug logs
  -e string
        git commit email (default "gitnfs@example.com")
  -f string
        private key file
  -h    this help
  -k string
        private key string
  -m string
        git commit name (default "gitnfs")
  -o    make nfs server readonly
  -p string
        basic auth password or GitHub personal access token
  -r string
        git reference name (default "refs/heads/main")
  -s duration
        interval when sync nfs files to git repo (default 5s)
  -u string
        basic auth user name
```

## Dependent

* nfs server: <https://github.com/willscott/go-nfs> (fixed some bugs, and change filesystem to afero)
* git client: <https://github.com/go-git/go-git>
* filesystem: <https://github.com/spf13/afero>

## TODO

- [x] feature: git repo auth: username, ssh
- [x] feature: cache files to local dir
- [ ] feature: ability to set number of commit history
- [x] feature: nfs readonly mode
- [ ] feature: docker and kubernetes depoly
- [ ] optimize: file diff before push
