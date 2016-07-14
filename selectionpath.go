package articletext

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

/*
Extract text by DOM path, aka jquery style
*/
func getTextByPathFromDocument(doc *goquery.Document, path string) (string, error) {
	sel := doc.Find(path)

	if sel != nil {

		return getTextFromHtml(sel), nil
	}

	return "nothing", nil
}

// this function returns some specific signature of a selection
// so it can be easy found to get data quickly next time
func getSelectionSignature(s *goquery.Selection) string {
	var signature string

	tag, _ := goquery.OuterHtml(s)

	pos := strings.Index(tag, ">")

	if pos > -1 {
		tag = tag[1:pos]
	} else {
		return ""
	}

	signature = convertTagToJqueryFormat(tag, s)

	s.Parents().Each(func(i int, sec *goquery.Selection) {
		ohtml, _ := goquery.OuterHtml(sec)

		pos := strings.Index(ohtml, ">")

		if pos > -1 {
			ohtml = ohtml[1:pos]
		}

		tag := convertTagToJqueryFormat(ohtml, sec)

		signature = tag + " " + signature
	})

	return signature
}

func convertTagToJqueryFormat(tag string, s *goquery.Selection) string {
	tagitself := tag

	pos := strings.Index(tag, " ")

	if pos > -1 {
		tagitself = tag[0:pos]
	} else {

		return tag
	}

	class, found := s.Attr("class")

	if found && class != "" {
		pos := strings.Index(class, " ")
		// leave only a first class from a list
		if pos > -1 {
			class = class[0:pos]
		}

		tagitself = tagitself + "." + class
	}

	return tagitself
}
