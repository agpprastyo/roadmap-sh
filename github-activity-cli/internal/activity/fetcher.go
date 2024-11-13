package activity

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type DefaultFetcher struct{}

func (f *DefaultFetcher) FetchGitHubActivity(url, username string) ([]GitHubActivity, error) {
	resp, err := http.Get(fmt.Sprintf("%s/users/%s/events", url, username))
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println("failed to close response body")
		}
	}(resp.Body)

	var activities []GitHubActivity
	if err := json.NewDecoder(resp.Body).Decode(&activities); err != nil {
		return nil, err
	}

	return activities, nil
}
