#!/bin/sh
set -e
awslocal sqs create-queue --queue-name bridgr-skill-gap >/dev/null
echo "LocalStack: queue bridgr-skill-gap ready"
