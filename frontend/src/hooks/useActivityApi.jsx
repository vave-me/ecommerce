// useActivityApi.jsx
import {useCallback} from 'react';
import {toast} from 'react-toastify';
import {addInteraction, createActivity, getActivity} from "../api/client/activityApi";
export default function useActivityApi() {
    const showToast = useCallback((type, message) => {
        const options = {theme: 'colored'};
        switch (type) {
            case 'success':
                toast.success(message, options);
                break;
            case 'info':
                toast.info(message, options);
                break;
            case 'error':
                toast.error(message, options);
                break;
            case 'warn':
                toast.warn(message, options);
                break;
            default:
                toast(message, options);
        }
    }, []);
    const getOrCreateActivityId = useCallback(async (userId) => {
        if (!userId) return null;
        // Check localStorage for an existing activity ID with user-specific key
        const activityKey = `activityId_${userId}`;
        let activityId = localStorage.getItem(activityKey);
        if (!activityId) {
            try {
                // Try to fetch an existing activity for this user
                // GET /api/activity/{userId}
                const data = await getActivity(userId);
                if (data?.activityId) {
                    // If the server returns an existing activity
                    activityId = data.activityId;
                    localStorage.setItem(activityKey, activityId);
                } else {
                    // Otherwise, create a new one
                    const newActivity = await createActivity(userId);
                    activityId = newActivity.id;
                    localStorage.setItem(activityKey, activityId);
                }
            } catch (error) {
                // If user doesn't have an activity, create one
                try {
                    const newActivity = await createActivity(userId);
                    activityId = newActivity.id;
                    localStorage.setItem(activityKey, activityId);
                } catch (createError) {
                    return null;
                }
            }
        }
        return activityId;
    }, []);
    const handleLike = useCallback(async (productId, userId) => {
        if (!userId) {
            showToast('warn', 'Please log in to like products.');
            return;
        }
        try {
            const activityId = await getOrCreateActivityId(userId);
            if (!activityId) throw new Error('Could not retrieve or create an activity ID');
            // POST /api/activity/interactions
            await addInteraction(activityId, productId, 'product', 'like');
            showToast('success', 'Product liked!');
        } catch (error) {
            showToast('error', 'Failed to like the product.');
        }
    }, [showToast, getOrCreateActivityId]);
    const handleDislike = useCallback(async (productId, userId) => {
        if (!userId) {
            showToast('warn', 'Please log in to dislike products.');
            return;
        }
        try {
            const activityId = await getOrCreateActivityId(userId);
            if (!activityId) throw new Error('Could not retrieve or create an activity ID');
            await addInteraction(activityId, productId, 'product', 'dislike');
            showToast('success', 'Product disliked!');
        } catch (error) {
            showToast('error', 'Failed to dislike the product.');
        }
    }, [showToast, getOrCreateActivityId]);
    return {
        handleLike,
        handleDislike,
    };
}
