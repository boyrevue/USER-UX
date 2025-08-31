import React, { useState } from 'react';
import { Card, Button, Badge, Alert } from 'flowbite-react';
import { 
  Upload, 
  Download, 
  Eye, 
  Scan, 
  FileText, 
  User, 
  CreditCard, 
  GraduationCap, 
  Briefcase, 
  Home, 
  Heart, 
  Plane, 
  Car, 
  Scale, 
  Shield,
  FolderOpen,
  CheckCircle,
  AlertCircle
} from 'lucide-react';
import { OntologyField } from '../../types';

interface DocumentCategory {
  id: string;
  name: string;
  icon: React.ReactNode;
  description: string;
  documents: DocumentType[];
}

interface DocumentType {
  id: string;
  name: string;
  description: string;
  ocrFields: string[];
  required: boolean;
  uploaded?: boolean;
  ocrProcessed?: boolean;
}

interface DocumentManagerProps {
  fields: OntologyField[];
  onUpload: (documentType: string, file: File) => void;
  onDownload: (documentType: string) => void;
  onOCR: (documentType: string, file: File) => Promise<any>;
}

const DocumentManager: React.FC<DocumentManagerProps> = ({
  fields,
  onUpload,
  onDownload,
  onOCR
}) => {
  const [activeCategory, setActiveCategory] = useState('identity');
  const [ocrResults, setOcrResults] = useState<{[key: string]: any}>({});
  const [processing, setProcessing] = useState<{[key: string]: boolean}>({});

  // Define document categories based on ontology
  const documentCategories: DocumentCategory[] = [
    {
      id: 'identity',
      name: 'Identity Documents',
      icon: <User className="w-5 h-5" />,
      description: 'Documents that establish or verify identity',
      documents: [
        {
          id: 'passport',
          name: 'Passport',
          description: 'Official travel document issued by government',
          ocrFields: ['passportNumber', 'firstName', 'lastName', 'dateOfBirth', 'nationality', 'expiryDate'],
          required: true
        },
        {
          id: 'driversLicense',
          name: "Driver's License",
          description: 'License permitting operation of motor vehicles',
          ocrFields: ['licenceNumber', 'firstName', 'lastName', 'dateOfBirth', 'address', 'licenceExpiry'],
          required: true
        },
        {
          id: 'nationalIdCard',
          name: 'National ID Card',
          description: 'Government-issued identity card',
          ocrFields: ['idNumber', 'firstName', 'lastName', 'dateOfBirth', 'address'],
          required: false
        },
        {
          id: 'birthCertificate',
          name: 'Birth Certificate',
          description: 'Official record of birth',
          ocrFields: ['firstName', 'lastName', 'dateOfBirth', 'placeOfBirth', 'parentNames'],
          required: false
        }
      ]
    },
    {
      id: 'financial',
      name: 'Financial Documents',
      icon: <CreditCard className="w-5 h-5" />,
      description: 'Financial accounts, transactions, and obligations',
      documents: [
        {
          id: 'bankStatement',
          name: 'Bank Statement',
          description: 'Monthly bank account statement',
          ocrFields: ['accountNumber', 'sortCode', 'accountHolderName', 'balance', 'statementDate'],
          required: true
        },
        {
          id: 'payslip',
          name: 'Payslip',
          description: 'Monthly salary statement',
          ocrFields: ['employerName', 'employeeName', 'grossSalary', 'netSalary', 'payPeriod'],
          required: true
        },
        {
          id: 'taxReturn',
          name: 'Tax Return',
          description: 'Annual tax return document',
          ocrFields: ['taxYear', 'totalIncome', 'taxPaid', 'refundAmount'],
          required: false
        }
      ]
    },
    {
      id: 'education',
      name: 'Educational Documents',
      icon: <GraduationCap className="w-5 h-5" />,
      description: 'Certificates, diplomas, and educational credentials',
      documents: [
        {
          id: 'degree',
          name: 'Degree Certificate',
          description: 'University degree certificate',
          ocrFields: ['degreeName', 'university', 'graduationDate', 'studentName', 'classification'],
          required: false
        },
        {
          id: 'transcript',
          name: 'Academic Transcript',
          description: 'Official academic record',
          ocrFields: ['studentName', 'courses', 'grades', 'gpa', 'institution'],
          required: false
        }
      ]
    },
    {
      id: 'employment',
      name: 'Employment Documents',
      icon: <Briefcase className="w-5 h-5" />,
      description: 'Work-related contracts, references, and records',
      documents: [
        {
          id: 'employmentContract',
          name: 'Employment Contract',
          description: 'Official employment agreement',
          ocrFields: ['employerName', 'employeeName', 'jobTitle', 'salary', 'startDate'],
          required: false
        },
        {
          id: 'reference',
          name: 'Employment Reference',
          description: 'Reference letter from employer',
          ocrFields: ['refereeTitle', 'companyName', 'employeeName', 'employmentPeriod'],
          required: false
        }
      ]
    },
    {
      id: 'property',
      name: 'Property Documents',
      icon: <Home className="w-5 h-5" />,
      description: 'Real estate and property ownership documents',
      documents: [
        {
          id: 'deed',
          name: 'Property Deed',
          description: 'Legal document of property ownership',
          ocrFields: ['propertyAddress', 'ownerName', 'purchaseDate', 'purchasePrice'],
          required: false
        },
        {
          id: 'mortgage',
          name: 'Mortgage Agreement',
          description: 'Home loan agreement',
          ocrFields: ['lenderName', 'borrowerName', 'loanAmount', 'interestRate', 'term'],
          required: false
        }
      ]
    },
    {
      id: 'medical',
      name: 'Medical Documents',
      icon: <Heart className="w-5 h-5" />,
      description: 'Health records and medical certificates',
      documents: [
        {
          id: 'medicalCertificate',
          name: 'Medical Certificate',
          description: 'Doctor-issued health certificate',
          ocrFields: ['patientName', 'doctorName', 'diagnosis', 'treatmentDate', 'restrictions'],
          required: false
        },
        {
          id: 'prescription',
          name: 'Prescription',
          description: 'Medical prescription document',
          ocrFields: ['patientName', 'doctorName', 'medication', 'dosage', 'prescriptionDate'],
          required: false
        }
      ]
    },
    {
      id: 'travel',
      name: 'Travel Documents',
      icon: <Plane className="w-5 h-5" />,
      description: 'Travel and immigration related documents',
      documents: [
        {
          id: 'visa',
          name: 'Visa',
          description: 'Entry permit for foreign country',
          ocrFields: ['visaNumber', 'country', 'validFrom', 'validUntil', 'visaType'],
          required: false
        },
        {
          id: 'travelInsurance',
          name: 'Travel Insurance',
          description: 'Travel insurance policy',
          ocrFields: ['policyNumber', 'insuredName', 'coverageAmount', 'validFrom', 'validUntil'],
          required: false
        }
      ]
    },
    {
      id: 'vehicle',
      name: 'Vehicle Documents',
      icon: <Car className="w-5 h-5" />,
      description: 'Vehicle ownership and operation documents',
      documents: [
        {
          id: 'registration',
          name: 'Vehicle Registration',
          description: 'Official vehicle registration document',
          ocrFields: ['registrationNumber', 'ownerName', 'vehicleMake', 'vehicleModel', 'registrationDate'],
          required: true
        },
        {
          id: 'mot',
          name: 'MOT Certificate',
          description: 'Ministry of Transport test certificate',
          ocrFields: ['registrationNumber', 'testDate', 'expiryDate', 'mileage', 'testResult'],
          required: true
        }
      ]
    },
    {
      id: 'legal',
      name: 'Legal Documents',
      icon: <Scale className="w-5 h-5" />,
      description: 'Contracts, wills, and legal instruments',
      documents: [
        {
          id: 'will',
          name: 'Last Will and Testament',
          description: 'Legal document specifying asset distribution',
          ocrFields: ['testatorName', 'executorName', 'beneficiaries', 'dateCreated'],
          required: false
        },
        {
          id: 'powerOfAttorney',
          name: 'Power of Attorney',
          description: 'Legal authorization document',
          ocrFields: ['principalName', 'agentName', 'powers', 'effectiveDate'],
          required: false
        }
      ]
    },
    {
      id: 'insurance',
      name: 'Insurance Documents',
      icon: <Shield className="w-5 h-5" />,
      description: 'Insurance policies and coverage documents',
      documents: [
        {
          id: 'homeInsurance',
          name: 'Home Insurance Policy',
          description: 'Property insurance policy',
          ocrFields: ['policyNumber', 'insuredName', 'propertyAddress', 'coverageAmount', 'premium'],
          required: false
        },
        {
          id: 'lifeInsurance',
          name: 'Life Insurance Policy',
          description: 'Life insurance coverage document',
          ocrFields: ['policyNumber', 'insuredName', 'beneficiaries', 'coverageAmount', 'premium'],
          required: false
        }
      ]
    }
  ];

  const handleFileUpload = async (documentType: string, event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (!file) return;

    // Upload file
    onUpload(documentType, file);

    // Process with OCR
    setProcessing(prev => ({ ...prev, [documentType]: true }));
    try {
      const ocrResult = await onOCR(documentType, file);
      setOcrResults(prev => ({ ...prev, [documentType]: ocrResult }));
    } catch (error) {
      console.error('OCR processing failed:', error);
    } finally {
      setProcessing(prev => ({ ...prev, [documentType]: false }));
    }
  };

  const renderDocumentCard = (doc: DocumentType, categoryId: string) => (
    <Card key={doc.id} className="mb-4">
      <div className="flex items-start justify-between">
        <div className="flex-1">
          <div className="flex items-center mb-2">
            <FileText className="w-4 h-4 mr-2 text-gray-600" />
            <h4 className="font-semibold text-gray-900">{doc.name}</h4>
            {doc.required && <Badge color="failure" size="sm" className="ml-2">Required</Badge>}
            {doc.uploaded && <Badge color="success" size="sm" className="ml-2">Uploaded</Badge>}
            {ocrResults[doc.id] && <Badge color="info" size="sm" className="ml-2">OCR Processed</Badge>}
          </div>
          <p className="text-sm text-gray-600 mb-3">{doc.description}</p>
          
          {/* OCR Fields Preview */}
          <div className="mb-3">
            <p className="text-xs font-medium text-gray-700 mb-1">OCR will extract:</p>
            <div className="flex flex-wrap gap-1">
              {doc.ocrFields.map(field => (
                <Badge key={field} color="gray" size="xs">{field}</Badge>
              ))}
            </div>
          </div>

          {/* OCR Results */}
          {ocrResults[doc.id] && (
            <div className="mt-3 p-3 bg-green-50 rounded-lg">
              <p className="text-sm font-medium text-green-800 mb-2">Extracted Data:</p>
              <div className="space-y-1">
                {Object.entries(ocrResults[doc.id]).map(([key, value]) => (
                  <div key={key} className="flex justify-between text-xs">
                    <span className="text-green-700 font-medium">{key}:</span>
                    <span className="text-green-600">{String(value)}</span>
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        <div className="flex flex-col gap-2 ml-4">
          {/* Upload Button */}
          <label className="cursor-pointer">
            <input
              type="file"
              className="hidden"
              accept=".pdf,.jpg,.jpeg,.png"
              onChange={(e) => handleFileUpload(doc.id, e)}
              disabled={processing[doc.id]}
            />
            <Button size="sm" color="blue" disabled={processing[doc.id]}>
              {processing[doc.id] ? (
                <>
                  <Scan className="w-4 h-4 mr-1 animate-spin" />
                  Processing...
                </>
              ) : (
                <>
                  <Upload className="w-4 h-4 mr-1" />
                  Upload
                </>
              )}
            </Button>
          </label>

          {/* Download Button */}
          {doc.uploaded && (
            <Button size="sm" color="gray" onClick={() => onDownload(doc.id)}>
              <Download className="w-4 h-4 mr-1" />
              Download
            </Button>
          )}

          {/* View Button */}
          {doc.uploaded && (
            <Button size="sm" color="light">
              <Eye className="w-4 h-4 mr-1" />
              View
            </Button>
          )}
        </div>
      </div>
    </Card>
  );

  const activeCategory_data = documentCategories.find(cat => cat.id === activeCategory);

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="text-center">
        <h2 className="text-2xl font-bold text-gray-900 mb-2">Document Management Center</h2>
        <p className="text-gray-600">Upload, process with OCR, and manage all your personal documents</p>
      </div>

      {/* Category Navigation */}
      <div className="flex flex-wrap gap-2 mb-6">
        {documentCategories.map((category) => (
          <Button
            key={category.id}
            color={activeCategory === category.id ? "blue" : "light"}
            size="sm"
            onClick={() => setActiveCategory(category.id)}
            className="flex items-center"
          >
            {category.icon}
            <span className="ml-2">{category.name}</span>
            <Badge color="gray" size="sm" className="ml-2">
              {category.documents.length}
            </Badge>
          </Button>
        ))}
      </div>

      {/* Active Category Content */}
      {activeCategory_data && (
        <div className="mt-6">
          {/* Category Description */}
          <Alert color="info" className="mb-6">
            <FolderOpen className="w-4 h-4" />
            <span className="ml-2">{activeCategory_data.description}</span>
          </Alert>

          {/* Documents Grid */}
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {activeCategory_data.documents.map(doc => renderDocumentCard(doc, activeCategory_data.id))}
          </div>

          {/* Category Stats */}
          <div className="mt-6 p-4 bg-gray-50 rounded-lg">
            <div className="flex justify-between items-center">
              <div className="flex items-center">
                <CheckCircle className="w-5 h-5 text-green-600 mr-2" />
                <span className="text-sm font-medium">
                  {activeCategory_data.documents.filter(d => d.uploaded).length} of {activeCategory_data.documents.length} uploaded
                </span>
              </div>
              <div className="flex items-center">
                <Scan className="w-5 h-5 text-blue-600 mr-2" />
                <span className="text-sm font-medium">
                  {Object.keys(ocrResults).filter(key => 
                    activeCategory_data.documents.some(d => d.id === key)
                  ).length} OCR processed
                </span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default DocumentManager;
