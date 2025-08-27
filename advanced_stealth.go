package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/dom"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/cdproto/runtime"
	"github.com/chromedp/chromedp"
)

// Advanced Anti-Bot Techniques inspired by Bitwarden
type AdvancedStealth struct {
	HumanBehavior   HumanBehaviorProfile `json:"humanBehavior"`
	SessionManager  SessionManager       `json:"sessionManager"`
	FingerprintMask FingerprintMask      `json:"fingerprintMask"`
	NetworkProfile  NetworkProfile       `json:"networkProfile"`
}

type HumanBehaviorProfile struct {
	TypingSpeed    TypingProfile  `json:"typingSpeed"`
	MouseMovement  MouseProfile   `json:"mouseMovement"`
	ScrollBehavior ScrollProfile  `json:"scrollBehavior"`
	FocusPatterns  FocusProfile   `json:"focusPatterns"`
	ReadingPauses  ReadingProfile `json:"readingPauses"`
}

type TypingProfile struct {
	BaseWPM         int     `json:"baseWpm"`         // Words per minute
	Variance        float64 `json:"variance"`        // Speed variation
	PauseFrequency  float64 `json:"pauseFrequency"`  // How often to pause
	BackspaceChance float64 `json:"backspaceChance"` // Chance to make "mistakes"
	BurstTyping     bool    `json:"burstTyping"`     // Type in bursts vs steady
}

type MouseProfile struct {
	MovementStyle string  `json:"movementStyle"` // "smooth", "jerky", "precise"
	ClickDelay    int     `json:"clickDelay"`    // ms before click after hover
	HoverBehavior bool    `json:"hoverBehavior"` // Hover before clicking
	MovementNoise float64 `json:"movementNoise"` // Random movement variation
}

type ScrollProfile struct {
	ScrollSpeed   int     `json:"scrollSpeed"`   // Pixels per scroll
	ScrollPauses  bool    `json:"scrollPauses"`  // Pause while scrolling
	ReadingScroll bool    `json:"readingScroll"` // Scroll like reading
	BackScroll    float64 `json:"backScroll"`    // Chance to scroll back up
}

type FocusProfile struct {
	TabSwitching bool `json:"tabSwitching"` // Switch tabs occasionally
	WindowBlur   bool `json:"windowBlur"`   // Simulate window focus loss
	IdleTime     int  `json:"idleTime"`     // Max idle time in seconds
}

type ReadingProfile struct {
	ReadingSpeed    int     `json:"readingSpeed"`    // Words per minute reading
	PauseOnKeywords bool    `json:"pauseOnKeywords"` // Pause on important words
	RereadChance    float64 `json:"rereadChance"`    // Chance to re-read sections
}

type SessionManager struct {
	SessionDuration time.Duration `json:"sessionDuration"`
	BreakFrequency  time.Duration `json:"breakFrequency"`
	SessionHistory  []SessionData `json:"sessionHistory"`
	CookieJar       []CookieData  `json:"cookieJar"`
}

type SessionData struct {
	StartTime    time.Time `json:"startTime"`
	EndTime      time.Time `json:"endTime"`
	PagesVisited int       `json:"pagesVisited"`
	ActionsCount int       `json:"actionsCount"`
	UserAgent    string    `json:"userAgent"`
}

type CookieData struct {
	Domain   string    `json:"domain"`
	Name     string    `json:"name"`
	Value    string    `json:"value"`
	Expires  time.Time `json:"expires"`
	Secure   bool      `json:"secure"`
	HttpOnly bool      `json:"httpOnly"`
}

type FingerprintMask struct {
	CanvasNoise    bool `json:"canvasNoise"`
	AudioNoise     bool `json:"audioNoise"`
	WebGLNoise     bool `json:"webglNoise"`
	FontMasking    bool `json:"fontMasking"`
	TimezoneSpoof  bool `json:"timezoneSpoof"`
	LanguageSpoof  bool `json:"languageSpoof"`
	ScreenResSpoof bool `json:"screenResSpoof"`
	HardwareMask   bool `json:"hardwareMask"`
}

type NetworkProfile struct {
	ConnectionType string        `json:"connectionType"` // "wifi", "cellular", "ethernet"
	Bandwidth      int           `json:"bandwidth"`      // Mbps
	Latency        time.Duration `json:"latency"`        // Network latency
	PacketLoss     float64       `json:"packetLoss"`     // Packet loss percentage
	DNSServers     []string      `json:"dnsServers"`     // DNS servers to use
}

// Initialize Advanced Stealth with Bitwarden-like profiles
func NewAdvancedStealth() *AdvancedStealth {
	return &AdvancedStealth{
		HumanBehavior: HumanBehaviorProfile{
			TypingSpeed: TypingProfile{
				BaseWPM:         45,   // Average human typing speed
				Variance:        0.3,  // 30% speed variation
				PauseFrequency:  0.15, // Pause 15% of the time
				BackspaceChance: 0.05, // 5% chance to make mistakes
				BurstTyping:     true,
			},
			MouseMovement: MouseProfile{
				MovementStyle: "smooth",
				ClickDelay:    150, // 150ms hover before click
				HoverBehavior: true,
				MovementNoise: 0.1, // 10% movement randomness
			},
			ScrollBehavior: ScrollProfile{
				ScrollSpeed:   120, // Pixels per scroll
				ScrollPauses:  true,
				ReadingScroll: true,
				BackScroll:    0.1, // 10% chance to scroll back
			},
			FocusPatterns: FocusProfile{
				TabSwitching: false, // Don't switch tabs during automation
				WindowBlur:   false, // Keep focus
				IdleTime:     30,    // Max 30 seconds idle
			},
			ReadingPauses: ReadingProfile{
				ReadingSpeed:    200, // 200 WPM reading
				PauseOnKeywords: true,
				RereadChance:    0.05, // 5% chance to re-read
			},
		},
		SessionManager: SessionManager{
			SessionDuration: 30 * time.Minute,
			BreakFrequency:  5 * time.Minute,
			SessionHistory:  []SessionData{},
			CookieJar:       []CookieData{},
		},
		FingerprintMask: FingerprintMask{
			CanvasNoise:    true,
			AudioNoise:     true,
			WebGLNoise:     true,
			FontMasking:    true,
			TimezoneSpoof:  true,
			LanguageSpoof:  false, // Keep consistent language
			ScreenResSpoof: true,
			HardwareMask:   true,
		},
		NetworkProfile: NetworkProfile{
			ConnectionType: "wifi",
			Bandwidth:      50, // 50 Mbps
			Latency:        20 * time.Millisecond,
			PacketLoss:     0.01, // 1% packet loss
			DNSServers:     []string{"8.8.8.8", "1.1.1.1"},
		},
	}
}

// Apply Bitwarden-style stealth techniques
func (as *AdvancedStealth) ApplyBitwardenStealth(ctx context.Context) error {
	var actions []chromedp.Action

	// 1. Advanced Browser Fingerprint Masking
	stealthScript := as.generateAdvancedStealthScript()
	actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
		_, _, err := runtime.Evaluate(stealthScript).Do(ctx)
		return err
	}))

	// 2. Hardware Fingerprint Masking
	if as.FingerprintMask.HardwareMask {
		hardwareScript := as.generateHardwareMaskScript()
		actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
			_, _, err := runtime.Evaluate(hardwareScript).Do(ctx)
			return err
		}))
	}

	// 3. Network Behavior Simulation
	networkScript := as.generateNetworkSimulationScript()
	actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
		_, _, err := runtime.Evaluate(networkScript).Do(ctx)
		return err
	}))

	// 4. Behavioral Timing Injection
	timingScript := as.generateTimingBehaviorScript()
	actions = append(actions, chromedp.ActionFunc(func(ctx context.Context) error {
		_, _, err := runtime.Evaluate(timingScript).Do(ctx)
		return err
	}))

	return chromedp.Run(ctx, actions...)
}

// Generate advanced stealth script that mimics Bitwarden's approach
func (as *AdvancedStealth) generateAdvancedStealthScript() string {
	return `
		// Bitwarden-style Advanced Stealth Techniques
		
		// 1. Deep Navigator Object Masking
		(function() {
			const originalNavigator = window.navigator;
			const navigatorProxy = new Proxy(originalNavigator, {
				get: function(target, prop) {
					switch(prop) {
						case 'webdriver':
							return undefined;
						case 'plugins':
							// Return realistic plugin array
							return [
								{name: 'Chrome PDF Plugin', filename: 'internal-pdf-viewer'},
								{name: 'Chrome PDF Viewer', filename: 'mhjfbmdgcfjbbpaeojofohoefgiehjai'},
								{name: 'Native Client', filename: 'internal-nacl-plugin'}
							];
						case 'languages':
							return ['en-US', 'en'];
						case 'platform':
							return 'MacIntel';
						case 'hardwareConcurrency':
							return 8;
						case 'deviceMemory':
							return 8;
						case 'maxTouchPoints':
							return 0;
						default:
							return target[prop];
					}
				}
			});
			
			Object.defineProperty(window, 'navigator', {
				value: navigatorProxy,
				writable: false,
				configurable: false
			});
		})();

		// 2. Advanced Chrome Runtime Masking
		if (!window.chrome) {
			window.chrome = {};
		}
		
		// Mock chrome.runtime with realistic responses
		window.chrome.runtime = {
			onConnect: {
				addListener: function() {},
				removeListener: function() {},
				hasListener: function() { return false; }
			},
			onMessage: {
				addListener: function() {},
				removeListener: function() {},
				hasListener: function() { return false; }
			},
			connect: function() {
				return {
					postMessage: function() {},
					onMessage: {
						addListener: function() {},
						removeListener: function() {}
					},
					onDisconnect: {
						addListener: function() {},
						removeListener: function() {}
					}
				};
			},
			sendMessage: function() {},
			getURL: function(path) {
				return 'chrome-extension://fake-extension-id/' + path;
			},
			getManifest: function() {
				return {
					name: 'Fake Extension',
					version: '1.0.0'
				};
			}
		};

		// 3. Permission API Masking (Bitwarden technique)
		const originalQuery = window.navigator.permissions.query;
		window.navigator.permissions.query = function(parameters) {
			const fakePermissions = {
				'notifications': 'default',
				'geolocation': 'denied',
				'camera': 'denied',
				'microphone': 'denied'
			};
			
			return Promise.resolve({
				state: fakePermissions[parameters.name] || 'denied',
				addEventListener: function() {},
				removeEventListener: function() {}
			});
		};

		// 4. Advanced Canvas Fingerprint Protection
		const originalGetContext = HTMLCanvasElement.prototype.getContext;
		HTMLCanvasElement.prototype.getContext = function(contextType, contextAttributes) {
			if (contextType === '2d') {
				const context = originalGetContext.call(this, contextType, contextAttributes);
				const originalGetImageData = context.getImageData;
				
				context.getImageData = function(sx, sy, sw, sh) {
					const imageData = originalGetImageData.call(this, sx, sy, sw, sh);
					
					// Add subtle noise that changes per session
					const noise = Math.sin(Date.now() / 1000) * 0.1;
					for (let i = 0; i < imageData.data.length; i += 4) {
						imageData.data[i] += Math.floor(noise * Math.random() * 10) - 5;
						imageData.data[i + 1] += Math.floor(noise * Math.random() * 10) - 5;
						imageData.data[i + 2] += Math.floor(noise * Math.random() * 10) - 5;
					}
					return imageData;
				};
				
				return context;
			}
			return originalGetContext.call(this, contextType, contextAttributes);
		};

		// 5. WebGL Fingerprint Protection
		const originalGetParameter = WebGLRenderingContext.prototype.getParameter;
		WebGLRenderingContext.prototype.getParameter = function(parameter) {
			const fakeValues = {
				37445: 'Intel Inc.',  // VENDOR
				37446: 'Intel Iris Pro OpenGL Engine',  // RENDERER
				7936: 'WebGL 1.0 (OpenGL ES 2.0 Chromium)',  // VERSION
				35724: 'WebGL GLSL ES 1.0 (OpenGL ES GLSL ES 1.0 Chromium)'  // SHADING_LANGUAGE_VERSION
			};
			
			return fakeValues[parameter] || originalGetParameter.call(this, parameter);
		};

		// 6. Audio Context Fingerprint Protection
		const originalCreateAnalyser = AudioContext.prototype.createAnalyser;
		AudioContext.prototype.createAnalyser = function() {
			const analyser = originalCreateAnalyser.call(this);
			const originalGetFloatFrequencyData = analyser.getFloatFrequencyData;
			
			analyser.getFloatFrequencyData = function(array) {
				originalGetFloatFrequencyData.call(this, array);
				// Add subtle audio noise
				for (let i = 0; i < array.length; i++) {
					array[i] += Math.random() * 0.001 - 0.0005;
				}
			};
			
			return analyser;
		};

		// 7. Screen Resolution Spoofing
		Object.defineProperty(window.screen, 'width', {
			get: function() { return 1920; }
		});
		Object.defineProperty(window.screen, 'height', {
			get: function() { return 1080; }
		});
		Object.defineProperty(window.screen, 'availWidth', {
			get: function() { return 1920; }
		});
		Object.defineProperty(window.screen, 'availHeight', {
			get: function() { return 1050; }
		});

		// 8. Timezone Consistency
		const originalGetTimezoneOffset = Date.prototype.getTimezoneOffset;
		Date.prototype.getTimezoneOffset = function() {
			return 0; // UTC
		};

		// 9. Battery API Masking
		if (navigator.getBattery) {
			navigator.getBattery = function() {
				return Promise.resolve({
					charging: true,
					chargingTime: Infinity,
					dischargingTime: Infinity,
					level: 1.0,
					addEventListener: function() {},
					removeEventListener: function() {}
				});
			};
		}

		// 10. Connection API Masking
		if (navigator.connection) {
			Object.defineProperty(navigator, 'connection', {
				value: {
					effectiveType: '4g',
					downlink: 10,
					rtt: 50,
					saveData: false,
					addEventListener: function() {},
					removeEventListener: function() {}
				},
				writable: false
			});
		}

		// 11. Memory Info Masking
		if (performance.memory) {
			Object.defineProperty(performance, 'memory', {
				value: {
					usedJSHeapSize: 10000000,
					totalJSHeapSize: 20000000,
					jsHeapSizeLimit: 2147483648
				},
				writable: false
			});
		}

		// 12. Advanced Event Listener Protection
		const originalAddEventListener = EventTarget.prototype.addEventListener;
		EventTarget.prototype.addEventListener = function(type, listener, options) {
			// Block certain automation detection events
			const blockedEvents = ['webkitvisibilitychange', 'mozvisibilitychange'];
			if (blockedEvents.includes(type)) {
				return;
			}
			return originalAddEventListener.call(this, type, listener, options);
		};

		console.log('🥷 Advanced Bitwarden-style stealth activated');
	`
}

// Generate hardware masking script
func (as *AdvancedStealth) generateHardwareMaskScript() string {
	return `
		// Hardware Fingerprint Masking
		
		// Mock realistic hardware specs
		Object.defineProperty(navigator, 'hardwareConcurrency', {
			get: function() { return 8; }
		});
		
		Object.defineProperty(navigator, 'deviceMemory', {
			get: function() { return 8; }
		});
		
		// Mock GPU info
		const canvas = document.createElement('canvas');
		const gl = canvas.getContext('webgl');
		if (gl) {
			const debugInfo = gl.getExtension('WEBGL_debug_renderer_info');
			if (debugInfo) {
				Object.defineProperty(gl, 'getParameter', {
					value: function(parameter) {
						if (parameter === debugInfo.UNMASKED_VENDOR_WEBGL) {
							return 'Intel Inc.';
						}
						if (parameter === debugInfo.UNMASKED_RENDERER_WEBGL) {
							return 'Intel Iris Pro OpenGL Engine';
						}
						return WebGLRenderingContext.prototype.getParameter.call(this, parameter);
					}
				});
			}
		}
	`
}

// Generate network simulation script
func (as *AdvancedStealth) generateNetworkSimulationScript() string {
	return `
		// Network Behavior Simulation
		
		// Simulate realistic connection timing
		const originalFetch = window.fetch;
		window.fetch = function(url, options) {
			const delay = Math.random() * 100 + 50; // 50-150ms delay
			return new Promise(resolve => {
				setTimeout(() => {
					resolve(originalFetch.call(this, url, options));
				}, delay);
			});
		};
		
		// Simulate realistic XMLHttpRequest timing
		const originalOpen = XMLHttpRequest.prototype.open;
		XMLHttpRequest.prototype.open = function(method, url, async, user, password) {
			const xhr = this;
			const originalSend = xhr.send;
			
			xhr.send = function(data) {
				const delay = Math.random() * 50 + 25; // 25-75ms delay
				setTimeout(() => {
					originalSend.call(xhr, data);
				}, delay);
			};
			
			return originalOpen.call(this, method, url, async, user, password);
		};
	`
}

// Generate timing behavior script
func (as *AdvancedStealth) generateTimingBehaviorScript() string {
	return `
		// Human-like Timing Behavior
		
		// Override setTimeout/setInterval to add human variance
		const originalSetTimeout = window.setTimeout;
		window.setTimeout = function(callback, delay) {
			const humanDelay = delay + (Math.random() * 100 - 50); // ±50ms variance
			return originalSetTimeout.call(this, callback, Math.max(0, humanDelay));
		};
		
		// Add realistic performance timing
		const originalNow = performance.now;
		let timeOffset = Math.random() * 1000;
		performance.now = function() {
			return originalNow.call(this) + timeOffset;
		};
	`
}

// Human-like typing with Bitwarden-style patterns
func (as *AdvancedStealth) HumanTypeText(ctx context.Context, selector, text string) error {
	// Click the element first
	err := chromedp.Run(ctx, chromedp.Click(selector))
	if err != nil {
		return err
	}

	// Add realistic pre-typing delay
	delay := as.generateHumanDelay(200, 800)
	time.Sleep(delay)

	// Type with human-like patterns
	for i, char := range text {
		// Simulate occasional mistakes
		if as.HumanBehavior.TypingSpeed.BackspaceChance > 0 &&
			as.randomFloat() < as.HumanBehavior.TypingSpeed.BackspaceChance {
			// Type wrong character then backspace
			wrongChar := as.generateRandomChar()
			err = chromedp.Run(ctx, chromedp.KeyEvent(string(wrongChar)))
			if err != nil {
				return err
			}

			time.Sleep(as.generateTypingDelay())

			// Backspace
			err = chromedp.Run(ctx, chromedp.KeyEvent("\b"))
			if err != nil {
				return err
			}

			time.Sleep(as.generateTypingDelay())
		}

		// Type the actual character
		err = chromedp.Run(ctx, chromedp.KeyEvent(string(char)))
		if err != nil {
			return err
		}

		// Human-like delay between characters
		if i < len(text)-1 {
			delay := as.generateTypingDelay()

			// Longer pauses at word boundaries
			if char == ' ' {
				delay = delay * 2
			}

			// Occasional longer pauses (thinking)
			if as.randomFloat() < as.HumanBehavior.TypingSpeed.PauseFrequency {
				delay = delay * 3
			}

			time.Sleep(delay)
		}
	}

	return nil
}

// Human-like clicking with movement
func (as *AdvancedStealth) HumanClick(ctx context.Context, selector string) error {
	// Get element bounds
	var nodes []*cdp.Node
	err := chromedp.Run(ctx, chromedp.Nodes(selector, &nodes))
	if err != nil {
		return err
	}

	if len(nodes) == 0 {
		return fmt.Errorf("element not found: %s", selector)
	}

	// Get element position
	var box *dom.BoxModel
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		box, err = dom.GetBoxModel().WithNodeID(nodes[0].NodeID).Do(ctx)
		return err
	}))
	if err != nil {
		return err
	}

	// Calculate click position with human-like offset
	centerX := (box.Content[0] + box.Content[2]) / 2
	centerY := (box.Content[1] + box.Content[5]) / 2

	// Add human-like randomness to click position
	offsetX := (as.randomFloat() - 0.5) * 20 // ±10px
	offsetY := (as.randomFloat() - 0.5) * 20 // ±10px

	clickX := centerX + offsetX
	clickY := centerY + offsetY

	// Simulate mouse movement to element (if hover behavior enabled)
	if as.HumanBehavior.MouseMovement.HoverBehavior {
		err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
			return input.DispatchMouseEvent(input.MouseMoved, clickX, clickY).Do(ctx)
		}))
		if err != nil {
			return err
		}

		// Hover delay
		hoverDelay := time.Duration(as.HumanBehavior.MouseMovement.ClickDelay) * time.Millisecond
		time.Sleep(hoverDelay)
	}

	// Perform the click
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return input.DispatchMouseEvent(input.MousePressed, clickX, clickY).
			WithButton(input.Left).
			WithClickCount(1).
			Do(ctx)
	}))
	if err != nil {
		return err
	}

	// Mouse up
	err = chromedp.Run(ctx, chromedp.ActionFunc(func(ctx context.Context) error {
		return input.DispatchMouseEvent(input.MouseReleased, clickX, clickY).
			WithButton(input.Left).
			WithClickCount(1).
			Do(ctx)
	}))

	return err
}

// Generate human-like delays
func (as *AdvancedStealth) generateHumanDelay(min, max time.Duration) time.Duration {
	if max <= min {
		return min
	}

	diff := max - min
	randomMs, _ := rand.Int(rand.Reader, big.NewInt(int64(diff/time.Millisecond)))

	return min + time.Duration(randomMs.Int64())*time.Millisecond
}

// Generate typing delay based on human behavior profile
func (as *AdvancedStealth) generateTypingDelay() time.Duration {
	// Base delay from WPM
	baseDelay := time.Duration(60000/as.HumanBehavior.TypingSpeed.BaseWPM/5) * time.Millisecond

	// Add variance
	variance := float64(baseDelay) * as.HumanBehavior.TypingSpeed.Variance
	randomVariance := (as.randomFloat() - 0.5) * variance

	finalDelay := time.Duration(float64(baseDelay) + randomVariance)

	// Ensure minimum delay
	if finalDelay < 50*time.Millisecond {
		finalDelay = 50 * time.Millisecond
	}

	return finalDelay
}

// Helper functions
func (as *AdvancedStealth) randomFloat() float64 {
	n, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	return float64(n.Int64()) / 1000000.0
}

func (as *AdvancedStealth) generateRandomChar() rune {
	chars := "abcdefghijklmnopqrstuvwxyz"
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(chars))))
	return rune(chars[n.Int64()])
}

// Session management like Bitwarden
func (as *AdvancedStealth) StartSession() {
	session := SessionData{
		StartTime:    time.Now(),
		PagesVisited: 0,
		ActionsCount: 0,
		UserAgent:    as.generateRealisticUserAgent(),
	}

	as.SessionManager.SessionHistory = append(as.SessionManager.SessionHistory, session)
}

func (as *AdvancedStealth) generateRealisticUserAgent() string {
	userAgents := []string{
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.1 Safari/605.1.15",
	}

	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(userAgents))))
	return userAgents[n.Int64()]
}
