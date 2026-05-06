package panes

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/bugsy6/recurlit/internal/models"
)

type ParamsPane struct {
	kv *KVPane
}

func NewParamsPane() *ParamsPane {
	return &ParamsPane{kv: NewKVPane("Params")}
}

func (p *ParamsPane) Params() []models.Param {
	var params []models.Param
	for _, pair := range p.kv.Pairs() {
		params = append(params, models.Param{Key: pair.Key, Value: pair.Value, Enabled: pair.Enabled})
	}
	return params
}

func (p *ParamsPane) SetParams(params []models.Param) {
	pairs := make([]KVPair, len(params))
	for i, param := range params {
		pairs[i] = KVPair{Key: param.Key, Value: param.Value, Enabled: param.Enabled}
	}
	p.kv.SetPairs(pairs)
}

func (pa *ParamsPane) Update(msg tea.Msg) (Pane, tea.Cmd) {
	p, cmd := pa.kv.Update(msg)
	pa.kv = p.(*KVPane)
	return pa, cmd
}

func (pa *ParamsPane) View() string               { return pa.kv.View() }
func (pa *ParamsPane) SetSize(width, height int)  { pa.kv.SetSize(width, height) }
func (pa *ParamsPane) SetActive(active bool)       { pa.kv.SetActive(active) }
func (pa *ParamsPane) SetInsert(insert bool)       { pa.kv.SetInsert(insert) }
func (pa *ParamsPane) IsDirty() bool               { return pa.kv.IsDirty() }
func (pa *ParamsPane) ClearDirty()                 { pa.kv.ClearDirty() }
