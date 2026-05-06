"use client";

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  FileText,
  Download,
  CheckCircle,
  Building2,
  TrendingUp,
  Users,
  Globe,
  Shield,
  Zap,
  ArrowRight,
  Mail,
  BarChart3,
  Package,
  DollarSign,
  Award,
  Briefcase,
  Target,
  Rocket,
  Clock,
  Lock
} from 'lucide-react';
import styles from './EnterpriseBrochure.module.css';

export default function EnterpriseBrochurePage() {
  const t = useTranslations('enterpriseBrochure');
  const router = useRouter();
  const [formData, setFormData] = useState({
    firstName: '',
    lastName: '',
    email: '',
    company: '',
    jobTitle: '',
    companySize: '',
    industry: '',
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
      window.open('/downloads/enterprise-brochure.pdf', '_blank');
    }, 2000);
  };

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData(prev => ({
      ...prev,
      [name]: type === 'checkbox' ? checked : value
    }));
  };

  const brochureHighlights = [
    {
      icon: Building2,
      title: t('highlight1Title', 'Enterprise Solutions'),
      description: t('highlight1Desc', 'Complete B2B/B2C platform designed for large-scale operations')
    },
    {
      icon: TrendingUp,
      title: t('highlight2Title', 'Proven ROI'),
      description: t('highlight2Desc', 'Average 40% cost reduction and 3x revenue growth in year one')
    },
    {
      icon: Globe,
      title: t('highlight3Title', 'Global Reach'),
      description: t('highlight3Desc', 'Multi-region deployment with 15+ data centers worldwide')
    },
    {
      icon: Shield,
      title: t('highlight4Title', 'Security First'),
      description: t('highlight4Desc', 'SOC2, ISO27001 certified with enterprise-grade protection')
    }
  ];

  const successMetrics = [
    {
      icon: DollarSign,
      metric: '€2.5B+',
      label: t('metric1', 'Annual GMV processed')
    },
    {
      icon: Users,
      metric: '5,000+',
      label: t('metric2', 'Enterprise vendors')
    },
    {
      icon: Package,
      metric: '10M+',
      label: t('metric3', 'Monthly transactions')
    },
    {
      icon: Award,
      metric: '99.99%',
      label: t('metric4', 'Platform uptime')
    }
  ];

  const industries = [
    t('industry1', 'Manufacturing'),
    t('industry2', 'Wholesale & Distribution'),
    t('industry3', 'Retail & E-commerce'),
    t('industry4', 'Technology & Software'),
    t('industry5', 'Fashion & Apparel'),
    t('industry6', 'Electronics'),
    t('industry7', 'Home & Garden'),
    t('industry8', 'Healthcare'),
    t('industry9', 'Automotive'),
    t('industry10', 'Food & Beverage')
  ];

  if (showThankYou) {
    return (
      <div className={styles.container}>
        <div className={styles.thankYouSection}>
          <div className={styles.thankYouContent}>
            <CheckCircle className={styles.successIcon} size={64} />
            <h1>{t('thankYouTitle', 'Thank You!')}</h1>
            <p>{t('thankYouMessage', 'Your enterprise brochure download will begin shortly. A member of our team will contact you within 24 hours.')}</p>
            <div className={styles.thankYouActions}>
              <button 
                onClick={() => router.push('/contact/sales')}
                className={styles.primaryButton}
              >
                {t('scheduleConsultation', 'Schedule Consultation')}
                <ArrowRight size={18} />
              </button>
              <button 
                onClick={() => router.push('/sell')}
                className={styles.secondaryButton}
              >
                {t('learnMore', 'Learn More About Platform')}
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
            {t('badge', 'Executive Overview')}
          </div>
          <h1 className={styles.heroTitle}>
            {t('heroTitle', 'Enterprise Commerce Platform')}
          </h1>
          <p className={styles.heroSubtitle}>
            {t('heroSubtitle', 'Discover how leading companies transform their commerce operations with our comprehensive B2B/B2C platform')}
          </p>
          <div className={styles.heroStats}>
            {successMetrics.map((metric, index) => (
              <div key={index} className={styles.statCard}>
                <metric.icon className={styles.statIcon} size={24} />
                <div className={styles.statContent}>
                  <div className={styles.statMetric}>{metric.metric}</div>
                  <div className={styles.statLabel}>{metric.label}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
        <div className={styles.heroVisual}>
          <div className={styles.brochurePreview}>
            <div className={styles.previewHeader}>
              <FileText size={32} />
              <span>{t('brochureType', 'Enterprise Brochure')}</span>
            </div>
            <div className={styles.previewDetails}>
              <span>{t('pages', '20 pages')}</span>
              <span>{t('format', 'PDF format')}</span>
              <span>{t('language', 'Available in 5 languages')}</span>
            </div>
          </div>
        </div>
      </section>

      {/* Highlights Section */}
      <section className={styles.highlightsSection}>
        <h2 className={styles.sectionTitle}>
          {t('highlightsTitle', 'What\'s Inside')}
        </h2>
        <div className={styles.highlightsGrid}>
          {brochureHighlights.map((highlight, index) => (
            <div key={index} className={styles.highlightCard}>
              <div className={styles.highlightIcon}>
                <highlight.icon size={24} />
              </div>
              <h3>{highlight.title}</h3>
              <p>{highlight.description}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Industries Section */}
      <section className={styles.industriesSection}>
        <h2 className={styles.sectionTitle}>
          {t('industriesTitle', 'Trusted Across Industries')}
        </h2>
        <p className={styles.sectionSubtitle}>
          {t('industriesSubtitle', 'Our platform powers commerce for leading companies in')}
        </p>
        <div className={styles.industriesGrid}>
          {industries.map((industry, index) => (
            <div key={index} className={styles.industryCard}>
              <Target className={styles.industryIcon} size={20} />
              <span>{industry}</span>
            </div>
          ))}
        </div>
      </section>

      {/* Brochure Content Preview */}
      <section className={styles.contentSection}>
        <h2 className={styles.sectionTitle}>
          {t('contentTitle', 'Comprehensive Business Overview')}
        </h2>
        <div className={styles.contentGrid}>
          <div className={styles.contentCard}>
            <Briefcase className={styles.contentIcon} />
            <h3>{t('content1Title', 'Executive Summary')}</h3>
            <p>{t('content1Desc', 'Platform vision, market opportunity, and strategic advantages')}</p>
          </div>
          <div className={styles.contentCard}>
            <BarChart3 className={styles.contentIcon} />
            <h3>{t('content2Title', 'Business Impact')}</h3>
            <p>{t('content2Desc', 'ROI analysis, case studies, and performance metrics')}</p>
          </div>
          <div className={styles.contentCard}>
            <Zap className={styles.contentIcon} />
            <h3>{t('content3Title', 'Platform Capabilities')}</h3>
            <p>{t('content3Desc', 'Core features, integrations, and technical infrastructure')}</p>
          </div>
          <div className={styles.contentCard}>
            <Users className={styles.contentIcon} />
            <h3>{t('content4Title', 'Success Stories')}</h3>
            <p>{t('content4Desc', 'Client testimonials and transformation journeys')}</p>
          </div>
          <div className={styles.contentCard}>
            <Rocket className={styles.contentIcon} />
            <h3>{t('content5Title', 'Implementation')}</h3>
            <p>{t('content5Desc', 'Onboarding process, timeline, and support structure')}</p>
          </div>
          <div className={styles.contentCard}>
            <DollarSign className={styles.contentIcon} />
            <h3>{t('content6Title', 'Investment & ROI')}</h3>
            <p>{t('content6Desc', 'Pricing models, cost savings, and value proposition')}</p>
          </div>
        </div>
      </section>

      {/* Download Form */}
      <section className={styles.formSection}>
        <div className={styles.formContainer}>
          <div className={styles.formHeader}>
            <h2>{t('formTitle', 'Get Your Enterprise Brochure')}</h2>
            <p>{t('formSubtitle', 'Complete the form for immediate access to our comprehensive enterprise overview')}</p>
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
                  placeholder={t('jobTitlePlaceholder', 'CEO / Director')}
                />
              </div>
              <div className={styles.formGroup}>
                <label htmlFor="companySize">{t('companySize', 'Company Size')} *</label>
                <select
                  id="companySize"
                  name="companySize"
                  value={formData.companySize}
                  onChange={handleChange}
                  required
                >
                  <option value="">{t('selectSize', 'Select company size')}</option>
                  <option value="1-50">1-50 employees</option>
                  <option value="51-200">51-200 employees</option>
                  <option value="201-500">201-500 employees</option>
                  <option value="501-1000">501-1000 employees</option>
                  <option value="1000+">1000+ employees</option>
                </select>
              </div>
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="industry">{t('industry', 'Industry')} *</label>
              <select
                id="industry"
                name="industry"
                value={formData.industry}
                onChange={handleChange}
                required
              >
                <option value="">{t('selectIndustry', 'Select your industry')}</option>
                {industries.map((industry, index) => (
                  <option key={index} value={industry}>{industry}</option>
                ))}
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
                  {t('consent', 'I agree to receive the enterprise brochure and occasional updates about platform features and best practices.')}
                </span>
              </label>
            </div>

            <button 
              type="submit" 
              className={styles.submitButton}
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <>{t('processing', 'Processing...')}</>
              ) : (
                <>
                  <Download size={18} />
                  {t('downloadButton', 'Download Brochure')}
                </>
              )}
            </button>
          </form>

          <div className={styles.formFooter}>
            <Lock size={16} />
            <span>{t('privacy', 'Your information is secure and will only be used to provide relevant resources')}</span>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className={styles.ctaSection}>
        <div className={styles.ctaContent}>
          <Clock className={styles.ctaIcon} size={48} />
          <h2>{t('ctaTitle', 'Ready to Transform Your Business?')}</h2>
          <p>{t('ctaSubtitle', 'Join industry leaders who have revolutionized their commerce operations')}</p>
          <div className={styles.ctaActions}>
            <button 
              onClick={() => router.push('/contact/sales')}
              className={styles.ctaButton}
            >
              {t('talkToExpert', 'Talk to an Expert')}
              <ArrowRight size={16} />
            </button>
            <button 
              onClick={() => router.push('/resources/whitepaper')}
              className={styles.ctaSecondaryButton}
            >
              <FileText size={16} />
              {t('technicalWhitepaper', 'Technical Whitepaper')}
            </button>
          </div>
        </div>
      </section>
    </div>
  );
}