import { ApiService } from './api';

export interface OCRResult {
  [key: string]: string | number | Date;
}

export interface DocumentOCRMapping {
  [documentType: string]: {
    fields: string[];
    patterns: { [field: string]: RegExp };
    postProcess?: (raw: any) => OCRResult;
  };
}

// Define OCR field mappings for each document type
const OCR_MAPPINGS: DocumentOCRMapping = {
  passport: {
    fields: ['passportNumber', 'firstName', 'lastName', 'dateOfBirth', 'nationality', 'expiryDate'],
    patterns: {
      passportNumber: /(?:Passport\s*No\.?\s*|P<[A-Z]{3})([A-Z0-9]{6,9})/i,
      firstName: /(?:Given\s*Names?\s*|First\s*Name\s*)([A-Z\s]+)/i,
      lastName: /(?:Surname\s*|Last\s*Name\s*)([A-Z\s]+)/i,
      dateOfBirth: /(?:Date\s*of\s*Birth\s*|DOB\s*)(\d{2}[\/\-]\d{2}[\/\-]\d{4})/i,
      nationality: /(?:Nationality\s*)([A-Z\s]+)/i,
      expiryDate: /(?:Date\s*of\s*Expiry\s*|Expires?\s*)(\d{2}[\/\-]\d{2}[\/\-]\d{4})/i
    },
    postProcess: (raw) => ({
      ...raw,
      dateOfBirth: raw.dateOfBirth ? new Date(raw.dateOfBirth) : null,
      expiryDate: raw.expiryDate ? new Date(raw.expiryDate) : null
    })
  },
  
  driversLicense: {
    fields: ['licenceNumber', 'firstName', 'lastName', 'dateOfBirth', 'address', 'licenceExpiry'],
    patterns: {
      licenceNumber: /(?:Licence\s*No\.?\s*|DL\s*No\.?\s*)([A-Z0-9]{5,16})/i,
      firstName: /(?:First\s*Name\s*|Given\s*Name\s*)([A-Z\s]+)/i,
      lastName: /(?:Last\s*Name\s*|Surname\s*)([A-Z\s]+)/i,
      dateOfBirth: /(?:Date\s*of\s*Birth\s*|DOB\s*)(\d{2}[\/\-]\d{2}[\/\-]\d{4})/i,
      address: /(?:Address\s*)([A-Z0-9\s,.-]+)/i,
      licenceExpiry: /(?:Valid\s*Until\s*|Expires?\s*)(\d{2}[\/\-]\d{2}[\/\-]\d{4})/i
    },
    postProcess: (raw) => ({
      ...raw,
      dateOfBirth: raw.dateOfBirth ? new Date(raw.dateOfBirth) : null,
      licenceExpiry: raw.licenceExpiry ? new Date(raw.licenceExpiry) : null
    })
  },

  bankStatement: {
    fields: ['accountNumber', 'sortCode', 'accountHolderName', 'balance', 'statementDate'],
    patterns: {
      accountNumber: /(?:Account\s*No\.?\s*|A\/C\s*No\.?\s*)(\d{8,12})/i,
      sortCode: /(?:Sort\s*Code\s*)(\d{2}[-\s]?\d{2}[-\s]?\d{2})/i,
      accountHolderName: /(?:Account\s*Holder\s*|Name\s*)([A-Z\s]+)/i,
      balance: /(?:Balance\s*|Current\s*Balance\s*)£?([0-9,]+\.?\d{0,2})/i,
      statementDate: /(?:Statement\s*Date\s*|Date\s*)(\d{2}[\/\-]\d{2}[\/\-]\d{4})/i
    },
    postProcess: (raw) => ({
      ...raw,
      balance: raw.balance ? parseFloat(raw.balance.replace(/,/g, '')) : null,
      statementDate: raw.statementDate ? new Date(raw.statementDate) : null
    })
  },

  payslip: {
    fields: ['employerName', 'employeeName', 'grossSalary', 'netSalary', 'payPeriod'],
    patterns: {
      employerName: /(?:Employer\s*|Company\s*)([A-Z\s&]+)/i,
      employeeName: /(?:Employee\s*|Name\s*)([A-Z\s]+)/i,
      grossSalary: /(?:Gross\s*Pay\s*|Gross\s*Salary\s*)£?([0-9,]+\.?\d{0,2})/i,
      netSalary: /(?:Net\s*Pay\s*|Take\s*Home\s*)£?([0-9,]+\.?\d{0,2})/i,
      payPeriod: /(?:Pay\s*Period\s*|Period\s*)(\d{2}[\/\-]\d{2}[\/\-]\d{4})/i
    },
    postProcess: (raw) => ({
      ...raw,
      grossSalary: raw.grossSalary ? parseFloat(raw.grossSalary.replace(/,/g, '')) : null,
      netSalary: raw.netSalary ? parseFloat(raw.netSalary.replace(/,/g, '')) : null,
      payPeriod: raw.payPeriod ? new Date(raw.payPeriod) : null
    })
  },

  registration: {
    fields: ['registrationNumber', 'ownerName', 'vehicleMake', 'vehicleModel', 'registrationDate'],
    patterns: {
      registrationNumber: /(?:Registration\s*No\.?\s*|Reg\s*No\.?\s*)([A-Z]{2}\d{2}\s?[A-Z]{3})/i,
      ownerName: /(?:Keeper\s*|Owner\s*)([A-Z\s]+)/i,
      vehicleMake: /(?:Make\s*)([A-Z\s]+)/i,
      vehicleModel: /(?:Model\s*)([A-Z0-9\s]+)/i,
      registrationDate: /(?:Date\s*of\s*First\s*Registration\s*)(\d{2}[\/\-]\d{2}[\/\-]\d{4})/i
    },
    postProcess: (raw) => ({
      ...raw,
      registrationDate: raw.registrationDate ? new Date(raw.registrationDate) : null
    })
  }
};

export class OCRService {
  private apiService: ApiService;

  constructor() {
    this.apiService = new ApiService();
  }

  /**
   * Process document with OCR and extract structured data
   */
  async processDocument(documentType: string, file: File): Promise<OCRResult> {
    try {
      // Step 1: Upload file and get OCR text
      const ocrText = await this.extractTextFromDocument(file);
      
      // Step 2: Apply document-specific field extraction
      const extractedData = this.extractFieldsFromText(documentType, ocrText);
      
      // Step 3: Validate and clean extracted data
      const validatedData = this.validateExtractedData(documentType, extractedData);
      
      // Step 4: Store in ontology format
      await this.storeInOntology(documentType, validatedData);
      
      return validatedData;
    } catch (error) {
      console.error('OCR processing failed:', error);
      const errorMessage = error instanceof Error ? error.message : 'Unknown error occurred';
      throw new Error(`Failed to process ${documentType}: ${errorMessage}`);
    }
  }

  /**
   * Extract raw text from document using backend OCR
   */
  private async extractTextFromDocument(file: File): Promise<string> {
    const formData = new FormData();
    formData.append('document', file);
    
    const response = await fetch('/api/ocr/extract', {
      method: 'POST',
      body: formData
    });
    
    if (!response.ok) {
      throw new Error('OCR extraction failed');
    }
    
    const result = await response.json();
    return result.text || '';
  }

  /**
   * Extract structured fields from OCR text using patterns
   */
  private extractFieldsFromText(documentType: string, text: string): OCRResult {
    const mapping = OCR_MAPPINGS[documentType];
    if (!mapping) {
      throw new Error(`No OCR mapping found for document type: ${documentType}`);
    }

    const extracted: OCRResult = {};
    
    // Apply regex patterns to extract each field
    for (const [fieldName, pattern] of Object.entries(mapping.patterns)) {
      const match = text.match(pattern);
      if (match && match[1]) {
        extracted[fieldName] = match[1].trim();
      }
    }

    // Apply post-processing if defined
    if (mapping.postProcess) {
      return mapping.postProcess(extracted);
    }

    return extracted;
  }

  /**
   * Validate extracted data against expected types and formats
   */
  private validateExtractedData(documentType: string, data: OCRResult): OCRResult {
    const validated: OCRResult = {};
    
    for (const [key, value] of Object.entries(data)) {
      if (value !== null && value !== undefined && value !== '') {
        // Basic validation - can be extended with more sophisticated rules
        if (key.includes('Date') && typeof value === 'string') {
          const date = new Date(value);
          validated[key] = isNaN(date.getTime()) ? value : date;
        } else if (key.includes('Number') && typeof value === 'string') {
          // Clean up numbers (remove spaces, special chars)
          validated[key] = value.replace(/[^\w]/g, '');
        } else if (key.includes('Name') && typeof value === 'string') {
          // Capitalize names properly
          validated[key] = value.toLowerCase().replace(/\b\w/g, l => l.toUpperCase());
        } else {
          validated[key] = value;
        }
      }
    }
    
    return validated;
  }

  /**
   * Store extracted data in ontology format
   */
  private async storeInOntology(documentType: string, data: OCRResult): Promise<void> {
    try {
      await fetch('/api/ontology/store-document-data', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          documentType,
          extractedData: data,
          timestamp: new Date().toISOString()
        })
      });
    } catch (error) {
      console.warn('Failed to store in ontology:', error);
      // Don't throw - OCR extraction succeeded, storage is secondary
    }
  }

  /**
   * Get available document types and their OCR fields
   */
  getDocumentTypes(): { [type: string]: string[] } {
    const types: { [type: string]: string[] } = {};
    
    for (const [docType, mapping] of Object.entries(OCR_MAPPINGS)) {
      types[docType] = mapping.fields;
    }
    
    return types;
  }

  /**
   * Preview what fields will be extracted for a document type
   */
  getExtractableFields(documentType: string): string[] {
    const mapping = OCR_MAPPINGS[documentType];
    return mapping ? mapping.fields : [];
  }
}

export default new OCRService();
