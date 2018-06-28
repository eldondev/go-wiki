package main

import (
	"bytes"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/format/diff"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"log"
	"strings"
)

type Commit struct {
	object.Commit
	FileNoExt string
}

func Contents(filename string) (string, error) {
	if r, err := git.PlainOpen(options.Dir); err != nil {
		return "", err
	} else {
		if head, err := r.Head(); err != nil {
			return "", err
		} else {
			if headCommit, err := r.CommitObject(head.Hash()); err != nil {
				return "", err
			} else {
				if tree, err := headCommit.Tree(); err != nil {
					return "", err
				} else {
					if entry, err := tree.FindEntry(filename); err != nil {
						return "", err
					} else {
						if file, err := tree.TreeEntryFile(entry); err != nil {
							return "", err
						} else {
							return file.Contents()
						}
					}
				}
			}
		}
	}
}

func Diff(file, hash string) ([]byte, error) {
	var err error
	if r, err := git.PlainOpen(options.Dir); err == nil {
		if co, err := r.CommitObject(plumbing.NewHash(hash)); err == nil {
			lp := co.ParentHashes[0]
			lpc, err := r.CommitObject(lp)
			p, err := lpc.Patch(co)
			var out bytes.Buffer
			for _, f := range p.FilePatches() {
				from, to := f.Files()
				if (from != nil && from.Path() == file) || (to != nil && to.Path() == file) {
					dw := diff.NewUnifiedEncoder(&out, 2)
					dw.Encode(p)
					return out.Bytes(), err
				}
			}
		}
	}
	return nil, err

}

func Commits(filename string, n int) ([]Commit, error) {
	var commits []Commit
	fnx := filename[:strings.LastIndex(filename, ".")]
	r, _ := git.PlainOpen(options.Dir)
	l, _ := r.Log(&git.LogOptions{})
	head, _ := r.Head()
	last, _ := r.CommitObject(head.Hash())
	l.ForEach(func(co *object.Commit) error {
		if co.NumParents() < 1 {
			return nil
		}
		cop, _ := co.Parent(0)
		p, _ := co.Patch(cop)
		for _, f := range p.FilePatches() {
			//log.Println(f);
			from, to := f.Files()
			if (from != nil && from.Path() == filename) || (to != nil && to.Path() == filename) {
				commits = append(commits, Commit{*co, fnx})
				break
			}
		}
		last = co
		return nil
	})

	return commits, nil
}

// Check if a path contains a Git repository
func IsGitRepository(path string) bool {
	r, err := git.PlainOpen(options.Dir)
	if err != nil {
		log.Println("ERROR", err)
		return false
	}

	r.Log(&git.LogOptions{})
	if err != nil {
		log.Println("ERROR", err)
		return false
	}

	return true
}
