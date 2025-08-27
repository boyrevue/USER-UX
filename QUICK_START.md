# 🚀 Quick Start Guide

## One Command to Start Everything
```bash
./start.sh
```

## What This Does
1. ✅ Kills any existing processes
2. ✅ Builds React frontend 
3. ✅ Copies build to static/
4. ✅ Fixes nested static/static/ issue
5. ✅ Builds Go backend
6. ✅ Starts server on http://localhost:3000

## Architecture Overview
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   ONTOLOGY      │    │   GO BACKEND    │    │  REACT FRONTEND │
│   FILES         │    │   (Interpreter) │    │   (UI Renderer) │
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ • autoins.ttl   │───▶│ • LoadOntology()│───▶│ • Sidebar Menu  │
│ • categories.json│   │ • API Endpoints │   │ • Form Builder  │
│ • fields.json   │   │ • SHACL Validation│   │ • State Mgmt    │
│ • subforms.json │   │ • Session Mgmt   │   │ • Validation    │
│ • settings.ttl  │   │ • File Serving   │   │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

## Key Features
- **Hierarchical Sidebar**: Left sidebar with collapsible sections
- **Ontology-Driven**: All forms and validation defined in ontology files
- **Multi-language**: English and German support
- **Document Processing**: Upload and process insurance documents

## Troubleshooting
- **Blank Screen**: Run `./start.sh` to rebuild everything
- **404 Errors**: Check that static files are properly copied
- **Sidebar Missing**: Ensure React build is up to date

## File Structure
```
insurance-quote-app/
├── start.sh              # 🚀 One-command startup
├── config.json           # 📋 Application configuration  
├── main.go              # 🖥️ Go server
├── ontology/            # 🧠 Application definition
├── insurance-frontend/  # 🎨 React app
└── static/              # 🌐 Served files
```

## Remember
- **Always use `./start.sh`** for fresh starts
- **Ontology files are the source of truth**
- **React builds to `insurance-frontend/build/` → copy to `static/`**
- **Go serves from `static/` directory**
