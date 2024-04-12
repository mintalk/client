# Mintalk

## Client

### Running

~~~
cd client
go run . <host> <username> <password>
~~~


## Server

### Running

~~~
cd server
go run .
~~~

### Configuration

The configuration file is located at `server/config.json`.

Example configuration:

~~~
database: root:2207@tcp(127.0.0.1:3306)/mintalk?charset=utf8mb4&parseTime=True&loc=Local
host: localhost:8000
session_lifetime: 1440
~~~

### Commands

When running the server, you can input some commands.

[List of commands](server/COMMANDS.md)
