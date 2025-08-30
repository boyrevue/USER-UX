import React from 'react';
import { DriverFormProps } from '../../types';

const DriverForm: React.FC<DriverFormProps> = ({ 
  driver, 
  index, 
  updateDriver, 
  removeDriver,
  validationErrors = []
}) => {
  return (
    <div className="driver-form">
      {/* Driver form content will be moved here */}
      <h3>Driver {index + 1}</h3>
      {/* Form fields */}
    </div>
  );
};


export default DriverForm;
