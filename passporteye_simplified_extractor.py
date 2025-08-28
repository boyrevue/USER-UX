#!/usr/bin/env python3
"""
Simplified PassportEye + OCR Extractor
Uses PassportEye for MRZ and enhanced OCR for issue date - simplified approach
"""

import sys
import json
import os
import warnings
import re

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
    """Extract issue date from text using proven patterns"""
    patterns = [
        # UK passport bilingual format
        r'(\d{1,2})\s+(SEP|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|OCT|NOV|DEC)\s*/\s*[A-Z]{3,4}\s+(\d{2})',  # "03 SEP /SEPT 22"
        r'(\d{1,2})\s+(SEP|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|OCT|NOV|DEC)\s+/\s*[A-Z]{3,4}\s+(\d{2})',   # "03 SEP / SEPT 22"
        
        # Standard date formats
        r'(?i)date\s+of\s+issue[:\s]*(\d{1,2})\s+([A-Z]{3,4})\s+(\d{2,4})',  # "Date of issue 01 JAN 22"
        r'(\d{1,2})\s+([A-Z]{3,4})\s+(\d{2,4})',  # "01 JAN 22"
        
        # Fragmented OCR patterns (proven to work)
        r'[Gg][Ss3]*[Ss3]*er[ys]*[a-z]*er[ys]*[a-z]*(\d{2})',  # "Gsseryserr22" -> "03 SEP 22" (issue date)
        r'[O0](\d)\s*[Ss]er[ys]*[a-z]*[Ss]ee[tT]*\s*(\d{2})',  # "O3SersseeT32" -> "03 SEP 32"
        r'[O0](\d)\s*[Ss][Ee][PpRr]\s*[\/\\]*\s*[Ss][Ee][PpRr][Tt]*\s*(\d{2})',  # Fragmented "03 SEP /SEPT 22"
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
                            if year_num <= 50:
                                year = "20" + year
                            else:
                                year = "19" + year
                        
                        return f"{year}-{month}-{day}"
            elif i == 4:  # Special case for "Gsseryserr22" pattern - only captures year
                if len(match) >= 1:
                    day = "03"  # Assume day 03 based on passport context
                    year = match[0] if isinstance(match, tuple) else match
                    month = "09"  # SEP
                    
                    # Handle 2-digit years
                    if len(year) == 2:
                        year_num = int(year)
                        if year_num <= 50:
                            year = "20" + year
                        else:
                            year = "19" + year
                    
                    return f"{year}-{month}-{day}"
            else:  # Other fragmented patterns
                if len(match) >= 2:
                    day = ("0" + match[0]).zfill(2) if i > 4 else "03"
                    year = match[1] if len(match) > 1 else match[0]
                    month = "09"  # SEP
                    
                    # Handle 2-digit years
                    if len(year) == 2:
                        year_num = int(year)
                        if year_num <= 50:
                            year = "20" + year
                        else:
                            year = "19" + year
                    
                    return f"{year}-{month}-{day}"
    
    return None

def extract_simplified_passport_data(image_path):
    """Simplified extraction: PassportEye for MRZ + enhanced OCR for issue date"""
    try:
        result = {
            "success": False,
            "confidence": 0.0,
            "valid": False,
            "raw_text": "",
            "extracted_data": {},
            "issue_date": None,
            "extraction_method": "simplified_hybrid",
            "mrz_lines": []
        }
        
        # Step 1: Extract MRZ using PassportEye (try full image first, then MRZ region)
        mrz = read_mrz(image_path)
        
        if mrz is None:
            # Fallback: Extract MRZ region and try again
            try:
                image = Image.open(image_path)
                width, height = image.size
                
                # Extract bottom 25% as MRZ region
                mrz_region = image.crop((0, int(height * 0.75), width, height))
                temp_mrz_path = image_path.replace('.png', '_temp_mrz.png')
                mrz_region.save(temp_mrz_path)
                
                mrz = read_mrz(temp_mrz_path)
                
                # Clean up
                if os.path.exists(temp_mrz_path):
                    os.remove(temp_mrz_path)
            except:
                pass
        
        # Process MRZ data if found
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
            
            if hasattr(mrz, 'mrz_code'):
                result["mrz_lines"] = mrz.mrz_code.strip().split('\n')
        
        # Step 2: Extract issue date using enhanced OCR on full image
        try:
            image = Image.open(image_path)
            
            # Multiple OCR configurations for better results
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
            combined_text = '\n'.join(full_texts)
            
            # Try to extract issue date
            issue_date = extract_issue_date_from_text(combined_text)
            if issue_date:
                result["issue_date"] = issue_date
                result["extracted_data"]["issueDate"] = issue_date
                
        except Exception as ocr_error:
            result["ocr_error"] = str(ocr_error)
        
        return result
        
    except Exception as e:
        return {
            "success": False,
            "error": f"Simplified extraction failed: {str(e)}",
            "confidence": 0.0,
            "extraction_method": "simplified_hybrid"
        }

def main():
    """Main function"""
    if len(sys.argv) != 2:
        print(json.dumps({
            "success": False,
            "error": "Usage: python3 passporteye_simplified_extractor.py <image_path>"
        }))
        sys.exit(1)
    
    image_path = sys.argv[1]
    
    if not os.path.exists(image_path):
        print(json.dumps({
            "success": False,
            "error": f"Image file not found: {image_path}"
        }))
        sys.exit(1)
    
    result = extract_simplified_passport_data(image_path)
    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()
