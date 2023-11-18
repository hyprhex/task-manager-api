# Task manager API 

Creating powerful and efficient REST API using Go + Gin framework

## What you can do 

- Get all tasks 
- Get a specific task 
- Update a specific task 
- Delete a specific task 
- Create a new task 

## API Endpoint 

- `GET /tasks`: Get a list of all tasks.
- `GET /tasks/id`: Get the details of a specific task.
- `PUT /tasks/id`: Update a specific task. 
- `DELETE /tasks/id`: Delete a specific task.
- `POST /tasks`: Create a new task.

## Task Parameter

Every Task has: `Title | Description | DueDate | Status`

## Usage 

Make sure that you have Go installed: `Go`
1. Clone the repo: `git clone https://github.com/AIpill/task-manager-api`
2. Go inside: `cd task-manager-api`
3. Install dependencies: `go get`
4. Run the server: `go run main.go`
5. Check: `http://0.0.0.0:8080`

