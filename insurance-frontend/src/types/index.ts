
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
