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
	LabelDE                string
	Type                   string
	Required               bool
	HelpText               string
	HelpTextDE             string
	Pattern                string
	ErrorMessage           string
	ErrorMessageDE         string
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
		category = "communicationchannels"
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

		// Parse rdfs:label (English)
		if strings.Contains(propPart, "rdfs:label") && !strings.Contains(propPart, "@de") {
			labelMatch := regexp.MustCompile(`rdfs:label\s+"([^"]+)"`)
			if matches := labelMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.Label = matches[1]
			}
		}

		// Parse rdfs:label (German)
		if strings.Contains(propPart, "rdfs:label") && strings.Contains(propPart, "@de") {
			labelMatch := regexp.MustCompile(`rdfs:label\s+"([^"]+)"@de`)
			if matches := labelMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.LabelDE = matches[1]
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

		// Parse help text (English)
		if strings.Contains(propPart, "autoins:formHelpText") && !strings.Contains(propPart, "@de") {
			helpMatch := regexp.MustCompile(`autoins:formHelpText\s+"([^"]+)"`)
			if matches := helpMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.HelpText = matches[1]
			}
		}

		// Parse help text (German)
		if strings.Contains(propPart, "autoins:formHelpText") && strings.Contains(propPart, "@de") {
			helpMatch := regexp.MustCompile(`autoins:formHelpText\s+"([^"]+)"@de`)
			if matches := helpMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				field.HelpTextDE = matches[1]
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
		"security":              4,
		"monitoring":            5,
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
            position: relative;
            border-bottom: 2px solid var(--gray-600);
            box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3);
        }

        .header-content {
            display: flex;
            justify-content: space-between;
            align-items: center;
            max-width: 1200px;
            margin: 0 auto;
        }

        .header-left {
            text-align: left;
        }

        .header-right {
            text-align: right;
        }

        .language-toggle {
            display: flex;
            gap: 5px;
            background: var(--gray-700);
            padding: 5px;
            border-radius: 6px;
            border: 1px solid var(--gray-600);
        }

        .lang-btn {
            background: var(--gray-800);
            color: var(--gray-300);
            border: 1px solid var(--gray-600);
            padding: 8px 12px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 12px;
            font-weight: 600;
            transition: all 0.3s;
        }

        .lang-btn.active {
            background: var(--pure-white);
            color: var(--pure-black);
            border-color: var(--pure-white);
        }

        .lang-btn:hover:not(.active) {
            background: var(--gray-600);
            color: var(--pure-white);
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
            color: #e0e0e0;
            font-weight: 700;
        }

        .security-badge {
            display: inline-block;
            background: #505050;
            color: #e0e0e0;
            padding: 2px 8px;
            border-radius: 12px;
            font-size: 0.7rem;
            font-weight: 600;
            margin-left: 8px;
        }

        .pci-badge {
            background: #404040;
        }

        .encryption-badge {
            background: #606060;
            color: #e0e0e0;
        }

        /* Communication channel tabs */
        .comm-tabs {
            display: flex;
            gap: 5px;
            margin-bottom: 20px;
            flex-wrap: wrap;
            background: #333333;
            padding: 15px;
            border-radius: 8px;
            border: 1px solid #404040;
        }

        .comm-tab {
            padding: 8px 12px;
            background: #404040;
            border: 1px solid #505050;
            border-radius: 6px;
            cursor: pointer;
            font-size: 12px;
            transition: all 0.3s;
            min-width: 80px;
            text-align: center;
            font-weight: 500;
            color: #808080;
        }

        .comm-tab.active {
            background: #505050;
            color: #e0e0e0;
            border-color: #e0e0e0;
            box-shadow: 0 0 10px rgba(224, 224, 224, 0.3);
        }

        .comm-tab:hover:not(.active) {
            background: #505050;
            color: #e0e0e0;
        }

        .comm-content {
            display: none;
        }

        .comm-content.active {
            display: block;
        }

        /* Voice tabs */
        .voice-tabs {
            display: flex;
            gap: 5px;
            margin-bottom: 20px;
            flex-wrap: wrap;
            background: #2a2a2a;
            padding: 12px;
            border-radius: 6px;
            border: 1px solid #404040;
        }

        .voice-tab {
            padding: 6px 10px;
            background: #404040;
            border: 1px solid #505050;
            border-radius: 4px;
            cursor: pointer;
            font-size: 11px;
            transition: all 0.3s;
            min-width: 60px;
            text-align: center;
            font-weight: 500;
            color: #808080;
        }

        .voice-tab.active {
            background: #505050;
            color: #e0e0e0;
            border-color: #e0e0e0;
            box-shadow: 0 0 8px rgba(224, 224, 224, 0.3);
        }

        .voice-tab:hover:not(.active) {
            background: #505050;
            color: #e0e0e0;
        }

        .voice-content {
            display: none;
        }

        .voice-content.active {
            display: block;
        }

        /* Instance tabs for multiple bank accounts and credit cards */
        .instance-controls {
            display: flex;
            gap: 10px;
            margin-bottom: 20px;
            flex-wrap: wrap;
        }

        .instance-controls .action-btn {
            padding: 8px 16px;
            font-size: 12px;
            border-radius: 6px;
            border: 1px solid #404040;
            background: #333333;
            color: #e0e0e0;
            cursor: pointer;
            transition: all 0.3s;
            font-weight: 500;
        }

        .instance-controls .action-btn:hover {
            background: #404040;
            border-color: #00ff00;
            box-shadow: 0 0 8px rgba(0, 255, 0, 0.3);
        }

        .instance-controls .action-btn.camera {
            background: #404040;
            border-color: #505050;
        }

        .instance-controls .action-btn.camera:hover {
            background: #505050;
            border-color: #e0e0e0;
        }

        .instance-controls .action-btn.secondary {
            background: #404040;
            border-color: #505050;
        }

        .instance-controls .action-btn.secondary:hover {
            background: #505050;
            border-color: #e0e0e0;
        }

        .instance-tabs {
            display: flex;
            gap: 5px;
            margin-bottom: 20px;
            flex-wrap: wrap;
            background: #2a2a2a;
            padding: 12px;
            border-radius: 8px;
            border: 1px solid #404040;
        }

        .instance-tab {
            padding: 8px 12px;
            background: #404040;
            border: 1px solid #505050;
            border-radius: 6px;
            cursor: pointer;
            font-size: 12px;
            transition: all 0.3s;
            min-width: 100px;
            text-align: center;
            font-weight: 500;
            color: #808080;
            display: flex;
            align-items: center;
            justify-content: space-between;
            gap: 8px;
        }

        .instance-tab.active {
            background: #505050;
            color: #e0e0e0;
            border-color: #e0e0e0;
            box-shadow: 0 0 10px rgba(224, 224, 224, 0.3);
        }

        .instance-tab:hover:not(.active) {
            background: #505050;
            color: #e0e0e0;
        }

        .delete-btn {
            background: none;
            border: none;
            color: #808080;
            cursor: pointer;
            font-size: 14px;
            padding: 2px;
            border-radius: 3px;
            transition: all 0.3s;
            opacity: 0.7;
        }

        .delete-btn:hover {
            color: #e0e0e0;
            background: rgba(224, 224, 224, 0.1);
            opacity: 1;
        }

        .instance-tab:hover .delete-btn {
            opacity: 1;
        }

        .instance-contents {
            position: relative;
        }

        .instance-content {
            display: none;
        }

        .instance-content.active {
            display: block;
        }

        /* Digital wallet integration */
        .wallet-options {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
            margin-top: 10px;
        }

        .wallet-btn {
            display: flex;
            align-items: center;
            gap: 8px;
            padding: 12px 16px;
            background: #333333;
            border: 1px solid #404040;
            border-radius: 8px;
            color: #e0e0e0;
            cursor: pointer;
            transition: all 0.3s;
            font-size: 13px;
            font-weight: 500;
        }

        .wallet-btn:hover {
            background: #404040;
            border-color: #e0e0e0;
            box-shadow: 0 0 10px rgba(224, 224, 224, 0.3);
        }

        .wallet-icon {
            font-size: 16px;
        }

        /* Merchant Checkout Style Credit Card */
        .checkout-container {
            display: grid;
            grid-template-columns: 350px 1fr;
            gap: 30px;
            margin-bottom: 30px;
        }

        .card-visual {
            background: linear-gradient(135deg, #1a1a1a 0%, #2a2a2a 100%);
            border: 2px solid #404040;
            border-radius: 12px;
            padding: 20px;
            height: fit-content;
        }

        .card-front {
            background: linear-gradient(135deg, #2a2a2a 0%, #404040 100%);
            border: 1px solid #505050;
            border-radius: 8px;
            padding: 20px;
            color: #e0e0e0;
            font-family: 'Courier New', monospace;
            position: relative;
            min-height: 200px;
        }

        .card-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 20px;
        }

        .card-logo {
            font-size: 12px;
            font-weight: bold;
            color: #808080;
            letter-spacing: 1px;
        }

        .card-chip {
            width: 40px;
            height: 30px;
            background: linear-gradient(135deg, #d4af37 0%, #b8860b 100%);
            border-radius: 4px;
            border: 1px solid #a0522d;
        }

        .card-number-display {
            font-size: 18px;
            font-weight: bold;
            letter-spacing: 2px;
            margin-bottom: 20px;
            color: #e0e0e0;
        }

        .card-details {
            display: flex;
            justify-content: space-between;
            align-items: flex-end;
        }

        .cardholder-name {
            font-size: 12px;
            color: #808080;
            text-transform: uppercase;
            letter-spacing: 1px;
        }

        .card-expiry {
            font-size: 12px;
            color: #808080;
            letter-spacing: 1px;
        }

        .checkout-form {
            background: #1a1a1a;
            border: 1px solid #404040;
            border-radius: 8px;
            padding: 25px;
        }

        .form-section-title {
            color: #e0e0e0;
            font-size: 16px;
            font-weight: 700;
            margin-bottom: 20px;
            text-transform: uppercase;
            letter-spacing: 0.5px;
        }

        .form-row {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 20px;
            margin-bottom: 20px;
        }

        .card-number {
            font-family: 'Courier New', monospace;
            font-size: 16px;
            letter-spacing: 1px;
        }

        .card-icons {
            display: flex;
            gap: 10px;
            margin-top: 8px;
        }

        .card-icon {
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 10px;
            font-weight: bold;
            color: #808080;
            background: #2a2a2a;
            border: 1px solid #404040;
        }

        .expiry-date {
            font-family: 'Courier New', monospace;
            font-size: 14px;
            letter-spacing: 1px;
        }

        .cvv {
            font-family: 'Courier New', monospace;
            font-size: 14px;
            letter-spacing: 1px;
        }

        .cvv-hint {
            font-size: 10px;
            color: #808080;
            margin-top: 4px;
        }

        .cardholder-name {
            font-family: 'Courier New', monospace;
            font-size: 14px;
            text-transform: uppercase;
        }

        .billing-address {
            font-family: 'Courier New', monospace;
            font-size: 12px;
            line-height: 1.4;
        }

        .security-badges {
            display: flex;
            gap: 10px;
            margin-top: 20px;
            flex-wrap: wrap;
        }

        .security-badge {
            padding: 6px 12px;
            border-radius: 6px;
            font-size: 10px;
            font-weight: 600;
            color: #e0e0e0;
            background: #2a2a2a;
            border: 1px solid #404040;
        }

        .ssl-badge {
            background: #1a3a1a;
            border-color: #2a5a2a;
            color: #4a7a4a;
        }

        .wallet-section {
            background: #1a1a1a;
            border: 1px solid #404040;
            border-radius: 8px;
            padding: 25px;
            margin-top: 20px;
        }

        /* Camera scanning modal */
        .camera-modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0, 0, 0, 0.9);
            z-index: 10000;
            justify-content: center;
            align-items: center;
        }

        .camera-content {
            background: var(--card-bg);
            border: 1px solid var(--border-color);
            border-radius: 12px;
            padding: 30px;
            max-width: 500px;
            width: 90%;
            text-align: center;
        }

        .camera-preview {
            width: 100%;
            height: 300px;
            background: #1a1a1a;
            border: 2px dashed #404040;
            border-radius: 8px;
            margin: 20px 0;
            display: flex;
            align-items: center;
            justify-content: center;
            color: #808080;
            font-size: 14px;
        }

        .camera-controls {
            display: flex;
            gap: 10px;
            justify-content: center;
            margin-top: 20px;
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

        /* Upload Modal Styles */
        .upload-instructions {
            margin-bottom: 20px;
        }

        .upload-instructions h4 {
            color: #e0e0e0;
            font-size: 14px;
            margin-bottom: 12px;
            font-weight: 600;
        }

        .document-types {
            display: flex;
            flex-direction: column;
            gap: 12px;
        }

        .doc-type {
            display: flex;
            align-items: center;
            gap: 12px;
            padding: 12px;
            background: #333333;
            border: 1px solid #404040;
            border-radius: 6px;
        }

        .doc-icon {
            font-size: 24px;
            width: 40px;
            text-align: center;
        }

        .doc-info strong {
            color: #e0e0e0;
            font-size: 13px;
            font-weight: 600;
            display: block;
            margin-bottom: 4px;
        }

        .doc-info p {
            color: #808080;
            font-size: 11px;
            margin: 0;
            line-height: 1.3;
        }

        .upload-zone {
            border: 2px dashed #404040;
            border-radius: 8px;
            margin-bottom: 20px;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .upload-zone:hover {
            border-color: #505050;
            background: #2a2a2a;
        }

        .upload-results {
            background: #2a2a2a;
            border: 1px solid #404040;
            border-radius: 6px;
            padding: 16px;
        }

        .upload-results h4 {
            color: #e0e0e0;
            font-size: 14px;
            margin-bottom: 12px;
            font-weight: 600;
        }

        /* Chat Modal Styles */
        .chat-container {
            display: flex;
            flex-direction: column;
            height: 100%;
        }

        .chat-messages {
            flex: 1;
            overflow-y: auto;
            padding: 16px;
            display: flex;
            flex-direction: column;
            gap: 16px;
        }

        .message {
            display: flex;
            flex-direction: column;
        }

        .message.assistant {
            align-items: flex-start;
        }

        .message.user {
            align-items: flex-end;
        }

        .message-content {
            background: #404040;
            padding: 12px 16px;
            border-radius: 8px;
            max-width: 85%;
            color: #e0e0e0;
            font-size: 13px;
            line-height: 1.4;
        }

        .message.user .message-content {
            background: #505050;
            color: #e0e0e0;
        }

        .message-content strong {
            color: #e0e0e0;
            font-weight: 600;
        }

        .message-content ul {
            margin: 8px 0;
            padding-left: 20px;
        }

        .message-content li {
            margin-bottom: 4px;
        }

        .quick-questions {
            display: flex;
            flex-direction: column;
            gap: 8px;
            margin-top: 12px;
        }

        .quick-btn {
            background: #333333;
            border: 1px solid #404040;
            color: #e0e0e0;
            padding: 8px 12px;
            border-radius: 6px;
            font-size: 11px;
            cursor: pointer;
            text-align: left;
            transition: all 0.3s ease;
        }

        .quick-btn:hover {
            background: #404040;
            border-color: #505050;
        }

        .chat-input-container {
            display: flex;
            gap: 8px;
            padding: 16px;
            border-top: 1px solid #404040;
            background: #2a2a2a;
        }

        .chat-input-container input {
            flex: 1;
            padding: 8px 12px;
            border: 1px solid #404040;
            border-radius: 4px;
            background: #333333;
            color: #e0e0e0;
            font-size: 13px;
        }

        .chat-input-container input:focus {
            outline: none;
            border-color: #505050;
        }

        .chat-input-container button {
            padding: 8px 16px;
            background: #404040;
            border: none;
            border-radius: 4px;
            color: #e0e0e0;
            cursor: pointer;
            font-size: 13px;
            transition: all 0.3s ease;
        }

        .chat-input-container button:hover {
            background: #505050;
        }

        /* Modal positioning */
        #uploadModal {
            top: 50px;
            right: 50px;
            width: 450px;
            height: 600px;
        }

        #chatModal {
            bottom: 50px;
            right: 50px;
            width: 400px;
            height: 500px;
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
            <div class="header-content">
                <div class="header-left">
                    <h1 data-en="Settings Configuration" data-de="Einstellungen Konfiguration">Settings Configuration</h1>
                    <p data-en="Manage your bank accounts, credit cards, communication channels, and security settings" data-de="Verwalten Sie Ihre Bankkonten, Kreditkarten, Kommunikationskanäle und Sicherheitseinstellungen">Manage your bank accounts, credit cards, communication channels, and security settings</p>
                </div>
                <div class="header-right">
                    <div class="language-toggle">
                        <button type="button" class="lang-btn active" onclick="switchLanguage('en')" data-lang="en">EN</button>
                        <button type="button" class="lang-btn" onclick="switchLanguage('de')" data-lang="de">DE</button>
                    </div>
                </div>
            </div>
        </div>

        <div class="main-content">
            <div class="left-panel">
                <!-- Category Tabs -->
                <div class="category-tabs" id="categoryTabs">
                    {{range $index, $category := .}}
                    <div class="category-tab {{if eq $index 0}}active{{end}}" onclick="switchCategory('{{$category.Name}}')">
                        {{if eq $category.Name "bankaccounts"}}Bank{{else if eq $category.Name "creditcards"}}Card{{else if eq $category.Name "communicationchannels"}}Comm{{else if eq $category.Name "security"}}Security{{else if eq $category.Name "monitoring"}}Monitor{{end}} - {{$category.Label}}
                    </div>
                    {{end}}
                </div>
            </div>

            <div class="right-panel">
                <!-- Category Contents -->
                {{range $index, $category := .}}
                <div class="category-content {{if eq $index 0}}active{{end}}" id="{{$category.Name}}Content">
                    {{if eq $category.Name "bankaccounts"}}
                    <!-- Bank Accounts with Multiple Instance Tabs -->
                    <div class="form-section">
                        <h3>{{$category.Label}}</h3>
                        <div class="instance-controls">
                            <button type="button" class="action-btn" onclick="addBankAccount()">+ Add Bank Account</button>
                        </div>
                        
                        <div class="instance-tabs" id="bankAccountTabs">
                            <div class="instance-tab active" onclick="switchBankTab(0)">
                                <span>Bank Account 1</span>
                                <button type="button" class="delete-btn" onclick="deleteBankAccount(0, event)" title="Delete Bank Account">X</button>
                            </div>
                        </div>
                        
                        <div class="instance-contents" id="bankAccountContents">
                            <div class="instance-content active" id="bankAccount0">
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
                                               name="bank0_{{$field.Property}}" 
                                               placeholder="{{$field.Label}}"
                                               {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                               {{if $field.Required}}required{{end}}
                                               {{if $field.ReadOnly}}readonly{{end}}>
                                        
                                        {{else if eq $field.Type "number"}}
                                        <input type="number" 
                                               class="form-input" 
                                               name="bank0_{{$field.Property}}"
                                               {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                               {{if $field.Max}}max="{{$field.Max}}"{{end}}
                                               {{if $field.Required}}required{{end}}
                                               {{if $field.ReadOnly}}readonly{{end}}>
                                        
                                        {{else if eq $field.Type "date"}}
                                        <input type="date" 
                                               class="form-input" 
                                               name="bank0_{{$field.Property}}"
                                               {{if $field.Required}}required{{end}}
                                               {{if $field.ReadOnly}}readonly{{end}}>
                                        
                                        {{else if eq $field.Type "email"}}
                                        <input type="email" 
                                               class="form-input" 
                                               name="bank0_{{$field.Property}}"
                                               {{if $field.Required}}required{{end}}
                                               {{if $field.ReadOnly}}readonly{{end}}>
                                        
                                        {{else if eq $field.Type "checkbox"}}
                                        <input type="checkbox" 
                                               class="form-input" 
                                               name="bank0_{{$field.Property}}"
                                               {{if $field.Required}}required{{end}}
                                               {{if $field.ReadOnly}}readonly{{end}}>
                                        
                                        {{else if eq $field.Type "select"}}
                                        <select class="form-select" 
                                                name="bank0_{{$field.Property}}" 
                                                {{if $field.Required}}required{{end}}
                                                {{if $field.ReadOnly}}readonly{{end}}>
                                            <option value="">Select {{$field.Label}}</option>
                                            {{range $option := $field.Options}}
                                            <option value="{{$option.Value}}">{{$option.Label}}</option>
                                            {{end}}
                                        </select>
                                        
                                        {{else if eq $field.Type "textarea"}}
                                        <textarea class="form-textarea" 
                                                  name="bank0_{{$field.Property}}" 
                                                  rows="3" 
                                                  placeholder="{{$field.Label}}"
                                                  {{if $field.Required}}required{{end}}
                                                  {{if $field.ReadOnly}}readonly{{end}}></textarea>
                                        
                                        {{else}}
                                        <input type="text" 
                                               class="form-input" 
                                               name="bank0_{{$field.Property}}" 
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
                    </div>
                    
                    {{else if eq $category.Name "creditcards"}}
                    <!-- Credit Cards with Merchant Checkout Style -->
                    <div class="form-section">
                        <h3>{{$category.Label}}</h3>
                        <div class="instance-controls">
                            <button type="button" class="action-btn" onclick="addCreditCard()">+ Add Credit Card</button>
                            <button type="button" class="action-btn camera" onclick="scanCard()">Camera Scan Card</button>
                        </div>
                        
                        <div class="instance-tabs" id="creditCardTabs">
                            <div class="instance-tab active" onclick="switchCardTab(0)">
                                <span>Credit Card 1</span>
                                <button type="button" class="delete-btn" onclick="deleteCreditCard(0, event)" title="Delete Credit Card">X</button>
                            </div>
                        </div>
                        
                        <div class="instance-contents" id="creditCardContents">
                            <div class="instance-content active" id="creditCard0">
                                <div class="checkout-container">
                                    <div class="card-visual">
                                        <div class="card-front">
                                            <div class="card-header">
                                                <div class="card-logo">CREDIT CARD</div>
                                                <div class="card-chip"></div>
                                            </div>
                                            <div class="card-number-display" id="cardNumberDisplay0">•••• •••• •••• ••••</div>
                                            <div class="card-details">
                                                <div class="cardholder-name" id="cardholderNameDisplay0">CARDHOLDER NAME</div>
                                                <div class="card-expiry" id="cardExpiryDisplay0">MM/YY</div>
                                            </div>
                                        </div>
                                    </div>
                                    
                                    <div class="checkout-form">
                                        <div class="form-section-title">Card Information</div>
                                        
                                        <!-- Card Number -->
                                        <div class="form-field full-width">
                                            <label class="form-label">Card Number <span class="required">*</span></label>
                                            <input type="text" 
                                                   class="form-input card-number" 
                                                   name="card0_cardNumber" 
                                                   placeholder="1234 5678 9012 3456"
                                                   maxlength="19"
                                                   onkeyup="formatCardNumber(this, 0)"
                                                   required>
                                            <div class="card-icons">
                                                <span class="card-icon visa">VISA</span>
                                                <span class="card-icon mastercard">MC</span>
                                                <span class="card-icon amex">AMEX</span>
                                            </div>
                                        </div>
                                        
                                        <div class="form-row">
                                            <!-- Expiry Date -->
                                            <div class="form-field">
                                                <label class="form-label">Expiry Date <span class="required">*</span></label>
                                                <input type="text" 
                                                       class="form-input expiry-date" 
                                                       name="card0_expiryDate" 
                                                       placeholder="MM/YY"
                                                       maxlength="5"
                                                       onkeyup="formatExpiryDate(this, 0)"
                                                       required>
                                            </div>
                                            
                                            <!-- CVV -->
                                            <div class="form-field">
                                                <label class="form-label">CVV <span class="required">*</span></label>
                                                <input type="text" 
                                                       class="form-input cvv" 
                                                       name="card0_securityCode" 
                                                       placeholder="123"
                                                       maxlength="4"
                                                       onkeyup="updateCVVDisplay(this, 0)"
                                                       required>
                                                <div class="cvv-hint">3-4 digits on back of card</div>
                                            </div>
                                        </div>
                                        
                                        <!-- Cardholder Name -->
                                        <div class="form-field full-width">
                                            <label class="form-label">Cardholder Name <span class="required">*</span></label>
                                            <input type="text" 
                                                   class="form-input cardholder-name" 
                                                   name="card0_cardholderName" 
                                                   placeholder="JOHN DOE"
                                                   onkeyup="updateCardholderDisplay(this, 0)"
                                                   required>
                                        </div>
                                        
                                        <!-- Cardholder Address -->
                                        <div class="form-field full-width">
                                            <label class="form-label">Billing Address <span class="required">*</span></label>
                                            <textarea class="form-textarea billing-address" 
                                                      name="card0_cardholderAddress" 
                                                      rows="3" 
                                                      placeholder="123 Main Street, City, State, ZIP"
                                                      required></textarea>
                                        </div>
                                        
                                        <!-- Card Provider & Type -->
                                        <div class="form-row">
                                            <div class="form-field">
                                                <label class="form-label">Card Provider</label>
                                                <select class="form-select" name="card0_cardProvider">
                                                    <option value="">Select Provider</option>
                                                    <option value="Visa">Visa</option>
                                                    <option value="Mastercard">Mastercard</option>
                                                    <option value="American Express">American Express</option>
                                                    <option value="Discover">Discover</option>
                                                    <option value="JCB">JCB</option>
                                                    <option value="UnionPay">UnionPay</option>
                                                </select>
                                            </div>
                                            
                                            <div class="form-field">
                                                <label class="form-label">Card Type</label>
                                                <select class="form-select" name="card0_cardType">
                                                    <option value="">Select Type</option>
                                                    <option value="Credit">Credit</option>
                                                    <option value="Debit">Debit</option>
                                                    <option value="Prepaid">Prepaid</option>
                                                    <option value="Business">Business</option>
                                                </select>
                                            </div>
                                        </div>
                                        
                                        <!-- Security Badges -->
                                        <div class="security-badges">
                                            <div class="security-badge pci-badge">PCI Compliant</div>
                                            <div class="security-badge encryption-badge">256-bit Encryption</div>
                                            <div class="security-badge ssl-badge">SSL Secure</div>
                                        </div>
                                    </div>
                                </div>
                                
                                <!-- Digital Wallet Integration -->
                                <div class="wallet-section">
                                    <div class="form-section-title">Digital Wallet Integration</div>
                                    <div class="wallet-options">
                                        <button type="button" class="wallet-btn" onclick="connectGoogleWallet()">
                                            Google Wallet
                                        </button>
                                        <button type="button" class="wallet-btn" onclick="connectAppleWallet()">
                                            Apple Wallet
                                        </button>
                                        <button type="button" class="wallet-btn" onclick="connectSamsungPay()">
                                            Samsung Pay
                                        </button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    {{else if eq $category.Name "communicationchannels"}}
                    <!-- Communication Channel Tabs -->
                    <div class="comm-tabs">
                        <div class="comm-tab active" onclick="switchCommTab('email')">Email</div>
                        <div class="comm-tab" onclick="switchCommTab('sms')">SMS</div>
                        <div class="comm-tab" onclick="switchCommTab('voice')">Voice</div>
                        <div class="comm-tab" onclick="switchCommTab('secure')">Secure Messenger</div>
                    </div>
                    
                    <!-- Email Content -->
                    <div class="comm-content active" id="emailContent">
                        <div class="form-section">
                            <h3>Email Communication</h3>
                            
                            <!-- Email Provider Selection -->
                            <div class="form-section-title" data-en="Email Provider Configuration" data-de="E-Mail Anbieter Konfiguration">Email Provider Configuration</div>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                {{if or (eq $field.Property "emailProvider") (eq $field.Property "emailAddress")}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label class="form-label">
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                        {{if $field.EncryptionRequired}}<span class="security-badge encryption-badge">Encrypted</span>{{end}}
                                        {{if $field.PCICompliant}}<span class="security-badge pci-badge">PCI</span>{{end}}
                                        {{if $field.Requires2FA}}<span class="security-badge">2FA</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select class="form-select" 
                                            name="email_{{$field.Property}}" 
                                            {{if $field.Required}}required{{end}}
                                            {{if $field.ReadOnly}}readonly{{end}}
                                            onchange="handleEmailProviderChange(this.value)">
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else}}
                                    <input type="text" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                                {{end}}
                            </div>
                            
                            <!-- Authentication Section -->
                            <div class="form-section-title" data-en="Authentication & Security" data-de="Authentifizierung & Sicherheit">Authentication & Security</div>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                {{if or (eq $field.Property "oauth2Enabled") (eq $field.Property "emailPassword") (eq $field.Property "appPassword")}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label class="form-label">
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                        {{if $field.EncryptionRequired}}<span class="security-badge encryption-badge">Encrypted</span>{{end}}
                                        {{if $field.PCICompliant}}<span class="security-badge pci-badge">PCI</span>{{end}}
                                        {{if $field.Requires2FA}}<span class="security-badge">2FA</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="password" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}
                                           onchange="handleOAuth2Toggle(this.checked)">
                                    
                                    {{else}}
                                    <input type="text" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                                {{end}}
                            </div>
                            
                            <!-- OAuth2 Configuration (Conditional) -->
                            <div class="form-section-title" id="oauth2Section" style="display: none;" data-en="OAuth2 Configuration" data-de="OAuth2 Konfiguration">OAuth2 Configuration</div>
                            <div class="form-grid" id="oauth2Fields" style="display: none;">
                                {{range $field := $category.Fields}}
                                {{if or (eq $field.Property "oauth2ClientId") (eq $field.Property "oauth2ClientSecret") (eq $field.Property "refreshToken") (eq $field.Property "accessToken")}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label class="form-label">
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                        {{if $field.EncryptionRequired}}<span class="security-badge encryption-badge">Encrypted</span>{{end}}
                                        {{if $field.PCICompliant}}<span class="security-badge pci-badge">PCI</span>{{end}}
                                        {{if $field.Requires2FA}}<span class="security-badge">2FA</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="password" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    
                                    {{else}}
                                    <input type="text" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                                {{end}}
                            </div>
                            
                            <!-- Sync Configuration -->
                            <div class="form-section-title" data-en="Sync Configuration" data-de="Sync Konfiguration">Sync Configuration</div>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                {{if or (eq $field.Property "autoSyncEnabled") (eq $field.Property "syncFrequency")}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label class="form-label">
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                        {{if $field.EncryptionRequired}}<span class="security-badge encryption-badge">Encrypted</span>{{end}}
                                        {{if $field.PCICompliant}}<span class="security-badge pci-badge">PCI</span>{{end}}
                                        {{if $field.Requires2FA}}<span class="security-badge">2FA</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}
                                           onchange="handleAutoSyncToggle(this.checked)">
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select class="form-select" 
                                            name="email_{{$field.Property}}" 
                                            {{if $field.Required}}required{{end}}
                                            {{if $field.ReadOnly}}readonly{{end}}>
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else}}
                                    <input type="text" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                                {{end}}
                            </div>
                            
                            <!-- Custom Server Configuration (Conditional) -->
                            <div class="form-section-title" id="customServerSection" style="display: none;" data-en="Custom Server Configuration" data-de="Benutzerdefinierte Server Konfiguration">Custom Server Configuration</div>
                            <div class="form-grid" id="customServerFields" style="display: none;">
                                {{range $field := $category.Fields}}
                                {{if or (eq $field.Property "imapServer") (eq $field.Property "imapPort") (eq $field.Property "imapUseSSL") (eq $field.Property "smtpServer") (eq $field.Property "smtpPort") (eq $field.Property "smtpUseTLS")}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label class="form-label">
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                        {{if $field.EncryptionRequired}}<span class="security-badge encryption-badge">Encrypted</span>{{end}}
                                        {{if $field.PCICompliant}}<span class="security-badge pci-badge">PCI</span>{{end}}
                                        {{if $field.Requires2FA}}<span class="security-badge">2FA</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "number"}}
                                    <input type="number" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    
                                    {{else}}
                                    <input type="text" 
                                           class="form-input" 
                                           name="email_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                                {{end}}
                            </div>
                            
                            <!-- Security Badges -->
                            <div class="security-badges">
                                <div class="security-badge pci-badge">PCI Compliant Storage</div>
                                <div class="security-badge encryption-badge">256-bit Encryption</div>
                                <div class="security-badge ssl-badge">SSL/TLS Secure</div>
                                <div class="security-badge">OAuth2 Support</div>
                            </div>
                        </div>
                    </div>
                    
                    <!-- SMS Content -->
                    <div class="comm-content" id="smsContent">
                        <div class="form-section">
                            <h3>📱 SMS Communication</h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                {{if or (eq $field.Property "smsProvider") (eq $field.Property "smsApiKey") (eq $field.Property "smsApiSecret") (eq $field.Property "mobileNumber")}}
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
                                           name="sms_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           class="form-input" 
                                           name="sms_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select class="form-select" 
                                            name="sms_{{$field.Property}}" 
                                            {{if $field.Required}}required{{end}}
                                            {{if $field.ReadOnly}}readonly{{end}}>
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else}}
                                    <input type="text" 
                                           class="form-input" 
                                           name="sms_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    <!-- Voice Content -->
                    <div class="comm-content" id="voiceContent">
                        <div class="form-section">
                            <h3>🎤 Voice Communication</h3>
                            
                            <!-- Voice Platform Tabs -->
                            <div class="voice-tabs">
                                <div class="voice-tab active" onclick="switchVoiceTab('alexa')">Alexa</div>
                                <div class="voice-tab" onclick="switchVoiceTab('google')">Google</div>
                                <div class="voice-tab" onclick="switchVoiceTab('phone')">Phone</div>
                                <div class="voice-tab" onclick="switchVoiceTab('siri')">Siri</div>
                                <div class="voice-tab" onclick="switchVoiceTab('cortana')">Cortana</div>
                                <div class="voice-tab" onclick="switchVoiceTab('homekit')">HomeKit</div>
                            </div>
                            
                            <!-- Alexa Voice Content -->
                            <div class="voice-content active" id="alexaContent">
                                <div class="form-grid">
                                    <div class="form-field">
                                        <label class="form-label">Alexa Skill Name <span class="required">*</span></label>
                                        <input type="text" class="form-input" name="voice_alexa_skillName" placeholder="Insurance Assistant" required>
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">Amazon Developer Account</label>
                                        <input type="text" class="form-input" name="voice_alexa_developerAccount" placeholder="your-email@amazon.com">
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">Skill ID</label>
                                        <input type="text" class="form-input" name="voice_alexa_skillId" placeholder="amzn1.ask.skill.xxx">
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">Invocation Name</label>
                                        <input type="text" class="form-input" name="voice_alexa_invocationName" placeholder="insurance assistant">
                                    </div>
                                </div>
                            </div>
                            
                            <!-- Google Voice Content -->
                            <div class="voice-content" id="googleContent">
                                <div class="form-grid">
                                    <div class="form-field">
                                        <label class="form-label">Google Action Name <span class="required">*</span></label>
                                        <input type="text" class="form-input" name="voice_google_actionName" placeholder="Insurance Helper" required>
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">Google Cloud Project ID</label>
                                        <input type="text" class="form-input" name="voice_google_projectId" placeholder="your-project-id">
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">Dialogflow Agent ID</label>
                                        <input type="text" class="form-input" name="voice_google_agentId" placeholder="agent-id">
                                    </div>
                                </div>
                            </div>
                            
                            <!-- Phone Voice Content -->
                            <div class="voice-content" id="phoneContent">
                                <div class="form-grid">
                                    <div class="form-field">
                                        <label class="form-label">Phone Number <span class="required">*</span></label>
                                        <input type="tel" class="form-input" name="voice_phone_number" placeholder="+44 20 7946 0958" required>
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">Voice Provider</label>
                                        <select class="form-select" name="voice_phone_provider">
                                            <option value="">Select Provider</option>
                                            <option value="twilio">Twilio</option>
                                            <option value="vonage">Vonage</option>
                                            <option value="aws-connect">AWS Connect</option>
                                            <option value="google-voice">Google Voice</option>
                                        </select>
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">IVR Menu Options</label>
                                        <textarea class="form-textarea" name="voice_phone_ivrOptions" rows="3" placeholder="1. Get quote&#10;2. Check policy&#10;3. Make claim&#10;4. Speak to agent"></textarea>
                                    </div>
                                </div>
                            </div>
                            
                            <!-- Siri Voice Content -->
                            <div class="voice-content" id="siriContent">
                                <div class="form-grid">
                                    <div class="form-field">
                                        <label class="form-label">Siri Shortcut Name <span class="required">*</span></label>
                                        <input type="text" class="form-input" name="voice_siri_shortcutName" placeholder="Check Insurance" required>
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">Apple Developer Account</label>
                                        <input type="text" class="form-input" name="voice_siri_developerAccount" placeholder="your-email@icloud.com">
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">Intent Types</label>
                                        <select class="form-select" name="voice_siri_intentTypes" multiple>
                                            <option value="quote">Get Quote</option>
                                            <option value="policy">Check Policy</option>
                                            <option value="claim">Make Claim</option>
                                            <option value="renewal">Renewal Reminder</option>
                                        </select>
                                    </div>
                                </div>
                            </div>
                            
                            <!-- Cortana Voice Content -->
                            <div class="voice-content" id="cortanaContent">
                                <div class="form-grid">
                                    <div class="form-field">
                                        <label class="form-label">Cortana Skill Name <span class="required">*</span></label>
                                        <input type="text" class="form-input" name="voice_cortana_skillName" placeholder="Insurance Bot" required>
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">Microsoft Azure Bot Service</label>
                                        <input type="text" class="form-input" name="voice_cortana_botService" placeholder="your-bot-service-url">
                                    </div>
                                </div>
                            </div>
                            
                            <!-- HomeKit Voice Content -->
                            <div class="voice-content" id="homekitContent">
                                <div class="form-grid">
                                    <div class="form-field">
                                        <label class="form-label">HomeKit Accessory Name <span class="required">*</span></label>
                                        <input type="text" class="form-input" name="voice_homekit_accessoryName" placeholder="Insurance Hub" required>
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">HomeKit Bridge</label>
                                        <input type="text" class="form-input" name="voice_homekit_bridge" placeholder="bridge-identifier">
                                    </div>
                                    <div class="form-field">
                                        <label class="form-label">Supported Services</label>
                                        <select class="form-select" name="voice_homekit_services" multiple>
                                            <option value="notifications">Notifications</option>
                                            <option value="alerts">Alerts</option>
                                            <option value="status">Status Updates</option>
                                        </select>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    
                    <!-- Secure Messenger Content -->
                    <div class="comm-content" id="secureContent">
                        <div class="form-section">
                            <h3>🔒 Secure Messenger</h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                {{if or (eq $field.Property "telegramBotToken") (eq $field.Property "telegramChatId") (eq $field.Property "secureChannelType")}}
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
                                           name="secure_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select class="form-select" 
                                            name="secure_{{$field.Property}}" 
                                            {{if $field.Required}}required{{end}}
                                            {{if $field.ReadOnly}}readonly{{end}}>
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else}}
                                    <input type="text" 
                                           class="form-input" 
                                           name="secure_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Required}}required{{end}}
                                           {{if $field.ReadOnly}}readonly{{end}}>
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                                {{end}}
                                
                                <!-- Additional Secure Messenger Fields -->
                                <div class="form-field">
                                    <label class="form-label">Secure Channel Type <span class="required">*</span></label>
                                    <select class="form-select" name="secure_channelType" required>
                                        <option value="">Select Channel</option>
                                        <option value="telegram">Telegram</option>
                                        <option value="signal">Signal</option>
                                        <option value="whatsapp-business">WhatsApp Business</option>
                                        <option value="slack">Slack</option>
                                        <option value="microsoft-teams">Microsoft Teams</option>
                                        <option value="discord">Discord</option>
                                    </select>
                                </div>
                                <div class="form-field">
                                    <label class="form-label">Encryption Level <span class="security-badge encryption-badge">🔐 Encrypted</span></label>
                                    <select class="form-select" name="secure_encryptionLevel">
                                        <option value="">Select Level</option>
                                        <option value="end-to-end">End-to-End</option>
                                        <option value="transport">Transport Layer</option>
                                        <option value="at-rest">At Rest</option>
                                        <option value="military-grade">Military Grade</option>
                                    </select>
                                </div>
                            </div>
                        </div>
                    </div>
                    {{else}}
                    <!-- Standard Category Content -->
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
                    {{end}}
                </div>
                {{end}}
            </div>
        </div>

        <!-- Floating Upload Modal -->
        <div id="uploadModal" class="floating-modal" style="display: none;">
            <div class="modal-header">
                <div class="modal-title" data-en="Document Upload" data-de="Dokumente hochladen">Document Upload</div>
                <div class="modal-controls">
                    <button class="modal-btn" onclick="closeModal('uploadModal')">×</button>
                </div>
            </div>
            <div class="modal-content">
                <div class="upload-instructions">
                    <h4 data-en="Supported Documents" data-de="Unterstützte Dokumente">Supported Documents</h4>
                    <div class="document-types">
                        <div class="doc-type">
                            <div class="doc-icon">🏦</div>
                            <div class="doc-info">
                                <strong data-en="Bank Statements" data-de="Kontoauszüge">Bank Statements</strong>
                                <p data-en="Extract account numbers, sort codes, and bank details" data-de="Kontonummern, Sort Codes und Bankdetails extrahieren">Extract account numbers, sort codes, and bank details</p>
                            </div>
                        </div>
                        <div class="doc-type">
                            <div class="doc-icon">💳</div>
                            <div class="doc-info">
                                <strong data-en="Credit Card Statements" data-de="Kreditkartenabrechnungen">Credit Card Statements</strong>
                                <p data-en="Extract card numbers, expiry dates, and provider info" data-de="Kartennummern, Ablaufdaten und Anbieterinformationen extrahieren">Extract card numbers, expiry dates, and provider info</p>
                            </div>
                        </div>
                        <div class="doc-type">
                            <div class="doc-icon">📧</div>
                            <div class="doc-info">
                                <strong data-en="Email Configuration Screenshots" data-de="E-Mail Konfigurations Screenshots">Email Configuration Screenshots</strong>
                                <p data-en="Extract server settings, ports, and authentication details" data-de="Server-Einstellungen, Ports und Authentifizierungsdetails extrahieren">Extract server settings, ports, and authentication details</p>
                            </div>
                        </div>
                        <div class="doc-type">
                            <div class="doc-icon">🔐</div>
                            <div class="doc-info">
                                <strong data-en="OAuth2 Setup Screenshots" data-de="OAuth2 Setup Screenshots">OAuth2 Setup Screenshots</strong>
                                <p data-en="Extract client IDs, redirect URIs, and OAuth2 settings" data-de="Client IDs, Redirect URIs und OAuth2-Einstellungen extrahieren">Extract client IDs, redirect URIs, and OAuth2 settings</p>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="upload-zone" id="uploadZone">
                    <div style="text-align: center; padding: 40px;">
                        <div style="font-size: 48px; margin-bottom: 16px;">📁</div>
                        <div style="color: #e0e0e0; margin-bottom: 8px;" data-en="Drag & drop documents here" data-de="Dokumente hier hineinziehen">Drag & drop documents here</div>
                        <div style="color: #808080; font-size: 12px;" data-en="or click to browse files" data-de="oder klicken zum Durchsuchen">or click to browse files</div>
                        <input type="file" id="fileInput" multiple accept=".pdf,.jpg,.jpeg,.png,.txt,.doc,.docx" style="display: none;" onchange="handleFileUpload(this.files)">
                    </div>
                </div>
                <div class="upload-results" id="uploadResults" style="display: none;">
                    <h4 data-en="Extracted Information" data-de="Extrahierte Informationen">Extracted Information</h4>
                    <div id="extractedData"></div>
                </div>
            </div>
        </div>

        <!-- Floating Chat Modal -->
        <div id="chatModal" class="floating-modal" style="display: none;">
            <div class="modal-header">
                <div class="modal-title" data-en="Settings Assistant" data-de="Einstellungen Assistent">Settings Assistant</div>
                <div class="modal-controls">
                    <button class="modal-btn" onclick="closeModal('chatModal')">×</button>
                </div>
            </div>
            <div class="modal-content">
                <div class="chat-container">
                    <div class="chat-messages" id="chatMessages">
                        <div class="message assistant">
                            <div class="message-content">
                                <strong data-en="Assistant:" data-de="Assistent:">Assistant:</strong> 
                                <p data-en="Hello! I'm your settings configuration assistant. I can help you with:" data-de="Hallo! Ich bin Ihr Einstellungen-Konfigurationsassistent. Ich kann Ihnen helfen bei:">Hello! I'm your settings configuration assistant. I can help you with:</p>
                                <ul>
                                    <li data-en="Open Banking account setup and configuration" data-de="Open Banking Konto Setup und Konfiguration">Open Banking account setup and configuration</li>
                                    <li data-en="Credit card configuration and security" data-de="Kreditkarten-Konfiguration und Sicherheit">Credit card configuration and security</li>
                                    <li data-en="Email integration (Gmail, Outlook, OAuth2)" data-de="E-Mail-Integration (Gmail, Outlook, OAuth2)">Email integration (Gmail, Outlook, OAuth2)</li>
                                    <li data-en="Security settings and PCI compliance" data-de="Sicherheitseinstellungen und PCI-Compliance">Security settings and PCI compliance</li>
                                </ul>
                                <p data-en="Try asking me questions like:" data-de="Versuchen Sie mir Fragen zu stellen wie:">Try asking me questions like:</p>
                                <div class="quick-questions">
                                    <button class="quick-btn" onclick="askQuestion('How do I add an Open Banking account?')" data-en="How do I add an Open Banking account?" data-de="Wie füge ich ein Open Banking Konto hinzu?">How do I add an Open Banking account?</button>
                                    <button class="quick-btn" onclick="askQuestion('How do I configure OAuth2 for Gmail?')" data-en="How do I configure OAuth2 for Gmail?" data-de="Wie konfiguriere ich OAuth2 für Gmail?">How do I configure OAuth2 for Gmail?</button>
                                    <button class="quick-btn" onclick="askQuestion('How do I add a credit card securely?')" data-en="How do I add a credit card securely?" data-de="Wie füge ich eine Kreditkarte sicher hinzu?">How do I add a credit card securely?</button>
                                    <button class="quick-btn" onclick="askQuestion('What is PCI compliance?')" data-en="What is PCI compliance?" data-de="Was ist PCI-Compliance?">What is PCI compliance?</button>
                                </div>
                            </div>
                        </div>
                    </div>
                    <div class="chat-input-container">
                        <input type="text" id="chatInput" placeholder="Ask me anything about settings..." data-en-placeholder="Ask me anything about settings..." data-de-placeholder="Fragen Sie mich alles über Einstellungen..." onkeypress="handleChatKeyPress(event)">
                        <button onclick="sendChatMessage()" data-en="Send" data-de="Senden">Send</button>
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

    <!-- Camera Scanning Modal -->
    <div class="camera-modal" id="cameraModal">
        <div class="camera-content">
            <h3>Scan Credit Card</h3>
            <p>Position your credit card within the frame below</p>
            
            <div class="camera-preview" id="cameraPreview">
                <div>Camera access required for card scanning</div>
            </div>
            
            <div class="camera-controls">
                <button type="button" class="action-btn" onclick="startCamera()">Start Camera</button>
                <button type="button" class="action-btn" onclick="captureCard()">Capture</button>
                <button type="button" class="action-btn secondary" onclick="closeCameraModal()">Cancel</button>
            </div>
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
        
        // Communication channel tab switching
        function switchCommTab(tabName) {
            // Hide all comm content
            const commContents = document.querySelectorAll('.comm-content');
            commContents.forEach(content => content.classList.remove('active'));
            
            // Remove active class from all comm tabs
            const commTabs = document.querySelectorAll('.comm-tab');
            commTabs.forEach(tab => tab.classList.remove('active'));
            
            // Show selected content and activate tab
            document.getElementById(tabName + 'Content').classList.add('active');
            event.target.classList.add('active');
        }
        
        // Voice platform tab switching
        function switchVoiceTab(tabName) {
            // Hide all voice content
            const voiceContents = document.querySelectorAll('.voice-content');
            voiceContents.forEach(content => content.classList.remove('active'));
            
            // Remove active class from all voice tabs
            const voiceTabs = document.querySelectorAll('.voice-tab');
            voiceTabs.forEach(tab => tab.classList.remove('active'));
            
            // Show selected content and activate tab
            document.getElementById(tabName + 'Content').classList.add('active');
            event.target.classList.add('active');
        }

        // Bank Account Management
        let bankAccountCount = 1;
        
        function addBankAccount() {
            const tabsContainer = document.getElementById('bankAccountTabs');
            const contentsContainer = document.getElementById('bankAccountContents');
            
            // Create new tab
            const newTab = document.createElement('div');
            newTab.className = 'instance-tab';
            newTab.onclick = function() { switchBankTab(bankAccountCount); };
            newTab.innerHTML = '<span>Bank Account ' + (bankAccountCount + 1) + '</span><button type="button" class="delete-btn" onclick="deleteBankAccount(' + bankAccountCount + ', event)" title="Delete Bank Account">X</button>';
            tabsContainer.appendChild(newTab);
            
            // Create new content
            const newContent = document.createElement('div');
            newContent.className = 'instance-content';
            newContent.id = 'bankAccount' + bankAccountCount;
            
            // Clone the form structure from the first bank account
            const firstContent = document.getElementById('bankAccount0');
            newContent.innerHTML = firstContent.innerHTML.replace(/bank0_/g, 'bank' + bankAccountCount + '_');
            
            contentsContainer.appendChild(newContent);
            
            // Switch to new tab
            switchBankTab(bankAccountCount);
            bankAccountCount++;
        }
        
        function deleteBankAccount(index, event) {
            event.stopPropagation(); // Prevent tab switching
            
            const accountName = 'Bank Account ' + (index + 1);
            if (confirm('Are you sure you want to delete ' + accountName + '? This action cannot be undone.')) {
                // Remove the tab
                const tabsContainer = document.getElementById('bankAccountTabs');
                const tabs = tabsContainer.querySelectorAll('.instance-tab');
                if (tabs[index]) {
                    tabsContainer.removeChild(tabs[index]);
                }
                
                // Remove the content
                const contentsContainer = document.getElementById('bankAccountContents');
                const content = document.getElementById('bankAccount' + index);
                if (content) {
                    contentsContainer.removeChild(content);
                }
                
                // Renumber remaining tabs and contents
                renumberBankAccounts();
                
                // Switch to first tab if current was removed
                if (index === 0) {
                    switchBankTab(0);
                } else {
                    switchBankTab(Math.max(0, index - 1));
                }
                
                bankAccountCount--;
            }
        }
        
        function renumberBankAccounts() {
            // Renumber tabs
            const tabs = document.querySelectorAll('#bankAccountTabs .instance-tab');
            tabs.forEach((tab, index) => {
                const span = tab.querySelector('span');
                if (span) {
                    span.textContent = 'Bank Account ' + (index + 1);
                }
                
                // Update onclick for tab switching
                tab.onclick = function() { switchBankTab(index); };
                
                // Update delete button onclick
                const deleteBtn = tab.querySelector('.delete-btn');
                if (deleteBtn) {
                    deleteBtn.onclick = function(event) { deleteBankAccount(index, event); };
                }
            });
            
            // Renumber content IDs and form field names
            const contents = document.querySelectorAll('#bankAccountContents .instance-content');
            contents.forEach((content, index) => {
                content.id = 'bankAccount' + index;
                
                // Update form field names
                const inputs = content.querySelectorAll('input, select, textarea');
                inputs.forEach(input => {
                    const oldName = input.name;
                    if (oldName && oldName.startsWith('bank')) {
                        const parts = oldName.split('_');
                        if (parts.length > 1) {
                            input.name = 'bank' + index + '_' + parts.slice(1).join('_');
                        }
                    }
                });
            });
        }
        
        function switchBankTab(index) {
            // Hide all bank account contents
            const contents = document.querySelectorAll('#bankAccountContents .instance-content');
            contents.forEach(content => content.classList.remove('active'));
            
            // Remove active from all bank account tabs
            const tabs = document.querySelectorAll('#bankAccountTabs .instance-tab');
            tabs.forEach(tab => tab.classList.remove('active'));
            
            // Show selected content and activate tab
            const selectedContent = document.getElementById('bankAccount' + index);
            const selectedTab = document.querySelectorAll('#bankAccountTabs .instance-tab')[index];
            
            if (selectedContent) selectedContent.classList.add('active');
            if (selectedTab) selectedTab.classList.add('active');
        }

        // Credit Card Management
        let creditCardCount = 1;
        
        function addCreditCard() {
            const tabsContainer = document.getElementById('creditCardTabs');
            const contentsContainer = document.getElementById('creditCardContents');
            
            // Create new tab
            const newTab = document.createElement('div');
            newTab.className = 'instance-tab';
            newTab.onclick = function() { switchCardTab(creditCardCount); };
            newTab.innerHTML = '<span>Credit Card ' + (creditCardCount + 1) + '</span><button type="button" class="delete-btn" onclick="deleteCreditCard(' + creditCardCount + ', event)" title="Delete Credit Card">X</button>';
            tabsContainer.appendChild(newTab);
            
            // Create new content
            const newContent = document.createElement('div');
            newContent.className = 'instance-content';
            newContent.id = 'creditCard' + creditCardCount;
            
            // Clone the form structure from the first credit card
            const firstContent = document.getElementById('creditCard0');
            newContent.innerHTML = firstContent.innerHTML.replace(/card0_/g, 'card' + creditCardCount + '_');
            
            contentsContainer.appendChild(newContent);
            
            // Switch to new tab
            switchCardTab(creditCardCount);
            creditCardCount++;
        }
        
        function deleteCreditCard(index, event) {
            event.stopPropagation(); // Prevent tab switching
            
            const cardName = 'Credit Card ' + (index + 1);
            if (confirm('Are you sure you want to delete ' + cardName + '? This action cannot be undone.')) {
                // Remove the tab
                const tabsContainer = document.getElementById('creditCardTabs');
                const tabs = tabsContainer.querySelectorAll('.instance-tab');
                if (tabs[index]) {
                    tabsContainer.removeChild(tabs[index]);
                }
                
                // Remove the content
                const contentsContainer = document.getElementById('creditCardContents');
                const content = document.getElementById('creditCard' + index);
                if (content) {
                    contentsContainer.removeChild(content);
                }
                
                // Renumber remaining tabs and contents
                renumberCreditCards();
                
                // Switch to first tab if current was removed
                if (index === 0) {
                    switchCardTab(0);
                } else {
                    switchCardTab(Math.max(0, index - 1));
                }
                
                creditCardCount--;
            }
        }
        
        function renumberCreditCards() {
            // Renumber tabs
            const tabs = document.querySelectorAll('#creditCardTabs .instance-tab');
            tabs.forEach((tab, index) => {
                const span = tab.querySelector('span');
                if (span) {
                    span.textContent = 'Credit Card ' + (index + 1);
                }
                
                // Update onclick for tab switching
                tab.onclick = function() { switchCardTab(index); };
                
                // Update delete button onclick
                const deleteBtn = tab.querySelector('.delete-btn');
                if (deleteBtn) {
                    deleteBtn.onclick = function(event) { deleteCreditCard(index, event); };
                }
            });
            
            // Renumber content IDs and form field names
            const contents = document.querySelectorAll('#creditCardContents .instance-content');
            contents.forEach((content, index) => {
                content.id = 'creditCard' + index;
                
                // Update form field names
                const inputs = content.querySelectorAll('input, select, textarea');
                inputs.forEach(input => {
                    const oldName = input.name;
                    if (oldName && oldName.startsWith('card')) {
                        const parts = oldName.split('_');
                        if (parts.length > 1) {
                            input.name = 'card' + index + '_' + parts.slice(1).join('_');
                        }
                    }
                });
            });
        }
        
        function switchCardTab(index) {
            // Hide all credit card contents
            const contents = document.querySelectorAll('#creditCardContents .instance-content');
            contents.forEach(content => content.classList.remove('active'));
            
            // Remove active from all credit card tabs
            const tabs = document.querySelectorAll('#creditCardTabs .instance-tab');
            tabs.forEach(tab => tab.classList.remove('active'));
            
            // Show selected content and activate tab
            const selectedContent = document.getElementById('creditCard' + index);
            const selectedTab = document.querySelectorAll('#creditCardTabs .instance-tab')[index];
            
            if (selectedContent) selectedContent.classList.add('active');
            if (selectedTab) selectedTab.classList.add('active');
        }

        // Camera Scanning
        let stream = null;
        let video = null;
        
        function scanCard() {
            document.getElementById('cameraModal').style.display = 'flex';
        }
        
        function closeCameraModal() {
            document.getElementById('cameraModal').style.display = 'none';
            if (stream) {
                stream.getTracks().forEach(track => track.stop());
                stream = null;
            }
        }
        
        async function startCamera() {
            try {
                stream = await navigator.mediaDevices.getUserMedia({ 
                    video: { 
                        facingMode: 'environment',
                        width: { ideal: 1280 },
                        height: { ideal: 720 }
                    } 
                });
                
                video = document.createElement('video');
                video.srcObject = stream;
                video.autoplay = true;
                video.style.width = '100%';
                video.style.height = '100%';
                video.style.objectFit = 'cover';
                
                const preview = document.getElementById('cameraPreview');
                preview.innerHTML = '';
                preview.appendChild(video);
                
            } catch (error) {
                console.error('Camera access denied:', error);
                document.getElementById('cameraPreview').innerHTML = 
                    '<div>Camera access denied. Please allow camera access and try again.</div>';
            }
        }
        
        function captureCard() {
            if (!video) {
                alert('Please start the camera first');
                return;
            }
            
            // Create canvas to capture frame
            const canvas = document.createElement('canvas');
            const context = canvas.getContext('2d');
            canvas.width = video.videoWidth;
            canvas.height = video.videoHeight;
            context.drawImage(video, 0, 0);
            
            // Simulate OCR processing
            setTimeout(() => {
                // Simulate extracted card data
                const cardData = {
                    cardNumber: '4532 **** **** 1234',
                    cardholderName: 'JOHN DOE',
                    expiryDate: '12/25',
                    cvv: '***'
                };
                
                // Fill the current credit card form
                fillCardForm(cardData);
                closeCameraModal();
                
                alert('Card details extracted successfully!');
            }, 2000);
        }
        
        function fillCardForm(cardData) {
            // Find the currently active credit card tab
            const activeContent = document.querySelector('#creditCardContents .instance-content.active');
            if (activeContent) {
                const cardNumberInput = activeContent.querySelector('input[name*="cardNumber"]');
                const cardholderInput = activeContent.querySelector('input[name*="cardholderName"]');
                const expiryInput = activeContent.querySelector('input[name*="expiryDate"]');
                
                if (cardNumberInput) cardNumberInput.value = cardData.cardNumber;
                if (cardholderInput) cardholderInput.value = cardData.cardholderName;
                if (expiryInput) expiryInput.value = cardData.expiryDate;
            }
        }

        // Digital Wallet Integration
        function connectGoogleWallet() {
            // Google Pay API integration
            if (window.google && window.google.payments) {
                const paymentDataRequest = {
                    apiVersion: 2,
                    apiVersionMinor: 0,
                    allowedPaymentMethods: [{
                        type: 'CARD',
                        parameters: {
                            allowedAuthMethods: ['PAN_ONLY', 'CRYPTOGRAM_3DS'],
                            allowedCardNetworks: ['VISA', 'MASTERCARD']
                        }
                    }],
                    transactionInfo: {
                        totalPriceStatus: 'FINAL',
                        totalPrice: '0.00',
                        currencyCode: 'USD',
                        countryCode: 'US'
                    },
                    merchantInfo: {
                        merchantName: 'Insurance App'
                    }
                };
                
                const paymentsClient = new google.payments.api.PaymentsClient({environment: 'TEST'});
                paymentsClient.loadPaymentData(paymentDataRequest)
                    .then(function(paymentData) {
                        console.log('Google Pay success:', paymentData);
                        alert('Google Wallet connected successfully!');
                    })
                    .catch(function(err) {
                        console.log('Google Pay error:', err);
                        alert('Google Wallet connection failed: ' + err.statusMessage);
                    });
            } else {
                alert('Google Pay API not available. Please ensure you are using a supported browser.');
            }
        }
        
        function connectAppleWallet() {
            // Apple Pay integration
            if (window.ApplePaySession && ApplePaySession.canMakePayments()) {
                const paymentRequest = {
                    countryCode: 'US',
                    currencyCode: 'USD',
                    supportedNetworks: ['visa', 'masterCard'],
                    merchantCapabilities: ['supports3DS'],
                    total: {
                        label: 'Insurance App',
                        amount: '0.00'
                    }
                };
                
                const session = new ApplePaySession(3, paymentRequest);
                session.onvalidatemerchant = function(event) {
                    // Validate merchant
                    session.completeMerchantValidation({});
                };
                
                session.onpaymentauthorized = function(event) {
                    session.completePayment(ApplePaySession.STATUS_SUCCESS);
                    alert('Apple Wallet connected successfully!');
                };
                
                session.begin();
            } else {
                alert('Apple Pay not available on this device/browser.');
            }
        }
        
        function connectSamsungPay() {
            // Samsung Pay integration
            if (window.SamsungPay) {
                SamsungPay.initialize({
                    merchantId: 'your-merchant-id',
                    merchantName: 'Insurance App'
                }).then(function() {
                    alert('Samsung Pay connected successfully!');
                }).catch(function(error) {
                    alert('Samsung Pay connection failed: ' + error.message);
                });
            } else {
                alert('Samsung Pay not available on this device/browser.');
            }
        }

        // Credit Card Display Functions
        function formatCardNumber(input, cardIndex) {
            let value = input.value.replace(/\D/g, '');
            let formattedValue = value.replace(/(\d{4})(?=\d)/g, '$1 ');
            input.value = formattedValue;
            
            // Update card display
            const display = document.getElementById('cardNumberDisplay' + cardIndex);
            if (value.length > 0) {
                let displayValue = value.replace(/(\d{4})(?=\d)/g, '$1 ');
                displayValue = displayValue + '•'.repeat(Math.max(0, 19 - value.length));
                display.textContent = displayValue;
            } else {
                display.textContent = '•••• •••• •••• ••••';
            }
        }
        
        function formatExpiryDate(input, cardIndex) {
            let value = input.value.replace(/\D/g, '');
            if (value.length >= 2) {
                value = value.substring(0, 2) + '/' + value.substring(2, 4);
            }
            input.value = value;
            
            // Update card display
            const display = document.getElementById('cardExpiryDisplay' + cardIndex);
            if (value.length > 0) {
                display.textContent = value;
            } else {
                display.textContent = 'MM/YY';
            }
        }
        
        function updateCardholderDisplay(input, cardIndex) {
            const display = document.getElementById('cardholderNameDisplay' + cardIndex);
            if (input.value.length > 0) {
                display.textContent = input.value.toUpperCase();
            } else {
                display.textContent = 'CARDHOLDER NAME';
            }
        }
        
        function updateCVVDisplay(input, cardIndex) {
            // CVV is not displayed on the card for security
            // This function can be used for validation or other purposes
        }

        // Email Configuration Functions
        function handleEmailProviderChange(provider) {
            const customServerSection = document.getElementById('customServerSection');
            const customServerFields = document.getElementById('customServerFields');
            
            if (provider === 'Custom IMAP') {
                customServerSection.style.display = 'block';
                customServerFields.style.display = 'block';
            } else {
                customServerSection.style.display = 'none';
                customServerFields.style.display = 'none';
            }
        }

        function handleOAuth2Toggle(enabled) {
            const oauth2Section = document.getElementById('oauth2Section');
            const oauth2Fields = document.getElementById('oauth2Fields');
            
            if (enabled) {
                oauth2Section.style.display = 'block';
                oauth2Fields.style.display = 'block';
            } else {
                oauth2Section.style.display = 'none';
                oauth2Fields.style.display = 'none';
            }
        }

        function handleAutoSyncToggle(enabled) {
            const syncFrequencyField = document.querySelector('select[name="email_syncFrequency"]');
            if (syncFrequencyField) {
                syncFrequencyField.disabled = !enabled;
            }
        }

        // Chatbot FAQ System
        const faqDatabase = {
            'how do i add an open banking account': {
                title: 'How to Add an Open Banking Account',
                content: '<strong>Step-by-step guide to add an Open Banking account:</strong>' +
                    '<ol>' +
                        '<li><strong>Navigate to Bank Accounts:</strong> Click on the "Bank Accounts" tab in the left panel</li>' +
                        '<li><strong>Add New Account:</strong> Click the "+ Add Bank Account" button</li>' +
                        '<li><strong>Select Bank:</strong> Choose your bank from the dropdown (NatWest, Barclays, HSBC, etc.)</li>' +
                        '<li><strong>Enter Account Details:</strong>' +
                            '<ul>' +
                                '<li>Account Number (8 digits)</li>' +
                                '<li>Sort Code (XX-XX-XX format)</li>' +
                                '<li>Account Holder Name</li>' +
                            '</ul>' +
                        '</li>' +
                        '<li><strong>Open Banking Setup:</strong>' +
                            '<ul>' +
                                '<li>Enable "Open Banking" option</li>' +
                                '<li>Enter your Open Banking Client ID</li>' +
                                '<li>Enter your Open Banking Client Secret</li>' +
                                '<li>Set up redirect URI for authentication</li>' +
                            '</ul>' +
                        '</li>' +
                        '<li><strong>Security:</strong> All credentials are encrypted and stored PCI-compliant</li>' +
                        '<li><strong>2FA Setup:</strong> Configure two-factor authentication if required by your bank</li>' +
                    '</ol>' +
                    '<p><strong>Tip:</strong> You can upload a bank statement to auto-fill account details!</p>'
            },
            'how do i configure oauth2 for gmail': {
                title: 'How to Configure OAuth2 for Gmail',
                content: '<strong>Complete OAuth2 setup for Gmail integration:</strong>' +
                    '<ol>' +
                        '<li><strong>Google Cloud Console Setup:</strong>' +
                            '<ul>' +
                                '<li>Go to <a href="https://console.cloud.google.com" target="_blank">Google Cloud Console</a></li>' +
                                '<li>Create a new project or select existing one</li>' +
                                '<li>Enable Gmail API</li>' +
                                '<li>Create OAuth2 credentials</li>' +
                                '<li>Set redirect URI to your application URL</li>' +
                            '</ul>' +
                        '</li>' +
                        '<li><strong>Application Configuration:</strong>' +
                            '<ul>' +
                                '<li>Navigate to "Communication Channels" → "Email" tab</li>' +
                                '<li>Select "Gmail" as email provider</li>' +
                                '<li>Enable "OAuth2 Enabled" checkbox</li>' +
                                '<li>Enter your OAuth2 Client ID from Google Console</li>' +
                                '<li>Enter your OAuth2 Client Secret</li>' +
                            '</ul>' +
                        '</li>' +
                        '<li><strong>Authentication Flow:</strong>' +
                            '<ul>' +
                                '<li>Click "Authorize" to start OAuth2 flow</li>' +
                                '<li>Grant permissions to your Gmail account</li>' +
                                '<li>Copy the authorization code</li>' +
                                '<li>Exchange code for access and refresh tokens</li>' +
                            '</ul>' +
                        '</li>' +
                        '<li><strong>Sync Configuration:</strong>' +
                            '<ul>' +
                                '<li>Enable "Auto Sync Enabled"</li>' +
                                '<li>Set sync frequency (every 15 minutes recommended)</li>' +
                                '<li>Configure which folders to sync</li>' +
                            '</ul>' +
                        '</li>' +
                    '</ol>' +
                    '<p><strong>Security Note:</strong> All OAuth2 tokens are encrypted and never stored in plain text.</p>'
            },
            'how do i add a credit card securely': {
                title: 'How to Add a Credit Card Securely',
                content: `
                    <strong>Secure credit card configuration guide:</strong>
                    <ol>
                        <li><strong>Navigate to Credit Cards:</strong> Click on the "Credit Cards" tab</li>
                        <li><strong>Add New Card:</strong> Click "+ Add Credit Card" button</li>
                        <li><strong>Card Information:</strong>
                            <ul>
                                <li>Enter card number (auto-formatted with spaces)</li>
                                <li>Select expiry date (MM/YY format)</li>
                                <li>Enter CVV (3-4 digits)</li>
                                <li>Enter cardholder name (as on card)</li>
                                <li>Enter billing address</li>
                            </ul>
                        </li>
                        <li><strong>Card Details:</strong>
                            <ul>
                                <li>Select card provider (Visa, Mastercard, Amex, etc.)</li>
                                <li>Choose card type (Credit, Debit, Business)</li>
                            </ul>
                        </li>
                        <li><strong>Security Features:</strong>
                            <ul>
                                <li>All data is PCI-compliant encrypted</li>
                                <li>Card numbers are tokenized</li>
                                <li>CVV is never stored</li>
                                <li>256-bit encryption used</li>
                            </ul>
                        </li>
                        <li><strong>Digital Wallet Integration:</strong>
                            <ul>
                                <li>Connect Google Wallet</li>
                                <li>Connect Apple Wallet</li>
                                <li>Connect Samsung Pay</li>
                            </ul>
                        </li>
                        <li><strong>Camera Scanning:</strong> Use "Camera Scan Card" for automatic data extraction</li>
                    </ol>
                    <p><strong>Tip:</strong> Upload a credit card statement to auto-fill card details!</p>
                `
            },
            'what is pci compliance': {
                title: 'What is PCI Compliance?',
                content: `
                    <strong>PCI DSS (Payment Card Industry Data Security Standard) Compliance:</strong>
                    <p>PCI compliance is a set of security standards designed to ensure that all companies that process, store, or transmit credit card information maintain a secure environment.</p>
                    
                    <h4>Key Requirements:</h4>
                    <ul>
                        <li><strong>Build and Maintain a Secure Network:</strong> Install and maintain firewall protection</li>
                        <li><strong>Protect Cardholder Data:</strong> Encrypt transmission of cardholder data</li>
                        <li><strong>Maintain Vulnerability Management:</strong> Use and regularly update anti-virus software</li>
                        <li><strong>Implement Strong Access Control:</strong> Restrict access to cardholder data</li>
                        <li><strong>Monitor and Test Networks:</strong> Track and monitor all access to network resources</li>
                        <li><strong>Maintain Information Security Policy:</strong> Maintain a policy that addresses information security</li>
                    </ul>
                    
                    <h4>Our Implementation:</h4>
                    <ul>
                        <li>256-bit AES encryption for all sensitive data</li>
                        <li>Tokenization of credit card numbers</li>
                        <li>Never store CVV codes</li>
                        <li>Secure key derivation using PBKDF2</li>
                        <li>Session-based access tokens with timeouts</li>
                        <li>Comprehensive audit logging</li>
                        <li>Regular security assessments</li>
                    </ul>
                    
                    <p><strong>Security Badges:</strong> Look for PCI, Encryption, and SSL badges throughout the interface.</p>
                `
            },
            'how do i configure outlook email': {
                title: 'How to Configure Outlook Email',
                content: `
                    <strong>Outlook email configuration guide:</strong>
                    <ol>
                        <li><strong>Microsoft Azure Setup:</strong>
                            <ul>
                                <li>Go to <a href="https://portal.azure.com" target="_blank">Azure Portal</a></li>
                                <li>Register a new application</li>
                                <li>Configure OAuth2 permissions for Microsoft Graph</li>
                                <li>Generate client ID and secret</li>
                            </ul>
                        </li>
                        <li><strong>Application Configuration:</strong>
                            <ul>
                                <li>Select "Outlook" as email provider</li>
                                <li>Enter your Outlook email address</li>
                                <li>Enable OAuth2 authentication</li>
                                <li>Enter Azure client ID and secret</li>
                            </ul>
                        </li>
                        <li><strong>Alternative: App Password (if 2FA enabled):</strong>
                            <ul>
                                <li>Generate app password in Microsoft account</li>
                                <li>Use app password instead of regular password</li>
                                <li>Disable OAuth2 if using app password</li>
                            </ul>
                        </li>
                        <li><strong>Server Settings (Manual):</strong>
                            <ul>
                                <li>IMAP Server: outlook.office365.com</li>
                                <li>IMAP Port: 993 (SSL)</li>
                                <li>SMTP Server: smtp.office365.com</li>
                                <li>SMTP Port: 587 (TLS)</li>
                            </ul>
                        </li>
                    </ol>
                `
            },
            'how do i configure pop3 gmail access': {
                title: 'How to Configure POP3 Gmail Access',
                content: `
                    <strong>POP3 Gmail configuration (Legacy method):</strong>
                    <ol>
                        <li><strong>Gmail Settings:</strong>
                            <ul>
                                <li>Enable POP3 in Gmail settings</li>
                                <li>Choose download behavior (keep/delete from server)</li>
                                <li>Enable "Less secure app access" (not recommended)</li>
                            </ul>
                        </li>
                        <li><strong>Application Configuration:</strong>
                            <ul>
                                <li>Select "Gmail" as email provider</li>
                                <li>Disable OAuth2 (use password authentication)</li>
                                <li>Enter your Gmail address and password</li>
                            </ul>
                        </li>
                        <li><strong>Server Settings:</strong>
                            <ul>
                                <li>POP3 Server: pop.gmail.com</li>
                                <li>POP3 Port: 995 (SSL)</li>
                                <li>SMTP Server: smtp.gmail.com</li>
                                <li>SMTP Port: 587 (TLS)</li>
                            </ul>
                        </li>
                        <li><strong>Security Warning:</strong>
                            <ul>
                                <li>POP3 is less secure than OAuth2</li>
                                <li>Passwords stored in plain text (encrypted)</li>
                                <li>No automatic token refresh</li>
                                <li>May require app passwords if 2FA enabled</li>
                            </ul>
                        </li>
                    </ol>
                    <p><strong>Recommendation:</strong> Use OAuth2 instead of POP3 for better security.</p>
                `
            }
        };

        // Chatbot Functions
        function askQuestion(question) {
            const chatInput = document.getElementById('chatInput');
            chatInput.value = question;
            sendChatMessage();
        }

        function handleChatKeyPress(event) {
            if (event.key === 'Enter') {
                sendChatMessage();
            }
        }

        function sendChatMessage() {
            const chatInput = document.getElementById('chatInput');
            const message = chatInput.value.trim();
            
            if (!message) return;
            
            // Add user message
            addChatMessage('user', message);
            chatInput.value = '';
            
            // Process and respond
            setTimeout(() => {
                const response = processChatMessage(message);
                addChatMessage('assistant', response);
            }, 500);
        }

        function addChatMessage(sender, content) {
            const chatMessages = document.getElementById('chatMessages');
            const messageDiv = document.createElement('div');
            messageDiv.className = `message ${sender}`;
            
            const contentDiv = document.createElement('div');
            contentDiv.className = 'message-content';
            
            if (sender === 'user') {
                contentDiv.innerHTML = `<strong>You:</strong> ${content}`;
            } else {
                contentDiv.innerHTML = content;
            }
            
            messageDiv.appendChild(contentDiv);
            chatMessages.appendChild(messageDiv);
            chatMessages.scrollTop = chatMessages.scrollHeight;
        }

        function processChatMessage(message) {
            const lowerMessage = message.toLowerCase();
            
            // Check for FAQ matches
            for (const [key, faq] of Object.entries(faqDatabase)) {
                if (lowerMessage.includes(key)) {
                    return `<strong>${faq.title}</strong>${faq.content}`;
                }
            }
            
            // Check for general help topics
            if (lowerMessage.includes('open banking') || lowerMessage.includes('bank')) {
                return `<strong>Open Banking Help:</strong>
                    <p>I can help you with Open Banking setup. Try asking:</p>
                    <ul>
                        <li>"How do I add an Open Banking account?"</li>
                        <li>"What banks support Open Banking?"</li>
                        <li>"How do I get Open Banking credentials?"</li>
                    </ul>`;
            }
            
            if (lowerMessage.includes('credit card') || lowerMessage.includes('card')) {
                return `<strong>Credit Card Help:</strong>
                    <p>I can help you with credit card configuration. Try asking:</p>
                    <ul>
                        <li>"How do I add a credit card securely?"</li>
                        <li>"What is PCI compliance?"</li>
                        <li>"How do I scan my credit card?"</li>
                    </ul>`;
            }
            
            if (lowerMessage.includes('email') || lowerMessage.includes('gmail') || lowerMessage.includes('outlook')) {
                return `<strong>Email Configuration Help:</strong>
                    <p>I can help you with email setup. Try asking:</p>
                    <ul>
                        <li>"How do I configure OAuth2 for Gmail?"</li>
                        <li>"How do I configure Outlook email?"</li>
                        <li>"How do I configure POP3 Gmail access?"</li>
                    </ul>`;
            }
            
            if (lowerMessage.includes('security') || lowerMessage.includes('pci')) {
                return `<strong>Security Help:</strong>
                    <p>I can help you with security questions. Try asking:</p>
                    <ul>
                        <li>"What is PCI compliance?"</li>
                        <li>"How secure is my data?"</li>
                        <li>"What encryption is used?"</li>
                    </ul>`;
            }
            
            // Default response
            return `<strong>I'm here to help!</strong>
                <p>I can assist you with:</p>
                <ul>
                    <li>Open Banking account setup</li>
                    <li>Credit card configuration</li>
                    <li>Email integration (Gmail, Outlook)</li>
                    <li>Security and PCI compliance</li>
                </ul>
                <p>Try asking a specific question or use the quick question buttons above!</p>`;
        }

        // Document Upload Functions
        function handleFileUpload(files) {
            if (!files || files.length === 0) return;
            
            const uploadResults = document.getElementById('uploadResults');
            const extractedData = document.getElementById('extractedData');
            
            uploadResults.style.display = 'block';
            extractedData.innerHTML = '';
            
            Array.from(files).forEach(file => {
                const fileInfo = document.createElement('div');
                fileInfo.style.cssText = 'margin-bottom: 16px; padding: 12px; background: #333333; border-radius: 6px;';
                
                const fileName = document.createElement('h5');
                fileName.textContent = `📄 ${file.name}`;
                fileName.style.cssText = 'color: #e0e0e0; margin: 0 0 8px 0; font-size: 13px;';
                
                const fileType = getFileType(file.name);
                const extractedInfo = simulateDocumentExtraction(fileType, file.name);
                
                const infoDiv = document.createElement('div');
                infoDiv.innerHTML = extractedInfo;
                infoDiv.style.cssText = 'color: #808080; font-size: 11px; line-height: 1.4;';
                
                fileInfo.appendChild(fileName);
                fileInfo.appendChild(infoDiv);
                extractedData.appendChild(fileInfo);
            });
            
            // Add message to chat
            addChatMessage('assistant', `<strong>Document Upload Complete!</strong>
                <p>I've analyzed your uploaded documents and extracted the following information. You can now:</p>
                <ul>
                    <li>Review the extracted data below</li>
                    <li>Auto-fill form fields with extracted data</li>
                    <li>Ask me questions about the extracted information</li>
                </ul>`);
        }

        function getFileType(fileName) {
            const ext = fileName.toLowerCase().split('.').pop();
            if (ext === 'pdf') return 'bank_statement';
            if (['jpg', 'jpeg', 'png'].includes(ext)) return 'screenshot';
            if (ext === 'txt') return 'config_file';
            return 'unknown';
        }

        function simulateDocumentExtraction(fileType, fileName) {
            switch (fileType) {
                case 'bank_statement':
                    return `
                        <strong>Detected: Bank Statement</strong><br>
                        ✅ Account Number: 12345678<br>
                        ✅ Sort Code: 12-34-56<br>
                        ✅ Bank Name: NatWest<br>
                        ✅ Account Holder: John Doe<br>
                        <button onclick="autoFillBankDetails()" style="margin-top: 8px; padding: 4px 8px; background: #404040; border: none; border-radius: 4px; color: #e0e0e0; cursor: pointer; font-size: 10px;">Auto-fill Bank Details</button>
                    `;
                case 'screenshot':
                    if (fileName.toLowerCase().includes('gmail') || fileName.toLowerCase().includes('google')) {
                        return `
                            <strong>Detected: Gmail Configuration Screenshot</strong><br>
                            ✅ Email Provider: Gmail<br>
                            ✅ OAuth2 Client ID: 123456789-abcdef.apps.googleusercontent.com<br>
                            ✅ Redirect URI: https://yourapp.com/oauth/callback<br>
                            <button onclick="autoFillGmailConfig()" style="margin-top: 8px; padding: 4px 8px; background: #404040; border: none; border-radius: 4px; color: #e0e0e0; cursor: pointer; font-size: 10px;">Auto-fill Gmail Config</button>
                        `;
                    }
                    return `
                        <strong>Detected: Configuration Screenshot</strong><br>
                        🔍 Analyzing image for configuration details...<br>
                        <button onclick="analyzeScreenshot()" style="margin-top: 8px; padding: 4px 8px; background: #404040; border: none; border-radius: 4px; color: #e0e0e0; cursor: pointer; font-size: 10px;">Analyze Screenshot</button>
                    `;
                case 'config_file':
                    return `
                        <strong>Detected: Configuration File</strong><br>
                        🔍 Analyzing configuration parameters...<br>
                        <button onclick="parseConfigFile()" style="margin-top: 8px; padding: 4px 8px; background: #404040; border: none; border-radius: 4px; color: #e0e0e0; cursor: pointer; font-size: 10px;">Parse Config File</button>
                    `;
                default:
                    return `
                        <strong>Unknown File Type</strong><br>
                        ⚠️ Unable to automatically extract information from this file type.<br>
                        Please ensure you're uploading supported document types.
                    `;
            }
        }

        // Auto-fill functions
        function autoFillBankDetails() {
            addChatMessage('assistant', `<strong>Auto-filling Bank Details</strong>
                <p>I've automatically filled in the bank account details from your uploaded statement:</p>
                <ul>
                    <li>✅ Account Number: 12345678</li>
                    <li>✅ Sort Code: 12-34-56</li>
                    <li>✅ Bank Name: NatWest</li>
                    <li>✅ Account Holder: John Doe</li>
                </ul>
                <p>Please review the information and click "Save" to confirm.</p>`);
        }

        function autoFillGmailConfig() {
            addChatMessage('assistant', `<strong>Auto-filling Gmail Configuration</strong>
                <p>I've automatically filled in the Gmail OAuth2 configuration:</p>
                <ul>
                    <li>✅ Email Provider: Gmail</li>
                    <li>✅ OAuth2 Client ID: 123456789-abcdef.apps.googleusercontent.com</li>
                    <li>✅ Redirect URI: https://yourapp.com/oauth/callback</li>
                </ul>
                <p>Next steps:</p>
                <ol>
                    <li>Complete the OAuth2 authorization flow</li>
                    <li>Grant permissions to your Gmail account</li>
                    <li>Configure sync settings</li>
                </ol>`);
        }

        // Initialize upload zone click handler
        document.addEventListener('DOMContentLoaded', function() {
            const uploadZone = document.getElementById('uploadZone');
            const fileInput = document.getElementById('fileInput');
            
            if (uploadZone && fileInput) {
                uploadZone.addEventListener('click', function() {
                    fileInput.click();
                });
                
                uploadZone.addEventListener('dragover', function(e) {
                    e.preventDefault();
                    uploadZone.style.borderColor = '#505050';
                    uploadZone.style.background = '#2a2a2a';
                });
                
                uploadZone.addEventListener('dragleave', function(e) {
                    e.preventDefault();
                    uploadZone.style.borderColor = '#404040';
                    uploadZone.style.background = 'transparent';
                });
                
                uploadZone.addEventListener('drop', function(e) {
                    e.preventDefault();
                    uploadZone.style.borderColor = '#404040';
                    uploadZone.style.background = 'transparent';
                    handleFileUpload(e.dataTransfer.files);
                });
            }
        });

        // Language switching
        function switchLanguage(lang) {
            // Update language buttons
            document.querySelectorAll('.lang-btn').forEach(btn => {
                btn.classList.remove('active');
            });
            document.querySelector('[data-lang="' + lang + '"]').classList.add('active');
            
            // Update all elements with data attributes
            document.querySelectorAll('[data-en][data-de]').forEach(element => {
                if (lang === 'de' && element.getAttribute('data-de')) {
                    element.textContent = element.getAttribute('data-de');
                } else if (lang === 'en' && element.getAttribute('data-en')) {
                    element.textContent = element.getAttribute('data-en');
                }
            });
            
            // Update form field labels and placeholders
            document.querySelectorAll('.form-label').forEach(label => {
                const fieldName = label.textContent.trim();
                const fieldNameLower = fieldName.toLowerCase();
                
                // Map English field names to German translations
                const translations = {
                    'email provider': 'E-Mail Anbieter',
                    'email address': 'E-Mail Adresse',
                    'oauth2 enabled': 'OAuth2 aktiviert',
                    'email password': 'E-Mail Passwort',
                    'app password': 'App Passwort',
                    'oauth2 client id': 'OAuth2 Client ID',
                    'oauth2 client secret': 'OAuth2 Client Secret',
                    'refresh token': 'Refresh Token',
                    'access token': 'Access Token',
                    'auto sync enabled': 'Auto-Sync aktiviert',
                    'sync frequency': 'Sync-Frequenz',
                    'imap server': 'IMAP Server',
                    'imap port': 'IMAP Port',
                    'imap use ssl': 'IMAP SSL verwenden',
                    'smtp server': 'SMTP Server',
                    'smtp port': 'SMTP Port',
                    'smtp use tls': 'SMTP TLS verwenden',
                    'bank name': 'Bankname',
                    'account number': 'Kontonummer',
                    'sort code': 'Sort Code',
                    'account holder name': 'Kontoinhaber Name',
                    'card provider': 'Kartenanbieter',
                    'card type': 'Kartentyp',
                    'card number': 'Kartennummer',
                    'expiry date': 'Ablaufdatum',
                    'security code': 'Sicherheitscode',
                    'cardholder name': 'Karteninhaber Name',
                    'cardholder address': 'Karteninhaber Adresse'
                };
                
                if (lang === 'de' && translations[fieldNameLower]) {
                    // Update label text
                    const labelText = label.childNodes[0];
                    if (labelText && labelText.nodeType === Node.TEXT_NODE) {
                        labelText.textContent = translations[fieldNameLower];
                    }
                    
                    // Update placeholder if input exists
                    const input = label.parentNode.querySelector('input, select, textarea');
                    if (input && input.placeholder) {
                        input.placeholder = translations[fieldNameLower];
                    }
                } else if (lang === 'en') {
                    // Revert to English
                    const originalNames = {
                        'E-Mail Anbieter': 'email provider',
                        'E-Mail Adresse': 'email address',
                        'OAuth2 aktiviert': 'oauth2 enabled',
                        'E-Mail Passwort': 'email password',
                        'App Passwort': 'app password',
                        'OAuth2 Client ID': 'oauth2 client id',
                        'OAuth2 Client Secret': 'oauth2 client secret',
                        'Refresh Token': 'refresh token',
                        'Access Token': 'access token',
                        'Auto-Sync aktiviert': 'auto sync enabled',
                        'Sync-Frequenz': 'sync frequency',
                        'IMAP Server': 'imap server',
                        'IMAP Port': 'imap port',
                        'IMAP SSL verwenden': 'imap use ssl',
                        'SMTP Server': 'smtp server',
                        'SMTP Port': 'smtp port',
                        'SMTP TLS verwenden': 'smtp use tls',
                        'Bankname': 'bank name',
                        'Kontonummer': 'account number',
                        'Sort Code': 'sort code',
                        'Kontoinhaber Name': 'account holder name',
                        'Kartenanbieter': 'card provider',
                        'Kartentyp': 'card type',
                        'Kartennummer': 'card number',
                        'Ablaufdatum': 'expiry date',
                        'Sicherheitscode': 'security code',
                        'Karteninhaber Name': 'cardholder name',
                        'Karteninhaber Adresse': 'cardholder address'
                    };
                    
                    const labelText = label.childNodes[0];
                    if (labelText && labelText.nodeType === Node.TEXT_NODE) {
                        const englishName = originalNames[labelText.textContent];
                        if (englishName) {
                            labelText.textContent = englishName;
                        }
                    }
                }
            });
            
            // Update help text
            document.querySelectorAll('.help-text').forEach(helpText => {
                const helpTextLower = helpText.textContent.toLowerCase();
                
                const helpTranslations = {
                    'email service provider (gmail, outlook, yahoo, etc.)': 'E-Mail-Dienstanbieter (Gmail, Outlook, Yahoo, etc.)',
                    'primary email address for notifications and communications': 'Primäre E-Mail-Adresse für Benachrichtigungen und Kommunikation',
                    'whether oauth2 authentication is enabled for this email provider': 'Ob OAuth2-Authentifizierung für diesen E-Mail-Anbieter aktiviert ist',
                    'email account password (for non-oauth2 authentication)': 'E-Mail-Kontopasswort (für nicht-OAuth2-Authentifizierung)',
                    'app-specific password for email access (gmail, outlook)': 'App-spezifisches Passwort für E-Mail-Zugriff (Gmail, Outlook)',
                    'oauth2 client id for email integration (from google/microsoft developer console)': 'OAuth2 Client ID für E-Mail-Integration (aus Google/Microsoft Developer Console)',
                    'oauth2 client secret for email integration (from google/microsoft developer console)': 'OAuth2 Client Secret für E-Mail-Integration (aus Google/Microsoft Developer Console)',
                    'oauth2 refresh token for automatic token renewal': 'OAuth2 Refresh Token für automatische Token-Erneuerung',
                    'oauth2 access token for email api access': 'OAuth2 Access Token für E-Mail-API-Zugriff',
                    'whether to automatically sync emails from this account': 'Ob E-Mails von diesem Konto automatisch synchronisiert werden sollen',
                    'how often to sync emails from this account': 'Wie oft E-Mails von diesem Konto synchronisiert werden sollen',
                    'imap server address for email access': 'IMAP-Server-Adresse für E-Mail-Zugriff',
                    'imap server port (usually 993 for ssl, 143 for non-ssl)': 'IMAP-Server-Port (normalerweise 993 für SSL, 143 für nicht-SSL)',
                    'whether to use ssl/tls for imap connection': 'Ob SSL/TLS für IMAP-Verbindung verwendet werden soll',
                    'smtp server address for sending emails': 'SMTP-Server-Adresse für das Senden von E-Mails',
                    'smtp server port (usually 587 for tls, 465 for ssl)': 'SMTP-Server-Port (normalerweise 587 für TLS, 465 für SSL)',
                    'whether to use tls for smtp connection': 'Ob TLS für SMTP-Verbindung verwendet werden soll'
                };
                
                if (lang === 'de' && helpTranslations[helpTextLower]) {
                    helpText.textContent = helpTranslations[helpTextLower];
                } else if (lang === 'en') {
                    // Revert to English help text
                    const englishHelpText = {
                        'E-Mail-Dienstanbieter (Gmail, Outlook, Yahoo, etc.)': 'Email service provider (Gmail, Outlook, Yahoo, etc.)',
                        'Primäre E-Mail-Adresse für Benachrichtigungen und Kommunikation': 'Primary email address for notifications and communications',
                        'Ob OAuth2-Authentifizierung für diesen E-Mail-Anbieter aktiviert ist': 'Whether OAuth2 authentication is enabled for this email provider',
                        'E-Mail-Kontopasswort (für nicht-OAuth2-Authentifizierung)': 'Email account password (for non-OAuth2 authentication)',
                        'App-spezifisches Passwort für E-Mail-Zugriff (Gmail, Outlook)': 'App-specific password for email access (Gmail, Outlook)',
                        'OAuth2 Client ID für E-Mail-Integration (aus Google/Microsoft Developer Console)': 'OAuth2 client ID for email integration (from Google/Microsoft developer console)',
                        'OAuth2 Client Secret für E-Mail-Integration (aus Google/Microsoft Developer Console)': 'OAuth2 client secret for email integration (from Google/Microsoft developer console)',
                        'OAuth2 Refresh Token für automatische Token-Erneuerung': 'OAuth2 refresh token for automatic token renewal',
                        'OAuth2 Access Token für E-Mail-API-Zugriff': 'OAuth2 access token for email API access',
                        'Ob E-Mails von diesem Konto automatisch synchronisiert werden sollen': 'Whether to automatically sync emails from this account',
                        'Wie oft E-Mails von diesem Konto synchronisiert werden sollen': 'How often to sync emails from this account',
                        'IMAP-Server-Adresse für E-Mail-Zugriff': 'IMAP server address for email access',
                        'IMAP-Server-Port (normalerweise 993 für SSL, 143 für nicht-SSL)': 'IMAP server port (usually 993 for SSL, 143 for non-SSL)',
                        'Ob SSL/TLS für IMAP-Verbindung verwendet werden soll': 'Whether to use SSL/TLS for IMAP connection',
                        'SMTP-Server-Adresse für das Senden von E-Mails': 'SMTP server address for sending emails',
                        'SMTP-Server-Port (normalerweise 587 für TLS, 465 für SSL)': 'SMTP server port (usually 587 for TLS, 465 for SSL)',
                        'Ob TLS für SMTP-Verbindung verwendet werden soll': 'Whether to use TLS for SMTP connection'
                    };
                    
                    const englishText = englishHelpText[helpText.textContent];
                    if (englishText) {
                        helpText.textContent = englishText;
                    }
                }
            });
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
