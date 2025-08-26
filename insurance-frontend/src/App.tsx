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
  Upload
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
  licenceIssueDate?: string;
  licenceExpiryDate?: string;
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
}

interface Vehicle {
  registration: string;
  make: string;
  model: string;
  year: number;
  mileage: number;
  value: number;
  overnightLocation: string;
  hasModifications: boolean;
  modifications: Record<string, any>;
}

interface Policy {
  coverType: string;
  startDate: string;
  voluntaryExcess: number;
  ncdYears: number;
  protectNCD: boolean;
}

interface ClaimsHistory {
  hasClaims: boolean;
  claims: any[];
  hasConvictions: boolean;
  convictions: any[];
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

const categories = [
  { id: 'drivers', title: 'Driver Details', icon: User, order: 1, description: 'Information about all drivers' },
  { id: 'vehicle', title: 'Vehicle Details', icon: Car, order: 2, description: 'Vehicle information and modifications' },
  { id: 'policy', title: 'Policy Details', icon: FileText, order: 3, description: 'Coverage and policy options' },
  { id: 'claims', title: 'Claims History', icon: Shield, order: 4, description: 'Previous claims and convictions' },
  { id: 'payment', title: 'Payment & Extras', icon: CreditCard, order: 5, description: 'Payment method and additional cover' },
  { id: 'marketing', title: 'Marketing Preferences', icon: Settings, order: 6, description: 'Communication preferences' }
];

function App() {
  const [activeCategory, setActiveCategory] = useState(0);
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
        relationship: '',
        sameAddress: true
      }
    ],
    vehicle: {
      registration: '',
      make: '',
      model: '',
      year: new Date().getFullYear(),
      mileage: 0,
      value: 0,
      overnightLocation: '',
      hasModifications: false,
      modifications: {}
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
      convictions: []
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
      relationship: '',
      sameAddress: false
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
        relationship: '',
        sameAddress: true,
        // Passport-specific fields
        passportNumber: passportData.PassportNumber,
        passportIssueDate: passportData.DateOfIssue,
        passportExpiryDate: passportData.DateOfExpiry,
        passportAuthority: passportData.Authority,
        placeOfBirth: passportData.PlaceOfBirth,
        gender: passportData.Gender,
        nationality: passportData.Nationality
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
                First Name
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
                Last Name
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
                Date of Birth
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
                Email Address
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
                Phone Number
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
                Address
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
                Postcode
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
                Licence Type
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
                Licence Number
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
                Years Licence Held
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
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="form-group">
                <Label className="form-label">Engine Modifications</Label>
                <Select className="form-select">
                  <option value="">Select modification</option>
                  <option value="none">None</option>
                  <option value="chip_tuning">Chip Tuning</option>
                  <option value="turbo">Turbo/Supercharger</option>
                  <option value="exhaust">Exhaust System</option>
                </Select>
              </div>

              <div className="form-group">
                <Label className="form-label">Body Modifications</Label>
                <Select className="form-select">
                  <option value="">Select modification</option>
                  <option value="none">None</option>
                  <option value="spoiler">Spoiler</option>
                  <option value="body_kit">Body Kit</option>
                  <option value="wheels">Alloy Wheels</option>
                </Select>
              </div>
            </div>
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
      </Card>
    </div>
  );

  const renderClaimsSection = () => (
    <div className="space-y-6">
      <h2 className="text-xl font-semibold">Claims History</h2>

      <Card className="form-section">
        <div className="space-y-4">
          <div className="form-group">
            <Label className="form-label">
              Previous Claims (Last 5 Years)
              {getStatusBadge(false, true)}
            </Label>
            <ToggleSwitch
              checked={session.claims.hasClaims}
              onChange={(checked) => updateClaims('hasClaims', checked)}
            />
          </div>

          <div className="form-group">
            <Label className="form-label">
              Convictions (Last 5 Years)
              {getStatusBadge(false, true)}
            </Label>
            <ToggleSwitch
              checked={session.claims.hasConvictions}
              onChange={(checked) => updateClaims('hasConvictions', checked)}
            />
          </div>

          {(session.claims.hasClaims || session.claims.hasConvictions) && (
            <Alert color="warning" className="mt-4">
              <div className="flex items-center gap-2">
                <AlertCircle className="w-5 h-5" />
                <span>Please provide details of any claims or convictions when prompted.</span>
              </div>
            </Alert>
          )}
        </div>
      </Card>
    </div>
  );

  const renderPaymentSection = () => (
    <div className="space-y-6">
      <h2 className="text-xl font-semibold">Payment & Extras</h2>

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

  const renderContent = () => {
    switch (activeCategory) {
      case 0:
        return renderDriversSection();
      case 1:
        return renderVehicleSection();
      case 2:
        return renderPolicySection();
      case 3:
        return renderClaimsSection();
      case 4:
        return renderPaymentSection();
      case 5:
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
                <HelpCircle className="w-4 h-4 mr-2" />
                AI Assistant
              </Button>
              <Button color="gray" onClick={() => setShowDocumentUpload(true)}>
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

      {/* Main Content */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex gap-8">
          {/* Sidebar */}
          <div className="w-80 flex-shrink-0">
            <div className="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
              <h2 className="text-lg font-semibold text-gray-900 mb-4">Quote Sections</h2>
              <div className="space-y-2">
                {categories.map((category, index) => (
                  <button
                    key={category.id}
                    onClick={() => setActiveCategory(index)}
                    className={`w-full flex items-center gap-3 p-3 rounded-lg text-left transition-colors ${
                      activeCategory === index
                        ? 'bg-blue-50 text-blue-700 border border-blue-200'
                        : 'hover:bg-gray-50 text-gray-700'
                    }`}
                  >
                    <category.icon className="w-5 h-5" />
                    <div className="flex-1">
                      <div className="font-medium">{category.title}</div>
                      <div className="text-sm text-gray-500">{category.description}</div>
                    </div>
                    {activeCategory === index && <ChevronRight className="w-4 h-4" />}
                  </button>
                ))}
              </div>
            </div>
          </div>

          {/* Main Content Area */}
          <div className="flex-1">
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
    </div>
  );
}

export default App;
