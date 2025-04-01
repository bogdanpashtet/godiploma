#!/usr/bin/env bash

# Determine the directory where the script is located
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"

# Execute run.sh and check compilation, saving the output to a variable
OUTPUT=$("$DIR/run.sh" go build ./protos/gen/go/... 2>&1)

# Check if the last executed command failed (non-zero exit status)
if [ $? -ne 0 ]; then
  echo "$OUTPUT"  # Print the output of the failed command
  echo "Error: Generated Go code does not compile. Please refer to the 'Golang compilation errors' section in TROUBLESHOOTING.md for fix steps."
  exit 1
fi
