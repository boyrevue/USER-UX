
const API_BASE_URL = 'http://localhost:3000/api';

export class ApiService {
  static async processDocument(files: FileList, uploadType: string) {
    const formData = new FormData();
    Array.from(files).forEach(file => {
      formData.append('files', file);
    });
    formData.append('uploadType', uploadType);

    const response = await fetch(`${API_BASE_URL}/process-document`, {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  static async saveSession(sessionData: any) {
    const response = await fetch(`${API_BASE_URL}/save-session`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(sessionData),
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }

  static async getOntology() {
    const response = await fetch(`${API_BASE_URL}/ontology`);
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    return response.json();
  }
}
