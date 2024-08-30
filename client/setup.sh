#!/bin/bash

CLIENT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )
ARG_PROM_IP="$1"

# Function to read input with a prompt
read_input() {
  local prompt="$1"
  local input_var="$2"
  read -p "$prompt: " "$input_var"
}

# Function to validate an IP address
validate_ip() {
  ip="$1"
  echo "$ip" | grep -E '^([0-9]{1,3}\.){3}[0-9]{1,3}$' >/dev/null
  if [ $? -eq 0 ]; then
    return 0
  else
    echo "Invalid IP address. Please enter a valid IP address."
    return 1
  fi
}

# Function to validate a domain name
validate_domain() {
  domain="$1"
  echo "$domain" | grep -E '^[a-zA-Z0-9.-]+$' >/dev/null
  if [ $? -eq 0 ]; then
    return 0
  else
    echo "Invalid domain name. Please enter a valid domain name."
    return 1
  fi
}

echo "Setting up client config"

cp $CLIENT_DIR/nginx.conf_sample $CLIENT_DIR/nginx.conf

WHITELIST_IP="$ARG_PROM_IP"

if [[ -z $WHITELIST_IP ]]; then
  # Get whitelisted IP address
  while true; do
    read_input "Enter Prometheus IP" whitelist_ip
    if validate_ip "$whitelist_ip"; then
      WHITELIST_IP="$whitelist_ip"
      break
    fi
  done
fi

sed -i 's/__PROMETHEUS_IP__/'"$WHITELIST_IP"'/g' $CLIENT_DIR/nginx.conf

echo "Setup complete!"
