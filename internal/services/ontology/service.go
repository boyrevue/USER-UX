package ontology

type Service struct {
	// Ontology service dependencies
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) GetFormDefinitions() (interface{}, error) {
	// TODO: Move TTL parsing logic from ttl_parser.go here
	return map[string]interface{}{
		"drivers":  map[string]interface{}{"fields": []interface{}{}},
		"vehicles": map[string]interface{}{"fields": []interface{}{}},
		"claims":   map[string]interface{}{"fields": []interface{}{}},
	}, nil
}
