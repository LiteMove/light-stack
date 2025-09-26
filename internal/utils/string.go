package utils

import (
	"regexp"
	"strings"
	"unicode"
)

// ToCamelCase 转换字符串为驼峰命名
func ToCamelCase(str string) string {
	if str == "" {
		return ""
	}

	// 分割字符串
	words := splitWords(str)
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	for i, word := range words {
		if i == 0 {
			// 第一个单词首字母大写
			result.WriteString(capitalize(word))
		} else {
			// 其他单词首字母大写
			result.WriteString(capitalize(word))
		}
	}

	return result.String()
}

// ToPascalCase 转换字符串为帕斯卡命名（首字母大写的驼峰）
func ToPascalCase(str string) string {
	if str == "" {
		return ""
	}

	words := splitWords(str)
	if len(words) == 0 {
		return ""
	}

	var result strings.Builder
	for _, word := range words {
		result.WriteString(capitalize(word))
	}

	return result.String()
}

// ToSnakeCase 转换字符串为下划线命名
func ToSnakeCase(str string) string {
	if str == "" {
		return ""
	}

	// 在大写字母前添加下划线（除了第一个字符）
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	result := re.ReplaceAllString(str, "${1}_${2}")

	// 转换为小写
	return strings.ToLower(result)
}

// ToKebabCase 转换字符串为短横线命名
func ToKebabCase(str string) string {
	if str == "" {
		return ""
	}

	// 在大写字母前添加短横线（除了第一个字符）
	re := regexp.MustCompile("([a-z0-9])([A-Z])")
	result := re.ReplaceAllString(str, "${1}-${2}")

	// 转换为小写
	return strings.ToLower(result)
}

// splitWords 分割单词
func splitWords(str string) []string {
	// 替换常见分隔符为空格
	separators := []string{"_", "-", ".", " ", "\t", "\n"}
	result := str
	for _, sep := range separators {
		result = strings.ReplaceAll(result, sep, " ")
	}

	// 在大写字母前添加空格（除了第一个字符）
	var words []rune
	for i, r := range result {
		if i > 0 && unicode.IsUpper(r) && unicode.IsLower(rune(result[i-1])) {
			words = append(words, ' ')
		}
		words = append(words, r)
	}

	// 分割并清理空白
	parts := strings.Fields(string(words))
	var cleanWords []string
	for _, part := range parts {
		if part != "" {
			cleanWords = append(cleanWords, strings.ToLower(part))
		}
	}

	return cleanWords
}

// capitalize 首字母大写
func capitalize(str string) string {
	if str == "" {
		return ""
	}

	runes := []rune(str)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

// Uncapitalize 首字母小写
func Uncapitalize(str string) string {
	if str == "" {
		return ""
	}

	runes := []rune(str)
	runes[0] = unicode.ToLower(runes[0])
	return string(runes)
}

// IsEmpty 检查字符串是否为空
func IsEmpty(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsNotEmpty 检查字符串是否不为空
func IsNotEmpty(str string) bool {
	return !IsEmpty(str)
}

// DefaultString 如果字符串为空则返回默认值
func DefaultString(str, defaultValue string) string {
	if IsEmpty(str) {
		return defaultValue
	}
	return str
}

// RemovePrefix 移除前缀
func RemovePrefix(str, prefix string) string {
	if strings.HasPrefix(str, prefix) {
		return str[len(prefix):]
	}
	return str
}

// RemoveSuffix 移除后缀
func RemoveSuffix(str, suffix string) string {
	if strings.HasSuffix(str, suffix) {
		return str[:len(str)-len(suffix)]
	}
	return str
}

// Pluralize 单词复数化（简单实现）
func Pluralize(word string) string {
	if word == "" {
		return ""
	}

	word = strings.ToLower(word)

	// 特殊情况
	irregulars := map[string]string{
		"child":  "children",
		"person": "people",
		"man":    "men",
		"woman":  "women",
		"tooth":  "teeth",
		"foot":   "feet",
		"mouse":  "mice",
		"goose":  "geese",
	}

	if plural, exists := irregulars[word]; exists {
		return plural
	}

	// 一般规则
	if strings.HasSuffix(word, "y") && !isVowel(word[len(word)-2]) {
		return word[:len(word)-1] + "ies"
	}

	if strings.HasSuffix(word, "s") || strings.HasSuffix(word, "x") ||
		strings.HasSuffix(word, "z") || strings.HasSuffix(word, "ch") ||
		strings.HasSuffix(word, "sh") {
		return word + "es"
	}

	if strings.HasSuffix(word, "f") {
		return word[:len(word)-1] + "ves"
	}

	if strings.HasSuffix(word, "fe") {
		return word[:len(word)-2] + "ves"
	}

	return word + "s"
}

// isVowel 检查字符是否为元音
func isVowel(b byte) bool {
	vowels := "aeiouAEIOU"
	return strings.ContainsRune(vowels, rune(b))
}

// GeneratePermissions 生成权限字符串
func GeneratePermissions(moduleName, businessName string) []string {
	operations := []string{"list", "add", "edit", "delete", "view"}
	permissions := make([]string, len(operations))

	for i, op := range operations {
		permissions[i] = strings.ToLower(moduleName) + ":" + strings.ToLower(businessName) + ":" + op
	}

	return permissions
}
