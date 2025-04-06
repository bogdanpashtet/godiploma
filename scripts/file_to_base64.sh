#!/bin/bash

if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <input_file>"
    echo "  input_file: Путь к файлу для кодирования."
    echo "  Output will be saved to ./files/byte/<input_filename_without_ext>.txt"
    exit 1
fi

INPUT_FILE="$1"
TARGET_DIR="./files/byte"

if [ ! -f "$INPUT_FILE" ]; then
    echo "Error: Input file '$INPUT_FILE' not found."
    exit 1
fi

mkdir -p "$TARGET_DIR"
if [ $? -ne 0 ]; then
    echo "Error: Could not create directory '$TARGET_DIR'."
    exit 1
fi

BASENAME=$(basename "$INPUT_FILE")
FILENAME="${BASENAME%.*}"
OUTPUT_FILE="${TARGET_DIR}/${FILENAME}.txt"

# Определение команды base64 (как в прошлый раз)
BASE64_CMD="base64"
if [[ "$(uname)" == "Linux" ]]; then
    BASE64_CMD="base64 -w 0"
elif [[ "$(uname)" == "Darwin" ]]; then
     if command -v gbase64 &> /dev/null; then
        BASE64_CMD="gbase64 -w 0"
     fi
fi

echo "Encoding '$INPUT_FILE' to '$OUTPUT_FILE'..."

$BASE64_CMD < "$INPUT_FILE" > "$OUTPUT_FILE"
ENCODE_EXIT_CODE=$?

if [ $ENCODE_EXIT_CODE -ne 0 ]; then
    echo "Error during Base64 encoding."
    rm -f "$OUTPUT_FILE"
    exit $ENCODE_EXIT_CODE
fi

echo "Encoding complete. Output saved to '$OUTPUT_FILE'."
exit 0