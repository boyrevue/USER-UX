# 🏗️ Insurance Quote Application - Ontology-Driven Architecture

## 🎯 Overview
This application uses **ontology-driven development** where the application structure, forms, validation, and UI are all defined in ontology files (RDF/TTL + JSON) and interpreted by a Go backend to serve a React frontend.

## 🏗️ Architecture

### Core Components
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

### Key Principles
1. **Ontology-First**: All application structure defined in RDF/TTL + JSON
2. **Go Interpreter**: Backend reads ontology, serves API, validates with SHACL
3. **React Renderer**: Frontend consumes API, builds UI dynamically
4. **Static File Structure**: React build → `static/` → Go serves correctly

## 📁 File Structure
```
insurance-quote-app/
├── ontology/           # 🧠 BRAIN - Application Definition
│   ├── autoins.ttl    # Main insurance ontology (RDF/OWL)
│   ├── categories.json # Menu structure & navigation
│   ├── fields.json    # Form field definitions
│   ├── subforms.json  # Dynamic form components
│   └── settings.ttl   # Configuration & validation
├── main.go            # 🚀 SERVER - Go backend
├── parser.go          # 📖 INTERPRETER - Ontology loader
├── types.go           # 🏗️ STRUCTURES - Go data models
├── insurance-frontend/ # 🎨 UI - React application
│   ├── src/App.tsx    # Main React component
│   └── build/         # Compiled React files
└── static/            # 🌐 SERVED - Go serves React build
    ├── index.html     # React entry point
    ├── js/            # React JavaScript
    └── css/           # React stylesheets
```

## 🚀 Quick Start

### 1. Build React Frontend
```bash
cd insurance-frontend
npm run build
cd ..
cp -r insurance-frontend/build/* static/
```

### 2. Start Go Backend
```bash
go build -o insurance-quote-app .
./insurance-quote-app
```

### 3. Access Application
Open http://localhost:3000

## 🔧 Critical Configuration

### Go Server Configuration
- **Root route**: Serves `./static/index.html` (React build)
- **Static files**: Serves `./static/js/` and `./static/css/`
- **API endpoints**: `/api/category/list`, `/api/category/{id}`

### React Build Process
1. React builds to `insurance-frontend/build/`
2. Copy build files to `static/` directory
3. Go server serves from `static/` directory

### File Serving Fix
**IMPORTANT**: React creates nested `static/static/` structure. Fix with:
```bash
cp -r static/static/js/* static/js/
cp -r static/static/css/* static/css/
rm -rf static/static
```

## 🎨 UI Components

### Hierarchical Sidebar Menu
- **Structure**: Defined in `menuStructure` array in App.tsx
- **Sections**: Car Insurance, Settings
- **Subcategories**: Expandable/collapsible
- **State**: `expandedSections` manages open/closed state

### Form Generation
- **Source**: Ontology files define form structure
- **API**: `/api/category/{id}` returns field definitions
- **Rendering**: React builds forms dynamically

## 🔍 Troubleshooting

### Blank Screen Issues
1. Check static file serving: `curl http://localhost:3000/static/js/main.*.js`
2. Verify file structure: `ls -la static/js/`
3. Check browser console for 404 errors

### Sidebar Not Showing
1. Verify React build: `npm run build` in insurance-frontend
2. Copy build files: `cp -r insurance-frontend/build/* static/`
3. Check file structure: No nested `static/static/` directories

### API Errors
1. Check ontology files exist: `ls -la ontology/`
2. Verify Go server running: `curl http://localhost:3000/api/category/list`
3. Check Go compilation: `go build .`

## 🧠 Ontology-Driven Development

### Adding New Features
1. **Define in Ontology**: Add to `categories.json`, `fields.json`, etc.
2. **Update Go Types**: Modify `types.go` if needed
3. **Update React**: Add UI components in `App.tsx`
4. **Rebuild**: `npm run build` → copy to `static/`

### Menu Structure
```javascript
const menuStructure = [
  {
    id: 'car-insurance',
    title: 'Car Insurance',
    icon: Shield,
    categories: [
      { id: 'drivers', title: 'Driver Details', icon: User },
      { id: 'vehicle', title: 'Vehicle Details', icon: Car },
      // ... more categories
    ]
  }
];
```

## 🔒 Security & Validation
- **SHACL Shapes**: Defined in TTL files for validation
- **Session Management**: Go handles user sessions
- **API Security**: CORS configured for development

## 📚 Dependencies
- **Go**: 1.21+
- **Node.js**: 16+
- **React**: 18+
- **Tailwind CSS**: Styling
- **Flowbite**: UI Components
- **Lucide**: Icons

---

**Remember**: This is an ontology-driven application. The ontology files are the source of truth for all application behavior!


