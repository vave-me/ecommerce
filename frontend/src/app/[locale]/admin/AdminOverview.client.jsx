"use client";

import React, { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import { 
  BarChart3, 
  Users, 
  Package, 
  ShoppingBag, 
  MessageSquare, 
  Settings, 
  Activity,
  Bell,
  FileText,
  CreditCard,
  Truck,
  Tag,
  Heart,
  Database,
  Server,
  Shield,
  Mail,
  Calendar,
  Image,
  Store,
  Layers,
  TrendingUp,
  AlertTriangle,
  CheckCircle,
  Clock,
  DollarSign,
  ArrowRight,
  Eye,
  Edit,
  Plus
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import { useQuery } from '@tanstack/react-query';
import { getPlatformMetrics } from '@/api/adminApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import styles from './AdminOverview.module.css';

// Admin Quick Actions
const QUICK_ACTIONS = [
  { id: 'add-user', label: 'Add User', icon: Users, href: '/admin/users/add-admin', color: 'primary' },
  { id: 'add-product', label: 'Add Product', icon: Package, href: '/admin/products/add', color: 'success' },
  { id: 'view-orders', label: 'Recent Orders', icon: ShoppingBag, href: '/admin/orders', color: 'info' },
  { id: 'system-health', label: 'System Health', icon: Activity, href: '/admin/analytics', color: 'warning' }
];

// Admin Sections Configuration
const ADMIN_SECTIONS = [
  {
    title: 'Core Management',
    sections: [
      { name: 'Dashboard', href: '/admin/dashboard', icon: BarChart3, description: 'System overview and metrics' },
      { name: 'Users', href: '/admin/users', icon: Users, description: 'User management and roles' },
      { name: 'Products', href: '/admin/products', icon: Package, description: 'Product catalog management' },
      { name: 'Orders', href: '/admin/orders', icon: ShoppingBag, description: 'Order processing and fulfillment' }
    ]
  },
  {
    title: 'Content & Communication',
    sections: [
      { name: 'Posts', href: '/admin/posts', icon: FileText, description: 'Content management system' },
      { name: 'Comments', href: '/admin/comments', icon: MessageSquare, description: 'Comment moderation' },
      { name: 'Messages', href: '/admin/messages', icon: Mail, description: 'Internal messaging system' },
      { name: 'Notifications', href: '/admin/notifications', icon: Bell, description: 'Notification management' }
    ]
  },
  {
    title: 'E-commerce & Sales',
    sections: [
      { name: 'Payments', href: '/admin/payments', icon: CreditCard, description: 'Payment processing' },
      { name: 'Shipping', href: '/admin/shipping', icon: Truck, description: 'Shipping management' },
      { name: 'Offers', href: '/admin/offers', icon: Tag, description: 'Promotions and discounts' },
      { name: 'Wishlists', href: '/admin/wishlists', icon: Heart, description: 'Customer wishlists' }
    ]
  },
  {
    title: 'System & Operations',
    sections: [
      { name: 'Analytics', href: '/admin/analytics', icon: TrendingUp, description: 'Platform analytics' },
      { name: 'Activity', href: '/admin/activity', icon: Activity, description: 'System activity monitoring' },
      { name: 'Database', href: '/admin/database', icon: Database, description: 'Database management tools' },
      { name: 'ERP Integration', href: '/admin/erp', icon: Server, description: 'ERP system integration' }
    ]
  },
  {
    title: 'Advanced Features',
    sections: [
      { name: 'Merchant Center', href: '/admin/merchant', icon: Store, description: 'Merchant management' },
      { name: 'Content Moderation', href: '/admin/moderation', icon: Shield, description: 'Content review and moderation' },
      { name: 'Scheduler', href: '/admin/scheduler', icon: Calendar, description: 'Task scheduling system' },
      { name: 'Media Library', href: '/admin/media', icon: Image, description: 'Media file management' }
    ]
  },
  {
    title: 'Configuration',
    sections: [
      { name: 'Categories', href: '/admin/categories', icon: Layers, description: 'Category hierarchy management' },
      { name: 'Reviews', href: '/admin/reviews', icon: MessageSquare, description: 'Review management' },
      { name: 'Newsletters', href: '/admin/newsletters', icon: Mail, description: 'Newsletter campaigns' },
      { name: 'Settings', href: '/admin/settings', icon: Settings, description: 'System configuration' }
    ]
  }
];

export default function AdminOverview() {
  const t = useTranslations('admin');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();
  const [recentActivity, setRecentActivity] = useState([]);

  // Fetch platform metrics
  const { data: metrics, isLoading: metricsLoading } = useQuery({
    queryKey: ['platform-metrics'],
    queryFn: getPlatformMetrics,
    staleTime: 30000 // 30 seconds
  });

  if (metricsLoading) {
    return (
      <div className={styles.loadingContainer}>
        <LoadingSpinner />
        <p className={styles.loadingText}>Loading admin overview...</p>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Header */}
      <div className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.headerText}>
            <h1 className={styles.title}>Admin Overview</h1>
            <p className={styles.subtitle}>
              Welcome back, {user?.name}. Manage your platform with enterprise-grade tools.
            </p>
          </div>
          <div className={styles.headerActions}>
            <Link href="/admin/dashboard" className={styles.dashboardButton}>
              <BarChart3 size={16} />
              View Dashboard
            </Link>
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className={styles.quickActions}>
        <h2 className={styles.sectionTitle}>Quick Actions</h2>
        <div className={styles.actionsGrid}>
          {QUICK_ACTIONS.map((action) => (
            <Link key={action.id} href={action.href} className={`${styles.actionCard} ${styles[action.color]}`}>
              <action.icon className={styles.actionIcon} size={20} />
              <span className={styles.actionLabel}>{action.label}</span>
              <ArrowRight className={styles.actionArrow} size={16} />
            </Link>
          ))}
        </div>
      </div>

      {/* System Metrics */}
      {metrics && (
        <div className={styles.metricsSection}>
          <h2 className={styles.sectionTitle}>System Overview</h2>
          <div className={styles.metricsGrid}>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <Users size={24} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{metrics.totalUsers || 0}</div>
                <div className={styles.metricLabel}>Total Users</div>
              </div>
            </div>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <Package size={24} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{metrics.totalProducts || 0}</div>
                <div className={styles.metricLabel}>Products</div>
              </div>
            </div>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <ShoppingBag size={24} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>{metrics.totalOrders || 0}</div>
                <div className={styles.metricLabel}>Orders</div>
              </div>
            </div>
            <div className={styles.metricCard}>
              <div className={styles.metricIcon}>
                <DollarSign size={24} />
              </div>
              <div className={styles.metricContent}>
                <div className={styles.metricValue}>${metrics.totalRevenue || 0}</div>
                <div className={styles.metricLabel}>Revenue</div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Admin Sections */}
      <div className={styles.sectionsContainer}>
        <h2 className={styles.sectionTitle}>Administration Sections</h2>
        {ADMIN_SECTIONS.map((group) => (
          <div key={group.title} className={styles.sectionGroup}>
            <h3 className={styles.groupTitle}>{group.title}</h3>
            <div className={styles.sectionsGrid}>
              {group.sections.map((section) => (
                <Link key={section.name} href={section.href} className={styles.sectionCard}>
                  <div className={styles.sectionIcon}>
                    <section.icon size={20} />
                  </div>
                  <div className={styles.sectionContent}>
                    <h4 className={styles.sectionName}>{section.name}</h4>
                    <p className={styles.sectionDescription}>{section.description}</p>
                  </div>
                  <ArrowRight className={styles.sectionArrow} size={16} />
                </Link>
              ))}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
} 