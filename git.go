package main

import (
	"bytes"
	git "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/hash"
	"strings"
	"log"
	"os/exec"
)

type Commit struct {
	object.Commit
	FileNoExt string
}

func Diff(file, hash string) ([]byte, error) {
	r, _ := git.PlainOpen(options.Dir)
	co, _ := r.CommitObject(hash.NewHash(hash))
	lp := co.ParentHashes()[0]
	p, _ := co.Patch(lp);
	for _, f := range(p.FilePatches()) {
			from, to := f.Files()
			if (from != nil && from.Path() == filename) || (to != nil && to.Path() == filename) {
				break;
			}
		}

	return out.Bytes(), err)
}

func Commits(filename string, n int) ([]Commit, error) {
	var commits []Commit
	fnx := filename[:strings.LastIndex(filename, ".")]
	r, _ := git.PlainOpen(options.Dir)
	l, _ := r.Log(&git.LogOptions{})
	head, _ := r.Head()
	last, _ := r.CommitObject(head.Hash())
	l.ForEach(func(co *object.Commit) error {
		p, _ := co.Patch(last);
		log.Printf("%+v", p.FilePatches());
		for _, f := range(p.FilePatches()) {
			//log.Println(f);
			from, to := f.Files()
			if (from != nil && from.Path() == filename) || (to != nil && to.Path() == filename) {
				log.Printf("%+v, %+v", last);
				commits = append(commits, Commit{*last,fnx} )
				break;
			}
		}
		log.Printf("done");
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
