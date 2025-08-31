# Complete Form Flow Example: EU Driver with Medical Conditions

This document demonstrates a complete end-to-end example of how the ontology-driven form system works in practice.

## Scenario
A user from Germany with an EU licence who has diabetes and wears glasses applies for UK car insurance.

## Step-by-Step Flow

### Step 1: Initial Form Load

**API Request:**
```bash
GET /api/ontology
```

**API Response (Relevant Fields):**
```json
{
  "sections": {
    "drivers": {
      "fields": [
        {
          "property": "licenceType",
          "label": "Licence Type",
          "type": "select",
          "required": true,
          "enumerationValues": ["FULL_UK", "PROVISIONAL_UK", "EU_EEA", "INTERNATIONAL", "OTHER_FOREIGN"],
          "formType": "select",
          "conditionalDisplay": null
        },
        {
          "property": "exchangedToUK",
          "label": "Exchanged To UK Licence?",
          "type": "radio",
          "required": false,
          "conditionalDisplay": "licenceType=EU_EEA OR licenceType=INTERNATIONAL OR licenceType=OTHER_FOREIGN",
          "enumerationValues": ["YES", "NO"],
          "formType": "radio"
        },
        {
          "property": "hasMedicalConditions",
          "label": "Any Medical Conditions?",
          "type": "radio",
          "required": true,
          "enumerationValues": ["YES", "NO"],
          "formType": "radio",
          "triggerSection": "medical"
        },
        {
          "property": "medicalConditionTypes",
          "label": "Medical Conditions",
          "type": "checkbox",
          "required": false,
          "conditionalDisplay": "hasMedicalConditions=YES",
          "isMultiSelect": true,
          "enumerationValues": ["Diabetes", "Epilepsy", "Heart_Condition", "Vision_Problems", "Hearing_Problems", "Mobility_Issues", "Mental_Health", "Other"]
        },
        {
          "property": "visionCorrectionRequired",
          "label": "Vision Correction Required?",
          "type": "radio",
          "required": true,
          "enumerationValues": ["YES", "NO"],
          "formType": "radio"
        }
      ]
    }
  }
}
```

**Initial Form State:**
```json
{
  "licenceType": null,
  "hasMedicalConditions": "NO",
  "visionCorrectionRequired": "NO"
}
```

**Rendered Form (Initial):**
```html
<form>
  <!-- Always visible -->
  <div class="field-group">
    <label>Licence Type *</label>
    <select name="licenceType">
      <option value="">Select Licence Type</option>
      <option value="FULL_UK">Full UK</option>
      <option value="PROVISIONAL_UK">Provisional UK</option>
      <option value="EU_EEA">EU EEA</option>
      <option value="INTERNATIONAL">International</option>
      <option value="OTHER_FOREIGN">Other Foreign</option>
    </select>
  </div>

  <!-- HIDDEN: exchangedToUK (conditionalDisplay not met) -->

  <div class="field-group">
    <label>Any Medical Conditions? *</label>
    <div class="radio-group">
      <input type="radio" name="hasMedicalConditions" value="NO" checked />
      <label>No</label>
      <input type="radio" name="hasMedicalConditions" value="YES" />
      <label>Yes</label>
    </div>
  </div>

  <!-- HIDDEN: medicalConditionTypes (conditionalDisplay not met) -->

  <div class="field-group">
    <label>Vision Correction Required? *</label>
    <div class="radio-group">
      <input type="radio" name="visionCorrectionRequired" value="NO" checked />
      <label>No</label>
      <input type="radio" name="visionCorrectionRequired" value="YES" />
      <label>Yes</label>
    </div>
  </div>
</form>
```

### Step 2: User Selects EU Licence

**User Action:** Selects "EU_EEA" from licence type dropdown

**Form State Update:**
```json
{
  "licenceType": "EU_EEA",
  "hasMedicalConditions": "NO",
  "visionCorrectionRequired": "NO"
}
```

**Conditional Logic Evaluation:**
```typescript
// Frontend evaluates: "licenceType=EU_EEA OR licenceType=INTERNATIONAL OR licenceType=OTHER_FOREIGN"
const condition = "licenceType=EU_EEA OR licenceType=INTERNATIONAL OR licenceType=OTHER_FOREIGN";
const orConditions = condition.split(' OR '); // ["licenceType=EU_EEA", "licenceType=INTERNATIONAL", "licenceType=OTHER_FOREIGN"]

// Check first condition: "licenceType=EU_EEA"
const [fieldName, value] = "licenceType=EU_EEA".split('='); // ["licenceType", "EU_EEA"]
const currentValue = formData["licenceType"]; // "EU_EEA"
const shouldShow = currentValue === "EU_EEA"; // true

// Result: exchangedToUK field becomes visible
```

**Updated Rendered Form:**
```html
<form>
  <div class="field-group">
    <label>Licence Type *</label>
    <select name="licenceType">
      <option value="EU_EEA" selected>EU EEA</option>
      <!-- other options -->
    </select>
  </div>

  <!-- NOW VISIBLE: exchangedToUK field appears -->
  <div class="field-group">
    <label>Exchanged To UK Licence?</label>
    <div class="radio-group">
      <input type="radio" name="exchangedToUK" value="NO" />
      <label>No</label>
      <input type="radio" name="exchangedToUK" value="YES" />
      <label>Yes</label>
    </div>
  </div>

  <!-- Rest of form unchanged -->
</form>
```

### Step 3: User Indicates Medical Conditions

**User Action:** Selects "YES" for medical conditions

**Form State Update:**
```json
{
  "licenceType": "EU_EEA",
  "exchangedToUK": "NO",
  "hasMedicalConditions": "YES",
  "visionCorrectionRequired": "NO"
}
```

**Conditional Logic Evaluation:**
```typescript
// Frontend evaluates: "hasMedicalConditions=YES"
const condition = "hasMedicalConditions=YES";
const [fieldName, value] = condition.split('='); // ["hasMedicalConditions", "YES"]
const currentValue = formData["hasMedicalConditions"]; // "YES"
const shouldShow = currentValue === "YES"; // true

// Result: medicalConditionTypes field becomes visible
```

**Updated Rendered Form:**
```html
<form>
  <!-- Previous fields unchanged -->

  <div class="field-group">
    <label>Any Medical Conditions? *</label>
    <div class="radio-group">
      <input type="radio" name="hasMedicalConditions" value="NO" />
      <label>No</label>
      <input type="radio" name="hasMedicalConditions" value="YES" checked />
      <label>Yes</label>
    </div>
  </div>

  <!-- NOW VISIBLE: Medical conditions multi-select appears -->
  <div class="field-group">
    <label>Medical Conditions</label>
    <div class="checkbox-group">
      <label class="checkbox-item">
        <input type="checkbox" name="medicalConditionTypes" value="Diabetes" />
        <span>Diabetes</span>
      </label>
      <label class="checkbox-item">
        <input type="checkbox" name="medicalConditionTypes" value="Epilepsy" />
        <span>Epilepsy</span>
      </label>
      <label class="checkbox-item">
        <input type="checkbox" name="medicalConditionTypes" value="Vision_Problems" />
        <span>Vision Problems</span>
      </label>
      <!-- More options... -->
    </div>
  </div>
</form>
```

### Step 4: User Completes Medical Information

**User Actions:**
- Selects "Diabetes" and "Vision_Problems" checkboxes
- Selects "YES" for vision correction required

**Final Form State:**
```json
{
  "licenceType": "EU_EEA",
  "exchangedToUK": "NO",
  "hasMedicalConditions": "YES",
  "medicalConditionTypes": ["Diabetes", "Vision_Problems"],
  "visionCorrectionRequired": "YES"
}
```

### Step 5: SHACL Validation

**Validation Rules Applied:**

1. **Basic Required Fields:**
```turtle
sh:property [
  sh:path autoins:licenceType ;
  sh:minCount 1 ;
  sh:message "Licence type is required"
] ;
```
✅ **PASS**: `licenceType = "EU_EEA"`

2. **Conditional Field Requirements:**
```turtle
autoins:NonUKLicenceExchangeShape a sh:NodeShape ;
  sh:condition [
    sh:property [
      sh:path autoins:licenceType ;
      sh:in ("EU_EEA" "INTERNATIONAL" "OTHER_FOREIGN")
    ]
  ] ;
  sh:property [
    sh:path autoins:exchangedToUK ;
    sh:minCount 1 ;
    sh:message "Non-UK licences must specify if exchanged to UK licence"
  ] .
```
✅ **PASS**: `exchangedToUK = "NO"` (required because licenceType is EU_EEA)

3. **Medical Conditions Validation:**
```turtle
autoins:MedicalDVLAConditionalShape a sh:NodeShape ;
  sh:condition [
    sh:property [
      sh:path autoins:hasMedicalConditions ;
      sh:hasValue "YES"
    ]
  ] ;
  sh:property [
    sh:path autoins:medicalDeclaredToDVLA ;
    sh:minCount 1 ;
    sh:message "Medical DVLA declaration required when hasMedicalConditions is YES"
  ] .
```
❌ **FAIL**: `medicalDeclaredToDVLA` is missing (would be required)

**Validation Result:**
```json
{
  "isValid": false,
  "errors": [
    {
      "field": "medicalDeclaredToDVLA",
      "message": "Medical DVLA declaration required when hasMedicalConditions is YES",
      "code": "CONDITIONAL_REQUIRED"
    }
  ]
}
```

### Step 6: Form Completion

**Additional Field Appears:**
The validation error triggers the display of the missing required field:

```html
<div class="field-group">
  <label>Medical Conditions Declared To DVLA? *</label>
  <div class="radio-group">
    <input type="radio" name="medicalDeclaredToDVLA" value="NO" />
    <label>No</label>
    <input type="radio" name="medicalDeclaredToDVLA" value="YES" />
    <label>Yes</label>
  </div>
</div>
```

**User selects "YES"**

**Final Valid Form State:**
```json
{
  "licenceType": "EU_EEA",
  "licenceCategory": ["B"],
  "manualOrAuto": "MANUAL",
  "dateFirstIssued": "2018-03-15",
  "countryOfIssue": "Germany",
  "exchangedToUK": "NO",
  "hasMedicalConditions": "YES",
  "medicalConditionTypes": ["Diabetes", "Vision_Problems"],
  "medicalDeclaredToDVLA": "YES",
  "visionCorrectionRequired": "YES"
}
```

## Backend Processing

### Step 7: Form Submission

**Frontend POST Request:**
```bash
POST /api/driver-licence
Content-Type: application/json

{
  "licenceType": "EU_EEA",
  "licenceCategory": ["B"],
  "manualOrAuto": "MANUAL",
  "dateFirstIssued": "2018-03-15",
  "countryOfIssue": "Germany",
  "exchangedToUK": "NO",
  "hasMedicalConditions": "YES",
  "medicalConditionTypes": ["Diabetes", "Vision_Problems"],
  "medicalDeclaredToDVLA": "YES",
  "visionCorrectionRequired": "YES"
}
```

**Backend Validation (Go):**
```go
func validateDriverLicence(data DriverLicenceData) []ValidationError {
    var errors []ValidationError
    
    // JSON Schema validation
    if data.LicenceType == "" {
        errors = append(errors, ValidationError{
            Field: "licenceType",
            Message: "Licence type is required",
        })
    }
    
    // Conditional validation
    if isNonUKLicence(data.LicenceType) && data.ExchangedToUK == "" {
        errors = append(errors, ValidationError{
            Field: "exchangedToUK", 
            Message: "Non-UK licences must specify if exchanged to UK licence",
        })
    }
    
    // Medical conditions validation
    if data.HasMedicalConditions == "YES" && data.MedicalDeclaredToDVLA == "" {
        errors = append(errors, ValidationError{
            Field: "medicalDeclaredToDVLA",
            Message: "Medical DVLA declaration required when hasMedicalConditions is YES",
        })
    }
    
    return errors
}
```

**Successful Response:**
```json
{
  "success": true,
  "driverId": "driver_12345",
  "validationPassed": true,
  "riskAssessment": {
    "medicalConditionsRisk": "MEDIUM",
    "foreignLicenceRisk": "LOW",
    "overallRisk": "MEDIUM"
  }
}
```

## Key Takeaways

1. **Dynamic Field Display**: Fields appear/disappear based on user input without page reloads
2. **Automatic Validation**: SHACL rules enforce business logic automatically
3. **JSON Schema Compliance**: Industry-standard validation patterns
4. **Single Source of Truth**: Ontology defines both UI and validation rules
5. **Type Safety**: Full end-to-end type checking from ontology to database
6. **Conditional Complexity**: Supports complex OR conditions and nested dependencies
7. **Real-time Feedback**: Immediate validation feedback as user progresses

This example demonstrates how the ontology-driven approach handles complex form logic while maintaining clean, maintainable code and ensuring compliance with UK insurance requirements.
