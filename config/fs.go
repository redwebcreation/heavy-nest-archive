package config

import (
	"io"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Filesystem interface {
	Get(path string) ([]byte, error)
	Put(path string, data []byte, metadata interface{}) error
	Exists(path string) bool
}

type Local struct {
	Path string
}

type Remote struct {
	Url  string
	Path string
	fs   billy.Filesystem
	repo *git.Repository
}

func (r *Remote) Load() error {
	if r.fs != nil {
		return nil
	}

	fs := memfs.New()

	repo, err := git.Clone(memory.NewStorage(), fs, &git.CloneOptions{
		URL: r.Url,
	})

	if err != nil {
		return err
	}

	r.repo = repo
	r.fs = fs

	return nil
}

func (l *Local) Get(path string) ([]byte, error) {
	return os.ReadFile(l.Path + "/" + path)
}

func (l *Local) Put(path string, data []byte, metadata map[string]interface{}) error {
	var perm os.FileMode = 0644

	if metadata["perm"].(os.FileMode) != 0 {
		perm = metadata["perm"].(os.FileMode)
	}

	return os.WriteFile(l.Path+"/"+path, data, perm)
}

func (l *Local) Exists(path string) bool {
	_, err := os.Stat(l.Path + "/" + path)

	return err == nil
}

func (r *Remote) Get(path string) ([]byte, error) {
	err := r.Load()
	if err != nil {
		return nil, err
	}

	f, err := r.fs.Open(r.Path + "/" + path)
	if err != nil {
		return nil, err
	}

	return io.ReadAll(f)
}

func (r *Remote) Put(path string, data []byte, metadata map[string]interface{}) error {
	err := r.Load()
	if err != nil {
		return err
	}

	f, err := r.fs.Create(r.Path + "/" + path)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}

	repo, err := r.repo.Worktree()
	if err != nil {
		return err
	}

	_, err = repo.Commit("add "+path, &git.CommitOptions{})

	return err
}
