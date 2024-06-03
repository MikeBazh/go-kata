#!/bin/bash
docker network create mynetwork
docker-compose up --force-recreate --build

