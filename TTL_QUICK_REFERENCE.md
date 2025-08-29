# TTL Quick Reference Card

## 🚨 GOLDEN RULES
1. **ALL** form fields MUST be defined in `autoins.ttl`
2. **NO** hardcoded field definitions in code
3. **SINGLE** source of truth = TTL ontology

---

## 📝 ADD NEW FIELD (3 Steps)

### Step 1: Add to TTL
```turtle
autoins:myNewField a owl:DatatypeProperty ;
  rdfs:domain autoins:Driver ;           # Driver/Vehicle/Claims
  rdfs:range xsd:string ;                # string/boolean/date
  rdfs:label "My New Field" ;            # Display label
  autoins:isRequired "true"^^xsd:boolean ; # Required?
  autoins:formHelpText "Help text here" . # Optional help
```

### Step 2: Restart App
```bash
pkill -f client-ux && ./client-ux &
```

### Step 3: Verify
```bash
curl -s http://localhost:3000/api/ontology | jq '.drivers.fields[] | select(.property == "myNewField")'
```

---

## 🎛️ FIELD TYPES

| TTL Range | Result | Use Case |
|-----------|--------|----------|
| `xsd:string` | Text input | Names, addresses |
| `xsd:boolean` | Radio (Yes/No) | True/false questions |
| `xsd:date` | Date picker | Dates |
| `xsd:string` + enum ≤3 | Radio buttons | Few options |
| `xsd:string` + enum >3 | Select dropdown | Many options |

---

## 📋 FIELD WITH OPTIONS

```turtle
autoins:mySelectField a owl:DatatypeProperty ;
  rdfs:domain autoins:Driver ;
  rdfs:range xsd:string ;
  rdfs:label "Choose Option" ;
  autoins:enumerationValues ("OPTION1" "OPTION2" "OPTION3") ;
  autoins:formHelpText "Pick one" .
```

---

## 🔍 TESTING COMMANDS

```bash
# Check all driver fields
curl -s http://localhost:3000/api/ontology | jq '.drivers.fields | map(.property)'

# Check specific field
curl -s http://localhost:3000/api/ontology | jq '.drivers.fields[] | select(.property == "hasConvictions")'

# Count fields by section
curl -s http://localhost:3000/api/ontology | jq '{drivers: (.drivers.fields|length), vehicles: (.vehicles.fields|length)}'

# Validate TTL syntax
rapper -i turtle -o ntriples ontology/autoins.ttl > /dev/null && echo "✅ Valid TTL"
```

---

## 🚫 FORBIDDEN ACTIONS

- ❌ Adding fields to JSON files
- ❌ Hardcoding field arrays in Go/JS
- ❌ Duplicating field definitions
- ❌ Bypassing the TTL ontology

---

## 🔧 TROUBLESHOOTING

| Problem | Solution |
|---------|----------|
| Field not appearing | Check TTL syntax, restart app |
| Wrong field type | Verify `rdfs:range` and enumerations |
| Field in wrong section | Check `rdfs:domain` |
| Missing label | Ensure `rdfs:label "text"` is quoted |

---

## 📍 KEY FILES

- `ontology/autoins.ttl` - **THE** source of truth
- `ttl_parser.go` - TTL → JSON converter
- `/api/ontology` - Dynamic API endpoint
- `SYSTEM_DOCTRINE.md` - Full doctrine

---

## ⚡ QUICK EXAMPLES

### Required Text Field
```turtle
autoins:email a owl:DatatypeProperty ;
  rdfs:domain autoins:Driver ;
  rdfs:range xsd:string ;
  rdfs:label "email address" ;
  autoins:isRequired "true"^^xsd:boolean .
```

### Optional Date Field
```turtle
autoins:birthDate a owl:DatatypeProperty ;
  rdfs:domain autoins:Driver ;
  rdfs:range xsd:date ;
  rdfs:label "date of birth" ;
  autoins:isRequired "false"^^xsd:boolean .
```

### Yes/No Radio Field
```turtle
autoins:hasLicense a owl:DatatypeProperty ;
  rdfs:domain autoins:Driver ;
  rdfs:range xsd:boolean ;
  rdfs:label "has driving license" ;
  autoins:enumerationValues ("YES" "NO") ;
  autoins:isRequired "true"^^xsd:boolean .
```

---

**Remember: If it's not in the TTL, it doesn't exist! 🎯**
