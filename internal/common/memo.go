package common

type Memo string

func (m *Memo) String() string {
	return string(*m)
}
