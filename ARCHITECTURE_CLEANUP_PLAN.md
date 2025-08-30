# 🏗️ CLIENT-UX Architecture Cleanup & Optimization Plan

## 📊 **Current Codebase Analysis**

### **File Size Issues (AI Management Concerns)**
| File | Lines | Issue | Recommendation |
|------|-------|-------|----------------|
| `App.tsx` | 5,452 | **CRITICAL** - Too large for AI | Split into components |
| `document_processor.go` | 2,265 | **HIGH** - Monolithic | Extract services |
| `user_ux.ttl` | 1,515 | **MEDIUM** - Large ontology | Modularize further |
| `AI_Claims_History.ttl` | 1,190 | **MEDIUM** - Complex | Split by domain |

### **Directory Bloat**
- **Sessions**: 46 files (176KB) - Needs cleanup
- **Static**: 54MB - Contains duplicates
- **Node_modules**: 432MB - Standard but excludable
- **Documentation**: 15+ MD files - Needs organization

---

## 🎯 **Recommended Clean Architecture**

### **1. Frontend Restructuring (Priority: CRITICAL)**

#### **Current Problem**: 
- `App.tsx` (5,452 lines) is unmanageable for AI
- Monolithic component with mixed concerns
- Poor separation of business logic

#### **Solution**: Micro-Frontend Architecture
```
src/
├── components/           # Reusable UI components
│   ├── forms/           # Form components (max 200 lines each)
│   ├── layout/          # Layout components
│   ├── ui/              # Basic UI elements
│   └── validation/      # Validation components
├── pages/               # Page-level components (max 300 lines each)
│   ├── DriverDetails/   # Driver management
│   ├── VehicleDetails/  # Vehicle management
│   ├── ClaimsHistory/   # Claims management
│   ├── PolicyDetails/   # Policy configuration
│   └── Documents/       # Document processing
├── services/            # Business logic services
│   ├── api/            # API communication
│   ├── validation/     # Validation logic
│   ├── ocr/           # OCR processing
│   └── gdpr/          # GDPR compliance
├── hooks/              # Custom React hooks
├── utils/              # Utility functions
├── types/              # TypeScript definitions
└── constants/          # Application constants
```

#### **Implementation Steps**:
1. **Extract Components** (Week 1)
   - Split `App.tsx` into 15-20 components (max 300 lines each)
   - Create dedicated form components for each section
   - Extract validation logic into custom hooks

2. **Service Layer** (Week 2)
   - Move API calls to service layer
   - Create GDPR compliance service
   - Extract OCR processing logic

3. **State Management** (Week 3)
   - Implement Context API or Zustand for state
   - Remove prop drilling
   - Add proper error boundaries

### **2. Backend Modularization (Priority: HIGH)**

#### **Current Problem**:
- `document_processor.go` (2,265 lines) handles multiple concerns
- `main.go` (713 lines) mixes routing with business logic

#### **Solution**: Clean Go Architecture
```
cmd/
└── server/
    └── main.go              # Entry point only (max 50 lines)

internal/
├── api/                     # HTTP handlers (max 200 lines each)
│   ├── handlers/
│   │   ├── driver.go
│   │   ├── vehicle.go
│   │   ├── claims.go
│   │   ├── documents.go
│   │   └── health.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logging.go
│   │   └── gdpr.go
│   └── routes/
│       └── router.go
├── services/               # Business logic (max 300 lines each)
│   ├── ocr/
│   │   ├── passport.go
│   │   ├── licence.go
│   │   └── document.go
│   ├── validation/
│   │   ├── age.go
│   │   ├── dates.go
│   │   └── shacl.go
│   ├── gdpr/
│   │   ├── compliance.go
│   │   ├── masking.go
│   │   └── audit.go
│   └── ontology/
│       ├── parser.go
│       └── validator.go
├── models/                 # Data structures
│   ├── driver.go
│   ├── vehicle.go
│   ├── claims.go
│   └── session.go
├── repository/             # Data access layer
│   ├── session.go
│   └── ontology.go
└── config/                 # Configuration
    └── config.go
```

### **3. TTL Ontology Optimization (Priority: MEDIUM)**

#### **Current Issues**:
- Some TTL files exceed 1,000 lines
- Mixed concerns in single files
- Difficult for AI to parse and modify

#### **Solution**: Micro-Ontology Architecture
```
ontology/
├── core/                   # Core definitions (max 200 lines each)
│   ├── base.ttl           # Base classes and properties
│   ├── validation.ttl     # SHACL shapes
│   └── i18n.ttl          # Internationalization
├── domains/               # Domain-specific ontologies (max 500 lines each)
│   ├── insurance/
│   │   ├── driver.ttl     # Driver-specific (300 lines)
│   │   ├── vehicle.ttl    # Vehicle-specific (300 lines)
│   │   ├── policy.ttl     # Policy-specific (300 lines)
│   │   ├── claims.ttl     # Claims-specific (400 lines)
│   │   └── payments.ttl   # Payments-specific (300 lines)
│   ├── documents/
│   │   ├── identity.ttl   # ID documents (200 lines)
│   │   ├── financial.ttl  # Financial docs (200 lines)
│   │   └── legal.ttl      # Legal documents (200 lines)
│   └── compliance/
│       ├── gdpr.ttl       # GDPR definitions (300 lines)
│       ├── uk_law.ttl     # UK legal requirements (200 lines)
│       └── validation.ttl # Validation rules (200 lines)
└── applications/          # Application-specific (max 300 lines each)
    ├── insurance_app.ttl  # Insurance application config
    ├── government_app.ttl # Government services config
    └── user_preferences.ttl # User settings
```

### **4. Documentation Restructuring (Priority: LOW)**

#### **Current Problem**:
- 15+ markdown files in root directory
- Mixed technical and business documentation
- Difficult to navigate

#### **Solution**: Organized Documentation
```
docs/
├── README.md              # Main project overview
├── QUICK_START.md         # Getting started guide
├── architecture/          # Technical documentation
│   ├── backend.md
│   ├── frontend.md
│   ├── ontology.md
│   └── deployment.md
├── compliance/            # Compliance documentation
│   ├── gdpr.md
│   ├── security.md
│   └── audit.md
├── api/                   # API documentation
│   ├── endpoints.md
│   ├── authentication.md
│   └── examples.md
├── development/           # Developer guides
│   ├── setup.md
│   ├── contributing.md
│   └── testing.md
└── business/             # Business documentation
    ├── product_overview.md
    ├── use_cases.md
    └── roadmap.md
```

---

## 🚀 **Implementation Roadmap**

### **Phase 1: Critical Fixes (Week 1-2)**
1. **Split App.tsx** into manageable components
2. **Extract document_processor.go** into services
3. **Clean up session files** (keep last 10 only)
4. **Remove duplicate static files**

### **Phase 2: Architecture Refactoring (Week 3-4)**
1. **Implement Go clean architecture**
2. **Create service layer separation**
3. **Add proper error handling and logging**
4. **Implement dependency injection**

### **Phase 3: Ontology Optimization (Week 5-6)**
1. **Split large TTL files** into micro-ontologies
2. **Create domain-specific modules**
3. **Optimize SHACL validation rules**
4. **Add automated ontology testing**

### **Phase 4: Documentation & Tooling (Week 7-8)**
1. **Reorganize documentation structure**
2. **Add automated code generation**
3. **Create development tooling**
4. **Implement CI/CD pipeline**

---

## 📏 **AI-Friendly File Size Guidelines**

### **Recommended Maximum Sizes**:
- **React Components**: 300 lines (200 preferred)
- **Go Services**: 300 lines (200 preferred)
- **Go Handlers**: 150 lines (100 preferred)
- **TTL Ontologies**: 500 lines (300 preferred)
- **Documentation**: 500 lines (300 preferred)

### **Complexity Metrics**:
- **Cyclomatic Complexity**: < 10 per function
- **Function Length**: < 50 lines
- **File Dependencies**: < 10 imports
- **Nesting Depth**: < 4 levels

---

## 🛠️ **Tooling & Automation**

### **Code Quality Tools**:
```bash
# Go tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
go install github.com/securego/gosec/v2/cmd/gosec@latest

# Frontend tools
npm install -g eslint prettier typescript

# TTL validation
pip install rdflib pyshacl

# Documentation
npm install -g @mermaid-js/mermaid-cli
```

### **Pre-commit Hooks**:
```yaml
# .pre-commit-config.yaml
repos:
  - repo: local
    hooks:
      - id: go-lint
        name: Go Lint
        entry: golangci-lint run
        language: system
        files: \.go$
      - id: ttl-validate
        name: TTL Validation
        entry: python scripts/validate_ttl.py
        language: system
        files: \.ttl$
      - id: file-size-check
        name: File Size Check
        entry: python scripts/check_file_sizes.py
        language: system
```

---

## 🎯 **Expected Benefits**

### **AI Manageability**:
- ✅ **95% reduction** in file complexity for AI processing
- ✅ **Clear separation** of concerns for focused AI assistance
- ✅ **Modular architecture** enables targeted AI code generation
- ✅ **Standardized patterns** improve AI code understanding

### **Developer Experience**:
- ✅ **Faster development** with smaller, focused files
- ✅ **Easier debugging** with clear service boundaries
- ✅ **Better testing** with isolated components
- ✅ **Improved maintainability** with clean architecture

### **System Performance**:
- ✅ **Faster compilation** with smaller Go files
- ✅ **Better caching** with modular frontend components
- ✅ **Optimized loading** with code splitting
- ✅ **Reduced memory usage** with lazy loading

---

## 📋 **Cleanup Checklist**

### **Immediate Actions (This Week)**:
- [ ] Split `App.tsx` into 10+ components
- [ ] Extract OCR service from `document_processor.go`
- [ ] Clean up session files (keep last 10)
- [ ] Remove duplicate static assets
- [ ] Create `.gitignore` for build artifacts

### **Short-term (Next 2 Weeks)**:
- [ ] Implement Go clean architecture
- [ ] Create service layer separation
- [ ] Add proper error handling
- [ ] Split large TTL files
- [ ] Reorganize documentation

### **Medium-term (Next Month)**:
- [ ] Add automated testing
- [ ] Implement CI/CD pipeline
- [ ] Create development tooling
- [ ] Add performance monitoring
- [ ] Optimize build process

---

## 🔧 **Implementation Scripts**

### **File Size Monitor**:
```bash
#!/bin/bash
# scripts/check_file_sizes.sh
find . -name "*.go" -o -name "*.tsx" -o -name "*.ttl" | while read file; do
    lines=$(wc -l < "$file")
    if [ "$lines" -gt 500 ]; then
        echo "WARNING: $file has $lines lines (>500)"
    fi
done
```

### **Component Splitter**:
```python
# scripts/split_components.py
import re
import os

def split_react_component(file_path, max_lines=300):
    """Split large React components into smaller ones"""
    # Implementation for automated component splitting
    pass
```

### **TTL Validator**:
```python
# scripts/validate_ttl.py
from rdflib import Graph
import sys

def validate_ttl_file(file_path):
    """Validate TTL syntax and size"""
    try:
        g = Graph()
        g.parse(file_path, format='turtle')
        
        with open(file_path, 'r') as f:
            lines = len(f.readlines())
            
        if lines > 500:
            print(f"WARNING: {file_path} has {lines} lines (>500)")
            
        return True
    except Exception as e:
        print(f"ERROR: {file_path} - {e}")
        return False
```

---

This cleanup plan will transform CLIENT-UX into an **AI-friendly, maintainable, and scalable architecture** while preserving all existing functionality and improving development velocity.

**🎯 Priority: Start with App.tsx splitting and document_processor.go modularization for immediate AI manageability improvements!**
