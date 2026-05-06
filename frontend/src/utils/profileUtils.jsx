/**
 * Profile utilities for generating correct profile links
 * Optimized for Next.js App Router with locale support
 */
/**
 * Generates the appropriate profile link based on the user ID and username
 * 
 * @param {string} userId - The ID of the user
 * @param {string|null} currentUserId - The ID of the current logged-in user (if any)
 * @param {string} username - The username for pretty URLs (optional but recommended)
 * @returns {string} The appropriate profile URL
 */
export const getProfileLink = (userId, currentUserId, username = null) => {
    // If viewing own profile, go to /user (private profile)
    if (userId === currentUserId) {
        return '/user';
    }
    // For public profiles, use username if available, otherwise userId
    const profileIdentifier = username || userId;
    return `/profile/${profileIdentifier}`;
};
/**
 * Generates a direct profile link using username (preferred method)
 * 
 * @param {string} username - The username
 * @returns {string} The profile URL
 */
export const getUsernameProfileLink = (username) => {
    if (!username) return '/profile/unknown';
    return `/profile/${username}`;
};
/**
 * Generates a legacy profile link using user ID (for backward compatibility)
 * 
 * @param {string} userId - The user ID  
 * @returns {string} The profile URL
 */
export const getUserIdProfileLink = (userId) => {
    if (!userId) return '/profile/unknown';
    return `/profile/slug/${userId}`;
}; 