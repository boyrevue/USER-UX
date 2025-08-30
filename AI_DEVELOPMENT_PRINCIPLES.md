# 🧠 AI Development Principles for CLIENT-UX

## Overarching Principle
All code must be designed with the understanding that it is part of an AI system integrated with semantic web technologies. This means prioritizing clarity for both humans and AI, modularity for easy experimentation, performance for scalable inference, and semantic consistency through ontologies with GDPR compliance at the field level.

---

## 1. 🌐 Ontology & Semantics Rules

### Rule 1.1: Ontology-Driven Development

**What**: All prompts, system strings, validation rules, and data models must be defined in and referenced from ontology files (.ttl format). Use standardized ontologies (Schema.org, Dublin Core, custom) rather than arbitrary strings.

**Why**: Ensures semantic consistency, enables reasoning, facilitates knowledge graph integration, and provides a single source of truth for all system knowledge.

**CLIENT-UX Implementation**: Our modular TTL architecture already implements this:
- `AI_Driver_Details.ttl` - Driver domain knowledge
- `AI_Vehicle_Details.ttl` - Vehicle specifications
- `AI_Claims_History.ttl` - Claims and validation rules
- `AI_Data_Compliance.ttl` - GDPR compliance definitions

**Cursor Prompt**: "Create a function that loads our prompt ontology from prompts.ttl and extracts the system prompt with ID :dataProcessingPrompt. Use the rdflib library to query the TTL file."

### Rule 1.2: Internationalization (i18n) Compliance

**What**: All user-facing strings (in prompts, UI, validation messages) must be stored in the ontology with language tags (e.g., "Hello"@en, "Hola"@es, "Bonjour"@fr). Never hardcode language-specific strings.

**Why**: Supports multilingual applications, separates content from logic, and enables dynamic language switching.

**CLIENT-UX Implementation**: 
```turtle
autoins:ValidationMessages a owl:Class ;
    rdfs:label "Validation Messages" .

autoins:ClaimDateTooOldMessage a autoins:ValidationMessage ;
    rdfs:label "Claim date too old"@en ;
    rdfs:label "Date de réclamation trop ancienne"@fr ;
    rdfs:label "Fecha de reclamo demasiado antigua"@es ;
    autoins:messageText "⚠️ CLAIM DATE ERROR: Claim date must be within the last 5 years"@en ;
    autoins:messageText "⚠️ ERREUR DE DATE: La date de réclamation doit être dans les 5 dernières années"@fr ;
    autoins:i18nKey "validation.claim_date.five_year_limit" .
```

**Cursor Prompt**: "Modify the response generation function to select the appropriate validation message based on the user's language preference from session data. Reference the multi-lingual messages in our i18n_temporal.ttl file."

### Rule 1.3: SHACL Validation Mandatory

**What**: Create SHACL shapes to validate all ontology data, prompts, configurations, and user input before processing. Implement validation checks at system initialization and data ingestion points.

**Why**: Ensures data quality, prevents runtime errors from malformed data, maintains ontology consistency, and enforces business rules.

**CLIENT-UX Implementation**: Already implemented in `AI_Claims_History.ttl`:
```turtle
autoins:ClaimsHistoryValidationShape a sh:NodeShape ;
    sh:targetClass autoins:Claim ;
    rdfs:label "Claims History 5-Year Validation" ;
    sh:property [
        sh:path autoins:claimDate ;
        sh:datatype xsd:date ;
        sh:minCount 1 ;
        sh:maxCount 1 ;
        sh:message "⚠️ CLAIM DATE ERROR: Claim date must be within the last 5 years" ;
        autoins:validationSeverity "ERROR" ;
        autoins:validationCode "CLAIM_DATE_TOO_OLD" ;
        autoins:validationRule "claimDate >= (TODAY - P5Y) AND claimDate <= TODAY"
    ] .
```

**Cursor Prompt**: "Write a function that validates new user input against the SHACL rules defined in AI_Claims_History.ttl before processing. Return meaningful validation errors with i18n support."

---

## 2. 🔒 GDPR Compliance Rules

### Rule 2.1: Field-Level GDPR Classification

**What**: Every data field in the ontology must be tagged with its GDPR classification (personal data, sensitive data, anonymous data, etc.) using properties from our `AI_Data_Compliance.ttl` ontology.

**Why**: Enables automated compliance checking, data handling policies, and audit trails for regulatory compliance.

**CLIENT-UX Implementation**: Already implemented with comprehensive field-level annotations:
```turtle
autoins:dateOfBirth a owl:DatatypeProperty ;
    rdfs:domain autoins:Driver ;
    rdfs:range xsd:date ;
    rdfs:label "date of birth" ;
    # GDPR Compliance - HIGHLY SENSITIVE
    autoins:dataClassification autoins:SensitivePersonalData ;
    autoins:personalDataCategory autoins:IdentityData ;
    autoins:accessLevel autoins:StaffAccess ;
    autoins:viewPermission "STAFF,MANAGER,ADMIN" ;
    autoins:editPermission "STAFF,MANAGER,ADMIN" ;
    autoins:deletePermission "ADMIN" ;
    autoins:obfuscationMethod autoins:PartialMasking ;
    autoins:maskingPattern "**/**/19**" ;
    autoins:consentRequired "true"^^xsd:boolean ;
    autoins:consentPurpose "Age verification and insurance underwriting" ;
    autoins:consentBasis autoins:Contract ;
    autoins:retentionPeriod "P7Y" ;
    autoins:retentionReason "Insurance regulatory requirements" ;
    autoins:auditRequired "true"^^xsd:boolean ;
    autoins:logLevel autoins:CriticalLogging .
```

**Cursor Prompt**: "Create a function that checks if a requested data processing operation is allowed based on the GDPR classification of each data field involved. Reference our AI_Data_Compliance.ttl file for access permissions."

### Rule 2.2: Data Processing Compliance Checks

**What**: Before any AI processing, validate that the operation complies with the GDPR classification of the data. Implement purpose limitation and data minimization by design.

**Why**: Prevents unintended GDPR violations and ensures privacy-by-design architecture.

**CLIENT-UX Implementation**: Compliance framework with access levels and consent management:
```turtle
autoins:DataProcessingOperation a owl:Class ;
    rdfs:label "Data Processing Operation" .

autoins:InsuranceUnderwriting a autoins:DataProcessingOperation ;
    autoins:allowedDataCategories autoins:IdentityData, autoins:FinancialData ;
    autoins:requiredConsentBasis autoins:Contract ;
    autoins:minimumAccessLevel autoins:StaffAccess ;
    autoins:auditRequired "true"^^xsd:boolean .
```

**Cursor Prompt**: "Implement a decorator function that checks GDPR compliance before executing any data processing function. The decorator should validate against our SHACL shapes for GDPR in AI_Data_Compliance.ttl."

### Rule 2.3: Right to Erasure Implementation

**What**: Ensure all data structures support erasure requests by implementing soft deletion patterns and data lineage tracking in your ontology.

**Why**: Required for GDPR Article 17 compliance and user privacy rights.

**CLIENT-UX Implementation**: Retention and erasure policies defined per field:
```turtle
autoins:DataRetentionPolicy a owl:Class ;
    rdfs:label "Data Retention Policy" .

autoins:InsuranceDataRetention a autoins:DataRetentionPolicy ;
    autoins:retentionPeriod "P7Y" ;
    autoins:retentionReason "Insurance regulatory requirements" ;
    autoins:erasureMethod autoins:SecureAnonymization ;
    autoins:auditTrail "true"^^xsd:boolean .
```

**Cursor Prompt**: "Create a function that processes GDPR erasure requests by identifying all data related to a user across our systems using the data lineage in our ontology, and anonymizes it according to our erasure protocols."

---

## 3. 🏗️ Code Structure & Architecture Rules

### Rule 3.1: Modular by Functionality with Ontology Backing

**What**: Break down the system into discrete, reusable modules where each module corresponds to an ontology domain (e.g., driver_service.py ↔ AI_Driver_Details.ttl).

**Why**: Allows for independent development, testing, and maintains semantic consistency between code and knowledge representation.

**CLIENT-UX Implementation**: Our modular architecture already follows this pattern:
- `DriverForm.tsx` ↔ `AI_Driver_Details.ttl`
- `VehicleForm.tsx` ↔ `AI_Vehicle_Details.ttl`
- `ClaimsForm.tsx` ↔ `AI_Claims_History.ttl`
- `services/validation.ts` ↔ SHACL shapes in TTL files

**Cursor Prompt**: "Create a function to load and preprocess driver data. Keep it in a separate module named driver_service.py. Ensure it validates against the SHACL shapes in AI_Driver_Details.ttl and returns semantically consistent data."

### Rule 3.2: Config-Driven Design with Ontology Configuration

**What**: Never hardcode parameters. Use ontology-based configuration where system settings are defined in TTL files with proper semantic relationships.

**Why**: Enables rapid experimentation, reproducibility, and semantic reasoning over configuration.

**CLIENT-UX Implementation**: Configuration in `user_ux.ttl`:
```turtle
autoins:ApplicationConfiguration a owl:Class ;
    rdfs:label "Application Configuration" .

autoins:ValidationSettings a autoins:ApplicationConfiguration ;
    autoins:ukMinDrivingAge "17"^^xsd:int ;
    autoins:maxHumanAge "130"^^xsd:int ;
    autoins:historicalDataLimit "P5Y"^^xsd:duration ;
    autoins:dvlaRecordsStart "1970-01-01"^^xsd:date .
```

**Cursor Prompt**: "Read the validation settings from user_ux.ttl instead of hardcoding them. Create a function to load the ontology configuration and extract validation parameters using SPARQL queries."

### Rule 3.3: Define Clear Interfaces with Semantic Contracts

**What**: Each function and class must have a single, clear purpose defined by its corresponding ontology concept. Use type hints and ontology-backed docstrings.

**Why**: Makes code self-documenting and ensures semantic consistency between implementation and domain model.

**Cursor Prompt**: "Write a function validate_driver_age(birth_date: str, licence_date: str) -> ValidationResult. Add a docstring that references the corresponding SHACL shape in AI_Driver_Details.ttl and explains the UK driving age validation rules."

---

## 4. 🤖 AI-Specific Development Rules with Semantics

### Rule 4.1: Semantic Prompt Management

**What**: Store all prompts in ontology files with semantic relationships, context constraints, and compliance annotations. Generate prompts dynamically by querying the ontology.

**Why**: Allows for semantic reasoning over prompts, automated compliance checking, and context-aware prompt selection.

**Implementation Example**:
```turtle
@prefix : <https://client-ux.example/prompts#> .

:DataExtractionPrompt a :AIPrompt ;
    rdfs:label "Document data extraction prompt"@en ;
    :promptText """
    You are a helpful assistant that extracts data from documents.
    IMPORTANT: You must comply with GDPR Article 6 when processing personal data.
    Only extract data that is explicitly visible in the document.
    
    Document type: {document_type}
    User query: {user_query}
    
    Return structured JSON with confidence scores.
    """@en ;
    :hasComplianceLevel :GDPR_Strict ;
    :allowedDataCategories :IdentityData, :DocumentData ;
    :requiredAccessLevel :UserAccess ;
    :contextWindowLimit "4096"^^xsd:int .
```

**Cursor Prompt**: "Create a function that constructs AI prompts by querying our prompt ontology for templates matching the current operation type and compliance level, then filling in the parameter slots with GDPR-compliant data."

### Rule 4.2: Ontology-Backed Response Validation

**What**: Validate AI responses against SHACL shapes defined in your ontology to ensure they comply with expected formats, data policies, and content restrictions.

**Why**: Maintains data quality, compliance, and semantic consistency in AI outputs.

**Cursor Prompt**: "After receiving a response from the AI model, validate it against the SHACL shape for :AIResponse in our response_shapes.ttl file before returning it to the user. Include GDPR compliance checks."

### Rule 4.3: Context Window Management with Semantic Prioritization

**What**: Use ontology-defined importance rankings and GDPR classifications to intelligently truncate context when approaching model limits.

**Why**: Ensures critical information is preserved while maintaining compliance and performance.

**Cursor Prompt**: "Before calling the model, check if the tokenized input exceeds the context limit. If it does, implement semantic truncation that preserves high-priority fields as defined in our ontology and respects GDPR data minimization principles."

---

## 5. 📊 Implementation in CLIENT-UX Architecture

### Current Ontology Structure
```
ontology/
├── AI_Driver_Details.ttl           # Driver domain with GDPR annotations
├── AI_Vehicle_Details.ttl          # Vehicle specifications
├── AI_Policy_Details.ttl           # Insurance policy rules
├── AI_Claims_History.ttl           # Claims with SHACL validation
├── AI_Insurance_Payments.ttl       # Payment processing
├── AI_Data_Compliance.ttl          # GDPR compliance framework
├── user_ux.ttl                    # Application configuration
├── user_documents.ttl             # Document processing rules
├── user_credentials_pci.ttl        # PCI compliance
├── i18n_temporal.ttl              # Internationalization
└── personal_documents_ontology.ttl # Document taxonomy
```

### Service Integration Pattern
```python
class OntologyBackedService:
    def __init__(self):
        self.ontology = Graph()
        self.load_ontologies()
        self.validator = SHACLValidator(self.ontology)
    
    def load_ontologies(self):
        """Load all TTL files into unified graph"""
        for ttl_file in self.get_ontology_files():
            self.ontology.parse(ttl_file, format="turtle")
    
    def validate_input(self, data: dict, shape_uri: str) -> ValidationResult:
        """Validate input against SHACL shapes"""
        return self.validator.validate(data, shape_uri)
    
    def get_gdpr_classification(self, field_name: str) -> GDPRClassification:
        """Get GDPR classification from ontology"""
        query = f"""
        PREFIX autoins: <https://autoins.example/ontology#>
        SELECT ?classification ?accessLevel ?retentionPeriod
        WHERE {{
            autoins:{field_name} autoins:dataClassification ?classification ;
                                autoins:accessLevel ?accessLevel ;
                                autoins:retentionPeriod ?retentionPeriod .
        }}
        """
        return self.ontology.query(query)
```

### Frontend Integration
```typescript
// services/ontology.ts
export class OntologyService {
  async getValidationRules(fieldName: string): Promise<ValidationRule[]> {
    const response = await ApiService.getOntology();
    return this.extractSHACLRules(response, fieldName);
  }
  
  async getI18nLabel(concept: string, language: string): Promise<string> {
    // Query ontology for language-specific labels
    return this.queryOntologyLabel(concept, language);
  }
  
  async checkGDPRCompliance(operation: string, dataFields: string[]): Promise<ComplianceResult> {
    // Validate operation against GDPR rules in ontology
    return this.validateGDPROperation(operation, dataFields);
  }
}
```

---

## 6. 🎯 Benefits of This Approach

### **Semantic Consistency**
- Single source of truth for all domain knowledge
- Automated consistency checking across code and data
- Reasoning capabilities over business rules

### **GDPR Compliance by Design**
- Field-level privacy annotations
- Automated compliance checking
- Audit trails and data lineage tracking
- Right to erasure implementation

### **AI Enhancement**
- Context-aware prompt generation
- Semantic validation of AI outputs
- Intelligent data prioritization
- Compliance-aware processing

### **Internationalization**
- Language-agnostic core logic
- Ontology-driven i18n content
- Cultural adaptation through semantic annotations

### **Maintainability**
- Self-documenting through semantic annotations
- Automated validation and testing
- Clear separation of concerns
- Future-proof architecture

---

## 7. 🚀 Next Steps for CLIENT-UX

1. **Enhance Prompt Ontology**: Create dedicated TTL files for AI prompts with compliance annotations
2. **Implement SHACL Validation**: Add runtime validation using pyshacl or similar libraries
3. **Expand GDPR Annotations**: Complete field-level GDPR classification across all ontologies
4. **Add Semantic Reasoning**: Implement inference rules for automated compliance checking
5. **Create Validation Pipeline**: Automated testing of ontology consistency and SHACL compliance

**🎯 This approach transforms CLIENT-UX into a semantically-aware, GDPR-compliant, AI-enhanced system that maintains knowledge consistency while enabling advanced reasoning and automated compliance checking.**
