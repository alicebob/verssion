// in-memory implementation of the DB interface for tests
package w

type Memory struct {
	hist []Page
}

func NewMemory() *Memory {
	return &Memory{}
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
	var seen = map[string]Page{}
	for _, p := range m.hist {
		if sp, ok := seen[p.Page]; ok && sp.T.After(p.T) {
			continue
		}
		seen[p.Page] = p
	}

	var ps []Page
	for _, p := range seen {
		ps = append(ps, p)
	}
	return ps, nil
}

func (m *Memory) Current(...string) ([]Page, error) {
	return nil, nil
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
	return nil
}

func (m *Memory) Known() ([]string, error) {
	return nil, nil
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
