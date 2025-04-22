package url

import (
	"ai-calls/internal/server/app"
	"fmt"
	"net/url"
	"strings"
)

type URL struct {
	App       *app.Application
	RoutesMap map[string]string
	url       *url.URL
	values    *url.Values
}

func New(app *app.Application) *URL {
	u := &url.URL{}
	v := u.Query()
	u.Scheme = "https"
	return &URL{
		App:       app,
		url:       u,
		values:    &v,
		RoutesMap: make(map[string]string),
	}
}
func (s *URL) AddRouteToUrl(relativePath string) {
	var pathNormalizer = strings.NewReplacer(
		"-", ".",
		"/", ".",
		":", "%",
	)
	pathKey := pathNormalizer.Replace(relativePath)
	s.RoutesMap[pathKey] = relativePath
}

func (s *URL) clone() *URL {
	urlCopy := *s.url

	valuesCopy := url.Values{}
	for k, v := range *s.values {
		valuesCopy[k] = append([]string(nil), v...)
	}

	return &URL{
		App:       s.App,
		RoutesMap: s.RoutesMap,
		url:       &urlCopy,
		values:    &valuesCopy,
	}
}

func (s *URL) SetParam(key, value string) *URL {
	cloned := s.clone()
	cloned.values.Set(key, value)
	return cloned
}

func (s *URL) GetRouteUrl(scheme, path string, args ...string) string {
	cloned := s.clone()
	pathConverted := cloned.RoutesMap[fmt.Sprintf(".%s", path)]
	parts := strings.Split(pathConverted, "/")
	argIndex := 0
	for i, part := range parts {
		if strings.HasPrefix(part, ":") && argIndex < len(args) {
			parts[i] = args[argIndex]
			argIndex++
		}
	}
	finalPath := strings.Join(parts, "/")

	cloned.url.Scheme = scheme
	cloned.url.Host = cloned.App.Config.Host
	cloned.url.Path = finalPath
	cloned.url.RawQuery = cloned.values.Encode()

	return cloned.url.String()
}
