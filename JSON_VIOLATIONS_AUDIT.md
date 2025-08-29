# JSON Files Audit - TTL Doctrine Violations

## 🚨 CRITICAL VIOLATIONS FOUND

### 1. Configuration Files
- **`config.json`** ✅ **FIXED** - Updated to reflect TTL doctrine
  - Removed references to eliminated JSON files
  - Updated application name and architecture description

### 2. Internationalization Files
- **`i18n/en.json`** ⚠️ **MAJOR VIOLATION** - Contains hardcoded field labels
- **`i18n/de.json`** ⚠️ **MAJOR VIOLATION** - Contains hardcoded field labels

**Violation Details:**
```json
"categories": {
  "drivers": {"title": "Driver Details", "description": "Add up to 4 drivers"},
  "claims": {"title": "Claims History", "description": "Claims and convictions"}
},
"fields": {
  "firstName": "First Name",
  "hasConvictions": "Any driving convictions?",
  "convictionType": "Conviction Type"
}
```

**These should come from TTL ontology with multi-language support!**

### 3. Session Files
- **`sessions/*.json`** ✅ **CLEANED** - Reduced from 218 to 10 files
  - These are legitimate runtime session data (not form definitions)

### 4. Unused Directories
- **`settings-frontend/`** ❌ **UNUSED** - Contains package.json, tsconfig.json
- **`settings_app/`** ❌ **UNUSED** - Contains package.json, tsconfig.json

### 5. Legitimate JSON Files (No Violations)
- **`insurance-frontend/package.json`** ✅ **OK** - NPM dependencies
- **`insurance-frontend/tsconfig.json`** ✅ **OK** - TypeScript config
- **`static/manifest.json`** ✅ **OK** - PWA manifest
- **`static/asset-manifest.json`** ✅ **OK** - React build manifest

---

## 🎯 REQUIRED ACTIONS

### HIGH PRIORITY
1. **Move i18n field labels to TTL ontology**
   - Add multi-language support to `autoins.ttl`
   - Use `rdfs:label` with language tags: `"first name"@en`, `"Vorname"@de`
   - Update TTL parser to handle language-specific labels

2. **Remove unused directories**
   - Delete `settings-frontend/` and `settings_app/`
   - They contain JSON files that could confuse developers

### MEDIUM PRIORITY
3. **Keep only essential i18n content**
   - Month/day mappings (used by OCR engine)
   - System messages not related to forms
   - Remove all field labels and categories

---

## 🔧 TTL MULTI-LANGUAGE SOLUTION

### Current TTL (Single Language)
```turtle
autoins:firstName rdfs:label "first name" ;
```

### Proposed TTL (Multi-Language)
```turtle
autoins:firstName rdfs:label "first name"@en ;
autoins:firstName rdfs:label "Vorname"@de ;
autoins:firstName rdfs:label "prénom"@fr ;
```

### Updated TTL Parser
```go
// Extract language-specific labels
labelPattern := regexp.MustCompile(`rdfs:label\s+"([^"]+)"@(\w+)\s*;`)
```

---

## 📊 VIOLATION SUMMARY

| File Type | Count | Status | Action Required |
|-----------|-------|--------|-----------------|
| Config JSON | 1 | ✅ Fixed | None |
| i18n JSON | 2 | ⚠️ Major Violation | Move to TTL |
| Session JSON | 10 | ✅ OK | None (runtime data) |
| Build JSON | 6 | ✅ OK | None (build artifacts) |
| Unused JSON | 4 | ❌ Remove | Delete directories |

---

## 🎯 DOCTRINE COMPLIANCE SCORE

**Current: 60% Compliant**
- ✅ Eliminated: fields.json, subforms.json, categories.json
- ✅ Fixed: config.json
- ✅ Cleaned: session files
- ⚠️ **Major Issue**: i18n field labels still hardcoded
- ❌ **Minor Issue**: Unused directories with JSON files

**Target: 100% Compliant**
- Move all form-related content to TTL ontology
- Implement multi-language TTL support
- Remove unused directories

---

**The TTL ontology must be the ONLY source for field definitions, labels, and form structure. No exceptions.**
