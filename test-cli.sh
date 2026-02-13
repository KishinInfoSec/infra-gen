#!/bin/bash

echo "=== Infra-Gen CLI Tool Test ==="
echo

# Test 1: Help command
echo "1. Testing help command..."
./infra-gen --help | head -10
echo

# Test 2: List presets
echo "2. Testing list presets..."
./infra-gen list presets | head -15
echo

# Test 3: List categories
echo "3. Testing list categories..."
./infra-gen list categories
echo

# Test 4: Initialize project
echo "4. Testing project initialization..."
rm -f infra-gen.yml *.yml *.tf
./infra-gen init web-app --name test-project --environment development
echo

# Test 5: List project
echo "5. Testing list project..."
./infra-gen list project | head -15
echo

# Test 6: Validate configuration
echo "6. Testing validation..."
./infra-gen validate | head -10
echo

# Test 7: Generate Docker Compose
echo "7. Testing Docker generation..."
./infra-gen generate docker
echo "Generated files:"
ls -la *.yml 2>/dev/null | head -5
echo

# Test 8: Generate Terraform
echo "8. Testing Terraform generation..."
./infra-gen generate terraform
echo "Generated files:"
ls -la *.tf 2>/dev/null | head -5
echo

# Test 9: Generate all
echo "9. Testing full generation..."
./infra-gen generate all
echo "Total generated files:"
ls -la *.yml *.tf 2>/dev/null | wc -l
echo

echo "=== All Tests Complete ==="