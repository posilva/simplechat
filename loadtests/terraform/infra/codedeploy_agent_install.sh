#!/bin/bash
sudo yum update -y 
sudo yum install ruby -y 
sudo yum install wget -y 

CODEDEPLOY_BIN="/opt/codedeploy-agent/bin/codedeploy-agent"
$CODEDEPLOY_BIN stop
yum erase codedeploy-agent -y

cd /home/ec2-user

wget https://aws-codedeploy-${aws_region}.s3.amazonaws.com/latest/install
chmod +x ./install

sudo ./install auto