#!/bin/bash

# Script to decode Base64 from file and save result to ./files/res/
# with a timestamp appended to the filename and the extension provided
# as the second argument.

# Check for required arguments
if [ "$#" -ne 2 ]; then
    echo "Usage: $0 <input_base64_txt_file> <output_file_type>"
    echo "  input_base64_txt_file: File with Base64 string (e.g. from ./files/byte/)."
    echo "  output_file_type:      Desired file extension for the output (e.g., 'png', 'jpeg', 'pdf')."
    echo "                         Output filename will be derived from input, timestamped,"
    echo "                         and saved with this extension in ./files/res/"
    exit 1
fi

INPUT_FILE="$1"
OUTPUT_TYPE="$2"          # e.g., png, jpg, pdf
TARGET_DIR="./files/$OUTPUT_TYPE/res"  # Target directory

# Check if input file exists
if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file '$INPUT_FILE' not found."
    exit 1
fi

# Check if output file type is provided
if [ -z "$OUTPUT_TYPE" ]; then
    echo "Error: Output file type cannot be empty."
    exit 1
fi

# Simple sanitization: remove leading dot if present from extension
OUTPUT_TYPE=$(echo "$OUTPUT_TYPE" | sed 's/^\.//')
if [ -z "$OUTPUT_TYPE" ]; then # Check again after removing dot
    echo "Error: Output file type cannot be empty or just '.'."
    exit 1
fi


# Create target directory if it doesn't exist
mkdir -p "$TARGET_DIR"
if [ $? -ne 0 ]; then
    echo "Error: Could not create directory '$TARGET_DIR'."
    exit 1
fi

# --- Construct timestamped filename based on input name and output type ---
TIMESTAMP=$(date +"%Y%m%d%H%M%S") # Format: YYYYMMDDHHMMSS
INPUT_BASENAME=$(basename "$INPUT_FILE") # Get input filename (e.g., image.txt)

# Remove .txt extension if present, otherwise use full name as base
if [[ "$INPUT_BASENAME" == *.txt ]]; then
  OUTPUT_FILENAME_NOEXT="${INPUT_BASENAME%.txt}"
else
  OUTPUT_FILENAME_NOEXT="$INPUT_BASENAME"
fi

# Assemble final name: name_timestamp.type
FINAL_FILENAME="${OUTPUT_FILENAME_NOEXT}_${TIMESTAMP}.${OUTPUT_TYPE}"
OUTPUT_FILE="${TARGET_DIR}/${FINAL_FILENAME}" # Full path
# --- End filename construction ---


# Determine base64 decode command (handle Linux vs macOS)
BASE64_DECODE_CMD="base64 --decode"
if [[ "$(uname)" == "Darwin" ]]; then
     if command -v gbase64 &> /dev/null; then # Check for GNU base64 on macOS
        BASE64_DECODE_CMD="gbase64 --decode"
     else # Use macOS default base64
        if base64 -D < /dev/null &> /dev/null; then # Check if -D option exists
            BASE64_DECODE_CMD="base64 -D"
        elif base64 -d < /dev/null &> /dev/null; then # Check if -d option exists
             BASE64_DECODE_CMD="base64 -d"
        else # Fail if no decode option found
             echo "Error: Cannot determine correct decode option (-D or -d) for default base64 on macOS. Try installing 'coreutils' (brew install coreutils) for gbase64."
             exit 1
        fi
     fi
elif ! command -v base64 &> /dev/null; then # Check if base64 command exists elsewhere
     echo "Error: 'base64' command not found."
     exit 1
fi


echo "Decoding Base64 from '$INPUT_FILE' to '$OUTPUT_FILE'..."

# Perform decoding
$BASE64_DECODE_CMD < "$INPUT_FILE" > "$OUTPUT_FILE"
DECODE_EXIT_CODE=$?

# Check for decoding errors
if [ $DECODE_EXIT_CODE -ne 0 ]; then
    # Clean up potentially incomplete output file on error
    rm -f "$OUTPUT_FILE"
    echo "Error during Base64 decoding."
    exit $DECODE_EXIT_CODE
fi

echo "Decoding complete. Output saved to '$OUTPUT_FILE'."
exit 0