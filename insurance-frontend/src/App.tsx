import React, { useState } from 'react';
import { 
  User, 
  Car, 
  FileText, 
  Shield, 
  CreditCard, 
  Settings,
  Plus,
  X,
  ChevronRight,
  ChevronLeft,
  AlertCircle,
  HelpCircle,
  Upload,
  Mail,
  CheckCircle,
  Navigation,
  Brain,
  Scan,
  Key,
  Zap,
  FileSearch,
  MapPin,
  Target,
  Wand2,
  Download,
  Globe
} from 'lucide-react';
import { 
  Card, 
  Button,
  Badge,
  Modal,
  TextInput,
  Label,
  Select,
  ToggleSwitch,
  Progress,
  Alert,
  Breadcrumb,
  BreadcrumbItem
} from 'flowbite-react';

// Types based on the ontology
interface Driver {
  id: string;
  classification: 'MAIN' | 'NAMED';
  firstName: string;
  lastName: string;
  dateOfBirth: string;
  email: string;
  phone: string;
  address: string;
  postcode: string;
  licenceType: string;
  licenceNumber: string;
  yearsHeld: number;
  pointsLost: number;
  licenceIssueDate: string;
  licenceExpiryDate: string;
  licenceValidUntil: string;
  relationship: string;
  sameAddress: boolean;
  // Passport-specific fields (document-derived)
  passportNumber?: string;
  passportIssueDate?: string;
  passportExpiryDate?: string;
  passportAuthority?: string;
  placeOfBirth?: string;
  gender?: string;
  nationality?: string;
  // Driving Licence-specific fields (document-derived)
  licenceAuthority?: string;
  entitlementA?: string;
  entitlementB?: string;
  entitlementC?: string;
  entitlementD?: string;
  entitlementBE?: string;
  entitlementCE?: string;
  entitlementDE?: string;
  licenceRestrictions?: string;
  licenceEndorsements?: string;
  // Disabilities and Restrictions
  hasDisability: boolean;
  disabilityTypes: string[];
  requiresAdaptations: boolean;
  adaptationTypes: string[];
  automaticOnly: boolean;
  // Licence Classes
  licenceClassA: boolean;
  licenceClassB: boolean;
  licenceClassC: boolean;
  licenceClassD: boolean;
  licenceClassBE: boolean;
  licenceClassCE: boolean;
  licenceClassDE: boolean;
  // Driver Communication
  driverEmail: string;
  emailSent: boolean;
  emailSentDate: string;
}

interface Vehicle {
  registration: string;
  make: string;
  model: string;
  year: number;
  mileage: number;
  value: number;
  daytimeLocation: string;
  overnightLocation: string;
  hasModifications: boolean;
  modifications: string[];
}

interface Policy {
  coverType: string;
  startDate: string;
  voluntaryExcess: number;
  ncdYears: number;
  protectNCD: boolean;
}

interface Claim {
  id: string;
  date: string;
  type: string;
  amount: number;
  status: string;
  faultStatus: string;
  policeReported: boolean;
  crimeReference?: string;
  description?: string;
}

interface Conviction {
  id: string;
  date: string;
  offenceCode: string;
  description?: string;
  penaltyPoints: number;
  fine?: number;
  disqualification?: number;
  status: string;
}

interface Accident {
  id: string;
  date: string;
  type: string;
  severity: string;
  claimMade: boolean;
  description?: string;
}

interface ClaimsHistory {
  hasClaims: boolean;
  claims: Claim[];
  hasConvictions: boolean;
  convictions: Conviction[];
  hasAccidents: boolean;
  accidents: Accident[];
}

interface Payment {
  frequency: string;
  method: string;
}

interface Extras {
  breakdownCover: string;
  legalExpenses: boolean;
  courtesyCar: boolean;
}

interface Marketing {
  emailMarketing: boolean;
  smsMarketing: boolean;
  postMarketing: boolean;
}

interface QuoteSession {
  id: string;
  language: string;
  drivers: Driver[];
  vehicle: Vehicle;
  policy: Policy;
  claims: ClaimsHistory;
  payment: Payment;
  extras: Extras;
  marketing: Marketing;
  progress: Record<string, boolean>;
  formData: Record<string, any>;
  createdAt: string;
  lastAccessed: string;
  completedAt?: string;
}



// Language translations
const translations: {[key: string]: {[key: string]: string}} = {
  en: {
    driverDetails: 'Driver Details',
    vehicleDetails: 'Vehicle Details',
    policyDetails: 'Policy Details',
    claimsHistory: 'Claims History',
    paymentExtras: 'Insurance Payments',
    marketingPreferences: 'Marketing Preferences',
    personalDocuments: 'Personal Documents',

    bitwardenIntegration: 'Bitwarden Integration',
    bitwardenIntegrationDesc: 'Auto-fill forms using your Bitwarden vault data',
    smartMapping: 'Smart Mapping',
    smartMappingDesc: 'Intelligent field mapping with SHACL transformation',
    // Web Spider Features
    webSpider: 'Web Spider',
    headlessMode: 'Headless Mode',
    takeScreenshots: 'Take Screenshots',
    waitForLoad: 'Wait for Page Load',
    enableJavaScript: 'Enable JavaScript',
    cssSelectors: 'CSS Selectors',

    extractionResults: 'Extraction Results',
    recentTasks: 'Recent Tasks',
    // Bitwarden & Stealth Browser
    bitwardenUnlocked: 'Bitwarden vault unlocked successfully',
    templatesCreated: 'Credential templates created successfully',
    browserSessionCreated: 'Stealth browser session created',
    loginSuccessful: 'Login successful',
    loginFailed: 'Login may have failed',
    stealthBrowser: 'Stealth Browser',
    unlockVault: 'Unlock Vault',
    setupTemplates: 'Setup Credential Templates',
    autoLogin: 'Auto Login',
    antiBotProtection: 'Anti-Bot Protection',
    // Money Supermarket Infiltration
    moneySupermarketInfiltration: 'Money Supermarket Infiltration',
    infiltrateMoneySupermarket: 'Infiltrate Money Supermarket',
    stealthInfiltration: 'Stealth Infiltration',
    targetAcquired: 'Target Acquired',
    infiltrationInProgress: 'Infiltration in Progress...',
    infiltrationComplete: 'Infiltration Complete!',
    quotesExtracted: 'Quotes Extracted',
    evidenceCaptured: 'Evidence Captured',
    drivingLicence: 'Driving Licence',
    identityCard: 'Identity Card',
    utilityBill: 'Utility Bill',
    vehicleRegistration: 'Vehicle Registration',
    bankStatement: 'Bank Statement',
    medicalCertificate: 'Medical Certificate',
    insuranceQuote: 'Insurance Quote',
    insurancePolicy: 'Insurance Policy',
    documentProcessing: 'Document Processing',
    uploadDocument: 'Upload Document',
    dragDropFiles: 'Drag and drop files here, or click to browse',
    dropFilesHere: 'Drop files here',
    chooseFile: 'Choose File',
    processing: 'Processing...',
    processDocument: 'Process Document',
    backToSelection: 'Back to Selection',
    cancel: 'Cancel',
    uploadedFiles: 'Uploaded Files',
    numberOfPassports: 'Number of Passports',
    passportUpload: 'Upload passport document',
    drivingLicenceUpload: 'Upload driving licence (front & back)',
    identityCardUpload: 'Upload national ID or identity card',
    utilityBillUpload: 'Upload proof of address',
    vehicleRegistrationUpload: 'Upload V5C or registration document',
    bankStatementUpload: 'Upload recent bank statement',
    medicalCertificateUpload: 'Upload medical or fitness certificate',
    insuranceQuoteUpload: 'Upload existing insurance quote',
    insurancePolicyUpload: 'Upload current/previous policy',
    frontSide: 'Front Side',
    backSide: 'Back Side',
    passportNumber: 'Passport Number',
    issuingCountry: 'Issuing Country',
    issueDate: 'Issue Date',
    expiryDate: 'Expiry Date',

    selectCountry: 'Select country',
    selectType: 'Select type',
    selectGender: 'Select gender',
    enterPassportNumber: 'Enter passport number',
    enterLicenceNumber: 'Enter licence number',
    asShownOnPassport: 'As shown on passport',
    asShownOnDocument: 'As shown on document',
    carInsurance: 'Car Insurance',
    settings: 'Settings',
    driverDetailsDesc: 'Information about all drivers',
    vehicleDetailsDesc: 'Vehicle information and modifications',
    policyDetailsDesc: 'Coverage and policy options',
    claimsHistoryDesc: 'Previous claims and convictions',
    paymentExtrasDesc: 'Payment method and additional cover',
    marketingPreferencesDesc: 'Communication preferences',
    mainDriver: 'Main Driver',
    additionalDriver: 'Additional Driver',
    firstName: 'First Name',
    lastName: 'Last Name',
    email: 'Email',
    phone: 'Phone',
    address: 'Address',
    postcode: 'Postcode',
    yearsHeld: 'Years Held',
    pointsLost: 'Points Lost',
    licenceIssueDate: 'Licence Issue Date',
    licenceExpiryDate: 'Licence Expiry Date',
    licenceValidUntil: 'Licence Valid Until',
    relationship: 'Relationship',
    sameAddress: 'Same Address',
    addDriver: 'Add Driver',
    removeDriver: 'Remove Driver',
    uploadDocuments: 'Upload Documents',
    aiAssistant: 'AI Assistant',
    language: 'Language',
    // Claims and Convictions
    convictions: 'Convictions',
    accidents: 'Accidents',
    addClaim: 'Add Claim',
    addConviction: 'Add Conviction',
    addAccident: 'Add Accident',
    claimDate: 'Claim Date',
    claimType: 'Claim Type',
    claimAmount: 'Claim Amount',
    claimStatus: 'Claim Status',
    faultStatus: 'Fault Status',
    policeReported: 'Police Reported',
    crimeReference: 'Crime Reference',
    convictionDate: 'Conviction Date',
    offenceCode: 'Offence Code',
    penaltyPoints: 'Penalty Points',
    fine: 'Fine',
    disqualification: 'Disqualification',
    accidentDate: 'Accident Date',
    accidentType: 'Accident Type',
    accidentSeverity: 'Accident Severity',
    claimMade: 'Claim Made',
    // Disabilities and Restrictions
    hasDisability: 'Has Disability',
    disabilityType: 'Disability Type',
    requiresAdaptations: 'Requires Adaptations',
    adaptationType: 'Adaptation Type',
    automaticOnly: 'Automatic Only',
    // Licence Classes
    licenceClasses: 'Licence Classes',
    licenceClassA: 'Class A (Motorcycles)',
    licenceClassB: 'Class B (Cars)',
    licenceClassC: 'Class C (Large Vehicles)',
    licenceClassD: 'Class D (Buses)',
    licenceClassBE: 'Class BE (Car + Trailer)',
    licenceClassCE: 'Class CE (Large Vehicle + Trailer)',
    licenceClassDE: 'Class DE (Bus + Trailer)',
    // Driver Communication
    driverEmail: 'Driver Email',
    emailSent: 'Email Sent',
    emailSentDate: 'Email Sent Date',
    sendEmailToDriver: 'Send Email to Driver',
    emailDriverForm: 'Email Driver Form'
  },
  de: {
    driverDetails: 'Fahrerdetails',
    vehicleDetails: 'Fahrzeugdetails',
    policyDetails: 'Policendetails',
    claimsHistory: 'Schadenshistorie',
    paymentExtras: 'Zahlung & Extras',
    marketingPreferences: 'Marketing-Einstellungen',

    bitwardenIntegration: 'Bitwarden-Integration',
    bitwardenIntegrationDesc: 'Formulare automatisch mit Ihren Bitwarden-Vault-Daten ausfüllen',
    smartMapping: 'Intelligente Zuordnung',
    smartMappingDesc: 'Intelligente Feldzuordnung mit SHACL-Transformation',
    // Web Spider Features
    webSpider: 'Web-Spider',
    headlessMode: 'Headless-Modus',
    takeScreenshots: 'Screenshots erstellen',
    waitForLoad: 'Auf Seitenladen warten',
    enableJavaScript: 'JavaScript aktivieren',
    cssSelectors: 'CSS-Selektoren',

    extractionResults: 'Extraktionsergebnisse',
    recentTasks: 'Letzte Aufgaben',
    // Bitwarden & Stealth Browser
    bitwardenUnlocked: 'Bitwarden-Tresor erfolgreich entsperrt',
    templatesCreated: 'Anmeldedaten-Vorlagen erfolgreich erstellt',
    browserSessionCreated: 'Stealth-Browser-Sitzung erstellt',
    loginSuccessful: 'Anmeldung erfolgreich',
    loginFailed: 'Anmeldung möglicherweise fehlgeschlagen',
    stealthBrowser: 'Stealth-Browser',
    unlockVault: 'Tresor entsperren',
    setupTemplates: 'Anmeldedaten-Vorlagen einrichten',
    autoLogin: 'Automatische Anmeldung',
    antiBotProtection: 'Anti-Bot-Schutz',
    // Money Supermarket Infiltration
    moneySupermarketInfiltration: 'Money Supermarket Infiltration',
    infiltrateMoneySupermarket: 'Money Supermarket infiltrieren',
    stealthInfiltration: 'Stealth-Infiltration',
    targetAcquired: 'Ziel erfasst',
    infiltrationInProgress: 'Infiltration läuft...',
    infiltrationComplete: 'Infiltration abgeschlossen!',
    quotesExtracted: 'Angebote extrahiert',
    evidenceCaptured: 'Beweise gesichert',
    additionalDriver: 'Zusatzfahrer',
    firstName: 'Vorname',
    lastName: 'Nachname',
    dateOfBirth: 'Geburtsdatum',
    email: 'E-Mail',
    phone: 'Telefon',
    address: 'Adresse',
    postcode: 'Postleitzahl',
    yearsHeld: 'Jahre gehalten',
    pointsLost: 'Punkte verloren',
    licenceIssueDate: 'Führerschein Ausstellungsdatum',
    licenceExpiryDate: 'Führerschein Ablaufdatum',
    licenceValidUntil: 'Führerschein gültig bis',
    relationship: 'Beziehung',
    sameAddress: 'Gleiche Adresse',
    addDriver: 'Fahrer hinzufügen',
    removeDriver: 'Fahrer entfernen',
    uploadDocuments: 'Dokumente hochladen',
    aiAssistant: 'KI-Assistent',
    language: 'Sprache',
    // Claims and Convictions
    convictions: 'Verurteilungen',
    accidents: 'Unfälle',
    addClaim: 'Schaden hinzufügen',
    addConviction: 'Verurteilung hinzufügen',
    addAccident: 'Unfall hinzufügen',
    claimDate: 'Schadensdatum',
    claimType: 'Schadenstyp',
    claimAmount: 'Schadenshöhe',
    claimStatus: 'Schadensstatus',
    faultStatus: 'Verschuldensstatus',
    policeReported: 'Polizei gemeldet',
    crimeReference: 'Strafverfahrensnummer',
    convictionDate: 'Verurteilungsdatum',
    offenceCode: 'Verstoßcode',
    penaltyPoints: 'Punkte',
    fine: 'Bußgeld',
    disqualification: 'Fahrverbot',
    accidentDate: 'Unfalldatum',
    accidentType: 'Unfalltyp',
    accidentSeverity: 'Unfallschwere',
    claimMade: 'Schaden gemeldet',
    // Disabilities and Restrictions
    hasDisability: 'Hat Behinderung',
    disabilityType: 'Behinderungstyp',
    requiresAdaptations: 'Benötigt Anpassungen',
    adaptationType: 'Anpassungstyp',
    automaticOnly: 'Nur Automatik',
    // Licence Classes
    licenceClasses: 'Führerscheinklassen',
    licenceClassA: 'Klasse A (Motorräder)',
    licenceClassB: 'Klasse B (Autos)',
    licenceClassC: 'Klasse C (Großfahrzeuge)',
    licenceClassD: 'Klasse D (Busse)',
    licenceClassBE: 'Klasse BE (Auto + Anhänger)',
    licenceClassCE: 'Klasse CE (Großfahrzeug + Anhänger)',
    licenceClassDE: 'Klasse DE (Bus + Anhänger)',
    // Driver Communication
    driverEmail: 'Fahrer E-Mail',
    emailSent: 'E-Mail gesendet',
    emailSentDate: 'E-Mail gesendet am',
    sendEmailToDriver: 'E-Mail an Fahrer senden',
    emailDriverForm: 'Fahrerformular per E-Mail',
    personalDocuments: 'Persönliche Dokumente',
    documentProcessing: 'Dokumentenverarbeitung',
    uploadDocument: 'Dokument hochladen',
    dragDropFiles: 'Dateien hier ablegen oder klicken zum Durchsuchen',
    dropFilesHere: 'Dateien hier ablegen',
    chooseFile: 'Datei auswählen',
    processing: 'Verarbeitung...',
    processDocument: 'Dokument verarbeiten',
    backToSelection: 'Zurück zur Auswahl',
    cancel: 'Abbrechen',
    uploadedFiles: 'Hochgeladene Dateien',
    numberOfPassports: 'Anzahl der Pässe',
    passport: 'Reisepass',
    passportUpload: 'Reisepass-Dokument hochladen',
    drivingLicence: 'Führerschein',
    drivingLicenceUpload: 'Führerschein hochladen (Vorder- und Rückseite)',
    identityCard: 'Personalausweis',
    identityCardUpload: 'Personalausweis oder Identitätskarte hochladen',
    utilityBill: 'Nebenkostenabrechnung',
    utilityBillUpload: 'Adressnachweis hochladen',
    vehicleRegistration: 'Fahrzeugschein',
    vehicleRegistrationUpload: 'Fahrzeugschein oder Zulassungsbescheinigung hochladen',
    bankStatement: 'Kontoauszug',
    bankStatementUpload: 'Aktuellen Kontoauszug hochladen',
    medicalCertificate: 'Ärztliches Attest',
    medicalCertificateUpload: 'Ärztliches Attest oder Tauglichkeitszeugnis hochladen',
    insuranceQuote: 'Versicherungsangebot',
    insuranceQuoteUpload: 'Bestehendes Versicherungsangebot hochladen',
    insurancePolicy: 'Versicherungspolice',
    insurancePolicyUpload: 'Aktuelle/vorherige Versicherungspolice hochladen',
    frontSide: 'Vorderseite',
    backSide: 'Rückseite',
    passportNumber: 'Reisepassnummer',
    issuingCountry: 'Ausstellendes Land',
    issueDate: 'Ausstellungsdatum',
    expiryDate: 'Ablaufdatum',

    selectCountry: 'Land auswählen',
    selectType: 'Typ auswählen',
    selectGender: 'Geschlecht auswählen',
    enterPassportNumber: 'Reisepassnummer eingeben',
    enterLicenceNumber: 'Führerscheinnummer eingeben',
    asShownOnPassport: 'Wie im Reisepass angegeben',
    asShownOnDocument: 'Wie im Dokument angegeben',
    carInsurance: 'Autoversicherung',
    settings: 'Einstellungen',
    driverDetailsDesc: 'Informationen über alle Fahrer',
    vehicleDetailsDesc: 'Fahrzeuginformationen und Modifikationen',
    policyDetailsDesc: 'Deckung und Policenoptionen',
    claimsHistoryDesc: 'Frühere Schäden und Verurteilungen',
    paymentExtrasDesc: 'Zahlungsmethode und zusätzliche Deckung',
    marketingPreferencesDesc: 'Kommunikationseinstellungen'
  }
};

function App() {
  const [activeCategory, setActiveCategory] = useState(0);
  const [expandedSections, setExpandedSections] = useState<Set<string>>(new Set(['car-insurance']));
  const [showDocumentModal, setShowDocumentModal] = useState(false);
  const [selectedDocumentType, setSelectedDocumentType] = useState<string>('');
  const [documentFiles, setDocumentFiles] = useState<{[key: string]: File[]}>({});
  const [isDragOverDocument, setIsDragOverDocument] = useState(false);
  const [extractedData, setExtractedData] = useState<{[key: string]: any}>({});
  const [isProcessing, setIsProcessing] = useState(false);
  const [passportCount, setPassportCount] = useState(1);
  
  // Browser state
  const [browserUrl, setBrowserUrl] = useState('');
  const [isBrowsing, setIsBrowsing] = useState(false);
  const [browserWindow, setBrowserWindow] = useState<Window | null>(null);
  
  // AI Form Analysis State
  const [showFormAnalyzer, setShowFormAnalyzer] = useState(false);
  const [analysisResults, setAnalysisResults] = useState<any>(null);
  const [fieldMappings, setFieldMappings] = useState<any[]>([]);
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [analysisType, setAnalysisType] = useState<'document' | 'html' | 'url' | 'stealth'>('document');
  
  // Web Spider State
  const [spiderTasks, setSpiderTasks] = useState<any[]>([]);
  const [isSpiderRunning, setIsSpiderRunning] = useState(false);
  const [extractionResults, setExtractionResults] = useState<any>(null);
  const [spiderConfig, setSpiderConfig] = useState({
    headless: true,
    screenshots: false,
    waitForLoad: true,
    javascript: true
  });

  // Stealth Browser State
  const [browserSession, setBrowserSession] = useState<any>(null);
  const [isBrowserRunning, setIsBrowserRunning] = useState(false);
  const [browserScreenshot, setBrowserScreenshot] = useState<string>('');
  
  // Bitwarden State
  const [bitwardenUnlocked, setBitwardenUnlocked] = useState(false);
  const [availableCredentials, setAvailableCredentials] = useState<any[]>([]);
  const [selectedSite, setSelectedSite] = useState<string>('moneysupermarket.com');

  // Money Supermarket Infiltration State
  const [infiltrationInProgress, setInfiltrationInProgress] = useState(false);
  const [infiltrationResult, setInfiltrationResult] = useState<any>(null);
  const [extractedQuotes, setExtractedQuotes] = useState<any[]>([]);
  const [infiltrationLogs, setInfiltrationLogs] = useState<string[]>([]);



  const toggleSection = (sectionId: string) => {
    setExpandedSections(prev => {
      const newSet = new Set(prev);
      if (newSet.has(sectionId)) {
        newSet.delete(sectionId);
      } else {
        newSet.add(sectionId);
      }
      return newSet;
    });
  };

  // Document upload handlers
  const handleDocumentDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOverDocument(true);
  };

  const handleDocumentDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOverDocument(false);
  };

  const handleDocumentDrop = (e: React.DragEvent, uploadType: string = 'main') => {
    e.preventDefault();
    setIsDragOverDocument(false);
    const files = Array.from(e.dataTransfer.files);
    handleDocumentFiles(files, uploadType);
  };

  const handleDocumentFileSelect = (e: React.ChangeEvent<HTMLInputElement>, uploadType: string = 'main') => {
    console.log('File selection triggered:', { uploadType, selectedDocumentType });
    const files = Array.from(e.target.files || []);
    console.log('Files selected:', files.map(f => f.name));
    
    // Auto-process files immediately
    handleDocumentFiles(files, uploadType);
    
    // Clear the input so the same file can be selected again
    e.target.value = '';
  };

  const handleDocumentFiles = (files: File[], uploadType: string = 'main') => {
    console.log('Processing files:', { files: files.map(f => f.name), uploadType, selectedDocumentType });
    
    if (!selectedDocumentType) {
      console.error('No document type selected!');
      return;
    }
    
    const key = `${selectedDocumentType}_${uploadType}`;
    setDocumentFiles(prev => ({
      ...prev,
      [key]: [...(prev[key] || []), ...files]
    }));
    
    // Auto-process if files are uploaded
    if (files.length > 0) {
      console.log('Starting OCR processing...');
      processDocumentOCR(files, uploadType);
    }
  };

  const removeDocumentFile = (fileIndex: number, uploadType: string = 'main') => {
    const key = `${selectedDocumentType}_${uploadType}`;
    setDocumentFiles(prev => ({
      ...prev,
      [key]: prev[key]?.filter((_, index) => index !== fileIndex) || []
    }));
  };

  // OCR Processing Functions
  const processDocumentOCR = async (files: File[], uploadType: string = 'main') => {
    console.log('processDocumentOCR called:', { files: files.map(f => f.name), uploadType, selectedDocumentType });
    
    if (!selectedDocumentType || files.length === 0) {
      console.error('Cannot process OCR:', { selectedDocumentType, filesLength: files.length });
      return;
    }
    
    setIsProcessing(true);
    
    try {
      for (const file of files) {
        console.log('Processing file:', file.name);
        const formData = new FormData();
        formData.append('file', file);
        formData.append('documentType', selectedDocumentType);
        formData.append('uploadType', uploadType);
        
        console.log('Sending OCR request to /api/process-document');
        const response = await fetch('/api/process-document', {
          method: 'POST',
          body: formData
        });
        
        console.log('OCR response status:', response.status);
        
        if (response.ok) {
          const result = await response.json();
          console.log('OCR result:', result);
          console.log('🔍 Extracted Fields Detail:', result.extractedFields);
          
          // Log clickable image URLs
          if (result.extractedFields) {
            const fields = result.extractedFields;
            console.group('📸 Extracted Images (Click to view):');
            
            if (fields._page1ImageUrl) {
              console.log('📄 Page 1 (Top Half):', `http://localhost:3000${fields._page1ImageUrl}`);
            }
            if (fields._page2UpperImageUrl) {
              console.log('📄 Page 2 Upper (Middle):', `http://localhost:3000${fields._page2UpperImageUrl}`);
            }
            if (fields._page2MrzImageUrl) {
              console.log('📄 Page 2 MRZ (Bottom):', `http://localhost:3000${fields._page2MrzImageUrl}`);
            }
            if (fields._page2MrzPreprocessedUrl) {
              console.log('🎨 Page 2 MRZ Preprocessed:', `http://localhost:3000${fields._page2MrzPreprocessedUrl}`);
            }
            
            console.groupEnd();
          }
          
          console.log('📊 OCR Metadata:', {
            engine: result.ocrEngine || result.extractedFields?._ocrEngine,
            confidence: result.confidence,
            processingPath: result.processingPath,
            ocrConfidence: result.extractedFields?._ocrConfidence,
            fieldConfidence: result.extractedFields?._fieldConfidence
          });
          
          setExtractedData(prev => ({
            ...prev,
            [`${selectedDocumentType}_${uploadType}_${file.name}`]: result
          }));
          
          // Auto-populate form fields with extracted data
          populateFormFields(result.extractedFields, uploadType);
        } else {
          const errorText = await response.text();
          console.error('OCR request failed:', response.status, errorText);
        }
      }
    } catch (error) {
      console.error('OCR processing failed:', error);
    } finally {
      setIsProcessing(false);
    }
  };

  // Map country name back to ISO code for form selection
  const getCountryCodeFromName = (countryName: string): string => {
    const countryMap: { [key: string]: string } = {
      'United Kingdom': 'GB',
      'United States': 'US',
      'Germany': 'DE',
      'France': 'FR',
      'Italy': 'IT',
      'Spain': 'ES',
      'Netherlands': 'NL',
      'Belgium': 'BE',
      'Austria': 'AT',
      'Switzerland': 'CH',
      'Ireland': 'IE',
      'Portugal': 'PT',
      'Greece': 'GR',
      'Poland': 'PL',
      'Czech Republic': 'CZ',
      'Hungary': 'HU',
      'Slovakia': 'SK',
      'Slovenia': 'SI',
      'Croatia': 'HR',
      'Romania': 'RO',
      'Bulgaria': 'BG',
      'Lithuania': 'LT',
      'Latvia': 'LV',
      'Estonia': 'EE',
      'Finland': 'FI',
      'Sweden': 'SE',
      'Denmark': 'DK',
      'Norway': 'NO',
      'Iceland': 'IS',
      'Luxembourg': 'LU',
      'Malta': 'MT',
      'Cyprus': 'CY',
      'Canada': 'CA',
      'Australia': 'AU',
      'New Zealand': 'NZ',
      'Japan': 'JP',
      'South Korea': 'KR',
      'China': 'CN',
      'India': 'IN',
      'Brazil': 'BR',
      'Mexico': 'MX',
      'Argentina': 'AR',
      'Chile': 'CL',
      'South Africa': 'ZA',
      'Russia': 'RU',
      'Turkey': 'TR',
      'Israel': 'IL',
      'United Arab Emirates': 'AE',
      'Saudi Arabia': 'SA',
      'Singapore': 'SG',
      'Malaysia': 'MY',
      'Thailand': 'TH',
      'Indonesia': 'ID',
      'Philippines': 'PH',
      'Vietnam': 'VN'
    };
    
    return countryMap[countryName] || countryName;
  };

  // Convert passport MRZ date format YYMMDD to YYYY-MM-DD (HTML date input format)
  const convertPassportDate = (passportDate: string): string => {
    if (!passportDate || passportDate.length !== 6) return '';
    
    const yy = passportDate.substring(0, 2);
    const mm = passportDate.substring(2, 4);
    const dd = passportDate.substring(4, 6);
    
    // Convert YY to YYYY according to MRZ specification:
    // For dates of birth: 00-30 = 20xx, 31-99 = 19xx
    // For expiry dates: typically all future dates, so 00-99 = 20xx
    // We'll use context-aware interpretation
    let year: string;
    const yyNum = parseInt(yy);
    
    // For years 00-99, determine century based on reasonable date ranges
    if (yyNum >= 0 && yyNum <= 99) {
      // If it's likely a birth date (before current year), use 19xx for 31-99, 20xx for 00-30
      // If it's likely an expiry date (future), use 20xx
      const currentYear = new Date().getFullYear();
      const currentYY = currentYear % 100;
      
      if (yyNum > currentYY + 10) {
        // Likely a birth date from previous century
        year = `19${yy}`;
      } else {
        // Likely current century (birth date for young people or expiry date)
        year = `20${yy}`;
      }
    } else {
      year = `20${yy}`;
    }
    
    return `${year}-${mm}-${dd}`;
  };

  const populateFormFields = (data: any, uploadType: string) => {
    console.log('Populating form fields:', { data, uploadType, selectedDocumentType });
    
    if (selectedDocumentType === 'driving-licence') {
      if (uploadType === 'front' && data.licenceNumber) {
        const input = document.getElementById('licenceNumber') as HTMLInputElement;
        if (input) input.value = data.licenceNumber;
      }
      if (data.expiryDate) {
        const input = document.getElementById('licenceExpiryDate') as HTMLInputElement;
        if (input) input.value = data.expiryDate;
      }
      if (data.issueDate) {
        const input = document.getElementById('licenceIssueDate') as HTMLInputElement;
        if (input) input.value = data.issueDate;
      }
    } else if (selectedDocumentType === 'passport') {
      // Extract passport number from uploadType (e.g., "passport_1" -> "1")
      const passportNum = uploadType.includes('passport_') ? uploadType.split('_')[1] : '1';
      console.log('Processing passport data:', { passportNum, data });
      
      if (data.passportNumber) {
        const input = document.getElementById(`passportNumber_${passportNum}`) as HTMLInputElement;
        console.log('Found passport number input:', input);
        if (input) {
          input.value = data.passportNumber;
          console.log('Set passport number:', data.passportNumber);
        }
      }
      if (data.issuingCountry) {
        const select = document.getElementById(`passportCountry_${passportNum}`) as HTMLSelectElement;
        if (select) {
          // Map full country name back to country code for form selection
          const countryCode = getCountryCodeFromName(data.issuingCountry);
          select.value = countryCode;
          console.log('Set issuing country:', data.issuingCountry, '->', countryCode);
        }
      }
      if (data.expiryDate) {
        const input = document.getElementById(`passportExpiryDate_${passportNum}`) as HTMLInputElement;
        if (input) {
          const formattedDate = convertPassportDate(data.expiryDate);
          input.value = formattedDate;
          console.log('Set expiry date:', data.expiryDate, '->', formattedDate);
        }
      }
      // Issue date is not available in MRZ - only populate if explicitly provided
      if (data.issueDate) {
        const input = document.getElementById(`passportIssueDate_${passportNum}`) as HTMLInputElement;
        if (input) {
          const formattedDate = convertPassportDate(data.issueDate);
          input.value = formattedDate;
          console.log('Set issue date:', data.issueDate, '->', formattedDate);
        }
      }
      if (data.givenNames) {
        const input = document.getElementById(`passportGivenNames_${passportNum}`) as HTMLInputElement;
        if (input) {
          // Clean up given names - remove extra padding and filler characters
          const cleanedNames = data.givenNames
            .replace(/\s+/g, ' ')  // Replace multiple spaces with single space
            .replace(/[K<]+/g, '') // Remove filler characters K and <
            .trim();               // Remove leading/trailing spaces
          input.value = cleanedNames;
          console.log('Set given names:', data.givenNames, '->', cleanedNames);
        }
      }
      if (data.surname) {
        const input = document.getElementById(`passportSurname_${passportNum}`) as HTMLInputElement;
        if (input) input.value = data.surname;
      }
      if (data.dateOfBirth) {
        const input = document.getElementById(`passportDateOfBirth_${passportNum}`) as HTMLInputElement;
        if (input) {
          const formattedDate = convertPassportDate(data.dateOfBirth);
          input.value = formattedDate;
          console.log('Set date of birth:', data.dateOfBirth, '->', formattedDate);
        }
      }
      if (data.gender) {
        const select = document.getElementById(`passportGender_${passportNum}`) as HTMLSelectElement;
        if (select) select.value = data.gender;
      }
    }
  };

  const validateWithSHACL = async (documentType: string, data: any) => {
    try {
      const response = await fetch('/api/validate-document', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          documentType,
          data
        })
      });
      
      if (response.ok) {
        const validation = await response.json();
        return validation;
      }
    } catch (error) {
      console.error('SHACL validation failed:', error);
    }
    return null;
  };

  // AI Form Analysis Functions
  const analyzeForm = async (file?: File, htmlContent?: string, url?: string) => {
    setIsAnalyzing(true);
    try {
      const formData = new FormData();
      formData.append('analysisType', analysisType);
      
      if (analysisType === 'document' && file) {
        formData.append('document', file);
      } else if (analysisType === 'html' && htmlContent) {
        formData.append('htmlContent', htmlContent);
      } else if (analysisType === 'url' && url) {
        formData.append('url', url);
      }

      const response = await fetch('/api/analyze-form', {
        method: 'POST',
        body: formData,
      });

      if (response.ok) {
        const result = await response.json();
        setAnalysisResults(result);
        
        // Automatically generate field mappings
        await generateFieldMappings(result);
      } else {
        console.error('Form analysis failed');
      }
    } catch (error) {
      console.error('Form analysis error:', error);
    } finally {
      setIsAnalyzing(false);
    }
  };

  const generateFieldMappings = async (formAnalysis: any) => {
    try {
      const response = await fetch('/api/map-fields', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(formAnalysis),
      });

      if (response.ok) {
        const mappingResult = await response.json();
        setFieldMappings(mappingResult.mappings);
        console.log('Generated field mappings:', mappingResult);
      }
    } catch (error) {
      console.error('Field mapping error:', error);
    }
  };

  const applyFieldMappings = async () => {
    if (!fieldMappings.length) return;

    try {
      // Transform the current session data using SHACL rules
      const transformRequest = {
        sourceData: session,
        targetShape: 'autoins:PersonShape',
        mappings: fieldMappings,
        options: {
          strictValidation: false,
          language: session.language,
          preserveOriginal: true
        }
      };

      const response = await fetch('/api/shacl-transform', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(transformRequest),
      });

      if (response.ok) {
        const result = await response.json();
        if (result.success) {
          // Apply transformed data to session
          const updatedSession = { ...session, ...result.transformedData };
          setSession(updatedSession);
          
          // Show success notification
          setNotification({
            type: 'success',
            message: translations[session.language].fieldMappingApplied || 'Field mappings applied successfully'
          });
        } else {
          console.error('Transformation errors:', result.validationErrors);
        }
      }
    } catch (error) {
      console.error('Field mapping application error:', error);
    }
  };

  // Web Spider Functions

  const mapExtractedDataToSession = async (extractedData: any) => {
    // Intelligent mapping of extracted data to session fields
    const mappings: any = {};
    
    // Common field mappings
    const fieldMappings: {[key: string]: string[]} = {
      'firstName': ['first_name', 'fname', 'given_name', 'forename'],
      'lastName': ['last_name', 'lname', 'surname', 'family_name'],
      'email': ['email', 'email_address', 'e_mail'],
      'phone': ['phone', 'telephone', 'mobile', 'phone_number'],
      'address': ['address', 'street_address', 'addr1', 'address_line_1'],
      'postcode': ['postcode', 'postal_code', 'zip', 'zip_code'],
      'dateOfBirth': ['dob', 'date_of_birth', 'birth_date', 'birthdate']
    };

    // Map extracted data to session fields
    Object.keys(extractedData).forEach(key => {
      const lowerKey = key.toLowerCase();
      Object.keys(fieldMappings).forEach(sessionField => {
        if (fieldMappings[sessionField].some(pattern => lowerKey.includes(pattern))) {
          mappings[sessionField] = extractedData[key];
        }
      });
    });

    // Update session with mapped data
    if (Object.keys(mappings).length > 0) {
      const updatedDrivers = [...session.drivers];
      if (updatedDrivers.length > 0) {
        updatedDrivers[0] = { ...updatedDrivers[0], ...mappings };
        setSession({ ...session, drivers: updatedDrivers });
      }
    }
  };

  const createSpiderTask = (type: string, config: any) => {
    const task = {
      id: `task_${Date.now()}`,
      type,
      status: 'pending',
      createdAt: new Date(),
      ...config
    };
    
    setSpiderTasks(prev => [...prev, task]);
    return task;
  };

  // Bitwarden Functions
  const unlockBitwarden = async (password: string) => {
    try {
      const response = await fetch('/api/bitwarden/unlock', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ password }),
      });

      if (response.ok) {
        const result = await response.json();
        setBitwardenUnlocked(result.isUnlocked);
        
        // Load available credentials
        await loadAvailableCredentials();
        
        setNotification({
          type: 'success',
          message: translations[session.language].bitwardenUnlocked || 'Bitwarden vault unlocked'
        });
      } else {
        console.error('Failed to unlock Bitwarden');
      }
    } catch (error) {
      console.error('Bitwarden unlock error:', error);
    }
  };

  // Smart Bitwarden Connection Handler
  const connectToBitwarden = async () => {
    try {
      // First check Bitwarden status
      const statusResponse = await fetch('/api/bitwarden/status');
      
      if (!statusResponse.ok) {
        throw new Error(`HTTP ${statusResponse.status}: ${statusResponse.statusText}`);
      }
      
      const responseText = await statusResponse.text();
      console.log('Raw Bitwarden status response:', responseText);
      
      let status;
      try {
        status = JSON.parse(responseText);
      } catch (parseError) {
        throw new Error(`Invalid JSON response: ${responseText.substring(0, 100)}...`);
      }
      
      if (!status.available) {
        alert('❌ Bitwarden CLI not available: ' + status.error);
        return;
      }
      
      const bitwardenStatus = status.status?.status;
      
      if (bitwardenStatus === 'unauthenticated') {
        // Ask user for login method
        const loginMethod = window.confirm('Choose login method:\n\nOK = Email/Password\nCancel = API Key (for 2FA accounts)');
        
        if (loginMethod) {
          // Email/Password login
          const email = prompt('🔐 Enter your Bitwarden email:');
          const password = prompt('🔑 Enter your Bitwarden password:');
          
          if (!email || !password) {
            alert('❌ Email and password required');
            return;
          }
          
          const loginResponse = await fetch('/api/bitwarden/login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ email, password })
          });
          
          const loginResult = await loginResponse.json();
          if (!loginResult.success) {
            let errorMsg = loginResult.error || 'Unknown error';
            
            // Provide helpful error messages and suggestions
            if (errorMsg.includes('invalid credentials') || errorMsg.includes('Username or password')) {
              errorMsg = 'Invalid email or password. Please check your Bitwarden credentials.';
            } else if (errorMsg.includes('two-factor') || errorMsg.includes('2FA')) {
              errorMsg = 'Two-factor authentication detected. Please use API Key login instead.';
            } else if (errorMsg.includes('server configuration')) {
              errorMsg = 'Server configuration error. Please check your Bitwarden server URL.';
            } else if (errorMsg.includes('network')) {
              errorMsg = 'Network error. Please check your internet connection.';
            }
            
            alert('❌ Login failed: ' + errorMsg);
            return;
          }
          
          setBitwardenUnlocked(true);
          alert('✅ Successfully logged in to Bitwarden!');
          
        } else {
          // API Key login
          const clientId = prompt('🔑 Enter your Bitwarden Client ID:');
          const clientSecret = prompt('🔐 Enter your Bitwarden Client Secret:');
          const masterPassword = prompt('🔓 Enter your Master Password:');
          
          if (!clientId || !clientSecret || !masterPassword) {
            alert('❌ All API key fields are required');
            return;
          }
          
          const apiLoginResponse = await fetch('/api/bitwarden/login-apikey', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ 
              clientId, 
              clientSecret, 
              password: masterPassword 
            })
          });
          
          const apiLoginResult = await apiLoginResponse.json();
          if (!apiLoginResult.success) {
            alert('❌ API Key login failed: ' + (apiLoginResult.error || 'Unknown error'));
            return;
          }
          
          setBitwardenUnlocked(true);
          alert('✅ Successfully logged in to Bitwarden with API key!');
        }
        
      } else if (bitwardenStatus === 'locked') {
        // Need to unlock
        const masterPassword = prompt('🔓 Enter your Bitwarden master password:');
        if (!masterPassword) {
          alert('❌ Master password required');
          return;
        }
        
        await unlockBitwarden(masterPassword);
        
      } else if (bitwardenStatus === 'unlocked') {
        setBitwardenUnlocked(true);
        alert('✅ Bitwarden is already unlocked!');
      }
      
      if (bitwardenUnlocked || bitwardenStatus === 'unlocked') {
        await loadAvailableCredentials();
      }
      
    } catch (error) {
      console.error('Bitwarden connection failed:', error);
      alert('❌ Bitwarden connection failed: ' + error);
    }
  };

  const loadAvailableCredentials = async () => {
    try {
      const response = await fetch('/api/bitwarden/list-credentials?category=site_login');
      if (response.ok) {
        const result = await response.json();
        setAvailableCredentials(result.credentials || []);
      }
    } catch (error) {
      console.error('Failed to load credentials:', error);
    }
  };

  const setupBitwardenTemplates = async () => {
    try {
      const response = await fetch('/api/bitwarden/setup-templates', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          setupSites: true,
          setupOpenBanking: true
        }),
      });

      if (response.ok) {
        await loadAvailableCredentials();
        setNotification({
          type: 'success',
          message: translations[session.language].templatesCreated || 'Credential templates created'
        });
      }
    } catch (error) {
      console.error('Failed to setup templates:', error);
    }
  };

  // Stealth Browser Functions
  const createStealthSession = async (url: string) => {
    setIsBrowserRunning(true);
    try {
      const response = await fetch('/api/stealth-browser', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          action: 'create',
          url: url,
          options: {
            randomizeUserAgent: true,
            randomizeViewport: true,
            humanizeTyping: true,
            humanizeClicks: true,
            randomDelays: true,
            stealthPlugins: true,
            webrtcBlock: true,
            canvasFingerprint: true,
            audioFingerprint: true,
            minDelay: 500,
            maxDelay: 2000
          }
        }),
      });

      if (response.ok) {
        const result = await response.json();
        setBrowserSession(result);
        setBrowserScreenshot(result.screenshot);
        
        setNotification({
          type: 'success',
          message: translations[session.language].browserSessionCreated || 'Stealth browser session created'
        });
      } else {
        console.error('Failed to create stealth browser session');
      }
    } catch (error) {
      console.error('Stealth browser error:', error);
    } finally {
      setIsBrowserRunning(false);
    }
  };

  const loginToSite = async (siteName: string) => {
    if (!browserSession) {
      await createStealthSession('');
    }

    setIsBrowserRunning(true);
    try {
      const response = await fetch('/api/stealth-browser', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          action: 'login',
          sessionId: browserSession?.sessionId || 'new',
          loginSite: siteName
        }),
      });

      if (response.ok) {
        const result = await response.json();
        setBrowserSession(result);
        setBrowserScreenshot(result.screenshot);
        
        if (result.isLoggedIn) {
          setNotification({
            type: 'success',
            message: translations[session.language].loginSuccessful || `Successfully logged in to ${siteName}`
          });
        } else {
          setNotification({
            type: 'info',
            message: translations[session.language].loginFailed || `Login to ${siteName} may have failed`
          });
        }
      } else {
        console.error('Failed to login to site');
      }
    } catch (error) {
      console.error('Site login error:', error);
    } finally {
      setIsBrowserRunning(false);
    }
  };

  // Money Supermarket Infiltration Function
  const infiltrateMoneySupermarket = async () => {
    setInfiltrationInProgress(true);
    setInfiltrationLogs([]);
    setExtractedQuotes([]);
    
    const logs: string[] = [];
    const addLog = (message: string) => {
      logs.push(`${new Date().toLocaleTimeString()}: ${message}`);
      setInfiltrationLogs([...logs]);
    };

    try {
      addLog('🥷 Initiating Money Supermarket infiltration...');
      addLog('🛡️ Activating maximum stealth protection...');
      
      const infiltrationData = {
        vehicleReg: session.vehicle?.registration || 'AB12 CDE',
        postCode: session.drivers?.[0]?.postcode || 'SW1A 1AA',
        dateOfBirth: session.drivers?.[0]?.dateOfBirth || '01/01/1990',
        licenceType: session.drivers?.[0]?.licenceType || 'full'
      };

      addLog('📝 Preparing infiltration data...');
      addLog(`🎯 Target: Money Supermarket (${Object.keys(infiltrationData).length} data points)`);
      
      const response = await fetch('/api/infiltrate/moneysupermarket', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(infiltrationData)
      });

      if (response.ok) {
        const result = await response.json();
        setInfiltrationResult(result);
        
        if (result.success) {
          addLog(`🏆 INFILTRATION SUCCESSFUL!`);
          addLog(`💰 Extracted ${result.quotes.length} insurance quotes`);
          addLog(`📸 Captured ${result.screenshots.length} evidence screenshots`);
          addLog(`⏱️ Mission duration: ${result.duration}`);
          
          setExtractedQuotes(result.quotes);
          
          // Display quotes
          result.quotes.forEach((quote: any, index: number) => {
            addLog(`📋 Quote ${index + 1}: ${quote.provider} - ${quote.price}`);
          });
          
        } else {
          addLog('❌ Infiltration failed');
          result.errors.forEach((error: string) => {
            addLog(`🚨 Error: ${error}`);
          });
        }
      } else {
        addLog('❌ Failed to connect to infiltration endpoint');
      }
      
    } catch (error) {
      console.error('Infiltration error:', error);
      addLog(`🚨 Critical error: ${error}`);
    } finally {
      setInfiltrationInProgress(false);
      addLog('🏁 Infiltration mission complete');
    }
  };
  const [session, setSession] = useState<QuoteSession>({
    id: '1',
    language: 'en',
    drivers: [
      {
        id: '1',
        classification: 'MAIN',
        firstName: '',
        lastName: '',
        dateOfBirth: '',
        email: '',
        phone: '',
        address: '',
        postcode: '',
        licenceType: '',
        licenceNumber: '',
        yearsHeld: 0,
        pointsLost: 0,
        licenceIssueDate: '',
        licenceExpiryDate: '',
        licenceValidUntil: '',
        relationship: '',
        sameAddress: true,
        // Disabilities and Restrictions
        hasDisability: false,
        disabilityTypes: [],
        requiresAdaptations: false,
        adaptationTypes: [],
        automaticOnly: false,
        // Licence Classes
        licenceClassA: false,
        licenceClassB: true,
        licenceClassC: false,
        licenceClassD: false,
        licenceClassBE: false,
        licenceClassCE: false,
        licenceClassDE: false,
        // Driver Communication
        driverEmail: '',
        emailSent: false,
        emailSentDate: ''
      }
    ],
    vehicle: {
      registration: '',
      make: '',
      model: '',
      year: new Date().getFullYear(),
      mileage: 0,
      value: 0,
      daytimeLocation: '',
      overnightLocation: '',
      hasModifications: false,
      modifications: []
    },
    policy: {
      coverType: '',
      startDate: '',
      voluntaryExcess: 0,
      ncdYears: 0,
      protectNCD: false
    },
    claims: {
      hasClaims: false,
      claims: [],
      hasConvictions: false,
      convictions: [],
      hasAccidents: false,
      accidents: []
    },
    payment: {
      frequency: '',
      method: ''
    },
    extras: {
      breakdownCover: '',
      legalExpenses: false,
      courtesyCar: false
    },
    marketing: {
      emailMarketing: false,
      smsMarketing: false,
      postMarketing: false
    },
    progress: {},
    formData: {},
    createdAt: new Date().toISOString(),
    lastAccessed: new Date().toISOString()
  });

  // Function to get menu structure with translations
  const getMenuStructure = () => [
    {
      id: 'car-insurance',
      title: translations[session.language].carInsurance,
      icon: Shield,
      categories: [
        { id: 'drivers', title: translations[session.language].driverDetails, icon: User, order: 1, description: translations[session.language].driverDetailsDesc },
        { id: 'vehicle', title: translations[session.language].vehicleDetails, icon: Car, order: 2, description: translations[session.language].vehicleDetailsDesc },
        { id: 'policy', title: translations[session.language].policyDetails, icon: FileText, order: 3, description: translations[session.language].policyDetailsDesc },
        { id: 'claims', title: translations[session.language].claimsHistory, icon: Shield, order: 4, description: translations[session.language].claimsHistoryDesc },
        { id: 'payment', title: translations[session.language].paymentExtras, icon: CreditCard, order: 5, description: translations[session.language].paymentExtrasDesc }
      ]
    },
    {
      id: 'personal-documents',
      title: translations[session.language].personalDocuments,
      icon: Upload,
      isModal: true,
      categories: []
    },
    {
      id: 'settings',
      title: translations[session.language].settings,
      icon: Settings,
      categories: [
        { id: 'marketing', title: translations[session.language].marketingPreferences, icon: Settings, order: 1, description: translations[session.language].marketingPreferencesDesc }
      ]
    }
  ];

  const menuStructure = getMenuStructure();
  // Flatten categories for backward compatibility
  const categories = menuStructure.flatMap(section => section.categories);

  const [isHelpOpen, setIsHelpOpen] = useState(false);
  const [showDocumentUpload, setShowDocumentUpload] = useState(false);
  const [uploadedFiles, setUploadedFiles] = useState<File[]>([]);
  const [isDragOver, setIsDragOver] = useState(false);
  const [chatMessages, setChatMessages] = useState<Array<{type: 'user' | 'bot', content: string, timestamp: Date}>>([]);
  const [showChatbot, setShowChatbot] = useState(false);
  const [notification, setNotification] = useState<{type: 'success' | 'info', message: string} | null>(null);

  const updateDriver = (index: number, field: keyof Driver, value: any) => {
    const updatedDrivers = [...session.drivers];
    updatedDrivers[index] = { ...updatedDrivers[index], [field]: value };
    setSession({ ...session, drivers: updatedDrivers });
  };

  const addDriver = () => {
    const newDriver: Driver = {
      id: (session.drivers.length + 1).toString(),
      classification: 'NAMED',
      firstName: '',
      lastName: '',
      dateOfBirth: '',
      email: '',
      phone: '',
      address: '',
      postcode: '',
      licenceType: '',
      licenceNumber: '',
      yearsHeld: 0,
      pointsLost: 0,
      licenceIssueDate: '',
      licenceExpiryDate: '',
      licenceValidUntil: '',
      relationship: '',
      sameAddress: false,
      // Disabilities and Restrictions
      hasDisability: false,
      disabilityTypes: [],
      requiresAdaptations: false,
      adaptationTypes: [],
      automaticOnly: false,
      // Licence Classes
      licenceClassA: false,
      licenceClassB: true,
      licenceClassC: false,
      licenceClassD: false,
      licenceClassBE: false,
      licenceClassCE: false,
      licenceClassDE: false,
      // Driver Communication
      driverEmail: '',
      emailSent: false,
      emailSentDate: ''
    };
    setSession({ ...session, drivers: [...session.drivers, newDriver] });
  };

  const removeDriver = (index: number) => {
    if (session.drivers.length > 1) {
      const updatedDrivers = session.drivers.filter((_, i) => i !== index);
      setSession({ ...session, drivers: updatedDrivers });
    }
  };

  const updateVehicle = (field: keyof Vehicle, value: any) => {
    setSession({ ...session, vehicle: { ...session.vehicle, [field]: value } });
  };

  const updatePolicy = (field: keyof Policy, value: any) => {
    setSession({ ...session, policy: { ...session.policy, [field]: value } });
  };

  const updateClaims = (field: keyof ClaimsHistory, value: any) => {
    setSession({ ...session, claims: { ...session.claims, [field]: value } });
  };

  const addClaim = () => {
    const newClaim: Claim = {
      id: Date.now().toString(),
      date: '',
      type: '',
      amount: 0,
      status: 'Open',
      faultStatus: '',
      policeReported: false,
      crimeReference: '',
      description: ''
    };
    setSession({
      ...session,
      claims: { ...session.claims, claims: [...session.claims.claims, newClaim] }
    });
  };

  const updateClaim = (index: number, field: keyof Claim, value: any) => {
    const updatedClaims = [...session.claims.claims];
    updatedClaims[index] = { ...updatedClaims[index], [field]: value };
    setSession({
      ...session,
      claims: { ...session.claims, claims: updatedClaims }
    });
  };

  const removeClaim = (index: number) => {
    const updatedClaims = session.claims.claims.filter((_, i) => i !== index);
    setSession({
      ...session,
      claims: { ...session.claims, claims: updatedClaims }
    });
  };

  const addConviction = () => {
    const newConviction: Conviction = {
      id: Date.now().toString(),
      date: '',
      offenceCode: '',
      description: '',
      penaltyPoints: 0,
      fine: 0,
      disqualification: 0,
      status: 'Active'
    };
    setSession({
      ...session,
      claims: { ...session.claims, convictions: [...session.claims.convictions, newConviction] }
    });
  };

  const updateConviction = (index: number, field: keyof Conviction, value: any) => {
    const updatedConvictions = [...session.claims.convictions];
    updatedConvictions[index] = { ...updatedConvictions[index], [field]: value };
    setSession({
      ...session,
      claims: { ...session.claims, convictions: updatedConvictions }
    });
  };

  const removeConviction = (index: number) => {
    const updatedConvictions = session.claims.convictions.filter((_, i) => i !== index);
    setSession({
      ...session,
      claims: { ...session.claims, convictions: updatedConvictions }
    });
  };

  const addAccident = () => {
    const newAccident: Accident = {
      id: Date.now().toString(),
      date: '',
      type: '',
      severity: '',
      claimMade: false,
      description: ''
    };
    setSession({
      ...session,
      claims: { ...session.claims, accidents: [...session.claims.accidents, newAccident] }
    });
  };

  const updateAccident = (index: number, field: keyof Accident, value: any) => {
    const updatedAccidents = [...session.claims.accidents];
    updatedAccidents[index] = { ...updatedAccidents[index], [field]: value };
    setSession({
      ...session,
      claims: { ...session.claims, accidents: updatedAccidents }
    });
  };

  const removeAccident = (index: number) => {
    const updatedAccidents = session.claims.accidents.filter((_, i) => i !== index);
    setSession({
      ...session,
      claims: { ...session.claims, accidents: updatedAccidents }
    });
  };

  // Email functionality for additional drivers
  const sendEmailToDriver = async (driverIndex: number) => {
    const driver = session.drivers[driverIndex];
    
    if (!driver.driverEmail) {
      setNotification({ type: 'info', message: 'Please enter the driver\'s email address first.' });
      return;
    }

    try {
      // In a real application, this would call your backend API
      // For now, we'll simulate the email sending
      const emailData = {
        to: driver.driverEmail,
        subject: 'Insurance Quote Form - Additional Driver Information Required',
        body: generateDriverEmailContent(driver),
        driverId: driver.id
      };

      // Simulate API call
      console.log('Sending email:', emailData);
      
      // Update driver with email sent status
      const updatedDrivers = [...session.drivers];
      updatedDrivers[driverIndex] = {
        ...updatedDrivers[driverIndex],
        emailSent: true,
        emailSentDate: new Date().toISOString()
      };
      
      setSession({ ...session, drivers: updatedDrivers });
      
      setNotification({ 
        type: 'success', 
        message: `Email sent to ${driver.firstName} ${driver.lastName} at ${driver.driverEmail}` 
      });
      
      addChatMessage('bot', `📧 **Email Sent:** Form has been emailed to ${driver.firstName} ${driver.lastName} (${driver.driverEmail})\n\n📋 **Email includes:**\n• Personal insurance form link\n• Instructions for completing their section\n• Contact information for support\n\n⏰ **Next Steps:**\n• Driver will receive email with form link\n• They can complete their details independently\n• You'll be notified when they submit`);
      
    } catch (error) {
      console.error('Error sending email:', error);
      setNotification({ 
        type: 'info', 
        message: 'Failed to send email. Please try again.' 
      });
    }
  };

  const generateDriverEmailContent = (driver: Driver) => {
    const formLink = `${window.location.origin}/driver-form/${driver.id}`;
    
    return `
Dear ${driver.firstName} ${driver.lastName},

You have been added as an additional driver to an insurance quote application.

Please complete your driver information by clicking the link below:
${formLink}

Required Information:
• Personal details (name, date of birth, address)
• Driving licence information
• Medical conditions or disabilities
• Claims and convictions history
• Vehicle usage details

If you have any questions, please contact the main policyholder.

Best regards,
Insurance Quote System
    `;
  };

  // Calculate years held based on issue date
  const calculateYearsHeld = (issueDate: string): number => {
    if (!issueDate) return 0;
    
    const issue = new Date(issueDate);
    const now = new Date();
    const diffTime = Math.abs(now.getTime() - issue.getTime());
    const diffYears = Math.floor(diffTime / (1000 * 60 * 60 * 24 * 365.25));
    
    return diffYears;
  };

  // Update years held when issue date changes
  const updateLicenceIssueDate = (driverIndex: number, issueDate: string) => {
    const yearsHeld = calculateYearsHeld(issueDate);
    updateDriver(driverIndex, 'licenceIssueDate', issueDate);
    updateDriver(driverIndex, 'yearsHeld', yearsHeld);
  };

  // Check licence expiry status
  const checkLicenceExpiry = (expiryDate: string): { status: 'valid' | 'expired' | 'expiring_soon', daysLeft: number } => {
    if (!expiryDate) return { status: 'valid', daysLeft: 999 };
    
    const expiry = new Date(expiryDate);
    const now = new Date();
    const diffTime = expiry.getTime() - now.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    
    if (diffDays < 0) {
      return { status: 'expired', daysLeft: Math.abs(diffDays) };
    } else if (diffDays <= 30) {
      return { status: 'expiring_soon', daysLeft: diffDays };
    } else {
      return { status: 'valid', daysLeft: diffDays };
    }
  };

  // Get expiry status badge
  const getExpiryStatusBadge = (expiryDate: string) => {
    const expiryStatus = checkLicenceExpiry(expiryDate);
    
    switch (expiryStatus.status) {
      case 'expired':
        return <Badge color="failure" className="ml-2 text-xs">Expired {expiryStatus.daysLeft} days ago</Badge>;
      case 'expiring_soon':
        return <Badge color="warning" className="ml-2 text-xs">Expires in {expiryStatus.daysLeft} days</Badge>;
      default:
        return <Badge color="success" className="ml-2 text-xs">Valid</Badge>;
    }
  };

  const updatePayment = (field: keyof Payment, value: any) => {
    setSession({ ...session, payment: { ...session.payment, [field]: value } });
  };

  const updateExtras = (field: keyof Extras, value: any) => {
    setSession({ ...session, extras: { ...session.extras, [field]: value } });
  };

  const updateMarketing = (field: keyof Marketing, value: any) => {
    setSession({ ...session, marketing: { ...session.marketing, [field]: value } });
  };

  const getStatusBadge = (isRequired: boolean, hasValue: boolean) => {
    if (isRequired && !hasValue) {
      return <Badge color="failure" className="status-badge status-missing">Missing</Badge>;
    } else if (isRequired) {
      return <Badge color="success" className="status-badge status-required">Required</Badge>;
    } else {
      return <Badge color="warning" className="status-badge status-optional">Optional</Badge>;
    }
  };

  const handleFileSelect = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files) {
      const fileArray = Array.from(files);
      setUploadedFiles(prev => [...prev, ...fileArray]);
      processFiles(fileArray);
    }
  };

  const handleDragOver = (event: React.DragEvent) => {
    event.preventDefault();
    setIsDragOver(true);
  };

  const handleDragLeave = (event: React.DragEvent) => {
    event.preventDefault();
    setIsDragOver(false);
  };

  const handleDrop = (event: React.DragEvent) => {
    event.preventDefault();
    setIsDragOver(false);
    
    const files = event.dataTransfer.files;
    if (files) {
      const fileArray = Array.from(files);
      setUploadedFiles(prev => [...prev, ...fileArray]);
      processFiles(fileArray);
    }
  };

  const processFiles = (files: File[]) => {
    files.forEach(file => {
      console.log(`Processing file: ${file.name} (${file.type})`);
      
      // Show processing start message
      addChatMessage('bot', `🔍 **Processing Document:** ${file.name}\n\n⏳ Analyzing document type and extracting fields...`);
      
      // Simulate OCR and document type recognition
      setTimeout(() => {
        const documentType = recognizeDocumentType(file);
        const extractedData = simulateDocumentExtraction(file, documentType);
        
        if (documentType === 'passport') {
          processPassportData(extractedData);
        } else if (documentType === 'driving_licence') {
          processDrivingLicenceData(extractedData);
        } else {
          addChatMessage('bot', `✅ **Document Processed:** ${file.name}\n\n📄 **Document Type:** ${documentType}\n\n💡 This document type is not yet fully supported for automatic field extraction.`);
        }
      }, 1000);
    });
  };

  const recognizeDocumentType = (file: File): string => {
    // Simulate document type recognition based on filename and content
    const fileName = file.name.toLowerCase();
    if (fileName.includes('passport') || fileName.includes('pass')) {
      return 'passport';
    } else if (fileName.includes('licence') || fileName.includes('license') || fileName.includes('driving')) {
      return 'driving_licence';
    } else if (fileName.includes('insurance')) {
      return 'insurance_policy';
    } else if (fileName.includes('vehicle') || fileName.includes('v5')) {
      return 'vehicle_registration';
    } else {
      return 'unknown';
    }
  };

  const simulateDocumentExtraction = (file: File, documentType: string) => {
    // Simulate OCR extraction based on document type
    if (documentType === 'passport') {
      return {
        PersonalName: 'John Michael Smith',
        GivenNames: 'John Michael',
        Surname: 'Smith',
        DateOfBirth: '1985-03-15',
        PlaceOfBirth: 'London, United Kingdom',
        Gender: 'M',
        Nationality: 'British',
        PassportNumber: 'GB123456789',
        DateOfIssue: '2020-01-15',
        DateOfExpiry: '2030-01-15',
        Authority: 'HM Passport Office',
        Address: '123 Main Street, London, SW1A 1AA',
        DocumentType: 'Passport',
        DocumentCountry: 'United Kingdom'
      };
    } else if (documentType === 'driving_licence') {
      return {
        PersonalName: 'Sarah Elizabeth Johnson',
        GivenNames: 'Sarah Elizabeth',
        Surname: 'Johnson',
        DateOfBirth: '1990-07-22',
        PlaceOfBirth: 'Manchester, United Kingdom',
        LicenceNumber: 'JOHN123456SJ9AB',
        DateOfIssue: '2015-03-10',
        DateOfExpiry: '2025-03-10',
        Authority: 'DVLA',
        CategoryA: '2020-01-15',
        CategoryB: '2015-03-10',
        CategoryC: '',
        CategoryD: '',
        CategoryBE: '2018-06-20',
        CategoryCE: '',
        CategoryDE: '',
        Restrictions: '01 - Eyesight correction',
        Endorsements: '',
        YearsHeld: 8,
        DocumentType: 'Driving Licence',
        DocumentCountry: 'United Kingdom'
      };
    }
    return {};
  };

  const normalizeName = (name: string): string => {
    // Remove extra spaces, convert to lowercase, and normalize
    return name.toLowerCase().replace(/\s+/g, ' ').trim();
  };

  const namesMatch = (name1: string, name2: string): boolean => {
    const normalized1 = normalizeName(name1);
    const normalized2 = normalizeName(name2);
    return normalized1 === normalized2;
  };

  const processPassportData = (passportData: any) => {
    // Create detailed extraction report
    const extractedFields = Object.entries(passportData)
      .filter(([key, value]) => value && value !== '')
      .map(([key, value]) => `• ${key}: ${value}`)
      .join('\n');

    addChatMessage('bot', `🔍 **Document Processing Complete**\n\n📄 **Extracted Fields:**\n${extractedFields}`);

    // Enhanced matching logic: try name + date of birth first, then name only
    let existingDriverIndex = -1;
    let matchType = '';
    
    // First, try to match by name AND date of birth (most precise)
    if (passportData.DateOfBirth) {
      existingDriverIndex = session.drivers.findIndex(driver => {
        const driverFullName = `${driver.firstName} ${driver.lastName}`;
        const passportFullName = `${passportData.GivenNames} ${passportData.Surname}`;
        const nameMatch = namesMatch(driverFullName, passportFullName);
        const dobMatch = driver.dateOfBirth === passportData.DateOfBirth;
        return nameMatch && dobMatch;
      });
      
      if (existingDriverIndex >= 0) {
        matchType = 'name_and_dob';
      }
    }
    
    // If no match by name + DOB, try name only
    if (existingDriverIndex === -1) {
      existingDriverIndex = session.drivers.findIndex(driver => {
        const driverFullName = `${driver.firstName} ${driver.lastName}`;
        const passportFullName = `${passportData.GivenNames} ${passportData.Surname}`;
        return namesMatch(driverFullName, passportFullName);
      });
      
      if (existingDriverIndex >= 0) {
        matchType = 'name_only';
      }
    }

    if (existingDriverIndex >= 0) {
      // Update existing driver with missing fields
      const updatedDriver = { ...session.drivers[existingDriverIndex] };
      const updatedFields: string[] = [];
      
      if (!updatedDriver.dateOfBirth && passportData.DateOfBirth) {
        updatedDriver.dateOfBirth = passportData.DateOfBirth;
        updatedFields.push('Date of Birth');
      }
      if (!updatedDriver.firstName && passportData.GivenNames) {
        updatedDriver.firstName = passportData.GivenNames;
        updatedFields.push('First Name');
      }
      if (!updatedDriver.lastName && passportData.Surname) {
        updatedDriver.lastName = passportData.Surname;
        updatedFields.push('Last Name');
      }
      
      // Add passport-specific fields
      updatedDriver.passportNumber = passportData.PassportNumber;
      updatedDriver.passportIssueDate = passportData.DateOfIssue;
      updatedDriver.passportExpiryDate = passportData.DateOfExpiry;
      updatedDriver.passportAuthority = passportData.Authority;
      updatedDriver.placeOfBirth = passportData.PlaceOfBirth;
      updatedDriver.gender = passportData.Gender;
      updatedDriver.nationality = passportData.Nationality;

      setSession(prev => {
        const updatedDrivers = prev.drivers.map((driver, index) => 
          index === existingDriverIndex ? updatedDriver : driver
        );
        console.log('Updated drivers:', updatedDrivers);
        return {
          ...prev,
          drivers: updatedDrivers
        };
      });

      const matchTypeText = matchType === 'name_and_dob' 
        ? '✅ **Exact Match Found:** Name + Date of Birth'
        : '⚠️ **Name Match Found:** Date of Birth differs or missing';
        
      const updateMessage = updatedFields.length > 0 
        ? `${matchTypeText}\n\n👤 **Driver Updated:** "${passportData.GivenNames} ${passportData.Surname}"\n\n📝 **Fields Updated:**\n${updatedFields.map(field => `• ${field}`).join('\n')}\n\n📋 **Passport Info Added:**\n• Passport Number, Issue/Expiry Dates\n• Place of Birth, Gender, Nationality\n• Issuing Authority`
        : `${matchTypeText}\n\n👤 **Driver Enhanced:** "${passportData.GivenNames} ${passportData.Surname}"\n\n📋 **Passport Information Added:**\n• Passport Number, Issue/Expiry Dates\n• Place of Birth, Gender, Nationality\n• Issuing Authority`;

      addChatMessage('bot', updateMessage);
      
      // Force UI update and show success notification
      setTimeout(() => {
        addChatMessage('bot', `✅ **UI Updated:** Driver form has been refreshed with the new data. You can now see the updated information in the driver section.`);
        setNotification({ type: 'success', message: `Driver "${passportData.GivenNames} ${passportData.Surname}" updated with passport data!` });
        // Switch to drivers tab to show the updated driver
        setActiveCategory(0);
      }, 500);
    } else {
      // Create new driver
      const newDriver: Driver = {
        id: `driver_${Date.now()}`,
        classification: session.drivers.length === 0 ? 'MAIN' : 'NAMED',
        firstName: passportData.GivenNames || '',
        lastName: passportData.Surname || '',
        dateOfBirth: passportData.DateOfBirth || '',
        email: '',
        phone: '',
        address: '',
        postcode: '',
        licenceType: '',
        licenceNumber: '',
        yearsHeld: 0,
        pointsLost: 0,
        licenceIssueDate: '',
        licenceExpiryDate: '',
        licenceValidUntil: '',
        relationship: '',
        sameAddress: true,
        // Passport-specific fields
        passportNumber: passportData.PassportNumber,
        passportIssueDate: passportData.DateOfIssue,
        passportExpiryDate: passportData.DateOfExpiry,
        passportAuthority: passportData.Authority,
        placeOfBirth: passportData.PlaceOfBirth,
        gender: passportData.Gender,
        nationality: passportData.Nationality,
        // Disabilities and Restrictions
        hasDisability: false,
        disabilityTypes: [],
        requiresAdaptations: false,
        adaptationTypes: [],
        automaticOnly: false,
        // Licence Classes
        licenceClassA: false,
        licenceClassB: true,
        licenceClassC: false,
        licenceClassD: false,
        licenceClassBE: false,
        licenceClassCE: false,
        licenceClassDE: false,
        // Driver Communication
        driverEmail: '',
        emailSent: false,
        emailSentDate: ''
      };

      setSession(prev => {
        const updatedDrivers = [...prev.drivers, newDriver];
        console.log('Added new driver:', newDriver);
        console.log('All drivers now:', updatedDrivers);
        return {
          ...prev,
          drivers: updatedDrivers
        };
      });

      // Show existing drivers for reference
      const existingDriversList = session.drivers.length > 0 
        ? session.drivers.map(driver => `• ${driver.firstName} ${driver.lastName}${driver.dateOfBirth ? ` (DOB: ${driver.dateOfBirth})` : ''}`).join('\n')
        : 'No drivers currently in the system';

      addChatMessage('bot', `🆕 **New Driver Created:** "${passportData.GivenNames} ${passportData.Surname}"\n\n🔍 **Matching Result:** No existing driver found with this name\n\n👥 **Existing Drivers:**\n${existingDriversList}\n\n📋 **Insurance Fields:**\n• First Name, Last Name, Date of Birth\n\n📋 **Passport Information:**\n• Passport Number, Issue/Expiry Dates\n• Place of Birth, Gender, Nationality\n• Issuing Authority\n\n💡 **Next Steps:**\nPlease complete the remaining insurance fields (email, phone, address, etc.)`);
      
      // Force UI update and show success notification
      setTimeout(() => {
        addChatMessage('bot', `✅ **UI Updated:** New driver has been added to the form. You can now see the new driver in the driver section.`);
        setNotification({ type: 'success', message: `New driver "${passportData.GivenNames} ${passportData.Surname}" created from passport data!` });
        // Switch to drivers tab to show the new driver
        setActiveCategory(0);
      }, 500);
    }
  };

  const processDrivingLicenceData = (licenceData: any) => {
    // Create detailed extraction report
    const extractedFields = Object.entries(licenceData)
      .filter(([key, value]) => value && value !== '')
      .map(([key, value]) => `• ${key}: ${value}`)
      .join('\n');

    addChatMessage('bot', `🔍 **Document Processing Complete**\n\n📄 **Extracted Fields:**\n${extractedFields}`);

    // Enhanced matching logic: try name + date of birth first, then name only
    let existingDriverIndex = -1;
    let matchType = '';
    
    // First, try to match by name AND date of birth (most precise)
    if (licenceData.DateOfBirth) {
      existingDriverIndex = session.drivers.findIndex(driver => {
        const driverFullName = `${driver.firstName} ${driver.lastName}`;
        const licenceFullName = `${licenceData.GivenNames} ${licenceData.Surname}`;
        const nameMatch = namesMatch(driverFullName, licenceFullName);
        const dobMatch = driver.dateOfBirth === licenceData.DateOfBirth;
        return nameMatch && dobMatch;
      });
      
      if (existingDriverIndex >= 0) {
        matchType = 'name_and_dob';
      }
    }
    
    // If no match by name + DOB, try name only
    if (existingDriverIndex === -1) {
      existingDriverIndex = session.drivers.findIndex(driver => {
        const driverFullName = `${driver.firstName} ${driver.lastName}`;
        const licenceFullName = `${licenceData.GivenNames} ${licenceData.Surname}`;
        return namesMatch(driverFullName, licenceFullName);
      });
      
      if (existingDriverIndex >= 0) {
        matchType = 'name_only';
      }
    }

    if (existingDriverIndex >= 0) {
      // Update existing driver with missing fields
      const updatedDriver = { ...session.drivers[existingDriverIndex] };
      const updatedFields: string[] = [];
      
      if (!updatedDriver.dateOfBirth && licenceData.DateOfBirth) {
        updatedDriver.dateOfBirth = licenceData.DateOfBirth;
        updatedFields.push('Date of Birth');
      }
      if (!updatedDriver.firstName && licenceData.GivenNames) {
        updatedDriver.firstName = licenceData.GivenNames;
        updatedFields.push('First Name');
      }
      if (!updatedDriver.lastName && licenceData.Surname) {
        updatedDriver.lastName = licenceData.Surname;
        updatedFields.push('Last Name');
      }
      
      // Add driving licence-specific fields
      updatedDriver.licenceNumber = licenceData.LicenceNumber;
      updatedDriver.licenceType = determineLicenceType(licenceData);
      updatedDriver.yearsHeld = licenceData.YearsHeld || 0;
      updatedDriver.pointsLost = licenceData.PointsLost || 0;
      updatedDriver.licenceIssueDate = licenceData.DateOfIssue || licenceData.LicenceIssueDate;
      updatedDriver.licenceExpiryDate = licenceData.DateOfExpiry || licenceData.LicenceExpiryDate;
      updatedDriver.licenceValidUntil = licenceData.ValidUntil || licenceData.LicenceValidUntil || '';
      updatedDriver.licenceAuthority = licenceData.Authority;
      updatedDriver.licenceRestrictions = licenceData.Restrictions;
      updatedDriver.licenceEndorsements = licenceData.Endorsements;
      
      // Update licence classes from driving licence data
      updatedDriver.licenceClassA = licenceData.CategoryA === 'Yes' || licenceData.CategoryA === true;
      updatedDriver.licenceClassB = licenceData.CategoryB === 'Yes' || licenceData.CategoryB === true;
      updatedDriver.licenceClassC = licenceData.CategoryC === 'Yes' || licenceData.CategoryC === true;
      updatedDriver.licenceClassD = licenceData.CategoryD === 'Yes' || licenceData.CategoryD === true;
      updatedDriver.licenceClassBE = licenceData.CategoryBE === 'Yes' || licenceData.CategoryBE === true;
      updatedDriver.licenceClassCE = licenceData.CategoryCE === 'Yes' || licenceData.CategoryCE === true;
      updatedDriver.licenceClassDE = licenceData.CategoryDE === 'Yes' || licenceData.CategoryDE === true;
      
      // Check for automatic-only restriction
      if (licenceData.Restrictions && licenceData.Restrictions.toLowerCase().includes('automatic')) {
        updatedDriver.automaticOnly = true;
      }

      setSession(prev => {
        const updatedDrivers = prev.drivers.map((driver, index) => 
          index === existingDriverIndex ? updatedDriver : driver
        );
        console.log('Updated drivers:', updatedDrivers);
        return {
          ...prev,
          drivers: updatedDrivers
        };
      });

      const matchTypeText = matchType === 'name_and_dob' 
        ? '✅ **Exact Match Found:** Name + Date of Birth'
        : '⚠️ **Name Match Found:** Date of Birth differs or missing';
        
      const updateMessage = updatedFields.length > 0 
        ? `${matchTypeText}\n\n👤 **Driver Updated:** "${licenceData.GivenNames} ${licenceData.Surname}"\n\n📝 **Fields Updated:**\n${updatedFields.map(field => `• ${field}`).join('\n')}\n\n📋 **Driving Licence Info Added:**\n• Licence Number, Issue/Expiry Dates\n• Licence Type, Years Held\n• Restrictions, Endorsements\n• Issuing Authority`
        : `${matchTypeText}\n\n👤 **Driver Enhanced:** "${licenceData.GivenNames} ${licenceData.Surname}"\n\n📋 **Driving Licence Information Added:**\n• Licence Number, Issue/Expiry Dates\n• Licence Type, Years Held\n• Restrictions, Endorsements\n• Issuing Authority`;

      addChatMessage('bot', updateMessage);
      
      // Force UI update and show success notification
      setTimeout(() => {
        addChatMessage('bot', `✅ **UI Updated:** Driver form has been refreshed with the new data. You can now see the updated information in the driver section.`);
        setNotification({ type: 'success', message: `Driver "${licenceData.GivenNames} ${licenceData.Surname}" updated with driving licence data!` });
        // Switch to drivers tab to show the updated driver
        setActiveCategory(0);
      }, 500);
    } else {
      // Create new driver
      const newDriver: Driver = {
        id: `driver_${Date.now()}`,
        classification: session.drivers.length === 0 ? 'MAIN' : 'NAMED',
        firstName: licenceData.GivenNames || '',
        lastName: licenceData.Surname || '',
        dateOfBirth: licenceData.DateOfBirth || '',
        email: '',
        phone: '',
        address: '',
        postcode: '',
        licenceType: determineLicenceType(licenceData),
        licenceNumber: licenceData.LicenceNumber || '',
        yearsHeld: licenceData.YearsHeld || 0,
        pointsLost: licenceData.PointsLost || 0,
        licenceIssueDate: licenceData.DateOfIssue || licenceData.LicenceIssueDate || '',
        licenceExpiryDate: licenceData.DateOfExpiry || licenceData.LicenceExpiryDate || '',
        licenceValidUntil: licenceData.ValidUntil || licenceData.LicenceValidUntil || '',
        relationship: '',
        sameAddress: true,
        // Driving licence-specific fields
        licenceAuthority: licenceData.Authority,
        licenceRestrictions: licenceData.Restrictions,
        licenceEndorsements: licenceData.Endorsements,
        // Disabilities and Restrictions
        hasDisability: false,
        disabilityTypes: [],
        requiresAdaptations: false,
        adaptationTypes: [],
        automaticOnly: false,
        // Licence Classes
        licenceClassA: false,
        licenceClassB: true,
        licenceClassC: false,
        licenceClassD: false,
        licenceClassBE: false,
        licenceClassCE: false,
        licenceClassDE: false,
        // Driver Communication
        driverEmail: '',
        emailSent: false,
        emailSentDate: ''
      };

      setSession(prev => {
        const updatedDrivers = [...prev.drivers, newDriver];
        console.log('Added new driver:', newDriver);
        console.log('All drivers now:', updatedDrivers);
        return {
          ...prev,
          drivers: updatedDrivers
        };
      });

      // Show existing drivers for reference
      const existingDriversList = session.drivers.length > 0 
        ? session.drivers.map(driver => `• ${driver.firstName} ${driver.lastName}${driver.dateOfBirth ? ` (DOB: ${driver.dateOfBirth})` : ''}`).join('\n')
        : 'No drivers currently in the system';

      addChatMessage('bot', `🆕 **New Driver Created:** "${licenceData.GivenNames} ${licenceData.Surname}"\n\n🔍 **Matching Result:** No existing driver found with this name\n\n👥 **Existing Drivers:**\n${existingDriversList}\n\n📋 **Insurance Fields:**\n• First Name, Last Name, Date of Birth\n\n📋 **Driving Licence Information:**\n• Licence Number, Issue/Expiry Dates\n• Licence Type, Years Held\n• Restrictions, Endorsements\n• Issuing Authority\n\n💡 **Next Steps:**\nPlease complete the remaining insurance fields (email, phone, address, etc.)`);
      
      // Force UI update and show success notification
      setTimeout(() => {
        addChatMessage('bot', `✅ **UI Updated:** New driver has been added to the form. You can now see the new driver in the driver section.`);
        setNotification({ type: 'success', message: `New driver "${licenceData.GivenNames} ${licenceData.Surname}" created from driving licence data!` });
        // Switch to drivers tab to show the new driver
        setActiveCategory(0);
      }, 500);
    }
  };

  const determineLicenceType = (licenceData: any): string => {
    // Determine licence type based on available categories
    const categories = [];
    if (licenceData.CategoryA) categories.push('A');
    if (licenceData.CategoryB) categories.push('B');
    if (licenceData.CategoryC) categories.push('C');
    if (licenceData.CategoryD) categories.push('D');
    if (licenceData.CategoryBE) categories.push('BE');
    if (licenceData.CategoryCE) categories.push('CE');
    if (licenceData.CategoryDE) categories.push('DE');
    
    return categories.length > 0 ? categories.join(', ') : 'B'; // Default to B if no categories found
  };

  const removeFile = (index: number) => {
    setUploadedFiles(prev => prev.filter((_, i) => i !== index));
  };

  const addChatMessage = (type: 'user' | 'bot', content: string) => {
    setChatMessages(prev => [...prev, {
      type,
      content,
      timestamp: new Date()
    }]);
    setShowChatbot(true);
  };

  const calculateProgress = () => {
    const totalFields = 25; // Approximate total fields
    const completedFields = Object.values(session.progress).filter(Boolean).length;
    return Math.round((completedFields / totalFields) * 100);
  };

  const renderDriversSection = () => (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h2 className="text-xl font-semibold">Driver Information</h2>
        <Button onClick={addDriver} disabled={session.drivers.length >= 4}>
          <Plus className="w-4 h-4 mr-2" />
          Add Driver
        </Button>
      </div>

      {session.drivers.map((driver, index) => (
        <Card key={driver.id} className="form-section">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-medium">
              {driver.classification === 'MAIN' ? 'Main Driver' : `Additional Driver ${index}`}
            </h3>
            {session.drivers.length > 1 && (
              <Button color="failure" size="sm" onClick={() => removeDriver(index)}>
                <X className="w-4 h-4" />
              </Button>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].firstName}
                {getStatusBadge(true, !!driver.firstName)}
                {driver.firstName && (driver as any).passportNumber && (
                  <Badge color="info" className="ml-2 text-xs">Auto-filled</Badge>
                )}
              </Label>
              <TextInput
                value={driver.firstName}
                onChange={(e) => updateDriver(index, 'firstName', e.target.value)}
                placeholder="Enter first name"
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].lastName}
                {getStatusBadge(true, !!driver.lastName)}
                {driver.lastName && (driver as any).passportNumber && (
                  <Badge color="info" className="ml-2 text-xs">Auto-filled</Badge>
                )}
              </Label>
              <TextInput
                value={driver.lastName}
                onChange={(e) => updateDriver(index, 'lastName', e.target.value)}
                placeholder="Enter last name"
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].dateOfBirth}
                {getStatusBadge(true, !!driver.dateOfBirth)}
                {driver.dateOfBirth && (driver as any).passportNumber && (
                  <Badge color="info" className="ml-2 text-xs">Auto-filled</Badge>
                )}
              </Label>
              <TextInput
                type="date"
                value={driver.dateOfBirth}
                onChange={(e) => updateDriver(index, 'dateOfBirth', e.target.value)}
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].email}
                {getStatusBadge(true, !!driver.email)}
              </Label>
              <TextInput
                type="email"
                value={driver.email}
                onChange={(e) => updateDriver(index, 'email', e.target.value)}
                placeholder="your@email.com"
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].phone}
                {getStatusBadge(true, !!driver.phone)}
              </Label>
              <TextInput
                value={driver.phone}
                onChange={(e) => updateDriver(index, 'phone', e.target.value)}
                placeholder="Enter phone number"
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].address}
                {getStatusBadge(true, !!driver.address)}
              </Label>
              <TextInput
                value={driver.address}
                onChange={(e) => updateDriver(index, 'address', e.target.value)}
                placeholder="Enter address"
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].postcode}
                {getStatusBadge(true, !!driver.postcode)}
              </Label>
              <TextInput
                value={driver.postcode}
                onChange={(e) => updateDriver(index, 'postcode', e.target.value)}
                placeholder="Enter postcode"
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].licenceType}
                {getStatusBadge(true, !!driver.licenceType)}
              </Label>
              <Select
                value={driver.licenceType}
                onChange={(e) => updateDriver(index, 'licenceType', e.target.value)}
                className="form-select"
              >
                <option value="">Select licence type</option>
                <option value="full">Full UK Licence</option>
                <option value="provisional">Provisional Licence</option>
                <option value="international">International Licence</option>
              </Select>
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].licenceNumber}
                {getStatusBadge(true, !!driver.licenceNumber)}
              </Label>
              <TextInput
                value={driver.licenceNumber}
                onChange={(e) => updateDriver(index, 'licenceNumber', e.target.value)}
                placeholder="Enter licence number"
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].yearsHeld}
                {getStatusBadge(true, driver.yearsHeld > 0)}
              </Label>
              <TextInput
                type="number"
                value={driver.yearsHeld}
                onChange={(e) => updateDriver(index, 'yearsHeld', parseInt(e.target.value) || 0)}
                placeholder="0"
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].pointsLost}
                {getStatusBadge(false, true)}
              </Label>
              <TextInput
                type="number"
                value={driver.pointsLost}
                onChange={(e) => updateDriver(index, 'pointsLost', parseInt(e.target.value) || 0)}
                placeholder="0"
                min="0"
                max="12"
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].licenceIssueDate}
                {getStatusBadge(true, !!driver.licenceIssueDate)}
              </Label>
              <TextInput
                type="date"
                value={driver.licenceIssueDate}
                onChange={(e) => updateLicenceIssueDate(index, e.target.value)}
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].licenceExpiryDate}
                {getStatusBadge(true, !!driver.licenceExpiryDate)}
                {driver.licenceExpiryDate && getExpiryStatusBadge(driver.licenceExpiryDate)}
              </Label>
              <TextInput
                type="date"
                value={driver.licenceExpiryDate}
                onChange={(e) => updateDriver(index, 'licenceExpiryDate', e.target.value)}
                className="form-input"
              />
            </div>

            <div className="form-group">
              <Label className="form-label">
                {translations[session.language as keyof typeof translations].licenceValidUntil}
                {getStatusBadge(false, true)}
              </Label>
              <TextInput
                type="date"
                value={driver.licenceValidUntil}
                onChange={(e) => updateDriver(index, 'licenceValidUntil', e.target.value)}
                className="form-input"
              />
            </div>

            {driver.classification !== 'MAIN' && (
              <>
                <div className="form-group">
                  <Label className="form-label">
                    Relationship to Main Driver
                    {getStatusBadge(true, !!driver.relationship)}
                  </Label>
                  <Select
                    value={driver.relationship}
                    onChange={(e) => updateDriver(index, 'relationship', e.target.value)}
                    className="form-select"
                  >
                    <option value="">Select relationship</option>
                    <option value="spouse">Spouse/Partner</option>
                    <option value="parent">Parent</option>
                    <option value="child">Child</option>
                    <option value="sibling">Sibling</option>
                    <option value="friend">Friend</option>
                    <option value="other">Other</option>
                  </Select>
                </div>

                <div className="form-group">
                  <Label className="form-label">
                    Lives at Same Address
                    {getStatusBadge(false, true)}
                  </Label>
                  <ToggleSwitch
                    checked={driver.sameAddress}
                    onChange={(checked) => updateDriver(index, 'sameAddress', checked)}
                  />
                </div>
              </>
            )}

            {/* Passport Information Section */}
            {(driver.passportNumber || driver.passportIssueDate || driver.passportExpiryDate) && (
              <div className="col-span-full mt-6 p-4 bg-blue-50 rounded-lg border border-blue-200">
                <h4 className="text-sm font-semibold text-blue-800 mb-3 flex items-center">
                  <FileText className="w-4 h-4 mr-2" />
                  Passport Information (Document-derived)
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {driver.passportNumber && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-blue-700">
                        Passport Number
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.passportNumber}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}
                  
                  {driver.passportIssueDate && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-blue-700">
                        Issue Date
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.passportIssueDate}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}

                  {driver.passportExpiryDate && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-blue-700">
                        Expiry Date
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.passportExpiryDate}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}

                  {driver.passportAuthority && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-blue-700">
                        Issuing Authority
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.passportAuthority}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}

                  {driver.placeOfBirth && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-blue-700">
                        Place of Birth
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.placeOfBirth}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}

                  {driver.gender && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-blue-700">
                        Gender
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.gender}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}

                  {driver.nationality && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-blue-700">
                        Nationality
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.nationality}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Driving Licence Information Section */}
            {(driver.licenceIssueDate || driver.licenceExpiryDate || driver.licenceAuthority || driver.licenceRestrictions || driver.licenceEndorsements) && (
              <div className="col-span-full mt-6 p-4 bg-green-50 rounded-lg border border-green-200">
                <h4 className="text-sm font-semibold text-green-800 mb-3 flex items-center">
                  <Car className="w-4 h-4 mr-2" />
                  Driving Licence Information (Document-derived)
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {driver.licenceIssueDate && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-green-700">
                        Issue Date
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.licenceIssueDate}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}

                  {driver.licenceExpiryDate && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-green-700">
                        Expiry Date
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                        {getExpiryStatusBadge(driver.licenceExpiryDate)}
                      </Label>
                      <TextInput
                        value={driver.licenceExpiryDate}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}

                  {driver.licenceAuthority && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-green-700">
                        Issuing Authority
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.licenceAuthority}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}

                  {driver.licenceRestrictions && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-green-700">
                        Restrictions
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.licenceRestrictions}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}

                  {driver.licenceEndorsements && (
                    <div className="form-group">
                      <Label className="form-label text-sm text-green-700">
                        Endorsements
                        <Badge color="info" className="ml-2 text-xs">From Document</Badge>
                      </Label>
                      <TextInput
                        value={driver.licenceEndorsements}
                        readOnly
                        className="form-input bg-gray-100"
                      />
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Disabilities and Restrictions Section */}
            <div className="col-span-full mt-6 p-4 bg-yellow-50 rounded-lg border border-yellow-200">
              <h4 className="text-sm font-semibold text-yellow-800 mb-3 flex items-center">
                <AlertCircle className="w-4 h-4 mr-2" />
                Disabilities and Restrictions
              </h4>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-group">
                  <Label className="form-label">
                    {translations[session.language as keyof typeof translations].hasDisability}
                    {getStatusBadge(false, true)}
                  </Label>
                  <ToggleSwitch
                    checked={driver.hasDisability}
                    onChange={(checked) => updateDriver(index, 'hasDisability', checked)}
                  />
                </div>

                {driver.hasDisability && (
                  <>
                    <div className="form-group">
                      <Label className="form-label">
                        {translations[session.language as keyof typeof translations].disabilityType}
                        {getStatusBadge(true, !!driver.disabilityTypes)}
                      </Label>
                      <Select
                        value={driver.disabilityTypes}
                        onChange={(e) => updateDriver(index, 'disabilityTypes', e.target.value)}
                        className="form-select"
                      >
                        <option value="">Select disability type</option>
                        <option value="Visual impairment">Visual impairment</option>
                        <option value="Hearing impairment">Hearing impairment</option>
                        <option value="Mobility impairment">Mobility impairment</option>
                        <option value="Cognitive impairment">Cognitive impairment</option>
                        <option value="Neurological condition">Neurological condition</option>
                        <option value="Cardiovascular condition">Cardiovascular condition</option>
                        <option value="Respiratory condition">Respiratory condition</option>
                        <option value="Diabetes">Diabetes</option>
                        <option value="Epilepsy">Epilepsy</option>
                        <option value="Mental health condition">Mental health condition</option>
                        <option value="Other medical condition">Other medical condition</option>
                      </Select>
                    </div>

                    <div className="form-group">
                      <Label className="form-label">
                        {translations[session.language as keyof typeof translations].requiresAdaptations}
                        {getStatusBadge(false, true)}
                      </Label>
                      <ToggleSwitch
                        checked={driver.requiresAdaptations}
                        onChange={(checked) => updateDriver(index, 'requiresAdaptations', checked)}
                      />
                    </div>

                    {driver.requiresAdaptations && (
                      <div className="form-group">
                        <Label className="form-label">
                          {translations[session.language as keyof typeof translations].adaptationType}
                          {getStatusBadge(true, !!driver.adaptationTypes)}
                        </Label>
                        <Select
                          value={driver.adaptationTypes}
                          onChange={(e) => updateDriver(index, 'adaptationTypes', e.target.value)}
                          className="form-select"
                        >
                          <option value="">Select adaptation type</option>
                          <option value="Hand controls">Hand controls</option>
                          <option value="Foot controls">Foot controls</option>
                          <option value="Steering wheel spinner">Steering wheel spinner</option>
                          <option value="Pedal extensions">Pedal extensions</option>
                          <option value="Seat modifications">Seat modifications</option>
                          <option value="Mirror adaptations">Mirror adaptations</option>
                          <option value="Brake modifications">Brake modifications</option>
                          <option value="Accelerator modifications">Accelerator modifications</option>
                          <option value="Gear lever modifications">Gear lever modifications</option>
                          <option value="Other adaptation">Other adaptation</option>
                        </Select>
                      </div>
                    )}
                  </>
                )}

                <div className="form-group">
                  <Label className="form-label">
                    {translations[session.language as keyof typeof translations].automaticOnly}
                    {getStatusBadge(false, true)}
                  </Label>
                  <ToggleSwitch
                    checked={driver.automaticOnly}
                    onChange={(checked) => updateDriver(index, 'automaticOnly', checked)}
                  />
                </div>
              </div>
            </div>

            {/* Licence Classes Section */}
            <div className="col-span-full mt-6 p-4 bg-green-50 rounded-lg border border-green-200">
              <h4 className="text-sm font-semibold text-green-800 mb-3 flex items-center">
                <Car className="w-4 h-4 mr-2" />
                {translations[session.language as keyof typeof translations].licenceClasses}
              </h4>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-group">
                  <Label className="form-label">
                    {translations[session.language as keyof typeof translations].licenceClassA}
                    {getStatusBadge(false, true)}
                  </Label>
                  <ToggleSwitch
                    checked={driver.licenceClassA}
                    onChange={(checked) => updateDriver(index, 'licenceClassA', checked)}
                  />
                </div>

                <div className="form-group">
                  <Label className="form-label">
                    {translations[session.language as keyof typeof translations].licenceClassB}
                    {getStatusBadge(false, true)}
                  </Label>
                  <ToggleSwitch
                    checked={driver.licenceClassB}
                    onChange={(checked) => updateDriver(index, 'licenceClassB', checked)}
                  />
                </div>

                <div className="form-group">
                  <Label className="form-label">
                    {translations[session.language as keyof typeof translations].licenceClassC}
                    {getStatusBadge(false, true)}
                  </Label>
                  <ToggleSwitch
                    checked={driver.licenceClassC}
                    onChange={(checked) => updateDriver(index, 'licenceClassC', checked)}
                  />
                </div>

                <div className="form-group">
                  <Label className="form-label">
                    {translations[session.language as keyof typeof translations].licenceClassD}
                    {getStatusBadge(false, true)}
                  </Label>
                  <ToggleSwitch
                    checked={driver.licenceClassD}
                    onChange={(checked) => updateDriver(index, 'licenceClassD', checked)}
                  />
                </div>

                <div className="form-group">
                  <Label className="form-label">
                    {translations[session.language as keyof typeof translations].licenceClassBE}
                    {getStatusBadge(false, true)}
                  </Label>
                  <ToggleSwitch
                    checked={driver.licenceClassBE}
                    onChange={(checked) => updateDriver(index, 'licenceClassBE', checked)}
                  />
                </div>

                <div className="form-group">
                  <Label className="form-label">
                    {translations[session.language as keyof typeof translations].licenceClassCE}
                    {getStatusBadge(false, true)}
                  </Label>
                  <ToggleSwitch
                    checked={driver.licenceClassCE}
                    onChange={(checked) => updateDriver(index, 'licenceClassCE', checked)}
                  />
                </div>

                <div className="form-group">
                  <Label className="form-label">
                    {translations[session.language as keyof typeof translations].licenceClassDE}
                    {getStatusBadge(false, true)}
                  </Label>
                  <ToggleSwitch
                    checked={driver.licenceClassDE}
                    onChange={(checked) => updateDriver(index, 'licenceClassDE', checked)}
                  />
                </div>
              </div>
            </div>

            {/* Driver Email Section (for additional drivers) */}
            {driver.classification !== 'MAIN' && (
              <div className="col-span-full mt-6 p-4 bg-purple-50 rounded-lg border border-purple-200">
                <h4 className="text-sm font-semibold text-purple-800 mb-3 flex items-center">
                  <Settings className="w-4 h-4 mr-2" />
                  Driver Communication
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="form-group">
                    <Label className="form-label">
                                              {translations[session.language as keyof typeof translations].driverEmail}
                      {getStatusBadge(true, !!driver.driverEmail)}
                    </Label>
                    <TextInput
                      type="email"
                      value={driver.driverEmail}
                      onChange={(e) => updateDriver(index, 'driverEmail', e.target.value)}
                      placeholder="driver@email.com"
                      className="form-input"
                    />
                  </div>

                  <div className="form-group flex items-end">
                    <Button
                      color="purple"
                      onClick={() => sendEmailToDriver(index)}
                      disabled={!driver.driverEmail || driver.emailSent}
                      className="w-full"
                    >
                      {driver.emailSent ? (
                        <>
                          <CheckCircle className="w-4 h-4 mr-2" />
                          Email Sent
                        </>
                      ) : (
                        <>
                          <Mail className="w-4 h-4 mr-2" />
                          {translations[session.language as keyof typeof translations].sendEmailToDriver}
                        </>
                      )}
                    </Button>
                  </div>

                  {driver.emailSent && (
                    <div className="col-span-full">
                      <Alert color="success" className="mt-2">
                        <div className="flex items-center gap-2">
                          <CheckCircle className="w-4 h-4" />
                          <span>Email sent to {driver.driverEmail} on {new Date(driver.emailSentDate).toLocaleDateString()}</span>
                        </div>
                      </Alert>
                    </div>
                  )}
                </div>
              </div>
            )}
          </div>
        </Card>
      ))}
    </div>
  );

  const renderVehicleSection = () => (
    <div className="space-y-6">
      <h2 className="text-xl font-semibold">Vehicle Information</h2>

      <Card className="form-section">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="form-group">
            <Label className="form-label">
              Registration Number
              {getStatusBadge(true, !!session.vehicle.registration)}
            </Label>
            <TextInput
              value={session.vehicle.registration}
              onChange={(e) => updateVehicle('registration', e.target.value.toUpperCase())}
              placeholder="Enter registration"
              className="form-input"
            />
          </div>

          <div className="form-group">
            <Label className="form-label">
              Make
              {getStatusBadge(true, !!session.vehicle.make)}
            </Label>
            <TextInput
              value={session.vehicle.make}
              onChange={(e) => updateVehicle('make', e.target.value)}
              placeholder="Enter make"
              className="form-input"
            />
          </div>

          <div className="form-group">
            <Label className="form-label">
              Model
              {getStatusBadge(true, !!session.vehicle.model)}
            </Label>
            <TextInput
              value={session.vehicle.model}
              onChange={(e) => updateVehicle('model', e.target.value)}
              placeholder="Enter model"
              className="form-input"
            />
          </div>

          <div className="form-group">
            <Label className="form-label">
              Year
              {getStatusBadge(true, session.vehicle.year > 0)}
            </Label>
            <TextInput
              type="number"
              value={session.vehicle.year}
              onChange={(e) => updateVehicle('year', parseInt(e.target.value) || 0)}
              placeholder="Enter year"
              className="form-input"
            />
          </div>

          <div className="form-group">
            <Label className="form-label">
              Mileage
              {getStatusBadge(true, session.vehicle.mileage > 0)}
            </Label>
            <TextInput
              type="number"
              value={session.vehicle.mileage}
              onChange={(e) => updateVehicle('mileage', parseInt(e.target.value) || 0)}
              placeholder="Enter mileage"
              className="form-input"
            />
          </div>

          <div className="form-group">
            <Label className="form-label">
              Vehicle Value
              {getStatusBadge(true, session.vehicle.value > 0)}
            </Label>
            <TextInput
              type="number"
              value={session.vehicle.value}
              onChange={(e) => updateVehicle('value', parseFloat(e.target.value) || 0)}
              placeholder="Enter value"
              className="form-input"
            />
          </div>

          <div className="form-group">
            <Label className="form-label">
              Daytime Location
              {getStatusBadge(false, !!session.vehicle.daytimeLocation)}
            </Label>
            <Select
              value={session.vehicle.daytimeLocation}
              onChange={(e) => updateVehicle('daytimeLocation', e.target.value)}
              className="form-select"
            >
              <option value="">Select location</option>
              <option value="driveway">Driveway</option>
              <option value="garage">Garage</option>
              <option value="street">Street</option>
              <option value="car_park">Car Park</option>
              <option value="work_car_park">Work Car Park</option>
              <option value="public_car_park">Public Car Park</option>
              <option value="private_car_park">Private Car Park</option>
            </Select>
          </div>

          <div className="form-group">
            <Label className="form-label">
              Overnight Location
              {getStatusBadge(true, !!session.vehicle.overnightLocation)}
            </Label>
            <Select
              value={session.vehicle.overnightLocation}
              onChange={(e) => updateVehicle('overnightLocation', e.target.value)}
              className="form-select"
            >
              <option value="">Select location</option>
              <option value="driveway">Driveway</option>
              <option value="garage">Garage</option>
              <option value="street">Street</option>
              <option value="car_park">Car Park</option>
            </Select>
          </div>

          <div className="form-group">
            <Label className="form-label">
              Vehicle Modifications
              {getStatusBadge(false, true)}
            </Label>
            <ToggleSwitch
              checked={session.vehicle.hasModifications}
              onChange={(checked) => updateVehicle('hasModifications', checked)}
            />
          </div>
        </div>

        {session.vehicle.hasModifications && (
          <div className="mt-6 p-4 bg-gray-50 rounded-lg">
            <h4 className="font-medium mb-4">Vehicle Modifications</h4>
            <p className="text-sm text-gray-600 mb-4">
              Select all modifications that apply to your vehicle. Multiple selections are allowed.
            </p>
            <div className="space-y-3">
              {[
                'None',
                'Engine Tuning',
                'Exhaust System', 
                'Air Intake',
                'Suspension',
                'Wheels/Tyres',
                'Body Kit',
                'Interior',
                'Audio System',
                'Performance Chip',
                'Turbo/Supercharger',
                'Other'
              ].map((modification) => (
                <div key={modification} className="flex items-center">
                  <input
                    type="checkbox"
                    id={`modification-${modification}`}
                    checked={session.vehicle.modifications.includes(modification)}
                    onChange={(e) => {
                      if (e.target.checked) {
                        updateVehicle('modifications', [...session.vehicle.modifications, modification]);
                      } else {
                        updateVehicle('modifications', session.vehicle.modifications.filter(m => m !== modification));
                      }
                    }}
                    className="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  />
                  <label htmlFor={`modification-${modification}`} className="ml-2 text-sm text-gray-700">
                    {modification}
                  </label>
                </div>
              ))}
            </div>
            {session.vehicle.modifications.length > 0 && (
              <div className="mt-4 p-3 bg-blue-50 rounded-lg">
                <p className="text-sm text-blue-800">
                  <strong>Selected Modifications:</strong> {session.vehicle.modifications.join(', ')}
                </p>
              </div>
            )}
          </div>
        )}
      </Card>
    </div>
  );

  const renderPolicySection = () => (
    <div className="space-y-6">
      <h2 className="text-xl font-semibold">Policy Details</h2>

      <Card className="form-section">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="form-group">
            <Label className="form-label">
              Cover Type
              {getStatusBadge(true, !!session.policy.coverType)}
            </Label>
            <Select
              value={session.policy.coverType}
              onChange={(e) => updatePolicy('coverType', e.target.value)}
              className="form-select"
            >
              <option value="">Select cover type</option>
              <option value="comprehensive">Comprehensive</option>
              <option value="third_party_fire_theft">Third Party, Fire & Theft</option>
              <option value="third_party">Third Party Only</option>
            </Select>
          </div>

          <div className="form-group">
            <Label className="form-label">
              Start Date
              {getStatusBadge(true, !!session.policy.startDate)}
            </Label>
            <TextInput
              type="date"
              value={session.policy.startDate}
              onChange={(e) => updatePolicy('startDate', e.target.value)}
              className="form-input"
            />
          </div>

          <div className="form-group">
            <Label className="form-label">
              Voluntary Excess
              {getStatusBadge(false, true)}
            </Label>
            <Select
              value={session.policy.voluntaryExcess}
              onChange={(e) => updatePolicy('voluntaryExcess', parseInt(e.target.value) || 0)}
              className="form-select"
            >
              <option value="0">£0</option>
              <option value="100">£100</option>
              <option value="250">£250</option>
              <option value="500">£500</option>
              <option value="1000">£1,000</option>
            </Select>
          </div>

          <div className="form-group">
            <Label className="form-label">
              No Claims Discount Years
              {getStatusBadge(false, true)}
            </Label>
            <TextInput
              type="number"
              value={session.policy.ncdYears}
              onChange={(e) => updatePolicy('ncdYears', parseInt(e.target.value) || 0)}
              placeholder="0"
              className="form-input"
            />
          </div>

          <div className="form-group">
            <Label className="form-label">
              Protect No Claims Discount
              {getStatusBadge(false, true)}
            </Label>
            <ToggleSwitch
              checked={session.policy.protectNCD}
              onChange={(checked) => updatePolicy('protectNCD', checked)}
            />
          </div>
        </div>

        {/* Additional Cover Section */}
        <div className="mt-6">
          <h4 className="font-medium mb-4">Additional Cover</h4>
          <div className="space-y-4">
            <div className="form-group">
              <Label className="form-label">
                Breakdown Cover
                {getStatusBadge(false, true)}
              </Label>
              <Select
                value={session.extras.breakdownCover}
                onChange={(e) => updateExtras('breakdownCover', e.target.value)}
                className="form-select"
              >
                <option value="">Select cover</option>
                <option value="none">No Breakdown Cover</option>
                <option value="basic">Basic Breakdown</option>
                <option value="comprehensive">Comprehensive Breakdown</option>
              </Select>
            </div>

            <div className="flex items-center gap-4">
              <div className="flex items-center gap-2">
                <ToggleSwitch
                  checked={session.extras.legalExpenses}
                  onChange={(checked) => updateExtras('legalExpenses', checked)}
                />
                <Label>Legal Expenses Cover</Label>
              </div>

              <div className="flex items-center gap-2">
                <ToggleSwitch
                  checked={session.extras.courtesyCar}
                  onChange={(checked) => updateExtras('courtesyCar', checked)}
                />
                <Label>Courtesy Car</Label>
              </div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );

  const renderClaimsSection = () => (
    <div className="space-y-6">
      <h2 className="text-xl font-semibold">{translations[session.language as keyof typeof translations].claimsHistory}</h2>

      <Card className="form-section">
        <div className="space-y-6">
          {/* Claims Toggle */}
          <div className="form-group">
            <Label className="form-label">
              {translations[session.language as keyof typeof translations].claimsHistory} (Last 5 Years)
              {getStatusBadge(false, true)}
            </Label>
            <ToggleSwitch
              checked={session.claims.hasClaims}
              onChange={(checked) => updateClaims('hasClaims', checked)}
            />
          </div>

          {/* Convictions Toggle */}
          <div className="form-group">
            <Label className="form-label">
              {translations[session.language as keyof typeof translations].convictions} (Last 5 Years)
              {getStatusBadge(false, true)}
            </Label>
            <ToggleSwitch
              checked={session.claims.hasConvictions}
              onChange={(checked) => updateClaims('hasConvictions', checked)}
            />
          </div>

          {/* Accidents Toggle */}
          <div className="form-group">
            <Label className="form-label">
              {translations[session.language as keyof typeof translations].accidents} (Last 5 Years)
              {getStatusBadge(false, true)}
            </Label>
            <ToggleSwitch
              checked={session.claims.hasAccidents}
              onChange={(checked) => updateClaims('hasAccidents', checked)}
            />
          </div>

          {/* Claims Form */}
          {session.claims.hasClaims && (
            <div className="mt-6 p-4 bg-gray-50 rounded-lg">
              <h3 className="text-lg font-medium mb-4">{translations[session.language as keyof typeof translations].claimsHistory}</h3>
              <div className="space-y-4">
                {session.claims.claims.map((claim, index) => (
                  <div key={claim.id} className="p-4 border border-gray-200 rounded-lg">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].claimDate}</Label>
                        <TextInput
                          type="date"
                          value={claim.date}
                          onChange={(e) => updateClaim(index, 'date', e.target.value)}
                          className="form-input"
                        />
                      </div>
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].claimType}</Label>
                        <Select
                          value={claim.type}
                          onChange={(e) => updateClaim(index, 'type', e.target.value)}
                          className="form-select"
                        >
                          <option value="">Select claim type</option>
                          <option value="Accident - At Fault">Accident - At Fault</option>
                          <option value="Accident - Not At Fault">Accident - Not At Fault</option>
                          <option value="Accident - Split Liability">Accident - Split Liability</option>
                          <option value="Theft - Vehicle">Theft - Vehicle</option>
                          <option value="Theft - Parts/Accessories">Theft - Parts/Accessories</option>
                          <option value="Fire Damage">Fire Damage</option>
                          <option value="Flood Damage">Flood Damage</option>
                          <option value="Storm Damage">Storm Damage</option>
                          <option value="Vandalism">Vandalism</option>
                          <option value="Windscreen Damage">Windscreen Damage</option>
                          <option value="Animal Strike">Animal Strike</option>
                          <option value="Medical Expenses">Medical Expenses</option>
                          <option value="Personal Injury">Personal Injury</option>
                          <option value="Legal Expenses">Legal Expenses</option>
                          <option value="Recovery Costs">Recovery Costs</option>
                          <option value="Other">Other</option>
                        </Select>
                      </div>
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].claimAmount}</Label>
                        <TextInput
                          type="number"
                          value={claim.amount}
                          onChange={(e) => updateClaim(index, 'amount', parseFloat(e.target.value) || 0)}
                          placeholder="0.00"
                          className="form-input"
                        />
                      </div>
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].faultStatus}</Label>
                        <Select
                          value={claim.faultStatus}
                          onChange={(e) => updateClaim(index, 'faultStatus', e.target.value)}
                          className="form-select"
                        >
                          <option value="">Select fault status</option>
                          <option value="At Fault">At Fault</option>
                          <option value="Not At Fault">Not At Fault</option>
                          <option value="Split Liability">Split Liability</option>
                          <option value="Uncertain">Uncertain</option>
                        </Select>
                      </div>
                    </div>
                    <div className="mt-4 flex justify-end">
                      <Button
                        color="failure"
                        size="sm"
                        onClick={() => removeClaim(index)}
                      >
                        Remove Claim
                      </Button>
                    </div>
                  </div>
                ))}
                <Button
                  color="gray"
                  onClick={addClaim}
                  className="w-full"
                >
                  {translations[session.language as keyof typeof translations].addClaim}
                </Button>
              </div>
            </div>
          )}

          {/* Convictions Form */}
          {session.claims.hasConvictions && (
            <div className="mt-6 p-4 bg-gray-50 rounded-lg">
              <h3 className="text-lg font-medium mb-4">{translations[session.language as keyof typeof translations].convictions}</h3>
              <div className="space-y-4">
                {session.claims.convictions.map((conviction, index) => (
                  <div key={conviction.id} className="p-4 border border-gray-200 rounded-lg">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].convictionDate}</Label>
                        <TextInput
                          type="date"
                          value={conviction.date}
                          onChange={(e) => updateConviction(index, 'date', e.target.value)}
                          className="form-input"
                        />
                      </div>
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].offenceCode}</Label>
                        <Select
                          value={conviction.offenceCode}
                          onChange={(e) => updateConviction(index, 'offenceCode', e.target.value)}
                          className="form-select"
                        >
                          <option value="">Select offence code</option>
                          <option value="AC10 - Failing to stop after an accident">AC10 - Failing to stop after an accident</option>
                          <option value="AC20 - Failing to give particulars or to report an accident within 24 hours">AC20 - Failing to give particulars or to report an accident within 24 hours</option>
                          <option value="CD10 - Driving without due care and attention">CD10 - Driving without due care and attention</option>
                          <option value="CD20 - Driving without reasonable consideration for other road users">CD20 - Driving without reasonable consideration for other road users</option>
                          <option value="DD10 - Causing death by dangerous driving">DD10 - Causing death by dangerous driving</option>
                          <option value="DD20 - Dangerous driving">DD20 - Dangerous driving</option>
                          <option value="DR10 - Driving or attempting to drive with alcohol level above limit">DR10 - Driving or attempting to drive with alcohol level above limit</option>
                          <option value="DR20 - Driving or attempting to drive while unfit through drink">DR20 - Driving or attempting to drive while unfit through drink</option>
                          <option value="IN10 - Using a vehicle uninsured against third party risks">IN10 - Using a vehicle uninsured against third party risks</option>
                          <option value="LC20 - Driving otherwise than in accordance with a licence">LC20 - Driving otherwise than in accordance with a licence</option>
                          <option value="SP10 - Exceeding goods vehicle speed limits">SP10 - Exceeding goods vehicle speed limits</option>
                          <option value="SP20 - Exceeding speed limit for type of vehicle">SP20 - Exceeding speed limit for type of vehicle</option>
                          <option value="SP30 - Exceeding statutory speed limit on a public road">SP30 - Exceeding statutory speed limit on a public road</option>
                          <option value="SP40 - Exceeding passenger vehicle speed limit">SP40 - Exceeding passenger vehicle speed limit</option>
                          <option value="SP50 - Exceeding speed limit on a motorway">SP50 - Exceeding speed limit on a motorway</option>
                          <option value="TS10 - Failing to comply with traffic light signals">TS10 - Failing to comply with traffic light signals</option>
                          <option value="TS20 - Failing to comply with double white lines">TS20 - Failing to comply with double white lines</option>
                          <option value="TS30 - Failing to comply with 'stop' sign">TS30 - Failing to comply with 'stop' sign</option>
                          <option value="Other - Other offence not listed">Other - Other offence not listed</option>
                        </Select>
                      </div>
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].penaltyPoints}</Label>
                        <TextInput
                          type="number"
                          value={conviction.penaltyPoints}
                          onChange={(e) => updateConviction(index, 'penaltyPoints', parseInt(e.target.value) || 0)}
                          min="0"
                          max="12"
                          className="form-input"
                        />
                      </div>
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].fine}</Label>
                        <TextInput
                          type="number"
                          value={conviction.fine || ''}
                          onChange={(e) => updateConviction(index, 'fine', parseFloat(e.target.value) || 0)}
                          placeholder="0.00"
                          className="form-input"
                        />
                      </div>
                    </div>
                    <div className="mt-4 flex justify-end">
                      <Button
                        color="failure"
                        size="sm"
                        onClick={() => removeConviction(index)}
                      >
                        Remove Conviction
                      </Button>
                    </div>
                  </div>
                ))}
                <Button
                  color="gray"
                  onClick={addConviction}
                  className="w-full"
                >
                  {translations[session.language as keyof typeof translations].addConviction}
                </Button>
              </div>
            </div>
          )}

          {/* Accidents Form */}
          {session.claims.hasAccidents && (
            <div className="mt-6 p-4 bg-gray-50 rounded-lg">
              <h3 className="text-lg font-medium mb-4">{translations[session.language as keyof typeof translations].accidents}</h3>
              <div className="space-y-4">
                {session.claims.accidents.map((accident, index) => (
                  <div key={accident.id} className="p-4 border border-gray-200 rounded-lg">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].accidentDate}</Label>
                        <TextInput
                          type="date"
                          value={accident.date}
                          onChange={(e) => updateAccident(index, 'date', e.target.value)}
                          className="form-input"
                        />
                      </div>
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].accidentType}</Label>
                        <Select
                          value={accident.type}
                          onChange={(e) => updateAccident(index, 'type', e.target.value)}
                          className="form-select"
                        >
                          <option value="">Select accident type</option>
                          <option value="Collision with another vehicle">Collision with another vehicle</option>
                          <option value="Collision with stationary object">Collision with stationary object</option>
                          <option value="Collision with pedestrian">Collision with pedestrian</option>
                          <option value="Collision with animal">Collision with animal</option>
                          <option value="Rollover">Rollover</option>
                          <option value="Fire">Fire</option>
                          <option value="Flood damage">Flood damage</option>
                          <option value="Storm damage">Storm damage</option>
                          <option value="Vandalism">Vandalism</option>
                          <option value="Theft">Theft</option>
                          <option value="Other">Other</option>
                        </Select>
                      </div>
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].accidentSeverity}</Label>
                        <Select
                          value={accident.severity}
                          onChange={(e) => updateAccident(index, 'severity', e.target.value)}
                          className="form-select"
                        >
                          <option value="">Select severity</option>
                          <option value="Minor - No injuries, minor damage">Minor - No injuries, minor damage</option>
                          <option value="Moderate - Minor injuries, moderate damage">Moderate - Minor injuries, moderate damage</option>
                          <option value="Serious - Serious injuries, major damage">Serious - Serious injuries, major damage</option>
                          <option value="Fatal - Fatal injuries">Fatal - Fatal injuries</option>
                        </Select>
                      </div>
                      <div className="form-group">
                        <Label className="form-label">{translations[session.language as keyof typeof translations].claimMade}</Label>
                        <ToggleSwitch
                          checked={accident.claimMade}
                          onChange={(checked) => updateAccident(index, 'claimMade', checked)}
                        />
                      </div>
                    </div>
                    <div className="mt-4 flex justify-end">
                      <Button
                        color="failure"
                        size="sm"
                        onClick={() => removeAccident(index)}
                      >
                        Remove Accident
                      </Button>
                    </div>
                  </div>
                ))}
                <Button
                  color="gray"
                  onClick={addAccident}
                  className="w-full"
                >
                  {translations[session.language as keyof typeof translations].addAccident}
                </Button>
              </div>
            </div>
          )}

          {(session.claims.hasClaims || session.claims.hasConvictions || session.claims.hasAccidents) && (
            <Alert color="warning" className="mt-4">
              <div className="flex items-center gap-2">
                <AlertCircle className="w-5 h-5" />
                <span>Please provide complete details of any claims, convictions, or accidents when prompted.</span>
              </div>
            </Alert>
          )}
        </div>
      </Card>
    </div>
  );

  const renderPaymentSection = () => (
    <div className="space-y-6">
      <h2 className="text-xl font-semibold">Payment Details</h2>

      <Card className="form-section">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div className="form-group">
            <Label className="form-label">
              Payment Frequency
              {getStatusBadge(true, !!session.payment.frequency)}
            </Label>
            <Select
              value={session.payment.frequency}
              onChange={(e) => updatePayment('frequency', e.target.value)}
              className="form-select"
            >
              <option value="">Select frequency</option>
              <option value="monthly">Monthly</option>
              <option value="quarterly">Quarterly</option>
              <option value="annually">Annually</option>
            </Select>
          </div>

          <div className="form-group">
            <Label className="form-label">
              Payment Method
              {getStatusBadge(true, !!session.payment.method)}
            </Label>
            <Select
              value={session.payment.method}
              onChange={(e) => updatePayment('method', e.target.value)}
              className="form-select"
            >
              <option value="">Select method</option>
              <option value="direct_debit">Direct Debit</option>
              <option value="card">Credit/Debit Card</option>
              <option value="bank_transfer">Bank Transfer</option>
            </Select>
          </div>
        </div>
      </Card>
    </div>
  );

  const renderMarketingSection = () => (
    <div className="space-y-6">
      <h2 className="text-xl font-semibold">Marketing Preferences</h2>

      <Card className="form-section">
        <div className="space-y-4">
          <div className="flex items-center gap-2">
            <ToggleSwitch
              checked={session.marketing.emailMarketing}
              onChange={(checked) => updateMarketing('emailMarketing', checked)}
            />
            <Label>Email Marketing</Label>
          </div>

          <div className="flex items-center gap-2">
            <ToggleSwitch
              checked={session.marketing.smsMarketing}
              onChange={(checked) => updateMarketing('smsMarketing', checked)}
            />
            <Label>SMS Marketing</Label>
          </div>

          <div className="flex items-center gap-2">
            <ToggleSwitch
              checked={session.marketing.postMarketing}
              onChange={(checked) => updateMarketing('postMarketing', checked)}
            />
            <Label>Post Marketing</Label>
          </div>
        </div>
      </Card>
    </div>
  );

  // Document Processing Sections
  const renderPassportSection = () => (
    <div className="space-y-6">
      <Card>
        <div className="p-6">
          <div className="flex items-center gap-3 mb-6">
            <FileText className="w-6 h-6 text-blue-600" />
            <h2 className="text-xl font-semibold text-gray-900">Passport Document Processing</h2>
          </div>
          
          <div className="space-y-6">
            <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center">
              <Upload className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">Upload Passport</h3>
              <p className="text-gray-600 mb-4">Upload a clear photo or scan of your passport</p>
              <Button color="blue">
                <Upload className="w-4 h-4 mr-2" />
                Choose File
              </Button>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label htmlFor="passportNumber">Passport Number</Label>
                <TextInput id="passportNumber" placeholder="Enter passport number" />
              </div>
              <div>
                <Label htmlFor="passportCountry">Issuing Country</Label>
                <Select id="passportCountry">
                  <option value="">{translations[session.language as keyof typeof translations].selectCountry}</option>
                  <option value="GB">United Kingdom</option>
                  <option value="US">United States</option>
                  <option value="DE">Germany</option>
                  <option value="FR">France</option>
                  <option value="IT">Italy</option>
                  <option value="ES">Spain</option>
                  <option value="NL">Netherlands</option>
                  <option value="BE">Belgium</option>
                  <option value="AT">Austria</option>
                  <option value="CH">Switzerland</option>
                  <option value="IE">Ireland</option>
                  <option value="PT">Portugal</option>
                  <option value="GR">Greece</option>
                  <option value="PL">Poland</option>
                  <option value="CZ">Czech Republic</option>
                  <option value="HU">Hungary</option>
                  <option value="SK">Slovakia</option>
                  <option value="SI">Slovenia</option>
                  <option value="HR">Croatia</option>
                  <option value="RO">Romania</option>
                  <option value="BG">Bulgaria</option>
                  <option value="LT">Lithuania</option>
                  <option value="LV">Latvia</option>
                  <option value="EE">Estonia</option>
                  <option value="FI">Finland</option>
                  <option value="SE">Sweden</option>
                  <option value="DK">Denmark</option>
                  <option value="NO">Norway</option>
                  <option value="IS">Iceland</option>
                  <option value="LU">Luxembourg</option>
                  <option value="MT">Malta</option>
                  <option value="CY">Cyprus</option>
                  <option value="CA">Canada</option>
                  <option value="AU">Australia</option>
                  <option value="NZ">New Zealand</option>
                  <option value="JP">Japan</option>
                  <option value="KR">South Korea</option>
                  <option value="CN">China</option>
                  <option value="IN">India</option>
                  <option value="BR">Brazil</option>
                  <option value="MX">Mexico</option>
                  <option value="AR">Argentina</option>
                  <option value="CL">Chile</option>
                  <option value="ZA">South Africa</option>
                  <option value="RU">Russia</option>
                  <option value="TR">Turkey</option>
                  <option value="IL">Israel</option>
                  <option value="AE">United Arab Emirates</option>
                  <option value="SA">Saudi Arabia</option>
                  <option value="SG">Singapore</option>
                  <option value="MY">Malaysia</option>
                  <option value="TH">Thailand</option>
                  <option value="ID">Indonesia</option>
                  <option value="PH">Philippines</option>
                  <option value="VN">Vietnam</option>
                  <option value="US">United States</option>
                  <option value="DE">Germany</option>
                  <option value="FR">France</option>
                </Select>
              </div>
              <div>
                <Label htmlFor="passportIssueDate">Issue Date</Label>
                <TextInput id="passportIssueDate" type="date" />
              </div>
              <div>
                <Label htmlFor="passportExpiryDate">Expiry Date</Label>
                <TextInput id="passportExpiryDate" type="date" />
              </div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );

  const renderDrivingLicenceSection = () => (
    <div className="space-y-6">
      <Card>
        <div className="p-6">
          <div className="flex items-center gap-3 mb-6">
            <CreditCard className="w-6 h-6 text-blue-600" />
            <h2 className="text-xl font-semibold text-gray-900">Driving Licence Processing</h2>
          </div>
          
          <div className="space-y-6">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
                <Upload className="w-10 h-10 text-gray-400 mx-auto mb-3" />
                <h4 className="font-medium text-gray-900 mb-2">Front Side</h4>
                <p className="text-sm text-gray-600 mb-3">Upload front of licence</p>
                <Button color="blue" size="sm">Choose File</Button>
              </div>
              <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
                <Upload className="w-10 h-10 text-gray-400 mx-auto mb-3" />
                <h4 className="font-medium text-gray-900 mb-2">Back Side</h4>
                <p className="text-sm text-gray-600 mb-3">Upload back of licence</p>
                <Button color="blue" size="sm">Choose File</Button>
              </div>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label htmlFor="licenceNumber">Licence Number</Label>
                <TextInput id="licenceNumber" placeholder="Enter licence number" />
              </div>
              <div>
                <Label htmlFor="licenceType">Licence Type</Label>
                <Select id="licenceType">
                  <option value="">Select type</option>
                  <option value="FULL">Full Licence</option>
                  <option value="PROVISIONAL">Provisional</option>
                  <option value="INTERNATIONAL">International</option>
                </Select>
              </div>
              <div>
                <Label htmlFor="licenceIssueDate">Issue Date</Label>
                <TextInput id="licenceIssueDate" type="date" />
              </div>
              <div>
                <Label htmlFor="licenceExpiryDate">Expiry Date</Label>
                <TextInput id="licenceExpiryDate" type="date" />
              </div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );

  const renderIdentityCardSection = () => (
    <div className="space-y-6">
      <Card>
        <div className="p-6">
          <div className="flex items-center gap-3 mb-6">
            <User className="w-6 h-6 text-blue-600" />
            <h2 className="text-xl font-semibold text-gray-900">Identity Card Processing</h2>
          </div>
          
          <div className="space-y-6">
            <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center">
              <Upload className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">Upload Identity Card</h3>
              <p className="text-gray-600 mb-4">Upload a clear photo or scan of your ID card</p>
              <Button color="blue">
                <Upload className="w-4 h-4 mr-2" />
                Choose File
              </Button>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label htmlFor="idNumber">ID Number</Label>
                <TextInput id="idNumber" placeholder="Enter ID number" />
              </div>
              <div>
                <Label htmlFor="idType">ID Type</Label>
                <Select id="idType">
                  <option value="">Select type</option>
                  <option value="NATIONAL_ID">National ID Card</option>
                  <option value="RESIDENCE_PERMIT">Residence Permit</option>
                  <option value="WORK_PERMIT">Work Permit</option>
                </Select>
              </div>
              <div>
                <Label htmlFor="idIssueDate">Issue Date</Label>
                <TextInput id="idIssueDate" type="date" />
              </div>
              <div>
                <Label htmlFor="idExpiryDate">Expiry Date</Label>
                <TextInput id="idExpiryDate" type="date" />
              </div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );

  const renderUtilityBillSection = () => (
    <div className="space-y-6">
      <Card>
        <div className="p-6">
          <div className="flex items-center gap-3 mb-6">
            <FileText className="w-6 h-6 text-blue-600" />
            <h2 className="text-xl font-semibold text-gray-900">Proof of Address</h2>
          </div>
          
          <div className="space-y-6">
            <div className="border-2 border-dashed border-gray-300 rounded-lg p-8 text-center">
              <Upload className="w-12 h-12 text-gray-400 mx-auto mb-4" />
              <h3 className="text-lg font-medium text-gray-900 mb-2">Upload Utility Bill</h3>
              <p className="text-gray-600 mb-4">Upload a recent utility bill or bank statement</p>
              <Button color="blue">
                <Upload className="w-4 h-4 mr-2" />
                Choose File
              </Button>
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <Label htmlFor="billType">Document Type</Label>
                <Select id="billType">
                  <option value="">Select type</option>
                  <option value="ELECTRICITY">Electricity Bill</option>
                  <option value="GAS">Gas Bill</option>
                  <option value="WATER">Water Bill</option>
                  <option value="COUNCIL_TAX">Council Tax</option>
                  <option value="BANK_STATEMENT">Bank Statement</option>
                  <option value="PHONE">Phone Bill</option>
                </Select>
              </div>
              <div>
                <Label htmlFor="billDate">Bill Date</Label>
                <TextInput id="billDate" type="date" />
              </div>
              <div>
                <Label htmlFor="billAddress">Address on Bill</Label>
                <TextInput id="billAddress" placeholder="Address as shown on bill" />
              </div>
              <div>
                <Label htmlFor="billPostcode">Postcode</Label>
                <TextInput id="billPostcode" placeholder="Postcode" />
              </div>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );



  // Navigate and Fill Section Renderers








  const renderAutoFillSection = () => (
    <div className="space-y-6">
      <div className="bg-white rounded-lg shadow-sm p-6">
        <div className="flex items-center space-x-3 mb-6">
          <Zap className="w-6 h-6 text-blue-600" />
          <h2 className="text-xl font-semibold text-gray-900">Auto-Fill Forms</h2>
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <Card className="p-6">
            <h3 className="font-medium mb-4">Quick Fill</h3>
            <p className="text-gray-600 mb-4">Instantly populate common form fields with your saved personal data.</p>
            <Button className="w-full">
              <Zap className="w-4 h-4 mr-2" />
              Fill Current Form
            </Button>
          </Card>
          
          <Card className="p-6">
            <h3 className="font-medium mb-4">Personal Data</h3>
            <p className="text-gray-600 mb-4">Auto-fill with email, phone, address, and personal details.</p>
            <Button outline className="w-full">
              <User className="w-4 h-4 mr-2" />
              Manage Personal Data
            </Button>
          </Card>

          <Card className="p-6">
            <h3 className="font-medium mb-4">Payment Details</h3>
            <p className="text-gray-600 mb-4">Securely auto-fill credit card and bank details.</p>
            <Button outline className="w-full">
              <CreditCard className="w-4 h-4 mr-2" />
              Manage Payment Info
            </Button>
          </Card>
        </div>

        <div className="mt-6 p-4 bg-blue-50 rounded-lg">
          <h4 className="font-medium text-blue-900 mb-2">🔒 Secure Data Storage</h4>
          <p className="text-blue-800 text-sm">
            Your personal information is encrypted and stored securely. Use auto-fill to quickly complete insurance applications, 
            web forms, and other documents while keeping your data safe.
          </p>
        </div>
      </div>
    </div>
  );

  // Web Browser Component with Real Browser Integration
  const renderWebBrowserSection = () => {
    const openRealBrowser = () => {
      const url = browserUrl || 'https://example.com';
      setIsBrowsing(true);
      
      // Open in new window with specific dimensions
      const newWindow = window.open(
        url,
        '_blank',
        'width=1200,height=800,scrollbars=yes,resizable=yes,toolbar=yes,menubar=yes,location=yes'
      );
      
      setBrowserWindow(newWindow);
      
      // Monitor if window is closed
      const checkClosed = setInterval(() => {
        if (newWindow?.closed) {
          setIsBrowsing(false);
          setBrowserWindow(null);
          clearInterval(checkClosed);
        }
      }, 1000);
    };

    const closeBrowser = () => {
      if (browserWindow) {
        browserWindow.close();
        setBrowserWindow(null);
      }
      setIsBrowsing(false);
    };

    return (
      <div className="space-y-6">
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center gap-3 mb-6">
            <Globe className="w-6 h-6 text-blue-600" />
            <h2 className="text-xl font-semibold text-gray-900">
              Web Browser
            </h2>
            {isBrowsing && (
              <Badge color="green" className="animate-pulse">
                <div className="flex items-center gap-1">
                  <div className="w-2 h-2 bg-green-500 rounded-full animate-ping"></div>
                  Browsing Active
                </div>
              </Badge>
            )}
          </div>

          {/* Browser Control Panel */}
          <div className="border border-gray-300 rounded-lg overflow-hidden bg-gray-50">
            {/* Browser Header */}
            <div className="bg-gray-200 px-4 py-3 flex items-center gap-3 border-b border-gray-300">
              <div className="flex gap-1">
                <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                <div className="w-3 h-3 bg-green-500 rounded-full"></div>
              </div>
              <div className="flex-1 mx-4">
                <TextInput
                  type="url"
                  placeholder="https://example.com"
                  value={browserUrl}
                  onChange={(e) => setBrowserUrl(e.target.value)}
                  className="text-sm"
                  onKeyPress={(e) => {
                    if (e.key === 'Enter') {
                      openRealBrowser();
                    }
                  }}
                />
              </div>
              <Button 
                size="sm" 
                color="blue"
                onClick={openRealBrowser}
                disabled={isBrowsing}
              >
                <Globe className="w-4 h-4 mr-2" />
                {isBrowsing ? 'Browsing...' : 'Navigate'}
              </Button>
            </div>
            
            {/* Browser Status */}
            <div className="p-6 text-center">
              {isBrowsing ? (
                <div className="space-y-4">
                  <div className="flex items-center justify-center gap-3">
                    <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
                    <Globe className="w-12 h-12 text-blue-600" />
                  </div>
                  <div>
                    <p className="text-lg font-medium text-gray-900">Browser Window Active</p>
                    <p className="text-sm text-gray-600">
                      Real browser opened for web navigation
                    </p>
                    <p className="text-xs text-gray-500 mt-2">
                      URL: {browserUrl || 'https://example.com'}
                    </p>
                  </div>
                  <Button 
                    color="failure" 
                    size="sm" 
                    onClick={closeBrowser}
                  >
                    <X className="w-4 h-4 mr-2" />
                    Close Browser
                  </Button>
                </div>
              ) : (
                <div className="space-y-4">
                  <Globe className="w-16 h-16 mx-auto text-gray-300" />
                  <div>
                    <p className="text-lg font-medium text-gray-900">Real Browser Integration</p>
                    <p className="text-sm text-gray-600">
                      Enter a URL above and click Navigate to open a real browser window
                    </p>
                    <p className="text-xs text-gray-500 mt-2">
                      Use the browser to navigate to insurance websites and forms
                    </p>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Floating Controls (when browsing) */}
          {isBrowsing && (
            <div className="mt-4 p-4 bg-blue-50 border border-blue-200 rounded-lg">
              <div className="flex items-center gap-2 mb-3">
                <Zap className="w-5 h-5 text-blue-600" />
                <h3 className="font-medium text-blue-900">Browser Tools</h3>
              </div>
              <div className="flex flex-wrap gap-2">
                <Button size="xs" color="blue">
                  <Wand2 className="w-3 h-3 mr-1" />
                  Auto Fill
                </Button>
                <Button size="xs" outline>
                  <User className="w-3 h-3 mr-1" />
                  Personal Data
                </Button>
                <Button size="xs" outline>
                  <CreditCard className="w-3 h-3 mr-1" />
                  Payment Info
                </Button>
              </div>
              <p className="text-xs text-blue-700 mt-2">
                Use these tools to quickly fill forms with your saved data
              </p>
            </div>
          )}
        </div>

        {/* Spider Configuration Panel */}
        <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
          <div className="flex items-center gap-3 mb-4">
            <Settings className="w-5 h-5 text-gray-600" />
            <h3 className="text-lg font-medium text-gray-900">
              {translations[session.language].spiderConfiguration}
            </h3>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <Label htmlFor="headless-mode" className="mb-2 block">
                {translations[session.language].headlessMode}
              </Label>
              <ToggleSwitch id="headless-mode" checked={false} onChange={() => {}} />
            </div>
            <div>
              <Label htmlFor="screenshots" className="mb-2 block">
                {translations[session.language].takeScreenshots}
              </Label>
              <ToggleSwitch id="screenshots" checked={true} onChange={() => {}} />
            </div>
            <div>
              <Label htmlFor="wait-load" className="mb-2 block">
                {translations[session.language].waitForLoad}
              </Label>
              <ToggleSwitch id="wait-load" checked={true} onChange={() => {}} />
            </div>
            <div>
              <Label htmlFor="javascript" className="mb-2 block">
                {translations[session.language].enableJavaScript}
              </Label>
              <ToggleSwitch id="javascript" checked={true} onChange={() => {}} />
            </div>
          </div>
        </div>
      </div>
    );
  };

  const renderContent = () => {
    const category = categories[activeCategory];
    if (!category) return renderDriversSection();

    switch (category.id) {
      // Insurance sections
      case 'drivers':
        return renderDriversSection();
      case 'vehicle':
        return renderVehicleSection();
      case 'policy':
        return renderPolicySection();
      case 'claims':
        return renderClaimsSection();
      case 'payment':
        return renderPaymentSection();
      
      // Auto-Fill section
      case 'form-filler':
        return renderAutoFillSection();
      
      // Settings
      case 'marketing':
        return renderMarketingSection();
      
      default:
        return renderDriversSection();
    }
  };

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <div className="bg-white shadow-sm border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex items-center justify-between h-16">
            <div className="flex items-center gap-3">
              <Shield className="w-8 h-8 text-blue-600" />
              <h1 className="text-xl font-semibold text-gray-900">Insurance Quote</h1>
            </div>
            <div className="flex items-center gap-2">
              <Button color="gray" onClick={() => setIsHelpOpen(true)}>
                <HelpCircle className="w-4 h-4 mr-2" />
                Help
              </Button>
              <Button color="gray" onClick={() => setShowChatbot(true)}>
                <AlertCircle className="w-4 h-4 mr-2" />
                AI Assistant
              </Button>
              <Button color="blue" onClick={() => setShowDocumentUpload(true)}>
                <Upload className="w-4 h-4 mr-2" />
                Upload Documents
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Progress Bar */}
      <div className="bg-white border-b border-gray-200">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
          <div className="flex items-center justify-between mb-2">
            <span className="text-sm font-medium text-gray-700">Progress</span>
            <span className="text-sm text-gray-500">{calculateProgress()}% Complete</span>
          </div>
          <Progress progress={calculateProgress()} color="blue" />
        </div>
      </div>

      {/* Main Layout with Left Sidebar */}
      <div className="flex h-screen bg-gray-50">
        {/* Left Sidebar */}
        <div className="w-80 bg-white shadow-lg border-r border-gray-200 flex-shrink-0 overflow-y-auto">
          <div className="p-4">
            <h2 className="text-lg font-semibold text-gray-900 mb-4">Navigation</h2>
            <div className="space-y-1">
              {menuStructure.map((section) => (
                <div key={section.id} className="space-y-1">
                  {/* Section Header */}
                  <button
                    onClick={() => {
                      if (section.id === 'personal-documents') {
                        setShowDocumentModal(true);
                      } else {
                        toggleSection(section.id);
                      }
                    }}
                    className={`w-full flex items-center gap-3 p-2 rounded-lg text-left transition-colors ${
                      expandedSections.has(section.id)
                        ? 'bg-gray-100 text-gray-900'
                        : 'hover:bg-gray-50 text-gray-700'
                    }`}
                  >
                    <section.icon className="w-4 h-4" />
                    <div className="flex-1 font-medium text-sm">{section.title}</div>
                    {expandedSections.has(section.id) ? (
                      <ChevronRight className="w-4 h-4 transform rotate-90" />
                    ) : (
                      <ChevronRight className="w-4 h-4" />
                    )}
                  </button>
                  
                  {/* Subcategories */}
                  {expandedSections.has(section.id) && (
                    <div className="ml-6 space-y-1">
                      {section.categories.map((category) => {
                        const categoryIndex = categories.findIndex(cat => cat.id === category.id);
                        return (
                          <button
                            key={category.id}
                            onClick={() => setActiveCategory(categoryIndex)}
                            className={`w-full flex items-center gap-3 p-2 rounded-lg text-left transition-colors text-sm ${
                              activeCategory === categoryIndex
                                ? 'bg-blue-50 text-blue-700 border border-blue-200'
                                : 'hover:bg-gray-50 text-gray-600'
                            }`}
                          >
                            <category.icon className="w-4 h-4" />
                            <div className="flex-1">{category.title}</div>
                            {activeCategory === categoryIndex && <ChevronRight className="w-4 h-4" />}
                          </button>
                        );
                      })}
                    </div>
                  )}
                </div>
              ))}
            </div>
          </div>
        </div>

        {/* Main Content Area */}
        <div className="flex-1 overflow-y-auto">
          <div className="p-8">
            <div className="mb-6">
              <Breadcrumb>
                <BreadcrumbItem>
                  <span className="text-gray-500">Insurance Quote</span>
                </BreadcrumbItem>
                <BreadcrumbItem>
                  <span className="text-gray-900">{categories[activeCategory].title}</span>
                </BreadcrumbItem>
              </Breadcrumb>
            </div>

            {renderContent()}

            {/* Navigation Buttons */}
            <div className="flex justify-between mt-8">
              <Button
                color="gray"
                disabled={activeCategory === 0}
                onClick={() => setActiveCategory(activeCategory - 1)}
              >
                <ChevronLeft className="w-4 h-4 mr-2" />
                Previous
              </Button>
              <Button
                color="blue"
                disabled={activeCategory === categories.length - 1}
                onClick={() => setActiveCategory(activeCategory + 1)}
              >
                Next
                <ChevronRight className="w-4 h-4 ml-2" />
              </Button>
            </div>
          </div>
        </div>
      </div>

      {/* Document Upload Modal */}
      <Modal show={showDocumentUpload} onClose={() => setShowDocumentUpload(false)} size="2xl">
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Upload Documents</h3>
            <Button color="gray" size="sm" onClick={() => setShowDocumentUpload(false)}>
              <X className="w-4 h-4" />
            </Button>
          </div>
          <div className="space-y-4">
            <div 
              className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
                isDragOver 
                  ? 'border-blue-500 bg-blue-50' 
                  : 'border-gray-300 hover:border-gray-400'
              }`}
              onDragOver={handleDragOver}
              onDragLeave={handleDragLeave}
              onDrop={handleDrop}
            >
              <Upload className="w-12 h-12 mx-auto text-gray-400 mb-4" />
              <p className="text-gray-600 mb-2">
                {isDragOver ? 'Drop files here' : 'Upload or drag and drop documents'}
              </p>
              <p className="text-sm text-gray-500 mb-4">Supports: PDF, JPG, PNG</p>
              <input
                type="file"
                multiple
                accept=".pdf,.jpg,.jpeg,.png"
                onChange={handleFileSelect}
                className="hidden"
                id="file-upload"
              />
              <label htmlFor="file-upload">
                <Button className="mt-4 cursor-pointer">
                  <Upload className="w-4 h-4 mr-2" />
                  Browse Files
                </Button>
              </label>
            </div>

            {uploadedFiles.length > 0 && (
              <div className="space-y-2">
                <h4 className="font-medium">Uploaded Files:</h4>
                <div className="space-y-2 max-h-40 overflow-y-auto">
                  {uploadedFiles.map((file, index) => (
                    <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded">
                      <span className="text-sm text-gray-700">{file.name}</span>
                      <Button 
                        color="failure" 
                        size="xs" 
                        onClick={() => removeFile(index)}
                      >
                        <X className="w-3 h-3" />
                      </Button>
                    </div>
                  ))}
                </div>
              </div>
            )}

            <div className="space-y-2">
              <h4 className="font-medium">Supported Documents:</h4>
              <ul className="text-sm text-gray-600 space-y-1">
                <li>• Driving License</li>
                <li>• Passport</li>
                <li>• Vehicle Registration Document</li>
                <li>• Insurance Certificates</li>
                <li>• Claims History</li>
              </ul>
            </div>
          </div>
          <div className="flex justify-end gap-2 mt-6">
            <Button color="gray" onClick={() => setShowDocumentUpload(false)}>Cancel</Button>
            <Button>Upload</Button>
          </div>
        </div>
      </Modal>

      {/* Help Modal */}
      <Modal show={isHelpOpen} onClose={() => setIsHelpOpen(false)} size="2xl">
        <div className="p-6">
          <div className="flex items-center justify-between mb-4">
            <h3 className="text-lg font-semibold">Help & Support</h3>
            <Button color="gray" size="sm" onClick={() => setIsHelpOpen(false)}>
              <X className="w-4 h-4" />
            </Button>
          </div>
          <div className="space-y-4">
            <div className="bg-blue-50 rounded-lg p-4">
              <p className="text-blue-800">
                Need help completing your insurance quote? Our AI assistant is here to help!
              </p>
            </div>
            <div className="space-y-2">
              <h4 className="font-medium">How to get help:</h4>
              <ul className="text-sm text-gray-600 space-y-1">
                <li>• Click the help icon next to any field</li>
                <li>• Upload documents to auto-fill information</li>
                <li>• Use the chat feature for real-time assistance</li>
              </ul>
            </div>
          </div>
          <div className="flex justify-end gap-2 mt-6">
            <Button color="gray" onClick={() => setIsHelpOpen(false)}>Cancel</Button>
            <Button>Get Help</Button>
          </div>
        </div>
      </Modal>

      {/* Notification */}
      {notification && (
        <div className="fixed top-4 right-4 z-50">
          <Alert color={notification.type === 'success' ? 'success' : 'info'} onDismiss={() => setNotification(null)}>
            {notification.message}
          </Alert>
        </div>
      )}

      {/* Floating Chatbot */}
      {showChatbot && (
        <div className="fixed bottom-4 right-4 w-96 h-96 bg-white border border-gray-300 rounded-lg shadow-lg z-50 flex flex-col">
          <div className="flex items-center justify-between p-3 bg-gray-100 border-b border-gray-300 rounded-t-lg">
            <h3 className="text-sm font-semibold">AI Assistant - Document Processing</h3>
            <Button 
              color="gray" 
              size="xs" 
              onClick={() => setShowChatbot(false)}
            >
              <X className="w-4 h-4" />
            </Button>
          </div>
          
          <div className="flex-1 overflow-y-auto p-3 space-y-3">
            {chatMessages.map((message, index) => (
              <div 
                key={index} 
                className={`flex ${message.type === 'user' ? 'justify-end' : 'justify-start'}`}
              >
                <div 
                  className={`max-w-xs p-2 rounded-lg text-sm ${
                    message.type === 'user' 
                      ? 'bg-blue-500 text-white' 
                      : 'bg-gray-100 text-gray-800'
                  }`}
                >
                  <div className="whitespace-pre-line">{message.content}</div>
                  <div className="text-xs opacity-70 mt-1">
                    {message.timestamp.toLocaleTimeString()}
                  </div>
                </div>
              </div>
            ))}
          </div>
          
          <div className="p-3 border-t border-gray-300">
            <div className="text-xs text-gray-500 text-center">
              Document processing results will appear here automatically
            </div>
          </div>
        </div>
      )}

      {/* Personal Documents Modal */}
      <Modal show={showDocumentModal} onClose={() => setShowDocumentModal(false)} size="7xl">
        <div className="p-6 bg-white rounded-lg">
          <div className="flex items-center justify-between mb-6">
            <div className="flex items-center gap-3">
              <Upload className="w-6 h-6 text-blue-600" />
              <h3 className="text-lg font-semibold text-gray-900">{translations[session.language].personalDocuments}</h3>
            </div>
            <Button color="gray" size="sm" onClick={() => setShowDocumentModal(false)}>
              <X className="w-4 h-4" />
            </Button>
          </div>
          <div className="space-y-6">
            {/* Document Upload Cards */}
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              {[
                { id: 'passport', title: translations[session.language].passport, icon: FileText, description: translations[session.language].passportUpload },
                { id: 'driving-licence', title: translations[session.language].drivingLicence, icon: CreditCard, description: translations[session.language].drivingLicenceUpload },
                { id: 'identity-card', title: translations[session.language].identityCard, icon: User, description: translations[session.language].identityCardUpload },
                { id: 'utility-bill', title: translations[session.language].utilityBill, icon: FileText, description: translations[session.language].utilityBillUpload },
                { id: 'vehicle-registration', title: translations[session.language].vehicleRegistration, icon: Car, description: translations[session.language].vehicleRegistrationUpload },
                { id: 'bank-statement', title: translations[session.language].bankStatement, icon: CreditCard, description: translations[session.language].bankStatementUpload },
                { id: 'medical-certificate', title: translations[session.language].medicalCertificate, icon: FileText, description: translations[session.language].medicalCertificateUpload },
                { id: 'insurance-quote', title: translations[session.language].insuranceQuote, icon: Shield, description: translations[session.language].insuranceQuoteUpload },
                { id: 'insurance-policy', title: translations[session.language].insurancePolicy, icon: Shield, description: translations[session.language].insurancePolicyUpload }
              ].map((docType) => (
                <Card 
                  key={docType.id} 
                  className={`transition-all hover:shadow-lg bg-white border-2 border-dashed ${
                    isDragOverDocument && selectedDocumentType === docType.id
                      ? 'border-blue-500 bg-blue-50' 
                      : isProcessing && selectedDocumentType === docType.id
                        ? 'border-green-400 bg-green-50'
                        : 'border-gray-300 hover:border-blue-400'
                  }`}
                  onDragOver={(e) => {
                    e.preventDefault();
                    setSelectedDocumentType(docType.id);
                    setIsDragOverDocument(true);
                  }}
                  onDragLeave={(e) => {
                    e.preventDefault();
                    setIsDragOverDocument(false);
                  }}
                  onDrop={(e) => {
                    e.preventDefault();
                    setSelectedDocumentType(docType.id);
                    setIsDragOverDocument(false);
                    handleDocumentDrop(e, docType.id);
                  }}
                >
                  <div className="p-4 text-center">
                    <docType.icon className="w-8 h-8 text-blue-600 mx-auto mb-3" />
                    <h3 className="font-medium text-gray-900 mb-2">{docType.title}</h3>
                    <p className="text-xs text-gray-600 mb-3">{docType.description}</p>
                    
                    {/* Upload Area */}
                    <div className="border-t pt-3">
                      <Upload className="w-6 h-6 text-gray-400 mx-auto mb-2" />
                      <p className="text-xs text-gray-500 mb-2">
                        {isDragOverDocument && selectedDocumentType === docType.id 
                          ? 'Drop files here' 
                          : 'Drag & drop or click'
                        }
                      </p>
                      <input
                        type="file"
                        multiple
                        accept=".pdf,.jpg,.jpeg,.png"
                        onChange={(e) => {
                          setSelectedDocumentType(docType.id);
                          handleDocumentFileSelect(e, docType.id);
                        }}
                        className="hidden"
                        id={`upload-${docType.id}`}
                      />
                      <Button 
                        size="xs" 
                        color="blue"
                        onClick={(e: React.MouseEvent) => {
                          e.stopPropagation();
                          setSelectedDocumentType(docType.id);
                          document.getElementById(`upload-${docType.id}`)?.click();
                        }}
                        disabled={isProcessing && selectedDocumentType === docType.id}
                      >
                        {isProcessing && selectedDocumentType === docType.id ? (
                          <>
                            <div className="animate-spin rounded-full h-3 w-3 border-b-2 border-white mr-1"></div>
                            Processing...
                          </>
                        ) : (
                          <>
                            <Upload className="w-3 h-3 mr-1" />
                            Upload
                          </>
                        )}
                      </Button>
                    </div>
                    
                    {/* File Count */}
                    {documentFiles[docType.id]?.length > 0 && (
                      <div className="mt-2">
                        <Badge color="green" size="xs">
                          {documentFiles[docType.id].length} file{documentFiles[docType.id].length > 1 ? 's' : ''} uploaded
                        </Badge>
                      </div>
                    )}
                  </div>
                </Card>
              ))}
            </div>

            {/* Document Upload Area */}
            {selectedDocumentType && (
              <div className="border-t pt-6">
                <div className="flex items-center gap-3 mb-4">
                  <FileText className="w-5 h-5 text-blue-600" />
                  <h3 className="text-lg font-semibold text-gray-900">
                    Upload {selectedDocumentType.charAt(0).toUpperCase() + selectedDocumentType.slice(1).replace('-', ' ')}
                  </h3>
                </div>
                
                {/* Upload Areas */}
                {selectedDocumentType === 'driving-licence' ? (
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                    {/* Front Side Upload */}
                    <div>
                      <div 
                        className={`border-2 border-dashed rounded-lg p-6 text-center transition-colors bg-white ${
                          isDragOverDocument ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-gray-400 hover:bg-gray-50'
                        }`}
                        onDragOver={handleDocumentDragOver}
                        onDragLeave={handleDocumentDragLeave}
                        onDrop={(e) => handleDocumentDrop(e, 'front')}
                      >
                        <Upload className="w-10 h-10 text-gray-400 mx-auto mb-3" />
                        <h4 className="font-medium text-gray-900 mb-2">{translations[session.language].frontSide}</h4>
                        <p className="text-sm text-gray-600 mb-3">
                          {isDragOverDocument ? translations[session.language].dropFilesHere : translations[session.language].dragDropFiles}
                        </p>
                        <input
                          type="file"
                          multiple
                          accept=".pdf,.jpg,.jpeg,.png"
                          onChange={(e) => handleDocumentFileSelect(e, 'front')}
                          className="hidden"
                          id="front-upload"
                        />
                        <label htmlFor="front-upload">
                          <Button color="blue" size="sm" className="cursor-pointer">
                            {isProcessing ? translations[session.language].processing : translations[session.language].chooseFile}
                          </Button>
                        </label>
                      </div>
                      {documentFiles[`${selectedDocumentType}_front`]?.length > 0 && (
                        <div className="mt-3 space-y-2">
                          {documentFiles[`${selectedDocumentType}_front`].map((file, index) => (
                            <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded text-sm">
                              <span className="text-gray-700">{file.name}</span>
                              <Button 
                                color="failure" 
                                size="xs" 
                                onClick={() => removeDocumentFile(index, 'front')}
                              >
                                <X className="w-3 h-3" />
                              </Button>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* Back Side Upload */}
                    <div>
                      <div 
                        className={`border-2 border-dashed rounded-lg p-6 text-center transition-colors bg-white ${
                          isDragOverDocument ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-gray-400 hover:bg-gray-50'
                        }`}
                        onDragOver={handleDocumentDragOver}
                        onDragLeave={handleDocumentDragLeave}
                        onDrop={(e) => handleDocumentDrop(e, 'back')}
                      >
                        <Upload className="w-10 h-10 text-gray-400 mx-auto mb-3" />
                        <h4 className="font-medium text-gray-900 mb-2">{translations[session.language].backSide}</h4>
                        <p className="text-sm text-gray-600 mb-3">
                          {isDragOverDocument ? translations[session.language].dropFilesHere : translations[session.language].dragDropFiles}
                        </p>
                        <input
                          type="file"
                          multiple
                          accept=".pdf,.jpg,.jpeg,.png"
                          onChange={(e) => handleDocumentFileSelect(e, 'back')}
                          className="hidden"
                          id="back-upload"
                        />
                        <label htmlFor="back-upload">
                          <Button color="blue" size="sm" className="cursor-pointer">
                            {isProcessing ? 'Processing...' : 'Choose File'}
                          </Button>
                        </label>
                      </div>
                      {documentFiles[`${selectedDocumentType}_back`]?.length > 0 && (
                        <div className="mt-3 space-y-2">
                          {documentFiles[`${selectedDocumentType}_back`].map((file, index) => (
                            <div key={index} className="flex items-center justify-between p-2 bg-gray-50 rounded text-sm">
                              <span className="text-gray-700">{file.name}</span>
                              <Button 
                                color="failure" 
                                size="xs" 
                                onClick={() => removeDocumentFile(index, 'back')}
                              >
                                <X className="w-3 h-3" />
                              </Button>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                ) : selectedDocumentType === 'passport' ? (
                  <div className="space-y-6">
                    {/* Passport Count Selector */}
                    <div className="flex items-center gap-4 mb-4">
                      <Label>Number of Passports:</Label>
                      <Select value={passportCount.toString()} onChange={(e) => setPassportCount(parseInt(e.target.value))}>
                        <option value="1">1 Passport</option>
                        <option value="2">2 Passports</option>
                        <option value="3">3 Passports</option>
                      </Select>
                    </div>
                    
                    {Array.from({length: passportCount}, (_, i) => (
                      <div key={i} className="border border-gray-200 rounded-lg p-4">
                        <h4 className="font-medium text-gray-900 mb-4">Passport {i + 1}</h4>
                        <div 
                          className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors bg-white ${
                            isDragOverDocument ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-gray-400 hover:bg-gray-50'
                          }`}
                          onDragOver={handleDocumentDragOver}
                          onDragLeave={handleDocumentDragLeave}
                          onDrop={(e) => handleDocumentDrop(e, `passport_${i + 1}`)}
                        >
                          <Upload className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                          <h4 className="text-lg font-medium text-gray-900 mb-2">Upload Passport {i + 1}</h4>
                          <p className="text-gray-600 mb-4">
                            {isDragOverDocument ? 'Drop files here' : 'Drag and drop your passport here, or click to browse'}
                          </p>
                          <input
                            type="file"
                            multiple
                            accept=".pdf,.jpg,.jpeg,.png"
                            onChange={(e) => handleDocumentFileSelect(e, `passport_${i + 1}`)}
                            className="hidden"
                            id={`passport-upload-${i + 1}`}
                          />
                          <label htmlFor={`passport-upload-${i + 1}`}>
                            <Button color="blue" className="cursor-pointer">
                              <Upload className="w-4 h-4 mr-2" />
                              {isProcessing ? 'Processing...' : 'Choose File'}
                            </Button>
                          </label>
                        </div>
                        {documentFiles[`${selectedDocumentType}_passport_${i + 1}`]?.length > 0 && (
                          <div className="mt-4 space-y-2">
                            <h5 className="font-medium">Uploaded Files:</h5>
                            {documentFiles[`${selectedDocumentType}_passport_${i + 1}`].map((file, fileIndex) => (
                              <div key={fileIndex} className="flex items-center justify-between p-2 bg-gray-50 rounded text-sm">
                                <span className="text-gray-700">{file.name}</span>
                                <Button 
                                  color="failure" 
                                  size="xs" 
                                  onClick={() => removeDocumentFile(fileIndex, `passport_${i + 1}`)}
                                >
                                  <X className="w-3 h-3" />
                                </Button>
                              </div>
                            ))}
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                ) : (
                  <div 
                    className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors bg-white ${
                      isDragOverDocument ? 'border-blue-500 bg-blue-50' : 'border-gray-300 hover:border-gray-400 hover:bg-gray-50'
                    }`}
                    onDragOver={handleDocumentDragOver}
                    onDragLeave={handleDocumentDragLeave}
                    onDrop={(e) => handleDocumentDrop(e, 'main')}
                  >
                    <Upload className="w-12 h-12 text-gray-400 mx-auto mb-4" />
                    <h4 className="text-lg font-medium text-gray-900 mb-2">
                      Upload {selectedDocumentType.charAt(0).toUpperCase() + selectedDocumentType.slice(1).replace('-', ' ')}
                    </h4>
                    <p className="text-gray-600 mb-4">
                      {isDragOverDocument ? 'Drop files here' : 'Drag and drop your file here, or click to browse'}
                    </p>
                    <input
                      type="file"
                      multiple
                      accept=".pdf,.jpg,.jpeg,.png"
                      onChange={(e) => handleDocumentFileSelect(e, 'main')}
                      className="hidden"
                      id="main-upload"
                    />
                    <label htmlFor="main-upload">
                      <Button color="blue" className="cursor-pointer">
                        <Upload className="w-4 h-4 mr-2" />
                        {isProcessing ? 'Processing...' : 'Choose File'}
                      </Button>
                    </label>
                  </div>
                )}

                {/* Document-specific fields */}
                <div className="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
                  {selectedDocumentType === 'passport' && (
                    <div className="space-y-6">
                      {Array.from({length: passportCount}, (_, i) => (
                        <div key={i} className="border border-gray-200 rounded-lg p-4">
                          <h4 className="font-medium text-gray-900 mb-4">Passport {i + 1} Details</h4>
                          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                            <div>
                              <Label htmlFor={`passportNumber_${i + 1}`}>Passport Number</Label>
                              <TextInput id={`passportNumber_${i + 1}`} placeholder="Enter passport number" />
                            </div>
                            <div>
                              <Label htmlFor={`passportCountry_${i + 1}`}>Issuing Country</Label>
                              <Select id={`passportCountry_${i + 1}`}>
                                <option value="">{translations[session.language as keyof typeof translations].selectCountry}</option>
                                <option value="GB">United Kingdom</option>
                                <option value="US">United States</option>
                                <option value="DE">Germany</option>
                                <option value="FR">France</option>
                                <option value="IT">Italy</option>
                                <option value="ES">Spain</option>
                                <option value="NL">Netherlands</option>
                                <option value="BE">Belgium</option>
                                <option value="AT">Austria</option>
                                <option value="CH">Switzerland</option>
                                <option value="IE">Ireland</option>
                                <option value="PT">Portugal</option>
                                <option value="GR">Greece</option>
                                <option value="PL">Poland</option>
                                <option value="CZ">Czech Republic</option>
                                <option value="HU">Hungary</option>
                                <option value="SK">Slovakia</option>
                                <option value="SI">Slovenia</option>
                                <option value="HR">Croatia</option>
                                <option value="RO">Romania</option>
                                <option value="BG">Bulgaria</option>
                                <option value="LT">Lithuania</option>
                                <option value="LV">Latvia</option>
                                <option value="EE">Estonia</option>
                                <option value="FI">Finland</option>
                                <option value="SE">Sweden</option>
                                <option value="DK">Denmark</option>
                                <option value="NO">Norway</option>
                                <option value="IS">Iceland</option>
                                <option value="LU">Luxembourg</option>
                                <option value="MT">Malta</option>
                                <option value="CY">Cyprus</option>
                                <option value="CA">Canada</option>
                                <option value="AU">Australia</option>
                                <option value="NZ">New Zealand</option>
                                <option value="JP">Japan</option>
                                <option value="KR">South Korea</option>
                                <option value="CN">China</option>
                                <option value="IN">India</option>
                                <option value="BR">Brazil</option>
                                <option value="MX">Mexico</option>
                                <option value="AR">Argentina</option>
                                <option value="CL">Chile</option>
                                <option value="ZA">South Africa</option>
                                <option value="RU">Russia</option>
                                <option value="TR">Turkey</option>
                                <option value="IL">Israel</option>
                                <option value="AE">United Arab Emirates</option>
                                <option value="SA">Saudi Arabia</option>
                                <option value="SG">Singapore</option>
                                <option value="MY">Malaysia</option>
                                <option value="TH">Thailand</option>
                                <option value="ID">Indonesia</option>
                                <option value="PH">Philippines</option>
                                <option value="VN">Vietnam</option>
                                <option value="US">United States</option>
                                <option value="DE">Germany</option>
                                <option value="FR">France</option>
                                <option value="CA">Canada</option>
                                <option value="AU">Australia</option>
                                <option value="IE">Ireland</option>
                                <option value="NZ">New Zealand</option>
                              </Select>
                            </div>
                            <div>
                              <Label htmlFor={`passportIssueDate_${i + 1}`}>Issue Date</Label>
                              <TextInput id={`passportIssueDate_${i + 1}`} type="date" />
                            </div>
                            <div>
                              <Label htmlFor={`passportExpiryDate_${i + 1}`}>Expiry Date</Label>
                              <TextInput id={`passportExpiryDate_${i + 1}`} type="date" />
                            </div>
                            <div>
                              <Label htmlFor={`passportGivenNames_${i + 1}`}>Given Names</Label>
                              <TextInput id={`passportGivenNames_${i + 1}`} placeholder="As shown on passport" />
                            </div>
                            <div>
                              <Label htmlFor={`passportSurname_${i + 1}`}>Surname</Label>
                              <TextInput id={`passportSurname_${i + 1}`} placeholder="As shown on passport" />
                            </div>
                            <div>
                              <Label htmlFor={`passportDateOfBirth_${i + 1}`}>Date of Birth</Label>
                              <TextInput id={`passportDateOfBirth_${i + 1}`} type="date" />
                            </div>
                            <div>
                              <Label htmlFor={`passportGender_${i + 1}`}>Gender</Label>
                              <Select id={`passportGender_${i + 1}`}>
                                <option value="">Select gender</option>
                                <option value="M">Male</option>
                                <option value="F">Female</option>
                                <option value="X">Other</option>
                              </Select>
                            </div>
                          </div>
                        </div>
                      ))}
                    </div>
                  )}

                  {selectedDocumentType === 'driving-licence' && (
                    <>
                      <div>
                        <Label htmlFor="licenceNumber">Licence Number</Label>
                        <TextInput id="licenceNumber" placeholder="Enter licence number" />
                      </div>
                      <div>
                        <Label htmlFor="licenceType">Licence Type</Label>
                        <Select id="licenceType">
                          <option value="">Select type</option>
                          <option value="FULL">Full Licence</option>
                          <option value="PROVISIONAL">Provisional</option>
                          <option value="INTERNATIONAL">International</option>
                        </Select>
                      </div>
                      <div>
                        <Label htmlFor="licenceIssueDate">Issue Date</Label>
                        <TextInput id="licenceIssueDate" type="date" />
                      </div>
                      <div>
                        <Label htmlFor="licenceExpiryDate">Expiry Date</Label>
                        <TextInput id="licenceExpiryDate" type="date" />
                      </div>
                    </>
                  )}

                  {selectedDocumentType === 'vehicle-registration' && (
                    <>
                      <div>
                        <Label htmlFor="registrationNumber">Registration Number</Label>
                        <TextInput id="registrationNumber" placeholder="Enter registration number" />
                      </div>
                      <div>
                        <Label htmlFor="vehicleMake">Make</Label>
                        <TextInput id="vehicleMake" placeholder="Vehicle make" />
                      </div>
                      <div>
                        <Label htmlFor="vehicleModel">Model</Label>
                        <TextInput id="vehicleModel" placeholder="Vehicle model" />
                      </div>
                      <div>
                        <Label htmlFor="yearOfManufacture">Year</Label>
                        <TextInput id="yearOfManufacture" type="number" placeholder="Year" />
                      </div>
                    </>
                  )}

                  {selectedDocumentType === 'bank-statement' && (
                    <>
                      <div>
                        <Label htmlFor="bankName">Bank Name</Label>
                        <Select id="bankName">
                          <option value="">Select bank</option>
                          <option value="HSBC">HSBC</option>
                          <option value="BARCLAYS">Barclays</option>
                          <option value="LLOYDS">Lloyds</option>
                          <option value="NATWEST">NatWest</option>
                        </Select>
                      </div>
                      <div>
                        <Label htmlFor="statementDate">Statement Date</Label>
                        <TextInput id="statementDate" type="date" />
                      </div>
                    </>
                  )}

                  {selectedDocumentType === 'insurance-policy' && (
                    <>
                      <div>
                        <Label htmlFor="policyProvider">Insurance Provider</Label>
                        <TextInput id="policyProvider" placeholder="Provider name" />
                      </div>
                      <div>
                        <Label htmlFor="policyNumber">Policy Number</Label>
                        <TextInput id="policyNumber" placeholder="Policy number" />
                      </div>
                      <div>
                        <Label htmlFor="ncdYears">No Claims Discount</Label>
                        <Select id="ncdYears">
                          <option value="">Select NCD years</option>
                          <option value="0">0 years</option>
                          <option value="1">1 year</option>
                          <option value="2">2 years</option>
                          <option value="3">3 years</option>
                          <option value="4">4 years</option>
                          <option value="5">5+ years</option>
                        </Select>
                      </div>
                    </>
                  )}
                </div>

                {/* Action Buttons */}
                <div className="flex justify-between mt-6">
                  <Button color="gray" onClick={() => setSelectedDocumentType('')}>
                    Back to Selection
                  </Button>
                  <div className="flex gap-3">
                    <Button color="gray" onClick={() => setShowDocumentModal(false)}>
                      Cancel
                    </Button>
                    <div className="text-sm text-gray-600 flex items-center">
                      <CheckCircle className="w-4 h-4 mr-2 text-green-600" />
                      Documents auto-process on upload
                    </div>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default App;
