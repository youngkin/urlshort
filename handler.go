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
	path, ok := mh.paths[r.URL.RequestURI()]
	if ok {
		http.Redirect(w, r, path, http.StatusFound)
		return
	}
	mh.fallback.ServeHTTP(w, r)
}

type yamlHandler struct {
	pathMaps []map[string]string
	fallback http.Handler
}

var yh yamlHandler

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
	yh = yamlHandler{}
	if err := yaml.Unmarshal(yml, &yh.pathMaps); err != nil {
		return nil, err
	}
	yh.fallback = fallback

	fmt.Printf("--- urlPaths:\n\t%v\n\n", yh.pathMaps)
	return handleYaml, nil
}

func handleYaml(w http.ResponseWriter, r *http.Request) {
	for _, pathMap := range yh.pathMaps {
		if pathMap["path"] == r.URL.RequestURI() {
			http.Redirect(w, r, pathMap["url"], http.StatusFound)
			return
		}
	}
	yh.fallback.ServeHTTP(w, r)
}
