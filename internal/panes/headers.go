package panes

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bugsy6/recurlit/internal/models"
)

type HeadersPane struct {
	kv *KVPane
}

func NewHeadersPane() *HeadersPane {
	return &HeadersPane{kv: NewKVPane("Headers")}
}

func (h *HeadersPane) Headers() []models.Header {
	var headers []models.Header
	for _, p := range h.kv.Pairs() {
		headers = append(headers, models.Header{Key: p.Key, Value: p.Value, Enabled: p.Enabled})
	}
	return headers
}

func (h *HeadersPane) SetHeaders(headers []models.Header) {
	pairs := make([]KVPair, len(headers))
	for i, hdr := range headers {
		pairs[i] = KVPair{Key: hdr.Key, Value: hdr.Value, Enabled: hdr.Enabled}
	}
	h.kv.SetPairs(pairs)
}

func (h *HeadersPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	p, cmd := h.kv.Update(msg)
	h.kv = p.(*KVPane)
	return h, cmd
}

func (h *HeadersPane) View() string               { return h.kv.View() }
func (h *HeadersPane) SetSize(width, height int)  { h.kv.SetSize(width, height) }
func (h *HeadersPane) SetActive(active bool)       { h.kv.SetActive(active) }
func (h *HeadersPane) SetInsert(insert bool)       { h.kv.SetInsert(insert) }
func (h *HeadersPane) IsDirty() bool               { return h.kv.IsDirty() }
func (h *HeadersPane) ClearDirty()                 { h.kv.ClearDirty() }
