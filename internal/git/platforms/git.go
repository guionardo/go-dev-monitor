package platforms

import "fmt"

type (
	GitPlatform struct {
		Name          string
		Owner         string
		Repository    string
		commitURLFunc commitURLFunc
	}
	parserFunc    func(url string) (name string, owner string, repository string, commitFunc commitURLFunc, ok bool)
	commitURLFunc func(commitHash string) string
)

var parsers = []parserFunc{
	githubSSHparser,
	githubHTTPSparser,
	gitlabSSHparser,
	gitlabHTTPSparser,
	keybaseParser,
}

func New(url string) *GitPlatform {
	for _, parser := range parsers {
		if name, owner, repository, commitFunc, ok := parser(url); ok {
			return &GitPlatform{
				name, owner, repository, commitFunc,
			}
		}
	}
	return &GitPlatform{Repository: url}
}

func (p *GitPlatform) String() string {
	if len(p.Name) == 0 {
		return "unknown:" + p.Repository
	}
	return fmt.Sprintf("%s: %s/%s", p.Name, p.Owner, p.Repository)
}

func (p *GitPlatform) CommitURL(hash string) string {
	if p.commitURLFunc != nil {
		return p.commitURLFunc(hash)
	}
	return ""
}
