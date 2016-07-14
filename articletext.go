package articletext

/*
The function extracts article text from a HTML page
It drops all additional elements from a html page (navigation, advertizing etc)

This file containes internal functions and a logic

Author: Roman Gelembjuk <roman@gelembjuk.com>
*/

import (
	"math"
	"sort"
	"unicode/utf8"

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

	// preprocess. Remove all tags that are not useful and can make parsing wrong
	cleanDocument(doc.Selection)

	selection := getPrimarySelection(doc.Selection)

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

/*
* This is the core function. It checks a selection object and finds if this is a text node
* or it is needed to go deeper , inside a node that has most of text
 */
func getPrimarySelection(s *goquery.Selection) *goquery.Selection {

	// if no children then return a text from this node
	if s.Children().Length() == 0 {
		return s
	}

	// variable to find a node with longest text inside it
	sort_by_text_len := 0
	// a node with longest text inside it
	var sort_by_text_node *goquery.Selection = nil
	// keep correlation of text to html in a node
	node_text_density := 0

	// variable to keep a previous "biggest" node
	// it can help in some cases when an article has many commends below
	// and comments block is bigger then an article itself
	// we have same set of variables as for biggest node
	sort_by_text_len_previous := 0
	var sort_by_text_node_previous *goquery.Selection = nil
	node_text_density_previous := 0

	// calcuate count of real symbols
	node_full_text_len := utf8.RuneCountInString(s.Text())

	// all subnodes lengths
	tlengths := []int{}
	densityes := []int{}

	s.Children().Each(func(i int, sec *goquery.Selection) {

		// node text length
		tlen := utf8.RuneCountInString(sec.Text())

		html, _ := sec.Html()
		hlen := utf8.RuneCountInString(html)

		tlengths = append(tlengths, tlen)

		if tlen == 0 {
			// process next subnode
			return
		}

		density := (hlen / tlen)

		densityes = append(densityes, density)

		// check if this is bigger and set to bigger if yes
		if tlen > sort_by_text_len {
			sort_by_text_len_previous = sort_by_text_len
			sort_by_text_node_previous = sort_by_text_node
			node_text_density_previous = node_text_density

			sort_by_text_len = tlen
			sort_by_text_node = sec
			node_text_density = density
		}

	})

	// if any nide with a text was found
	if sort_by_text_len > 0 {
		// calculate mean deviation
		lvar := getMeanDeviation(tlengths)

		// get relative value of a mean deviation agains full text length in a node
		lvarproc := (100 * lvar) / float64(node_full_text_len)

		// during tests we found that if this value is less 5
		// the a node is what we are looking for
		// it is the node with "main" text of a page
		if lvarproc < 5 && s.Children().Length() > 3 {

			// we found that a text is equally distributed between subnodes
			// no need to go deeper

			return s
		}
		// go deeper inside a node with most of text

		// but there is an adge case, when an article has commens below. and comments section is bigger from an aticle

		// in previous biggest node there is more text relative to html
		// commenst list usually has a lot of html formatting
		// we consider only really long previous part.  not less 30% of max
		if node_text_density_previous*2 <= node_text_density &&
			float32(sort_by_text_len)*0.3 < float32(sort_by_text_len_previous) {
			// there is much more text in previous node
			// we will continue to work with previous node inn this case
			sort_by_text_node = sort_by_text_node_previous

		} else if float32(sort_by_text_len)*0.7 < float32(sort_by_text_len_previous) &&
			node_text_density_previous < node_text_density {
			// length of previous node is not so less then maximum next
			// so, there is high probability it is an article and maximum is comments section
			sort_by_text_node = sort_by_text_node_previous

		}

		return getPrimarySelection(sort_by_text_node)
	}
	// no subnodes found. return a node itself
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

// The function to calculate a Mean Deviation for a list of integer values
func getMeanDeviation(list []int) float64 {

	if len(list) < 1 {
		return 0.0
	}

	sum := 0

	for i := range list {
		sum += list[i]
	}

	// calculate arithmetic mean
	avg := float64(sum / len(list))

	number1 := 0.0

	for i := range list {
		number1 += math.Abs(float64(list[i]) - avg)
	}
	// calculate mean deviation
	meandeviation := number1 / float64(len(list))

	return meandeviation
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
