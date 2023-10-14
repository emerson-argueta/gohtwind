#!/bin/bash

# Check if feature name is provided
if [ -z "$1" ]; then
    echo "Please provide a feature name."
    exit 1
fi

FEATURE_NAME="$1"

# Create the feature directory and sub-directories
mkdir $FEATURE_NAME
mkdir "$FEATURE_NAME/static"
mkdir "$FEATURE_NAME/static/js"
mkdir "$FEATURE_NAME/static/css"
mkdir "$FEATURE_NAME/templates"

# Create basic template files for the feature (optional)
touch "$FEATURE_NAME/templates/create.html"
touch "$FEATURE_NAME/templates/read.html"
touch "$FEATURE_NAME/templates/update.html"
touch "$FEATURE_NAME/templates/delete.html"
touch "$FEATURE_NAME/templates/list.html"

# Generate the handler.go file with base content
cat > "$FEATURE_NAME/handler.go" <<EOL
package $FEATURE_NAME

import (
    "net/http"
)

// List displays the list of all items
func List(w http.ResponseWriter, r *http.Request) {
    // TODO: Fetch items and render the list template
}

// Create handles the creation of a new item
func Create(w http.ResponseWriter, r *http.Request) {
    // TODO: Handle item creation and render the create template
}

// Read displays details of a specific item
func Read(w http.ResponseWriter, r *http.Request) {
    // TODO: Fetch item details and render the read template
}

// Update handles updating an existing item
func Update(w http.ResponseWriter, r *http.Request) {
    // TODO: Handle item update and render the update template
}

// Delete handles deleting an item
func Delete(w http.ResponseWriter, r *http.Request) {
    // TODO: Handle item deletion
}
EOL

# Generate the routes.go file with base content
cat > "$FEATURE_NAME/routes.go" <<EOL
package $FEATURE_NAME

import (
    "net/http"
)

func SetupRoutes() {
    http.HandleFunc("/$FEATURE_NAME/", List)
    http.HandleFunc("/$FEATURE_NAME/create", Create)
    http.HandleFunc("/$FEATURE_NAME/read", Read)
    http.HandleFunc("/$FEATURE_NAME/update", Update)
    http.HandleFunc("/$FEATURE_NAME/delete", Delete)

    // Serve static files for the $FEATURE_NAME feature
    http.Handle("/static/$FEATURE_NAME/", http.StripPrefix("/static/$FEATURE_NAME/", http.FileServer(http.Dir("./$FEATURE_NAME/static/"))))
}
EOL

# Generate the view.go file with base content
cat > "$FEATURE_NAME/view.go" <<EOL
package $FEATURE_NAME

import (
    "html/template"
    "net/http"
)

var templates = template.Must(template.ParseGlob("$FEATURE_NAME/templates/*.html"))

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    err := templates.ExecuteTemplate(w, tmpl, data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}
EOL

echo "Feature '$FEATURE_NAME' has been generated!"