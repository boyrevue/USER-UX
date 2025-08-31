import React from 'react';
import { Card, Button, TextInput, Label, Select, Badge } from 'flowbite-react';
import { X, User } from 'lucide-react';
import { DriverFormProps, ValidationError, OntologyField } from '../../types';

interface DynamicDriverFormProps extends DriverFormProps {
  fields?: OntologyField[];
}

const DriverForm: React.FC<DynamicDriverFormProps> = ({ 
  driver, 
  index, 
  updateDriver, 
  removeDriver,
  validationErrors = [],
  fields = []
}) => {
  const getStatusBadge = (isRequired: boolean, hasValue: boolean) => {
    if (isRequired && !hasValue) {
      return <Badge color="failure" className="ml-2 text-xs">Required</Badge>;
    }
    if (hasValue) {
      return <Badge color="success" className="ml-2 text-xs">Complete</Badge>;
    }
    return null;
  };

  const getFieldError = (fieldName: string): ValidationError | undefined => {
    return validationErrors.find(error => 
      error.field === `driver_${index}_${fieldName}` || error.field === fieldName
    );
  };

  const handleDriverUpdate = (field: keyof typeof driver, value: any) => {
    updateDriver(index, field, value);
  };

  return (
    <Card className="mb-6">
      <div className="flex justify-between items-center mb-4">
        <div className="flex items-center">
          <User className="w-5 h-5 mr-2 text-blue-600" />
          <h3 className="text-lg font-semibold">
            Driver {index + 1} {driver.classification === 'MAIN' ? '(Main Driver)' : '(Named Driver)'}
          </h3>
        </div>
        {index > 0 && (
          <Button
            color="failure"
            size="sm"
            onClick={() => removeDriver(index)}
          >
            <X className="w-4 h-4" />
          </Button>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {/* First Name */}
        <div>
          <Label htmlFor={`firstName-${index}`} className="flex items-center">
            First Name
            {getStatusBadge(true, !!driver.firstName)}
          </Label>
          <TextInput
            id={`firstName-${index}`}
            value={driver.firstName || ''}
            onChange={(e) => handleDriverUpdate('firstName', e.target.value)}
            placeholder="Enter first name"
            color={getFieldError('firstName') ? 'failure' : 'gray'}
          />
          {getFieldError('firstName') && (
            <p className="text-sm text-red-600 mt-1">{getFieldError('firstName')?.message}</p>
          )}
        </div>

        {/* Last Name */}
        <div>
          <Label htmlFor={`lastName-${index}`} className="flex items-center">
            Last Name
            {getStatusBadge(true, !!driver.lastName)}
          </Label>
          <TextInput
            id={`lastName-${index}`}
            value={driver.lastName || ''}
            onChange={(e) => handleDriverUpdate('lastName', e.target.value)}
            placeholder="Enter last name"
            color={getFieldError('lastName') ? 'failure' : 'gray'}

          />
        </div>

        {/* Date of Birth */}
        <div>
          <Label htmlFor={`dateOfBirth-${index}`} className="flex items-center">
            Date of Birth
            {getStatusBadge(true, !!driver.dateOfBirth)}
          </Label>
          <TextInput
            id={`dateOfBirth-${index}`}
            type="date"
            value={driver.dateOfBirth || ''}
            onChange={(e) => handleDriverUpdate('dateOfBirth', e.target.value)}
            color={getFieldError('dateOfBirth') ? 'failure' : 'gray'}

          />
        </div>

        {/* Email */}
        <div>
          <Label htmlFor={`email-${index}`} className="flex items-center">
            Email
            {getStatusBadge(true, !!driver.email)}
          </Label>
          <TextInput
            id={`email-${index}`}
            type="email"
            value={driver.email || ''}
            onChange={(e) => handleDriverUpdate('email', e.target.value)}
            placeholder="Enter email address"
            color={getFieldError('email') ? 'failure' : 'gray'}

          />
        </div>

        {/* Phone */}
        <div>
          <Label htmlFor={`phone-${index}`} className="flex items-center">
            Phone Number
            {getStatusBadge(true, !!driver.phone)}
          </Label>
          <TextInput
            id={`phone-${index}`}
            type="tel"
            value={driver.phone || ''}
            onChange={(e) => handleDriverUpdate('phone', e.target.value)}
            placeholder="Enter phone number"
            color={getFieldError('phone') ? 'failure' : 'gray'}

          />
        </div>

        {/* Licence Number */}
        <div>
          <Label htmlFor={`licenceNumber-${index}`} className="flex items-center">
            Licence Number
            {getStatusBadge(true, !!driver.licenceNumber)}
          </Label>
          <TextInput
            id={`licenceNumber-${index}`}
            value={driver.licenceNumber || ''}
            onChange={(e) => handleDriverUpdate('licenceNumber', e.target.value)}
            placeholder="Enter licence number"
            color={getFieldError('licenceNumber') ? 'failure' : 'gray'}

          />
        </div>

        {/* Licence Issue Date */}
        <div>
          <Label htmlFor={`licenceIssueDate-${index}`} className="flex items-center">
            Licence Issue Date
            {getStatusBadge(false, !!driver.licenceIssueDate)}
          </Label>
          <TextInput
            id={`licenceIssueDate-${index}`}
            type="date"
            value={driver.licenceIssueDate || ''}
            onChange={(e) => handleDriverUpdate('licenceIssueDate', e.target.value)}
            color={getFieldError('licenceIssueDate') ? 'failure' : 'gray'}

          />
        </div>

        {/* Licence Expiry Date */}
        <div>
          <Label htmlFor={`licenceExpiryDate-${index}`} className="flex items-center">
            Licence Expiry Date
            {getStatusBadge(false, !!driver.licenceExpiryDate)}
          </Label>
          <TextInput
            id={`licenceExpiryDate-${index}`}
            type="date"
            value={driver.licenceExpiryDate || ''}
            onChange={(e) => handleDriverUpdate('licenceExpiryDate', e.target.value)}
            color={getFieldError('licenceExpiryDate') ? 'failure' : 'gray'}

          />
        </div>
      </div>

      {/* Classification Badge */}
      <div className="mt-4 pt-4 border-t border-gray-200">
        <div className="flex items-center justify-between">
          <div className="flex items-center space-x-2">
            <span className="text-sm text-gray-600">Driver Type:</span>
            <Badge color={driver.classification === 'MAIN' ? 'success' : 'info'}>
              {driver.classification === 'MAIN' ? 'Main Driver' : 'Named Driver'}
            </Badge>
          </div>
        </div>
      </div>
    </Card>
  );
};

export default DriverForm;
