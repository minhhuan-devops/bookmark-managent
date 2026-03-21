// VPC
output "vpc_id" {
  description = "VPC ID"
  value       = module.vpc.vpc_id
}

output "public_subnets" {
  description = "Public subnet IDs"
  value       = module.vpc.public_subnets
}

output "private_subnets" {
  description = "Private subnet IDs"
  value       = module.vpc.private_subnets
}

// ALB
output "alb_dns_name" {
  description = "DNS name of the ALB - use this to access your app"
  value       = module.loadbalancer.dns_name
}

output "alb_arn" {
  description = "ARN of the ALB"
  value       = module.loadbalancer.arn
}

output "alb_security_group_id" {
  description = "Security group ID of the ALB"
  value       = module.loadbalancer.security_group_id
}

output "target_group_arns" {
  description = "Target group ARNs"
  value       = { for k, v in module.loadbalancer.target_groups : k => v.arn }
}

// Valkey (Redis)
output "valkey_endpoint" {
  description = "Valkey cache endpoint"
  value       = module.valkey_cache.cluster_cache_nodes
}

// ECS
output "ecs_cluster_name" {
  description = "ECS cluster name"
  value       = module.ecs.cluster_name
}

output "ecs_cluster_arn" {
  description = "ECS cluster ARN"
  value       = module.ecs.cluster_arn
}

output "ecs_services" {
  description = "ECS service details"
  value       = { for k, v in module.ecs.services : k => v.name }
}

// API Gateway
output "api_gateway_url" {
  description = "API Gateway URL - use this to access your app"
  value       = module.api_gateway.api_endpoint
}
