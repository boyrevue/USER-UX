package formgenerator

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func main() {
	outputFile := flag.String("output", "settings_form.html", "Output HTML file")
	flag.Parse()

	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Settings Configuration</title>
    <style>
        :root {
            --bg: #1a1a1a;
            --fg: #e0e0e0;
            --card-bg: #2a2a2a;
            --border: #404040;
            --accent: #505050;
            --accent-green: #00ff88;
            --pure-black: #000000;
        }

        * {
            margin: 0;
            padding: 0;
            box-sizing: border-box;
        }

        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: var(--bg);
            color: var(--fg);
            line-height: 1.6;
        }

        .container {
            display: grid;
            grid-template-columns: 250px 1fr;
            height: 100vh;
            overflow: hidden;
        }

        .left-panel {
            background: var(--card-bg);
            border-right: 1px solid var(--border);
            padding: 20px 0;
            overflow-y: auto;
        }

        .category-tab {
            padding: 12px 20px;
            cursor: pointer;
            border-bottom: 1px solid var(--border);
            transition: all 0.3s ease;
            font-size: 14px;
            font-weight: 500;
        }

        .category-tab:hover {
            background: var(--accent);
        }

        .category-tab.active {
            background: var(--accent-green);
            color: var(--pure-black);
        }

        .right-panel {
            padding: 30px;
            overflow-y: auto;
        }

        .category-content {
            display: none;
        }

        .category-content.active {
            display: block;
        }

        .form-section {
            background: var(--card-bg);
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 24px;
            margin-bottom: 24px;
        }

        .form-section-title {
            font-size: 18px;
            font-weight: 600;
            margin-bottom: 20px;
            color: var(--fg);
        }

        .form-group {
            margin-bottom: 20px;
        }

        .form-label {
            display: block;
            margin-bottom: 8px;
            font-weight: 500;
            color: var(--fg);
            font-size: 14px;
        }

        .form-input {
            width: 100%;
            padding: 12px 16px;
            border: 1px solid var(--border);
            border-radius: 6px;
            background: var(--bg);
            color: var(--fg);
            font-size: 14px;
            transition: all 0.3s ease;
        }

        .form-input:focus {
            outline: none;
            border-color: var(--accent-green);
            box-shadow: 0 0 0 2px rgba(0, 255, 136, 0.2);
        }

        .form-select {
            width: 100%;
            padding: 12px 16px;
            border: 1px solid var(--border);
            border-radius: 6px;
            background: var(--bg);
            color: var(--fg);
            font-size: 14px;
            cursor: pointer;
        }

        .help-text {
            font-size: 12px;
            color: #808080;
            margin-top: 4px;
            line-height: 1.4;
        }

        .language-toggle {
            position: fixed;
            top: 20px;
            right: 20px;
            background: var(--card-bg);
            border: 1px solid var(--border);
            border-radius: 6px;
            padding: 8px 12px;
            font-size: 12px;
            cursor: pointer;
            z-index: 1000;
        }

        /* Floating modals */
        .floating-modal {
            position: fixed;
            background: var(--card-bg);
            border: 1px solid var(--border);
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
            border-bottom: 1px solid var(--border);
            cursor: move;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }

        .modal-title {
            font-size: 14px;
            font-weight: 600;
            color: var(--fg);
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
            background: var(--accent);
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
            color: var(--fg);
        }

        .modal-content {
            padding: 16px;
            height: calc(100% - 50px);
            overflow-y: auto;
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
            background: var(--accent);
            padding: 12px 16px;
            border-radius: 8px;
            max-width: 85%;
            color: var(--fg);
            font-size: 13px;
            line-height: 1.4;
        }

        .message.user .message-content {
            background: #505050;
            color: var(--fg);
        }

        .message-content strong {
            color: var(--fg);
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
            border: 1px solid var(--border);
            color: var(--fg);
            padding: 8px 12px;
            border-radius: 6px;
            font-size: 11px;
            cursor: pointer;
            text-align: left;
            transition: all 0.3s ease;
        }

        .quick-btn:hover {
            background: var(--accent);
            border-color: #505050;
        }

        .chat-input-container {
            display: flex;
            gap: 8px;
            padding: 16px;
            border-top: 1px solid var(--border);
            background: var(--card-bg);
        }

        .chat-input-container input {
            flex: 1;
            padding: 8px 12px;
            border: 1px solid var(--border);
            border-radius: 4px;
            background: #333333;
            color: var(--fg);
            font-size: 13px;
        }

        .chat-input-container input:focus {
            outline: none;
            border-color: #505050;
        }

        .chat-input-container button {
            padding: 8px 16px;
            background: var(--accent);
            border: none;
            border-radius: 4px;
            color: var(--fg);
            cursor: pointer;
            font-size: 13px;
            transition: all 0.3s ease;
        }

        .chat-input-container button:hover {
            background: #505050;
        }

        /* Upload Modal Styles */
        .upload-instructions {
            margin-bottom: 20px;
        }

        .upload-instructions h4 {
            color: var(--fg);
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
            border: 1px solid var(--border);
            border-radius: 6px;
        }

        .doc-icon {
            font-size: 24px;
            width: 40px;
            text-align: center;
        }

        .doc-info strong {
            color: var(--fg);
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
            border: 2px dashed var(--border);
            border-radius: 8px;
            margin-bottom: 20px;
            cursor: pointer;
            transition: all 0.3s ease;
        }

        .upload-zone:hover {
            border-color: #505050;
            background: var(--card-bg);
        }

        .upload-results {
            background: var(--card-bg);
            border: 1px solid var(--border);
            border-radius: 6px;
            padding: 16px;
        }

        .upload-results h4 {
            color: var(--fg);
            font-size: 14px;
            margin-bottom: 12px;
            font-weight: 600;
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
            font-size: 13px;
            transition: all 0.3s ease;
        }

        .action-btn:hover {
            background: #00cc6a;
            transform: translateY(-1px);
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
    </style>
</head>
<body>
    <div class="container">
        <div class="left-panel">
            <div class="category-tab active" onclick="switchCategory('bankaccounts')">Bank Accounts</div>
            <div class="category-tab" onclick="switchCategory('creditcards')">Credit Cards</div>
            <div class="category-tab" onclick="switchCategory('communicationchannels')">Communication Channels</div>
            <div class="category-tab" onclick="switchCategory('security')">Security</div>
        </div>

        <div class="right-panel">
            <div id="bankaccountsContent" class="category-content active">
                <div class="form-section">
                    <h3 class="form-section-title">Bank Account Configuration</h3>
                    <div class="form-group">
                        <label class="form-label">Bank Name</label>
                        <select class="form-select">
                            <option value="">Select Bank</option>
                            <option value="natwest">NatWest</option>
                            <option value="barclays">Barclays</option>
                            <option value="hsbc">HSBC</option>
                            <option value="lloyds">Lloyds</option>
                        </select>
                        <div class="help-text">Select your bank from the list</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Account Number</label>
                        <input type="text" class="form-input" placeholder="12345678" maxlength="8">
                        <div class="help-text">Enter your 8-digit account number</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Sort Code</label>
                        <input type="text" class="form-input" placeholder="12-34-56" maxlength="8">
                        <div class="help-text">Enter your sort code in XX-XX-XX format</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Account Holder Name</label>
                        <input type="text" class="form-input" placeholder="John Doe">
                        <div class="help-text">Enter the account holder's full name</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Open Banking Enabled</label>
                        <select class="form-select">
                            <option value="yes">Yes</option>
                            <option value="no">No</option>
                        </select>
                        <div class="help-text">Enable Open Banking for this account</div>
                    </div>
                </div>
            </div>

            <div id="creditcardsContent" class="category-content">
                <div class="form-section">
                    <h3 class="form-section-title">Credit Card Configuration</h3>
                    <div class="form-group">
                        <label class="form-label">Card Provider</label>
                        <select class="form-select">
                            <option value="">Select Provider</option>
                            <option value="visa">Visa</option>
                            <option value="mastercard">Mastercard</option>
                            <option value="amex">American Express</option>
                        </select>
                        <div class="help-text">Select your card provider</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Card Number</label>
                        <input type="text" class="form-input" placeholder="1234 5678 9012 3456" maxlength="19">
                        <div class="help-text">Enter your 16-digit card number</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Expiry Date</label>
                        <input type="text" class="form-input" placeholder="MM/YY" maxlength="5">
                        <div class="help-text">Enter expiry date in MM/YY format</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">CVV</label>
                        <input type="text" class="form-input" placeholder="123" maxlength="4">
                        <div class="help-text">Enter the 3 or 4 digit security code</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Cardholder Name</label>
                        <input type="text" class="form-input" placeholder="John Doe">
                        <div class="help-text">Enter the name as it appears on the card</div>
                    </div>
                </div>
            </div>

            <div id="communicationchannelsContent" class="category-content">
                <div class="form-section">
                    <h3 class="form-section-title">Email Configuration</h3>
                    <div class="form-group">
                        <label class="form-label">Email Provider</label>
                        <select class="form-select">
                            <option value="">Select Provider</option>
                            <option value="gmail">Gmail</option>
                            <option value="outlook">Outlook</option>
                            <option value="yahoo">Yahoo</option>
                        </select>
                        <div class="help-text">Select your email provider</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Email Address</label>
                        <input type="email" class="form-input" placeholder="user@example.com">
                        <div class="help-text">Enter your email address</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">OAuth2 Enabled</label>
                        <select class="form-select">
                            <option value="yes">Yes</option>
                            <option value="no">No</option>
                        </select>
                        <div class="help-text">Enable OAuth2 authentication</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Client ID</label>
                        <input type="text" class="form-input" placeholder="your-client-id.apps.googleusercontent.com">
                        <div class="help-text">Enter your OAuth2 client ID</div>
                    </div>
                </div>
            </div>

            <div id="securityContent" class="category-content">
                <div class="form-section">
                    <h3 class="form-section-title">Security Settings</h3>
                    <div class="form-group">
                        <label class="form-label">Two-Factor Authentication</label>
                        <select class="form-select">
                            <option value="enabled">Enabled</option>
                            <option value="disabled">Disabled</option>
                        </select>
                        <div class="help-text">Enable 2FA for enhanced security</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Session Timeout</label>
                        <select class="form-select">
                            <option value="15">15 minutes</option>
                            <option value="30">30 minutes</option>
                            <option value="60">1 hour</option>
                        </select>
                        <div class="help-text">Set session timeout duration</div>
                    </div>
                    <div class="form-group">
                        <label class="form-label">Encryption Level</label>
                        <select class="form-select">
                            <option value="256">256-bit AES</option>
                            <option value="128">128-bit AES</option>
                        </select>
                        <div class="help-text">Select encryption level for data protection</div>
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Language Toggle -->
    <div class="language-toggle" onclick="switchLanguage()">DE</div>

    <!-- Floating Chat Modal -->
    <div id="chatModal" class="floating-modal" style="display: none;">
        <div class="modal-header">
            <div class="modal-title">Settings Assistant</div>
            <div class="modal-controls">
                <button class="modal-btn" onclick="closeModal('chatModal')">×</button>
            </div>
        </div>
        <div class="modal-content">
            <div class="chat-container">
                <div class="chat-messages" id="chatMessages">
                    <div class="message assistant">
                        <div class="message-content">
                            <strong>Assistant:</strong> 
                            <p>Hello! I'm your settings configuration assistant. I can help you with:</p>
                            <ul>
                                <li>Open Banking account setup and configuration</li>
                                <li>Credit card configuration and security</li>
                                <li>Email integration (Gmail, Outlook, OAuth2)</li>
                                <li>Security settings and PCI compliance</li>
                            </ul>
                            <p>Try asking me questions like:</p>
                            <div class="quick-questions">
                                <button class="quick-btn" onclick="askQuestion('How do I add an Open Banking account?')">How do I add an Open Banking account?</button>
                                <button class="quick-btn" onclick="askQuestion('How do I configure OAuth2 for Gmail?')">How do I configure OAuth2 for Gmail?</button>
                                <button class="quick-btn" onclick="askQuestion('How do I add a credit card securely?')">How do I add a credit card securely?</button>
                                <button class="quick-btn" onclick="askQuestion('What is PCI compliance?')">What is PCI compliance?</button>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="chat-input-container">
                    <input type="text" id="chatInput" placeholder="Ask me anything about settings..." onkeypress="handleChatKeyPress(event)">
                    <button onclick="sendChatMessage()">Send</button>
                </div>
            </div>
        </div>
    </div>

    <!-- Floating Upload Modal -->
    <div id="uploadModal" class="floating-modal" style="display: none;">
        <div class="modal-header">
            <div class="modal-title">Document Upload</div>
            <div class="modal-controls">
                <button class="modal-btn" onclick="closeModal('uploadModal')">×</button>
            </div>
        </div>
        <div class="modal-content">
            <div class="upload-instructions">
                <h4>Supported Documents</h4>
                <div class="document-types">
                    <div class="doc-type">
                        <div class="doc-icon">🏦</div>
                        <div class="doc-info">
                            <strong>Bank Statements</strong>
                            <p>Extract account numbers, sort codes, and bank details</p>
                        </div>
                    </div>
                    <div class="doc-type">
                        <div class="doc-icon">💳</div>
                        <div class="doc-info">
                            <strong>Credit Card Statements</strong>
                            <p>Extract card numbers, expiry dates, and provider info</p>
                        </div>
                    </div>
                    <div class="doc-type">
                        <div class="doc-icon">📧</div>
                        <div class="doc-info">
                            <strong>Email Configuration Screenshots</strong>
                            <p>Extract server settings, ports, and authentication details</p>
                        </div>
                    </div>
                </div>
            </div>
            <div class="upload-zone" id="uploadZone">
                <div style="text-align: center; padding: 40px;">
                    <div style="font-size: 48px; margin-bottom: 16px;">📁</div>
                    <div style="color: #e0e0e0; margin-bottom: 8px;">Drag & drop documents here</div>
                    <div style="color: #808080; font-size: 12px;">or click to browse files</div>
                    <input type="file" id="fileInput" multiple accept=".pdf,.jpg,.jpeg,.png,.txt,.doc,.docx" style="display: none;" onchange="handleFileUpload(this.files)">
                </div>
            </div>
            <div class="upload-results" id="uploadResults" style="display: none;">
                <h4>Extracted Information</h4>
                <div id="extractedData"></div>
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

        // Chatbot FAQ System
        const faqDatabase = {
            'how do i add an open banking account': {
                title: 'How to Add an Open Banking Account',
                content: 'Step-by-step guide to add an Open Banking account: 1. Navigate to Bank Accounts tab 2. Click Add New Account 3. Select your bank 4. Enter account details 5. Enable Open Banking 6. Enter OAuth2 credentials 7. Configure 2FA if required. Tip: Upload a bank statement to auto-fill details!'
            },
            'how do i configure oauth2 for gmail': {
                title: 'How to Configure OAuth2 for Gmail',
                content: 'Complete OAuth2 setup: 1. Go to Google Cloud Console 2. Create project and enable Gmail API 3. Create OAuth2 credentials 4. In app, select Gmail provider 5. Enable OAuth2 checkbox 6. Enter Client ID and Secret 7. Complete authorization flow 8. Configure sync settings. All tokens are encrypted!'
            },
            'how do i add a credit card securely': {
                title: 'How to Add a Credit Card Securely',
                content: 'Secure credit card setup: 1. Go to Credit Cards tab 2. Click Add Credit Card 3. Enter card number, expiry, CVV 4. Add cardholder name and address 5. Select provider and type 6. All data is PCI-compliant encrypted 7. Card numbers are tokenized 8. CVV never stored. Use camera scan for auto-extraction!'
            },
            'what is pci compliance': {
                title: 'What is PCI Compliance?',
                content: 'PCI DSS ensures secure handling of credit card data. Key requirements: Secure network, encrypted data, vulnerability management, access control, monitoring, security policy. Our implementation: 256-bit AES encryption, tokenization, no CVV storage, secure key derivation, audit logging, regular assessments.'
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
            messageDiv.className = 'message ' + sender;
            
            const contentDiv = document.createElement('div');
            contentDiv.className = 'message-content';
            
            if (sender === 'user') {
                contentDiv.innerHTML = '<strong>You:</strong> ' + content;
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
                    return '<strong>' + faq.title + '</strong><br>' + faq.content;
                }
            }
            
            // Check for general help topics
            if (lowerMessage.includes('open banking') || lowerMessage.includes('bank')) {
                return '<strong>Open Banking Help:</strong><br>I can help you with Open Banking setup. Try asking: "How do I add an Open Banking account?" or "What banks support Open Banking?"';
            }
            
            if (lowerMessage.includes('credit card') || lowerMessage.includes('card')) {
                return '<strong>Credit Card Help:</strong><br>I can help you with credit card configuration. Try asking: "How do I add a credit card securely?" or "What is PCI compliance?"';
            }
            
            if (lowerMessage.includes('email') || lowerMessage.includes('gmail') || lowerMessage.includes('outlook')) {
                return '<strong>Email Configuration Help:</strong><br>I can help you with email setup. Try asking: "How do I configure OAuth2 for Gmail?" or "How do I configure Outlook email?"';
            }
            
            // Default response
            return '<strong>I\'m here to help!</strong><br>I can assist you with Open Banking, credit cards, email integration, and security. Try asking a specific question or use the quick question buttons above!';
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
                fileName.textContent = '📄 ' + file.name;
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
            addChatMessage('assistant', '<strong>Document Upload Complete!</strong><br>I\'ve analyzed your uploaded documents and extracted information. You can now review the data and auto-fill form fields.');
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
                    return '<strong>Detected: Bank Statement</strong><br>' +
                           '✅ Account Number: 12345678<br>' +
                           '✅ Sort Code: 12-34-56<br>' +
                           '✅ Bank Name: NatWest<br>' +
                           '✅ Account Holder: John Doe<br>' +
                           '<button onclick="autoFillBankDetails()" style="margin-top: 8px; padding: 4px 8px; background: #404040; border: none; border-radius: 4px; color: #e0e0e0; cursor: pointer; font-size: 10px;">Auto-fill Bank Details</button>';
                case 'screenshot':
                    return '<strong>Detected: Configuration Screenshot</strong><br>' +
                           '🔍 Analyzing image for configuration details...<br>' +
                           '<button onclick="analyzeScreenshot()" style="margin-top: 8px; padding: 4px 8px; background: #404040; border: none; border-radius: 4px; color: #e0e0e0; cursor: pointer; font-size: 10px;">Analyze Screenshot</button>';
                default:
                    return '<strong>Unknown File Type</strong><br>' +
                           '⚠️ Unable to automatically extract information from this file type.';
            }
        }

        function autoFillBankDetails() {
            addChatMessage('assistant', '<strong>Auto-filling Bank Details</strong><br>I\'ve automatically filled in the bank account details from your uploaded statement. Please review the information and click "Save" to confirm.');
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
        function switchLanguage() {
            const toggle = document.querySelector('.language-toggle');
            if (toggle.textContent === 'DE') {
                toggle.textContent = 'EN';
            } else {
                toggle.textContent = 'DE';
            }
        }
    </script>
</body>
</html>`

	tmpl, err := template.New("settings").Parse(htmlTemplate)
	if err != nil {
		fmt.Printf("❌ Error parsing template: %v\n", err)
		os.Exit(1)
	}

	var buf strings.Builder
	err = tmpl.Execute(&buf, nil)
	if err != nil {
		fmt.Printf("❌ Error executing template: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(*outputFile, []byte(buf.String()), 0644)
	if err != nil {
		fmt.Printf("❌ Error writing HTML file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Settings form generated successfully: %s\n", *outputFile)
}



