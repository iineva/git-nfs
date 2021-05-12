package syncfile

import (
	"io"
	"os"
	"path/filepath"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-git/v5/plumbing/filemode"
	"github.com/spf13/afero"
)

func Afero2Billy(a afero.Fs, b billy.Filesystem) error {
	// remove all
	list, err := b.ReadDir(".")
	if err != nil {
		return err
	}
	for _, f := range list {
		b.Remove(f.Name())
	}

	// sync
	return doSyncAfero2Billy(a, b, ".")
}

func doSyncAfero2Billy(from afero.Fs, to billy.Filesystem, dir string) error {
	f, err := from.Open(dir)
	if err != nil {
		return err
	}
	list, err := f.Readdir(-1)
	if err != nil {
		return err
	}
	for _, f := range list {
		fileName := filepath.Join(dir, f.Name())
		if f.IsDir() {
			mode, _ := filemode.Dir.ToOSFileMode()
			err := to.MkdirAll(fileName, mode)
			if err != nil {
				return err
			}
			doSyncAfero2Billy(from, to, fileName)
		} else {
			bf, err := from.Open(fileName)
			defer bf.Close()
			if err != nil {
				return err
			}
			mode, _ := filemode.Regular.ToOSFileMode()
			af, err := to.OpenFile(fileName, os.O_RDWR|os.O_CREATE, mode)
			defer af.Close()
			if err != nil {
				return err
			}
			if _, err := io.Copy(af, bf); err != nil {
				return nil
			}
		}
	}
	return nil
}

func Billy2Afero(b billy.Filesystem, a afero.Fs) error {
	return doSyncBilly2Afero(b, a, ".")
}

func doSyncBilly2Afero(from billy.Filesystem, to afero.Fs, dir string) error {
	list, err := from.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range list {
		fileName := filepath.Join(dir, f.Name())
		if f.IsDir() {
			err := to.MkdirAll(fileName, f.Mode())
			if err != nil {
				return err
			}
			doSyncBilly2Afero(from, to, fileName)
		} else {
			bf, err := from.Open(fileName)
			defer bf.Close()
			if err != nil {
				return err
			}
			af, err := to.OpenFile(fileName, os.O_RDWR|os.O_CREATE, f.Mode())
			defer af.Close()
			if err != nil {
				return err
			}
			if _, err := io.Copy(af, bf); err != nil {
				return nil
			}
		}
	}
	return nil
}
