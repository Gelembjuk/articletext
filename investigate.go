package articletext

/*
This file contains a function to investigate a list of urls and chooose optimal
path (selector) to use later for quick extracting a text from HTML document

Author: Roman Gelembjuk <roman@gelembjuk.com>
*/

import (
	"errors"
)

// the functions finds a path (selector, signature) for each url and returns one that was found most often
func getOptimalArticleSignatureByUrls(urls []string) (string, error) {

	if len(urls) < 1 {
		return "", errors.New("No urls provided")
	}

	var paths map[string]int
	paths = make(map[string]int)

	for _, url := range urls {

		path, err := GetArticleSignatureFromUrl(url)

		if err != nil {
			return "", err
		}

		if count, ok := paths[path]; ok {
			paths[path] = count + 1
		} else {
			paths[path] = 1
		}
	}

	// find what path has maximum of occurences
	maxpath := ""
	maxval := 0

	for k, v := range paths {
		if v > maxval {
			maxval = v
			maxpath = k
		}
	}

	return maxpath, nil
}
