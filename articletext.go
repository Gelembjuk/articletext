package articletext

/*
The function extracts article text from a HTML page
It drops all additional elements from a html page (navigation, advertizing etc)

Author: Roman Gelembjuk <roman@gelembjuk.com>
*/

import (
	"io"
	"log"
	"math"
	"os"
	"sort"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"github.com/jaytaylor/html2text"
)

// liist of tags to ignore, as they dones't contain useful data
var skiphtmltags []string = []string{"script", "style", "noscript", "head"}

func init() {
	// to make lookup faster
	sort.Strings(skiphtmltags)
}

// extracts useful text from a html file
func GetArticleTextFromFile(filepath string) (string, error) {
	// create reader from file
	reader, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return GetArticleText(reader)
}

// extracts useful text from a html page presented by an url
func GetArticleTextFromUrl(url string) (string, error) {
	doc, err := goquery.NewDocument(url)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return processArticle(doc)
}

// extracts useful text from a html document presented as a Reader object
func GetArticleText(input io.Reader) (string, error) {

	doc, err := goquery.NewDocumentFromReader(input)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return processArticle(doc)
}

// the function prepares a document for analysing
// cleans a DOM object and starts analysing
func processArticle(doc *goquery.Document) (string, error) {

	if doc == nil {
		return "", nil
	}

	// preprocess. Remove all tags that are not useful and can make parsing wrong
	cleanDocument(doc.Selection)

	return getTextFromSelection(doc.Selection), nil
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

func getTextFromSelection(s *goquery.Selection) string {
	// if no children then return a text from this node
	if s.Children().Length() == 0 {
		return getTextFromHtml(s)
	}

	// variable to fid a node with longest text inside it
	sort_by_text_len := 0
	// a node with longest text inside it
	var sort_by_text_node *goquery.Selection = nil

	// calcuate count of real symbols
	node_full_text_len := utf8.RuneCountInString(s.Text())

	// all subnodes lengths
	tlengths := []int{}

	s.Children().Each(func(i int, sec *goquery.Selection) {
		// node text length
		tlen := utf8.RuneCountInString(sec.Text())

		tlengths = append(tlengths, tlen)

		if tlen == 0 {
			// process next subnode
			return
		}

		// check if this is bigger and set to bigger if yes
		if tlen > sort_by_text_len {
			sort_by_text_len = tlen
			sort_by_text_node = sec
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
			return getTextFromHtml(s)
		}
		// go deeper inside a node with most of text
		return getTextFromSelection(sort_by_text_node)
	}
	// no subnodes found. return a node itself
	return getTextFromHtml(s)
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
