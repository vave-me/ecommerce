"use client";

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  Video,
  Calendar,
  Clock,
  Users,
  PlayCircle,
  CheckCircle,
  Award,
  Globe,
  Zap,
  TrendingUp,
  MessageCircle,
  Download,
  Star,
  ChevronRight,
  Filter,
  User,
  Mic,
  Monitor,
  BarChart3,
  ShoppingBag,
  Heart,
  Bell,
  ArrowRight,
  Info,
  Sparkles,
  Building2,
  Target,
  Package,
  Code,
  Megaphone
} from 'lucide-react';
import styles from './Webinars.module.css';

export default function WebinarsPage() {
  const t = useTranslations('webinars');
  const router = useRouter();
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [showRegistrationForm, setShowRegistrationForm] = useState(false);
  const [selectedWebinar, setSelectedWebinar] = useState(null);

  const categories = [
    { id: 'all', label: t('categoryAll', 'All Webinars'), icon: Globe },
    { id: 'getting-started', label: t('categoryGettingStarted', 'Getting Started'), icon: Rocket },
    { id: 'growth', label: t('categoryGrowth', 'Growth & Scale'), icon: TrendingUp },
    { id: 'marketing', label: t('categoryMarketing', 'Marketing & Sales'), icon: Megaphone },
    { id: 'technical', label: t('categoryTechnical', 'Technical & API'), icon: Code }
  ];

  const upcomingWebinars = [
    {
      id: 1,
      title: t('webinar1Title', 'Platform Basics: Your First 30 Days'),
      description: t('webinar1Desc', 'Everything you need to know to launch successfully on our platform. Perfect for new vendors.'),
      date: '2024-02-15',
      time: '14:00 UTC',
      duration: '60 min',
      category: 'getting-started',
      level: 'Beginner',
      speaker: {
        name: 'Sarah Chen',
        role: t('speaker1Role', 'Head of Vendor Success'),
        image: '/speakers/sarah.jpg'
      },
      topics: [
        t('topic1a', 'Account setup and verification'),
        t('topic1b', 'Product listing best practices'),
        t('topic1c', 'Understanding platform features'),
        t('topic1d', 'Q&A session')
      ],
      attendees: 247,
      maxAttendees: 500
    },
    {
      id: 2,
      title: t('webinar2Title', 'Scaling to €1M: Growth Strategies That Work'),
      description: t('webinar2Desc', 'Learn from successful vendors who have scaled their businesses to 7-figures on our platform.'),
      date: '2024-02-20',
      time: '16:00 UTC',
      duration: '90 min',
      category: 'growth',
      level: 'Advanced',
      speaker: {
        name: 'Michael Torres',
        role: t('speaker2Role', 'Power Seller (€2M+ ARR)'),
        image: '/speakers/michael.jpg'
      },
      topics: [
        t('topic2a', 'Automation and efficiency'),
        t('topic2b', 'Multi-channel expansion'),
        t('topic2c', 'Team building strategies'),
        t('topic2d', 'Live case studies')
      ],
      attendees: 412,
      maxAttendees: 500
    },
    {
      id: 3,
      title: t('webinar3Title', 'Community Building: Turn Followers into Customers'),
      description: t('webinar3Desc', 'Master the art of community engagement and convert your followers into loyal customers.'),
      date: '2024-02-25',
      time: '15:00 UTC',
      duration: '75 min',
      category: 'marketing',
      level: 'Intermediate',
      speaker: {
        name: 'Emma Rodriguez',
        role: t('speaker3Role', 'Community Strategy Expert'),
        image: '/speakers/emma.jpg'
      },
      topics: [
        t('topic3a', 'Content strategy for engagement'),
        t('topic3b', 'Newsletter marketing tactics'),
        t('topic3c', 'Social proof and reviews'),
        t('topic3d', 'Building brand advocates')
      ],
      attendees: 189,
      maxAttendees: 300
    }
  ];

  const pastWebinars = [
    {
      id: 4,
      title: t('past1Title', 'Holiday Sales Mastery'),
      description: t('past1Desc', 'Strategies for maximizing sales during peak shopping seasons'),
      date: '2023-12-10',
      duration: '90 min',
      category: 'marketing',
      views: '3.2K',
      rating: 4.8,
      recordingAvailable: true
    },
    {
      id: 5,
      title: t('past2Title', 'International Expansion Guide'),
      description: t('past2Desc', 'How to successfully sell in multiple countries and currencies'),
      date: '2023-12-05',
      duration: '75 min',
      category: 'growth',
      views: '2.1K',
      rating: 4.9,
      recordingAvailable: true
    },
    {
      id: 6,
      title: t('past3Title', 'AI Tools for Sellers'),
      description: t('past3Desc', 'Leverage AI for product descriptions, pricing, and customer service'),
      date: '2023-11-28',
      duration: '60 min',
      category: 'technical',
      views: '4.5K',
      rating: 4.7,
      recordingAvailable: true
    },
    {
      id: 7,
      title: t('past4Title', 'Product Photography Workshop'),
      description: t('past4Desc', 'Professional product photos on any budget'),
      date: '2023-11-20',
      duration: '120 min',
      category: 'getting-started',
      views: '5.8K',
      rating: 4.9,
      recordingAvailable: true
    }
  ];

  const webinarSeries = [
    {
      id: 'seller-success',
      title: t('series1Title', 'Seller Success Series'),
      description: t('series1Desc', 'Monthly deep-dives into advanced selling strategies'),
      schedule: t('series1Schedule', 'Every 3rd Tuesday'),
      upcoming: 'Feb 20, 2024',
      topics: ['Growth Hacking', 'Automation', 'Analytics', 'Team Building']
    },
    {
      id: 'new-vendor',
      title: t('series2Title', 'New Vendor Bootcamp'),
      description: t('series2Desc', 'Weekly sessions for vendors in their first 90 days'),
      schedule: t('series2Schedule', 'Every Thursday'),
      upcoming: 'Feb 8, 2024',
      topics: ['Setup', 'First Sale', 'Marketing Basics', 'Customer Service']
    }
  ];

  const benefits = [
    {
      icon: Award,
      title: t('benefit1', 'Expert Insights'),
      description: t('benefit1Desc', 'Learn from successful sellers and platform experts')
    },
    {
      icon: Users,
      title: t('benefit2', 'Live Q&A'),
      description: t('benefit2Desc', 'Get your questions answered in real-time')
    },
    {
      icon: Download,
      title: t('benefit3', 'Resources & Templates'),
      description: t('benefit3Desc', 'Download exclusive guides and worksheets')
    },
    {
      icon: Globe,
      title: t('benefit4', 'Global Community'),
      description: t('benefit4Desc', 'Connect with vendors from around the world')
    }
  ];

  const handleRegister = (webinar) => {
    setSelectedWebinar(webinar);
    setShowRegistrationForm(true);
  };

  const filteredUpcoming = selectedCategory === 'all' 
    ? upcomingWebinars 
    : upcomingWebinars.filter(w => w.category === selectedCategory);

  const filteredPast = selectedCategory === 'all'
    ? pastWebinars
    : pastWebinars.filter(w => w.category === selectedCategory);

  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.headerInfo}>
            <h1 className={styles.title}>
              {t('pageTitle', 'Vendor Webinars & Training')}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', 'Free live training sessions to help you succeed on our platform')}
            </p>
            
            {/* Stats */}
            <div className={styles.headerStats}>
              <div className={styles.stat}>
                <Video size={20} />
                <span>{t('stat1', '50+ webinars/year')}</span>
              </div>
              <div className={styles.stat}>
                <Users size={20} />
                <span>{t('stat2', '10K+ attendees')}</span>
              </div>
              <div className={styles.stat}>
                <Star size={20} />
                <span>{t('stat3', '4.8 avg rating')}</span>
              </div>
            </div>
          </div>
          
          <div className={styles.headerImage}>
            <div className={styles.webinarPreview}>
              <Monitor className={styles.previewIcon} />
              <PlayCircle className={styles.playIcon} />
            </div>
          </div>
        </div>
      </header>

      {/* Benefits */}
      <section className={styles.benefits}>
        <div className={styles.benefitsContainer}>
          {benefits.map((benefit, index) => (
            <div key={index} className={styles.benefitCard}>
              <benefit.icon className={styles.benefitIcon} size={24} />
              <h3>{benefit.title}</h3>
              <p>{benefit.description}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Category Filter */}
      <section className={styles.filterSection}>
        <div className={styles.filterContainer}>
          <h3>
            <Filter size={20} />
            {t('filterBy', 'Filter by Category')}
          </h3>
          <div className={styles.categoryFilter}>
            {categories.map(category => (
              <button
                key={category.id}
                onClick={() => setSelectedCategory(category.id)}
                className={`${styles.categoryButton} ${selectedCategory === category.id ? styles.active : ''}`}
              >
                <category.icon size={18} />
                <span>{category.label}</span>
              </button>
            ))}
          </div>
        </div>
      </section>

      {/* Main Content */}
      <main className={styles.mainContent}>
        {/* Upcoming Webinars */}
        <section className={styles.upcomingSection}>
          <div className={styles.sectionHeader}>
            <h2>
              <Calendar size={28} />
              {t('upcomingWebinars', 'Upcoming Webinars')}
            </h2>
            <span className={styles.count}>{filteredUpcoming.length} {t('scheduled', 'scheduled')}</span>
          </div>

          <div className={styles.webinarGrid}>
            {filteredUpcoming.map(webinar => (
              <div key={webinar.id} className={styles.webinarCard}>
                <div className={styles.webinarHeader}>
                  <div className={styles.webinarMeta}>
                    <span className={styles.category}>{webinar.category.replace('-', ' ')}</span>
                    <span className={styles.level}>{webinar.level}</span>
                  </div>
                  <div className={styles.attendeeCount}>
                    <Users size={16} />
                    <span>{webinar.attendees}/{webinar.maxAttendees}</span>
                  </div>
                </div>

                <h3 className={styles.webinarTitle}>{webinar.title}</h3>
                <p className={styles.webinarDescription}>{webinar.description}</p>

                <div className={styles.webinarDetails}>
                  <div className={styles.dateTime}>
                    <div>
                      <Calendar size={16} />
                      <span>{new Date(webinar.date).toLocaleDateString()}</span>
                    </div>
                    <div>
                      <Clock size={16} />
                      <span>{webinar.time}</span>
                    </div>
                    <div>
                      <Zap size={16} />
                      <span>{webinar.duration}</span>
                    </div>
                  </div>

                  <div className={styles.speaker}>
                    <div className={styles.speakerAvatar}>
                      <User size={20} />
                    </div>
                    <div>
                      <p className={styles.speakerName}>{webinar.speaker.name}</p>
                      <p className={styles.speakerRole}>{webinar.speaker.role}</p>
                    </div>
                  </div>

                  <div className={styles.topics}>
                    <h4>{t('whatYouWillLearn', "What you'll learn:")}</h4>
                    <ul>
                      {webinar.topics.map((topic, index) => (
                        <li key={index}>
                          <CheckCircle size={14} />
                          <span>{topic}</span>
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>

                <button 
                  onClick={() => handleRegister(webinar)}
                  className={styles.registerButton}
                >
                  {t('registerNow', 'Register Now')}
                  <ArrowRight size={16} />
                </button>
              </div>
            ))}
          </div>

          {filteredUpcoming.length === 0 && (
            <div className={styles.noResults}>
              <Calendar size={48} />
              <p>{t('noUpcoming', 'No upcoming webinars in this category')}</p>
            </div>
          )}
        </section>

        {/* Webinar Series */}
        <section className={styles.seriesSection}>
          <div className={styles.sectionHeader}>
            <h2>
              <Sparkles size={28} />
              {t('webinarSeries', 'Webinar Series')}
            </h2>
          </div>

          <div className={styles.seriesGrid}>
            {webinarSeries.map(series => (
              <div key={series.id} className={styles.seriesCard}>
                <div className={styles.seriesHeader}>
                  <h3>{series.title}</h3>
                  <span className={styles.schedule}>{series.schedule}</span>
                </div>
                <p className={styles.seriesDescription}>{series.description}</p>
                
                <div className={styles.seriesTopics}>
                  {series.topics.map((topic, index) => (
                    <span key={index} className={styles.topicTag}>{topic}</span>
                  ))}
                </div>
                
                <div className={styles.seriesFooter}>
                  <span className={styles.nextSession}>
                    <Calendar size={14} />
                    {t('next', 'Next')}: {series.upcoming}
                  </span>
                  <button className={styles.subscribeButton}>
                    <Bell size={16} />
                    {t('subscribe', 'Subscribe to Series')}
                  </button>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Past Webinars */}
        <section className={styles.pastSection}>
          <div className={styles.sectionHeader}>
            <h2>
              <PlayCircle size={28} />
              {t('pastWebinars', 'Past Webinar Recordings')}
            </h2>
            <span className={styles.count}>{filteredPast.length} {t('available', 'available')}</span>
          </div>

          <div className={styles.recordingsGrid}>
            {filteredPast.map(recording => (
              <div key={recording.id} className={styles.recordingCard}>
                <div className={styles.recordingThumbnail}>
                  <PlayCircle className={styles.playButton} size={40} />
                  <span className={styles.duration}>{recording.duration}</span>
                </div>
                
                <div className={styles.recordingContent}>
                  <h4>{recording.title}</h4>
                  <p>{recording.description}</p>
                  
                  <div className={styles.recordingMeta}>
                    <span className={styles.views}>
                      <Monitor size={14} />
                      {recording.views} views
                    </span>
                    <span className={styles.rating}>
                      <Star size={14} />
                      {recording.rating}/5
                    </span>
                    <span className={styles.date}>
                      {new Date(recording.date).toLocaleDateString()}
                    </span>
                  </div>
                  
                  <button className={styles.watchButton}>
                    {t('watchRecording', 'Watch Recording')}
                    <ChevronRight size={16} />
                  </button>
                </div>
              </div>
            ))}
          </div>

          {filteredPast.length === 0 && (
            <div className={styles.noResults}>
              <PlayCircle size={48} />
              <p>{t('noRecordings', 'No recordings in this category')}</p>
            </div>
          )}
        </section>

        {/* CTA Section */}
        <section className={styles.ctaSection}>
          <div className={styles.ctaContent}>
            <Building2 className={styles.ctaIcon} />
            <h2>{t('ctaTitle', 'Want a Custom Training for Your Team?')}</h2>
            <p>{t('ctaDesc', 'We offer private webinars and training sessions tailored to your business needs')}</p>
            <button 
              onClick={() => router.push('/contact/training')}
              className={styles.ctaButton}
            >
              {t('requestTraining', 'Request Custom Training')}
              <ArrowRight size={20} />
            </button>
          </div>
        </section>
      </main>

      {/* Registration Modal */}
      {showRegistrationForm && selectedWebinar && (
        <div className={styles.modal} onClick={() => setShowRegistrationForm(false)}>
          <div className={styles.modalContent} onClick={e => e.stopPropagation()}>
            <button 
              onClick={() => setShowRegistrationForm(false)}
              className={styles.closeButton}
            >
              ×
            </button>
            
            <h2>{t('registerFor', 'Register for Webinar')}</h2>
            <h3>{selectedWebinar.title}</h3>
            
            <form className={styles.registrationForm}>
              <div className={styles.formGroup}>
                <label>{t('fullName', 'Full Name')}</label>
                <input type="text" required />
              </div>
              
              <div className={styles.formGroup}>
                <label>{t('email', 'Email Address')}</label>
                <input type="email" required />
              </div>
              
              <div className={styles.formGroup}>
                <label>{t('businessName', 'Business Name')}</label>
                <input type="text" />
              </div>
              
              <div className={styles.formGroup}>
                <label>{t('businessSize', 'Business Size')}</label>
                <select>
                  <option>{t('size1', 'Just starting')}</option>
                  <option>{t('size2', '1-10 products')}</option>
                  <option>{t('size3', '11-100 products')}</option>
                  <option>{t('size4', '100+ products')}</option>
                </select>
              </div>
              
              <div className={styles.formInfo}>
                <Info size={16} />
                <p>{t('registrationInfo', "You'll receive a confirmation email with the webinar link")}</p>
              </div>
              
              <button type="submit" className={styles.submitButton}>
                {t('confirmRegistration', 'Confirm Registration')}
              </button>
            </form>
          </div>
        </div>
      )}
    </div>
  );
}