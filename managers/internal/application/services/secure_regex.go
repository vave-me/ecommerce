package services // Assuming it's part of the same domain package

import (
	"context" // Needed for context-aware operations potentially
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	// "github.com/stackus/errors" // If using this for errors
	stdErrors "errors"
)

// Secure regex constants
const (
	DefaultRegexTimeout          = 200 * time.Millisecond // Increased default timeout
	MaxRegexInputLengthDefault   = 10240                  // Default max input length (10KB)
	MaxRegexPatternLengthDefault = 2048                   // Default max pattern length (2KB)
	DefaultMaxRegexMatches       = 100                    // Default limit for FindAllStringSubmatchSecure
	logPrefixSecureRegex         = "[SecureRegex]"
)

// Secure regex errors
var (
	ErrRegexTimeout          = stdErrors.New("secure regex: operation timed out")
	ErrRegexInputTooLong     = stdErrors.New("secure regex: input too long for operation")
	ErrRegexPatternTooLong   = stdErrors.New("secure regex: pattern too long")
	ErrRegexUnsafePattern    = stdErrors.New("secure regex: potentially unsafe pattern detected")
	ErrRegexCompilationPanic = stdErrors.New("secure regex: panic during compilation")
	ErrRegexExecutionPanic   = stdErrors.New("secure regex: panic during execution")
)

// SecureRegexEngine provides timeout-protected regex operations and basic ReDoS prevention.
type SecureRegexEngine struct {
	timeout          time.Duration
	maxInputLength   int
	maxPatternLength int
	// TODO: Consider adding a cache for compiled regex patterns if performance is critical
	// compiledPatternCache map[string]*regexp.Regexp
	// cacheMutex           sync.RWMutex
}

// NewSecureRegexEngine creates a new secure regex engine with default settings.
func NewSecureRegexEngine() *SecureRegexEngine {
	return &SecureRegexEngine{
		timeout:          DefaultRegexTimeout,
		maxInputLength:   MaxRegexInputLengthDefault,
		maxPatternLength: MaxRegexPatternLengthDefault,
	}
}

// NewSecureRegexEngineWithConfig creates a secure regex engine with custom configuration.
func NewSecureRegexEngineWithConfig(timeout time.Duration, maxInput int, maxPattern int) *SecureRegexEngine {
	if timeout <= 0 {
		timeout = DefaultRegexTimeout
	}
	if maxInput <= 0 {
		maxInput = MaxRegexInputLengthDefault
	}
	if maxPattern <= 0 {
		maxPattern = MaxRegexPatternLengthDefault
	}
	return &SecureRegexEngine{
		timeout:          timeout,
		maxInputLength:   maxInput,
		maxPatternLength: maxPattern,
	}
}

// CompileSecure compiles a regex pattern with safety checks and timeout.
func (sre *SecureRegexEngine) CompileSecure(pattern string) (re *regexp.Regexp, err error) {
	if len(pattern) > sre.maxPatternLength {
		return nil, fmt.Errorf("%w: length %d, max %d", ErrRegexPatternTooLong, len(pattern), sre.maxPatternLength)
	}
	if err := sre.validatePatternHeuristics(pattern); err != nil {
		return nil, err // ErrRegexUnsafePattern
	}

	// Use a channel to get result from goroutine, enabling timeout for compilation
	resultChan := make(chan *regexp.Regexp, 1)
	errChan := make(chan error, 1)

	compileCtx, cancel := context.WithTimeout(context.Background(), sre.timeout) // Apply timeout to compilation itself
	defer cancel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("%s CRITICAL: Panic during regex compilation for pattern '%s': %v", logPrefixSecureRegex, pattern, r)
				select {
				case errChan <- fmt.Errorf("%w: %v", ErrRegexCompilationPanic, r):
				case <-compileCtx.Done(): // Avoid blocking if context already done
				}
			}
		}()
		compiledRe, compileErr := regexp.Compile(pattern)
		if compileErr != nil {
			select {
			case errChan <- compileErr:
			case <-compileCtx.Done():
			}
			return
		}
		select {
		case resultChan <- compiledRe:
		case <-compileCtx.Done():
		}
	}()

	select {
	case re = <-resultChan:
		log.Printf("%s INFO: Successfully compiled regex pattern (first 50 chars): '%.50s'", logPrefixSecureRegex, pattern)
		return re, nil
	case err = <-errChan:
		log.Printf("%s ERROR: Failed to compile regex pattern '%.50s': %v", logPrefixSecureRegex, pattern, err)
		return nil, err
	case <-compileCtx.Done(): // Handles sre.timeout for compilation
		log.Printf("%s ERROR: Regex compilation timed out for pattern '%.50s' after %s", logPrefixSecureRegex, pattern, sre.timeout)
		return nil, ErrRegexTimeout
	}
}

// executeRegexOperation safely executes a regex function that returns a result and an error.
// It handles input length validation and timeout.
func (sre *SecureRegexEngine) executeRegexOperation(
	input string,
	operationName string,
	opFunc func() (interface{}, error),
) (result interface{}, err error) {
	if len(input) > sre.maxInputLength {
		return nil, fmt.Errorf("%w: operation '%s' input length %d, max %d", ErrRegexInputTooLong, operationName, len(input), sre.maxInputLength)
	}

	resultChan := make(chan interface{}, 1)
	errChan := make(chan error, 1)

	opCtx, cancel := context.WithTimeout(context.Background(), sre.timeout)
	defer cancel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("%s CRITICAL: Panic during regex operation '%s': %v", logPrefixSecureRegex, operationName, r)
				select {
				case errChan <- fmt.Errorf("%w: operation '%s': %v", ErrRegexExecutionPanic, operationName, r):
				case <-opCtx.Done():
				}
			}
		}()
		res, opErr := opFunc()
		if opErr != nil {
			select {
			case errChan <- opErr:
			case <-opCtx.Done():
			}
			return
		}
		select {
		case resultChan <- res:
		case <-opCtx.Done():
		}
	}()

	select {
	case result = <-resultChan:
		return result, nil
	case err = <-errChan:
		return nil, err
	case <-opCtx.Done():
		log.Printf("%s ERROR: Regex operation '%s' timed out after %s", logPrefixSecureRegex, operationName, sre.timeout)
		return nil, ErrRegexTimeout
	}
}

// FindStringSubmatchSecure performs regex matching with timeout and input validation.
func (sre *SecureRegexEngine) FindStringSubmatchSecure(re *regexp.Regexp, input string) ([]string, error) {
	if re == nil {
		return nil, fmt.Errorf("regex is nil")
	}
	res, err := sre.executeRegexOperation(input, "FindStringSubmatch", func() (interface{}, error) {
		return re.FindStringSubmatch(input), nil
	})
	if err != nil {
		return nil, err
	}
	return res.([]string), nil
}

// FindAllStringSubmatchSecure performs regex matching with timeout, input validation, and match limit.
func (sre *SecureRegexEngine) FindAllStringSubmatchSecure(re *regexp.Regexp, input string, n int) ([][]string, error) {
	if re == nil {
		return nil, fmt.Errorf("regex is nil")
	}
	// Apply a sensible default/max for n if it's too large or negative
	if n < 0 || n > DefaultMaxRegexMatches {
		log.Printf("%s WARN: FindAllStringSubmatchSecure 'n' was %d, adjusted to %d", logPrefixSecureRegex, n, DefaultMaxRegexMatches)
		n = DefaultMaxRegexMatches
	}
	res, err := sre.executeRegexOperation(input, "FindAllStringSubmatch", func() (interface{}, error) {
		return re.FindAllStringSubmatch(input, n), nil
	})
	if err != nil {
		return nil, err
	}
	return res.([][]string), nil
}

// MatchStringSecure performs regex matching with timeout and input validation.
func (sre *SecureRegexEngine) MatchStringSecure(re *regexp.Regexp, input string) (bool, error) {
	if re == nil {
		return false, fmt.Errorf("regex is nil")
	}
	res, err := sre.executeRegexOperation(input, "MatchString", func() (interface{}, error) {
		return re.MatchString(input), nil
	})
	if err != nil {
		return false, err
	}
	return res.(bool), nil
}

// ReplaceAllStringSecure performs regex replacement with timeout protection.
func (sre *SecureRegexEngine) ReplaceAllStringSecure(re *regexp.Regexp, input, replacement string) (string, error) {
	if re == nil {
		return "", fmt.Errorf("regex is nil")
	}
	// Input length check for 'input' string. Replacement string length is usually not a ReDoS vector.
	if len(input) > sre.maxInputLength {
		return "", fmt.Errorf("%w: operation 'ReplaceAllString' input length %d, max %d", ErrRegexInputTooLong, len(input), sre.maxInputLength)
	}

	resultChan := make(chan string, 1)
	errChan := make(chan error, 1)

	opCtx, cancel := context.WithTimeout(context.Background(), sre.timeout)
	defer cancel()

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("%s CRITICAL: Panic during regex operation 'ReplaceAllString': %v", logPrefixSecureRegex, r)
				select {
				case errChan <- fmt.Errorf("%w: operation 'ReplaceAllString': %v", ErrRegexExecutionPanic, r):
				case <-opCtx.Done():
				}
			}
		}()
		replacedStr := re.ReplaceAllString(input, replacement)
		select {
		case resultChan <- replacedStr:
		case <-opCtx.Done():
		}
	}()

	select {
	case resStr := <-resultChan:
		return resStr, nil
	case err := <-errChan:
		return "", err
	case <-opCtx.Done():
		log.Printf("%s ERROR: Regex operation 'ReplaceAllString' timed out after %s", logPrefixSecureRegex, sre.timeout)
		return "", ErrRegexTimeout
	}
}

// validatePatternHeuristics checks for some common ReDoS-vulnerable patterns.
// This is not exhaustive and real ReDoS static analysis is complex.
func (sre *SecureRegexEngine) validatePatternHeuristics(pattern string) error {
	// Basic checks for patterns known to be catastrophically bad (from original + refinements)
	// (A*)*, (A+)+, (A?)*, (A?) +
	// ([a-zA-Z]+)*
	// (a|aa)+
	// (a|b|c|d|e|f|g|h|i|j|k|l|m|n|o|p|q|r|s|t|u|v|w|x|y|z)*
	// These are examples, a more comprehensive list or library would be better.
	// Refined checks for nested quantifiers:
	// Look for patterns like `(.*\w*)*`, `(a+)+`, `(a*)*`, `(a?)*`, `(a{M,N})*` (where M,N are numbers)
	// or `(a*)+`, `(a?)+`, `(a{M,N})+`
	// A simple check: multiple Kleene stars or pluses separated only by simple groups or optionals.
	if matched, _ := regexp.MatchString(`\((.*[*+?\{\d,}]){2,}\)[*+?\{\d,}]`, pattern); matched {
		// This regex looks for something like `((...*)...*)*` or `((...+)...+)+`
		// Example: `(a*b*)*` or `(a+b+)+`
		// This is a heuristic for "evil regex" patterns involving nested quantifiers on potentially overlapping expressions.
		if strings.Count(pattern, "*")+strings.Count(pattern, "+") > 3 && strings.Count(pattern, "(") > 2 { // Very basic heuristic
			log.Printf("%s WARN: Pattern '%.50s...' has multiple nested-like quantifiers, potential ReDoS risk.", logPrefixSecureRegex, pattern)
			return fmt.Errorf("%w: pattern contains potentially problematic nested quantifiers", ErrRegexUnsafePattern)
		}
	}

	// Alternation with repeated groups (e.g., (a|b)* problematic if a and b can match same char)
	// (X|Y)*Z where X and Y can match the same prefix
	if strings.Count(pattern, "|") > 0 && (strings.Count(pattern, "*") > 0 || strings.Count(pattern, "+") > 0) {
		// This is a very broad heuristic and can have false positives.
		// A real static analyzer for regex safety is needed for true ReDoS prevention.
		log.Printf("%s WARN: Pattern '%.50s...' uses alternation with quantifiers. Review for ReDoS potential.", logPrefixSecureRegex, pattern)
	}

	// Discourage overly complex patterns based on length or too many groups/quantifiers if not using a static analyzer.
	if len(pattern) > 256 && (strings.Count(pattern, "*") > 5 || strings.Count(pattern, "+") > 5 || strings.Count(pattern, "?") > 10 || strings.Count(pattern, "(") > 10) {
		log.Printf("%s WARN: Pattern '%.50s...' is long and complex. Review for ReDoS potential.", logPrefixSecureRegex, pattern)
		// return fmt.Errorf("%w: pattern is overly complex and may lead to ReDoS", ErrRegexUnsafePattern) // Optionally make this an error
	}

	return nil
}

// --- SecureExtractors (Depends on SecureRegexEngine and MultiLanguageSecurityFilter) ---

// SecureExtractors provides secure implementations of various extraction functions.
type SecureExtractors struct {
	regexEngine *SecureRegexEngine
	textFilter  *MultiLanguageSecurityFilter // Assuming this filter handles sanitization & prohibited content
}

// NewSecureExtractors creates a new secure extractors instance.
func NewSecureExtractors(engine *SecureRegexEngine, filter *MultiLanguageSecurityFilter) (*SecureExtractors, error) {
	if engine == nil {
		return nil, fmt.Errorf("SecureRegexEngine cannot be nil")
	}
	if filter == nil {
		return nil, fmt.Errorf("MultiLanguageSecurityFilter cannot be nil")
	}
	return &SecureExtractors{
		regexEngine: engine,
		textFilter:  filter,
	}, nil
}

// ExtractPageNumberSecure securely extracts page number.
func (se *SecureExtractors) ExtractPageNumberSecure(message string) (string, error) {
	sanitizedMessage, err := se.textFilter.SanitizeAndValidateInput(message, MaxRegexInputLengthDefault) // Use a reasonable length
	if err != nil {
		return "", fmt.Errorf("input sanitization/validation failed for page number extraction: %w", err)
	}

	// Prefer simple string logic if possible to avoid regex overhead/risk
	words := strings.Fields(strings.ToLower(sanitizedMessage))
	for i, word := range words {
		if (word == "page" || word == "p." || word == "go" || word == "goto") && i+1 < len(words) {
			numStr := words[i+1]
			// Validate if it's purely numeric and within a sensible range
			if pageNum, err := strconv.Atoi(numStr); err == nil && pageNum > 0 && pageNum <= 10000 { // Max 10k pages
				return numStr, nil
			}
		}
	}

	// Fallback to a safe regex if simple parsing fails
	// Pattern: page (optional space) number (1-4 digits)
	pattern := `(?i)\b(?:page|pg\.?|p\.?|go\s*to)\s+(\d{1,4})\b`
	re, err := se.regexEngine.CompileSecure(pattern)
	if err != nil { // Should be rare if pattern is fixed
		log.Printf("%s ERROR: Failed to compile secure regex for page number: %v", logPrefixSecureRegex, err)
		return "", fmt.Errorf("internal regex error for page number: %w", err)
	}

	matches, err := se.regexEngine.FindStringSubmatchSecure(re, sanitizedMessage)
	if err != nil { // Timeout or execution error
		return "", fmt.Errorf("regex execution failed for page number: %w", err)
	}
	if len(matches) > 1 {
		// Additional validation on the matched number if needed (already done by regex \d{1,4})
		return matches[1], nil
	}
	return "", nil // No page number found
}

// ExtractPageSizeSecure securely extracts page size.
func (se *SecureExtractors) ExtractPageSizeSecure(message string) (string, error) {
	sanitizedMessage, err := se.textFilter.SanitizeAndValidateInput(message, MaxRegexInputLengthDefault)
	if err != nil {
		return "", fmt.Errorf("input sanitization/validation failed for page size extraction: %w", err)
	}

	// Regex patterns, ordered from more specific to more general
	patterns := []string{
		`(?i)\b(?:show|display|limit)\s+(\d{1,3})\s+(?:items|results|entries|records)\b`, // "show 20 items"
		`(?i)\b(\d{1,3})\s+per\s+page\b`,                                                 // "20 per page"
		`(?i)\bpage(?:-|_)size\s*[:=]?\s*(\d{1,3})\b`,                                    // "page_size: 20"
		`(?i)\blimit\s*[:=]?\s*(\d{1,3})\b`,                                              // "limit: 20"
	}
	maxPageSize := 200 // Sensible max page size

	for _, pattern := range patterns {
		re, compileErr := se.regexEngine.CompileSecure(pattern)
		if compileErr != nil {
			log.Printf("%s ERROR: Failed to compile secure regex for page size ('%s'): %v", logPrefixSecureRegex, pattern, compileErr)
			continue // Try next pattern
		}
		matches, execErr := se.regexEngine.FindStringSubmatchSecure(re, sanitizedMessage)
		if execErr != nil { // Timeout or execution error
			log.Printf("%s WARN: Regex execution error for page size ('%s'): %v", logPrefixSecureRegex, pattern, execErr)
			continue
		}
		if len(matches) > 1 {
			numStr := matches[1]
			if pageSize, errConv := strconv.Atoi(numStr); errConv == nil && pageSize > 0 && pageSize <= maxPageSize {
				return numStr, nil
			}
		}
	}
	return "", nil // No page size found
}

// ExtractSearchTermSecure securely extracts search terms.
func (se *SecureExtractors) ExtractSearchTermSecure(message string) (string, error) {
	sanitizedMessage, err := se.textFilter.SanitizeAndValidateInput(message, MaxRegexInputLengthDefault)
	if err != nil {
		return "", fmt.Errorf("input sanitization/validation failed for search term extraction: %w", err)
	}
	lowerMessage := strings.ToLower(sanitizedMessage)

	// More robust extraction: look for keywords then take subsequent text,
	// stopping at common delimiters or sentence end.
	// This version improves upon the original's simple prefix check.
	// Patterns: "search for X [delimiter]", "find X [delimiter]", "query: X"
	// Using a regex that captures content after specific keywords.
	// This regex is an example and would need refinement for various phrasings.
	// It looks for "search for/find/query" followed by captured group, until a preposition or too many words.
	pattern := `(?i)(?:search\s*(?:for)?|find|query\s*(?:for|is|:)?)[:\s]*((?:[\w'-]+\s*){1,10})`
	// The `((?:[\w'-]+\s*){1,10})` part captures 1 to 10 "words" (alphanumeric, hyphen, apostrophe).

	re, compileErr := se.regexEngine.CompileSecure(pattern)
	if compileErr != nil {
		log.Printf("%s ERROR: Failed to compile search term regex: %v", logPrefixSecureRegex, compileErr)
		return "", fmt.Errorf("internal regex error for search term: %w", compileErr)
	}

	matches, execErr := se.regexEngine.FindStringSubmatchSecure(re, sanitizedMessage) // Use sanitizedMessage
	if execErr != nil {
		return "", fmt.Errorf("regex execution failed for search term: %w", execErr)
	}

	if len(matches) > 1 {
		term := strings.TrimSpace(matches[1])  // The captured group
		if len(term) > 2 && len(term) <= 150 { // Min/max length for a reasonable search term
			log.Printf("%s INFO: Extracted search term: '%s'", logPrefixSecureRegex, term)
			return term, nil
		}
	}

	// Fallback: if no keyword pattern matches, but message is short, consider it all as search term.
	if len(strings.Fields(lowerMessage)) <= 10 && len(lowerMessage) > 2 && len(lowerMessage) <= 150 {
		// Avoid using very short or very long messages as default search terms without keywords.
		// Check if it looks like a question or command instead of a search term.
		if !strings.HasPrefix(lowerMessage, "what") && !strings.HasPrefix(lowerMessage, "how") && !strings.HasPrefix(lowerMessage, "is ") && !strings.HasPrefix(lowerMessage, "can ") {
			log.Printf("%s INFO: Using fallback search term (short message): '%s'", logPrefixSecureRegex, sanitizedMessage)
			return sanitizedMessage, nil
		}
	}

	return "", nil // No search term extracted
}

// ExtractPriceRangeSecure securely extracts min and max price.
func (se *SecureExtractors) ExtractPriceRangeSecure(message string) (minPrice, maxPrice int64, err error) {
	sanitizedMessage, err := se.textFilter.SanitizeAndValidateInput(message, MaxRegexInputLengthDefault)
	if err != nil {
		return 0, 0, fmt.Errorf("input sanitization/validation failed for price range: %w", err)
	}

	// Regex for currency amount (dollars, euros, pounds, yen - can be expanded)
	// Allows for optional cents, commas for thousands.
	// Captures the numeric part.
	currencyAmountPattern := `(?:[$€£¥]?\s*)(\d{1,7}(?:[,.]\d{3})*(?:\.\d{1,2})?)` // Up to 9,999,999.99

	// Pattern 1: "between X and Y" or "X to Y"
	// e.g., "between $100 and $200", "1000-5000 euros"
	rangePattern := fmt.Sprintf(`(?i)\b(?:between\s+%s\s+and\s+%s|%s\s*(?:to|-)\s*%s)\b`,
		currencyAmountPattern, currencyAmountPattern, currencyAmountPattern, currencyAmountPattern)

	reRange, compileErr := se.regexEngine.CompileSecure(rangePattern)
	if compileErr != nil {
		return 0, 0, fmt.Errorf("internal regex error (range): %w", compileErr)
	}

	matches, execErr := se.regexEngine.FindStringSubmatchSecure(reRange, sanitizedMessage)
	if execErr != nil {
		return 0, 0, fmt.Errorf("regex execution error (range): %w", execErr)
	}

	if len(matches) == 5 { // (full match, cur1_val1, cur1_val2, cur2_val1, cur2_val2) -> pattern structure dependent
		// This regex structure suggests matches[1] and matches[2] or similar if structured carefully
		// For: between val1 and val2 -> matches[1] is val1, matches[2] is val2
		// For: val1 to val2 -> matches[3] is val1, matches[4] is val2
		val1Str := matches[1]
		val2Str := matches[2]
		if val1Str == "" && val2Str == "" { // Matched the second part of OR
			val1Str = matches[3]
			val2Str = matches[4]
		}

		p1, e1 := parsePriceToCents(val1Str)
		p2, e2 := parsePriceToCents(val2Str)
		if e1 == nil && e2 == nil {
			return p1, p2, nil
		}
		log.Printf("%s WARN: Could not parse range prices '%s', '%s'", logPrefixSecureRegex, val1Str, val2Str)
	}

	// Pattern 2: "under/less than/below X" or "over/more than/above X"
	// e.g., "under $500", "more than 200 EUR"
	underPattern := fmt.Sprintf(`(?i)\b(?:under|less\s+than|below|up\s+to|max(?:imum)?)\s+%s\b`, currencyAmountPattern)
	overPattern := fmt.Sprintf(`(?i)\b(?:over|more\s+than|above|min(?:imum)?)\s+%s\b`, currencyAmountPattern)

	reUnder, _ := se.regexEngine.CompileSecure(underPattern) // Ignore compile error for brevity, handle in prod
	reOver, _ := se.regexEngine.CompileSecure(overPattern)

	underMatches, _ := se.regexEngine.FindStringSubmatchSecure(reUnder, sanitizedMessage)
	if len(underMatches) > 1 {
		p, e := parsePriceToCents(underMatches[1])
		if e == nil {
			maxPrice = p
		}
	}

	overMatches, _ := se.regexEngine.FindStringSubmatchSecure(reOver, sanitizedMessage)
	if len(overMatches) > 1 {
		p, e := parsePriceToCents(overMatches[1])
		if e == nil {
			minPrice = p
		}
	}

	// Pattern 3: Single price mentioned (e.g., "around $X", "for X dollars")
	// This is ambiguous for a range, but could set min and max to the same if intent is "exactly X"
	if minPrice == 0 && maxPrice == 0 {
		singlePricePattern := fmt.Sprintf(`(?i)\b(?:around|about|for|price(?:d\s+at)?)\s+%s\b`, currencyAmountPattern)
		reSingle, _ := se.regexEngine.CompileSecure(singlePricePattern)
		singleMatches, _ := se.regexEngine.FindStringSubmatchSecure(reSingle, sanitizedMessage)
		if len(singleMatches) > 1 {
			p, e := parsePriceToCents(singleMatches[1])
			if e == nil {
				minPrice = p
				maxPrice = p
			}
		}
	}

	if minPrice > 0 || maxPrice > 0 {
		return minPrice, maxPrice, nil
	}

	return 0, 0, nil // No price range found
}

// parsePriceToCents converts a string price (e.g., "1,200.50") to cents (int64).
func parsePriceToCents(priceStr string) (int64, error) {
	if priceStr == "" {
		return 0, fmt.Errorf("empty price string")
	}
	// Remove currency symbols, commas
	cleanedStr := strings.ReplaceAll(priceStr, ",", "")
	cleanedStr = strings.TrimLeft(cleanedStr, "$€£¥ ") // Common currency symbols

	// Check for cents / decimal part
	var major, minor int64
	var err error
	if strings.Contains(cleanedStr, ".") {
		parts := strings.Split(cleanedStr, ".")
		major, err = strconv.ParseInt(parts[0], 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid major price part '%s': %w", parts[0], err)
		}
		if len(parts[1]) > 2 {
			parts[1] = parts[1][:2]
		} // Max 2 decimal places for cents
		minorStr := parts[1]
		if len(minorStr) == 1 {
			minorStr += "0"
		} // e.g. .5 -> .50
		minor, err = strconv.ParseInt(minorStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid minor price part '%s': %w", parts[1], err)
		}
	} else {
		major, err = strconv.ParseInt(cleanedStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid price value '%s': %w", cleanedStr, err)
		}
	}

	totalCents := major*100 + minor
	if totalCents < 0 || totalCents > 100000000000 { // Validate range (e.g., max 1 billion dollars in cents)
		return 0, fmt.Errorf("price value %d cents out of sensible range", totalCents)
	}
	return totalCents, nil
}

// ExtractVehicleModelSecure securely extracts vehicle model using a predefined list.
func (se *SecureExtractors) ExtractVehicleModelSecure(message string) (string, error) {
	sanitizedMessage, err := se.textFilter.SanitizeAndValidateInput(message, MaxRegexInputLengthDefault)
	if err != nil {
		return "", fmt.Errorf("input sanitization/validation failed for vehicle model: %w", err)
	}
	lowerMessage := strings.ToLower(sanitizedMessage)

	// TODO: This list should be configurable and much larger for production.
	// Consider sourcing from a database or configuration file.
	knownModels := []string{
		"civic", "camry", "corolla", "accord", "rav4", "cr-v", "escape", "explorer",
		"f-150", "silverado", "ram 1500", "sierra", "tacoma", "tundra",
		"model s", "model 3", "model x", "model y", "id.4", "ioniq 5", "ev6",
		"mustang", "challenger", "charger", "camaro",
		"wrangler", "cherokee", "grand cherokee", "forester", "outback",
		"3 series", "5 series", "c-class", "e-class", "a4", "a6", "x3", "x5", "gle", "glc",
		// Add more models, including multi-word models
	}

	// For multi-word models, ensure they are checked effectively.
	// Sorting by length (desc) can help match longer phrases first.
	// For simplicity here, direct iteration.
	var foundModel string
	for _, model := range knownModels {
		// Use word boundaries for more precise matching if model is a distinct word
		// For multi-word models like "f-150" or "model s", direct Contains is okay.
		// Using a simple strings.Contains for this example.
		if strings.Contains(lowerMessage, strings.ToLower(model)) {
			// If multiple models match, this will pick the first one in the list.
			// More sophisticated logic could score matches or look for surrounding context.
			if len(model) > len(foundModel) { // Prefer longer match if multiple sub-strings match
				foundModel = model
			}
		}
	}
	if foundModel != "" {
		log.Printf("%s INFO: Extracted vehicle model: '%s'", logPrefixSecureRegex, foundModel)
		return foundModel, nil
	}

	return "", nil // No model found
}
