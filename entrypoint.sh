#!/bin/sh -eu

DEBUG=${DEBUG:-"0"}

if [ "${DEBUG}" = "1" ]; then
    set -x
fi

export PATH=/usr/local/bin:/usr/sbin:/sbin:${PATH}

#
# Main
#

# if command starts with an option, prepend daikin-ac-exporter
if [ "${1:0:1}" = '-' ]; then
     set -- daikin-ac-exporter "$@"
fi

exec "$@"
