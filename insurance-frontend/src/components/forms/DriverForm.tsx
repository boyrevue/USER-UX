import React from 'react';
import { Driver } from '../../types';


interface DriverFormProps {
  driver: Driver;
  index: number;
  updateDriver: (index: number, field: string, value: string) => void;
  removeDriver: (index: number) => void;
  // Add other props as needed
}

const DriverForm: React.FC<DriverFormProps> = ({ 
  driver, 
  index, 
  updateDriver, 
  removeDriver 
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
