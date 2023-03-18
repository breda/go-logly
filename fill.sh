#!/bin/bash


for I in {1..1000000}
do
	DATA=`cat /dev/urandom | base64 | head -c 50000 | tr -d "\n"`
	curl -XPOST -d "{\"data\":\"$I $DATA\"}" https://localhost:3332/append
done
