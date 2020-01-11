# go-todo-server

This repository contains an example of Rest API for simple Todo managing written on Go.

## General

The main goal and usage are described [here](https://github.com/VS-work/go-todo-client).

## Important notes regarding install

* The latest [Go](https://golang.org/) language version should be installed.
* We use [Go dependency management tool](https://github.com/golang/dep). So, it should be installed also.
* Use `git clone git@github.com:VS-work/go-todo-server.git` to get the repository.
* Use `cd go-todo-server` and `dep ensure` to install all expected dependencies of the project. 
* This application works with SQLite database

## Test

This project contains a minimum of required API tests. You can see these tests [here](https://github.com/VS-work/go-todo-server/blob/master/main_test.go).
Run `go test` or `go test -v` if you want to get the detailed information regarding tests passing.

## Build

Run `go build`. You will see `go-todo-server` binary file in the root directory of the project. Important: The file above does NOT include to the repository and it will be ignored on future commits. More info [here](https://github.com/VS-work/go-todo-server/blob/master/.gitignore).

## Run locally

You just need to run `go-todo-server`. It has the following parameters:

1. Path to SQLite database file which will be used by the application.
2. A host that will be used as `Access-Control-Allow-Origin` header value (CORS requirement). You can see related point int the following part [of code](https://github.com/VS-work/go-todo-server/blob/master/app.go#L211)

For example, the following command 

`./go-todo-server ./db/todos_test.db http://localhost:3000` 

will start API server and allow to get requests from `http://localhost:3000` and will use `./db/todos_test.db` file as the current database.

Important note: this application use port `3001` on local use.

## Deploy

There is only [Heroku](https://heroku.com) deployment supporting.

There are following features regarding deployment:

1. Please follow [instructions](https://devcenter.heroku.com/articles/getting-started-with-go) to deploy the application on `Heroku`
2. Main settings regarding `Heroku` are placed [here](https://github.com/VS-work/go-todo-server/blob/master/Procfile). This file contains command-runner that includes the database file name and allowed origin (see `Run locally` above)
3. There is an [empty DB file](https://github.com/VS-work/go-todo-server/blob/master/todos.db) (SQLite) that used as a database on Heroku's server.
4. Only one frontend is allowed to works with the API: [https://vs-work.github.io](https://vs-work.github.io).
5. If you want to use the API on Heroku go to the following URL: [https://dry-woodland-14649.herokuapp.com/](https://dry-woodland-14649.herokuapp.com/) 

## Sendgrid integration

This application has integration with [Sendgrid cloud email service](https://sendgrid.com) and it sends an email after todo creation, modification, and removing. Please set `SENDGRID_API_KEY` environment variable on your local PC or Heroku application if you want to activate this feature.