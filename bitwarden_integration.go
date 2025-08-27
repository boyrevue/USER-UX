package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// Bitwarden Integration for Open Banking & Site Credentials
type BitwardenManager struct {
	CLIPath      string `json:"cliPath"`
	ServerURL    string `json:"serverUrl"`
	IsUnlocked   bool   `json:"isUnlocked"`
	SessionToken string `json:"sessionToken"`
}

// Open Banking Credential Structure
type OpenBankingCredential struct {
	BankName        string            `json:"bankName"`
	ClientID        string            `json:"clientId"`
	ClientSecret    string            `json:"clientSecret"`
	APIKey          string            `json:"apiKey"`
	CertificatePath string            `json:"certificatePath"`
	PrivateKeyPath  string            `json:"privateKeyPath"`
	BaseURL         string            `json:"baseUrl"`
	AuthURL         string            `json:"authUrl"`
	TokenURL        string            `json:"tokenUrl"`
	Scopes          []string          `json:"scopes"`
	RedirectURI     string            `json:"redirectUri"`
	Environment     string            `json:"environment"` // "sandbox", "production"
	ExtraFields     map[string]string `json:"extraFields"`
}

// Site Login Credential Structure
type SiteCredential struct {
	SiteName         string            `json:"siteName"`
	LoginURL         string            `json:"loginUrl"`
	Username         string            `json:"username"`
	Password         string            `json:"password"`
	UsernameSelector string            `json:"usernameSelector"`
	PasswordSelector string            `json:"passwordSelector"`
	SubmitSelector   string            `json:"submitSelector"`
	TwoFactorSecret  string            `json:"twoFactorSecret"`
	ExtraFields      map[string]string `json:"extraFields"`
	LoginSteps       []LoginStep       `json:"loginSteps"`
}

type LoginStep struct {
	Action   string `json:"action"` // "wait", "click", "type", "select"
	Selector string `json:"selector"`
	Value    string `json:"value"`
	Delay    int    `json:"delay"` // milliseconds
	Required bool   `json:"required"`
}

// Bitwarden Item Categories
const (
	CategoryOpenBanking = "open_banking"
	CategorySiteLogin   = "site_login"
	CategoryInsurance   = "insurance_site"
	CategoryComparison  = "comparison_site"
)

// Initialize Bitwarden Manager
func NewBitwardenManager() *BitwardenManager {
	return &BitwardenManager{
		CLIPath:   "bw", // Assumes Bitwarden CLI is in PATH
		ServerURL: "https://vault.bitwarden.com",
	}
}

// Authenticate with Bitwarden
func (bw *BitwardenManager) Login(email, password string) error {
	// Login to Bitwarden
	cmd := exec.Command(bw.CLIPath, "login", email, password, "--raw")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("bitwarden login failed: %v", err)
	}

	bw.SessionToken = strings.TrimSpace(string(output))
	return nil
}

// Unlock Bitwarden vault
func (bw *BitwardenManager) Unlock(password string) error {
	cmd := exec.Command(bw.CLIPath, "unlock", password, "--raw")
	cmd.Env = append(cmd.Env, fmt.Sprintf("BW_SESSION=%s", bw.SessionToken))

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("bitwarden unlock failed: %v", err)
	}

	bw.SessionToken = strings.TrimSpace(string(output))
	bw.IsUnlocked = true
	return nil
}

// Store Open Banking Credentials
func (bw *BitwardenManager) StoreOpenBankingCredentials(bankName string, creds OpenBankingCredential) error {
	if !bw.IsUnlocked {
		return fmt.Errorf("bitwarden vault is locked")
	}

	// Create Bitwarden item for open banking
	item := map[string]interface{}{
		"type":  1, // Login type
		"name":  fmt.Sprintf("Open Banking - %s", bankName),
		"notes": fmt.Sprintf("Open Banking credentials for %s\nEnvironment: %s", bankName, creds.Environment),
		"login": map[string]interface{}{
			"username": creds.ClientID,
			"password": creds.ClientSecret,
			"uris": []map[string]interface{}{
				{
					"uri":   creds.BaseURL,
					"match": 0, // Domain match
				},
			},
		},
		"fields": []map[string]interface{}{
			{
				"name":  "API_Key",
				"value": creds.APIKey,
				"type":  1, // Hidden field
			},
			{
				"name":  "Certificate_Path",
				"value": creds.CertificatePath,
				"type":  0, // Text field
			},
			{
				"name":  "Private_Key_Path",
				"value": creds.PrivateKeyPath,
				"type":  1, // Hidden field
			},
			{
				"name":  "Auth_URL",
				"value": creds.AuthURL,
				"type":  0, // Text field
			},
			{
				"name":  "Token_URL",
				"value": creds.TokenURL,
				"type":  0, // Text field
			},
			{
				"name":  "Scopes",
				"value": strings.Join(creds.Scopes, ","),
				"type":  0, // Text field
			},
			{
				"name":  "Redirect_URI",
				"value": creds.RedirectURI,
				"type":  0, // Text field
			},
			{
				"name":  "Environment",
				"value": creds.Environment,
				"type":  0, // Text field
			},
			{
				"name":  "Category",
				"value": CategoryOpenBanking,
				"type":  0, // Text field
			},
		},
	}

	// Add extra fields
	for key, value := range creds.ExtraFields {
		field := map[string]interface{}{
			"name":  key,
			"value": value,
			"type":  0, // Text field
		}
		item["fields"] = append(item["fields"].([]map[string]interface{}), field)
	}

	// Convert to JSON and create item
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %v", err)
	}

	cmd := exec.Command(bw.CLIPath, "create", "item", string(itemJSON))
	cmd.Env = append(cmd.Env, fmt.Sprintf("BW_SESSION=%s", bw.SessionToken))

	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to create bitwarden item: %v", err)
	}

	return nil
}

// Store Site Login Credentials
func (bw *BitwardenManager) StoreSiteCredentials(siteName string, creds SiteCredential) error {
	if !bw.IsUnlocked {
		return fmt.Errorf("bitwarden vault is locked")
	}

	// Determine category based on site
	category := CategorySiteLogin
	if strings.Contains(strings.ToLower(siteName), "insurance") {
		category = CategoryInsurance
	} else if strings.Contains(strings.ToLower(siteName), "money") ||
		strings.Contains(strings.ToLower(siteName), "compare") ||
		strings.Contains(strings.ToLower(siteName), "market") {
		category = CategoryComparison
	}

	// Create Bitwarden item for site login
	item := map[string]interface{}{
		"type":  1, // Login type
		"name":  fmt.Sprintf("Site Login - %s", siteName),
		"notes": fmt.Sprintf("Login credentials for %s\nAuto-login enabled with stealth browser", siteName),
		"login": map[string]interface{}{
			"username": creds.Username,
			"password": creds.Password,
			"totp":     creds.TwoFactorSecret,
			"uris": []map[string]interface{}{
				{
					"uri":   creds.LoginURL,
					"match": 0, // Domain match
				},
			},
		},
		"fields": []map[string]interface{}{
			{
				"name":  "Username_Selector",
				"value": creds.UsernameSelector,
				"type":  0, // Text field
			},
			{
				"name":  "Password_Selector",
				"value": creds.PasswordSelector,
				"type":  0, // Text field
			},
			{
				"name":  "Submit_Selector",
				"value": creds.SubmitSelector,
				"type":  0, // Text field
			},
			{
				"name":  "Category",
				"value": category,
				"type":  0, // Text field
			},
		},
	}

	// Add extra fields
	for key, value := range creds.ExtraFields {
		field := map[string]interface{}{
			"name":  key,
			"value": value,
			"type":  0, // Text field
		}
		item["fields"] = append(item["fields"].([]map[string]interface{}), field)
	}

	// Add login steps if any
	if len(creds.LoginSteps) > 0 {
		stepsJSON, _ := json.Marshal(creds.LoginSteps)
		field := map[string]interface{}{
			"name":  "Login_Steps",
			"value": string(stepsJSON),
			"type":  1, // Hidden field
		}
		item["fields"] = append(item["fields"].([]map[string]interface{}), field)
	}

	// Convert to JSON and create item
	itemJSON, err := json.Marshal(item)
	if err != nil {
		return fmt.Errorf("failed to marshal item: %v", err)
	}

	cmd := exec.Command(bw.CLIPath, "create", "item", string(itemJSON))
	cmd.Env = append(cmd.Env, fmt.Sprintf("BW_SESSION=%s", bw.SessionToken))

	_, err = cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to create bitwarden item: %v", err)
	}

	return nil
}

// Retrieve Open Banking Credentials
func (bw *BitwardenManager) GetOpenBankingCredentials(bankName string) (*OpenBankingCredential, error) {
	if !bw.IsUnlocked {
		return nil, fmt.Errorf("bitwarden vault is locked")
	}

	// Search for open banking item
	searchTerm := fmt.Sprintf("Open Banking - %s", bankName)
	cmd := exec.Command(bw.CLIPath, "list", "items", "--search", searchTerm)
	cmd.Env = append(cmd.Env, fmt.Sprintf("BW_SESSION=%s", bw.SessionToken))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to search bitwarden items: %v", err)
	}

	var items []BitwardenItem
	err = json.Unmarshal(output, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bitwarden items: %v", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no open banking credentials found for %s", bankName)
	}

	item := items[0]
	creds := &OpenBankingCredential{
		BankName:     bankName,
		ClientID:     item.Login.Username,
		ClientSecret: item.Login.Password,
	}

	// Extract fields
	for _, field := range item.Fields {
		switch field.Name {
		case "API_Key":
			creds.APIKey = field.Value
		case "Certificate_Path":
			creds.CertificatePath = field.Value
		case "Private_Key_Path":
			creds.PrivateKeyPath = field.Value
		case "Auth_URL":
			creds.AuthURL = field.Value
		case "Token_URL":
			creds.TokenURL = field.Value
		case "Scopes":
			creds.Scopes = strings.Split(field.Value, ",")
		case "Redirect_URI":
			creds.RedirectURI = field.Value
		case "Environment":
			creds.Environment = field.Value
		default:
			if creds.ExtraFields == nil {
				creds.ExtraFields = make(map[string]string)
			}
			creds.ExtraFields[field.Name] = field.Value
		}
	}

	// Extract base URL from URIs
	if len(item.URIs) > 0 {
		creds.BaseURL = item.URIs[0].URI
	}

	return creds, nil
}

// Retrieve Site Login Credentials
func (bw *BitwardenManager) GetSiteCredentials(siteName string) (*SiteCredential, error) {
	if !bw.IsUnlocked {
		return nil, fmt.Errorf("bitwarden vault is locked")
	}

	// Search for site login item
	searchTerm := fmt.Sprintf("Site Login - %s", siteName)
	cmd := exec.Command(bw.CLIPath, "list", "items", "--search", searchTerm)
	cmd.Env = append(cmd.Env, fmt.Sprintf("BW_SESSION=%s", bw.SessionToken))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to search bitwarden items: %v", err)
	}

	var items []BitwardenItem
	err = json.Unmarshal(output, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bitwarden items: %v", err)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no site credentials found for %s", siteName)
	}

	item := items[0]
	creds := &SiteCredential{
		SiteName:        siteName,
		Username:        item.Login.Username,
		Password:        item.Login.Password,
		TwoFactorSecret: item.Login.TOTP,
		ExtraFields:     make(map[string]string),
	}

	// Extract login URL from URIs
	if len(item.URIs) > 0 {
		creds.LoginURL = item.URIs[0].URI
	}

	// Extract fields
	for _, field := range item.Fields {
		switch field.Name {
		case "Username_Selector":
			creds.UsernameSelector = field.Value
		case "Password_Selector":
			creds.PasswordSelector = field.Value
		case "Submit_Selector":
			creds.SubmitSelector = field.Value
		case "Login_Steps":
			var steps []LoginStep
			json.Unmarshal([]byte(field.Value), &steps)
			creds.LoginSteps = steps
		case "Category":
			// Skip category field
		default:
			creds.ExtraFields[field.Name] = field.Value
		}
	}

	return creds, nil
}

// List all stored credentials by category
func (bw *BitwardenManager) ListCredentialsByCategory(category string) ([]string, error) {
	if !bw.IsUnlocked {
		return nil, fmt.Errorf("bitwarden vault is locked")
	}

	cmd := exec.Command(bw.CLIPath, "list", "items")
	cmd.Env = append(cmd.Env, fmt.Sprintf("BW_SESSION=%s", bw.SessionToken))

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list bitwarden items: %v", err)
	}

	var items []BitwardenItem
	err = json.Unmarshal(output, &items)
	if err != nil {
		return nil, fmt.Errorf("failed to parse bitwarden items: %v", err)
	}

	var credentials []string
	for _, item := range items {
		// Check if item has the specified category
		for _, field := range item.Fields {
			if field.Name == "Category" && field.Value == category {
				credentials = append(credentials, item.Name)
				break
			}
		}
	}

	return credentials, nil
}

// Sync Bitwarden vault
func (bw *BitwardenManager) Sync() error {
	cmd := exec.Command(bw.CLIPath, "sync")
	cmd.Env = append(cmd.Env, fmt.Sprintf("BW_SESSION=%s", bw.SessionToken))

	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("bitwarden sync failed: %v", err)
	}

	return nil
}

// Lock Bitwarden vault
func (bw *BitwardenManager) Lock() error {
	cmd := exec.Command(bw.CLIPath, "lock")
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("bitwarden lock failed: %v", err)
	}

	bw.IsUnlocked = false
	bw.SessionToken = ""
	return nil
}

// Pre-configured site credentials for common insurance/comparison sites
func (bw *BitwardenManager) SetupCommonSites() error {
	commonSites := map[string]SiteCredential{
		"moneysupermarket.com": {
			SiteName:         "Money Supermarket",
			LoginURL:         "https://www.moneysupermarket.com/my-account/sign-in/",
			UsernameSelector: "input[name='email'], input[id='email']",
			PasswordSelector: "input[name='password'], input[id='password']",
			SubmitSelector:   "button[type='submit'], .btn-primary",
			LoginSteps: []LoginStep{
				{Action: "wait", Selector: "input[name='email']", Delay: 1000, Required: true},
				{Action: "click", Selector: "input[name='email']", Delay: 500, Required: true},
				{Action: "wait", Selector: "input[name='password']", Delay: 500, Required: true},
			},
		},
		"comparethemarket.com": {
			SiteName:         "Compare the Market",
			LoginURL:         "https://www.comparethemarket.com/my-account/sign-in/",
			UsernameSelector: "input[name='email']",
			PasswordSelector: "input[name='password']",
			SubmitSelector:   "button[type='submit']",
		},
		"gocompare.com": {
			SiteName:         "Go Compare",
			LoginURL:         "https://www.gocompare.com/my-account/sign-in/",
			UsernameSelector: "input[name='email']",
			PasswordSelector: "input[name='password']",
			SubmitSelector:   "button[type='submit']",
		},
		"confused.com": {
			SiteName:         "Confused.com",
			LoginURL:         "https://www.confused.com/my-account/sign-in/",
			UsernameSelector: "input[name='email']",
			PasswordSelector: "input[name='password']",
			SubmitSelector:   "button[type='submit']",
		},
	}

	for siteName, creds := range commonSites {
		// Only create template if credentials don't exist
		existing, err := bw.GetSiteCredentials(siteName)
		if err != nil || existing == nil {
			// Store template (user will need to add actual username/password)
			err = bw.StoreSiteCredentials(siteName, creds)
			if err != nil {
				return fmt.Errorf("failed to setup %s: %v", siteName, err)
			}
		}
	}

	return nil
}

// Setup Open Banking templates for major UK banks
func (bw *BitwardenManager) SetupOpenBankingTemplates() error {
	banks := map[string]OpenBankingCredential{
		"Lloyds Bank": {
			BankName:    "Lloyds Bank",
			BaseURL:     "https://api.lloydsbank.com",
			AuthURL:     "https://api.lloydsbank.com/auth/oauth2/authorize",
			TokenURL:    "https://api.lloydsbank.com/auth/oauth2/token",
			Scopes:      []string{"accounts", "payments"},
			Environment: "sandbox",
		},
		"Barclays": {
			BankName:    "Barclays",
			BaseURL:     "https://api.barclays.com",
			AuthURL:     "https://api.barclays.com/auth/oauth2/authorize",
			TokenURL:    "https://api.barclays.com/auth/oauth2/token",
			Scopes:      []string{"accounts", "payments"},
			Environment: "sandbox",
		},
		"HSBC": {
			BankName:    "HSBC",
			BaseURL:     "https://api.hsbc.com",
			AuthURL:     "https://api.hsbc.com/auth/oauth2/authorize",
			TokenURL:    "https://api.hsbc.com/auth/oauth2/token",
			Scopes:      []string{"accounts", "payments"},
			Environment: "sandbox",
		},
		"Santander": {
			BankName:    "Santander",
			BaseURL:     "https://api.santander.co.uk",
			AuthURL:     "https://api.santander.co.uk/auth/oauth2/authorize",
			TokenURL:    "https://api.santander.co.uk/auth/oauth2/token",
			Scopes:      []string{"accounts", "payments"},
			Environment: "sandbox",
		},
	}

	for bankName, creds := range banks {
		// Only create template if credentials don't exist
		existing, err := bw.GetOpenBankingCredentials(bankName)
		if err != nil || existing == nil {
			// Store template (user will need to add actual credentials)
			err = bw.StoreOpenBankingCredentials(bankName, creds)
			if err != nil {
				return fmt.Errorf("failed to setup %s open banking: %v", bankName, err)
			}
		}
	}

	return nil
}
