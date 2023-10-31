
func Create{{MODEL_NAME}}Repo(dbs map[string]*sql.DB, m model.{{MODEL_NAME}}) error {
	stmt := {{MODEL_NAME}}.INSERT().VALUES(m)
	log.Println(stmt.DebugSql())
	_, err := stmt.Exec(dbs["{{DB_NAME}}"])
	return err
}

func All{{MODEL_NAME}}Repo(dbs map[string]*sql.DB, pg *infra.Pagination) ([]struct{ model.{{MODEL_NAME}} }, *infra.Pagination) {
	if pg == nil {
		stmt := SELECT(COUNT(STAR).AS("count")).FROM({{MODEL_NAME}})
		var dest struct{ Count int64 }
		log.Println(stmt.DebugSql())
		err := stmt.Query(dbs["{{DB_NAME}}"], &dest)
		if err != nil {
			log.Println(err)
		}
		pg = infra.NewPagination(1, 10, dest.Count)
	}
	offset := (pg.CurrentPage - 1) * pg.PerPage
	stmt := SELECT({{MODEL_NAME}}.AllColumns).FROM({{MODEL_NAME}}).LIMIT(pg.PerPage).OFFSET(offset)
	var dest []struct {
		model.{{MODEL_NAME}}
	}
	log.Println(stmt.DebugSql())
	err := stmt.Query(dbs["{{DB_NAME}}"], &dest)
	if err != nil {
		log.Println(err)
	}
	return dest, pg
}

func Fetch{{MODEL_NAME}}Repo(dbs map[string]*sql.DB, id int64) model.{{MODEL_NAME}} {
	stmt := SELECT({{MODEL_NAME}}.AllColumns).FROM({{MODEL_NAME}}).WHERE({{MODEL_NAME}}.ID.EQ(Int(id)))
	var dest model.{{MODEL_NAME}}
	log.Println(stmt.DebugSql())
	err := stmt.Query(dbs["{{DB_NAME}}"], &dest)
	if err != nil {
		log.Println(err)
	}
	return dest
}

func Update{{MODEL_NAME}}Repo(dbs map[string]*sql.DB, m model.{{MODEL_NAME}}) error {
	stmt := {{MODEL_NAME}}.
		UPDATE({{MODEL_NAME}}.AllColumns).
		MODEL(m).
		WHERE({{MODEL_NAME}}.ID.EQ(Int(m.ID)))
	log.Println(stmt.DebugSql())
	_, err := stmt.Exec(dbs["{{DB_NAME}}"])
	return err
}

func Delete{{MODEL_NAME}}Repo(dbs map[string]*sql.DB, id int64) error {
	stmt := {{MODEL_NAME}}.
		DELETE().
		WHERE({{MODEL_NAME}}.ID.EQ(Int(id)))
	log.Println(stmt.DebugSql())
	_, err := stmt.Exec(dbs["{{DB_NAME}}"])
	return err
}
