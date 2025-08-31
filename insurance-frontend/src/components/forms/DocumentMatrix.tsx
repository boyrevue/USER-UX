import React from 'react';
import { Card, Button, Badge } from 'flowbite-react';
import { Upload, FileText, Check, X } from 'lucide-react';

import { OntologyField } from '../../types';

interface DocumentMatrixProps {
  fields: OntologyField[];
  onFileUpload: (files: FileList, documentType: string, fieldName: string) => void;
  isProcessing: boolean;
  documents: Array<{
    id: string;
    name: string;
    type: string;
    fieldName?: string;
  }>;
  onRemoveDocument: (id: string) => void;
}

const DocumentMatrix: React.FC<DocumentMatrixProps> = ({
  fields,
  onFileUpload,
  isProcessing,
  documents,
  onRemoveDocument
}) => {
  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>, field: OntologyField) => {
    const files = event.target.files;
    if (files && files.length > 0) {
      onFileUpload(files, field.property, field.property);
    }
  };

  const getUploadedDocument = (fieldName: string) => {
    return documents.find(doc => doc.fieldName === fieldName);
  };

  const getDocumentIcon = (fieldName: string) => {
    const iconMap: { [key: string]: string } = {
      'drivingLicence': '🪪',
      'passport': '📘',
      'utilityBill': '📄',
      'bankStatement': '🏦',
      'payslip': '💰',
      'p60': '📊',
      'medicalCertificate': '🏥',
      'insuranceCertificate': '🛡️',
      'vehicleRegistration': '🚗',
      'identityCard': '🆔',
      'proofOfAddress': '🏠',
      'employmentLetter': '💼',
      'taxReturn': '📋',
      'creditReport': '📈',
      'mortgageStatement': '🏡',
      'pensionStatement': '👴',
      'studentLoan': '🎓',
      'courtOrder': '⚖️',
      'deathCertificate': '⚱️'
    };
    return iconMap[fieldName] || '📄';
  };

  if (!fields || fields.length === 0) {
    return (
      <Card>
        <div className="text-center py-8">
          <FileText className="w-12 h-12 mx-auto mb-4 text-gray-400" />
          <p className="text-gray-600">No document fields available</p>
        </div>
      </Card>
    );
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-semibold">Document Upload Matrix</h3>
        <Badge color="info">
          {documents.length} of {fields.length} uploaded
        </Badge>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {fields.map((field) => {
          const uploadedDoc = getUploadedDocument(field.property);
          const isUploaded = !!uploadedDoc;
          
          return (
            <Card key={field.property} className={`relative ${isUploaded ? 'ring-2 ring-green-500' : ''}`}>
              <div className="space-y-3">
                {/* Header */}
                <div className="flex items-center justify-between">
                  <div className="flex items-center">
                    <span className="text-2xl mr-2">{getDocumentIcon(field.property)}</span>
                    <div>
                      <h4 className="font-medium text-sm">{field.label}</h4>
                      {field.required && (
                        <Badge color="failure" size="xs">Required</Badge>
                      )}
                    </div>
                  </div>
                  {isUploaded && (
                    <Check className="w-5 h-5 text-green-600" />
                  )}
                </div>

                {/* Help Text */}
                {field.helpText && (
                  <p className="text-xs text-gray-600">{field.helpText}</p>
                )}

                {/* Upload Area */}
                {!isUploaded ? (
                  <div className="border-2 border-dashed border-gray-300 rounded-lg p-4 text-center hover:border-blue-400 transition-colors">
                    <Upload className="w-8 h-8 mx-auto mb-2 text-gray-400" />
                    <input
                      type="file"
                      accept="image/*,.pdf"
                      onChange={(e) => handleFileChange(e, field)}
                      className="hidden"
                      id={`file-upload-${field.property}`}
                      disabled={isProcessing}
                    />
                    <label 
                      htmlFor={`file-upload-${field.property}`} 
                      className={`cursor-pointer ${isProcessing ? 'cursor-not-allowed' : ''}`}
                    >
                      <Button size="xs" as="span">
                        {isProcessing ? 'Processing...' : 'Upload'}
                      </Button>
                    </label>
                  </div>
                ) : (
                  <div className="bg-green-50 border border-green-200 rounded-lg p-3">
                    <div className="flex items-center justify-between">
                      <div className="flex items-center">
                        <FileText className="w-4 h-4 mr-2 text-green-600" />
                        <span className="text-sm font-medium text-green-800 truncate">
                          {uploadedDoc.name}
                        </span>
                      </div>
                      <Button
                        size="xs"
                        color="failure"
                        onClick={() => onRemoveDocument(uploadedDoc.id)}
                      >
                        <X className="w-3 h-3" />
                      </Button>
                    </div>
                  </div>
                )}
              </div>
            </Card>
          );
        })}
      </div>

      {/* Summary */}
      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <div className="flex items-center">
          <FileText className="w-5 h-5 mr-2 text-blue-600" />
          <div>
            <p className="text-sm font-medium text-blue-800">
              Upload Progress: {documents.length} of {fields.length} documents
            </p>
            <p className="text-xs text-blue-600">
              {fields.filter(f => f.required).length} required documents, {fields.filter(f => !f.required).length} optional
            </p>
          </div>
        </div>
      </div>
    </div>
  );
};

export default DocumentMatrix;
