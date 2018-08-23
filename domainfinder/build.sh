#!/bin/bash

echo "build domainfinder..."
go build -o domainfinder

echo "build synonyms..."
cd ../synonyms
go build -o ../domainfinder/lib/synonyms

echo "build available..."
cd ../available
go build -o ../domainfinder/lib/available

echo "build sprinkle..."
cd ../sprinkle
go build -o ../domainfinder/lib/sprinkle

echo "build coolify..."
cd ../coolify
go build -o ../domainfinder/lib/coolify

echo "build domainify..."
cd ../domainify
go build -o ../domainfinder/lib/domainify
