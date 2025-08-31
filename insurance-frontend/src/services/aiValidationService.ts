interface ValidationRequest {
  fieldName: string;
  userInput: string;
  validationPrompt: string;
}

interface ValidationResponse {
  isValid: boolean;
  message: string;
  suggestions?: string;
  requiredInfo?: string;
}

class AIValidationService {
  private baseUrl = '/api';

  async validateInput(request: ValidationRequest): Promise<ValidationResponse> {
    try {
      const response = await fetch(`${this.baseUrl}/validate-ai-input`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify(request),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      return await response.json();
    } catch (error) {
      console.error('AI validation error:', error);
      // Return a default "valid" response if service is unavailable
      return {
        isValid: true,
        message: 'Validation service temporarily unavailable. Your input has been accepted.',
      };
    }
  }
}

export const aiValidationService = new AIValidationService();
export type { ValidationRequest, ValidationResponse };
