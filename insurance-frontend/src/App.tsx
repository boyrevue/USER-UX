import React, { useState, useEffect } from 'react';
import { Card, Button, Progress, Alert } from 'flowbite-react';
import { ChevronLeft, ChevronRight, CheckCircle, AlertCircle } from 'lucide-react';

// Import new modular components
import DriverForm from './components/forms/DriverForm';
import VehicleForm from './components/forms/VehicleForm';
import ClaimsForm from './components/forms/ClaimsForm';
import DocumentUpload from './components/forms/DocumentUpload';
import Navigation from './components/layout/Navigation';

// Import services
import { ApiService } from './services/api';
import { validateBirthDate, validateLicenceDate, validateHistoricalDate } from './services/validation';

// Import types
import { 
  AppState, 
  Session, 
  Driver, 
  Vehicle, 
  Claim, 
  Accident, 
  Policy, 
  Document,
  ValidationError 
} from './types';

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

  const steps = [
    'Driver Details',
    'Vehicle Details', 
    'Claims History',
    'Policy Details',
    'Documents',
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
        return (
          <div className="space-y-6">
            {appState.session.drivers.map((driver, index) => (
              <DriverForm
                key={driver.id}
                driver={driver}
                index={index}
                updateDriver={updateDriver}
                removeDriver={removeDriver}
                validationErrors={appState.validationErrors}
              />
            ))}
            <Button onClick={addDriver} color="light" className="w-full">
              Add Named Driver
            </Button>
          </div>
        );

      case 1: // Vehicle Details
        return (
          <div className="space-y-6">
            {appState.session.vehicles.map((vehicle, index) => (
              <VehicleForm
                key={vehicle.id}
                vehicle={vehicle}
                index={index}
                updateVehicle={updateVehicle}
                removeVehicle={removeVehicle}
                validationErrors={appState.validationErrors}
              />
            ))}
            <Button onClick={addVehicle} color="light" className="w-full">
              Add Vehicle
            </Button>
          </div>
        );

      case 2: // Claims History
        return (
          <ClaimsForm
            claims={appState.session.claims.claims}
            accidents={appState.session.claims.accidents}
            updateClaim={updateClaim}
            updateAccident={updateAccident}
            addClaim={addClaim}
            addAccident={addAccident}
            removeClaim={removeClaim}
            removeAccident={removeAccident}
            validationErrors={appState.validationErrors}
          />
        );

      case 3: // Policy Details
        return (
          <Card>
            <h3 className="text-lg font-semibold mb-4">Policy Configuration</h3>
            <p className="text-gray-600">Policy details form will be implemented here</p>
          </Card>
        );

      case 4: // Documents
        return (
          <DocumentUpload
            onFileUpload={handleFileUpload}
            isProcessing={isProcessing}
            documents={appState.session.documents}
            onRemoveDocument={removeDocument}
          />
        );

      case 5: // Review & Submit
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
    <div className="min-h-screen bg-gray-50">
      <div className="container mx-auto px-4 py-8">
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

        {/* Navigation */}
        <Navigation
          currentStep={appState.currentStep}
          setCurrentStep={setCurrentStep}
          steps={steps}
          canNavigate={canNavigate}
          completedSteps={getCompletedSteps()}
        />

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
  );
};

export default App;
