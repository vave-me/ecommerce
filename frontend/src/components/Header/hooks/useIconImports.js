/**
 * Optimized Icon Imports for Header Components
 * Centralizes icon imports to improve tree shaking and reduce bundle size
 */
// Core navigation icons - most commonly used
export const NavigationIcons = {
    Bell: () => import('@/icons').then(module => ({ default: module.Bell })),
    Heart: () => import('@/icons').then(module => ({ default: module.Heart })),
    MessageCircle: () => import('@/icons').then(module => ({ default: module.MessageCircle })),
    ShoppingBag: () => import('@/icons').then(module => ({ default: module.ShoppingBag })),
    User: () => import('@/icons').then(module => ({ default: module.User })),
    ChevronDown: () => import('@/icons').then(module => ({ default: module.ChevronDown })),
    X: () => import('@/icons').then(module => ({ default: module.X })),
    Home: () => import('@/icons').then(module => ({ default: module.Home })),
    Search: () => import('@/icons').then(module => ({ default: module.Search })),
    Filter: () => import('@/icons').then(module => ({ default: module.Filter })),
};
// Action icons - for add/create functionality
export const ActionIcons = {
    Tag: () => import('@/icons').then(module => ({ default: module.Tag })),
    LogOut: () => import('@/icons').then(module => ({ default: module.LogOut })),
    Mail: () => import('@/icons').then(module => ({ default: module.Mail })),
    Briefcase: () => import('@/icons').then(module => ({ default: module.Briefcase })),
    Car: () => import('@/icons').then(module => ({ default: module.Car })),
    Wrench: () => import('@/icons').then(module => ({ default: module.Wrench })),
    Video: () => import('@/icons').then(module => ({ default: module.Video })),
};
// Search and filter icons
export const SearchIcons = {
    MapPin: () => import('@/icons').then(module => ({ default: module.MapPin })),
    Clock: () => import('@/icons').then(module => ({ default: module.Clock })),
    Loader: () => import('@/icons').then(module => ({ default: module.Loader })),
    Building2: () => import('@/icons').then(module => ({ default: module.Building2 })),
    Map: () => import('@/icons').then(module => ({ default: module.Map })),
};
// Mode and feature icons
export const FeatureIcons = {
    Bot: () => import('@/icons').then(module => ({ default: module.Bot })),
    Grid3X3: () => import('@/icons').then(module => ({ default: module.Grid3X3 })),
    Sparkles: () => import('@/icons').then(module => ({ default: module.Sparkles })),
    TrendingUp: () => import('@/icons').then(module => ({ default: module.TrendingUp })),
};
// Hook for lazy loading icons
export const useLazyIcon = (iconName, iconCategory = 'NavigationIcons') => {
    const [IconComponent, setIconComponent] = React.useState(null);
    const [isLoading, setIsLoading] = React.useState(true);
    React.useEffect(() => {
        const loadIcon = async () => {
            try {
                setIsLoading(true);
                const iconMap = {
                    NavigationIcons,
                    ActionIcons,
                    SearchIcons,
                    FeatureIcons,
                };
                const iconLoader = iconMap[iconCategory]?.[iconName];
                if (iconLoader) {
                    const { default: Icon } = await iconLoader();
                    setIconComponent(() => Icon);
                }
            } catch (error) {
        // Error logged for debugging
        if (process.env.NODE_ENV === 'development') {
            console.error('Error:', error);
        }
    } finally {
                setIsLoading(false);
            }
        };
        loadIcon();
    }, [iconName, iconCategory]);
    return { IconComponent, isLoading };
}; 