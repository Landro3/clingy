#!/bin/bash
echo "Copying files..."
mkdir client-1
mkdir client-2
cp -r clingy-client/* client-1
cp -r clingy-client/* client-2
echo "Finished copying files"

echo "Setting UI env variables"
touch client-1/ui/.env
echo "API_URL=http://localhost:8888/api" >client-1/ui/.env
touch client-2/ui/.env
echo "API_URL=http://localhost:8989/api" >client-2/ui/.env

