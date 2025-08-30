import React from 'react';
import { Vehicle } from '../../types';


interface VehicleFormProps {
  vehicle: Vehicle;
  index: number;
  updateVehicle: (index: number, field: string, value: string) => void;
  removeVehicle: (index: number) => void;
}

const VehicleForm: React.FC<VehicleFormProps> = ({ 
  vehicle, 
  index, 
  updateVehicle, 
  removeVehicle 
}) => {
  return (
    <div className="vehicle-form">
      {/* Vehicle form content will be moved here */}
      <h3>Vehicle {index + 1}</h3>
      {/* Form fields */}
    </div>
  );
};


export default VehicleForm;
