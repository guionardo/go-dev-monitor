package repository

import (
	"errors"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	git "github.com/go-git/go-git/v6"
	"github.com/go-git/go-git/v6/config"
	"github.com/guionardo/go-dev-monitor/internal/git/platforms"
	"github.com/guionardo/go-dev-monitor/internal/projects"
)

type (
	Local struct {
		r *git.Repository // The git repository instance

		Origin            string     `json:"origin"`
		FolderName        string     `json:"folder_name"`    // The folder name of the repository
		CurrentBranch     string     `json:"current_branch"` // The current branch of the repository
		Clean             bool       `json:"clean"`          // True if all files are in unmodified status
		UntrackedFiles    []FileInfo `json:"untracked"`      // Untracked files
		ChangedFiles      []FileInfo `json:"changed"`        // Changed files
		LastModifiedFiles []FileInfo `json:"last_modified"`  // Last all modified files
		Status            string     `json:"current_status"`
		LastCommit        *Commit    `json:"last_commit"`
		Hostname          string     `json:"hostname"`
		FetchTime         time.Time  `json:"fetch_time"`
		Description       string     `json:"description"`
		Language          string     `json:"language"`
		LanguageIcon      string     `json:"language_icon"`
	}
	FileInfo struct {
		Name   string
		Time   time.Time
		Status git.StatusCode
	}
	Commit struct {
		When        time.Time `json:"when"`
		Message     string    `json:"message"`
		Author      string    `json:"author"`
		AuthorEmail string    `json:"author_email"`
		Hash        string    `json:"hash"`
		URL         string    `json:"url"`
	}
)

const lastModifiedCount = 10

func New(folderName string, hostName string) (repo *Local, err error) {
	r, err := git.PlainOpen(folderName)
	if err != nil {
		return nil, err
	}
	currentBranch, err := r.Head()
	if err != nil {
		return nil, err
	}

	c, err := r.Config()
	if err != nil {
		return nil, err
	}

	origin, err := getOrigin(c.Remotes)
	if err != nil {
		return nil, err
	}

	// Project
	project, _ := projects.New(folderName)

	platform := platforms.New(origin)
	// Get description
	var description string
	if descriptionContent, err := os.ReadFile(path.Join(folderName, ".git", "description")); err == nil {
		description = string(descriptionContent)
	}
	if len(description) == 0 || strings.HasPrefix(description, "Unnamed repository") {
		if len(project.ProjectName) > 0 {
			description = project.ProjectName
		} else {
			description = platform.String()
		}
	}

	// Last commit
	// Get the commit object pointed to by HEAD
	var lastCommit *Commit
	if commitObject, err := r.CommitObject(currentBranch.Hash()); err == nil {
		lastCommit = &Commit{
			When:        commitObject.Author.When,
			Message:     commitObject.Message,
			Author:      commitObject.Author.Name,
			AuthorEmail: commitObject.Author.Email,
			Hash:        commitObject.Hash.String(),
			URL:         platform.CommitURL(commitObject.Hash.String()),
		}
	}

	// fmt.Printf("Last Commit Hash: %s\n", lastCommit.Hash)
	// fmt.Printf("Last Commit Message: %s\n", lastCommit.Message)
	// fmt.Printf("Last Commit Author: %s <%s>\n", lastCommit.Author.Name, lastCommit.Author.Email)
	// fmt.Printf("Last Commit Date: %s\n", lastCommit.Author.When)

	wt, err := r.Worktree()
	if err != nil {
		return nil, err
	}
	st, err := wt.Status()
	if err != nil {
		return nil, err
	}

	var (
		untracked = make([]FileInfo, 0)
		changed   = make([]FileInfo, 0)
		modified  = make([]FileInfo, 0)
	)
	for fileName, status := range st {
		fs, err := os.Stat(path.Join(folderName, fileName))
		if os.IsNotExist(err) {
			fi := FileInfo{fileName, time.Time{}, git.Deleted}
			changed = append(changed, fi)
			continue
		}
		if err != nil {
			return nil, err
		}
		fi := FileInfo{fileName, fs.ModTime(), status.Staging}
		modified = append(modified, fi)
		switch status.Staging {
		case git.Unmodified:
			continue
		case git.Untracked:
			untracked = append(untracked, fi)
			continue
		}
		changed = append(changed, fi)
	}

	sort.Slice(modified, func(i, j int) bool {
		return modified[i].Time.After(modified[j].Time)
	})

	return &Local{
		FolderName:        folderName,
		r:                 r,
		Origin:            origin,
		CurrentBranch:     currentBranch.Name().Short(),
		Clean:             st.IsClean(),
		UntrackedFiles:    untracked,
		ChangedFiles:      changed,
		LastModifiedFiles: modified[0:min(lastModifiedCount, len(modified))],
		Status:            st.String(),
		LastCommit:        lastCommit,
		Hostname:          hostName,
		FetchTime:         time.Now(),
		Description:       description,
		Language:          project.Language,
		LanguageIcon:      project.Icon,
	}, nil
}

func getOrigin(remoteConfig map[string]*config.RemoteConfig) (origin string, err error) {
	for name, remote := range remoteConfig {
		if name == "origin" && len(remote.URLs) > 0 {
			origin = remote.URLs[0]
			return
		}
		if len(remote.URLs) > 0 {
			origin = remote.URLs[0]
			return
		}
	}
	return "", errors.New("there is no remote in this repository")

}
