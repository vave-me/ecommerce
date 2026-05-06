/**
 * Dynamic Component Imports for Code Splitting
 * 
 * This file implements dynamic imports for large components to reduce initial bundle size
 * and improve performance through code splitting.
 */
import dynamic from 'next/dynamic';
// Modal components (large and not needed on initial load)
export const ProductModal = dynamic(
  () => import('../features/CreateProductModal/ProductModal'),
  {
    loading: () => <div className="modal-loading">Loading...</div>,
    ssr: false
  }
);
export const VehicleModal = dynamic(
  () => import('../features/CreateVehicleModal/VehicleModal'),
  {
    loading: () => <div className="modal-loading">Loading...</div>,
    ssr: false
  }
);
export const ServiceModal = dynamic(
  () => import('../features/CreateServiceModal/ServiceModal'),
  {
    loading: () => <div className="modal-loading">Loading...</div>,
    ssr: false
  }
);
export const JobModal = dynamic(
  () => import('../features/CreateJobModal/JobModal'),
  {
    loading: () => <div className="modal-loading">Loading...</div>,
    ssr: false
  }
);
export const DealModal = dynamic(
  () => import('../features/CreateDealModal/DealModal'),
  {
    loading: () => <div className="modal-loading">Loading...</div>,
    ssr: false
  }
);
export const PropertyModal = dynamic(
  () => import('../features/CreatePropertyModal/PropertyModal'),
  {
    loading: () => <div className="modal-loading">Loading...</div>,
    ssr: false
  }
);
// Large components that can be lazy loaded
export const SearchBar = dynamic(
  () => import('../Header/SearchBar'),
  {
    loading: () => <div className="search-loading">🔍</div>,
    ssr: true // Keep SSR for SEO
  }
);
export const CategoryTree = dynamic(
  () => import('../Category/CategoryTree'),
  {
    loading: () => <div className="category-loading">Loading categories...</div>,
    ssr: false
  }
);
export const FeedProvider = dynamic(
  () => import('../Feed/FeedProvider.client'),
  {
    loading: () => <div className="feed-loading">Loading feed...</div>,
    ssr: false
  }
);
export const TextEditor = dynamic(
  () => import('../TextEditor/TextEditor'),
  {
    loading: () => <div className="editor-loading">Loading editor...</div>,
    ssr: false
  }
);
// Chart components (heavy libraries)
export const PerformanceChart = dynamic(
  () => import('../Charts/PerformanceChart'),
  {
    loading: () => <div className="chart-loading">Loading chart...</div>,
    ssr: false
  }
);
// Map components (heavy libraries)
export const LocationMap = dynamic(
  () => import('../Maps/LocationMap'),
  {
    loading: () => <div className="map-loading">Loading map...</div>,
    ssr: false
  }
);
// Video player (heavy library)
export const VideoPlayer = dynamic(
  () => import('../Media/VideoPlayer'),
  {
    loading: () => <div className="video-loading">Loading player...</div>,
    ssr: false
  }
);
// Image gallery (can be large with many images)
export const ImageGallery = dynamic(
  () => import('../Media/ImageGallery'),
  {
    loading: () => <div className="gallery-loading">Loading gallery...</div>,
    ssr: false
  }
);
// Admin components (only needed for admin users)
export const AdminPanel = dynamic(
  () => import('../Admin/AdminPanel'),
  {
    loading: () => <div className="admin-loading">Loading admin panel...</div>,
    ssr: false
  }
);
export const UserManagement = dynamic(
  () => import('../Admin/UserManagement'),
  {
    loading: () => <div className="admin-loading">Loading user management...</div>,
    ssr: false
  }
);
// Analytics components (not critical for initial load)
export const AnalyticsDashboard = dynamic(
  () => import('../Analytics/AnalyticsDashboard'),
  {
    loading: () => <div className="analytics-loading">Loading analytics...</div>,
    ssr: false
  }
);
// Social features (can be loaded later)
export const SocialShare = dynamic(
  () => import('../Shared/ShareComponent'),
  {
    loading: () => <div className="social-loading">Loading share options...</div>,
    ssr: false
  }
);
export const ShareComponent = dynamic(
  () => import('../Shared/ShareComponent'),
  {
    loading: () => <div className="share-loading">Loading share component...</div>,
    ssr: false
  }
);
export const CommentSystem = dynamic(
  () => import('../Comments/CommentSystem'),
  {
    loading: () => <div className="comments-loading">Loading comments...</div>,
    ssr: false
  }
);
// Payment components (only needed during checkout)
export const PaymentForm = dynamic(
  () => import('../Payment/PaymentForm'),
  {
    loading: () => <div className="payment-loading">Loading payment form...</div>,
    ssr: false
  }
);
export const StripeCheckout = dynamic(
  () => import('../Payment/StripeCheckout'),
  {
    loading: () => <div className="payment-loading">Loading checkout...</div>,
    ssr: false
  }
);
// Notification components
export const NotificationCenter = dynamic(
  () => import('../Notifications/NotificationCenter'),
  {
    loading: () => <div className="notification-loading">Loading notifications...</div>,
    ssr: false
  }
);
// Settings components
export const UserSettings = dynamic(
  () => import('../Settings/UserSettings'),
  {
    loading: () => <div className="settings-loading">Loading settings...</div>,
    ssr: false
  }
);
export const PrivacySettings = dynamic(
  () => import('../Settings/PrivacySettings'),
  {
    loading: () => <div className="settings-loading">Loading privacy settings...</div>,
    ssr: false
  }
);
// Help and support components
export const HelpCenter = dynamic(
  () => import('../Help/HelpCenter'),
  {
    loading: () => <div className="help-loading">Loading help center...</div>,
    ssr: false
  }
);
export const ContactForm = dynamic(
  () => import('../Contact/ContactForm'),
  {
    loading: () => <div className="contact-loading">Loading contact form...</div>,
    ssr: false
  }
);
// Export all dynamic components for easy importing
export default {
  // Modals
  ProductModal,
  VehicleModal,
  ServiceModal,
  JobModal,
  DealModal,
  PropertyModal,
  // Large components
  SearchBar,
  CategoryTree,
  FeedProvider,
  TextEditor,
  // Heavy libraries
  PerformanceChart,
  LocationMap,
  VideoPlayer,
  ImageGallery,
  // Admin
  AdminPanel,
  UserManagement,
  // Analytics
  AnalyticsDashboard,
  // Social
  SocialShare,
  ShareComponent,
  CommentSystem,
  // Payment
  PaymentForm,
  StripeCheckout,
  // Notifications
  NotificationCenter,
  // Settings
  UserSettings,
  PrivacySettings,
  // Help
  HelpCenter,
  ContactForm,
}; 