#!/bin/sh

# credit: https://gist.github.com/cherti/61ec48deaaab7d288c9fcf17e700853a

name="intergration-testing" # Generate a random number
url='http://alertmanager:9093/api/v2/alerts'

echo "firing up alert $name"

# Get the current time in RFC 3339 format (UTC)
now=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Manually calculate the time 1 minute ago
# Extract hours, minutes, and seconds separately for adjustment
minute_ago=$(date -u +%s) # Get current time in seconds since epoch
minute_ago=$(($minute_ago - 60)) # Subtract 60 seconds
one_minute_ago=$(date -u -d @$minute_ago +%Y-%m-%dT%H:%M:%SZ) # Convert back to RFC 3339 format

# Firing alert with v2 API
curl -XPOST "$url" -H "Content-Type: application/json" -d "[{
    \"status\": \"firing\",
    \"labels\": {
        \"alertname\": \"$name\",
        \"service\": \"dummy-service\",
        \"severity\": \"warning\",
        \"job\": \"node_exporter\",
        \"instance\": \"$name\"
    },
    \"annotations\": {
        \"summary\": \"Testing local nanda bawana\",
        \"title\": \"Title testing local\",
        \"description\": \"Description testing local\"
    },
    \"startsAt\": \"$now\",
    \"endsAt\": \"0001-01-01T00:00:00Z\",
    \"generatorURL\": \"http://localhost\"
}]"

echo ""

echo "press enter to resolve alert"
read _input

echo "sending resolve"

# Resolving alert with v2 API
curl -XPOST "$url" -H "Content-Type: application/json" -d "[{
    \"status\": \"resolved\",
    \"labels\": {
        \"alertname\": \"$name\",
        \"service\": \"dummy-service\",
        \"severity\": \"warning\",
        \"instance\": \"$name\"
    },
    \"annotations\": {
        \"summary\": \"High latency is high!\"
    },
    \"startsAt\": \"$one_minute_ago\",
    \"endsAt\": \"$now\",
    \"generatorURL\": \"http://prometheus.int.example.net/<generating_expression>\"
}]"

echo ""
