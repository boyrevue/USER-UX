import React from 'react';
import { Badge } from 'flowbite-react';
import { CheckCircle, User, Car, FileText, Shield, Upload, Settings, Eye } from 'lucide-react';
import { NavigationProps } from '../../types';

const Navigation: React.FC<NavigationProps> = ({ 
  currentStep,
  setCurrentStep,
  steps,
  canNavigate,
  completedSteps
}) => {
  const getStepIcon = (stepName: string, index: number) => {
    const iconClass = "w-5 h-5";
    switch (stepName) {
      case 'Driver Details': return <User className={iconClass} />;
      case 'Vehicle Details': return <Car className={iconClass} />;
      case 'Claims History': return <FileText className={iconClass} />;
      case 'Policy Details': return <Shield className={iconClass} />;
      case 'Documents': return <Upload className={iconClass} />;
      case 'Settings': return <Settings className={iconClass} />;
      case 'Review & Submit': return <Eye className={iconClass} />;
      default: return <FileText className={iconClass} />;
    }
  };

  return (
    <nav className="w-64 bg-white border-r border-gray-200 h-screen fixed left-0 top-0 overflow-y-auto">
      <div className="p-6">
        <h2 className="text-xl font-bold text-gray-800 mb-6">CLIENT-UX Insurance</h2>
        
        <div className="space-y-2">
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
                  w-full flex items-center px-4 py-3 rounded-lg text-sm font-medium transition-all duration-200
                  ${isActive 
                    ? 'bg-blue-600 text-white shadow-lg' 
                    : isCompleted 
                      ? 'bg-green-50 text-green-700 hover:bg-green-100 border border-green-200' 
                      : isClickable
                        ? 'bg-gray-50 text-gray-700 hover:bg-gray-100 border border-gray-200'
                        : 'bg-gray-25 text-gray-400 cursor-not-allowed'
                  }
                `}
              >
                <div className="flex items-center w-full">
                  <div className="flex-shrink-0 mr-3">
                    {isCompleted && !isActive ? (
                      <CheckCircle className="w-5 h-5 text-green-600" />
                    ) : (
                      getStepIcon(step, index)
                    )}
                  </div>
                  
                  <div className="flex-1 text-left">
                    <div className="font-medium">{step}</div>
                    {isActive && (
                      <div className="text-xs opacity-75 mt-1">Current Step</div>
                    )}
                    {isCompleted && !isActive && (
                      <div className="text-xs text-green-600 mt-1">Completed</div>
                    )}
                  </div>
                  
                  {isActive && (
                    <Badge color="info" size="sm">
                      {index + 1}
                    </Badge>
                  )}
                </div>
              </button>
            );
          })}
        </div>
        
        <div className="mt-8 pt-6 border-t border-gray-200">
          <div className="text-xs text-gray-500 mb-2">Progress</div>
          <div className="w-full bg-gray-200 rounded-full h-2">
            <div 
              className="bg-blue-600 h-2 rounded-full transition-all duration-300"
              style={{ width: `${(completedSteps.length / steps.length) * 100}%` }}
            ></div>
          </div>
          <div className="text-xs text-gray-500 mt-2">
            {completedSteps.length} of {steps.length} completed
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navigation;
