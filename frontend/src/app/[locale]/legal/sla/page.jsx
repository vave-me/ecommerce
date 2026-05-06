"use client";

import React from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  Shield,
  Clock,
  CheckCircle,
  AlertTriangle,
  Activity,
  TrendingUp,
  Server,
  Zap,
  AlertCircle,
  BarChart3,
  Phone,
  Mail,
  Download,
  ArrowLeft,
  CreditCard,
  Calendar,
  Target,
  Settings,
  Gauge,
  Users,
  HeadphonesIcon,
  FileText
} from 'lucide-react';
import styles from './SLA.module.css';

export default function SLAPage() {
  const t = useTranslations('sla');
  const router = useRouter();

  const serviceCategories = [
    {
      icon: Server,
      title: t('category1', 'Platform Availability'),
      target: '99.9%',
      description: t('category1Desc', 'Core platform services and API endpoints')
    },
    {
      icon: Zap,
      title: t('category2', 'API Response Time'),
      target: '<50ms',
      description: t('category2Desc', 'Average response time for API calls')
    },
    {
      icon: Activity,
      title: t('category3', 'Transaction Processing'),
      target: '<2s',
      description: t('category3Desc', 'End-to-end transaction completion time')
    },
    {
      icon: Shield,
      title: t('category4', 'Security Incident Response'),
      target: '<1hr',
      description: t('category4Desc', 'Initial response to critical security incidents')
    }
  ];

  const uptimeTargets = [
    {
      plan: t('planStarter', 'Starter'),
      monthly: '99.0%',
      quarterly: '99.5%',
      yearly: '99.5%',
      allowedDowntime: t('downtime1', '7.2 hours/month')
    },
    {
      plan: t('planProfessional', 'Professional'),
      monthly: '99.9%',
      quarterly: '99.9%',
      yearly: '99.9%',
      allowedDowntime: t('downtime2', '43.8 minutes/month')
    },
    {
      plan: t('planEnterprise', 'Enterprise'),
      monthly: '99.95%',
      quarterly: '99.95%',
      yearly: '99.99%',
      allowedDowntime: t('downtime3', '21.9 minutes/month')
    }
  ];

  const responseTimeTargets = [
    {
      severity: t('severity1', 'Critical'),
      icon: AlertTriangle,
      color: '#ef4444',
      description: t('critical', 'Complete service outage or data loss'),
      initialResponse: t('time1', '15 minutes'),
      updateFrequency: t('freq1', 'Every 30 minutes'),
      resolution: t('resolve1', '2 hours')
    },
    {
      severity: t('severity2', 'High'),
      icon: AlertCircle,
      color: '#f97316',
      description: t('high', 'Major functionality impaired'),
      initialResponse: t('time2', '1 hour'),
      updateFrequency: t('freq2', 'Every 2 hours'),
      resolution: t('resolve2', '8 hours')
    },
    {
      severity: t('severity3', 'Medium'),
      icon: Activity,
      color: '#f59e0b',
      description: t('medium', 'Partial functionality impaired'),
      initialResponse: t('time3', '4 hours'),
      updateFrequency: t('freq3', 'Daily'),
      resolution: t('resolve3', '24 hours')
    },
    {
      severity: t('severity4', 'Low'),
      icon: Settings,
      color: '#3b82f6',
      description: t('low', 'Minor issues or questions'),
      initialResponse: t('time4', '24 hours'),
      updateFrequency: t('freq4', 'Weekly'),
      resolution: t('resolve4', '5 business days')
    }
  ];

  const creditSchedule = [
    { downtime: '0-99.9%', credit: '0%' },
    { downtime: '99.5-99.9%', credit: '10%' },
    { downtime: '99.0-99.5%', credit: '25%' },
    { downtime: '95.0-99.0%', credit: '50%' },
    { downtime: '<95.0%', credit: '100%' }
  ];

  const exclusions = [
    t('exclusion1', 'Scheduled maintenance (with 48 hours notice)'),
    t('exclusion2', 'Force majeure events'),
    t('exclusion3', 'Customer-caused issues or misconfigurations'),
    t('exclusion4', 'Third-party service provider outages'),
    t('exclusion5', 'Beta features and experimental APIs'),
    t('exclusion6', 'Actions of third parties or denial of service attacks')
  ];

  const supportChannels = [
    {
      icon: Mail,
      channel: t('channel1', 'Email Support'),
      availability: t('availability1', '24/7'),
      plans: ['Starter', 'Professional', 'Enterprise']
    },
    {
      icon: HeadphonesIcon,
      channel: t('channel2', 'Live Chat'),
      availability: t('availability2', 'Business hours'),
      plans: ['Professional', 'Enterprise']
    },
    {
      icon: Phone,
      channel: t('channel3', 'Phone Support'),
      availability: t('availability3', '24/7'),
      plans: ['Enterprise']
    },
    {
      icon: Users,
      channel: t('channel4', 'Dedicated Account Manager'),
      availability: t('availability4', 'Business hours'),
      plans: ['Enterprise']
    }
  ];

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
              {t('pageTitle', 'Service Level Agreement')}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', 'Our commitment to platform reliability and performance')}
            </p>
            <div className={styles.metadata}>
              <span>{t('version', 'Version 3.0')}</span>
              <span>{t('effectiveDate', 'Effective: January 1, 2024')}</span>
              <span>{t('lastReview', 'Last Review: December 2023')}</span>
            </div>
          </div>
        </div>
      </header>

      {/* Key Metrics Overview */}
      <section className={styles.metricsOverview}>
        <div className={styles.metricsContainer}>
          <h2>{t('keyMetrics', 'Service Level Objectives')}</h2>
          <div className={styles.metricsGrid}>
            {serviceCategories.map((category, index) => (
              <div key={index} className={styles.metricCard}>
                <div className={styles.metricIcon}>
                  <category.icon size={24} />
                </div>
                <h3>{category.title}</h3>
                <div className={styles.metricTarget}>{category.target}</div>
                <p>{category.description}</p>
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Main Content */}
      <main className={styles.mainContent}>
        {/* Service Availability */}
        <section className={styles.section}>
          <h2>{t('availabilityTitle', 'Service Availability Commitment')}</h2>
          <p className={styles.sectionIntro}>
            {t('availabilityIntro', 'We guarantee the following uptime percentages for our platform services based on your subscription plan:')}
          </p>
          
          <div className={styles.uptimeTable}>
            <table>
              <thead>
                <tr>
                  <th>{t('plan', 'Plan')}</th>
                  <th>{t('monthly', 'Monthly')}</th>
                  <th>{t('quarterly', 'Quarterly')}</th>
                  <th>{t('yearly', 'Yearly')}</th>
                  <th>{t('allowedDowntime', 'Allowed Downtime')}</th>
                </tr>
              </thead>
              <tbody>
                {uptimeTargets.map((target, index) => (
                  <tr key={index}>
                    <td className={styles.planName}>{target.plan}</td>
                    <td>{target.monthly}</td>
                    <td>{target.quarterly}</td>
                    <td>{target.yearly}</td>
                    <td className={styles.downtime}>{target.allowedDowntime}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className={styles.calculation}>
            <h3>{t('calculationTitle', 'Availability Calculation')}</h3>
            <div className={styles.formula}>
              <code>
                Availability % = (Total Minutes - Downtime Minutes) / Total Minutes × 100
              </code>
            </div>
            <p>{t('calculationNote', 'Measured on a calendar month basis, excluding scheduled maintenance')}</p>
          </div>
        </section>

        {/* Response Time Commitments */}
        <section className={styles.section}>
          <h2>{t('responseTitle', 'Support Response Time Commitments')}</h2>
          <p className={styles.sectionIntro}>
            {t('responseIntro', 'Our support response times are based on the severity of the issue:')}</p>
          
          <div className={styles.responseGrid}>
            {responseTimeTargets.map((target, index) => (
              <div key={index} className={styles.responseCard}>
                <div className={styles.responseHeader} style={{ borderColor: target.color }}>
                  <target.icon size={24} style={{ color: target.color }} />
                  <h3>{target.severity}</h3>
                </div>
                <p className={styles.responseDescription}>{target.description}</p>
                <div className={styles.responseMetrics}>
                  <div>
                    <span className={styles.label}>{t('initialResponse', 'Initial Response')}</span>
                    <span className={styles.value}>{target.initialResponse}</span>
                  </div>
                  <div>
                    <span className={styles.label}>{t('updates', 'Updates')}</span>
                    <span className={styles.value}>{target.updateFrequency}</span>
                  </div>
                  <div>
                    <span className={styles.label}>{t('resolution', 'Resolution')}</span>
                    <span className={styles.value}>{target.resolution}</span>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Performance Targets */}
        <section className={styles.section}>
          <h2>{t('performanceTitle', 'Performance Targets')}</h2>
          <div className={styles.performanceGrid}>
            <div className={styles.performanceCard}>
              <Gauge className={styles.performanceIcon} />
              <h3>{t('apiLatency', 'API Latency')}</h3>
              <div className={styles.performanceMetric}>
                <span className={styles.bigNumber}>50ms</span>
                <span className={styles.unit}>{t('p95', '95th percentile')}</span>
              </div>
              <p>{t('apiLatencyDesc', 'For standard API endpoints under normal load')}</p>
            </div>
            
            <div className={styles.performanceCard}>
              <Target className={styles.performanceIcon} />
              <h3>{t('errorRate', 'Error Rate')}</h3>
              <div className={styles.performanceMetric}>
                <span className={styles.bigNumber}>&lt;0.1%</span>
                <span className={styles.unit}>{t('of', 'of requests')}</span>
              </div>
              <p>{t('errorRateDesc', 'Server-side errors (5xx status codes)')}</p>
            </div>
            
            <div className={styles.performanceCard}>
              <Activity className={styles.performanceIcon} />
              <h3>{t('throughput', 'Throughput')}</h3>
              <div className={styles.performanceMetric}>
                <span className={styles.bigNumber}>10K</span>
                <span className={styles.unit}>{t('rps', 'requests/second')}</span>
              </div>
              <p>{t('throughputDesc', 'Sustained throughput per customer')}</p>
            </div>
          </div>
        </section>

        {/* Service Credits */}
        <section className={styles.section}>
          <h2>{t('creditsTitle', 'Service Credits')}</h2>
          <p className={styles.sectionIntro}>
            {t('creditsIntro', 'If we fail to meet our SLA commitments, you are eligible for service credits:')}</p>
          
          <div className={styles.creditTable}>
            <table>
              <thead>
                <tr>
                  <th>{t('availability', 'Monthly Availability')}</th>
                  <th>{t('creditPercentage', 'Service Credit')}</th>
                </tr>
              </thead>
              <tbody>
                {creditSchedule.map((credit, index) => (
                  <tr key={index}>
                    <td>{credit.downtime}</td>
                    <td className={styles.creditAmount}>{credit.credit}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          <div className={styles.creditNotes}>
            <h3>{t('creditProcess', 'Credit Request Process')}</h3>
            <ol>
              <li>{t('creditStep1', 'Submit credit request within 30 days of the incident')}</li>
              <li>{t('creditStep2', 'Include your account information and calculation of downtime')}</li>
              <li>{t('creditStep3', 'Credits will be applied to your next billing cycle')}</li>
              <li>{t('creditStep4', 'Maximum credit per month is 100% of monthly fees')}</li>
            </ol>
          </div>
        </section>

        {/* Support Channels */}
        <section className={styles.section}>
          <h2>{t('supportTitle', 'Support Channels by Plan')}</h2>
          <div className={styles.supportGrid}>
            {supportChannels.map((channel, index) => (
              <div key={index} className={styles.supportCard}>
                <div className={styles.supportIcon}>
                  <channel.icon size={24} />
                </div>
                <h3>{channel.channel}</h3>
                <p className={styles.availability}>{channel.availability}</p>
                <div className={styles.planBadges}>
                  {channel.plans.map((plan, idx) => (
                    <span key={idx} className={styles.planBadge}>{plan}</span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </section>

        {/* Exclusions */}
        <section className={styles.section}>
          <h2>{t('exclusionsTitle', 'SLA Exclusions')}</h2>
          <p className={styles.sectionIntro}>
            {t('exclusionsIntro', 'The following circumstances are excluded from SLA calculations:')}</p>
          
          <ul className={styles.exclusionsList}>
            {exclusions.map((exclusion, index) => (
              <li key={index}>
                <AlertCircle size={16} />
                <span>{exclusion}</span>
              </li>
            ))}
          </ul>
        </section>

        {/* Monitoring and Reporting */}
        <section className={styles.section}>
          <h2>{t('monitoringTitle', 'Monitoring and Reporting')}</h2>
          <div className={styles.monitoringGrid}>
            <div className={styles.monitoringCard}>
              <BarChart3 className={styles.monitoringIcon} />
              <h3>{t('statusPage', 'Real-time Status Page')}</h3>
              <p>{t('statusPageDesc', 'Live platform status and incident history at status.platform.com')}</p>
              <a href="https://status.platform.com" className={styles.statusLink}>
                {t('viewStatus', 'View Status Page')} →
              </a>
            </div>
            
            <div className={styles.monitoringCard}>
              <Activity className={styles.monitoringIcon} />
              <h3>{t('performanceDashboard', 'Performance Dashboard')}</h3>
              <p>{t('dashboardDesc', 'Access detailed performance metrics in your account dashboard')}</p>
            </div>
            
            <div className={styles.monitoringCard}>
              <Calendar className={styles.monitoringIcon} />
              <h3>{t('monthlyReports', 'Monthly SLA Reports')}</h3>
              <p>{t('reportsDesc', 'Detailed availability and performance reports delivered monthly')}</p>
            </div>
          </div>
        </section>

        {/* Contact Information */}
        <section className={styles.contactSection}>
          <h2>{t('contactTitle', 'SLA Support Contacts')}</h2>
          <div className={styles.contactGrid}>
            <div className={styles.contactCard}>
              <Mail className={styles.contactIcon} />
              <h3>{t('slaTeam', 'SLA Team')}</h3>
              <p>redacted-email@example.com</p>
              <p>{t('slaResponse', 'For SLA credit requests')}</p>
            </div>
            <div className={styles.contactCard}>
              <Phone className={styles.contactIcon} />
              <h3>{t('emergencySupport', 'Emergency Support')}</h3>
              <p>+1 (555) 911-2345</p>
              <p>{t('emergencyAvailable', '24/7 for critical issues')}</p>
            </div>
            <div className={styles.contactCard}>
              <HeadphonesIcon className={styles.contactIcon} />
              <h3>{t('technicalSupport', 'Technical Support')}</h3>
              <p>redacted-email@example.com</p>
              <p>{t('techSupportHours', 'Based on your plan')}</p>
            </div>
          </div>
        </section>

        {/* Agreement Actions */}
        <section className={styles.actions}>
          <button className={styles.downloadButton}>
            <Download size={20} />
            {t('downloadSLA', 'Download SLA Document')}
          </button>
          <button 
            onClick={() => router.push('/contact/support')}
            className={styles.contactButton}
          >
            <HeadphonesIcon size={20} />
            {t('contactSupport', 'Contact Support')}
          </button>
          <button className={styles.printButton}>
            <FileText size={20} />
            {t('print', 'Print SLA')}
          </button>
        </section>
      </main>
    </div>
  );
}