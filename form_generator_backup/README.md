# TTL to HTML Form Generator

This Go application reads the TTL ontology file and generates dynamic HTML forms with proper validation, field types, and categorization.

## Features

- **TTL Parser**: Reads and parses Turtle (TTL) ontology files
- **Dynamic Form Generation**: Creates HTML forms based on ontology field definitions
- **Field Categorization**: Organizes fields into logical categories (Drivers, Vehicle, Coverage, etc.)
- **Validation Support**: Includes HTML5 validation patterns from the ontology
- **Responsive Design**: Mobile-friendly form layout
- **Interactive Tabs**: Category-based navigation

## Usage

### Quick Start

```bash
# Run the form generator
./run.sh

# Or manually
go run main.go
```

### Command Line Options

```bash
go run main.go -ttl ../ontology/autoins.ttl -output my_forms.html
```

- `-ttl`: Path to the TTL ontology file (default: `../ontology/autoins.ttl`)
- `-output`: Path for the generated HTML file (default: `generated_forms.html`)

## Generated Categories

The form generator creates 7 main categories:

1. **👥 Driver Details** - Personal information, license details
2. **🚗 Vehicle Information** - Make, model, VIN, specifications
3. **🛡️ Coverage Options** - Insurance coverage and deductibles
4. **📋 Claims History** - Accident and violation history
5. **💳 Payment Information** - Payment methods and billing
6. **⚙️ Preferences** - Communication and contact preferences
7. **📊 Summary** - Policy summary and final details

## Field Types Supported

- **Text**: Standard text input with validation patterns
- **Number**: Numeric input with min/max constraints
- **Date**: Date picker for birth dates, policy dates
- **Email**: Email validation
- **Tel**: Phone number with pattern validation
- **Select**: Dropdown with predefined options
- **Textarea**: Multi-line text for descriptions
- **Checkbox**: Boolean values

## TTL Ontology Structure

The parser expects the TTL file to contain:

```turtle
# Field definitions
autoins:firstName a owl:DatatypeProperty ;
    rdfs:label "First Name" ;
    rdfs:range xsd:string ;
    autoins:isRequired true ;
    autoins:validationPattern "^[A-Za-z\\s\\-']{1,50}$" ;
    autoins:formHelpText "Enter your legal first name" .
```

## Output

The generator creates a complete HTML file with:

- Modern, responsive CSS styling
- Interactive JavaScript for tab switching
- Form validation based on ontology rules
- Proper field labeling and help text
- Mobile-friendly layout

## Example Output

```html
<!DOCTYPE html>
<html>
<head>
    <title>Ontology-Driven Insurance Forms</title>
    <!-- CSS and meta tags -->
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🚗📄 Ontology-Driven Insurance Forms</h1>
        </div>
        <div class="main-content">
            <!-- Category tabs and form fields -->
        </div>
    </div>
    <!-- JavaScript for interactivity -->
</body>
</html>
```

## Development

### Adding New Field Types

1. Update the `parseFieldType()` function in `main.go`
2. Add corresponding HTML template logic
3. Update category mappings if needed

### Extending Categories

1. Add new category to `categoryConfigs` in `OrganizeIntoCategories()`
2. Update field mappings in `categoryMappings`
3. Add corresponding HTML template section

## Requirements

- Go 1.21 or later
- TTL ontology file with proper structure
- Web browser to view generated forms

## License

This project is part of the AUTOINS insurance application suite.
