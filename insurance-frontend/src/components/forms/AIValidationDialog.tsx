import React, { useState, useEffect } from 'react';
import { Modal, Button, TextInput, Alert } from 'flowbite-react';
import { MessageCircle, AlertCircle, CheckCircle, HelpCircle } from 'lucide-react';
import { aiValidationService, ValidationResponse } from '../../services/aiValidationService';

interface AIValidationDialogProps {
  isOpen: boolean;
  onClose: () => void;
  fieldName: string;
  fieldLabel: string;
  initialValue: string;
  validationPrompt: string;
  onValidatedValue: (value: string) => void;
}

const AIValidationDialog: React.FC<AIValidationDialogProps> = ({
  isOpen,
  onClose,
  fieldName,
  fieldLabel,
  initialValue,
  validationPrompt,
  onValidatedValue,
}) => {
  const [currentValue, setCurrentValue] = useState(initialValue);
  const [validationResult, setValidationResult] = useState<ValidationResponse | null>(null);
  const [isValidating, setIsValidating] = useState(false);
  const [showHelp, setShowHelp] = useState(false);
  const [assistantReply, setAssistantReply] = useState<string>('');

  useEffect(() => {
    if (!isOpen) return;
    setAssistantReply('');
    // fire-and-forget assistant suggestion on open
    (async () => {
      try {
        const res = await fetch('/api/assistant/generate', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ field: fieldName, prompt: validationPrompt, userInput: currentValue || '' })
        });
        if (res.ok) {
          const j = await res.json();
          if (j?.reply) setAssistantReply(j.reply);
        }
      } catch (_) {}
    })();
  }, [isOpen]);

  const handleValidate = async () => {
    if (!currentValue.trim()) {
      setValidationResult({
        isValid: false,
        message: 'Please provide some information before validating.',
      });
      return;
    }

    setIsValidating(true);
    try {
      const result = await aiValidationService.validateInput({
        fieldName,
        userInput: currentValue,
        validationPrompt,
      });
      setValidationResult(result);
    } catch (error) {
      setValidationResult({
        isValid: false,
        message: 'Validation service error. Please try again.',
      });
    } finally {
      setIsValidating(false);
    }
  };

  const handleAccept = () => {
    if (validationResult?.isValid) {
      onValidatedValue(currentValue);
      onClose();
    }
  };

  const handleClose = () => {
    setValidationResult(null);
    setCurrentValue(initialValue);
    onClose();
  };

  const getValidationIcon = () => {
    if (!validationResult) return null;
    
    if (validationResult.isValid) {
      return <CheckCircle className="w-5 h-5 text-green-600" />;
    } else {
      return <AlertCircle className="w-5 h-5 text-red-600" />;
    }
  };

  const getValidationColor = () => {
    if (!validationResult) return 'info';
    return validationResult.isValid ? 'success' : 'failure';
  };

  return (
    <Modal show={isOpen} onClose={handleClose} size="lg">
      <div className="p-6">
        <div className="flex items-center mb-4">
          <MessageCircle className="w-5 h-5 mr-2 text-blue-600" />
          <h3 className="text-lg font-semibold">Insurance Information Validation</h3>
        </div>
        <div className="space-y-4">
          <div>
            <h4 className="font-medium text-gray-900 mb-2">{fieldLabel}</h4>
            <p className="text-sm text-gray-600 mb-4">
              Please provide specific details that would be relevant to an insurance company. 
              Vague or nonsensical answers will not be accepted.
            </p>
          </div>

          {assistantReply && (
            <Alert color="info" className="whitespace-pre-wrap">{assistantReply}</Alert>
          )}

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">
              Your Response:
            </label>
            <textarea
              value={currentValue}
              onChange={(e) => setCurrentValue(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter' && !e.shiftKey) {
                  e.preventDefault();
                  handleValidate();
                }
              }}
              className="w-full p-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500"
              rows={4}
              placeholder="Please provide specific details about your situation..."
            />
          </div>

          {validationResult && (
            <Alert color={getValidationColor()} icon={getValidationIcon}>
              <div>
                <p className="font-medium">{validationResult.message}</p>
                {validationResult.requiredInfo && (
                  <div className="mt-2">
                    <p className="text-sm font-medium">Required Information:</p>
                    <p className="text-sm">{validationResult.requiredInfo}</p>
                  </div>
                )}
                {validationResult.suggestions && (
                  <div className="mt-2">
                    <p className="text-sm font-medium">Suggestions:</p>
                    <p className="text-sm">{validationResult.suggestions}</p>
                  </div>
                )}
              </div>
            </Alert>
          )}

          <div className="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <button
              type="button"
              onClick={() => setShowHelp(!showHelp)}
              className="flex items-center text-blue-800 hover:text-blue-900"
            >
              <HelpCircle className="w-4 h-4 mr-2" />
              What kind of information do you need?
            </button>
            
            {showHelp && (
              <div className="mt-3 text-sm text-blue-700">
                <p className="font-medium mb-2">Insurance companies need specific details such as:</p>
                <ul className="list-disc list-inside space-y-1">
                  <li>Official reasons given by DVLA or other authorities</li>
                  <li>Specific dates when events occurred</li>
                  <li>Medical conditions or circumstances involved</li>
                  <li>Current status of your licence or situation</li>
                  <li>Any official codes, reference numbers, or documentation</li>
                </ul>
                <p className="mt-2 font-medium">
                  Avoid vague responses like "personal reasons" or "don't want to say"
                </p>
              </div>
            )}
          </div>
        </div>
        
        <div className="flex justify-between pt-4 border-t border-gray-200">
          <Button color="gray" onClick={handleClose}>
            Cancel
          </Button>
          
          <div className="flex space-x-2">
            <Button
              color="blue"
              onClick={handleValidate}
              disabled={isValidating || !currentValue.trim()}
            >
              {isValidating ? 'Validating...' : 'Press Enter to Validate'}
            </Button>
            
            {validationResult?.isValid && (
              <Button color="success" onClick={handleAccept}>
                Accept & Continue
              </Button>
            )}
          </div>
        </div>
      </div>
    </Modal>
  );
};

export default AIValidationDialog;
