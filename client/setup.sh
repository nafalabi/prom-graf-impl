#!/bin/bash

# Function to read input with a prompt
read_input() {
    local prompt="$1"
    local input_var="$2"
    read -p "$prompt: " "$input_var"
}

# Function to validate an IP address
validate_ip() {
    ip="$1"
    echo "$ip" | grep -E '^([0-9]{1,3}\.){3}[0-9]{1,3}$' > /dev/null
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
    echo "$domain" | grep -E '^[a-zA-Z0-9.-]+$' > /dev/null
    if [ $? -eq 0 ]; then
        return 0
    else
        echo "Invalid domain name. Please enter a valid domain name."
        return 1
    fi
}

echo "Setting up client config"

cp nginx.conf_sample nginx.conf

# Get whitelisted IP address
while true; do
    read_input "Enter Prometheus IP" whitelist_ip
    if validate_ip "$whitelist_ip"; then
        sed -i 's/__PROMETHEUS_IP__/'"$whitelist_ip"'/g' nginx.conf
        break
    fi
done

# Get current domain name
while true; do
    read_input "Enter domain name" domain_name
    if validate_domain "$domain_name"; then
        sed -i 's/__DOMAIN_NAME__/'"$domain_name"'/g' nginx.conf
        break
    fi
done

echo "Setup complete!"
