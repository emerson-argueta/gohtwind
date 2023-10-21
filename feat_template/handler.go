package {{FEATURE_NAME}}

import (
    "{{PROJECT_NAME}}/.gen/test_db/model"
    . "{{PROJECT_NAME}}/.gen/test_db/table"
    "database/sql"
    . "github.com/go-jet/jet/v2/mysql"
    "log"
    "net/http"
)

// List displays the list of all items
func List(db *sql.DB) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: Fetch items from the database or service, then put them in 'items'
        items := []string{"test0", "test1", "test2"} // This is a placeholder. Replace 'Item' with your actual data structure.

        // Render the list view with the fetched items
        renderTemplate(w, "list.html", map[string]interface{}{
            "Items": items, // Pass the items as data to the template
        })
    })
}

// Create handles the creation of a new item
func Create(db *sql.DB) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            // TODO: Handle item creation logic

            // Redirect to the list view or display a success message
            http.Redirect(w, r, "/{{FEATURE_NAME}}/list", http.StatusSeeOther)
            return
        }
    })
}

// Read displays details of a specific item
func Read(db *sql.DB) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // TODO: Fetch item details based on an identifier from 'r', then put it in 'item'
        item := "test" // This is a placeholder. Replace 'Item' with your actual data structure.

        // Render the read view with the fetched item details
        renderPartialTemplate(w, "read.html", map[string]interface{}{
            "Item": item, // Pass the item as data to the template
        })
    })
}

// Update handles updating an existing item
func Update(db *sql.DB) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPatch {
            // TODO: Handle item update logic

            // Redirect to the list view or display a success message
            http.Redirect(w, r, "/{{FEATURE_NAME}}/list", http.StatusSeeOther)
            return
        }
    })
}

// Delete handles deleting an item
func Delete(db *sql.DB) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method == http.MethodPost {
            // TODO: Handle item deletion logic

            // Redirect to the list view or display a success message
            http.Redirect(w, r, "/{{FEATURE_NAME}}/list", http.StatusSeeOther)
            return
        }
    })
}
