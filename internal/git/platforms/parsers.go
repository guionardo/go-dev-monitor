package platforms

import (
	"fmt"
	"strings"
)

// COMMIT: https://github.com/guionardo/gs-dev/commit/aabbf1201831f68e4b20a8fcfcc9c4bc2166e99f

// git@github.com:gin-contrib/cors.git
func githubSSHparser(url string) (name string, owner string, repository string, commitFunc commitURLFunc, ok bool) {
	if !strings.HasPrefix(url, "git@github.com:") || !strings.HasSuffix(url, ".git") {
		return
	}
	s := strings.TrimSuffix(strings.TrimPrefix(url, "git@github.com:"), ".git")
	w := strings.SplitN(s, "/", 2)
	name = "github.com"
	owner = w[0]
	repository = w[1]
	commitFunc = func(hash string) string {
		return fmt.Sprintf("https://github.com/%s/%s/commit/%s", owner, repository, hash)
	}
	ok = true
	return
}

// https://github.com/gin-contrib/cors.git
func githubHTTPSparser(url string) (name string, owner string, repository string, commitFunc commitURLFunc, ok bool) {
	if !strings.HasPrefix(url, "https://github.com/") || !strings.HasSuffix(url, ".git") {
		return
	}
	s := strings.TrimSuffix(strings.TrimPrefix(url, "https://github.com/"), ".git")
	w := strings.SplitN(s, "/", 2)
	name = "github.com"
	owner = w[0]
	repository = w[1]
	commitFunc = func(hash string) string {
		return fmt.Sprintf("https://github.com/%s/%s/commit/%s", owner, repository, hash)
	}
	ok = true
	return
}

// git@gitlab.com:gitlab-org/gitlab.git

// https://gitlab.com/gitlab-org/gitlab/-/commit/d29196d3a2f5c3b7fc2f091f40139daf424bf93c

func gitlabSSHparser(url string) (name string, owner string, repository string, commitFunc commitURLFunc, ok bool) {
	if !strings.HasPrefix(url, "git@gitlab.com:") || !strings.HasSuffix(url, ".git") {
		return
	}
	s := strings.TrimSuffix(strings.TrimPrefix(url, "git@gitlab.com:"), ".git")
	w := strings.SplitN(s, "/", 2)
	name = "gitlab.com"
	owner = w[0]
	repository = w[1]
	commitFunc = func(hash string) string {
		return fmt.Sprintf("https://gitlab.com/%s/%s/-/commit/%s", owner, repository, hash)
	}
	ok = true
	return
}

// https://gitlab.com/gitlab-org/gitlab
func gitlabHTTPSparser(url string) (name string, owner string, repository string, commitFunc commitURLFunc, ok bool) {
	if !strings.HasPrefix(url, "https://gitlab.com/") || !strings.HasSuffix(url, ".git") {
		return
	}
	s := strings.TrimSuffix(strings.TrimPrefix(url, "https://gitlab.com/"), ".git")
	w := strings.SplitN(s, "/", 2)
	name = "gitlab.com"
	owner = w[0]
	repository = w[1]
	commitFunc = func(hash string) string {
		return fmt.Sprintf("https://gitlab.com/%s/%s/-/commit/%s", owner, repository, hash)
	}
	ok = true
	return
}

// keybase://private/guionardo/guionardo
func keybaseParser(url string) (name string, owner string, repository string, commitFunc commitURLFunc, ok bool) {
	if !strings.HasPrefix(url, "keybase://") {
		return
	}
	s := strings.TrimSuffix(strings.TrimPrefix(url, "keybase://"), ".git")
	w := strings.SplitN(s, "/", 3)
	name = fmt.Sprintf("keybase:%s", w[0])
	owner = w[1]
	repository = w[2]
	commitFunc = func(hash string) string {
		return ""
	}
	ok = true
	return
}
