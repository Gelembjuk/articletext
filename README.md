## ArticleText

Golang package with a function to extract useful text from a HTML document.

A function analyses a html code and drops everything related to navigation, advertising etc. 
Extracts only useful contents of a document, text of a central element.

### Installation

go get github.com/gelembjuk/articletext

### Manual

There are 3 types of exported functions.

1. Functions to get a text from a HTML document. From 3 different types of sources

#### GetArticleText(input io.Reader)

#### GetArticleTextFromFile(filepath string)

#### GetArticleTextFromUrl(url string)

2. Functions to return a path (signature) for a text location block. The path is a JQuery style selector - tags with classes.

Also 3 functions for input form different sources

#### GetArticleSignature(input io.Reader)

#### GetArticleSignatureFromFile(filepath string)

#### GetArticleSignatureFromUrl(url string)

Result of these functions is somethign like "body div div div.content div.article div.text" . And then this path can be used to get a text with one of following functions

3. Functions to get a text from a HTML document using a path (signature) in a JQuery style. A path can be get by using one of functions from  blcok 2, or prepared manually

#### GetArticleTextByPath(input io.Reader, path string) 

#### GetArticleTextFromFileByPath(filepath string, path string)

#### GetArticleTextFromUrlByPath(url string, path string)

### Example 

```
package main

import (
	"fmt"
	"os"
	"github.com/gelembjuk/articletext"
)

func main() {

	url := os.Args[1]
	text, err := articletext.GetArticleTextFromUrl(url)
	
	fmt.Println(text)
}
```

### Author

Roman Gelembjuk (@gelembjuk)

