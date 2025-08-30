
// Date validation utilities for 5-year historical data limit
const getDateLimits = () => {
  const today = new Date();
  const fiveYearsAgo = new Date();
  fiveYearsAgo.setFullYear(today.getFullYear() - 5);
  
  return {
    today: today.toISOString().split('T')[0],
    fiveYearsAgo: fiveYearsAgo.toISOString().split('T')[0],
    maxBirthDate: new Date(today.getFullYear() - 17, today.getMonth(), today.getDate()).toISOString().split('T')[0],
    minBirthDate: new Date(today.getFullYear() - 130, today.getMonth(), today.getDate()).toISOString().split('T')[0],
    earliestLicenceDate: new Date(1970, 0, 1).toISOString().split('T')[0],
    ukMinDrivingAge: 17,
    maxHumanAge: 130
  };
};

// Validation functions
export const validateBirthDate = (dateString: string): { valid: boolean; error?: string } => {
  // Validation logic here
  return { valid: true };
};

export const validateHistoricalDate = (dateString: string): { valid: boolean; error?: string } => {
  // Validation logic here
  return { valid: true };
};

export const validateLicenceDate = (dateString: string, birthDate: string): { valid: boolean; error?: string } => {
  // Validation logic here
  return { valid: true };
};

export { getDateLimits };
