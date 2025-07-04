package platforms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantName   string
		wantOwner  string
		wantRepo   string
		wantCommit string
	}{
		{"github_ssh", "git@github.com:gin-contrib/cors.git", "github.com", "gin-contrib", "cors", "https://github.com/gin-contrib/cors/commit/1234"},
		{"github_https", "https://github.com/gin-contrib/cors.git", "github.com", "gin-contrib", "cors", "https://github.com/gin-contrib/cors/commit/1234"},
		{"gitlab_ssh", "git@gitlab.com:gitlab-org/gitlab.git", "gitlab.com", "gitlab-org", "gitlab", "https://gitlab.com/gitlab-org/gitlab/-/commit/1234"},
		{"gitlab_ssh", "https://gitlab.com/gitlab-org/gitlab.git", "gitlab.com", "gitlab-org", "gitlab", "https://gitlab.com/gitlab-org/gitlab/-/commit/1234"},
		{"keybase_private", "keybase://private/guionardo/guionardo", "keybase:private", "guionardo", "guionardo", ""},
		{"keybase_public", "keybase://public/guionardo/guionardo", "keybase:public", "guionardo", "guionardo", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.url)
			if !assert.Equal(t, tt.wantName, got.Name) {
				return
			}
			if !assert.Equal(t, tt.wantOwner, got.Owner) {
				return
			}
			if !assert.Equal(t, tt.wantRepo, got.Repository) {
				return
			}
			if !assert.Equal(t, tt.wantCommit, got.CommitURL("1234")) {
				return
			}
		})
	}
}
