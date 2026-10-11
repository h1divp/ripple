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
    set -euo pipefail

    workdir=$(mktemp -d)
    trap 'rm -rf "$workdir"' EXIT
    cd "$workdir"

    dnf install -y docker
    command -v curl >/dev/null 2>&1 || dnf install -y curl

    COMPOSE_VERSION=v5.6.0
    COMPOSE_ASSET=docker-compose-linux-x86_64
    COMPOSE_CHECKSUMS=checksums.txt
    COMPOSE_URL=https://github.com/docker/compose/releases/download/$COMPOSE_VERSION

    curl -fsSL "$COMPOSE_URL/$COMPOSE_ASSET" -o "$COMPOSE_ASSET"
    curl -fsSL "$COMPOSE_URL/$COMPOSE_CHECKSUMS" -o "$COMPOSE_CHECKSUMS"
    sha256sum --ignore-missing -c "$COMPOSE_CHECKSUMS"
    install -D -o root -g root -m 0755 "$COMPOSE_ASSET" \
      /usr/local/lib/docker/cli-plugins/docker-compose

    SOPS_VERSION=3.13.3
    SOPS_ASSET=sops-v$SOPS_VERSION.linux.amd64
    SOPS_CHECKSUMS=sops-v$SOPS_VERSION.checksums.txt
    SOPS_URL=https://github.com/getsops/sops/releases/download/v$SOPS_VERSION

    curl -fsSL "$SOPS_URL/$SOPS_ASSET" -o "$SOPS_ASSET"
    curl -fsSL "$SOPS_URL/$SOPS_CHECKSUMS" -o "$SOPS_CHECKSUMS"
    sha256sum --ignore-missing -c "$SOPS_CHECKSUMS"
    install -o root -g root -m 0755 "$SOPS_ASSET" /usr/local/bin/sops

    systemctl enable --now docker
    usermod -aG docker ec2-user
    docker compose version
    sops --version
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
