# TTL Implementation Guide

## Overview
This guide provides technical details for implementing and maintaining the TTL-as-single-source-of-truth architecture in CLIENT-UX.

---

## 1. TTL ONTOLOGY STRUCTURE

### 1.1 Core Ontology Pattern
```turtle
@prefix autoins: <http://autoins.example.org/ontology#> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .
@prefix xsd: <http://www.w3.org/2001/XMLSchema#> .
@prefix owl: <http://www.w3.org/2002/07/owl#> .

# Property Definition Template
autoins:propertyName a owl:DatatypeProperty ;
  rdfs:domain autoins:DomainClass ;
  rdfs:range xsd:dataType ;
  rdfs:label "human readable label" ;
  autoins:isRequired "true"^^xsd:boolean ;
  autoins:formHelpText "Help text for users" ;
  autoins:enumerationValues ("OPTION1" "OPTION2" "OPTION3") ;
  autoins:validationPattern "^regex$" ;
  autoins:errorMessage "Error message for validation" .
```

### 1.2 Domain Classes
- `autoins:Driver` - Driver-related properties
- `autoins:Vehicle` - Vehicle-related properties  
- `autoins:InsurancePolicy` - Policy-related properties
- `autoins:Coverage` - Coverage-related properties
- `autoins:Claims` - Claims-related properties (if needed)

### 1.3 Custom Properties
- `autoins:isRequired` - Boolean flag for required fields
- `autoins:formHelpText` - User-facing help text
- `autoins:enumerationValues` - List of valid options
- `autoins:validationPattern` - Regex validation pattern
- `autoins:errorMessage` - Validation error message
- `autoins:convictionRegion` - Regional tagging (e.g., "UK")

---

## 2. TTL PARSER IMPLEMENTATION

### 2.1 Parser Architecture (`ttl_parser.go`)
```go
// Dynamic TTL parsing using regex patterns
func ParseTTLOntology() (map[string]OntologySection, error) {
    // 1. Read TTL file
    ttlData, err := ioutil.ReadFile("ontology/autoins.ttl")
    
    // 2. Extract properties using regex
    propertyPattern := regexp.MustCompile(`autoins:(\w+)\s+a\s+owl:DatatypeProperty\s*;`)
    
    // 3. Build form structure dynamically
    // 4. Return sections organized by domain
}
```

### 2.2 Regex Patterns
```go
// Core patterns for TTL parsing
propertyPattern := regexp.MustCompile(`autoins:(\w+)\s+a\s+owl:DatatypeProperty\s*;`)
labelPattern := regexp.MustCompile(`rdfs:label\s+"([^"]+)"\s*;`)
domainPattern := regexp.MustCompile(`rdfs:domain\s+autoins:(\w+)\s*;`)
rangePattern := regexp.MustCompile(`rdfs:range\s+xsd:(\w+)\s*;`)
requiredPattern := regexp.MustCompile(`autoins:isRequired\s+"(true|false)"\^\^xsd:boolean\s*;`)
helpTextPattern := regexp.MustCompile(`autoins:formHelpText\s+"([^"]+)"\s*;`)
enumPattern := regexp.MustCompile(`autoins:enumerationValues\s+\(([^)]+)\)\s*;`)
```

### 2.3 Field Type Resolution
```go
// Automatic field type detection
switch range_ {
case "boolean":
    fieldType = "radio"
case "date":
    fieldType = "date"
default:
    if strings.Contains(propName, "email") {
        fieldType = "email"
    } else if strings.Contains(propName, "phone") {
        fieldType = "tel"
    }
}

// Handle enumeration values
if len(enumValues) > 0 {
    if len(enumValues) <= 3 {
        fieldType = "radio"
    } else {
        fieldType = "select"
    }
}
```

---

## 3. API ENDPOINT STRUCTURE

### 3.1 Ontology API (`/api/ontology`)
```json
{
  "drivers": {
    "id": "drivers",
    "label": "Driver Details",
    "fields": [
      {
        "property": "hasConvictions",
        "label": "has convictions",
        "type": "radio",
        "required": true,
        "helpText": "Any driving convictions in the last 5 years",
        "options": [
          {"value": "YES", "label": "YES"},
          {"value": "NO", "label": "NO"}
        ]
      }
    ]
  },
  "vehicles": { /* ... */ },
  "claims": { /* ... */ }
}
```

### 3.2 Field Structure
```go
type OntologyField struct {
    Property string        `json:"property"`    // TTL property name
    Label    string        `json:"label"`       // rdfs:label
    Type     string        `json:"type"`        // Resolved field type
    Required bool          `json:"required"`    // autoins:isRequired
    HelpText string        `json:"helpText"`    // autoins:formHelpText
    Options  []FieldOption `json:"options"`     // autoins:enumerationValues
    Domain   string        `json:"domain"`      // rdfs:domain
}
```

---

## 4. FRONTEND INTEGRATION

### 4.1 Dynamic Form Loading
```typescript
// React hook for ontology loading
useEffect(() => {
  const fetchOntology = async () => {
    const response = await fetch('/api/ontology');
    const ontologyData = await response.json();
    setOntology(ontologyData);
  };
  fetchOntology();
}, []);
```

### 4.2 Dynamic Field Rendering
```typescript
// Render fields based on ontology structure
const renderField = (field: OntologyField) => {
  switch (field.type) {
    case 'radio':
      return <RadioGroup options={field.options} required={field.required} />;
    case 'select':
      return <Select options={field.options} required={field.required} />;
    case 'date':
      return <DatePicker required={field.required} />;
    default:
      return <TextInput type={field.type} required={field.required} />;
  }
};
```

---

## 5. ADDING NEW FIELDS

### 5.1 Step-by-Step Process
1. **Add to TTL ontology**:
```turtle
autoins:newField a owl:DatatypeProperty ;
  rdfs:domain autoins:Driver ;
  rdfs:range xsd:string ;
  rdfs:label "New Field Label" ;
  autoins:isRequired "false"^^xsd:boolean ;
  autoins:formHelpText "Help text for new field" .
```

2. **Restart application** (TTL is parsed at startup)
3. **Field automatically appears** in API and frontend
4. **No code changes required**

### 5.2 Field with Options
```turtle
autoins:newSelectField a owl:DatatypeProperty ;
  rdfs:domain autoins:Driver ;
  rdfs:range xsd:string ;
  rdfs:label "Select Field" ;
  autoins:isRequired "true"^^xsd:boolean ;
  autoins:enumerationValues ("OPTION1" "OPTION2" "OPTION3") ;
  autoins:formHelpText "Choose one option" .
```

---

## 6. UK CONVICTION CODES EXAMPLE

### 6.1 Complete Implementation
```turtle
autoins:convictionType a owl:DatatypeProperty ;
  rdfs:domain autoins:Driver ;
  rdfs:range xsd:string ;
  rdfs:label "conviction type" ;
  autoins:isRequired "false"^^xsd:boolean ;
  autoins:formHelpText "Select the type of UK driving conviction" ;
  autoins:enumerationValues (
    "SP10" "SP20" "SP30" "SP40" "SP50" "SP60" 
    "DR10" "DR20" "DR30" "DR40" "DR50" "DR60" "DR70" "DR80"
    "IN10" "LC20" "LC30" "LC40" "LC50"
    "MS10" "MS20" "MS30" "MS50" "MS60" "MS70" "MS80" "MS90"
    "UT50" "CD10" "CD20" "CD30" "CD40" "CD50" "CD60" "CD70" "CD71"
    "BA10" "BA20" "BA30" "BA40" "BA60"
    "DD40" "DD60" "DD80" "DD90"
  ) ;
  autoins:convictionRegion "UK" ;
  autoins:errorMessage "Please select a valid UK driving conviction code" .
```

---

## 7. VALIDATION AND TESTING

### 7.1 TTL Syntax Validation
```bash
# Validate TTL syntax
rapper -i turtle -o ntriples ontology/autoins.ttl > /dev/null
```

### 7.2 API Testing
```bash
# Test ontology API
curl -s http://localhost:3000/api/ontology | jq '.drivers.fields | length'

# Test specific field extraction
curl -s http://localhost:3000/api/ontology | jq '.drivers.fields[] | select(.property == "hasConvictions")'
```

### 7.3 Field Count Verification
```bash
# Count fields by domain
curl -s http://localhost:3000/api/ontology | jq '{
  drivers: (.drivers.fields | length),
  vehicles: (.vehicles.fields | length),
  claims: (.claims.fields | length)
}'
```

---

## 8. TROUBLESHOOTING

### 8.1 Common Issues
- **Empty fields returned**: Check regex patterns match TTL format
- **Missing labels**: Ensure `rdfs:label` is properly quoted
- **Wrong field types**: Verify `rdfs:range` and enumeration values
- **Fields in wrong section**: Check `rdfs:domain` assignment

### 8.2 Debug Mode
```go
// Add debug logging to TTL parser
fmt.Printf("Found property: %s, domain: %s, label: %s\n", propName, domain, label)
```

### 8.3 TTL Format Requirements
- Properties must end with ` .` (space + period)
- Strings must be quoted: `"value"`
- Boolean values: `"true"^^xsd:boolean`
- Lists: `("ITEM1" "ITEM2" "ITEM3")`

---

## 9. PERFORMANCE CONSIDERATIONS

### 9.1 Caching Strategy
- TTL parsed once at startup
- Consider caching parsed ontology in memory
- Reload on TTL file changes (file watcher)

### 9.2 Large Ontologies
- Use streaming parser for very large TTL files
- Consider splitting into multiple domain-specific TTL files
- Implement lazy loading for sections

---

## 10. FUTURE ENHANCEMENTS

### 10.1 SHACL Integration
```turtle
# SHACL validation shapes
autoins:DriverShape a sh:NodeShape ;
  sh:targetClass autoins:Driver ;
  sh:property [
    sh:path autoins:hasConvictions ;
    sh:datatype xsd:boolean ;
    sh:minCount 1 ;
    sh:maxCount 1 ;
  ] .
```

### 10.2 Multi-language Support
```turtle
autoins:firstName rdfs:label "first name"@en ;
autoins:firstName rdfs:label "Vorname"@de ;
autoins:firstName rdfs:label "prénom"@fr .
```

---

*This guide provides the technical foundation for maintaining and extending the TTL-based architecture in CLIENT-UX.*
