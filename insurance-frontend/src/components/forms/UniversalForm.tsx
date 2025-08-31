import React, { useEffect, useState } from 'react';
import { Card, Button, TextInput, Label, Select, Badge } from 'flowbite-react';
import { X, Plus, MessageCircle } from 'lucide-react';
import { OntologyField } from '../../types';
import AIValidationDialog from './AIValidationDialog';

interface UniversalFormProps {
  title: string;
  icon?: React.ReactNode;
  fields: OntologyField[];
  data: any;
  onUpdate: (field: string, value: any) => void;
  onRemove?: () => void;
  showRemoveButton?: boolean;
  validationErrors?: Array<{ field: string; message: string }>;
  formType?: 'driver' | 'vehicle' | 'claims' | 'policy' | 'settings' | 'documents';
}

const UniversalForm: React.FC<UniversalFormProps> = ({
  title,
  icon,
  fields,
  data,
  onUpdate,
  onRemove,
  showRemoveButton = false,
  validationErrors = [],
  formType = 'general'
}) => {
  const [collapsedSections, setCollapsedSections] = useState<{[key: string]: boolean}>({
    'Convictions & Penalty Points': true, // Default closed
    'Disabilities & Adaptations': true, // Default closed
    'Medical Conditions & Restrictions': true // Default closed
  });

  const [aiValidationDialog, setAiValidationDialog] = useState<{
    isOpen: boolean;
    fieldName: string;
    fieldLabel: string;
    initialValue: string;
    validationPrompt: string;
  }>({
    isOpen: false,
    fieldName: '',
    fieldLabel: '',
    initialValue: '',
    validationPrompt: '',
  });

  // Track which AI fields we already auto-opened to avoid loops
  const [autoOpenedFields, setAutoOpenedFields] = useState<Set<string>>(new Set());

  // Auto-open AI dialog when an AI-validated field becomes visible with no value
  useEffect(() => {
    const candidate = fields.find((f) => (f as any).requiresAIValidation && shouldDisplayField(f) && !(data as any)[f.property]);
    if (candidate && !autoOpenedFields.has(candidate.property)) {
      setAutoOpenedFields(prev => new Set(prev).add(candidate.property));
      openAIValidationDialog(candidate, '');
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data, fields]);

  // Listen for global assistant trigger to open the first relevant AI field
  useEffect(() => {
    const handler = () => {
      const aiFields = fields.filter((f) => (f as any).requiresAIValidation);
      const candidate = aiFields.find((f) => shouldDisplayField(f) && !(data as any)[f.property]) || aiFields.find(shouldDisplayField);
      if (candidate) {
        openAIValidationDialog(candidate, (data as any)[candidate.property] || '');
      } else {
        // No AI field visible right now; open a general assistant dialog
        setAiValidationDialog({
          isOpen: true,
          fieldName: 'generalHelp',
          fieldLabel: 'Assistant',
          initialValue: '',
          validationPrompt: 'You are an insurance assistant. Help the user provide precise, insurance-relevant details required by underwriters. Ask short follow-up questions and keep responses structured.'
        });
      }
    };
    window.addEventListener('open-ai-validation', handler as EventListener);
    return () => window.removeEventListener('open-ai-validation', handler as EventListener);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [fields, data]);

  const toggleSection = (sectionName: string) => {
    setCollapsedSections(prev => ({
      ...prev,
      [sectionName]: !prev[sectionName]
    }));
  };

  const openAIValidationDialog = (field: OntologyField, currentValue: string) => {
    setAiValidationDialog({
      isOpen: true,
      fieldName: field.property,
      fieldLabel: field.label,
      initialValue: currentValue || '',
      validationPrompt: (field as any).aiValidationPrompt || 'Please provide detailed, insurance-relevant information.',
    });
  };

  const triggerAIForNewData = (newData: any, triggerField?: OntologyField) => {
    // Find any newly-visible AI-validated field with no value
    const aiCandidate = fields.find((f) => (f as any).requiresAIValidation && shouldDisplayField(f) && !newData[f.property]);
    if (aiCandidate) {
      openAIValidationDialog(aiCandidate, '');
      return;
    }
    // If a trigger field is provided, try dependent fields first
    if (triggerField && isTriggerField(triggerField)) {
      const dependent = getTriggeredFields(triggerField.property)
        .find((dep) => (dep as any).requiresAIValidation && shouldDisplayField(dep) && !newData[(dep as any).property]);
      if (dependent) {
        openAIValidationDialog(dependent, '');
      }
    }
  };

  const closeAIValidationDialog = () => {
    setAiValidationDialog({
      isOpen: false,
      fieldName: '',
      fieldLabel: '',
      initialValue: '',
      validationPrompt: '',
    });
  };

  const handleValidatedValue = (value: string) => {
    onUpdate(aiValidationDialog.fieldName, value);
    closeAIValidationDialog();
  };

  const shouldDisplayField = (field: OntologyField): boolean => {
    // Check if field has conditional display rules
    if (field.conditionalDisplay) {
      const condition = field.conditionalDisplay;
      
      // Handle AND conditions like "hasAccidents=YES AND driverType!=MAIN_DRIVER"
      if (condition.includes(' AND ')) {
        const andConditions = condition.split(' AND ');
        return andConditions.every(andCondition => evaluateCondition(andCondition.trim(), data));
      }
      
      // Handle OR conditions like "licenceType=EU_EEA OR licenceType=INTERNATIONAL OR licenceType=OTHER_FOREIGN"
      if (condition.includes(' OR ')) {
        const orConditions = condition.split(' OR ');
        return orConditions.some(orCondition => evaluateCondition(orCondition.trim(), data));
      }
      
      // Handle single conditions
      return evaluateCondition(condition, data);
    }
    
    return true; // Show field by default
  };

  const evaluateCondition = (condition: string, data: any): boolean => {
    if (condition.includes('_includes=')) {
      const [fieldName, value] = condition.split('_includes=');
      const currentValue = data[fieldName.trim()];
      const expectedValue = value.trim();
      
      // Handle array/multi-select fields
      if (Array.isArray(currentValue)) {
        return currentValue.includes(expectedValue);
      }
      // Handle comma-separated string values
      if (typeof currentValue === 'string' && currentValue.includes(',')) {
        return currentValue.split(',').map(v => v.trim()).includes(expectedValue);
      }
      return false;
    } else if (condition.includes('!=')) {
      const [fieldName, value] = condition.split('!=');
      const currentValue = data[fieldName.trim()];
      const expectedValue = value.trim();
      
      if (expectedValue === 'null') {
        return currentValue && currentValue !== '' && currentValue !== null && currentValue !== undefined;
      } else if (expectedValue === 'MAIN_DRIVER') {
        return currentValue !== 'MAIN_DRIVER';
      } else {
        return currentValue !== expectedValue;
      }
    } else if (condition.includes('=')) {
      const [fieldName, value] = condition.split('=');
      const currentValue = data[fieldName.trim()];
      const expectedValue = value.trim();
      
      if (expectedValue === 'null') {
        return !currentValue || currentValue === '' || currentValue === null || currentValue === undefined;
      } else if (expectedValue === 'YES') {
        return currentValue === 'YES' || currentValue === true;
      } else if (expectedValue === 'NO') {
        return currentValue === 'NO' || currentValue === false;
      } else {
        return currentValue === expectedValue;
      }
    }
    
    return true;
  };

  // Check if a field is a trigger field that shows additional fields
  const isTriggerField = (field: OntologyField): boolean => {
    // A field is a trigger if other fields depend on it
    return hasConditionalFields(field.property);
  };

  // Check if there are fields that depend on this trigger field
  const hasConditionalFields = (triggerProperty: string): boolean => {
    return fields.some(field => {
      if (!field.conditionalDisplay) return false;
      
      // Check if this field's conditional display references the trigger property
      const condition = field.conditionalDisplay;
      
      // Handle various condition patterns:
      // Simple: "triggerProperty=VALUE"
      // Complex: "field1=VALUE AND field2=VALUE"
      // OR: "field1=VALUE OR field2=VALUE"
      // Includes: "triggerProperty_includes=VALUE"
      // Not equals: "triggerProperty!=VALUE"
      
      return condition.includes(triggerProperty + '=') || 
             condition.includes(triggerProperty + '!=') ||
             condition.includes(triggerProperty + '_includes=');
    });
  };

  // Get fields that are triggered by a specific field
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

  // Render fields with proper grouping for trigger fields
  const renderFieldsWithTriggerGrouping = (fieldsToRender: OntologyField[]) => {
    const renderedFields: React.ReactElement[] = [];
    const processedFields = new Set<string>();

    // Pre-process: Mark all conditional fields as processed so they don't render in their original position
    fieldsToRender.forEach((field) => {
      if (isTriggerField(field)) {
        const triggeredFields = getTriggeredFields(field.property);
        triggeredFields.forEach(f => processedFields.add(f.property));
      }
    });

    fieldsToRender.forEach((field) => {
      if (processedFields.has(field.property)) return;

      if (isTriggerField(field)) {
        // This is a trigger field - render it with its conditional fields in a bordered section
        const triggeredFields = getTriggeredFields(field.property);
        
        // Use the existing shouldDisplayField logic to determine if any triggered fields should show
        const visibleTriggeredFields = triggeredFields.filter(shouldDisplayField);
        const showTriggeredFields = visibleTriggeredFields.length > 0;

        renderedFields.push(
          <div key={field.property} className="space-y-4">
            {/* Render the trigger field */}
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              {renderField(field)}
            </div>
            
            {/* Render triggered fields in a bordered section if conditions are met */}
            {showTriggeredFields && (
              <div className="border-2 border-blue-200 rounded-lg p-4 bg-blue-50/30">
                <div className="mb-3">
                  <h5 className="text-sm font-semibold text-blue-800 uppercase tracking-wide">
                    {field.label?.replace(/\?$/, '')} Details
                  </h5>
                  <div className="w-12 h-0.5 bg-blue-400 mt-1"></div>
                </div>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {visibleTriggeredFields.map((triggeredField) => {
                    return renderField(triggeredField);
                  })}
                </div>
              </div>
            )}
          </div>
        );

        // Mark the trigger field as processed
        processedFields.add(field.property);
      } else {
        // Regular field - render normally
        renderedFields.push(
          <div key={field.property} className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {renderField(field)}
          </div>
        );
        processedFields.add(field.property);
      }
    });

    return renderedFields;
  };

  const getStatusBadge = (isRequired: boolean, hasValue: boolean) => {
    if (isRequired && !hasValue) {
      return <Badge color="failure" className="ml-2 text-xs">Required</Badge>;
    }
    if (hasValue) {
      return <Badge color="success" className="ml-2 text-xs">✓</Badge>;
    }
    return null;
  };

  const getFieldError = (fieldName: string) => {
    return validationErrors.find(error => error.field === fieldName);
  };

  // Group fields by category based on ontology formSection property first, then fallback to field analysis
  const groupFields = (fields: OntologyField[]) => {
    const groups: { [key: string]: OntologyField[] } = {};
    
    // Deduplicate fields by property name (keep the first occurrence)
    const uniqueFields = fields.filter((field, index, self) => 
      index === self.findIndex(f => f.property === field.property)
    );
    
    uniqueFields.forEach(field => {
      let category = 'General';
      
      // First, check if field has a formSection property from ontology
      if (field.formSection) {
        switch (field.formSection) {
          case 'driver_identity':
            category = 'Driver Identity';
            break;
          case 'mandatory':
            category = 'General';
            break;
          case 'licence_core':
          case 'licence_history':
          case 'licence_optional':
            category = 'Licence Information';
            break;
          case 'no_claims':
            category = 'No Claims';
            break;
          case 'medical':
          case 'medical_history':
            category = 'Medical Conditions & Restrictions';
            break;
          case 'endorsements':
          case 'endorsements_history':
            category = 'Convictions & Penalty Points';
            break;
          case 'accidents':
            category = 'Claims & Accidents';
            break;
          case 'disqualifications':
          case 'disqualifications_history':
            category = 'Disqualifications';
            break;
          case 'modifications':
            category = 'Modifications';
            break;
          default:
            // Use the formSection as category name (capitalize first letter)
            category = field.formSection.charAt(0).toUpperCase() + field.formSection.slice(1).replace(/_/g, ' ');
        }
      } else {
        // Fallback to field name analysis for fields without formSection
        if (['firstName', 'lastName', 'middleName', 'dateOfBirth', 'email', 'phone', 'address', 'postcode', 'title'].includes(field.property)) {
          // Personal details only show in Settings form, skip in Driver form
          if (formType === 'settings') {
            category = 'Personal Details';
          } else {
            return; // Skip personal details in non-settings forms
          }
        } else if (field.property.toLowerCase().includes('licence') || field.property.toLowerCase().includes('license')) {
          category = 'Licence Information';
        } else if (['hasConvictions', 'offenceCode', 'penaltyPoints', 'pointsLost', 'totalActivePoints', 'disqualificationRisk', 'convictionDate', 'convictionType', 'convictionPenalty', 'convictionPoints', 'convictionFine', 'disqualificationPeriod', 'convictionActive'].includes(field.property)) {
          category = 'Convictions & Penalty Points';
        } else if (['disabilityTypes', 'adaptationTypes', 'automaticOnly'].includes(field.property)) {
          category = 'Disabilities & Adaptations';
        } else if (['registration', 'make', 'model', 'year', 'engineSize', 'fuelType', 'transmission', 'bodyType', 'doors', 'seats'].includes(field.property)) {
          category = 'Basic Information';
        } else if (['value', 'purchasePrice', 'purchaseDate', 'ownershipType', 'financingCompany', 'estimatedValue'].includes(field.property)) {
          category = 'Value & Ownership';
        } else if (['daytimeLocation', 'overnightLocation', 'annualMileage', 'businessUse', 'businessUseType'].includes(field.property)) {
          category = 'Usage & Location';
        } else if (field.property.toLowerCase().includes('modification') || field.property === 'hasModifications') {
          category = 'Modifications';
        } else if (field.property.toLowerCase().includes('alarm') || field.property.toLowerCase().includes('immobiliser') || field.property.toLowerCase().includes('tracking') || field.property.toLowerCase().includes('security')) {
          category = 'Security Features';
        } else if (field.property.toLowerCase().includes('mot') || field.property.toLowerCase().includes('tax') || field.property.toLowerCase().includes('service') || ['vin', 'engineNumber', 'chassisNumber'].includes(field.property)) {
          category = 'Documents & Compliance';
        } else if (['hasClaims', 'hasAccidents', 'claimsCount', 'accidentsCount'].includes(field.property)) {
          category = 'Claims & Accidents';
        } else if (field.property.toLowerCase().includes('claim') || field.property.toLowerCase().includes('accident')) {
          category = 'Claim Details';
        } else if (field.property.toLowerCase().includes('cover') || field.property.toLowerCase().includes('policy') || field.property.toLowerCase().includes('premium')) {
          category = 'Policy Information';
        }
      }
      
      if (!groups[category]) {
        groups[category] = [];
      }
      groups[category].push(field);
    });
    
    return groups;
  };

  const renderField = (field: OntologyField) => {
    // Get field value, applying default if empty and field has a default value
    let fieldValue = (data as any)[field.property];
    if (!fieldValue && field.defaultValue) {
      fieldValue = field.defaultValue;
      // Update the data with the default value
      onUpdate(field.property, fieldValue);
    }
    
    const fieldError = getFieldError(field.property);
    
    return (
      <div key={field.property} className="space-y-2">
        <Label htmlFor={field.property} className="flex items-center">
          {field.label.split(' ').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')}
          {getStatusBadge(field.required, !!fieldValue)}
        </Label>
        
        {field.formType === 'radio' || field.formType === 'slider' ? (
          <div className="flex items-center space-x-4">
            {(field.enumerationValues || ['NO', 'YES']).map((option, index) => {
              const isChecked = fieldValue === option;
              
              // Get display label from options or use the option value
              const displayLabel = field.options?.find(opt => opt.value === option)?.label || 
                                 (option === 'YES' ? 'Yes' : 
                                  option === 'NO' ? 'No' : 
                                  option === 'MANUAL' ? 'Manual' :
                                  option === 'AUTOMATIC' ? 'Automatic' :
                                  option.replace(/_/g, ' ').toLowerCase().replace(/\b\w/g, l => l.toUpperCase()));
              
              return (
                <label key={option} className="flex items-center">
                  <input
                    type="radio"
                    name={field.property}
                    value={option}
                    checked={isChecked}
                    onChange={() => onUpdate(field.property, option)}
                    className="mr-2 text-blue-600 focus:ring-blue-500"
                  />
                  <span className="text-sm font-medium text-gray-700">{displayLabel}</span>
                </label>
              );
            })}
          </div>
        ) : field.formType === 'endorsement_array' ? (
          // Ontology-driven endorsement array with individual points and dates for each selected code
          <div className="space-y-4">
            {(field.enumerationValues || []).map((value) => {
              const currentValues = Array.isArray(fieldValue) ? fieldValue : (fieldValue ? [fieldValue] : []);
              const isChecked = currentValues.includes(value);
              const endorsementKey = value.split('_')[0]; // Get the code part (e.g., SP30)
              
              return (
                <div key={value} className="border border-gray-200 rounded-lg p-3">
                  <label className="flex items-center mb-2">
                    <input
                      type="checkbox"
                      checked={isChecked}
                      onChange={(e) => {
                        let newValues;
                        if (e.target.checked) {
                          newValues = [...currentValues, value];
                        } else {
                          newValues = currentValues.filter(v => v !== value);
                          // Also clear the related points and date fields
                          onUpdate(`${endorsementKey}_points`, '');
                          onUpdate(`${endorsementKey}_date`, '');
                        }
                        onUpdate(field.property, newValues);
                      }}
                      className="mr-2 text-blue-600 focus:ring-blue-500"
                    />
                    <span className="text-sm font-medium text-gray-700">{value.replace(/_/g, ' ')}</span>
                  </label>
                  
                  {isChecked && (
                    <div className="ml-6 grid grid-cols-1 md:grid-cols-2 gap-3 mt-2">
                      <div>
                        <label className="block text-xs font-medium text-gray-600 mb-1">
                          Penalty Points
                        </label>
                        <input
                          type="number"
                          min="0"
                          max="12"
                          placeholder="Points (0-12)"
                          value={(data as any)[`${endorsementKey}_points`] || ''}
                          onChange={(e) => onUpdate(`${endorsementKey}_points`, e.target.value)}
                          className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                        />
                      </div>
                      <div>
                        <label className="block text-xs font-medium text-gray-600 mb-1">
                          Offence Date
                        </label>
                        <input
                          type="date"
                          value={(data as any)[`${endorsementKey}_date`] || ''}
                          onChange={(e) => onUpdate(`${endorsementKey}_date`, e.target.value)}
                          className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                        />
                      </div>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        ) : field.isMultiSelect ? (
          <div className="space-y-2">
            {field.property === 'endorsementCode' ? (
              // Special handling for endorsement codes with individual points and dates
              <div className="space-y-4">
                {(field.enumerationValues || []).map((value) => {
                  const currentValues = Array.isArray(fieldValue) ? fieldValue : (fieldValue ? [fieldValue] : []);
                  const isChecked = currentValues.includes(value);
                  const endorsementKey = value.split('_')[0]; // Get the code part (e.g., SP30)
                  
                  return (
                    <div key={value} className="border border-gray-200 rounded-lg p-3">
                      <label className="flex items-center mb-2">
                        <input
                          type="checkbox"
                          checked={isChecked}
                          onChange={(e) => {
                            let newValues;
                            if (e.target.checked) {
                              newValues = [...currentValues, value];
                            } else {
                              newValues = currentValues.filter(v => v !== value);
                              // Also clear the related points and date fields
                              onUpdate(`${endorsementKey}_points`, '');
                              onUpdate(`${endorsementKey}_date`, '');
                            }
                            onUpdate(field.property, newValues);
                          }}
                          className="mr-2 text-blue-600 focus:ring-blue-500"
                        />
                        <span className="text-sm font-medium text-gray-700">{value.replace(/_/g, ' ')}</span>
                      </label>
                      
                      {isChecked && (
                        <div className="ml-6 grid grid-cols-1 md:grid-cols-2 gap-3 mt-2">
                          <div>
                            <label className="block text-xs font-medium text-gray-600 mb-1">
                              Penalty Points
                            </label>
                            <input
                              type="number"
                              min="0"
                              max="12"
                              placeholder="Points (0-12)"
                              value={(data as any)[`${endorsementKey}_points`] || ''}
                              onChange={(e) => onUpdate(`${endorsementKey}_points`, e.target.value)}
                              className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                            />
                          </div>
                          <div>
                            <label className="block text-xs font-medium text-gray-600 mb-1">
                              Offence Date
                            </label>
                            <input
                              type="date"
                              value={(data as any)[`${endorsementKey}_date`] || ''}
                              onChange={(e) => onUpdate(`${endorsementKey}_date`, e.target.value)}
                              className="w-full px-2 py-1 text-sm border border-gray-300 rounded focus:ring-1 focus:ring-blue-500 focus:border-blue-500"
                            />
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            ) : (
              // Standard multi-select for other fields
              <div className="space-y-2">
                {(field.enumerationValues || []).map((value) => {
                  const currentValues = Array.isArray(fieldValue) ? fieldValue : (fieldValue ? [fieldValue] : []);
                  const isChecked = currentValues.includes(value);
                  
                  return (
                    <label key={value} className="flex items-center">
                      <input
                        type="checkbox"
                        checked={isChecked}
                        onChange={(e) => {
                          let newValues;
                          if (e.target.checked) {
                            newValues = [...currentValues, value];
                          } else {
                            newValues = currentValues.filter(v => v !== value);
                          }
                          onUpdate(field.property, newValues);

                          // Auto-open AI dialog for dependent fields when conditions become true (e.g., selecting "Other")
                          const newData = { ...(data as any), [field.property]: newValues };
                          if (isTriggerField(field)) {
                            const dependentFields = getTriggeredFields(field.property).filter((dep) => {
                              if (!dep.conditionalDisplay) return false;
                              const condition = dep.conditionalDisplay;
                              if (condition.includes(' AND ')) {
                                const andConditions = condition.split(' AND ');
                                return andConditions.every((c) => evaluateCondition(c.trim(), newData));
                              }
                              if (condition.includes(' OR ')) {
                                const orConditions = condition.split(' OR ');
                                return orConditions.some((c) => evaluateCondition(c.trim(), newData));
                              }
                              return evaluateCondition(condition, newData);
                            });

                            const aiDependent = dependentFields.find((dep) => (dep as any).requiresAIValidation);
                            if (aiDependent) {
                              openAIValidationDialog(aiDependent, '');
                            }
                          }
                        }}
                        className="mr-2 text-blue-600 focus:ring-blue-500"
                      />
                      <span className="text-sm text-gray-700">{value.replace(/_/g, ' ')}</span>
                    </label>
                  );
                })}
              </div>
            )}
          </div>
        ) : field.options && field.options.length > 0 ? (
          <Select
            id={field.property}
            value={fieldValue}
            onChange={(e) => {
              const newValue = e.target.value;
              const newData = { ...(data as any), [field.property]: newValue };
              onUpdate(field.property, newValue);
              // If selecting this option reveals dependent fields, auto-open AI dialog
              if (isTriggerField(field)) {
                const dependentFields = getTriggeredFields(field.property).filter((dep) => {
                  if (!dep.conditionalDisplay) return false;
                  const condition = dep.conditionalDisplay;
                  if (condition.includes(' AND ')) {
                    const andConditions = condition.split(' AND ');
                    return andConditions.every((c) => evaluateCondition(c.trim(), newData));
                  }
                  if (condition.includes(' OR ')) {
                    const orConditions = condition.split(' OR ');
                    return orConditions.some((c) => evaluateCondition(c.trim(), newData));
                  }
                  return evaluateCondition(condition, newData);
                });

                const aiDependent = dependentFields.find((dep) => (dep as any).requiresAIValidation);
                if (aiDependent) {
                  openAIValidationDialog(aiDependent, '');
                }
              }
              // Generic fallback: if any newly-visible AI field has no value, open it
              triggerAIForNewData(newData);
            }}
            color={fieldError ? 'failure' : 'gray'}
          >
            <option value="">Select {field.label}</option>
            {field.options.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </Select>
        ) : field.type === 'radio' && field.enumerationValues ? (
          <div className="flex gap-4">
            {field.enumerationValues.map((value) => (
              <label key={value} className="flex items-center">
                <input
                  type="radio"
                  name={field.property}
                  value={value}
                  checked={fieldValue === value}
                  onChange={(e) => onUpdate(field.property, e.target.value)}
                  className="mr-2"
                />
                {value}
              </label>
            ))}
          </div>
        ) : field.property === 'hasConvictions' || field.property === 'hasModifications' || field.property === 'hasClaims' || field.property === 'hasAccidents' ? (
          <div className="flex gap-4">
            <label className="flex items-center">
              <input
                type="radio"
                name={field.property}
                value="true"
                checked={fieldValue === 'true' || fieldValue === true}
                onChange={() => onUpdate(field.property, true)}
                className="mr-2"
              />
              Yes
            </label>
            <label className="flex items-center">
              <input
                type="radio"
                name={field.property}
                value="false"
                checked={fieldValue === 'false' || fieldValue === false}
                onChange={() => onUpdate(field.property, false)}
                className="mr-2"
              />
              No
            </label>
          </div>
        ) : field.property === 'modifications' ? (
          <textarea
            id={field.property}
            value={Array.isArray(fieldValue) ? fieldValue.join('\n') : fieldValue}
            onChange={(e) => onUpdate(field.property, e.target.value.split('\n').filter(line => line.trim()))}
            placeholder="List all modifications (one per line)&#10;e.g.&#10;Alloy wheels&#10;Lowered suspension&#10;Performance exhaust&#10;Turbo/Supercharger&#10;Body kit"
            rows={6}
            className="w-full p-2 border border-gray-300 rounded-md focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        ) : field.type === 'file' ? (
          <div className="border-2 border-dashed border-gray-300 rounded-lg p-4 text-center">
            <input
              type="file"
              id={field.property}
              accept="image/*,.pdf"
              className="hidden"
            />
            <label htmlFor={field.property} className="cursor-pointer">
              <div className="text-gray-600">
                <Plus className="w-8 h-8 mx-auto mb-2" />
                <p>Upload {field.label}</p>
              </div>
            </label>
          </div>
        ) : field.formType === 'textarea' ? (
          <div className="space-y-2">
            <div className="relative">
              <textarea
                id={field.property}
                value={fieldValue || ''}
                onChange={(e) => {
                  if ((field as any).requiresAIValidation) return; // block manual edits
                  onUpdate(field.property, e.target.value);
                }}
                onKeyDown={(e) => {
                  if ((field as any).requiresAIValidation && e.key === 'Enter' && !e.shiftKey) {
                    e.preventDefault();
                    openAIValidationDialog(field, (fieldValue || '') as string);
                  }
                }}
                placeholder={field.helpText || `Enter ${field.label.toLowerCase()}`}
                rows={4}
                className={`w-full p-3 border rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 ${
                  fieldError ? 'border-red-500' : 'border-gray-300'
                }`}
                style={(field as any).requiresAIValidation && !fieldValue ? { display: 'none' } : undefined}
                readOnly={(field as any).requiresAIValidation}
                onFocus={(e) => {
                  if ((field as any).requiresAIValidation) {
                    e.currentTarget.blur();
                    openAIValidationDialog(field, (fieldValue || '') as string);
                  }
                }}
              />
              {(field as any).requiresAIValidation && fieldValue && (
                <Button
                  size="xs"
                  color="blue"
                  className="absolute top-2 right-2"
                  onClick={() => openAIValidationDialog(field, (fieldValue || '') as string)}
                >
                  Edit via AI
                </Button>
              )}
            </div>
            {(field as any).requiresAIValidation && (
              <p className="text-xs text-blue-600">
                ℹ️ This field requires specific insurance-relevant information. Press Enter to open validation. Use Shift+Enter for a new line.
              </p>
            )}
          </div>
        ) : (
          <TextInput
            id={field.property}
            type={
              field.type === 'date' ? 'date' : 
              field.type === 'email' ? 'email' : 
              field.type === 'tel' ? 'tel' : 
              field.type === 'number' || 
              field.property.includes('year') || 
              field.property.includes('Year') ||
              field.property.includes('value') || 
              field.property.includes('Value') ||
              field.property.includes('price') || 
              field.property.includes('Price') ||
              field.property.includes('doors') || 
              field.property.includes('seats') ||
              field.property.includes('mileage') ||
              field.property.includes('points') ||
              field.property.includes('Points') ? 'number' : 
              'text'
            }
            value={fieldValue}
            onChange={(e) => onUpdate(field.property, e.target.value)}
            placeholder={field.helpText ? field.helpText.substring(0, 50) + '...' : `Enter ${field.label.toLowerCase()}`}
            color={fieldError ? 'failure' : 'gray'}
            min={field.minInclusive !== undefined ? field.minInclusive : undefined}
            max={field.maxInclusive !== undefined ? field.maxInclusive : undefined}
          />
        )}
        
        {fieldError && (
          <p className="text-sm text-red-600 mt-1">{fieldError.message}</p>
        )}
        
        {field.helpText && !fieldError && !(field as any).requiresAIValidation && (
          <p className="text-xs text-gray-500 mt-1">{field.helpText}</p>
        )}
      </div>
    );
  };

  const fieldGroups = groupFields(fields);

  // Define section order to ensure Driver Identity appears first
  const sectionOrder = [
    'Driver Identity',
    'General', 
    'Licence Information',
    'No Claims',
    'Medical Conditions & Restrictions',
    'Convictions & Penalty Points',
    'Claims & Accidents',
    'Disqualifications',
    'Claim Details',
    'Modifications',
    'Disabilities & Adaptations'
  ];

  // Sort sections according to defined order
  const sortedSections: [string, OntologyField[]][] = sectionOrder
    .filter(sectionName => fieldGroups[sectionName])
    .map(sectionName => [sectionName, fieldGroups[sectionName]] as [string, OntologyField[]]);

  return (
    <Card className="mb-6">
      <div className="flex justify-between items-center mb-4">
        <div className="flex items-center">
          {icon}
          <h3 className="text-lg font-semibold">{title}</h3>
        </div>
        {showRemoveButton && onRemove && (
          <Button color="failure" size="sm" onClick={onRemove}>
            <X className="w-4 h-4" />
          </Button>
        )}
      </div>

      {sortedSections.map(([groupName, groupFields]) => (
        <div key={groupName} className="mb-6">
          {/* Collapsible header for optional sections */}
          {(groupName === 'Convictions & Penalty Points' || groupName === 'Disabilities & Adaptations' || groupName === 'Medical Conditions & Restrictions' || groupName === 'Claims & Accidents') ? (
            <div className="border border-gray-200 rounded-lg">
              <button
                type="button"
                onClick={() => toggleSection(groupName)}
                className="w-full flex items-center justify-between p-4 text-left hover:bg-gray-50 rounded-lg"
              >
                <div className="flex items-center">
                  <h4 className="text-md font-medium text-gray-700">
                    {groupName}
                  </h4>
                  <Badge color="gray" size="sm" className="ml-2">Optional</Badge>
                </div>
                <div className="flex items-center">
                  {((groupName === 'Convictions & Penalty Points' && (data as any)?.hasConvictions === 'YES') ||
                    (groupName === 'Disabilities & Adaptations' && ((data as any)?.disabilityTypes || (data as any)?.adaptationTypes)) ||
                    (groupName === 'Medical Conditions & Restrictions' && (data as any)?.hasMedicalConditions === 'YES') ||
                    (groupName === 'Claims & Accidents' && (data as any)?.hasAccidents === 'YES')) && (
                    <Badge color="warning" size="sm" className="mr-2">Active</Badge>
                  )}
                  <span className="text-gray-400">
                    {collapsedSections[groupName] ? '▼' : '▲'}
                  </span>
                </div>
              </button>
              
              {!collapsedSections[groupName] && (
                <div className="p-4 border-t border-gray-200">
                  {/* Warning message */}
                  <div className="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4">
                    <p className="text-sm text-yellow-800">
                      {groupName === 'Convictions & Penalty Points' ? (
                        <>
                          <strong>Important:</strong> You must declare all convictions and penalty points. 
                          Failure to disclose may invalidate your insurance.
                        </>
                      ) : groupName === 'Medical Conditions & Restrictions' ? (
                        <>
                          <strong>Important:</strong> You must declare all medical conditions that affect your driving. 
                          This includes conditions declared to DVLA and any vision/hearing requirements.
                        </>
                      ) : groupName === 'Claims & Accidents' ? (
                        <>
                          <strong>Important:</strong> You must declare all accidents and claims in the last 5 years. 
                          This includes both fault and non-fault accidents, regardless of whether you claimed.
                        </>
                      ) : (
                        <>
                          <strong>Important:</strong> Please declare any disabilities that may affect your driving. 
                          This helps us provide appropriate coverage and may qualify you for discounts.
                        </>
                      )}
                    </p>
                  </div>
                  
                  {/* Render section fields */}
                  <div className="space-y-4">
                    {renderFieldsWithTriggerGrouping(groupFields.filter(shouldDisplayField))}
                  </div>
                </div>
              )}
            </div>
          ) : (
            <>
              <h4 className="text-md font-medium mb-3 text-gray-700">{groupName}</h4>
              
              {/* Special warnings for other categories */}
          
          {groupName === 'Modifications' && (
            <div className="bg-orange-50 border border-orange-200 rounded-lg p-4 mb-4">
              <p className="text-sm text-orange-800">
                <strong>Important:</strong> You must declare ALL modifications to your vehicle. 
                Undeclared modifications may invalidate your insurance policy.
              </p>
            </div>
          )}
          
          <div className="space-y-4">
            {renderFieldsWithTriggerGrouping(groupFields.filter(shouldDisplayField))}
          </div>
          
          {/* Show modification examples if hasModifications is true */}
          {groupName === 'Modifications' && (data as any).hasModifications === true && (
            <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
              <h5 className="font-medium text-blue-800 mb-2">Common Vehicle Modifications:</h5>
              <div className="text-sm text-blue-700 grid grid-cols-2 gap-2">
                <div>
                  • Alloy wheels
                  • Lowered suspension
                  • Performance exhaust
                  • Turbo/Supercharger
                  • Cold air intake
                </div>
                <div>
                  • Body kit/spoilers
                  • Window tinting
                  • Roll cage
                  • Performance brakes
                  • Engine remapping
                </div>
              </div>
            </div>
          )}
            </>
          )}
        </div>
      ))}
      
      <AIValidationDialog
        isOpen={aiValidationDialog.isOpen}
        onClose={closeAIValidationDialog}
        fieldName={aiValidationDialog.fieldName}
        fieldLabel={aiValidationDialog.fieldLabel}
        initialValue={aiValidationDialog.initialValue}
        validationPrompt={aiValidationDialog.validationPrompt}
        onValidatedValue={handleValidatedValue}
      />
    </Card>
  );
};

export default UniversalForm;
