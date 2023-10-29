package generator

import (
	"docapi/generator/collector"
	"docapi/types"
)

type IntermediateGen struct {
	Api   types.Api
	Types map[string]types.Value
}

type Generator struct {
	api   *collector.ApiCollector
	types *collector.TypesCollector
}

func New() *Generator {
	return &Generator{}
}

func (g *Generator) Run(path string) (IntermediateGen, error) {
	apiCollector := collector.NewAPICollector()
	err := apiCollector.Run(path)
	if err != nil {
		return IntermediateGen{}, err
	}
	g.api = apiCollector

	typesCollector := collector.NewTypesCollector()
	err = typesCollector.Run(path)
	if err != nil {
		return IntermediateGen{}, err
	}
	g.types = typesCollector

	return g.Output()
}

func (g *Generator) Output() (IntermediateGen, error) {
	out := IntermediateGen{
		Types: map[string]types.Value{},
	}

	apis, err := g.api.Output()
	if err != nil {
		return IntermediateGen{}, err
	}
	out.Api = apis

	allTypes, err := g.types.Output()
	if err != nil {
		return IntermediateGen{}, err
	}

	keptTypes := map[string]types.Value{}
	for _, r := range g.api.Handlers {
		for _, resp := range r.Responses {
			if resp.Ref != "" {
				keptTypes[resp.Ref] = allTypes[resp.Ref]
			}
		}
	}
	out.Types = keptTypes
	return out, err
}
