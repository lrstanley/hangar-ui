package types

type FilterMsg struct {
	Filter string
}

func (m FilterMsg) Value() string {
	return m.Filter
}
