// in-memory implementation of the DB interface for tests

package core

import (
	"sort"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Memory struct {
	hist    []Page
	current map[string]Page
	mu      sync.Mutex
	curated map[string]Curated
}

func NewMemory() *Memory {
	return &Memory{
		current: map[string]Page{},
		curated: map[string]Curated{},
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
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	m.mu.Lock()
	defer m.mu.Unlock()

	t := time.Now()
	ids := id.String()
	m.curated[ids] = Curated{
		Created:     t,
		LastUpdated: t,
		LastUsed:    t,
	}
	return ids, nil
}

func (m *Memory) LoadCurated(id string) (*Curated, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c, ok := m.curated[id]
	if !ok {
		return nil, nil
	}
	return &c, nil
}

func (m *Memory) CuratedSetPages(id string, pages []string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.curated[id]
	if !ok {
		return ErrCuratedNotFound
	}
	c.Pages = pages
	sort.Strings(c.Pages)
	m.curated[id] = c
	return nil
}

func (m *Memory) CuratedSetUsed(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.curated[id]
	if !ok {
		return ErrCuratedNotFound
	}
	c.LastUsed = time.Now().UTC()
	m.curated[id] = c
	return nil
}

func (m *Memory) CuratedSetTitle(id string, title string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	c, ok := m.curated[id]
	if !ok {
		return ErrCuratedNotFound
	}
	c.CustomTitle = title
	m.curated[id] = c
	return nil
}
