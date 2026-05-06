"use client";

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  FileText,
  Download,
  CheckCircle,
  Globe,
  Shield,
  TrendingUp,
  Zap,
  Database,
  Cloud,
  Lock,
  ArrowRight,
  Mail,
  Building2,
  Users,
  BarChart3,
  Cpu,
  GitBranch,
  Bot,
  Activity,
  Layers
} from 'lucide-react';
import styles from './Whitepaper.module.css';

export default function WhitepaperPage() {
  const t = useTranslations('whitepaper');
  const router = useRouter();
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    company: '',
    jobTitle: '',
    phoneNumber: '',
    country: '',
    consent: false
  });
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [showThankYou, setShowThankYou] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSubmitting(true);
    
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500));
    
    setIsSubmitting(false);
    setShowThankYou(true);
    
    // In production, this would trigger email with download link
    setTimeout(() => {
      window.open('/downloads/platform-whitepaper.pdf', '_blank');
    }, 2000);
  };

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const whitepaperTopics = [
    {
      icon: Layers,
      title: t('topic1Title', 'Microservices Architecture'),
      description: t('topic1Desc', 'Deep dive into our cloud-native microservices design patterns and implementation strategies')
    },
    {
      icon: Bot,
      title: t('topic2Title', 'AI & Machine Learning'),
      description: t('topic2Desc', 'How we leverage GPT-4, custom models, and predictive analytics for commerce optimization')
    },
    {
      icon: Shield,
      title: t('topic3Title', 'Security & Compliance'),
      description: t('topic3Desc', 'Enterprise-grade security measures, zero-trust architecture, and regulatory compliance')
    },
    {
      icon: GitBranch,
      title: t('topic4Title', 'API Design Philosophy'),
      description: t('topic4Desc', 'RESTful and GraphQL API patterns, versioning strategies, and developer experience')
    },
    {
      icon: Cloud,
      title: t('topic5Title', 'Scalability & Performance'),
      description: t('topic5Desc', 'Auto-scaling strategies, caching layers, and global distribution architecture')
    },
    {
      icon: Database,
      title: t('topic6Title', 'Data Architecture'),
      description: t('topic6Desc', 'Event sourcing, CQRS patterns, and real-time data processing pipelines')
    }
  ];

  const keyInsights = [
    {
      icon: TrendingUp,
      stat: '10x',
      label: t('insight1', 'Faster deployment than traditional platforms')
    },
    {
      icon: Zap,
      stat: '50ms',
      label: t('insight2', 'Average API response time globally')
    },
    {
      icon: Users,
      stat: '99.99%',
      label: t('insight3', 'Platform uptime SLA guarantee')
    },
    {
      icon: Globe,
      stat: '15+',
      label: t('insight4', 'Global data center locations')
    }
  ];

  if (showThankYou) {
    return (
      <div className={styles.container}>
        <div className={styles.thankYouSection}>
          <div className={styles.thankYouContent}>
            <CheckCircle className={styles.successIcon} size={64} />
            <h1>{t('thankYouTitle', 'Thank You for Your Interest!')}</h1>
            <p>{t('thankYouMessage', 'Your whitepaper download will begin shortly. Check your email for additional resources.')}</p>
            <div className={styles.thankYouActions}>
              <button 
                onClick={() => router.push('/contact/sales')}
                className={styles.primaryButton}
              >
                {t('scheduleDemoButton', 'Schedule a Demo')}
                <ArrowRight size={18} />
              </button>
              <button 
                onClick={() => router.push('/resources')}
                className={styles.secondaryButton}
              >
                {t('moreResourcesButton', 'Explore More Resources')}
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className={styles.container}>
      {/* Hero Section */}
      <section className={styles.hero}>
        <div className={styles.heroContent}>
          <div className={styles.badge}>
            {t('badge', 'Technical Whitepaper')}
          </div>
          <h1 className={styles.heroTitle}>
            {t('heroTitle', 'Building Next-Generation Commerce Infrastructure')}
          </h1>
          <p className={styles.heroSubtitle}>
            {t('heroSubtitle', 'A comprehensive guide to our enterprise-grade platform architecture, technology stack, and implementation strategies')}
          </p>
          <div className={styles.heroMeta}>
            <span>{t('pages', '45 pages')}</span>
            <span>{t('readTime', '20 min read')}</span>
            <span>{t('updated', 'Updated Q1 2024')}</span>
          </div>
        </div>
        <div className={styles.heroVisual}>
          <div className={styles.documentPreview}>
            <FileText size={48} />
            <span>{t('format', 'PDF Document')}</span>
          </div>
        </div>
      </section>

      {/* Topics Section */}
      <section className={styles.topicsSection}>
        <h2 className={styles.sectionTitle}>
          {t('topicsTitle', 'What You\'ll Learn')}
        </h2>
        <div className={styles.topicsGrid}>
          {whitepaperTopics.map((topic, index) => (
            <div key={index} className={styles.topicCard}>
              <div className={styles.topicIcon}>
                <topic.icon size={24} />
              </div>
              <h3>{topic.title}</h3>
              <p>{topic.description}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Key Insights */}
      <section className={styles.insightsSection}>
        <h2 className={styles.sectionTitle}>
          {t('insightsTitle', 'Key Platform Metrics')}
        </h2>
        <div className={styles.insightsGrid}>
          {keyInsights.map((insight, index) => (
            <div key={index} className={styles.insightCard}>
              <insight.icon className={styles.insightIcon} size={32} />
              <div className={styles.insightStat}>{insight.stat}</div>
              <div className={styles.insightLabel}>{insight.label}</div>
            </div>
          ))}
        </div>
      </section>

      {/* Download Form */}
      <section className={styles.formSection}>
        <div className={styles.formContainer}>
          <div className={styles.formHeader}>
            <h2>{t('formTitle', 'Download Your Free Copy')}</h2>
            <p>{t('formSubtitle', 'Fill out the form below to receive instant access to our technical whitepaper')}</p>
          </div>
          
          <form onSubmit={handleSubmit} className={styles.form}>
            <div className={styles.formRow}>
              <div className={styles.formGroup}>
                <label htmlFor="firstName">{t('firstName', 'First Name')} *</label>
                <input
                  type="text"
                  id="firstName"
                  name="firstName"
                  value={formData.firstName}
                  onChange={handleChange}
                  required
                  placeholder={t('firstNamePlaceholder', 'John')}
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="lastName">{t('lastName', 'Last Name')} *</label>
                <input
                  type="text"
                  id="lastName"
                  name="lastName"
                  value={formData.lastName}
                  onChange={handleChange}
                  required
                  placeholder={t('lastNamePlaceholder', 'Doe')}
                />
              </div>
            </div>

            <div className={styles.formRow}>
              <div className={styles.formGroup}>
                <label htmlFor="email">{t('businessEmail', 'Business Email')} *</label>
                <input
                  type="email"
                  id="email"
                  name="email"
                  value={formData.email}
                  onChange={handleChange}
                  required
                  placeholder={t('emailPlaceholder', 'redacted-email@example.com')}
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="company">{t('company', 'Company')} *</label>
                <input
                  type="text"
                  id="company"
                  name="company"
                  value={formData.company}
                  onChange={handleChange}
                  required
                  placeholder={t('companyPlaceholder', 'Your Company')}
                />
              </div>
            </div>

            <div className={styles.formRow}>
              <div className={styles.formGroup}>
                <label htmlFor="jobTitle">{t('jobTitle', 'Job Title')} *</label>
                <input
                  type="text"
                  id="jobTitle"
                  name="jobTitle"
                  value={formData.jobTitle}
                  onChange={handleChange}
                  required
                  placeholder={t('jobTitlePlaceholder', 'CTO / Tech Lead')}
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="phoneNumber">{t('phone', 'Phone Number')}</label>
                <input
                  type="tel"
                  id="phoneNumber"
                  name="phoneNumber"
                  value={formData.phoneNumber}
                  onChange={handleChange}
                  placeholder={t('phonePlaceholder', '+1 (555) 123-4567')}
                />
              </div>
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="country">{t('country', 'Country')} *</label>
              <select
                id="country"
                name="country"
                value={formData.country}
                onChange={handleChange}
                required
              >
                <option value="">{t('selectCountry', 'Select your country')}</option>
                <option value="US">United States</option>
                <option value="UK">United Kingdom</option>
                <option value="DE">Germany</option>
                <option value="FR">France</option>
                <option value="IT">Italy</option>
                <option value="ES">Spain</option>
                <option value="PL">Poland</option>
                <option value="NL">Netherlands</option>
                <option value="other">{t('other', 'Other')}</option>
              </select>
            </div>

            <div className={styles.consentGroup}>
              <label>
                <input
                  type="checkbox"
                  name="consent"
                  checked={formData.consent}
                  onChange={handleChange}
                  required
                />
                <span>
                  {t('consent', 'I agree to receive communications about platform updates, best practices, and relevant offers. I can unsubscribe at any time.')}
                </span>
              </label>
            </div>

            <button 
              type="submit" 
              className={styles.submitButton}
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <>{t('downloading', 'Processing...')}</>
              ) : (
                <>
                  <Download size={18} />
                  {t('downloadButton', 'Download Whitepaper')}
                </>
              )}
            </button>
          </form>

          <div className={styles.formFooter}>
            <Lock size={16} />
            <span>{t('privacy', 'Your information is secure and will never be shared')}</span>
          </div>
        </div>
      </section>

      {/* Additional Resources */}
      <section className={styles.resourcesSection}>
        <h2 className={styles.sectionTitle}>
          {t('relatedResourcesTitle', 'Related Resources')}
        </h2>
        <div className={styles.resourcesGrid}>
          <a href="/resources/enterprise-brochure" className={styles.resourceCard}>
            <Building2 className={styles.resourceIcon} />
            <h3>{t('resource1Title', 'Enterprise Brochure')}</h3>
            <p>{t('resource1Desc', 'Executive overview of platform capabilities')}</p>
            <span className={styles.resourceLink}>
              {t('downloadNow', 'Download Now')}
              <ArrowRight size={16} />
            </span>
          </a>
          <a href="/docs/api" className={styles.resourceCard}>
            <Cpu className={styles.resourceIcon} />
            <h3>{t('resource2Title', 'API Documentation')}</h3>
            <p>{t('resource2Desc', 'Complete technical reference and guides')}</p>
            <span className={styles.resourceLink}>
              {t('explore', 'Explore Docs')}
              <ArrowRight size={16} />
            </span>
          </a>
          <a href="/webinars" className={styles.resourceCard}>
            <Activity className={styles.resourceIcon} />
            <h3>{t('resource3Title', 'Technical Webinars')}</h3>
            <p>{t('resource3Desc', 'Live demos and implementation workshops')}</p>
            <span className={styles.resourceLink}>
              {t('register', 'Register Now')}
              <ArrowRight size={16} />
            </span>
          </a>
        </div>
      </section>
    </div>
  );
}