#!/bin/sh
cd ..
docker build --no-cache -t teaftp -f docker/teaftp/Dockerfile .
