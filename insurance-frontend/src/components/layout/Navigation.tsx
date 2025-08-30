import React from 'react';



interface NavigationProps {
  currentStep: number;
  setCurrentStep: (step: number) => void;
  steps: string[];
}

const Navigation: React.FC<NavigationProps> = ({ 
  currentStep, 
  setCurrentStep, 
  steps 
}) => {
  return (
    <nav className="navigation">
      {steps.map((step, index) => (
        <button
          key={index}
          onClick={() => setCurrentStep(index)}
          className={`nav-button ${currentStep === index ? 'active' : ''}`}
        >
          {step}
        </button>
      ))}
    </nav>
  );
};


export default Navigation;
