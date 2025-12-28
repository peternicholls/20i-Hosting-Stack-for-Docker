// Project: 20i Stack Manager TUI
// File: client_test.go
// Purpose: Unit tests for Docker client wrapper (connection, mapping, list)
// Version: 0.1.0
// Updated: 2025-12-28

package docker

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
)

type fakeDockerClient struct {
	containers     []container.Summary
	listErr        error
	listOptions    []container.ListOptions
	startErr       error
	startOptions   []container.StartOptions
	stopErr        error
	stopOptions    []container.StopOptions
	restartErr     error
	restartOptions []container.StopOptions
	connectionErr  error
}

func (f *fakeDockerClient) Ping(ctx context.Context) (types.Ping, error) {
	return types.Ping{}, f.connectionErr
}

func (f *fakeDockerClient) ContainerList(ctx context.Context, options container.ListOptions) ([]container.Summary, error) {
	f.listOptions = append(f.listOptions, options)
	if f.listErr != nil {
		return nil, f.listErr
	}
	return f.containers, nil
}

func (f *fakeDockerClient) ContainerStart(ctx context.Context, containerID string, options container.StartOptions) error {
	f.startOptions = append(f.startOptions, options)
	return f.startErr
}

func (f *fakeDockerClient) ContainerStop(ctx context.Context, containerID string, options container.StopOptions) error {
	f.stopOptions = append(f.stopOptions, options)
	return f.stopErr
}

func (f *fakeDockerClient) ContainerRestart(ctx context.Context, containerID string, options container.StopOptions) error {
	f.restartOptions = append(f.restartOptions, options)
	return f.restartErr
}

func (f *fakeDockerClient) ContainerStats(ctx context.Context, containerID string, stream bool) (container.StatsResponseReader, error) {
	// Return nil for now - stats not tested in existing tests
	return container.StatsResponseReader{}, nil
}

func TestMapConnectionError(t *testing.T) {
	cases := []struct {
		name string
		err  error
		want error
	}{
		{"permission denied", errors.New("permission denied"), ErrPermissionDenied},
		{"connection refused", errors.New("connect: connection refused"), ErrDaemonUnreachable},
		{"cannot connect to daemon", errors.New("Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"), ErrDaemonUnreachable},
		{"windows pipe missing", errors.New("open //./pipe/docker_engine: The system cannot find the file specified."), ErrDaemonUnreachable},
		{"no such file", errors.New("no such file or directory"), ErrDaemonUnreachable},
		{"timeout", errors.New("i/o timeout while dialing"), ErrTimeout},
		{"deadline exceeded", context.DeadlineExceeded, ErrTimeout},
		{"not found", errors.New("not found"), ErrNotFound},
		{"port conflict", errors.New("address already in use"), ErrConflict},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := mapConnectionError(c.err)
			if !errors.Is(got, c.want) {
				t.Fatalf("mapConnectionError(%v) = %v, want %v", c.err, got, c.want)
			}
		})
	}
}

func TestNewClient_Integration(t *testing.T) {
	// Integration test: only run when explicitly requested
	if os.Getenv("DOCKER_INTEGRATION_TEST") != "1" {
		t.Skip("skipping integration test; set DOCKER_INTEGRATION_TEST=1 to run")
	}

	ctx := context.Background()
	c, err := NewClient(ctx)
	if err != nil {
		t.Fatalf("NewClient returned error: %v", err)
	}
	if c == nil {
		t.Fatal("expected non-nil client")
	}
}

func TestNewClient_CanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := NewClient(ctx)
	if err == nil {
		t.Fatal("expected error when creating client with canceled context")
	}
}

func TestMapDockerState(t *testing.T) {
	cases := []struct {
		name  string
		state string
		want  ContainerStatus
	}{
		{"running", "running", StatusRunning},
		{"restarting", "restarting", StatusRestarting},
		{"stopped-exited", "exited", StatusStopped},
		{"stopped-created", "created", StatusStopped},
		{"stopped-paused", "paused", StatusStopped},
		{"stopped-dead", "dead", StatusStopped},
		{"unknown", "weird", StatusError},
		{"empty", "", StatusError},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := mapDockerState(c.state)
			if got != c.want {
				t.Fatalf("mapDockerState(%q) = %v, want %v", c.state, got, c.want)
			}
		})
	}
}

func TestListContainers(t *testing.T) {
	mock := &fakeDockerClient{
		containers: []container.Summary{
			{
				ID:     "abc123",
				Names:  []string{"/myproject-apache-1"},
				Image:  "php:8.2-apache",
				State:  "running",
				Status: "Up 2 minutes",
				Labels: map[string]string{
					"com.docker.compose.project": "myproject",
					"com.docker.compose.service": "apache",
				},
			},
		},
	}

	client := &Client{cli: mock, ctx: context.Background()}
	containers, err := client.ListContainers("myproject")
	if err != nil {
		t.Fatalf("ListContainers returned error: %v", err)
	}

	if len(containers) != 1 {
		t.Fatalf("expected 1 container, got %d", len(containers))
	}

	if containers[0].Service != "apache" {
		t.Fatalf("expected service 'apache', got %q", containers[0].Service)
	}

	if containers[0].Name != "myproject-apache-1" {
		t.Fatalf("expected name 'myproject-apache-1', got %q", containers[0].Name)
	}

	if len(mock.listOptions) != 1 {
		t.Fatalf("expected 1 list call, got %d", len(mock.listOptions))
	}

	options := mock.listOptions[0]
	if !options.All {
		t.Fatalf("expected ListContainers to set All=true")
	}

	if !options.Filters.Match("label", "com.docker.compose.project=myproject") {
		t.Fatalf("expected project label filter to be set")
	}
}

func TestListContainers_Empty(t *testing.T) {
	mock := &fakeDockerClient{}
	client := &Client{cli: mock, ctx: context.Background()}

	containers, err := client.ListContainers("myproject")
	if err != nil {
		t.Fatalf("ListContainers returned error: %v", err)
	}

	if len(containers) != 0 {
		t.Fatalf("expected 0 containers, got %d", len(containers))
	}
}

func TestListContainers_Error(t *testing.T) {
	mock := &fakeDockerClient{listErr: errors.New("boom")}
	client := &Client{cli: mock, ctx: context.Background()}

	_, err := client.ListContainers("myproject")
	if err == nil {
		t.Fatal("expected error from ListContainers")
	}
}

func TestStartContainer(t *testing.T) {
	mock := &fakeDockerClient{}
	client := &Client{cli: mock, ctx: context.Background()}

	if err := client.StartContainer("abc123"); err != nil {
		t.Fatalf("StartContainer returned error: %v", err)
	}

	if len(mock.startOptions) != 1 {
		t.Fatalf("expected 1 start call, got %d", len(mock.startOptions))
	}
}

func TestStartContainer_Error(t *testing.T) {
	mock := &fakeDockerClient{startErr: errors.New("start failed")}
	client := &Client{cli: mock, ctx: context.Background()}

	if err := client.StartContainer("abc123"); err == nil {
		t.Fatal("expected error from StartContainer")
	}
}

func TestStopContainer_DefaultTimeout(t *testing.T) {
	mock := &fakeDockerClient{}
	client := &Client{cli: mock, ctx: context.Background()}

	if err := client.StopContainer("abc123", 0); err != nil {
		t.Fatalf("StopContainer returned error: %v", err)
	}

	if len(mock.stopOptions) != 1 {
		t.Fatalf("expected 1 stop call, got %d", len(mock.stopOptions))
	}

	if mock.stopOptions[0].Timeout == nil || *mock.stopOptions[0].Timeout != 10 {
		t.Fatalf("expected default timeout 10, got %v", mock.stopOptions[0].Timeout)
	}
}

func TestStopContainer_CustomTimeout(t *testing.T) {
	mock := &fakeDockerClient{}
	client := &Client{cli: mock, ctx: context.Background()}

	if err := client.StopContainer("abc123", 5); err != nil {
		t.Fatalf("StopContainer returned error: %v", err)
	}

	if len(mock.stopOptions) != 1 {
		t.Fatalf("expected 1 stop call, got %d", len(mock.stopOptions))
	}

	if mock.stopOptions[0].Timeout == nil || *mock.stopOptions[0].Timeout != 5 {
		t.Fatalf("expected timeout 5, got %v", mock.stopOptions[0].Timeout)
	}
}

func TestStopContainer_Error(t *testing.T) {
	mock := &fakeDockerClient{stopErr: errors.New("stop failed")}
	client := &Client{cli: mock, ctx: context.Background()}

	if err := client.StopContainer("abc123", 5); err == nil {
		t.Fatal("expected error from StopContainer")
	}
}

func TestRestartContainer_DefaultTimeout(t *testing.T) {
	mock := &fakeDockerClient{}
	client := &Client{cli: mock, ctx: context.Background()}

	if err := client.RestartContainer("abc123", 0); err != nil {
		t.Fatalf("RestartContainer returned error: %v", err)
	}

	if len(mock.restartOptions) != 1 {
		t.Fatalf("expected 1 restart call, got %d", len(mock.restartOptions))
	}

	if mock.restartOptions[0].Timeout == nil || *mock.restartOptions[0].Timeout != 10 {
		t.Fatalf("expected default timeout 10, got %v", mock.restartOptions[0].Timeout)
	}
}

func TestRestartContainer_CustomTimeout(t *testing.T) {
	mock := &fakeDockerClient{}
	client := &Client{cli: mock, ctx: context.Background()}

	if err := client.RestartContainer("abc123", 7); err != nil {
		t.Fatalf("RestartContainer returned error: %v", err)
	}

	if len(mock.restartOptions) != 1 {
		t.Fatalf("expected 1 restart call, got %d", len(mock.restartOptions))
	}

	if mock.restartOptions[0].Timeout == nil || *mock.restartOptions[0].Timeout != 7 {
		t.Fatalf("expected timeout 7, got %v", mock.restartOptions[0].Timeout)
	}
}

func TestRestartContainer_Error(t *testing.T) {
	mock := &fakeDockerClient{restartErr: errors.New("restart failed")}
	client := &Client{cli: mock, ctx: context.Background()}

	if err := client.RestartContainer("abc123", 5); err == nil {
		t.Fatal("expected error from RestartContainer")
	}
}
