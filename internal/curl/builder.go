package curl

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/bugsy6/recurlit/internal/models"
)

func Build(req models.Request) string {
	var parts []string

	parts = append(parts, "curl")

	if req.Method != models.GET && req.Method != "" {
		parts = append(parts, fmt.Sprintf("-X %s", req.Method))
	}

	// Build URL with query params
	rawURL := req.URL
	if rawURL == "" {
		rawURL = "https://example.com"
	}
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
	parts = append(parts, fmt.Sprintf(`"%s"`, rawURL))

	// Auth
	switch req.Auth.Type {
	case models.AuthBearer:
		if req.Auth.Token != "" {
			parts = append(parts, fmt.Sprintf(`-H "Authorization: Bearer %s"`, req.Auth.Token))
		}
	case models.AuthBasic:
		if req.Auth.Username != "" {
			parts = append(parts, fmt.Sprintf(`-u "%s:%s"`, req.Auth.Username, req.Auth.Password))
		}
	}

	// Headers
	for _, h := range req.Headers {
		if h.Enabled && h.Key != "" {
			parts = append(parts, fmt.Sprintf(`-H "%s: %s"`, h.Key, h.Value))
		}
	}

	// Body
	if req.Body != "" {
		body := strings.ReplaceAll(req.Body, `'`, `'\''`)
		body = strings.ReplaceAll(body, "\r\n", "")
		body = strings.ReplaceAll(body, "\n", "")
		body = strings.ReplaceAll(body, "\r", "")
		parts = append(parts, fmt.Sprintf(`-d '%s'`, body))
	}

	return strings.Join(parts, " ")
}
