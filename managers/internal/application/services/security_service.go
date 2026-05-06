package services

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"
	"middleman/managers/internal/domain"
)

const (
	logPrefixSecurityService = "[SecurityService]"
)

// Predefined Trust Levels (example, can be enum or consts)
const (
	TrustLevelLow    = "low"
	TrustLevelMedium = "medium"
	TrustLevelHigh   = "high"
)

// SecurityService handles security validation, data protection, and endpoint access control.
type SecurityService struct {
	privateEndpointsConfig map[string]domain.PrivateEndpointConfig
	allowedPublicEndpoints map[string]bool // Added to make public endpoint check explicit
}

// NewSecurityService creates a new SecurityService.
func NewSecurityService(
	privateEndpoints map[string]domain.PrivateEndpointConfig,
	publicEndpoints map[string]bool, // e.g., VaverAllowedEndpoints
	_ interface{}, // Unused parameter kept for compatibility
) (*SecurityService, error) {
	if privateEndpoints == nil {
		log.Printf("%s WARN: No private endpoint configurations provided, initializing empty.", logPrefixSecurityService)
		privateEndpoints = make(map[string]domain.PrivateEndpointConfig)
	}
	if publicEndpoints == nil {
		log.Printf("%s WARN: No public endpoint configurations provided, initializing empty.", logPrefixSecurityService)
		publicEndpoints = make(map[string]bool)
	}

	return &SecurityService{
		privateEndpointsConfig: privateEndpoints,
		allowedPublicEndpoints: publicEndpoints,
	}, nil
}

// ValidateSecurityContext validates the security context for operations.
func (s *SecurityService) ValidateSecurityContext(secCtx domain.SecurityContext) error {
	if secCtx.UserID == "" {
		return fmt.Errorf("user ID is required in security context")
	}
	if secCtx.ExpiresAt.IsZero() {
		// log.Printf("%s WARN: SecurityContext for UserID '%s' has no expiration time set.", logPrefixSecurityService, secCtx.UserID)
		// Depending on policy, could be an error or just a warning. For now, allow.
	} else if time.Now().UTC().After(secCtx.ExpiresAt) {
		return fmt.Errorf("security context for UserID '%s' has expired at %s", secCtx.UserID, secCtx.ExpiresAt)
	}
	if secCtx.TrustLevel == "" { // Trust level should be a defined enum or set of consts
		return fmt.Errorf("trust level is required in security context")
	}
	// Basic validation for permissions and data scope, more specific checks in HasPermission/hasScope
	if len(secCtx.Permissions) == 0 && len(secCtx.DataScope) == 0 {
		log.Printf("%s WARN: SecurityContext for UserID '%s' has no permissions or data scopes defined.", logPrefixSecurityService, secCtx.UserID)
		// This might be acceptable for very low-trust operations, but generally indicates an issue.
	}
	return nil
}

// HasPermission checks if the security context has a specific permission, supporting wildcards.
func (s *SecurityService) HasPermission(secCtx domain.SecurityContext, requiredPermission string) bool {
	if requiredPermission == "" {
		return true
	} // No specific permission required
	normalizedReqPerm := strings.ToLower(strings.TrimSpace(requiredPermission))

	for _, userPerm := range secCtx.Permissions {
		normalizedUserPerm := strings.ToLower(strings.TrimSpace(userPerm))
		if normalizedUserPerm == normalizedReqPerm || normalizedUserPerm == "*" { // "*" grants all permissions
			return true
		}
		if strings.HasSuffix(normalizedUserPerm, "*") {
			prefix := strings.TrimSuffix(normalizedUserPerm, "*")
			if strings.HasPrefix(normalizedReqPerm, prefix) {
				return true
			}
		}
	}
	log.Printf("%s Permission denied: User '%s' lacks required permission '%s'. Has: %v",
		logPrefixSecurityService, secCtx.UserID, requiredPermission, secCtx.Permissions)
	return false
}

// CanAccessEndpoint checks if the security context allows access to a configured private endpoint.
func (s *SecurityService) CanAccessEndpoint(endpointPath string, secCtx domain.SecurityContext) (bool, error) {
	if s.privateEndpointsConfig == nil {
		log.Printf("%s ERROR: privateEndpointsConfig not initialized in SecurityService.", logPrefixSecurityService)
		return false, errors.New("internal security configuration error")
	}

	endpointConfig, exists := s.privateEndpointsConfig[endpointPath]
	if !exists {
		// This might be a public endpoint or an unknown one. This function is for *private* endpoints.
		// If it's a known public endpoint, access might be true by default (checked elsewhere).
		// If it's truly unknown, it should probably be denied.
		// For now, assume if not in private list, this specific check passes (or is N/A).
		// A more robust system would have an IsPrivateEndpoint method.
		log.Printf("%s INFO: Endpoint '%s' not found in private endpoint configurations. Access not restricted by this check.", logPrefixSecurityService, endpointPath)
		return true, nil // Or false if all non-listed endpoints should be denied by default
	}

	// 0. Validate the SecurityContext itself first
	if err := s.ValidateSecurityContext(secCtx); err != nil {
		return false, fmt.Errorf("invalid security context for endpoint access: %w", err)
	}

	// 1. Check Trust Level
	// Define trust levels with ordinal values for easier comparison
	trustLevelHierarchy := map[string]int{TrustLevelLow: 1, TrustLevelMedium: 2, TrustLevelHigh: 3, "verified": 3, "admin": 4} // Example

	userTrustOrdinal, userTrustOk := trustLevelHierarchy[strings.ToLower(secCtx.TrustLevel)]
	requiredTrustOrdinal, reqTrustOk := trustLevelHierarchy[strings.ToLower(endpointConfig.MinTrustLevel)]

	if !userTrustOk {
		log.Printf("%s WARN: User '%s' has an unrecognized trust level: '%s'", logPrefixSecurityService, secCtx.UserID, secCtx.TrustLevel)
		return false, errors.New("user has unrecognized trust level")
	}
	if !reqTrustOk && endpointConfig.MinTrustLevel != "" { // Only error if MinTrustLevel was specified but is unknown
		log.Printf("%s WARN: Endpoint '%s' has an unrecognized MinTrustLevel requirement: '%s'", logPrefixSecurityService, endpointPath, endpointConfig.MinTrustLevel)
		return false, errors.New("endpoint has unrecognized minimum trust level requirement")
	}

	if endpointConfig.MinTrustLevel != "" && userTrustOrdinal < requiredTrustOrdinal {
		log.Printf("%s Access Denied: User '%s' (trust: %s/%d) does not meet MinTrustLevel '%s' (%d) for endpoint '%s'",
			logPrefixSecurityService, secCtx.UserID, secCtx.TrustLevel, userTrustOrdinal, endpointConfig.MinTrustLevel, requiredTrustOrdinal, endpointPath)
		return false, nil
	}

	// 2. Check MFA Requirement
	if endpointConfig.RequiresMFA && !secCtx.MFAVerified {
		log.Printf("%s Access Denied: User '%s' MFA not verified for MFA-required endpoint '%s'", logPrefixSecurityService, secCtx.UserID, endpointPath)
		return false, nil
	}

	// 3. Check Required Scopes (using HasPermission for consistency, assuming scopes are a form of permission)
	for _, requiredScope := range endpointConfig.RequiredScopes {
		if !s.HasPermission(secCtx, requiredScope) { // Using HasPermission logic which includes wildcard support for scopes
			log.Printf("%s Access Denied: User '%s' lacks required scope '%s' for endpoint '%s'. Has scopes/permissions: %v",
				logPrefixSecurityService, secCtx.UserID, requiredScope, endpointPath, secCtx.Permissions) // Or secCtx.DataScope if they are distinct
			return false, nil
		}
	}

	log.Printf("%s INFO: User '%s' granted access to private endpoint '%s'", logPrefixSecurityService, secCtx.UserID, endpointPath)
	return true, nil
}

// MaskSensitiveData masks known sensitive fields in a map.
// This is a basic implementation; real-world PII masking is complex.
func (s *SecurityService) MaskSensitiveData(data map[string]interface{}, secCtx domain.SecurityContext) map[string]interface{} {
	if data == nil {
		return nil
	}
	maskedData := make(map[string]interface{}, len(data))

	// It's better to have a configurable list of sensitive fields and their masking rules.
	// fieldType here is a simplified concept. Permissions should be more granular.
	for key, value := range data {
		lowerKey := strings.ToLower(key)
		var maskedValue interface{} = value // Default to original value

		switch {
		// Always mask credentials regardless of permissions
		case lowerKey == "password" || lowerKey == "token" || lowerKey == "secret" || strings.Contains(lowerKey, "apikey") || strings.Contains(lowerKey, "sessionid"):
			maskedValue = "[MASKED_CREDENTIAL]"
		case lowerKey == "email" || lowerKey == "e-mail":
			if s.shouldMaskField(secCtx, "email_address", key) { // More specific permission
				if emailStr, ok := value.(string); ok {
					maskedValue = s.maskEmail(emailStr)
				}
			}
		case lowerKey == "phone" || lowerKey == "phone_number" || lowerKey == "mobile":
			if s.shouldMaskField(secCtx, "phone_number", key) {
				if phoneStr, ok := value.(string); ok {
					maskedValue = s.maskPhone(phoneStr)
				}
			}
		case strings.Contains(lowerKey, "card_number") || strings.Contains(lowerKey, "creditcard") || lowerKey == "pan":
			if s.shouldMaskField(secCtx, "payment_instrument", key) {
				if cardStr, ok := value.(string); ok {
					maskedValue = s.maskCreditCard(cardStr)
				}
			}
		case lowerKey == "ssn" || strings.Contains(lowerKey, "social_security") || strings.Contains(lowerKey, "national_id"):
			if s.shouldMaskField(secCtx, "national_identifier", key) {
				maskedValue = "***-**-****" // Generic SSN mask
			}
		case lowerKey == "address" || strings.Contains(lowerKey, "street"): // Example: Addresses might be sensitive
			if s.shouldMaskField(secCtx, "physical_address", key) {
				if addrStr, ok := value.(string); ok {
					maskedValue = s.maskAddress(addrStr)
				}
			}
		// TODO: Add more cases for common PII fields (date_of_birth, etc.)
		// Default: Check if the key itself is in a list of generally sensitive field names
		default:
			if s.isGenerallySensitiveKey(lowerKey) && s.shouldMaskField(secCtx, "generic_sensitive_data", key) {
				maskedValue = "[MASKED_SENSITIVE_DATA]"
			}
		}
		maskedData[key] = maskedValue
	}
	return maskedData
}

// shouldMaskField determines if a field should be masked.
// fieldType could be "email_address", "phone_number", "payment_instrument", etc.
// originalKey is logged for context.
func (s *SecurityService) shouldMaskField(secCtx domain.SecurityContext, fieldType string, originalKey string) bool {
	// By default, mask sensitive data unless specific "read_sensitive_FIELDTYPE" permission exists.
	// Example permission: "data:pii:email_address:read"
	// Or a more general one: "data:pii:read_all_unmasked"
	permissionToReadUnmasked := fmt.Sprintf("data:pii:%s:read_unmasked", fieldType)
	hasExplicitPermission := s.HasPermission(secCtx, permissionToReadUnmasked)

	if !hasExplicitPermission {
		log.Printf("%s INFO: Masking field '%s' (type: '%s') for UserID '%s' due to missing permission '%s'.",
			logPrefixSecurityService, originalKey, fieldType, secCtx.UserID, permissionToReadUnmasked)
		return true // Mask if no explicit permission to read unmasked
	}
	log.Printf("%s INFO: Field '%s' (type: '%s') NOT masked for UserID '%s' due to permission '%s'.",
		logPrefixSecurityService, originalKey, fieldType, secCtx.UserID, permissionToReadUnmasked)
	return false // Do not mask if explicit permission exists
}

func (s *SecurityService) isGenerallySensitiveKey(key string) bool {
	// Simple heuristic, expand as needed
	sensitiveKeywords := []string{"dob", "birthdate", "gender", "race", "ethnicity", "salary", "income"}
	for _, kw := range sensitiveKeywords {
		if strings.Contains(key, kw) {
			return true
		}
	}
	return false
}

// --- Masking Helper Functions ---
// Regexps are compiled once for efficiency if SecurityService is a long-lived object.
// For this file structure, they are local. Consider making them package-level vars or part of the struct.
var (
	nonDigitsRegex = regexp.MustCompile(`\D`)
)

func (s *SecurityService) maskEmail(email string) string {
	if email == "" {
		return "[EMPTY_EMAIL_MASKED]"
	}
	parts := strings.Split(email, "@")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "[INVALID_EMAIL_FORMAT_MASKED]"
	}
	username, domainPart := parts[0], parts[1]
	maskChar := "*"

	maskedUsername := username
	if len(username) > 3 {
		maskedUsername = string(username[0]) + strings.Repeat(maskChar, len(username)-2) + string(username[len(username)-1])
	} else if len(username) > 0 {
		maskedUsername = string(username[0]) + strings.Repeat(maskChar, maxInt(0, len(username)-1))
	}

	// Mask domain part less aggressively, e.g., only middle part of domain name
	domainParts := strings.Split(domainPart, ".")
	if len(domainParts) > 0 && domainParts[0] != "" {
		dName := domainParts[0]
		if len(dName) > 3 {
			domainParts[0] = string(dName[0]) + strings.Repeat(maskChar, len(dName)-2) + string(dName[len(dName)-1])
		} else {
			domainParts[0] = strings.Repeat(maskChar, len(dName))
		}
	}
	return maskedUsername + "@" + strings.Join(domainParts, ".")
}

func (s *SecurityService) maskPhone(phone string) string {
	if phone == "" {
		return "[EMPTY_PHONE_MASKED]"
	}
	digits := nonDigitsRegex.ReplaceAllString(phone, "")
	if len(digits) <= 4 {
		return strings.Repeat("*", len(digits))
	}
	// Example: mask all but last 4: ***-***-1234
	return strings.Repeat("*", len(digits)-4) + digits[len(digits)-4:]
}

func (s *SecurityService) maskCreditCard(card string) string {
	if card == "" {
		return "[EMPTY_CC_MASKED]"
	}
	digits := nonDigitsRegex.ReplaceAllString(card, "")
	if len(digits) < 4 {
		return "****"
	} // Not enough digits to show last 4
	if len(digits) <= 4 {
		return digits
	} // Show all if 4 or less (edge case)

	return strings.Repeat("*", len(digits)-4) + digits[len(digits)-4:]
}

func (s *SecurityService) maskAddress(address string) string {
	if address == "" {
		return "[EMPTY_ADDRESS_MASKED]"
	}
	// Simple masking: keep first word (number + street name part) and mask the rest.
	// This is very basic. Real address masking is complex.
	parts := strings.Fields(address)
	if len(parts) > 1 {
		return parts[0] + " " + strings.Repeat("X", len(strings.Join(parts[1:], " "))-len(parts)+1)
	} else if len(parts) == 1 {
		return strings.Repeat("X", len(parts[0]))
	}
	return "[MASKED_ADDRESS]"
}

// HashMessage creates a SHA256 hash of a message.
func (s *SecurityService) HashMessage(message string) string {
	if message == "" {
		return ""
	}
	hash := sha256.Sum256([]byte(message))
	return fmt.Sprintf("%x", hash)
}

// AuditSecureOperation logs secure operations.
// TODO: Integrate with a proper audit logging service/system.
func (s *SecurityService) AuditSecureOperation(
	ctx context.Context, // Pass context for tracing/cancellation
	operationName string, // Specific operation being audited
	secCtx domain.SecurityContext,
	requestDetails map[string]interface{}, // e.g., key parameters from the request
	outcome SecurityOperationOutcome, // Struct to hold success/failure and error message
) {
	// This is a placeholder for actual audit logging.
	// A real implementation would structure the audit log entry more comprehensively
	// and send it to a dedicated audit store (e.g., ELK stack, audit database).

	// Example of data to include in an audit entry:
	auditEntry := map[string]interface{}{
		"timestamp_utc":         time.Now().UTC().Format(time.RFC3339Nano),
		"service_name":          "managers_domain_security_service", // Or more specific service
		"operation_name":        operationName,
		"user_id":               secCtx.UserID,
		"trust_level":           secCtx.TrustLevel,
		"mfa_verified":          secCtx.MFAVerified,
		"ip_address":            secCtx.IPAddress,
		"user_agent":            "unknown",          // SecurityContext doesn't have UserAgent field
		"permissions_effective": secCtx.Permissions, // Could be filtered to relevant ones
		"data_scope_effective":  secCtx.DataScope,
		"request_details_hash":  s.HashMessage(fmt.Sprintf("%v", requestDetails)), // Hash details for brevity/security
		"outcome_success":       outcome.Success,
		"outcome_error":         "",
		"outcome_details":       outcome.Details,
	}
	if outcome.Error != nil {
		auditEntry["outcome_error"] = outcome.Error.Error()
	}

	log.Printf("%s AUDIT: UserID: %s, Operation: %s, Success: %t, Details: %v",
		logPrefixSecurityService, secCtx.UserID, operationName, outcome.Success, auditEntry)
	// In production, send `auditEntry` to a dedicated audit log system.
}

// SecurityOperationOutcome helps structure audit outcomes.
type SecurityOperationOutcome struct {
	Success bool
	Error   error
	Details map[string]interface{} // Additional outcome details
}

// GetPrivateEndpointConfig returns the configuration for a specific private endpoint.
func (s *SecurityService) GetPrivateEndpointConfig(endpointPath string) (domain.PrivateEndpointConfig, bool) {
	if s.privateEndpointsConfig == nil {
		log.Printf("%s WARN: GetPrivateEndpointConfig called but privateEndpointsConfig is nil.", logPrefixSecurityService)
		return domain.PrivateEndpointConfig{}, false
	}
	config, exists := s.privateEndpointsConfig[endpointPath]
	return config, exists
}

// IsPublicEndpoint checks if an endpoint is in the configured list of public endpoints.
func (s *SecurityService) IsPublicEndpoint(endpointPath string) bool {
	if s.allowedPublicEndpoints == nil {
		log.Printf("%s WARN: IsPublicEndpoint called but allowedPublicEndpoints is nil.", logPrefixSecurityService)
		return false // Secure default: if not configured, assume not public
	}
	// Normalize endpoint path for lookup (e.g., trim trailing slash if any)
	normalizedPath := strings.TrimSuffix(endpointPath, "/")
	return s.allowedPublicEndpoints[normalizedPath]
}

// GetAllowedEndpoints returns a list of all known public and private endpoints.
// This might be used to inform an LLM about its capabilities if it's allowed to see endpoint names.
func (s *SecurityService) GetAllowedEndpoints(secCtx domain.SecurityContext, capabilities []domain.ManagerCapability) []string {
	// This method's utility depends on whether LLMs are directly given endpoint names.
	// If an LLM uses abstract actions (e.g., "search_products"), this list might not be directly for it.
	// However, it can be used internally for validation or by a more privileged system component.

	allowed := make(map[string]struct{}) // Use a map to avoid duplicates

	// Check for general public and private access capabilities first
	hasPublicAPIAccess := s.hasCapability(capabilities, domain.CapabilityPublicAPIAccess)
	hasPrivateAPIAccess := s.hasCapability(capabilities, domain.CapabilityPrivateAPIAccess)

	if hasPublicAPIAccess && s.allowedPublicEndpoints != nil {
		for endpoint := range s.allowedPublicEndpoints {
			allowed[endpoint] = struct{}{}
		}
	}

	// For private endpoints, iterate and check CanAccessEndpoint for each.
	// This is more accurate than just checking CapabilityPrivateAPIAccess, as individual
	// private endpoints have their own scope/trust/MFA requirements.
	if hasPrivateAPIAccess && s.privateEndpointsConfig != nil {
		for endpointPath := range s.privateEndpointsConfig {
			// Perform a CanAccessEndpoint check here for the current user context.
			// This is more complex as it requires the SecurityContext.
			// The original code just listed all private endpoints if CapabilityPrivateAPIAccess was present.
			// A more secure approach:
			canAccess, err := s.CanAccessEndpoint(endpointPath, secCtx)
			if err == nil && canAccess {
				allowed[endpointPath] = struct{}{}
			} else if err != nil {
				log.Printf("%s WARN: Error checking access for endpoint '%s' in GetAllowedEndpoints: %v", logPrefixSecurityService, endpointPath, err)
			}
		}
	}

	endpointsList := make([]string, 0, len(allowed))
	for ep := range allowed {
		endpointsList = append(endpointsList, ep)
	}
	return endpointsList
}

// Helper to check capability, assuming it's part of Manager domain or passed in.
func (s *SecurityService) hasCapability(capabilities []domain.ManagerCapability, target domain.ManagerCapability) bool {
	for _, cap := range capabilities {
		if cap == target {
			return true
		}
	}
	return false
}

// Helper needed by MaskSensitiveData
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
