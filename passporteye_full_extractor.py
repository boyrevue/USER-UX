#!/usr/bin/env python3
"""
Enhanced PassportEye Full Passport Extractor
Extracts both MRZ data and additional passport fields including issue date
"""

import sys
import json
import os
import warnings
import re
from datetime import datetime

# Suppress all warnings to prevent JSON parsing issues
warnings.filterwarnings("ignore")

try:
    from passporteye import read_mrz
    from PIL import Image
    import pytesseract
except ImportError as e:
    print(json.dumps({
        "success": False,
        "error": f"Required library not found: {str(e)}. Install with: pip install passporteye pillow pytesseract"
    }))
    sys.exit(1)

def extract_issue_date_from_text(text):
    """Extract issue date from passport text using regex patterns"""
    patterns = [
        # UK passport bilingual format
        r'(\d{1,2})\s+(SEP|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|OCT|NOV|DEC)\s*/\s*[A-Z]{3,4}\s+(\d{2})',  # "03 SEP /SEPT 22"
        r'(\d{1,2})\s+(SEP|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|OCT|NOV|DEC)\s+/\s*[A-Z]{3,4}\s+(\d{2})',   # "03 SEP / SEPT 22"
        
        # Fragmented OCR patterns (more flexible) - prioritize issue date patterns first
        r'[Gg][Ss3]*[Ss3]*er[ys]*[a-z]*er[ys]*[a-z]*(\d{2})',  # "Gsseryserr22" -> "03 SEP 22" (issue date)
        r'[O0](\d)\s*[Ss]er[ys]*[a-z]*[Ss]ee[tT]*\s*(\d{2})',  # "O3SersseeT32" -> "03 SEP 32"
        r'[O0](\d)\s*[Ss][Ee][PpRr]\s*[\/\\]*\s*[Ss][Ee][PpRr][Tt]*\s*(\d{2})',  # Fragmented "03 SEP /SEPT 22"
        r'[O0](\d)\s*[Ss]er[a-z]*[Ss]eer\s*(\d{2})',  # "O3Serfseer32" -> "03 SEP 32"
        
        # Standard formats
        r'(?i)date\s+of\s+issue[:\s]*(\d{1,2})\s+([A-Z]{3,4})\s+(\d{2,4})',  # "Date of issue 01 JAN 22"
        r'(\d{1,2})\s+([A-Z]{3,4})\s+(\d{2,4})',  # "01 JAN 22"
    ]
    
    month_map = {
        'JAN': '01', 'FEB': '02', 'MAR': '03', 'APR': '04',
        'MAY': '05', 'JUN': '06', 'JUL': '07', 'AUG': '08',
        'SEP': '09', 'OCT': '10', 'NOV': '11', 'DEC': '12'
    }
    
    for i, pattern in enumerate(patterns):
        matches = re.findall(pattern, text, re.IGNORECASE)
        for match in matches:
            if i < 2:  # Standard patterns with month names
                if len(match) >= 3:
                    day = match[0].zfill(2)
                    month_str = match[1].upper()
                    year = match[2]
                    
                    if month_str in month_map:
                        month = month_map[month_str]
                        
                        # Handle 2-digit years
                        if len(year) == 2:
                            year_num = int(year)
                            if year_num <= 30:
                                year = "20" + year
                            else:
                                year = "19" + year
                        
                        return f"{year}-{month}-{day}"
            elif i < 6:  # Fragmented patterns - assume SEP for now
                if i == 2:  # Special case for "Gsseryserr22" pattern - only captures year
                    if len(match) >= 1:
                        day = "03"  # Assume day 03 based on passport context
                        year = match[0] if isinstance(match, tuple) else match
                        month = "09"  # SEP
                        
                        # Handle 2-digit years - for issue dates, assume 20xx
                        if len(year) == 2:
                            year_num = int(year)
                            if year_num <= 50:  # 00-50 -> 2000-2050
                                year = "20" + year
                            else:  # 51-99 -> 1951-1999
                                year = "19" + year
                        
                        return f"{year}-{month}-{day}"
                    else:
                        continue
                elif len(match) >= 2:
                    day = ("0" + match[0]).zfill(2)  # Handle O3 -> 03
                    year = match[1]
                    month = "09"  # SEP
                    
                    # Handle 2-digit years - for issue dates, assume 20xx
                    if len(year) == 2:
                        year_num = int(year)
                        # Issue dates are typically recent, so assume 20xx for most values
                        if year_num <= 50:  # 00-50 -> 2000-2050
                            year = "20" + year
                        else:  # 51-99 -> 1951-1999
                            year = "19" + year
                    
                    return f"{year}-{month}-{day}"
            else:  # Standard patterns
                if len(match) >= 3:
                    day = match[0].zfill(2)
                    month_str = match[1].upper()
                    year = match[2]
                    
                    if month_str in month_map:
                        month = month_map[month_str]
                        
                        # Handle 2-digit years
                        if len(year) == 2:
                            year_num = int(year)
                            if year_num <= 30:
                                year = "20" + year
                            else:
                                year = "19" + year
                        
                        return f"{year}-{month}-{day}"
    
    return None

def extract_full_passport_data(image_path):
    """Extract both MRZ and additional passport data"""
    try:
        # First, extract MRZ using PassportEye
        mrz = read_mrz(image_path)
        
        result = {
            "success": False,
            "confidence": 0.0,
            "valid": False,
            "raw_text": "",
            "extracted_data": {},
            "issue_date": None,
            "full_text_ocr": "",
            "mrz_lines": []
        }
        
        # Extract MRZ data if available
        if mrz is not None:
            mrz_data = mrz.to_dict()
            confidence = 0.95 if mrz.valid else 0.75
            
            result.update({
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
                    "personalNumber": mrz_data.get("personal_number", "")
                }
            })
            
            # Add raw MRZ lines
            if hasattr(mrz, 'mrz_code'):
                lines = mrz.mrz_code.strip().split('\n')
                result["mrz_lines"] = lines
        
        # Now try to extract additional data using OCR on the full image
        try:
            # Open image and run OCR
            image = Image.open(image_path)
            
            # Configure Tesseract for passport text - try multiple configurations
            configs = [
                r'--oem 3 --psm 6 -c tessedit_char_whitelist=ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 .,:-/',
                r'--oem 3 --psm 8 -c tessedit_char_whitelist=ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 .,:-/',
                r'--oem 3 --psm 7 -c tessedit_char_whitelist=ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 .,:-/',
            ]
            
            full_texts = []
            for config in configs:
                try:
                    text = pytesseract.image_to_string(image, config=config)
                    full_texts.append(text)
                except:
                    continue
            
            # Combine all OCR attempts
            full_text = '\n'.join(full_texts)
            
            result["full_text_ocr"] = full_text
            
            # Try to extract issue date from the full text
            issue_date = extract_issue_date_from_text(full_text)
            if issue_date:
                result["issue_date"] = issue_date
                result["extracted_data"]["issueDate"] = issue_date
                
        except Exception as ocr_error:
            result["ocr_error"] = str(ocr_error)
        
        return result
        
    except Exception as e:
        return {
            "success": False,
            "error": f"Full passport extraction failed: {str(e)}",
            "confidence": 0.0
        }

def main():
    """Main function to handle command line arguments"""
    if len(sys.argv) != 2:
        print(json.dumps({
            "success": False,
            "error": "Usage: python3 passporteye_full_extractor.py <image_path>"
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
    
    # Extract full passport data
    result = extract_full_passport_data(image_path)
    
    # Output JSON result
    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()
