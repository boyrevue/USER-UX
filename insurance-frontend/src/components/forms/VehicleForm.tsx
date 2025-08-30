import React from 'react';
import { Card, Button, TextInput, Label, Select, Badge } from 'flowbite-react';
import { X, Car } from 'lucide-react';
import { VehicleFormProps, ValidationError } from '../../types';

const VehicleForm: React.FC<VehicleFormProps> = ({ 
  vehicle, 
  index, 
  updateVehicle, 
  removeVehicle,
  validationErrors = []
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
      error.field === `vehicle_${index}_${fieldName}` || error.field === fieldName
    );
  };

  const handleVehicleUpdate = (field: keyof typeof vehicle, value: any) => {
    updateVehicle(index, field, value);
  };

  return (
    <Card className="mb-6">
      <div className="flex justify-between items-center mb-4">
        <div className="flex items-center">
          <Car className="w-5 h-5 mr-2 text-blue-600" />
          <h3 className="text-lg font-semibold">Vehicle {index + 1}</h3>
        </div>
        {index > 0 && (
          <Button
            color="failure"
            size="sm"
            onClick={() => removeVehicle(index)}
          >
            <X className="w-4 h-4" />
          </Button>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div className="form-group">
          <Label className="form-label">
            Registration Number
            {getStatusBadge(true, !!vehicle.registration)}
          </Label>
          <TextInput
            type="text"
            value={vehicle.registration}
            onChange={(e) => handleVehicleUpdate('registration', e.target.value.toUpperCase())}
            className="form-input"
            placeholder="e.g. AB12 CDE"
            color={getFieldError('registration') ? 'failure' : undefined}
          />
          {getFieldError('registration') && (
            <p className="text-red-500 text-sm mt-1">{getFieldError('registration')?.message}</p>
          )}
        </div>

        <div className="form-group">
          <Label className="form-label">
            Make
            {getStatusBadge(true, !!vehicle.make)}
          </Label>
          <Select
            value={vehicle.make}
            onChange={(e) => handleVehicleUpdate('make', e.target.value)}
            className="form-input"
            color={getFieldError('make') ? 'failure' : undefined}
          >
            <option value="">Select make</option>
            <option value="Audi">Audi</option>
            <option value="BMW">BMW</option>
            <option value="Ford">Ford</option>
            <option value="Honda">Honda</option>
            <option value="Mercedes-Benz">Mercedes-Benz</option>
            <option value="Nissan">Nissan</option>
            <option value="Toyota">Toyota</option>
            <option value="Vauxhall">Vauxhall</option>
            <option value="Volkswagen">Volkswagen</option>
            <option value="Other">Other</option>
          </Select>
          {getFieldError('make') && (
            <p className="text-red-500 text-sm mt-1">{getFieldError('make')?.message}</p>
          )}
        </div>

        <div className="form-group">
          <Label className="form-label">
            Model
            {getStatusBadge(true, !!vehicle.model)}
          </Label>
          <TextInput
            type="text"
            value={vehicle.model}
            onChange={(e) => handleVehicleUpdate('model', e.target.value)}
            className="form-input"
            placeholder="Enter model"
            color={getFieldError('model') ? 'failure' : undefined}
          />
          {getFieldError('model') && (
            <p className="text-red-500 text-sm mt-1">{getFieldError('model')?.message}</p>
          )}
        </div>

        <div className="form-group">
          <Label className="form-label">
            Year
            {getStatusBadge(true, !!vehicle.year)}
          </Label>
          <TextInput
            type="number"
            min="1900"
            max={new Date().getFullYear() + 1}
            value={vehicle.year}
            onChange={(e) => handleVehicleUpdate('year', parseInt(e.target.value) || new Date().getFullYear())}
            className="form-input"
            placeholder="Enter year"
            color={getFieldError('year') ? 'failure' : undefined}
          />
          {getFieldError('year') && (
            <p className="text-red-500 text-sm mt-1">{getFieldError('year')?.message}</p>
          )}
        </div>

        <div className="form-group">
          <Label className="form-label">
            Engine Size
            {getStatusBadge(true, !!vehicle.engineSize)}
          </Label>
          <Select
            value={vehicle.engineSize}
            onChange={(e) => handleVehicleUpdate('engineSize', e.target.value)}
            className="form-input"
            color={getFieldError('engineSize') ? 'failure' : undefined}
          >
            <option value="">Select engine size</option>
            <option value="1.0L">1.0L</option>
            <option value="1.2L">1.2L</option>
            <option value="1.4L">1.4L</option>
            <option value="1.6L">1.6L</option>
            <option value="1.8L">1.8L</option>
            <option value="2.0L">2.0L</option>
            <option value="2.5L">2.5L</option>
            <option value="3.0L">3.0L</option>
            <option value="3.5L+">3.5L+</option>
          </Select>
          {getFieldError('engineSize') && (
            <p className="text-red-500 text-sm mt-1">{getFieldError('engineSize')?.message}</p>
          )}
        </div>

        <div className="form-group">
          <Label className="form-label">
            Fuel Type
            {getStatusBadge(true, !!vehicle.fuelType)}
          </Label>
          <Select
            value={vehicle.fuelType}
            onChange={(e) => handleVehicleUpdate('fuelType', e.target.value)}
            className="form-input"
            color={getFieldError('fuelType') ? 'failure' : undefined}
          >
            <option value="">Select fuel type</option>
            <option value="Petrol">Petrol</option>
            <option value="Diesel">Diesel</option>
            <option value="Hybrid">Hybrid</option>
            <option value="Electric">Electric</option>
            <option value="LPG">LPG</option>
          </Select>
          {getFieldError('fuelType') && (
            <p className="text-red-500 text-sm mt-1">{getFieldError('fuelType')?.message}</p>
          )}
        </div>

        <div className="form-group">
          <Label className="form-label">
            Transmission
            {getStatusBadge(true, !!vehicle.transmission)}
          </Label>
          <Select
            value={vehicle.transmission}
            onChange={(e) => handleVehicleUpdate('transmission', e.target.value)}
            className="form-input"
            color={getFieldError('transmission') ? 'failure' : undefined}
          >
            <option value="">Select transmission</option>
            <option value="Manual">Manual</option>
            <option value="Automatic">Automatic</option>
            <option value="Semi-Automatic">Semi-Automatic</option>
          </Select>
          {getFieldError('transmission') && (
            <p className="text-red-500 text-sm mt-1">{getFieldError('transmission')?.message}</p>
          )}
        </div>

        <div className="form-group">
          <Label className="form-label">
            Estimated Value (£)
            {getStatusBadge(true, vehicle.estimatedValue > 0)}
          </Label>
          <TextInput
            type="number"
            min="0"
            step="100"
            value={vehicle.estimatedValue}
            onChange={(e) => handleVehicleUpdate('estimatedValue', parseFloat(e.target.value) || 0)}
            className="form-input"
            placeholder="Enter estimated value"
            color={getFieldError('estimatedValue') ? 'failure' : undefined}
          />
          {getFieldError('estimatedValue') && (
            <p className="text-red-500 text-sm mt-1">{getFieldError('estimatedValue')?.message}</p>
          )}
        </div>
      </div>

      {/* Modifications Section */}
      <div className="mt-6">
        <h4 className="text-md font-semibold mb-4">Modifications</h4>
        <div className="space-y-2">
          {vehicle.modifications.map((modification, modIndex) => (
            <div key={modIndex} className="flex items-center space-x-2">
              <TextInput
                type="text"
                value={modification}
                onChange={(e) => {
                  const updatedMods = vehicle.modifications.map((mod, idx) => 
                    idx === modIndex ? e.target.value : mod
                  );
                  handleVehicleUpdate('modifications', updatedMods);
                }}
                className="flex-1"
                placeholder="Enter modification"
              />
              <Button
                color="failure"
                size="sm"
                onClick={() => {
                  const updatedMods = vehicle.modifications.filter((_, idx) => idx !== modIndex);
                  handleVehicleUpdate('modifications', updatedMods);
                }}
              >
                <X className="w-4 h-4" />
              </Button>
            </div>
          ))}
          <Button
            color="light"
            size="sm"
            onClick={() => {
              const updatedMods = [...vehicle.modifications, ''];
              handleVehicleUpdate('modifications', updatedMods);
            }}
          >
            Add Modification
          </Button>
        </div>
      </div>
    </Card>
  );
};

export default VehicleForm;