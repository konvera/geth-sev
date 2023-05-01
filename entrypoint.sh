#!/bin/bash
set -e

mount -t securityfs securityfs /sys/kernel/security

exec /geth/geth "$@"
