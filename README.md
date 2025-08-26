# AutoIns - Insurance Quote Application

A comprehensive insurance quote system with ontology-driven form generation, intelligent document processing, and modern React frontends.

## 🏗️ Project Structure

```
insurance-quote-app/
├── ontology/                    # TTL Ontologies
│   ├── autoins.ttl             # Main insurance ontology
│   ├── passport.ttl            # Passport document ontology
│   ├── driving_licence.ttl     # Driving licence document ontology
│   └── settings.ttl            # Settings application ontology
├── insurance-frontend/          # React insurance application
├── settings-frontend/           # React settings application
├── form_generator/             # Go backend for form generation
├── sessions/                   # Session management
├── static/                     # Static assets
└── templates/                  # HTML templates
```

## 🚀 Applications

### 1. Insurance Frontend (`insurance-frontend/`)

**Modern React application for insurance quote processing with intelligent document upload and processing.**

#### Features:
- **Ontology-Driven Forms**: Dynamic form generation based on TTL ontologies
- **Intelligent Document Processing**: 
  - Passport recognition and field extraction
  - Driving licence front/back processing
  - Name matching and driver creation/updates
- **Real-time Chatbot**: AI assistant for document processing feedback
- **Multi-language Support**: English and German
- **Progress Tracking**: Visual progress indicators
- **Responsive Design**: Modern UI with Tailwind CSS and Flowbite

#### Technology Stack:
- React 18 + TypeScript
- Tailwind CSS
- Flowbite React Components
- Lucide React Icons

#### Key Features:
- **Document Upload**: Drag & drop file upload with OCR simulation
- **Smart Matching**: Name + DOB matching for existing drivers
- **Auto-fill**: Automatic form population from documents
- **Visual Feedback**: Real-time processing status and notifications
- **Passport Processing**: Complete passport field extraction
- **Driving Licence Processing**: Front and back side processing

### 2. Settings Frontend (`settings-frontend/`)

**Comprehensive settings management application for insurance configuration.**

#### Features:
- **Bank Account Management**: Up to 8 accounts with Open Banking support
- **Credit Card Management**: Up to 8 cards with digital wallet integration
- **Communication Channels**: Email, SMS, Voice, Secure Messenger
- **Security Settings**: PCI-compliant credential management
- **Document Upload**: Configuration document processing
- **Multi-language Support**: i18n ready

#### Technology Stack:
- React 18 + TypeScript
- Tailwind CSS
- Flowbite React Components
- Lucide React Icons

## 🧠 Ontology System

### Core Ontologies:

#### `autoins.ttl` - Main Insurance Ontology
- Complete driver, vehicle, and policy definitions
- SHACL validation rules
- Document processing rules
- Multi-language support

#### `passport.ttl` - Passport Document Ontology
- All passport field definitions
- Insurance field mapping
- Processing rules and validation
- Front/back side differentiation

#### `driving_licence.ttl` - Driving Licence Document Ontology
- Front side fields (personal info, licence details)
- Back side fields (entitlements, restrictions)
- Category A-Z entitlement mapping
- Processing rules for both sides

#### `settings.ttl` - Settings Application Ontology
- Bank account configurations
- Credit card management
- Communication channel settings
- Security and compliance rules

## 🛠️ Setup Instructions

### Prerequisites
- Node.js 18+
- Go 1.21+
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd insurance-quote-app
   ```

2. **Install Insurance Frontend Dependencies**
   ```bash
   cd insurance-frontend
   npm install
   ```

3. **Install Settings Frontend Dependencies**
   ```bash
   cd ../settings-frontend
   npm install
   ```

4. **Install Go Dependencies**
   ```bash
   cd ..
   go mod tidy
   ```

### Running the Applications

#### Insurance Frontend
```bash
cd insurance-frontend
npm start
```
Access at: http://localhost:3000

#### Settings Frontend
```bash
cd settings-frontend
npm start
```
Access at: http://localhost:3001

#### Go Backend (Form Generator)
```bash
go run main.go
```

## 🔧 Development

### Adding New Document Types

1. Create new ontology file in `ontology/` directory
2. Define document fields and processing rules
3. Update `autoins.ttl` with new document class
4. Add processing logic in React frontend

### Ontology Development

The system uses TTL (Turtle) format for ontologies with:
- RDF Schema for class definitions
- SHACL for validation rules
- Custom properties for UI mapping
- Multi-language support

### Frontend Development

Both React applications use:
- TypeScript for type safety
- Tailwind CSS for styling
- Flowbite for UI components
- Lucide React for icons

## 🎯 Key Features

### Document Processing
- **OCR Simulation**: Document type recognition
- **Field Extraction**: Automatic data extraction
- **Smart Matching**: Name and DOB-based driver matching
- **Conflict Resolution**: Handle data conflicts intelligently
- **Real-time Feedback**: Chatbot integration for processing status

### Form Generation
- **Dynamic Forms**: Ontology-driven form creation
- **Validation**: SHACL-based field validation
- **Conditional Logic**: Show/hide fields based on conditions
- **Multi-language**: Full i18n support

### Settings Management
- **Bank Integration**: Open Banking support
- **Digital Wallets**: Apple Pay, Google Pay integration
- **Communication**: Multi-channel communication setup
- **Security**: PCI-compliant credential storage

## 📝 License

This project is licensed under the MIT License.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

## 📞 Support

For support and questions, please open an issue in the GitHub repository.
