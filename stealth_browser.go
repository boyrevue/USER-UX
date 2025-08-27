package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"strings"
	"time"

	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/cdproto/network"

	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// Stealth Browser Engine
type StealthBrowser struct {
	UserAgents      []string            `json:"userAgents"`
	Proxies         []string            `json:"proxies"`
	ViewportSizes   []ViewportSize      `json:"viewportSizes"`
	BitwardenClient *BitwardenManager   `json:"bitwardenClient"`
	AntiDetection   AntiDetectionConfig `json:"antiDetection"`
	AdvancedStealth *AdvancedStealth    `json:"advancedStealth"`
}

// Bitwarden Integration
type BitwardenClient struct {
	ServerURL    string `json:"serverUrl"`
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
	SessionToken string `json:"sessionToken"`
	IsUnlocked   bool   `json:"isUnlocked"`
}

type BitwardenItem struct {
	ID     string           `json:"id"`
	Name   string           `json:"name"`
	Login  BitwardenLogin   `json:"login"`
	Fields []BitwardenField `json:"fields"`
	URIs   []BitwardenURI   `json:"uris"`
}

type BitwardenLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
	TOTP     string `json:"totp"`
}

type BitwardenField struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  int    `json:"type"` // 0=text, 1=hidden, 2=boolean
}

type BitwardenURI struct {
	URI   string `json:"uri"`
	Match int    `json:"match"` // 0=domain, 1=host, 2=startsWith, 3=exact, 4=regex, 5=never
}

type ViewportSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

type LoginCredential struct {
	Username         string            `json:"username"`
	Password         string            `json:"password"`
	LoginURL         string            `json:"loginUrl"`
	UsernameSelector string            `json:"usernameSelector"`
	PasswordSelector string            `json:"passwordSelector"`
	SubmitSelector   string            `json:"submitSelector"`
	ExtraFields      map[string]string `json:"extraFields"`
	TwoFactorEnabled bool              `json:"twoFactorEnabled"`
}

type AntiDetectionConfig struct {
	RandomizeUserAgent bool          `json:"randomizeUserAgent"`
	RandomizeViewport  bool          `json:"randomizeViewport"`
	HumanizeTyping     bool          `json:"humanizeTyping"`
	HumanizeClicks     bool          `json:"humanizeClicks"`
	RandomDelays       bool          `json:"randomDelays"`
	StealthPlugins     bool          `json:"stealthPlugins"`
	WebRTCBlock        bool          `json:"webrtcBlock"`
	CanvasFingerprint  bool          `json:"canvasFingerprint"`
	AudioFingerprint   bool          `json:"audioFingerprint"`
	MinDelay           time.Duration `json:"minDelay"`
	MaxDelay           time.Duration `json:"maxDelay"`
}

type BrowserSession struct {
	ID             string            `json:"id"`
	URL            string            `json:"url"`
	Title          string            `json:"title"`
	Screenshot     string            `json:"screenshot"`
	Cookies        []network.Cookie  `json:"cookies"`
	LocalStorage   map[string]string `json:"localStorage"`
	SessionStorage map[string]string `json:"sessionStorage"`
	IsLoggedIn     bool              `json:"isLoggedIn"`
	LoginSite      string            `json:"loginSite"`
	CreatedAt      time.Time         `json:"createdAt"`
	LastActivity   time.Time         `json:"lastActivity"`
}

type BrowserAction struct {
	Type    string                 `json:"type"` // "navigate", "click", "type", "scroll", "wait", "login", "extract"
	Target  string                 `json:"target"`
	Value   string                 `json:"value"`
	Options map[string]interface{} `json:"options"`
	Delay   time.Duration          `json:"delay"`
}

type BrowserResult struct {
	Success       bool                   `json:"success"`
	SessionID     string                 `json:"sessionId"`
	URL           string                 `json:"url"`
	Title         string                 `json:"title"`
	Screenshot    string                 `json:"screenshot"`
	IsLoggedIn    bool                   `json:"isLoggedIn"`
	ExtractedData map[string]interface{} `json:"extractedData"`
	Cookies       []*network.Cookie      `json:"cookies"`
	Errors        []string               `json:"errors"`
	Warnings      []string               `json:"warnings"`
	ProcessedAt   time.Time              `json:"processedAt"`
}

// Stealth Browser Handler
func StealthBrowserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		Action    string              `json:"action"` // "create", "navigate", "interact", "login", "extract"
		SessionID string              `json:"sessionId,omitempty"`
		URL       string              `json:"url,omitempty"`
		Actions   []BrowserAction     `json:"actions,omitempty"`
		LoginSite string              `json:"loginSite,omitempty"`
		Selectors []Selector          `json:"selectors,omitempty"`
		Options   AntiDetectionConfig `json:"options"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	browser := NewStealthBrowser()

	var result *BrowserResult

	switch request.Action {
	case "create":
		result, err = browser.CreateSession(request.URL, request.Options)
	case "navigate":
		result, err = browser.Navigate(request.SessionID, request.URL)
	case "interact":
		result, err = browser.ExecuteActions(request.SessionID, request.Actions)
	case "login":
		result, err = browser.Login(request.SessionID, request.LoginSite)
	case "extract":
		result, err = browser.ExtractData(request.SessionID, request.Selectors)
	default:
		http.Error(w, "Invalid action", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, fmt.Sprintf("Browser action failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Initialize Stealth Browser
func NewStealthBrowser() *StealthBrowser {
	return &StealthBrowser{
		UserAgents: []string{
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
			"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
			"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0",
		},
		ViewportSizes: []ViewportSize{
			{Width: 1920, Height: 1080},
			{Width: 1366, Height: 768},
			{Width: 1440, Height: 900},
			{Width: 1536, Height: 864},
			{Width: 1280, Height: 720},
		},
		BitwardenClient: NewBitwardenManager(),
		AdvancedStealth: NewAdvancedStealth(),
		AntiDetection: AntiDetectionConfig{
			RandomizeUserAgent: true,
			RandomizeViewport:  true,
			HumanizeTyping:     true,
			HumanizeClicks:     true,
			RandomDelays:       true,
			StealthPlugins:     true,
			WebRTCBlock:        true,
			CanvasFingerprint:  true,
			AudioFingerprint:   true,
			MinDelay:           500 * time.Millisecond,
			MaxDelay:           2000 * time.Millisecond,
		},
	}
}

// Create new stealth browser session
func (sb *StealthBrowser) CreateSession(url string, config AntiDetectionConfig) (*BrowserResult, error) {
	// Generate random session ID
	sessionID := fmt.Sprintf("session_%d", time.Now().Unix())

	// Create stealth Chrome context
	opts := sb.buildStealthOptions(config)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	result := &BrowserResult{
		SessionID:     sessionID,
		ExtractedData: make(map[string]interface{}),
		Errors:        []string{},
		Warnings:      []string{},
		ProcessedAt:   time.Now(),
	}

	// Apply stealth techniques
	err := sb.applyStealthTechniques(ctx, config)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Stealth setup failed: %v", err))
	}

	// Navigate to URL
	if url != "" {
		err = chromedp.Run(ctx, chromedp.Navigate(url))
		if err != nil {
			result.Success = false
			result.Errors = append(result.Errors, fmt.Sprintf("Navigation failed: %v", err))
			return result, err
		}

		// Wait for page load with human-like delay
		if config.RandomDelays {
			delay := sb.randomDelay(config.MinDelay, config.MaxDelay)
			time.Sleep(delay)
		}

		// Get page info
		err = chromedp.Run(ctx,
			chromedp.Location(&result.URL),
			chromedp.Title(&result.Title),
		)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Could not get page info: %v", err))
		}

		// Take screenshot
		var screenshot []byte
		err = chromedp.Run(ctx, chromedp.Screenshot("body", &screenshot))
		if err == nil {
			// In production, save screenshot and return path
			result.Screenshot = fmt.Sprintf("screenshot_%s.png", sessionID)
		}

		// Get cookies
		cookies, err := network.GetCookies().Do(ctx)
		if err == nil {
			result.Cookies = cookies
		}
	}

	result.Success = len(result.Errors) == 0
	return result, nil
}

// Build stealth Chrome options
func (sb *StealthBrowser) buildStealthOptions(config AntiDetectionConfig) []chromedp.ExecAllocatorOption {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		// Basic stealth options
		chromedp.Flag("headless", false), // Run in visible mode for better stealth
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),

		// Anti-detection flags
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("exclude-switches", "enable-automation"),
		chromedp.Flag("disable-extensions-except", ""),
		chromedp.Flag("disable-extensions", false),
		chromedp.Flag("disable-plugins-discovery", true),
		chromedp.Flag("disable-default-apps", true),

		// Fingerprinting protection
		chromedp.Flag("disable-features", "VizDisplayCompositor,TranslateUI"),
		chromedp.Flag("disable-ipc-flooding-protection", true),
		chromedp.Flag("disable-renderer-backgrounding", true),
		chromedp.Flag("disable-backgrounding-occluded-windows", true),
		chromedp.Flag("disable-client-side-phishing-detection", true),
		chromedp.Flag("disable-sync", true),
		chromedp.Flag("metrics-recording-only", true),
		chromedp.Flag("no-report-upload", true),

		// WebRTC protection
		chromedp.Flag("force-webrtc-ip-handling-policy", "disable_non_proxied_udp"),

		// Language and locale
		chromedp.Flag("lang", "en-US,en"),
		chromedp.Flag("accept-lang", "en-US,en;q=0.9"),
	)

	// Randomize user agent
	if config.RandomizeUserAgent {
		userAgent := sb.getRandomUserAgent()
		opts = append(opts, chromedp.UserAgent(userAgent))
	}

	return opts
}

// Apply advanced stealth techniques
func (sb *StealthBrowser) applyStealthTechniques(ctx context.Context, config AntiDetectionConfig) error {
	var actions []chromedp.Action

	// Set random viewport
	if config.RandomizeViewport {
		viewport := sb.getRandomViewport()
		actions = append(actions, emulation.SetDeviceMetricsOverride(
			int64(viewport.Width), int64(viewport.Height), 1.0, false,
		))
	}

	// Inject stealth scripts
	if config.StealthPlugins {
		stealthScript := sb.getStealthScript()
		actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, err := runtime.Evaluate(stealthScript).Do(ctx)
			return err
		}))
	}

	// Override canvas fingerprinting
	if config.CanvasFingerprint {
		canvasScript := `
			const getContext = HTMLCanvasElement.prototype.getContext;
			HTMLCanvasElement.prototype.getContext = function(contextType, contextAttributes) {
				if (contextType === '2d') {
					const context = getContext.call(this, contextType, contextAttributes);
					const getImageData = context.getImageData;
					context.getImageData = function(sx, sy, sw, sh) {
						const imageData = getImageData.call(this, sx, sy, sw, sh);
						for (let i = 0; i < imageData.data.length; i += 4) {
							imageData.data[i] += Math.floor(Math.random() * 10) - 5;
							imageData.data[i + 1] += Math.floor(Math.random() * 10) - 5;
							imageData.data[i + 2] += Math.floor(Math.random() * 10) - 5;
						}
						return imageData;
					};
					return context;
				}
				return getContext.call(this, contextType, contextAttributes);
			};
		`
		actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, err := runtime.Evaluate(canvasScript).Do(ctx)
			return err
		}))
	}

	// Override audio fingerprinting
	if config.AudioFingerprint {
		audioScript := `
			const audioContext = window.AudioContext || window.webkitAudioContext;
			if (audioContext) {
				const getChannelData = audioContext.prototype.getChannelData;
				audioContext.prototype.getChannelData = function(channel) {
					const originalChannelData = getChannelData.call(this, channel);
					for (let i = 0; i < originalChannelData.length; i++) {
						originalChannelData[i] = originalChannelData[i] + Math.random() * 0.0001;
					}
					return originalChannelData;
				};
			}
		`
		actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, err := runtime.Evaluate(audioScript).Do(ctx)
			return err
		}))
	}

	return chromedp.Run(ctx, actions...)
}

// Get comprehensive stealth script
func (sb *StealthBrowser) getStealthScript() string {
	return `
		// Remove webdriver property
		Object.defineProperty(navigator, 'webdriver', {
			get: () => undefined,
		});

		// Override plugins
		Object.defineProperty(navigator, 'plugins', {
			get: () => [1, 2, 3, 4, 5],
		});

		// Override languages
		Object.defineProperty(navigator, 'languages', {
			get: () => ['en-US', 'en'],
		});

		// Override permissions
		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = (parameters) => (
			parameters.name === 'notifications' ?
				Promise.resolve({ state: Notification.permission }) :
				originalQuery(parameters)
		);

		// Override chrome runtime
		if (!window.chrome) {
			window.chrome = {};
		}
		if (!window.chrome.runtime) {
			window.chrome.runtime = {};
		}

		// Override toString methods
		const elementDescriptor = Object.getOwnPropertyDescriptor(HTMLElement.prototype, 'offsetHeight');
		Object.defineProperty(HTMLDivElement.prototype, 'offsetHeight', {
			...elementDescriptor,
			get: function() {
				if (this.id === 'modernizr') {
					return 1;
				}
				return elementDescriptor.get.apply(this);
			},
		});

		// Mock chrome.app
		if (!window.chrome.app) {
			window.chrome.app = {
				isInstalled: false,
				InstallState: {
					DISABLED: 'disabled',
					INSTALLED: 'installed',
					NOT_INSTALLED: 'not_installed'
				},
				RunningState: {
					CANNOT_RUN: 'cannot_run',
					READY_TO_RUN: 'ready_to_run',
					RUNNING: 'running'
				}
			};
		}

		// Mock chrome.csi
		if (!window.chrome.csi) {
			window.chrome.csi = function() {};
		}

		// Mock chrome.loadTimes
		if (!window.chrome.loadTimes) {
			window.chrome.loadTimes = function() {
				return {
					requestTime: Date.now() / 1000,
					startLoadTime: Date.now() / 1000,
					commitLoadTime: Date.now() / 1000,
					finishDocumentLoadTime: Date.now() / 1000,
					finishLoadTime: Date.now() / 1000,
					firstPaintTime: Date.now() / 1000,
					firstPaintAfterLoadTime: 0,
					navigationType: 'Other',
					wasFetchedViaSpdy: false,
					wasNpnNegotiated: false,
					npnNegotiatedProtocol: 'unknown',
					wasAlternateProtocolAvailable: false,
					connectionInfo: 'unknown'
				};
			};
		}
	`
}

// Login to a website using Bitwarden stored credentials
func (sb *StealthBrowser) Login(sessionID, loginSite string) (*BrowserResult, error) {
	result := &BrowserResult{
		SessionID:     sessionID,
		ExtractedData: make(map[string]interface{}),
		Errors:        []string{},
		Warnings:      []string{},
		ProcessedAt:   time.Now(),
	}

	// Get credentials from Bitwarden
	creds, err := sb.BitwardenClient.GetSiteCredentials(loginSite)
	if err != nil {
		return result, fmt.Errorf("failed to get credentials from Bitwarden for %s: %v", loginSite, err)
	}

	if creds.Username == "" || creds.Password == "" {
		return result, fmt.Errorf("credentials not configured in Bitwarden for %s", loginSite)
	}

	// Create context (in production, reuse existing session)
	opts := sb.buildStealthOptions(sb.AntiDetection)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	// Navigate to login page
	err = chromedp.Run(ctx, chromedp.Navigate(creds.LoginURL))
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to navigate to login page: %v", err))
		return result, err
	}

	// Wait for page load
	time.Sleep(sb.randomDelay(1*time.Second, 3*time.Second))

	// Apply basic stealth techniques
	err = sb.applyStealthTechniques(ctx, sb.AntiDetection)
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Basic stealth setup warning: %v", err))
	}

	// Apply advanced Bitwarden-style stealth
	err = sb.AdvancedStealth.ApplyBitwardenStealth(ctx)
	if err != nil {
		result.Warnings = append(result.Warnings, fmt.Sprintf("Advanced stealth setup warning: %v", err))
	}

	// Start session tracking
	sb.AdvancedStealth.StartSession()

	// Wait for username field to be visible
	err = chromedp.Run(ctx, chromedp.WaitVisible(creds.UsernameSelector))
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Username field not found: %v", err))
		return result, err
	}

	// Use advanced human-like typing for username
	err = sb.AdvancedStealth.HumanTypeText(ctx, creds.UsernameSelector, creds.Username)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to type username: %v", err))
		return result, err
	}

	// Human-like pause between fields
	time.Sleep(sb.randomDelay(800*time.Millisecond, 2000*time.Millisecond))

	// Use advanced human-like typing for password
	err = sb.AdvancedStealth.HumanTypeText(ctx, creds.PasswordSelector, creds.Password)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to type password: %v", err))
		return result, err
	}

	// Fill extra fields with human-like behavior
	for selector, value := range creds.ExtraFields {
		time.Sleep(sb.randomDelay(300*time.Millisecond, 800*time.Millisecond))
		err = sb.AdvancedStealth.HumanTypeText(ctx, selector, value)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Failed to fill extra field %s: %v", selector, err))
		}
	}

	// Human-like pause before submitting (thinking time)
	time.Sleep(sb.randomDelay(1500*time.Millisecond, 3000*time.Millisecond))

	// Submit form with human-like click
	err = sb.AdvancedStealth.HumanClick(ctx, creds.SubmitSelector)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Sprintf("Failed to click submit button: %v", err))
		return result, err
	}

	// Wait for login to complete
	time.Sleep(sb.randomDelay(3*time.Second, 5*time.Second))

	// Check if login was successful
	var currentURL, title string
	err = chromedp.Run(ctx,
		chromedp.Location(&currentURL),
		chromedp.Title(&title),
	)
	if err == nil {
		result.URL = currentURL
		result.Title = title

		// Simple check: if we're not on the login page anymore, assume success
		if !strings.Contains(currentURL, "sign-in") && !strings.Contains(currentURL, "login") {
			result.IsLoggedIn = true
			result.ExtractedData["loginStatus"] = "success"
		} else {
			result.ExtractedData["loginStatus"] = "failed"
		}
	}

	// Take screenshot
	var screenshot []byte
	err = chromedp.Run(ctx, chromedp.Screenshot("body", &screenshot))
	if err == nil {
		result.Screenshot = fmt.Sprintf("login_screenshot_%s.png", sessionID)
	}

	result.Success = len(result.Errors) == 0
	return result, nil
}

// Extract data from current page
func (sb *StealthBrowser) ExtractData(sessionID string, selectors []Selector) (*BrowserResult, error) {
	result := &BrowserResult{
		SessionID:     sessionID,
		ExtractedData: make(map[string]interface{}),
		Errors:        []string{},
		Warnings:      []string{},
		ProcessedAt:   time.Now(),
	}

	// In production, reuse existing session context
	opts := sb.buildStealthOptions(sb.AntiDetection)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Extract data using selectors (similar to existing extractData function)
	for _, selector := range selectors {
		var text string
		err := chromedp.Run(ctx, chromedp.Text(selector.CSS, &text))
		if err != nil {
			if selector.Required {
				result.Errors = append(result.Errors, fmt.Sprintf("Required selector '%s' failed: %v", selector.Name, err))
			} else {
				result.Warnings = append(result.Warnings, fmt.Sprintf("Optional selector '%s' failed: %v", selector.Name, err))
			}
			continue
		}
		result.ExtractedData[selector.Name] = strings.TrimSpace(text)
	}

	result.Success = len(result.Errors) == 0
	return result, nil
}

// Navigate to URL in existing session
func (sb *StealthBrowser) Navigate(sessionID, url string) (*BrowserResult, error) {
	// Implementation similar to CreateSession but for existing session
	return sb.CreateSession(url, sb.AntiDetection)
}

// Execute multiple actions in sequence
func (sb *StealthBrowser) ExecuteActions(sessionID string, actions []BrowserAction) (*BrowserResult, error) {
	result := &BrowserResult{
		SessionID:     sessionID,
		ExtractedData: make(map[string]interface{}),
		Errors:        []string{},
		Warnings:      []string{},
		ProcessedAt:   time.Now(),
	}

	// Implementation for executing browser actions
	// This would handle click, type, scroll, wait, etc.

	result.Success = len(result.Errors) == 0
	return result, nil
}

// Helper functions
func (sb *StealthBrowser) getRandomUserAgent() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(sb.UserAgents))))
	return sb.UserAgents[n.Int64()]
}

func (sb *StealthBrowser) getRandomViewport() ViewportSize {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(sb.ViewportSizes))))
	return sb.ViewportSizes[n.Int64()]
}

func (sb *StealthBrowser) randomDelay(min, max time.Duration) time.Duration {
	if max <= min {
		return min
	}
	diff := max - min
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	return min + time.Duration(n.Int64())
}

// Store credentials for a site (now handled by Bitwarden)
func (sb *StealthBrowser) StoreCredentials(site string, creds SiteCredential) error {
	return sb.BitwardenClient.StoreSiteCredentials(site, creds)
}

// Get stored credentials for a site (now handled by Bitwarden)
func (sb *StealthBrowser) GetCredentials(site string) (*SiteCredential, error) {
	return sb.BitwardenClient.GetSiteCredentials(site)
}
