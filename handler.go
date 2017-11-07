package urlshort

import (
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

type mapHandler struct {
	paths    map[string]string
	fallback http.Handler
}

var mh mapHandler

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	mh = mapHandler{
		pathsToUrls,
		fallback,
	}
	return handleMap
}

func handleMap(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	path, ok := mh.paths[url.RequestURI()]
	if ok {
		http.Redirect(w, r, path, http.StatusFound)
		return
	}
	mh.fallback.ServeHTTP(w, r)
}

/*
yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
*/

type paths []map[string]string

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//     - path: /some-path
//       url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var urlPaths paths
	if err := yaml.Unmarshal(yml, &urlPaths); err != nil {
		return nil, err
	}

	fmt.Printf("--- urlPaths:\n\t%v\n\n", urlPaths)
	return fallback.ServeHTTP, nil
}
