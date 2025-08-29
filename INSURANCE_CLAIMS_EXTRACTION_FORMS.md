# Insurance Claims & Risk Assessment - Extraction Forms

## [SECTION 2: EXTRACTION FORMS FOR UNDERWRITING]

These JSON-LD forms represent the exact data an insurance company's system needs to extract from customers or prior insurers to process new policy applications. The forms map directly to the ontology properties defined in `insurance_claims_ontology.ttl`.

### Form 1: Insurance Claim History (Formal Claims)

```json
{
  "@context": {
    "@vocab": "http://example.org/ontology/insurance_claims#",
    "xsd": "http://www.w3.org/2001/XMLSchema#"
  },
  "@type": "Claim",
  "claimNumber": "CLM-2023-456789",
  "claimDate": "2023-11-15",
  "incidentDate": "2023-11-12", 
  "settlementDate": "2024-01-20",
  "claimType": {
    "@type": "skos:Concept",
    "@id": "CollisionClaim",
    "prefLabel": "Collision"
  },
  "claimAmount": 4250.00,
  "paidAmount": 3850.00,
  "deductibleAmount": 500.00,
  "reserveAmount": 4500.00,
  "faultStatus": {
    "@type": "skos:Concept", 
    "@id": "AtFault",
    "prefLabel": "At-fault"
  },
  "claimStatus": {
    "@type": "skos:Concept",
    "@id": "ClosedPaid", 
    "prefLabel": "Closed-Paid"
  },
  "description": "Reversed into a concrete pole in shopping center parking lot. Damage to rear bumper, trunk lid, and rear lights. No injuries reported.",
  "adjusterNotes": "Clear liability - policyholder admitted fault. Damage consistent with low-speed reverse collision. No suspicious circumstances.",
  "paidByInsurer": {
    "@type": "InsuranceCompany",
    "companyName": "Previous Auto Insurance Co.",
    "naic_code": "12345"
  },
  "madeAgainstPolicy": {
    "@type": "Policy", 
    "policyNumber": "POL-2023-ABC123",
    "policyEffectiveDate": "2023-06-01",
    "policyExpirationDate": "2024-06-01"
  },
  "involvesDamage": {
    "@type": "Damage",
    "damageDescription": "Rear bumper cracked and dented, trunk lid misaligned, both rear taillights cracked",
    "estimatedDamageAmount": 4250.00,
    "actualRepairCost": 3850.00
  }
}
```

### Form 2: Accident/Incident (No Claim Filed)

```json
{
  "@context": {
    "@vocab": "http://example.org/ontology/insurance_claims#",
    "xsd": "http://www.w3.org/2001/XMLSchema#"
  },
  "@type": "Accident",
  "accidentDate": "2024-01-20T08:15:00",
  "accidentLocation": "Intersection of 5th Avenue and Main Street, Anytown, State 12345",
  "locationAddress": {
    "@type": "Address",
    "streetAddress": "5th Avenue & Main Street",
    "city": "Anytown", 
    "state": "State",
    "postalCode": "12345"
  },
  "damageDescription": "Scratches and dent to front passenger side door and quarter panel. Paint transfer visible.",
  "estimatedDamageAmount": 1200.00,
  "actualRepairCost": 0.00,
  "policeReportFiled": false,
  "policeReportNumber": null,
  "citationIssued": false,
  "faultPercentage": 100,
  "numberOfVehiclesInvolved": 1,
  "numberOfInjuries": 0,
  "injurySeverity": "None",
  "speedAtImpact": 15,
  "claimFiled": false,
  "reasonNoClaimFiled": "Damage below deductible threshold. Chose to pay out-of-pocket to avoid potential premium increase and claims history impact.",
  "description": "Side-swiped a metal guard rail while swerving to avoid road debris (fallen tree branch) during morning commute. Single vehicle incident with no other parties involved. Damage cosmetic only - vehicle remained drivable.",
  "involvedVehicle": {
    "@type": "Vehicle",
    "vehicleVIN": "1HGBH41JXMN109186", 
    "vehicleMake": "Honda",
    "vehicleModel": "Accord",
    "vehicleYear": "2021"
  },
  "involvedDriver": {
    "@type": "Person",
    "name": "John Smith",
    "licenseNumber": "DL123456789",
    "licenseState": "State"
  }
}
```

### Form 3: Multi-Vehicle Collision with Claim

```json
{
  "@context": {
    "@vocab": "http://example.org/ontology/insurance_claims#",
    "xsd": "http://www.w3.org/2001/XMLSchema#"
  },
  "@type": "Claim",
  "claimNumber": "CLM-2023-789012",
  "claimDate": "2023-08-05",
  "incidentDate": "2023-08-03",
  "settlementDate": "2023-12-15", 
  "claimType": {
    "@type": "skos:Concept",
    "@id": "LiabilityClaim",
    "prefLabel": "Liability"
  },
  "claimAmount": 15750.00,
  "paidAmount": 12600.00,
  "deductibleAmount": 1000.00,
  "faultStatus": {
    "@type": "skos:Concept",
    "@id": "PartialFault", 
    "prefLabel": "Partial fault"
  },
  "claimStatus": {
    "@type": "skos:Concept",
    "@id": "ClosedPaid",
    "prefLabel": "Closed-Paid"
  },
  "description": "Three-vehicle chain reaction collision during rush hour traffic. Failed to stop in time when traffic suddenly slowed, rear-ended vehicle ahead, which was pushed into vehicle in front of it.",
  "relatedAccident": {
    "@type": "Accident",
    "accidentDate": "2023-08-03T17:30:00",
    "accidentLocation": "Interstate 95 Southbound, Mile Marker 127, near Exit 45",
    "policeReportFiled": true,
    "policeReportNumber": "TR-2023-080312",
    "citationIssued": true,
    "faultPercentage": 70,
    "numberOfVehiclesInvolved": 3,
    "numberOfInjuries": 2,
    "injurySeverity": "Minor",
    "speedAtImpact": 25
  },
  "paidByInsurer": {
    "@type": "InsuranceCompany", 
    "companyName": "State Farm Insurance",
    "naic_code": "25178"
  }
}
```

### Form 4: Theft Claim (Comprehensive Coverage)

```json
{
  "@context": {
    "@vocab": "http://example.org/ontology/insurance_claims#",
    "xsd": "http://www.w3.org/2001/XMLSchema#"
  },
  "@type": "Claim",
  "claimNumber": "CLM-2023-345678",
  "claimDate": "2023-09-22",
  "incidentDate": "2023-09-20",
  "settlementDate": "2023-11-30",
  "claimType": {
    "@type": "skos:Concept",
    "@id": "TheftClaim", 
    "prefLabel": "Theft"
  },
  "claimAmount": 28500.00,
  "paidAmount": 27500.00,
  "deductibleAmount": 1000.00,
  "faultStatus": {
    "@type": "skos:Concept",
    "@id": "NotAtFault",
    "prefLabel": "Not-at-fault"
  },
  "claimStatus": {
    "@type": "skos:Concept",
    "@id": "ClosedPaid",
    "prefLabel": "Closed-Paid"
  },
  "description": "Vehicle stolen from apartment complex parking garage overnight. Vehicle was locked and parked in assigned space. Theft discovered when leaving for work the following morning.",
  "relatedAccident": {
    "@type": "Accident",
    "accidentDate": "2023-09-20T23:00:00",
    "accidentLocation": "Riverside Apartments, 1234 Oak Street, Parking Level B, Space 47",
    "policeReportFiled": true,
    "policeReportNumber": "THEFT-2023-092001",
    "numberOfVehiclesInvolved": 1,
    "numberOfInjuries": 0,
    "claimFiled": true
  },
  "involvedVehicle": {
    "@type": "Vehicle",
    "vehicleVIN": "JM1BK32F781234567",
    "vehicleMake": "Mazda", 
    "vehicleModel": "CX-5",
    "vehicleYear": "2022"
  }
}
```

### Form 5: Windshield Claim (Glass Coverage)

```json
{
  "@context": {
    "@vocab": "http://example.org/ontology/insurance_claims#",
    "xsd": "http://www.w3.org/2001/XMLSchema#"
  },
  "@type": "Claim",
  "claimNumber": "CLM-2024-123456",
  "claimDate": "2024-03-10",
  "incidentDate": "2024-03-08",
  "settlementDate": "2024-03-15",
  "claimType": {
    "@type": "skos:Concept",
    "@id": "WindshieldClaim",
    "prefLabel": "Windshield/Glass"
  },
  "claimAmount": 450.00,
  "paidAmount": 450.00,
  "deductibleAmount": 0.00,
  "faultStatus": {
    "@type": "skos:Concept",
    "@id": "NotAtFault", 
    "prefLabel": "Not-at-fault"
  },
  "claimStatus": {
    "@type": "skos:Concept",
    "@id": "ClosedPaid",
    "prefLabel": "Closed-Paid"
  },
  "description": "Windshield cracked by rock thrown up by truck on highway. Crack started small but spread across driver's field of vision, requiring full windshield replacement.",
  "relatedAccident": {
    "@type": "Accident",
    "accidentDate": "2024-03-08T14:20:00",
    "accidentLocation": "Highway 101 Northbound, approximately 2 miles south of Exit 23",
    "policeReportFiled": false,
    "numberOfVehiclesInvolved": 1,
    "numberOfInjuries": 0,
    "claimFiled": true,
    "reasonNoClaimFiled": null
  }
}
```

## Key Extraction Points for Underwriting Systems

### Critical Risk Indicators
1. **Fault Percentage** - Direct impact on risk scoring
2. **Claim Frequency** - Number of claims over time period  
3. **Claim Severity** - Average claim amounts
4. **Accident-to-Claim Ratio** - Incidents that didn't result in claims
5. **Police Report Correlation** - Official documentation patterns
6. **Injury Involvement** - Bodily injury claim history
7. **Multi-Vehicle Incidents** - Complexity and liability patterns

### Underwriting Decision Factors
- **At-Fault Claims**: Highest risk impact
- **Claim Amount Trends**: Escalating damage costs
- **Unreported Accidents**: Risk of future claims
- **Citation History**: Traffic violation correlation
- **Damage Patterns**: Consistent risk behaviors
- **Settlement Speed**: Claims handling complexity

### Data Quality Indicators
- **Complete Police Reports**: Higher data reliability
- **Consistent Fault Determinations**: Clear liability patterns  
- **Detailed Damage Descriptions**: Accurate loss assessment
- **Temporal Accuracy**: Precise incident timing
- **Geographic Patterns**: Location-based risk factors

This extraction framework enables insurance companies to perform comprehensive risk assessment based on both claimed and unclaimed incidents, providing a complete picture of an applicant's driving risk profile.
