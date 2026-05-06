// Quick debug utility for testing catalog API
// Run in browser console: testCatalogAPI()
window.testCatalogAPI = async function() {
    try {
        // Get current user info
        const authContextElement = document.querySelector('[data-auth-context]');
        const userId = authContextElement?.dataset?.userId;
        if (!userId) {
            return;
        }
        // Import the API function dynamically
        const { getUnifiedCatalog } = await import('../api/client/searchApi');
        // Test the API call
        const response = await getUnifiedCatalog(userId, {});
        if (response?.items && response.items.length > 0) {
            response.items.slice(0, 3).forEach((item, i) => {
            });
        } else {
        }
        return response;
    } catch (error) {
        return { error: error.message };
    }
};
 