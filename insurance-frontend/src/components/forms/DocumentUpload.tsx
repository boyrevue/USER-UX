import React from 'react';



interface DocumentUploadProps {
  onFileUpload: (files: FileList, type: string) => void;
  isProcessing: boolean;
  extractedData?: any;
}

const DocumentUpload: React.FC<DocumentUploadProps> = ({ 
  onFileUpload, 
  isProcessing, 
  extractedData 
}) => {
  return (
    <div className="document-upload">
      {/* Document upload content will be moved here */}
      <h3>Document Upload</h3>
      {/* Upload interface */}
    </div>
  );
};


export default DocumentUpload;
