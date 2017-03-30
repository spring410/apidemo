
-- add user name test to db
grant all privileges on *.* to test@localhost Identified by "mysql";

-- only for test	
-- drop user test1@localhost;


show databases;

-- create db name: Accounts
create database if not exists Accounts;

use Accounts;

show tables;

-- only for test
-- drop database Accountss;

-- show databases;

-- end.

