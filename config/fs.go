package config

import (
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
	"io"
	"os"
)

type Filesystem interface {
	Get(path string) ([]byte, error)
	Put(path string, data []byte) error
	Exists(path string) bool
}

type Local struct {
	Path string
}

type Remote struct {
	Url  string
	Path string
}

func (l *Local) Get(path string) ([]byte, error) {
	return os.ReadFile(l.Path + "/" + path)
}

func (l *Local) Put(path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(l.Path+"/"+path, data, perm)
}

func (l *Local) Exists(path string) bool {
	_, err := os.Stat(l.Path + "/" + path)

	return err == nil
}

func (r *Remote) Get(path string) ([]byte, error) {
	fs := memfs.New()

	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: r.Url,
	})
	if err != nil {
		return nil, err
	}

	f, err := fs.Open(r.Path + "/" + path)
	if err != nil {
		return nil, err
	}

	all, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	fmt.Println(string(all))

	return nil, err
}
