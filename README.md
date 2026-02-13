# Infra-Gen

A CLI tool for generating Docker Compose, Ansible, and Terraform configurations to deploy any type of project.

## Features

- **Project Presets**: Ready-to-use templates for web apps, microservices, databases, ML projects
- **Multi-Target Generation**: Generate Docker Compose, Ansible, and Terraform configurations
- **Simple Presets**: Easy-to-use preset system for quick project setup
- **Validation**: Built-in validation to catch configuration errors early
- **Hybrid Templates**: Embedded templates with support for custom templates

## Installation

```bash
go build -o infra-gen
sudo mv infra-gen /usr/local/bin/
```

## Quick Start

### 1. List Available Presets

```bash
infra-gen list presets
```

### 2. Initialize a New Project

```bash
# Initialize a web application project
infra-gen init web-app --name my-web-app --environment development

# Initialize a microservice project
infra-gen init microservice --name my-microservice --environment staging
```

### 3. Generate Infrastructure

```bash
# Generate all configurations
infra-gen generate all

# Generate specific targets
infra-gen generate docker
infra-gen generate ansible
infra-gen generate terraform
```

### 4. Validate Configuration

```bash
infra-gen validate
```

## Available Commands

### `init <preset>`
Initialize a new project from a preset template.

```bash
infra-gen init web-app --name my-project --environment production --output ./my-project
```

**Flags:**
- `--name, -n`: Project name (required)
- `--environment, -e`: Environment (development, staging, production)
- `--output, -o`: Output directory

### `generate [target]`
Generate infrastructure configurations.

```bash
infra-gen generate docker    # Generate Docker Compose only
infra-gen generate all      # Generate all configurations
```

**Flags:**
- `--config, -c`: Project configuration file (default: infra-gen.yml)
- `--output, -o`: Output directory

### `list [type]`
List available presets and project information.

```bash
infra-gen list presets      # Show all presets
infra-gen list categories   # Show preset categories
infra-gen list project      # Show current project details
```

### `validate`
Validate project configuration.

```bash
infra-gen validate --target all    # Validate all targets
infra-gen validate --target docker # Validate Docker only
```

## Project Presets

### Web Application (`web-app`)
- Frontend (Nginx)
- Backend API (Node.js)
- Database (PostgreSQL)
- Environment variables for development

### Microservice (`microservice`)
- API Gateway (Nginx)
- Multiple microservices
- Redis cache
- Service discovery

### Database (`database`)
- PostgreSQL
- MySQL (optional)
- MongoDB (optional)
- Persistent volumes

## Configuration File

The `infra-gen.yml` file contains your project configuration:

```yaml
name: my-web-app
type: web-app
description: Basic web application with frontend, backend, and database
environment: development
version: "1.0.0"
services:
  - name: frontend
    type: frontend
    image: nginx:alpine
    ports:
      - container: 80
        protocol: tcp
    enabled: true
  - name: api
    type: api
    image: node:18-alpine
    ports:
      - container: 8080
        protocol: tcp
    environment:
      NODE_ENV: development
      DB_HOST: database
    depends_on:
      - database
    enabled: true
  - name: database
    type: database
    image: postgres:15
    ports:
      - container: 5432
        protocol: tcp
    volumes:
      - source: db_data
        target: /var/lib/postgresql/data
        type: volume
    environment:
      POSTGRES_DB: webapp
      POSTGRES_USER: admin
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    enabled: true
variables:
  DB_PASSWORD: "change-me"
created_at: "2024-01-01T00:00:00Z"
updated_at: "2024-01-01T00:00:00Z"
```

## Generated Files

### Docker Compose
- `docker-compose.yml` - Main Docker Compose configuration
- `.env` - Environment variables file

### Ansible
- `playbook.yml` - Main Ansible playbook
- `inventory.yml` - Ansible inventory
- `requirements.yml` - Ansible requirements

### Terraform
- `main.tf` - Main Terraform configuration
- `variables.tf` - Input variables
- `outputs.tf` - Output values
- `provider.tf` - Provider configuration

## Examples

### Web Application Example

```bash
# Initialize project
infra-gen init web-app --name blog --name development

# Generate configurations
infra-gen generate all

# Deploy with Docker
docker-compose up -d

# Deploy with Ansible
ansible-playbook -i inventory.yml playbook.yml

# Deploy with Terraform
terraform init
terraform apply
```

### Microservice Example

```bash
# Initialize microservice project
infra-gen init microservice --name api-service --name production

# Generate and deploy
infra-gen generate terraform
cd terraform
terraform init
terraform apply
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## License

MIT License# infra-gen
