package types

type ConstantValues struct {
	Int64Values  map[string]int64  `json:"int_64_values"`
	BoolValues   map[string]bool   `json:"bool_values"`
	StringValues map[string]string `json:"string_values"`
}

func (c ConstantValues) IsEmpty() bool {
	return len(c.Int64Values) == 0 && len(c.BoolValues) == 0 && len(c.StringValues) == 0
}
