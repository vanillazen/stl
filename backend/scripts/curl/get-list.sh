#!/bin/bash

################################################################################
# Script: get-list.sh
# Description: A shell script to call `get-list` API endpoint using curl.
#
# Usage:
#   ./api_call.sh [OPTIONS]
#
# Options:
#   -h, --host HOST       Specify the host (default: localhost)
#   -p, --port PORT       Specify the port (default: 8081)
#   -l, --list-id LIST_ID Specify the list ID (default: fa7af80c-63e0-413c-aa3d-b7e417459a69)
#
# Example usage:
#   ./scripts/curl/get-list.sh -h localhost -p 8080 -l fa7af80c-63e0-413c-aa3d-b7e417459a69
#
################################################################################

# Default values
default_host="localhost"
default_port="8080"
default_list_id="cdc7a443-3c6a-431b-b45a-b14735953a19"

# Parse command-line arguments
while [[ $# -gt 0 ]]; do
  case "$1" in
    -h|--host)
      host="$2"
      shift 2
      ;;
    -p|--port)
      port="$2"
      shift 2
      ;;
    -l|--list-id)
      list_id="$2"
      shift 2
      ;;
    *)
      echo "Invalid argument: $1"
      exit 1
      ;;
  esac
done

# Values
host="${host:-$default_host}"
port="${port:-$default_port}"
list_id="${list_id:-$default_list_id}"

# API endpoint
url="http://$host:$port/api/v1/lists/$list_id"

# Debug
echo "Calling API endpoint: GET $url"

# Call
curl "$url"
