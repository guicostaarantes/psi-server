#!/bin/bash

while :; do
    sleep 30
    ERROR_PROCESS_PENDING_MAIL=$(curl "$1" -H "Content-Type: application/json" -H "Authorization: $2" \
    --data-binary '{"query":"mutation { processPendingMail }"}' 2>/dev/null \
    | jq -r '.errors')
    if [[ $ERROR_PROCESS_PENDING_MAIL != "null" ]]; then
        echo "Error while processing pending mail: "$ERROR_PROCESS_PENDING_MAIL
        exit 1
    fi
done