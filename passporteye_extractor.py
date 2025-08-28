#!/usr/bin/env python3
"""
PassportEye MRZ Extractor
Specialized MRZ extraction using PassportEye library
"""

import sys
import json
import os
import warnings

# Suppress all warnings to prevent JSON parsing issues
warnings.filterwarnings("ignore")

from passporteye import read_mrz

def extract_mrz_with_passporteye(image_path):
    """Extract MRZ using PassportEye library"""
    try:
        # Read MRZ from image
        mrz = read_mrz(image_path)
        
        if mrz is None:
            return {
                "success": False,
                "error": "No MRZ found in image",
                "confidence": 0.0
            }
        
        # Extract data from MRZ
        mrz_data = mrz.to_dict()
        
        # Calculate confidence based on validity
        confidence = 0.95 if mrz.valid else 0.75
        
        # Format the response
        result = {
            "success": True,
            "confidence": confidence,
            "valid": mrz.valid,
            "raw_text": mrz.mrz_code if hasattr(mrz, 'mrz_code') else "",
            "extracted_data": {
                "documentType": mrz_data.get("type", ""),
                "issuingCountry": mrz_data.get("country", ""),
                "surname": mrz_data.get("surname", ""),
                "givenNames": mrz_data.get("names", ""),
                "passportNumber": mrz_data.get("number", ""),
                "nationality": mrz_data.get("nationality", ""),
                "dateOfBirth": mrz_data.get("date_of_birth", ""),
                "gender": mrz_data.get("sex", ""),
                "expiryDate": mrz_data.get("expiration_date", ""),
                "personalNumber": mrz_data.get("personal_number", ""),
                "checkDigits": {
                    "number": mrz_data.get("check_number", ""),
                    "dateOfBirth": mrz_data.get("check_date_of_birth", ""),
                    "expiryDate": mrz_data.get("check_expiration_date", ""),
                    "personalNumber": mrz_data.get("check_personal_number", ""),
                    "composite": mrz_data.get("check_composite", "")
                }
            },
            "validation": {
                "valid": mrz.valid,
                "errors": []
            }
        }
        
        # Add validation errors if any
        if hasattr(mrz, 'report') and mrz.report:
            result["validation"]["errors"] = [str(error) for error in mrz.report]
        
        # Add raw MRZ lines
        if hasattr(mrz, 'mrz_code'):
            lines = mrz.mrz_code.strip().split('\n')
            result["mrz_lines"] = lines
        
        return result
        
    except Exception as e:
        return {
            "success": False,
            "error": f"PassportEye extraction failed: {str(e)}",
            "confidence": 0.0
        }

def main():
    """Main function to handle command line arguments"""
    if len(sys.argv) != 2:
        print(json.dumps({
            "success": False,
            "error": "Usage: python3 passporteye_extractor.py <image_path>"
        }))
        sys.exit(1)
    
    image_path = sys.argv[1]
    
    # Check if file exists
    if not os.path.exists(image_path):
        print(json.dumps({
            "success": False,
            "error": f"Image file not found: {image_path}"
        }))
        sys.exit(1)
    
    # Extract MRZ
    result = extract_mrz_with_passporteye(image_path)
    
    # Output JSON result
    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()
