Authentication Workflow and Token Lifecycle
Table of Contents
Authentication Methods
Token Architecture
Authentication Flows
Token Lifecycle Management
Security Implementation Details
Event Tracking and Audit
1. Authentication Methods
   The system supports three primary authentication methods:
   1.1 Email/Password Authentication
   Standard username/password login
   Credential validation against stored (hashed) passwords
   Returns JWT access token and refresh token
   1.2 Web Google OAuth Authentication
   OAuth 2.0 flow for web applications
   Uses Google's authentication provider
   Verifies ID token on the server side
   1.3 Mobile Google OAuth Authentication
   Similar to web OAuth but with mobile-specific configuration
   Handles mobile app authentication scenarios
   Uses different client IDs and verification settings
2. Token Architecture
   2.1 Token Types
   Access Tokens
   Format: JWT (JSON Web Token)
   Contents: User ID, username, email, first/last name, expiration
   Signature: Signed with server's private key
   Expiration: Short-lived (typically 15-60 minutes)
   Usage: Sent with each API request in Authorization header
   Refresh Tokens
   Format: Secure random string or JWT
   Contents: User ID, token ID, expiration
   Storage: Stored securely in database with user reference
   Expiration: Longer-lived (days to weeks)
   Usage: Used only to obtain new access/refresh token pairs
   2.2 Token Storage
   Server-side
   Refresh tokens are tracked in database with:
   User ID
   Token ID (partial hash)
   Creation timestamp
   Expiration timestamp
   Last used timestamp
   Device/client information
   Client-side
   Access token: LocalStorage or memory
   Refresh token: HttpOnly cookie or secure storage
   Tokens are never exposed to JavaScript when using cookies
3. Authentication Flows
   3.1 Standard Login Flow
   Frontend:
'async function login(email, password) {
   const response = await fetch('/api/users/login', {
   method: 'POST',
   headers: { 'Content-Type': 'application/json' },
   body: JSON.stringify({ email, password })
   });

if (!response.ok) throw new Error('Authentication failed');

const data = await response.json();
// Store tokens securely
storeTokens(data.token, data.accessToken, data.userName);

return data;
}'


Backend Process:
LoginUser gRPC method receives credentials
Credentials are validated against stored user data
On success, token pair is generated:

tokens, err := h.auth.GenerateTokenPair(&auth.JwtUser{
ID:        user.ID,
Email:     user.Email,
FirstName: user.FirstName,
LastName:  user.LastName,
Username:  user.Username,
})
UserLoggedIn domain event is triggered
Tokens are returned to client
Server records token metadata for the user
3.2 Google OAuth Login Flow (Web)
Frontend:

// Using Google's OAuth client library
async function initGoogleLogin() {
const auth2 = await gapi.auth2.init({
client_id: 'YOUR_GOOGLE_CLIENT_ID.apps.googleusercontent.com'
});

auth2.attachClickHandler(document.getElementById('googleSignIn'), {},
async (googleUser) => {
// Get ID token
const idToken = googleUser.getAuthResponse().id_token;

      // Send to backend
      const response = await fetch('/api/users/google-login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ idToken })
      });
      
      if (!response.ok) throw new Error('Google authentication failed');
      
      const data = await response.json();
      storeTokens(data.token, data.accessToken, data.userName);
      
      return data;
    }
);
}
ackend Process:
WebLoginWithGoogle gRPC method receives ID token
Token is verified with Google's services:

verifiedToken, err := googleOIDCClient.VerifyIDToken(ctx, idToken)

var claims oidcclient.GoogleUserClaims
if err := googleOIDCClient.ParseClaims(verifiedToken, &claims); err != nil {
return nil, status.Errorf(codes.InvalidArgument, "failed to parse Google token claims")
}

System checks if user exists by Google ID
If not found, new user is created with Google profile data
UserLoggedIn domain event is triggered
JWT tokens are generated and returned to client
3.3 Google OAuth Login Flow (Mobile)
Similar to web flow but uses mobile-specific client configuration and verification:
Backend Process:
MobileLoginWithGoogle gRPC method receives ID token
Uses mobile-specific Google verifier:


googleOIDCClient, ok := di.Get(ctx, constants.MobileGoogleVerifierKey)

4. Token Lifecycle Management
   4.1 Using Access Tokens
   Frontend:

// Attach access token to all API requests
axios.interceptors.request.use(config => {
const token = getAccessToken();
if (token) {
config.headers.Authorization = `Bearer ${token}`;
}
return config;
});

Backend Process:
Middleware extracts token from Authorization header
Token is validated:
Signature verification
Expiration check
Additional claims validation
User context is added to request
Request proceeds to protected endpoint
4.2 Token Refresh Flow
Frontend:

// Intercept 401 errors and attempt token refresh
axios.interceptors.response.use(
response => response,
async error => {
const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;
      
      try {
        // Refresh tokens
        const refreshToken = getRefreshToken();
        const userId = extractUserIdFromToken(getAccessToken());
        
        const response = await axios.post('/api/users/refresh-token', {
          refreshToken,
          userId // Optional but helps with validation
        });
        
        // Update stored tokens
        storeTokens(response.data.token, response.data.refreshToken);
        
        // Retry original request with new token
        originalRequest.headers.Authorization = `Bearer ${response.data.token}`;
        return axios(originalRequest);
      } catch (refreshError) {
        // Force re-login on refresh failure
        clearTokens();
        redirectToLogin();
        return Promise.reject(refreshError);
      }
    }
    
    return Promise.reject(error);
}
);

Backend Process:
RefreshAuthToken gRPC method receives refresh token
Token is validated:

jwtUser, err := h.auth.ValidateRefreshToken(cmd.RefreshToken)

User is verified as activ

user, err := h.middleman.Find(ctx, jwtUser.ID)
if !user.Enabled {
return "", "", errors.New("user account is disabled")
}

New token pair is generated
Old refresh token is invalidated
Token refresh is recorded:


oldTokenID := safeSubstring(cmd.RefreshToken, 8)
newTokenID := safeSubstring(tokens.RefreshToken, 8)
event, err := domainUser.TokenRefreshed(oldTokenID, newTokenID)


UserTokenRefreshed domain event is triggered
New tokens are returned to client
4.3 Logout Process
Frontend:
async function logout() {
const userId = extractUserIdFromToken(getAccessToken());
const accessToken = getAccessToken();
const refreshToken = getRefreshToken();

try {
await axios.post('/api/users/logout', {
id: userId,
authToken: accessToken,
refreshToken: refreshToken
});
} catch (error) {
console.error('Logout failed:', error);
} finally {
// Always clear tokens regardless of success/failure
clearTokens();
redirectToLogin();
}
}
Backend Process:
LogoutUser gRPC method receives user ID and tokens
User record is retrieved
Tokens are explicitly invalidated:

tokenID := "unknown"
if cmd.RefreshToken != "" {
tokenID = safeSubstring(cmd.RefreshToken, 8)
}
tokenEvent, err = user.InvalidateTokens(tokenID, "user logout")

UserLoggedOut domain event is triggered
UserTokenInvalidated domain event is triggered
Success response is returned to client
4.4 Token Invalidation for Security

async function clearUserTokens(reason = "security concern") {
const userId = extractUserIdFromToken(getAccessToken());
const refreshToken = getRefreshToken();

try {
await axios.post('/api/users/clear-tokens', {
userId,
refreshToken,
reason
});
} catch (error) {
console.error('Token invalidation failed:', error);
} finally {
clearTokens();
redirectToLogin();
}
}
Backend Process:
ClearTokens gRPC method receives user ID and tokens
User record is retrieved
Token invalidation is executed:

event, err := user.InvalidateTokens(tokenID, reason)

5. Security Implementation Details
   5.1 Token Security Features
   Token Rotation: Refresh tokens are single-use and rotated with each use
   Short-lived Access Tokens: Minimize damage from token theft
   Token Tracking: All tokens are tracked with partial identifiers
   Explicit Invalidation: Support for immediate invalidation of compromised tokens
   Device Tracking: Optional linking of tokens to specific devices/browsers
   5.2 Defense Against Common Attacks
   Token Theft:
   HttpOnly, Secure cookies for refresh tokens
   Short-lived access tokens
   Token rotation policy
   CSRF Protection:
   SameSite cookie attributes
   CSRF tokens for sensitive operations
   Token Replay:
   Single-use refresh tokens
   Token tracking in database
   Man-in-the-Middle:
   TLS for all communications
   Token signatures
   Brute Force:
   Rate limiting on login and token endpoints
   Account lockout after multiple failed attempts










