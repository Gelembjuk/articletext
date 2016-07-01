## ArticleText

Golang package with a function to extract useful text from a HTML document.

A function analyses a html code and drops everything related to navigation, advertising etc. 
Extracts only useful contents of a document, text of a central element.

### Installation

go get github.com/gelembjuk/articletext

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

