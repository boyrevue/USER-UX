package formgenerator

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
	Property                  string
	Label                     string
	Type                      string
	Required                  bool
	Pattern                   string
	HelpText                  string
	Options                   []FieldOption
	Min                       string
	Max                       string
	DefaultValue              string
	Conditional               string
	ConditionalVisibility     string
	ConditionalVisibilityRule string
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
	var currentTriple strings.Builder
	var inTriple bool

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

		// Handle multi-line triples
		if strings.Contains(line, " a owl:DatatypeProperty") {
			// Start of a new triple
			if inTriple {
				p.parseTriple(currentTriple.String())
			}
			currentTriple.Reset()
			currentTriple.WriteString(line)
			inTriple = true
		} else if inTriple {
			// Continue the current triple
			currentTriple.WriteString(" ")
			currentTriple.WriteString(line)

			// Check if triple ends
			if strings.HasSuffix(line, ".") {
				p.parseTriple(currentTriple.String())
				currentTriple.Reset()
				inTriple = false
			}
		}
	}

	// Handle any remaining triple
	if inTriple {
		p.parseTriple(currentTriple.String())
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
	// Extract the subject (field name) from the first part
	parts := strings.Fields(line)
	if len(parts) < 3 {
		return
	}

	subject := parts[0]
	if !strings.Contains(subject, "autoins:") {
		return
	}

	fieldName := strings.TrimPrefix(subject, "autoins:")

	// Initialize field if not exists
	if _, exists := p.Fields[fieldName]; !exists {
		p.Fields[fieldName] = &TTLField{Property: fieldName}
	}

	// Parse the rest of the triple for properties
	// Split by semicolons to handle multiple properties
	propertyParts := strings.Split(line, ";")

	for _, propPart := range propertyParts {
		propPart = strings.TrimSpace(propPart)
		if propPart == "" {
			continue
		}

		// Parse rdfs:label
		if strings.Contains(propPart, "rdfs:label") {
			labelMatch := regexp.MustCompile(`rdfs:label\s+"([^"]+)"`)
			if matches := labelMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				p.Fields[fieldName].Label = matches[1]
			}
		}

		// Parse rdfs:range (field type)
		if strings.Contains(propPart, "rdfs:range") {
			rangeMatch := regexp.MustCompile(`rdfs:range\s+([^\s;]+)`)
			if matches := rangeMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				p.Fields[fieldName].Type = p.parseFieldType(matches[1])
			}
		}

		// Parse validation pattern
		if strings.Contains(propPart, "autoins:validationPattern") {
			patternMatch := regexp.MustCompile(`autoins:validationPattern\s+"([^"]+)"`)
			if matches := patternMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				p.Fields[fieldName].Pattern = matches[1]
			}
		}

		// Parse help text
		if strings.Contains(propPart, "autoins:formHelpText") {
			helpMatch := regexp.MustCompile(`autoins:formHelpText\s+"([^"]+)"`)
			if matches := helpMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				p.Fields[fieldName].HelpText = matches[1]
			}
		}

		// Parse required field
		if strings.Contains(propPart, "autoins:isRequired") {
			requiredMatch := regexp.MustCompile(`autoins:isRequired\s+"([^"]+)"`)
			if matches := requiredMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				p.Fields[fieldName].Required = matches[1] == "true"
			}
		}

		// Parse enumeration values
		if strings.Contains(propPart, "autoins:enumerationValues") {
			enumMatch := regexp.MustCompile(`autoins:enumerationValues\s+\(([^)]+)\)`)
			if matches := enumMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				enumValues := strings.TrimSpace(matches[1])
				// Remove quotes and split by spaces
				enumValues = strings.ReplaceAll(enumValues, "\"", "")
				values := strings.Fields(enumValues)
				for _, value := range values {
					p.Fields[fieldName].Options = append(p.Fields[fieldName].Options, FieldOption{
						Value: value,
						Label: value,
					})
				}
			}
		}

		// Parse conditional visibility
		if strings.Contains(propPart, "autoins:conditionalVisibility") {
			visibilityMatch := regexp.MustCompile(`autoins:conditionalVisibility\s+"([^"]+)"`)
			if matches := visibilityMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				p.Fields[fieldName].ConditionalVisibility = matches[1]
			}
		}

		// Parse conditional visibility rule
		if strings.Contains(propPart, "autoins:conditionalVisibilityRule") {
			ruleMatch := regexp.MustCompile(`autoins:conditionalVisibilityRule\s+"([^"]+)"`)
			if matches := ruleMatch.FindStringSubmatch(propPart); len(matches) > 1 {
				p.Fields[fieldName].ConditionalVisibilityRule = matches[1]
			}
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
	// Add enumeration options to fields
	p.addEnumerationOptions()

	// Define category mappings based on field prefixes
	categoryMappings := map[string]string{
		// Driver fields
		"firstName":                "drivers",
		"lastName":                 "drivers",
		"middleName":               "drivers",
		"dateOfBirth":              "drivers",
		"gender":                   "drivers",
		"title":                    "drivers",
		"email":                    "drivers",
		"mobilePhone":              "drivers",
		"homePhone":                "drivers",
		"addressLine1":             "drivers",
		"addressLine2":             "drivers",
		"city":                     "drivers",
		"postcode":                 "drivers",
		"licenceType":              "drivers",
		"licenceNumber":            "drivers",
		"driverClassification":     "drivers",
		"relationshipToMainDriver": "drivers",
		"livingAtSameAddress":      "drivers",
		"ukResidentSince":          "drivers",
		"isHomeowner":              "drivers",
		"maritalStatus":            "drivers",
		"employmentStatus":         "drivers",
		"occupation":               "drivers",
		"industry":                 "drivers",

		// Vehicle fields
		"make":                 "vehicle",
		"model":                "vehicle",
		"variant":              "vehicle",
		"yearOfManufacture":    "vehicle",
		"registrationNumber":   "vehicle",
		"vinNumber":            "vehicle",
		"colour":               "vehicle",
		"bodyType":             "vehicle",
		"engineSize":           "vehicle",
		"fuelType":             "vehicle",
		"transmission":         "vehicle",
		"numberOfSeats":        "vehicle",
		"numberOfDoors":        "vehicle",
		"vehicleValue":         "vehicle",
		"purchasePrice":        "vehicle",
		"dateOfPurchase":       "vehicle",
		"purchaseType":         "vehicle",
		"annualMileage":        "vehicle",
		"overnightLocation":    "vehicle",
		"overnightPostcode":    "vehicle",
		"daytimeLocation":      "vehicle",
		"daytimePostcode":      "vehicle",
		"imported":             "vehicle",
		"importType":           "vehicle",
		"leftHandDrive":        "vehicle",
		"alarm":                "vehicle",
		"immobiliser":          "vehicle",
		"tracker":              "vehicle",
		"dashCam":              "vehicle",
		"vehicleModifications": "vehicle",

		// Vehicle Modification Specific Fields
		"engineTuningBrand": "vehicle",
		"engineTuningPower": "vehicle",
		"engineTuningDate":  "vehicle",
		"exhaustBrand":      "vehicle",
		"exhaustType":       "vehicle",
		"exhaustDate":       "vehicle",
		"suspensionBrand":   "vehicle",
		"suspensionType":    "vehicle",
		"suspensionDrop":    "vehicle",
		"suspensionDate":    "vehicle",
		"wheelsBrand":       "vehicle",
		"wheelsSize":        "vehicle",
		"tyresBrand":        "vehicle",
		"wheelsDate":        "vehicle",
		"bodyKitBrand":      "vehicle",
		"bodyKitType":       "vehicle",
		"bodyKitDate":       "vehicle",
		"audioBrand":        "vehicle",
		"audioValue":        "vehicle",
		"audioDate":         "vehicle",
		"chipBrand":         "vehicle",
		"chipPower":         "vehicle",
		"chipDate":          "vehicle",
		"turboBrand":        "vehicle",
		"turboType":         "vehicle",
		"turboDate":         "vehicle",

		// Coverage fields
		"coverType":                   "coverage",
		"classOfUse":                  "coverage",
		"thirdPartyPropertyLimit":     "coverage",
		"thirdPartyBodilyInjuryLimit": "coverage",
		"accidentalDamageCovered":     "coverage",
		"fireTheftCovered":            "coverage",
		"maliciousDamageCovered":      "coverage",
		"courtesyCar":                 "coverage",
		"approvedRepairerRequired":    "coverage",
		"windscreenCover":             "coverage",
		"windscreenRepairExcess":      "coverage",
		"windscreenReplacementExcess": "coverage",
		"audioNonStandardLimit":       "coverage",

		// Claims fields
		"atFaultAccidents":       "claims",
		"notAtFaultAccidents":    "claims",
		"trafficViolations":      "claims",
		"duiConvictions":         "claims",
		"licenseSuspensions":     "claims",
		"insuranceCancellations": "claims",

		// Payment fields
		"paymentMethod":    "payment",
		"paymentFrequency": "payment",
		"cardNumber":       "payment",
		"expiryDate":       "payment",
		"cvv":              "payment",
		"billingAddress":   "payment",
		"autoPay":          "payment",

		// Settings fields
		"emailAccount":                 "settings",
		"bankAccount":                  "settings",
		"creditCard":                   "settings",
		"notificationSettings":         "settings",
		"privacySettings":              "settings",
		"languagePreference":           "settings",
		"timezone":                     "settings",
		"currency":                     "settings",
		"communicationPreference":      "settings",
		"paperlessDocuments":           "settings",
		"marketingCommunications":      "settings",
		"emergencyContactName":         "settings",
		"emergencyContactPhone":        "settings",
		"emergencyContactRelationship": "settings",
		"preferredAgent":               "settings",

		// Summary fields
		"policyType":       "summary",
		"policyTerm":       "summary",
		"effectiveDate":    "summary",
		"expirationDate":   "summary",
		"estimatedPremium": "summary",
		"discountsApplied": "summary",
		"policyNumber":     "summary",
		"productName":      "summary",
		"quoteReference":   "summary",
	}

	// Initialize categories
	categoryConfigs := map[string]struct {
		title string
		icon  string
		order int
	}{
		"drivers":  {"Driver Details", "👥", 1},
		"vehicle":  {"Vehicle Information", "🚗", 2},
		"coverage": {"Coverage Options", "🛡️", 3},
		"claims":   {"Claims History", "📋", 4},
		"payment":  {"Payment Information", "💳", 5},
		"settings": {"Settings", "⚙️", 6},
		"summary":  {"Summary", "📊", 7},
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

	// German translations for common fields
	germanTranslations := map[string]string{
		"title":                        "Titel",
		"firstName":                    "Vorname",
		"lastName":                     "Nachname",
		"middleName":                   "Zweiter Vorname",
		"dateOfBirth":                  "Geburtsdatum",
		"gender":                       "Geschlecht",
		"email":                        "E-Mail-Adresse",
		"mobilePhone":                  "Mobiltelefon",
		"homePhone":                    "Festnetz",
		"addressLine1":                 "Adresszeile 1",
		"addressLine2":                 "Adresszeile 2",
		"city":                         "Stadt",
		"postcode":                     "Postleitzahl",
		"licenceType":                  "Führerscheintyp",
		"licenceNumber":                "Führerscheinnummer",
		"maritalStatus":                "Familienstand",
		"employmentStatus":             "Beschäftigungsstatus",
		"occupation":                   "Beruf",
		"industry":                     "Branche",
		"registrationNumber":           "Kennzeichen",
		"make":                         "Marke",
		"model":                        "Modell",
		"yearOfManufacture":            "Baujahr",
		"engineSize":                   "Hubraum",
		"fuelType":                     "Kraftstoffart",
		"transmission":                 "Getriebe",
		"bodyType":                     "Karosserieform",
		"colour":                       "Farbe",
		"annualMileage":                "Jahreskilometer",
		"overnightLocation":            "Nachtparkplatz",
		"overnightPostcode":            "Nachtpostleitzahl",
		"daytimeLocation":              "Tagparkplatz",
		"daytimePostcode":              "Tagpostleitzahl",
		"alarm":                        "Alarmanlage",
		"immobiliser":                  "Wegfahrsperre",
		"tracker":                      "GPS-Tracker",
		"dashCam":                      "Dashcam",
		"vehicleModifications":         "Fahrzeugmodifikationen",
		"paymentMethod":                "Zahlungsmethode",
		"paymentFrequency":             "Zahlungsfrequenz",
		"cardNumber":                   "Kartennummer",
		"expiryDate":                   "Ablaufdatum",
		"cvv":                          "CVV",
		"billingAddress":               "Rechnungsadresse",
		"autoPay":                      "Automatische Zahlung",
		"communicationPreference":      "Kommunikationspräferenz",
		"paperlessDocuments":           "Papierlose Dokumente",
		"marketingCommunications":      "Marketing-Kommunikation",
		"emergencyContactName":         "Notfallkontakt Name",
		"emergencyContactPhone":        "Notfallkontakt Telefon",
		"emergencyContactRelationship": "Notfallkontakt Beziehung",
		"preferredAgent":               "Bevorzugter Agent",
		"policyType":                   "Policentyp",
		"policyTerm":                   "Policenlaufzeit",
		"effectiveDate":                "Gültigkeitsdatum",
		"expirationDate":               "Ablaufdatum",
		"estimatedPremium":             "Geschätzte Prämie",
		"discountsApplied":             "Angewandte Rabatte",
		"policyNumber":                 "Policennummer",
		"productName":                  "Produktname",
		"quoteReference":               "Angebotsreferenz",
	}

	// Generate HTML template
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Ontology-Driven Insurance Forms</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
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

        /* Reset & Base */
        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        *:focus {
            outline: 2px solid var(--input-focus);
            outline-offset: 2px;
        }

        body {
            background: var(--bg);
            color: var(--fg);
            font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen', 'Ubuntu', 'Helvetica Neue', Arial, sans-serif;
            font-size: 14px;
            line-height: 1.6;
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
            min-height: 100vh;
            padding: 20px;
        }
        .container { 
            max-width: 1400px; 
            margin: 0 auto; 
            background: var(--card-bg); 
            border: 1px solid var(--card-border);
            border-radius: 0;
            box-shadow: 0 20px 40px var(--shadow-strong);
            overflow: hidden;
        }
        .header { 
            background: var(--gray-800); 
            color: var(--fg); 
            padding: 30px; 
            text-align: center;
            position: relative;
            border-bottom: 2px solid var(--card-border);
            box-shadow: 0 1px 3px var(--shadow);
        }
        .header h1 { 
            margin: 0; 
            font-size: 24px; 
            font-weight: 700; 
            letter-spacing: -0.02em;
        }
        .header p { 
            margin: 10px 0 0 0; 
            color: var(--muted); 
            font-size: 14px;
        }
        
        .language-toggle {
            position: absolute;
            top: 20px;
            right: 20px;
            display: flex;
            gap: 10px;
            align-items: center;
            background: var(--card-bg);
            padding: 8px;
            border-radius: 0;
            border: 2px solid var(--card-border);
            box-shadow: 0 2px 4px var(--shadow);
        }
        
        .language-btn {
            background: var(--button-bg);
            color: var(--button-fg);
            border: 2px solid var(--button-bg);
            border-radius: 0;
            padding: 8px 16px;
            cursor: pointer;
            font-size: 12px;
            font-weight: 600;
            letter-spacing: 0.02em;
            text-transform: uppercase;
            transition: all 0.15s ease;
        }
        
        .language-btn.active {
            background: var(--button-fg);
            color: var(--button-bg);
            font-weight: 700;
        }
        
        .language-btn:hover:not(.active) {
            transform: translateY(-1px);
            box-shadow: 0 4px 8px var(--shadow-strong);
        }
        
        .main-content {
            display: grid;
            grid-template-columns: 300px 1fr;
            min-height: 600px;
            background: var(--card-bg);
        }
        
        .left-panel {
            background: var(--gray-800);
            border-right: 2px solid var(--card-border);
            padding: 20px 0;
        }

        .category-tabs {
            display: flex;
            flex-direction: column;
            gap: 2px;
        }

        .category-tab {
            background: var(--gray-700);
            color: var(--muted);
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
            background: var(--fg);
            color: var(--bg);
            font-weight: 700;
            filter: grayscale(0%);
        }

        .category-tab:hover:not(.active) {
            background: var(--gray-600);
            color: var(--fg);
            transform: translateY(-1px);
        }

        .category-content {
            display: none;
        }

        .category-content.active {
            display: block;
        }

        .right-panel {
            background: var(--card-bg);
            padding: 30px;
            overflow-y: auto;
        }

        .driver-tabs {
            display: flex;
            gap: 5px;
            margin-bottom: 20px;
            flex-wrap: wrap;
            background: #333333;
            padding: 15px;
            border-radius: 8px;
            border: 1px solid #404040;
        }

        .driver-tab {
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

        .driver-tab.active {
            background: #505050;
            color: #e0e0e0;
            border-color: #00ff00;
            box-shadow: 0 0 10px rgba(0, 255, 0, 0.3);
        }

        .driver-tab:hover:not(.active) {
            background: #505050;
            color: #e0e0e0;
        }

        .driver-content {
            display: none;
        }

        .driver-content.active {
            display: block;
        }

        .right-panel {
            background: var(--card-bg);
            padding: 30px;
            overflow-y: auto;
        }

        .form-section {
            margin-bottom: 30px;
        }

        .form-section h3 {
            color: var(--fg);
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
            color: var(--fg);
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
            background: var(--input-bg);
            color: var(--fg);
            border: 2px solid var(--input-border);
            border-radius: 0;
            padding: 10px 12px;
            font-size: 14px;
            font-family: inherit;
            transition: all 0.15s ease;
            -webkit-appearance: none;
            -moz-appearance: none;
            appearance: none;
        }

        /* Grayscale icons */
        .category-tab,
        .driver-tab {
            filter: grayscale(100%);
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

        .modal-content::-webkit-scrollbar {
            width: 6px;
        }

        .modal-content::-webkit-scrollbar-track {
            background: #333333;
            border-radius: 3px;
        }

        .modal-content::-webkit-scrollbar-thumb {
            background: #404040;
            border-radius: 3px;
        }

        .modal-content::-webkit-scrollbar-thumb:hover {
            background: #00ff00;
        }

        /* Upload modal specific styles */
        .upload-area {
            border: 2px dashed #404040;
            border-radius: 6px;
            padding: 30px;
            text-align: center;
            margin-bottom: 20px;
            transition: all 0.3s ease;
        }

        .upload-area:hover {
            border-color: #00ff00;
            background: #333333;
        }

        .upload-area.dragover {
            border-color: #00ff00;
            background: #2a2a2a;
        }

        .upload-icon {
            font-size: 48px;
            color: #808080;
            margin-bottom: 16px;
        }

        .upload-text {
            color: #e0e0e0;
            margin-bottom: 8px;
        }

        .upload-help {
            color: #808080;
            font-size: 12px;
        }

        .document-types {
            background: #333333;
            border-radius: 6px;
            padding: 16px;
            margin-top: 16px;
        }

        .document-type {
            display: flex;
            align-items: center;
            margin-bottom: 12px;
            padding: 8px;
            border-radius: 4px;
            background: #404040;
        }

        .document-icon {
            font-size: 20px;
            color: #808080;
            margin-right: 12px;
        }

        .document-info h4 {
            color: #e0e0e0;
            font-size: 13px;
            margin: 0 0 4px 0;
        }

        .document-info p {
            color: #808080;
            font-size: 11px;
            margin: 0;
        }

        /* Chat modal specific styles */
        .chat-messages {
            height: calc(100% - 60px);
            overflow-y: auto;
            margin-bottom: 16px;
        }

        .chat-input-area {
            display: flex;
            gap: 8px;
        }

        .chat-input {
            flex: 1;
            padding: 8px 12px;
            border: 1px solid #404040;
            border-radius: 4px;
            background: #333333;
            color: #e0e0e0;
            font-size: 13px;
        }

        .chat-send {
            padding: 8px 16px;
            background: #404040;
            border: none;
            border-radius: 4px;
            color: #e0e0e0;
            cursor: pointer;
            font-size: 13px;
            transition: all 0.3s ease;
        }

        .chat-send:hover {
            background: #505050;
        }

        .chat-message {
            margin-bottom: 12px;
            padding: 10px;
            border-radius: 6px;
            font-size: 13px;
        }

        .chat-message.user {
            background: #404040;
            margin-left: 20px;
        }

        .chat-message.assistant {
            background: #333333;
            margin-right: 20px;
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

        .form-field input:hover,
        .form-field select:hover,
        .form-field textarea:hover {
            border-color: var(--gray-500);
        }

        .form-field input:focus,
        .form-field select:focus,
        .form-field textarea:focus {
            background: var(--bg);
            border-color: var(--fg);
            outline: none;
            box-shadow: 0 0 0 3px var(--shadow);
        }

        /* Custom Select Arrow */
        .form-field select {
            background-image: url("data:image/svg+xml,%3Csvg width='12' height='8' viewBox='0 0 12 8' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M1 1.5L6 6.5L11 1.5' stroke='%23000000' stroke-width='2' stroke-linecap='square'/%3E%3C/svg%3E");
            background-repeat: no-repeat;
            background-position: right 12px center;
            padding-right: 40px;
            cursor: pointer;
        }

        @media (prefers-color-scheme: dark) {
            .form-field select {
                background-image: url("data:image/svg+xml,%3Csvg width='12' height='8' viewBox='0 0 12 8' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M1 1.5L6 6.5L11 1.5' stroke='%23FFFFFF' stroke-width='2' stroke-linecap='square'/%3E%3C/svg%3E");
            }
        }

        .form-field .help-text {
            color: var(--muted);
            font-size: 11px;
            margin-top: 4px;
            font-style: italic;
        }

        .required {
            color: var(--pure-black);
            font-weight: 700;
        }

        @media (prefers-color-scheme: dark) {
            .required {
                color: var(--pure-white);
            }
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
            <div class="language-toggle">
                <div class="language-btn active" onclick="switchLanguage('en')">🇬🇧 English</div>
                <div class="language-btn" onclick="switchLanguage('de')">🇩🇪 Deutsch</div>
            </div>
            <h1 data-en="🚗📄 Ontology-Driven Insurance Forms" data-de="🚗📄 Ontologie-gesteuerte Versicherungsformulare">🚗📄 Ontology-Driven Insurance Forms</h1>
            <p data-en="Dynamically generated from TTL ontology" data-de="Dynamisch aus TTL-Ontologie generiert">Dynamically generated from TTL ontology</p>
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
            </div>

            <div class="right-panel">
                <!-- Category Contents -->
                {{range $index, $category := .Categories}}
                <div class="category-content {{if eq $index 0}}active{{end}}" id="{{$category.ID}}Content">
                    {{if eq $category.ID "drivers"}}
                    <!-- Driver tabs for multiple drivers -->
                    <div class="driver-tabs">
                        <div class="driver-tab active" onclick="switchDriver(1)">👤 <span data-en="Driver 1" data-de="Fahrer 1">Driver 1</span></div>
                        <div class="driver-tab" onclick="switchDriver(2)">👤 <span data-en="Driver 2" data-de="Fahrer 2">Driver 2</span></div>
                        <div class="driver-tab" onclick="switchDriver(3)">👤 <span data-en="Driver 3" data-de="Fahrer 3">Driver 3</span></div>
                        <div class="driver-tab" onclick="switchDriver(4)">👤 <span data-en="Driver 4" data-de="Fahrer 4">Driver 4</span></div>
                    </div>
                    
                    <!-- Driver 1 content -->
                    <div class="driver-content active" id="driver1Content">
                        <div class="form-section">
                            <h3>👤 <span data-en="Driver 1 Details" data-de="Fahrer 1 Details">Driver 1 Details</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver1_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver1_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver1_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver1_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver1_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver1_{{$field.Property}}" {{if $field.Required}}required{{end}}>
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver1_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"
                                              {{if $field.Required}}required{{end}}></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver1_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver1_{{$field.Property}}" 
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
                    
                    <!-- Driver 2 content -->
                    <div class="driver-content" id="driver2Content">
                        <div class="form-section">
                            <h3>👤 <span data-en="Driver 2 Details" data-de="Fahrer 2 Details">Driver 2 Details</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver2_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver2_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver2_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver2_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver2_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver2_{{$field.Property}}">
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver2_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver2_{{$field.Property}}">
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver2_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}">
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    <!-- Driver 3 content -->
                    <div class="driver-content" id="driver3Content">
                        <div class="form-section">
                            <h3>👤 <span data-en="Driver 3 Details" data-de="Fahrer 3 Details">Driver 3 Details</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver3_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver3_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver3_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver3_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver3_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver3_{{$field.Property}}">
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver3_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver3_{{$field.Property}}">
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver3_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}">
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    <!-- Driver 4 content -->
                    <div class="driver-content" id="driver4Content">
                        <div class="form-section">
                            <h3>👤 <span data-en="Driver 4 Details" data-de="Fahrer 4 Details">Driver 4 Details</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver4_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver4_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver4_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver4_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver4_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver4_{{$field.Property}}">
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver4_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver4_{{$field.Property}}">
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver4_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}">
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    {{else if eq $category.ID "vehicle"}}
                    <!-- Vehicle tabs -->
                    <div class="driver-tabs">
                        <div class="driver-tab active" onclick="switchVehicleTab('details')">🚗 <span data-en="Vehicle Details" data-de="Fahrzeugdetails">Vehicle Details</span></div>
                        <div class="driver-tab" onclick="switchVehicleTab('modifications')">🔧 <span data-en="Modifications" data-de="Modifikationen">Modifications</span></div>
                    </div>
                    
                    <!-- Vehicle Details content -->
                    <div class="driver-content active" id="vehicleDetailsContent">
                        <div class="form-section">
                            <h3>🚗 <span data-en="Vehicle Details" data-de="Fahrzeugdetails">Vehicle Details</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                {{if not (eq $field.Property "vehicleModifications")}}
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
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    <!-- Vehicle Modifications content -->
                    <div class="driver-content" id="vehicleModificationsContent">
                        <div class="form-section">
                            <h3>🔧 <span data-en="Vehicle Modifications" data-de="Fahrzeugmodifikationen">Vehicle Modifications</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                {{if eq $field.Property "vehicleModifications"}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "select"}}
                                    <select name="{{$field.Property}}" {{if $field.Required}}required{{end}}>
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
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
                    
                    <!-- Driver 2 content -->
                    <div class="driver-content" id="driver2Content">
                        <div class="form-section">
                            <h3>👤 Driver 2 Details</h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver2_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver2_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver2_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver2_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver2_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver2_{{$field.Property}}">
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver2_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver2_{{$field.Property}}">
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver2_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}">
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    <!-- Driver 3 content -->
                    <div class="driver-content" id="driver3Content">
                        <div class="form-section">
                            <h3>👤 Driver 3 Details</h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver3_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver3_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver3_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver3_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver3_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver3_{{$field.Property}}">
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver3_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver3_{{$field.Property}}">
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver3_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}">
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    <!-- Driver 4 content -->
                    <div class="driver-content" id="driver4Content">
                        <div class="form-section">
                            <h3>👤 Driver 4 Details</h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver4_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver4_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver4_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver4_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver4_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver4_{{$field.Property}}">
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver4_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver4_{{$field.Property}}">
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver4_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}">
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    {{else if eq $category.ID "vehicle"}}
                    <!-- Vehicle Details content -->
                    <div class="driver-content active" id="vehicleDetailsContent">
                        <div class="form-section">
                            <h3>🚗 <span data-en="Vehicle Details" data-de="Fahrzeugdetails">Vehicle Details</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                {{if not (eq $field.Property "vehicleModifications")}}
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
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    <!-- Vehicle Modifications content -->
                    <div class="driver-content" id="vehicleModificationsContent">
                        <div class="form-section">
                            <h3>🔧 <span data-en="Vehicle Modifications" data-de="Fahrzeugmodifikationen">Vehicle Modifications</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                {{if eq $field.Property "vehicleModifications"}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "select"}}
                                    <select name="{{$field.Property}}" id="vehicleModificationsSelect" onchange="showModificationForm(this.value)" {{if $field.Required}}required{{end}}>
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                                {{end}}
                            </div>
                            
                            <!-- Dynamic Modification Forms -->
                            <div id="modificationForms" style="display: none; margin-top: 20px;">
                                <!-- Engine Tuning Form -->
                                <div id="engineTuningForm" class="modification-form" style="display: none;">
                                    <h4>🔧 <span data-en="Engine Tuning Details" data-de="Motor-Tuning Details">Engine Tuning Details</span></h4>
                                    <div class="form-grid">
                                        {{range $field := $category.Fields}}
                                        {{if or (eq $field.Property "engineTuningBrand") (eq $field.Property "engineTuningPower") (eq $field.Property "engineTuningDate")}}
                                        <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                            <label>{{$field.Label}}</label>
                                            {{if eq $field.Type "select"}}
                                            <select name="{{$field.Property}}">
                                                <option value="">Select {{$field.Label}}</option>
                                                {{range $option := $field.Options}}
                                                <option value="{{$option.Value}}">{{$option.Label}}</option>
                                                {{end}}
                                            </select>
                                            {{else if eq $field.Type "date"}}
                                            <input type="date" name="{{$field.Property}}">
                                            {{else}}
                                            <input type="text" name="{{$field.Property}}" placeholder="{{$field.Label}}">
                                            {{end}}
                                            {{if $field.HelpText}}
                                            <div class="help-text">{{$field.HelpText}}</div>
                                            {{end}}
                                        </div>
                                        {{end}}
                                        {{end}}
                                    </div>
                                </div>
                                
                                <!-- Exhaust System Form -->
                                <div id="exhaustForm" class="modification-form" style="display: none;">
                                    <h4>🔧 <span data-en="Exhaust System Details" data-de="Auspuffanlage Details">Exhaust System Details</span></h4>
                                    <div class="form-grid">
                                        {{range $field := $category.Fields}}
                                        {{if or (eq $field.Property "exhaustBrand") (eq $field.Property "exhaustType") (eq $field.Property "exhaustDate")}}
                                        <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                            <label>{{$field.Label}}</label>
                                            {{if eq $field.Type "select"}}
                                            <select name="{{$field.Property}}">
                                                <option value="">Select {{$field.Label}}</option>
                                                {{range $option := $field.Options}}
                                                <option value="{{$option.Value}}">{{$option.Label}}</option>
                                                {{end}}
                                            </select>
                                            {{else if eq $field.Type "date"}}
                                            <input type="date" name="{{$field.Property}}">
                                            {{else}}
                                            <input type="text" name="{{$field.Property}}" placeholder="{{$field.Label}}">
                                            {{end}}
                                            {{if $field.HelpText}}
                                            <div class="help-text">{{$field.HelpText}}</div>
                                            {{end}}
                                        </div>
                                        {{end}}
                                        {{end}}
                                    </div>
                                </div>
                                
                                <!-- Suspension Form -->
                                <div id="suspensionForm" class="modification-form" style="display: none;">
                                    <h4>🔧 <span data-en="Suspension Details" data-de="Federung Details">Suspension Details</span></h4>
                                    <div class="form-grid">
                                        {{range $field := $category.Fields}}
                                        {{if or (eq $field.Property "suspensionBrand") (eq $field.Property "suspensionType") (eq $field.Property "suspensionDrop") (eq $field.Property "suspensionDate")}}
                                        <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                            <label>{{$field.Label}}</label>
                                            {{if eq $field.Type "select"}}
                                            <select name="{{$field.Property}}">
                                                <option value="">Select {{$field.Label}}</option>
                                                {{range $option := $field.Options}}
                                                <option value="{{$option.Value}}">{{$option.Label}}</option>
                                                {{end}}
                                            </select>
                                            {{else if eq $field.Type "date"}}
                                            <input type="date" name="{{$field.Property}}">
                                            {{else}}
                                            <input type="text" name="{{$field.Property}}" placeholder="{{$field.Label}}">
                                            {{end}}
                                            {{if $field.HelpText}}
                                            <div class="help-text">{{$field.HelpText}}</div>
                                            {{end}}
                                        </div>
                                        {{end}}
                                        {{end}}
                                    </div>
                                </div>
                                
                                <!-- Wheels/Tyres Form -->
                                <div id="wheelsForm" class="modification-form" style="display: none;">
                                    <h4>🔧 <span data-en="Wheels & Tyres Details" data-de="Räder & Reifen Details">Wheels & Tyres Details</span></h4>
                                    <div class="form-grid">
                                        {{range $field := $category.Fields}}
                                        {{if or (eq $field.Property "wheelsBrand") (eq $field.Property "wheelsSize") (eq $field.Property "tyresBrand") (eq $field.Property "wheelsDate")}}
                                        <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                            <label>{{$field.Label}}</label>
                                            {{if eq $field.Type "select"}}
                                            <select name="{{$field.Property}}">
                                                <option value="">Select {{$field.Label}}</option>
                                                {{range $option := $field.Options}}
                                                <option value="{{$option.Value}}">{{$option.Label}}</option>
                                                {{end}}
                                            </select>
                                            {{else if eq $field.Type "date"}}
                                            <input type="date" name="{{$field.Property}}">
                                            {{else}}
                                            <input type="text" name="{{$field.Property}}" placeholder="{{$field.Label}}">
                                            {{end}}
                                            {{if $field.HelpText}}
                                            <div class="help-text">{{$field.HelpText}}</div>
                                            {{end}}
                                        </div>
                                        {{end}}
                                        {{end}}
                                    </div>
                                </div>
                                
                                <!-- Body Kit Form -->
                                <div id="bodyKitForm" class="modification-form" style="display: none;">
                                    <h4>🔧 <span data-en="Body Kit Details" data-de="Bodykit Details">Body Kit Details</span></h4>
                                    <div class="form-grid">
                                        {{range $field := $category.Fields}}
                                        {{if or (eq $field.Property "bodyKitBrand") (eq $field.Property "bodyKitType") (eq $field.Property "bodyKitDate")}}
                                        <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                            <label>{{$field.Label}}</label>
                                            {{if eq $field.Type "select"}}
                                            <select name="{{$field.Property}}">
                                                <option value="">Select {{$field.Label}}</option>
                                                {{range $option := $field.Options}}
                                                <option value="{{$option.Value}}">{{$option.Label}}</option>
                                                {{end}}
                                            </select>
                                            {{else if eq $field.Type "date"}}
                                            <input type="date" name="{{$field.Property}}">
                                            {{else}}
                                            <input type="text" name="{{$field.Property}}" placeholder="{{$field.Label}}">
                                            {{end}}
                                            {{if $field.HelpText}}
                                            <div class="help-text">{{$field.HelpText}}</div>
                                            {{end}}
                                        </div>
                                        {{end}}
                                        {{end}}
                                    </div>
                                </div>
                                
                                <!-- Audio System Form -->
                                <div id="audioForm" class="modification-form" style="display: none;">
                                    <h4>🔧 <span data-en="Audio System Details" data-de="Audiosystem Details">Audio System Details</span></h4>
                                    <div class="form-grid">
                                        {{range $field := $category.Fields}}
                                        {{if or (eq $field.Property "audioBrand") (eq $field.Property "audioValue") (eq $field.Property "audioDate")}}
                                        <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                            <label>{{$field.Label}}</label>
                                            {{if eq $field.Type "select"}}
                                            <select name="{{$field.Property}}">
                                                <option value="">Select {{$field.Label}}</option>
                                                {{range $option := $field.Options}}
                                                <option value="{{$option.Value}}">{{$option.Label}}</option>
                                                {{end}}
                                            </select>
                                            {{else if eq $field.Type "date"}}
                                            <input type="date" name="{{$field.Property}}">
                                            {{else}}
                                            <input type="text" name="{{$field.Property}}" placeholder="{{$field.Label}}">
                                            {{end}}
                                            {{if $field.HelpText}}
                                            <div class="help-text">{{$field.HelpText}}</div>
                                            {{end}}
                                        </div>
                                        {{end}}
                                        {{end}}
                                    </div>
                                </div>
                                
                                <!-- Performance Chip Form -->
                                <div id="chipForm" class="modification-form" style="display: none;">
                                    <h4>🔧 <span data-en="Performance Chip Details" data-de="Performance-Chip Details">Performance Chip Details</span></h4>
                                    <div class="form-grid">
                                        {{range $field := $category.Fields}}
                                        {{if or (eq $field.Property "chipBrand") (eq $field.Property "chipPower") (eq $field.Property "chipDate")}}
                                        <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                            <label>{{$field.Label}}</label>
                                            {{if eq $field.Type "select"}}
                                            <select name="{{$field.Property}}">
                                                <option value="">Select {{$field.Label}}</option>
                                                {{range $option := $field.Options}}
                                                <option value="{{$option.Value}}">{{$option.Label}}</option>
                                                {{end}}
                                            </select>
                                            {{else if eq $field.Type "date"}}
                                            <input type="date" name="{{$field.Property}}">
                                            {{else}}
                                            <input type="text" name="{{$field.Property}}" placeholder="{{$field.Label}}">
                                            {{end}}
                                            {{if $field.HelpText}}
                                            <div class="help-text">{{$field.HelpText}}</div>
                                            {{end}}
                                        </div>
                                        {{end}}
                                        {{end}}
                                    </div>
                                </div>
                                
                                <!-- Turbo/Supercharger Form -->
                                <div id="turboForm" class="modification-form" style="display: none;">
                                    <h4>🔧 <span data-en="Turbo/Supercharger Details" data-de="Turbo/Lader Details">Turbo/Supercharger Details</span></h4>
                                    <div class="form-grid">
                                        {{range $field := $category.Fields}}
                                        {{if or (eq $field.Property "turboBrand") (eq $field.Property "turboType") (eq $field.Property "turboDate")}}
                                        <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                            <label>{{$field.Label}}</label>
                                            {{if eq $field.Type "select"}}
                                            <select name="{{$field.Property}}">
                                                <option value="">Select {{$field.Label}}</option>
                                                {{range $option := $field.Options}}
                                                <option value="{{$option.Value}}">{{$option.Label}}</option>
                                                {{end}}
                                            </select>
                                            {{else if eq $field.Type "date"}}
                                            <input type="date" name="{{$field.Property}}">
                                            {{else}}
                                            <input type="text" name="{{$field.Property}}" placeholder="{{$field.Label}}">
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
                        </div>
                    </div>
                    
                    {{else if eq $category.ID "drivers"}}
                    <!-- Driver 1 content -->
                    <div class="driver-content active" id="driver1Content">
                        <div class="form-section">
                            <h3>👤 <span data-en="Driver 1 Details" data-de="Fahrer 1 Details">Driver 1 Details</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver1_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver1_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver1_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver1_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver1_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver1_{{$field.Property}}" {{if $field.Required}}required{{end}}>
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver1_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"
                                              {{if $field.Required}}required{{end}}></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver1_{{$field.Property}}"
                                           {{if $field.Required}}required{{end}}>
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver1_{{$field.Property}}" 
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
                    
                    <!-- Driver 2 content -->
                    <div class="driver-content" id="driver2Content">
                        <div class="form-section">
                            <h3>👤 <span data-en="Driver 2 Details" data-de="Fahrer 2 Details">Driver 2 Details</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver2_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver2_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver2_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver2_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver2_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver2_{{$field.Property}}">
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver2_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver2_{{$field.Property}}">
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver2_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}">
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    <!-- Driver 3 content -->
                    <div class="driver-content" id="driver3Content">
                        <div class="form-section">
                            <h3>👤 <span data-en="Driver 3 Details" data-de="Fahrer 3 Details">Driver 3 Details</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver3_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver3_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver3_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver3_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver3_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver3_{{$field.Property}}">
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver3_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver3_{{$field.Property}}">
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver3_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}">
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    <!-- Driver 4 content -->
                    <div class="driver-content" id="driver4Content">
                        <div class="form-section">
                            <h3>👤 <span data-en="Driver 4 Details" data-de="Fahrer 4 Details">Driver 4 Details</span></h3>
                            <div class="form-grid">
                                {{range $field := $category.Fields}}
                                <div class="form-field {{if eq $field.Type "textarea"}}full-width{{end}}">
                                    <label>
                                        {{$field.Label}}
                                        {{if $field.Required}}<span class="required">*</span>{{end}}
                                    </label>
                                    
                                    {{if eq $field.Type "text"}}
                                    <input type="text" 
                                           name="driver4_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "number"}}
                                    <input type="number" 
                                           name="driver4_{{$field.Property}}"
                                           {{if $field.Min}}min="{{$field.Min}}"{{end}}
                                           {{if $field.Max}}max="{{$field.Max}}"{{end}}>
                                    
                                    {{else if eq $field.Type "date"}}
                                    <input type="date" 
                                           name="driver4_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "email"}}
                                    <input type="email" 
                                           name="driver4_{{$field.Property}}">
                                    
                                    {{else if eq $field.Type "tel"}}
                                    <input type="tel" 
                                           name="driver4_{{$field.Property}}"
                                           {{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>
                                    
                                    {{else if eq $field.Type "select"}}
                                    <select name="driver4_{{$field.Property}}">
                                        <option value="">Select {{$field.Label}}</option>
                                        {{range $option := $field.Options}}
                                        <option value="{{$option.Value}}">{{$option.Label}}</option>
                                        {{end}}
                                    </select>
                                    
                                    {{else if eq $field.Type "textarea"}}
                                    <textarea name="driver4_{{$field.Property}}" 
                                              rows="3" 
                                              placeholder="{{$field.Label}}"></textarea>
                                    
                                    {{else if eq $field.Type "checkbox"}}
                                    <input type="checkbox" 
                                           name="driver4_{{$field.Property}}">
                                    
                                    {{else}}
                                    <input type="text" 
                                           name="driver4_{{$field.Property}}" 
                                           placeholder="{{$field.Label}}">
                                    {{end}}
                                    
                                    {{if $field.HelpText}}
                                    <div class="help-text">{{$field.HelpText}}</div>
                                    {{end}}
                                </div>
                                {{end}}
                            </div>
                        </div>
                    </div>
                    
                    {{else}}
                    <!-- Regular category content for non-driver categories -->
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
                    {{end}}
                </div>
                {{end}}
            </div>
        </div>
    </div>
</div>

    <script>
        let currentLanguage = 'en';
        
        function switchCategory(category) {
            console.log('Switching to category:', category);
            
            // Hide all category contents
            const categoryContents = document.querySelectorAll('.category-content');
            categoryContents.forEach(content => {
                content.classList.remove('active');
            });

            // Show the selected category content
            const selectedContent = document.getElementById(category + 'Content');
            if (selectedContent) {
                selectedContent.classList.add('active');
                console.log('Activated content for:', category);
            } else {
                console.error('Content not found for category:', category);
            }

            // Update tab states
            const tabs = document.querySelectorAll('.category-tab');
            tabs.forEach(tab => {
                tab.classList.remove('active');
            });

            const activeTab = document.querySelector("[onclick=\"switchCategory('" + category + "')\"]");
            if (activeTab) {
                activeTab.classList.add('active');
                console.log('Activated tab for:', category);
            } else {
                console.error('Tab not found for category:', category);
            }
        }
        
        function switchDriver(driverNumber) {
            console.log('Switching to driver:', driverNumber);
            
            // Hide all driver contents
            const driverContents = document.querySelectorAll('.driver-content');
            driverContents.forEach(content => {
                content.classList.remove('active');
            });

            // Show the selected driver content
            const selectedContent = document.getElementById('driver' + driverNumber + 'Content');
            if (selectedContent) {
                selectedContent.classList.add('active');
                console.log('Activated driver content for:', driverNumber);
            } else {
                console.error('Driver content not found for:', driverNumber);
            }

            // Update driver tab states
            const driverTabs = document.querySelectorAll('.driver-tab');
            driverTabs.forEach(tab => {
                tab.classList.remove('active');
            });

            const activeDriverTab = document.querySelector("[onclick=\"switchDriver(" + driverNumber + ")\"]");
            if (activeDriverTab) {
                activeDriverTab.classList.add('active');
                console.log('Activated driver tab for:', driverNumber);
            } else {
                console.error('Driver tab not found for:', driverNumber);
            }
        }

        function switchVehicleTab(tabName) {
            console.log('Switching to vehicle tab:', tabName);
            
            // Hide all vehicle contents
            const vehicleContents = document.querySelectorAll('#vehicleDetailsContent, #vehicleModificationsContent');
            vehicleContents.forEach(content => {
                content.classList.remove('active');
            });

            // Show the selected vehicle content
            const selectedContent = document.getElementById('vehicle' + tabName.charAt(0).toUpperCase() + tabName.slice(1) + 'Content');
            if (selectedContent) {
                selectedContent.classList.add('active');
                console.log('Activated vehicle content for:', tabName);
            } else {
                console.error('Vehicle content not found for:', tabName);
            }

            // Update vehicle tab states
            const vehicleTabs = document.querySelectorAll('.driver-tab');
            vehicleTabs.forEach(tab => {
                tab.classList.remove('active');
            });

            const activeVehicleTab = document.querySelector("[onclick=\"switchVehicleTab('" + tabName + "')\"]");
            if (activeVehicleTab) {
                activeVehicleTab.classList.add('active');
                console.log('Activated vehicle tab for:', tabName);
            } else {
                console.error('Vehicle tab not found for:', tabName);
            }
        }
        
        function showModificationForm(selectedModification) {
            console.log('Selected modification:', selectedModification);
            
            // Hide all modification forms
            const modificationForms = document.querySelectorAll('.modification-form');
            modificationForms.forEach(form => {
                form.style.display = 'none';
            });
            
            // Hide the modification forms container
            const formsContainer = document.getElementById('modificationForms');
            formsContainer.style.display = 'none';
            
            // Show the appropriate form based on selection
            if (selectedModification && selectedModification !== 'None') {
                formsContainer.style.display = 'block';
                
                switch (selectedModification) {
                    case 'Engine Tuning':
                        document.getElementById('engineTuningForm').style.display = 'block';
                        break;
                    case 'Exhaust System':
                        document.getElementById('exhaustForm').style.display = 'block';
                        break;
                    case 'Suspension':
                        document.getElementById('suspensionForm').style.display = 'block';
                        break;
                    case 'Wheels/Tyres':
                        document.getElementById('wheelsForm').style.display = 'block';
                        break;
                    case 'Body Kit':
                        document.getElementById('bodyKitForm').style.display = 'block';
                        break;
                    case 'Audio System':
                        document.getElementById('audioForm').style.display = 'block';
                        break;
                    case 'Performance Chip':
                        document.getElementById('chipForm').style.display = 'block';
                        break;
                    case 'Turbo/Supercharger':
                        document.getElementById('turboForm').style.display = 'block';
                        break;
                    default:
                        console.log('No specific form for:', selectedModification);
                        break;
                }
            }
        }
        
        // Track main driver selection
        let mainDriverSelected = null;
        
        function handleDriverClassificationChange(driverNumber, selectedValue) {
            console.log('Driver', driverNumber, 'classification changed to:', selectedValue);
            
            if (selectedValue === 'MAIN') {
                if (mainDriverSelected && mainDriverSelected !== driverNumber) {
                    const previousMainSelect = document.querySelector('select[name="driver' + mainDriverSelected + '_driverClassification"]');
                    if (previousMainSelect) {
                        previousMainSelect.value = 'NAMED';
                        console.log('Reset driver', mainDriverSelected, 'to NAMED');
                    }
                }
                mainDriverSelected = driverNumber;
                console.log('Main driver set to:', driverNumber);
            } else if (selectedValue === 'NAMED' && mainDriverSelected === driverNumber) {
                mainDriverSelected = null;
                console.log('Main driver cleared');
            }
            
            updateDriverClassificationOptions();
            updateConditionalFieldVisibility(driverNumber, selectedValue);
        }
        
        function updateDriverClassificationOptions() {
            const driverNumbers = [1, 2, 3, 4];
            
            driverNumbers.forEach(function(driverNum) {
                const selectElement = document.querySelector('select[name="driver' + driverNum + '_driverClassification"]');
                if (selectElement) {
                    selectElement.innerHTML = '<option value="">Select Driver Classification</option>';
                    
                    const namedOption = document.createElement('option');
                    namedOption.value = 'NAMED';
                    namedOption.textContent = 'Named Driver';
                    selectElement.appendChild(namedOption);
                    
                    if (!mainDriverSelected || mainDriverSelected === driverNum) {
                        const mainOption = document.createElement('option');
                        mainOption.value = 'MAIN';
                        mainOption.textContent = 'Main Driver';
                        selectElement.appendChild(mainOption);
                    }
                    
                    if (mainDriverSelected === driverNum) {
                        selectElement.value = 'MAIN';
                    } else if (selectElement.getAttribute('data-current-value') === 'NAMED') {
                        selectElement.value = 'NAMED';
                    }
                }
            });
        }
        
        function updateConditionalFieldVisibility(driverNumber, selectedValue) {
            console.log('Updating conditional visibility for driver', driverNumber, 'with value:', selectedValue);
            
            // Get all form fields for this driver
            const driverFields = document.querySelectorAll('[name^="driver' + driverNumber + '_"]');
            
            driverFields.forEach(function(field) {
                const fieldName = field.name;
                const fieldContainer = field.closest('.form-field');
                
                if (!fieldContainer) return;
                
                // Handle relationship to main driver field
                if (fieldName.includes('relationshipToMainDriver')) {
                    if (selectedValue === 'MAIN') {
                        fieldContainer.style.display = 'none';
                        field.value = '';
                        console.log('Hiding relationship field for main driver');
                    } else {
                        fieldContainer.style.display = 'block';
                        console.log('Showing relationship field for named driver');
                    }
                }
                
                // Handle living at same address field
                if (fieldName.includes('livingAtSameAddress')) {
                    if (selectedValue === 'MAIN') {
                        fieldContainer.style.display = 'none';
                        field.value = '';
                        console.log('Hiding living at same address field for main driver');
                    } else {
                        fieldContainer.style.display = 'block';
                        console.log('Showing living at same address field for named driver');
                    }
                }
            });
        }
        
        function switchLanguage(lang) {
            console.log('Switching language to:', lang);
            currentLanguage = lang;
            
            // Update language button states
            const languageBtns = document.querySelectorAll('.language-btn');
            languageBtns.forEach(btn => {
                btn.classList.remove('active');
            });
            
            const activeBtn = document.querySelector("[onclick=\"switchLanguage('" + lang + "')\"]");
            if (activeBtn) {
                activeBtn.classList.add('active');
            }
            
            // Update all translatable elements
            const translatableElements = document.querySelectorAll('[data-en][data-de]');
            translatableElements.forEach(element => {
                if (lang === 'en' && element.getAttribute('data-en')) {
                    element.textContent = element.getAttribute('data-en');
                } else if (lang === 'de' && element.getAttribute('data-de')) {
                    element.textContent = element.getAttribute('data-de');
                }
            });
            
            // Update category tab labels
            const categoryTabs = document.querySelectorAll('.category-tab');
            categoryTabs.forEach(tab => {
                const text = tab.textContent.trim();
                if (lang === 'de') {
                    // Simple German translations
                    if (text.includes('Driver Details')) tab.textContent = '👥 Fahrerdetails';
                    else if (text.includes('Vehicle Information')) tab.textContent = '🚗 Fahrzeuginformationen';
                    else if (text.includes('Coverage Options')) tab.textContent = '🛡️ Deckungsoptionen';
                    else if (text.includes('Claims History')) tab.textContent = '📋 Schadenshistorie';
                    else if (text.includes('Payment Information')) tab.textContent = '💳 Zahlungsinformationen';
                    else if (text.includes('Preferences')) tab.textContent = '⚙️ Einstellungen';
                    else if (text.includes('Summary')) tab.textContent = '📊 Zusammenfassung';
                } else {
                    // English translations
                    if (text.includes('Fahrerdetails')) tab.textContent = '👥 Driver Details';
                    else if (text.includes('Fahrzeuginformationen')) tab.textContent = '🚗 Vehicle Information';
                    else if (text.includes('Deckungsoptionen')) tab.textContent = '🛡️ Coverage Options';
                    else if (text.includes('Schadenshistorie')) tab.textContent = '📋 Claims History';
                    else if (text.includes('Zahlungsinformationen')) tab.textContent = '💳 Payment Information';
                    else if (text.includes('Einstellungen')) tab.textContent = '⚙️ Preferences';
                    else if (text.includes('Zusammenfassung')) tab.textContent = '📊 Summary';
                }
            });
            
            // Update driver tab labels
            const driverTabs = document.querySelectorAll('.driver-tab');
            driverTabs.forEach(tab => {
                const text = tab.textContent.trim();
                if (lang === 'de') {
                    if (text.includes('Driver')) {
                        const driverNum = text.match(/\d+/)[0];
                        tab.textContent = '👤 Fahrer ' + driverNum;
                    }
                } else {
                    if (text.includes('Fahrer')) {
                        const driverNum = text.match(/\d+/)[0];
                        tab.textContent = '👤 Driver ' + driverNum;
                    }
                }
            });
        }

        // Initialize on page load
        document.addEventListener('DOMContentLoaded', function() {
            console.log('Form loaded with categories:', document.querySelectorAll('.category-tab').length);
            console.log('Available categories:', Array.from(document.querySelectorAll('.category-tab')).map(tab => tab.textContent.trim()));
            console.log('Driver tabs available:', document.querySelectorAll('.driver-tab').length);
            
            // Initialize conditional field visibility for all drivers
            initializeConditionalFieldVisibility();
        });

        function initializeConditionalFieldVisibility() {
            // Check all driver classification dropdowns and update visibility
            for (let driverNum = 1; driverNum <= 4; driverNum++) {
                const classificationSelect = document.querySelector('select[name="driver' + driverNum + '_driverClassification"]');
                if (classificationSelect && classificationSelect.value) {
                    updateConditionalFieldVisibility(driverNum, classificationSelect.value);
                }
            }
        }
    </script>

    <!-- Floating Upload Modal -->
    <div class="floating-modal" id="uploadModal">
        <div class="modal-header" id="uploadModalHeader">
            <div class="modal-title">📄 <span data-en="Document Upload" data-de="Dokument-Upload">Document Upload</span></div>
            <div class="modal-controls">
                <button class="modal-btn" onclick="minimizeModal('uploadModal')">−</button>
                <button class="modal-btn" onclick="closeModal('uploadModal')">×</button>
            </div>
        </div>
        <div class="modal-content">
            <div class="upload-area" id="uploadArea">
                <div class="upload-icon">📁</div>
                <div class="upload-text" data-en="Drag & drop files here or click to browse" data-de="Dateien hier ablegen oder klicken zum Durchsuchen">Drag & drop files here or click to browse</div>
                <div class="upload-help" data-en="Supports PDF, JPG, PNG up to 10MB" data-de="Unterstützt PDF, JPG, PNG bis 10MB">Supports PDF, JPG, PNG up to 10MB</div>
                <input type="file" id="fileInput" multiple accept=".pdf,.jpg,.jpeg,.png" style="display: none;">
            </div>
            
            <div class="document-types">
                <h4 data-en="Supported Documents" data-de="Unterstützte Dokumente">Supported Documents</h4>
                <div class="document-type">
                    <div class="document-icon">🛂</div>
                    <div class="document-info">
                        <h4 data-en="Passport" data-de="Reisepass">Passport</h4>
                        <p data-en="Extracts: Name, Date of Birth, Nationality, Passport Number" data-de="Extrahiert: Name, Geburtsdatum, Nationalität, Reisepassnummer">Extracts: Name, Date of Birth, Nationality, Passport Number</p>
                    </div>
                </div>
                <div class="document-type">
                    <div class="document-icon">🚗</div>
                    <div class="document-info">
                        <h4 data-en="Previous Quote" data-de="Vorheriges Angebot">Previous Quote</h4>
                        <p data-en="Extracts: Policy Details, Premium, Coverage, Discounts" data-de="Extrahiert: Policendetails, Prämie, Deckung, Rabatte">Extracts: Policy Details, Premium, Coverage, Discounts</p>
                    </div>
                </div>
                <div class="document-type">
                    <div class="document-icon">📋</div>
                    <div class="document-info">
                        <h4 data-en="Existing Insurance" data-de="Bestehende Versicherung">Existing Insurance</h4>
                        <p data-en="Extracts: Policy Number, Expiry Date, Claims History" data-de="Extrahiert: Policennummer, Ablaufdatum, Schadenshistorie">Extracts: Policy Number, Expiry Date, Claims History</p>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Floating Chat Modal -->
    <div class="floating-modal" id="chatModal">
        <div class="modal-header" id="chatModalHeader">
            <div class="modal-title">💬 <span data-en="Insurance Assistant" data-de="Versicherungsassistent">Insurance Assistant</span></div>
            <div class="modal-controls">
                <button class="modal-btn" onclick="minimizeModal('chatModal')">−</button>
                <button class="modal-btn" onclick="closeModal('chatModal')">×</button>
            </div>
        </div>
        <div class="modal-content">
            <div class="chat-messages" id="chatMessages">
                <div class="chat-message assistant">
                    <span data-en="Hello! I'm your insurance assistant. How can I help you today?" data-de="Hallo! Ich bin Ihr Versicherungsassistent. Wie kann ich Ihnen heute helfen?">Hello! I'm your insurance assistant. How can I help you today?</span>
                </div>
            </div>
            <div class="chat-input-area">
                <input type="text" class="chat-input" id="chatInput" placeholder="Type your message..." data-en-placeholder="Type your message..." data-de-placeholder="Nachricht eingeben...">
                <button class="chat-send" onclick="sendMessage()" data-en="Send" data-de="Senden">Send</button>
            </div>
        </div>
    </div>

    <script>
        // Modal functions
        function closeModal(modalId) {
            document.getElementById(modalId).style.display = 'none';
        }

        function minimizeModal(modalId) {
            const modal = document.getElementById(modalId);
            const content = modal.querySelector('.modal-content');
            if (content.style.display === 'none') {
                content.style.display = 'block';
            } else {
                content.style.display = 'none';
            }
        }

        // Drag and drop functionality
        function makeDraggable(modalId) {
            const modal = document.getElementById(modalId);
            const header = document.getElementById(modalId + 'Header');
            let isDragging = false;
            let currentX;
            let currentY;
            let initialX;
            let initialY;
            let xOffset = 0;
            let yOffset = 0;

            header.addEventListener('mousedown', dragStart);
            document.addEventListener('mousemove', drag);
            document.addEventListener('mouseup', dragEnd);

            function dragStart(e) {
                initialX = e.clientX - xOffset;
                initialY = e.clientY - yOffset;
                if (e.target === header) {
                    isDragging = true;
                }
            }

            function drag(e) {
                if (isDragging) {
                    e.preventDefault();
                    currentX = e.clientX - initialX;
                    currentY = e.clientY - initialY;
                    xOffset = currentX;
                    yOffset = currentY;
                    setTranslate(currentX, currentY, modal);
                }
            }

            function setTranslate(xPos, yPos, el) {
                el.style.transform = "translate3d(" + xPos + "px, " + yPos + "px, 0)";
            }

            function dragEnd(e) {
                initialX = currentX;
                initialY = currentY;
                isDragging = false;
            }
        }

        // Upload functionality
        function setupUpload() {
            const uploadArea = document.getElementById('uploadArea');
            const fileInput = document.getElementById('fileInput');

            uploadArea.addEventListener('click', () => fileInput.click());
            
            uploadArea.addEventListener('dragover', (e) => {
                e.preventDefault();
                uploadArea.classList.add('dragover');
            });
            
            uploadArea.addEventListener('dragleave', () => {
                uploadArea.classList.remove('dragover');
            });
            
            uploadArea.addEventListener('drop', (e) => {
                e.preventDefault();
                uploadArea.classList.remove('dragover');
                const files = e.dataTransfer.files;
                handleFiles(files);
            });
            
            fileInput.addEventListener('change', (e) => {
                handleFiles(e.target.files);
            });
        }

        function handleFiles(files) {
            Array.from(files).forEach(file => {
                console.log('Uploaded:', file.name);
                
                // Add document upload message to chat
                addMessage('user', '📄 Uploaded: ' + file.name);
                
                // Simulate document processing and extraction
                setTimeout(() => {
                    const fileType = getFileType(file.name);
                    const extractedInfo = getExtractedInfo(fileType, file.name);
                    addMessage('assistant', extractedInfo);
                }, 1500);
            });
        }

        function getFileType(filename) {
            const ext = filename.split('.').pop().toLowerCase();
            if (ext === 'pdf') return 'pdf';
            if (['jpg', 'jpeg', 'png'].includes(ext)) return 'image';
            return 'unknown';
        }

        function getExtractedInfo(fileType, filename) {
            const timestamp = new Date().toLocaleTimeString();
            
            if (fileType === 'pdf') {
                const extractedData = {
                    'driver1_firstName': 'John',
                    'driver1_lastName': 'Smith',
                    'driver1_dateOfBirth': '1985-03-15',
                    'driver1_addressLine1': '123 Main Street',
                    'driver1_city': 'London',
                    'driver1_postcode': 'SW1A 1AA',
                    'driver1_licenceNumber': 'SMITH123456',
                    'registrationNumber': 'AB12 CDE'
                };
                
                fillEmptyFormFields(extractedData);
                
                return "📋 **Document Processed**: " + filename + "\n⏰ **Time**: " + timestamp + "\n\n**Extracted Information**:\n• **Name**: John Smith\n• **Date of Birth**: 15/03/1985\n• **Address**: 123 Main Street, London\n• **License Number**: SMITH123456\n• **Vehicle Registration**: AB12 CDE\n\n**Status**: ✅ Successfully processed and data extracted\n\n**Form Fields Updated**: " + Object.keys(extractedData).length + " fields filled";
            } else if (fileType === 'image') {
                const extractedData = {
                    'driver1_firstName': 'John',
                    'driver1_lastName': 'Smith',
                    'driver1_dateOfBirth': '1985-03-15',
                    'driver1_title': 'Mr',
                    'driver1_gender': 'Male'
                };
                
                fillEmptyFormFields(extractedData);
                
                return "🖼️ **Image Processed**: " + filename + "\n⏰ **Time**: " + timestamp + "\n\n**Extracted Information**:\n• **Document Type**: Passport\n• **Name**: John Smith\n• **Title**: Mr\n• **Date of Birth**: 15/03/1985\n• **Gender**: Male\n\n**Status**: ✅ Successfully processed and data extracted\n\n**Form Fields Updated**: " + Object.keys(extractedData).length + " fields filled";
            } else {
                return "❓ **Unknown File Type**: " + filename + "\n⏰ **Time**: " + timestamp + "\n\n**Status**: ⚠️ Unable to process this file type. Please upload PDF or image files.";
            }
        }

        function fillEmptyFormFields(extractedData) {
            let filledCount = 0;
            let skippedCount = 0;
            let filledFields = [];
            let skippedFields = [];
            
            for (const [fieldName, value] of Object.entries(extractedData)) {
                const field = document.querySelector('[name="' + fieldName + '"]');
                if (field) {
                    // Get the field label for better display
                    const label = field.closest('.form-field')?.querySelector('label')?.textContent?.trim() || fieldName;
                    
                    // Only fill if the field is empty
                    if (!field.value || field.value.trim() === '') {
                        field.value = value;
                        field.style.borderColor = '#00ff00';
                        field.style.boxShadow = '0 0 10px rgba(0, 255, 0, 0.3)';
                        
                        // Remove the green highlight after 3 seconds
                        setTimeout(() => {
                            field.style.borderColor = '';
                            field.style.boxShadow = '';
                        }, 3000);
                        
                        filledCount++;
                        filledFields.push(label);
                    } else {
                        skippedCount++;
                        skippedFields.push(label);
                    }
                }
            }
            
            // Add a summary message to chat
            if (filledCount > 0) {
                setTimeout(() => {
                    let message = "🎯 **Form Fields Updated**:\n• **Filled**: " + filledCount + " empty fields\n• **Skipped**: " + skippedCount + " already filled fields\n\n";
                    
                    if (filledFields.length > 0) {
                        message += "**✅ Filled Fields**:\n";
                        filledFields.forEach(field => {
                            message += "• " + field + "\n";
                        });
                    }
                    
                    if (skippedFields.length > 0) {
                        message += "\n**⏭️ Skipped Fields (already filled)**:\n";
                        skippedFields.forEach(field => {
                            message += "• " + field + "\n";
                        });
                    }
                    
                    message += "\n✅ Form has been automatically updated with extracted data!";
                    
                    addMessage('assistant', message);
                }, 500);
            }
        }

        // Chat functionality
        function sendMessage() {
            const input = document.getElementById('chatInput');
            const message = input.value.trim();
            if (message) {
                addMessage('user', message);
                input.value = '';
                // Simulate response
                setTimeout(function() {
                    addMessage('assistant', 'Thank you for your message. I\'ll help you with that.');
                }, 1000);
            }
        }

        function addMessage(type, text) {
            const messages = document.getElementById('chatMessages');
            const messageDiv = document.createElement('div');
            messageDiv.className = "chat-message " + type;
            
            // Handle formatted text with line breaks and bold text
            const formattedText = text
                .replace(/\*\*(.*?)\*\*/g, '<strong>$1</strong>')
                .replace(/\n/g, '<br>');
            
            messageDiv.innerHTML = formattedText;
            messages.appendChild(messageDiv);
            messages.scrollTop = messages.scrollHeight;
        }

        // Initialize modals
        document.addEventListener('DOMContentLoaded', function() {
            makeDraggable('uploadModal');
            makeDraggable('chatModal');
            setupUpload();
            
            // Enter key for chat
            document.getElementById('chatInput').addEventListener('keypress', function(e) {
                if (e.key === 'Enter') {
                    sendMessage();
                }
            });
        });

        // Document processing functionality
        let processedDocuments = [];
        let extractedFields = {};

        function processDocument(file) {
            console.log('Processing document:', file.name);
            
            // Simulate document recognition and field extraction
            const documentType = recognizeDocumentType(file.name);
            const extractedData = simulateFieldExtraction(documentType, file);
            
            // Add to processed documents
            processedDocuments.push({
                name: file.name,
                type: documentType,
                extractedData: extractedData,
                timestamp: new Date().toISOString()
            });
            
            // Check for conflicts and update chat
            checkFieldConflicts(extractedData);
            
            // Update document upload area
            updateDocumentUploadArea();
        }

        function recognizeDocumentType(filename) {
            const lowerName = filename.toLowerCase();
            
            if (lowerName.includes('driving') || lowerName.includes('licence') || lowerName.includes('license')) {
                return 'DrivingLicence';
            } else if (lowerName.includes('passport')) {
                return 'Passport';
            } else if (lowerName.includes('quote') || lowerName.includes('insurance')) {
                return 'InsuranceQuote';
            } else if (lowerName.includes('policy')) {
                return 'InsurancePolicy';
            } else if (lowerName.includes('registration') || lowerName.includes('v5c')) {
                return 'VehicleRegistration';
            } else if (lowerName.includes('bank') || lowerName.includes('statement')) {
                return 'BankStatement';
            } else if (lowerName.includes('utility') || lowerName.includes('bill')) {
                return 'UtilityBill';
            } else if (lowerName.includes('medical') || lowerName.includes('health')) {
                return 'MedicalCertificate';
            } else {
                return 'Unknown';
            }
        }

        function simulateFieldExtraction(documentType, file) {
            // Simulate extracted fields based on document type
            const extractionMap = {
                'DrivingLicence': {
                    title: 'Mr',
                    firstName: 'John',
                    lastName: 'Smith',
                    dateOfBirth: '1985-03-15',
                    addressLine1: '123 Main Street',
                    addressLine2: 'Apartment 4B',
                    city: 'London',
                    postcode: 'SW1A 1AA',
                    licenceNumber: 'SMITH123456789',
                    licenceType: 'Full UK',
                    expiryDate: '2030-12-31'
                },
                'Passport': {
                    title: 'Mr',
                    firstName: 'John',
                    lastName: 'Smith',
                    dateOfBirth: '1985-03-15',
                    nationality: 'British',
                    passportNumber: '123456789',
                    expiryDate: '2030-12-31',
                    placeOfBirth: 'London'
                },
                'InsuranceQuote': {
                    policyType: 'Comprehensive',
                    estimatedPremium: '850.00',
                    policyTerm: '12 months',
                    effectiveDate: '2024-01-01',
                    expirationDate: '2024-12-31',
                    discountsApplied: 'No Claims Bonus',
                    productName: 'Premium Auto Insurance',
                    quoteReference: 'QTE-2024-001'
                },
                'VehicleRegistration': {
                    registrationNumber: 'AB12 CDE',
                    make: 'Ford',
                    model: 'Focus',
                    yearOfManufacture: '2020',
                    engineSize: '1.5',
                    fuelType: 'Petrol',
                    transmission: 'Manual',
                    bodyType: 'Hatchback',
                    colour: 'Blue',
                    vinNumber: '1HGBH41JXMN109186'
                }
            };
            
            return extractionMap[documentType] || {};
        }

        function checkFieldConflicts(extractedData) {
            const conflicts = [];
            
            for (const [fieldName, extractedValue] of Object.entries(extractedData)) {
                // Check if field exists in form and has a value
                const formField = document.querySelector('[name*="' + fieldName + '"]');
                if (formField && formField.value && formField.value !== extractedValue) {
                    conflicts.push({
                        fieldName: fieldName,
                        currentValue: formField.value,
                        extractedValue: extractedValue,
                        fieldElement: formField
                    });
                }
            }
            
            if (conflicts.length > 0) {
                showConflictDialog(conflicts);
            } else {
                // Auto-fill empty fields
                autoFillEmptyFields(extractedData);
            }
        }

        function showConflictDialog(conflicts) {
            const chatMessages = document.getElementById('chatMessages');
            
            conflicts.forEach(conflict => {
                const conflictMessage = document.createElement('div');
                conflictMessage.className = 'chat-message assistant conflict-message';
                conflictMessage.innerHTML = 
                    '<div class="conflict-alert">' +
                        '<strong>⚠️ Field Conflict Detected</strong><br>' +
                        '<strong>Field:</strong> ' + conflict.fieldName + '<br>' +
                        '<strong>Current Value:</strong> ' + conflict.currentValue + '<br>' +
                        '<strong>Extracted Value:</strong> ' + conflict.extractedValue + '<br>' +
                        '<div class="conflict-actions">' +
                            '<button onclick="resolveConflict(\'' + conflict.fieldName + '\', \'' + conflict.extractedValue + '\', true)" class="btn-overwrite">Overwrite</button>' +
                            '<button onclick="resolveConflict(\'' + conflict.fieldName + '\', \'' + conflict.extractedValue + '\', false)" class="btn-keep">Keep Current</button>' +
                        '</div>' +
                    '</div>';
                chatMessages.appendChild(conflictMessage);
            });
            
            // Scroll to bottom
            chatMessages.scrollTop = chatMessages.scrollHeight;
        }

        function resolveConflict(fieldName, extractedValue, overwrite) {
            if (overwrite) {
                const formField = document.querySelector('[name*="' + fieldName + '"]');
                if (formField) {
                    formField.value = extractedValue;
                    formField.style.borderColor = '#00ff41';
                    setTimeout(() => {
                        formField.style.borderColor = '';
                    }, 2000);
                }
            }
            
            // Remove conflict message
            const conflictMessages = document.querySelectorAll('.conflict-message');
            conflictMessages.forEach(msg => {
                if (msg.textContent.includes(fieldName)) {
                    msg.remove();
                }
            });
            
            // Add resolution message
            const chatMessages = document.getElementById('chatMessages');
            const resolutionMessage = document.createElement('div');
            resolutionMessage.className = 'chat-message assistant';
            resolutionMessage.textContent = '✅ Field "' + fieldName + '" ' + (overwrite ? 'updated' : 'kept unchanged') + '.';
            chatMessages.appendChild(resolutionMessage);
            chatMessages.scrollTop = chatMessages.scrollHeight;
        }

        function autoFillEmptyFields(extractedData) {
            let filledCount = 0;
            
            for (const [fieldName, extractedValue] of Object.entries(extractedData)) {
                const formField = document.querySelector('[name*="' + fieldName + '"]');
                if (formField && !formField.value) {
                    formField.value = extractedValue;
                    formField.style.borderColor = '#00ff41';
                    setTimeout(() => {
                        formField.style.borderColor = '';
                    }, 2000);
                    filledCount++;
                }
            }
            
            if (filledCount > 0) {
                const chatMessages = document.getElementById('chatMessages');
                const autoFillMessage = document.createElement('div');
                autoFillMessage.className = 'chat-message assistant';
                autoFillMessage.textContent = '✅ Automatically filled ' + filledCount + ' empty fields from document.';
                chatMessages.appendChild(autoFillMessage);
                chatMessages.scrollTop = chatMessages.scrollHeight;
            }
        }

        function updateDocumentUploadArea() {
            const uploadArea = document.getElementById('uploadArea');
            const documentList = document.createElement('div');
            documentList.className = 'document-list';
            
            processedDocuments.forEach(doc => {
                const docItem = document.createElement('div');
                docItem.className = 'document-item';
                docItem.innerHTML = 
                    '<div class="doc-icon">📄</div>' +
                    '<div class="doc-info">' +
                        '<strong>' + doc.name + '</strong><br>' +
                        '<small>Type: ' + doc.type + '</small><br>' +
                        '<small>Extracted: ' + Object.keys(doc.extractedData).length + ' fields</small>' +
                    '</div>';
                documentList.appendChild(docItem);
            });
            
            // Replace existing document list
            const existingList = uploadArea.querySelector('.document-list');
            if (existingList) {
                existingList.remove();
            }
            uploadArea.appendChild(documentList);
        }

        // Enhanced file upload handling
        document.getElementById('fileInput').addEventListener('change', function(e) {
            const files = e.target.files;
            for (let file of files) {
                processDocument(file);
            }
        });

        // Enhanced drag and drop
        const uploadArea = document.getElementById('uploadArea');
        
        uploadArea.addEventListener('dragover', function(e) {
            e.preventDefault();
            uploadArea.style.borderColor = '#00ff41';
            uploadArea.style.backgroundColor = 'rgba(0, 255, 65, 0.1)';
        });
        
        uploadArea.addEventListener('dragleave', function(e) {
            e.preventDefault();
            uploadArea.style.borderColor = '';
            uploadArea.style.backgroundColor = '';
        });
        
        uploadArea.addEventListener('drop', function(e) {
            e.preventDefault();
            uploadArea.style.borderColor = '';
            uploadArea.style.backgroundColor = '';
            
            const files = e.dataTransfer.files;
            for (let file of files) {
                processDocument(file);
            }
        });
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
		"Categories":         sortedCategories,
		"GermanTranslations": germanTranslations,
	})
	if err != nil {
		return "", fmt.Errorf("failed to execute template: %v", err)
	}

	return buf.String(), nil
}

func (p *TTLParser) addEnumerationOptions() {
	// Define enumeration options for various fields
	enumerationOptions := map[string][]FieldOption{
		"title": {
			{Value: "Mr", Label: "Mr"},
			{Value: "Mrs", Label: "Mrs"},
			{Value: "Miss", Label: "Miss"},
			{Value: "Ms", Label: "Ms"},
			{Value: "Dr", Label: "Dr"},
			{Value: "Prof", Label: "Prof"},
			{Value: "Rev", Label: "Rev"},
		},
		"gender": {
			{Value: "Male", Label: "Male"},
			{Value: "Female", Label: "Female"},
			{Value: "Other", Label: "Other"},
			{Value: "Prefer not to say", Label: "Prefer not to say"},
		},
		"maritalStatus": {
			{Value: "Single", Label: "Single"},
			{Value: "Married", Label: "Married"},
			{Value: "Divorced", Label: "Divorced"},
			{Value: "Widowed", Label: "Widowed"},
			{Value: "Civil Partnership", Label: "Civil Partnership"},
			{Value: "Separated", Label: "Separated"},
		},
		"licenceType": {
			{Value: "Full UK", Label: "Full UK Licence"},
			{Value: "Provisional UK", Label: "Provisional UK Licence"},
			{Value: "EU/EEA", Label: "EU/EEA Licence"},
			{Value: "International", Label: "International Licence"},
		},
		"driverClassification": {
			{Value: "MAIN", Label: "Main Driver"},
			{Value: "NAMED", Label: "Named Driver"},
		},
		"relationshipToMainDriver": {
			{Value: "Spouse/Partner", Label: "Spouse/Partner"},
			{Value: "Child", Label: "Child"},
			{Value: "Parent", Label: "Parent"},
			{Value: "Sibling", Label: "Sibling"},
			{Value: "Friend", Label: "Friend"},
			{Value: "Employee", Label: "Employee"},
			{Value: "Other", Label: "Other"},
		},
		"employmentStatus": {
			{Value: "Employed", Label: "Employed"},
			{Value: "Self-Employed", Label: "Self-Employed"},
			{Value: "Student", Label: "Student"},
			{Value: "Retired", Label: "Retired"},
			{Value: "Unemployed", Label: "Unemployed"},
			{Value: "Homemaker", Label: "Homemaker"},
		},
		"fuelType": {
			{Value: "Petrol", Label: "Petrol"},
			{Value: "Diesel", Label: "Diesel"},
			{Value: "Hybrid", Label: "Hybrid"},
			{Value: "Electric", Label: "Electric"},
			{Value: "LPG", Label: "LPG"},
			{Value: "Biofuel", Label: "Biofuel"},
		},
		"transmission": {
			{Value: "Manual", Label: "Manual"},
			{Value: "Automatic", Label: "Automatic"},
			{Value: "Semi-Auto", Label: "Semi-Automatic"},
		},
		"bodyType": {
			{Value: "Hatchback", Label: "Hatchback"},
			{Value: "Saloon", Label: "Saloon"},
			{Value: "Estate", Label: "Estate"},
			{Value: "SUV", Label: "SUV"},
			{Value: "MPV", Label: "MPV"},
			{Value: "Coupe", Label: "Coupe"},
			{Value: "Convertible", Label: "Convertible"},
			{Value: "Van", Label: "Van"},
			{Value: "Pickup", Label: "Pickup"},
		},
		"purchaseType": {
			{Value: "NEW", Label: "New"},
			{Value: "USED", Label: "Used"},
		},
		"importType": {
			{Value: "None", Label: "Not Imported"},
			{Value: "Personal Import", Label: "Personal Import"},
			{Value: "Parallel Import", Label: "Parallel Import"},
			{Value: "Grey Import", Label: "Grey Import"},
		},
		"alarm": {
			{Value: "None", Label: "No Alarm"},
			{Value: "Standard", Label: "Standard Alarm"},
			{Value: "Thatcham Cat 1", Label: "Thatcham Category 1"},
			{Value: "Thatcham Cat 2", Label: "Thatcham Category 2"},
			{Value: "Thatcham Cat 2-1", Label: "Thatcham Category 2-1"},
		},
		"immobiliser": {
			{Value: "None", Label: "No Immobiliser"},
			{Value: "Standard", Label: "Standard Immobiliser"},
			{Value: "Thatcham Cat 1", Label: "Thatcham Category 1"},
			{Value: "Thatcham Cat 2", Label: "Thatcham Category 2"},
			{Value: "Thatcham Cat 2-1", Label: "Thatcham Category 2-1"},
		},
		"tracker": {
			{Value: "None", Label: "No Tracker"},
			{Value: "Standard", Label: "Standard Tracker"},
			{Value: "Thatcham Cat 5", Label: "Thatcham Category 5"},
			{Value: "Thatcham Cat 6", Label: "Thatcham Category 6"},
			{Value: "Thatcham Cat 7", Label: "Thatcham Category 7"},
		},
		"coverType": {
			{Value: "COMPREHENSIVE", Label: "Comprehensive"},
			{Value: "THIRD_PARTY_FIRE_THEFT", Label: "Third Party Fire & Theft"},
			{Value: "THIRD_PARTY", Label: "Third Party Only"},
		},
		"classOfUse": {
			{Value: "Social Domestic Pleasure", Label: "Social, Domestic & Pleasure"},
			{Value: "Social Domestic Pleasure and Commuting", Label: "Social, Domestic, Pleasure & Commuting"},
			{Value: "Business Use", Label: "Business Use"},
			{Value: "Commercial", Label: "Commercial"},
		},
		"paymentType": {
			{Value: "ANNUAL", Label: "Pay Annually"},
			{Value: "MONTHLY", Label: "Pay Monthly"},
		},
		"paymentMethod": {
			{Value: "Direct Debit", Label: "Direct Debit"},
			{Value: "Credit Card", Label: "Credit Card"},
			{Value: "Debit Card", Label: "Debit Card"},
			{Value: "Bank Transfer", Label: "Bank Transfer"},
		},
		"overnightLocation": {
			{Value: "Garage", Label: "Garage"},
			{Value: "Driveway", Label: "Driveway"},
			{Value: "Street", Label: "Street"},
			{Value: "Car Park", Label: "Car Park"},
			{Value: "Private Land", Label: "Private Land"},
			{Value: "Work Premises", Label: "Work Premises"},
			{Value: "Other", Label: "Other"},
		},
		"daytimeLocation": {
			{Value: "Home", Label: "Home"},
			{Value: "Work", Label: "Work"},
			{Value: "Street", Label: "Street"},
			{Value: "Car Park", Label: "Car Park"},
			{Value: "Shopping Centre", Label: "Shopping Centre"},
			{Value: "Train Station", Label: "Train Station"},
			{Value: "Other", Label: "Other"},
		},
		"overnightPostcode": {
			{Value: "SW1A 1AA", Label: "SW1A 1AA (Westminster)"},
			{Value: "M1 1AA", Label: "M1 1AA (Manchester)"},
			{Value: "B1 1AA", Label: "B1 1AA (Birmingham)"},
			{Value: "L1 1AA", Label: "L1 1AA (Liverpool)"},
			{Value: "EH1 1AA", Label: "EH1 1AA (Edinburgh)"},
			{Value: "CF1 1AA", Label: "CF1 1AA (Cardiff)"},
			{Value: "Other", Label: "Other"},
		},
		"daytimePostcode": {
			{Value: "SW1A 1AA", Label: "SW1A 1AA (Westminster)"},
			{Value: "M1 1AA", Label: "M1 1AA (Manchester)"},
			{Value: "B1 1AA", Label: "B1 1AA (Birmingham)"},
			{Value: "L1 1AA", Label: "L1 1AA (Liverpool)"},
			{Value: "EH1 1AA", Label: "EH1 1AA (Edinburgh)"},
			{Value: "CF1 1AA", Label: "CF1 1AA (Cardiff)"},
			{Value: "Other", Label: "Other"},
		},
		"vehicleModifications": {
			{Value: "None", Label: "No Modifications"},
			{Value: "Engine Tuning", Label: "Engine Tuning"},
			{Value: "Exhaust System", Label: "Exhaust System"},
			{Value: "Air Intake", Label: "Air Intake"},
			{Value: "Suspension", Label: "Suspension"},
			{Value: "Wheels/Tyres", Label: "Wheels/Tyres"},
			{Value: "Body Kit", Label: "Body Kit"},
			{Value: "Interior", Label: "Interior Modifications"},
			{Value: "Audio System", Label: "Audio System"},
			{Value: "Performance Chip", Label: "Performance Chip"},
			{Value: "Turbo/Supercharger", Label: "Turbo/Supercharger"},
			{Value: "Other", Label: "Other Modifications"},
		},
	}

	// Add options to fields
	for fieldName, options := range enumerationOptions {
		if field, exists := p.Fields[fieldName]; exists {
			field.Options = options
			field.Type = "select" // Change type to select for fields with options
		}
	}
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

// Main function removed to avoid conflicts with settings generator
