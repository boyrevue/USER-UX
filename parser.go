package main

import (
    "encoding/json"
    "log"
    "os"
    "path/filepath"
)

func LoadOntology() *OntologyData {
    return &OntologyData{
        Categories: loadCategories(),
        Fields:     loadFields(),
        Subforms:   loadSubforms(),
        Schemes:    map[string]ConceptScheme{},
    }
}

func loadCategories() map[string]Category {
    data, err := os.ReadFile(filepath.Join("ontology", "categories.json"))
    if err != nil {
        log.Printf("Could not load categories.json: %v", err)
        return getDefaultCategories()
    }
    var categories map[string]Category
    if err := json.Unmarshal(data, &categories); err != nil {
        log.Printf("Could not parse categories.json: %v", err)
        return getDefaultCategories()
    }
    return categories
}

func loadFields() map[string][]Field {
    data, err := os.ReadFile(filepath.Join("ontology", "fields.json"))
    if err != nil {
        log.Printf("Could not load fields.json: %v", err)
        return map[string][]Field{}
    }
    var fields map[string][]Field
    if err := json.Unmarshal(data, &fields); err != nil {
        log.Printf("Could not parse fields.json: %v", err)
        return map[string][]Field{}
    }
    return fields
}

func loadSubforms() map[string]Subform {
    data, err := os.ReadFile(filepath.Join("ontology", "subforms.json"))
    if err != nil {
        log.Printf("Could not load subforms.json: %v", err)
        return map[string]Subform{}
    }
    var subforms map[string]Subform
    if err := json.Unmarshal(data, &subforms); err != nil {
        log.Printf("Could not parse subforms.json: %v", err)
        return map[string]Subform{}
    }
    return subforms
}

func getDefaultCategories() map[string]Category {
    return map[string]Category{
        "drivers": Category{ID: "drivers", Title: "Driver Details", Icon: "👤", Order: 1},
        "vehicle": Category{ID: "vehicle", Title: "Vehicle Details", Icon: "🚗", Order: 2},
        "policy":  Category{ID: "policy", Title: "Policy Details", Icon: "📋", Order: 3},
        "claims":  Category{ID: "claims", Title: "Claims History", Icon: "📊", Order: 4},
    }
}
