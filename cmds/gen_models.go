package cmds

import (
	"bufio"
	"embed"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

//go:embed bin/jet-*
var jetBinary embed.FS

func genModelsUsageString() string {
	return `
Usage: gohtwind gen-models [options]
    Options:
		-adapter string
			Database adapter (mysql, postgres)
		-dsn string
			Database connection string
			postgres ex: "<username>:<password>@tcp(<host>:<port>)/<dbname>"
			mysql ex: "<username>:<password>@tcp(<host>:<port>)/<dbname>"
		-schema string
			Database schema (postgres adapter only)
	Info:
		* It wraps the go-jet/jet run the jet command
		* Generated models are placed in the .gen directory at the root of your project directory.

`
}

func GenModels() {
	genModelsFlags := flag.NewFlagSet("gohtwind gen-models", flag.ExitOnError)
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
	executeJetCmd(*modelsSchema, *modelsDsn, *modelsAdapter)
}

func executeJetCmd(modelsSchema string, modelsDsn string, modelsAdapter string) {
	tmpBinPath, err := createTempExecutable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create temp executable: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(tmpBinPath) // Clean up
	var cmd *exec.Cmd
	dsnArg := fmt.Sprintf("-dsn=%s", modelsDsn)
	adapterArg := fmt.Sprintf("-source=%s", modelsAdapter)
	if modelsSchema == "" {
		schemaArg := fmt.Sprintf("-schema=%s", modelsSchema)
		cmd = exec.Command(tmpBinPath, "-path=./.gen", dsnArg, schemaArg, adapterArg)
	} else {
		cmd = exec.Command(tmpBinPath, "-path=./.gen", dsnArg, adapterArg)
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create stdout pipe: %v\n", err)
		os.Exit(1)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create stderr pipe: %v\n", err)
		os.Exit(1)
	}
	if err := cmd.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start command: %v\n", err)
		os.Exit(1)
	}
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading standard output: %v\n", err)
		}
	}()
	go func() {
		errScanner := bufio.NewScanner(stderr)
		for errScanner.Scan() {
			fmt.Fprintln(os.Stderr, errScanner.Text())
		}
		if err := errScanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading standard error: %v\n", err)
		}
	}()

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Command finished with error: %v\n", err)
		os.Exit(1)
	}
}

func createTempExecutable() (string, error) {
	o := strings.ToLower(runtime.GOOS)
	a := strings.ToLower(runtime.GOARCH)
	f, err := jetBinary.Open(fmt.Sprintf("bin/jet-%s-%s", o, a))
	if err != nil {
		return "", err
	}
	defer f.Close()

	// Create a temporary file
	tmpFile, err := os.CreateTemp(os.TempDir(), "jet-")
	if err != nil {
		return "", err
	}
	defer tmpFile.Close()

	// Copy the embedded file's content to the temporary file
	_, err = io.Copy(tmpFile, f)
	if err != nil {
		return "", err
	}

	// Make the temporary file executable
	err = os.Chmod(tmpFile.Name(), 0700)
	if err != nil {
		return "", err
	}

	return tmpFile.Name(), nil
}
