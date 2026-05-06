// src/components/UserProfile/ItemList.jsx
import React, { memo } from 'react';
import PropTypes from 'prop-types';
import { useQuery } from '@tanstack/react-query';
import { getUnifiedCatalog } from '../../api/client/searchApi';
import styles from './ItemList.module.css';
import { useAuth } from '../../context/AuthContext';
const ItemList = memo(function ItemList({ userId }) {
    const { user } = useAuth();
    // Use the authenticated user's ID, or fall back to the passed userId
    const effectiveUserId = user?.userId || userId;
    // 1) Use React Query to fetch the user's catalog using unified API
    const {
        data,
        isLoading,
        isError,
        error
    } = useQuery({
        queryKey: ['unifiedCatalog', effectiveUserId],
        queryFn: () => getUnifiedCatalog(effectiveUserId, {}),
        // Only fetch if we have a user ID and user is authenticated
        enabled: !!effectiveUserId && !!user,
        retry: (failureCount, error) => {
            // Don't retry on 401/403 errors
            if (error?.response?.status === 401 || error?.response?.status === 403) {
                return false;
            }
            return failureCount < 3; // retry up to 3 times for other errors
        }
    });
    if (isLoading) {
        return <div>Loading user products...</div>;
    }
    if (!user) {
        return <div className={styles.authError}>Authentication required to view catalog.</div>;
    }
    if (!effectiveUserId) {
        return <div className={styles.authError}>No user ID available.</div>;
    }
    if (isError) {
        const errorMessage = error?.response?.status === 401 
            ? "You are not authorized to view this catalog." 
            : "Error loading user products.";
        return <div className={styles.errorMessage}>{errorMessage}</div>;
    }
    const userProducts = data?.items || data?.products || [];
    if (!userProducts.length) {
        return (
            <section className={styles.sectionContainer}>
                <div className={styles.noDataMessage}>
                    No products found for user {effectiveUserId}.
                </div>
            </section>
        );
    }
    return (
        <section className={styles.sectionContainer}>
            {/* If you want a title, you can include the h2 below
          <h2 className={styles.sectionTitle}>User Products</h2>
      */}
            <div className={styles.listContainer}>
            </div>
        </section>
    );
});
ItemList.displayName = 'ItemList';
ItemList.propTypes = {
    userId: PropTypes.string.isRequired,
};
export default ItemList;
