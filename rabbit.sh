#!/bin/bash

start_or_run () {
    docker inspect battleship_rabbitmq > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo "Starting Battleship RabbitMQ container..."
        docker start battleship_rabbitmq
    else
        echo "Battleship RabbitMQ container not found, creating a new one..."
        docker run -d --name battleship_rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.13-management
    fi
}

case "$1" in
    start)
        start_or_run
        ;;
    stop)
        echo "Stopping Battleship RabbitMQ container..."
        docker stop battleship_rabbitmq
        ;;
    logs)
        echo "Fetching logs for Battleship RabbitMQ container..."
        docker logs -f battleship_rabbitmq
        ;;
    *)
        echo "Usage: $0 {start|stop|logs}"
        exit 1
esac
