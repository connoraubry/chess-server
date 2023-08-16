#!/bin/bash


wget https://dev.mysql.com/get/mysql-apt-config_0.8.26-1_all.deb
sudo dpkg -i mysql-apt-config_0.8.26-1_all.deb
rm mysql-apt-config_0.8.26-1_all.deb

sudo apt-get update

sudo apt-get install mysql-server
