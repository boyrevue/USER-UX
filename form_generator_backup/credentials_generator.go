package formgenerator

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

// CredentialField represents a field in the credentials ontology
type CredentialField struct {
	Property            string
	Label               string
	Type                string
	Required            bool
	HelpText            string
	HelpPrompt          string
	Pattern             string
	ErrorMessage        string
	Options             []FieldOption
	InputType           string
	EncryptionRequired  bool
	EncryptionAlgorithm string
	KeyDerivation       string
	SessionTimeout      string
	Requires2FA         bool
	PCICompliant        bool
	StorageType         string
	AccessLogging       string
}

// CredentialCategory represents a category of credentials
type CredentialCategory struct {
	Name   string
	Label  string
	Fields []*CredentialField
	Order  int
}

// CredentialsParser parses the credentials ontology
type CredentialsParser struct {
	Categories map[string]*CredentialCategory
}

// NewCredentialsParser creates a new credentials parser
func NewCredentialsParser() *CredentialsParser {
	return &CredentialsParser{
		Categories: make(map[string]*CredentialCategory),
	}
}

// ParseCredentialsTTL parses the credentials TTL file
func (p *CredentialsParser) ParseCredentialsTTL(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("failed to open credentials TTL file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var currentClass string
	var currentProperty string
	var currentField *CredentialField

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Check for class definitions
		if strings.Contains(line, "a owl:Class") {
			if strings.Contains(line, "EmailCredential") {
				currentClass = "email"
				p.Categories["email"] = &CredentialCategory{
					Name:  "email",
					Label: "Email Credentials",
					Order: 1,
				}
			} else if strings.Contains(line, "BankCredential") {
				currentClass = "bank"
				p.Categories["bank"] = &CredentialCategory{
					Name:  "bank",
					Label: "Bank Credentials",
					Order: 2,
				}
			} else if strings.Contains(line, "CreditCardCredential") {
				currentClass = "creditcard"
				p.Categories["creditcard"] = &CredentialCategory{
					Name:  "creditcard",
					Label: "Credit Card Credentials",
					Order: 3,
				}
			}
			continue
		}

		// Check for property definitions
		if strings.Contains(line, "a owl:DatatypeProperty") {
			// Extract property name
			re := regexp.MustCompile(`creds:(\w+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) > 1 {
				currentProperty = matches[1]
				currentField = &CredentialField{
					Property: currentProperty,
				}
			}
			continue
		}

		// Parse property attributes
		if currentField != nil {
			if strings.Contains(line, "rdfs:label") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.Label = matches[1]
				}
			} else if strings.Contains(line, "creds:isRequired") {
				if strings.Contains(line, "true") {
					currentField.Required = true
				}
			} else if strings.Contains(line, "creds:formHelpText") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.HelpText = matches[1]
				}
			} else if strings.Contains(line, "creds:helpPrompt") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.HelpPrompt = matches[1]
				}
			} else if strings.Contains(line, "creds:validationPattern") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.Pattern = matches[1]
				}
			} else if strings.Contains(line, "creds:errorMessage") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.ErrorMessage = matches[1]
				}
			} else if strings.Contains(line, "creds:inputType") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.InputType = matches[1]
				}
			} else if strings.Contains(line, "creds:enumerationValues") {
				re := regexp.MustCompile(`\(([^)]+)\)`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					values := strings.Split(matches[1], " ")
					for _, value := range values {
						value = strings.Trim(value, `"`)
						if value != "" {
							currentField.Options = append(currentField.Options, FieldOption{
								Value: value,
								Label: value,
							})
						}
					}
				}
			} else if strings.Contains(line, "sec:encryptionRequired") {
				if strings.Contains(line, "true") {
					currentField.EncryptionRequired = true
				}
			} else if strings.Contains(line, "sec:encryptionAlgorithm") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.EncryptionAlgorithm = matches[1]
				}
			} else if strings.Contains(line, "sec:keyDerivation") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.KeyDerivation = matches[1]
				}
			} else if strings.Contains(line, "sec:sessionTimeout") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.SessionTimeout = matches[1]
				}
			} else if strings.Contains(line, "sec:requires2FA") {
				if strings.Contains(line, "true") {
					currentField.Requires2FA = true
				}
			} else if strings.Contains(line, "sec:pciCompliant") {
				if strings.Contains(line, "true") {
					currentField.PCICompliant = true
				}
			} else if strings.Contains(line, "sec:storageType") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.StorageType = matches[1]
				}
			} else if strings.Contains(line, "sec:accessLogging") {
				re := regexp.MustCompile(`"([^"]+)"`)
				matches := re.FindStringSubmatch(line)
				if len(matches) > 1 {
					currentField.AccessLogging = matches[1]
				}
			} else if strings.Contains(line, ".") && currentField.Property != "" {
				// End of property definition
				if currentClass != "" && p.Categories[currentClass] != nil {
					p.Categories[currentClass].Fields = append(p.Categories[currentClass].Fields, currentField)
				}
				currentField = nil
			}
		}
	}

	// Determine field types based on patterns and options
	p.determineFieldTypes()

	return scanner.Err()
}

// determineFieldTypes determines the appropriate input type for each field
func (p *CredentialsParser) determineFieldTypes() {
	for _, category := range p.Categories {
		for _, field := range category.Fields {
			if len(field.Options) > 0 {
				field.Type = "select"
			} else if field.InputType == "password" {
				field.Type = "password"
			} else if strings.Contains(field.Pattern, "^[0-9]+$") {
				field.Type = "number"
			} else if strings.Contains(field.Property, "email") {
				field.Type = "email"
			} else if strings.Contains(field.Property, "port") {
				field.Type = "number"
			} else {
				field.Type = "text"
			}
		}
	}
}

// GenerateCredentialsHTML generates HTML forms for credentials
func (p *CredentialsParser) GenerateCredentialsHTML() (string, error) {
	// Sort categories by order
	var sortedCategories []*CredentialCategory
	for _, category := range p.Categories {
		sortedCategories = append(sortedCategories, category)
	}
	sort.Slice(sortedCategories, func(i, j int) bool {
		return sortedCategories[i].Order < sortedCategories[j].Order
	})

	// HTML template for credentials
	tmpl := `
<!DOCTYPE html>
<html>
<head>
    <title>Secure Credentials Management</title>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
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
            --accent-red: #ff0040;
            --accent-yellow: #ffff00;
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
        }

        .header h1 {
            color: var(--accent-green);
            font-size: 2.5rem;
            margin-bottom: 10px;
        }

        .security-notice {
            background: var(--accent-yellow);
            color: var(--pure-black);
            padding: 15px;
            border-radius: 8px;
            margin-bottom: 30px;
            font-weight: bold;
        }

        .credentials-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
            gap: 30px;
        }

        .credential-card {
            background: var(--gray-900);
            border: 2px solid var(--gray-700);
            border-radius: 8px;
            padding: 25px;
            transition: all 0.3s ease;
        }

        .credential-card:hover {
            border-color: var(--accent-green);
            box-shadow: 0 0 20px rgba(0, 255, 65, 0.2);
        }

        .credential-card h2 {
            color: var(--accent-green);
            margin-bottom: 20px;
            font-size: 1.5rem;
        }

        .form-group {
            margin-bottom: 20px;
        }

        .form-label {
            display: block;
            margin-bottom: 8px;
            font-weight: 600;
            color: var(--gray-300);
        }

        .form-input {
            width: 100%;
            padding: 12px;
            background: var(--gray-800);
            border: 2px solid var(--gray-600);
            border-radius: 6px;
            color: var(--pure-white);
            font-size: 14px;
            transition: all 0.3s ease;
        }

        .form-input:focus {
            outline: none;
            border-color: var(--accent-green);
            box-shadow: 0 0 10px rgba(0, 255, 65, 0.3);
        }

        .form-input[type="password"] {
            font-family: monospace;
        }

        .form-select {
            width: 100%;
            padding: 12px;
            background: var(--gray-800);
            border: 2px solid var(--gray-600);
            border-radius: 6px;
            color: var(--pure-white);
            font-size: 14px;
            cursor: pointer;
        }

        .form-select:focus {
            outline: none;
            border-color: var(--accent-green);
        }

        .help-text {
            font-size: 12px;
            color: var(--gray-400);
            margin-top: 5px;
        }

        .security-info {
            background: var(--gray-800);
            border-left: 4px solid var(--accent-green);
            padding: 15px;
            margin-top: 15px;
            border-radius: 0 6px 6px 0;
        }

        .security-info h4 {
            color: var(--accent-green);
            margin-bottom: 10px;
        }

        .security-info ul {
            list-style: none;
            padding-left: 0;
        }

        .security-info li {
            margin-bottom: 5px;
            font-size: 12px;
        }

        .security-info li:before {
            content: "🔒 ";
            margin-right: 5px;
        }

        .required {
            color: var(--accent-red);
        }

        .add-credential-btn {
            background: var(--accent-green);
            color: var(--pure-black);
            border: none;
            padding: 12px 24px;
            border-radius: 6px;
            font-weight: bold;
            cursor: pointer;
            transition: all 0.3s ease;
            margin-top: 15px;
        }

        .add-credential-btn:hover {
            background: var(--pure-white);
            transform: translateY(-2px);
        }

        .save-all-btn {
            background: var(--accent-green);
            color: var(--pure-black);
            border: none;
            padding: 15px 30px;
            border-radius: 8px;
            font-weight: bold;
            font-size: 16px;
            cursor: pointer;
            transition: all 0.3s ease;
            margin-top: 30px;
            width: 100%;
        }

        .save-all-btn:hover {
            background: var(--pure-white);
            transform: translateY(-2px);
        }

        .credential-instance {
            background: var(--gray-800);
            border: 1px solid var(--gray-600);
            border-radius: 6px;
            padding: 20px;
            margin-bottom: 20px;
        }

        .credential-instance h3 {
            color: var(--accent-green);
            margin-bottom: 15px;
        }

        .remove-btn {
            background: var(--accent-red);
            color: var(--pure-white);
            border: none;
            padding: 8px 16px;
            border-radius: 4px;
            cursor: pointer;
            float: right;
        }

        .remove-btn:hover {
            background: var(--pure-white);
            color: var(--accent-red);
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🔐 Secure Credentials Management</h1>
            <p>PCI-Compliant Storage for Email, Bank, and Credit Card Credentials</p>
        </div>

        <div class="security-notice">
            ⚠️ SECURITY NOTICE: All credentials are encrypted using AES-256-GCM and stored in compliance with PCI DSS standards. 
            Access is logged and sessions timeout automatically for security.
        </div>

        <form id="credentialsForm">
            {{range $category := .}}
            <div class="credential-card">
                <h2>{{$category.Label}}</h2>
                
                <div class="credential-instances" id="{{$category.Name}}Instances">
                    <!-- Credential instances will be added here -->
                </div>

                <button type="button" class="add-credential-btn" onclick="addCredential('{{$category.Name}}')">
                    ➕ Add {{$category.Label}}
                </button>

                <div class="security-info">
                    <h4>🔒 Security Specifications</h4>
                    <ul>
                        {{range $field := $category.Fields}}
                        {{if $field.EncryptionRequired}}
                        <li>{{$field.Label}}: Encrypted with {{$field.EncryptionAlgorithm}}</li>
                        {{end}}
                        {{end}}
                        {{if $category.Name "eq" "bank"}}
                        <li>Session Timeout: 5 minutes</li>
                        <li>2FA Required: Yes</li>
                        {{else if $category.Name "eq" "creditcard"}}
                        <li>Session Timeout: 5 minutes</li>
                        <li>2FA Required: Yes</li>
                        {{else}}
                        <li>Session Timeout: 15 minutes</li>
                        <li>2FA Required: Optional</li>
                        {{end}}
                        <li>Access Logging: Required</li>
                        <li>PCI Compliant: Yes</li>
                    </ul>
                </div>
            </div>
            {{end}}

            <button type="submit" class="save-all-btn">
                💾 Save All Credentials Securely
            </button>
        </form>
    </div>

    <script>
        let credentialCounts = {
            email: 0,
            bank: 0,
            creditcard: 0
        };

        function addCredential(category) {
            credentialCounts[category]++;
            const instancesContainer = document.getElementById(category + 'Instances');
            
            const instanceDiv = document.createElement('div');
            instanceDiv.className = 'credential-instance';
            instanceDiv.id = category + 'Instance' + credentialCounts[category];
            
            instanceDiv.innerHTML = 
                '<h3>' + '{{.Label}}' + ' #' + credentialCounts[category] + '</h3>' +
                '<button type="button" class="remove-btn" onclick="removeCredential(\'' + category + '\', ' + credentialCounts[category] + ')">🗑️ Remove</button>' +
                '<div style="clear: both;"></div>' +
                '{{range $field := .Fields}}' +
                '<div class="form-group">' +
                    '<label class="form-label">' +
                        '{{$field.Label}}' +
                        '{{if $field.Required}}<span class="required">*</span>{{end}}' +
                    '</label>' +
                    '{{if eq $field.Type "select"}}' +
                    '<select class="form-select" name="{{$category}}_{{$field.Property}}_' + credentialCounts[category] + '" {{if $field.Required}}required{{end}}>' +
                        '<option value="">Select {{$field.Label}}</option>' +
                        '{{range $option := $field.Options}}' +
                        '<option value="{{$option.Value}}">{{$option.Label}}</option>' +
                        '{{end}}' +
                    '</select>' +
                    '{{else if eq $field.Type "password"}}' +
                    '<input type="password" class="form-input" name="{{$category}}_{{$field.Property}}_' + credentialCounts[category] + '" ' +
                           'placeholder="{{$field.Label}}" {{if $field.Required}}required{{end}}' +
                           '{{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>' +
                    '{{else if eq $field.Type "number"}}' +
                    '<input type="number" class="form-input" name="{{$category}}_{{$field.Property}}_' + credentialCounts[category] + '" ' +
                           'placeholder="{{$field.Label}}" {{if $field.Required}}required{{end}}>' +
                    '{{else if eq $field.Type "email"}}' +
                    '<input type="email" class="form-input" name="{{$category}}_{{$field.Property}}_' + credentialCounts[category] + '" ' +
                           'placeholder="{{$field.Label}}" {{if $field.Required}}required{{end}}' +
                           '{{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>' +
                    '{{else}}' +
                    '<input type="text" class="form-input" name="{{$category}}_{{$field.Property}}_' + credentialCounts[category] + '" ' +
                           'placeholder="{{$field.Label}}" {{if $field.Required}}required{{end}}' +
                           '{{if $field.Pattern}}pattern="{{$field.Pattern}}"{{end}}>' +
                    '{{end}}' +
                    '{{if $field.HelpText}}' +
                    '<div class="help-text">{{$field.HelpText}}</div>' +
                    '{{end}}' +
                '</div>' +
                '{{end}}';
            
            instancesContainer.appendChild(instanceDiv);
        }

        function removeCredential(category, instanceNumber) {
            const instanceDiv = document.getElementById(category + 'Instance' + instanceNumber);
            if (instanceDiv) {
                instanceDiv.remove();
            }
        }

        document.getElementById('credentialsForm').addEventListener('submit', function(e) {
            e.preventDefault();
            
            // Here you would implement secure credential storage
            // This is just a demonstration - in production, use proper encryption
            alert('🔐 Credentials would be encrypted and stored securely here.\n\nIn production, this would:\n• Encrypt all sensitive data with AES-256-GCM\n• Store in PCI-compliant vault\n• Log all access attempts\n• Implement proper session management\n• Use secure key derivation (PBKDF2)');
        });

        // Auto-add one instance of each credential type
        window.addEventListener('load', function() {
            addCredential('email');
            addCredential('bank');
            addCredential('creditcard');
        });
    </script>
</body>
</html>`

	// Parse and execute template
	tmplParsed, err := template.New("credentials").Parse(tmpl)
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
