# Gohtwind Full-stack Framework

## Introduction
### work in progress, feedback welcomed
Gohtwind is an opinionated and lightweight full-stack framework designed for rapid web application development using Go, TailwindCSS, and htmx. Streamline your development process, from backend to frontend, with Gohtwind's fast and simple approach!

## Features

- **Backend**: Go backend setup with a focus on performance and simplicity
- **Frontend**: Integrated with TailwindCSS for utility-first CSS styling and htmx for efficient frontend enhancements
- **Live Reloading**: Includes tooling for live-reloading during development
- **Deploying**: Comes with Docker configurations tailored for production deployments

## Tech Stack
### Backend:
- Language: Go (Golang)
- Routing: Standard net/http library
- Templating: Standard html/template library
- Authentication: Authboss (coming soon)
- Authorization: Casbin (coming soon)
- Database Interaction: go-jet/jet
### Frontend:
- Dynamic Behavior: htmx
- Styling: Tailwind CSS
### Development Environment(Suggested):
- IDE: GoLand by JetBrains
- Version Control: GitHub

### CI/CD:
- Automation: GitHub Actions (coming soon)
- Automated deployment to Google Cloud Run when merging a PR from the main branch to the prod branch.
- Deployment Setup: Docker (with a production-specific Dockerfile)
- Deployment & Hosting:
  - Containerization: Docker
  - Hosting: Google Cloud Run
  - Database:
    - Database Service: Google Cloud SQL
    - Database Engine: MySQL

## Installation

To get started with Gohtwind:

1. Ensure you have Go installed on your machine.
* Optional: Ensure you have Docker installed on your machine. (For development database)
2. Clone/download the Gohtwind repository.
3. Navigate to the Gohtwind directory and run:

```bash
go build
sudo cp gohtwind /usr/local/bin/
```

Now, you can use the `gohtwind` command from anywhere in your terminal!

## Quick Start
1. Create a new Gohtwind project:
```bash
# gohtwind name your_project_name
# ex:
gohtwind new town
```
2. Navigate to your project directory:
```bash
# cd your_project_name 
# ex:
cd town
```
3. Create a page (a.k.a. feature) within your project:
```bash
# gohtwind gen-feature [feature_name]
# ex:
gohtwind gen-feature market 
```
4. Add feature routes to the `main.go` in the root of your project directory:
```go
package main

// the ide will automatically import the packages for you
// if not, you can manually import them like so:
import(
  // ...
  // "<project_name>/infra"
  "town/infra"
  // "<project_name>/<feature_name>"
  "town/market"
  // ...
)

func main() {
  // ...
  http.Handle("/static/", infra.LoggingMiddleware(http.StripPrefix("/static/", http.FileServer(http.Dir("./frontend/static/")))))
  /**
    Replace feature_name with the name of your feature like so:
    feature_name_1.SetupRoutes(dbs, infra.LoggingMiddleware)
            ...
    feature_name_2.SetupRoutes(dbs, infra.LoggingMiddleware)
            ...
    feature_name_n.SetupRoutes(dbs, infra.LoggingMiddleware)
     **/
  // ex: 
  market.SetupRoutes(dbs, infra.LoggingMiddleware)

  log.Printf("Server started on :%s\n", port)
  err = http.ListenAndServe(":"+port, nil)
  if err != nil {
    log.Fatal(err)
  }
}
````
5. Optional: Run containerized database for development:
```bash
docker build -t gohtwind-db -f Dockerfile.db .
docker run -d -p 3306:3306 --name gohtwind-db gohtwind-db
```
5. Generate sql models using:
```bash
# gohtwind gen-models -adapter=mysql -dsn="<username>:<password>@tcp(<host>:<port>)/<dbname>"
# or
# gohtwind gen-models -adapter=postgres -dsn="postgresql://<user>:<password>@<host>:<port>/<dbname>?sslmode=disable -schema=<schema>"
# using the containerized database
gohtwind gen-models -adapter=mysql -dsn="root:root@tcp(localhost:3306)/dev"
```
6. Generate a repository file for the feature:
```bash
# Make sure that the model_name is the same as the generated model name (Usually the table name in TitleCase) 
# gohtwind gen-repository -feature-name=<feature_name> -model-name=<model_name> -db-name=<dbname> -adapter=<mysql | postgres> -schema=<schema postgres only>
# using the model from containerized database
gohtwind gen-repository -feature-name=market -model-name=Products -db-name=dev -adapter=mysql
```
7. Copy the example.env file and rename it to .env:
```bash
cp example.env .env
```
8. Start the development server:
```bash
./dev-run.sh
```
9. Start developing your application!
## Directory Structure
```

/myapp

|-- dev-run.sh

|-- Dockerfile.prod

|-- example.env

|-- .gitignore

|-- .gen/ # Generated sql models

|-- config/ # Configuration files 
  

|-- go.mod

|-- go.sum

  
|-- templates/ # base template and shared templates

|   |-- shared/

|   |-- base.html 

|-- frontend/

|   |-- static/

|   |   |-- css/

|   |   |   |-- main.css  # Base CSS file for Tailwind

|   |

|   |-- output.css  # Generated CSS after processing with Tailwind

|   |-- tailwind.config.js

  

|-- author-books/  # Example feature module

|   |-- handler.go

|   |-- repository.go

|   |-- view.go

|   |-- routes.go

|   |-- static/

|   |   |-- js/

|   |   |-- css/

|   |-- templates/

|       |-- create.html

|       |-- read.html

|       |-- update.html

|       |-- delete.html

|       |-- list.html

  

|-- other-feature/  # Another example feature module

|   ...  # Similar structure as above

  

|-- ...  # Other project-wide files, utilities, shared components, etc.

```



## Utility Scripts

1. To generates a new feature within your project run the gohtwind command with the gen-feature flag:
```bash
gohtwind gen-feature [feature_name]
```
* This will create a new feature module with the name `feature_name` in the root of your project directory.
* A feature is a page within the web application.
* The feature script generates boilerplate code for basic CRUD operations.
* All the CRUD operations are done within the context of a single page.
  * Traditionally, CRUD operations are done across multiple pages.
  * Gohtwind's approach is to keep all the CRUD operations within a single page using dialog modals.
  * This approach is more efficient and provides a better user experience (IMO).
2. To generate sql models, use the following gohtwind command:
```bash
gohtwind gen-models -adapter=<mysql | postgres> -dsn=<dsn> -schema=<schema>
```
* It wraps the go-jet/jet run the jet command
* Generated models are placed in the `.gen` directory at the root of your project directory.
* The `-adapter` flag specifies the database source. Currently, only MySQL and Postgres is supported.
* The `-dsn` flag specifies the database connection string.
* The `-schema` flag specifies the database schema to generate models for. (Only applicable for Postgres)
3. To generate repository boilerplate code, use the following gohtwind command:
```bash
gohtwind gen-repository -feature-name=<feature_name> -model-name=<model_name> -db-name=<dbname> -adapter=<mysql | postgres> -schema=<schema postgres only>
```
* This will generate a repository file for the specified feature and model.
* The repository file contains boilerplate code for basic CRUD operations.
* The repository file is used by the feature's handler to interact with the database.
* The `-feature-name` flag specifies the name of the feature the repository is for.
* The `-model-name` flag specifies the name of the model (sql table) the repository is for.
* The `-db-name` flag specifies the name of the database
* The `-schema` schema (postgres only) the model is in.
* The `-adapter` flag specifies the database adapter to use. Currently, only MySQL and Postgres is supported.

4. To generate a form inside a template, use the following gohtwind command:
```bash
gohtwind gen-form -feature-name=<feature_name> -model-name=<model_name>
```
* This command replaces {{GEN_FORM}} in the specified template with a form for the specified model
* The form is generated using the model's form tags
* When an instance name is provided, the form is generated with the instance's values
* The `-feature-name` flag specifies the name of the feature the form is for.
* The `-model-name` flag specifies the name of the model the form is for.
* The  `-template-name` flag specifies the name of the template the form is for.
* The  `-instance-name` flag specifies the name of the instance the form is for. Use this flag for update forms. Omit this flag for create forms.
* The `-action` flag specifies the action of the form.

3. To start the development server run the script:
```bash
./dev-run.sh
```

## Contributing

Contributions to Gohtwind are welcome. Please read our [Contributing Guide](<link-to-contributing-guide>) (coming soon) for more information.

## License

Gohtwind is licensed under the [MIT License](<link-to-license-file>).

---

Feel free to modify or expand any sections to better fit the specifics of your project.
