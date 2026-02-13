# Input variables
variable "project_name" {
  description = "Name of the project"
  type        = string
  default     = "test-project"
}

variable "environment" {
  description = "Environment (dev, staging, prod)"
  type        = string
  default     = "development"
}

variable "frontend_image" {
  description = "Docker image for frontend"
  type        = string
  default     = "nginx:alpine"
}

variable "api_image" {
  description = "Docker image for api"
  type        = string
  default     = "node:18-alpine"
}

variable "database_image" {
  description = "Docker image for database"
  type        = string
  default     = "postgres:15"
}

