#!/bin/bash

set -e

declare -a arr=("./api" "./router" "./utils")
for i in "${arr[@]}"
do
	OUTPUT="$(goimports -l -d $i)"
	if [[ $OUTPUT ]]; then
		echo "[WARNING] The following files contain goimports errors"
		echo "[WARNING] $OUTPUT"
		echo "[WARNING] The goimports command must be run for these files"
	fi
done
