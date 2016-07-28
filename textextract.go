package articletext

/*
The function finds a DOM node containing majority of a text
in HTML document

Author: Roman Gelembjuk <roman@gelembjuk.com>
*/

import (
	"math"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/neurosnap/sentences.v1"
	"gopkg.in/neurosnap/sentences.v1/data"
)

type TextDescription struct {
	CountSentences        int
	AverageWords          int
	CountLongSentences    int
	CountGoodSentences    int
	CountCorrectSentences int
}

var tokenizer *sentences.DefaultSentenceTokenizer

func init() {
	// prepare tokenizer
	b, _ := data.Asset("data/english.json")

	// load the training data
	training, _ := sentences.LoadTraining(b)

	// create the default sentence tokenizer
	tokenizer = sentences.NewSentenceTokenizer(training)
}

func getPrimarySelection(s *goquery.Selection) *goquery.Selection {

	// prepare a selection for search. Add some descriptions for DOM nodes
	// to find most correct node with a text
	describeDocumentNode(s)

	// now find a node with a text and return it
	return findSelectionWithPrimaryText(s)
}

/*
* This is the core function. It checks a selection object and finds if this is a text node
* or it is needed to go deeper , inside a node that has most of text
 */
func findSelectionWithPrimaryText(s *goquery.Selection) *goquery.Selection {

	// if no children then return a text from this node
	if s.Children().Length() == 0 {
		return s
	}

	// variable to find a node with longest text inside it
	sort_by_count_sentences := 0
	// a node with longest text inside it
	var sort_by_text_node *goquery.Selection = nil

	// keep count of nodes containing more 2 sentences
	count_of_nodes_with_sentences := 0

	max_count_of_correct_sentences := 0

	// calcuate count of real symbols
	node_full_text_len := utf8.RuneCountInString(s.Text())

	top_total_count_of_correct_sentences := getNumbericAttribute(s, "totalcountofcorrectsentences")

	// all subnodes lengths
	tlengths := []int{}
	densityes := []int{}

	s.Children().Each(func(i int, sec *goquery.Selection) {
		totalcountofcorrectsentences := getNumbericAttribute(sec, "totalcountofcorrectsentences")

		if totalcountofcorrectsentences > 1 {
			count_of_nodes_with_sentences++

			if totalcountofcorrectsentences > max_count_of_correct_sentences {
				max_count_of_correct_sentences = totalcountofcorrectsentences
			}
		}

		// node text length
		tlen := utf8.RuneCountInString(sec.Text())

		html, _ := sec.Html()
		hlen := utf8.RuneCountInString(html)

		if tlen == 0 {
			// process next subnode
			return
		}

		tlengths = append(tlengths, tlen)

		density := (hlen / tlen)

		densityes = append(densityes, density)

		// check if this block is better then previous
		// choose better block only if previous is empty or
		// has less then 10 real sentences
		if totalcountofcorrectsentences > sort_by_count_sentences && sort_by_count_sentences < 10 {

			sort_by_count_sentences = totalcountofcorrectsentences
			sort_by_text_node = sec
		}

	})

	// if any nide with a text was found
	if sort_by_count_sentences > 0 {
		// calculate mean deviation
		lvar := getMeanDeviation(tlengths)

		// get relative value of a mean deviation agains full text length in a node
		lvarproc := (100 * lvar) / float64(node_full_text_len)

		// during tests we found that if this value is less 5
		// the a node is what we are looking for
		// it is the node with "main" text of a page
		if lvarproc < 15 && len(tlengths) > 3 ||
			(count_of_nodes_with_sentences > 2 &&
				float32(max_count_of_correct_sentences) < float32(top_total_count_of_correct_sentences)*0.8) {

			// we found that a text is equally distributed between subnodes
			// no need to go deeper

			return s
		}
		// go deeper inside a node with most of text

		return findSelectionWithPrimaryText(sort_by_text_node)
	}
	// no subnodes found. return a node itself
	return s
}

// describe a text inside a node and add description as pseudo attributes
func describeDocumentNode(s *goquery.Selection) *goquery.Selection {
	var totalcountofgoodsentences int
	var totalcountofcorrectsentences int
	var maxcountofflatsentences int

	countchildren := s.Children().Length()

	var sd TextDescription

	if countchildren > 0 {
		// for each child node check if to remove or not
		s.Children().Each(func(i int, sec *goquery.Selection) {

			// go deeper recursively
			describeDocumentNode(sec)

			// aggregate data to set to a node

			totalcountofgoodsentences += getNumbericAttribute(sec, "totalcountofgoodsentences")
			totalcountofcorrectsentences += getNumbericAttribute(sec, "totalcountofcorrectsentences")

			countsentences := getNumbericAttribute(sec, "maxcountofflatsentences")

			if countsentences > maxcountofflatsentences {
				maxcountofflatsentences = countsentences
			}

		})

		// describe sentences in this html tag only, drop child nodes
		secclone := getSelectionWihoutChildren(s)

		sd = describeSentences(secclone)

		totalcountofgoodsentences += sd.CountGoodSentences
		totalcountofcorrectsentences += sd.CountCorrectSentences

		if sd.CountGoodSentences > maxcountofflatsentences {
			maxcountofflatsentences = sd.CountGoodSentences
		}

	} else {
		// no child nodes
		//fmt.Println(s.Text())

		sd = describeSentences(s)
		totalcountofgoodsentences = sd.CountGoodSentences
		maxcountofflatsentences = sd.CountGoodSentences
		totalcountofcorrectsentences = sd.CountCorrectSentences
	}
	//fmt.Printf("set totalcountofgoodsentences ")
	// set attributes for the node
	s.SetAttr("countsentences", strconv.Itoa(sd.CountSentences))
	s.SetAttr("averagewords", strconv.Itoa(sd.AverageWords))
	s.SetAttr("countgoodsentences", strconv.Itoa(sd.CountGoodSentences))
	s.SetAttr("countlongsentences", strconv.Itoa(sd.CountLongSentences))
	s.SetAttr("totalcountofgoodsentences", strconv.Itoa(totalcountofgoodsentences))
	s.SetAttr("totalcountofcorrectsentences", strconv.Itoa(totalcountofcorrectsentences))
	s.SetAttr("maxcountofflatsentences", strconv.Itoa(maxcountofflatsentences))

	return s
}

/*
*
 */
func describeSentences(s *goquery.Selection) TextDescription {
	var d TextDescription

	var text string
	// get text of this node and then split for sentences
	if s.Children().Length() > 0 {
		text = getTextFromHtml(s)
	} else {
		text = s.Text()
	}

	sentences := tokenizer.Tokenize(text)

	d.CountSentences = len(sentences)
	//fmt.Println("==============================================")
	for _, s := range sentences {
		sentence := s.Text

		if len(sentence) == 0 {
			continue
		}

		c := len(get_words_from(sentence))
		//fmt.Println(sentence)

		d.AverageWords += c

		if c > 3 {
			// presume normal sentence usually has more 3 words
			d.CountLongSentences++

			if c < 25 {
				// but a sentence should not have nore 25 words. We will not
				// consider such sentence as a good one
				d.CountGoodSentences++

			}
			lastsymbol := sentence[len(sentence)-1:]

			if strings.ContainsAny(lastsymbol, ".?!") {
				d.CountCorrectSentences++
			}
		}

	}

	if d.CountSentences > 0 {
		d.AverageWords = int(d.AverageWords / d.CountSentences)
	}

	return d
}

func get_words_from(text string) []string {
	words := regexp.MustCompile("[^\\s]+")
	return words.FindAllString(text, -1)
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

// this function returns text in a selection with ignoring child nodes
// for "xxx <a>YYY</a>" result wil be "xxx "

func getSelectionWihoutChildren(s *goquery.Selection) *goquery.Selection {
	clone := s.Clone()

	// remove all child nodes in this selection
	clone.Children().Each(func(i int, sec *goquery.Selection) {
		sec.Remove()
	})

	return clone
}

func getNumbericAttribute(s *goquery.Selection, attr string) int {
	a, f := s.Attr(attr)

	if f {
		ai, _ := strconv.Atoi(a)
		return ai
	}
	return 0
}
