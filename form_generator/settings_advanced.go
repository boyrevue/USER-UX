package main

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
            --bg: #0f0f0f;
            --fg: #ffffff;
            --card-bg: #1a1a1a;
            --border: #333333;
            --accent: #404040;
            --accent-green: #00ff88;
            --accent-blue: #0088ff;
            --pure-black: #000000;
            --text-muted: #888888;
            --error: #ff4444;
            --warning: #ffaa00;
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
            overflow: hidden;
        }

        .container {
            display: grid;
            grid-template-columns: 1fr 400px;
            height: 100vh;
            overflow: hidden;
        }

        .left-panel {
            background: var(--card-bg);
            border-right: 1px solid var(--border);
            padding: 0;
            overflow-y: auto;
            display: flex;
            flex-direction: column;
        }

        .panel-header {
            padding: 20px;
            border-bottom: 1px solid var(--border);
            background: var(--bg);
        }

        .panel-title {
            font-size: 18px;
            font-weight: 600;
            color: var(--fg);
            margin-bottom: 8px;
        }

        .panel-subtitle {
            font-size: 14px;
            color: var(--text-muted);
        }

        .category-tabs {
            display: flex;
            flex-direction: column;
            padding: 0;
        }

        .category-tab {
            padding: 16px 20px;
            cursor: pointer;
            border-bottom: 1px solid var(--border);
            transition: all 0.3s ease;
            font-size: 14px;
            font-weight: 500;
            display: flex;
            align-items: center;
            justify-content: space-between;
        }

        .category-tab:hover {
            background: var(--accent);
        }

        .category-tab.active {
            background: var(--accent-green);
            color: var(--pure-black);
        }

        .category-tab .tab-icon {
            font-size: 16px;
            margin-right: 12px;
        }

        .category-tab .tab-count {
            background: var(--accent);
            color: var(--fg);
            padding: 2px 8px;
            border-radius: 12px;
            font-size: 12px;
            font-weight: 500;
        }

        .category-tab.active .tab-count {
            background: var(--pure-black);
            color: var(--accent-green);
        }

        .right-panel {
            background: var(--bg);
            padding: 0;
            overflow-y: auto;
            display: flex;
            flex-direction: column;
        }

        .content-header {
            padding: 20px;
            border-bottom: 1px solid var(--border);
            background: var(--card-bg);
        }

        .content-title {
            font-size: 20px;
            font-weight: 600;
            color: var(--fg);
            margin-bottom: 8px;
        }

        .content-subtitle {
            font-size: 14px;
            color: var(--text-muted);
        }

        .content-body {
            padding: 20px;
            flex: 1;
        }

        .category-content {
            display: none;
        }

        .category-content.active {
            display: flex;
            flex-direction: column;
            height: 100%;
        }

        .form-section {
            background: var(--card-bg);
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 20px;
            margin-bottom: 16px;
        }

        .form-section-title {
            font-size: 16px;
            font-weight: 600;
            margin-bottom: 16px;
            color: var(--fg);
            display: flex;
            align-items: center;
        }

        .form-section-title::before {
            content: "📋";
            margin-right: 8px;
            font-size: 18px;
        }

        .form-group {
            margin-bottom: 16px;
            position: relative;
        }

        .form-label {
            display: flex;
            align-items: center;
            justify-content: space-between;
            margin-bottom: 8px;
            font-weight: 500;
            color: var(--fg);
            font-size: 14px;
        }

        .form-label .status {
            font-size: 12px;
            padding: 2px 8px;
            border-radius: 12px;
            font-weight: 500;
        }

        .form-label .status.required {
            background: var(--accent-green);
            color: var(--pure-black);
        }

        .form-label .status.optional {
            background: var(--accent);
            color: var(--fg);
        }

        .form-label .status.missing {
            background: var(--error);
            color: var(--fg);
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
            background: #000000;
            color: #ffffff;
            font-size: 14px;
            cursor: pointer;
        }

        .form-select option {
            background: #000000;
            color: #ffffff;
            padding: 8px;
        }

        .form-select:focus {
            outline: none;
            border-color: var(--accent-green);
            box-shadow: 0 0 0 2px rgba(0, 255, 136, 0.2);
        }

        /* Instance Tabs */
        .instance-tabs {
            display: flex;
            border-bottom: 1px solid var(--border);
            margin-bottom: 20px;
            overflow-x: auto;
        }

        .instance-tab {
            padding: 12px 16px;
            background: none;
            border: none;
            color: var(--fg);
            cursor: pointer;
            border-bottom: 2px solid transparent;
            transition: all 0.3s ease;
            font-size: 13px;
            white-space: nowrap;
            position: relative;
        }

        .instance-tab.active {
            border-bottom-color: var(--accent-green);
            color: var(--accent-green);
        }

        .instance-tab:hover {
            background: var(--accent);
        }

        .instance-tab .delete-btn {
            position: absolute;
            top: 2px;
            right: 2px;
            width: 16px;
            height: 16px;
            background: #ff4444;
            border: none;
            border-radius: 50%;
            color: white;
            font-size: 10px;
            cursor: pointer;
            display: none;
        }

        .instance-tab:hover .delete-btn {
            display: block;
        }

        .instance-tab-content {
            display: none;
        }

        .instance-tab-content.active {
            display: block;
        }

        .add-instance-btn {
            padding: 12px 16px;
            background: var(--accent-green);
            color: var(--pure-black);
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-weight: 600;
            font-size: 13px;
            margin-left: 8px;
            transition: all 0.3s ease;
        }

        .add-instance-btn:hover {
            background: #00cc6a;
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

        /* Upload Tabs */
        .upload-tabs {
            display: flex;
            border-bottom: 1px solid var(--border);
            margin-bottom: 20px;
        }

        .upload-tab {
            padding: 12px 20px;
            background: none;
            border: none;
            color: var(--fg);
            cursor: pointer;
            border-bottom: 2px solid transparent;
            transition: all 0.3s ease;
            font-size: 14px;
        }

        .upload-tab.active {
            border-bottom-color: var(--accent-green);
            color: var(--accent-green);
        }

        .upload-tab:hover {
            background: var(--accent);
        }

        .upload-tab-content {
            display: none;
        }

        .upload-tab-content.active {
            display: block;
        }

        /* Document History */
        .document-history {
            height: 100%;
        }

        .history-list {
            max-height: 400px;
            overflow-y: auto;
            padding: 16px;
        }

        .history-item {
            background: var(--card-bg);
            border: 1px solid var(--border);
            border-radius: 8px;
            padding: 16px;
            margin-bottom: 12px;
        }

        .history-item-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 12px;
        }

        .history-item-title {
            font-weight: 600;
            color: var(--fg);
            font-size: 14px;
        }

        .history-item-time {
            color: #808080;
            font-size: 12px;
        }

        .history-item-type {
            background: var(--accent);
            padding: 4px 8px;
            border-radius: 4px;
            font-size: 11px;
            color: var(--fg);
        }

        .history-item-data {
            background: #1a1a1a;
            border: 1px solid var(--border);
            border-radius: 4px;
            padding: 12px;
            margin-top: 8px;
        }

        .history-item-data h5 {
            margin-bottom: 8px;
            color: var(--accent-green);
            font-size: 12px;
        }

        .history-item-data ul {
            list-style: none;
            padding: 0;
            margin: 0;
        }

        .history-item-data li {
            padding: 4px 0;
            font-size: 12px;
            color: var(--fg);
        }

        .history-item-data strong {
            color: #808080;
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

        /* Credit Card Display Styles */
        .credit-card-display {
            margin-bottom: 30px;
        }

        .card-visual {
            width: 100%;
            max-width: 400px;
            height: 250px;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            border-radius: 15px;
            padding: 20px;
            position: relative;
            color: white;
            font-family: 'Courier New', monospace;
            box-shadow: 0 10px 30px rgba(0, 0, 0, 0.3);
        }

        .card-chip {
            width: 50px;
            height: 40px;
            background: linear-gradient(45deg, #ffd700, #ffed4e);
            border-radius: 8px;
            margin-bottom: 20px;
        }

        .card-number {
            font-size: 24px;
            font-weight: bold;
            letter-spacing: 2px;
            margin-bottom: 20px;
        }

        .card-details {
            display: flex;
            justify-content: space-between;
            align-items: flex-end;
        }

        .cardholder-name {
            font-size: 14px;
            text-transform: uppercase;
        }

        .expiry-date {
            font-size: 14px;
        }

        .card-logo {
            position: absolute;
            top: 20px;
            right: 20px;
            font-size: 18px;
            font-weight: bold;
        }

        /* Wallet Integration */
        .wallet-integration {
            margin-top: 20px;
            padding: 20px;
            background: #333333;
            border-radius: 8px;
        }

        .wallet-integration h4 {
            margin-bottom: 15px;
            color: var(--fg);
        }

        .wallet-buttons {
            display: flex;
            gap: 10px;
            flex-wrap: wrap;
        }

        .wallet-btn {
            padding: 8px 16px;
            background: var(--accent);
            border: none;
            border-radius: 6px;
            color: var(--fg);
            cursor: pointer;
            font-size: 12px;
            transition: all 0.3s ease;
        }

        .wallet-btn:hover {
            background: #505050;
        }

        .camera-scan {
            margin-top: 20px;
        }

        .scan-btn {
            width: 100%;
            padding: 12px;
            background: var(--accent-green);
            color: var(--pure-black);
            border: none;
            border-radius: 6px;
            cursor: pointer;
            font-weight: 600;
            transition: all 0.3s ease;
        }

        .scan-btn:hover {
            background: #00cc6a;
        }

        /* Communication Channel Tabs */
        .comm-tabs {
            display: flex;
            gap: 2px;
            margin-bottom: 20px;
            background: var(--border);
            border-radius: 6px;
            padding: 4px;
        }

        .comm-tab {
            flex: 1;
            padding: 10px 16px;
            background: transparent;
            border: none;
            border-radius: 4px;
            color: var(--fg);
            cursor: pointer;
            font-size: 13px;
            transition: all 0.3s ease;
        }

        .comm-tab.active {
            background: var(--accent-green);
            color: var(--pure-black);
        }

        .comm-content {
            display: none;
        }

        .comm-content.active {
            display: block;
        }

        /* Voice Platform Tabs */
        .voice-tabs {
            display: flex;
            gap: 2px;
            margin-bottom: 20px;
            background: var(--border);
            border-radius: 6px;
            padding: 4px;
        }

        .voice-tab {
            flex: 1;
            padding: 8px 12px;
            background: transparent;
            border: none;
            border-radius: 4px;
            color: var(--fg);
            cursor: pointer;
            font-size: 11px;
            transition: all 0.3s ease;
        }

        .voice-tab.active {
            background: var(--accent-green);
            color: var(--pure-black);
        }

        .voice-content {
            display: none;
        }

        .voice-content.active {
            display: block;
        }

        /* OAuth2 and Sync Sections */
        .oauth2-section {
            display: none;
        }

        .sync-section {
            display: none;
        }

        .custom-server-section {
            display: none;
        }

        /* Camera Modal Styles */
        #cameraModal {
            top: 50px;
            left: 50px;
            width: 500px;
            height: 600px;
        }

        .camera-container {
            display: flex;
            flex-direction: column;
            height: 100%;
        }

        .camera-preview {
            flex: 1;
            background: #000;
            border-radius: 8px;
            margin-bottom: 16px;
            position: relative;
            overflow: hidden;
            display: flex;
            align-items: center;
            justify-content: center;
        }

        .camera-placeholder {
            text-align: center;
            color: #808080;
        }

        .camera-overlay {
            position: absolute;
            top: 0;
            left: 0;
            right: 0;
            bottom: 0;
            pointer-events: none;
        }

        .scan-frame {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 280px;
            height: 180px;
            border: 2px solid var(--accent-green);
            border-radius: 12px;
            box-shadow: 0 0 0 9999px rgba(0, 0, 0, 0.5);
        }

        .scan-instructions {
            position: absolute;
            bottom: 20px;
            left: 50%;
            transform: translateX(-50%);
            color: white;
            font-size: 14px;
            text-align: center;
            background: rgba(0, 0, 0, 0.7);
            padding: 8px 16px;
            border-radius: 20px;
        }

        .camera-controls {
            display: flex;
            gap: 10px;
            margin-bottom: 16px;
        }

        .camera-btn {
            flex: 1;
            padding: 10px 16px;
            background: var(--accent);
            border: none;
            border-radius: 6px;
            color: var(--fg);
            cursor: pointer;
            font-size: 13px;
            transition: all 0.3s ease;
        }

        .camera-btn:hover {
            background: #505050;
        }

        .camera-btn.capture {
            background: var(--accent-green);
            color: var(--pure-black);
        }

        .camera-btn.capture:hover {
            background: #00cc6a;
        }

        .camera-btn.secondary {
            background: #333333;
        }

        .scan-results {
            background: var(--card-bg);
            border: 1px solid var(--border);
            border-radius: 6px;
            padding: 16px;
        }

        .scan-results h4 {
            color: var(--fg);
            font-size: 14px;
            margin-bottom: 12px;
            font-weight: 600;
        }

        #scannedData {
            margin-bottom: 16px;
        }

        .scanned-field {
            display: flex;
            justify-content: space-between;
            align-items: center;
            padding: 8px 0;
            border-bottom: 1px solid var(--border);
        }

        .scanned-field:last-child {
            border-bottom: none;
        }

        .scanned-label {
            color: var(--fg);
            font-size: 13px;
            font-weight: 500;
        }

        .scanned-value {
            color: #808080;
            font-size: 13px;
            font-family: 'Courier New', monospace;
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
            <div class="panel-header">
                <div class="panel-title">Settings Configuration</div>
                <div class="panel-subtitle">Manage your application settings</div>
            </div>
            <div class="category-tabs">
                <div class="category-tab active" onclick="switchCategory('bankaccounts')">
                    <span><span class="tab-icon">🏦</span>Bank Accounts</span>
                    <span class="tab-count">4</span>
                </div>
                <div class="category-tab" onclick="switchCategory('creditcards')">
                    <span><span class="tab-icon">💳</span>Credit Cards</span>
                    <span class="tab-count">3</span>
                </div>
                <div class="category-tab" onclick="switchCategory('communicationchannels')">
                    <span><span class="tab-icon">📡</span>Communication Channels</span>
                    <span class="tab-count">4</span>
                </div>
                <div class="category-tab" onclick="switchCategory('security')">
                    <span><span class="tab-icon">🔒</span>Security</span>
                    <span class="tab-count">3</span>
                </div>
            </div>
        </div>

        <div class="right-panel">
            <div id="bankaccountsContent" class="category-content active">
                <div class="content-header">
                    <div class="content-title">Bank Accounts</div>
                    <div class="content-subtitle">Configure your bank accounts for Open Banking integration</div>
                </div>
                <div class="content-body">
                <div class="form-section">
                    <h3 class="form-section-title">Bank Account Configuration</h3>
                    
                    <div class="instance-tabs" id="bankTabs">
                        <button class="instance-tab active" onclick="switchBankTab(1)">Bank Account 1</button>
                        <button class="add-instance-btn" onclick="addBankAccount()">+ Add Bank</button>
                    </div>

                    <div id="bankTab1" class="instance-tab-content active">
                        <div class="form-group">
                            <label class="form-label">
                                Bank Name
                                <span class="status required">REQ</span>
                            </label>
                            <select class="form-select" name="bankName1">
                                <option value="">Select Bank</option>
                                <option value="allied_irish_bank">Allied Irish Bank (GB)</option>
                                <option value="bank_of_ireland">Bank of Ireland (UK)</option>
                                <option value="bank_of_scotland">Bank of Scotland</option>
                                <option value="barclays">Barclays</option>
                                <option value="co_operative_bank">Co-operative Bank</option>
                                <option value="first_direct">First Direct</option>
                                <option value="halifax">Halifax</option>
                                <option value="hsbc">HSBC</option>
                                <option value="lloyds_bank">Lloyds Bank</option>
                                <option value="metrobank">Metro Bank</option>
                                <option value="monzo">Monzo</option>
                                <option value="natwest">NatWest</option>
                                <option value="nationwide">Nationwide Building Society</option>
                                <option value="rbs">Royal Bank of Scotland</option>
                                <option value="santander">Santander</option>
                                <option value="starling_bank">Starling Bank</option>
                                <option value="tsb">TSB</option>
                                <option value="ulster_bank">Ulster Bank</option>
                                <option value="virgin_money">Virgin Money</option>
                                <option value="yorkshire_bank">Yorkshire Bank</option>
                            </select>
                            <div class="help-text">Select your bank from the TrueLayer Open Banking supported list</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">
                                Account Number
                                <span class="status required">REQ</span>
                            </label>
                            <input type="text" class="form-input" name="accountNumber1" placeholder="12345678" maxlength="8">
                            <div class="help-text">Enter your 8-digit account number</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">
                                Sort Code
                                <span class="status required">REQ</span>
                            </label>
                            <input type="text" class="form-input" name="sortCode1" placeholder="12-34-56" maxlength="8">
                            <div class="help-text">Enter your sort code in XX-XX-XX format</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">
                                Account Holder Name
                                <span class="status required">REQ</span>
                            </label>
                            <input type="text" class="form-input" name="accountHolderName1" placeholder="John Doe">
                            <div class="help-text">Enter the account holder's full name</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">
                                Open Banking Enabled
                                <span class="status optional">Optional</span>
                            </label>
                            <select class="form-select" name="openBankingEnabled1" onchange="toggleOpenBankingConfig(this.value, 1)">
                                <option value="yes">Yes</option>
                                <option value="no">No</option>
                            </select>
                            <div class="help-text">Enable TrueLayer Open Banking for this account</div>
                        </div>
                        
                        <div class="open-banking-config" id="openBankingConfig1" style="display: none;">
                            <h4 style="color: var(--fg); margin-bottom: 15px; font-size: 14px;">TrueLayer Open Banking Configuration</h4>
                            
                            <div class="form-group">
                                <label class="form-label">
                                    TrueLayer Client ID
                                    <span class="status missing">Missing</span>
                                </label>
                                <input type="text" class="form-input" name="truelayerClientId1" placeholder="your-truelayer-client-id">
                                <div class="help-text">Your TrueLayer application client ID</div>
                            </div>
                            
                            <div class="form-group">
                                <label class="form-label">TrueLayer Client Secret</label>
                                <input type="password" class="form-input" name="truelayerClientSecret1" placeholder="your-truelayer-client-secret">
                                <div class="help-text">Your TrueLayer application client secret</div>
                            </div>
                            
                            <div class="form-group">
                                <label class="form-label">Redirect URI</label>
                                <input type="url" class="form-input" name="truelayerRedirectUri1" placeholder="https://yourapp.com/oauth/callback">
                                <div class="help-text">OAuth2 redirect URI for TrueLayer authentication</div>
                            </div>
                            
                            <div class="form-group">
                                <label class="form-label">Environment</label>
                                <select class="form-select" name="truelayerEnvironment1">
                                    <option value="sandbox">Sandbox (Testing)</option>
                                    <option value="live">Live (Production)</option>
                                </select>
                                <div class="help-text">TrueLayer environment (sandbox for testing, live for production)</div>
                            </div>
                            
                            <div class="form-group">
                                <label class="form-label">Scopes</label>
                                <select class="form-select" name="truelayerScopes1" multiple>
                                    <option value="accounts">Accounts</option>
                                    <option value="balance">Balance</option>
                                    <option value="transactions">Transactions</option>
                                    <option value="cards">Cards</option>
                                    <option value="direct_debits">Direct Debits</option>
                                    <option value="standing_orders">Standing Orders</option>
                                </select>
                                <div class="help-text">Select the data scopes you need access to</div>
                            </div>

                            <div class="form-group">
                                <label class="form-label">Bank Username</label>
                                <input type="text" class="form-input" name="bankUsername1" placeholder="Enter bank login username">
                                <div class="help-text">Username for bank login (stored securely)</div>
                            </div>
                            <div class="form-group">
                                <label class="form-label">Bank Password</label>
                                <input type="password" class="form-input" name="bankPassword1" placeholder="Enter bank login password">
                                <div class="help-text">Password for bank login (encrypted and PCI compliant)</div>
                            </div>
                            <div class="form-group">
                                <label class="form-label">2FA Method</label>
                                <select class="form-select" name="bank2faMethod1">
                                    <option value="none">None</option>
                                    <option value="sms">SMS</option>
                                    <option value="authenticator">Authenticator App</option>
                                    <option value="hardware_token">Hardware Token</option>
                                    <option value="biometric">Biometric</option>
                                </select>
                                <div class="help-text">Two-factor authentication method for bank access</div>
                            </div>
                            <div class="form-group">
                                <label class="form-label">2FA Code/Token</label>
                                <input type="text" class="form-input" name="bank2faCode1" placeholder="Enter 2FA code or token">
                                <div class="help-text">2FA code, token, or device identifier</div>
                            </div>
                            <div class="form-group">
                                <label class="form-label">Memorable Information</label>
                                <input type="text" class="form-input" name="bankMemorableInfo1" placeholder="Enter memorable information">
                                <div class="help-text">Memorable information, security questions, or PIN</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <div id="creditcardsContent" class="category-content">
                <div class="form-section">
                    <h3 class="form-section-title">Credit Card Configuration</h3>
                    
                    <div class="instance-tabs" id="cardTabs">
                        <button class="instance-tab active" onclick="switchCardTab(1)">Credit Card 1</button>
                        <button class="add-instance-btn" onclick="addCreditCard()">+ Add Card</button>
                    </div>

                    <div id="cardTab1" class="instance-tab-content active">
                        <!-- Interactive Credit Card Display -->
                        <div class="credit-card-display">
                            <div class="card-visual">
                                <div class="card-chip"></div>
                                <div class="card-number" id="cardNumberDisplay1">**** **** **** ****</div>
                                <div class="card-details">
                                    <div class="cardholder-name" id="cardholderDisplay1">CARDHOLDER NAME</div>
                                    <div class="expiry-date" id="expiryDisplay1">MM/YY</div>
                                </div>
                                <div class="card-logo" id="cardLogo1">VISA</div>
                            </div>
                        </div>

                        <div class="form-group">
                            <label class="form-label">Card Provider</label>
                            <select class="form-select" name="cardProvider1" onchange="updateCardLogo(this.value, 1)">
                                <option value="">Select Provider</option>
                                <option value="visa">Visa</option>
                                <option value="mastercard">Mastercard</option>
                                <option value="amex">American Express</option>
                                <option value="discover">Discover</option>
                            </select>
                            <div class="help-text">Select your card provider</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">
                                Card Number
                                <span class="status required">REQ</span>
                            </label>
                            <input type="text" class="form-input" name="cardNumber1" id="cardNumber1" placeholder="1234 5678 9012 3456" maxlength="19" oninput="formatCardNumber(this.value, 1)">
                            <div class="help-text">Enter your 16-digit card number</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">
                                Expiry Date
                                <span class="status required">REQ</span>
                            </label>
                            <input type="text" class="form-input" name="expiryDate1" id="expiryDate1" placeholder="MM/YY" maxlength="5" oninput="formatExpiryDate(this.value, 1)">
                            <div class="help-text">Enter expiry date in MM/YY format</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">
                                CVV
                                <span class="status required">REQ</span>
                            </label>
                            <input type="text" class="form-input" name="cvv1" id="cvv1" placeholder="123" maxlength="4" oninput="updateCVVDisplay(this.value, 1)">
                            <div class="help-text">Enter the 3 or 4 digit security code</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">
                                Cardholder Name
                                <span class="status required">REQ</span>
                            </label>
                            <input type="text" class="form-input" name="cardholderName1" id="cardholderName1" placeholder="John Doe" oninput="updateCardholderDisplay(this.value, 1)">
                            <div class="help-text">Enter the name as it appears on the card</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">
                                Billing Address
                                <span class="status optional">Optional</span>
                            </label>
                            <textarea class="form-input" name="billingAddress1" placeholder="Enter your billing address"></textarea>
                            <div class="help-text">Enter the billing address for this card</div>
                        </div>

                        <div class="form-group">
                            <label class="form-label">Open Banking Enabled</label>
                            <select class="form-select" name="cardOpenBankingEnabled1" onchange="toggleCardOpenBankingConfig(this.value, 1)">
                                <option value="yes">Yes</option>
                                <option value="no">No</option>
                            </select>
                            <div class="help-text">Enable TrueLayer Open Banking for this card</div>
                        </div>
                        
                        <div class="open-banking-config" id="cardOpenBankingConfig1" style="display: none;">
                            <h4 style="color: var(--fg); margin-bottom: 15px; font-size: 14px;">TrueLayer Open Banking Configuration</h4>
                            
                            <div class="form-group">
                                <label class="form-label">TrueLayer Client ID</label>
                                <input type="text" class="form-input" name="cardTruelayerClientId1" placeholder="your-truelayer-client-id">
                                <div class="help-text">Your TrueLayer application client ID</div>
                            </div>
                            
                            <div class="form-group">
                                <label class="form-label">TrueLayer Client Secret</label>
                                <input type="password" class="form-input" name="cardTruelayerClientSecret1" placeholder="your-truelayer-client-secret">
                                <div class="help-text">Your TrueLayer application client secret</div>
                            </div>
                            
                            <div class="form-group">
                                <label class="form-label">Redirect URI</label>
                                <input type="url" class="form-input" name="cardTruelayerRedirectUri1" placeholder="https://yourapp.com/oauth/callback">
                                <div class="help-text">OAuth2 redirect URI for TrueLayer authentication</div>
                            </div>
                            
                            <div class="form-group">
                                <label class="form-label">Environment</label>
                                <select class="form-select" name="cardTruelayerEnvironment1">
                                    <option value="sandbox">Sandbox (Testing)</option>
                                    <option value="live">Live (Production)</option>
                                </select>
                                <div class="help-text">TrueLayer environment (sandbox for testing, live for production)</div>
                            </div>
                            
                            <div class="form-group">
                                <label class="form-label">Scopes</label>
                                <select class="form-select" name="cardTruelayerScopes1" multiple>
                                    <option value="accounts">Accounts</option>
                                    <option value="balance">Balance</option>
                                    <option value="transactions">Transactions</option>
                                    <option value="cards">Cards</option>
                                    <option value="direct_debits">Direct Debits</option>
                                    <option value="standing_orders">Standing Orders</option>
                                </select>
                                <div class="help-text">Select the data scopes you need access to</div>
                            </div>

                            <div class="form-group">
                                <label class="form-label">Card Username</label>
                                <input type="text" class="form-input" name="cardUsername1" placeholder="Enter card login username">
                                <div class="help-text">Username for card login (stored securely)</div>
                            </div>
                            <div class="form-group">
                                <label class="form-label">Card Password</label>
                                <input type="password" class="form-input" name="cardPassword1" placeholder="Enter card login password">
                                <div class="help-text">Password for card login (encrypted and PCI compliant)</div>
                            </div>
                            <div class="form-group">
                                <label class="form-label">2FA Method</label>
                                <select class="form-select" name="card2faMethod1">
                                    <option value="none">None</option>
                                    <option value="sms">SMS</option>
                                    <option value="authenticator">Authenticator App</option>
                                    <option value="hardware_token">Hardware Token</option>
                                    <option value="biometric">Biometric</option>
                                </select>
                                <div class="help-text">Two-factor authentication method for card access</div>
                            </div>
                            <div class="form-group">
                                <label class="form-label">2FA Code/Token</label>
                                <input type="text" class="form-input" name="card2faCode1" placeholder="Enter 2FA code or token">
                                <div class="help-text">2FA code, token, or device identifier</div>
                            </div>
                            <div class="form-group">
                                <label class="form-label">Memorable Information</label>
                                <input type="text" class="form-input" name="cardMemorableInfo1" placeholder="Enter memorable information">
                                <div class="help-text">Memorable information, security questions, or PIN</div>
                            </div>
                        </div>
                        
                        <!-- Digital Wallet Integration -->
                        <div class="wallet-integration">
                            <h4>Digital Wallet Integration</h4>
                            <div class="wallet-buttons">
                                <button class="wallet-btn" onclick="connectGoogleWallet()">Google Wallet</button>
                                <button class="wallet-btn" onclick="connectAppleWallet()">Apple Wallet</button>
                                <button class="wallet-btn" onclick="connectSamsungPay()">Samsung Pay</button>
                            </div>
                        </div>
                        
                        <!-- Camera Scanning -->
                        <div class="camera-scan">
                            <button class="scan-btn" onclick="scanCard()">📷 Scan Card with Camera</button>
                            <button class="scan-btn" onclick="uploadCardImage()" style="margin-left: 10px;">📁 Upload Card Photo</button>
                        </div>
                </div>
            </div>

            <div id="communicationchannelsContent" class="category-content">
                <!-- Communication Channel Tabs -->
                <div class="comm-tabs">
                    <button class="comm-tab active" onclick="switchCommTab('email')">Email</button>
                    <button class="comm-tab" onclick="switchCommTab('sms')">SMS</button>
                    <button class="comm-tab" onclick="switchCommTab('voice')">Voice</button>
                    <button class="comm-tab" onclick="switchCommTab('secure')">Secure Messenger</button>
                </div>

                <!-- Email Configuration -->
                <div id="emailContent" class="comm-content active">
                    <div class="form-section">
                        <h3 class="form-section-title">Email Configuration</h3>
                        
                        <div class="form-group">
                            <label class="form-label">Email Provider</label>
                            <select class="form-select" onchange="handleEmailProviderChange(this.value)">
                                <option value="">Select Provider</option>
                                <option value="gmail">Gmail</option>
                                <option value="outlook">Outlook</option>
                                <option value="yahoo">Yahoo</option>
                                <option value="custom">Custom Server</option>
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
                            <select class="form-select" onchange="handleOAuth2Toggle(this.value)">
                                <option value="yes">Yes</option>
                                <option value="no">No</option>
                            </select>
                            <div class="help-text">Enable OAuth2 authentication</div>
                        </div>
                        
                        <div class="form-group oauth2-section">
                            <label class="form-label">OAuth2 Client ID</label>
                            <input type="text" class="form-input" placeholder="your-client-id.apps.googleusercontent.com">
                            <div class="help-text">Enter your OAuth2 client ID</div>
                        </div>
                        
                        <div class="form-group oauth2-section">
                            <label class="form-label">OAuth2 Client Secret</label>
                            <input type="password" class="form-input" placeholder="your-client-secret">
                            <div class="help-text">Enter your OAuth2 client secret</div>
                        </div>
                        
                        <div class="form-group">
                            <label class="form-label">Auto Sync Enabled</label>
                            <select class="form-select" onchange="handleAutoSyncToggle(this.value)">
                                <option value="yes">Yes</option>
                                <option value="no">No</option>
                            </select>
                            <div class="help-text">Enable automatic email synchronization</div>
                        </div>
                        
                        <div class="form-group sync-section">
                            <label class="form-label">Sync Frequency</label>
                            <select class="form-select">
                                <option value="15">Every 15 minutes</option>
                                <option value="30">Every 30 minutes</option>
                                <option value="60">Every hour</option>
                            </select>
                            <div class="help-text">How often to sync emails</div>
                        </div>
                        
                        <div class="form-group custom-server-section" style="display: none;">
                            <label class="form-label">IMAP Server</label>
                            <input type="text" class="form-input" placeholder="imap.example.com">
                            <div class="help-text">Enter your IMAP server address</div>
                        </div>
                        
                        <div class="form-group custom-server-section" style="display: none;">
                            <label class="form-label">IMAP Port</label>
                            <input type="number" class="form-input" placeholder="993">
                            <div class="help-text">Enter your IMAP port (usually 993 for SSL)</div>
                        </div>
                    </div>
                </div>

                <!-- SMS Configuration -->
                <div id="smsContent" class="comm-content">
                    <div class="form-section">
                        <h3 class="form-section-title">SMS Configuration</h3>
                        <div class="form-group">
                            <label class="form-label">Mobile Number</label>
                            <input type="tel" class="form-input" placeholder="+44 123 456 7890">
                            <div class="help-text">Enter your mobile number for SMS notifications</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">SMS Provider</label>
                            <select class="form-select">
                                <option value="">Select Provider</option>
                                <option value="twilio">Twilio</option>
                                <option value="nexmo">Nexmo</option>
                                <option value="custom">Custom</option>
                            </select>
                            <div class="help-text">Select your SMS service provider</div>
                        </div>
                    </div>
                </div>

                <!-- Voice Configuration -->
                <div id="voiceContent" class="comm-content">
                    <div class="form-section">
                        <h3 class="form-section-title">Voice Configuration</h3>
                        
                        <!-- Voice Platform Tabs -->
                        <div class="voice-tabs">
                            <button class="voice-tab active" onclick="switchVoiceTab('alexa')">Alexa</button>
                            <button class="voice-tab" onclick="switchVoiceTab('google')">Google</button>
                            <button class="voice-tab" onclick="switchVoiceTab('phone')">Phone</button>
                            <button class="voice-tab" onclick="switchVoiceTab('siri')">Siri</button>
                            <button class="voice-tab" onclick="switchVoiceTab('cortana')">Cortana</button>
                        </div>

                        <div id="alexaContent" class="voice-content active">
                            <div class="form-group">
                                <label class="form-label">Alexa Skill Name</label>
                                <input type="text" class="form-input" placeholder="Insurance Assistant">
                                <div class="help-text">Name of your Alexa skill</div>
                            </div>
                            <div class="form-group">
                                <label class="form-label">Alexa Account</label>
                                <input type="email" class="form-input" placeholder="your-alexa@amazon.com">
                                <div class="help-text">Your Amazon Alexa account email</div>
                            </div>
                        </div>

                        <div id="googleContent" class="voice-content">
                            <div class="form-group">
                                <label class="form-label">Google Assistant Project</label>
                                <input type="text" class="form-input" placeholder="insurance-assistant">
                                <div class="help-text">Your Google Assistant project ID</div>
                            </div>
                        </div>

                        <div id="phoneContent" class="voice-content">
                            <div class="form-group">
                                <label class="form-label">Phone Number</label>
                                <input type="tel" class="form-input" placeholder="+44 123 456 7890">
                                <div class="help-text">Phone number for voice calls</div>
                            </div>
                        </div>

                        <div id="siriContent" class="voice-content">
                            <div class="form-group">
                                <label class="form-label">Siri Shortcuts</label>
                                <input type="text" class="form-input" placeholder="Check Insurance">
                                <div class="help-text">Siri shortcut phrase</div>
                            </div>
                        </div>

                        <div id="cortanaContent" class="voice-content">
                            <div class="form-group">
                                <label class="form-label">Cortana Skill</label>
                                <input type="text" class="form-input" placeholder="Insurance Helper">
                                <div class="help-text">Your Cortana skill name</div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Secure Messenger Configuration -->
                <div id="secureContent" class="comm-content">
                    <div class="form-section">
                        <h3 class="form-section-title">Secure Messenger Configuration</h3>
                        <div class="form-group">
                            <label class="form-label">Telegram Bot Token</label>
                            <input type="text" class="form-input" placeholder="123456789:ABCdefGHIjklMNOpqrsTUVwxyz">
                            <div class="help-text">Your Telegram bot token</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">Telegram Chat ID</label>
                            <input type="text" class="form-input" placeholder="123456789">
                            <div class="help-text">Your Telegram chat ID</div>
                        </div>
                        <div class="form-group">
                            <label class="form-label">Signal Number</label>
                            <input type="tel" class="form-input" placeholder="+44 123 456 7890">
                            <div class="help-text">Your Signal phone number</div>
                        </div>
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
            <div class="modal-title">Document Upload & History</div>
            <div class="modal-controls">
                <button class="modal-btn" onclick="closeModal('uploadModal')">×</button>
            </div>
        </div>
        <div class="modal-content">
            <div class="upload-tabs">
                <button class="upload-tab active" onclick="switchUploadTab('upload')">Upload Documents</button>
                <button class="upload-tab" onclick="switchUploadTab('history')">Document History</button>
            </div>
            
            <div id="uploadTab" class="upload-tab-content active">
                <div class="upload-instructions">
                    <h4>Supported Documents for Configuration</h4>
                    <div class="document-types">
                        <div class="doc-type">
                            <div class="doc-icon">🏦</div>
                            <div class="doc-info">
                                <strong>Bank Statements</strong>
                                <p>Extract account numbers, sort codes, bank names, and account holder details</p>
                            </div>
                        </div>
                        <div class="doc-type">
                            <div class="doc-icon">💳</div>
                            <div class="doc-info">
                                <strong>Credit Card Photos/Statements</strong>
                                <p>Extract card numbers, expiry dates, cardholder names, and provider info</p>
                            </div>
                        </div>
                        <div class="doc-type">
                            <div class="doc-icon">📧</div>
                            <div class="doc-info">
                                <strong>Email Configuration Screenshots</strong>
                                <p>Extract server settings, ports, authentication details, and OAuth2 config</p>
                            </div>
                        </div>
                        <div class="doc-type">
                            <div class="doc-icon">📱</div>
                            <div class="doc-info">
                                <strong>Mobile Banking Screenshots</strong>
                                <p>Extract app settings, login details, and Open Banking configurations</p>
                            </div>
                        </div>
                        <div class="doc-type">
                            <div class="doc-icon">🔐</div>
                            <div class="doc-info">
                                <strong>Security Setup Screenshots</strong>
                                <p>Extract 2FA settings, security questions, and authentication methods</p>
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
            </div>
            
            <div id="historyTab" class="upload-tab-content">
                <div class="document-history">
                    <h4>Uploaded Documents History</h4>
                    <div class="history-list" id="documentHistory">
                        <!-- Document history will be populated here -->
                    </div>
                </div>
            </div>
        </div>
    </div>

    <!-- Camera Scanning Modal -->
    <div id="cameraModal" class="floating-modal" style="display: none;">
        <div class="modal-header">
            <div class="modal-title">Scan Credit Card</div>
            <div class="modal-controls">
                <button class="modal-btn" onclick="closeCameraModal()">×</button>
            </div>
        </div>
        <div class="modal-content">
            <div class="camera-container">
                <div class="camera-preview" id="cameraPreview">
                    <div class="camera-placeholder">
                        <div style="font-size: 48px; margin-bottom: 16px;">📷</div>
                        <div style="color: #e0e0e0; margin-bottom: 8px;">Camera access required for card scanning</div>
                        <div style="color: #808080; font-size: 12px;">Click "Start Camera" to begin</div>
                    </div>
                </div>
                <div class="camera-overlay">
                    <div class="scan-frame"></div>
                    <div class="scan-instructions">Position your credit card within the frame</div>
                </div>
                <div class="camera-controls">
                    <button class="camera-btn" onclick="startCamera()">Start Camera</button>
                    <button class="camera-btn" onclick="uploadCardImage()">Upload Photo</button>
                    <button class="camera-btn capture" onclick="captureCard()" style="display: none;">Capture Card</button>
                    <button class="camera-btn secondary" onclick="closeCameraModal()">Cancel</button>
                </div>
                <div class="scan-results" id="scanResults" style="display: none;">
                    <h4>Scanned Card Details</h4>
                    <div id="scannedData"></div>
                    <button class="camera-btn" onclick="applyScannedData()">Apply to Form</button>
                </div>
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

        // Upload tab switching
        function switchUploadTab(tabName) {
            // Hide all tab contents
            document.querySelectorAll('.upload-tab-content').forEach(content => {
                content.classList.remove('active');
            });
            
            // Remove active from all tabs
            document.querySelectorAll('.upload-tab').forEach(tab => {
                tab.classList.remove('active');
            });
            
            // Show selected tab content
            document.getElementById(tabName + 'Tab').classList.add('active');
            
            // Add active to clicked tab
            event.target.classList.add('active');
        }

        // Document history storage
        let documentHistory = [];

        // Add document to history
        function addToDocumentHistory(file, fileType, extractedData) {
            const historyItem = {
                id: Date.now(),
                fileName: file.name,
                fileType: fileType,
                uploadTime: new Date().toLocaleString(),
                extractedData: extractedData
            };
            
            documentHistory.unshift(historyItem); // Add to beginning
            updateDocumentHistoryDisplay();
        }

        // Update document history display
        function updateDocumentHistoryDisplay() {
            const historyContainer = document.getElementById('documentHistory');
            if (!historyContainer) return;
            
            historyContainer.innerHTML = '';
            
            if (documentHistory.length === 0) {
                historyContainer.innerHTML = '<div style="text-align: center; color: #808080; padding: 40px;">No documents uploaded yet</div>';
                return;
            }
            
            documentHistory.forEach(item => {
                const historyDiv = document.createElement('div');
                historyDiv.className = 'history-item';
                historyDiv.innerHTML = 
                    '<div class="history-item-header">' +
                        '<div class="history-item-title">' + item.fileName + '</div>' +
                        '<div class="history-item-time">' + item.uploadTime + '</div>' +
                    '</div>' +
                    '<div class="history-item-type">' + item.fileType + '</div>' +
                    '<div class="history-item-data">' +
                        '<h5>Extracted Information:</h5>' +
                        '<ul>' +
                            Object.entries(item.extractedData).map(([key, value]) => 
                                '<li><strong>' + key + ':</strong> ' + value + '</li>'
                            ).join('') +
                        '</ul>' +
                    '</div>';
                historyContainer.appendChild(historyDiv);
            });
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

        // Enhanced Document Upload Functions
        function handleFileUpload(files) {
            if (!files || files.length === 0) return;
            
            Array.from(files).forEach(file => {
                const fileType = getFileType(file.name);
                const extractedData = simulateDocumentExtraction(fileType, file.name);
                
                // Add to document history
                addToDocumentHistory(file, fileType, extractedData);
                
                // Add to chat history
                addChatMessage('system', '📄 Document uploaded: ' + file.name + ' (' + fileType + ')');
                addChatMessage('assistant', 'I\'ve extracted the following information from your ' + fileType + ':\n' + Object.entries(extractedData).map(([key, value]) => '• ' + key + ': ' + value).join('\n'));
                
                // Auto-create forms based on document type
                autoCreateFormFromDocument(fileType, extractedData);
                
                // Auto-fill empty form fields
                autoFillEmptyFields(extractedData);
            });
            
            // Add message to chat
            addChatMessage('assistant', '<strong>Document Upload Complete!</strong><br>I\'ve analyzed your uploaded documents and extracted information. I\'ve also created new forms where appropriate.');
        }

        function autoCreateFormFromDocument(fileType, extractedData) {
            switch (fileType) {
                case 'credit_card':
                    // Auto-add new credit card
                    if (creditCardCount < maxCreditCards) {
                        addCreditCard();
                        addChatMessage('assistant', '💳 <strong>New Credit Card Added!</strong><br>I\'ve automatically created a new credit card form based on your uploaded card image.');
                    } else {
                        addChatMessage('assistant', '⚠️ <strong>Maximum Credit Cards Reached</strong><br>You\'ve reached the maximum of ' + maxCreditCards + ' credit cards. Please delete an existing card first.');
                    }
                    break;
                    
                case 'bank_statement':
                    // Auto-add new bank account
                    if (bankAccountCount < maxBankAccounts) {
                        addBankAccount();
                        addChatMessage('assistant', '🏦 <strong>New Bank Account Added!</strong><br>I\'ve automatically created a new bank account form based on your uploaded statement.');
                    } else {
                        addChatMessage('assistant', '⚠️ <strong>Maximum Bank Accounts Reached</strong><br>You\'ve reached the maximum of ' + maxBankAccounts + ' bank accounts. Please delete an existing account first.');
                    }
                    break;
                    
                case 'screenshot':
                    addChatMessage('assistant', '📱 <strong>Configuration Screenshot Processed</strong><br>I\'ve extracted configuration details from your screenshot. Check the form fields for auto-filled information.');
                    break;
                    
                default:
                    addChatMessage('assistant', '📄 <strong>Document Processed</strong><br>I\'ve extracted information from your document. Check the form fields for auto-filled information.');
            }
        }

        function getFileType(fileName) {
            const ext = fileName.toLowerCase().split('.').pop();
            const name = fileName.toLowerCase();
            
            // Check for card-related keywords in filename
            if (['jpg', 'jpeg', 'png'].includes(ext) && 
                (name.includes('card') || name.includes('credit') || name.includes('visa') || 
                 name.includes('mastercard') || name.includes('amex') || name.includes('discover'))) {
                return 'credit_card';
            }
            
            // Check for bank-related keywords
            if (['pdf', 'jpg', 'jpeg', 'png'].includes(ext) && 
                (name.includes('bank') || name.includes('statement') || name.includes('account') || 
                 name.includes('natwest') || name.includes('barclays') || name.includes('hsbc'))) {
                return 'bank_statement';
            }
            
            // Default file type detection
            if (ext === 'pdf') return 'bank_statement';
            if (['jpg', 'jpeg', 'png'].includes(ext)) return 'screenshot';
            if (ext === 'txt') return 'config_file';
            return 'unknown';
        }

        function simulateDocumentExtraction(fileType, fileName) {
            switch (fileType) {
                case 'credit_card':
                    const cardTypes = ['visa', 'mastercard', 'amex', 'discover'];
                    const randomType = cardTypes[Math.floor(Math.random() * cardTypes.length)];
                    
                    switch (randomType) {
                        case 'visa':
                            return {
                                'Document Type': 'Credit Card',
                                'Card Number': '4532 1234 5678 9012',
                                'Expiry Date': '12/25',
                                'Cardholder Name': 'JOHN DOE',
                                'Card Type': 'Visa',
                                'Card Provider': 'visa'
                            };
                        case 'mastercard':
                            return {
                                'Document Type': 'Credit Card',
                                'Card Number': '5555 4444 3333 2222',
                                'Expiry Date': '09/26',
                                'Cardholder Name': 'JANE SMITH',
                                'Card Type': 'Mastercard',
                                'Card Provider': 'mastercard'
                            };
                        case 'amex':
                            return {
                                'Document Type': 'Credit Card',
                                'Card Number': '3782 822463 10005',
                                'Expiry Date': '08/25',
                                'Cardholder Name': 'MIKE JOHNSON',
                                'Card Type': 'American Express',
                                'Card Provider': 'amex'
                            };
                        case 'discover':
                            return {
                                'Document Type': 'Credit Card',
                                'Card Number': '6011 1111 1111 1117',
                                'Expiry Date': '03/27',
                                'Cardholder Name': 'SARAH WILSON',
                                'Card Type': 'Discover',
                                'Card Provider': 'discover'
                            };
                    }
                    break;
                    
                case 'bank_statement':
                    const banks = ['NatWest', 'Barclays', 'HSBC', 'Lloyds', 'Santander'];
                    const randomBank = banks[Math.floor(Math.random() * banks.length)];
                    
                    return {
                        'Document Type': 'Bank Statement',
                        'Account Number': '12345678',
                        'Sort Code': '12-34-56',
                        'Bank Name': randomBank,
                        'Account Holder': 'John Doe',
                        'Statement Date': '2024-01-15'
                    };
                    
                case 'screenshot':
                    return {
                        'Document Type': 'Configuration Screenshot',
                        'Email Provider': 'Gmail',
                        'Server Type': 'IMAP',
                        'Server Address': 'imap.gmail.com',
                        'Port': '993',
                        'Security': 'SSL/TLS',
                        'Authentication': 'OAuth2'
                    };
                    
                case 'config_file':
                    return {
                        'Document Type': 'Configuration File',
                        'File Format': 'Text Configuration',
                        'Settings Detected': 'Email, Banking, Security',
                        'Configuration Type': 'Application Settings'
                    };
                    
                default:
                    return {
                        'Document Type': 'Unknown File Type',
                        'Status': 'Unable to automatically extract information',
                        'Recommendation': 'Please upload a supported document type'
                    };
            }
        }

        function autoFillBankDetails() {
            addChatMessage('assistant', '<strong>Auto-filling Bank Details</strong><br>I\'ve automatically filled in the bank account details from your uploaded statement. Please review the information and click "Save" to confirm.');
        }

        function autoFillEmptyFields(extractedData) {
            // Map extracted data to form fields
            const fieldMappings = {
                'Account Number': 'accountNumber',
                'Sort Code': 'sortCode',
                'Bank Name': 'bankName',
                'Account Holder': 'accountHolderName',
                'Card Number': 'cardNumber',
                'Expiry Date': 'expiryDate',
                'Cardholder Name': 'cardholderName',
                'Card Type': 'cardProvider',
                'Email Provider': 'emailProvider',
                'Server Address': 'imapServer',
                'Port': 'imapPort',
                'Security': 'imapUseSSL'
            };
            
            Object.entries(extractedData).forEach(([key, value]) => {
                const fieldName = fieldMappings[key];
                if (fieldName) {
                    // Handle both bank and card fields
                    if (fieldName.startsWith('card')) {
                        // Find the currently active card tab
                        const activeCardTab = document.querySelector('#cardTabs .instance-tab.active');
                        if (activeCardTab) {
                            const tabText = activeCardTab.textContent;
                            const cardNumber = tabText.match(/Credit Card (\d+)/)[1];
                            const field = document.querySelector('[name="' + fieldName + cardNumber + '"]');
                            if (field && !field.value) {
                                field.value = value;
                                field.style.borderColor = '#00ff88';
                                setTimeout(() => {
                                    field.style.borderColor = '';
                                }, 2000);
                                
                                // Update visual elements for cards
                                if (fieldName === 'cardNumber') {
                                    formatCardNumber(value, cardNumber);
                                } else if (fieldName === 'expiryDate') {
                                    formatExpiryDate(value, cardNumber);
                                } else if (fieldName === 'cardholderName') {
                                    updateCardholderDisplay(value, cardNumber);
                                } else if (fieldName === 'cardProvider') {
                                    updateCardLogo(value.toLowerCase(), cardNumber);
                                }
                            }
                        }
                    } else if (fieldName.startsWith('bank') || fieldName.startsWith('account')) {
                        // Find the currently active bank tab
                        const activeBankTab = document.querySelector('#bankTabs .instance-tab.active');
                        if (activeBankTab) {
                            const tabText = activeBankTab.textContent;
                            const bankNumber = tabText.match(/Bank Account (\d+)/)[1];
                            const field = document.querySelector('[name="' + fieldName + bankNumber + '"]');
                            if (field && !field.value) {
                                field.value = value;
                                field.style.borderColor = '#00ff88';
                                setTimeout(() => {
                                    field.style.borderColor = '';
                                }, 2000);
                            }
                        }
                    } else {
                        // Handle other fields (email, etc.)
                        const field = document.querySelector('[name="' + fieldName + '"]');
                        if (field && !field.value) {
                            field.value = value;
                            field.style.borderColor = '#00ff88';
                            setTimeout(() => {
                                field.style.borderColor = '';
                            }, 2000);
                        }
                    }
                }
            });
            
            addChatMessage('assistant', '<strong>Auto-fill Complete!</strong><br>I\'ve automatically filled empty form fields with data from your uploaded document. Please review and save the changes.');
        }

        // Initialize upload zone click handler and modal functionality
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

            // Initialize draggable modals
            initializeDraggableModals();
        });

        // Make modals draggable and resizable
        function initializeDraggableModals() {
            const modals = document.querySelectorAll('.floating-modal');
            
            modals.forEach(modal => {
                const header = modal.querySelector('.modal-header');
                let isDragging = false;
                let currentX;
                let currentY;
                let initialX;
                let initialY;
                let xOffset = 0;
                let yOffset = 0;

                // Dragging functionality
                header.addEventListener('mousedown', function(e) {
                    if (e.target.classList.contains('modal-btn')) return; // Don't drag when clicking close button
                    
                    isDragging = true;
                    initialX = e.clientX - xOffset;
                    initialY = e.clientY - yOffset;
                    header.style.cursor = 'grabbing';
                });

                document.addEventListener('mousemove', function(e) {
                    if (isDragging) {
                        e.preventDefault();
                        currentX = e.clientX - initialX;
                        currentY = e.clientY - initialY;
                        xOffset = currentX;
                        yOffset = currentY;

                        modal.style.transform = 'translate(' + currentX + 'px, ' + currentY + 'px)';
                    }
                });

                document.addEventListener('mouseup', function() {
                    if (isDragging) {
                        isDragging = false;
                        header.style.cursor = 'grab';
                    }
                });

                // Resize functionality
                let isResizing = false;
                let startWidth, startHeight, startX, startY;

                modal.addEventListener('mousedown', function(e) {
                    const rect = modal.getBoundingClientRect();
                    const isNearRight = e.clientX > rect.right - 10;
                    const isNearBottom = e.clientY > rect.bottom - 10;
                    
                    if (isNearRight && isNearBottom) {
                        isResizing = true;
                        startWidth = rect.width;
                        startHeight = rect.height;
                        startX = e.clientX;
                        startY = e.clientY;
                        e.preventDefault();
                    }
                });

                document.addEventListener('mousemove', function(e) {
                    if (isResizing) {
                        const width = startWidth + (e.clientX - startX);
                        const height = startHeight + (e.clientY - startY);
                        
                        if (width > 300 && height > 200) {
                            modal.style.width = width + 'px';
                            modal.style.height = height + 'px';
                        }
                    }
                });

                document.addEventListener('mouseup', function() {
                    isResizing = false;
                });
            });
        }

        // Credit Card Functions
        function formatCardNumber(value, cardNumber = 1) {
            const cleaned = value.replace(/\s/g, '');
            const groups = cleaned.match(/.{1,4}/g);
            const formatted = groups ? groups.join(' ') : cleaned;
            document.getElementById('cardNumber' + cardNumber).value = formatted;
            document.getElementById('cardNumberDisplay' + cardNumber).textContent = formatted || '**** **** **** ****';
        }

        function formatExpiryDate(value, cardNumber = 1) {
            const cleaned = value.replace(/\D/g, '');
            if (cleaned.length >= 2) {
                const formatted = cleaned.slice(0, 2) + '/' + cleaned.slice(2, 4);
                document.getElementById('expiryDate' + cardNumber).value = formatted;
                document.getElementById('expiryDisplay' + cardNumber).textContent = formatted;
            } else {
                document.getElementById('expiryDisplay' + cardNumber).textContent = 'MM/YY';
            }
        }

        function updateCardholderDisplay(value, cardNumber = 1) {
            document.getElementById('cardholderDisplay' + cardNumber).textContent = value.toUpperCase() || 'CARDHOLDER NAME';
        }

        function updateCVVDisplay(value, cardNumber = 1) {
            // CVV is not displayed on the card for security
        }

        function updateCardLogo(provider, cardNumber = 1) {
            const logo = document.getElementById('cardLogo' + cardNumber);
            switch (provider) {
                case 'visa':
                    logo.textContent = 'VISA';
                    break;
                case 'mastercard':
                    logo.textContent = 'MASTERCARD';
                    break;
                case 'amex':
                    logo.textContent = 'AMEX';
                    break;
                case 'discover':
                    logo.textContent = 'DISCOVER';
                    break;
                default:
                    logo.textContent = 'VISA';
            }
        }

        // Digital Wallet Functions
        function connectGoogleWallet() {
            alert('Connecting to Google Wallet...');
        }

        function connectAppleWallet() {
            alert('Connecting to Apple Wallet...');
        }

        function connectSamsungPay() {
            alert('Connecting to Samsung Pay...');
        }

        // Camera scanning variables
        let stream = null;
        let video = null;
        let canvas = null;
        let capturedImage = null;

        function scanCard() {
            toggleModal('cameraModal');
        }

        function closeCameraModal() {
            if (stream) {
                stream.getTracks().forEach(track => track.stop());
                stream = null;
            }
            if (video) {
                video.srcObject = null;
                video = null;
            }
            toggleModal('cameraModal');
        }

        async function startCamera() {
            try {
                // Request camera access
                stream = await navigator.mediaDevices.getUserMedia({ 
                    video: { 
                        facingMode: 'environment',
                        width: { ideal: 1280 },
                        height: { ideal: 720 }
                    } 
                });

                // Create video element
                const cameraPreview = document.getElementById('cameraPreview');
                cameraPreview.innerHTML = '';
                
                video = document.createElement('video');
                video.autoplay = true;
                video.playsInline = true;
                video.style.width = '100%';
                video.style.height = '100%';
                video.style.objectFit = 'cover';
                
                video.srcObject = stream;
                cameraPreview.appendChild(video);

                // Add overlay
                const overlay = document.createElement('div');
                overlay.className = 'camera-overlay';
                overlay.innerHTML = '<div class="scan-frame"></div><div class="scan-instructions">Position your credit card within the frame</div>';
                cameraPreview.appendChild(overlay);

                // Show capture button
                document.querySelector('.camera-btn.capture').style.display = 'block';
                
            } catch (error) {
                console.error('Error accessing camera:', error);
                alert('Unable to access camera. Please check permissions and try again.');
            }
        }

        function captureCard() {
            if (!video || !stream) {
                alert('Please start the camera first');
                return;
            }

            try {
                // Create canvas to capture frame
                canvas = document.createElement('canvas');
                canvas.width = video.videoWidth;
                canvas.height = video.videoHeight;
                
                const ctx = canvas.getContext('2d');
                ctx.drawImage(video, 0, 0);
                
                // Get image data
                capturedImage = canvas.toDataURL('image/jpeg');
                
                // Stop camera
                stream.getTracks().forEach(track => track.stop());
                stream = null;
                
                // Show captured image
                const cameraPreview = document.getElementById('cameraPreview');
                cameraPreview.innerHTML = '<img src="' + capturedImage + '" style="width: 100%; height: 100%; object-fit: cover; border-radius: 8px;">';
                
                // Process the image for OCR
                processCardImage(capturedImage);
                
            } catch (error) {
                console.error('Error capturing image:', error);
                alert('Error capturing image. Please try again.');
            }
        }

        async function processCardImage(imageData) {
            try {
                // Show processing message
                const scanResults = document.getElementById('scanResults');
                const scannedData = document.getElementById('scannedData');
                scannedData.innerHTML = '<div style="text-align: center; color: #808080;">Processing card image...</div>';
                scanResults.style.display = 'block';

                // Simulate OCR processing with realistic delays
                await new Promise(resolve => setTimeout(resolve, 2000));

                // For demo purposes, we'll simulate OCR results
                // In a real implementation, you would send the image to an OCR service
                const mockOCRResults = {
                    cardNumber: '4532 1234 5678 9012',
                    expiryDate: '12/25',
                    cardholderName: 'JOHN DOE',
                    cardType: 'visa'
                };

                // Display results
                displayScannedResults(mockOCRResults);

            } catch (error) {
                console.error('Error processing card image:', error);
                scannedData.innerHTML = '<div style="color: #ff6b6b;">Error processing card image. Please try again.</div>';
            }
        }

        function displayScannedResults(results) {
            const scannedData = document.getElementById('scannedData');
            scannedData.innerHTML = '<div class="scanned-field"><span class="scanned-label">Card Number:</span><span class="scanned-value">' + results.cardNumber + '</span></div>' +
                '<div class="scanned-field"><span class="scanned-label">Expiry Date:</span><span class="scanned-value">' + results.expiryDate + '</span></div>' +
                '<div class="scanned-field"><span class="scanned-label">Cardholder Name:</span><span class="scanned-value">' + results.cardholderName + '</span></div>' +
                '<div class="scanned-field"><span class="scanned-label">Card Type:</span><span class="scanned-value">' + results.cardType.toUpperCase() + '</span></div>';
        }

        function applyScannedData() {
            // Get the currently active card tab
            const activeCardTab = document.querySelector('#cardTabs .instance-tab.active');
            if (!activeCardTab) {
                alert('Please select a credit card tab first');
                return;
            }
            
            // Extract card number from active tab
            const tabText = activeCardTab.textContent;
            const cardNumber = tabText.match(/Credit Card (\d+)/)[1];
            
            // Apply scanned data to form fields
            const scannedData = document.getElementById('scannedData');
            const fields = scannedData.querySelectorAll('.scanned-field');
            
            fields.forEach(field => {
                const label = field.querySelector('.scanned-label').textContent;
                const value = field.querySelector('.scanned-value').textContent;
                
                switch (label) {
                    case 'Card Number:':
                        document.getElementById('cardNumber' + cardNumber).value = value;
                        formatCardNumber(value, cardNumber);
                        break;
                    case 'Expiry Date:':
                        document.getElementById('expiryDate' + cardNumber).value = value;
                        formatExpiryDate(value, cardNumber);
                        break;
                    case 'Cardholder Name:':
                        document.getElementById('cardholderName' + cardNumber).value = value;
                        updateCardholderDisplay(value, cardNumber);
                        break;
                    case 'Card Type:':
                        const providerSelect = document.querySelector('select[name="cardProvider' + cardNumber + '"]');
                        providerSelect.value = value.toLowerCase();
                        updateCardLogo(value.toLowerCase(), cardNumber);
                        break;
                }
            });

            // Close camera modal
            closeCameraModal();
            
            // Show success message
            alert('Card details applied to Credit Card ' + cardNumber + ' successfully!');
        }

        // Enhanced card recognition and upload functionality
        function uploadCardImage() {
            const input = document.createElement('input');
            input.type = 'file';
            input.accept = 'image/*';
            input.onchange = function(e) {
                const file = e.target.files[0];
                if (file) {
                    processUploadedCardImage(file);
                }
            };
            input.click();
        }

        function processUploadedCardImage(file) {
            const reader = new FileReader();
            reader.onload = function(e) {
                const imageData = e.target.result;
                
                // Show processing message
                const scanResults = document.getElementById('scanResults');
                const scannedData = document.getElementById('scannedData');
                scannedData.innerHTML = '<div style="text-align: center; color: #808080;">Processing uploaded card image...</div>';
                scanResults.style.display = 'block';
                
                // Simulate OCR processing
                setTimeout(() => {
                    // Enhanced OCR simulation with different card types
                    const cardTypes = ['visa', 'mastercard', 'amex', 'discover'];
                    const randomType = cardTypes[Math.floor(Math.random() * cardTypes.length)];
                    
                    let mockOCRResults;
                    switch (randomType) {
                        case 'visa':
                            mockOCRResults = {
                                cardNumber: '4532 1234 5678 9012',
                                expiryDate: '12/25',
                                cardholderName: 'JOHN DOE',
                                cardType: 'visa',
                                confidence: 0.95
                            };
                            break;
                        case 'mastercard':
                            mockOCRResults = {
                                cardNumber: '5555 4444 3333 2222',
                                expiryDate: '09/26',
                                cardholderName: 'JANE SMITH',
                                cardType: 'mastercard',
                                confidence: 0.92
                            };
                            break;
                        case 'amex':
                            mockOCRResults = {
                                cardNumber: '3782 822463 10005',
                                expiryDate: '08/25',
                                cardholderName: 'MIKE JOHNSON',
                                cardType: 'amex',
                                confidence: 0.88
                            };
                            break;
                        case 'discover':
                            mockOCRResults = {
                                cardNumber: '6011 1111 1111 1117',
                                expiryDate: '03/27',
                                cardholderName: 'SARAH WILSON',
                                cardType: 'discover',
                                confidence: 0.90
                            };
                            break;
                    }
                    
                    displayScannedResults(mockOCRResults);
                }, 2000);
            };
            reader.readAsDataURL(file);
        }

        // Communication Channel Functions
        function switchCommTab(tabName) {
            // Hide all comm content
            const commContents = document.querySelectorAll('.comm-content');
            commContents.forEach(content => content.classList.remove('active'));
            
            // Remove active class from all comm tabs
            const commTabs = document.querySelectorAll('.comm-tab');
            commTabs.forEach(tab => tab.classList.remove('active'));
            
            // Show selected content and activate tab
            document.getElementById(tabName + 'Content').classList.add('active');
            event.target.classList.add('active');
        }

        function switchVoiceTab(tabName) {
            // Hide all voice content
            const voiceContents = document.querySelectorAll('.voice-content');
            voiceContents.forEach(content => content.classList.remove('active'));
            
            // Remove active class from all voice tabs
            const voiceTabs = document.querySelectorAll('.voice-tab');
            voiceTabs.forEach(tab => tab.classList.remove('active'));
            
            // Show selected content and activate tab
            document.getElementById(tabName + 'Content').classList.add('active');
            event.target.classList.add('active');
        }

        function handleEmailProviderChange(provider) {
            const customSections = document.querySelectorAll('.custom-server-section');
            if (provider === 'custom') {
                customSections.forEach(section => section.style.display = 'block');
            } else {
                customSections.forEach(section => section.style.display = 'none');
            }
        }

        function handleOAuth2Toggle(enabled) {
            const oauth2Sections = document.querySelectorAll('.oauth2-section');
            if (enabled === 'yes') {
                oauth2Sections.forEach(section => section.style.display = 'block');
            } else {
                oauth2Sections.forEach(section => section.style.display = 'none');
            }
        }

        function handleAutoSyncToggle(enabled) {
            const syncSections = document.querySelectorAll('.sync-section');
            if (enabled === 'yes') {
                syncSections.forEach(section => section.style.display = 'block');
            } else {
                syncSections.forEach(section => section.style.display = 'none');
            }
        }

        // Bank account management
        let bankAccountCount = 1;
        const maxBankAccounts = 8;

        function addBankAccount() {
            if (bankAccountCount >= maxBankAccounts) {
                alert('Maximum of ' + maxBankAccounts + ' bank accounts allowed');
                return;
            }
            
            bankAccountCount++;
            const bankTabs = document.getElementById('bankTabs');
            const addButton = bankTabs.querySelector('.add-instance-btn');
            
            // Create new tab
            const newTab = document.createElement('button');
            newTab.className = 'instance-tab';
            newTab.onclick = function() { switchBankTab(bankAccountCount); };
            newTab.innerHTML = 'Bank Account ' + bankAccountCount + '<button class="delete-btn" onclick="deleteBankAccount(' + bankAccountCount + ')">×</button>';
            
            // Insert before add button
            bankTabs.insertBefore(newTab, addButton);
            
            // Create new tab content
            const newTabContent = document.createElement('div');
            newTabContent.id = 'bankTab' + bankAccountCount;
            newTabContent.className = 'instance-tab-content';
            newTabContent.innerHTML = generateBankAccountHTML(bankAccountCount);
            
            // Add to container
            const container = document.getElementById('bankTab1').parentNode;
            container.appendChild(newTabContent);
            
            // Switch to new tab
            switchBankTab(bankAccountCount);
        }

        function deleteBankAccount(accountNumber) {
            if (confirm('Are you sure you want to delete Bank Account ' + accountNumber + '?')) {
                // Remove tab
                const tabs = document.querySelectorAll('#bankTabs .instance-tab');
                tabs[accountNumber - 1].remove();
                
                // Remove content
                const content = document.getElementById('bankTab' + accountNumber);
                content.remove();
                
                // Renumber remaining accounts
                renumberBankAccounts();
            }
        }

        function renumberBankAccounts() {
            const tabs = document.querySelectorAll('#bankTabs .instance-tab');
            const contents = document.querySelectorAll('[id^="bankTab"]');
            
            tabs.forEach((tab, index) => {
                const newNumber = index + 1;
                tab.innerHTML = 'Bank Account ' + newNumber + '<button class="delete-btn" onclick="deleteBankAccount(' + newNumber + ')">×</button>';
                tab.onclick = function() { switchBankTab(newNumber); };
            });
            
            contents.forEach((content, index) => {
                const newNumber = index + 1;
                content.id = 'bankTab' + newNumber;
                // Update all form field names
                const inputs = content.querySelectorAll('input, select');
                inputs.forEach(input => {
                    if (input.name) {
                        input.name = input.name.replace(/\\d+$/, newNumber);
                    }
                });
            });
            
            bankAccountCount = tabs.length;
        }

        function switchBankTab(accountNumber) {
            // Hide all tab contents
            document.querySelectorAll('[id^="bankTab"]').forEach(content => {
                content.classList.remove('active');
            });
            
            // Remove active from all tabs
            document.querySelectorAll('#bankTabs .instance-tab').forEach(tab => {
                tab.classList.remove('active');
            });
            
            // Show selected tab content
            document.getElementById('bankTab' + accountNumber).classList.add('active');
            
            // Add active to clicked tab
            event.target.classList.add('active');
        }

        function generateBankAccountHTML(accountNumber) {
            return '<div class="form-group">' +
                    '<label class="form-label">Bank Name</label>' +
                    '<select class="form-select" name="bankName' + accountNumber + '">' +
                        '<option value="">Select Bank</option>' +
                        '<option value="allied_irish_bank">Allied Irish Bank (GB)</option>' +
                        '<option value="bank_of_ireland">Bank of Ireland (UK)</option>' +
                        '<option value="bank_of_scotland">Bank of Scotland</option>' +
                        '<option value="barclays">Barclays</option>' +
                        '<option value="co_operative_bank">Co-operative Bank</option>' +
                        '<option value="first_direct">First Direct</option>' +
                        '<option value="halifax">Halifax</option>' +
                        '<option value="hsbc">HSBC</option>' +
                        '<option value="lloyds_bank">Lloyds Bank</option>' +
                        '<option value="metrobank">Metro Bank</option>' +
                        '<option value="monzo">Monzo</option>' +
                        '<option value="natwest">NatWest</option>' +
                        '<option value="nationwide">Nationwide Building Society</option>' +
                        '<option value="rbs">Royal Bank of Scotland</option>' +
                        '<option value="santander">Santander</option>' +
                        '<option value="starling_bank">Starling Bank</option>' +
                        '<option value="tsb">TSB</option>' +
                        '<option value="ulster_bank">Ulster Bank</option>' +
                        '<option value="virgin_money">Virgin Money</option>' +
                        '<option value="yorkshire_bank">Yorkshire Bank</option>' +
                    '</select>' +
                    '<div class="help-text">Select your bank from the TrueLayer Open Banking supported list</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">Account Number</label>' +
                    '<input type="text" class="form-input" name="accountNumber' + accountNumber + '" placeholder="12345678" maxlength="8">' +
                    '<div class="help-text">Enter your 8-digit account number</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">Sort Code</label>' +
                    '<input type="text" class="form-input" name="sortCode' + accountNumber + '" placeholder="12-34-56" maxlength="8">' +
                    '<div class="help-text">Enter your sort code in XX-XX-XX format</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">Account Holder Name</label>' +
                    '<input type="text" class="form-input" name="accountHolderName' + accountNumber + '" placeholder="John Doe">' +
                    '<div class="help-text">Enter the account holder\'s full name</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">Open Banking Enabled</label>' +
                    '<select class="form-select" name="openBankingEnabled' + accountNumber + '" onchange="toggleOpenBankingConfig(this.value, ' + accountNumber + ')">' +
                        '<option value="yes">Yes</option>' +
                        '<option value="no">No</option>' +
                    '</select>' +
                    '<div class="help-text">Enable TrueLayer Open Banking for this account</div>' +
                '</div>' +
                '<div class="open-banking-config" id="openBankingConfig' + accountNumber + '" style="display: none;">' +
                    '<h4 style="color: var(--fg); margin-bottom: 15px; font-size: 14px;">TrueLayer Open Banking Configuration</h4>' +
                    '<div class="form-group">' +
                        '<label class="form-label">TrueLayer Client ID</label>' +
                        '<input type="text" class="form-input" name="truelayerClientId' + accountNumber + '" placeholder="your-truelayer-client-id">' +
                        '<div class="help-text">Your TrueLayer application client ID</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">TrueLayer Client Secret</label>' +
                        '<input type="password" class="form-input" name="truelayerClientSecret' + accountNumber + '" placeholder="your-truelayer-client-secret">' +
                        '<div class="help-text">Your TrueLayer application client secret</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Redirect URI</label>' +
                        '<input type="url" class="form-input" name="truelayerRedirectUri' + accountNumber + '" placeholder="https://yourapp.com/oauth/callback">' +
                        '<div class="help-text">OAuth2 redirect URI for TrueLayer authentication</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Environment</label>' +
                        '<select class="form-select" name="truelayerEnvironment' + accountNumber + '">' +
                            '<option value="sandbox">Sandbox (Testing)</option>' +
                            '<option value="live">Live (Production)</option>' +
                        '</select>' +
                        '<div class="help-text">TrueLayer environment (sandbox for testing, live for production)</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Scopes</label>' +
                        '<select class="form-select" name="truelayerScopes' + accountNumber + '" multiple>' +
                            '<option value="accounts">Accounts</option>' +
                            '<option value="balance">Balance</option>' +
                            '<option value="transactions">Transactions</option>' +
                            '<option value="cards">Cards</option>' +
                            '<option value="direct_debits">Direct Debits</option>' +
                            '<option value="standing_orders">Standing Orders</option>' +
                        '</select>' +
                        '<div class="help-text">Select the data scopes you need access to</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Bank Username</label>' +
                        '<input type="text" class="form-input" name="bankUsername' + accountNumber + '" placeholder="Enter bank login username">' +
                        '<div class="help-text">Username for bank login (stored securely)</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Bank Password</label>' +
                        '<input type="password" class="form-input" name="bankPassword' + accountNumber + '" placeholder="Enter bank login password">' +
                        '<div class="help-text">Password for bank login (encrypted and PCI compliant)</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">2FA Method</label>' +
                        '<select class="form-select" name="bank2faMethod' + accountNumber + '">' +
                            '<option value="none">None</option>' +
                            '<option value="sms">SMS</option>' +
                            '<option value="authenticator">Authenticator App</option>' +
                            '<option value="hardware_token">Hardware Token</option>' +
                            '<option value="biometric">Biometric</option>' +
                        '</select>' +
                        '<div class="help-text">Two-factor authentication method for bank access</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">2FA Code/Token</label>' +
                        '<input type="text" class="form-input" name="bank2faCode' + accountNumber + '" placeholder="Enter 2FA code or token">' +
                        '<div class="help-text">2FA code, token, or device identifier</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Memorable Information</label>' +
                        '<input type="text" class="form-input" name="bankMemorableInfo' + accountNumber + '" placeholder="Enter memorable information">' +
                        '<div class="help-text">Memorable information, security questions, or PIN</div>' +
                    '</div>' +
                '</div>';
        }

        function toggleOpenBankingConfig(enabled, accountNumber) {
            const configDiv = document.getElementById('openBankingConfig' + accountNumber);
            if (enabled === 'yes') {
                configDiv.style.display = 'block';
            } else {
                configDiv.style.display = 'none';
            }
        }

        // Credit card management
        let creditCardCount = 1;
        const maxCreditCards = 8;

        function addCreditCard() {
            if (creditCardCount >= maxCreditCards) {
                alert('Maximum of ' + maxCreditCards + ' credit cards allowed');
                return;
            }
            
            creditCardCount++;
            const cardTabs = document.getElementById('cardTabs');
            const addButton = cardTabs.querySelector('.add-instance-btn');
            
            // Create new tab
            const newTab = document.createElement('button');
            newTab.className = 'instance-tab';
            newTab.onclick = function() { switchCardTab(creditCardCount); };
            newTab.innerHTML = 'Credit Card ' + creditCardCount + '<button class="delete-btn" onclick="deleteCreditCard(' + creditCardCount + ')">×</button>';
            
            // Insert before add button
            cardTabs.insertBefore(newTab, addButton);
            
            // Create new tab content
            const newTabContent = document.createElement('div');
            newTabContent.id = 'cardTab' + creditCardCount;
            newTabContent.className = 'instance-tab-content';
            newTabContent.innerHTML = generateCreditCardHTML(creditCardCount);
            
            // Add to container
            const container = document.getElementById('cardTab1').parentNode;
            container.appendChild(newTabContent);
            
            // Switch to new tab
            switchCardTab(creditCardCount);
        }

        function deleteCreditCard(cardNumber) {
            if (confirm('Are you sure you want to delete Credit Card ' + cardNumber + '?')) {
                // Remove tab
                const tabs = document.querySelectorAll('#cardTabs .instance-tab');
                tabs[cardNumber - 1].remove();
                
                // Remove content
                const content = document.getElementById('cardTab' + cardNumber);
                content.remove();
                
                // Renumber remaining cards
                renumberCreditCards();
            }
        }

        function renumberCreditCards() {
            const tabs = document.querySelectorAll('#cardTabs .instance-tab');
            const contents = document.querySelectorAll('[id^="cardTab"]');
            
            tabs.forEach((tab, index) => {
                const newNumber = index + 1;
                tab.innerHTML = 'Credit Card ' + newNumber + '<button class="delete-btn" onclick="deleteCreditCard(' + newNumber + ')">×</button>';
                tab.onclick = function() { switchCardTab(newNumber); };
            });
            
            contents.forEach((content, index) => {
                const newNumber = index + 1;
                content.id = 'cardTab' + newNumber;
                // Update all form field names
                const inputs = content.querySelectorAll('input, select, textarea');
                inputs.forEach(input => {
                    if (input.name) {
                        input.name = input.name.replace(/\\d+$/, newNumber);
                    }
                    if (input.id) {
                        input.id = input.id.replace(/\\d+$/, newNumber);
                    }
                });
                // Update display elements
                const displays = content.querySelectorAll('[id^="cardNumberDisplay"], [id^="cardholderDisplay"], [id^="expiryDisplay"], [id^="cardLogo"]');
                displays.forEach(display => {
                    display.id = display.id.replace(/\\d+$/, newNumber);
                });
            });
            
            creditCardCount = tabs.length;
        }

        function switchCardTab(cardNumber) {
            // Hide all tab contents
            document.querySelectorAll('[id^="cardTab"]').forEach(content => {
                content.classList.remove('active');
            });
            
            // Remove active from all tabs
            document.querySelectorAll('#cardTabs .instance-tab').forEach(tab => {
                tab.classList.remove('active');
            });
            
            // Show selected tab content
            document.getElementById('cardTab' + cardNumber).classList.add('active');
            
            // Add active to clicked tab
            event.target.classList.add('active');
        }

        function generateCreditCardHTML(cardNumber) {
            return '<!-- Interactive Credit Card Display -->' +
                '<div class="credit-card-display">' +
                    '<div class="card-visual">' +
                        '<div class="card-chip"></div>' +
                        '<div class="card-number" id="cardNumberDisplay' + cardNumber + '">**** **** **** ****</div>' +
                        '<div class="card-details">' +
                            '<div class="cardholder-name" id="cardholderDisplay' + cardNumber + '">CARDHOLDER NAME</div>' +
                            '<div class="expiry-date" id="expiryDisplay' + cardNumber + '">MM/YY</div>' +
                        '</div>' +
                        '<div class="card-logo" id="cardLogo' + cardNumber + '">VISA</div>' +
                    '</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">Card Provider</label>' +
                    '<select class="form-select" name="cardProvider' + cardNumber + '" onchange="updateCardLogo(this.value, ' + cardNumber + ')">' +
                        '<option value="">Select Provider</option>' +
                        '<option value="visa">Visa</option>' +
                        '<option value="mastercard">Mastercard</option>' +
                        '<option value="amex">American Express</option>' +
                        '<option value="discover">Discover</option>' +
                    '</select>' +
                    '<div class="help-text">Select your card provider</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">Card Number</label>' +
                    '<input type="text" class="form-input" name="cardNumber' + cardNumber + '" id="cardNumber' + cardNumber + '" placeholder="1234 5678 9012 3456" maxlength="19" oninput="formatCardNumber(this.value, ' + cardNumber + ')">' +
                    '<div class="help-text">Enter your 16-digit card number</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">Expiry Date</label>' +
                    '<input type="text" class="form-input" name="expiryDate' + cardNumber + '" id="expiryDate' + cardNumber + '" placeholder="MM/YY" maxlength="5" oninput="formatExpiryDate(this.value, ' + cardNumber + ')">' +
                    '<div class="help-text">Enter expiry date in MM/YY format</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">CVV</label>' +
                    '<input type="text" class="form-input" name="cvv' + cardNumber + '" id="cvv' + cardNumber + '" placeholder="123" maxlength="4" oninput="updateCVVDisplay(this.value, ' + cardNumber + ')">' +
                    '<div class="help-text">Enter the 3 or 4 digit security code</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">Cardholder Name</label>' +
                    '<input type="text" class="form-input" name="cardholderName' + cardNumber + '" id="cardholderName' + cardNumber + '" placeholder="John Doe" oninput="updateCardholderDisplay(this.value, ' + cardNumber + ')">' +
                    '<div class="help-text">Enter the name as it appears on the card</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">Billing Address</label>' +
                    '<textarea class="form-input" name="billingAddress' + cardNumber + '" placeholder="Enter your billing address"></textarea>' +
                    '<div class="help-text">Enter the billing address for this card</div>' +
                '</div>' +
                '<div class="form-group">' +
                    '<label class="form-label">Open Banking Enabled</label>' +
                    '<select class="form-select" name="cardOpenBankingEnabled' + cardNumber + '" onchange="toggleCardOpenBankingConfig(this.value, ' + cardNumber + ')">' +
                        '<option value="yes">Yes</option>' +
                        '<option value="no">No</option>' +
                    '</select>' +
                    '<div class="help-text">Enable TrueLayer Open Banking for this card</div>' +
                '</div>' +
                '<div class="open-banking-config" id="cardOpenBankingConfig' + cardNumber + '" style="display: none;">' +
                    '<h4 style="color: var(--fg); margin-bottom: 15px; font-size: 14px;">TrueLayer Open Banking Configuration</h4>' +
                    '<div class="form-group">' +
                        '<label class="form-label">TrueLayer Client ID</label>' +
                        '<input type="text" class="form-input" name="cardTruelayerClientId' + cardNumber + '" placeholder="your-truelayer-client-id">' +
                        '<div class="help-text">Your TrueLayer application client ID</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">TrueLayer Client Secret</label>' +
                        '<input type="password" class="form-input" name="cardTruelayerClientSecret' + cardNumber + '" placeholder="your-truelayer-client-secret">' +
                        '<div class="help-text">Your TrueLayer application client secret</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Redirect URI</label>' +
                        '<input type="url" class="form-input" name="cardTruelayerRedirectUri' + cardNumber + '" placeholder="https://yourapp.com/oauth/callback">' +
                        '<div class="help-text">OAuth2 redirect URI for TrueLayer authentication</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Environment</label>' +
                        '<select class="form-select" name="cardTruelayerEnvironment' + cardNumber + '">' +
                            '<option value="sandbox">Sandbox (Testing)</option>' +
                            '<option value="live">Live (Production)</option>' +
                        '</select>' +
                        '<div class="help-text">TrueLayer environment (sandbox for testing, live for production)</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Scopes</label>' +
                        '<select class="form-select" name="cardTruelayerScopes' + cardNumber + '" multiple>' +
                            '<option value="accounts">Accounts</option>' +
                            '<option value="balance">Balance</option>' +
                            '<option value="transactions">Transactions</option>' +
                            '<option value="cards">Cards</option>' +
                            '<option value="direct_debits">Direct Debits</option>' +
                            '<option value="standing_orders">Standing Orders</option>' +
                        '</select>' +
                        '<div class="help-text">Select the data scopes you need access to</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Card Username</label>' +
                        '<input type="text" class="form-input" name="cardUsername' + cardNumber + '" placeholder="Enter card login username">' +
                        '<div class="help-text">Username for card login (stored securely)</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Card Password</label>' +
                        '<input type="password" class="form-input" name="cardPassword' + cardNumber + '" placeholder="Enter card login password">' +
                        '<div class="help-text">Password for card login (encrypted and PCI compliant)</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">2FA Method</label>' +
                        '<select class="form-select" name="card2faMethod' + cardNumber + '">' +
                            '<option value="none">None</option>' +
                            '<option value="sms">SMS</option>' +
                            '<option value="authenticator">Authenticator App</option>' +
                            '<option value="hardware_token">Hardware Token</option>' +
                            '<option value="biometric">Biometric</option>' +
                        '</select>' +
                        '<div class="help-text">Two-factor authentication method for card access</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">2FA Code/Token</label>' +
                        '<input type="text" class="form-input" name="card2faCode' + cardNumber + '" placeholder="Enter 2FA code or token">' +
                        '<div class="help-text">2FA code, token, or device identifier</div>' +
                    '</div>' +
                    '<div class="form-group">' +
                        '<label class="form-label">Memorable Information</label>' +
                        '<input type="text" class="form-input" name="cardMemorableInfo' + cardNumber + '" placeholder="Enter memorable information">' +
                        '<div class="help-text">Memorable information, security questions, or PIN</div>' +
                    '</div>' +
                '</div>' +
                '<!-- Digital Wallet Integration -->' +
                '<div class="wallet-integration">' +
                    '<h4>Digital Wallet Integration</h4>' +
                    '<div class="wallet-buttons">' +
                        '<button class="wallet-btn" onclick="connectGoogleWallet()">Google Wallet</button>' +
                        '<button class="wallet-btn" onclick="connectAppleWallet()">Apple Wallet</button>' +
                        '<button class="wallet-btn" onclick="connectSamsungPay()">Samsung Pay</button>' +
                    '</div>' +
                '</div>' +
                '<!-- Camera Scanning -->' +
                '<div class="camera-scan">' +
                    '<button class="scan-btn" onclick="scanCard()">📷 Scan Card with Camera</button>' +
                '</div>';
        }

        function toggleCardOpenBankingConfig(enabled, cardNumber) {
            const configDiv = document.getElementById('cardOpenBankingConfig' + cardNumber);
            if (enabled === 'yes') {
                configDiv.style.display = 'block';
            } else {
                configDiv.style.display = 'none';
            }
        }

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
