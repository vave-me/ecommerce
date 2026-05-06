'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  Sparkles,
  Globe,
  Zap,
  Users,
  BarChart3,
  Server,
  Shield,
  Layers,
  ArrowRight,
  Check,
  MessageSquare,
  ShoppingBag,
  Brain,
  Building2,
  Code,
  Palette,
  Package,
  Headphones,
  LineChart,
  Cloud,
  Activity,
  Briefcase,
  Award,
  TrendingUp,
  HelpCircle,
  FileText,
  GitBranch,
  Database,
  Lock,
  Key,
  Bot,
  Rocket,
  Target,
  Settings,
  PieChart,
  Calendar,
  CreditCard,
  PhoneCall,
  Mail,
  Clock,
  CheckCircle,
  Cpu,
  ShieldCheck,
  Wrench,
  Users2,
  Network,
  Workflow
} from 'lucide-react';
import styles from './SellPage.module.css';

export default function MainSaaSPage() {
  const t = useTranslations('demo');
  const router = useRouter();

  const platformCapabilities = [
    {
      icon: Globe,
      title: t('capability1Title', 'Multi-Domain Ecosystem'),
      description: t('capability1Desc', 'Build platforms for any industry: from green energy to healthcare, from B2B marketplaces to community-driven networks')
    },
    {
      icon: Layers,
      title: t('capability2Title', 'Modular Architecture'),
      description: t('capability2Desc', 'Choose exactly what you need: commerce engine, social layer, AI intelligence, or combine all for a complete solution')
    },
    {
      icon: Bot,
      title: t('capability3Title', 'AI-Native Platform'),
      description: t('capability3Desc', 'Built-in AI co-pilot that understands your business, automates operations, and provides intelligent insights 24/7')
    },
    {
      icon: Shield,
      title: t('capability4Title', 'Enterprise Security'),
      description: t('capability4Desc', 'Bank-grade security with end-to-end encryption, compliance certifications, and continuous monitoring')
    }
  ];

  const ecosystemFeatures = [
    {
      icon: ShoppingBag,
      title: t('commerceTitle', 'Complete Commerce Engine'),
      description: t('commerceDesc', 'Everything for selling: products, services, subscriptions, digital goods. Multi-vendor support, global payments, intelligent logistics'),
      features: [
        'Multi-vendor marketplace',
        'Global payment orchestration',
        'Inventory & order management',
        'Dynamic pricing engine',
        'Subscription billing',
        'B2B/B2C/D2C modes'
      ]
    },
    {
      icon: Users,
      title: t('socialTitle', 'Social & Community Layer'),
      description: t('socialDesc', 'Transform transactions into relationships with real-time messaging, forums, reviews, and social engagement'),
      features: [
        'Real-time messaging',
        'Community forums',
        'User profiles & following',
        'Reviews & ratings',
        'Content sharing',
        'Live streaming'
      ]
    },
    {
      icon: Brain,
      title: t('intelligenceTitle', 'AI Intelligence Suite'),
      description: t('intelligenceDesc', '24/7 AI assistant that handles customer queries, automates tasks, and provides predictive analytics'),
      features: [
        'Conversational AI assistant',
        'Automated customer support',
        'Predictive analytics',
        'Personalization engine',
        'Voice & image search',
        'Smart recommendations'
      ]
    },
    {
      icon: Building2,
      title: t('verticalTitle', 'Industry Solutions'),
      description: t('verticalDesc', 'Pre-built templates and workflows for specific industries with compliance and best practices'),
      features: [
        'Green energy ecosystem',
        'Healthcare marketplace',
        'Real estate platform',
        'Education & courses',
        'Professional services',
        'Custom verticals'
      ]
    }
  ];

  const buildingBlocks = [
    {
      category: t('frontendCategory', 'Frontend & UX'),
      icon: Palette,
      items: [
        { name: 'React/Next.js framework', premium: true },
        { name: 'Mobile-first responsive design', premium: true },
        { name: 'Component library (100+ components)', premium: true },
        { name: 'White-label theming engine', premium: true },
        { name: 'Multi-language support', premium: true },
        { name: 'Progressive Web App', premium: true }
      ]
    },
    {
      category: t('backendCategory', 'Backend & APIs'),
      icon: Server,
      items: [
        { name: 'Microservices architecture', premium: true },
        { name: 'GraphQL & REST APIs', premium: true },
        { name: 'Event-driven architecture', premium: true },
        { name: 'Real-time WebSockets', premium: true },
        { name: 'API rate limiting & security', premium: true },
        { name: 'Webhooks & integrations', premium: true }
      ]
    },
    {
      category: t('dataCategory', 'Data & Analytics'),
      icon: Database,
      items: [
        { name: 'Multi-tenant architecture', premium: true },
        { name: 'Real-time analytics dashboard', premium: true },
        { name: 'Data lake & warehousing', premium: true },
        { name: 'Machine learning pipelines', premium: true },
        { name: 'Custom reporting engine', premium: true },
        { name: 'Export & API access', premium: true }
      ]
    },
    {
      category: t('infrastructureCategory', 'Infrastructure'),
      icon: Cloud,
      items: [
        { name: 'Auto-scaling Kubernetes', premium: true },
        { name: 'Multi-region deployment', premium: true },
        { name: 'CDN & edge computing', premium: true },
        { name: '99.99% uptime SLA', premium: true },
        { name: 'Disaster recovery', premium: true },
        { name: 'On-premise option', premium: true }
      ]
    }
  ];

  const implementationPlans = [
    {
      name: t('starterPlan', 'Starter Ecosystem'),
      price: t('starterPrice', '€9,999'),
      timeline: t('starterTimeline', '4-6 weeks'),
      features: [
        'Basic commerce engine',
        'User management system',
        'Payment integration',
        'Admin dashboard',
        'Mobile responsive design',
        'Basic analytics',
        '3 months support'
      ],
      best: false
    },
    {
      name: t('professionalPlan', 'Professional Platform'),
      price: t('professionalPrice', '€29,999'),
      timeline: t('professionalTimeline', '8-12 weeks'),
      features: [
        'Full commerce suite',
        'Social & community features',
        'AI assistant integration',
        'Multi-vendor support',
        'Advanced analytics',
        'Custom branding',
        'API access',
        '6 months support'
      ],
      best: true
    },
    {
      name: t('enterprisePlan', 'Enterprise Ecosystem'),
      price: t('enterprisePrice', 'Custom'),
      timeline: t('enterpriseTimeline', '3-6 months'),
      features: [
        'Complete custom solution',
        'All platform modules',
        'Custom AI training',
        'White-label everything',
        'Dedicated infrastructure',
        'Source code access',
        'Custom integrations',
        'Ongoing partnership'
      ],
      best: false
    }
  ];

  const successStories = [
    {
      industry: t('story1Industry', 'Green Energy'),
      metric: t('story1Metric', '500% Growth'),
      description: t('story1Desc', 'Solar panel marketplace connecting installers with homeowners, integrated with energy monitoring')
    },
    {
      industry: t('story2Industry', 'Healthcare'),
      metric: t('story2Metric', '2M+ Users'),
      description: t('story2Desc', 'Telemedicine platform with integrated pharmacy, appointment booking, and health records')
    },
    {
      industry: t('story3Industry', 'B2B Commerce'),
      metric: t('story3Metric', '€50M GMV'),
      description: t('story3Desc', 'Industrial supplies marketplace with quotation system, bulk ordering, and credit management')
    },
    {
      industry: t('story4Industry', 'Education'),
      metric: t('story4Metric', '100K+ Students'),
      description: t('story4Desc', 'Online learning platform with live classes, certifications, and community forums')
    }
  ];

  const processSteps = [
    {
      number: '01',
      title: t('process1Title', 'Discovery Workshop'),
      description: t('process1Desc', 'Deep dive into your vision, target market, and unique requirements')
    },
    {
      number: '02',
      title: t('process2Title', 'Architecture Design'),
      description: t('process2Desc', 'Custom solution blueprint with technology stack and integration plan')
    },
    {
      number: '03',
      title: t('process3Title', 'Agile Development'),
      description: t('process3Desc', 'Iterative development with weekly demos and continuous feedback')
    },
    {
      number: '04',
      title: t('process4Title', 'Launch & Training'),
      description: t('process4Desc', 'Smooth deployment with comprehensive training for your team')
    },
    {
      number: '05',
      title: t('process5Title', 'Growth Partnership'),
      description: t('process5Desc', 'Ongoing support, optimization, and feature development')
    }
  ];

  const technologies = [
    { name: 'Next.js 14', category: 'Frontend' },
    { name: 'React 19', category: 'Frontend' },
    { name: 'Node.js', category: 'Backend' },
    { name: 'Go', category: 'Backend' },
    { name: 'PostgreSQL', category: 'Database' },
    { name: 'Redis', category: 'Cache' },
    { name: 'Kubernetes', category: 'Infrastructure' },
    { name: 'OpenAI GPT-4', category: 'AI' },
    { name: 'Stripe', category: 'Payments' },
    { name: 'AWS/GCP/Azure', category: 'Cloud' },
    { name: 'GraphQL', category: 'API' },
    { name: 'NATS', category: 'Messaging' }
  ];

  const metrics = [
    { label: t('metric1', 'Platforms Built'), value: '250+' },
    { label: t('metric2', 'Total Users'), value: '10M+' },
    { label: t('metric3', 'Countries'), value: '45' },
    { label: t('metric4', 'Uptime'), value: '99.99%' }
  ];

  return (
    <div className={styles.container}>
      {/* Hero Section */}
      <section className={styles.hero}>
        <div className={styles.heroContent}>
          <div className={styles.enterpriseBadge}>
            <Sparkles size={16} />
            {t('badge', 'Enterprise Platform Solutions')}
          </div>
          <h1 className={styles.heroTitle}>
            {t('heroTitle', 'Build Your Own Digital Ecosystem with Our Battle-Tested Platform')}
          </h1>
          <p className={styles.heroSubtitle}>
            {t('heroSubtitle', 'From marketplace to social network, from SaaS to community platform - we provide the complete infrastructure and expertise to launch your vision in weeks, not years')}
          </p>
          <div className={styles.heroActions}>
            <button 
              onClick={() => router.push('/contact/platform')}
              className={styles.primaryButton}
            >
              {t('startBuilding', 'Start Building')}
              <ArrowRight size={18} />
            </button>
            <button 
              onClick={() => router.push('/demo/case-studies')}
              className={styles.secondaryButton}
            >
              <Award size={18} />
              {t('viewCases', 'View Case Studies')}
            </button>
            <button 
              onClick={() => router.push('/demo/architecture')}
              className={styles.tertiaryButton}
            >
              <GitBranch size={18} />
              {t('techDocs', 'Technical Documentation')}
            </button>
          </div>
          <div className={styles.heroStats}>
            {metrics.map((metric, index) => (
              <div key={index} className={styles.stat}>
                <strong>{metric.value}</strong>
                <span>{metric.label}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Platform Capabilities */}
      <section className={styles.benefits}>
        <h2 className={styles.sectionTitle}>
          {t('capabilitiesTitle', 'Why Build on Our Platform')}
        </h2>
        <div className={styles.benefitsGrid}>
          {platformCapabilities.map((capability, index) => (
            <div key={index} className={styles.benefitCard}>
              <div className={styles.benefitIcon}>
                <capability.icon size={24} />
              </div>
              <h3>{capability.title}</h3>
              <p>{capability.description}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Ecosystem Features */}
      <section className={styles.ecosystemFeatures}>
        <h2 className={styles.sectionTitle}>
          {t('ecosystemTitle', 'Complete Ecosystem Modules')}
        </h2>
        <p className={styles.sectionSubtitle}>
          {t('ecosystemSubtitle', 'Mix and match modules to create your perfect platform')}
        </p>
        <div className={styles.ecosystemGrid}>
          {ecosystemFeatures.map((feature, index) => (
            <div key={index} className={styles.ecosystemCard}>
              <div className={styles.ecosystemHeader}>
                <feature.icon size={32} className={styles.ecosystemIcon} />
                <h3>{feature.title}</h3>
              </div>
              <p className={styles.ecosystemDesc}>{feature.description}</p>
              <ul className={styles.ecosystemFeatures}>
                {feature.features.map((item, idx) => (
                  <li key={idx}>
                    <Check size={16} />
                    {item}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </section>

      {/* Building Blocks */}
      <section className={styles.buildingBlocks}>
        <h2 className={styles.sectionTitle}>
          {t('blocksTitle', 'Technical Building Blocks')}
        </h2>
        <p className={styles.sectionSubtitle}>
          {t('blocksSubtitle', 'Enterprise-grade components ready for scale')}
        </p>
        <div className={styles.blocksGrid}>
          {buildingBlocks.map((block, index) => (
            <div key={index} className={styles.blockCard}>
              <div className={styles.blockHeader}>
                <block.icon size={24} />
                <h3>{block.category}</h3>
              </div>
              <ul className={styles.blockItems}>
                {block.items.map((item, idx) => (
                  <li key={idx} className={item.premium ? styles.premium : ''}>
                    <CheckCircle size={16} />
                    {item.name}
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </section>

      {/* Success Stories */}
      <section className={styles.successStories}>
        <h2 className={styles.sectionTitle}>
          {t('storiesTitle', 'Success Stories')}
        </h2>
        <div className={styles.storiesGrid}>
          {successStories.map((story, index) => (
            <div key={index} className={styles.storyCard}>
              <div className={styles.storyIndustry}>{story.industry}</div>
              <div className={styles.storyMetric}>{story.metric}</div>
              <p className={styles.storyDesc}>{story.description}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Process */}
      <section className={styles.process}>
        <h2 className={styles.sectionTitle}>
          {t('processTitle', 'Our Implementation Process')}
        </h2>
        <div className={styles.stepsContainer}>
          {processSteps.map((step, index) => (
            <div key={index} className={styles.step}>
              <div className={styles.stepNumber}>{step.number}</div>
              <h3>{step.title}</h3>
              <p>{step.description}</p>
              {index < processSteps.length - 1 && <div className={styles.stepConnector} />}
            </div>
          ))}
        </div>
      </section>

      {/* Implementation Plans */}
      <section className={styles.pricing}>
        <h2 className={styles.sectionTitle}>
          {t('plansTitle', 'Implementation Packages')}
        </h2>
        <p className={styles.pricingSubtitle}>
          {t('plansSubtitle', 'Transparent pricing for every scale')}
        </p>
        <div className={styles.pricingGrid}>
          {implementationPlans.map((plan, index) => (
            <div 
              key={index} 
              className={`${styles.pricingCard} ${plan.best ? styles.popular : ''}`}
            >
              {plan.best && (
                <div className={styles.popularBadge}>
                  Most Popular
                </div>
              )}
              <h3>{plan.name}</h3>
              <div className={styles.price}>
                <strong>{plan.price}</strong>
                <span>{plan.timeline}</span>
              </div>
              <ul className={styles.features}>
                {plan.features.map((feature, idx) => (
                  <li key={idx}>
                    <Check size={16} />
                    {feature}
                  </li>
                ))}
              </ul>
              <button 
                onClick={() => router.push('/contact/platform')}
                className={plan.best ? styles.primaryButton : styles.outlineButton}
              >
                {t('getStarted', 'Get Started')}
              </button>
            </div>
          ))}
        </div>
      </section>

      {/* Technologies */}
      <section className={styles.technologies}>
        <h2 className={styles.sectionTitle}>
          {t('techTitle', 'Built with Best-in-Class Technologies')}
        </h2>
        <div className={styles.techGrid}>
          {technologies.map((tech, index) => (
            <div key={index} className={styles.techCard}>
              <span className={styles.techName}>{tech.name}</span>
              <span className={styles.techCategory}>{tech.category}</span>
            </div>
          ))}
        </div>
      </section>

      {/* Additional Services */}
      <section className={styles.services}>
        <h2 className={styles.sectionTitle}>
          {t('servicesTitle', 'Additional Services')}
        </h2>
        <div className={styles.servicesGrid}>
          <div className={styles.serviceCard}>
            <Wrench className={styles.serviceIcon} />
            <h3>{t('service1', 'Custom Development')}</h3>
            <p>{t('service1Desc', 'Unique features and integrations tailored to your specific needs')}</p>
          </div>
          <div className={styles.serviceCard}>
            <Users2 className={styles.serviceIcon} />
            <h3>{t('service2', 'Team Training')}</h3>
            <p>{t('service2Desc', 'Comprehensive training for your team to manage and grow the platform')}</p>
          </div>
          <div className={styles.serviceCard}>
            <LineChart className={styles.serviceIcon} />
            <h3>{t('service3', 'Growth Consulting')}</h3>
            <p>{t('service3Desc', 'Strategic guidance to scale your platform and acquire users')}</p>
          </div>
          <div className={styles.serviceCard}>
            <Shield className={styles.serviceIcon} />
            <h3>{t('service4', 'Security Audit')}</h3>
            <p>{t('service4Desc', 'Regular security assessments and compliance certifications')}</p>
          </div>
          <div className={styles.serviceCard}>
            <Cloud className={styles.serviceIcon} />
            <h3>{t('service5', 'DevOps Support')}</h3>
            <p>{t('service5Desc', '24/7 infrastructure monitoring and optimization')}</p>
          </div>
          <div className={styles.serviceCard}>
            <Sparkles className={styles.serviceIcon} />
            <h3>{t('service6', 'AI Training')}</h3>
            <p>{t('service6Desc', 'Custom AI model training for your specific use cases')}</p>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className={styles.cta}>
        <div className={styles.ctaContent}>
          <Rocket className={styles.ctaIcon} />
          <h2>{t('ctaTitle', 'Ready to Build Your Digital Empire?')}</h2>
          <p>{t('ctaSubtitle', 'Join hundreds of successful platforms powered by our technology')}</p>
          <div className={styles.ctaActions}>
            <button 
              onClick={() => router.push('/contact/platform')}
              className={styles.ctaButton}
            >
              {t('scheduleConsultation', 'Schedule Free Consultation')}
              <ArrowRight size={16} />
            </button>
            <button 
              onClick={() => router.push('/demo/roi-calculator')}
              className={styles.ctaSecondaryButton}
            >
              <PieChart size={16} />
              {t('calculateROI', 'Calculate Your ROI')}
            </button>
          </div>
        </div>
      </section>

      {/* Footer Links */}
      <section className={styles.footerLinks}>
        <div className={styles.linksContainer}>
          <a href="/demo/technical" className={styles.footerLink}>
            <Code size={20} />
            {t('techDocsLink', 'Technical Documentation')}
          </a>
          <a href="/demo/api" className={styles.footerLink}>
            <GitBranch size={20} />
            {t('apiDocsLink', 'API Reference')}
          </a>
          <a href="/demo/case-studies" className={styles.footerLink}>
            <Award size={20} />
            {t('caseStudiesLink', 'Case Studies')}
          </a>
          <a href="/demo/faq" className={styles.footerLink}>
            <HelpCircle size={20} />
            {t('faqLink', 'Platform FAQ')}
          </a>
        </div>
      </section>
    </div>
  );
}