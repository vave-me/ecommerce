package services

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// ParameterExtractor consolidates all parameter extraction logic
type ParameterExtractor struct {
	priceRegex   *regexp.Regexp
	idRegex      *regexp.Regexp
	yearRegex    *regexp.Regexp
	mileageRegex *regexp.Regexp
	rangeRegex   *regexp.Regexp
}

// NewParameterExtractor creates a new parameter extractor
func NewParameterExtractor() *ParameterExtractor {
	return &ParameterExtractor{
		priceRegex:   regexp.MustCompile(`\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)`),
		idRegex:      regexp.MustCompile(`(?:id|ID)\s*[:\-]?\s*([a-zA-Z0-9\-_]+)`),
		yearRegex:    regexp.MustCompile(`(?:19|20)\d{2}`),
		mileageRegex: regexp.MustCompile(`(\d+(?:,\d{3})*)\s*(?:miles?|km|kilometers?)`),
		rangeRegex:   regexp.MustCompile(`between\s+(\d+)\s+and\s+(\d+)`),
	}
}

// ExtractParameters extracts all relevant parameters for the given entity type
func (pe *ParameterExtractor) ExtractParameters(message string, context map[string]interface{}, entityType string) map[string]interface{} {
	params := make(map[string]interface{})

	// Common parameters for all entity types
	pe.extractCommonParams(message, context, params)

	// Entity-specific parameters
	switch entityType {
	case "products":
		pe.extractProductParams(message, context, params)
	case "vehicles":
		pe.extractVehicleParams(message, context, params)
	case "properties":
		pe.extractPropertyParams(message, context, params)
	case "jobs":
		pe.extractJobParams(message, context, params)
	case "services":
		pe.extractServiceParams(message, context, params)
	case "deals":
		pe.extractDealsParams(message, context, params)
	}

	return params
}

// extractCommonParams extracts parameters common to all entity types
func (pe *ParameterExtractor) extractCommonParams(message string, context map[string]interface{}, params map[string]interface{}) {
	// Search term
	if searchTerm := pe.extractSearchTerm(message); searchTerm != "" {
		params["search_term"] = searchTerm
	}

	// Price range
	if minPrice, maxPrice := pe.extractPriceRange(message); minPrice > 0 || maxPrice > 0 {
		if minPrice > 0 {
			params["min_price"] = minPrice
		}
		if maxPrice > 0 {
			params["max_price"] = maxPrice
		}
	}

	// Category
	if category := pe.extractCategory(message); category != "" {
		params["category"] = category
	}

	// Location from context
	if lat, lng, radius := pe.extractLocationFromContext(context); lat != 0 && lng != 0 {
		params["lat"] = lat
		params["lng"] = lng
		if radius > 0 {
			params["radius"] = radius
		}
	}

	// Sorting
	if sortBy := pe.extractSortBy(message); sortBy != "" {
		params["sort_by"] = sortBy
	}
	if sortOrder := pe.extractSortOrder(message); sortOrder != "" {
		params["sort_order"] = sortOrder
	}

	// Pagination
	if page := pe.extractPageNumber(message); page > 0 {
		params["page"] = page
	}
	if pageSize := pe.extractPageSize(message); pageSize > 0 {
		params["page_size"] = pageSize
	}

	// User ID from context
	if userID, exists := context["user_id"]; exists {
		params["user_id"] = userID
	}
}

// extractProductParams extracts product-specific parameters
func (pe *ParameterExtractor) extractProductParams(message string, context map[string]interface{}, params map[string]interface{}) {
	// Brand
	if brand := pe.extractBrand(message); brand != "" {
		params["brand"] = brand
	}

	// Condition
	if condition := pe.extractCondition(message); condition != "" {
		params["condition"] = condition
	}

	// Product ID
	if id := pe.extractID(message, "product"); id != "" {
		params["product_id"] = id
	}

	// Size/Weight for products
	if size := pe.extractSize(message); size != "" {
		params["size"] = size
	}
}

// extractVehicleParams extracts vehicle-specific parameters
func (pe *ParameterExtractor) extractVehicleParams(message string, context map[string]interface{}, params map[string]interface{}) {
	// Make and model
	if make := pe.extractVehicleMake(message); make != "" {
		params["make"] = make
	}
	if model := pe.extractVehicleModel(message); model != "" {
		params["model"] = model
	}

	// Year range
	if minYear, maxYear := pe.extractYearRange(message); minYear > 0 || maxYear > 0 {
		if minYear > 0 {
			params["min_year"] = minYear
		}
		if maxYear > 0 {
			params["max_year"] = maxYear
		}
	}

	// Mileage range
	if minMileage, maxMileage := pe.extractMileageRange(message); minMileage > 0 || maxMileage > 0 {
		if minMileage > 0 {
			params["min_mileage"] = minMileage
		}
		if maxMileage > 0 {
			params["max_mileage"] = maxMileage
		}
	}

	// Fuel type
	if fuelType := pe.extractFuelType(message); fuelType != "" {
		params["fuel_type"] = fuelType
	}

	// Transmission
	if transmission := pe.extractTransmission(message); transmission != "" {
		params["transmission"] = transmission
	}
}

// extractPropertyParams extracts property-specific parameters
func (pe *ParameterExtractor) extractPropertyParams(message string, context map[string]interface{}, params map[string]interface{}) {
	// Property type
	if propType := pe.extractPropertyType(message); propType != "" {
		params["property_type"] = propType
	}

	// Bedrooms/bathrooms
	if bedrooms := pe.extractBedrooms(message); bedrooms > 0 {
		params["bedrooms"] = bedrooms
	}
	if bathrooms := pe.extractBathrooms(message); bathrooms > 0 {
		params["bathrooms"] = bathrooms
	}

	// Square footage
	if minSqft, maxSqft := pe.extractSquareFootageRange(message); minSqft > 0 || maxSqft > 0 {
		if minSqft > 0 {
			params["min_sqft"] = minSqft
		}
		if maxSqft > 0 {
			params["max_sqft"] = maxSqft
		}
	}
}

// extractJobParams extracts job-specific parameters
func (pe *ParameterExtractor) extractJobParams(message string, context map[string]interface{}, params map[string]interface{}) {
	// Job type
	if jobType := pe.extractJobType(message); jobType != "" {
		params["job_type"] = jobType
	}

	// Salary range
	if minSalary, maxSalary := pe.extractSalaryRange(message); minSalary > 0 || maxSalary > 0 {
		if minSalary > 0 {
			params["min_salary"] = minSalary
		}
		if maxSalary > 0 {
			params["max_salary"] = maxSalary
		}
	}

	// Remote work
	if pe.isRemoteJob(message) {
		params["remote"] = true
	}

	// Experience level
	if expLevel := pe.extractExperienceLevel(message); expLevel != "" {
		params["experience_level"] = expLevel
	}
}

// extractServiceParams extracts service-specific parameters
func (pe *ParameterExtractor) extractServiceParams(message string, context map[string]interface{}, params map[string]interface{}) {
	// Service type
	if serviceType := pe.extractServiceType(message); serviceType != "" {
		params["service_type"] = serviceType
	}

	// Availability
	if availability := pe.extractAvailability(message); availability != "" {
		params["availability"] = availability
	}
}

// extractDealsParams extracts deal-specific parameters
func (pe *ParameterExtractor) extractDealsParams(message string, context map[string]interface{}, params map[string]interface{}) {
	// Discount percentage
	if discount := pe.extractDiscountPercentage(message); discount > 0 {
		params["min_discount"] = discount
	}

	// Deal type
	if dealType := pe.extractDealType(message); dealType != "" {
		params["deal_type"] = dealType
	}
}

// Core extraction methods
func (pe *ParameterExtractor) extractSearchTerm(message string) string {
	// Remove common stop words and action words
	stopWords := []string{"find", "search", "show", "get", "for", "me", "with", "having", "that", "are", "is"}
	words := strings.Fields(message)

	var cleanWords []string
	for _, word := range words {
		isStopWord := false
		for _, stopWord := range stopWords {
			if strings.EqualFold(word, stopWord) {
				isStopWord = true
				break
			}
		}
		if !isStopWord && len(word) > 2 {
			cleanWords = append(cleanWords, word)
		}
	}

	if len(cleanWords) > 0 {
		return strings.Join(cleanWords, " ")
	}

	return ""
}

func (pe *ParameterExtractor) extractPriceRange(message string) (int64, int64) {
	var minPrice, maxPrice int64

	// Look for "between X and Y" pattern
	matches := pe.rangeRegex.FindStringSubmatch(message)
	if len(matches) == 3 {
		if min, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			minPrice = min
		}
		if max, err := strconv.ParseInt(matches[2], 10, 64); err == nil {
			maxPrice = max
		}
		return minPrice, maxPrice
	}

	// Look for "under X" or "below X"
	if strings.Contains(message, "under") || strings.Contains(message, "below") {
		matches := pe.priceRegex.FindAllString(message, -1)
		if len(matches) > 0 {
			if price, err := strconv.ParseInt(strings.ReplaceAll(matches[0], ",", ""), 10, 64); err == nil {
				maxPrice = price
			}
		}
	}

	// Look for "over X" or "above X"
	if strings.Contains(message, "over") || strings.Contains(message, "above") {
		matches := pe.priceRegex.FindAllString(message, -1)
		if len(matches) > 0 {
			if price, err := strconv.ParseInt(strings.ReplaceAll(matches[0], ",", ""), 10, 64); err == nil {
				minPrice = price
			}
		}
	}

	return minPrice, maxPrice
}

func (pe *ParameterExtractor) extractCategory(message string) string {
	categories := map[string]string{
		"electronics": "electronics",
		"clothing":    "clothing",
		"books":       "books",
		"home":        "home",
		"garden":      "garden",
		"sports":      "sports",
		"automotive":  "automotive",
		"health":      "health",
		"beauty":      "beauty",
		"toys":        "toys",
		"jewelry":     "jewelry",
		"music":       "music",
		"movies":      "movies",
		"games":       "games",
	}

	for keyword, category := range categories {
		if strings.Contains(message, keyword) {
			return category
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractBrand(message string) string {
	brands := []string{
		"apple", "samsung", "sony", "nike", "adidas", "honda", "toyota", "ford",
		"bmw", "microsoft", "google", "amazon", "dell", "hp", "lenovo",
	}

	for _, brand := range brands {
		if strings.Contains(message, brand) {
			return brand
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractCondition(message string) string {
	conditions := map[string]string{
		"new":         "new",
		"used":        "used",
		"like new":    "like_new",
		"excellent":   "excellent",
		"good":        "good",
		"fair":        "fair",
		"poor":        "poor",
		"damaged":     "damaged",
		"refurbished": "refurbished",
	}

	for keyword, condition := range conditions {
		if strings.Contains(message, keyword) {
			return condition
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractVehicleMake(message string) string {
	makes := []string{
		"honda", "toyota", "ford", "chevrolet", "nissan", "bmw", "mercedes",
		"audi", "volkswagen", "hyundai", "kia", "mazda", "subaru", "lexus",
	}

	for _, make := range makes {
		if strings.Contains(message, make) {
			return make
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractVehicleModel(message string) string {
	models := []string{
		"civic", "accord", "camry", "corolla", "f150", "silverado", "3 series",
		"c class", "a4", "jetta", "elantra", "optima", "cx5", "forester",
	}

	for _, model := range models {
		if strings.Contains(message, model) {
			return model
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractYearRange(message string) (int, int) {
	years := pe.yearRegex.FindAllString(message, -1)
	if len(years) == 0 {
		return 0, 0
	}

	var minYear, maxYear int
	for _, yearStr := range years {
		if year, err := strconv.Atoi(yearStr); err == nil {
			if minYear == 0 || year < minYear {
				minYear = year
			}
			if maxYear == 0 || year > maxYear {
				maxYear = year
			}
		}
	}

	// If only one year found, use it as both min and max
	if minYear == maxYear && minYear > 0 {
		return minYear, minYear
	}

	return minYear, maxYear
}

func (pe *ParameterExtractor) extractMileageRange(message string) (int64, int64) {
	matches := pe.mileageRegex.FindAllStringSubmatch(message, -1)
	if len(matches) == 0 {
		return 0, 0
	}

	var minMileage, maxMileage int64
	for _, match := range matches {
		if len(match) > 1 {
			mileageStr := strings.ReplaceAll(match[1], ",", "")
			if mileage, err := strconv.ParseInt(mileageStr, 10, 64); err == nil {
				if minMileage == 0 || mileage < minMileage {
					minMileage = mileage
				}
				if maxMileage == 0 || mileage > maxMileage {
					maxMileage = mileage
				}
			}
		}
	}

	return minMileage, maxMileage
}

func (pe *ParameterExtractor) extractFuelType(message string) string {
	fuelTypes := map[string]string{
		"gas":      "gasoline",
		"gasoline": "gasoline",
		"diesel":   "diesel",
		"electric": "electric",
		"hybrid":   "hybrid",
		"ev":       "electric",
	}

	for keyword, fuelType := range fuelTypes {
		if strings.Contains(message, keyword) {
			return fuelType
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractTransmission(message string) string {
	transmissions := map[string]string{
		"manual":    "manual",
		"automatic": "automatic",
		"cvt":       "cvt",
		"stick":     "manual",
		"auto":      "automatic",
	}

	for keyword, transmission := range transmissions {
		if strings.Contains(message, keyword) {
			return transmission
		}
	}

	return ""
}

// Additional extraction methods for other entity types
func (pe *ParameterExtractor) extractPropertyType(message string) string {
	propertyTypes := map[string]string{
		"house":     "house",
		"apartment": "apartment",
		"condo":     "condo",
		"townhouse": "townhouse",
		"studio":    "studio",
		"loft":      "loft",
	}

	for keyword, propType := range propertyTypes {
		if strings.Contains(message, keyword) {
			return propType
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractBedrooms(message string) int {
	bedroomRegex := regexp.MustCompile(`(\d+)\s*(?:bed|bedroom|br)`)
	matches := bedroomRegex.FindStringSubmatch(message)
	if len(matches) > 1 {
		if bedrooms, err := strconv.Atoi(matches[1]); err == nil {
			return bedrooms
		}
	}
	return 0
}

func (pe *ParameterExtractor) extractBathrooms(message string) int {
	bathroomRegex := regexp.MustCompile(`(\d+(?:\.\d+)?)\s*(?:bath|bathroom|ba)`)
	matches := bathroomRegex.FindStringSubmatch(message)
	if len(matches) > 1 {
		if bathrooms, err := strconv.ParseFloat(matches[1], 64); err == nil {
			return int(bathrooms)
		}
	}
	return 0
}

func (pe *ParameterExtractor) extractSquareFootageRange(message string) (int, int) {
	sqftRegex := regexp.MustCompile(`(\d+(?:,\d{3})*)\s*(?:sq\s*ft|sqft|square\s*feet)`)
	matches := sqftRegex.FindAllStringSubmatch(message, -1)

	var minSqft, maxSqft int
	for _, match := range matches {
		if len(match) > 1 {
			sqftStr := strings.ReplaceAll(match[1], ",", "")
			if sqft, err := strconv.Atoi(sqftStr); err == nil {
				if minSqft == 0 || sqft < minSqft {
					minSqft = sqft
				}
				if maxSqft == 0 || sqft > maxSqft {
					maxSqft = sqft
				}
			}
		}
	}

	return minSqft, maxSqft
}

func (pe *ParameterExtractor) extractJobType(message string) string {
	jobTypes := map[string]string{
		"full time":  "fulltime",
		"part time":  "parttime",
		"contract":   "contract",
		"freelance":  "freelance",
		"temporary":  "temporary",
		"internship": "internship",
	}

	for keyword, jobType := range jobTypes {
		if strings.Contains(message, keyword) {
			return jobType
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractSalaryRange(message string) (int64, int64) {
	salaryRegex := regexp.MustCompile(`\$?(\d{1,3}(?:,\d{3})*(?:\.\d{2})?)\s*(?:k|thousand|per\s*year|annually)?`)
	matches := salaryRegex.FindAllStringSubmatch(message, -1)

	var minSalary, maxSalary int64
	for _, match := range matches {
		if len(match) > 1 {
			salaryStr := strings.ReplaceAll(match[1], ",", "")
			if salary, err := strconv.ParseInt(salaryStr, 10, 64); err == nil {
				// If contains 'k' or 'thousand', multiply by 1000
				if strings.Contains(match[0], "k") || strings.Contains(match[0], "thousand") {
					salary *= 1000
				}

				if minSalary == 0 || salary < minSalary {
					minSalary = salary
				}
				if maxSalary == 0 || salary > maxSalary {
					maxSalary = salary
				}
			}
		}
	}

	return minSalary, maxSalary
}

func (pe *ParameterExtractor) isRemoteJob(message string) bool {
	remoteKeywords := []string{"remote", "work from home", "telecommute", "virtual"}
	for _, keyword := range remoteKeywords {
		if strings.Contains(message, keyword) {
			return true
		}
	}
	return false
}

func (pe *ParameterExtractor) extractExperienceLevel(message string) string {
	levels := map[string]string{
		"entry level": "entry",
		"junior":      "junior",
		"mid level":   "mid",
		"senior":      "senior",
		"lead":        "lead",
		"principal":   "principal",
	}

	for keyword, level := range levels {
		if strings.Contains(message, keyword) {
			return level
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractServiceType(message string) string {
	serviceTypes := map[string]string{
		"cleaning":     "cleaning",
		"repair":       "repair",
		"maintenance":  "maintenance",
		"tutoring":     "tutoring",
		"consultation": "consultation",
		"photography":  "photography",
		"design":       "design",
		"development":  "development",
	}

	for keyword, serviceType := range serviceTypes {
		if strings.Contains(message, keyword) {
			return serviceType
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractAvailability(message string) string {
	availabilities := map[string]string{
		"immediately": "immediate",
		"asap":        "immediate",
		"flexible":    "flexible",
		"weekends":    "weekends",
		"weekdays":    "weekdays",
		"evenings":    "evenings",
	}

	for keyword, availability := range availabilities {
		if strings.Contains(message, keyword) {
			return availability
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractDiscountPercentage(message string) int {
	discountRegex := regexp.MustCompile(`(\d+)%\s*(?:off|discount)`)
	matches := discountRegex.FindStringSubmatch(message)
	if len(matches) > 1 {
		if discount, err := strconv.Atoi(matches[1]); err == nil {
			return discount
		}
	}
	return 0
}

func (pe *ParameterExtractor) extractDealType(message string) string {
	dealTypes := map[string]string{
		"flash sale":      "flash",
		"clearance":       "clearance",
		"bogo":            "buy_one_get_one",
		"buy one get one": "buy_one_get_one",
		"limited time":    "limited_time",
		"daily deal":      "daily",
	}

	for keyword, dealType := range dealTypes {
		if strings.Contains(message, keyword) {
			return dealType
		}
	}

	return ""
}

// Utility methods
func (pe *ParameterExtractor) extractLocationFromContext(context map[string]interface{}) (float64, float64, float64) {
	lat, latOk := context["lat"].(float64)
	lng, lngOk := context["lng"].(float64)
	radius, radiusOk := context["radius"].(float64)

	if !latOk || !lngOk {
		return 0, 0, 0
	}

	if !radiusOk {
		radius = 10.0 // Default 10km radius
	}

	return lat, lng, radius
}

func (pe *ParameterExtractor) extractSortBy(message string) string {
	sortOptions := map[string]string{
		"price":     "price",
		"date":      "created_at",
		"newest":    "created_at",
		"oldest":    "created_at",
		"popular":   "popularity",
		"rating":    "rating",
		"distance":  "distance",
		"relevance": "relevance",
	}

	for keyword, sortBy := range sortOptions {
		if strings.Contains(message, keyword) {
			return sortBy
		}
	}

	return ""
}

func (pe *ParameterExtractor) extractSortOrder(message string) string {
	if strings.Contains(message, "desc") || strings.Contains(message, "descending") ||
		strings.Contains(message, "highest") || strings.Contains(message, "newest") {
		return "desc"
	}

	if strings.Contains(message, "asc") || strings.Contains(message, "ascending") ||
		strings.Contains(message, "lowest") || strings.Contains(message, "oldest") {
		return "asc"
	}

	return "desc" // Default to descending
}

func (pe *ParameterExtractor) extractPageNumber(message string) int {
	pageRegex := regexp.MustCompile(`page\s*(\d+)`)
	matches := pageRegex.FindStringSubmatch(message)
	if len(matches) > 1 {
		if page, err := strconv.Atoi(matches[1]); err == nil {
			return page
		}
	}
	return 0
}

func (pe *ParameterExtractor) extractPageSize(message string) int {
	sizeRegex := regexp.MustCompile(`(?:show|limit|size)\s*(\d+)`)
	matches := sizeRegex.FindStringSubmatch(message)
	if len(matches) > 1 {
		if size, err := strconv.Atoi(matches[1]); err == nil {
			return size
		}
	}
	return 0
}

func (pe *ParameterExtractor) extractID(message, entityType string) string {
	idPattern := fmt.Sprintf(`(?:%s[_\s]*)?(?:id|ID)\s*[:\-]?\s*([a-zA-Z0-9\-_]+)`, entityType)
	idRegex := regexp.MustCompile(idPattern)
	matches := idRegex.FindStringSubmatch(message)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func (pe *ParameterExtractor) extractSize(message string) string {
	sizes := []string{"xs", "s", "m", "l", "xl", "xxl", "small", "medium", "large", "extra large"}
	for _, size := range sizes {
		if strings.Contains(message, size) {
			return size
		}
	}
	return ""
}
