package nodes

type Node struct {
	Name   string `yaml:"name"`
	Url    string `yaml:"url"`
	Secret string `yaml:"secret"`
}

type Scope int

const (
	ScopeGlobal Scope = iota
	ScopeProject
)

type Entry struct {
	Node
	Scope Scope
}

type fileContent struct {
	Nodes []Node `yaml:"nodes"`
}
