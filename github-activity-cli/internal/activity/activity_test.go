package activity

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchGitHubActivity(t *testing.T) {
	// Mock server to simulate GitHub API
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/users/testuser/events" {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`[{"type":"PushEvent","repo":{"name":"testrepo"},"created_at":"2023-10-01T00:00:00Z","payload":{"commits":[{"message":"Initial commit"}]}}]`))
			if err != nil {
				t.Fatalf("failed to write response: %v", err)
			}
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer mockServer.Close()

	// Test case: valid user
	activities, err := FetchGitHubActivity(mockServer.URL, "testuser")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(activities) != 1 {
		t.Fatalf("expected 1 activity, got %d", len(activities))
	}

	// Test case: invalid user
	_, err = FetchGitHubActivity(mockServer.URL, "invaliduser")
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}

func TestDisplayActivity(t *testing.T) {
	activities := []GitHubActivity{
		{
			Type: "PushEvent",
			Repo: Repo{Name: "testrepo"},
			Payload: struct {
				Action  string `json:"action"`
				Ref     string `json:"ref"`
				RefType string `json:"ref_type"`
				Commits []struct {
					Message string `json:"message"`
				} `json:"commits"`
			}{
				Commits: []struct {
					Message string `json:"message"`
				}{
					{Message: "Initial commit"},
				},
			},
		},
	}

	// Test case: valid activities
	err := DisplayActivity("testuser", activities)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Test case: no activities
	err = DisplayActivity("testuser", []GitHubActivity{})
	if err == nil {
		t.Fatalf("expected error, got none")
	}
}
