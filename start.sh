#!/bin/sh

if [ -z $NEW_RELIC_LICENSE_KEY ]; then
   echo "Please configure NEW_RELIC_LICENSE_KEY"
   exit
fi

echo "Starting infra-lite"
exec /infra-lite
