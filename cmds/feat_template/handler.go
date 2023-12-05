package {{FEATURE_NAME}}

import (
    "{{PROJECT_NAME}}/infra"
    "database/sql"
    "net/http"
    "fmt"
)

type Handle struct {
    dbs map[string]*sql.DB
}

// List displays the list of all items
func (h *Handle) List(w http.ResponseWriter, r *http.Request) {
    // TODO: Fetch items from the database or service, then put them in 'items'
    items := []string{"test0", "test1", "test2"} // This is a placeholder. Replace 'Item' with your actual data structure.

    // Render the list view with the fetched items
    vt := infra.ViewTemplate{
        BasePath: "templates",
        Path:     "{{FEATURE_NAME}}/templates/list.html",
        PartialPaths: []string{
            "{{FEATURE_NAME}}/templates/_read.html",
            "{{FEATURE_NAME}}/templates/_update.html",
            "{{FEATURE_NAME}}/templates/_delete.html",
            "{{FEATURE_NAME}}/templates/_create.html",
        },
    }
    renderTemplate(w, &vt, map[string]interface{}{
        "Items": items, // Pass the items as data to the template
        "Model": struct{Test string}{}, // Pass the model dto here
    })
}

// Create handles the creation of a new item
func (h *Handle) Create(w http.ResponseWriter, r *http.Request) {
    // TODO: Handle item creation logic

    // Redirect to the list view or display a success message
    http.Redirect(w, r, "/{{FEATURE_NAME}}", http.StatusSeeOther)
    return
}

// Read displays details of a specific item
func (h *Handle) Read(w http.ResponseWriter, r *http.Request) {
    idStr := r.Context().Value("id").(string)
    fmt.Println(idStr)
    // TODO: Fetch item details based on an identifier from 'r', then put it in 'item'
    item := "test" // This is a placeholder. Replace 'Item' with your actual data structure.

    vt := infra.ViewTemplate{
        BasePath: "templates",
        Path:     "{{FEATURE_NAME}}/templates/read.html",
    }
    renderTemplate(w, &vt, map[string]interface{}{
        "Item": item, // Pass the item as data to the template
        "Model": struct{Test string}{}, // Pass the model dto here
    })
}

// Update handles updating an existing item
func (h *Handle) Update(w http.ResponseWriter, r *http.Request) {
    idStr := r.Context().Value("id").(string)
    fmt.Println(idStr)
    // TODO: Handle item update logic

    // Redirect to the list view or display a success message
    http.Redirect(w, r, "/{{FEATURE_NAME}}", http.StatusSeeOther)
    return
}

// Delete handles deleting an item
func (h *Handle) Delete(w http.ResponseWriter, r *http.Request) {
    idStr := r.Context().Value("id").(string)
    fmt.Println(idStr)
    // TODO: Handle item deletion logic

    // Redirect to the list view or display a success message
    http.Redirect(w, r, "/{{FEATURE_NAME}}", http.StatusSeeOther)
    return
}
