import React, { useEffect, useState } from 'react';
import { Modal, Button, Alert } from 'flowbite-react';
import { MessageCircle } from 'lucide-react';

interface GlobalAssistantProps {
  isOpen: boolean;
  onClose: () => void;
  onOpen?: () => void;
  onTriggerField?: (fieldName: string, prompt?: string) => void;
}

const GlobalAssistant: React.FC<GlobalAssistantProps> = ({ isOpen, onClose, onOpen, onTriggerField }) => {
  const [message, setMessage] = useState<string>('');
  const [noticeShown, setNoticeShown] = useState<boolean>(false);

  useEffect(() => {
    if (isOpen && !noticeShown) {
      setMessage('Hi! I\'m here to help you complete required fields. Select “Other” anywhere and I\'ll guide you with the right prompt.');
      setNoticeShown(true);
    }
  }, [isOpen, noticeShown]);

  return (
    <>
      {/* Floating presence bubble */}
      <div className="fixed bottom-4 right-4 z-50">
        <Button color="blue" onClick={() => { if (onOpen) { onOpen(); };
          // also broadcast a global event to open the first AI field
          window.dispatchEvent(new Event('open-ai-validation')); }}>
          <span className="flex items-center">
            <MessageCircle className="w-4 h-4 mr-2" />
            Assistant
          </span>
        </Button>
      </div>

      {/* Intro modal */}
      <Modal show={isOpen} size="md" onClose={onClose}>
        <div className="p-5">
          <div className="flex items-center mb-3">
            <MessageCircle className="w-5 h-5 mr-2 text-blue-600" />
            <h3 className="text-lg font-semibold">Form Assistant</h3>
          </div>
          <Alert color="info" className="mb-3">
            {message}
          </Alert>
          <div className="text-sm text-gray-600 mb-4">
            - Press Enter in any assistant dialog to validate.
            <br />- Shift+Enter makes a new line.
          </div>
          <div className="flex justify-end">
            <Button color="blue" onClick={onClose}>Got it</Button>
          </div>
        </div>
      </Modal>
    </>
  );
};

export default GlobalAssistant;


