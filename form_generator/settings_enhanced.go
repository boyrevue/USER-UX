package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// FieldOption represents a select option
type FieldOption struct {
	Value string
	Label string
}

// SettingsField represents a field in the settings ontology
type SettingsField struct {
	Property               string
	Label                  string
	Type                   string
	Required               bool
	HelpText               string
	Pattern                string
	ErrorMessage           string
	Options                []FieldOption
	DefaultValue           string
	Min                    string
	Max                    string
	ConditionalRequirement string
	EncryptionRequired     bool
	EncryptionAlgorithm    string
	KeyDerivation          string
	SessionTimeout         string
	Requires2FA            bool
	PCICompliant           bool
	StorageType            string
	AccessLogging          string
	HealthCheckEnabled     bool
	HealthCheckFrequency   string
	AlertOnFailure         bool
	RetryAttempts          string
	ReadOnly               bool
	AlertOnChange          bool
}

// SettingsCategory represents a category of settings
type SettingsCategory struct {
	Name   string
	Label  string
	Fields []*SettingsField
	Order  int
}

// SettingsParser parses the settings ontology
type SettingsParser struct {
	Categories map[string]*SettingsCategory
}

// NewSettingsParser creates a new settings parser
func NewSettingsParser() *SettingsParser {
	return &SettingsParser{
		Categories: make(map[string]*SettingsCategory),
	}
}

// ParseSettingsTTL parses the settings TTL file
func (p *SettingsParser) ParseSettingsTTL(filepath string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("failed to open settings TTL file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentTriple strings.Builder
	var inTriple bool

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Handle multi-line triples
		if strings.Contains(line, " a owl:DatatypeProperty") {
			if inTriple {
				p.parseSettingsTriple(currentTriple.String())
			}
			currentTriple.Reset()
			currentTriple.WriteString(line)
			inTriple = true
		} else if inTriple {
			currentTriple.WriteString(" ")
			currentTriple.WriteString(line)

			if strings.HasSuffix(line, ".") {
				p.parseSettingsTriple(currentTriple.String())
				currentTriple.Reset()
				inTriple = false
			}
		}
	}

	if inTriple {
		p.parseSettingsTriple(currentTriple.String())
	}

	return nil
}

func (p *SettingsParser) parseSettingsTriple(triple string) {
	parts := strings.Fields(triple)
	if len(parts) < 3 {
		return
	}

	subject := parts[0]
	if !strings.Contains(subject, ":") {
		return
	}

	// Extract field name from subject
	fieldName := strings.TrimPrefix(subject, "banking:")
	fieldName = strings.TrimPrefix(fieldName, "comms:")
	fieldName = strings.TrimPrefix(fieldName, "security:")
	fieldName = strings.TrimPrefix(fieldName, "monitoring:")

	// Determine category based on prefix
	var category string
	if strings.HasPrefix(subject, "banking:") {
		if strings.Contains(triple, "card") {
			category = "creditcards"
		} else {
			category = "bankaccounts"
		}
	} else if strings.HasPrefix(subject, "comms:") {
		if strings.Contains(triple, "email") {
			category = "emailintegration"
		} else {
			category = "communicationchannels"
		}
	} else if strings.HasPrefix(subject, "security:") {
		category = "security"
	} else if strings.HasPrefix(subject, "monitoring:") {
		category = "monitoring"
	}

	if category == "" {
		return
	}

	// Initialize category if not exists
	if _, exists := p.Categories[category]; !exists {
		p.Categories[category] = &SettingsCategory{
			Name:  category,
			Label: p.getCategoryLabel(category),
			Order: p.getCategoryOrder(category),
		}
	}

	// Initialize field if not exists
	var field *SettingsField
	for _, f := range p.Categories[category].Fields {
		if f.Property == fieldName {
			field = f
			break
		}
	}
	if field == nil {
		field = &SettingsField{Property: fieldName}
		p.Categories[category].Fields = append(p.Categories[category].Fields, field)
	}

	// Parse properties
	propertyParts := strings.Split(triple, ";")
	for _, propPart := range propertyParts {
		propPart = strings.TrimSpace(propPart)
		if propPart == "" {
			continue
		}

		// Parse rdfs:label
		if strings.Contains(propPart, "rdfs:label") {
			labelMatch := regexp.MustCompile(`rdfs:label\s+"([^"]+)"`)
			if matches := labelMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.Label = matches[1]
			}
		}

		// Parse rdfs:range (field type)
		if strings.Contains(propPart, "rdfs:range") {
			rangeMatch := regexp.MustCompile(`rdfs:range\s+([^\s;]+)`)
			if matches := rangeMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.Type = p.parseSettingsFieldType(matches[1])
			}
		}

		// Parse validation pattern
		if strings.Contains(propPart, "autoins:validationPattern") {
			patternMatch := regexp.MustCompile(`autoins:validationPattern\s+"([^"]+)"`)
			if matches := patternMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.Pattern = matches[1]
			}
		}

		// Parse help text
		if strings.Contains(propPart, "autoins:formHelpText") {
			helpMatch := regexp.MustCompile(`autoins:formHelpText\s+"([^"]+)"`)
			if matches := helpMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.HelpText = matches[1]
			}
		}

		// Parse required field
		if strings.Contains(propPart, "autoins:isRequired") {
			requiredMatch := regexp.MustCompile(`autoins:isRequired\s+"([^"]+)"`)
			if matches := requiredMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.Required = matches[1] == "true"
			}
		}

		// Parse enumeration values
		if strings.Contains(propPart, "autoins:enumerationValues") {
			enumMatch := regexp.MustCompile(`autoins:enumerationValues\s+\(([^)]+)\)`)
			if matches := enumMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				enumValues := strings.TrimSpace(matches[1])
				enumValues = strings.ReplaceAll(enumValues, "\"", "")
				values := strings.Fields(enumValues)
				for _, value := range values {
					field.Options = append(field.Options, FieldOption{
						Value: value,
						Label: value,
					})
				}
			}
		}

		// Parse security properties
		if strings.Contains(propPart, "security:encryptionRequired") {
			encMatch := regexp.MustCompile(`security:encryptionRequired\s+"([^"]+)"`)
			if matches := encMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.EncryptionRequired = matches[1] == "true"
			}
		}

		if strings.Contains(propPart, "security:requires2FA") {
			tfaMatch := regexp.MustCompile(`security:requires2FA\s+"([^"]+)"`)
			if matches := tfaMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.Requires2FA = matches[1] == "true"
			}
		}

		if strings.Contains(propPart, "security:pciCompliant") {
			pciMatch := regexp.MustCompile(`security:pciCompliant\s+"([^"]+)"`)
			if matches := pciMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.PCICompliant = matches[1] == "true"
			}
		}

		// Parse monitoring properties
		if strings.Contains(propPart, "monitoring:healthCheckEnabled") {
			healthMatch := regexp.MustCompile(`monitoring:healthCheckEnabled\s+"([^"]+)"`)
			if matches := healthMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.HealthCheckEnabled = matches[1] == "true"
			}
		}

		if strings.Contains(propPart, "monitoring:alertOnFailure") {
			alertMatch := regexp.MustCompile(`monitoring:alertOnFailure\s+"([^"]+)"`)
			if matches := alertMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.AlertOnFailure = matches[1] == "true"
			}
		}

		if strings.Contains(propPart, "monitoring:readOnly") {
			readMatch := regexp.MustCompile(`monitoring:readOnly\s+"([^"]+)"`)
			if matches := readMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.ReadOnly = matches[1] == "true"
			}
		}

		if strings.Contains(propPart, "monitoring:alertOnChange") {
			changeMatch := regexp.MustCompile(`monitoring:alertOnChange\s+"([^"]+)"`)
			if matches := changeMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.AlertOnChange = matches[1] == "true"
			}
		}
	}
}

func (p *SettingsParser) parseSettingsFieldType(object string) string {
	switch {
	case strings.Contains(object, "xsd:string"):
		return "text"
	case strings.Contains(object, "xsd:integer"):
		return "number"
	case strings.Contains(object, "xsd:date"):
		return "date"
	case strings.Contains(object, "xsd:boolean"):
		return "checkbox"
	case strings.Contains(object, "xsd:dateTime"):
		return "datetime-local"
	default:
		return "text"
	}
}

func (p *SettingsParser) getCategoryLabel(category string) string {
	labels := map[string]string{
		"bankaccounts":          "Bank Accounts",
		"creditcards":           "Credit Cards",
		"communicationchannels": "Communication Channels",
		"emailintegration":      "Email Integration",
		"security":              "Security Settings",
		"monitoring":            "Monitoring & Alerts",
	}
	return labels[category]
}

func (p *SettingsParser) getCategoryOrder(category string) int {
	orders := map[string]int{
		"bankaccounts":          1,
		"creditcards":           2,
		"communicationchannels": 3,
		"emailintegration":      4,
		"security":              5,
		"monitoring":            6,
	}
	return orders[category]
}

func (p *SettingsParser) determineFieldTypes() {
	for _, category := range p.Categories {
		for _, field := range category.Fields {
			if len(field.Options) > 0 {
				field.Type = "select"
			}
		}
	}
}

// GenerateSettingsHTML generates the HTML form for settings
func (p *SettingsParser) GenerateSettingsHTML() (string, error) {
	p.determineFieldTypes()

	// Sort categories by order
	var sortedCategories []*SettingsCategory
	for _, category := range p.Categories {
		sortedCategories = append(sortedCategories, category)
	}
	sort.Slice(sortedCategories, func(i, j int) bool {
		return sortedCategories[i].Order < sortedCategories[j].Order
	})

	const tmpl = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Settings Configuration</title>
    <style>
        /* Sharp black/white/gray palette */
        :root {
            --pure-black: #000000;
            --pure-white: #ffffff;
            --gray-100: #f8f8f8;
            --gray-200: #e8e8e8;
            --gray-300: #d4d4d4;
            --gray-400: #a8a8a8;
            --gray-500: #7a7a7a;
            --gray-600: #5a5a5a;
            --gray-700: #3a3a3a;
            --gray-800: #2a2a2a;
            --gray-900: #1a1a1a;
            
            /* Semantic colors */
            --bg: var(--pure-white);
            --fg: var(--pure-black);
            --card-bg: var(--pure-white);
            --card-border: var(--gray-300);
            --input-bg: var(--gray-100);
            --input-border: var(--gray-400);
            --input-focus: var(--pure-black);
            --button-bg: var(--pure-black);
            --button-fg: var(--pure-white);
            --muted: var(--gray-600);
            --shadow: rgba(0, 0, 0, 0.1);
            --shadow-strong: rgba(0, 0, 0, 0.2);
            --accent-green: #00ff41;
            --accent-yellow: #ffff00;
            --accent-red: #ff0040;
            --accent-blue: #0080ff;
        }

        /* Dark mode */
        @media (prefers-color-scheme: dark) {
            :root {
                --bg: var(--pure-black);
                --fg: var(--pure-white);
                --card-bg: var(--gray-900);
                --card-border: var(--gray-700);
                --input-bg: var(--gray-800);
                --input-border: var(--gray-600);
                --input-focus: var(--pure-white);
                --button-bg: var(--pure-white);
                --button-fg: var(--pure-black);
                --muted: var(--gray-400);
                --shadow: rgba(255, 255, 255, 0.05);
                --shadow-strong: rgba(255, 255, 255, 0.1);
            }
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--pure-black);
            color: var(--pure-white);
            line-height: 1.6;
            font-size: 14px;
        }

        .container {
            max-width: 1200px;
            margin: 0 auto;
            padding: 20px;
        }

        .header {
            text-align: center;
            margin-bottom: 40px;
            padding: 20px;
            background: var(--gray-900);
            border-radius: 8px;
            border: 1px solid var(--gray-600);
        }

        .header h1 {
            font-size: 2.5rem;
            font-weight: 700;
            color: var(--accent-green);
            margin-bottom: 10px;
        }

        .header p {
            color: var(--gray-300);
            font-size: 1.1rem;
        }

        .settings-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
            gap: 20px;
            margin-bottom: 40px;
        }

        .settings-card {
            background: var(--gray-900);
            border: 1px solid var(--gray-600);
            border-radius: 8px;
            padding: 25px;
            transition: all 0.3s ease;
        }

        .settings-card:hover {
            border-color: var(--accent-green);
            box-shadow: 0 0 20px rgba(0, 255, 65, 0.1);
        }

        .card-header {
            display: flex;
            align-items: center;
            margin-bottom: 20px;
            padding-bottom: 15px;
            border-bottom: 1px solid var(--gray-600);
        }

        .card-icon {
            font-size: 24px;
            margin-right: 15px;
        }

        .card-title {
            font-size: 1.5rem;
            font-weight: 600;
            color: var(--accent-green);
        }

        .form-grid {
            display: grid;
            gap: 15px;
        }

        .form-field {
            display: flex;
            flex-direction: column;
        }

        .form-field.full-width {
            grid-column: 1 / -1;
        }

        .form-label {
            font-weight: 600;
            margin-bottom: 8px;
            color: var(--gray-200);
            font-size: 0.9rem;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .required {
            color: var(--accent-red);
            margin-left: 4px;
        }

        .form-input, .form-select, .form-textarea {
            background: var(--gray-800);
            border: 1px solid var(--gray-600);
            border-radius: 4px;
            padding: 12px 15px;
            color: var(--pure-white);
            font-size: 14px;
            transition: all 0.3s ease;
        }

        .form-input:focus, .form-select:focus, .form-textarea:focus {
            outline: none;
            border-color: var(--accent-green);
            box-shadow: 0 0 10px rgba(0, 255, 65, 0.3);
        }

        .form-input[readonly] {
            background: var(--gray-700);
            color: var(--gray-400);
            cursor: not-allowed;
        }

        .help-text {
            font-size: 0.8rem;
            color: var(--gray-400);
            margin-top: 5px;
            font-style: italic;
        }

        .security-badge {
            display: inline-block;
            background: var(--accent-green);
            color: var(--pure-black);
            padding: 2px 8px;
            border-radius: 12px;
            font-size: 0.7rem;
            font-weight: 600;
            margin-left: 8px;
        }

        .pci-badge {
            background: var(--accent-blue);
        }

        .encryption-badge {
            background: var(--accent-yellow);
            color: var(--pure-black);
        }

        .add-button {
            background: var(--accent-green);
            color: var(--pure-black);
            border: none;
            padding: 12px 20px;
            border-radius: 4px;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
            margin-top: 15px;
        }

        .add-button:hover {
            background: var(--pure-white);
            transform: translateY(-2px);
        }

        .remove-button {
            background: var(--accent-red);
            color: var(--pure-white);
            border: none;
            padding: 8px 12px;
            border-radius: 4px;
            font-size: 0.8rem;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .remove-button:hover {
            background: var(--pure-white);
            color: var(--accent-red);
        }

        .instance-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 15px;
            padding: 10px;
            background: var(--gray-800);
            border-radius: 4px;
            border-left: 4px solid var(--accent-green);
        }

        .instance-title {
            font-weight: 600;
            color: var(--accent-green);
        }

        .health-status {
            display: inline-block;
            padding: 4px 8px;
            border-radius: 12px;
            font-size: 0.7rem;
            font-weight: 600;
        }

        .status-healthy {
            background: var(--accent-green);
            color: var(--pure-black);
        }

        .status-warning {
            background: var(--accent-yellow);
            color: var(--pure-black);
        }

        .status-error {
            background: var(--accent-red);
            color: var(--pure-white);
        }

        .status-unknown {
            background: var(--gray-500);
            color: var(--pure-white);
        }

        .submit-section {
            text-align: center;
            margin-top: 40px;
            padding: 30px;
            background: var(--gray-900);
            border-radius: 8px;
            border: 1px solid var(--gray-600);
        }

        .submit-button {
            background: var(--accent-green);
            color: var(--pure-black);
            border: none;
            padding: 15px 30px;
            border-radius: 4px;
            font-size: 1.1rem;
            font-weight: 600;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .submit-button:hover {
            background: var(--pure-white);
            transform: translateY(-2px);
        }

        .alert {
            padding: 15px;
            border-radius: 4px;
            margin-bottom: 20px;
            border-left: 4px solid;
        }

        .alert-info {
            background: rgba(0, 128, 255, 0.1);
            border-color: var(--accent-blue);
            color: var(--accent-blue);
        }

        .alert-warning {
            background: rgba(255, 255, 0, 0.1);
            border-color: var(--accent-yellow);
            color: var(--accent-yellow);
        }

        .alert-error {
            background: rgba(255, 0, 64, 0.1);
            border-color: var(--accent-red);
            color: var(--accent-red);
        }

        .alert-success {
            background: rgba(0, 255, 65, 0.1);
            border-color: var(--accent-green);
            color: var(--accent-green);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔧 Settings Configuration</h1>
            <p>Manage your bank accounts, credit cards, communication channels, and security settings</p>
        </div>

        <div class="alert alert-info">
            <strong>🔐 Security Notice:</strong> All sensitive data will be encrypted and stored in PCI-compliant storage. 
            Access attempts and changes are logged for audit purposes.
        </div>

        <form id="settingsForm">
            {{range $category := .}}
            <div class="settings-card">
                <div class="card-header">
                    <div class="card-icon">
                        {{if eq $category.Name "bankaccounts"}}🏦{{else if eq $category.Name "creditcards"}}💳{{else if eq $category.Name "communicationchannels"}}📱{{else if eq $category.Name "emailintegration"}}📧{{else if eq $category.Name "security"}}🔒{{else if eq $category.Name "monitoring"}}📊{{end}}
                    </div>
                    <div class="card-title">{{$category.Label}}</div>
                </div>

                <div class="form-grid">
                    {{range $field := $category.Fields}}
                    <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                        <label class="form-label">
                            {{$field.Label}}
                            {{if $field.Required}}<span class="required">*</span>{{end}}
                            {{if $field.EncryptionRequired}}<span class="security-badge encryption-badge">🔐 Encrypted</span>{{end}}
                            {{if $field.PCICompliant}}<span class="security-badge pci-badge">PCI</span>{{end}}
                            {{if $field.Requires2FA}}<span class="security-badge">2FA</span>{{end}}
                        </label>
                        
                        {{if eq $field.Type "text"}}
                        <input type="text" 
                               class="form-input" 
                               name="{{$category.Name}}_{{$field.Property}}" 
                               placeholder="{{$field.Label}}"
                               {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                               {{if $field.Required}}required{{end}}
                               {{if $field.ReadOnly}}readonly{{end}}>
                        
                        {{else if eq $field.Type "number"}}
                        <input type="number" 
                               class="form-input" 
                               name="{{$category.Name}}_{{$field.Property}}"
                               {{if $field.Min}}min="{{$field.Min}}"{{end}}
                               {{if $field.Max}}max="{{$field.Max}}"{{end}}
                               {{if $field.Required}}required{{end}}
                               {{if $field.ReadOnly}}readonly{{end}}>
                        
                        {{else if eq $field.Type "date"}}
                        <input type="date" 
                               class="form-input" 
                               name="{{$category.Name}}_{{$field.Property}}"
                               {{if $field.Required}}required{{end}}
                               {{if $field.ReadOnly}}readonly{{end}}>
                        
                        {{else if eq $field.Type "datetime-local"}}
                        <input type="datetime-local" 
                               class="form-input" 
                               name="{{$category.Name}}_{{$field.Property}}"
                               {{if $field.Required}}required{{end}}
                               {{if $field.ReadOnly}}readonly{{end}}>
                        
                        {{else if eq $field.Type "email"}}
                        <input type="email" 
                               class="form-input" 
                               name="{{$category.Name}}_{{$field.Property}}"
                               {{if $field.Required}}required{{end}}
                               {{if $field.ReadOnly}}readonly{{end}}>
                        
                        {{else if eq $field.Type "checkbox"}}
                        <input type="checkbox" 
                               class="form-input" 
                               name="{{$category.Name}}_{{$field.Property}}"
                               {{if $field.Required}}required{{end}}
                               {{if $field.ReadOnly}}readonly{{end}}>
                        
                        {{else if eq $field.Type "select"}}
                        <select class="form-select" 
                                name="{{$category.Name}}_{{$field.Property}}" 
                                {{if $field.Required}}required{{end}}
                                {{if $field.ReadOnly}}readonly{{end}}>
                            <option value="">Select {{$field.Label}}</option>
                            {{range $option := $field.Options}}
                            <option value="{{$option.Value}}">{{$option.Label}}</option>
                            {{end}}
                        </select>
                        
                        {{else if eq $field.Type "textarea"}}
                        <textarea class="form-textarea" 
                                  name="{{$category.Name}}_{{$field.Property}}" 
                                  rows="3" 
                                  placeholder="{{$field.Label}}"
                                  {{if $field.Required}}required{{end}}
                                  {{if $field.ReadOnly}}readonly{{end}}></textarea>
                        
                        {{else}}
                        <input type="text" 
                               class="form-input" 
                               name="{{$category.Name}}_{{$field.Property}}" 
                               placeholder="{{$field.Label}}"
                               {{if $field.Required}}required{{end}}
                               {{if $field.ReadOnly}}readonly{{end}}>
                        {{end}}
                        
                        {{if $field.HelpText}}
                        <div class="help-text">{{$field.HelpText}}</div>
                        {{end}}
                    </div>
                    {{end}}
                </div>

                {{if or (eq $category.Name "bankaccounts") (eq $category.Name "creditcards") (eq $category.Name "communicationchannels")}}
                <button type="button" class="add-button" onclick="addInstance('{{$category.Name}}')">
                    ➕ Add {{$category.Label}} Instance
                </button>
                {{end}}
            </div>
            {{end}}

            <div class="submit-section">
                <button type="submit" class="submit-button">
                    💾 Save Settings Configuration
                </button>
            </div>
        </form>
    </div>

    <script>
        let instanceCounts = {
            bankaccounts: 0,
            creditcards: 0,
            communicationchannels: 0
        };

        function addInstance(category) {
            instanceCounts[category]++;
            const card = document.querySelector('.settings-card');
            const formGrid = card.querySelector('.form-grid');
            
            const instanceDiv = document.createElement('div');
            instanceDiv.className = 'instance-header';
            instanceDiv.innerHTML = 
                '<div class="instance-title">' + category + ' Instance #' + instanceCounts[category] + '</div>' +
                '<button type="button" class="remove-button" onclick="removeInstance(this)">🗑️ Remove</button>';
            
            formGrid.appendChild(instanceDiv);
            
            // Clone form fields for this instance
            const fields = card.querySelectorAll('.form-field');
            fields.forEach(field => {
                const clonedField = field.cloneNode(true);
                const inputs = clonedField.querySelectorAll('input, select, textarea');
                inputs.forEach(input => {
                    input.name = input.name + '_' + instanceCounts[category];
                    input.value = '';
                });
                formGrid.appendChild(clonedField);
            });
        }

        function removeInstance(button) {
            const instanceDiv = button.closest('.instance-header');
            const card = instanceDiv.closest('.settings-card');
            const formGrid = card.querySelector('.form-grid');
            
            // Remove the instance header
            instanceDiv.remove();
            
            // Remove all form fields after this instance
            const fields = formGrid.querySelectorAll('.form-field');
            const instanceIndex = Array.from(formGrid.children).indexOf(instanceDiv);
            
            for (let i = instanceIndex; i < fields.length; i++) {
                fields[i].remove();
            }
        }

        document.getElementById('settingsForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            // Here you would implement secure settings storage
            alert('🔐 Settings would be encrypted and stored securely here.\n\nIn production, this would:\n• Encrypt all sensitive data with AES-256-GCM\n• Store in PCI-compliant vault\n• Log all access attempts\n• Implement proper session management\n• Use secure key derivation (PBKDF2)\n• Set up health monitoring and alerts');
        });

        // Initialize health status indicators
        function updateHealthStatus() {
            const statusElements = document.querySelectorAll('[name*="healthStatus"]');
            statusElements.forEach(element => {
                const status = element.value || 'Unknown';
                const badge = document.createElement('span');
                badge.className = 'health-status status-' + status.toLowerCase();
                badge.textContent = status;
                element.parentNode.appendChild(badge);
            });
        }

        // Auto-add one instance of each type
        window.addEventListener('load', function() {
            addInstance('bankaccounts');
            addInstance('creditcards');
            addInstance('communicationchannels');
            updateHealthStatus();
        });
    </script>
</body>
</html>`

	tmplParsed, err := template.New("settings").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %v", err)
	}

	var buf strings.Builder
	err = tmplParsed.Execute(&buf, sortedCategories)
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return buf.String(), nil
}

func main() {
	var ttlPath = flag.String("ttl", "../ontology/settings.ttl", "Path to settings TTL file")
	var outputPath = flag.String("output", "settings_form.html", "Output HTML file path")
	flag.Parse()

	parser := NewSettingsParser()

	// Parse TTL file
	err := parser.ParseSettingsTTL(*ttlPath)
	if err != nil {
		fmt.Printf("❌ Error parsing settings TTL file: %v\n", err)
		os.Exit(1)
	}

	// Generate HTML
	html, err := parser.GenerateSettingsHTML()
	if err != nil {
		fmt.Printf("❌ Error generating settings HTML: %v\n", err)
		os.Exit(1)
	}

	// Write to file
	err = os.WriteFile(*outputPath, []byte(html), 0644)
	if err != nil {
		fmt.Printf("❌ Error writing settings HTML file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Generated settings form from TTL ontology: %s\n", *outputPath)
	fmt.Printf("📊 Categories found: %d\n", len(parser.Categories))
	for _, category := range parser.Categories {
		fmt.Printf("   • %s: %d fields\n", category.Label, len(category.Fields))
	}
}
