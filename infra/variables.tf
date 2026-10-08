variable "aws_region" {
  type    = string
  default = "us-east-1"
}

variable "project" {
  type    = string
  default = "ripple"
}

variable "vpc_cidr" {
  type    = string
  default = "10.20.0.0/16"
}

variable "public_subnet_cidr" {
  type    = string
  default = "10.20.1.0/24"
}

variable "api_instance_type" {
  type    = string
  default = "t3.small"
}

variable "web_instance_type" {
  type    = string
  default = "t3.small"
}
