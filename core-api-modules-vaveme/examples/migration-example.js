// EXAMPLE: Migrating from old userApi.jsx to new library

// ============ OLD CODE (userApi.jsx) ============
/*
import axiosInstance from './axiosInstance';

export const loginUser = async (email, password) => {
  try {
    const response = await axiosInstance.post('/users/login', { email, password });
    if (response.data.accessToken) {
      localStorage.setItem('accessToken', response.data.accessToken);
      localStorage.setItem('refreshToken', response.data.refreshToken);
    }
    return response.data;
  } catch (error) {
    throw error;
  }
};

export const getUserById = async (userId) => {
  try {
    const response = await axiosInstance.get(`/users/${userId}`);
    return response.data;
  } catch (error) {
    console.error('Error fetching user:', error);
    throw error;
  }
};

export const updateUser = async (userId, userData) => {
  try {
    const response = await axiosInstance.put(`/users/${userId}`, userData);
    return response.data;
  } catch (error) {
    console.error('Error updating user:', error);
    throw error;
  }
};
*/

// ============ NEW CODE (using library) ============
import { createApiServices } from '@vaveme/core-api-modules';

// Create API instance
const api = createApiServices();

// Example: Login component
export function LoginExample() {
  const handleLogin = async (email, password) => {
    try {
      // OLD: const result = await loginUser(email, password);
      
      // NEW: Token management is automatic
      const { data } = await api.auth.login({ email, password });
      
      // User data is in data.user
      console.log('Logged in user:', data.user);
      
      // Tokens are automatically stored by the library
      // No need for manual localStorage.setItem
      
      return data;
    } catch (error) {
      // Error is already formatted with user-friendly message
      console.error('Login failed:', error.userMessage);
      throw error;
    }
  };
}

// Example: User profile component
export function UserProfileExample({ userId }) {
  const [user, setUser] = useState(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchUser = async () => {
      try {
        // OLD: const userData = await getUserById(userId);
        
        // NEW: Consistent response format
        const { data } = await api.users.getUser(userId);
        setUser(data);
      } catch (error) {
        // Error handling is more consistent
        toast.error(error.userMessage || 'Failed to load user');
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [userId]);

  return (
    <div>
      {loading ? 'Loading...' : user?.displayName}
    </div>
  );
}

// Example: Update user with React Query
export function useUpdateUser() {
  const queryClient = useQueryClient();
  
  return useMutation({
    mutationFn: async ({ userId, data }) => {
      // OLD: return await updateUser(userId, data);
      
      // NEW: Better error handling and response format
      const response = await api.users.updateUser(userId, data);
      return response.data;
    },
    onSuccess: (data, variables) => {
      // Invalidate relevant queries
      queryClient.invalidateQueries(['user', variables.userId]);
      toast.success('Profile updated successfully');
    },
    onError: (error) => {
      // Consistent error format
      toast.error(error.userMessage || 'Update failed');
    }
  });
}

// Example: Custom hook with the new API
export function useUser(userId) {
  return useQuery({
    queryKey: ['user', userId],
    queryFn: async () => {
      // OLD: return await getUserById(userId);
      
      // NEW: Extract data from response
      const { data } = await api.users.getUser(userId);
      return data;
    },
    enabled: !!userId,
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
}

// Example: File upload
export function AvatarUploadExample() {
  const handleUpload = async (file) => {
    try {
      // OLD: 
      // const formData = new FormData();
      // formData.append('avatar', file);
      // const response = await axiosInstance.post('/users/me/avatar', formData);
      
      // NEW: Simpler API
      const { data } = await api.users.uploadAvatar(file);
      
      console.log('Avatar URL:', data.url);
      return data.url;
    } catch (error) {
      console.error('Upload failed:', error.userMessage);
      throw error;
    }
  };
}

// Example: Search with complex filters
export function SearchExample() {
  const [results, setResults] = useState([]);
  
  const handleSearch = async (query, filters) => {
    try {
      // OLD:
      // const response = await searchWithFilters({
      //   searchTerm: query,
      //   categories: filters.categories,
      //   priceRange: { min: filters.minPrice, max: filters.maxPrice },
      //   page: filters.page
      // });
      
      // NEW: More intuitive API
      const { data } = await api.search.search({
        query,
        category: filters.categories?.[0],
        priceMin: filters.minPrice,
        priceMax: filters.maxPrice,
        page: filters.page,
        sortBy: filters.sortBy
      });
      
      setResults(data.items);
      
      // Access pagination info
      console.log(`Page ${data.page} of ${data.totalPages}`);
      
      // Access facets for filter counts
      if (data.facets) {
        console.log('Available categories:', data.facets.categories);
      }
    } catch (error) {
      console.error('Search failed:', error.userMessage);
    }
  };
}

// Example: Using different axios instances
export function MultiInstanceExample() {
  // For server-side rendering
  const getServerSideProps = async () => {
    const api = createApiServices(AxiosClient.getSSR());
    const { data } = await api.products.getFeatured();
    return { props: { products: data } };
  };

  // For public endpoints (no auth)
  const fetchPublicData = async () => {
    const api = createApiServices(AxiosClient.getPublic());
    const { data } = await api.categories.getAll();
    return data;
  };

  // For custom configuration
  const fetchWithCustomTimeout = async () => {
    const customAxios = AxiosClient.create({
      config: { timeout: 30000 }
    });
    const api = createApiServices(customAxios);
    const { data } = await api.reports.generateLarge();
    return data;
  };
}