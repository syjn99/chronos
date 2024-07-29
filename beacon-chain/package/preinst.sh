#!/bin/bash

set -e

SERVICE_USER=overprotocol-beacon

# Create the service account, if needed
getent passwd $SERVICE_USER > /dev/null || useradd -s /bin/false --no-create-home --system --user-group $SERVICE_USER

# Create directories
mkdir -p /etc/overprotocol
mkdir -p /var/lib/overprotocol
install -d -m 0700 -o $SERVICE_USER -g $SERVICE_USER /var/lib/overprotocol/beacon-chain