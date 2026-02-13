# Output values
output "frontend_url" {
  description = "URL for frontend"
  value = "http://${aws_instance.frontend.public_ip}:80"
}

output "api_url" {
  description = "URL for api"
  value = "http://${aws_instance.api.public_ip}:8080"
}

output "database_url" {
  description = "URL for database"
  value = "http://${aws_instance.database.public_ip}:5432"
}

