// Package strings 提供字符串操作工具。
package strings

// 这里所谓的口语连接对于一些西方语言效果很好。欢迎为其他语言提交 PR，但请注意
// 某些语言的实现将比这里使用的实现更复杂。

import (
	"strings"
)

// Language 是一种口语语言。
type Language int

// 可用的口语语言。
const (
	DE Language = iota
	DK
	EN
	ES
	FR
	IT
	NO
	PT
	SE
)

// String 返回 [Language] 代码的英文名称。
func (l Language) String() string {
	return map[Language]string{
		DE: "German",
		DK: "Danish",
		EN: "English",
		ES: "Spanish",
		FR: "French",
		IT: "Italian",
		NO: "Norwegian",
		PT: "Portuguese",
		SE: "Swedish",
	}[l]
}

func (l Language) conjunction() string {
	switch l {
	case DE:
		return "und"
	case DK:
		return "og"
	case EN:
		return "and"
	case ES:
		return "y"
	case FR:
		return "et"
	case NO:
		return "og"
	case IT:
		return "e"
	case PT:
		return "e"
	case SE:
		return "och"
	default:
		return ""
	}
}

func (l Language) separator() string {
	switch l {
	case DE, DK, EN, ES, FR, NO, IT, PT, SE:
		return ", "
	default:
		return " "
	}
}

// EnglishJoin 使用逗号连接字符串切片，并在最后一项前使用 "and" 连接词。
// 可以选择应用牛津逗号。
//
// 示例：
//
//	str := EnglishJoin([]string{"meow", "purr", "raow"}, true)
//	fmt.Println(str) // meow, purr, and raow
func EnglishJoin(words []string, oxfordComma bool) string {
	return spokenLangJoin(words, EN, oxfordComma)
}

// SpokenLanguageJoin 使用逗号连接字符串切片，并在最后一项前使用连接词。
// 您可以使用 [Language] 指定语言。
//
// 如果您使用英语并且需要牛津逗号，请使用 [EnglishJoin]。
//
// 示例：
//
//	str := SpokenLanguageJoin([]string{"eins", "zwei", "drei"}, DE)
//	fmt.Println(str) // eins, zwei und drei
func SpokenLanguageJoin(words []string, language Language) string {
	return spokenLangJoin(words, language, false)
}

func spokenLangJoin(words []string, language Language, oxfordComma bool) string {
	conjunction := language.conjunction() + " "
	separator := language.separator()

	b := strings.Builder{}
	for i, word := range words {
		if word == "" {
			continue
		}

		if i == 0 {
			b.WriteString(word)
			continue
		}

		// 这是最后一个单词吗？
		if len(words) > 1 && i == len(words)-1 {
			// 如果请求并且语言是英语，则应用牛津逗号。
			if language == EN && oxfordComma && i > 1 {
				b.WriteString(separator)
			} else {
				b.WriteRune(' ')
			}

			b.WriteString(conjunction + word)
			continue
		}

		b.WriteString(separator + word)
	}
	return b.String()
}
