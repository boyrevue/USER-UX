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
            --accent-green: #00ff41;
            --accent-yellow: #ffff00;
            --accent-red: #ff0040;
            --accent-blue: #0080ff;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            background: var(--pure-black);
            color: var(--pure-white);
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            line-height: 1.6;
            font-size: 14px;
        }

        .container {
            max-width: 1400px;
            margin: 0 auto;
            background: var(--gray-900);
            border: 1px solid var(--gray-600);
            border-radius: 0;
            box-shadow: 0 20px 40px rgba(0, 0, 0, 0.8);
            overflow: hidden;
        }

        .header {
            background: var(--gray-800);
            color: var(--pure-white);
            padding: 30px;
            text-align: center;
            position: relative;
            border-bottom: 2px solid var(--gray-600);
            box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
        }

        .header h1 {
            margin: 0;
            font-size: 24px;
            font-weight: 700;
            letter-spacing: -0.02em;
        }

        .header p {
            margin: 10px 0 0 0;
            color: var(--gray-300);
            font-size: 14px;
        }

        .main-content {
            display: grid;
            grid-template-columns: 300px 1fr;
            min-height: 600px;
            background: var(--gray-900);
        }

        .left-panel {
            background: var(--gray-800);
            border-right: 2px solid var(--gray-600);
            padding: 20px 0;
        }

        .category-tabs {
            display: flex;
            flex-direction: column;
            gap: 2px;
        }

        .category-tab {
            background: var(--gray-700);
            color: var(--gray-300);
            border: none;
            padding: 15px 20px;
            text-align: left;
            cursor: pointer;
            font-size: 13px;
            font-weight: 600;
            letter-spacing: 0.02em;
            transition: all 0.15s ease;
            filter: grayscale(100%);
        }

        .category-tab.active {
            background: var(--pure-white);
            color: var(--pure-black);
            font-weight: 700;
            filter: grayscale(0%);
        }

        .category-tab:hover:not(.active) {
            background: var(--gray-600);
            color: var(--pure-white);
            transform: translateY(-1px);
        }

        .category-content {
            display: none;
        }

        .category-content.active {
            display: block;
        }

        .right-panel {
            background: var(--gray-900);
            padding: 30px;
            overflow-y: auto;
        }

        .form-section {
            margin-bottom: 30px;
        }

        .form-section h3 {
            color: var(--pure-white);
            font-size: 18px;
            margin-bottom: 20px;
            font-weight: 700;
            letter-spacing: -0.01em;
        }

        .form-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
            gap: 20px;
        }

        .form-field {
            display: flex;
            flex-direction: column;
        }

        .form-field.full-width {
            grid-column: 1 / -1;
        }

        .form-field label {
            color: var(--pure-white);
            margin-bottom: 8px;
            font-size: 12px;
            font-weight: 700;
            text-transform: uppercase;
            letter-spacing: 0.05em;
        }

        .form-field input,
        .form-field select,
        .form-field textarea {
            width: 100%;
            background: var(--gray-800);
            color: var(--pure-white);
            border: 2px solid var(--gray-600);
            border-radius: 0;
            padding: 10px 12px;
            font-size: 14px;
            font-family: inherit;
            transition: all 0.15s ease;
        }

        .form-field input:focus,
        .form-field select:focus,
        .form-field textarea:focus {
            outline: none;
            border-color: var(--accent-green);
            box-shadow: 0 0 10px rgba(0, 255, 65, 0.3);
        }

        .form-field .help-text {
            color: var(--gray-400);
            font-size: 11px;
            margin-top: 4px;
            font-style: italic;
        }

        .required {
            color: var(--accent-red);
            font-weight: 700;
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

        /* Floating modals */
        .floating-modal {
            position: fixed;
            background: #2a2a2a;
            border: 1px solid #404040;
            border-radius: 8px;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.8);
            z-index: 1000;
            min-width: 300px;
            min-height: 200px;
            resize: both;
            overflow: hidden;
        }

        .modal-header {
            background: #333333;
            padding: 12px 16px;
            border-bottom: 1px solid #404040;
            cursor: move;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .modal-title {
            font-size: 14px;
            font-weight: 600;
            color: #e0e0e0;
        }

        .modal-controls {
            display: flex;
            gap: 8px;
        }

        .modal-btn {
            width: 20px;
            height: 20px;
            border: none;
            border-radius: 4px;
            background: #404040;
            color: #808080;
            cursor: pointer;
            display: flex;
            align-items: center;
            justify-content: center;
            font-size: 12px;
            transition: all 0.3s ease;
        }

        .modal-btn:hover {
            background: #505050;
            color: #e0e0e0;
        }

        .modal-content {
            padding: 16px;
            height: calc(100% - 50px);
            overflow-y: auto;
        }

        /* Modal positioning */
        #uploadModal {
            top: 50px;
            right: 50px;
            width: 400px;
            height: 500px;
        }

        #chatModal {
            bottom: 50px;
            right: 50px;
            width: 350px;
            height: 400px;
        }

        /* Floating Action Buttons */
        .floating-actions {
            position: fixed;
            bottom: 20px;
            left: 20px;
            display: flex;
            gap: 10px;
            z-index: 999;
        }

        .action-btn {
            background: var(--accent-green);
            color: var(--pure-black);
            border: none;
            padding: 12px 16px;
            border-radius: 6px;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.3s ease;
        }

        .action-btn:hover {
            transform: translateY(-2px);
            box-shadow: 0 4px 8px rgba(0, 0, 0, 0.3);
        }

        .action-btn.help {
            background: var(--accent-blue);
            color: var(--pure-white);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔧 Settings Configuration</h1>
            <p>Manage your bank accounts, credit cards, communication channels, and security settings</p>
        </div>

        <div class="main-content">
            <div class="left-panel">
                <!-- Category Tabs -->
                <div class="category-tabs" id="categoryTabs">
                    {{range $index, $category := .}}
                    <div class="category-tab {{if eq $index 0}}active{{end}}" onclick="switchCategory('{{$category.Name}}')">
                        {{if eq $category.Name "bankaccounts"}}🏦{{else if eq $category.Name "creditcards"}}💳{{else if eq $category.Name "communicationchannels"}}📱{{else if eq $category.Name "emailintegration"}}📧{{else if eq $category.Name "security"}}🔒{{else if eq $category.Name "monitoring"}}📊{{end}} {{$category.Label}}
                    </div>
                    {{end}}
                </div>
            </div>

            <div class="right-panel">
                <!-- Category Contents -->
                {{range $index, $category := .}}
                <div class="category-content {{if eq $index 0}}active{{end}}" id="{{$category.Name}}Content">
                    <div class="form-section">
                        <h3>{{$category.Label}}</h3>
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
                    </div>
                </div>
                {{end}}
            </div>
        </div>

        <!-- Floating Upload Modal -->
        <div id="uploadModal" class="floating-modal" style="display: none;">
            <div class="modal-header">
                <div class="modal-title">📄 Document Upload</div>
                <div class="modal-controls">
                    <button class="modal-btn" onclick="closeModal('uploadModal')">×</button>
                </div>
            </div>
            <div class="modal-content">
                <div style="text-align: center; padding: 40px;">
                    <div style="font-size: 48px; margin-bottom: 16px;">📁</div>
                    <div style="color: #e0e0e0; margin-bottom: 8px;">Drag & drop documents here</div>
                    <div style="color: #808080; font-size: 12px;">or click to browse files</div>
                </div>
            </div>
        </div>

        <!-- Floating Chat Modal -->
        <div id="chatModal" class="floating-modal" style="display: none;">
            <div class="modal-header">
                <div class="modal-title">🤖 Settings Assistant</div>
                <div class="modal-controls">
                    <button class="modal-btn" onclick="closeModal('chatModal')">×</button>
                </div>
            </div>
            <div class="modal-content">
                <div style="padding: 20px;">
                    <div style="background: #404040; padding: 10px; border-radius: 6px; margin-bottom: 16px;">
                        Hello! I'm here to help you configure your settings. I can assist with:
                        <br>• Bank account setup and Open Banking
                        <br>• Credit card configuration
                        <br>• Communication channel setup
                        <br>• Security settings
                        <br>• Email integration
                        <br><br>What would you like help with?
                    </div>
                    <div style="display: flex; gap: 8px;">
                        <input type="text" style="flex: 1; padding: 8px 12px; border: 1px solid #404040; border-radius: 4px; background: #333333; color: #e0e0e0; font-size: 13px;" placeholder="Ask me anything about settings...">
                        <button style="padding: 8px 16px; background: #404040; border: none; border-radius: 4px; color: #e0e0e0; cursor: pointer; font-size: 13px;">Send</button>
                    </div>
                </div>
            </div>
        </div>

        <!-- Floating Action Buttons -->
        <div class="floating-actions">
            <button onclick="toggleModal('uploadModal')" class="action-btn">
                📄 Upload Documents
            </button>
            <button onclick="toggleModal('chatModal')" class="action-btn help">
                🤖 Get Help
            </button>
        </div>
    </div>

    <script>
        // Category switching
        function switchCategory(categoryName) {
            // Hide all category contents
            document.querySelectorAll('.category-content').forEach(content => {
                content.classList.remove('active');
            });
            
            // Remove active from all category tabs
            document.querySelectorAll('.category-tab').forEach(tab => {
                tab.classList.remove('active');
            });
            
            // Show selected category content
            document.getElementById(categoryName + 'Content').classList.add('active');
            
            // Add active to clicked tab
            event.target.classList.add('active');
        }

        // Modal functions
        function toggleModal(modalId) {
            const modal = document.getElementById(modalId);
            if (modal.style.display === 'none' || modal.style.display === '') {
                modal.style.display = 'block';
            } else {
                modal.style.display = 'none';
            }
        }

        function closeModal(modalId) {
            document.getElementById(modalId).style.display = 'none';
        }

        // Initialize
        window.addEventListener('load', function() {
            console.log('Settings form loaded successfully!');
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


