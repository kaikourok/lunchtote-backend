package service

import (
	"html"
	"math/rand"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

var newLineReplacer *strings.Replacer
var hrRegex *regexp.Regexp
var colorStartRegex *regexp.Regexp
var d100Regex *regexp.Regexp
var d6Regex *regexp.Regexp
var untagRegex *regexp.Regexp

/*-------------------------------------------------------------------------------------------------
  init
-------------------------------------------------------------------------------------------------*/

func init() {
	newLineReplacer = strings.NewReplacer(
		"\r\n", "<br>",
		"\r", "<br>",
		"\n", "<br>",
	)
	hrRegex = regexp.MustCompile(`(?i)(<br>)?\[hr\]`)
	colorStartRegex = regexp.MustCompile(`\[#[0-9a-f]{6}\]`)
	d100Regex = regexp.MustCompile(`(?i)\[d100\]`)
	d6Regex = regexp.MustCompile(`(?i)\[d6\]`)
	untagRegex = regexp.MustCompile(`<.*?>`)
}

/*-------------------------------------------------------------------------------------------------
	Inner Util Functions
-------------------------------------------------------------------------------------------------*/

func trimStart(s string, remove string) string {
	if strings.HasPrefix(s, remove) {
		return s[len(remove):]
	} else {
		return s
	}
}

func trimEnd(s string, remove string) string {
	if strings.HasSuffix(s, remove) {
		return s[:len(s)-len(remove)]
	} else {
		return s
	}
}

func isImageFileName(s string) bool {
	_, err := uuid.Parse(s)
	return err == nil
}

func isImagePath(url string, uploadPath string) bool {
	allowedExtensions := []string{"png", "gif", "jpg", "jpeg"}
	uploadPathDepth := strings.Count(uploadPath, "/")

	prefixIndex := strings.Index(url, uploadPath)
	if prefixIndex != 0 {
		return false
	}

	paths := strings.Split(url, "/")
	if len(paths) != 3+uploadPathDepth {
		return false
	}

	characterDirectory := paths[1+uploadPathDepth]
	_, err := strconv.Atoi(characterDirectory)
	if err != nil {
		return false
	}

	fileFullname := paths[2+uploadPathDepth]
	extensionIndex := strings.LastIndex(fileFullname, ".")
	if extensionIndex == -1 {
		return false
	}
	extension := fileFullname[extensionIndex+1:]

	found := false
	for _, allowedExtension := range allowedExtensions {
		if extension == allowedExtension {
			found = true
			break
		}
	}
	if !found {
		return false
	}

	filename := fileFullname[:extensionIndex]

	return isImageFileName(filename)
}

func splitTagSections(target string, startTag string, endTag string) (before string, inner string, after string, index int) {
	t := strings.ToLower(target)
	sp := len(t)

	for {
		startIndex := strings.LastIndex(t[:sp], startTag)
		if startIndex == -1 {
			break
		}

		sp = startIndex
		endRangeStartIndex := strings.Index(t[sp+len(startTag):], startTag)
		endRangeEndIndex := strings.Index(t[sp+len(startTag):], endTag)

		if endRangeEndIndex == -1 || (endRangeStartIndex != -1 && endRangeStartIndex < endRangeEndIndex) {
			continue
		}

		return target[:sp], target[sp+len(startTag) : sp+len(startTag)+endRangeEndIndex], target[sp+len(startTag)+endRangeEndIndex+len(endTag):], sp
	}

	return "", "", "", -1
}

func replaceImageTag(target string, uploadPath string, startTag string, endTag string, class string) (result string, found bool) {
	sp := len(target)

	for {
		before, inner, after, index := splitTagSections(target[:sp], startTag, endTag)
		if index == -1 {
			break
		}

		sp = index
		url := strings.TrimSpace(inner)
		if !isImagePath(url, uploadPath) {
			continue
		}

		return trimEnd(before, "<br>") + `<img class="` + class + `" src="` + url + `">` + trimStart(after+target[sp:], "<br>"), true
	}

	return target, false
}

/*-------------------------------------------------------------------------------------------------
	GeneralTag
-------------------------------------------------------------------------------------------------*/

func replaceGeneralTag(target string, startTag string, endTag string, startTagTo string, endTagTo string) (result string, found bool) {
	before, inner, after, index := splitTagSections(target, startTag, endTag)
	if index == -1 {
		return target, false
	}

	return before + startTagTo + inner + endTagTo + after, true
}

func replaceGeneralTagAll(target string, startTag string, endTag string, startTagTo string, endTagTo string) string {
	result := target
	var found bool

	for {
		result, found = replaceGeneralTag(result, startTag, endTag, startTagTo, endTagTo)
		if !found {
			break
		}
	}

	return result
}

/*-------------------------------------------------------------------------------------------------
	ImageTag
-------------------------------------------------------------------------------------------------*/

func replaceImageCenterTag(target string, uploadPath string) (result string, found bool) {
	result, found = replaceImageTag(target, uploadPath, "[img]", "[/img]", "cutin")
	return
}

func replaceImageLeftTag(target string, uploadPath string) (result string, found bool) {
	result, found = replaceImageTag(target, uploadPath, "[img-l]", "[/img-l]", "cutin cutin-left")
	return
}

func replaceImageRightTag(target string, uploadPath string) (result string, found bool) {
	result, found = replaceImageTag(target, uploadPath, "[img-r]", "[/img-r]", "cutin cutin-right")
	return
}

func replaceImageTagAll(target string, uploadPath string) string {
	result := target
	var found bool

	for {
		result, found = replaceImageCenterTag(result, uploadPath)
		if !found {
			break
		}
	}

	for {
		result, found = replaceImageLeftTag(result, uploadPath)
		if !found {
			break
		}
	}

	for {
		result, found = replaceImageRightTag(result, uploadPath)
		if !found {
			break
		}
	}

	return result
}

/*-------------------------------------------------------------------------------------------------
	MessageTag
-------------------------------------------------------------------------------------------------*/

func replaceMessageTag(target string, uploadPath string) (result string, found bool) {
	nameStartTag := "[name]"
	nameEndTag := "[/name]"
	iconStartTag := "[icon]"
	iconEndTag := "[/icon]"

	sp := len(target)

	for {
		before, inner, after, index := splitTagSections(target[:sp], "[message]", "[/message]")
		if index == -1 {
			break
		}

		lines := strings.Split(inner, "<br>")
		var startEmptyLineCount, endEmptyLineCount int

		for i := 0; i < len(lines); i++ {
			if lines[i] == "" {
				startEmptyLineCount++
			} else {
				break
			}
		}

		for i := len(lines) - 1; i >= 0; i-- {
			if lines[i] == "" {
				endEmptyLineCount++
			} else {
				break
			}
		}

		if len(lines) < startEmptyLineCount+endEmptyLineCount {
			continue
		}

		var b strings.Builder
		b.Grow(len(inner) + 128)

		name := ""
		icon := ""
		added := false

		for _, line := range lines {
			if added {
				b.WriteString("<br>")
			}

			if strings.HasPrefix(line, nameStartTag) && strings.HasSuffix(line, nameEndTag) {
				name = `<div class="message-name">` + line[len(nameStartTag):len(line)-len(nameEndTag)] + `</div>`
				continue
			}

			if strings.HasPrefix(line, iconStartTag) && strings.HasSuffix(line, iconEndTag) {
				parsed := strings.TrimSpace(line[len(iconStartTag) : len(line)-len(iconEndTag)])

				if isImagePath(parsed, uploadPath) {
					icon = `<div class="message-icon-wrapper"><img class="message-icon" src="` + icon + `"></div>`
					continue
				}
			}

			b.WriteString(line)
			added = true
		}

		var resultBuilder strings.Builder
		b.Grow(len(target) + 512)

		resultBuilder.WriteString(trimEnd(before, "<br>"))
		resultBuilder.WriteString(`<section class="message">`)
		resultBuilder.WriteString(icon)
		resultBuilder.WriteString(`<div class="message-content">`)
		resultBuilder.WriteString(name)
		resultBuilder.WriteString(`<div class="">`)
		resultBuilder.WriteString(b.String())
		resultBuilder.WriteString(`</div>`)
		resultBuilder.WriteString(`</div>`)
		resultBuilder.WriteString(`</section>`)
		resultBuilder.WriteString(trimStart(after+target[sp:], "<br>"))

		return resultBuilder.String(), true
	}

	return target, false
}

func replaceMessageTagAll(target string, uploadPath string) string {
	result := target
	var found bool

	for {
		result, found = replaceMessageTag(result, uploadPath)
		if !found {
			break
		}
	}

	return result
}

/*-------------------------------------------------------------------------------------------------
	RubyTag
-------------------------------------------------------------------------------------------------*/

func replaceRubyTag(target string) (result string, found bool) {
	textStartTag := "[rt]"
	textEndTag := "[/rt]"
	rubyStartTag := "[rb]"
	rubyEndTag := "[/rb]"

	sp := len(target)

	for {
		before, inner, after, index := splitTagSections(target[:sp], textStartTag, rubyEndTag)
		if index == -1 {
			break
		}
		sp = index

		if strings.Count(inner, textEndTag) != 1 {
			continue
		}

		if strings.Count(inner, rubyStartTag) != 1 {
			continue
		}

		separatorIndex := strings.Index(inner, textEndTag+rubyStartTag)
		if separatorIndex == -1 {
			continue
		}

		text := inner[:separatorIndex]
		ruby := inner[separatorIndex+len(textEndTag+rubyStartTag):]

		return before + `<ruby>` + text + `<rp>(</rp><rt>` + ruby + `</rt><rp>)</rp></ruby>` + after + target[sp:], true
	}

	return target, false
}

func replaceRubyTagAll(target string) string {
	result := target
	var found bool

	for {
		result, found = replaceRubyTag(result)
		if !found {
			break
		}
	}

	return result
}

/*-------------------------------------------------------------------------------------------------
	ColorTag
-------------------------------------------------------------------------------------------------*/

func replaceColorTag(target string) (result string, found bool) {
	t := strings.ToLower(target)
	sp := len(t)

	startIndexes := colorStartRegex.FindAllStringIndex(t, -1)

	for i := len(startIndexes) - 1; i >= 0; i-- {
		startIndex := startIndexes[i][0]
		startColor := t[startIndex+1 : startIndex+1+6]

		startTag := "[#" + startColor + "]"
		endTag := "[/#" + startColor + "]"

		endRangeStartIndex := strings.Index(t[sp+len(startTag):], startTag)
		endRangeEndIndex := strings.Index(t[sp+len(startTag):], endTag)

		if endRangeEndIndex == -1 || (endRangeStartIndex != -1 && endRangeStartIndex < endRangeEndIndex) {
			continue
		}

		return target[:sp] + startTag + target[sp+len(startTag):sp+len(startTag)+endRangeEndIndex] + endTag + target[sp+len(startTag)+endRangeEndIndex+len(endTag):], true
	}

	return target, false
}

func replaceColorTagAll(target string) string {
	result := target
	var found bool

	for {
		result, found = replaceColorTag(result)
		if !found {
			break
		}
	}

	return result
}

/*-------------------------------------------------------------------------------------------------
	Publics
-------------------------------------------------------------------------------------------------*/

func stylizeBasic(message string, uploadPath string) string {
	var s string
	s = html.EscapeString(message)
	s = newLineReplacer.Replace(s)
	s = replaceGeneralTagAll(s, "[+]", "[/+]", `<span class="larger">`, `</span>`)
	s = replaceGeneralTagAll(s, "[-]", "[/-]", `<span class="smaller">`, `</span>`)
	s = replaceGeneralTagAll(s, "[b]", "[/b]", `<span class="bold">`, `</span>`)
	s = replaceGeneralTagAll(s, "[s]", "[/s]", `<span class="strike">`, `</span>`)
	s = replaceGeneralTagAll(s, "[i]", "[/i]", `<span class="italic">`, `</span>`)
	s = replaceGeneralTagAll(s, "[u]", "[/u]", `<span class="underline">`, `</span>`)
	s = replaceColorTagAll(s)
	s = replaceRubyTagAll(s)
	s = replaceImageTagAll(s, uploadPath)
	return s
}

func StylizeMessage(message string, uploadPath string) string {
	s := stylizeBasic(message, uploadPath)
	s = d100Regex.ReplaceAllStringFunc(s, func(target string) string {
		return `<span class="dice d100">` + strconv.Itoa(rand.Intn(100)+1) + `</span>`
	})
	s = d6Regex.ReplaceAllStringFunc(s, func(target string) string {
		return `<span class="dice d6">` + strconv.Itoa(rand.Intn(6)+1) + `</span>`
	})
	return s
}

func StylizeTextEntry(profile string, uploadPath string) string {
	s := stylizeBasic(profile, uploadPath)
	s = hrRegex.ReplaceAllString(s, `<hr class="message-hr">`)
	s = replaceMessageTagAll(s, uploadPath)
	return s
}

func ConvertMessageToSearchText(message string) string {
	return untagRegex.ReplaceAllString(message, " ")
}
