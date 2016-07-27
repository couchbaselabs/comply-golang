# Couchbase Compliance Demo with GoLang

This project is meant to demonstrate a full stack application using GoLang, Angular 2, and Couchbase Server.  This project shows separation of the frontend and backend code as well as highly relational, yet schema-less NoSQL document data.

## The Requirements

There are a few requirements that must be met before trying to run this project:

* The Angular 2 CLI
* GoLang 1.5+
* Couchbase Server 4.1+

The Angular CLI is necessary to build the Angular 2 frontend that is provided with this project.  GoLang is required for building the backend.  Couchbase Server is required as the database layer.

## Configuring the Project

There are three things that must be done to run this project:

1. The Angular 2 frontend must be built
2. The GoLang project must be compiled
3. The database must be configured

The Angular 2 code is not dependent on GoLang or the database.  To build the Angular 2 project, execute the following from the **angular** directory of the project:

```sh
npm install
ng build
```

The above commands will install the Angular 2 dependencies and build the **angular/src** to the **angular/dist** directory.  Copy this directory to the parent directory like the following:

```sh
cp -r angular/dist public
```

Notice that the directory must be renamed to **public** at the root of the project.

Before compiling the GoLang code, the Couchbase Server information must be specified first.  This project depends on a **comply** bucket to be created.  Because this project makes use of N1QL, the bucket must have at least one index.

To build the GoLang project, execute the following from your Command Prompt (Windows) or Terminal (Linux and Mac):

```sh
go build
```

## Running the Project

Provided that the Angular 2 code has been built and you've configured your database information, execute the following:

```sh
./comply-golang
```

From the web browser, navigate to **http://localhost:3000** to see the application in action.

## Resources

Couchbase - [http://www.couchbase.com](http://www.couchbase.com)

Couchbase Compliance Demo with Node.js - [https://github.com/couchbaselabs/comply-nodejs](https://github.com/couchbaselabs/comply-nodejs)

Couchbase Compliance Demo with Java - [https://github.com/couchbaselabs/comply-java](https://github.com/couchbaselabs/comply-java)
