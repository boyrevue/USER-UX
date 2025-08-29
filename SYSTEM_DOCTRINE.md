# CLIENT-UX System Doctrine

## Core Principle: TTL as Single Source of Truth

### **DOCTRINE STATEMENT**
> **The TTL ontology files (`autoins.ttl`) are the SINGLE, AUTHORITATIVE source of truth for all form definitions, field types, validation rules, labels, and UI behavior in the CLIENT-UX system. No exceptions.**

---

## 1. ARCHITECTURAL COMMANDMENTS

### 1.1 The TTL Supremacy Rule
- **ALL** form fields, labels, types, and validation MUST be defined in TTL ontology files
- **NO** hardcoded field definitions in Go, JavaScript, or JSON files
- **NO** duplicate field definitions across multiple sources
- **ALL** UI components MUST derive their structure from the TTL-based API

### 1.2 The Elimination Principle
The following are **FORBIDDEN** in the CLIENT-UX system:
- ❌ `fields.json` - Eliminated
- ❌ `subforms.json` - Eliminated  
- ❌ `categories.json` - Eliminated
- ❌ `parser.go` (JSON-based) - Eliminated
- ❌ Hardcoded field arrays in any source code
- ❌ Duplicate form definitions

### 1.3 The Dynamic Extraction Rule
- **ALL** form metadata MUST be extracted dynamically from TTL at runtime
- Field properties (name, label, type, required, help text, options) MUST come from semantic triples
- The system MUST adapt automatically to TTL changes without code modifications

---

## 2. TECHNICAL IMPLEMENTATION

### 2.1 TTL Parser Architecture
```
autoins.ttl (Single Source) → ttl_parser.go → /api/ontology → React Frontend
```

### 2.2 Semantic Property Mapping
| TTL Property | System Function |
|--------------|-----------------|
| `rdfs:label` | Field display label |
| `rdfs:domain` | Section assignment (Driver/Vehicle/Claims) |
| `rdfs:range` | Field data type (xsd:string, xsd:boolean, xsd:date) |
| `autoins:isRequired` | Field validation requirement |
| `autoins:formHelpText` | Field help/description text |
| `autoins:enumerationValues` | Select/radio options |

### 2.3 Dynamic Field Type Resolution
- `xsd:boolean` + enumeration → Radio buttons
- `xsd:date` → Date picker
- `xsd:string` + enumeration (≤3 options) → Radio buttons
- `xsd:string` + enumeration (>3 options) → Select dropdown
- `xsd:string` (email pattern) → Email input
- `xsd:string` (phone pattern) → Tel input
- `xsd:string` (default) → Text input

---

## 3. OPERATIONAL PROCEDURES

### 3.1 Adding New Fields
1. **ONLY** add field definitions to `autoins.ttl`
2. Define semantic properties: domain, range, label, required, help text
3. System automatically detects and renders new fields
4. **NO** code changes required

### 3.2 Modifying Existing Fields
1. **ONLY** modify field properties in `autoins.ttl`
2. Changes are immediately reflected via API
3. **NO** frontend or backend code changes required

### 3.3 Field Validation
1. All validation rules MUST be defined in TTL using SHACL or custom properties
2. **NO** hardcoded validation in application code
3. Validation logic MUST be derived from ontology

---

## 4. COMPLIANCE VERIFICATION

### 4.1 System Health Checks
```bash
# Verify TTL parsing is working
curl -s http://localhost:3000/api/ontology | jq '.drivers.fields | length'

# Verify no hardcoded fields remain
grep -r "hardcoded\|fields\.json\|subforms\.json" src/ || echo "✅ Clean"

# Verify dynamic extraction
curl -s http://localhost:3000/api/ontology | jq '.drivers.fields[] | select(.property == "hasConvictions")'
```

### 4.2 Doctrine Compliance Checklist
- [ ] All fields extracted from TTL ✅
- [ ] No JSON field definitions ✅
- [ ] No hardcoded field arrays ✅
- [ ] Dynamic type resolution working ✅
- [ ] API serves TTL-derived structure ✅
- [ ] Frontend consumes dynamic API ✅

---

## 5. BENEFITS OF TTL DOCTRINE

### 5.1 Semantic Web Compliance
- True RDF/OWL ontology-driven architecture
- Interoperable with semantic web tools
- Machine-readable domain knowledge

### 5.2 Maintainability
- Single point of truth eliminates inconsistencies
- Changes propagate automatically
- Reduced code complexity

### 5.3 Extensibility
- New domains easily added via TTL
- Multi-language support through ontology
- SHACL validation integration ready

### 5.4 Developer Experience
- No need to modify multiple files for field changes
- Ontology-driven development workflow
- Clear separation of concerns

---

## 6. VIOLATION CONSEQUENCES

### 6.1 Code Review Requirements
- **ANY** hardcoded field definition MUST be rejected
- **ANY** duplicate field source MUST be eliminated
- **ALL** form changes MUST go through TTL first

### 6.2 System Integrity
- Violations compromise the single source of truth
- Creates maintenance nightmares
- Breaks semantic web compliance

---

## 7. FUTURE EVOLUTION

### 7.1 Multi-Domain Support
- Insurance (`autoins.ttl`)
- Finance (`finance.ttl`)
- Legal (`legal.ttl`)
- Health (`health.ttl`)

### 7.2 Advanced Features
- SHACL validation integration
- Multi-language ontology support
- Voice/TTS ontology-driven interfaces
- AI-assisted form completion from ontology

---

## 8. DOCTRINE ENFORCEMENT

### 8.1 Development Rules
1. **BEFORE** adding any field: Check if it exists in TTL
2. **BEFORE** modifying UI: Check if change should be in TTL
3. **BEFORE** hardcoding anything: Ask "Should this be in the ontology?"

### 8.2 Code Standards
- All form-related code MUST consume `/api/ontology`
- No direct field definitions outside TTL
- Dynamic rendering based on ontology structure

---

## CONCLUSION

**The TTL-as-single-source-of-truth doctrine is not just a technical choice—it's the foundational principle that makes CLIENT-UX a true semantic web application. This doctrine ensures consistency, maintainability, and semantic compliance while enabling powerful ontology-driven features.**

**Adherence to this doctrine is mandatory for all system development and maintenance.**

---

*Document Version: 1.0*  
*Last Updated: 2025-01-28*  
*Status: ACTIVE DOCTRINE*
