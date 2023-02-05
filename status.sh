#!/bin/sh

echo "STABLE_MAIN_VERSION $(make --quiet --no-print-directory tag TAG=3.72.3 )"
echo "STABLE_COLLECTOR_VERSION $(cat COLLECTOR_VERSION)"
echo "STABLE_SCANNER_VERSION $(cat SCANNER_VERSION)"
echo "STABLE_GIT_SHORT_SHA $(make --quiet --no-print-directory shortcommit)"
echo "BUILD_TIMESTAMP $(date '+%s')"
