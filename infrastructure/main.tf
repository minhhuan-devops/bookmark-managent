// network

module "vpc" {
  source = "terraform-aws-modules/vpc/aws"

  name = "VBAC-VPC"
  cidr = "10.0.0.0/16"

  azs             = ["ap-southeast-1a", "ap-southeast-1b"]
  private_subnets = ["10.0.1.0/24", "10.0.2.0/24"]
  public_subnets  = ["10.0.101.0/24", "10.0.102.0/24"]

  enable_nat_gateway = true
  single_nat_gateway = true

  enable_dns_hostnames = true
  enable_dns_support   = true

  tags = {
    Project = "bookmark-management"
  }
}

// redis
module "valkey_cache" {
  source  = "terraform-aws-modules/elasticache/aws"
  version = "~> 1.11.0"

  replication_group_id = "vbac-valkey-group"
  description          = "Valkey cluster for VBAC project"
  cluster_id           = "vbac-valkey-cluster"
  engine               = "valkey"
  engine_version       = "7.2"
  node_type            = "cache.t4g.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.valkey7"
  port                 = 6379

  vpc_id     = module.vpc.vpc_id
  subnet_ids = [module.vpc.private_subnets[0]]

  security_group_rules = {
    ingress_vpc = {
      description = "Allow inbound from VPC"
      cidr_ipv4   = module.vpc.vpc_cidr_block
    }
  }

  tags = {
    Project = "bookmark-management"
  }
}
// loadbalancer (internal - chỉ API Gateway truy cập được)
module "loadbalancer" {
  source = "terraform-aws-modules/alb/aws"

  name     = "api-loadbalancer"
  vpc_id   = module.vpc.vpc_id
  subnets  = module.vpc.private_subnets
  internal = true

  enable_deletion_protection = false

  create_security_group = true
  security_group_ingress_rules = {
    vpc-http = {
      from_port   = 80
      to_port     = 80
      ip_protocol = "tcp"
      description = "HTTP from VPC (API Gateway VPC Link)"
      cidr_ipv4   = module.vpc.vpc_cidr_block
    }
  }
  security_group_egress_rules = {
    all = {
      ip_protocol = "-1"
      cidr_ipv4   = "0.0.0.0/0"
    }
  }

  target_groups = {
    ecs-target = {
      name_prefix       = "bm-"
      port              = 8080
      protocol          = "HTTP"
      target_type       = "ip"
      create_attachment = false

      health_check = {
        enabled             = true
        interval            = 30
        path                = "/health-check"
        port                = "traffic-port"
        healthy_threshold   = 3
        unhealthy_threshold = 3
        timeout             = 5
        protocol            = "HTTP"
        matcher             = "200"
      }
    }
  }

  listeners = {
    http = {
      port     = 80
      protocol = "HTTP"
      forward = {
        target_group_key = "ecs-target"
      }
    }
  }

  tags = {
    Project = "bookmark-management"
  }
}

// ecs
module "ecs" {
  source = "terraform-aws-modules/ecs/aws"

  cluster_name = "bookmark-management-cluster"

  create_cloudwatch_log_group = true

  cluster_capacity_providers = ["FARGATE"]
  default_capacity_provider_strategy = {
    FARGATE = {
      weight = 1
      base   = 0
    }
  }

  services = {
    backend-services = {
      cpu    = 256
      memory = 512

      subnet_ids            = module.vpc.private_subnets
      create_security_group = true

      security_group_ingress_rules = {
        ingress_alb = {
          description                  = "Allow traffic from ALB"
          from_port                    = 8080
          to_port                      = 8080
          ip_protocol                  = "tcp"
          referenced_security_group_id = module.loadbalancer.security_group_id
        }
      }

      security_group_egress_rules = {
        egress_all = {
          description = "Allow all outbound traffic"
          ip_protocol = "-1"
          cidr_ipv4   = "0.0.0.0/0"
        }
      }

      load_balancer = {
        service = {
          target_group_arn = module.loadbalancer.target_groups["ecs-target"].arn
          container_name   = "backend"
          container_port   = 8080
        }
      }

      container_definitions = {
        backend = {
          name = "backend"

          cpu       = 256
          memory    = 512
          essential = true

          image                  = "senn404/bookmark-management:v3"
          readonlyRootFilesystem = false

          environment = [
            {
              name  = "REDIS_ADDRESS"
              value = module.valkey_cache.replication_group_primary_endpoint_address
            },
            {
              name  = "REDIS_PORT"
              value = tostring(module.valkey_cache.replication_group_port)
            },
            {
              name  = "PORT"
              value = "8080"
            },
            {
              name  = "REDIS_USE_TLS"
              value = "true"
            },
            {
              name  = "SWAGGER_HOST"
              value = replace(module.api_gateway.api_endpoint, "https://", "")
            }
          ]

          portMappings = [
            {
              name          = "backend"
              containerPort = 8080
              protocol      = "tcp"
            }
          ]
        }
      }
    }
  }
}

// api gateway
resource "aws_security_group" "vpc_link" {
  name_prefix = "apigw-vpc-link-"
  description = "Security group for API Gateway VPC Link"
  vpc_id      = module.vpc.vpc_id

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Project = "bookmark-management"
  }
}

module "api_gateway" {
  source = "terraform-aws-modules/apigateway-v2/aws"

  name          = "bookmark-api"
  description   = "Bookmark Management API Gateway"
  protocol_type = "HTTP"

  create_domain_name = false

  cors_configuration = {
    allow_headers = ["content-type", "authorization"]
    allow_methods = ["*"]
    allow_origins = ["*"]
  }

  # VPC Link
  vpc_links = {
    ecs-vpc-link = {
      name               = "ecs-vpc-link"
      security_group_ids = [aws_security_group.vpc_link.id]
      subnet_ids         = module.vpc.private_subnets
    }
  }

  # Routes — forward tất cả traffic đến ALB
  routes = {
    "GET /gen-pass" = {
      integration = {
        connection_type = "VPC_LINK"
        uri             = module.loadbalancer.listeners["http"].arn
        type            = "HTTP_PROXY"
        method          = "GET"
        vpc_link_key    = "ecs-vpc-link"
      }
    }
    "GET /health-check" = {
      integration = {
        connection_type = "VPC_LINK"
        uri             = module.loadbalancer.listeners["http"].arn
        type            = "HTTP_PROXY"
        method          = "GET"
        vpc_link_key    = "ecs-vpc-link"
      }
    }
    "POST /links/shorten" = {
      integration = {
        connection_type = "VPC_LINK"
        uri             = module.loadbalancer.listeners["http"].arn
        type            = "HTTP_PROXY"
        method          = "POST"
        vpc_link_key    = "ecs-vpc-link"
      }
    }
    "GET /links/redirect/{code}" = { # ← {code} thay vì :code
      integration = {
        connection_type = "VPC_LINK"
        uri             = module.loadblance.listeners["http"].arn
        type            = "HTTP_PROXY"
        method          = "GET"
        vpc_link_key    = "ecs-vpc-link"
      }
    }
    "GET /swagger/{proxy+}" = {
      integration = {
        connection_type = "VPC_LINK"
        uri             = module.loadbalancer.listeners["http"].arn
        type            = "HTTP_PROXY"
        method          = "GET"
        vpc_link_key    = "ecs-vpc-link"
      }
    }
  }

  # Access logs
  stage_access_log_settings = {
    create_log_group            = true
    log_group_retention_in_days = 7
    format = jsonencode({
      requestId      = "$context.requestId"
      ip             = "$context.identity.sourceIp"
      httpMethod     = "$context.httpMethod"
      path           = "$context.path"
      status         = "$context.status"
      responseLength = "$context.responseLength"
    })
  }

  tags = {
    Project = "bookmark-management"
  }
}
