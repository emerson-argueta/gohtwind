package utils

import (
	"embed"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadFile(url string, dest string, projectName string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath.Join(projectName, dest))
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

// loadEmbeddedEnv reads the .env file from the embedded file system.
func LoadEmbeddedEnv(envFile embed.FS) (map[string]string, error) {
	// Read the embedded .env file
	env, err := envFile.ReadFile(".env") // make sure the path is correct relative to the embedding directive
	if err != nil {
		return nil, err
	}

	// Parse the environment variables from the byte content
	return godotenv.Unmarshal(string(env))
}
