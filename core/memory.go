// in-memory implementation of the DB interface for tests

package core

import (
	"sync"
)

type Memory struct {
	hist    []Page
	current map[string]Page
	mu      sync.Mutex
}

func NewMemory() *Memory {
	return &Memory{
		current: map[string]Page{},
	}
}

var _ DB = NewMemory()

func (m *Memory) Last(page string) (*Page, error) {
	var last Page
	for _, p := range m.hist {
		if p.Page == page && p.T.After(last.T) {
			last = p
		}
	}
	if last.Page != "" {
		return &last, nil
	}
	return nil, ErrNotFound{Page: page}
}

func (m *Memory) Recent(n int) ([]Page, error) {
	// TODO: this is not right
	h, err := m.CurrentAll()
	if err != nil {
		return nil, err
	}
	if len(h) > n {
		h = h[:n]
	}
	return h, nil
}

func (m *Memory) CurrentAll() ([]Page, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var ps []Page
	for _, p := range m.current {
		ps = append(ps, p)
	}
	return ps, nil
}

func (m *Memory) Current(pages ...string) ([]Page, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var ps []Page
	for _, p := range pages {
		if c, ok := m.current[p]; ok {
			ps = append(ps, c)
		}
	}
	return ps, nil
}

func (m *Memory) History(pages ...string) ([]Page, error) {
	var ps []Page
	for _, p := range m.hist {
		for _, pn := range pages {
			if p.Page == pn {
				ps = append(ps, p)
				break
			}
		}
	}
	return ps, nil
}

func (m *Memory) Store(p Page) error {
	m.hist = append(m.hist, p)

	m.mu.Lock()
	defer m.mu.Unlock()
	old, ok := m.current[p.Page]
	if !ok || old.StableVersion != p.StableVersion {
		m.current[p.Page] = p
	}

	return nil
}

func (m *Memory) Known() ([]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var ps []string
	for k := range m.current {
		ps = append(ps, k)
	}
	return ps, nil
}

func (m *Memory) CreateCurated() (string, error) {
	return "fake", nil
}

func (m *Memory) LoadCurated(string) (*Curated, error) {
	return nil, nil
}

func (m *Memory) CuratedPages(string, []string) error {
	return nil
}

func (m *Memory) CuratedUsed(string) error {
	return nil
}

func (m *Memory) CuratedTitle(string, string) error {
	return nil
}
