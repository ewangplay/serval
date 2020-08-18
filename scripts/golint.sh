#!/bin/bash

set -e

PKGS=`go list github.com/ewangplay/serval/...`

for i in "$PKGS"
do
	OUTPUT="$(golint $i)"
	if [[ $OUTPUT ]]; then
		echo "$OUTPUT"
	fi
done