package cmd

import (
	"errors"
	"testing"

	"github-activity-cli/internal/activity"
	"github.com/stretchr/testify/assert"
)

type MockFetcher struct {
	mockFetch func(url, username string) ([]activity.GitHubActivity, error)
}

func (m *MockFetcher) FetchGitHubActivity(url, username string) ([]activity.GitHubActivity, error) {
	return m.mockFetch(url, username)
}

func TestRunDisplayActivityCmd(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		mockFetch func(url, username string) ([]activity.GitHubActivity, error)
		wantErr   bool
		errMsg    string
	}{
		{
			name:    "no arguments",
			args:    []string{},
			wantErr: true,
			errMsg:  "please provide a username",
		},
		{
			name: "valid username",
			args: []string{"testuser"},
			mockFetch: func(url, username string) ([]activity.GitHubActivity, error) {
				return []activity.GitHubActivity{
					{Type: "PushEvent", Repo: activity.Repo{Name: "testrepo"}},
				}, nil
			},
			wantErr: false,
		},
		{
			name: "invalid username",
			args: []string{"invaliduser"},
			mockFetch: func(url, username string) ([]activity.GitHubActivity, error) {
				return nil, errors.New("user not found")
			},
			wantErr: true,
			errMsg:  "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fetcher = &MockFetcher{mockFetch: tt.mockFetch}

			err := RunDisplayActivityCmd(tt.args)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, tt.errMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
