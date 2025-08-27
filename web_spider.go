package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// Web Spider and Form Filler Engine
type WebSpider struct {
	UserAgent    string            `json:"userAgent"`
	Timeout      time.Duration     `json:"timeout"`
	MaxDepth     int               `json:"maxDepth"`
	RateLimiting time.Duration     `json:"rateLimiting"`
	Headers      map[string]string `json:"headers"`
}

type SpiderTask struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"` // "extract", "fill", "navigate"
	URL         string            `json:"url"`
	Selectors   []Selector        `json:"selectors"`
	FormData    map[string]string `json:"formData"`
	Navigation  []NavigationStep  `json:"navigation"`
	Options     SpiderOptions     `json:"options"`
	Status      string            `json:"status"`
	CreatedAt   time.Time         `json:"createdAt"`
	CompletedAt *time.Time        `json:"completedAt,omitempty"`
}

type Selector struct {
	Name     string `json:"name"`
	CSS      string `json:"css"`
	XPath    string `json:"xpath"`
	Attr     string `json:"attr"` // attribute to extract
	Text     bool   `json:"text"` // extract text content
	Required bool   `json:"required"`
	Multiple bool   `json:"multiple"` // extract multiple elements
}

type NavigationStep struct {
	Action    string `json:"action"` // "click", "type", "select", "wait", "scroll"
	Selector  string `json:"selector"`
	Value     string `json:"value"`
	WaitFor   string `json:"waitFor"`
	Timeout   int    `json:"timeout"`
	Condition string `json:"condition"`
}

type SpiderOptions struct {
	Headless        bool              `json:"headless"`
	Screenshots     bool              `json:"screenshots"`
	WaitForLoad     bool              `json:"waitForLoad"`
	HandleCookies   bool              `json:"handleCookies"`
	FollowRedirects bool              `json:"followRedirects"`
	CustomHeaders   map[string]string `json:"customHeaders"`
	Proxy           string            `json:"proxy"`
	JavaScript      bool              `json:"javascript"`
}

type SpiderResult struct {
	TaskID        string                 `json:"taskId"`
	Success       bool                   `json:"success"`
	URL           string                 `json:"url"`
	ExtractedData map[string]interface{} `json:"extractedData"`
	FormsFilled   []FormFillResult       `json:"formsFilled"`
	Screenshots   []string               `json:"screenshots"`
	Errors        []string               `json:"errors"`
	Warnings      []string               `json:"warnings"`
	Metadata      SpiderMetadata         `json:"metadata"`
	ProcessedAt   time.Time              `json:"processedAt"`
}

type FormFillResult struct {
	FormSelector string            `json:"formSelector"`
	FieldsFilled map[string]string `json:"fieldsFilled"`
	Success      bool              `json:"success"`
	Error        string            `json:"error,omitempty"`
}

type SpiderMetadata struct {
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Keywords    []string      `json:"keywords"`
	LoadTime    time.Duration `json:"loadTime"`
	PageSize    int64         `json:"pageSize"`
	StatusCode  int           `json:"statusCode"`
	FinalURL    string        `json:"finalUrl"`
}

// Web Spider Handler
func WebSpiderHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var task SpiderTask
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Initialize spider
	spider := NewWebSpider()

	// Execute task
	result, err := spider.ExecuteTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Spider execution failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Data Extraction Handler
func ExtractDataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		URL       string        `json:"url"`
		Selectors []Selector    `json:"selectors"`
		Options   SpiderOptions `json:"options"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	spider := NewWebSpider()
	task := SpiderTask{
		ID:        generateTaskID(),
		Type:      "extract",
		URL:       request.URL,
		Selectors: request.Selectors,
		Options:   request.Options,
		Status:    "pending",
		CreatedAt: time.Now(),
	}

	result, err := spider.ExecuteTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Data extraction failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Form Fill Handler
func FillFormHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request struct {
		URL        string            `json:"url"`
		FormData   map[string]string `json:"formData"`
		Navigation []NavigationStep  `json:"navigation"`
		Options    SpiderOptions     `json:"options"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	spider := NewWebSpider()
	task := SpiderTask{
		ID:         generateTaskID(),
		Type:       "fill",
		URL:        request.URL,
		FormData:   request.FormData,
		Navigation: request.Navigation,
		Options:    request.Options,
		Status:     "pending",
		CreatedAt:  time.Now(),
	}

	result, err := spider.ExecuteTask(task)
	if err != nil {
		http.Error(w, fmt.Sprintf("Form filling failed: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Initialize Web Spider
func NewWebSpider() *WebSpider {
	return &WebSpider{
		UserAgent:    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		Timeout:      30 * time.Second,
		MaxDepth:     3,
		RateLimiting: 1 * time.Second,
		Headers: map[string]string{
			"Accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8",
			"Accept-Language":           "en-US,en;q=0.5",
			"Accept-Encoding":           "gzip, deflate",
			"DNT":                       "1",
			"Connection":                "keep-alive",
			"Upgrade-Insecure-Requests": "1",
		},
	}
}

// Execute Spider Task
func (ws *WebSpider) ExecuteTask(task SpiderTask) (*SpiderResult, error) {
	result := &SpiderResult{
		TaskID:        task.ID,
		URL:           task.URL,
		ExtractedData: make(map[string]interface{}),
		FormsFilled:   []FormFillResult{},
		Screenshots:   []string{},
		Errors:        []string{},
		Warnings:      []string{},
		ProcessedAt:   time.Now(),
	}

	// Create Chrome context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", task.Options.Headless),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.UserAgent(ws.UserAgent),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, ws.Timeout)
	defer cancel()

	startTime := time.Now()

	switch task.Type {
	case "extract":
		err := ws.extractData(ctx, task, result)
		if err != nil {
			result.Success = false
			result.Errors = append(result.Errors, err.Error())
			return result, err
		}
	case "fill":
		err := ws.fillForms(ctx, task, result)
		if err != nil {
			result.Success = false
			result.Errors = append(result.Errors, err.Error())
			return result, err
		}
	case "navigate":
		err := ws.navigateAndProcess(ctx, task, result)
		if err != nil {
			result.Success = false
			result.Errors = append(result.Errors, err.Error())
			return result, err
		}
	default:
		return nil, fmt.Errorf("unknown task type: %s", task.Type)
	}

	// Collect metadata
	result.Metadata.LoadTime = time.Since(startTime)
	ws.collectMetadata(ctx, result)

	result.Success = len(result.Errors) == 0
	return result, nil
}

// Extract data from webpage
func (ws *WebSpider) extractData(ctx context.Context, task SpiderTask, result *SpiderResult) error {
	// Navigate to URL
	err := chromedp.Run(ctx, chromedp.Navigate(task.URL))
	if err != nil {
		return fmt.Errorf("failed to navigate to %s: %v", task.URL, err)
	}

	// Wait for page load
	if task.Options.WaitForLoad {
		err = chromedp.Run(ctx, chromedp.WaitReady("body"))
		if err != nil {
			result.Warnings = append(result.Warnings, "Page may not have fully loaded")
		}
	}

	// Extract data using selectors
	for _, selector := range task.Selectors {
		var extractedValue interface{}

		if selector.Multiple {
			// Extract multiple elements
			var texts []string
			var nodes []*cdp.Node

			if selector.CSS != "" {
				err = chromedp.Run(ctx, chromedp.Nodes(selector.CSS, &nodes))
			} else if selector.XPath != "" {
				err = chromedp.Run(ctx, chromedp.Nodes(selector.XPath, &nodes))
			}

			if err != nil {
				if selector.Required {
					return fmt.Errorf("required selector '%s' not found: %v", selector.Name, err)
				}
				result.Warnings = append(result.Warnings, fmt.Sprintf("Optional selector '%s' not found", selector.Name))
				continue
			}

			for _, node := range nodes {
				var text string
				if selector.Text {
					text = node.Children[0].NodeValue
				} else if selector.Attr != "" {
					for _, attr := range node.Attributes {
						if attr == selector.Attr && len(node.Attributes) > 1 {
							text = node.Attributes[1] // Attribute value
							break
						}
					}
				}
				texts = append(texts, text)
			}
			extractedValue = texts
		} else {
			// Extract single element
			var text string

			if selector.Text {
				if selector.CSS != "" {
					err = chromedp.Run(ctx, chromedp.Text(selector.CSS, &text))
				} else if selector.XPath != "" {
					err = chromedp.Run(ctx, chromedp.Text(selector.XPath, &text))
				}
			} else if selector.Attr != "" {
				if selector.CSS != "" {
					err = chromedp.Run(ctx, chromedp.AttributeValue(selector.CSS, selector.Attr, &text, nil))
				}
			}

			if err != nil {
				if selector.Required {
					return fmt.Errorf("required selector '%s' not found: %v", selector.Name, err)
				}
				result.Warnings = append(result.Warnings, fmt.Sprintf("Optional selector '%s' not found", selector.Name))
				continue
			}

			extractedValue = text
		}

		result.ExtractedData[selector.Name] = extractedValue
	}

	return nil
}

// Fill forms on webpage
func (ws *WebSpider) fillForms(ctx context.Context, task SpiderTask, result *SpiderResult) error {
	// Navigate to URL
	err := chromedp.Run(ctx, chromedp.Navigate(task.URL))
	if err != nil {
		return fmt.Errorf("failed to navigate to %s: %v", task.URL, err)
	}

	// Wait for page load
	if task.Options.WaitForLoad {
		err = chromedp.Run(ctx, chromedp.WaitReady("body"))
		if err != nil {
			result.Warnings = append(result.Warnings, "Page may not have fully loaded")
		}
	}

	// Execute navigation steps first
	for _, step := range task.Navigation {
		err = ws.executeNavigationStep(ctx, step)
		if err != nil {
			result.Warnings = append(result.Warnings, fmt.Sprintf("Navigation step failed: %v", err))
		}
	}

	// Fill form fields
	formResult := FormFillResult{
		FormSelector: "form", // Default form selector
		FieldsFilled: make(map[string]string),
		Success:      true,
	}

	for fieldName, value := range task.FormData {
		// Try multiple common selectors for the field
		selectors := []string{
			fmt.Sprintf(`input[name="%s"]`, fieldName),
			fmt.Sprintf(`input[id="%s"]`, fieldName),
			fmt.Sprintf(`select[name="%s"]`, fieldName),
			fmt.Sprintf(`select[id="%s"]`, fieldName),
			fmt.Sprintf(`textarea[name="%s"]`, fieldName),
			fmt.Sprintf(`textarea[id="%s"]`, fieldName),
		}

		filled := false
		for _, selector := range selectors {
			// Check if element exists
			var nodes []*cdp.Node
			err = chromedp.Run(ctx, chromedp.Nodes(selector, &nodes))
			if err != nil || len(nodes) == 0 {
				continue
			}

			// Fill the field
			err = chromedp.Run(ctx, chromedp.SetValue(selector, value))
			if err == nil {
				formResult.FieldsFilled[fieldName] = value
				filled = true
				break
			}
		}

		if !filled {
			formResult.Success = false
			formResult.Error = fmt.Sprintf("Could not fill field: %s", fieldName)
			result.Warnings = append(result.Warnings, formResult.Error)
		}
	}

	result.FormsFilled = append(result.FormsFilled, formResult)
	return nil
}

// Navigate and process multi-step workflows
func (ws *WebSpider) navigateAndProcess(ctx context.Context, task SpiderTask, result *SpiderResult) error {
	// Navigate to URL
	err := chromedp.Run(ctx, chromedp.Navigate(task.URL))
	if err != nil {
		return fmt.Errorf("failed to navigate to %s: %v", task.URL, err)
	}

	// Execute navigation steps
	for i, step := range task.Navigation {
		err = ws.executeNavigationStep(ctx, step)
		if err != nil {
			return fmt.Errorf("navigation step %d failed: %v", i+1, err)
		}

		// Take screenshot if requested
		if task.Options.Screenshots {
			var screenshot []byte
			err = chromedp.Run(ctx, chromedp.Screenshot("body", &screenshot))
			if err == nil {
				// In a real implementation, save screenshot and add path to result
				result.Screenshots = append(result.Screenshots, fmt.Sprintf("screenshot_%d.png", i+1))
			}
		}
	}

	return nil
}

// Execute individual navigation step
func (ws *WebSpider) executeNavigationStep(ctx context.Context, step NavigationStep) error {
	timeout := 5 * time.Second
	if step.Timeout > 0 {
		timeout = time.Duration(step.Timeout) * time.Second
	}

	stepCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	switch step.Action {
	case "click":
		return chromedp.Run(stepCtx, chromedp.Click(step.Selector))
	case "type":
		return chromedp.Run(stepCtx, chromedp.SetValue(step.Selector, step.Value))
	case "select":
		return chromedp.Run(stepCtx, chromedp.SetAttributeValue(step.Selector, "value", step.Value))
	case "wait":
		if step.WaitFor != "" {
			return chromedp.Run(stepCtx, chromedp.WaitVisible(step.WaitFor))
		}
		time.Sleep(timeout)
		return nil
	case "scroll":
		return chromedp.Run(stepCtx, chromedp.ScrollIntoView(step.Selector))
	default:
		return fmt.Errorf("unknown navigation action: %s", step.Action)
	}
}

// Collect page metadata
func (ws *WebSpider) collectMetadata(ctx context.Context, result *SpiderResult) {
	var title string
	chromedp.Run(ctx, chromedp.Title(&title))
	result.Metadata.Title = title

	var description string
	chromedp.Run(ctx, chromedp.AttributeValue(`meta[name="description"]`, "content", &description, nil))
	result.Metadata.Description = description

	var finalURL string
	chromedp.Run(ctx, chromedp.Location(&finalURL))
	result.Metadata.FinalURL = finalURL
}

// Helper functions
func generateTaskID() string {
	return fmt.Sprintf("task_%d", time.Now().Unix())
}

// URL validation
func isValidURL(rawURL string) bool {
	_, err := url.ParseRequestURI(rawURL)
	return err == nil
}

// Clean extracted text
func cleanText(text string) string {
	// Remove extra whitespace and normalize
	re := regexp.MustCompile(`\s+`)
	return strings.TrimSpace(re.ReplaceAllString(text, " "))
}
