#!/bin/bash

echo "🚀 Generating HTML forms from TTL ontology..."

# Run the form generator
go run main.go

if [ $? -eq 0 ]; then
    echo "✅ Forms generated successfully!"
    echo "🌐 Opening generated_forms.html in browser..."
    
    # Open in default browser
    if command -v open >/dev/null 2>&1; then
        open generated_forms.html
    elif command -v xdg-open >/dev/null 2>&1; then
        xdg-open generated_forms.html
    else
        echo "📄 Generated file: $(pwd)/generated_forms.html"
        echo "Please open this file in your browser manually."
    fi
else
    echo "❌ Failed to generate forms"
    exit 1
fi
