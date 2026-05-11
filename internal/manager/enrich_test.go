package manager

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/nobe4/gh-not/internal/config"
	"github.com/nobe4/gh-not/internal/gh"
	"github.com/nobe4/gh-not/internal/notifications"
)

var errUnexpectedPath = errors.New("unexpected path")

type enrichRequestor struct {
	mu sync.Mutex

	active      int
	maxActive   int
	requestMade int
}

func (e *enrichRequestor) Request(_ string, path string, _ io.Reader) (*http.Response, error) {
	e.mu.Lock()
	e.active++

	e.requestMade++
	if e.active > e.maxActive {
		e.maxActive = e.active
	}
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		e.active--
		e.mu.Unlock()
	}()

	time.Sleep(10 * time.Millisecond)

	switch {
	case strings.HasPrefix(path, "https://subject.url/"):
		id := strings.TrimPrefix(path, "https://subject.url/")

		return jsonResponse(struct {
			User    notifications.User `json:"user"`
			HTMLURL string             `json:"html_url"`
		}{
			User:    notifications.User{Login: "author-" + id, Type: "User"},
			HTMLURL: "https://html.url/" + id,
		})
	case strings.HasPrefix(path, "https://latest.comment.url/"):
		id := strings.TrimPrefix(path, "https://latest.comment.url/")

		return jsonResponse(struct {
			User notifications.User `json:"user"`
		}{
			User: notifications.User{Login: "commentor-" + id, Type: "User"},
		})
	default:
		return nil, fmt.Errorf("%w: %s", errUnexpectedPath, path)
	}
}

func jsonResponse(body any) (*http.Response, error) {
	b, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal body: %w", err)
	}

	return &http.Response{Body: io.NopCloser(bytes.NewReader(b))}, nil
}

func testNotification(id int) *notifications.Notification {
	idS := strconv.Itoa(id)

	return &notifications.Notification{
		ID: idS,
		Subject: notifications.Subject{
			URL:              "https://subject.url/" + idS,
			LatestCommentURL: "https://latest.comment.url/" + idS,
		},
	}
}

func testManager(requestor *enrichRequestor) *Manager {
	return &Manager{
		client:        gh.NewClient(requestor, nil, gh.Endpoint{}),
		config:        &config.Data{},
		ForceStrategy: 0,
	}
}

func TestEnrichParallel(t *testing.T) {
	t.Parallel()

	requestor := &enrichRequestor{}
	m := testManager(requestor)

	ns := notifications.Notifications{}
	for i := range 30 {
		ns = append(ns, testNotification(i))
	}

	got := m.Enrich(ns)

	if got[13].Author.Login != "author-13" {
		t.Fatalf("expected author to be enriched, got %q", got[13].Author.Login)
	}

	if got[13].LatestCommentor.Login != "commentor-13" {
		t.Fatalf("expected latest commentor to be enriched, got %q", got[13].LatestCommentor.Login)
	}

	if requestor.maxActive <= 1 {
		t.Fatalf("expected enrich requests to run in parallel, max concurrency = %d", requestor.maxActive)
	}
}

func TestEnrichSkipsDoneWithoutForce(t *testing.T) {
	t.Parallel()

	requestor := &enrichRequestor{}
	m := testManager(requestor)

	done := testNotification(1)
	done.Meta.Done = true

	notDone := testNotification(2)

	_ = m.Enrich(notifications.Notifications{done, notDone})

	if requestor.requestMade != 2 {
		t.Fatalf("expected only non-done notification to be enriched (2 requests), got %d", requestor.requestMade)
	}

	if done.Author.Login != "" {
		t.Fatalf("expected done notification to be skipped, got author %q", done.Author.Login)
	}

	if notDone.Author.Login != "author-2" {
		t.Fatalf("expected non-done notification to be enriched, got author %q", notDone.Author.Login)
	}
}

func TestEnrichSkipsOnlyFullyCachedNotifications(t *testing.T) {
	t.Parallel()

	requestor := &enrichRequestor{}
	m := testManager(requestor)

	cached := testNotification(1)
	cached.Subject.HTMLURL = "https://cached.html/1"
	cached.LatestCommentor = notifications.User{Login: "cached-commentor", Type: "User"}

	partial := testNotification(2)
	partial.Subject.HTMLURL = "https://cached.html/2"

	_ = m.Enrich(notifications.Notifications{cached, partial})

	if requestor.requestMade != 2 {
		t.Fatalf("expected only partial notification to be enriched (2 requests), got %d", requestor.requestMade)
	}

	if cached.LatestCommentor.Login != "cached-commentor" {
		t.Fatalf("expected cached notification preserved, got %q", cached.LatestCommentor.Login)
	}

	if partial.LatestCommentor.Login != "commentor-2" {
		t.Fatalf("expected partial notification to be enriched, got %q", partial.LatestCommentor.Login)
	}
}
