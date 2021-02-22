#!/bin/bash

URL=$(printenv PSI_BACKEND_URL)
USER=$(printenv PSI_CONNECT_USER)
PASS=$(printenv PSI_CONNECT_PASSWORD)

BASEDIR=$( dirname $(readlink -f "$0") )
COUNTER=0

while [[ $COUNTER -lt 10 ]]; do
    TOKEN=$(curl "$URL" -H "Content-Type: application/json" \
    --data-binary "{\"query\":\"{ authenticateUser( input: { email: \\\"$USER\\\" password: \\\"$PASS\\\" } ) { token } } \"}" 2>/dev/null \
    | jq -r '.data.authenticateUser.token')
    if [[ $TOKEN == "" ]]; then
        sleep 5
    else
        break
    fi
done

if [[ $TOKEN == "" ]]; then
    echo "Token retrieval failed. Exiting..."
    exit 1
fi

$BASEDIR/process_pending_mail.sh $URL $TOKEN