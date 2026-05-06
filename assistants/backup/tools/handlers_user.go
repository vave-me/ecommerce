package tools

import (
	"context"
	"fmt"
)

// ==================== USER HANDLERS ====================
func (r *ComprehensiveToolRegistry) initializeUserHandlers() {
	r.handlers["user_get_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return reg.userRepo.GetUserByID(ctx, userID)
	}

	r.handlers["user_get_base_by_id"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return reg.userRepo.GetBaseUserByID(ctx, userID)
	}

	r.handlers["user_get_multiple_by_ids"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userIDs := getStringArrayParam(params, "user_ids")
		if len(userIDs) == 0 {
			return nil, fmt.Errorf("user_ids array cannot be empty")
		}
		// Validate each user ID
		for i, userID := range userIDs {
			if err := ValidateIDParam(fmt.Sprintf("user_ids[%d]", i), userID); err != nil {
				return nil, fmt.Errorf("invalid user ID at index %d: %w", i, err)
			}
		}
		return reg.userRepo.GetMultipleUsersByIDs(ctx, userIDs)
	}

	r.handlers["user_get_all_participating"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		return reg.userRepo.GetAllParticipatingUsers(ctx)
	}

	r.handlers["user_create_new"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		email := getStringParam(params, "email")
		password := getStringParam(params, "password")
		username := getStringParam(params, "username")
		firstName := getStringParam(params, "first_name")
		lastName := getStringParam(params, "last_name")
		location := getStringParam(params, "location")
		lat := getFloat64Param(params, "lat", 0)
		lng := getFloat64Param(params, "lng", 0)
		thumbnail := getStringParam(params, "thumbnail")
		language := getStringParam(params, "language")

		// Validate required fields
		if err := ValidateEmailParam(email); err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}
		if password == "" {
			return nil, fmt.Errorf("password is required")
		}
		if username == "" {
			return nil, fmt.Errorf("username is required")
		}
		if firstName == "" {
			return nil, fmt.Errorf("first_name is required")
		}
		if lastName == "" {
			return nil, fmt.Errorf("last_name is required")
		}

		// Validate coordinates if provided
		v := NewValidator()
		if lat != 0 {
			v.ValidateLatitude("lat", lat)
		}
		if lng != 0 {
			v.ValidateLongitude("lng", lng)
		}
		if err := v.GetError(); err != nil {
			return nil, fmt.Errorf("invalid coordinates: %w", err)
		}

		// Sanitize string inputs
		username = SanitizeString(username)
		firstName = SanitizeString(firstName)
		lastName = SanitizeString(lastName)
		location = SanitizeString(location)

		return reg.userRepo.CreateNewUser(ctx,
			email,
			password,
			username,
			firstName,
			lastName,
			location,
			float32(lat),
			float32(lng),
			thumbnail,
			language)
	}

	r.handlers["user_update_profile"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		username := getStringParam(params, "username")
		firstName := getStringParam(params, "first_name")
		lastName := getStringParam(params, "last_name")
		bio := getStringParam(params, "bio")
		privacy := getStringParam(params, "privacy")
		background := getStringParam(params, "background")
		location := getStringParam(params, "location")
		lat := getFloat64Param(params, "lat", 0)
		lng := getFloat64Param(params, "lng", 0)
		thumbnail := getStringParam(params, "thumbnail")

		// Validate required fields
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid user id: %w", err)
		}

		// Validate coordinates if provided
		v := NewValidator()
		if lat != 0 {
			v.ValidateLatitude("lat", lat)
		}
		if lng != 0 {
			v.ValidateLongitude("lng", lng)
		}
		if err := v.GetError(); err != nil {
			return nil, fmt.Errorf("invalid coordinates: %w", err)
		}

		// Sanitize string inputs
		username = SanitizeString(username)
		firstName = SanitizeString(firstName)
		lastName = SanitizeString(lastName)
		bio = SanitizeString(bio)
		location = SanitizeString(location)

		return reg.userRepo.UpdateUserProfile(ctx,
			id,
			username,
			firstName,
			lastName,
			bio,
			privacy,
			background,
			location,
			float32(lat),
			float32(lng),
			thumbnail)
	}

	r.handlers["user_change_username"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		username := getStringParam(params, "username")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid user id: %w", err)
		}
		if username == "" {
			return nil, fmt.Errorf("username is required")
		}
		username = SanitizeString(username)
		return reg.userRepo.ChangeUsername(ctx, id, username)
	}

	r.handlers["user_activate_account"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		verificationToken := getStringParam(params, "verification_token")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid user id: %w", err)
		}
		if verificationToken == "" {
			return nil, fmt.Errorf("verification_token is required")
		}
		return nil, reg.userRepo.ActivateUserAccount(ctx, id, verificationToken)
	}

	r.handlers["user_deactivate_account"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid user id: %w", err)
		}
		return nil, reg.userRepo.DeactivateUserAccount(ctx, id)
	}

	r.handlers["user_authenticate"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		email := getStringParam(params, "email")
		password := getStringParam(params, "password")
		if err := ValidateEmailParam(email); err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}
		if password == "" {
			return nil, fmt.Errorf("password is required")
		}
		return reg.userRepo.AuthenticateUser(ctx, email, password)
	}

	r.handlers["user_authenticate_google_web"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		idToken := getStringParam(params, "id_token")
		if idToken == "" {
			return nil, fmt.Errorf("id_token is required")
		}
		return reg.userRepo.AuthenticateWithGoogleWeb(ctx, idToken)
	}

	r.handlers["user_authenticate_google_mobile"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		idToken := getStringParam(params, "id_token")
		if idToken == "" {
			return nil, fmt.Errorf("id_token is required")
		}
		return reg.userRepo.AuthenticateWithGoogleMobile(ctx, idToken)
	}

	r.handlers["user_logout"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		id := getStringParam(params, "id")
		authToken := getStringParam(params, "auth_token")
		refreshToken := getStringParam(params, "refresh_token")
		if err := ValidateIDParam("id", id); err != nil {
			return nil, fmt.Errorf("invalid user id: %w", err)
		}
		if authToken == "" {
			return nil, fmt.Errorf("auth_token is required")
		}
		if refreshToken == "" {
			return nil, fmt.Errorf("refresh_token is required")
		}
		return nil, reg.userRepo.LogoutUser(ctx, id, authToken, refreshToken)
	}

	r.handlers["user_refresh_auth_token"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		refreshToken := getStringParam(params, "refresh_token")
		userID := getStringParam(params, "user_id")
		if refreshToken == "" {
			return nil, fmt.Errorf("refresh_token is required")
		}
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		return reg.userRepo.RefreshUserAuthToken(ctx, refreshToken, userID)
	}

	r.handlers["user_revoke_tokens"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		userID := getStringParam(params, "user_id")
		tokenID := getStringParam(params, "token_id")
		refreshToken := getStringParam(params, "refresh_token")
		reason := getStringParam(params, "reason")
		if err := ValidateIDParam("user_id", userID); err != nil {
			return nil, fmt.Errorf("invalid user_id: %w", err)
		}
		if tokenID == "" {
			return nil, fmt.Errorf("token_id is required")
		}
		if refreshToken == "" {
			return nil, fmt.Errorf("refresh_token is required")
		}
		if reason == "" {
			return nil, fmt.Errorf("reason is required")
		}
		reason = SanitizeString(reason)
		return reg.userRepo.RevokeUserTokens(ctx, userID, tokenID, refreshToken, reason)
	}

	r.handlers["user_send_password_reset_email"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		email := getStringParam(params, "email")
		if err := ValidateEmailParam(email); err != nil {
			return nil, fmt.Errorf("invalid email: %w", err)
		}
		return reg.userRepo.SendPasswordResetEmail(ctx, email)
	}

	r.handlers["user_reset_password"] = func(ctx context.Context, reg *ComprehensiveToolRegistry, params map[string]interface{}) (interface{}, error) {
		token := getStringParam(params, "token")
		newPassword := getStringParam(params, "new_password")
		if token == "" {
			return nil, fmt.Errorf("token is required")
		}
		if newPassword == "" {
			return nil, fmt.Errorf("new_password is required")
		}
		return reg.userRepo.ResetUserPassword(ctx, token, newPassword)
	}
}