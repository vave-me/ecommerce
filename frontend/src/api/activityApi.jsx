// Activity types and configurations
export const getActivityTypes = () => [
    { value: 'like', label: 'Likes', icon: 'Heart' },
    { value: 'dislike', label: 'Dislikes', icon: 'ThumbsDown' },
    { value: 'view', label: 'Views', icon: 'Eye' },
    { value: 'purchase', label: 'Purchases', icon: 'ShoppingCart' },
    { value: 'wishlist', label: 'Wishlist', icon: 'Bookmark' },
    { value: 'review', label: 'Reviews', icon: 'MessageSquare' },
    { value: 'share', label: 'Shares', icon: 'Share' },
    { value: 'follow', label: 'Follows', icon: 'UserPlus' },
    { value: 'comment', label: 'Comments', icon: 'MessageCircle' },
    { value: 'rating', label: 'Ratings', icon: 'Star' }
];
export const getActionTypes = () => [
    { value: 'like', label: 'Liked', icon: 'ThumbsUp' },
    { value: 'dislike', label: 'Disliked', icon: 'ThumbsDown' },
    { value: 'love', label: 'Loved', icon: 'Heart' },
    { value: 'haha', label: 'Laughed', icon: 'Laugh' },
    { value: 'sad', label: 'Sad', icon: 'Frown' },
    { value: 'angry', label: 'Angry', icon: 'Angry' },
    { value: 'view', label: 'Viewed', icon: 'Eye' },
    { value: 'click', label: 'Clicked', icon: 'MousePointer' },
    { value: 'add_to_cart', label: 'Added to Cart', icon: 'ShoppingCart' },
    { value: 'purchase', label: 'Purchased', icon: 'CreditCard' },
    { value: 'wishlist_add', label: 'Added to Wishlist', icon: 'Bookmark' },
    { value: 'wishlist_remove', label: 'Removed from Wishlist', icon: 'BookmarkX' },
    { value: 'follow', label: 'Followed', icon: 'UserPlus' },
    { value: 'unfollow', label: 'Unfollowed', icon: 'UserMinus' },
    { value: 'share', label: 'Shared', icon: 'Share' },
    { value: 'comment', label: 'Commented', icon: 'MessageSquare' },
    { value: 'review', label: 'Reviewed', icon: 'Edit' },
    { value: 'rating', label: 'Rated', icon: 'Star' }
];
export const getItemTypes = () => [
    { value: 'product', label: 'Product', icon: 'Package' },
    { value: 'store', label: 'Store', icon: 'Store' },
    { value: 'user', label: 'User', icon: 'User' },
    { value: 'category', label: 'Category', icon: 'Folder' },
    { value: 'brand', label: 'Brand', icon: 'Tag' },
    { value: 'review', label: 'Review', icon: 'MessageSquare' },
    { value: 'article', label: 'Article', icon: 'FileText' },
    { value: 'promotion', label: 'Promotion', icon: 'Percent' }
];
/**
 * Get activity type configuration by value
 */
export const getActivityTypeConfig = (type) => {
    return getActivityTypes().find(t => t.value === type) || 
           { value: type, label: type, icon: 'Activity' };
};
/**
 * Get action type configuration by value
 */
export const getActionTypeConfig = (actionType) => {
    return getActionTypes().find(t => t.value === actionType) || 
           { value: actionType, label: actionType, icon: 'Activity' };
};
/**
 * Get item type configuration by value
 */
export const getItemTypeConfig = (itemType) => {
    return getItemTypes().find(t => t.value === itemType) || 
           { value: itemType, label: itemType, icon: 'Package' };
};
/**
 * Format activity message based on action and item types
 */
export const formatActivityMessage = (activity) => {
    const actionConfig = getActionTypeConfig(activity.actionType);
    const itemConfig = getItemTypeConfig(activity.itemType);
    if (activity.message) {
        return activity.message;
    }
    // Generate default message based on action and item type
    return `${actionConfig.label} ${itemConfig.label.toLowerCase()}`;
}; 