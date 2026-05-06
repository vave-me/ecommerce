import { createApiServices, AxiosClient, TokenManager } from '@vaveme/core-api-modules';

// Quick Start Guide - Common Usage Patterns

async function quickStartExamples() {
  // 1. Basic Setup - Use default configuration
  const api = createApiServices();

  // 2. Authentication Flow
  try {
    // Login
    const { data: authData } = await api.auth.login({
      email: 'redacted-email@example.com',
      password: 'password123'
    });
    console.log('Logged in as:', authData.user.email);

    // Tokens are automatically stored, but you can check:
    console.log('Is authenticated:', api.auth.isAuthenticated());
    console.log('User ID:', api.auth.getUserId());

  } catch (error) {
    console.error('Login failed:', error.userMessage);
  }

  // 3. User Operations
  try {
    // Get current user
    const { data: currentUser } = await api.users.getCurrentUser();
    console.log('Current user:', currentUser);

    // Update profile
    const { data: updatedUser } = await api.users.updateCurrentUser({
      firstName: 'John',
      lastName: 'Doe',
      bio: 'Software Developer'
    });

    // Upload avatar
    const avatarFile = new File(['...'], 'avatar.jpg', { type: 'image/jpeg' });
    const { data: avatarResult } = await api.users.uploadAvatar(avatarFile);
    console.log('Avatar URL:', avatarResult.url);

  } catch (error) {
    console.error('User operation failed:', error.userMessage);
  }

  // 4. Search Operations
  try {
    // Advanced search
    const { data: searchResults } = await api.search.search({
      query: 'laptop',
      category: 'electronics',
      priceMin: 500,
      priceMax: 1500,
      sortBy: 'price_asc',
      page: 1,
      pageSize: 20
    });

    console.log(`Found ${searchResults.total} results`);
    console.log('First result:', searchResults.items[0]);

    // Get search suggestions
    const { data: suggestions } = await api.search.getSuggestions('mac');
    console.log('Suggestions:', suggestions);

  } catch (error) {
    console.error('Search failed:', error.userMessage);
  }

  // 5. Different Client Configurations
  
  // For SSR/Server-side rendering
  const ssrApi = createApiServices(AxiosClient.getSSR());
  
  // For public endpoints (no auth)
  const publicApi = createApiServices(AxiosClient.getPublic());
  
  // Custom configuration
  const customApi = createApiServices(
    AxiosClient.create({
      config: {
        baseUrl: 'https://api.staging.example.com',
        timeout: 30000,
        enableLogging: true
      }
    })
  );

  // 6. Error Handling Patterns
  async function robustApiCall() {
    try {
      const { data } = await api.users.getUser('user-123');
      return data;
    } catch (error) {
      // The error is already formatted with useful info
      console.error({
        message: error.userMessage,     // User-friendly message
        statusCode: error.statusCode,   // HTTP status
        severity: error.severity,       // Error severity
        endpoint: error.endpoint,       // Which endpoint failed
        timestamp: error.timestamp      // When it happened
      });
      
      // Handle specific errors
      if (error.statusCode === 404) {
        console.log('User not found');
      } else if (error.statusCode === 401) {
        console.log('Need to login again');
        // Redirect to login
      }
      
      throw error;
    }
  }

  // 7. Using with React Query
  function useUserProfile(userId: string) {
    return useQuery({
      queryKey: ['user', userId],
      queryFn: async () => {
        const { data } = await api.users.getUser(userId);
        return data;
      },
      enabled: !!userId
    });
  }

  // 8. Batch Operations
  async function batchOperations() {
    // Run multiple API calls in parallel
    const [users, products, categories] = await Promise.all([
      api.users.searchUsers({ query: 'john' }),
      api.products.getFeatured(),
      api.categories.getAll()
    ]);

    console.log('Batch results:', { users, products, categories });
  }

  // 9. File Upload Patterns
  async function handleFileUpload(files: FileList) {
    const uploadPromises = Array.from(files).map(file => 
      api.media.upload(file, { 
        folder: 'products',
        optimize: true 
      })
    );

    const results = await Promise.all(uploadPromises);
    const urls = results.map(r => r.data.url);
    console.log('Uploaded files:', urls);
  }

  // 10. Pagination Pattern
  async function paginatedSearch(page = 1) {
    const { data } = await api.search.search({
      query: 'shoes',
      page,
      pageSize: 20
    });

    console.log({
      currentPage: data.page,
      totalPages: data.totalPages,
      hasNext: data.hasNextPage,
      items: data.items
    });

    // Load next page
    if (data.hasNextPage) {
      await paginatedSearch(page + 1);
    }
  }
}

// React Hook Examples
export function useApiHooks() {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Login hook
  const login = useCallback(async (email: string, password: string) => {
    setLoading(true);
    setError(null);
    
    try {
      const api = createApiServices();
      const { data } = await api.auth.login({ email, password });
      // Handle successful login (e.g., redirect)
      return data;
    } catch (err) {
      setError(err.userMessage);
      throw err;
    } finally {
      setLoading(false);
    }
  }, []);

  // Search hook with debouncing
  const search = useMemo(() => {
    const api = createApiServices();
    
    return debounce(async (query: string) => {
      if (!query.trim()) return [];
      
      try {
        const { data } = await api.search.getSuggestions(query);
        return data;
      } catch (err) {
        console.error('Search failed:', err);
        return [];
      }
    }, 300);
  }, []);

  return { login, search, loading, error };
}

// Next.js API Route Example
export async function apiRouteExample(req: NextApiRequest, res: NextApiResponse) {
  // Use SSR client for server-side calls
  const api = createApiServices(AxiosClient.getSSR());

  try {
    const { data } = await api.products.getProduct(req.query.id as string);
    res.status(200).json(data);
  } catch (error) {
    res.status(error.statusCode || 500).json({
      error: error.userMessage || 'Internal server error'
    });
  }
}

// Environment-specific configuration
export function getEnvironmentApi() {
  const isDev = process.env.NODE_ENV === 'development';
  const isTest = process.env.NODE_ENV === 'test';

  if (isTest) {
    // Return mocked API for tests
    return createMockedApi();
  }

  if (isDev) {
    // Development configuration
    return createApiServices(
      AxiosClient.create({
        config: {
          baseUrl: 'http://localhost:3001',
          timeout: 60000, // Longer timeout for debugging
          enableLogging: true
        }
      })
    );
  }

  // Production configuration
  return createApiServices();
}