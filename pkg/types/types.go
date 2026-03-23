package types

type Property struct {
	Name       string `json:"name"`
	ColumnName string `json:"-"`
	Alias      string `json:"alias,omitempty"`
}
type Layer struct {
	Name       string     `json:"name"`
	GroupName  string     `json:"groupName,omitempty"`
	Properties []Property `json:"properties"`
}
type Scheme struct {
	AutomaticCasing bool    `json:"automaticCasing"`
	Version         int     `json:"version"`
	Layers          []Layer `json:"layers"`
}
