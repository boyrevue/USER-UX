import React, { useState, useEffect } from 'react';
import { Card, Alert } from 'flowbite-react';
import { AlertCircle, CheckCircle } from 'lucide-react';
import { 
  DriverLicence, 
  LicenceFormState, 
  LicenceValidationError,
  LICENCE_VALIDATION_PATTERNS,
  LICENCE_CONSTRAINTS,
  requiresExchangeField,
  requiresEndorsementDetails,
  requiresDisqualificationDetails,
  requiresMedicalDVLADeclaration
} from '../../types/licence';

interface LicenceFormHandlerProps {
  initialData?: Partial<DriverLicence>;
  onSubmit: (data: DriverLicence) => void;
  onValidationChange?: (isValid: boolean, errors: LicenceValidationError[]) => void;
}

const LicenceFormHandler: React.FC<LicenceFormHandlerProps> = ({
  initialData = {},
  onSubmit,
  onValidationChange
}) => {
  const [formState, setFormState] = useState<LicenceFormState>({
    data: initialData,
    errors: [],
    touched: new Set(),
    isSubmitting: false
  });

  // JSON Schema Validation Implementation
  const validateField = (fieldName: string, value: any): LicenceValidationError[] => {
    const errors: LicenceValidationError[] = [];
    
    switch (fieldName) {
      case 'licenceType':
        if (!value) {
          errors.push({
            field: fieldName,
            message: 'Licence type is required',
            code: 'REQUIRED'
          });
        } else if (!['FULL_UK', 'PROVISIONAL_UK', 'EU_EEA', 'INTERNATIONAL', 'OTHER_FOREIGN'].includes(value)) {
          errors.push({
            field: fieldName,
            message: 'Invalid licence type',
            code: 'INVALID_ENUM'
          });
        }
        break;
        
      case 'licenceCategory':
        if (!value || !Array.isArray(value) || value.length === 0) {
          errors.push({
            field: fieldName,
            message: 'At least one licence category is required',
            code: 'MIN_ITEMS'
          });
        } else {
          const validCategories = ['B', 'B1', 'A', 'AM', 'Q', 'C', 'C1', 'D', 'D1', 'BE', 'C1E', 'CE', 'D1E', 'DE'];
          const invalidCategories = value.filter(cat => !validCategories.includes(cat));
          if (invalidCategories.length > 0) {
            errors.push({
              field: fieldName,
              message: `Invalid licence categories: ${invalidCategories.join(', ')}`,
              code: 'INVALID_ENUM'
            });
          }
        }
        break;
        
      case 'dateFirstIssued':
        if (!value) {
          errors.push({
            field: fieldName,
            message: 'Date first issued is required',
            code: 'REQUIRED'
          });
        } else if (!LICENCE_VALIDATION_PATTERNS.DATE.test(value)) {
          errors.push({
            field: fieldName,
            message: 'Date must be in YYYY-MM-DD format',
            code: 'INVALID_PATTERN'
          });
        }
        break;
        
      case 'yearsHeldFull':
        if (value !== undefined && value !== null) {
          if (typeof value !== 'number' || value < LICENCE_CONSTRAINTS.YEARS_HELD_MIN || value > LICENCE_CONSTRAINTS.YEARS_HELD_MAX) {
            errors.push({
              field: fieldName,
              message: `Years held must be between ${LICENCE_CONSTRAINTS.YEARS_HELD_MIN} and ${LICENCE_CONSTRAINTS.YEARS_HELD_MAX}`,
              code: 'OUT_OF_RANGE'
            });
          }
        }
        break;
        
      case 'countryOfIssue':
        if (!value) {
          errors.push({
            field: fieldName,
            message: 'Country of issue is required',
            code: 'REQUIRED'
          });
        } else if (value.length < LICENCE_CONSTRAINTS.COUNTRY_MIN_LENGTH) {
          errors.push({
            field: fieldName,
            message: `Country must be at least ${LICENCE_CONSTRAINTS.COUNTRY_MIN_LENGTH} characters`,
            code: 'MIN_LENGTH'
          });
        }
        break;
        
      case 'licenceNumber':
        if (value && !LICENCE_VALIDATION_PATTERNS.LICENCE_NUMBER.test(value)) {
          errors.push({
            field: fieldName,
            message: 'Licence number must match DVLA format (16 characters)',
            code: 'INVALID_PATTERN'
          });
        }
        break;
        
      case 'endorsements':
        if (requiresEndorsementDetails(formState.data.hasEndorsements || 'NO')) {
          if (!value || !Array.isArray(value) || value.length === 0) {
            errors.push({
              field: fieldName,
              message: 'Endorsement details are required when you have endorsements',
              code: 'CONDITIONAL_REQUIRED'
            });
          } else {
            value.forEach((endorsement, index) => {
              if (!endorsement.code || !LICENCE_VALIDATION_PATTERNS.DVLA_CODE.test(endorsement.code)) {
                errors.push({
                  field: `${fieldName}[${index}].code`,
                  message: 'DVLA code must be in format XX00 (e.g., SP30)',
                  code: 'INVALID_PATTERN'
                });
              }
              if (endorsement.points < LICENCE_CONSTRAINTS.POINTS_MIN || endorsement.points > LICENCE_CONSTRAINTS.POINTS_MAX) {
                errors.push({
                  field: `${fieldName}[${index}].points`,
                  message: `Points must be between ${LICENCE_CONSTRAINTS.POINTS_MIN} and ${LICENCE_CONSTRAINTS.POINTS_MAX}`,
                  code: 'OUT_OF_RANGE'
                });
              }
              if (!endorsement.offenceDate || !LICENCE_VALIDATION_PATTERNS.DATE.test(endorsement.offenceDate)) {
                errors.push({
                  field: `${fieldName}[${index}].offenceDate`,
                  message: 'Offence date must be in YYYY-MM-DD format',
                  code: 'INVALID_PATTERN'
                });
              }
            });
          }
        }
        break;
    }
    
    return errors;
  };

  // Comprehensive form validation using JSON Schema rules
  const validateForm = (data: Partial<DriverLicence>): LicenceValidationError[] => {
    let allErrors: LicenceValidationError[] = [];
    
    // Validate all fields
    Object.keys(data).forEach(fieldName => {
      const fieldErrors = validateField(fieldName, data[fieldName as keyof DriverLicence]);
      allErrors = allErrors.concat(fieldErrors);
    });
    
    // JSON Schema "allOf" conditional validations
    
    // 1. UK licence must have UK country
    if (data.licenceType === 'FULL_UK' && data.countryOfIssue !== 'UK') {
      allErrors.push({
        field: 'countryOfIssue',
        message: 'UK licence must have UK as country of issue',
        code: 'CONDITIONAL_CONSTRAINT'
      });
    }
    
    // 2. Non-UK licences require exchangedToUK field
    if (requiresExchangeField(data.licenceType as any) && !data.exchangedToUK) {
      allErrors.push({
        field: 'exchangedToUK',
        message: 'Non-UK licences must specify if exchanged to UK licence',
        code: 'CONDITIONAL_REQUIRED'
      });
    }
    
    // 3. Medical conditions require DVLA declaration
    if (requiresMedicalDVLADeclaration(data.hasMedicalConditions as any) && !data.medicalDeclaredToDVLA) {
      allErrors.push({
        field: 'medicalDeclaredToDVLA',
        message: 'Medical DVLA declaration required when you have medical conditions',
        code: 'CONDITIONAL_REQUIRED'
      });
    }
    
    // 4. Disqualifications require details
    if (requiresDisqualificationDetails(data.hasDisqualifications as any)) {
      if (!data.disqualifications || data.disqualifications.length === 0) {
        allErrors.push({
          field: 'disqualifications',
          message: 'Disqualification details are required when you have disqualifications',
          code: 'CONDITIONAL_REQUIRED'
        });
      }
    }
    
    return allErrors;
  };

  // Update form data with validation
  const updateField = (fieldName: string, value: any) => {
    const newData = { ...formState.data, [fieldName]: value };
    const errors = validateForm(newData);
    
    setFormState(prev => ({
      ...prev,
      data: newData,
      errors,
      touched: new Set(Array.from(prev.touched).concat([fieldName]))
    }));
    
    // Notify parent of validation changes
    if (onValidationChange) {
      onValidationChange(errors.length === 0, errors);
    }
  };

  // Handle form submission
  const handleSubmit = () => {
    const errors = validateForm(formState.data);
    
    if (errors.length === 0) {
      setFormState(prev => ({ ...prev, isSubmitting: true }));
      onSubmit(formState.data as DriverLicence);
    } else {
      setFormState(prev => ({
        ...prev,
        errors,
        touched: new Set(Object.keys(formState.data))
      }));
    }
  };

  // Get field-specific errors
  const getFieldErrors = (fieldName: string): LicenceValidationError[] => {
    return formState.errors.filter(error => error.field === fieldName);
  };

  // Check if field should be displayed based on conditional logic
  const shouldDisplayField = (fieldName: string): boolean => {
    switch (fieldName) {
      case 'exchangedToUK':
        return requiresExchangeField(formState.data.licenceType as any);
      case 'endorsements':
        return requiresEndorsementDetails(formState.data.hasEndorsements as any);
      case 'disqualifications':
        return requiresDisqualificationDetails(formState.data.hasDisqualifications as any);
      case 'medicalDeclaredToDVLA':
        return requiresMedicalDVLADeclaration(formState.data.hasMedicalConditions as any);
      default:
        return true;
    }
  };

  return (
    <Card className="w-full">
      <div className="space-y-6">
        <div className="flex items-center justify-between">
          <h3 className="text-lg font-semibold text-gray-900">Driver Licence Information</h3>
          <div className="flex items-center space-x-2">
            {formState.errors.length === 0 ? (
              <div className="flex items-center text-green-600">
                <CheckCircle className="w-4 h-4 mr-1" />
                <span className="text-sm">Valid</span>
              </div>
            ) : (
              <div className="flex items-center text-red-600">
                <AlertCircle className="w-4 h-4 mr-1" />
                <span className="text-sm">{formState.errors.length} errors</span>
              </div>
            )}
          </div>
        </div>

        {/* Validation Errors Summary */}
        {formState.errors.length > 0 && (
          <Alert color="failure">
            <AlertCircle className="h-4 w-4" />
            <span className="ml-2 font-medium">Please fix the following errors:</span>
            <ul className="mt-2 ml-4 list-disc">
              {formState.errors.map((error, index) => (
                <li key={index} className="text-sm">
                  <strong>{error.field}:</strong> {error.message}
                </li>
              ))}
            </ul>
          </Alert>
        )}

        {/* Form fields would be rendered here using the ontology-driven UniversalForm */}
        <div className="text-sm text-gray-600">
          <p><strong>JSON Schema Validation Active:</strong></p>
          <ul className="list-disc ml-4 mt-1">
            <li>Real-time field validation</li>
            <li>Conditional field requirements</li>
            <li>Pattern matching (dates, DVLA codes, licence numbers)</li>
            <li>Array validation (categories, endorsements)</li>
            <li>Cross-field validation rules</li>
          </ul>
        </div>

        {/* Debug Info (remove in production) */}
        <details className="text-xs text-gray-500">
          <summary>Debug: Current Form State</summary>
          <pre className="mt-2 p-2 bg-gray-100 rounded overflow-auto">
            {JSON.stringify(formState, null, 2)}
          </pre>
        </details>
      </div>
    </Card>
  );
};

export default LicenceFormHandler;
