package generator

type Generator struct {
	apis []*apiTemp
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Run(path string) error {
	a := newAPI()
	err := a.walk(path)
	if err != nil {
		return err
	}
	g.apis = append(g.apis, a)
	return nil
}

func (g *Generator) Output() (string, error) {
	s := ""
	for _, a := range g.apis {
		out, err := a.Output()
		if err != nil {
			return "", err
		}
		s += out
	}
	return s, nil
}
