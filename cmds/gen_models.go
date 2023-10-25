package cmds

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func genModelsUsageString() string {
	return `
Usage: gohtwind gen-models [options]
    Options:
		-adapter string
			Database adapter (mysql, postgres)
		-dsn string
			Database connection string
			postgres ex: <username>:<password>@tcp(<host>:<port>)/<dbname>
			mysql ex: <username>:<password>@tcp(<host>:<port>)/<dbname
		-schema string
			Database schema (postgres adapter only)
`
}

func GenModels() {
	genModelsFlags := flag.NewFlagSet("models", flag.ExitOnError)
	modelsAdapter := genModelsFlags.String("adapter", "", "Database adapter (mysql, postgres)")
	u := `Database connection string
			postgres ex: <username>:<password>@tcp(<host>:<port>)/<dbname>
			mysql ex: <username>:<password>@tcp(<host>:<port>)/<dbname`
	modelsDsn := genModelsFlags.String("dsn", "", u)
	modelsSchema := genModelsFlags.String("schema", "", "Database schema (postgres adapter only)")
	args := os.Args[2:]
	genModelsFlags.Parse(args)
	if *modelsAdapter == "" || *modelsDsn == "" {
		fmt.Println(genModelsUsageString())
		os.Exit(1)
	}
	if *modelsAdapter == "postgres" && *modelsSchema == "" {
		fmt.Println(genModelsUsageString())
		os.Exit(1)
	}
	dsnArg := fmt.Sprintf("-dsn=%s", *modelsDsn)
	adapterArg := fmt.Sprintf("-adapter=%s", *modelsAdapter)
	var schemaArg string
	if *modelsSchema == "" {
		schemaArg = fmt.Sprintf("-schema=%s", *modelsSchema)
		exec.Command("./bin/jet", "-path=./.gen", dsnArg, schemaArg, adapterArg).Run()
	} else {
		exec.Command("./bin/jet", "-path=./.gen", dsnArg, adapterArg).Run()
	}
}
