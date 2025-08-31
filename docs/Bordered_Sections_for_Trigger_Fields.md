# Bordered Sections for Yes/No Trigger Fields

## Overview

Implemented visual grouping with bordered sections for Yes/No fields that trigger additional form fields. This improves user experience by clearly showing which fields are related to each other.

## The Problem

Previously, when a user selected "YES" to a question like "Any Medical Conditions?", the additional fields would appear but it wasn't visually clear that they were related to the trigger question. This created confusion about which fields belonged together.

## The Solution

### ✅ **Visual Grouping with Bordered Sections**

When a Yes/No radio button field triggers additional fields, they are now grouped together in a visually distinct bordered section with:

1. **Blue border** around the conditional fields
2. **Light blue background** to distinguish the section
3. **Section header** with the trigger field name
4. **Visual separator line** under the header

## Implementation Details

### **Frontend Logic (UniversalForm.tsx)**

**New Helper Functions:**
```typescript
// Check if a field is a Yes/No trigger field that shows additional fields
const isTriggerField = (field: OntologyField): boolean => {
  return field.formType === 'radio' && 
         (field.enumerationValues?.includes('YES') || false) && 
         (field.enumerationValues?.includes('NO') || false) &&
         hasConditionalFields(field.property);
};

// Check if there are fields that depend on this trigger field
const hasConditionalFields = (triggerProperty: string): boolean => {
  return fields.some(field => 
    field.conditionalDisplay?.includes(`${triggerProperty}=YES`)
  );
};

// Get fields that are triggered by a specific field
const getTriggeredFields = (triggerProperty: string): OntologyField[] => {
  return fields.filter(field => 
    field.conditionalDisplay?.includes(`${triggerProperty}=YES`) &&
    shouldDisplayField(field)
  );
};
```

**New Rendering Function:**
```typescript
const renderFieldsWithTriggerGrouping = (fieldsToRender: OntologyField[]) => {
  const renderedFields: React.ReactElement[] = [];
  const processedFields = new Set<string>();

  fieldsToRender.forEach((field) => {
    if (processedFields.has(field.property)) return;

    if (isTriggerField(field)) {
      // Render trigger field with its conditional fields in a bordered section
      const triggeredFields = getTriggeredFields(field.property);
      const triggerValue = (data as any)[field.property];
      const showTriggeredFields = triggerValue === 'YES' || triggerValue === true;

      renderedFields.push(
        <div key={field.property} className="space-y-4">
          {/* Render the trigger field */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {renderField(field)}
          </div>
          
          {/* Render triggered fields in a bordered section if YES is selected */}
          {showTriggeredFields && triggeredFields.length > 0 && (
            <div className="border-2 border-blue-200 rounded-lg p-4 bg-blue-50/30">
              <div className="mb-3">
                <h5 className="text-sm font-semibold text-blue-800 uppercase tracking-wide">
                  {field.label?.replace(/\?$/, '')} Details
                </h5>
                <div className="w-12 h-0.5 bg-blue-400 mt-1"></div>
              </div>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                {triggeredFields.map((triggeredField) => {
                  processedFields.add(triggeredField.property);
                  return renderField(triggeredField);
                })}
              </div>
            </div>
          )}
        </div>
      );
    }
  });

  return renderedFields;
};
```

## Visual Examples

### **Before (Confusing):**
```
□ Any Medical Conditions?  ○ No  ● Yes

Medical Condition Types
☐ Diabetes  ☐ Epilepsy  ☐ Heart Condition

Medical Conditions Declared To DVLA?  ○ No  ● Yes

Vision Correction Required?  ○ No  ● Yes
```

### **After (Clear Grouping):**
```
□ Any Medical Conditions?  ○ No  ● Yes

┌─────────────────────────────────────────────────────────┐
│ ANY MEDICAL CONDITIONS DETAILS                          │
│ ────                                                    │
│                                                         │
│ Medical Condition Types                                 │
│ ☐ Diabetes  ☐ Epilepsy  ☐ Heart Condition             │
│                                                         │
│ Medical Conditions Declared To DVLA?  ○ No  ● Yes      │
└─────────────────────────────────────────────────────────┘

□ Vision Correction Required?  ○ No  ● Yes
```

## Trigger Fields Affected

The following Yes/No fields now create bordered sections when "YES" is selected:

### **1. Medical Conditions**
- **Trigger:** `hasMedicalConditions` = YES
- **Shows:** Medical condition types, DVLA declaration

### **2. Endorsements/Penalty Points**
- **Trigger:** `hasEndorsements` = YES  
- **Shows:** DVLA codes, points, offence dates

### **3. Disqualifications**
- **Trigger:** `hasDisqualifications` = YES
- **Shows:** Start/end dates, reason, duration

### **4. Accidents/Claims**
- **Trigger:** `hasAccidents` = YES
- **Shows:** Accident date, fault determination

### **5. Vehicle Modifications**
- **Trigger:** `hasModifications` = YES
- **Shows:** Modification types, values, declarations

## CSS Classes Used

```css
/* Bordered section container */
.border-2.border-blue-200.rounded-lg.p-4.bg-blue-50/30

/* Section header */
.text-sm.font-semibold.text-blue-800.uppercase.tracking-wide

/* Header underline */
.w-12.h-0.5.bg-blue-400.mt-1

/* Field grid inside section */
.grid.grid-cols-1.md:grid-cols-2.gap-4
```

## Benefits

### **✅ Improved User Experience:**
1. **Visual Clarity** - Users can immediately see which fields are related
2. **Reduced Confusion** - Clear separation between different question groups
3. **Better Flow** - Logical progression through related questions
4. **Professional Appearance** - Clean, modern UI design

### **✅ Technical Benefits:**
1. **Maintainable Code** - Automatic detection of trigger fields
2. **Scalable Solution** - Works with any Yes/No → conditional fields pattern
3. **Consistent Styling** - Uniform appearance across all trigger sections
4. **Responsive Design** - Works on mobile and desktop

## Usage Rule

**Remember this rule:** Wherever a Yes/No field introduces more fields, create a separate div and border around this in styling - it's confusing otherwise.

This feature automatically implements this rule for all qualifying fields in the ontology-driven form system.

## Server Status

**✅ Feature Active:** http://localhost:3000

The bordered sections are now live and automatically applied to all Yes/No trigger fields throughout the insurance application forms.
