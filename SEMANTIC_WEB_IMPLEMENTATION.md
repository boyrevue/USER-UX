# Semantic Web Implementation - CLIENT-UX

## Overview

This document describes the implementation of semantic web standards, SHACL validation, and grounded AI patterns in CLIENT-UX, following the ontology-driven development principles.

## 🎯 Implemented Features

### A) FNOL (First Notice of Loss) Completeness SHACL Validation

**Location**: `ontology/AI_Claims_History.ttl` (lines 1190-1252)

**Implementation**:
```turtle
autoins:FNOLCompletenessShape a sh:NodeShape ;
  sh:targetClass autoins:Claim ;
  rdfs:label "FNOL Completeness Validation" ;
  sh:property [
    sh:path autoins:hasIncident ;
    sh:minCount 1 ;
    sh:class autoins:Accident ;
    sh:message "🚨 FNOL ERROR: Claim must link to an incident/accident record for complete First Notice of Loss processing." ;
  ] ;
  sh:property [
    sh:path [ sh:inversePath autoins:relatesToPolicy ] ;
    sh:minCount 1 ;
    sh:class autoins:Policy ;
    sh:message "📋 FNOL ERROR: Claim must relate to a valid insurance policy for processing." ;
  ] .
```

**Validation Requirements**:
- ✅ Claim must link to incident (`autoins:hasIncident`)
- ✅ Claim must relate to policy (`autoins:relatesToPolicy`)
- ✅ Claim date required (`autoins:claimDate`)
- ✅ Incident date required (`autoins:incidentDate`)
- ✅ Claim type required (`autoins:claimType`)
- ✅ Detailed description required (min 10 characters)

### B) Reserve Calculator (Verifier) Service

**Location**: `internal/services/reserve/calculator.go`

**Implementation**:
```go
func (rc *ReserveCalculator) CalculateReserve(claim ClaimData) ReserveResult {
    // Step 1: Get base reserve from severity table
    base := rc.getSeverityTableAmount(claim.LossType, claim.VehicleACV)
    
    // Step 2: Apply modifiers
    var mods float64 = 0
    if claim.HasFraudSignals {
        mods += 0.1 * base  // +10% for fraud signals
    }
    if claim.LiabilityUncertain {
        mods += 0.15 * base // +15% for liability uncertainty
    }
    if claim.PartsBackorder {
        mods += 0.05 * base // +5% for parts availability
    }
    
    // Step 3: Calculate final reserve and quantize to band
    finalReserve := base + mods
    reserveBand := rc.quantizeToBand(finalReserve)
    
    return ReserveResult{...}
}
```

**Severity Tables**:
- ✅ Collision: £1,500 - £12,000 (by vehicle ACV)
- ✅ Theft: £3,000 - £25,000 (by vehicle ACV)
- ✅ Vandalism: £800 - £5,000 (by vehicle ACV)
- ✅ Fire: £2,500 - £22,000 (by vehicle ACV)
- ✅ Third Party: £5,000 - £35,000 (by vehicle ACV)

**Reserve Bands**:
- £0-2k, £2k-5k, £5k-10k, £10k-25k, £25k-50k, £50k+

### C) Grounded Prompt Pattern with SPARQL Tools

**Location**: `internal/services/grounded/prompt_engine.go`

**System Prompt** (from ontology):
```
You are a grounded AI assistant for insurance claims processing. 
You must answer only from the knowledge graph and attached calculations.

MANDATORY TOOLS AVAILABLE:
- SPARQL_SELECT: Query the insurance ontology graph
- COVERAGE_CALC: Calculate policy coverage amounts
- RESERVE_CALC: Calculate claim reserves using severity tables
- FRAUD_SCORER: Assess fraud risk indicators

GROUNDING REQUIREMENTS:
1. You MUST use SPARQL_SELECT to retrieve facts from the graph before answering
2. You MUST use appropriate calculation tools for any numerical assessments
3. If a fact is missing from the graph, ask ONE targeted follow-up question, then re-check
4. Your final answer MUST list the IRIs of all facts used in your reasoning
5. You MUST cite specific calculation results with their input parameters
```

**Response Format**:
```json
{
  "answer": "Based on the facts retrieved from the knowledge graph...",
  "factsUsed": [
    {
      "iri": "autoins:claim_001#claimType",
      "value": "Collision",
      "property": "autoins:claimType",
      "source": "SPARQL_SELECT",
      "confidence": 1.0
    }
  ],
  "calculationsPerformed": [
    {
      "type": "RESERVE_CALC",
      "input": {"lossType": "Collision", "vehicleACV": 15000},
      "result": {"finalReserve": 4600, "reserveBand": "£2k-5k"},
      "confidence": 0.95
    }
  ],
  "followUpNeeded": false,
  "confidenceLevel": 0.92
}
```

## 🔧 API Endpoints

### Grounded AI Processing
- `POST /api/grounded/query` - Process grounded AI queries
- `POST /api/grounded/reserve` - Calculate claim reserves
- `POST /api/grounded/fraud` - Assess fraud risk
- `POST /api/grounded/fnol` - Validate FNOL completeness
- `GET /api/grounded/prompt` - Get system prompt template

### Example Usage

#### Reserve Calculation
```bash
curl -X POST http://localhost:3000/api/grounded/reserve \
  -H "Content-Type: application/json" \
  -d '{
    "lossType": "Collision",
    "vehicleACV": 15000,
    "hasFraudSignals": false,
    "liabilityUncertain": true,
    "partsBackorder": false
  }'
```

Response:
```json
{
  "baseReserve": 4000,
  "fraudModifier": 0,
  "liabilityModifier": 600,
  "partsModifier": 0,
  "finalReserve": 4600,
  "reserveBand": "£2k-5k",
  "breakdown": "Base Reserve: £4000.00 + Liability Uncertainty: £600.00 = Final Reserve: £4600.00"
}
```

#### FNOL Validation
```bash
curl -X POST http://localhost:3000/api/grounded/fnol \
  -H "Content-Type: application/json" \
  -d '{
    "claimDate": "2025-01-24",
    "incidentDate": "2025-01-20",
    "claimType": "Collision",
    "claimDescription": "Vehicle collision at intersection"
  }'
```

Response:
```json
{
  "isComplete": false,
  "missingFields": ["hasIncident", "relatesToPolicy"],
  "validationErrors": [
    "🚨 FNOL ERROR: Incident/accident record is required for First Notice of Loss processing",
    "📋 FNOL ERROR: Policy reference is required for processing"
  ],
  "completeness": 66.67,
  "recommendations": [
    "Link this claim to the underlying incident/accident record",
    "Verify and link the claim to the appropriate insurance policy"
  ]
}
```

#### Grounded Query Processing
```bash
curl -X POST http://localhost:3000/api/grounded/query \
  -H "Content-Type: application/json" \
  -d '{
    "userQuery": "What is the reserve estimate for this collision claim?",
    "graphContext": {
      "claim": {
        "claimType": "Collision",
        "claimAmount": 5000
      },
      "vehicle": {
        "actualCashValue": 15000
      },
      "fraudIndicators": {
        "hasFraudSignals": false,
        "liabilityUncertain": true,
        "partsBackorder": false
      }
    },
    "requiredTools": ["SPARQL_SELECT", "RESERVE_CALC"]
  }'
```

## 🏗️ Architecture Integration

### Ontology-Driven Development
- All prompts stored in TTL format with i18n support
- SHACL shapes validate data integrity
- GDPR compliance annotations at field level
- Semantic relationships drive business logic

### Service Layer Architecture
```
┌─────────────────────────────────────────┐
│           API Handlers                   │
│  (grounded_ai.go)                       │
├─────────────────────────────────────────┤
│           Service Layer                  │
│  ┌─────────────┐  ┌─────────────────┐   │
│  │   Reserve   │  │    Grounded     │   │
│  │ Calculator  │  │ Prompt Engine   │   │
│  └─────────────┘  └─────────────────┘   │
├─────────────────────────────────────────┤
│           Ontology Layer                 │
│  ┌─────────────┐  ┌─────────────────┐   │
│  │ SHACL Rules │  │  TTL Ontology   │   │
│  │ Validation  │  │   Knowledge     │   │
│  └─────────────┘  └─────────────────┘   │
└─────────────────────────────────────────┘
```

### Data Flow
1. **Request** → API Handler validates against SHACL shapes
2. **SPARQL Query** → Extract facts from ontology graph
3. **Calculation** → Apply business rules via services
4. **Response** → Grounded answer with fact citations

## 🔍 Validation & Compliance

### SHACL Validation Rules
- **FNOL Completeness**: Ensures all required claim elements
- **5-Year Historical Limits**: Claims/accidents within underwriting period
- **Age Validation**: UK driving age limits (17-130 years)
- **Data Type Validation**: Proper ranges and formats

### GDPR Compliance
- Field-level data classification
- Access control annotations
- Retention period specifications
- Audit trail requirements
- Right to erasure support

### Example GDPR Annotation
```turtle
autoins:claimAmount a owl:DatatypeProperty ;
  # ... property definition ...
  # GDPR Compliance - SENSITIVE FINANCIAL DATA
  autoins:dataClassification autoins:SensitivePersonalData ;
  autoins:personalDataCategory autoins:FinancialData ;
  autoins:accessLevel autoins:StaffAccess ;
  autoins:obfuscationMethod autoins:PartialMasking ;
  autoins:retentionPeriod "P10Y" ;
  autoins:auditRequired "true"^^xsd:boolean .
```

## 🚀 Benefits Achieved

### 1. **Semantic Consistency**
- Single source of truth in ontologies
- Standardized vocabulary across system
- Automated validation and error detection

### 2. **AI Transparency**
- All AI responses grounded in verifiable facts
- Complete audit trail of reasoning process
- Fact citations with confidence scores

### 3. **Regulatory Compliance**
- GDPR compliance by design
- Automated data retention management
- Field-level access control

### 4. **Business Agility**
- Rule changes via ontology updates
- No code changes for business logic
- Multilingual support through i18n annotations

### 5. **Quality Assurance**
- SHACL validation prevents data errors
- Automated reserve calculations
- Fraud detection with explainable scoring

## 📈 Next Steps

1. **Enhanced SPARQL Engine**: Implement full RDF store integration
2. **Advanced Fraud Detection**: ML-based risk scoring
3. **Real-time Validation**: Live SHACL validation in UI
4. **Audit Dashboard**: GDPR compliance monitoring
5. **API Documentation**: OpenAPI/Swagger integration

## 🔧 Development Guidelines

### Adding New SHACL Rules
1. Define shape in appropriate TTL file
2. Add validation messages with i18n keys
3. Implement backend validation logic
4. Update API documentation

### Extending Reserve Calculator
1. Update severity tables in `calculator.go`
2. Add new modifier logic
3. Update ontology with new properties
4. Test with various scenarios

### Creating Grounded Prompts
1. Define prompt in TTL with language tags
2. Specify required tools and compliance level
3. Implement validation shapes
4. Test with fact verification

This implementation provides a solid foundation for semantic web standards while maintaining practical usability and regulatory compliance.
