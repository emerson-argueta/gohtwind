
// CRUD operation

func Create{{MODEL_NAME}}Repo(db *sql.DB) *{{MODEL_NAME}}Repo {
	return &{{MODEL_NAME}}Repo{db: db}
}