import { useAuth } from '../context/AuthContext';

export const useUserRole = () => {
  const { user } = useAuth();
  
  return {
    role: user?.role || 'customer',
    isAdmin: user?.role === 'admin',
    isBusiness: user?.role === 'business',
    isCustomer: user?.role === 'customer',
    hasRole: (requiredRole) => user?.role === requiredRole,
    hasAnyRole: (roles) => roles.includes(user?.role)
  };
};