"use client";
export const dynamic = 'force-dynamic';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  BookOpen,
  Rocket,
  Store,
  Package,
  Camera,
  DollarSign,
  TrendingUp,
  Users,
  MessageCircle,
  BarChart3,
  Settings,
  Shield,
  ChevronRight,
  ChevronDown,
  Search,
  Download,
  PlayCircle,
  CheckCircle,
  AlertCircle,
  Lightbulb,
  Zap,
  Globe,
  Heart,
  Star,
  ArrowRight,
  FileText,
  Edit3,
  ShoppingCart,
  Truck,
  CreditCard,
  Mail,
  Hash,
  Target,
  Award,
  HelpCircle
} from 'lucide-react';
import styles from './VendorGuide.module.css';

export default function VendorGuidePage() {
  const t = useTranslations('vendorGuide');
  const router = useRouter();
  const [activeSection, setActiveSection] = useState('getting-started');
  const [expandedTopics, setExpandedTopics] = useState({});
  const [searchQuery, setSearchQuery] = useState('');

  const toggleTopic = (topicId) => {
    setExpandedTopics(prev => ({
      ...prev,
      [topicId]: !prev[topicId]
    }));
  };

  const sections = [
    {
      id: 'getting-started',
      icon: Rocket,
      title: t('section1', 'Getting Started'),
      topics: [
        {
          id: 'account-setup',
          title: t('topic1', 'Account Setup & Verification'),
          icon: Settings,
          content: [
            t('content1a', 'Creating your business account'),
            t('content1b', 'Identity verification process'),
            t('content1c', 'Business documentation requirements'),
            t('content1d', 'Setting up payment methods'),
            t('content1e', 'Tax information and compliance')
          ]
        },
        {
          id: 'store-customization',
          title: t('topic2', 'Store Customization'),
          icon: Store,
          content: [
            t('content2a', 'Choosing your store theme'),
            t('content2b', 'Customizing colors and branding'),
            t('content2c', 'Adding your logo and banner'),
            t('content2d', 'Creating store categories'),
            t('content2e', 'Setting up store policies')
          ]
        },
        {
          id: 'first-product',
          title: t('topic3', 'Adding Your First Product'),
          icon: Package,
          content: [
            t('content3a', 'Product information requirements'),
            t('content3b', 'Writing compelling descriptions'),
            t('content3c', 'Pricing strategies'),
            t('content3d', 'Inventory management basics'),
            t('content3e', 'Product variations and options')
          ]
        }
      ]
    },
    {
      id: 'product-management',
      icon: Package,
      title: t('section2', 'Product Management'),
      topics: [
        {
          id: 'product-photos',
          title: t('topic4', 'Product Photography'),
          icon: Camera,
          content: [
            t('content4a', 'Image requirements and specifications'),
            t('content4b', 'Photography tips for products'),
            t('content4c', 'Image editing best practices'),
            t('content4d', 'Using 360° product views'),
            t('content4e', 'Video content guidelines')
          ]
        },
        {
          id: 'seo-optimization',
          title: t('topic5', 'SEO & Search Optimization'),
          icon: Target,
          content: [
            t('content5a', 'Keyword research for products'),
            t('content5b', 'Optimizing product titles'),
            t('content5c', 'Meta descriptions and tags'),
            t('content5d', 'Category selection strategies'),
            t('content5e', 'Internal linking best practices')
          ]
        },
        {
          id: 'inventory-management',
          title: t('topic6', 'Advanced Inventory Management'),
          icon: BarChart3,
          content: [
            t('content6a', 'Stock level monitoring'),
            t('content6b', 'Automated reorder points'),
            t('content6c', 'Multi-channel inventory sync'),
            t('content6d', 'Seasonal inventory planning'),
            t('content6e', 'Dead stock management')
          ]
        }
      ]
    },
    {
      id: 'selling-strategies',
      icon: TrendingUp,
      title: t('section3', 'Selling Strategies'),
      topics: [
        {
          id: 'pricing-strategies',
          title: t('topic7', 'Pricing & Promotions'),
          icon: DollarSign,
          content: [
            t('content7a', 'Competitive pricing analysis'),
            t('content7b', 'Dynamic pricing strategies'),
            t('content7c', 'Creating effective promotions'),
            t('content7d', 'Bundle deals and discounts'),
            t('content7e', 'Flash sales and limited offers')
          ]
        },
        {
          id: 'cross-selling',
          title: t('topic8', 'Cross-selling & Upselling'),
          icon: ShoppingCart,
          content: [
            t('content8a', 'Product recommendation setup'),
            t('content8b', 'Creating product bundles'),
            t('content8c', 'Upsell strategies at checkout'),
            t('content8d', 'Related products optimization'),
            t('content8e', 'Post-purchase upselling')
          ]
        },
        {
          id: 'international-selling',
          title: t('topic9', 'International Expansion'),
          icon: Globe,
          content: [
            t('content9a', 'Currency and language setup'),
            t('content9b', 'International shipping options'),
            t('content9c', 'Customs and duties handling'),
            t('content9d', 'Regional pricing strategies'),
            t('content9e', 'Local payment methods')
          ]
        }
      ]
    },
    {
      id: 'community-engagement',
      icon: Users,
      title: t('section4', 'Community & Engagement'),
      topics: [
        {
          id: 'building-followers',
          title: t('topic10', 'Building Your Following'),
          icon: Heart,
          content: [
            t('content10a', 'Creating engaging content'),
            t('content10b', 'Follower growth strategies'),
            t('content10c', 'Community interaction best practices'),
            t('content10d', 'Exclusive follower benefits'),
            t('content10e', 'Loyalty program setup')
          ]
        },
        {
          id: 'content-marketing',
          title: t('topic11', 'Content Marketing'),
          icon: Edit3,
          content: [
            t('content11a', 'Blog post creation guidelines'),
            t('content11b', 'Video content strategies'),
            t('content11c', 'Live streaming best practices'),
            t('content11d', 'Newsletter campaigns'),
            t('content11e', 'Social media integration')
          ]
        },
        {
          id: 'customer-communication',
          title: t('topic12', 'Customer Communication'),
          icon: MessageCircle,
          content: [
            t('content12a', 'Live chat best practices'),
            t('content12b', 'Response time optimization'),
            t('content12c', 'Handling customer inquiries'),
            t('content12d', 'Review management strategies'),
            t('content12e', 'Building customer relationships')
          ]
        }
      ]
    },
    {
      id: 'fulfillment-shipping',
      icon: Truck,
      title: t('section5', 'Fulfillment & Shipping'),
      topics: [
        {
          id: 'shipping-setup',
          title: t('topic13', 'Shipping Configuration'),
          icon: Truck,
          content: [
            t('content13a', 'Shipping zones and rates'),
            t('content13b', 'Carrier integration setup'),
            t('content13c', 'Free shipping strategies'),
            t('content13d', 'Express shipping options'),
            t('content13e', 'Shipping calculator setup')
          ]
        },
        {
          id: 'order-fulfillment',
          title: t('topic14', 'Order Processing'),
          icon: Package,
          content: [
            t('content14a', 'Order management workflow'),
            t('content14b', 'Automated order processing'),
            t('content14c', 'Tracking number generation'),
            t('content14d', 'Custom packaging options'),
            t('content14e', 'Returns and exchanges')
          ]
        }
      ]
    },
    {
      id: 'analytics-growth',
      icon: BarChart3,
      title: t('section6', 'Analytics & Growth'),
      topics: [
        {
          id: 'performance-metrics',
          title: t('topic15', 'Understanding Your Metrics'),
          icon: BarChart3,
          content: [
            t('content15a', 'Key performance indicators'),
            t('content15b', 'Sales analytics dashboard'),
            t('content15c', 'Customer behavior analysis'),
            t('content15d', 'Conversion rate optimization'),
            t('content15e', 'Revenue forecasting')
          ]
        },
        {
          id: 'growth-strategies',
          title: t('topic16', 'Scaling Your Business'),
          icon: TrendingUp,
          content: [
            t('content16a', 'Growth planning framework'),
            t('content16b', 'Automation opportunities'),
            t('content16c', 'Team building and delegation'),
            t('content16d', 'Multi-channel expansion'),
            t('content16e', 'Investment and funding options')
          ]
        }
      ]
    }
  ];

  const quickActions = [
    {
      icon: PlayCircle,
      title: t('quickAction1', 'Watch Video Tutorial'),
      description: t('quickAction1Desc', '15-minute platform overview'),
      link: '/tutorials/platform-overview'
    },
    {
      icon: Download,
      title: t('quickAction2', 'Download Vendor Handbook'),
      description: t('quickAction2Desc', 'Complete PDF guide (120 pages)'),
      link: '/resources/vendor-handbook.pdf'
    },
    {
      icon: HelpCircle,
      title: t('quickAction3', 'Contact Vendor Support'),
      description: t('quickAction3Desc', 'Get help from our team'),
      link: '/support/vendors'
    }
  ];

  const bestPractices = [
    {
      icon: Star,
      title: t('bestPractice1', 'Complete Your Profile'),
      description: t('bestPractice1Desc', 'Stores with complete profiles see 40% more traffic')
    },
    {
      icon: Camera,
      title: t('bestPractice2', 'Use High-Quality Images'),
      description: t('bestPractice2Desc', 'Products with 5+ images convert 3x better')
    },
    {
      icon: MessageCircle,
      title: t('bestPractice3', 'Respond Quickly'),
      description: t('bestPractice3Desc', 'Reply to customer messages within 1 hour')
    },
    {
      icon: Heart,
      title: t('bestPractice4', 'Build Your Community'),
      description: t('bestPractice4Desc', 'Engaged followers become repeat customers')
    }
  ];

  // Filter content based on search
  const filteredSections = sections.map(section => ({
    ...section,
    topics: section.topics.filter(topic => 
      topic.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      topic.content.some(item => item.toLowerCase().includes(searchQuery.toLowerCase()))
    )
  })).filter(section => section.topics.length > 0);

  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.headerInfo}>
            <h1 className={styles.title}>
              {t('pageTitle', 'Complete Vendor Guide')}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', 'Everything you need to know to succeed on our platform')}
            </p>
            
            {/* Search Bar */}
            <div className={styles.searchBar}>
              <Search className={styles.searchIcon} size={20} />
              <input
                type="text"
                placeholder={t('searchPlaceholder', 'Search the guide...')}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className={styles.searchInput}
              />
            </div>
          </div>
        </div>
      </header>

      {/* Quick Actions */}
      <section className={styles.quickActions}>
        <div className={styles.quickActionsContainer}>
          {quickActions.map((action, index) => (
            <button
              key={index}
              onClick={() => router.push(action.link)}
              className={styles.quickActionCard}
            >
              <action.icon className={styles.quickActionIcon} size={24} />
              <h3>{action.title}</h3>
              <p>{action.description}</p>
            </button>
          ))}
        </div>
      </section>

      {/* Main Content */}
      <main className={styles.mainContent}>
        <div className={styles.contentGrid}>
          {/* Navigation Sidebar */}
          <aside className={styles.sidebar}>
            <div className={styles.sidebarContent}>
              <h3 className={styles.sidebarTitle}>
                <BookOpen size={20} />
                {t('tableOfContents', 'Table of Contents')}
              </h3>
              <nav className={styles.navigation}>
                {sections.map(section => (
                  <button
                    key={section.id}
                    onClick={() => setActiveSection(section.id)}
                    className={`${styles.navItem} ${activeSection === section.id ? styles.active : ''}`}
                  >
                    <section.icon size={18} />
                    <span>{section.title}</span>
                    <ChevronRight size={16} />
                  </button>
                ))}
              </nav>

              {/* Best Practices */}
              <div className={styles.bestPractices}>
                <h3 className={styles.bestPracticesTitle}>
                  <Lightbulb size={20} />
                  {t('bestPractices', 'Best Practices')}
                </h3>
                {bestPractices.map((practice, index) => (
                  <div key={index} className={styles.practiceCard}>
                    <practice.icon className={styles.practiceIcon} size={20} />
                    <div>
                      <h4>{practice.title}</h4>
                      <p>{practice.description}</p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </aside>

          {/* Guide Content */}
          <div className={styles.guideContent}>
            {searchQuery && (
              <div className={styles.searchResults}>
                <p>{t('searchResultsText', 'Search results for')}: <strong>{searchQuery}</strong></p>
                <p className={styles.resultCount}>
                  {filteredSections.reduce((acc, section) => acc + section.topics.length, 0)} {t('topicsFound', 'topics found')}
                </p>
              </div>
            )}

            {(searchQuery ? filteredSections : sections.filter(s => s.id === activeSection)).map(section => (
              <section key={section.id} className={styles.section}>
                <div className={styles.sectionHeader}>
                  <section.icon className={styles.sectionIcon} size={32} />
                  <div>
                    <h2>{section.title}</h2>
                    <p className={styles.sectionDescription}>
                      {t(`${section.id}Desc`, `Learn the fundamentals of ${section.title.toLowerCase()}`)}
                    </p>
                  </div>
                </div>

                <div className={styles.topicsList}>
                  {section.topics.map(topic => (
                    <div key={topic.id} className={styles.topicCard}>
                      <button
                        onClick={() => toggleTopic(topic.id)}
                        className={styles.topicHeader}
                      >
                        <div className={styles.topicTitle}>
                          <topic.icon size={20} />
                          <h3>{topic.title}</h3>
                        </div>
                        <ChevronDown 
                          size={20} 
                          className={`${styles.toggleIcon} ${expandedTopics[topic.id] ? styles.expanded : ''}`}
                        />
                      </button>

                      {expandedTopics[topic.id] && (
                        <div className={styles.topicContent}>
                          <ul className={styles.contentList}>
                            {topic.content.map((item, index) => (
                              <li key={index}>
                                <CheckCircle size={16} />
                                <span>{item}</span>
                              </li>
                            ))}
                          </ul>
                          
                          <div className={styles.topicActions}>
                            <button className={styles.learnMoreButton}>
                              {t('learnMore', 'View Detailed Guide')}
                              <ArrowRight size={16} />
                            </button>
                            <button className={styles.videoButton}>
                              <PlayCircle size={16} />
                              {t('watchVideo', 'Watch Tutorial')}
                            </button>
                          </div>
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </section>
            ))}

            {/* Success Tips */}
            <section className={styles.successTips}>
              <div className={styles.tipsHeader}>
                <Award className={styles.tipsIcon} />
                <h2>{t('successTips', 'Pro Tips from Successful Vendors')}</h2>
              </div>
              
              <div className={styles.tipsGrid}>
                <div className={styles.tipCard}>
                  <Zap className={styles.tipIcon} />
                  <h3>{t('tip1Title', 'Automate Everything')}</h3>
                  <p>{t('tip1Desc', 'Use our automation tools to save time on repetitive tasks like inventory updates and order processing.')}</p>
                </div>
                
                <div className={styles.tipCard}>
                  <Users className={styles.tipIcon} />
                  <h3>{t('tip2Title', 'Engage Daily')}</h3>
                  <p>{t('tip2Desc', 'Spend 30 minutes daily engaging with your community through comments, messages, and content.')}</p>
                </div>
                
                <div className={styles.tipCard}>
                  <TrendingUp className={styles.tipIcon} />
                  <h3>{t('tip3Title', 'Test & Optimize')}</h3>
                  <p>{t('tip3Desc', 'Regularly test different pricing, descriptions, and images to find what converts best.')}</p>
                </div>
                
                <div className={styles.tipCard}>
                  <Mail className={styles.tipIcon} />
                  <h3>{t('tip4Title', 'Build Your List')}</h3>
                  <p>{t('tip4Desc', 'Use newsletters to keep customers informed about new products and exclusive offers.')}</p>
                </div>
              </div>
            </section>

            {/* Need Help CTA */}
            <section className={styles.helpCta}>
              <HelpCircle className={styles.helpIcon} />
              <div className={styles.helpContent}>
                <h2>{t('needHelp', 'Still Have Questions?')}</h2>
                <p>{t('helpDesc', 'Our vendor support team is here to help you succeed')}</p>
                <div className={styles.helpActions}>
                  <button 
                    onClick={() => router.push('/support/vendors')}
                    className={styles.contactButton}
                  >
                    <MessageCircle size={20} />
                    {t('contactSupport', 'Contact Support')}
                  </button>
                  <button 
                    onClick={() => router.push('/webinars')}
                    className={styles.webinarButton}
                  >
                    <PlayCircle size={20} />
                    {t('joinWebinar', 'Join Live Webinar')}
                  </button>
                </div>
              </div>
            </section>
          </div>
        </div>
      </main>
    </div>
  );
}