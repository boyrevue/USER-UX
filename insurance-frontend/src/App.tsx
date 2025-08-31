import React, { useState, useEffect } from 'react';
import { Card, Button, Progress, Alert } from 'flowbite-react';
import { ChevronLeft, ChevronRight, CheckCircle, AlertCircle, User, Car, FileText, Shield, Upload } from 'lucide-react';

// Import universal components
import UniversalForm from './components/forms/UniversalForm';
import GlobalAssistant from './components/forms/GlobalAssistant';
import DocumentMatrix from './components/forms/DocumentMatrix';
import DocumentManager from './components/documents/DocumentManager';
import Navigation from './components/layout/Navigation';

// Import services
import { ApiService } from './services/api';
import ocrService from './services/ocrService';
import { validateBirthDate, validateLicenceDate, validateHistoricalDate } from './services/validation';

// Import types
import { 
  AppState, 
  Session, 
  Driver, 
  Vehicle, 
  Claim, 
  Accident, 
  Document,
  ValidationError 
} from './types/index';

const App: React.FC = () => {
  // Main application state
  const [appState, setAppState] = useState<AppState>({
    currentStep: 0,
    session: {
      id: `session_${Date.now()}`,
      language: 'en',
      drivers: [{
        id: 'driver_main',
        classification: 'MAIN',
        firstName: '',
        lastName: '',
        dateOfBirth: '',
        email: '',
        phone: '',
        licenceNumber: '',
        licenceIssueDate: '',
        licenceExpiryDate: '',
        licenceValidUntil: '',
        convictions: []
      }],
      vehicles: [{
        id: 'vehicle_1',
        registration: '',
        make: '',
        model: '',
        year: new Date().getFullYear(),
        engineSize: '',
        fuelType: '',
        transmission: '',
        estimatedValue: 0,
        modifications: []
      }],
      claims: {
        claims: [],
        accidents: []
      },
      policy: {
        startDate: '',
        coverType: '',
        excess: 0,
        ncdYears: 0,
        ncdProtected: false
      },
      documents: []
    },
    loading: false,
    error: null,
    ontology: null,
    validationErrors: []
  });

  const [isProcessing, setIsProcessing] = useState(false);
  const [showAssistant, setShowAssistant] = useState(true);

  const steps = [
    'Driver Details',
    'Vehicle Details', 
    'Claims History',
    'Policy Details',
    'Documents',
    'Settings',
    'Review & Submit'
  ];

  // Load ontology on component mount
  useEffect(() => {
    const loadOntology = async () => {
      try {
        setAppState(prev => ({ ...prev, loading: true }));
        const ontologyData = await ApiService.getOntology();
        setAppState(prev => ({ 
          ...prev, 
          ontology: ontologyData,
          loading: false 
        }));
      } catch (error) {
        console.error('Failed to load ontology:', error);
        setAppState(prev => ({ 
          ...prev, 
          error: 'Failed to load form configuration',
          loading: false 
        }));
      }
    };

    loadOntology();
  }, []);

  // Driver management functions
  const updateDriver = (index: number, field: string, value: any) => {
    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        drivers: prev.session.drivers.map((driver, idx) => 
          idx === index ? { ...driver, [field]: value } : driver
        )
      }
    }));
  };

  const addDriver = () => {
    const newDriver: Driver = {
      id: `driver_${Date.now()}`,
      classification: 'NAMED',
      firstName: '',
      lastName: '',
      dateOfBirth: '',
      email: '',
      phone: '',
      licenceNumber: '',
      licenceIssueDate: '',
      licenceExpiryDate: '',
      licenceValidUntil: '',
      convictions: []
    };

    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        drivers: [...prev.session.drivers, newDriver]
      }
    }));
  };

  const removeDriver = (index: number) => {
    if (appState.session.drivers[index].classification !== 'MAIN') {
      setAppState(prev => ({
        ...prev,
        session: {
          ...prev.session,
          drivers: prev.session.drivers.filter((_, idx) => idx !== index)
        }
      }));
    }
  };

  // Vehicle management functions
  const updateVehicle = (index: number, field: string, value: any) => {
    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        vehicles: prev.session.vehicles.map((vehicle, idx) => 
          idx === index ? { ...vehicle, [field]: value } : vehicle
        )
      }
    }));
  };

  const addVehicle = () => {
    const newVehicle: Vehicle = {
      id: `vehicle_${Date.now()}`,
      registration: '',
      make: '',
      model: '',
      year: new Date().getFullYear(),
      engineSize: '',
      fuelType: '',
      transmission: '',
      estimatedValue: 0,
      modifications: []
    };

    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        vehicles: [...prev.session.vehicles, newVehicle]
      }
    }));
  };

  const removeVehicle = (index: number) => {
    if (appState.session.vehicles.length > 1) {
      setAppState(prev => ({
        ...prev,
        session: {
          ...prev.session,
          vehicles: prev.session.vehicles.filter((_, idx) => idx !== index)
        }
      }));
    }
  };

  // Claims management functions
  const updateClaim = (index: number, field: string, value: any) => {
    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        claims: {
          ...prev.session.claims,
          claims: prev.session.claims.claims.map((claim, idx) => 
            idx === index ? { ...claim, [field]: value } : claim
          )
        }
      }
    }));
  };

  const updateAccident = (index: number, field: string, value: any) => {
    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        claims: {
          ...prev.session.claims,
          accidents: prev.session.claims.accidents.map((accident, idx) => 
            idx === index ? { ...accident, [field]: value } : accident
          )
        }
      }
    }));
  };

  const addClaim = () => {
    const newClaim: Claim = {
      id: `claim_${Date.now()}`,
      date: '',
      type: '',
      amount: 0,
      description: '',
      settled: false
    };

    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        claims: {
          ...prev.session.claims,
          claims: [...prev.session.claims.claims, newClaim]
        }
      }
    }));
  };

  const addAccident = () => {
    const newAccident: Accident = {
      id: `accident_${Date.now()}`,
      date: '',
      type: '',
      description: '',
      estimatedCost: 0,
      faultClaim: false
    };

    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        claims: {
          ...prev.session.claims,
          accidents: [...prev.session.claims.accidents, newAccident]
        }
      }
    }));
  };

  const removeClaim = (index: number) => {
    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        claims: {
          ...prev.session.claims,
          claims: prev.session.claims.claims.filter((_, idx) => idx !== index)
        }
      }
    }));
  };

  const removeAccident = (index: number) => {
    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        claims: {
          ...prev.session.claims,
          accidents: prev.session.claims.accidents.filter((_, idx) => idx !== index)
        }
      }
    }));
  };

  // Document management functions
  const handleFileUpload = async (files: FileList, uploadType: string, selectedDocumentType?: string) => {
    try {
      setIsProcessing(true);
      const result = await ApiService.processDocument(files, uploadType, selectedDocumentType);
      
      // Create document record
      const newDocument: Document = {
        id: `doc_${Date.now()}`,
        type: selectedDocumentType as any || 'other',
        name: files[0].name,
        size: files[0].size,
        uploadedAt: new Date().toISOString(),
        processed: true,
        fieldName: selectedDocumentType || uploadType,
        extractedData: result.extractedFields,
        confidence: result.confidence,
        imagePaths: result.images
      };

      setAppState(prev => ({
        ...prev,
        session: {
          ...prev.session,
          documents: [...prev.session.documents, newDocument]
        }
      }));

      // Auto-populate form fields if extraction was successful
      if (result.extractedFields) {
        populateFormFromExtraction(result.extractedFields, uploadType);
      }

    } catch (error) {
      console.error('Document processing failed:', error);
      setAppState(prev => ({
        ...prev,
        error: `Document processing failed: ${error}`
      }));
    } finally {
      setIsProcessing(false);
    }
  };

  // New document manager handlers
  const handleFileDownload = (documentType: string) => {
    // Find the document in the session
    const document = appState.session.documents.find(doc => doc.fieldName === documentType);
    if (document) {
      // In a real implementation, this would download from the server
      console.log(`Downloading document: ${document.name}`);
      // For now, just show an alert
      alert(`Download functionality for ${document.name} would be implemented here`);
    }
  };

  const handleOCRProcessing = async (documentType: string, file: File) => {
    try {
      setIsProcessing(true);
      const ocrResult = await ocrService.processDocument(documentType, file);
      
      // Update the session with extracted data
      setAppState(prev => ({
        ...prev,
        session: {
          ...prev.session,
          // Map OCR results to appropriate session fields
          ...ocrResult
        }
      }));
      
      return ocrResult;
    } catch (error) {
      console.error('OCR processing failed:', error);
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
      setAppState(prev => ({
        ...prev,
        error: `OCR processing failed: ${errorMessage}`
      }));
      throw error;
    } finally {
      setIsProcessing(false);
    }
  };

  const populateFormFromExtraction = (extractedData: Record<string, any>, uploadType: string) => {
    if (uploadType === 'passport' && extractedData.passportNumber) {
      // Populate driver details from passport
      updateDriver(0, 'firstName', extractedData.givenNames || '');
      updateDriver(0, 'lastName', extractedData.surname || '');
      updateDriver(0, 'dateOfBirth', extractedData.dateOfBirth || '');
    }
    // Add more extraction mappings as needed
  };

  const removeDocument = (id: string) => {
    setAppState(prev => ({
      ...prev,
      session: {
        ...prev.session,
        documents: prev.session.documents.filter(doc => doc.id !== id)
      }
    }));
  };

  // Navigation functions
  const canNavigate = (step: number): boolean => {
    // Add validation logic here
    return true;
  };

  const getCompletedSteps = (): number[] => {
    const completed: number[] = [];
    // Add completion logic here
    return completed;
  };

  const nextStep = () => {
    if (appState.currentStep < steps.length - 1) {
      setAppState(prev => ({ ...prev, currentStep: prev.currentStep + 1 }));
    }
  };

  const prevStep = () => {
    if (appState.currentStep > 0) {
      setAppState(prev => ({ ...prev, currentStep: prev.currentStep - 1 }));
    }
  };

  const setCurrentStep = (step: number) => {
    if (canNavigate(step)) {
      setAppState(prev => ({ ...prev, currentStep: step }));
    }
  };

  // Save session
  const saveSession = async () => {
    try {
      setAppState(prev => ({ ...prev, loading: true }));
      await ApiService.saveSession(appState.session);
      setAppState(prev => ({ ...prev, loading: false }));
    } catch (error) {
      console.error('Failed to save session:', error);
      setAppState(prev => ({ 
        ...prev, 
        error: 'Failed to save session',
        loading: false 
      }));
    }
  };

  // Render current step content
  const renderStepContent = () => {
    switch (appState.currentStep) {
      case 0: // Driver Details
        const driverFields = appState.ontology?.sections?.drivers?.fields || [];
        return (
          <div className="space-y-6">
            {appState.session.drivers.map((driver, index) => (
              <UniversalForm
                key={driver.id}
                title={`Driver ${index + 1} ${driver.classification === 'MAIN' ? '(Main Driver)' : '(Named Driver)'}`}
                icon={<User className="w-5 h-5 mr-2 text-blue-600" />}
                fields={driverFields}
                data={driver}
                onUpdate={(field, value) => updateDriver(index, field, value)}
                onRemove={index > 0 ? () => removeDriver(index) : undefined}
                showRemoveButton={index > 0}
                validationErrors={appState.validationErrors}
                formType="driver"
              />
            ))}
            <Button onClick={addDriver} color="light" className="w-full">
              Add Named Driver
            </Button>
          </div>
        );

      case 1: // Vehicle Details
        const vehicleFields = appState.ontology?.sections?.vehicles?.fields || [];
        return (
          <div className="space-y-6">
            {appState.session.vehicles.map((vehicle, index) => (
              <UniversalForm
                key={vehicle.id}
                title={`Vehicle ${index + 1}`}
                icon={<Car className="w-5 h-5 mr-2 text-blue-600" />}
                fields={vehicleFields}
                data={vehicle}
                onUpdate={(field, value) => updateVehicle(index, field, value)}
                onRemove={index > 0 ? () => removeVehicle(index) : undefined}
                showRemoveButton={index > 0}
                validationErrors={appState.validationErrors}
                formType="vehicle"
              />
            ))}
            <Button onClick={addVehicle} color="light" className="w-full">
              Add Vehicle
            </Button>
          </div>
        );

      case 2: // Claims History
        const claimsFields = appState.ontology?.sections?.claims?.fields || [];
        return (
          <div className="space-y-6">
            <UniversalForm
              title="Claims History"
              icon={<FileText className="w-5 h-5 mr-2 text-blue-600" />}
              fields={claimsFields}
              data={appState.session.claims}
              onUpdate={(field, value) => {
                setAppState(prev => ({
                  ...prev,
                  session: {
                    ...prev.session,
                    claims: { ...prev.session.claims, [field]: value }
                  }
                }));
              }}
              validationErrors={appState.validationErrors}
              formType="claims"
            />
          </div>
        );

      case 3: // Policy Details
        const policyFields = appState.ontology?.sections?.policy?.fields || [];
        return (
          <div className="space-y-6">
            <UniversalForm
              title="Policy Details"
              icon={<Shield className="w-5 h-5 mr-2 text-blue-600" />}
              fields={policyFields}
              data={appState.session.policy}
              onUpdate={(field, value) => {
                setAppState(prev => ({
                  ...prev,
                  session: {
                    ...prev.session,
                    policy: { ...prev.session.policy, [field]: value }
                  }
                }));
              }}
              validationErrors={appState.validationErrors}
            />
          </div>
        );

      case 4: // Documents
        const documentFields = appState.ontology?.sections?.documents?.fields || [];
        return (
          <DocumentManager
            fields={documentFields}
            onUpload={(documentType: string, file: File) => {
              const fileList = new DataTransfer();
              fileList.items.add(file);
              handleFileUpload(fileList.files, documentType, documentType);
            }}
            onDownload={handleFileDownload}
            onOCR={handleOCRProcessing}
          />
        );

      case 5: // Settings
        const settingsFields = appState.ontology?.sections?.settings?.fields || [];
        return (
          <div className="space-y-6">
            <UniversalForm
              title="Personal Details & Settings"
              icon={<User className="w-5 h-5 mr-2 text-blue-600" />}
              fields={settingsFields}
              data={appState.session}
              onUpdate={(field, value) => {
                setAppState(prev => ({
                  ...prev,
                  session: { ...prev.session, [field]: value }
                }));
              }}
              validationErrors={appState.validationErrors}
              formType="settings"
            />
          </div>
        );

      case 6: // Review & Submit
        return (
          <Card>
            <h3 className="text-lg font-semibold mb-4">Review & Submit</h3>
            <p className="text-gray-600">Review summary will be implemented here</p>
            <Button onClick={saveSession} className="mt-4" disabled={appState.loading}>
              {appState.loading ? 'Saving...' : 'Submit Application'}
            </Button>
          </Card>
        );

      default:
        return <div>Invalid step</div>;
    }
  };

  if (appState.loading && !appState.ontology) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-600 mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading application...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 flex">
      {/* Left Sidebar Navigation */}
      <Navigation
        currentStep={appState.currentStep}
        setCurrentStep={setCurrentStep}
        steps={steps}
        canNavigate={canNavigate}
        completedSteps={getCompletedSteps()}
      />
      
      {/* Main Content Area */}
      <div className="flex-1 ml-64">
        <div className="container mx-auto px-8 py-8">
          {/* Header */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold text-gray-900 mb-2">CLIENT-UX Insurance Application</h1>
            <p className="text-gray-600">Complete your insurance application step by step</p>
          </div>

          {/* Progress Bar */}
          <div className="mb-8">
            <div className="flex justify-between items-center mb-2">
              <span className="text-sm font-medium text-gray-700">
                Step {appState.currentStep + 1} of {steps.length}
              </span>
              <span className="text-sm text-gray-500">
                {Math.round(((appState.currentStep + 1) / steps.length) * 100)}% Complete
              </span>
            </div>
            <Progress 
              progress={((appState.currentStep + 1) / steps.length) * 100} 
              className="mb-4"
            />
          </div>

        {/* Error Display */}
        {appState.error && (
          <Alert color="failure" className="mb-6">
            <AlertCircle className="w-4 h-4" />
            <span>{appState.error}</span>
          </Alert>
        )}

        {/* Main Content */}
        <div className="mb-8">
          {renderStepContent()}
        </div>

        {/* Navigation Buttons */}
        <div className="flex justify-between">
          <Button
            color="light"
            onClick={prevStep}
            disabled={appState.currentStep === 0}
          >
            <ChevronLeft className="w-4 h-4 mr-2" />
            Previous
          </Button>

          <Button
            onClick={nextStep}
            disabled={appState.currentStep === steps.length - 1}
          >
            Next
            <ChevronRight className="w-4 h-4 ml-2" />
          </Button>
          </div>
        </div>
      </div>
      <GlobalAssistant
        isOpen={showAssistant}
        onClose={() => setShowAssistant(false)}
        onOpen={() => setShowAssistant(true)}
      />
    </div>
  );
};

export default App;
