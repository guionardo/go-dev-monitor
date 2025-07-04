package api

import (
	"time"

	"github.com/guionardo/go-dev-monitor/internal/repository"
)

type (
	AgentRequest struct {
		Hostname     string              `json:"hostname" validate:"required"`
		Repositories []*repository.Local `json:"repositories" validate:"required,min=1,dive"`
	}
	ServerResponse struct {
		Message string `json:"message"`
		Error   string `json:"error,omitempty"`
	}
	SummaryResponse struct {
		Origins map[string][]LocalRepositoryResponse `json:"origins"`
	}
	LocalRepositoryResponse struct {
		Hostname          string              `json:"host_name"`
		FolderName        string              `json:"folder_name"`
		CurrentBranch     string              `json:"current_branch"`
		Clean             bool                `json:"clean"`
		LastCommit        *repository.Commit  `json:"last_commit"`
		UntrackedFiles    []*FileInfoResponse `json:"untracked_files"`
		ChangedFiles      []*FileInfoResponse `json:"changed_files"`
		LastModifiedFiles []*FileInfoResponse `json:"last_changed_files"`
		FetchTime         time.Time           `json:"fetch_time"`
		Description       string              `json:"description"`
		Language          string              `json:"language"`
		LanguageIcon      string              `json:"language_icon"`
	}
	FileInfoResponse struct {
		Name   string    `json:"name"`
		Time   time.Time `json:"time"`
		Status string    `json:"status"`
	}
)

func ToLocalRepositoryResponse(r *repository.Local) LocalRepositoryResponse {
	return LocalRepositoryResponse{
		Hostname:          r.Hostname,
		FolderName:        r.FolderName,
		CurrentBranch:     r.CurrentBranch,
		Clean:             r.Clean,
		LastCommit:        r.LastCommit,
		UntrackedFiles:    ToFileInfoResponses(r.UntrackedFiles),
		ChangedFiles:      ToFileInfoResponses(r.ChangedFiles),
		LastModifiedFiles: ToFileInfoResponses(r.LastModifiedFiles),
		FetchTime:         r.FetchTime,
		Description:       r.Description,
		Language:          r.Language,
		LanguageIcon:      r.LanguageIcon,
	}
}

func ToFileInfoResponses(resp []repository.FileInfo) []*FileInfoResponse {
	var responses = make([]*FileInfoResponse, len(resp))
	for i, r := range resp {
		responses[i] = &FileInfoResponse{
			Name:   r.Name,
			Time:   r.Time,
			Status: string(rune(r.Status)),
		}
	}
	return responses
}
