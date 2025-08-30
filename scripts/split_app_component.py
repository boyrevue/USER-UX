#!/usr/bin/env python3
"""
Split the massive App.tsx into manageable components for better AI handling.
This script extracts logical sections into separate components.
"""

import re
import os
from pathlib import Path

def create_component_file(name, content, imports=""):
    """Create a new component file with proper structure"""
    template = f"""import React from 'react';
{imports}

{content}

export default {name};
"""
    return template

def extract_driver_form_component():
    """Extract driver form into separate component"""
    return """
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
"""

def extract_vehicle_form_component():
    """Extract vehicle form into separate component"""
    return """
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
"""

def extract_claims_form_component():
    """Extract claims form into separate component"""
    return """
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
"""

def extract_document_upload_component():
    """Extract document upload into separate component"""
    return """
interface DocumentUploadProps {
  onFileUpload: (files: FileList, type: string) => void;
  isProcessing: boolean;
  extractedData?: any;
}

const DocumentUpload: React.FC<DocumentUploadProps> = ({ 
  onFileUpload, 
  isProcessing, 
  extractedData 
}) => {
  return (
    <div className="document-upload">
      {/* Document upload content will be moved here */}
      <h3>Document Upload</h3>
      {/* Upload interface */}
    </div>
  );
};
"""

def create_navigation_component():
    """Create navigation component"""
    return """
interface NavigationProps {
  currentStep: number;
  setCurrentStep: (step: number) => void;
  steps: string[];
}

const Navigation: React.FC<NavigationProps> = ({ 
  currentStep, 
  setCurrentStep, 
  steps 
}) => {
  return (
    <nav className="navigation">
      {steps.map((step, index) => (
        <button
          key={index}
          onClick={() => setCurrentStep(index)}
          className={`nav-button ${currentStep === index ? 'active' : ''}`}
        >
          {step}
        </button>
      ))}
    </nav>
  );
};
"""

def create_validation_service():
    """Create validation service"""
    return """
// Date validation utilities for 5-year historical data limit
const getDateLimits = () => {
  const today = new Date();
  const fiveYearsAgo = new Date();
  fiveYearsAgo.setFullYear(today.getFullYear() - 5);
  
  return {
    today: today.toISOString().split('T')[0],
    fiveYearsAgo: fiveYearsAgo.toISOString().split('T')[0],
    maxBirthDate: new Date(today.getFullYear() - 17, today.getMonth(), today.getDate()).toISOString().split('T')[0],
    minBirthDate: new Date(today.getFullYear() - 130, today.getMonth(), today.getDate()).toISOString().split('T')[0],
    earliestLicenceDate: new Date(1970, 0, 1).toISOString().split('T')[0],
    ukMinDrivingAge: 17,
    maxHumanAge: 130
  };
};

// Validation functions
export const validateBirthDate = (dateString: string): { valid: boolean; error?: string } => {
  // Validation logic here
  return { valid: true };
};

export const validateHistoricalDate = (dateString: string): { valid: boolean; error?: string } => {
  // Validation logic here
  return { valid: true };
};

export const validateLicenceDate = (dateString: string, birthDate: string): { valid: boolean; error?: string } => {
  // Validation logic here
  return { valid: true };
};

export { getDateLimits };
"""

def create_api_service():
    """Create API service"""
    return """
const API_BASE_URL = 'http://localhost:3000/api';

export class ApiService {
  static async processDocument(files: FileList, uploadType: string) {
    const formData = new FormData();
    Array.from(files).forEach(file => {
      formData.append('files', file);
    });
    formData.append('uploadType', uploadType);

    const response = await fetch(`${API_BASE_URL}/process-document`, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  static async saveSession(sessionData: any) {
    const response = await fetch(`${API_BASE_URL}/save-session`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(sessionData),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  static async getOntology() {
    const response = await fetch(`${API_BASE_URL}/ontology`);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }
}
"""

def create_types_file():
    """Create TypeScript types file"""
    return """
// Core data types for CLIENT-UX application

export interface Driver {
  id: string;
  classification: 'MAIN' | 'NAMED';
  firstName: string;
  lastName: string;
  dateOfBirth: string;
  email: string;
  phone: string;
  licenceNumber: string;
  licenceIssueDate: string;
  licenceExpiryDate: string;
  licenceValidUntil: string;
  convictions: Conviction[];
}

export interface Vehicle {
  id: string;
  registration: string;
  make: string;
  model: string;
  year: number;
  engineSize: string;
  fuelType: string;
  transmission: string;
  estimatedValue: number;
  modifications: string[];
}

export interface Claim {
  id: string;
  date: string;
  type: string;
  amount: number;
  description: string;
  settled: boolean;
}

export interface Accident {
  id: string;
  date: string;
  type: string;
  description: string;
  estimatedCost: number;
  faultClaim: boolean;
}

export interface Conviction {
  id: string;
  date: string;
  offenceCode: string;
  description: string;
  penaltyPoints: number;
  fineAmount: number;
}

export interface Policy {
  startDate: string;
  coverType: string;
  excess: number;
  ncdYears: number;
  ncdProtected: boolean;
}

export interface Session {
  id: string;
  language: string;
  drivers: Driver[];
  vehicles: Vehicle[];
  claims: {
    claims: Claim[];
    accidents: Accident[];
  };
  policy: Policy;
  documents: any[];
}

export interface ValidationResult {
  valid: boolean;
  error?: string;
}

export interface DateLimits {
  today: string;
  fiveYearsAgo: string;
  maxBirthDate: string;
  minBirthDate: string;
  earliestLicenceDate: string;
  ukMinDrivingAge: number;
  maxHumanAge: number;
}
"""

def main():
    """Main function to create the component structure"""
    base_path = Path("insurance-frontend/src")
    
    # Create component files
    components = {
        "components/forms/DriverForm.tsx": (extract_driver_form_component(), "import { Driver } from '../../types';"),
        "components/forms/VehicleForm.tsx": (extract_vehicle_form_component(), "import { Vehicle } from '../../types';"),
        "components/forms/ClaimsForm.tsx": (extract_claims_form_component(), "import { Claim } from '../../types';"),
        "components/forms/DocumentUpload.tsx": (extract_document_upload_component(), ""),
        "components/layout/Navigation.tsx": (create_navigation_component(), ""),
        "services/validation.ts": (create_validation_service(), ""),
        "services/api.ts": (create_api_service(), ""),
        "types/index.ts": (create_types_file(), "")
    }
    
    print("🔧 Creating modular component structure...")
    
    for file_path, (content, imports) in components.items():
        full_path = base_path / file_path
        full_path.parent.mkdir(parents=True, exist_ok=True)
        
        if file_path.endswith('.ts'):
            # For TypeScript files, don't use the React component template
            with open(full_path, 'w') as f:
                f.write(content)
        else:
            # For React components, use the template
            component_name = Path(file_path).stem
            with open(full_path, 'w') as f:
                f.write(create_component_file(component_name, content, imports))
        
        print(f"✅ Created {file_path}")
    
    print("\n📋 Component Structure Created:")
    print("  ├── components/forms/")
    print("  │   ├── DriverForm.tsx")
    print("  │   ├── VehicleForm.tsx") 
    print("  │   ├── ClaimsForm.tsx")
    print("  │   └── DocumentUpload.tsx")
    print("  ├── components/layout/")
    print("  │   └── Navigation.tsx")
    print("  ├── services/")
    print("  │   ├── api.ts")
    print("  │   └── validation.ts")
    print("  └── types/")
    print("      └── index.ts")
    
    print("\n🎯 Next Steps:")
    print("1. Move corresponding code from App.tsx to these components")
    print("2. Update imports in App.tsx")
    print("3. Test each component individually")
    print("4. Remove original code from App.tsx")
    
    print("\n⚠️  Manual Steps Required:")
    print("- Extract actual form JSX from App.tsx into component files")
    print("- Update prop interfaces with correct types")
    print("- Add proper error handling and validation")
    print("- Test component integration")

if __name__ == "__main__":
    main()
