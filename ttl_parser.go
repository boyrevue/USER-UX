package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// TTL Parser structures
type TTLField struct {
	Property     string
	Label        string
	Type         string
	Required     bool
	Pattern      string
	HelpText     string
	Options      []FieldOption
	Min          string
	Max          string
	DefaultValue string
	Conditional  string
}

type FieldOption struct {
	Value string
	Label string
}

type TTLCategory struct {
	ID          string
	Title       string
	Icon        string
	Order       int
	Description string
	Fields      []TTLField
}

type TTLParser struct {
	Categories map[string]*TTLCategory
	Fields     map[string]*TTLField
	Prefixes   map[string]string
}

func NewTTLParser() *TTLParser {
	return &TTLParser{
		Categories: make(map[string]*TTLCategory),
		Fields:     make(map[string]*TTLField),
		Prefixes:   make(map[string]string),
	}
}

func (p *TTLParser) ParseTTLFile(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open TTL file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// Parse prefixes
		if strings.HasPrefix(line, "@prefix") {
			p.parsePrefix(line)
			continue
		}

		// Parse triples
		if strings.Contains(line, " ;") || strings.Contains(line, " .") {
			p.parseTriple(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading TTL file: %v", err)
	}

	return nil
}

func (p *TTLParser) parsePrefix(line string) {
	// Parse @prefix autoins: <https://autoins.example/ontology#> .
	re := regexp.MustCompile(`@prefix\s+(\w+):\s+<([^>]+)>`)
	matches := re.FindStringSubmatch(line)
	if len(matches) == 3 {
		p.Prefixes[matches[1]] = matches[2]
	}
}

func (p *TTLParser) parseTriple(line string) {
	// Parse triples like: autoins:firstName a owl:DatatypeProperty ;
	parts := strings.Split(line, " ")
	if len(parts) < 3 {
		return
	}

	subject := parts[0]
	predicate := parts[1]
	object := strings.Join(parts[2:], " ")

	// Remove trailing semicolon or period
	object = strings.TrimSuffix(object, ";")
	object = strings.TrimSuffix(object, ".")

	// Parse field definitions
	if strings.Contains(subject, "autoins:") && strings.Contains(predicate, "a") {
		fieldName := strings.TrimPrefix(subject, "autoins:")
		if _, exists := p.Fields[fieldName]; !exists {
			p.Fields[fieldName] = &TTLField{Property: fieldName}
		}
	}

	// Parse field properties
	if strings.Contains(predicate, "rdfs:label") {
		fieldName := strings.TrimPrefix(subject, "autoins:")
		if field, exists := p.Fields[fieldName]; exists {
			field.Label = strings.Trim(object, `"`)
		}
	}

	// Parse field types
	if strings.Contains(predicate, "rdfs:range") {
		fieldName := strings.TrimPrefix(subject, "autoins:")
		if field, exists := p.Fields[fieldName]; exists {
			field.Type = p.parseFieldType(object)
		}
	}

	// Parse validation patterns
	if strings.Contains(predicate, "autoins:validationPattern") {
		fieldName := strings.TrimPrefix(subject, "autoins:")
		if field, exists := p.Fields[fieldName]; exists {
			field.Pattern = strings.Trim(object, `"`)
		}
	}

	// Parse help text
	if strings.Contains(predicate, "autoins:formHelpText") {
		fieldName := strings.TrimPrefix(subject, "autoins:")
		if field, exists := p.Fields[fieldName]; exists {
			field.HelpText = strings.Trim(object, `"`)
		}
	}

	// Parse required fields
	if strings.Contains(predicate, "autoins:isRequired") {
		fieldName := strings.TrimPrefix(subject, "autoins:")
		if field, exists := p.Fields[fieldName]; exists {
			field.Required = object == "true"
		}
	}
}

func (p *TTLParser) parseFieldType(object string) string {
	switch {
	case strings.Contains(object, "xsd:string"):
		return "text"
	case strings.Contains(object, "xsd:integer"):
		return "number"
	case strings.Contains(object, "xsd:date"):
		return "date"
	case strings.Contains(object, "xsd:boolean"):
		return "checkbox"
	case strings.Contains(object, "autoins:"):
		return "select"
	default:
		return "text"
	}
}

func (p *TTLParser) OrganizeIntoCategories() {
	// Define category mappings based on field prefixes
	categoryMappings := map[string]string{
		"firstName":     "drivers",
		"lastName":      "drivers",
		"dateOfBirth":   "drivers",
		"licenseNumber": "drivers",
		"licenseType":   "drivers",
		"relationship":  "drivers",

		"make":            "vehicle",
		"model":           "vehicle",
		"year":            "vehicle",
		"vin":             "vehicle",
		"color":           "vehicle",
		"bodyType":        "vehicle",
		"engineSize":      "vehicle",
		"fuelType":        "vehicle",
		"annualMileage":   "vehicle",
		"primaryUse":      "vehicle",
		"garageLocation":  "vehicle",
		"antiTheftDevice": "vehicle",

		"liabilityCoverage":       "coverage",
		"collisionDeductible":     "coverage",
		"comprehensiveDeductible": "coverage",
		"medicalPayments":         "coverage",
		"uninsuredMotorist":       "coverage",
		"rentalReimbursement":     "coverage",
		"roadsideAssistance":      "coverage",
		"gapInsurance":            "coverage",

		"atFaultAccidents":       "claims",
		"notAtFaultAccidents":    "claims",
		"trafficViolations":      "claims",
		"duiConvictions":         "claims",
		"licenseSuspensions":     "claims",
		"insuranceCancellations": "claims",

		"paymentMethod":    "payment",
		"paymentFrequency": "payment",
		"cardNumber":       "payment",
		"expiryDate":       "payment",
		"cvv":              "payment",
		"billingAddress":   "payment",
		"autoPay":          "payment",

		"communicationPreference":      "preferences",
		"paperlessDocuments":           "preferences",
		"marketingCommunications":      "preferences",
		"emergencyContactName":         "preferences",
		"emergencyContactPhone":        "preferences",
		"emergencyContactRelationship": "preferences",
		"preferredAgent":               "preferences",

		"policyType":       "summary",
		"policyTerm":       "summary",
		"effectiveDate":    "summary",
		"expirationDate":   "summary",
		"estimatedPremium": "summary",
		"discountsApplied": "summary",
	}

	// Initialize categories
	categoryConfigs := map[string]struct {
		title string
		icon  string
		order int
	}{
		"drivers":     {"Driver Details", "👥", 1},
		"vehicle":     {"Vehicle Information", "🚗", 2},
		"coverage":    {"Coverage Options", "🛡️", 3},
		"claims":      {"Claims History", "📋", 4},
		"payment":     {"Payment Information", "💳", 5},
		"preferences": {"Preferences", "⚙️", 6},
		"summary":     {"Summary", "📊", 7},
	}

	for categoryID, config := range categoryConfigs {
		p.Categories[categoryID] = &TTLCategory{
			ID:          categoryID,
			Title:       config.title,
			Icon:        config.icon,
			Order:       config.order,
			Description: fmt.Sprintf("Configure %s", config.title),
			Fields:      []TTLField{},
		}
	}

	// Organize fields into categories
	for fieldName, field := range p.Fields {
		if categoryID, exists := categoryMappings[fieldName]; exists {
			if category, exists := p.Categories[categoryID]; exists {
				category.Fields = append(category.Fields, *field)
			}
		}
	}

	// Sort fields within each category
	for _, category := range p.Categories {
		sort.Slice(category.Fields, func(i, j int) bool {
			return category.Fields[i].Property < category.Fields[j].Property
		})
	}
}

func (p *TTLParser) GenerateHTMLForms() (string, error) {
	// Sort categories by order
	var sortedCategories []*TTLCategory
	for _, category := range p.Categories {
		sortedCategories = append(sortedCategories, category)
	}
	sort.Slice(sortedCategories, func(i, j int) bool {
		return sortedCategories[i].Order < sortedCategories[j].Order
	})

	// Generate HTML template
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Ontology-Driven Insurance Forms</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        * { box-sizing: border-box; }
        body { 
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; 
            margin: 0; 
            background: linear-gradient(135deg, #333333 0%, #666666 100%);
            min-height: 100vh;
            padding: 20px;
            font-size: 18px;
        }
        .container { 
            max-width: 1400px; 
            margin: 0 auto; 
            background: white; 
            border-radius: 15px; 
            box-shadow: 0 20px 40px rgba(0,0,0,0.3);
            overflow: hidden;
        }
        .header { 
            background: linear-gradient(135deg, #333333, #666666); 
            color: white; 
            padding: 30px; 
            text-align: center;
        }
        .header h1 { margin: 0; font-size: 2.8em; font-weight: 300; }
        .header p { margin: 10px 0 0 0; opacity: 0.9; font-size: 1.3em; }
        
        .main-content {
            display: flex;
            flex-direction: column;
            min-height: 600px;
        }
        
        .left-panel {
            padding: 20px;
            background: #fafafa;
            overflow-y: auto;
            max-height: 70vh;
        }

        .category-tabs {
            display: flex;
            gap: 5px;
            margin-bottom: 20px;
            flex-wrap: wrap;
            background: white;
            padding: 15px;
            border-radius: 10px;
            border: 1px solid #cccccc;
        }

        .category-tab {
            padding: 12px 16px;
            background: #f0f0f0;
            border: 1px solid #cccccc;
            border-radius: 8px;
            cursor: pointer;
            font-size: 14px;
            transition: all 0.3s;
            min-width: 100px;
            text-align: center;
            font-weight: 500;
        }

        .category-tab.active {
            background: linear-gradient(135deg, #666666, #333333);
            color: white;
            border-color: #333333;
        }

        .category-tab:hover:not(.active) {
            background: #e0e0e0;
        }

        .category-content {
            display: none;
        }

        .category-content.active {
            display: block;
        }

        .form-section {
            background: white;
            padding: 20px;
            border-radius: 15px;
            border: 2px solid #cccccc;
            margin: 20px 0;
        }

        .form-section h3 {
            margin: 0 0 20px 0;
            font-size: 1.4em;
            color: #333333;
        }

        .form-grid {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 15px;
        }

        .form-field {
            display: flex;
            flex-direction: column;
        }

        .form-field.full-width {
            grid-column: 1 / -1;
        }

        .form-field label {
            font-weight: bold;
            margin-bottom: 5px;
            color: #333333;
        }

        .form-field input,
        .form-field select,
        .form-field textarea {
            padding: 12px 15px;
            border: 2px solid #cccccc;
            border-radius: 8px;
            font-size: 16px;
            background: white;
            color: #333333;
        }

        .form-field input:focus,
        .form-field select:focus,
        .form-field textarea:focus {
            outline: none;
            border-color: #666666;
        }

        .form-field .help-text {
            font-size: 12px;
            color: #666666;
            margin-top: 4px;
        }

        .required {
            color: #ff4444;
        }

        @media (max-width: 768px) {
            .form-grid {
                grid-template-columns: 1fr;
            }
            .category-tabs {
                flex-direction: column;
            }
            .category-tab {
                min-width: auto;
            }
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚗📄 Ontology-Driven Insurance Forms</h1>
            <p>Dynamically generated from TTL ontology</p>
        </div>

        <div class="main-content">
            <div class="left-panel">
                <!-- Category Tabs -->
                <div class="category-tabs" id="categoryTabs">
                    {{range $index, $category := .Categories}}
                    <div class="category-tab {{if eq $index 0}}active{{end}}" onclick="switchCategory('{{$category.ID}}')">
                        {{$category.Icon}} {{$category.Title}}
                    </div>
                    {{end}}
                </div>

                <!-- Category Contents -->
                {{range $index, $category := .Categories}}
                <div class="category-content {{if eq $index 0}}active{{end}}" id="{{$category.ID}}Content">
                    <div class="form-section">
                        <h3>{{$category.Icon}} {{$category.Title}}</h3>
                        <div class="form-grid">
                            {{range $field := $category.Fields}}
                            <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                <label>
                                    {{$field.Label}}
                                    {{if $field.Required}}<span class="required">*</span>{{end}}
                                </label>
                                
                                {{if eq $field.Type "text"}}
                                <input type="text" 
                                       name="{{$field.Property}}" 
                                       placeholder="{{$field.Label}}"
                                       {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                       {{if $field.Required}}required{{end}}>
                                
                                {{else if eq $field.Type "number"}}
                                <input type="number" 
                                       name="{{$field.Property}}"
                                       {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                       {{if $field.Max}}max="{{$field.Max}}"{{end}}
                                       {{if $field.Required}}required{{end}}>
                                
                                {{else if eq $field.Type "date"}}
                                <input type="date" 
                                       name="{{$field.Property}}"
                                       {{if $field.Required}}required{{end}}>
                                
                                {{else if eq $field.Type "email"}}
                                <input type="email" 
                                       name="{{$field.Property}}"
                                       {{if $field.Required}}required{{end}}>
                                
                                {{else if eq $field.Type "tel"}}
                                <input type="tel" 
                                       name="{{$field.Property}}"
                                       {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                       {{if $field.Required}}required{{end}}>
                                
                                {{else if eq $field.Type "select"}}
                                <select name="{{$field.Property}}" {{if $field.Required}}required{{end}}>
                                    <option value="">Select {{$field.Label}}</option>
                                    {{range $option := $field.Options}}
                                    <option value="{{$option.Value}}">{{$option.Label}}</option>
                                    {{end}}
                                </select>
                                
                                {{else if eq $field.Type "textarea"}}
                                <textarea name="{{$field.Property}}" 
                                          rows="3" 
                                          placeholder="{{$field.Label}}"
                                          {{if $field.Required}}required{{end}}></textarea>
                                
                                {{else if eq $field.Type "checkbox"}}
                                <input type="checkbox" 
                                       name="{{$field.Property}}"
                                       {{if $field.Required}}required{{end}}>
                                
                                {{else}}
                                <input type="text" 
                                       name="{{$field.Property}}" 
                                       placeholder="{{$field.Label}}"
                                       {{if $field.Required}}required{{end}}>
                                {{end}}
                                
                                {{if $field.HelpText}}
                                <div class="help-text">{{$field.HelpText}}</div>
                                {{end}}
                            </div>
                            {{end}}
                        </div>
                    </div>
                </div>
                {{end}}
            </div>
        </div>
    </div>

    <script>
        function switchCategory(category) {
            // Hide all category contents
            const categoryContents = document.querySelectorAll('.category-content');
            categoryContents.forEach(content => {
                content.classList.remove('active');
            });

            // Show the selected category content
            const selectedContent = document.getElementById(category + 'Content');
            if (selectedContent) {
                selectedContent.classList.add('active');
            }

            // Update tab states
            const tabs = document.querySelectorAll('.category-tab');
            tabs.forEach(tab => {
                tab.classList.remove('active');
            });

            const activeTab = document.querySelector("[onclick=\"switchCategory('" + category + "')\"]");
            if (activeTab) {
                activeTab.classList.add('active');
            }
        }
    </script>
</body>
</html>`

	// Parse and execute template
	t, err := template.New("forms").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var buf strings.Builder
	err = t.Execute(&buf, map[string]interface{}{
		"Categories": sortedCategories,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return buf.String(), nil
}

// Main function to generate forms from TTL
func GenerateFormsFromTTL(ttlPath, outputPath string) error {
	parser := NewTTLParser()

	// Parse TTL file
	err := parser.ParseTTLFile(ttlPath)
	if err != nil {
		return fmt.Errorf("failed to parse TTL file: %v", err)
	}

	// Organize fields into categories
	parser.OrganizeIntoCategories()

	// Generate HTML
	html, err := parser.GenerateHTMLForms()
	if err != nil {
		return fmt.Errorf("failed to generate HTML: %v", err)
	}

	// Write to file
	err = os.WriteFile(outputPath, []byte(html), 0644)
	if err != nil {
		return fmt.Errorf("failed to write HTML file: %v", err)
	}

	log.Printf("✅ Generated forms from TTL ontology: %s", outputPath)
	return nil
}
