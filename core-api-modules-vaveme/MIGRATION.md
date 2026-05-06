# Migration Guide

This guide helps you migrate from the existing API modules to the new `@vaveme/core-api-modules` library.

## Overview

The new library consolidates all API functionality into a single, well-organized package with:
- Centralized configuration
- Consistent error handling
- Built-in token management
- TypeScript support
- Better testing capabilities

## Installation

1. Install the library:
```bash
cd frontend
npm install ./core-api-modules-vaveme
```

2. Build the library:
```bash
cd core-api-modules-vaveme
npm install
npm run build
```

## Migration Steps

### 1. Update Imports

Replace individual API imports with the new library:

**Before:**
```javascript
import { loginUser, registerUser } from '@/api/userApi';
import { searchProducts } from '@/api/searchApi';
import axiosInstance from '@/api/axiosInstance';
```

**After:**
```javascript
import { createApiServices } from '@vaveme/core-api-modules';

const api = createApiServices();
```

### 2. Update API Calls

#### Authentication

**Before:**
```javascript
// Login
const response = await loginUser(email, password);

// Register
const response = await registerUser({
  email,
  password,
  username,
  acceptedTerms
});

// Logout
await logoutUser();
```

**After:**
```javascript
// Login
const { data } = await api.auth.login({ email, password });

// Register
const { data } = await api.auth.register({
  email,
  password,
  username,
  acceptTerms
});

// Logout
await api.auth.logout();
```

#### User Operations

**Before:**
```javascript
// Get user
const user = await getUserById(userId);

// Update user
const updated = await updateUser(userId, userData);

// Upload avatar
const result = await uploadAvatar(file);
```

**After:**
```javascript
// Get user
const { data: user } = await api.users.getUser(userId);

// Update user
const { data: updated } = await api.users.updateUser(userId, userData);

// Upload avatar
const { data: result } = await api.users.uploadAvatar(file);
```

#### Search

**Before:**
```javascript
// Search with filters
const results = await searchWithFilters({
  searchTerm: query,
  categories,
  priceRange,
  page
});

// Get suggestions
const suggestions = await getSuggestions(query);
```

**After:**
```javascript
// Search with filters
const { data: results } = await api.search.search({
  query,
  category: categories[0],
  priceMin: priceRange.min,
  priceMax: priceRange.max,
  page
});

// Get suggestions
const { data: suggestions } = await api.search.getSuggestions(query);
```

### 3. Update Error Handling

**Before:**
```javascript
try {
  const result = await someApiCall();
} catch (error) {
  if (error.response?.status === 401) {
    // Handle auth error
  } else {
    console.error('API Error:', error.message);
  }
}
```

**After:**
```javascript
try {
  const { data } = await api.someService.someMethod();
} catch (error) {
  // Error is already formatted by ApiErrorHandler
  console.error(error.userMessage);
  
  if (error.statusCode === 401) {
    // Handle auth error
  }
}
```

### 4. Update React Query Hooks

**Before:**
```javascript
import { useQuery } from '@tanstack/react-query';
import { getUserById } from '@/api/userApi';

function useUser(userId) {
  return useQuery({
    queryKey: ['user', userId],
    queryFn: () => getUserById(userId),
    enabled: !!userId
  });
}
```

**After:**
```javascript
import { useQuery } from '@tanstack/react-query';
import { createApiServices } from '@vaveme/core-api-modules';

const api = createApiServices();

function useUser(userId) {
  return useQuery({
    queryKey: ['user', userId],
    queryFn: async () => {
      const { data } = await api.users.getUser(userId);
      return data;
    },
    enabled: !!userId
  });
}
```

### 5. Update Form Submissions

**Before:**
```javascript
const handleSubmit = async (formData) => {
  try {
    setLoading(true);
    const result = await updateUserProfile(formData);
    toast.success('Profile updated');
  } catch (error) {
    toast.error(error.response?.data?.message || 'Update failed');
  } finally {
    setLoading(false);
  }
};
```

**After:**
```javascript
const handleSubmit = async (formData) => {
  try {
    setLoading(true);
    const { data } = await api.users.updateCurrentUser(formData);
    toast.success('Profile updated');
  } catch (error) {
    toast.error(error.userMessage);
  } finally {
    setLoading(false);
  }
};
```

### 6. Update File Uploads

**Before:**
```javascript
const formData = new FormData();
formData.append('file', file);
formData.append('type', 'avatar');

const response = await axiosInstance.post('/upload', formData, {
  headers: { 'Content-Type': 'multipart/form-data' }
});
```

**After:**
```javascript
const { data } = await api.users.uploadAvatar(file);
// Or for generic uploads:
const { data } = await api.media.upload(file, { type: 'avatar' });
```

### 7. Update Token Handling

**Before:**
```javascript
// Manual token management
localStorage.setItem('accessToken', token);
const token = localStorage.getItem('accessToken');
localStorage.removeItem('accessToken');

// Check if logged in
const isLoggedIn = !!localStorage.getItem('accessToken');
```

**After:**
```javascript
import { TokenManager } from '@vaveme/core-api-modules';

// Automatic token management (handled by auth service)
// Manual access if needed:
const token = TokenManager.getAccessToken();
const isLoggedIn = TokenManager.isAccessTokenValid();
```

## Common Patterns

### Creating a Custom API Instance

```javascript
import { AxiosClient, createApiServices } from '@vaveme/core-api-modules';

// Custom axios instance for specific needs
const customAxios = AxiosClient.create({
  config: {
    baseUrl: process.env.NEXT_PUBLIC_CUSTOM_API_URL,
    timeout: 30000,
    retryAttempts: 5
  }
});

const customApi = createApiServices(customAxios);
```

### SSR-Safe API Calls

```javascript
import { AxiosClient, createApiServices } from '@vaveme/core-api-modules';

// In getServerSideProps or server components
export async function getServerSideProps() {
  const ssrAxios = AxiosClient.getSSR();
  const api = createApiServices(ssrAxios);
  
  try {
    const { data } = await api.products.getProduct('123');
    return { props: { product: data } };
  } catch (error) {
    return { notFound: true };
  }
}
```

### Global API Instance

Create a singleton instance for your app:

```javascript
// lib/api.js
import { createApiServices } from '@vaveme/core-api-modules';

export const api = createApiServices();
```

Then import everywhere:
```javascript
import { api } from '@/lib/api';

const { data } = await api.users.getCurrentUser();
```

## Gradual Migration

You can migrate gradually by using both old and new APIs:

1. Install the library
2. Start using it for new features
3. Migrate existing features one by one
4. Remove old API files once fully migrated

## Testing

Update your tests to use the new API:

```javascript
import { createApiServices } from '@vaveme/core-api-modules';
import { render, waitFor } from '@testing-library/react';

// Mock the API
jest.mock('@vaveme/core-api-modules', () => ({
  createApiServices: () => ({
    users: {
      getUser: jest.fn().mockResolvedValue({
        data: { id: '1', name: 'Test User' }
      })
    }
  })
}));

test('loads user data', async () => {
  const { getByText } = render(<UserProfile userId="1" />);
  
  await waitFor(() => {
    expect(getByText('Test User')).toBeInTheDocument();
  });
});
```

## Troubleshooting

### Build Errors

If you encounter build errors:
1. Ensure the library is built: `cd core-api-modules-vaveme && npm run build`
2. Clear Next.js cache: `rm -rf .next`
3. Restart dev server

### TypeScript Errors

Add to your `tsconfig.json`:
```json
{
  "compilerOptions": {
    "paths": {
      "@vaveme/core-api-modules": ["./core-api-modules-vaveme/src"]
    }
  }
}
```

### Import Errors

If imports aren't resolving, update your Next.js config:
```javascript
// next.config.js
module.exports = {
  transpilePackages: ['@vaveme/core-api-modules'],
};
```

## Benefits After Migration

1. **Consistent API** - All services follow the same patterns
2. **Better Error Handling** - User-friendly error messages
3. **Automatic Retries** - Built-in retry logic for failed requests
4. **Token Management** - Automatic refresh and storage
5. **Type Safety** - Full TypeScript support
6. **Smaller Bundle** - Tree-shaking removes unused code
7. **Easier Testing** - Mockable service interfaces
8. **Better DX** - IntelliSense and documentation