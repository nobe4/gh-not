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

const userTypeUser = "User"

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
			User:    notifications.User{Login: "author-" + id, Type: userTypeUser},
			HTMLURL: "https://html.url/" + id,
		})
	case strings.HasPrefix(path, "https://latest.comment.url/"):
		id := strings.TrimPrefix(path, "https://latest.comment.url/")

		return jsonResponse(struct {
			User notifications.User `json:"user"`
		}{
			User: notifications.User{Login: "commentor-" + id, Type: userTypeUser},
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

func testEnrichManager(requestor *enrichRequestor, workers int) *Manager {
	return &Manager{
		client:        gh.NewClient(requestor, nil, gh.Endpoint{}),
		config:        &config.Data{Enrichment: config.Enrichment{Workers: workers}},
		ForceStrategy: 0,
	}
}

func TestEnrichDefaultsToSequential(t *testing.T) {
	t.Parallel()

	requestor := &enrichRequestor{}
	m := testEnrichManager(requestor, 0)

	ns := make(notifications.Notifications, 0, 3)
	for i := range 3 {
		ns = append(ns, testNotification(i))
	}

	m.Enrich(ns)

	if ns[2].Author.Login != "author-2" {
		t.Fatalf("expected author to be enriched, got %q", ns[2].Author.Login)
	}

	if ns[2].LatestCommentor.Login != "commentor-2" {
		t.Fatalf("expected latest commentor to be enriched, got %q", ns[2].LatestCommentor.Login)
	}

	if !ns[2].Meta.Enriched {
		t.Fatal("expected notification to be marked enriched")
	}

	if requestor.maxActive > 1 {
		t.Fatalf("expected default enrichment to stay sequential, max concurrency = %d", requestor.maxActive)
	}
}

func TestEnrichParallelWhenConfigured(t *testing.T) {
	t.Parallel()

	requestor := &enrichRequestor{}
	m := testEnrichManager(requestor, 10)

	ns := make(notifications.Notifications, 0, 30)
	for i := range 30 {
		ns = append(ns, testNotification(i))
	}

	m.Enrich(ns)

	if ns[13].Author.Login != "author-13" {
		t.Fatalf("expected author to be enriched, got %q", ns[13].Author.Login)
	}

	if ns[13].LatestCommentor.Login != "commentor-13" {
		t.Fatalf("expected latest commentor to be enriched, got %q", ns[13].LatestCommentor.Login)
	}

	if !ns[13].Meta.Enriched {
		t.Fatal("expected notification to be marked enriched")
	}

	if requestor.maxActive <= 1 {
		t.Fatalf("expected configured workers to run in parallel, max concurrency = %d", requestor.maxActive)
	}
}

func TestEnrichSkipsDoneWithoutForce(t *testing.T) {
	t.Parallel()

	requestor := &enrichRequestor{}
	m := testEnrichManager(requestor, 0)

	done := testNotification(1)
	done.Meta.Done = true

	notDone := testNotification(2)

	m.Enrich(notifications.Notifications{done, notDone})

	if requestor.requestMade != 2 {
		t.Fatalf("expected only non-done notification to be enriched (2 requests), got %d", requestor.requestMade)
	}

	if done.Author.Login != "" {
		t.Fatalf("expected done notification to be skipped, got author %q", done.Author.Login)
	}

	if notDone.Author.Login != "author-2" {
		t.Fatalf("expected non-done notification to be enriched, got author %q", notDone.Author.Login)
	}

	if !notDone.Meta.Enriched {
		t.Fatal("expected non-done notification to be marked enriched")
	}
}

func TestEnrichSkipsEnrichedNotifications(t *testing.T) {
	t.Parallel()

	requestor := &enrichRequestor{}
	m := testEnrichManager(requestor, 0)

	cached := testNotification(1)
	cached.Meta.Enriched = true
	cached.Subject.HTMLURL = "https://cached.html/1"
	cached.LatestCommentor = notifications.User{Login: "cached-commentor", Type: userTypeUser}

	notCached := testNotification(2)

	m.Enrich(notifications.Notifications{cached, notCached})

	if cached.LatestCommentor.Login != "cached-commentor" {
		t.Fatalf("expected cached notification preserved, got %q", cached.LatestCommentor.Login)
	}

	if requestor.requestMade != 2 {
		t.Fatalf("expected only uncached notification to be enriched (2 requests), got %d", requestor.requestMade)
	}

	if notCached.LatestCommentor.Login != "commentor-2" {
		t.Fatalf("expected uncached notification to be enriched, got %q", notCached.LatestCommentor.Login)
	}

	if !notCached.Meta.Enriched {
		t.Fatal("expected uncached notification to be marked enriched")
	}
}

func TestEnrichForceBypassesCachedAndDone(t *testing.T) {
	t.Parallel()

	requestor := &enrichRequestor{}
	m := testEnrichManager(requestor, 0)
	m.ForceStrategy = ForceEnrich

	doneCached := testNotification(1)
	doneCached.Meta.Done = true
	doneCached.Meta.Enriched = true

	m.Enrich(notifications.Notifications{doneCached})

	if requestor.requestMade != 2 {
		t.Fatalf("expected force to enrich done cached notification (2 requests), got %d", requestor.requestMade)
	}

	if doneCached.Author.Login != "author-1" {
		t.Fatalf("expected forced notification to be enriched, got author %q", doneCached.Author.Login)
	}

	if !doneCached.Meta.Enriched {
		t.Fatal("expected forced notification to remain marked enriched")
	}
}

func TestEnrichContinuesAfterFailure(t *testing.T) {
	t.Parallel()

	requestor := &enrichRequestor{}
	m := testEnrichManager(requestor, 0)

	failed := testNotification(1)
	failed.Subject.URL = "https://unexpected.url/1"
	successful := testNotification(2)

	m.Enrich(notifications.Notifications{failed, successful})

	if failed.Meta.Enriched {
		t.Fatal("expected failed notification to remain unenriched")
	}

	if !successful.Meta.Enriched {
		t.Fatal("expected later notification to be enriched")
	}
}
