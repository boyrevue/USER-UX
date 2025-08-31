
// Core data types for CLIENT-UX application

export interface OntologyField {
  property: string;
  label: string;
  type: string;
  required: boolean;
  helpText?: string;
  validationPattern?: string;
  enumerationValues?: string[];
  conditionalDisplay?: string;
  conditionalRequirement?: string;
  isMultiSelect?: boolean;
  formType?: string;
  options?: Array<{
    value: string;
    label: string;
  }>;
  formInfoText?: string;
  defaultValue?: string;
  formSection?: string;
  requiresAIValidation?: boolean;
  aiValidationPrompt?: string;
  minInclusive?: number;
  maxInclusive?: number;
}

export interface OntologySection {
  id: string;
  title: string;
  icon: string;
  fields: OntologyField[];
}

export interface OntologyResponse {
  status: string;
  categories?: Record<string, any>;
  sections: {
    drivers: OntologySection;
    vehicles: OntologySection;
    claims: OntologySection;
    settings: OntologySection;
    documents: OntologySection;
    policy?: OntologySection;
  };
}

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
  // Additional fields from ontology
  title?: string;
  middleName?: string;
  occupation?: string;
  employmentStatus?: string;
  maritalStatus?: string;
  residentialStatus?: string;
  yearsAtAddress?: number;
  previousAddress?: string;
  licenceType?: string;
  endorsements?: string[];
  medicalConditions?: string[];
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
  // Additional fields from ontology
  bodyType?: string;
  colour?: string;
  doors?: number;
  seats?: number;
  mileage?: number;
  purchaseDate?: string;
  purchasePrice?: number;
  financeType?: string;
  keeper?: string;
  registeredKeeper?: string;
  previousOwners?: number;
  serviceHistory?: boolean;
  motExpiry?: string;
  taxExpiry?: string;
  insuranceGroup?: number;
  security?: string[];
  parking?: string;
  usage?: string;
  businessUse?: boolean;
}

export interface Claim {
  id: string;
  date: string;
  type: string;
  amount: number;
  description: string;
  settled: boolean;
  // Additional fields
  policyNumber?: string;
  claimNumber?: string;
  incidentDate?: string;
  reportedDate?: string;
  status?: 'Open' | 'Closed' | 'Pending' | 'Rejected';
  faultPercentage?: number;
  thirdPartyInvolved?: boolean;
  injuries?: boolean;
  policeReported?: boolean;
  witnessDetails?: string;
}

export interface Accident {
  id: string;
  date: string;
  type: string;
  description: string;
  estimatedCost: number;
  faultClaim: boolean;
  // Additional fields
  location?: string;
  weatherConditions?: string;
  roadConditions?: string;
  timeOfDay?: string;
  vehiclesDamaged?: number;
  injuriesSustained?: boolean;
  emergencyServices?: boolean;
  witnessPresent?: boolean;
  dashcamFootage?: boolean;
  policeAttended?: boolean;
  breathalyzerTest?: boolean;
}

export interface Conviction {
  id: string;
  date: string;
  offenceCode: string;
  description: string;
  penaltyPoints: number;
  fineAmount: number;
  // Additional fields
  courtDate?: string;
  disqualificationPeriod?: number;
  endorsementCode?: string;
  spentDate?: string;
  rehabilitationPeriod?: number;
  drivingBan?: boolean;
  communityService?: boolean;
  prisonSentence?: boolean;
}

export interface Policy {
  startDate: string;
  coverType: string;
  excess: number;
  ncdYears: number;
  ncdProtected: boolean;
  // Additional fields
  endDate?: string;
  policyNumber?: string;
  previousInsurer?: string;
  renewalDate?: string;
  paymentMethod?: 'Annual' | 'Monthly';
  directDebit?: boolean;
  voluntaryExcess?: number;
  compulsoryExcess?: number;
  coverLevel?: 'Third Party' | 'Third Party Fire & Theft' | 'Comprehensive';
  europeanCover?: boolean;
  breakdown?: boolean;
  legalExpenses?: boolean;
  personalAccident?: boolean;
  keyReplacement?: boolean;
  windscreenCover?: boolean;
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
  documents: Document[];
  // Additional session metadata
  createdAt?: string;
  updatedAt?: string;
  status?: 'Draft' | 'Complete' | 'Submitted';
  currentStep?: number;
  totalSteps?: number;
  validationErrors?: ValidationError[];
}

export interface Document {
  id: string;
  type: 'passport' | 'driving_licence' | 'utility_bill' | 'bank_statement' | 'other';
  name: string;
  size: number;
  uploadedAt: string;
  processed: boolean;
  fieldName?: string;
  extractedData?: Record<string, any>;
  confidence?: number;
  imagePaths?: {
    original?: string;
    page1?: string;
    page2Upper?: string;
    page2MRZ?: string;
    preprocessed?: string;
  };
}

export interface ValidationResult {
  valid: boolean;
  error?: string;
  warnings?: string[];
}

export interface ValidationError {
  field: string;
  message: string;
  code: string;
  severity: 'error' | 'warning' | 'info';
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

export interface FormField {
  property: string;
  label: string;
  type: 'text' | 'email' | 'tel' | 'date' | 'number' | 'select' | 'checkbox' | 'textarea';
  required: boolean;
  helpText?: string;
  validationPattern?: string;
  enumerationValues?: string[];
  placeholder?: string;
  min?: number | string;
  max?: number | string;
  step?: number;
}

export interface FormSection {
  label: string;
  fields: FormField[];
}

export interface AppState {
  currentStep: number;
  session: Session;
  loading: boolean;
  error: string | null;
  ontology: OntologyResponse | null;
  validationErrors: ValidationError[];
}

// Component Props Interfaces
export interface DriverFormProps {
  driver: Driver;
  index: number;
  updateDriver: (index: number, field: string, value: any) => void;
  removeDriver: (index: number) => void;
  validationErrors?: ValidationError[];
}

export interface VehicleFormProps {
  vehicle: Vehicle;
  index: number;
  updateVehicle: (index: number, field: string, value: any) => void;
  removeVehicle: (index: number) => void;
  validationErrors?: ValidationError[];
}

export interface ClaimsFormProps {
  claims: Claim[];
  accidents: Accident[];
  updateClaim: (index: number, field: string, value: any) => void;
  updateAccident: (index: number, field: string, value: any) => void;
  addClaim: () => void;
  addAccident: () => void;
  removeClaim: (index: number) => void;
  removeAccident: (index: number) => void;
  validationErrors?: ValidationError[];
}

export interface DocumentUploadProps {
  onFileUpload: (files: FileList, type: string, selectedDocumentType?: string) => void;
  isProcessing: boolean;
  extractedData?: Record<string, any>;
  documents: Document[];
  onRemoveDocument: (id: string) => void;
}

export interface NavigationProps {
  currentStep: number;
  setCurrentStep: (step: number) => void;
  steps: string[];
  canNavigate: (step: number) => boolean;
  completedSteps: number[];
}
