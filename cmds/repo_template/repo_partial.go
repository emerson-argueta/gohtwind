
func Create{{MODEL_NAME}}Repo(dbs map[string]*sql.DB, m *{{MODEL_NAME}}) error {
	v := reflect.ValueOf(m).Elem()
	field := v.FieldByName("UpdatedAt")
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(time.Now()))
	}
	v = reflect.ValueOf(m).Elem()
	field = v.FieldByName("CreatedAt")
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(time.Now()))
	}
	stmt := jet.{{MODEL_NAME}}.INSERT().VALUES(*m)
	log.Println(stmt.DebugSql())
	_, err := stmt.Exec(dbs["{{DB_NAME}}"])
	return err
}

func All{{MODEL_NAME}}Repo(dbs map[string]*sql.DB, pg *infra.Pagination) ([]struct{ {{MODEL_NAME}} }, *infra.Pagination) {
	if pg == nil {
		stmt := SELECT(COUNT(STAR).AS("count")).FROM(jet.{{MODEL_NAME}})
		var dest struct{ Count int64 }
		log.Println(stmt.DebugSql())
		err := stmt.Query(dbs["{{DB_NAME}}"], &dest)
		if err != nil {
			log.Println(err)
		}
		pg = infra.NewPagination(1, 10, dest.Count)
	}
	offset := (pg.CurrentPage - 1) * pg.PerPage
	stmt := SELECT(jet.{{MODEL_NAME}}.AllColumns).FROM(jet.{{MODEL_NAME}}).LIMIT(pg.PerPage).OFFSET(offset)
	var dest []struct {
		{{MODEL_NAME}}
	}
	log.Println(stmt.DebugSql())
	err := stmt.Query(dbs["{{DB_NAME}}"], &dest)
	if err != nil {
		log.Println(err)
	}
	return dest, pg
}

func Fetch{{MODEL_NAME}}Repo(dbs map[string]*sql.DB, id int64) {{MODEL_NAME}} {
	stmt := SELECT(jet.{{MODEL_NAME}}.AllColumns).FROM(jet.{{MODEL_NAME}}).WHERE(jet.{{MODEL_NAME}}.ID.EQ(Int(id)))
	var dest {{MODEL_NAME}}
	log.Println(stmt.DebugSql())
	err := stmt.Query(dbs["{{DB_NAME}}"], &dest)
	if err != nil {
		log.Println(err)
	}
	return dest
}

func Update{{MODEL_NAME}}Repo(dbs map[string]*sql.DB, m *{{MODEL_NAME}}) error {
	v := reflect.ValueOf(m).Elem()
	field := v.FieldByName("UpdatedAt")
	if field.IsValid() && field.CanSet() {
		field.Set(reflect.ValueOf(time.Now()))
	}
	stmt := jet.{{MODEL_NAME}}.
		UPDATE(jet.{{MODEL_NAME}}.AllColumns).
		MODEL(*m).
		WHERE(jet.{{MODEL_NAME}}.ID.EQ(Int(m.ID)))
	log.Println(stmt.DebugSql())
	_, err := stmt.Exec(dbs["{{DB_NAME}}"])
	return err
}

func Delete{{MODEL_NAME}}Repo(dbs map[string]*sql.DB, id int64) error {
	stmt := jet.{{MODEL_NAME}}.
		DELETE().
		WHERE(jet.{{MODEL_NAME}}.ID.EQ(Int(id)))
	log.Println(stmt.DebugSql())
	_, err := stmt.Exec(dbs["{{DB_NAME}}"])
	return err
}
