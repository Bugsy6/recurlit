package executor

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/bugsy6/recurlit/internal/models"
)

type ResponseMsg struct {
	Response *models.Response
}

func Execute(req models.Request) tea.Cmd {
	return func() tea.Msg {
		resp, err := execute(req)
		if err != nil {
			return ResponseMsg{Response: &models.Response{Error: err.Error()}}
		}
		return ResponseMsg{Response: resp}
	}
}

func execute(req models.Request) (*models.Response, error) {
	start := time.Now()

	rawURL := req.URL
	if rawURL == "" {
		return nil, fmt.Errorf("URL is empty")
	}

	// Apply query params
	if len(req.Params) > 0 {
		vals := url.Values{}
		for _, p := range req.Params {
			if p.Enabled && p.Key != "" {
				vals.Set(p.Key, p.Value)
			}
		}
		if len(vals) > 0 {
			if strings.Contains(rawURL, "?") {
				rawURL += "&" + vals.Encode()
			} else {
				rawURL += "?" + vals.Encode()
			}
		}
	}

	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = bytes.NewBufferString(req.Body)
	}

	method := string(req.Method)
	if method == "" {
		method = "GET"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(ctx, method, rawURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	// Auth
	switch req.Auth.Type {
	case models.AuthBearer:
		if req.Auth.Token != "" {
			httpReq.Header.Set("Authorization", "Bearer "+req.Auth.Token)
		}
	case models.AuthBasic:
		httpReq.SetBasicAuth(req.Auth.Username, req.Auth.Password)
	}

	// Headers
	for _, h := range req.Headers {
		if h.Enabled && h.Key != "" {
			httpReq.Header.Set(h.Key, h.Value)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	return &models.Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    map[string][]string(resp.Header),
		Body:       string(bodyBytes),
		Duration:   time.Since(start),
	}, nil
}
