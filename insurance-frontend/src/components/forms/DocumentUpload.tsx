import React from 'react';
import { Card, Button } from 'flowbite-react';
import { Upload } from 'lucide-react';
import { DocumentUploadProps } from '../../types';

const DocumentUpload: React.FC<DocumentUploadProps> = ({ 
  onFileUpload,
  isProcessing,
  documents,
  onRemoveDocument
}) => {
  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const files = event.target.files;
    if (files && files.length > 0) {
      onFileUpload(files, 'passport', 'passport');
    }
  };

  return (
    <Card>
      <div className="flex items-center mb-4">
        <Upload className="w-5 h-5 mr-2 text-blue-600" />
        <h3 className="text-lg font-semibold">Document Upload</h3>
      </div>
      
      <div className="space-y-4">
        <div className="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
          <Upload className="w-12 h-12 mx-auto mb-4 text-gray-400" />
          <p className="text-gray-600 mb-4">Upload your documents for automatic data extraction</p>
          <input
            type="file"
            accept="image/*,.pdf"
            onChange={handleFileChange}
            className="hidden"
            id="file-upload"
            disabled={isProcessing}
          />
          <label htmlFor="file-upload" className={isProcessing ? 'cursor-not-allowed' : 'cursor-pointer'}>
            <Button as="span">
              {isProcessing ? 'Processing...' : 'Choose Files'}
            </Button>
          </label>
        </div>

        {documents.length > 0 && (
          <div>
            <h4 className="font-medium mb-2">Uploaded Documents ({documents.length})</h4>
            <div className="space-y-2">
              {documents.map((doc) => (
                <div key={doc.id} className="flex justify-between items-center p-2 bg-gray-50 rounded">
                  <span className="text-sm">{doc.name}</span>
                  <Button size="sm" color="failure" onClick={() => onRemoveDocument(doc.id)}>
                    Remove
                  </Button>
                </div>
              ))}
            </div>
          </div>
        )}
      </div>
    </Card>
  );
};

export default DocumentUpload;
