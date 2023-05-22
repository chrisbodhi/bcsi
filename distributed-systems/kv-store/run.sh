#!/bin/sh

# run.sh starts the load balancer and the two servers
# It periodically checks if the servers are still running
# If they are not, it restarts any downed servers

# Start the load balancer
go run lb/lb.go &

# Start the two servers
go run server/server.go 8889 &
go run server/server.go 8890 &

# Ping the healthcheck endpoint of each server every 5 seconds
while true; do
    sleep 5
    curl localhost:8889/ping
    curl localhost:8890/ping
done
