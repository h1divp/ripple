locals {
  instances = {
    api = var.api_instance_type
    web = var.web_instance_type
  }
}

resource "aws_instance" "app" {
  for_each = local.instances

  ami                         = data.aws_ami.al2023.id
  instance_type               = each.value
  subnet_id                   = aws_subnet.public.id
  vpc_security_group_ids      = [aws_security_group.app.id]
  iam_instance_profile        = aws_iam_instance_profile.ec2.name
  associate_public_ip_address = true

  user_data_replace_on_change = true
  user_data                   = <<-EOF
    #!/bin/bash
    set -euxo pipefail

    COMPOSE_VERSION='v5.6.0'
    ARCH=x86_64

    tmp=$(mktemp)
    dnf install -y docker
    curl -fsSL \
      "https://github.com/docker/compose/releases/download/$COMPOSE_VERSION/docker-compose-linux-$ARCH" \
      -o "$tmp"
    install -D -o root -g root -m 0755 \
      "$tmp" /usr/local/lib/docker/cli-plugins/docker-compose
    rm -rf "$tmp"

    systemctl enable --now docker
    usermod -aG docker ec2-user
    docker compose version
  EOF

  root_block_device {
    encrypted   = true
    volume_type = "gp3"
    volume_size = 30
  }

  metadata_options {
    http_endpoint = "enabled"
    http_tokens   = "required"
  }

  tags = {
    Name    = "${var.project}-${each.key}"
    Project = var.project
  }
}

resource "aws_eip" "app" {
  for_each = aws_instance.app

  domain   = "vpc"
  instance = each.value.id

  tags = {
    Name    = "${var.project}-${each.key}"
    Project = var.project
  }
}
