package url

import (
	"fmt"
	"strings"
)

func Handle(path string, iterBlck func(string, func(string))) (params map[string]string, err error) {
	params = map[string]string{
		"$path": path,
	}
	fragments := strings.Split(path, "/")
	for _, fragment := range fragments {
		if err != nil {
			params = nil
			return
		}
		iterBlck(
			fragment,
			func(k string) {
				_, exists := params[k]
				if exists {
					err = fmt.Errorf("parameter \"%s\" already marked", k)
					return
				}
				params[k] = fragment
			})
	}
	return
}