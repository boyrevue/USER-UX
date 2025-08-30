
// Date validation utilities with UK driving age and realistic human age limits
export const getDateLimits = () => {
  const today = new Date();
  const fiveYearsAgo = new Date();
  fiveYearsAgo.setFullYear(today.getFullYear() - 5);
  
  // UK driving age restrictions and human age limits
  const ukMinDrivingAge = 17; // UK minimum driving age
  const maxHumanAge = 130; // Maximum realistic human age
  const minDrivingDate = new Date(today.getFullYear() - ukMinDrivingAge, today.getMonth(), today.getDate());
  const maxAgeDate = new Date(today.getFullYear() - maxHumanAge, today.getMonth(), today.getDate());
  
  return {
    today: today.toISOString().split('T')[0],
    fiveYearsAgo: fiveYearsAgo.toISOString().split('T')[0],
    
    // Birth date limits (17-130 years old)
    maxBirthDate: minDrivingDate.toISOString().split('T')[0], // Must be at least 17 years old
    minBirthDate: maxAgeDate.toISOString().split('T')[0], // Cannot be older than 130 years
    
    // UK driving licence issue date limits
    earliestLicenceDate: new Date(1970, 0, 1).toISOString().split('T')[0], // UK DVLA records start ~1970
    
    // Age validation constants
    ukMinDrivingAge,
    maxHumanAge
  };
};

// Validation function for birth dates (UK driving age and human age limits)
export const validateBirthDate = (dateString: string): { valid: boolean; error?: string } => {
  if (!dateString) return { valid: true }; // Optional dates are valid when empty
  
  const inputDate = new Date(dateString);
  const { today, maxBirthDate, minBirthDate, ukMinDrivingAge, maxHumanAge } = getDateLimits();
  const todayDate = new Date(today);
  
  // Calculate age
  const age = todayDate.getFullYear() - inputDate.getFullYear();
  const monthDiff = todayDate.getMonth() - inputDate.getMonth();
  const dayDiff = todayDate.getDate() - inputDate.getDate();
  const actualAge = monthDiff < 0 || (monthDiff === 0 && dayDiff < 0) ? age - 1 : age;
  
  if (inputDate > todayDate) {
    return { 
      valid: false, 
      error: "📅 Future birth dates are not allowed. Please enter a valid birth date." 
    };
  }
  
  if (actualAge < ukMinDrivingAge) {
    return { 
      valid: false, 
      error: `🚗 Driver must be at least ${ukMinDrivingAge} years old to drive in the UK. This person is only ${actualAge} years old.` 
    };
  }
  
  if (actualAge > maxHumanAge) {
    return { 
      valid: false, 
      error: `👴 Age cannot exceed ${maxHumanAge} years. Please check the birth date (calculated age: ${actualAge} years).` 
    };
  }
  
  return { valid: true };
};

// Validation function for historical dates (claims, accidents, convictions)
export const validateHistoricalDate = (dateString: string): { valid: boolean; error?: string } => {
  if (!dateString) return { valid: true }; // Optional dates are valid when empty
  
  const inputDate = new Date(dateString);
  const { today, fiveYearsAgo } = getDateLimits();
  const todayDate = new Date(today);
  const fiveYearsAgoDate = new Date(fiveYearsAgo);
  
  if (inputDate > todayDate) {
    return { 
      valid: false, 
      error: "📅 Future dates are not allowed for historical events. Please select a date in the past." 
    };
  }
  
  if (inputDate < fiveYearsAgoDate) {
    return { 
      valid: false, 
      error: `⚠️ This date is more than 5 years old (before ${fiveYearsAgo}). For insurance purposes, we only consider events from the last 5 years.` 
    };
  }
  
  return { valid: true };
};

// Validation function for licence dates
export const validateLicenceDate = (dateString: string, birthDate: string): { valid: boolean; error?: string } => {
  if (!dateString) return { valid: true }; // Optional dates are valid when empty
  
  const inputDate = new Date(dateString);
  const { today, earliestLicenceDate, ukMinDrivingAge } = getDateLimits();
  const todayDate = new Date(today);
  const earliestDate = new Date(earliestLicenceDate);
  
  if (inputDate > todayDate) {
    return { 
      valid: false, 
      error: "📅 Future licence dates are not allowed. Please select a date in the past." 
    };
  }
  
  if (inputDate < earliestDate) {
    return { 
      valid: false, 
      error: `📋 UK DVLA records start from 1970. Please enter a licence date after ${earliestLicenceDate}.` 
    };
  }
  
  // Check if licence was issued when person was old enough to drive
  if (birthDate) {
    const birthDateObj = new Date(birthDate);
    const licenceDateObj = new Date(dateString);
    const ageAtLicence = licenceDateObj.getFullYear() - birthDateObj.getFullYear();
    const monthDiff = licenceDateObj.getMonth() - birthDateObj.getMonth();
    const dayDiff = licenceDateObj.getDate() - birthDateObj.getDate();
    const actualAgeAtLicence = monthDiff < 0 || (monthDiff === 0 && dayDiff < 0) ? ageAtLicence - 1 : ageAtLicence;
    
    if (actualAgeAtLicence < ukMinDrivingAge) {
      return { 
        valid: false, 
        error: `🚗 Licence cannot be issued before age ${ukMinDrivingAge}. This licence was issued when the person was ${actualAgeAtLicence} years old.` 
      };
    }
  }
  
  return { valid: true };
};
