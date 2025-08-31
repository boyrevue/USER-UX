# Ontology-Driven Trigger Fields System

## Overview

The CLIENT-UX application uses a completely ontology-driven system for managing trigger fields and conditional form display. This system automatically detects trigger fields, manages their conditional dependencies, and renders appropriate UI sections based on ontology definitions.

## How It Works

### 1. Trigger Field Detection

A field is automatically identified as a "trigger field" if other fields in the ontology depend on it through `conditionalDisplay` properties.

```typescript
// Automatic detection - no hardcoding required
const isTriggerField = (field: OntologyField): boolean => {
  return hasConditionalFields(field.property);
};
```

### 2. Supported Conditional Patterns

The system supports all conditional display patterns defined in the ontology:

#### Simple Conditions
```ttl
autoins:conditionalDisplay "isMainDriver=NO" ;
autoins:conditionalDisplay "hasAccidents=YES" ;
```

#### Complex AND Conditions
```ttl
autoins:conditionalDisplay "licenceEverRevoked=YES AND licenceReinstated=YES" ;
autoins:conditionalDisplay "hasAccidents=YES AND isMainDriver=NO" ;
```

#### OR Conditions
```ttl
autoins:conditionalDisplay "licenceType=EU_EEA OR licenceType=INTERNATIONAL OR licenceType=OTHER_FOREIGN" ;
```

#### Includes Conditions (for multi-select fields)
```ttl
autoins:conditionalDisplay "reinstatementConditions_includes=Other" ;
```

#### Not Equals Conditions
```ttl
autoins:conditionalDisplay "driverType!=MAIN_DRIVER" ;
```

### 3. Automatic UI Rendering

When a trigger field is detected, the system:

1. **Renders the trigger field normally**
2. **Monitors its value changes**
3. **Automatically shows/hides dependent fields** based on ontology conditions
4. **Groups dependent fields in bordered sections** for visual clarity
5. **Maintains trigger field visibility** so users can always change their selection

## Ontology Properties for Enhanced Control

### Core Properties
- `autoins:conditionalDisplay` - Defines when a field should be visible
- `autoins:conditionalRequirement` - Defines when a field becomes required

### Advanced Trigger Control
- `autoins:triggerSection` - Groups triggered fields under a named section
- `autoins:sectionCollapsible` - Makes the triggered section collapsible
- `autoins:sectionDefaultClosed` - Sets default collapsed state

### Example: Medical Conditions Trigger
```ttl
autoins:hasMedicalConditions a owl:DatatypeProperty ;
  rdfs:label "Any Medical Conditions Affecting Driving?" ;
  autoins:enumerationValues ("YES" "NO") ;
  autoins:formType "radio" ;
  autoins:defaultValue "NO" ;
  autoins:triggerSection "medical" ;
  autoins:sectionCollapsible "true"^^xsd:boolean ;
  autoins:sectionDefaultClosed "true"^^xsd:boolean .

autoins:medicalConditionTypes a owl:DatatypeProperty ;
  rdfs:label "Medical Conditions" ;
  autoins:conditionalDisplay "hasMedicalConditions=YES" ;
  autoins:isMultiSelect "true"^^xsd:boolean ;
  autoins:enumerationValues ("Diabetes" "Epilepsy" "Heart_Condition" ...) .
```

## Benefits

### 1. **Completely Ontology-Driven**
- No hardcoded field names or conditions in the frontend
- All logic defined in TTL files
- Easy to add new trigger fields without code changes

### 2. **Supports All Condition Types**
- Simple equality/inequality
- Complex AND/OR logic
- Multi-select includes conditions
- Nested dependencies

### 3. **Automatic UI Management**
- Trigger fields always remain visible and accessible
- Conditional fields appear in clearly bordered sections
- Smooth show/hide transitions
- Proper field grouping and labeling

### 4. **Consistent User Experience**
- Users can always see and change their trigger field selections
- Clear visual indication of conditional sections
- Intuitive field organization

## Implementation Details

### Frontend Logic (UniversalForm.tsx)

```typescript
// Detect trigger fields automatically
const isTriggerField = (field: OntologyField): boolean => {
  return hasConditionalFields(field.property);
};

// Find all fields that depend on a trigger
const getTriggeredFields = (triggerProperty: string): OntologyField[] => {
  return fields.filter(field => {
    if (!field.conditionalDisplay) return false;
    
    const condition = field.conditionalDisplay;
    const referencesThisField = condition.includes(triggerProperty + '=') || 
                               condition.includes(triggerProperty + '!=') ||
                               condition.includes(triggerProperty + '_includes=');
    
    return referencesThisField && shouldDisplayField(field);
  });
};

// Render with proper grouping
const renderFieldsWithTriggerGrouping = (fieldsToRender: OntologyField[]) => {
  // Automatically groups trigger fields with their dependents
  // Uses existing shouldDisplayField() logic for all conditions
  // Renders trigger field + bordered section for dependents
};
```

### Ontology Parsing (ttl_parser.go)

The system automatically extracts all conditional display properties:

```go
conditionalDisplayPattern := regexp.MustCompile(`(autoins|docs|):conditionalDisplay\s+"([^"]+)"\s*;`)
```

## Adding New Trigger Fields

To add a new trigger field, simply define it in the ontology:

1. **Create the trigger field:**
```ttl
autoins:hasNewFeature a owl:DatatypeProperty ;
  rdfs:label "Do you have this new feature?" ;
  autoins:enumerationValues ("YES" "NO") ;
  autoins:formType "radio" ;
  autoins:defaultValue "NO" .
```

2. **Add dependent fields:**
```ttl
autoins:newFeatureDetails a owl:DatatypeProperty ;
  rdfs:label "Feature Details" ;
  autoins:conditionalDisplay "hasNewFeature=YES" ;
  autoins:formType "textarea" .
```

3. **That's it!** The system automatically:
   - Detects the trigger relationship
   - Renders the UI properly
   - Manages show/hide logic
   - Maintains accessibility

## Testing

The system has been tested with all existing trigger fields:

- ✅ `isMainDriver` (NO condition) → `relationship`, `driverType`
- ✅ `hasMedicalConditions` (YES condition) → `medicalConditionTypes`
- ✅ `hasAccidents` (YES condition) → `accidentDate`, `whoWasDriving`, etc.
- ✅ `hasEndorsements` (YES condition) → `endorsementOccurrences`
- ✅ `hasDisqualifications` (YES condition) → disqualification details
- ✅ `licenceEverRevoked` (YES condition) → revocation details
- ✅ `licenceEverRestricted` (YES condition) → `licenceRestrictionCodes`
- ✅ `hasModifications` (YES condition) → modification details
- ✅ Complex conditions with AND/OR logic

## Future Enhancements

1. **Nested Trigger Fields**: Support for triggers within triggered sections
2. **Dynamic Section Titles**: Use ontology properties for section naming
3. **Conditional Validation**: Extend to conditional field validation rules
4. **Animation Controls**: Ontology properties for transition animations
