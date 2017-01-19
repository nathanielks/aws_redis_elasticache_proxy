#!/bin/bash
hostport=$1
secret=$2
signature=$(echo -n "${hostport}${secret}" | shasum -a 256 | cut -d ' ' -f 1)
echo $(echo -n "$hostport $signature" | base64)
