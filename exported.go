package articletext

/*
The function extracts article text from a HTML page
It drops all additional elements from a html page (navigation, advertizing etc)

This file contains exported functiosn of a package

Author: Roman Gelembjuk <roman@gelembjuk.com>
*/

import (
	"io"
	"log"
	"os"

	"github.com/PuerkitoBio/goquery"
)

// extracts useful html part from a html document presented as a Reader object
func GetArticleHtmlFromReader(input io.Reader) (string, error) {
	doc, err := goquery.NewDocumentFromReader(input)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return processArticleToHtml(doc)
}

// extracts useful html part from a html page presented by an url
func GetArticleHtmlFromUrl(url string) (string, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return processArticleToHtml(doc)
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

	return processArticleToText(doc)
}

// extracts useful text from a html document presented as a Reader object
func GetArticleText(input io.Reader) (string, error) {

	doc, err := goquery.NewDocumentFromReader(input)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return processArticleToText(doc)
}

// extracts useful text from a html file
// returns a DOM signature
func GetArticleSignatureFromFile(filepath string) (string, error) {
	// create reader from file
	reader, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return GetArticleSignature(reader)
}

// extracts useful text from a html page presented by an url
func GetArticleSignatureFromUrl(url string) (string, error) {
	doc, err := goquery.NewDocument(url)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return processArticleToSignature(doc)
}

// extracts useful text from a html document presented as a Reader object
func GetArticleSignature(input io.Reader) (string, error) {

	doc, err := goquery.NewDocumentFromReader(input)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return processArticleToSignature(doc)
}

// extracts useful text from a html file
func GetArticleTextFromFileByPath(filepath string, path string) (string, error) {
	// create reader from file
	reader, err := os.Open(filepath)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return GetArticleTextByPath(reader, path)
}

// extracts useful text from a html page presented by an url
func GetArticleTextFromUrlByPath(url string, path string) (string, error) {
	doc, err := goquery.NewDocument(url)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return getTextByPathFromDocument(doc, path)
}

// extracts useful text from a html document presented as a Reader object
func GetArticleTextByPath(input io.Reader, path string) (string, error) {

	doc, err := goquery.NewDocumentFromReader(input)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return getTextByPathFromDocument(doc, path)
}
