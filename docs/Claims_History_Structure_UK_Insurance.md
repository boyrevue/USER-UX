# Claims History Structure for UK Insurance

## The Question: Main Driver vs Actual Driver Claims

**User Question**: "Regarding claims history, what happens when the main driver was not the driver generating the claim? Should we keep claims history in same place as main driver history? What's the insurance firm really after here - the dirt on the offender or the non dirt on the main driver?"

## The Answer: Both!

### UK Insurance Reality

**Claims Follow the Policy, Not the Driver**
- In the UK, insurance follows the **vehicle/policy**, not the individual driver
- When a named driver causes an accident, it affects the **main policyholder's record**
- The main driver's No Claims Bonus (NCB) is impacted regardless of who was driving
- Premium increases apply to the policy holder, not just the actual driver

### What Insurance Firms Actually Need

1. **Policy-Level Claims History**
   - All claims made under this policy (regardless of who was driving)
   - Affects the main driver's NCB and future premiums
   - Required for underwriting the policy renewal

2. **Individual Driver Attribution**
   - Who was actually driving when each incident occurred
   - Individual risk assessment for each driver
   - Required for accurate risk pricing

3. **Cross-Policy Claims History**
   - Each driver's claims history from OTHER policies
   - Personal driving record regardless of whose policy it was under

## Updated Ontology Structure

### Main Claims Question
```turtle
autoins:hasAccidents a owl:DatatypeProperty ;
  rdfs:label "Any Accidents/Claims In Last 5 Years?" ;
  autoins:formHelpText "Have you or any named driver had accidents/claims while driving ANY vehicle in the last 5 years? This includes claims made under other policies." ;
```

### Driver Attribution Field
```turtle
autoins:whoWasDriving a owl:DatatypeProperty ;
  rdfs:label "Who Was Driving?" ;
  autoins:enumerationValues ("Main Driver" "Named Driver" "Other Permitted Driver" "Unknown") ;
  autoins:formHelpText "Who was driving the vehicle at the time of the accident/claim?" ;
```

### SHACL Validation Added
```turtle
sh:property [
  sh:path autoins:yearsHeldFull ;
  sh:datatype xsd:integer ;
  sh:minInclusive 0 ;
  sh:maxInclusive 80 ;
  sh:message "Years held full licence must be a positive number from 0 to 80" ;
] ;
```

## Insurance Firm Priorities

1. **Risk Assessment**: They want to know the total risk exposure of the policy
2. **Claims Attribution**: They need to understand who caused what for individual risk profiling
3. **Policy Pricing**: Both policy-level and individual-level risk factors affect pricing
4. **Regulatory Compliance**: FCA requires accurate disclosure of all material facts

## Form Structure Decision

**Keep claims history with the main driver section** because:
- Claims legally belong to the policy holder
- Main driver's NCB is affected regardless of who was driving
- Underwriting decisions are made at the policy level
- But capture driver attribution for individual risk assessment

## Key Insight

Insurance firms want **both the dirt on the offender AND the impact on the main driver** because:
- The claim affects the main driver's policy and NCB
- The individual driver's risk profile affects future pricing
- Cross-referencing helps detect fraud and assess overall risk
