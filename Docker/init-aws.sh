#!/usr/bin/env bash 

#curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscliv2.zip"
#unzip awscliv2.zip
#./aws/install

echo "Listing existing dynamodb tables"
python --version
aws --version
aws --region us-east-1  dynamodb list-tables --endpoint-url http://localhost:4566


