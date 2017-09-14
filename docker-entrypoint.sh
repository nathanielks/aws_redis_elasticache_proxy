#!/bin/sh

set -euxo pipefail

function get_value_from_result () {
  echo "${1}" | jq --arg name "${2}" -r '.[] | select(.Name==$name) | .Value'
}

DEFAULT_HOST_PORT="0.0.0.0:6379"
HOST_PORT="${1:-$DEFAULT_HOST_PORT}"; shift
CERT="server.crt"
KEY="server.key"

SSM_PARAMS="${SSM_TLS_KEY} ${SSM_TLS_CERT} ${SSM_SHARED_KEY}"
RESULT="$(aws ssm get-parameters --names ${SSM_PARAMS} --with-decryption --query 'Parameters[*]' | jq -rc '.')"
KEY_TEXT="$(get_value_from_result "${RESULT}" "${SSM_TLS_KEY}" | base64 -d)"
CERT_TEXT="$(get_value_from_result "${RESULT}" "${SSM_TLS_CERT}" | base64 -d)"
PASSWORD="$(get_value_from_result "${RESULT}" "${SSM_SHARED_KEY}" | base64 -d)"

echo "${CERT_TEXT}" | dd of="${CERT}" status=none
echo "${KEY_TEXT}" | dd of="${KEY}" status=none

redis_auth_proxy "$HOST_PORT" "$PASSWORD" "$CERT" "$KEY"
