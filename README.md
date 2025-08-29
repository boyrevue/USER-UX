# 🏗️ CLIENT-UX Personal Data Manager - TTL-Driven Architecture

## 🎯 Overview
CLIENT-UX is a **semantic web application** where ALL form definitions, field types, validation rules, and UI behavior are defined in TTL ontology files and dynamically interpreted by the system. **The TTL ontology is the single source of truth.**

## 📜 SYSTEM DOCTRINE
> **⚠️ CRITICAL**: Read [`SYSTEM_DOCTRINE.md`](SYSTEM_DOCTRINE.md) before making ANY changes. The TTL-as-single-source-of-truth principle is mandatory.

**Quick References:**
- 📖 [`SYSTEM_DOCTRINE.md`](SYSTEM_DOCTRINE.md) - Core principles and rules
- 🔧 [`TTL_IMPLEMENTATION_GUIDE.md`](TTL_IMPLEMENTATION_GUIDE.md) - Technical implementation
- ⚡ [`TTL_QUICK_REFERENCE.md`](TTL_QUICK_REFERENCE.md) - Developer cheat sheet

## 🏗️ Architecture

### TTL-Driven Architecture
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   TTL ONTOLOGY  │    │   GO BACKEND    │    │  REACT FRONTEND │
│  (Single Source)│    │ (TTL Interpreter)│   │ (Dynamic Renderer)│
├─────────────────┤    ├─────────────────┤    ├─────────────────┤
│ • autoins.ttl   │───▶│ • ttl_parser.go │───▶│ • Dynamic Forms │
│   - Fields      │    │ • /api/ontology │    │ • Auto-generated│
│   - Labels      │    │ • Field Types   │    │ • TTL-driven UI │
│   - Validation  │    │ • Validation    │    │ • No hardcoding │
│   - Help Text   │    │ • OCR Engine    │    │                 │
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Key Principles
1. **TTL Supremacy**: ALL form definitions MUST be in autoins.ttl
2. **Dynamic Extraction**: Go parses TTL at runtime, builds API dynamically  
3. **Zero Hardcoding**: No field definitions in Go/JavaScript code
4. **Semantic Web**: True RDF/OWL ontology-driven architecture

## 📁 File Structure
```
client-ux/
├── ontology/              # 🧠 SINGLE SOURCE OF TRUTH
│   └── autoins.ttl       # ⭐ THE ontology - ALL fields defined here
├── main.go               # 🚀 SERVER - Go backend + API
├── ttl_parser.go         # 🔍 TTL PARSER - Dynamic ontology interpreter  
├── types.go              # 🏗️ STRUCTURES - Go data models
├── document_processor.go # 📄 OCR ENGINE - Passport/document processing
├── insurance-frontend/   # 🎨 UI - React application
│   ├── src/App.tsx      # Dynamic form renderer
│   └── build/           # Compiled React files
├── static/              # 🌐 SERVED - Go serves React build
│   ├── index.html       # React entry point
│   ├── js/              # React JavaScript
│   └── css/             # React stylesheets
├── SYSTEM_DOCTRINE.md   # 📜 CORE PRINCIPLES - READ FIRST
├── TTL_IMPLEMENTATION_GUIDE.md # 🔧 Technical guide
└── TTL_QUICK_REFERENCE.md      # ⚡ Developer cheat sheet
```

### 🚨 ELIMINATED FILES (TTL Doctrine Compliance)
- ❌ `categories.json` - Removed (redundant with TTL)
- ❌ `fields.json` - Removed (redundant with TTL)  
- ❌ `subforms.json` - Removed (redundant with TTL)
- ❌ `parser.go` - Removed (replaced with ttl_parser.go)

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

## 🎯 TTL DOCTRINE SUMMARY

### ✅ WHAT WE ACHIEVED
- **Single Source of Truth**: `autoins.ttl` is the ONLY place where fields are defined
- **Dynamic Extraction**: 82 driver fields, 45 UK conviction codes, all extracted from TTL
- **Zero Hardcoding**: No field definitions in Go/JavaScript code
- **Semantic Web Compliance**: True RDF/OWL ontology-driven architecture

### 🚨 DEVELOPER RULES
1. **BEFORE** adding any field → Add to `autoins.ttl` first
2. **NEVER** hardcode field definitions in code
3. **ALWAYS** use `/api/ontology` for form structure
4. **READ** `SYSTEM_DOCTRINE.md` before making changes

### 🔧 QUICK FIELD ADDITION
```turtle
# Add to autoins.ttl
autoins:newField a owl:DatatypeProperty ;
  rdfs:domain autoins:Driver ;
  rdfs:range xsd:string ;
  rdfs:label "New Field" ;
  autoins:isRequired "true"^^xsd:boolean .
```
```bash
# Restart & verify
pkill -f client-ux && ./client-ux &
curl -s http://localhost:3000/api/ontology | jq '.drivers.fields[] | select(.property == "newField")'
```

**The TTL ontology is the single source of truth. This is not negotiable.** 🎯


