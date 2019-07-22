package router

import (
	"log"
	"os"
	"path/filepath"
)

// Returns file walk function
// WalkFunc apply a filter for a given file name

func Lookup(files *[]string, routeFileName string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if len(routeFileName) == 0 {
			routeFileName = "route_config.json.template"
		}

		if info.Name() == routeFileName {
			*files = append(*files, path)
		}

		return nil
	}
}
