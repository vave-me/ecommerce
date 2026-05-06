"use client";

export const dynamic = 'force-dynamic';

import React from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  Store,
  MessageCircle,
  Users,
  Video,
  Heart,
  Star,
  Bell,
  Newspaper,
  ShoppingBag,
  Wrench,
  Camera,
  PenTool,
  Building2,
  CheckCircle, 
  TrendingUp, 
  Shield, 
  Zap,
  ArrowRight,
  Package,
  CreditCard,
  Headphones,
  BarChart3,
  Globe,
  Award,
  Check,
  Bot,
  Clock,
  Mail,
  MessageSquare,
  LineChart,
  Database,
  Lock,
  Layers,
  PhoneCall,
  FileText,
  Activity,
  Cloud,
  Key,
  Briefcase,
  GitBranch,
  ShieldCheck,
  Cpu,
  HelpCircle,
  Rocket,
  Server,
  UserPlus,
  Megaphone,
  Sparkles,
  Share2,
  Send,
  BookOpen,
  Palette,
  Coffee,
  Hammer
} from 'lucide-react';
import styles from './SellPage.module.css';

export default function SellPage() {
  const t = useTranslations('sell');
  const router = useRouter();
  const locale = router.locale || 'en';

  const businessTypes = [
    {
      icon: Store,
      title: t('businessType1', 'Retailers & Resellers'),
      description: t('businessType1Desc', 'Whether you have 1 product or 10,000 - from boutique shops to large retailers'),
      examples: t('businessType1Examples', 'Online stores • Physical shops • Distributors • Wholesalers')
    },
    {
      icon: Wrench,
      title: t('businessType2', 'Service Providers'),
      description: t('businessType2Desc', 'Offer your expertise and services directly to businesses and consumers'),
      examples: t('businessType2Examples', 'Installers • Consultants • Maintenance • Technical support')
    },
    {
      icon: PenTool,
      title: t('businessType3', 'Content Creators'),
      description: t('businessType3Desc', 'Share knowledge, build audience, and monetize your expertise'),
      examples: t('businessType3Examples', 'Bloggers • Educators • Influencers • Industry experts')
    },
    {
      icon: Hammer,
      title: t('businessType4', 'Manufacturers'),
      description: t('businessType4Desc', 'Connect directly with your B2B and B2C customers'),
      examples: t('businessType4Examples', 'Producers • Craftsmen • Custom manufacturers • OEM suppliers')
    }
  ];

  const ecosystemFeatures = [
    {
      icon: MessageCircle,
      title: t('ecosystem1Title', 'Live Messaging & Chat'),
      description: t('ecosystem1Desc', 'Real-time conversations with customers, instant notifications, and AI-powered chat assistant available 24/7')
    },
    {
      icon: Users,
      title: t('ecosystem2Title', 'Build Your Community'),
      description: t('ecosystem2Desc', 'Gain followers, create loyalty programs, and build lasting relationships with your customers')
    },
    {
      icon: Newspaper,
      title: t('ecosystem3Title', 'Content Publishing Platform'),
      description: t('ecosystem3Desc', 'Share blogs, articles, videos, tutorials - educate your audience and boost SEO rankings')
    },
    {
      icon: Star,
      title: t('ecosystem4Title', 'Reviews & Social Proof'),
      description: t('ecosystem4Desc', 'Collect authentic reviews, showcase testimonials, and build trust with social validation')
    },
    {
      icon: Mail,
      title: t('ecosystem5Title', 'Newsletter & Marketing'),
      description: t('ecosystem5Desc', 'Send targeted newsletters to your followers, announce new products, and share exclusive offers')
    },
    {
      icon: Video,
      title: t('ecosystem6Title', 'Live Streaming & Videos'),
      description: t('ecosystem6Desc', 'Host live product demos, Q&A sessions, and share video content with your audience')
    },
    {
      icon: Bell,
      title: t('ecosystem7Title', 'Smart Notifications'),
      description: t('ecosystem7Desc', 'Keep customers engaged with personalized updates about products, content, and special offers')
    },
    {
      icon: Share2,
      title: t('ecosystem8Title', 'Social Commerce Integration'),
      description: t('ecosystem8Desc', 'Connect with social media platforms, enable social sharing, and viral marketing features')
    }
  ];

  const joinReasons = [
    {
      icon: ShoppingBag,
      title: t('reason1Title', 'No Platform Dependencies'),
      text: t('reason1', 'Sell directly without relying on external marketplaces - own your customer relationships')
    },
    {
      icon: Sparkles,
      title: t('reason2Title', 'Start Small, Scale Big'),
      text: t('reason2', 'Whether you have 1 product or 10,000 - grow at your own pace with flexible plans')
    },
    {
      icon: UserPlus,
      title: t('reason3Title', 'Built-in Customer Base'),
      text: t('reason3', 'Access millions of active buyers looking for products and services in your category')
    },
    {
      icon: Megaphone,
      title: t('reason4Title', 'Complete Marketing Suite'),
      text: t('reason4', 'Everything you need to promote your business - from SEO tools to social features')
    },
    {
      icon: Coffee,
      title: t('reason5Title', 'Join a Thriving Ecosystem'),
      text: t('reason5', 'Connect with other businesses, share knowledge, and grow together')
    },
    {
      icon: Rocket,
      title: t('reason6Title', 'Future-Proof Technology'),
      text: t('reason6', 'AI-powered tools, automation, and continuous platform innovation')
    }
  ];

  const steps = [
    {
      number: '01',
      title: t('step1Title', 'Quick Sign-Up'),
      description: t('step1Desc', 'Register your business in minutes - solo entrepreneur or enterprise, all are welcome'),
      icon: UserPlus
    },
    {
      number: '02',
      title: t('step2Title', 'Build Your Presence'),
      description: t('step2Desc', 'Create your store, add products/services, publish content, customize your brand'),
      icon: Palette
    },
    {
      number: '03',
      title: t('step3Title', 'Engage & Connect'),
      description: t('step3Desc', 'Start conversations, gain followers, share expertise, build your community'),
      icon: MessageCircle
    },
    {
      number: '04',
      title: t('step4Title', 'Grow Your Business'),
      description: t('step4Desc', 'Use AI tools, analytics, and marketing features to scale your success'),
      icon: TrendingUp
    }
  ];

  const enterpriseFeatures = [
    {
      icon: Server,
      title: t('erpIntegrationTitle', 'Microservices Architecture'),
      description: t('erpIntegrationDesc', 'Kubernetes-based microservices with service mesh, auto-scaling, and distributed tracing. Docker containerization with CI/CD pipelines')
    },
    {
      icon: GitBranch,
      title: t('apiPlatformTitle', 'API Gateway & Management'),
      description: t('apiPlatformDesc', 'Kong/AWS API Gateway with rate limiting, OAuth2/JWT auth, GraphQL federation, webhook management, and OpenAPI documentation')
    },
    {
      icon: Database,
      title: t('dataLakeTitle', 'Data Infrastructure'),
      description: t('dataLakeDesc', 'PostgreSQL clusters, Redis caching, Elasticsearch, Apache Kafka streaming, TimescaleDB for time-series, and S3-compatible object storage')
    },
    {
      icon: Bot,
      title: t('aiTitle', 'AI/ML Platform'),
      description: t('aiDesc', 'TensorFlow/PyTorch models, GPT-4/Claude API integration, vector databases (Pinecone/Weaviate), ML pipelines with Kubeflow')
    },
    {
      icon: ShieldCheck,
      title: t('securityTitle', 'Security Infrastructure'),
      description: t('securityDesc', 'Zero-trust architecture, HashiCorp Vault, mTLS, SIEM integration, penetration testing, and automated security scanning')
    },
    {
      icon: Cloud,
      title: t('cloudNativeTitle', 'Cloud-Native Stack'),
      description: t('cloudNativeDesc', 'Multi-cloud support (AWS/GCP/Azure), Terraform IaC, GitOps with ArgoCD, Prometheus/Grafana monitoring, and ELK stack')
    },
    {
      icon: Zap,
      title: t('realtimeTitle', 'Real-Time Engine'),
      description: t('realtimeDesc', 'WebSocket infrastructure, Server-Sent Events, NATS/RabbitMQ messaging, Redis Pub/Sub, and WebRTC for video/audio')
    },
    {
      icon: Package,
      title: t('devToolsTitle', 'Developer Experience'),
      description: t('devToolsDesc', 'SDK libraries (Node.js, Python, Go, Java), CLI tools, Postman collections, sandbox environments, and comprehensive documentation')
    }
  ];

  const pricingPlans = [
    {
      name: t('planStarter', 'Starter'),
      price: t('planStarterPrice', 'FREE'),
      commission: t('planStarterCommission', '+ 8% transaction fee'),
      features: [
        t('feature1', 'Up to 100 products/services'),
        t('feature2', 'Basic store customization'),
        t('feature3', 'Live messaging & chat'),
        t('feature4', 'Blog & content publishing'),
        t('feature5', 'Email support'),
        t('feature6', 'Basic analytics'),
        t('feature7', 'Community features')
      ],
      cta: t('starterCta', 'Start Free')
    },
    {
      name: t('planProfessional', 'Professional'),
      price: t('planProfessionalPrice', '€49/month'),
      commission: t('planProfessionalCommission', '+ 5% transaction fee'),
      features: [
        t('feature8', 'Up to 1,000 products/services'),
        t('feature9', 'Advanced customization'),
        t('feature10', 'AI assistant & automation'),
        t('feature11', 'Newsletter to 5,000 contacts'),
        t('feature12', 'Priority support'),
        t('feature13', 'Advanced analytics'),
        t('feature14', 'API access'),
        t('feature15', 'Remove platform branding')
      ],
      popular: true,
      cta: t('proCta', 'Start 14-day Trial')
    },
    {
      name: t('planEnterprise', 'Enterprise'),
      price: t('planEnterprisePrice', 'Custom'),
      commission: t('planEnterpriseCommission', 'From 3% transaction fee'),
      features: [
        t('feature16', 'Unlimited products/services'),
        t('feature17', 'White-label solution'),
        t('feature18', 'Dedicated infrastructure'),
        t('feature19', 'Custom integrations'),
        t('feature20', 'Unlimited newsletters'),
        t('feature21', 'Dedicated account manager'),
        t('feature22', '24/7 phone support'),
        t('feature23', 'SLA guarantees'),
        t('feature24', 'Custom development')
      ],
      cta: t('enterpriseCta', 'Contact Sales')
    }
  ];

  // Platform metrics for enterprise credibility
  const platformMetrics = [
    {
      label: t('metric1Label', 'API Response Time'),
      value: t('metric1Value', '<50ms')
    },
    {
      label: t('metric2Label', 'Platform Uptime'),
      value: t('metric2Value', '99.99%')
    },
    {
      label: t('metric3Label', 'Daily Transactions'),
      value: t('metric3Value', '10M+')
    },
    {
      label: t('metric4Label', 'Global Regions'),
      value: t('metric4Value', '12')
    }
  ];

  // Performance Advantages
  const performanceMetrics = [
    { name: '10x Faster than Magento', icon: Zap },
    { name: '15x More Scalable', icon: TrendingUp },
    { name: '50ms Response Time', icon: Activity },
    { name: '99.99% Uptime Guaranteed', icon: Shield },
    { name: 'Auto-Scaling Infrastructure', icon: Layers },
    { name: 'Real-Time Processing', icon: Clock },
    { name: 'Enterprise-Grade Caching', icon: Database },
    { name: 'Global Edge Network', icon: Globe },
    { name: 'Microservices Architecture', icon: GitBranch },
    { name: 'AI-Optimized Performance', icon: Bot },
    { name: 'Zero Downtime Deployments', icon: Rocket },
    { name: 'Infinite Scalability', icon: Cloud }
  ];

  // Technical Capabilities
  const technicalCapabilities = [
    {
      icon: GitBranch,
      title: t('cicdTitle', 'CI/CD & DevOps'),
      description: t('cicdDesc', 'GitHub Actions/GitLab CI, automated testing, blue-green deployments, canary releases, and infrastructure as code with Terraform')
    },
    {
      icon: Activity,
      title: t('monitoringTitle', 'Observability Stack'),
      description: t('monitoringDesc', 'Prometheus metrics, Grafana dashboards, Jaeger distributed tracing, ELK stack for logs, and custom alerting rules')
    },
    {
      icon: Lock,
      title: t('authTitle', 'Authentication & Authorization'),
      description: t('authDesc', 'OAuth2/OIDC, SAML, multi-factor authentication, role-based access control (RBAC), and API key management')
    },
    {
      icon: Globe,
      title: t('cdnTitle', 'CDN & Edge Computing'),
      description: t('cdnDesc', 'Cloudflare/Fastly integration, edge workers, global content delivery, DDoS protection, and WAF rules')
    },
    {
      icon: MessageSquare,
      title: t('messagingTitle', 'Event-Driven Architecture'),
      description: t('messagingDesc', 'Apache Kafka streams, RabbitMQ/NATS messaging, event sourcing, CQRS pattern, and saga orchestration')
    }
  ];

  return (
    <div className={styles.container}>
      {/* Hero Section */}
      <section className={styles.hero}>
        <div className={styles.heroContent}>
          <div className={styles.enterpriseBadge}>
            {t('pageTitle', 'Business Ecosystem Platform')}
          </div>
          <h1 className={styles.heroTitle}>
            {t('heroTitle', 'Your All-in-One Business Platform: Sell, Connect, and Grow')}
          </h1>
          <p className={styles.heroSubtitle}>
            {t('heroSubtitle', 'Join thousands of businesses - from solo entrepreneurs to global brands. Sell products, offer services, publish content, and build your community. No dependencies, no intermediaries - just direct business growth.')}
          </p>
          <div className={styles.heroActions}>
            <button 
              onClick={() => router.push('/contact/sales')}
              className={styles.primaryButton}
            >
              {t('startSelling', 'Request Demo')}
              <ArrowRight size={18} />
            </button>
            <button 
              onClick={() => router.push('/resources/whitepaper')}
              className={styles.secondaryButton}
            >
              <FileText size={18} />
              {t('learnMore', 'Download Whitepaper')}
            </button>
            <button 
              onClick={() => router.push('/contact/consultation')}
              className={styles.tertiaryButton}
            >
              <PhoneCall size={18} />
              {t('scheduleCall', 'Schedule Consultation')}
            </button>
          </div>
          <div className={styles.heroStats}>
            <div className={styles.stat}>
              <strong>15,000+</strong>
              <span>{t('activeSellers', 'Active Businesses')}</span>
            </div>
            <div className={styles.stat}>
              <strong>2M+</strong>
              <span>{t('monthlyUsers', 'Monthly Active Users')}</span>
            </div>
            <div className={styles.stat}>
              <strong>500K+</strong>
              <span>{t('satisfaction', 'Daily Messages')}</span>
            </div>
            <div className={styles.stat}>
              <strong>4.8/5</strong>
              <span>{t('gmv', 'Seller Rating')}</span>
            </div>
          </div>
        </div>
        <div className={styles.heroVisual}>
          <div className={styles.metricCards}>
            {platformMetrics.map((metric, index) => (
              <div key={index} className={styles.metricCard}>
                <span className={styles.metricValue}>{metric.value}</span>
                <span className={styles.metricLabel}>{metric.label}</span>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Who Can Join Section */}
      <section className={styles.benefits}>
        <h2 className={styles.sectionTitle}>
          {t('whoCanJoinTitle', 'Every Business Has a Place Here')}
        </h2>
        <p className={styles.sectionSubtitle}>
          {t('whoCanJoinSubtitle', 'No matter your size or industry - if you have something to offer, you belong here')}
        </p>
        <div className={styles.benefitsGrid}>
          {businessTypes.map((type, index) => (
            <div key={index} className={styles.benefitCard}>
              <div className={styles.benefitIcon}>
                <type.icon size={24} />
              </div>
              <h3>{type.title}</h3>
              <p>{type.description}</p>
              <span className={styles.examples}>{type.examples}</span>
            </div>
          ))}
        </div>
      </section>

      {/* Ecosystem Features Section */}
      <section className={styles.enterpriseFeatures}>
        <h2 className={styles.sectionTitle}>
          {t('ecosystemTitle', 'More Than a Marketplace - A Complete Business Ecosystem')}
        </h2>
        <p className={styles.sectionSubtitle}>
          {t('ecosystemSubtitle', 'Everything you need to run, grow, and scale your business in one platform')}
        </p>
        <div className={styles.featuresGrid}>
          {ecosystemFeatures.map((feature, index) => (
            <div key={index} className={styles.featureCard}>
              <div className={styles.featureIcon}>
                <feature.icon size={20} />
              </div>
              <h3>{feature.title}</h3>
              <p>{feature.description}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Technical Capabilities */}
      <section className={styles.contentCommunity}>
        <h2 className={styles.sectionTitle}>
          {t('technicalCapabilitiesTitle', 'Advanced Technical Capabilities')}
        </h2>
        <div className={styles.featuresGrid}>
          {technicalCapabilities.map((capability, index) => (
            <div key={index} className={styles.featureCard}>
              <div className={styles.featureIcon}>
                <capability.icon size={20} />
              </div>
              <h3>{capability.title}</h3>
              <p>{capability.description}</p>
            </div>
          ))}
        </div>
      </section>

      {/* Performance Advantages */}
      <section className={styles.integrations}>
        <h2 className={styles.sectionTitle}>
          {t('performanceTitle', 'Unmatched Platform Performance')}
        </h2>
        <p className={styles.sectionSubtitle}>
          {t('performanceSubtitle', 'Engineered for speed, built for scale - outperforming traditional platforms')}
        </p>
        <div className={styles.integrationGrid}>
          {performanceMetrics.map((metric, index) => (
            <div key={index} className={styles.integrationCard}>
              <metric.icon size={32} className={styles.integrationIcon} />
              <span>{metric.name}</span>
            </div>
          ))}
        </div>
      </section>

      {/* Why Join Section */}
      <section className={styles.requirements}>
        <div className={styles.requirementsContent}>
          <h2 className={styles.sectionTitle}>
            {t('whyJoinTitle', 'Why Businesses Choose Our Platform')}
          </h2>
          <p className={styles.sectionSubtitle}>
            {t('whyJoinSubtitle', 'Join a thriving ecosystem designed for modern business success')}
          </p>
          <div className={styles.requirementsList}>
            {joinReasons.map((reason, index) => (
              <div key={index} className={styles.requirementItem}>
                <reason.icon className={styles.requirementIcon} size={20} />
                <div>
                  <strong>{reason.title}</strong>
                  <span>{reason.text}</span>
                </div>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Getting Started Process */}
      <section id="process" className={styles.process}>
        <h2 className={styles.sectionTitle}>
          {t('processTitle', 'Start Your Journey in Minutes')}
        </h2>
        <p className={styles.sectionSubtitle}>
          {t('processSubtitle', 'Simple onboarding process - be live in less than 24 hours')}
        </p>
        <div className={styles.stepsContainer}>
          {steps.map((step, index) => (
            <div key={index} className={styles.step}>
              <div className={styles.stepIcon}>
                <step.icon size={24} />
              </div>
              <div className={styles.stepNumber}>{step.number}</div>
              <h3>{step.title}</h3>
              <p>{step.description}</p>
              {index < steps.length - 1 && <div className={styles.stepConnector} />}
            </div>
          ))}
        </div>
      </section>

      {/* Enterprise Pricing */}
      <section className={styles.pricing}>
        <h2 className={styles.sectionTitle}>
          {t('pricingTitle', 'Enterprise Pricing Plans')}
        </h2>
        <p className={styles.pricingSubtitle}>
          {t('customPricing', 'Custom pricing based on volume and requirements')}
        </p>
        <div className={styles.pricingGrid}>
          {pricingPlans.map((plan, index) => (
            <div 
              key={index} 
              className={`${styles.pricingCard} ${plan.popular ? styles.popular : ''}`}
            >
              {plan.popular && (
                <div className={styles.popularBadge}>
                  Recommended
                </div>
              )}
              <h3>{plan.name}</h3>
              <div className={styles.price}>
                <strong>{plan.price}</strong>
                <span>{plan.commission}</span>
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
                onClick={() => {
                  if (plan.name === t('planStarter', 'Starter')) {
                    router.push('/register');
                  } else if (plan.name === t('planEnterprise', 'Enterprise')) {
                    router.push('/contact/sales');
                  } else {
                    router.push('/register?plan=professional');
                  }
                }}
                className={plan.popular ? styles.primaryButton : styles.outlineButton}
              >
                {plan.cta}
              </button>
            </div>
          ))}
        </div>
      </section>

      {/* Professional Services */}
      <section className={styles.services}>
        <h2 className={styles.sectionTitle}>
          {t('servicesTitle', 'Professional Services & Support')}
        </h2>
        <div className={styles.servicesGrid}>
          <div className={styles.serviceCard}>
            <Cpu className={styles.serviceIcon} />
            <h3>{t('service1', 'Platform Development')}</h3>
            <p>{t('service1Desc', 'Custom microservices, API development, frontend components, and system integrations')}</p>
          </div>
          <div className={styles.serviceCard}>
            <Cloud className={styles.serviceIcon} />
            <h3>{t('service2', 'DevOps & Infrastructure')}</h3>
            <p>{t('service2Desc', 'Kubernetes setup, CI/CD pipelines, monitoring, auto-scaling, and disaster recovery')}</p>
          </div>
          <div className={styles.serviceCard}>
            <Bot className={styles.serviceIcon} />
            <h3>{t('service3', 'AI/ML Solutions')}</h3>
            <p>{t('service3Desc', 'Custom model training, NLP implementations, recommendation engines, and predictive analytics')}</p>
          </div>
          <div className={styles.serviceCard}>
            <ShieldCheck className={styles.serviceIcon} />
            <h3>{t('service4', 'Security & Compliance')}</h3>
            <p>{t('service4Desc', 'Security audits, penetration testing, compliance implementation, and incident response')}</p>
          </div>
          <div className={styles.serviceCard}>
            <Database className={styles.serviceIcon} />
            <h3>{t('service5', 'Data Engineering')}</h3>
            <p>{t('service5Desc', 'Data pipelines, ETL processes, real-time analytics, and data warehouse design')}</p>
          </div>
          <div className={styles.serviceCard}>
            <Headphones className={styles.serviceIcon} />
            <h3>{t('service6', '24/7 Technical Support')}</h3>
            <p>{t('service6Desc', 'Round-the-clock monitoring, incident management, performance optimization, and SLA guarantees')}</p>
          </div>
        </div>
      </section>

      {/* Security & Compliance */}
      <section className={styles.security}>
        <h2 className={styles.sectionTitle}>
          {t('securityTitle', 'Enterprise Security & Compliance')}
        </h2>
        <p className={styles.sectionSubtitle}>
          {t('securitySubtitle', 'Your data protected by industry-leading standards')}
        </p>
        <div className={styles.certificationGrid}>
          <div className={styles.certCard}>
            <ShieldCheck size={48} className={styles.certIcon} />
            <h4>{t('cert1', 'ISO 27001:2013')}</h4>
          </div>
          <div className={styles.certCard}>
            <Lock size={48} className={styles.certIcon} />
            <h4>{t('cert2', 'SOC 2 Type II')}</h4>
          </div>
          <div className={styles.certCard}>
            <CreditCard size={48} className={styles.certIcon} />
            <h4>{t('cert3', 'PCI DSS Level 1')}</h4>
          </div>
          <div className={styles.certCard}>
            <Shield size={48} className={styles.certIcon} />
            <h4>{t('cert4', 'GDPR Compliant')}</h4>
          </div>
          <div className={styles.certCard}>
            <FileText size={48} className={styles.certIcon} />
            <h4>{t('cert5', 'CCPA Compliant')}</h4>
          </div>
          <div className={styles.certCard}>
            <Key size={48} className={styles.certIcon} />
            <h4>{t('cert6', 'HIPAA Ready')}</h4>
          </div>
        </div>
      </section>

      {/* Client Testimonials */}
      <section className={styles.testimonials}>
        <h2 className={styles.sectionTitle}>
          {t('trustTitle', 'Trusted by Industry Leaders')}
        </h2>
        <div className={styles.testimonialGrid}>
          <div className={styles.testimonialCard}>
            <blockquote>
              {t('clientTestimonial1', 'The seamless SAP integration saved us months of development time and reduced operational costs by 40%.')}
            </blockquote>
            <cite>{t('clientName1', 'CTO, Global Manufacturing Leader')}</cite>
          </div>
          <div className={styles.testimonialCard}>
            <blockquote>
              {t('clientTestimonial2', 'Best-in-class API documentation and support. We migrated 1M+ SKUs without any downtime.')}
            </blockquote>
            <cite>{t('clientName2', 'Head of E-commerce, Fashion Retailer')}</cite>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className={styles.cta}>
        <div className={styles.ctaContent}>
          <Zap className={styles.ctaIcon} />
          <h2>{t('ctaTitle', 'Ready to Transform Your Commerce Operations?')}</h2>
          <p>{t('ctaSubtitle', 'Join Fortune 500 companies and industry leaders on our platform')}</p>
          <div className={styles.ctaActions}>
            <button 
              onClick={() => router.push('/contact/sales')}
              className={styles.ctaButton}
            >
              {t('scheduleDemo', 'Schedule Executive Demo')}
              <ArrowRight size={16} />
            </button>
            <button 
              onClick={() => router.push('/resources/enterprise-brochure')}
              className={styles.ctaSecondaryButton}
            >
              <FileText size={16} />
              {t('downloadBrochure', 'Download Enterprise Brochure')}
            </button>
          </div>
        </div>
      </section>

      {/* Footer Links */}
      <section className={styles.footerLinks}>
        <div className={styles.linksContainer}>
          <a href="/faq/vendors" className={styles.footerLink}>
            <HelpCircle size={20} />
            {t('faqLinkText', 'View Enterprise FAQ')}
          </a>
          <a href="/legal/vendors" className={styles.footerLink}>
            <Shield size={20} />
            {t('legalLinkText', 'Legal & Compliance')}
          </a>
          <a href="/legal/dpa" className={styles.footerLink}>
            <Lock size={20} />
            {t('privacyLinkText', 'Data Processing Agreement')}
          </a>
          <a href="/legal/sla" className={styles.footerLink}>
            <FileText size={20} />
            {t('slaLinkText', 'Service Level Agreement')}
          </a>
        </div>
      </section>
    </div>
  );
}