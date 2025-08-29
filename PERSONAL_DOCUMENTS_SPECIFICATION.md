# Personal Documents Ontology Specification

## Overview
This comprehensive ontology models the complete domain of personal, legal, and official documents for the CLIENT-UX Personal Data Manager system. It provides a semantic framework for document classification, metadata extraction, and relationship mapping.

## SECTION 2: COMPLETE PROPERTY LISTING

### Datatype Properties (with XSD ranges)

#### Document Identification Properties
- `documentID` → `xsd:string` - Unique identifier for the document
- `documentNumber` → `xsd:string` - Official number assigned to the document
- `documentTitle` → `xsd:string` - Official title or name of the document
- `documentType` → `xsd:string` - Category or classification of the document
- `documentStatus` → `xsd:string` - Current status (active, expired, suspended, etc.)

#### Temporal Properties
- `issueDate` → `xsd:date` - Date when the document was issued
- `expiryDate` → `xsd:date` - Date when the document expires
- `validFrom` → `xsd:date` - Date from which the document is valid
- `validUntil` → `xsd:date` - Date until which the document is valid
- `graduationDate` → `xsd:date` - Date of graduation or completion

#### Personal Information Properties
- `fullName` → `xsd:string` - Complete legal name of the person
- `firstName` → `xsd:string` - Given name of the person
- `middleName` → `xsd:string` - Middle name or initial of the person
- `lastName` → `xsd:string` - Family name or surname of the person
- `dateOfBirth` → `xsd:date` - Birth date of the person
- `placeOfBirth` → `xsd:string` - Location where the person was born
- `gender` → `xsd:string` - Gender of the person
- `nationality` → `xsd:string` - Nationality or citizenship of the person

#### Financial Properties
- `amount` → `xsd:decimal` - Monetary amount
- `currency` → `xsd:string` - Currency code (ISO 4217)
- `accountNumber` → `xsd:string` - Financial account identifier

#### Educational Properties
- `degreeName` → `xsd:string` - Name of the academic degree or qualification
- `fieldOfStudy` → `xsd:string` - Academic discipline or major
- `gpa` → `xsd:decimal` - Grade Point Average

#### Vehicle Properties
- `vehicleIdentificationNumber` → `xsd:string` - Vehicle Identification Number (VIN)
- `licensePlate` → `xsd:string` - Vehicle registration plate number
- `vehicleMake` → `xsd:string` - Manufacturer of the vehicle
- `vehicleModel` → `xsd:string` - Model of the vehicle
- `vehicleYear` → `xsd:gYear` - Year the vehicle was manufactured

#### Address Properties
- `streetAddress` → `xsd:string` - Street number and name
- `city` → `xsd:string` - City or locality
- `stateProvince` → `xsd:string` - State, province, or region
- `postalCode` → `xsd:string` - ZIP code or postal code
- `countryCode` → `xsd:string` - ISO country code

### Object Properties (with Class ranges)

#### Document Relationships
- `documentHolder` → `:Person` - Person who holds or owns the document
- `issuedBy` → `:Organization` - Organization or authority that issued the document
- `issuedIn` → `:Country` - Country where the document was issued
- `relatedTo` → `:PersonalDocument` - Another document that this document is related to
- `supersedes` → `:PersonalDocument` - Previous version of document that this document replaces

#### Person Relationships
- `hasAddress` → `:Address` - Address associated with the person
- `citizenOf` → `:Country` - Country of citizenship
- `residentOf` → `:Country` - Country of residence

#### Educational Relationships
- `awardedBy` → `:EducationalInstitution` - Institution that awarded the degree or certificate
- `attendedBy` → `:Person` - Person who attended the educational program

#### Employment Relationships
- `employedBy` → `:Organization` - Organization that employs the person
- `employee` → `:Person` - Person who is the employee

#### Financial Relationships
- `accountHolderOf` → `:FinancialInstitution` - Financial institution where account is held
- `accountHolder` → `:Person` - Person who holds the financial account

#### Property Relationships
- `propertyOwner` → `:Person` - Person who owns the property
- `propertyLocation` → `:Address` - Address of the property

#### Medical Relationships
- `treatedBy` → `:HealthcareProvider` - Healthcare provider who provided treatment
- `patient` → `:Person` - Person who is the patient

#### Vehicle Relationships
- `vehicleOwner` → `:Person` - Person who owns the vehicle
- `registeredIn` → `:Country` - Country where the vehicle is registered

#### Insurance Relationships
- `insuredPerson` → `:Person` - Person covered by the insurance
- `insuranceProvider` → `:Organization` - Company providing the insurance coverage

#### Address Relationships
- `locatedIn` → `:Country` - Country where the address is located

## SECTION 3: JSON-LD EXTRACTION FORM EXAMPLES

### Example 1: Passport Document
```json
{
  "@context": "http://example.org/ontology/personal_documents#",
  "@type": "Passport",
  "documentID": "passport_001",
  "documentNumber": "123456789",
  "documentTitle": "United Kingdom Passport",
  "documentType": "Travel Document",
  "documentStatus": "Active",
  "issueDate": "2020-05-15",
  "expiryDate": "2030-05-14",
  "validFrom": "2020-05-15",
  "validUntil": "2030-05-14",
  "documentHolder": {
    "@type": "Person",
    "fullName": "Vincent Gerard Power",
    "firstName": "Vincent",
    "middleName": "Gerard", 
    "lastName": "Power",
    "dateOfBirth": "1961-04-09",
    "placeOfBirth": "London, United Kingdom",
    "gender": "Male",
    "nationality": "British",
    "hasAddress": {
      "@type": "Address",
      "streetAddress": "123 Main Street",
      "city": "London",
      "stateProvince": "England",
      "postalCode": "SW1A 1AA",
      "locatedIn": {
        "@type": "Country",
        "countryCode": "GB"
      }
    }
  },
  "issuedBy": {
    "@type": "GoverningBody",
    "name": "HM Passport Office"
  },
  "issuedIn": {
    "@type": "Country",
    "countryCode": "GB"
  }
}
```

### Example 2: University Degree Certificate
```json
{
  "@context": "http://example.org/ontology/personal_documents#",
  "@type": "UniversityDegree",
  "documentID": "degree_001",
  "documentNumber": "DEG-2018-CS-001234",
  "documentTitle": "Bachelor of Science in Computer Science",
  "documentType": "Academic Degree",
  "documentStatus": "Awarded",
  "issueDate": "2018-06-15",
  "graduationDate": "2018-06-15",
  "degreeName": "Bachelor of Science",
  "fieldOfStudy": "Computer Science",
  "gpa": 3.75,
  "documentHolder": {
    "@type": "Person",
    "fullName": "Vincent Gerard Power",
    "firstName": "Vincent",
    "lastName": "Power",
    "dateOfBirth": "1961-04-09"
  },
  "awardedBy": {
    "@type": "EducationalInstitution",
    "name": "University of London",
    "hasAddress": {
      "@type": "Address",
      "city": "London",
      "countryCode": "GB"
    }
  },
  "issuedIn": {
    "@type": "Country",
    "countryCode": "GB"
  }
}
```

### Example 3: Driver's License
```json
{
  "@context": "http://example.org/ontology/personal_documents#",
  "@type": "DriversLicense",
  "documentID": "license_001",
  "documentNumber": "POWER123456VG9GH",
  "documentTitle": "UK Driving Licence",
  "documentType": "Identity Document",
  "documentStatus": "Active",
  "issueDate": "2019-03-10",
  "expiryDate": "2029-03-09",
  "validFrom": "2019-03-10",
  "validUntil": "2029-03-09",
  "documentHolder": {
    "@type": "Person",
    "fullName": "Vincent Gerard Power",
    "firstName": "Vincent",
    "lastName": "Power",
    "dateOfBirth": "1961-04-09",
    "hasAddress": {
      "@type": "Address",
      "streetAddress": "123 Main Street",
      "city": "London",
      "postalCode": "SW1A 1AA"
    }
  },
  "issuedBy": {
    "@type": "GoverningBody",
    "name": "Driver and Vehicle Licensing Agency"
  },
  "issuedIn": {
    "@type": "Country",
    "countryCode": "GB"
  }
}
```

### Example 4: Bank Statement
```json
{
  "@context": "http://example.org/ontology/personal_documents#",
  "@type": "BankStatement",
  "documentID": "statement_001",
  "documentNumber": "STMT-2024-08-001",
  "documentTitle": "Monthly Bank Statement",
  "documentType": "Financial Document",
  "documentStatus": "Final",
  "issueDate": "2024-08-31",
  "validFrom": "2024-08-01",
  "validUntil": "2024-08-31",
  "accountNumber": "12345678",
  "amount": 5247.83,
  "currency": "GBP",
  "accountHolder": {
    "@type": "Person",
    "fullName": "Vincent Gerard Power",
    "firstName": "Vincent",
    "lastName": "Power"
  },
  "accountHolderOf": {
    "@type": "FinancialInstitution",
    "name": "Barclays Bank PLC"
  },
  "issuedBy": {
    "@type": "FinancialInstitution", 
    "name": "Barclays Bank PLC"
  }
}
```

### Example 5: Vehicle Registration
```json
{
  "@context": "http://example.org/ontology/personal_documents#",
  "@type": "VehicleRegistration",
  "documentID": "registration_001",
  "documentNumber": "V5C-ABC123DEF",
  "documentTitle": "Vehicle Registration Certificate",
  "documentType": "Vehicle Document",
  "documentStatus": "Active",
  "issueDate": "2022-01-15",
  "vehicleIdentificationNumber": "WBA12345678901234",
  "licensePlate": "AB12 CDE",
  "vehicleMake": "BMW",
  "vehicleModel": "320i",
  "vehicleYear": "2021",
  "vehicleOwner": {
    "@type": "Person",
    "fullName": "Vincent Gerard Power",
    "firstName": "Vincent",
    "lastName": "Power",
    "hasAddress": {
      "@type": "Address",
      "streetAddress": "123 Main Street",
      "city": "London",
      "postalCode": "SW1A 1AA"
    }
  },
  "issuedBy": {
    "@type": "GoverningBody",
    "name": "Driver and Vehicle Licensing Agency"
  },
  "registeredIn": {
    "@type": "Country",
    "countryCode": "GB"
  }
}
```

### Example 6: Medical Certificate
```json
{
  "@context": "http://example.org/ontology/personal_documents#",
  "@type": "MedicalCertificate",
  "documentID": "medical_cert_001",
  "documentNumber": "MED-2024-08-001",
  "documentTitle": "Fitness to Drive Medical Certificate",
  "documentType": "Medical Certificate",
  "documentStatus": "Valid",
  "issueDate": "2024-08-15",
  "expiryDate": "2025-08-14",
  "validFrom": "2024-08-15",
  "validUntil": "2025-08-14",
  "patient": {
    "@type": "Person",
    "fullName": "Vincent Gerard Power",
    "firstName": "Vincent",
    "lastName": "Power",
    "dateOfBirth": "1961-04-09"
  },
  "treatedBy": {
    "@type": "HealthcareProvider",
    "name": "Dr. Sarah Johnson",
    "hasAddress": {
      "@type": "Address",
      "streetAddress": "456 Medical Centre",
      "city": "London",
      "postalCode": "W1A 0AX"
    }
  },
  "issuedBy": {
    "@type": "HealthcareProvider",
    "name": "London Medical Centre"
  }
}
```

## Integration with CLIENT-UX System

This ontology is designed to integrate seamlessly with the existing CLIENT-UX Personal Data Manager system:

1. **TTL-Driven Architecture**: Follows the established pattern of using TTL files as the single source of truth
2. **OCR Integration**: Supports the existing passport/document OCR pipeline with structured extraction targets
3. **Form Generation**: Can drive dynamic form generation through the existing `ttl_parser.go` system
4. **Semantic Validation**: Enables SHACL-based validation of extracted document data
5. **Multi-domain Support**: Extends the insurance domain with comprehensive personal document management

## Usage in Document Processing Pipeline

1. **Document Upload**: User uploads document via existing upload interface
2. **OCR Processing**: Existing Tesseract/PassportEye pipeline extracts text
3. **Semantic Classification**: Document type identified using ontology classes
4. **Structured Extraction**: OCR results mapped to appropriate JSON-LD structure
5. **Validation**: SHACL shapes validate extracted data completeness and accuracy
6. **Storage**: Structured data stored with semantic relationships preserved
7. **Search & Retrieval**: Semantic queries enable powerful document search capabilities

This comprehensive ontology provides the foundation for a complete personal document management system within the CLIENT-UX platform.
