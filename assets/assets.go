package assets

import "net/http"
import "github.com/rakyll/statik/fs"

var (
	// store static files in memory by statik
	FileSystem http.FileSystem

	// if prefix is not empty, we get file content from disk
	prefixPath string
)

func Load(path string) (err error) {
	prefixPath = path
	if prefixPath != "" {
		FileSystem = http.Dir(prefixPath)
		return nil
	} else {
		FileSystem, err = fs.New()
	}
	return err
}
