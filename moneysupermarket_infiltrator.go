package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/chromedp/chromedp"
)

// Money Supermarket Stealth Infiltrator
type MoneySupermarketInfiltrator struct {
	StealthBrowser  *StealthBrowser   `json:"stealthBrowser"`
	TargetURL       string            `json:"targetUrl"`
	FormSelectors   MSFormSelectors   `json:"formSelectors"`
	AntiDetection   MSAntiDetection   `json:"antiDetection"`
	ExtractionRules MSExtractionRules `json:"extractionRules"`
}

type MSFormSelectors struct {
	// Car Insurance Form Selectors
	VehicleReg      string `json:"vehicleReg"`
	PostCode        string `json:"postCode"`
	DateOfBirth     string `json:"dateOfBirth"`
	LicenceType     string `json:"licenceType"`
	YearsNoClaims   string `json:"yearsNoClaims"`
	CoverType       string `json:"coverType"`
	VoluntaryExcess string `json:"voluntaryExcess"`
	AnnualMileage   string `json:"annualMileage"`

	// Navigation Selectors
	GetQuotesButton string `json:"getQuotesButton"`
	ContinueButton  string `json:"continueButton"`
	NextButton      string `json:"nextButton"`

	// Results Selectors
	QuoteResults string `json:"quoteResults"`
	ProviderName string `json:"providerName"`
	QuotePrice   string `json:"quotePrice"`
	CoverDetails string `json:"coverDetails"`
}

type MSAntiDetection struct {
	UserAgentRotation   bool     `json:"userAgentRotation"`
	ViewportVariation   bool     `json:"viewportVariation"`
	TimingRandomization bool     `json:"timingRandomization"`
	MouseJitter         bool     `json:"mouseJitter"`
	ScrollSimulation    bool     `json:"scrollSimulation"`
	CookieManagement    bool     `json:"cookieManagement"`
	ReferrerSpoofing    bool     `json:"referrerSpoofing"`
	CustomHeaders       []string `json:"customHeaders"`
}

type MSExtractionRules struct {
	QuoteProviders []string `json:"quoteProviders"`
	PricePatterns  []string `json:"pricePatterns"`
	DataFields     []string `json:"dataFields"`
	ScreenshotMode bool     `json:"screenshotMode"`
}

// Initialize Money Supermarket Infiltrator
func NewMoneySupermarketInfiltrator() *MoneySupermarketInfiltrator {
	return &MoneySupermarketInfiltrator{
		StealthBrowser: NewStealthBrowser(),
		TargetURL:      "https://www.moneysupermarket.com/car-insurance/",
		FormSelectors: MSFormSelectors{
			VehicleReg:      "input[name='vehicleReg'], input[id*='reg'], input[placeholder*='reg']",
			PostCode:        "input[name='postcode'], input[id*='postcode'], input[placeholder*='postcode']",
			DateOfBirth:     "input[name='dob'], input[id*='birth'], select[name*='birth']",
			LicenceType:     "select[name*='licence'], input[name*='licence']",
			YearsNoClaims:   "select[name*='claims'], input[name*='claims']",
			CoverType:       "select[name*='cover'], input[name*='cover']",
			VoluntaryExcess: "select[name*='excess'], input[name*='excess']",
			AnnualMileage:   "select[name*='mileage'], input[name*='mileage']",

			GetQuotesButton: "button[type='submit'], input[type='submit'], button:contains('Get quotes'), button:contains('Compare')",
			ContinueButton:  "button:contains('Continue'), button:contains('Next'), .continue-btn",
			NextButton:      "button:contains('Next'), .next-btn, .btn-next",

			QuoteResults: ".quote-result, .insurance-quote, .provider-quote",
			ProviderName: ".provider-name, .insurer-name, .company-name",
			QuotePrice:   ".quote-price, .premium, .price",
			CoverDetails: ".cover-details, .policy-details, .coverage",
		},
		AntiDetection: MSAntiDetection{
			UserAgentRotation:   true,
			ViewportVariation:   true,
			TimingRandomization: true,
			MouseJitter:         true,
			ScrollSimulation:    true,
			CookieManagement:    true,
			ReferrerSpoofing:    true,
			CustomHeaders: []string{
				"Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
				"Accept-Language: en-GB,en;q=0.5",
				"Accept-Encoding: gzip, deflate, br",
				"DNT: 1",
				"Connection: keep-alive",
				"Upgrade-Insecure-Requests: 1",
			},
		},
		ExtractionRules: MSExtractionRules{
			QuoteProviders: []string{
				"Admiral", "Aviva", "AXA", "Churchill", "Direct Line",
				"Hastings", "LV=", "More Than", "RAC", "Tesco",
			},
			PricePatterns: []string{
				`£[\d,]+\.?\d*`, `\d+\.\d{2}`, `[\d,]+\s*per\s*year`,
			},
			DataFields: []string{
				"provider", "price", "excess", "cover_type", "rating",
			},
			ScreenshotMode: true,
		},
	}
}

// Execute Money Supermarket Infiltration
func (msi *MoneySupermarketInfiltrator) ExecuteInfiltration(ctx context.Context, quoteData map[string]interface{}) (*InfiltrationResult, error) {
	result := &InfiltrationResult{
		Target:       "Money Supermarket",
		StartTime:    time.Now(),
		StealthLevel: "MAXIMUM",
		Quotes:       []QuoteData{},
		Screenshots:  []string{},
		Errors:       []string{},
		Warnings:     []string{},
	}

	// Phase 1: Stealth Approach
	fmt.Println("🥷 PHASE 1: Stealth Approach - Activating all anti-bot protection...")
	err := msi.activateStealthMode(ctx)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Stealth activation failed: %v", err))
		return result, err
	}

	// Phase 2: Target Navigation
	fmt.Println("🎯 PHASE 2: Target Navigation - Approaching Money Supermarket...")
	err = msi.navigateToTarget(ctx)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Navigation failed: %v", err))
		return result, err
	}

	// Phase 3: Form Infiltration
	fmt.Println("📝 PHASE 3: Form Infiltration - Filling quote form with human-like behavior...")
	err = msi.fillQuoteForm(ctx, quoteData)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Form filling failed: %v", err))
		return result, err
	}

	// Phase 4: Data Extraction
	fmt.Println("💰 PHASE 4: Data Extraction - Harvesting insurance quotes...")
	quotes, err := msi.extractQuotes(ctx)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Quote extraction failed: %v", err))
		return result, err
	}
	result.Quotes = quotes

	// Phase 5: Evidence Collection
	fmt.Println("📸 PHASE 5: Evidence Collection - Taking stealth screenshots...")
	screenshots, err := msi.captureEvidence(ctx)
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Screenshot capture warning: %v", err))
	}
	result.Screenshots = screenshots

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)
	result.Success = len(result.Errors) == 0 && len(result.Quotes) > 0

	if result.Success {
		fmt.Printf("🏆 INFILTRATION SUCCESSFUL! Extracted %d quotes in %v\n", len(result.Quotes), result.Duration)
	}

	return result, nil
}

// Activate maximum stealth mode
func (msi *MoneySupermarketInfiltrator) activateStealthMode(ctx context.Context) error {
	// Apply all Bitwarden-style stealth techniques
	err := msi.StealthBrowser.AdvancedStealth.ApplyBitwardenStealth(ctx)
	if err != nil {
		return fmt.Errorf("advanced stealth activation failed: %v", err)
	}

	// Money Supermarket specific evasion techniques
	stealthScript := `
		// Money Supermarket Anti-Detection Suite
		console.log('🥷 Money Supermarket Stealth Mode Activated');
		
		// 1. Disable automation detection
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined,
		});
		
		// 2. Spoof realistic browser metrics
		Object.defineProperty(navigator, 'hardwareConcurrency', {
			get: () => 8,
		});
		
		// 3. Mock realistic connection
		Object.defineProperty(navigator, 'connection', {
			get: () => ({
				effectiveType: '4g',
				downlink: 10,
				rtt: 50
			}),
		});
		
		// 4. Hide Chrome automation
		window.chrome = {
			runtime: {},
			loadTimes: function() {},
			csi: function() {},
		};
		
		// 5. Realistic timing
		const originalDate = Date;
		Date = class extends originalDate {
			constructor(...args) {
				if (args.length === 0) {
					super(originalDate.now() + Math.random() * 1000 - 500);
				} else {
					super(...args);
				}
			}
			static now() {
				return originalDate.now() + Math.random() * 10 - 5;
			}
		};
		
		console.log('🛡️ Money Supermarket defenses bypassed');
	`

	return chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return chromedp.Evaluate(stealthScript, nil).Do(ctx)
	}))
}

// Navigate to Money Supermarket with stealth
func (msi *MoneySupermarketInfiltrator) navigateToTarget(ctx context.Context) error {
	actions := []chromedp.Action{
		chromedp.Navigate(msi.TargetURL),
		chromedp.WaitReady("body"),
		chromedp.Sleep(2 * time.Second), // Human-like pause
	}

	return chromedp.Run(ctx, actions...)
}

// Fill quote form with human-like behavior
func (msi *MoneySupermarketInfiltrator) fillQuoteForm(ctx context.Context, data map[string]interface{}) error {
	// Extract data with defaults
	vehicleReg := getStringValue(data, "vehicleReg", "AB12 CDE")
	postCode := getStringValue(data, "postCode", "SW1A 1AA")
	dateOfBirth := getStringValue(data, "dateOfBirth", "01/01/1990")

	// Fill form fields with human-like typing
	fields := []struct {
		selector string
		value    string
		delay    time.Duration
	}{
		{msi.FormSelectors.VehicleReg, vehicleReg, 2 * time.Second},
		{msi.FormSelectors.PostCode, postCode, 1 * time.Second},
		{msi.FormSelectors.DateOfBirth, dateOfBirth, 1 * time.Second},
	}

	for _, field := range fields {
		// Wait for field to be visible
		err := chromedp.Run(ctx, chromedp.WaitVisible(field.selector))
		if err != nil {
			continue // Skip if field not found
		}

		// Human-like typing
		err = msi.StealthBrowser.AdvancedStealth.HumanTypeText(ctx, field.selector, field.value)
		if err != nil {
			return fmt.Errorf("failed to fill field %s: %v", field.selector, err)
		}

		// Human-like pause between fields
		time.Sleep(field.delay)
	}

	// Submit form with human-like click
	err := msi.StealthBrowser.AdvancedStealth.HumanClick(ctx, msi.FormSelectors.GetQuotesButton)
	if err != nil {
		return fmt.Errorf("failed to submit form: %v", err)
	}

	// Wait for results to load
	return chromedp.Run(ctx, chromedp.WaitVisible(msi.FormSelectors.QuoteResults, chromedp.ByQuery))
}

// Extract insurance quotes from results
func (msi *MoneySupermarketInfiltrator) extractQuotes(ctx context.Context) ([]QuoteData, error) {
	var quotes []QuoteData

	// Wait for results to load
	time.Sleep(5 * time.Second)

	// Extract quote data using JavaScript
	extractScript := `
		const quotes = [];
		const quoteElements = document.querySelectorAll('` + msi.FormSelectors.QuoteResults + `');
		
		quoteElements.forEach((element, index) => {
			const providerEl = element.querySelector('` + msi.FormSelectors.ProviderName + `');
			const priceEl = element.querySelector('` + msi.FormSelectors.QuotePrice + `');
			const detailsEl = element.querySelector('` + msi.FormSelectors.CoverDetails + `');
			
			if (providerEl && priceEl) {
				quotes.push({
					provider: providerEl.textContent.trim(),
					price: priceEl.textContent.trim(),
					details: detailsEl ? detailsEl.textContent.trim() : '',
					position: index + 1,
					extracted_at: new Date().toISOString()
				});
			}
		});
		
		return quotes;
	`

	var result []map[string]interface{}
	err := chromedp.Run(ctx, chromedp.Evaluate(extractScript, &result))
	if err != nil {
		return quotes, fmt.Errorf("quote extraction failed: %v", err)
	}

	// Convert to QuoteData structs
	for _, item := range result {
		quote := QuoteData{
			Provider:    fmt.Sprintf("%v", item["provider"]),
			Price:       fmt.Sprintf("%v", item["price"]),
			Details:     fmt.Sprintf("%v", item["details"]),
			Position:    int(item["position"].(float64)),
			ExtractedAt: time.Now(),
		}
		quotes = append(quotes, quote)
	}

	return quotes, nil
}

// Capture evidence screenshots
func (msi *MoneySupermarketInfiltrator) captureEvidence(ctx context.Context) ([]string, error) {
	var screenshots []string

	// Take screenshot of results page
	var screenshot []byte
	err := chromedp.Run(ctx, chromedp.CaptureScreenshot(&screenshot))
	if err != nil {
		return screenshots, err
	}

	filename := fmt.Sprintf("moneysupermarket_infiltration_%d.png", time.Now().Unix())
	screenshots = append(screenshots, filename)

	// Save screenshot (in real implementation)
	fmt.Printf("📸 Screenshot captured: %s\n", filename)

	return screenshots, nil
}

// Result structures
type InfiltrationResult struct {
	Target       string        `json:"target"`
	Success      bool          `json:"success"`
	StartTime    time.Time     `json:"startTime"`
	EndTime      time.Time     `json:"endTime"`
	Duration     time.Duration `json:"duration"`
	StealthLevel string        `json:"stealthLevel"`
	Quotes       []QuoteData   `json:"quotes"`
	Screenshots  []string      `json:"screenshots"`
	Errors       []string      `json:"errors"`
	Warnings     []string      `json:"warnings"`
}

type QuoteData struct {
	Provider    string    `json:"provider"`
	Price       string    `json:"price"`
	Details     string    `json:"details"`
	Position    int       `json:"position"`
	ExtractedAt time.Time `json:"extractedAt"`
}

// API Handler for Money Supermarket Infiltration
func MoneySupermarketInfiltrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		VehicleReg  string `json:"vehicleReg"`
		PostCode    string `json:"postCode"`
		DateOfBirth string `json:"dateOfBirth"`
		LicenceType string `json:"licenceType"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Create infiltrator
	infiltrator := NewMoneySupermarketInfiltrator()

	// Create Chrome context with stealth options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // Visible for demonstration
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Execute infiltration
	quoteData := map[string]interface{}{
		"vehicleReg":  request.VehicleReg,
		"postCode":    request.PostCode,
		"dateOfBirth": request.DateOfBirth,
		"licenceType": request.LicenceType,
	}

	result, err := infiltrator.ExecuteInfiltration(ctx, quoteData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Utility function
func getStringValue(data map[string]interface{}, key, defaultValue string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok && str != "" {
			return str
		}
	}
	return defaultValue
}
