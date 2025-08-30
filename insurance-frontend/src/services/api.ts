
import { OntologyResponse } from '../types';

const API_BASE_URL = 'http://localhost:3000/api';

export interface ProcessDocumentResponse {
  documentType: string;
  uploadType: string;
  extractedFields: Record<string, any>;
  confidence: number;
  processedAt: string;
  images?: {
    page1?: string;
    page2Upper?: string;
    page2MRZ?: string;
    preprocessed?: string;
  };
  metadata?: {
    engine: string;
    ocrConfidence: number;
    fieldConfidence: number;
  };
}

export interface SessionResponse {
  success: boolean;
  sessionId: string;
  message?: string;
}



export class ApiService {
  static async processDocument(
    files: FileList, 
    uploadType: string, 
    selectedDocumentType?: string
  ): Promise<ProcessDocumentResponse> {
    const formData = new FormData();
    
    Array.from(files).forEach(file => {
      formData.append('files', file);
    });
    
    formData.append('uploadType', uploadType);
    
    if (selectedDocumentType) {
      formData.append('selectedDocumentType', selectedDocumentType);
    }

    const response = await fetch(`${API_BASE_URL}/process-document`, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`OCR processing failed: ${errorText}`);
    }

    return response.json();
  }

  static async saveSession(sessionData: any): Promise<SessionResponse> {
    const response = await fetch(`${API_BASE_URL}/save-session`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(sessionData),
    });

    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to save session: ${errorText}`);
    }

    return response.json();
  }

  static async getSession(sessionId: string): Promise<any> {
    const response = await fetch(`${API_BASE_URL}/sessions/${sessionId}`);
    
    if (!response.ok) {
      if (response.status === 404) {
        throw new Error('Session not found');
      }
      const errorText = await response.text();
      throw new Error(`Failed to get session: ${errorText}`);
    }

    return response.json();
  }

  static async getOntology(): Promise<OntologyResponse> {
    const response = await fetch(`${API_BASE_URL}/ontology`);
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to load ontology: ${errorText}`);
    }

    return response.json();
  }

  static async exportData(sessionId: string, format: 'json' | 'pdf' = 'json'): Promise<Blob> {
    const response = await fetch(`${API_BASE_URL}/export/${sessionId}?format=${format}`);
    
    if (!response.ok) {
      const errorText = await response.text();
      throw new Error(`Failed to export data: ${errorText}`);
    }

    return response.blob();
  }

  static async healthCheck(): Promise<{ status: string; timestamp: string; version: string }> {
    const response = await fetch(`${API_BASE_URL}/health`);
    
    if (!response.ok) {
      throw new Error('Health check failed');
    }

    return response.json();
  }
}
