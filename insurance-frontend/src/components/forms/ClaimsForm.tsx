import React from 'react';
import { Claim } from '../../types';


interface ClaimsFormProps {
  claims: Claim[];
  updateClaim: (index: number, field: string, value: string) => void;
  addClaim: () => void;
  removeClaim: (index: number) => void;
}

const ClaimsForm: React.FC<ClaimsFormProps> = ({ 
  claims, 
  updateClaim, 
  addClaim, 
  removeClaim 
}) => {
  return (
    <div className="claims-form">
      {/* Claims form content will be moved here */}
      <h3>Claims History</h3>
      {/* Form fields */}
    </div>
  );
};


export default ClaimsForm;
