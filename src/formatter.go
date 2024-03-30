package src

type Highlighter interface {
	Highlight() string
}

type JSON struct{}

func (j JSON) Highlight() string {
	return "json"
}

type YAML struct{}

func (y YAML) Highlight() string {
	return "yaml"
}

type XML struct{}

func (x XML) Highlight() string {
	return "xml"
}
