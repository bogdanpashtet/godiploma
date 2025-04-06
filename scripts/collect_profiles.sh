#!/bin/bash

SERVICE_NAME=$1
DATE=$(date +"%Y%m")
DIR="${SERVICE_NAME}-${DATE}"

mkdir -p "$DIR"

PORT=8081
BASE_URL="http://localhost:$PORT/debug/pprof"

# Fetch profiling
curl "${BASE_URL}/profile?seconds=120" -o "${DIR}/cpu.pprof" &
curl "${BASE_URL}/trace?seconds=30" -o "${DIR}/trace.trace" &
curl "${BASE_URL}/heap" -o "${DIR}/heap.pprof" &
curl "${BASE_URL}/goroutine" -o "${DIR}/goroutine.pprof" &
curl "${BASE_URL}/block" -o "${DIR}/block.pprof" &
curl "${BASE_URL}/mutex" -o "${DIR}/mutex.pprof" &
curl "${BASE_URL}/allocs" -o "${DIR}/allocs.pprof" &
curl "${BASE_URL}/threadcreate" -o "${DIR}/threadcreate.pprof" &

wait

# Start profiling servers
go tool pprof -http=:8080 "${DIR}/cpu.pprof" &
go tool trace "${DIR}/trace.trace" &
go tool pprof -http=:8081 "${DIR}/heap.pprof" &
go tool pprof -http=:8082 "${DIR}/goroutine.pprof" &
go tool pprof -http=:8083 "${DIR}/block.pprof" &
go tool pprof -http=:8084 "${DIR}/mutex.pprof" &
go tool pprof -http=:8085 "${DIR}/allocs.pprof" &
go tool pprof -http=:8086 "${DIR}/threadcreate.pprof" &

sleep 2

# Open URLs in default browser with readable naming
if [[ "$OSTYPE" == "darwin"* ]]; then
  open "http://localhost:8090" # CPU Profile
  open "http://localhost:8091" # Heap Profile
  open "http://localhost:8092" # Goroutine Profile
  open "http://localhost:8093" # Block Profile
  open "http://localhost:8094" # Mutex Profile
  open "http://localhost:8095" # Allocations Profile
  open "http://localhost:8096" # Thread Create Profile
else
  xdg-open "http://localhost:8090" # CPU Profile
  xdg-open "http://localhost:8091" # Heap Profile
  xdg-open "http://localhost:8092" # Goroutine Profile
  xdg-open "http://localhost:8093" # Block Profile
  xdg-open "http://localhost:8094" # Mutex Profile
  xdg-open "http://localhost:8095" # Allocations Profile
  xdg-open "http://localhost:8096" # Thread Create Profile
fi