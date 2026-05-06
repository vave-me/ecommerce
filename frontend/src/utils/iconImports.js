/**
 * Centralized Icon Imports for Bundle Optimization
 * 
 * This file centralizes ALL icon imports across the application to improve:
 * - Tree-shaking and bundle size optimization
 * - Consistent icon usage
 * - Easy maintenance and updates
 * - Performance optimization
 * 
 * USAGE: Always import icons from this file instead of directly from icon libraries
 * Example: import { Search, User, FaCheckCircle } from '../utils/iconImports';
 */
// ===================================
// LUCIDE REACT ICONS (Primary choice for modern, lightweight icons)
// ===================================
import {
  // Navigation & Movement
  Search,
  ArrowLeft,
  ArrowRight,
  ChevronDown,
  ChevronUp,
  ChevronLeft,
  ChevronRight,
  ExternalLink,
  RefreshCw,
  // User Interface
  X,
  Plus,
  Minus,
  Edit,
  Check,
  CheckCircle,
  Eye,
  Upload,
  Download,
  Share2,
  MoreHorizontal,
  Settings,
  // User & Profile
  User,
  UserCheck,
  UserPlus,
  Users,
  Bell,
  // Content & Media
  Image,
  Video,
  Camera,
  Play,
  MessageCircle,
  MessagesSquare,
  Mail,
  Phone,
  // Business & Commerce
  ShoppingCart,
  ShoppingBag,
  Euro,
  CircleDollarSign,
  CreditCard,
  Briefcase,
  BriefcaseBusiness,
  // Location & Time
  MapPin,
  Map,
  Clock,
  Clock3,
  Calendar,
  CalendarRange,
  // Status & Feedback
  Heart,
  Star,
  ThumbsUp,
  ThumbsDown,
  AlertCircle,
  AlertTriangle,
  Trash2,
  Flag,
  // Security & Privacy
  Lock,
  Shield,
  ShieldCheck,
  LogOut,
  // Technology & Tools
  Bot,
  Smartphone,
  Laptop,
  Globe,
  Terminal,
  Wrench,
  // Categories & Tags
  Car,
  Home,
  Truck,
  Factory,
  Building,
  Building2,
  Tag,
  Hash,
  // Specialized Icons
  Flame,
  Award,
  Bookmark,
  Grid3X3,
  Loader,
  Loader2,
  Mic,
  MicOff,
  Send,
  Zap,
  Info,
  // Additional specific icons found in usage
  AlarmClock,
  BadgeCheck,
  BarChart,
  Clapperboard,
  Fuel,
  GraduationCap,
  HelpCircle,
  KeyRound,
  Link2,
  Milestone,
  Paperclip,
  PiggyBank,
  Repeat2,
  Scissors,
  Server,
  Snowflake,
  Target,
  Utensils,
  // Icon aliases used in components
  Settings as SettingsIcon,
  X as CloseIcon,
  Paperclip as PaperclipIcon,
  Send as SendIcon,
  Smile as SmileIcon,
  Plane as PlaneIcon
} from '@/icons';
// ===================================
// FONT AWESOME ICONS (For specific UI patterns and branded icons)
// ===================================
import {
  // Status & Feedback Icons
  FaCheckCircle,
  FaTimesCircle,
  FaExclamationCircle,
  FaExclamationTriangle,
  FaInfoCircle,
  FaSpinner,
  // UI Actions
  FaPaperPlane,
  FaPaperclip,
  FaSmile,
  FaSave,
  FaRedo,
  FaUndo,
  FaExpand,
  FaPlay,
  FaPause,
  FaGripVertical,
  FaBars,
  // Text Editor Icons
  FaBold,
  FaItalic,
  FaUnderline,
  FaListUl,
  FaListOl,
  FaQuoteRight,
  FaAlignLeft,
  FaAlignCenter,
  FaAlignRight,
  FaTable,
  FaHeading,
  FaUnlink,
  // Business & Commerce
  FaBoxOpen,
  FaHandshake,
  FaPercentage,
  FaClipboardList,
  // Emotions & Reactions
  FaAngry,
  FaLaugh,
  FaSadTear,
  // Media & Content
  FaNewspaper,
  FaRegCommentDots,
  FaStarHalfAlt,
  FaRegStar,
  // Display & Appearance
  FaColumns,
  FaMoon,
  FaPaintBrush,
  FaSun,
  FaTextHeight,
  FaPalette,
  FaEye,
  // Security & Privacy
  FaBan,
  FaShieldAlt,
  FaUserShield,
  // External & Navigation
  FaExternalLinkAlt,
  FaHistory,
  // Social Media Brands
  FaFacebook,
  FaTwitter,
  FaInstagram,
  FaLinkedin,
  FaGithub
} from '../icons'; // Migrated from react-icons/fa to centralized icons
// ===================================
// OTHER REACT ICONS
// ===================================
// Migrated to centralized icons system
// These icons are now available from '../icons'
import { AiOutlineHeart, TiWaves, FiFilter } from '../icons';
// ===================================
// EXPORTS - Organized by Category
// ===================================
// Navigation & Movement
export {
  Search,
  ArrowLeft,
  ArrowRight,
  ChevronDown,
  ChevronUp,
  ChevronLeft,
  ChevronRight,
  ExternalLink,
  RefreshCw,
  FaExternalLinkAlt,
  FaHistory
};
// User Interface Controls
export {
  X,
  CloseIcon,
  Plus,
  Minus,
  Edit,
  Check,
  CheckCircle,
  Eye,
  Upload,
  Download,
  Share2,
  MoreHorizontal,
  Settings,
  SettingsIcon,
  FaBars,
  FaExpand
};
// User & Profile
export {
  User,
  UserCheck,
  UserPlus,
  Users,
  Bell,
  FaUserShield
};
// Content & Media
export {
  Image,
  Video,
  Camera,
  Play,
  FaPlay,
  FaPause,
  MessageCircle,
  MessagesSquare,
  Mail,
  Phone,
  FaRegCommentDots,
  FaNewspaper
};
// Business & Commerce
export {
  ShoppingCart,
  ShoppingBag,
  Euro,
  CircleDollarSign,
  CreditCard,
  Briefcase,
  BriefcaseBusiness,
  FaBoxOpen,
  FaHandshake,
  FaPercentage,
  FaClipboardList
};
// Location & Time
export {
  MapPin,
  Map,
  Clock,
  Clock3,
  Calendar,
  CalendarRange,
  AlarmClock,
  Target
};
// Status & Feedback
export {
  Heart,
  AiOutlineHeart,
  Star,
  FaStarHalfAlt,
  FaRegStar,
  ThumbsUp,
  ThumbsDown,
  AlertCircle,
  AlertTriangle,
  Trash2,
  Flag,
  FaCheckCircle,
  FaTimesCircle,
  FaExclamationCircle,
  FaExclamationTriangle,
  FaInfoCircle
};
// Security & Privacy
export {
  Lock,
  Shield,
  ShieldCheck,
  LogOut,
  FaBan,
  FaShieldAlt
};
// Technology & Tools
export {
  Bot,
  Smartphone,
  Laptop,
  Globe,
  Terminal,
  Wrench,
  Server,
  KeyRound,
  Link2,
  Scissors
};
// Categories & Classifications
export {
  Car,
  Home,
  Truck,
  Factory,
  Building,
  Building2,
  Tag,
  Hash,
  Utensils,
  GraduationCap,
  Fuel
};
// Actions & Interactions
export {
  Send,
  SendIcon,
  Paperclip,
  PaperclipIcon,
  FaPaperPlane,
  FaPaperclip,
  FaSmile,
  SmileIcon,
  FaSave,
  FaRedo,
  FaUndo,
  FaGripVertical
};
// Text Editor & Formatting
export {
  FaBold,
  FaItalic,
  FaUnderline,
  FaListUl,
  FaListOl,
  FaQuoteRight,
  FaAlignLeft,
  FaAlignCenter,
  FaAlignRight,
  FaTable,
  FaHeading,
  FaUnlink
};
// Emotions & Reactions
export {
  FaAngry,
  FaLaugh,
  FaSadTear
};
// Display & Appearance
export {
  FaColumns,
  FaMoon,
  FaPaintBrush,
  FaSun,
  FaTextHeight,
  FaPalette,
  FaEye
};
// Loading & Progress
export {
  Loader,
  Loader2,
  FaSpinner
};
// Media Controls
export {
  Mic,
  MicOff,
  PlaneIcon
};
// Special Effects & Decorative
export {
  Flame,
  TiWaves,
  Snowflake,
  Zap,
  Award,
  Bookmark,
  Grid3X3,
  Milestone,
  PiggyBank,
  Repeat2
};
// Information & Help
export {
  Info,
  HelpCircle
};
// Social Media Brands
export {
  FaFacebook,
  FaTwitter,
  FaInstagram,
  FaLinkedin,
  FaGithub
};
// Filters & Organization
export {
  FiFilter,
  Clapperboard,
  BarChart,
  BadgeCheck
};
// ===================================
// ICON MAPPING FOR DYNAMIC USAGE
// ===================================
export const iconMap = {
  // Navigation
  search: Search,
  arrowLeft: ArrowLeft,
  arrowRight: ArrowRight,
  chevronDown: ChevronDown,
  chevronUp: ChevronUp,
  chevronLeft: ChevronLeft,
  chevronRight: ChevronRight,
  // Actions
  close: X,
  add: Plus,
  remove: Minus,
  edit: Edit,
  check: Check,
  delete: Trash2,
  // Status
  success: FaCheckCircle,
  error: FaTimesCircle,
  warning: FaExclamationTriangle,
  info: FaInfoCircle,
  loading: FaSpinner,
  // User
  user: User,
  users: Users,
  profile: User,
  // Communication
  message: MessageCircle,
  mail: Mail,
  phone: Phone,
  notification: Bell,
  // Media
  image: Image,
  video: Video,
  camera: Camera,
  play: Play,
  // Business
  cart: ShoppingCart,
  money: Euro,
  business: Briefcase,
  // Location
  location: MapPin,
  map: Map,
  // Security
  lock: Lock,
  shield: Shield,
  // Technology
  bot: Bot,
  mobile: Smartphone,
  laptop: Laptop,
  globe: Globe
};
// ===================================
// UTILITY FUNCTIONS
// ===================================
/**
 * Get icon component by name from the icon map
 * @param {string} iconName - The name of the icon
 * @returns {Component|null} - The icon component or null if not found
 */
export const getIcon = (iconName) => {
  return iconMap[iconName] || null;
};
/**
 * Check if an icon exists in the icon map
 * @param {string} iconName - The name of the icon
 * @returns {boolean} - Whether the icon exists
 */
export const hasIcon = (iconName) => {
  return iconName in iconMap;
};
/**
 * Get all available icon names
 * @returns {string[]} - Array of all icon names
 */
export const getAvailableIcons = () => {
  return Object.keys(iconMap);
}; 