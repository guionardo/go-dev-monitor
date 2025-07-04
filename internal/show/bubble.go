package show

import (
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/guionardo/go-dev-monitor/internal/api"
	"github.com/guionardo/go-dev-monitor/internal/config"
	"github.com/guionardo/go-dev-monitor/internal/debug"
	"github.com/guionardo/go-dev-monitor/internal/repository"
	"github.com/guionardo/go-dev-monitor/internal/utils/git"
)

/*
STEPS:
 1. Get git root of the repository
 2. Read local repository
 3. Check availability of server API
    4.
*/
type (
	model struct {
		config        *config.Agent
		localFolder   string
		serverAddress string
		serverOnline  bool

		repository *repository.Local
	}
)

func display(cfg *config.Agent, localFolder string) error {
	m := model{
		config:        cfg,
		serverAddress: cfg.ServerAddress,
		localFolder:   localFolder,
	}
	if err := m.getGitRoot(); err != nil {
		return err
	}
	if err := m.readRepository(); err != nil {
		return err
	}
	if err := m.checkServer(); err != nil {
		return err
	}
	if err := m.fetchRemotes(); err != nil {
		return err
	}
	debug.Log().Debug("display", slog.Any("result", m))

	return nil
}

func (m *model) checkServer() error {
	c := &http.Client{
		Timeout: 10 * time.Second,
	}
	res, err := c.Get(fmt.Sprintf("%s/hc", m.serverAddress))
	if err == nil {
		res.Body.Close() // nolint:errcheck
		m.serverOnline = res.StatusCode == http.StatusOK
	} else {
		if e, ok := err.(net.Error); ok {
			err = fmt.Errorf("failed to connect timeout:%v", e.Timeout())
		}
		fmt.Println(err.Error())
	}
	return nil
}

func (m *model) getGitRoot() error {
	root, err := git.GetGitRoot(m.localFolder)
	if err != nil {
		return err
	}
	if m.localFolder != root {
		fmt.Printf("Repository root: %s\n", root)
		m.localFolder = root
	}
	return nil
}

func (m *model) readRepository() error {
	fmt.Printf("Reading repository... ")
	r, err := repository.New(m.localFolder, m.config.Hostname)
	if err != nil {
		return err
	}
	m.repository = r
	local := api.ToLocalRepositoryResponse(r)
	fmt.Print("\r")
	displayRepo(r.Origin, &local, true)

	return nil

}

func (m *model) fetchRemotes() error {
	if !m.serverOnline {
		return nil
	}
	c := &http.Client{
		Timeout: 10 * time.Second,
	}

	res, err := c.Get(fmt.Sprintf("%s/data?origin=%s", m.serverAddress, m.repository.Origin))
	if err == nil {
		defer res.Body.Close() // nolint:errcheck
		content, err := io.ReadAll(res.Body)
		if err != nil {
			return err
		}
		var response api.SummaryResponse
		err = response.UnmarshalJSON(content)
		if err != nil {
			return err
		}
		for origin, locals := range response.Origins {
			for _, local := range locals {
				if local.Hostname != m.config.Hostname {
					displayRepo(origin, &local, false)
				}
			}
		}

		m.serverOnline = res.StatusCode == http.StatusOK
	} else {
		if e, ok := err.(net.Error); ok {
			err = fmt.Errorf("failed to connect timeout:%v", e.Timeout())
		}
		fmt.Println(err.Error())
	}
	return nil
}

/*
 */
func displayRepo(origin string, r *api.LocalRepositoryResponse, header bool) {

	box := NewBox()
	if header {
		box.AddPanel("Origin: %s", color.BlueString(origin)).
			AddContent("Description: %s", r.Description).
			AddContent("Language: %s", r.Language)
	}
	box.AddPanel("Host: %s", r.Hostname).
		AddContent("Fetch time: %s", r.FetchTime.Format(time.RFC3339)).
		AddContent("Folder: %s", r.FolderName)
	st := box.AddPanel("Status")
	if r.Clean {
		st.AddContent("Clean: %s", color.GreenString("YES"))
	} else {
		st.AddContent("Clean: %s", color.YellowString("NO"))
	}

	if r.LastCommit != nil {
		box.AddPanel("Commit").
			AddContent("Branch: %s", color.YellowString(r.CurrentBranch)).
			AddContent("Hash %s %s", color.YellowString(r.LastCommit.Hash), color.YellowString(r.LastCommit.URL)).
			AddContent("When: %s", color.YellowString(r.LastCommit.When.Format(time.RFC3339))).
			AddContent("Message: %s", color.CyanString(r.LastCommit.Message)).
			AddContent("User: %s %s", color.YellowString(r.LastCommit.Author), color.YellowString(r.LastCommit.AuthorEmail))
	}
	if len(r.ChangedFiles) > 0 {
		cf := box.AddPanel("[%s] Changed files", color.BlueString("%d", len(r.ChangedFiles)))
		for _, f := range r.ChangedFiles {
			cf.AddContent("%s %20s %s", f.Status, humanTime(f.Time), f.Name)
		}
	}
	if len(r.LastModifiedFiles) > 0 {
		mf := box.AddPanel("[%s] Last modified files", color.BlueString("%d", len(r.LastModifiedFiles)))
		for _, f := range r.LastModifiedFiles {
			mf.AddContent("%s %20s %s", f.Status, humanTime(f.Time), f.Name)
		}
	}
	if len(r.UntrackedFiles) > 0 {
		uf := box.AddPanel("[%s] Untracked files", color.BlueString("%d", len(r.UntrackedFiles)))
		for i, f := range r.UntrackedFiles {
			if i > 10 {
				uf.AddContent("...")
				break
			}
			uf.AddContent("%s %20s %s", f.Status, humanTime(f.Time), f.Name)
		}
	}
	box.Write(os.Stdout)
}

func humanTime(t time.Time) string {
	d := time.Since(t)
	if d > time.Hour*48 {
		return t.Format(time.DateTime)
	}
	if d > time.Hour*24 {
		return t.Format("yesterday 15:04")
	}
	if d > time.Hour {
		return fmt.Sprintf("%.1f hours  ", d.Hours())
	}
	return fmt.Sprintf("%.1f minutes", d.Minutes())

}
