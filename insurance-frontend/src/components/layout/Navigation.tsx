import React from 'react';
import { Badge } from 'flowbite-react';
import { CheckCircle } from 'lucide-react';
import { NavigationProps } from '../../types';

const Navigation: React.FC<NavigationProps> = ({ 
  currentStep,
  setCurrentStep,
  steps,
  canNavigate,
  completedSteps
}) => {
  return (
    <nav className="mb-8">
      <div className="flex flex-wrap justify-center gap-2">
        {steps.map((step, index) => {
          const isActive = currentStep === index;
          const isCompleted = completedSteps.includes(index);
          const isClickable = canNavigate(index);
          
          return (
            <button
              key={index}
              onClick={() => isClickable && setCurrentStep(index)}
              disabled={!isClickable}
              className={`
                flex items-center px-4 py-2 rounded-lg text-sm font-medium transition-colors
                ${isActive 
                  ? 'bg-blue-600 text-white' 
                  : isCompleted 
                    ? 'bg-green-100 text-green-800 hover:bg-green-200' 
                    : isClickable
                      ? 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                      : 'bg-gray-50 text-gray-400 cursor-not-allowed'
                }
              `}
            >
              {isCompleted && <CheckCircle className="w-4 h-4 mr-2" />}
              <span>{step}</span>
              {isActive && <Badge color="info" className="ml-2">Current</Badge>}
            </button>
          );
        })}
      </div>
    </nav>
  );
};

export default Navigation;
