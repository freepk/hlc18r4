#!/bin/bash
APP=stor.highloadcup.ru/accounts/cute_barracuda
echo "Building $APP"
docker build -t $APP .
docker push $APP
