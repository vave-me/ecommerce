package services

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/text/language"
	"golang.org/x/text/unicode/norm"
)

// Security constants
const (
	MaxInputLength         = 10000 // Maximum input message length
	MaxContextSize         = 100   // Maximum context entries
	MaxSecurityFieldLength = 1000  // Maximum field length for security validation
	MixedScriptLimit       = 2     // Maximum allowed mixed scripts
)

// Security errors
var (
	ErrInputTooLong             = fmt.Errorf("input exceeds maximum length")
	ErrContextTooLarge          = fmt.Errorf("context exceeds maximum size")
	ErrSecurityMaliciousContent = fmt.Errorf("malicious content detected")
	ErrUnicodeAttack            = fmt.Errorf("unicode attack detected")
	ErrHomographAttack          = fmt.Errorf("homograph attack detected")
)

// MultiLanguageSecurityFilter provides comprehensive multi-language security filtering
type MultiLanguageSecurityFilter struct {
	prohibitedPatterns map[language.Tag][]string
	normalizer         norm.Form
	homographDetector  *HomographDetector
}

// HomographDetector detects homograph and script mixing attacks
type HomographDetector struct {
	suspiciousScripts map[string]bool
	mixedScriptLimit  int
}

// NewMultiLanguageSecurityFilter creates a new comprehensive security filter
func NewMultiLanguageSecurityFilter() *MultiLanguageSecurityFilter {
	return &MultiLanguageSecurityFilter{
		prohibitedPatterns: map[language.Tag][]string{
			language.English: {
				"admin", "administrator", "delete all", "drop table", "system", "root",
				"exec", "execute", "script", "eval", "function", "constructor",
				"__proto__", "prototype", "import", "require", "process",
				"rm -rf", "sudo", "chmod", "chown", "passwd", "shadow",
			},
			language.Spanish: {
				"administrador", "eliminar todo", "borrar todo", "sistema", "raíz",
				"ejecutar", "tabla", "función", "proceso", "eliminar",
			},
			language.French: {
				"administrateur", "supprimer tout", "effacer tout", "système", "racine",
				"exécuter", "table", "fonction", "processus", "supprimer",
			},
			language.German: {
				"administrator", "alle löschen", "alles löschen", "system", "wurzel",
				"ausführen", "tabelle", "funktion", "prozess", "löschen",
			},
			language.Russian: {
				"админ", "администратор", "удалить все", "стереть все", "система", "корень",
				"выполнить", "таблица", "функция", "процесс", "удалить",
			},
			language.Chinese: {
				"管理员", "删除所有", "清除所有", "系统", "根", "执行",
				"表格", "函数", "进程", "管理", "删除", "清除",
			},
			language.Arabic: {
				"مدير", "مدير النظام", "احذف الكل", "مسح الكل", "نظام", "جذر",
				"تنفيذ", "جدول", "وظيفة", "عملية", "حذف",
			},
			language.Japanese: {
				"管理者", "すべて削除", "全削除", "システム", "ルート", "実行",
				"テーブル", "関数", "プロセス", "削除",
			},
			language.Korean: {
				"관리자", "모두 삭제", "전체 삭제", "시스템", "루트", "실행",
				"테이블", "함수", "프로세스", "삭제",
			},
		},
		normalizer:        norm.NFKC, // Canonical decomposition + compatibility
		homographDetector: NewHomographDetector(),
	}
}

// NewHomographDetector creates a new homograph detector
func NewHomographDetector() *HomographDetector {
	return &HomographDetector{
		suspiciousScripts: map[string]bool{
			"Cyrillic": true,
			"Greek":    true,
			"Arabic":   true,
			"Hebrew":   true,
		},
		mixedScriptLimit: MixedScriptLimit,
	}
}

// ContainsProhibitedContent checks for prohibited content in multiple languages
func (f *MultiLanguageSecurityFilter) ContainsProhibitedContent(input string) (bool, string) {
	// Input validation
	if len(input) > MaxInputLength {
		return true, fmt.Sprintf("Input too long: %d chars", len(input))
	}

	// UTF-8 validation
	if !utf8.ValidString(input) {
		return true, "Invalid UTF-8 encoding detected"
	}

	// Normalize Unicode to prevent normalization attacks
	normalized := f.normalizer.String(input)

	// Remove zero-width characters and other bypass attempts
	cleaned := f.removeZeroWidthChars(normalized)

	// Check for homograph attacks
	if suspicious, reason := f.homographDetector.DetectHomographAttack(cleaned); suspicious {
		return true, fmt.Sprintf("Homograph attack detected: %s", reason)
	}

	// Check against all language patterns
	for lang, patterns := range f.prohibitedPatterns {
		for _, pattern := range patterns {
			if f.containsPattern(cleaned, pattern) {
				return true, fmt.Sprintf("Prohibited pattern '%s' detected in %s", pattern, lang)
			}
		}
	}

	// Check for encoding bypass attempts
	if f.detectEncodingBypass(cleaned) {
		return true, "Encoding bypass attempt detected"
	}

	return false, ""
}

// removeZeroWidthChars removes zero-width characters that can be used for bypass
func (f *MultiLanguageSecurityFilter) removeZeroWidthChars(input string) string {
	zeroWidthChars := []rune{
		'\u200B', // Zero Width Space
		'\u200C', // Zero Width Non-Joiner
		'\u200D', // Zero Width Joiner
		'\u2060', // Word Joiner
		'\uFEFF', // Zero Width No-Break Space
		'\u180E', // Mongolian Vowel Separator
		'\u061C', // Arabic Letter Mark
		'\u202A', // Left-to-Right Embedding
		'\u202B', // Right-to-Left Embedding
		'\u202C', // Pop Directional Formatting
		'\u202D', // Left-to-Right Override
		'\u202E', // Right-to-Left Override
		'\u2066', // Left-to-Right Isolate
		'\u2067', // Right-to-Left Isolate
		'\u2068', // First Strong Isolate
		'\u2069', // Pop Directional Isolate
	}

	result := input
	for _, char := range zeroWidthChars {
		result = strings.ReplaceAll(result, string(char), "")
	}

	// Remove other control characters
	result = strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t' {
			return -1 // Remove control character
		}
		return r
	}, result)

	return result
}

// containsPattern checks if input contains a pattern (case-insensitive, Unicode-aware)
func (f *MultiLanguageSecurityFilter) containsPattern(input, pattern string) bool {
	inputLower := strings.ToLower(input)
	patternLower := strings.ToLower(pattern)

	// Direct match
	if strings.Contains(inputLower, patternLower) {
		return true
	}

	// Check for character substitutions
	return f.checkCharacterSubstitutions(inputLower, patternLower)
}

// checkCharacterSubstitutions checks for common character substitutions
func (f *MultiLanguageSecurityFilter) checkCharacterSubstitutions(input, pattern string) bool {
	// Common lookalike substitutions
	substitutions := map[rune][]rune{
		'a': {'а', 'α', 'ａ', 'à', 'á', 'â', 'ã', 'ä', 'å'}, // Latin a + variants
		'e': {'е', 'ε', 'ｅ', 'è', 'é', 'ê', 'ë'},           // Latin e + variants
		'o': {'о', 'ο', 'ｏ', 'ò', 'ó', 'ô', 'õ', 'ö'},      // Latin o + variants
		'p': {'р', 'ρ', 'ｐ'},                               // Latin p + variants
		'c': {'с', 'ｃ', 'ç'},                               // Latin c + variants
		'x': {'х', 'χ', 'ｘ'},                               // Latin x + variants
		'y': {'у', 'γ', 'ｙ', 'ý', 'ÿ'},                     // Latin y + variants
		'i': {'і', 'ι', 'ｉ', 'ì', 'í', 'î', 'ï'},           // Latin i + variants
		's': {'ѕ', 'σ', 'ｓ', 'ş', 'š'},                     // Latin s + variants
		'n': {'п', 'η', 'ｎ', 'ñ', 'ň'},                     // Latin n + variants
		'u': {'υ', 'ｕ', 'ù', 'ú', 'û', 'ü'},                // Latin u + variants
		't': {'т', 'τ', 'ｔ'},                               // Latin t + variants
		'r': {'г', 'ρ', 'ｒ'},                               // Latin r + variants
		'm': {'м', 'μ', 'ｍ'},                               // Latin m + variants
		'h': {'һ', 'η', 'ｈ'},                               // Latin h + variants
		'd': {'ԁ', 'δ', 'ｄ'},                               // Latin d + variants
		'l': {'ӏ', 'λ', 'ｌ'},                               // Latin l + variants
	}

	// Generate variations of the pattern with substitutions
	variations := f.generatePatternVariations(pattern, substitutions)

	for _, variation := range variations {
		if strings.Contains(input, variation) {
			return true
		}
	}

	return false
}

// generatePatternVariations generates pattern variations with character substitutions
func (f *MultiLanguageSecurityFilter) generatePatternVariations(pattern string, substitutions map[rune][]rune) []string {
	variations := []string{pattern}
	maxVariations := 100 // Limit to prevent explosion

	for i, char := range pattern {
		if len(variations) >= maxVariations {
			break
		}

		if subs, exists := substitutions[char]; exists {
			newVariations := make([]string, 0)
			for _, variation := range variations {
				if len(newVariations) >= maxVariations {
					break
				}
				for _, sub := range subs {
					if len(newVariations) >= maxVariations {
						break
					}
					newVar := variation[:i] + string(sub) + variation[i+1:]
					newVariations = append(newVariations, newVar)
				}
			}
			variations = append(variations, newVariations...)
		}
	}

	return variations
}

// DetectHomographAttack detects homograph attacks and script mixing
func (h *HomographDetector) DetectHomographAttack(input string) (bool, string) {
	scripts := make(map[string]int)

	for _, r := range input {
		if unicode.IsLetter(r) {
			script := h.getScript(r)
			scripts[script]++
		}
	}

	// Check for mixed scripts
	if len(scripts) > h.mixedScriptLimit {
		scriptNames := make([]string, 0, len(scripts))
		for script := range scripts {
			scriptNames = append(scriptNames, script)
		}
		return true, fmt.Sprintf("Mixed scripts detected: %v", scriptNames)
	}

	// Check for suspicious scripts in Latin-looking text
	if scripts["Latin"] > 0 {
		for script := range scripts {
			if h.suspiciousScripts[script] {
				return true, fmt.Sprintf("Suspicious script '%s' mixed with Latin", script)
			}
		}
	}

	// Check for fullwidth character abuse
	if scripts["Fullwidth"] > 0 && len(scripts) > 1 {
		return true, "Fullwidth character mixing detected"
	}

	return false, ""
}

// getScript returns the script of a Unicode character
func (h *HomographDetector) getScript(r rune) string {
	switch {
	case r >= 0x0041 && r <= 0x005A, r >= 0x0061 && r <= 0x007A:
		return "Latin"
	case r >= 0x0400 && r <= 0x04FF:
		return "Cyrillic"
	case r >= 0x0370 && r <= 0x03FF:
		return "Greek"
	case r >= 0x0590 && r <= 0x05FF:
		return "Hebrew"
	case r >= 0x0600 && r <= 0x06FF:
		return "Arabic"
	case r >= 0x4E00 && r <= 0x9FFF:
		return "CJK"
	case r >= 0xFF01 && r <= 0xFF5E:
		return "Fullwidth"
	case r >= 0x3040 && r <= 0x309F:
		return "Hiragana"
	case r >= 0x30A0 && r <= 0x30FF:
		return "Katakana"
	case r >= 0xAC00 && r <= 0xD7AF:
		return "Hangul"
	default:
		return "Other"
	}
}

// detectEncodingBypass detects various encoding bypass attempts
func (f *MultiLanguageSecurityFilter) detectEncodingBypass(input string) bool {
	// Check for URL encoding patterns
	if strings.Contains(input, "%") {
		// Simple URL encoding detection
		urlEncodedPatterns := []string{"%20", "%2F", "%3C", "%3E", "%22", "%27"}
		for _, pattern := range urlEncodedPatterns {
			if strings.Contains(strings.ToUpper(input), pattern) {
				return true
			}
		}
	}

	// Check for HTML entity encoding
	if strings.Contains(input, "&") && strings.Contains(input, ";") {
		htmlEntities := []string{"&lt;", "&gt;", "&amp;", "&quot;", "&#", "&apos;"}
		for _, entity := range htmlEntities {
			if strings.Contains(strings.ToLower(input), entity) {
				return true
			}
		}
	}

	// Check for hex encoding
	if strings.Contains(input, "\\x") {
		return true
	}

	// Check for unicode escape sequences
	if strings.Contains(input, "\\u") {
		return true
	}

	return false
}

// SanitizeInput sanitizes input by removing dangerous content
func (f *MultiLanguageSecurityFilter) SanitizeInput(input string) string {
	// Normalize Unicode
	normalized := f.normalizer.String(input)

	// Remove zero-width characters
	cleaned := f.removeZeroWidthChars(normalized)

	// Limit length
	if len(cleaned) > MaxInputLength {
		cleaned = cleaned[:MaxInputLength]
	}

	return cleaned
}

// ValidateContextEntry validates a context entry for security
func (f *MultiLanguageSecurityFilter) ValidateContextEntry(key string, value interface{}) error {
	// Key validation
	if len(key) > 100 {
		return fmt.Errorf("context key too long: %d chars", len(key))
	}

	if prohibited, reason := f.ContainsProhibitedContent(key); prohibited {
		return fmt.Errorf("prohibited content in context key: %s", reason)
	}

	// Value validation
	switch v := value.(type) {
	case string:
		if len(v) > MaxSecurityFieldLength {
			return fmt.Errorf("context string value too long: %d chars", len(v))
		}
		if prohibited, reason := f.ContainsProhibitedContent(v); prohibited {
			return fmt.Errorf("prohibited content in context value: %s", reason)
		}
	case map[string]interface{}:
		return fmt.Errorf("nested objects not allowed in context")
	case []interface{}:
		if len(v) > 100 {
			return fmt.Errorf("context array too large: %d items", len(v))
		}
		// Validate array elements
		for i, item := range v {
			if str, ok := item.(string); ok {
				if prohibited, reason := f.ContainsProhibitedContent(str); prohibited {
					return fmt.Errorf("prohibited content in context array item %d: %s", i, reason)
				}
			}
		}
	}

	return nil
}

// SanitizeAndValidateInput combines input sanitization and validation in one call
func (f *MultiLanguageSecurityFilter) SanitizeAndValidateInput(input string, maxLength int) (string, error) {
	// First validate input length
	if len(input) > maxLength {
		return "", fmt.Errorf("input exceeds maximum length of %d characters", maxLength)
	}

	// Check for prohibited content before sanitization
	if hasProhibited, reason := f.ContainsProhibitedContent(input); hasProhibited {
		return "", fmt.Errorf("prohibited content detected: %s", reason)
	}

	// Sanitize the input
	sanitized := f.SanitizeInput(input)

	// Validate the sanitized input doesn't exceed length after sanitization
	if len(sanitized) > maxLength {
		return "", fmt.Errorf("sanitized input exceeds maximum length of %d characters", maxLength)
	}

	return sanitized, nil
}

// FilterThinkingBlocks removes thinking blocks from AI responses before showing to users
func FilterThinkingBlocks(response string) string {
	if response == "" {
		return response
	}

	// Remove thinking blocks using regex
	// This pattern matches <thinking>...</thinking> blocks including nested content
	thinkingPattern := `(?i)<thinking>[\s\S]*?</thinking>`

	// Compile regex with error handling
	re, err := regexp.Compile(thinkingPattern)
	if err != nil {
		// If regex compilation fails, fallback to manual removal
		return manualThinkingBlockRemoval(response)
	}

	// Remove all thinking blocks
	filtered := re.ReplaceAllString(response, "")

	// Clean up extra whitespace and newlines left behind
	filtered = strings.TrimSpace(filtered)
	filtered = regexp.MustCompile(`\n\s*\n\s*\n`).ReplaceAllString(filtered, "\n\n")

	return filtered
}

// manualThinkingBlockRemoval provides fallback thinking block removal
func manualThinkingBlockRemoval(response string) string {
	result := response

	for {
		startPos := strings.Index(strings.ToLower(result), "<thinking>")
		if startPos == -1 {
			break
		}

		// Find the matching closing tag
		endPos := strings.Index(strings.ToLower(result[startPos:]), "</thinking>")
		if endPos == -1 {
			// No closing tag found, remove from start position to end
			result = result[:startPos]
			break
		}

		// Remove the thinking block
		endPos += startPos + len("</thinking>")
		result = result[:startPos] + result[endPos:]
	}

	return strings.TrimSpace(result)
}

// ValidateAndFilterResponse validates and filters AI responses for user safety
func (f *MultiLanguageSecurityFilter) ValidateAndFilterResponse(response string) (string, error) {
	if response == "" {
		return "", nil
	}

	// First, remove thinking blocks (this should ALWAYS happen)
	filtered := FilterThinkingBlocks(response)

	// Then apply security filtering
	if prohibited, reason := f.ContainsProhibitedContent(filtered); prohibited {
		return "", fmt.Errorf("response contains prohibited content: %s", reason)
	}

	// Apply input sanitization to clean up any remaining issues
	sanitized := f.SanitizeInput(filtered)

	return sanitized, nil
}
