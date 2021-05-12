# git-nfs

Make git repo as a nfs server storage backend

## Features

* NFS v3
* In memory filesystem to cache files
* Sync files to git repository every 5s (option to set)

## Getting Started

* install and run

```shell
go get github.com/iineva/gitnfs
# change port number whatever you want
gitnfs -d -a ":1234" https://github.com/iineva/gitnfs
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
Usage: gitnfs [options] <YOUR GIT REPO>

Options:
  -a string
        nfs listen addr (default ":0")
  -d    enable debug logs
  -e string
        git commit email (default "gitnfs@example.com")
  -h    this help
  -m string
        git commit name (default "gitnfs")
  -o    make nfs server readonly
  -r string
        git reference name (default "refs/heads/master")
  -s duration
        interval when sync nfs files to git repo (default 5s)
```

## Dependent

* nfs server: <https://github.com/willscott/go-nfs> (fixed some bugs, and change filesystem to afero)
* git client: <https://github.com/go-git/go-git>
* filesystem: <https://github.com/spf13/afero>

## TODO

* feature: ability to set number of commit history
* feature: nfs readonly mode
* feature: docker and kubernetes depoly
* feature: git repo auth: username, ssh
