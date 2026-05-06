"use client";

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  ChevronDown,
  Search,
  HelpCircle,
  DollarSign,
  Package,
  Truck,
  Shield,
  BarChart3,
  Users,
  Settings,
  Bot,
  Mail,
  Clock,
  Globe,
  FileText,
  CreditCard,
  AlertCircle,
  ArrowLeft,
  Phone,
  MessageSquare,
  Zap,
  Store,
  TrendingUp,
  Heart,
  Sparkles,
  Megaphone,
  Building2,
  CheckCircle,
  BookOpen,
  HeadphonesIcon,
  Target,
  Award,
  Lightbulb
} from 'lucide-react';
import styles from './VendorsFAQ.module.css';

export default function VendorsFAQPage() {
  const t = useTranslations('vendorsFAQ');
  const router = useRouter();
  const [searchTerm, setSearchTerm] = useState('');
  const [activeCategory, setActiveCategory] = useState('all');
  const [expandedItems, setExpandedItems] = useState({});

  const toggleExpand = (id) => {
    setExpandedItems(prev => ({
      ...prev,
      [id]: !prev[id]
    }));
  };

  const faqCategories = [
    { id: 'all', name: t('categoryAll', 'All Questions'), icon: HelpCircle, count: 45 },
    { id: 'getting-started', name: t('categoryGettingStarted', 'Getting Started'), icon: Sparkles, count: 8 },
    { id: 'selling', name: t('categorySelling', 'Selling & Growth'), icon: TrendingUp, count: 10 },
    { id: 'payments', name: t('categoryPayments', 'Payments & Fees'), icon: DollarSign, count: 7 },
    { id: 'shipping', name: t('categoryShipping', 'Shipping & Fulfillment'), icon: Truck, count: 6 },
    { id: 'marketing', name: t('categoryMarketing', 'Marketing & Community'), icon: Megaphone, count: 8 },
    { id: 'technical', name: t('categoryTechnical', 'Technical & API'), icon: Globe, count: 6 }
  ];

  const popularQuestions = [
    { id: 'pop1', icon: Zap, text: t('popular1', 'How much does it cost to start selling?') },
    { id: 'pop2', icon: Clock, text: t('popular2', 'How quickly can I start selling?') },
    { id: 'pop3', icon: Users, text: t('popular3', 'Do I need a business to join?') },
    { id: 'pop4', icon: CreditCard, text: t('popular4', 'When do I get paid?') }
  ];

  const faqItems = [
    // Getting Started
    {
      id: 'gs1',
      category: 'getting-started',
      question: t('gs1Question', 'Can anyone join the platform?'),
      answer: t('gs1Answer', 'Yes! Our platform welcomes all types of businesses - from solo entrepreneurs and small businesses to large enterprises. Whether you sell physical products, digital goods, or services, you can start with just one product or thousands. We have plans designed for every business size.'),
      helpful: 156
    },
    {
      id: 'gs2',
      category: 'getting-started',
      question: t('gs2Question', 'How long does the verification process take?'),
      answer: t('gs2Answer', 'Most applications are approved within 24-48 hours. During peak periods, it may take up to 72 hours. We\'ll notify you via email at each step of the process. Business verification ensures a trusted marketplace for all users.'),
      helpful: 142
    },
    {
      id: 'gs3',
      category: 'getting-started',
      question: t('gs3Question', 'What documents do I need to register?'),
      answer: t('gs3Answer', 'For businesses: Business registration document, Tax ID, Bank account verification. For individuals (in supported countries): Government ID, Tax information, Bank account. We also accept sole proprietors and freelancers.'),
      helpful: 128
    },
    {
      id: 'gs4',
      category: 'getting-started',
      question: t('gs4Question', 'Is there a free plan available?'),
      answer: t('gs4Answer', 'Yes! Our Starter plan is completely FREE with no monthly fees. You only pay a transaction fee of 8% when you make a sale. This includes hosting, SSL, payment processing, and all basic features. Perfect for testing the platform or starting small.'),
      helpful: 189
    },
    {
      id: 'gs5',
      category: 'getting-started',
      question: t('gs5Question', 'How do I set up my first store?'),
      answer: t('gs5Answer', 'After approval, you\'ll get access to our intuitive setup wizard that guides you through: 1) Customizing your store design, 2) Adding your first products, 3) Setting up payment methods, 4) Configuring shipping options, 5) Activating marketing tools. Most sellers are live within 2-3 hours!'),
      helpful: 134
    },

    // Selling & Growth
    {
      id: 'sell1',
      category: 'selling',
      question: t('sell1Question', 'How many products can I list?'),
      answer: t('sell1Answer', 'Product limits by plan: Starter - up to 100 products, Professional - up to 1,000 products, Enterprise - unlimited products. Each product can have unlimited variations (sizes, colors) and up to 10 high-quality images or videos.'),
      helpful: 112
    },
    {
      id: 'sell2',
      category: 'selling',
      question: t('sell2Question', 'Can I sell both products and services?'),
      answer: t('sell2Answer', 'Absolutely! Our platform supports physical products, digital downloads, services, subscriptions, bookings, and consultations. You can mix different types in your store. Many sellers offer products with installation services or digital products with consulting.'),
      helpful: 98
    },
    {
      id: 'sell3',
      category: 'selling',
      question: t('sell3Question', 'How do I grow my customer base?'),
      answer: t('sell3Answer', 'We provide powerful tools to grow: 1) Built-in SEO optimization, 2) Social media integration, 3) Email marketing (newsletters), 4) Community features (followers, reviews), 5) Content marketing (blogs, videos), 6) Live streaming for product demos, 7) Loyalty programs and referrals.'),
      helpful: 167
    },
    {
      id: 'sell4',
      category: 'selling',
      question: t('sell4Question', 'Can I import products from other platforms?'),
      answer: t('sell4Answer', 'Yes! We offer seamless import from: Amazon, eBay, Shopify, WooCommerce, Etsy, and more. You can import via CSV/XML files or use our API. Our migration team provides free assistance for large catalogs (100+ products).'),
      helpful: 89
    },

    // Payments & Fees
    {
      id: 'pay1',
      category: 'payments',
      question: t('pay1Question', 'What are the transaction fees?'),
      answer: t('pay1Answer', 'Transaction fees by plan: Starter - 8% (FREE plan, no monthly fee), Professional - 5% + €49/month, Enterprise - from 3% + custom pricing. Fees include payment processing, hosting, SSL, and platform features. No hidden charges!'),
      helpful: 203
    },
    {
      id: 'pay2',
      category: 'payments',
      question: t('pay2Question', 'When do I receive payments?'),
      answer: t('pay2Answer', 'Standard payout cycle: 7 business days after delivery confirmation. Express payouts available: 2-3 days for sellers with 6+ months history. Instant payouts: Available for Enterprise plans. All payouts are automatic to your registered bank account.'),
      helpful: 178
    },
    {
      id: 'pay3',
      category: 'payments',
      question: t('pay3Question', 'What payment methods are supported?'),
      answer: t('pay3Answer', 'We accept all major payment methods: Credit/debit cards (Visa, Mastercard, Amex), Digital wallets (PayPal, Apple Pay, Google Pay), Bank transfers, Buy now pay later (Klarna, Afterpay), Cryptocurrencies (Bitcoin, Ethereum), Local payment methods (100+ options).'),
      helpful: 145
    },
    {
      id: 'pay4',
      category: 'payments',
      question: t('pay4Question', 'Are there any setup or hidden fees?'),
      answer: t('pay4Answer', 'No hidden fees! Free: Product listings, Store customization, SSL certificate, Basic analytics, Email support. Optional paid services: Professional photography, Premium themes, Advanced marketing tools, Priority support. Everything is transparent in your dashboard.'),
      helpful: 132
    },

    // Shipping & Fulfillment
    {
      id: 'ship1',
      category: 'shipping',
      question: t('ship1Question', 'Do I have to handle shipping myself?'),
      answer: t('ship1Answer', 'You have three options: 1) Self-fulfillment - You pack and ship orders, 2) Our Fulfillment Service - We store, pack, and ship for you (Professional plan+), 3) Dropshipping - Ship directly from suppliers. Choose what works best for your business model.'),
      helpful: 156
    },
    {
      id: 'ship2',
      category: 'shipping',
      question: t('ship2Question', 'Which shipping carriers are integrated?'),
      answer: t('ship2Answer', 'Major carriers worldwide: DHL, UPS, FedEx, USPS, Royal Mail, DPD, and 50+ regional carriers. Features include: Discounted shipping rates, Automatic label printing, Real-time tracking, International shipping, Same-day delivery options in major cities.'),
      helpful: 123
    },
    {
      id: 'ship3',
      category: 'shipping',
      question: t('ship3Question', 'Can I sell internationally?'),
      answer: t('ship3Answer', 'Yes! Sell globally with: Automatic currency conversion (170+ currencies), Multi-language stores (30+ languages), International shipping calculation, Customs documentation, Tax compliance tools, Local payment methods. We make global selling simple!'),
      helpful: 167
    },

    // Marketing & Community
    {
      id: 'mkt1',
      category: 'marketing',
      question: t('mkt1Question', 'How do the community features work?'),
      answer: t('mkt1Answer', 'Build your community with: Followers system - customers follow your store, Live chat - real-time customer conversations, Content publishing - blogs, videos, tutorials, Reviews & ratings - social proof, Comments - engage on products, Live streaming - product demos and Q&A sessions.'),
      helpful: 145
    },
    {
      id: 'mkt2',
      category: 'marketing',
      question: t('mkt2Question', 'Can I send newsletters to customers?'),
      answer: t('mkt2Answer', 'Yes! Email marketing features: Drag-and-drop email builder, Automated campaigns (welcome, abandoned cart), Segmentation tools, A/B testing, Analytics and reporting. Limits: Starter - 1,000 contacts, Professional - 5,000, Enterprise - unlimited.'),
      helpful: 112
    },
    {
      id: 'mkt3',
      category: 'marketing',
      question: t('mkt3Question', 'How does the AI assistant help with sales?'),
      answer: t('mkt3Answer', 'Your AI assistant works 24/7 to: Answer customer questions instantly, Recommend products based on needs, Handle basic support queries, Capture leads when you\'re offline, Speak 30+ languages automatically, Learn from your products and policies. It\'s like having a smart salesperson always available!'),
      helpful: 189
    },

    // Technical & API
    {
      id: 'tech1',
      category: 'technical',
      question: t('tech1Question', 'Can I customize my store design?'),
      answer: t('tech1Answer', 'Absolutely! Customization options: 50+ professional themes, Drag-and-drop builder, Custom CSS/HTML (Professional+), Mobile-responsive designs, Brand colors and fonts, Custom domains. No coding required, but developers can access advanced customization.'),
      helpful: 134
    },
    {
      id: 'tech2',
      category: 'technical',
      question: t('tech2Question', 'Is there an API for integrations?'),
      answer: t('tech2Answer', 'Yes! Our RESTful API allows: Product management, Order processing, Inventory sync, Customer data access, Custom integrations. We also offer webhooks, SDKs for popular languages, and Zapier integration for 3000+ apps. Documentation at docs.platform.com/api.'),
      helpful: 98
    },
    {
      id: 'tech3',
      category: 'technical',
      question: t('tech3Question', 'What about SEO and marketing tools?'),
      answer: t('tech3Answer', 'Built-in SEO features: Automatic sitemap generation, Meta tag optimization, Schema markup, Fast page load speeds, Mobile optimization, Blog platform for content marketing. Plus integrations with Google Analytics, Facebook Pixel, and marketing platforms.'),
      helpful: 156
    }
  ];

  const filteredItems = faqItems.filter(item => {
    const matchesCategory = activeCategory === 'all' || item.category === activeCategory;
    const matchesSearch = item.question.toLowerCase().includes(searchTerm.toLowerCase()) ||
                         item.answer.toLowerCase().includes(searchTerm.toLowerCase());
    return matchesCategory && matchesSearch;
  });

  const getCategoryIcon = (categoryId) => {
    const category = faqCategories.find(cat => cat.id === categoryId);
    return category ? category.icon : HelpCircle;
  };

  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <button 
            onClick={() => router.back()}
            className={styles.backButton}
          >
            <ArrowLeft size={20} />
            <span>{t('back', 'Back')}</span>
          </button>
          <div className={styles.headerInfo}>
            <h1 className={styles.title}>
              {t('pageTitle', 'Vendor Help Center')}
            </h1>
            <p className={styles.subtitle}>
              {t('pageSubtitle', 'Everything you need to know about selling on our platform')}
            </p>
          </div>
        </div>
      </header>

      {/* Search Section */}
      <section className={styles.searchSection}>
        <div className={styles.searchContainer}>
          <div className={styles.searchWrapper}>
            <Search className={styles.searchIcon} size={24} />
            <input
              type="text"
              value={searchTerm}
              onChange={(e) => setSearchTerm(e.target.value)}
              placeholder={t('searchPlaceholder', 'Search for answers...')}
              className={styles.searchInput}
            />
            {searchTerm && (
              <button
                onClick={() => setSearchTerm('')}
                className={styles.clearButton}
              >
                ×
              </button>
            )}
          </div>
          
          {/* Popular Questions */}
          <div className={styles.popularQuestions}>
            <h3>{t('popularQuestions', 'Popular questions')}</h3>
            <div className={styles.popularGrid}>
              {popularQuestions.map(q => (
                <button
                  key={q.id}
                  onClick={() => setSearchTerm(q.text)}
                  className={styles.popularButton}
                >
                  <q.icon size={16} />
                  <span>{q.text}</span>
                </button>
              ))}
            </div>
          </div>
        </div>
      </section>

      {/* Main Content */}
      <main className={styles.mainContent}>
        <div className={styles.contentGrid}>
          {/* Categories Sidebar */}
          <aside className={styles.sidebar}>
            <h3 className={styles.sidebarTitle}>
              <BookOpen size={20} />
              {t('categories', 'Categories')}
            </h3>
            <nav className={styles.categoryNav}>
              {faqCategories.map(category => (
                <button
                  key={category.id}
                  onClick={() => setActiveCategory(category.id)}
                  className={`${styles.categoryButton} ${
                    activeCategory === category.id ? styles.active : ''
                  }`}
                >
                  <category.icon size={18} />
                  <span>{category.name}</span>
                  <span className={styles.count}>{category.count}</span>
                </button>
              ))}
            </nav>

            {/* Need More Help */}
            <div className={styles.needHelp}>
              <HeadphonesIcon className={styles.helpIcon} />
              <h3>{t('needHelp', 'Need more help?')}</h3>
              <p>{t('helpDesc', 'Our support team is here for you')}</p>
              <button 
                onClick={() => router.push('/contact/support')}
                className={styles.contactButton}
              >
                {t('contactSupport', 'Contact Support')}
              </button>
            </div>
          </aside>

          {/* FAQ Content */}
          <div className={styles.faqContent}>
            {searchTerm && (
              <div className={styles.searchResults}>
                <p>
                  {t('showingResults', 'Showing results for')}: <strong>{searchTerm}</strong>
                </p>
                <p className={styles.resultCount}>
                  {filteredItems.length} {t('questionsFound', 'questions found')}
                </p>
              </div>
            )}

            {filteredItems.length === 0 ? (
              <div className={styles.noResults}>
                <AlertCircle size={48} />
                <h3>{t('noResults', 'No results found')}</h3>
                <p>{t('tryDifferentSearch', 'Try different keywords or browse categories')}</p>
              </div>
            ) : (
              <div className={styles.faqList}>
                {filteredItems.map(item => {
                  const CategoryIcon = getCategoryIcon(item.category);
                  return (
                    <div 
                      key={item.id} 
                      className={`${styles.faqItem} ${
                        expandedItems[item.id] ? styles.expanded : ''
                      }`}
                    >
                      <button
                        onClick={() => toggleExpand(item.id)}
                        className={styles.faqQuestion}
                      >
                        <div className={styles.questionContent}>
                          <CategoryIcon size={20} className={styles.categoryIcon} />
                          <span>{item.question}</span>
                        </div>
                        <ChevronDown 
                          size={20} 
                          className={`${styles.toggleIcon} ${
                            expandedItems[item.id] ? styles.expanded : ''
                          }`}
                        />
                      </button>
                      {expandedItems[item.id] && (
                        <div className={styles.faqAnswer}>
                          <p>{item.answer}</p>
                          <div className={styles.answerFooter}>
                            <span className={styles.helpful}>
                              <Lightbulb size={16} />
                              {item.helpful} {t('foundHelpful', 'found this helpful')}
                            </span>
                            <div className={styles.answerActions}>
                              <button className={styles.helpfulButton}>
                                <CheckCircle size={16} />
                                {t('yes', 'Yes')}
                              </button>
                              <button className={styles.helpfulButton}>
                                <AlertCircle size={16} />
                                {t('no', 'No')}
                              </button>
                            </div>
                          </div>
                        </div>
                      )}
                    </div>
                  );
                })}
              </div>
            )}
          </div>
        </div>
      </main>

      {/* Resources Section */}
      <section className={styles.resourcesSection}>
        <div className={styles.resourcesContent}>
          <h2>{t('moreResources', 'More Resources for Success')}</h2>
          <div className={styles.resourcesGrid}>
            <a href="/docs/vendor-guide" className={styles.resourceCard}>
              <div className={styles.resourceIcon}>
                <FileText size={32} />
              </div>
              <h3>{t('vendorGuide', 'Vendor Guide')}</h3>
              <p>{t('vendorGuideDesc', 'Complete guide to selling successfully')}</p>
              <span className={styles.resourceLink}>
                {t('readGuide', 'Read Guide')}
                <ArrowLeft size={16} style={{ transform: 'rotate(180deg)' }} />
              </span>
            </a>
            
            <a href="/webinars" className={styles.resourceCard}>
              <div className={styles.resourceIcon}>
                <Award size={32} />
              </div>
              <h3>{t('webinars', 'Live Webinars')}</h3>
              <p>{t('webinarsDesc', 'Learn from experts and successful sellers')}</p>
              <span className={styles.resourceLink}>
                {t('viewWebinars', 'View Schedule')}
                <ArrowLeft size={16} style={{ transform: 'rotate(180deg)' }} />
              </span>
            </a>
            
            <a href="/blog/vendors" className={styles.resourceCard}>
              <div className={styles.resourceIcon}>
                <TrendingUp size={32} />
              </div>
              <h3>{t('successStories', 'Success Stories')}</h3>
              <p>{t('successStoriesDesc', 'Get inspired by other vendors')}</p>
              <span className={styles.resourceLink}>
                {t('readStories', 'Read Stories')}
                <ArrowLeft size={16} style={{ transform: 'rotate(180deg)' }} />
              </span>
            </a>
            
            <a href="/docs/api" className={styles.resourceCard}>
              <div className={styles.resourceIcon}>
                <Globe size={32} />
              </div>
              <h3>{t('apiDocs', 'API Documentation')}</h3>
              <p>{t('apiDocsDesc', 'Integrate and automate your business')}</p>
              <span className={styles.resourceLink}>
                {t('viewDocs', 'View Docs')}
                <ArrowLeft size={16} style={{ transform: 'rotate(180deg)' }} />
              </span>
            </a>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className={styles.ctaSection}>
        <div className={styles.ctaContent}>
          <Building2 className={styles.ctaIcon} />
          <h2>{t('readyToStart', 'Ready to Start Selling?')}</h2>
          <p>{t('ctaDesc', 'Join thousands of successful businesses on our platform')}</p>
          <div className={styles.ctaButtons}>
            <button 
              onClick={() => router.push('/sell')}
              className={styles.startButton}
            >
              <Store size={20} />
              {t('startSelling', 'Start Selling Now')}
            </button>
            <button 
              onClick={() => router.push('/contact/sales')}
              className={styles.talkButton}
            >
              <Phone size={20} />
              {t('talkToSales', 'Talk to Sales')}
            </button>
          </div>
        </div>
      </section>
    </div>
  );
}