// UK Insurance Driver Licence Types - JSON Schema Compliant
// Generated from JSON Schema v1.0.0-2025-01-24

export type LicenceType = "FULL_UK" | "PROVISIONAL_UK" | "EU_EEA" | "INTERNATIONAL" | "OTHER_FOREIGN";
export type YesNo = "YES" | "NO";
export type ManualAuto = "MANUAL" | "AUTOMATIC";
export type LicenceCategory = "B" | "B1" | "A" | "AM" | "Q" | "C" | "C1" | "D" | "D1" | "BE" | "C1E" | "CE" | "D1E" | "DE";

export interface Endorsement {
  code: string;           // DVLA code e.g., "SP30", pattern: ^[A-Z]{2}\d{2}$
  points: number;         // 0–12 penalty points
  offenceDate: string;    // YYYY-MM-DD format
  expiryDate?: string;    // YYYY-MM-DD format (optional)
}

export interface Disqualification {
  from: string;           // YYYY-MM-DD start date
  to: string;             // YYYY-MM-DD end date  
  reason: string;         // max 200 characters
}

export interface DriverLicence {
  // Core Required Fields
  licenceType: LicenceType;
  licenceCategory: LicenceCategory[];  // Array, minimum 1 item, unique items
  manualOrAuto: ManualAuto;
  dateFirstIssued: string;             // YYYY-MM-DD pattern
  countryOfIssue: string;              // minimum 2 characters
  hasEndorsements: YesNo;
  hasDisqualifications: YesNo;
  hasMedicalConditions: YesNo;
  visionCorrectionRequired: YesNo;
  
  // Optional Fields
  yearsHeldFull?: number;              // 0-80 range
  exchangedToUK?: YesNo;               // Required if licenceType is EU_EEA, INTERNATIONAL, or OTHER_FOREIGN
  
  // Conditional Arrays (required if corresponding "has" field is YES)
  endorsements?: Endorsement[];        // Required if hasEndorsements = YES, minimum 1 item
  disqualifications?: Disqualification[]; // Required if hasDisqualifications = YES, minimum 1 item
  
  // Medical Conditional
  medicalDeclaredToDVLA?: YesNo;       // Required if hasMedicalConditions = YES
  
  // Optional Additional Fields
  pendingProsecutions?: YesNo;
  licenceNumber?: string;              // Optional DVLA pattern: ^[A-Z9]{5}\d{6}[A-Z9]{2}\d{2}$
}

// Validation Helper Types
export interface LicenceValidationError {
  field: string;
  message: string;
  code: string;
}

export interface LicenceValidationResult {
  isValid: boolean;
  errors: LicenceValidationError[];
}

// Form State Management
export interface LicenceFormState {
  data: Partial<DriverLicence>;
  errors: LicenceValidationError[];
  touched: Set<string>;
  isSubmitting: boolean;
}

// JSON Schema Validation Constants
export const LICENCE_VALIDATION_PATTERNS = {
  DATE: /^\d{4}-\d{2}-\d{2}$/,
  DVLA_CODE: /^[A-Z]{2}\d{2}$/,
  LICENCE_NUMBER: /^[A-Z9]{5}\d{6}[A-Z9]{2}\d{2}$/
} as const;

export const LICENCE_CONSTRAINTS = {
  YEARS_HELD_MIN: 0,
  YEARS_HELD_MAX: 80,
  POINTS_MIN: 0,
  POINTS_MAX: 12,
  REASON_MAX_LENGTH: 200,
  COUNTRY_MIN_LENGTH: 2
} as const;

// Conditional Logic Helpers
export const isNonUKLicence = (licenceType: LicenceType): boolean => {
  return ['EU_EEA', 'INTERNATIONAL', 'OTHER_FOREIGN'].includes(licenceType);
};

export const requiresExchangeField = (licenceType: LicenceType): boolean => {
  return isNonUKLicence(licenceType);
};

export const requiresEndorsementDetails = (hasEndorsements: YesNo): boolean => {
  return hasEndorsements === 'YES';
};

export const requiresDisqualificationDetails = (hasDisqualifications: YesNo): boolean => {
  return hasDisqualifications === 'YES';
};

export const requiresMedicalDVLADeclaration = (hasMedicalConditions: YesNo): boolean => {
  return hasMedicalConditions === 'YES';
};

// UK Licence Category Descriptions
export const LICENCE_CATEGORY_DESCRIPTIONS = {
  'B': 'Cars & small vans up to 3,500kg (most common)',
  'B1': 'Light quadricycles and motor tricycles',
  'A': 'Motorcycles (various power restrictions)',
  'AM': 'Mopeds (speed limited to 28mph)',
  'Q': 'Two or three-wheeled vehicles without pedals',
  'C': 'Large goods vehicles over 3,500kg',
  'C1': 'Medium goods vehicles 3,500-7,500kg',
  'D': 'Buses with more than 8 passenger seats',
  'D1': 'Minibuses 9-16 passenger seats',
  'BE': 'Car with trailer over 3,500kg',
  'C1E': 'Medium lorry with trailer',
  'CE': 'Large lorry with trailer',
  'D1E': 'Minibus with trailer',
  'DE': 'Bus with trailer'
} as const;

// Common DVLA Endorsement Codes
export const DVLA_ENDORSEMENT_CODES = {
  'SP30': 'Speeding (3-6 points)',
  'SP10': 'Exceeding goods vehicle speed limit (3-6 points)',
  'SP20': 'Exceeding speed limit for type of vehicle (3-6 points)',
  'DR10': 'Driving or attempting to drive with alcohol level above limit (3-11 points)',
  'DR20': 'Driving or attempting to drive while unfit through drink (3-11 points)',
  'IN10': 'Using vehicle uninsured against third party risks (6-8 points)',
  'CD10': 'Driving without due care and attention (3-9 points)',
  'DD40': 'Dangerous driving (3-11 points)',
  'LC20': 'Driving otherwise than in accordance with licence (3-6 points)',
  'TS10': 'Failing to comply with traffic light signals (3 points)',
  'AC10': 'Failing to stop after accident (5-10 points)'
} as const;
