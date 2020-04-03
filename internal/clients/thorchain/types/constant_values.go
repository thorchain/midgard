package types

type ConstantValues struct {
	Int64Values  map[string]int64  `json:"int_64_values"`
	BoolValues   map[string]bool   `json:"bool_values"`
	StringValues map[string]string `json:"string_values"`
}
