"use client";

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  TrendingUp,
  Users,
  ShoppingBag,
  MessageCircle,
  Calendar,
  Clock,
  ArrowRight,
  Search,
  Filter,
  BookOpen,
  Video,
  Lightbulb,
  Award,
  Target,
  Zap,
  Globe,
  ChevronRight,
  Tag,
  User,
  BarChart3,
  Package,
  Megaphone,
  Sparkles,
  Building2,
  Heart,
  ArrowUp
} from 'lucide-react';
import styles from './VendorBlog.module.css';

export default function VendorBlogPage() {
  const t = useTranslations('vendorBlog');
  const router = useRouter();
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [selectedType, setSelectedType] = useState('all');

  const categories = [
    { id: 'all', label: t('categoryAll', 'All Articles'), count: 42 },
    { id: 'success-stories', label: t('categorySuccess', 'Success Stories'), count: 12 },
    { id: 'tips-tricks', label: t('categoryTips', 'Tips & Tricks'), count: 15 },
    { id: 'platform-updates', label: t('categoryUpdates', 'Platform Updates'), count: 8 },
    { id: 'industry-insights', label: t('categoryInsights', 'Industry Insights'), count: 7 }
  ];

  const contentTypes = [
    { id: 'all', label: t('typeAll', 'All Types'), icon: BookOpen },
    { id: 'article', label: t('typeArticle', 'Articles'), icon: BookOpen },
    { id: 'video', label: t('typeVideo', 'Videos'), icon: Video },
    { id: 'guide', label: t('typeGuide', 'Guides'), icon: Lightbulb }
  ];

  const featuredArticles = [
    {
      id: 1,
      type: 'success-story',
      contentType: 'article',
      title: t('featured1Title', 'How Sarah Built a €1M Business from Her Living Room'),
      excerpt: t('featured1Excerpt', 'Starting with just 5 handmade products, Sarah leveraged our platform\'s community features to build a loyal following of 50,000 customers in just 18 months.'),
      author: 'Platform Team',
      date: '2024-01-15',
      readTime: '8 min read',
      image: '/blog/success-sarah.jpg',
      tags: ['Success Story', 'Small Business', 'Community Building'],
      metrics: {
        revenue: '€1.2M ARR',
        followers: '50K+',
        products: '150+'
      }
    },
    {
      id: 2,
      type: 'platform-update',
      contentType: 'video',
      title: t('featured2Title', 'New AI-Powered Analytics Dashboard'),
      excerpt: t('featured2Excerpt', 'Discover how our new AI analytics can help you understand customer behavior, predict trends, and optimize your product listings for maximum conversion.'),
      author: 'Product Team',
      date: '2024-01-10',
      readTime: '5 min watch',
      image: '/blog/ai-analytics.jpg',
      tags: ['Product Update', 'AI', 'Analytics']
    }
  ];

  const articles = [
    {
      id: 3,
      type: 'tips-tricks',
      contentType: 'guide',
      title: t('article3Title', '10 Ways to Boost Your Product Visibility'),
      excerpt: t('article3Excerpt', 'Learn proven strategies to improve your search rankings and get your products in front of more customers.'),
      author: 'Marketing Team',
      date: '2024-01-08',
      readTime: '6 min read',
      tags: ['SEO', 'Marketing', 'Growth']
    },
    {
      id: 4,
      type: 'industry-insights',
      contentType: 'article',
      title: t('article4Title', 'E-commerce Trends 2024: What Sellers Need to Know'),
      excerpt: t('article4Excerpt', 'Stay ahead of the curve with our comprehensive analysis of emerging e-commerce trends and consumer behaviors.'),
      author: 'Research Team',
      date: '2024-01-05',
      readTime: '10 min read',
      tags: ['Trends', 'Industry Analysis', '2024']
    },
    {
      id: 5,
      type: 'success-stories',
      contentType: 'video',
      title: t('article5Title', 'From Local Shop to Global Brand: The Martinez Journey'),
      excerpt: t('article5Excerpt', 'Watch how a family business expanded internationally using our platform\'s multi-language and currency features.'),
      author: 'Platform Team',
      date: '2024-01-03',
      readTime: '12 min watch',
      tags: ['Success Story', 'International', 'Growth']
    },
    {
      id: 6,
      type: 'tips-tricks',
      contentType: 'guide',
      title: t('article6Title', 'Master Live Streaming: Engage Your Audience in Real-Time'),
      excerpt: t('article6Excerpt', 'Everything you need to know about using live streaming to showcase products and connect with customers.'),
      author: 'Content Team',
      date: '2023-12-28',
      readTime: '7 min read',
      tags: ['Live Streaming', 'Engagement', 'Sales']
    },
    {
      id: 7,
      type: 'platform-updates',
      contentType: 'article',
      title: t('article7Title', 'Introducing Advanced Inventory Management'),
      excerpt: t('article7Excerpt', 'Our new inventory system helps you track stock across multiple channels and automate reordering.'),
      author: 'Product Team',
      date: '2023-12-20',
      readTime: '4 min read',
      tags: ['Product Update', 'Inventory', 'Automation']
    },
    {
      id: 8,
      type: 'tips-tricks',
      contentType: 'guide',
      title: t('article8Title', 'Building Customer Loyalty Through Community'),
      excerpt: t('article8Excerpt', 'Discover how to use followers, newsletters, and exclusive content to create a loyal customer base.'),
      author: 'Community Team',
      date: '2023-12-15',
      readTime: '8 min read',
      tags: ['Community', 'Loyalty', 'Retention']
    }
  ];

  const popularTags = [
    { name: 'Growth Strategies', count: 24 },
    { name: 'Marketing Tips', count: 18 },
    { name: 'Success Stories', count: 15 },
    { name: 'Product Updates', count: 12 },
    { name: 'Community Building', count: 10 },
    { name: 'Analytics', count: 8 }
  ];

  const contributorSpotlight = {
    name: 'Emma Rodriguez',
    role: t('spotlightRole', 'Power Seller & Community Leader'),
    bio: t('spotlightBio', 'Emma has grown her sustainable fashion brand to €2M in annual revenue and shares her journey and tips with our community.'),
    articles: 12,
    followers: '25K'
  };

  const filteredArticles = [...featuredArticles, ...articles].filter(article => {
    const matchesSearch = searchQuery === '' || 
      article.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
      article.excerpt.toLowerCase().includes(searchQuery.toLowerCase());
    
    const matchesCategory = selectedCategory === 'all' || article.type === selectedCategory;
    const matchesType = selectedType === 'all' || article.contentType === selectedType;
    
    return matchesSearch && matchesCategory && matchesType;
  });

  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.headerInfo}>
            <h1 className={styles.title}>
              {t('pageTitle', 'Vendor Success Hub')}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', 'Stories, insights, and resources to help you grow your business')}
            </p>
            
            {/* Search Bar */}
            <div className={styles.searchBar}>
              <Search className={styles.searchIcon} size={20} />
              <input
                type="text"
                placeholder={t('searchPlaceholder', 'Search articles, guides, and videos...')}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className={styles.searchInput}
              />
            </div>
          </div>
        </div>
      </header>

      {/* Quick Stats */}
      <section className={styles.quickStats}>
        <div className={styles.statsContainer}>
          <div className={styles.statCard}>
            <BookOpen className={styles.statIcon} />
            <div className={styles.statValue}>42</div>
            <div className={styles.statLabel}>{t('statArticles', 'Articles & Guides')}</div>
          </div>
          <div className={styles.statCard}>
            <Video className={styles.statIcon} />
            <div className={styles.statValue}>18</div>
            <div className={styles.statLabel}>{t('statVideos', 'Video Tutorials')}</div>
          </div>
          <div className={styles.statCard}>
            <Users className={styles.statIcon} />
            <div className={styles.statValue}>156</div>
            <div className={styles.statLabel}>{t('statContributors', 'Contributors')}</div>
          </div>
          <div className={styles.statCard}>
            <TrendingUp className={styles.statIcon} />
            <div className={styles.statValue}>€2.4B</div>
            <div className={styles.statLabel}>{t('statSuccess', 'Seller Success GMV')}</div>
          </div>
        </div>
      </section>

      {/* Main Content */}
      <main className={styles.mainContent}>
        <div className={styles.contentGrid}>
          {/* Left Sidebar - Filters */}
          <aside className={styles.sidebar}>
            {/* Categories */}
            <div className={styles.filterSection}>
              <h3 className={styles.filterTitle}>
                <Filter size={18} />
                {t('categories', 'Categories')}
              </h3>
              <div className={styles.categoryList}>
                {categories.map(category => (
                  <button
                    key={category.id}
                    onClick={() => setSelectedCategory(category.id)}
                    className={`${styles.categoryButton} ${selectedCategory === category.id ? styles.active : ''}`}
                  >
                    <span>{category.label}</span>
                    <span className={styles.count}>{category.count}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Content Types */}
            <div className={styles.filterSection}>
              <h3 className={styles.filterTitle}>
                <Sparkles size={18} />
                {t('contentType', 'Content Type')}
              </h3>
              <div className={styles.typeList}>
                {contentTypes.map(type => (
                  <button
                    key={type.id}
                    onClick={() => setSelectedType(type.id)}
                    className={`${styles.typeButton} ${selectedType === type.id ? styles.active : ''}`}
                  >
                    <type.icon size={16} />
                    <span>{type.label}</span>
                  </button>
                ))}
              </div>
            </div>

            {/* Popular Tags */}
            <div className={styles.filterSection}>
              <h3 className={styles.filterTitle}>
                <Tag size={18} />
                {t('popularTags', 'Popular Tags')}
              </h3>
              <div className={styles.tagCloud}>
                {popularTags.map((tag, index) => (
                  <button key={index} className={styles.tagButton}>
                    {tag.name} ({tag.count})
                  </button>
                ))}
              </div>
            </div>

            {/* Newsletter CTA */}
            <div className={styles.newsletterCard}>
              <Megaphone className={styles.newsletterIcon} />
              <h3>{t('newsletterTitle', 'Weekly Vendor Tips')}</h3>
              <p>{t('newsletterDesc', 'Get the latest success stories and growth strategies delivered to your inbox')}</p>
              <button className={styles.subscribeButton}>
                {t('subscribe', 'Subscribe Now')}
              </button>
            </div>
          </aside>

          {/* Main Article Grid */}
          <div className={styles.articlesSection}>
            {/* Featured Articles */}
            {selectedCategory === 'all' && selectedType === 'all' && searchQuery === '' && (
              <div className={styles.featuredSection}>
                <h2 className={styles.sectionTitle}>
                  <Award size={24} />
                  {t('featured', 'Featured This Week')}
                </h2>
                <div className={styles.featuredGrid}>
                  {featuredArticles.map(article => (
                    <article key={article.id} className={styles.featuredCard}>
                      <div className={styles.featuredImage}>
                        <div className={styles.contentTypeTag}>
                          {article.contentType === 'video' ? <Video size={16} /> : <BookOpen size={16} />}
                          {article.contentType}
                        </div>
                      </div>
                      <div className={styles.featuredContent}>
                        <div className={styles.articleMeta}>
                          <span className={styles.author}>
                            <User size={14} />
                            {article.author}
                          </span>
                          <span className={styles.date}>
                            <Calendar size={14} />
                            {new Date(article.date).toLocaleDateString()}
                          </span>
                          <span className={styles.readTime}>
                            <Clock size={14} />
                            {article.readTime}
                          </span>
                        </div>
                        <h3 className={styles.featuredTitle}>{article.title}</h3>
                        <p className={styles.featuredExcerpt}>{article.excerpt}</p>
                        
                        {article.metrics && (
                          <div className={styles.successMetrics}>
                            <div className={styles.metric}>
                              <BarChart3 size={16} />
                              {article.metrics.revenue}
                            </div>
                            <div className={styles.metric}>
                              <Users size={16} />
                              {article.metrics.followers}
                            </div>
                            <div className={styles.metric}>
                              <Package size={16} />
                              {article.metrics.products}
                            </div>
                          </div>
                        )}
                        
                        <div className={styles.articleTags}>
                          {article.tags.map((tag, index) => (
                            <span key={index} className={styles.tag}>{tag}</span>
                          ))}
                        </div>
                        
                        <button className={styles.readMoreButton}>
                          {t('readMore', 'Read Full Story')}
                          <ArrowRight size={16} />
                        </button>
                      </div>
                    </article>
                  ))}
                </div>
              </div>
            )}

            {/* All Articles */}
            <div className={styles.allArticles}>
              <h2 className={styles.sectionTitle}>
                <BookOpen size={24} />
                {searchQuery || selectedCategory !== 'all' || selectedType !== 'all' 
                  ? t('searchResults', 'Search Results') 
                  : t('allArticles', 'All Articles')}
                <span className={styles.resultCount}>({filteredArticles.length})</span>
              </h2>
              
              <div className={styles.articleGrid}>
                {filteredArticles.map(article => (
                  <article key={article.id} className={styles.articleCard}>
                    <div className={styles.articleHeader}>
                      <div className={styles.contentTypeIcon}>
                        {article.contentType === 'video' ? <Video size={20} /> : 
                         article.contentType === 'guide' ? <Lightbulb size={20} /> : 
                         <BookOpen size={20} />}
                      </div>
                      <div className={styles.articleMeta}>
                        <span className={styles.author}>{article.author}</span>
                        <span className={styles.date}>{new Date(article.date).toLocaleDateString()}</span>
                      </div>
                    </div>
                    
                    <h3 className={styles.articleTitle}>{article.title}</h3>
                    <p className={styles.articleExcerpt}>{article.excerpt}</p>
                    
                    <div className={styles.articleFooter}>
                      <div className={styles.articleTags}>
                        {article.tags.slice(0, 2).map((tag, index) => (
                          <span key={index} className={styles.tag}>{tag}</span>
                        ))}
                      </div>
                      <span className={styles.readTime}>
                        <Clock size={14} />
                        {article.readTime}
                      </span>
                    </div>
                    
                    <button className={styles.articleLink}>
                      {t('read', 'Read')}
                      <ChevronRight size={16} />
                    </button>
                  </article>
                ))}
              </div>
              
              {filteredArticles.length === 0 && (
                <div className={styles.noResults}>
                  <Search size={48} />
                  <h3>{t('noResults', 'No articles found')}</h3>
                  <p>{t('tryDifferent', 'Try adjusting your filters or search terms')}</p>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Contributor Spotlight */}
        <section className={styles.spotlightSection}>
          <div className={styles.spotlightContent}>
            <div className={styles.spotlightHeader}>
              <h2>
                <Users size={28} />
                {t('contributorSpotlight', 'Contributor Spotlight')}
              </h2>
            </div>
            <div className={styles.spotlightCard}>
              <div className={styles.spotlightInfo}>
                <div className={styles.spotlightAvatar}>
                  <User size={40} />
                </div>
                <div>
                  <h3>{contributorSpotlight.name}</h3>
                  <p className={styles.spotlightRole}>{contributorSpotlight.role}</p>
                  <p className={styles.spotlightBio}>{contributorSpotlight.bio}</p>
                  <div className={styles.spotlightStats}>
                    <span>
                      <BookOpen size={16} />
                      {contributorSpotlight.articles} {t('articles', 'Articles')}
                    </span>
                    <span>
                      <Heart size={16} />
                      {contributorSpotlight.followers} {t('followers', 'Followers')}
                    </span>
                  </div>
                </div>
              </div>
              <button className={styles.viewProfileButton}>
                {t('viewProfile', 'View Profile')}
                <ArrowRight size={16} />
              </button>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className={styles.ctaSection}>
          <div className={styles.ctaContent}>
            <Building2 className={styles.ctaIcon} />
            <h2>{t('ctaTitle', 'Share Your Success Story')}</h2>
            <p>{t('ctaDesc', 'Have a great story about growing your business on our platform? We\'d love to feature you!')}</p>
            <div className={styles.ctaButtons}>
              <button 
                onClick={() => router.push('/contact/editorial')}
                className={styles.submitStoryButton}
              >
                {t('submitStory', 'Submit Your Story')}
                <ArrowUp size={20} />
              </button>
              <button 
                onClick={() => router.push('/docs/vendor-guide')}
                className={styles.guideButton}
              >
                {t('vendorGuide', 'Read Vendor Guide')}
                <ChevronRight size={20} />
              </button>
            </div>
          </div>
        </section>
      </main>
    </div>
  );
}