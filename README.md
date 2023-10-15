# Gohtwind Full-stack Framework 

## Introduction
### work in progress, feedback welcomed
Gohtwind is an opinionated and lightweight full-stack framework designed for rapid web application development using Go, TailwindCSS, and htmx. Streamline your development process, from backend to frontend, with Gohtwind's fast and simple approach!

## Features

- **Backend**: Robust Go backend setup with a focus on performance and simplicity 
- **Frontend**: Integrated with TailwindCSS for utility-first CSS styling and htmx for efficient frontend enhancements
- **Live Reloading**: Includes tooling for live-reloading during development
- **Production Ready**: Comes with Docker configurations tailored for production deployments

## Tech Stack
### Backend:
- Language: Go (Golang)
- Routing: Standard net/http library
- Templating: Standard html/template library
- Authentication: Auth0
- Authorization: Casbin
- Database Interaction: sqlx
### Frontend:
- Dynamic Behavior: htmx
- Styling: Tailwind CSS
### Development Environment(Suggested):
- IDE: GoLand by JetBrains
Version Control: GitHub

### CI/CD:
- Automation: GitHub Actions
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

1. Ensure you have Go and Nodejs(for tailwind, frontend tooling deps) installed on your machine.
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
gohtwind -name your_project_name
```

2. Navigate to your project directory:

```bash
cd your_project_name
```

3. Start developing your application!

## Directory Structure
```

/myapp

|-- dev-setup-<linux | macos | windows>.sh

|-- Dockerfile.prod

|-- example.env

|-- .gitignore

  

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

|   |-- postcss.config.js

|   |-- package.json

|   |-- yarn-lock.json

  

|-- author-books/  # Example feature module

|   |-- handler.go

|   |-- repository.go

|   |-- view.go

|   |-- models.go

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
gohtwind -gen-feature [feature_name]
```
* This will create a new feature module with the name `feature_name` in the root of your project directory. 
* A feature is a page within the web application. 
* The feature script generates boilerplate code for basic CRUD operations.
* All the CRUD operations are done within the context of a single page. 
  * Traditionally, CRUD operations are done across multiple pages.
  * Gohtwind's approach is to keep all the CRUD operations within a single page using dialog modals. 
  * This approach is more efficient and provides a better user experience (IMHO).
2. To start the development server run the script:
```bash
./dev-setup-<linux | macos | windows>.sh
```

## Contributing

Contributions to Gohtwind are welcome. Please read our [Contributing Guide](<link-to-contributing-guide>) for more information.

## License

Gohtwind is licensed under the [LGPL License](<link-to-license-file>).

---

Feel free to modify or expand any sections to better fit the specifics of your project.
