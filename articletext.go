package articletext

/*
The function extracts article text from a HTML page
It drops all additional elements from a html page (navigation, advertizing etc)

This file containes internal functions and a logic

Author: Roman Gelembjuk <roman@gelembjuk.com>
*/

import (
	"sort"

	"github.com/PuerkitoBio/goquery"
	"github.com/jaytaylor/html2text"
)

// liist of tags to ignore, as they dones't contain useful data
var skiphtmltags []string = []string{"script", "style", "noscript", "head",
	"header", "footer", "nav"}

func init() {
	// to make lookup faster
	sort.Strings(skiphtmltags)
}

// the function prepares a document for analysing
// cleans a DOM object and starts analysing
func processArticle(doc *goquery.Document, responsetype int) (string, error) {

	if doc == nil {
		return "", nil
	}

	// get clone of a selection. Clone is neede,d because we willdo some transformations

	docselection := doc.Selection.Clone()

	// preprocess. Remove all tags that are not useful and can make parsing wrong
	cleanDocument(docselection)

	// get a selection that contains a text of a page (only primary or article text)
	selection := getPrimarySelection(docselection)

	if responsetype == 2 {
		// return parent node path and attributes
		return getSelectionSignature(selection), nil
	}

	return getTextFromHtml(selection), nil
}

// clean HTML document. Removes all tags that are not useful
func cleanDocument(s *goquery.Selection) *goquery.Selection {
	tagname := goquery.NodeName(s)

	if checkTagsToSkip(tagname) {
		s.Remove()
		return nil
	}
	// for each child node check if to remove or not
	s.Children().Each(func(i int, sec *goquery.Selection) {
		tagname := goquery.NodeName(sec)

		if checkTagsToSkip(tagname) {

			sec.Remove()

			return
		}
		// go deeper recursively
		cleanDocument(sec)
	})

	return s
}

// convert HTML to text from a DOM node
// we ignore errors in this function
func getTextFromHtml(s *goquery.Selection) string {
	// gethtml from a node
	html, _ := s.Html()
	// convert to text
	text, err := html2text.FromString(html)

	if err != nil {
		return ""
	}

	return text
}

// check if aword (string) is in an array of tags
// we have list of tags to ignore some not useful tags
func checkTagsToSkip(tag string) bool {
	for _, v := range skiphtmltags {
		if v == tag {
			return true
		}
	}
	return false
}
