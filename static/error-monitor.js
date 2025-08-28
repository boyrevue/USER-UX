// Chrome Error Monitor for Insurance Quote App
// Automatically captures and reports JavaScript errors, network failures, and console messages

(function() {
    'use strict';
    
    const ERROR_ENDPOINT = '/api/log-error';
    const APP_NAME = 'Insurance Quote App';
    
    // Error collection
    const errors = [];
    let errorCount = 0;
    
    // Create error display panel
    function createErrorPanel() {
        const panel = document.createElement('div');
        panel.id = 'error-monitor-panel';
        panel.style.cssText = `
            position: fixed;
            top: 10px;
            right: 10px;
            width: 400px;
            max-height: 300px;
            background: #1a1a1a;
            color: #fff;
            border: 2px solid #ff4444;
            border-radius: 8px;
            padding: 10px;
            font-family: 'Courier New', monospace;
            font-size: 12px;
            z-index: 10000;
            overflow-y: auto;
            display: none;
        `;
        
        panel.innerHTML = `
            <div style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px;">
                <strong>🚨 Error Monitor</strong>
                <button id="error-monitor-close" style="background: #ff4444; color: white; border: none; padding: 2px 6px; border-radius: 3px; cursor: pointer;">×</button>
            </div>
            <div id="error-monitor-content"></div>
            <div style="margin-top: 10px; text-align: center;">
                <button id="error-monitor-clear" style="background: #444; color: white; border: none; padding: 4px 8px; border-radius: 3px; cursor: pointer; margin-right: 5px;">Clear</button>
                <button id="error-monitor-export" style="background: #0066cc; color: white; border: none; padding: 4px 8px; border-radius: 3px; cursor: pointer;">Export</button>
            </div>
        `;
        
        document.body.appendChild(panel);
        
        // Event listeners
        document.getElementById('error-monitor-close').onclick = () => panel.style.display = 'none';
        document.getElementById('error-monitor-clear').onclick = clearErrors;
        document.getElementById('error-monitor-export').onclick = exportErrors;
        
        return panel;
    }
    
    // Create floating error indicator
    function createErrorIndicator() {
        const indicator = document.createElement('div');
        indicator.id = 'error-indicator';
        indicator.style.cssText = `
            position: fixed;
            bottom: 20px;
            right: 20px;
            width: 50px;
            height: 50px;
            background: #ff4444;
            color: white;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            font-weight: bold;
            cursor: pointer;
            z-index: 9999;
            box-shadow: 0 4px 8px rgba(0,0,0,0.3);
            display: none;
        `;
        indicator.innerHTML = '0';
        indicator.title = 'Click to view errors';
        
        indicator.onclick = () => {
            const panel = document.getElementById('error-monitor-panel');
            panel.style.display = panel.style.display === 'none' ? 'block' : 'none';
        };
        
        document.body.appendChild(indicator);
        return indicator;
    }
    
    // Log error function
    function logError(type, message, details = {}) {
        errorCount++;
        const timestamp = new Date().toISOString();
        
        const error = {
            id: errorCount,
            type,
            message,
            timestamp,
            url: window.location.href,
            userAgent: navigator.userAgent,
            ...details
        };
        
        errors.push(error);
        updateErrorDisplay();
        
        // Send to server (optional)
        sendErrorToServer(error);
        
        console.group(`🚨 ${APP_NAME} Error #${errorCount}`);
        console.error(`Type: ${type}`);
        console.error(`Message: ${message}`);
        console.error('Details:', details);
        console.error('Timestamp:', timestamp);
        console.groupEnd();
    }
    
    // Update error display
    function updateErrorDisplay() {
        const indicator = document.getElementById('error-indicator');
        const panel = document.getElementById('error-monitor-panel');
        const content = document.getElementById('error-monitor-content');
        
        if (errorCount > 0) {
            indicator.style.display = 'flex';
            indicator.innerHTML = errorCount;
            indicator.style.background = errorCount > 5 ? '#cc0000' : '#ff4444';
            
            content.innerHTML = errors.slice(-10).reverse().map(error => `
                <div style="margin-bottom: 8px; padding: 6px; background: #333; border-radius: 4px; border-left: 3px solid ${getErrorColor(error.type)};">
                    <div style="font-weight: bold; color: ${getErrorColor(error.type)};">#${error.id} ${error.type.toUpperCase()}</div>
                    <div style="margin: 2px 0; word-break: break-all;">${error.message}</div>
                    <div style="font-size: 10px; color: #888;">${new Date(error.timestamp).toLocaleTimeString()}</div>
                </div>
            `).join('');
        }
    }
    
    // Get error color
    function getErrorColor(type) {
        switch(type) {
            case 'network': return '#ff9900';
            case 'javascript': return '#ff4444';
            case 'console': return '#ffff44';
            case 'resource': return '#ff6600';
            case 'static-js-404': return '#ff0066';
            case 'static-css-404': return '#ff0066';
            case 'mime-js': return '#cc00ff';
            case 'mime-css': return '#cc00ff';
            default: return '#ff4444';
        }
    }
    
    // Clear errors
    function clearErrors() {
        errors.length = 0;
        errorCount = 0;
        updateErrorDisplay();
        document.getElementById('error-indicator').style.display = 'none';
    }
    
    // Export errors
    function exportErrors() {
        const data = {
            app: APP_NAME,
            timestamp: new Date().toISOString(),
            url: window.location.href,
            userAgent: navigator.userAgent,
            errors: errors
        };
        
        const blob = new Blob([JSON.stringify(data, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `errors-${new Date().toISOString().split('T')[0]}.json`;
        a.click();
        URL.revokeObjectURL(url);
    }
    
    // Send error to server
    function sendErrorToServer(error) {
        fetch(ERROR_ENDPOINT, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(error)
        }).catch(() => {}); // Ignore server errors for error logging
    }
    
    // Report static file issues to server for auto-fixing
    function reportStaticFileIssue(url, issueType) {
        const report = {
            url: url,
            issueType: issueType,
            timestamp: new Date().toISOString(),
            userAgent: navigator.userAgent,
            currentPage: window.location.href
        };
        
        fetch('/api/static-file-issue', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(report)
        }).catch(() => {}); // Ignore if server endpoint doesn't exist
        
        console.log('📡 Auto-reported static file issue:', report);
    }
    
    // Initialize when DOM is ready
    function initialize() {
        createErrorPanel();
        createErrorIndicator();
        
        // JavaScript errors
        window.addEventListener('error', (e) => {
            logError('javascript', e.message, {
                filename: e.filename,
                lineno: e.lineno,
                colno: e.colno,
                stack: e.error?.stack
            });
        });
        
        // Promise rejections
        window.addEventListener('unhandledrejection', (e) => {
            logError('javascript', `Unhandled Promise Rejection: ${e.reason}`, {
                reason: e.reason,
                stack: e.reason?.stack
            });
        });
        
        // Network errors (fetch)
        const originalFetch = window.fetch;
        window.fetch = function(...args) {
            return originalFetch.apply(this, args).catch(error => {
                logError('network', `Fetch failed: ${args[0]}`, {
                    url: args[0],
                    error: error.message
                });
                throw error;
            }).then(response => {
                if (!response.ok) {
                    logError('network', `HTTP ${response.status}: ${args[0]}`, {
                        url: args[0],
                        status: response.status,
                        statusText: response.statusText
                    });
                }
                return response;
            });
        };
        
        // Resource loading errors with static file analysis
        document.addEventListener('error', (e) => {
            if (e.target !== window) {
                const src = e.target.src || e.target.href;
                let message = `Failed to load: ${src}`;
                let errorType = 'resource';
                
                // Enhanced static file error detection
                if (src && src.includes('/static/')) {
                    if (src.includes('.js')) {
                        message += '\n🚨 JavaScript file missing - check build output';
                        errorType = 'static-js-404';
                    } else if (src.includes('.css')) {
                        message += '\n🚨 CSS file missing - check build output';
                        errorType = 'static-css-404';
                    }
                    
                    // Auto-report static file issues
                    reportStaticFileIssue(src, '404_not_found');
                }
                
                logError(errorType, message, {
                    element: e.target.tagName,
                    src: src,
                    isStaticFile: src && src.includes('/static/')
                });
            }
        }, true);
        
        // Console errors with MIME type detection
        const originalConsoleError = console.error;
        console.error = function(...args) {
            const message = args.join(' ');
            
            // Auto-detect MIME type errors for static files
            if (message.includes('MIME type') && message.includes('not executable')) {
                const match = message.match(/from '([^']+)'/);
                if (match) {
                    const url = match[1];
                    logError('mime-js', `❌ JavaScript MIME Error: ${url}\n🔧 Server serving 'text/plain' instead of 'application/javascript'`, {
                        url: url,
                        expectedMime: 'application/javascript',
                        actualMime: 'text/plain'
                    });
                    reportStaticFileIssue(url, 'wrong_mime_js');
                }
            } else if (message.includes('MIME type') && message.includes('stylesheet')) {
                const match = message.match(/from '([^']+)'/);
                if (match) {
                    const url = match[1];
                    logError('mime-css', `❌ CSS MIME Error: ${url}\n🔧 Server serving 'text/plain' instead of 'text/css'`, {
                        url: url,
                        expectedMime: 'text/css',
                        actualMime: 'text/plain'
                    });
                    reportStaticFileIssue(url, 'wrong_mime_css');
                }
            } else {
                logError('console', message, { args });
            }
            
            originalConsoleError.apply(console, args);
        };
        
        console.log('🚨 Error Monitor initialized for', APP_NAME);
    }
    
    // Initialize when DOM is ready
    if (document.readyState === 'loading') {
        document.addEventListener('DOMContentLoaded', initialize);
    } else {
        initialize();
    }
})();
