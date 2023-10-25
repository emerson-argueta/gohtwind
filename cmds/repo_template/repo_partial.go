
// CRUD operation

func All{{MODEL_NAME}}Repo(dbs map[string]*sql.DB) []struct{ model.{{MODEL_NAME}} } {
	stmt := SELECT({{MODEL_NAME}}.AllColumns).FROM({{MODEL_NAME}}).LIMIT(10)
	var dest []struct { model.{{MODEL_NAME}} }
	err := stmt.Query(dbs["{{DB_NAME}}"], &dest)
	if err != nil {
		log.Fatal(err)
	}
	return dest
}
