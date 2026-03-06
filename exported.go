package articletext

/*
The package is used extracts article text from a HTML page
It drops all additional elements from a html page (navigation, advertizing etc)

This file contains exported functiosn of a package. It is entry point of the package

Author: Roman Gelembjuk <roman@gelembjuk.com>
*/

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// extracts useful text from a html file
func GetArticleTextFromFile(filepath string) (string, error) {
	// create reader from file
	reader, err := os.Open(filepath)

	if err != nil {
		return "", err
	}
	defer reader.Close()

	return GetArticleText(reader)
}

// extracts useful text from a html page presented by an url
func GetArticleTextFromUrl(url string) (string, error) {
	doc, err := newDocumentFromURL(url)

	if err != nil {
		return "", err
	}

	return processArticle(doc, 1)
}

// extracts useful text from a html document presented as a Reader object
func GetArticleText(input io.Reader) (string, error) {

	doc, err := goquery.NewDocumentFromReader(input)

	if err != nil {
		return "", err
	}

	return processArticle(doc, 1)
}

// extracts useful text from a html file
// returns a DOM signature
func GetArticleSignatureFromFile(filepath string) (string, error) {
	// create reader from file
	reader, err := os.Open(filepath)

	if err != nil {
		return "", err
	}
	defer reader.Close()

	return GetArticleSignature(reader)
}

// extracts useful text from a html page presented by an url
func GetArticleSignatureFromUrl(url string) (string, error) {
	doc, err := newDocumentFromURL(url)

	if err != nil {
		return "", err
	}

	return processArticle(doc, 2)
}

// extracts useful text from a html document presented as a Reader object
func GetArticleSignature(input io.Reader) (string, error) {

	doc, err := goquery.NewDocumentFromReader(input)

	if err != nil {
		return "", err
	}

	return processArticle(doc, 2)
}

// extracts useful text from a html file
func GetArticleTextFromFileByPath(filepath string, path string) (string, error) {
	// create reader from file
	reader, err := os.Open(filepath)

	if err != nil {
		return "", err
	}
	defer reader.Close()

	return GetArticleTextByPath(reader, path)
}

// extracts useful text from a html page presented by an url
func GetArticleTextFromUrlByPath(url string, path string) (string, error) {
	doc, err := newDocumentFromURL(url)

	if err != nil {
		return "", err
	}

	return getTextByPathFromDocument(doc, path)
}

// extracts useful text from a html document presented as a Reader object
func GetArticleTextByPath(input io.Reader, path string) (string, error) {

	doc, err := goquery.NewDocumentFromReader(input)

	if err != nil {
		return "", err
	}

	return getTextByPathFromDocument(doc, path)
}

// the functions finds a path (selector, signature) for each url and returns one that was found most often
func GetOptimalArticleSignatureByUrls(urls []string) (string, error) {

	return getOptimalArticleSignatureByUrls(urls)
}

func newDocumentFromURL(url string) (*goquery.Document, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return nil, fmt.Errorf("unexpected HTTP status %d for %s", resp.StatusCode, url)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}
