import React from 'react';
import { Card, Button } from 'flowbite-react';
import { ClaimsFormProps } from '../../types';

const ClaimsForm: React.FC<ClaimsFormProps> = ({ 
  claims,
  accidents,
  updateClaim,
  updateAccident,
  addClaim,
  addAccident,
  removeClaim,
  removeAccident
}) => {
  return (
    <Card>
      <h3 className="text-lg font-semibold mb-4">Claims & Accidents History</h3>
      <div className="space-y-4">
        <div>
          <h4 className="font-medium mb-2">Claims ({claims.length})</h4>
          <Button onClick={addClaim} size="sm">Add Claim</Button>
        </div>
        <div>
          <h4 className="font-medium mb-2">Accidents ({accidents.length})</h4>
          <Button onClick={addAccident} size="sm">Add Accident</Button>
        </div>
      </div>
    </Card>
  );
};

export default ClaimsForm;
