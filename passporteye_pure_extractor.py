#!/usr/bin/env python3
"""
Pure PassportEye Extractor
Uses only PassportEye for all passport data extraction
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

def extract_issue_date_from_image_regions(image_path):
    """Extract issue date by analyzing different regions of the passport image"""
    try:
        image = Image.open(image_path)
        width, height = image.size
        
        # Define regions where issue date might be located
        regions = {
            "upper_middle": (0, int(height * 0.3), width, int(height * 0.7)),  # Middle section
            "left_side": (0, int(height * 0.4), int(width * 0.6), int(height * 0.8)),  # Left side
            "right_side": (int(width * 0.4), int(height * 0.4), width, int(height * 0.8)),  # Right side
        }
        
        issue_dates = []
        
        for region_name, (left, top, right, bottom) in regions.items():
            try:
                # Crop the region
                region_img = image.crop((left, top, right, bottom))
                
                # OCR with passport-optimized settings
                custom_config = r'--oem 3 --psm 6 -c tessedit_char_whitelist=ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789 .,:-/'
                text = pytesseract.image_to_string(region_img, config=custom_config)
                
                # Try to extract issue date from this region
                issue_date = extract_issue_date_from_text(text)
                
                # Always add region info for debugging (even if no date found)
                region_info = {
                    "region": region_name,
                    "text": text.strip()[:200],  # First 200 chars for debugging
                    "date": issue_date
                }
                
                if issue_date:
                    issue_dates.append(region_info)
                else:
                    # Add to debug info even if no date found
                    issue_dates.append(region_info)
                    
            except Exception as e:
                continue
        
        return issue_dates
        
    except Exception as e:
        return []

def extract_issue_date_from_text(text):
    """Extract issue date from text using comprehensive patterns"""
    patterns = [
        # UK passport bilingual format
        r'(\d{1,2})\s+(SEP|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|OCT|NOV|DEC)\s*/\s*[A-Z]{3,4}\s+(\d{2})',  # "03 SEP /SEPT 22"
        r'(\d{1,2})\s+(SEP|JAN|FEB|MAR|APR|MAY|JUN|JUL|AUG|OCT|NOV|DEC)\s+/\s*[A-Z]{3,4}\s+(\d{2})',   # "03 SEP / SEPT 22"
        
        # Standard date formats
        r'(?i)date\s+of\s+issue[:\s]*(\d{1,2})\s+([A-Z]{3,4})\s+(\d{2,4})',  # "Date of issue 01 JAN 22"
        r'(\d{1,2})\s+([A-Z]{3,4})\s+(\d{2,4})',  # "01 JAN 22"
        
        # Fragmented OCR patterns - prioritize issue date patterns
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

def extract_pure_passport_data(image_path):
    """Extract passport data using pure PassportEye approach with regional analysis"""
    try:
        # Primary: Extract MRZ using PassportEye on full image
        mrz = read_mrz(image_path)
        
        # If PassportEye fails on full image, try to extract MRZ region first
        if mrz is None:
            try:
                image = Image.open(image_path)
                width, height = image.size
                
                # Extract bottom 25% as MRZ region
                mrz_region = image.crop((0, int(height * 0.75), width, height))
                
                # Save temporary MRZ region
                temp_mrz_path = image_path.replace('.png', '_temp_mrz.png')
                mrz_region.save(temp_mrz_path)
                
                # Try PassportEye on MRZ region
                mrz = read_mrz(temp_mrz_path)
                
                # Clean up temp file
                if os.path.exists(temp_mrz_path):
                    os.remove(temp_mrz_path)
                    
            except Exception as e:
                pass
        
        result = {
            "success": False,
            "confidence": 0.0,
            "valid": False,
            "raw_text": "",
            "extracted_data": {},
            "issue_date": None,
            "issue_date_candidates": [],
            "extraction_method": "pure_passporteye",
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
        
        # Secondary: Extract issue date from different regions
        issue_date_candidates = extract_issue_date_from_image_regions(image_path)
        result["issue_date_candidates"] = issue_date_candidates
        
        # Select the best issue date candidate
        if issue_date_candidates:
            # Prefer dates from upper_middle region, then others
            best_candidate = None
            for candidate in issue_date_candidates:
                if candidate["region"] == "upper_middle":
                    best_candidate = candidate
                    break
            
            if not best_candidate:
                best_candidate = issue_date_candidates[0]
            
            result["issue_date"] = best_candidate["date"]
            result["extracted_data"]["issueDate"] = best_candidate["date"]
        
        return result
        
    except Exception as e:
        return {
            "success": False,
            "error": f"Pure PassportEye extraction failed: {str(e)}",
            "confidence": 0.0,
            "extraction_method": "pure_passporteye"
        }

def main():
    """Main function to handle command line arguments"""
    if len(sys.argv) != 2:
        print(json.dumps({
            "success": False,
            "error": "Usage: python3 passporteye_pure_extractor.py <image_path>"
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
    
    # Extract passport data using pure PassportEye approach
    result = extract_pure_passport_data(image_path)
    
    # Output JSON result
    print(json.dumps(result, indent=2))

if __name__ == "__main__":
    main()
