# @vaveme/core-api-modules

Core API modules for the Vaveme platform, providing a unified interface for all API interactions.

## Installation

```bash
npm install @vaveme/core-api-modules
```

## Features

- 🔐 **Built-in Authentication** - JWT token management with automatic refresh
- 🔄 **Retry Logic** - Intelligent retry mechanism for transient failures
- 🛡️ **Type Safety** - Full TypeScript support with detailed interfaces
- 📦 **Modular Design** - Use only what you need
- 🚀 **Optimized for SSR** - Special handling for server-side rendering
- 🔍 **Comprehensive Error Handling** - Consistent error responses
- 🧪 **Easy Testing** - Mockable services and dependency injection

## Quick Start

```typescript
import { createApiServices, AxiosClient } from '@vaveme/core-api-modules';

// Use default configuration
const api = createApiServices();

// Or create with custom axios instance
const customAxios = AxiosClient.create({
  config: {
    baseUrl: 'https://api.sfx-markt.de',
    timeout: 10000,
  }
});
const api = createApiServices(customAxios);

// Use the services
const { data: user } = await api.auth.login({
  email: 'user@example.com',
  password: 'password123'
});
```

## Core Concepts

### Axios Clients

The library provides pre-configured axios instances:

```typescript
import { AxiosClient } from '@vaveme/core-api-modules';

// Default authenticated client
const authClient = AxiosClient.getDefault();

// Public endpoints (no auth)
const publicClient = AxiosClient.getPublic();

// SSR-optimized client
const ssrClient = AxiosClient.getSSR();

// Custom configuration
const customClient = AxiosClient.create({
  config: {
    baseUrl: 'https://api.example.com',
    timeout: 5000,
    retryAttempts: 3,
  },
  includeAuth: true,
  isSSR: false,
});
```

### Token Management

Automatic token handling with the TokenManager:

```typescript
import { TokenManager } from '@vaveme/core-api-modules';

// Check authentication status
if (TokenManager.isAccessTokenValid()) {
  // User is authenticated
}

// Get user ID from token
const userId = TokenManager.getUserIdFromToken(TokenManager.getAccessToken());

// Manual token management
TokenManager.setTokens(accessToken, refreshToken);
TokenManager.clearTokens();
```

### Error Handling

Consistent error handling across all services:

```typescript
import { ApiErrorHandler, ErrorSeverity } from '@vaveme/core-api-modules';

try {
  await api.users.getUser('123');
} catch (error) {
  const apiError = ApiErrorHandler.handle(error);
  
  console.log(apiError.userMessage); // User-friendly message
  console.log(apiError.severity); // Error severity level
  console.log(apiError.statusCode); // HTTP status code
}
```

## Available Services

### Authentication Service

```typescript
const api = createApiServices();

// Login
const { data } = await api.auth.login({
  email: 'user@example.com',
  password: 'password123',
  rememberMe: true
});

// Register
await api.auth.register({
  email: 'user@example.com',
  password: 'password123',
  username: 'johndoe',
  acceptTerms: true
});

// Logout
await api.auth.logout();

// Password reset
await api.auth.requestPasswordReset({ email: 'user@example.com' });
await api.auth.confirmPasswordReset({ token: 'reset-token', newPassword: 'newPassword123' });
```

### User Service

```typescript
// Get user profile
const { data: user } = await api.users.getUser('user-id');

// Update profile
await api.users.updateCurrentUser({
  firstName: 'John',
  lastName: 'Doe',
  bio: 'Software Developer'
});

// Upload avatar
const file = new File(['...'], 'avatar.jpg');
await api.users.uploadAvatar(file);

// Search users
const { data: results } = await api.users.searchUsers({
  query: 'john',
  page: 1,
  pageSize: 20
});
```

### Search Service

```typescript
// Advanced search
const { data: results } = await api.search.search({
  query: 'laptop',
  category: 'electronics',
  priceMin: 500,
  priceMax: 1500,
  sortBy: 'price_asc'
});

// Quick search
const { data: items } = await api.search.quickSearch('macbook pro');

// Get suggestions
const { data: suggestions } = await api.search.getSuggestions('lapt');

// Search by image
const imageFile = new File(['...'], 'product.jpg');
const { data: similar } = await api.search.searchByImage(imageFile);
```

## Utilities

### Validators

```typescript
import { Validators } from '@vaveme/core-api-modules';

// Validate email
const emailResult = Validators.email('user@example.com');
if (!emailResult.isValid) {
  console.log(emailResult.errors);
}

// Clean parameters
const cleaned = Validators.cleanParams({
  name: 'John',
  age: null,
  email: undefined
}); // { name: 'John' }
```

### Encoders

```typescript
import { Encoders } from '@vaveme/core-api-modules';

// Build URL with query params
const url = Encoders.buildUrl('/api/users', {
  page: 1,
  sort: 'name',
  filters: ['active', 'verified']
});

// Encode path parameters
const path = Encoders.buildPath('/users/:id/posts/:postId', {
  id: 'user-123',
  postId: 'post-456'
});
```

### Mappers

```typescript
import { Mappers } from '@vaveme/core-api-modules';

// Map paginated response
const paginated = Mappers.mapPaginatedResponse(response);

// Convert case styles
const camelCase = Mappers.snakeToCamel(snake_case_obj);
const snake_case = Mappers.camelToSnake(camelCaseObj);

// Extract error message
const message = Mappers.extractErrorMessage(error);
```

## Configuration

### Environment Variables

```env
NEXT_PUBLIC_API_URL=https://api.sfx-markt.de
NODE_ENV=production
```

### Custom Configuration

```typescript
import { defaultConfig } from '@vaveme/core-api-modules';

const customConfig = {
  ...defaultConfig,
  timeout: 15000,
  retryAttempts: 5,
  enableLogging: true
};
```

## TypeScript Support

The library is written in TypeScript and provides comprehensive type definitions:

```typescript
import type {
  User,
  SearchFilters,
  ApiResponse,
  PaginatedResponse
} from '@vaveme/core-api-modules';
```

## Testing

Mock services for testing:

```typescript
import { createApiServices } from '@vaveme/core-api-modules';
import MockAdapter from 'axios-mock-adapter';
import axios from 'axios';

const mock = new MockAdapter(axios);
const api = createApiServices(axios);

mock.onPost('/auth/login').reply(200, {
  success: true,
  accessToken: 'mock-token',
  refreshToken: 'mock-refresh',
  user: { id: '1', email: 'test@example.com' }
});
```

## Error Codes

Common error codes and their meanings:

- `400` - Bad Request (invalid input)
- `401` - Unauthorized (token expired/invalid)
- `403` - Forbidden (insufficient permissions)
- `404` - Not Found
- `429` - Too Many Requests
- `500` - Internal Server Error
- `503` - Service Unavailable

## Contributing

1. Fork the repository
2. Create your feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## License

MIT © Vaveme Team