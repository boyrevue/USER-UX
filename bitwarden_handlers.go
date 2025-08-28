package main

import (
	"encoding/json"
	"net/http"
	"os/exec"
)

var globalBitwardenManager *BitwardenManager

func init() {
	globalBitwardenManager = NewBitwardenManager()
}

// Bitwarden Status Handler
func BitwardenStatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check Bitwarden CLI status
	cmd := exec.Command(globalBitwardenManager.CLIPath, "status")
	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	var status map[string]interface{}
	if err != nil {
		status = map[string]interface{}{
			"available": false,
			"error":     "Bitwarden CLI error: " + err.Error(),
			"message":   outputStr,
		}
	} else {
		// Try to parse as JSON first
		var statusData map[string]interface{}
		jsonErr := json.Unmarshal(output, &statusData)

		if jsonErr != nil {
			// If not JSON, treat as plain text (likely an error message)
			status = map[string]interface{}{
				"available": false,
				"error":     "Bitwarden CLI returned non-JSON response",
				"message":   outputStr,
				"rawOutput": outputStr,
			}
		} else {
			status = map[string]interface{}{
				"available": true,
				"status":    statusData,
				"message":   "Bitwarden CLI available",
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}

// Bitwarden Login Handler
func BitwardenLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = globalBitwardenManager.Login(request.Email, request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Successfully logged in to Bitwarden",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Bitwarden API Key Login Handler
func BitwardenAPIKeyLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		ClientId     string `json:"clientId"`
		ClientSecret string `json:"clientSecret"`
		Password     string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = globalBitwardenManager.LoginWithAPIKey(request.ClientId, request.ClientSecret, request.Password)
	if err != nil {
		response := map[string]interface{}{
			"success": false,
			"error":   err.Error(),
			"message": "Bitwarden API key login failed",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Successfully logged in to Bitwarden with API key",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Bitwarden Unlock Handler
func BitwardenUnlockHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Password string `json:"password"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = globalBitwardenManager.Unlock(request.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	response := map[string]interface{}{
		"success":    true,
		"message":    "Successfully unlocked Bitwarden vault",
		"isUnlocked": globalBitwardenManager.IsUnlocked,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Store Site Credentials Handler
func BitwardenStoreSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		SiteName    string         `json:"siteName"`
		Credentials SiteCredential `json:"credentials"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = globalBitwardenManager.StoreSiteCredentials(request.SiteName, request.Credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Site credentials stored successfully in Bitwarden",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Store Open Banking Credentials Handler
func BitwardenStoreBankingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		BankName    string                `json:"bankName"`
		Credentials OpenBankingCredential `json:"credentials"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = globalBitwardenManager.StoreOpenBankingCredentials(request.BankName, request.Credentials)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Open Banking credentials stored successfully in Bitwarden",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Get Site Credentials Handler
func BitwardenGetSiteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	siteName := r.URL.Query().Get("site")
	if siteName == "" {
		http.Error(w, "Site name is required", http.StatusBadRequest)
		return
	}

	credentials, err := globalBitwardenManager.GetSiteCredentials(siteName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Don't return the actual password for security
	safeCredentials := map[string]interface{}{
		"siteName":         credentials.SiteName,
		"loginUrl":         credentials.LoginURL,
		"username":         credentials.Username,
		"hasPassword":      credentials.Password != "",
		"usernameSelector": credentials.UsernameSelector,
		"passwordSelector": credentials.PasswordSelector,
		"submitSelector":   credentials.SubmitSelector,
		"hasTwoFactor":     credentials.TwoFactorSecret != "",
		"extraFields":      credentials.ExtraFields,
		"loginSteps":       credentials.LoginSteps,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(safeCredentials)
}

// Get Open Banking Credentials Handler
func BitwardenGetBankingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	bankName := r.URL.Query().Get("bank")
	if bankName == "" {
		http.Error(w, "Bank name is required", http.StatusBadRequest)
		return
	}

	credentials, err := globalBitwardenManager.GetOpenBankingCredentials(bankName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Don't return sensitive credentials for security
	safeCredentials := map[string]interface{}{
		"bankName":        credentials.BankName,
		"baseUrl":         credentials.BaseURL,
		"authUrl":         credentials.AuthURL,
		"tokenUrl":        credentials.TokenURL,
		"hasClientId":     credentials.ClientID != "",
		"hasClientSecret": credentials.ClientSecret != "",
		"hasApiKey":       credentials.APIKey != "",
		"certificatePath": credentials.CertificatePath,
		"privateKeyPath":  credentials.PrivateKeyPath,
		"scopes":          credentials.Scopes,
		"redirectUri":     credentials.RedirectURI,
		"environment":     credentials.Environment,
		"extraFields":     credentials.ExtraFields,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(safeCredentials)
}

// List Credentials by Category Handler
func BitwardenListCredentialsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	category := r.URL.Query().Get("category")
	if category == "" {
		category = CategorySiteLogin // Default category
	}

	credentials, err := globalBitwardenManager.ListCredentialsByCategory(category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success":     true,
		"category":    category,
		"credentials": credentials,
		"count":       len(credentials),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Sync Bitwarden Vault Handler
func BitwardenSyncHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := globalBitwardenManager.Sync()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]interface{}{
		"success": true,
		"message": "Bitwarden vault synced successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Setup Templates Handler
func BitwardenSetupTemplatesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		SetupSites       bool `json:"setupSites"`
		SetupOpenBanking bool `json:"setupOpenBanking"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	var messages []string

	if request.SetupSites {
		err = globalBitwardenManager.SetupCommonSites()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, "Common site templates created")
	}

	if request.SetupOpenBanking {
		err = globalBitwardenManager.SetupOpenBankingTemplates()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		messages = append(messages, "Open Banking templates created")
	}

	response := map[string]interface{}{
		"success":  true,
		"messages": messages,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
