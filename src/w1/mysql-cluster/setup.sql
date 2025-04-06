CREATE USER 'icadmin'@'%' IDENTIFIED BY 'icadmin';

GRANT ALL privileges ON *.* TO 'icadmin'@'%' with grant option;

reset master;