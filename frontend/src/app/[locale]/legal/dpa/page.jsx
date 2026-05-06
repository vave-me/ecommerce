"use client";

import React from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  Shield,
  Lock,
  FileText,
  CheckCircle,
  Globe,
  Server,
  Key,
  AlertCircle,
  Users,
  Database,
  Cloud,
  ShieldCheck,
  Mail,
  Phone,
  Download,
  ArrowLeft,
  ExternalLink,
  Scale,
  Eye,
  UserCheck,
  RefreshCw,
  Trash2,
  Clock
} from 'lucide-react';
import styles from './DPA.module.css';

export default function DPAPage() {
  const t = useTranslations('dpa');
  const router = useRouter();

  const sections = [
    { id: 'definitions', title: t('section1', 'Definitions'), icon: FileText },
    { id: 'scope', title: t('section2', 'Scope and Application'), icon: Globe },
    { id: 'obligations', title: t('section3', 'Processor Obligations'), icon: Shield },
    { id: 'security', title: t('section4', 'Security Measures'), icon: Lock },
    { id: 'subprocessing', title: t('section5', 'Sub-processing'), icon: Users },
    { id: 'rights', title: t('section6', 'Data Subject Rights'), icon: UserCheck },
    { id: 'breach', title: t('section7', 'Data Breach Notification'), icon: AlertCircle },
    { id: 'audit', title: t('section8', 'Audit and Inspection'), icon: Eye },
    { id: 'return', title: t('section9', 'Return and Deletion'), icon: Trash2 },
    { id: 'liability', title: t('section10', 'Liability and Indemnity'), icon: Scale }
  ];

  const securityMeasures = [
    {
      icon: Lock,
      title: t('measure1Title', 'Encryption'),
      description: t('measure1Desc', 'AES-256 encryption at rest and TLS 1.3 in transit for all personal data')
    },
    {
      icon: Key,
      title: t('measure2Title', 'Access Control'),
      description: t('measure2Desc', 'Role-based access control with multi-factor authentication and principle of least privilege')
    },
    {
      icon: Shield,
      title: t('measure3Title', 'Network Security'),
      description: t('measure3Desc', 'Firewalls, intrusion detection systems, and DDoS protection on all endpoints')
    },
    {
      icon: Database,
      title: t('measure4Title', 'Data Segregation'),
      description: t('measure4Desc', 'Logical separation of customer data with dedicated encryption keys per tenant')
    },
    {
      icon: RefreshCw,
      title: t('measure5Title', 'Backup & Recovery'),
      description: t('measure5Desc', 'Automated daily backups with point-in-time recovery and geo-redundant storage')
    },
    {
      icon: Eye,
      title: t('measure6Title', 'Monitoring & Logging'),
      description: t('measure6Desc', '24/7 security monitoring with comprehensive audit logs and anomaly detection')
    }
  ];

  const dataSubjectRights = [
    {
      icon: Eye,
      right: t('right1', 'Right of Access'),
      description: t('right1Desc', 'Access personal data we process about them')
    },
    {
      icon: FileText,
      right: t('right2', 'Right to Rectification'),
      description: t('right2Desc', 'Correct inaccurate or incomplete personal data')
    },
    {
      icon: Trash2,
      right: t('right3', 'Right to Erasure'),
      description: t('right3Desc', 'Request deletion of personal data in certain circumstances')
    },
    {
      icon: Lock,
      right: t('right4', 'Right to Restriction'),
      description: t('right4Desc', 'Restrict processing of personal data in specific situations')
    },
    {
      icon: Download,
      right: t('right5', 'Right to Data Portability'),
      description: t('right5Desc', 'Receive personal data in a structured, machine-readable format')
    },
    {
      icon: AlertCircle,
      right: t('right6', 'Right to Object'),
      description: t('right6Desc', 'Object to processing based on legitimate interests')
    }
  ];

  const scrollToSection = (sectionId) => {
    const element = document.getElementById(sectionId);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
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
              {t('pageTitle', 'Data Processing Agreement')}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', 'GDPR-compliant agreement for processing personal data')}
            </p>
            <div className={styles.metadata}>
              <span>{t('version', 'Version 2.0')}</span>
              <span>{t('updated', 'Last updated: January 2024')}</span>
              <span>{t('language', 'Available in 5 languages')}</span>
            </div>
          </div>
        </div>
      </header>

      {/* Navigation */}
      <nav className={styles.navigation}>
        <div className={styles.navContent}>
          <h3>{t('tableOfContents', 'Table of Contents')}</h3>
          <ul className={styles.navList}>
            {sections.map((section) => (
              <li key={section.id}>
                <button
                  onClick={() => scrollToSection(section.id)}
                  className={styles.navItem}
                >
                  <section.icon size={16} />
                  <span>{section.title}</span>
                </button>
              </li>
            ))}
          </ul>
          <button className={styles.downloadButton}>
            <Download size={18} />
            {t('downloadPDF', 'Download PDF')}
          </button>
        </div>
      </nav>

      {/* Main Content */}
      <main className={styles.mainContent}>
        {/* Introduction */}
        <section className={styles.introduction}>
          <div className={styles.introCard}>
            <ShieldCheck className={styles.introIcon} size={48} />
            <div>
              <h2>{t('introTitle', 'Data Protection Commitment')}</h2>
              <p>{t('introText', 'This Data Processing Agreement ("DPA") forms part of the Contract for Services between the Platform (Processor) and the Customer (Controller) to reflect the parties\' agreement with regard to the Processing of Personal Data in accordance with GDPR requirements.')}</p>
            </div>
          </div>
        </section>

        {/* Definitions Section */}
        <section id="definitions" className={styles.section}>
          <h2>{t('definitionsTitle', '1. Definitions')}</h2>
          <div className={styles.definitionsList}>
            <div className={styles.definition}>
              <strong>{t('term1', 'Personal Data')}</strong>
              <p>{t('term1Def', 'Any information relating to an identified or identifiable natural person')}</p>
            </div>
            <div className={styles.definition}>
              <strong>{t('term2', 'Processing')}</strong>
              <p>{t('term2Def', 'Any operation performed on personal data, whether or not by automated means')}</p>
            </div>
            <div className={styles.definition}>
              <strong>{t('term3', 'Controller')}</strong>
              <p>{t('term3Def', 'The entity that determines the purposes and means of processing personal data')}</p>
            </div>
            <div className={styles.definition}>
              <strong>{t('term4', 'Processor')}</strong>
              <p>{t('term4Def', 'The entity that processes personal data on behalf of the Controller')}</p>
            </div>
          </div>
        </section>

        {/* Scope Section */}
        <section id="scope" className={styles.section}>
          <h2>{t('scopeTitle', '2. Scope and Application')}</h2>
          <p>{t('scopeIntro', 'This DPA applies to all processing of personal data by the Processor on behalf of the Controller in connection with the Platform Services.')}</p>
          <div className={styles.scopeGrid}>
            <div className={styles.scopeCard}>
              <h3>{t('scope1', 'Categories of Data Subjects')}</h3>
              <ul>
                <li>{t('dataSubject1', 'Customer employees and contractors')}</li>
                <li>{t('dataSubject2', 'End users and customers')}</li>
                <li>{t('dataSubject3', 'Vendors and suppliers')}</li>
                <li>{t('dataSubject4', 'Website visitors')}</li>
              </ul>
            </div>
            <div className={styles.scopeCard}>
              <h3>{t('scope2', 'Types of Personal Data')}</h3>
              <ul>
                <li>{t('dataType1', 'Contact information (names, emails, phones)')}</li>
                <li>{t('dataType2', 'Account and profile data')}</li>
                <li>{t('dataType3', 'Transaction and payment information')}</li>
                <li>{t('dataType4', 'Usage data and analytics')}</li>
              </ul>
            </div>
          </div>
        </section>

        {/* Processor Obligations */}
        <section id="obligations" className={styles.section}>
          <h2>{t('obligationsTitle', '3. Processor Obligations')}</h2>
          <div className={styles.obligationsList}>
            <div className={styles.obligation}>
              <CheckCircle className={styles.obligationIcon} />
              <div>
                <h4>{t('obligation1', 'Process only on documented instructions')}</h4>
                <p>{t('obligation1Desc', 'The Processor shall process Personal Data only on documented instructions from the Controller')}</p>
              </div>
            </div>
            <div className={styles.obligation}>
              <CheckCircle className={styles.obligationIcon} />
              <div>
                <h4>{t('obligation2', 'Ensure confidentiality')}</h4>
                <p>{t('obligation2Desc', 'Ensure all personnel are subject to appropriate confidentiality obligations')}</p>
              </div>
            </div>
            <div className={styles.obligation}>
              <CheckCircle className={styles.obligationIcon} />
              <div>
                <h4>{t('obligation3', 'Implement security measures')}</h4>
                <p>{t('obligation3Desc', 'Implement appropriate technical and organizational measures to ensure security')}</p>
              </div>
            </div>
            <div className={styles.obligation}>
              <CheckCircle className={styles.obligationIcon} />
              <div>
                <h4>{t('obligation4', 'Assist with compliance')}</h4>
                <p>{t('obligation4Desc', 'Assist the Controller in ensuring compliance with GDPR obligations')}</p>
              </div>
            </div>
          </div>
        </section>

        {/* Security Measures */}
        <section id="security" className={styles.section}>
          <h2>{t('securityTitle', '4. Technical and Organizational Security Measures')}</h2>
          <p className={styles.sectionIntro}>{t('securityIntro', 'The Processor implements and maintains the following security measures:')}</p>
          <div className={styles.measuresGrid}>
            {securityMeasures.map((measure, index) => (
              <div key={index} className={styles.measureCard}>
                <div className={styles.measureIcon}>
                  <measure.icon size={24} />
                </div>
                <h3>{measure.title}</h3>
                <p>{measure.description}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Sub-processing */}
        <section id="subprocessing" className={styles.section}>
          <h2>{t('subprocessingTitle', '5. Sub-processing')}</h2>
          <div className={styles.subprocessingContent}>
            <p>{t('subprocessingIntro', 'The Controller provides general authorization for the Processor to engage sub-processors, subject to:')}</p>
            <ul className={styles.requirementsList}>
              <li>{t('subReq1', 'Prior written notification of any intended changes')}</li>
              <li>{t('subReq2', 'Opportunity for the Controller to object to such changes')}</li>
              <li>{t('subReq3', 'Flow-down of all data protection obligations')}</li>
              <li>{t('subReq4', 'Processor remains fully liable for sub-processor performance')}</li>
            </ul>
            <div className={styles.subprocessorList}>
              <h3>{t('currentSubprocessors', 'Current Sub-processors')}</h3>
              <table className={styles.subprocessorTable}>
                <thead>
                  <tr>
                    <th>{t('name', 'Name')}</th>
                    <th>{t('service', 'Service')}</th>
                    <th>{t('location', 'Location')}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>AWS</td>
                    <td>{t('cloudHosting', 'Cloud Infrastructure')}</td>
                    <td>{t('euRegion', 'EU (Frankfurt)')}</td>
                  </tr>
                  <tr>
                    <td>Stripe</td>
                    <td>{t('paymentProcessing', 'Payment Processing')}</td>
                    <td>{t('euRegion', 'EU (Dublin)')}</td>
                  </tr>
                  <tr>
                    <td>SendGrid</td>
                    <td>{t('emailService', 'Email Services')}</td>
                    <td>{t('euRegion', 'EU (Ireland)')}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </section>

        {/* Data Subject Rights */}
        <section id="rights" className={styles.section}>
          <h2>{t('rightsTitle', '6. Data Subject Rights')}</h2>
          <p className={styles.sectionIntro}>{t('rightsIntro', 'The Processor shall assist the Controller in fulfilling its obligations to respond to data subject requests:')}</p>
          <div className={styles.rightsGrid}>
            {dataSubjectRights.map((item, index) => (
              <div key={index} className={styles.rightCard}>
                <div className={styles.rightIcon}>
                  <item.icon size={20} />
                </div>
                <h4>{item.right}</h4>
                <p>{item.description}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Data Breach */}
        <section id="breach" className={styles.section}>
          <h2>{t('breachTitle', '7. Data Breach Notification')}</h2>
          <div className={styles.breachProcess}>
            <div className={styles.breachStep}>
              <div className={styles.stepNumber}>1</div>
              <h3>{t('breach1', 'Immediate Notification')}</h3>
              <p>{t('breach1Desc', 'Notify Controller without undue delay and within 24 hours of becoming aware')}</p>
            </div>
            <div className={styles.breachStep}>
              <div className={styles.stepNumber}>2</div>
              <h3>{t('breach2', 'Detailed Information')}</h3>
              <p>{t('breach2Desc', 'Provide nature of breach, categories affected, and likely consequences')}</p>
            </div>
            <div className={styles.breachStep}>
              <div className={styles.stepNumber}>3</div>
              <h3>{t('breach3', 'Mitigation Measures')}</h3>
              <p>{t('breach3Desc', 'Describe measures taken or proposed to address the breach')}</p>
            </div>
            <div className={styles.breachStep}>
              <div className={styles.stepNumber}>4</div>
              <h3>{t('breach4', 'Documentation')}</h3>
              <p>{t('breach4Desc', 'Document all breaches and measures taken for compliance records')}</p>
            </div>
          </div>
        </section>

        {/* Audit Rights */}
        <section id="audit" className={styles.section}>
          <h2>{t('auditTitle', '8. Audit and Inspection Rights')}</h2>
          <div className={styles.auditContent}>
            <p>{t('auditIntro', 'The Controller has the right to audit Processor compliance:')}</p>
            <div className={styles.auditOptions}>
              <div className={styles.auditCard}>
                <FileText className={styles.auditIcon} />
                <h3>{t('audit1', 'Certifications')}</h3>
                <p>{t('audit1Desc', 'Annual SOC2 Type II and ISO 27001 certification reports')}</p>
              </div>
              <div className={styles.auditCard}>
                <Eye className={styles.auditIcon} />
                <h3>{t('audit2', 'On-site Audits')}</h3>
                <p>{t('audit2Desc', 'Annual right to conduct on-site audits with 30 days notice')}</p>
              </div>
              <div className={styles.auditCard}>
                <Shield className={styles.auditIcon} />
                <h3>{t('audit3', 'Security Assessments')}</h3>
                <p>{t('audit3Desc', 'Quarterly security assessment reports and penetration test results')}</p>
              </div>
            </div>
          </div>
        </section>

        {/* Return and Deletion */}
        <section id="return" className={styles.section}>
          <h2>{t('returnTitle', '9. Return and Deletion of Data')}</h2>
          <div className={styles.returnProcess}>
            <p>{t('returnIntro', 'Upon termination of services, the Processor shall:')}</p>
            <div className={styles.returnSteps}>
              <div className={styles.returnStep}>
                <Clock className={styles.returnIcon} />
                <h4>{t('return1', 'Within 30 days')}</h4>
                <p>{t('return1Desc', 'Return all personal data in commonly used format')}</p>
              </div>
              <div className={styles.returnStep}>
                <Trash2 className={styles.returnIcon} />
                <h4>{t('return2', 'Delete all copies')}</h4>
                <p>{t('return2Desc', 'Securely delete all copies unless retention is required by law')}</p>
              </div>
              <div className={styles.returnStep}>
                <FileText className={styles.returnIcon} />
                <h4>{t('return3', 'Provide certification')}</h4>
                <p>{t('return3Desc', 'Certify in writing that deletion has been completed')}</p>
              </div>
            </div>
          </div>
        </section>

        {/* Liability */}
        <section id="liability" className={styles.section}>
          <h2>{t('liabilityTitle', '10. Liability and Indemnity')}</h2>
          <div className={styles.liabilityContent}>
            <div className={styles.liabilityCard}>
              <h3>{t('liability1', 'Processor Liability')}</h3>
              <p>{t('liability1Desc', 'The Processor shall be liable for damages caused by processing where it has not complied with GDPR obligations specifically directed to processors or acted outside or contrary to lawful instructions.')}</p>
            </div>
            <div className={styles.liabilityCard}>
              <h3>{t('liability2', 'Indemnification')}</h3>
              <p>{t('liability2Desc', 'Each party shall indemnify the other against all losses arising from breach of this DPA or applicable data protection laws.')}</p>
            </div>
          </div>
        </section>

        {/* Contact Information */}
        <section className={styles.contactSection}>
          <h2>{t('contactTitle', 'Data Protection Contact')}</h2>
          <div className={styles.contactGrid}>
            <div className={styles.contactCard}>
              <Mail className={styles.contactIcon} />
              <h3>{t('dpoTitle', 'Data Protection Officer')}</h3>
              <p>redacted-email@example.com</p>
              <p>{t('responseTime', 'Response within 24 hours')}</p>
            </div>
            <div className={styles.contactCard}>
              <Phone className={styles.contactIcon} />
              <h3>{t('emergencyTitle', 'Emergency Hotline')}</h3>
              <p>+1 (555) 123-4567</p>
              <p>{t('available247', 'Available 24/7 for data incidents')}</p>
            </div>
            <div className={styles.contactCard}>
              <Shield className={styles.contactIcon} />
              <h3>{t('securityTitle', 'Security Team')}</h3>
              <p>redacted-email@example.com</p>
              <p>{t('securityResponse', 'For security concerns')}</p>
            </div>
          </div>
        </section>

        {/* Agreement Actions */}
        <section className={styles.actions}>
          <button className={styles.downloadFullButton}>
            <Download size={20} />
            {t('downloadFull', 'Download Full Agreement')}
          </button>
          <button 
            onClick={() => router.push('/contact/legal')}
            className={styles.contactButton}
          >
            <Mail size={20} />
            {t('contactLegal', 'Contact Legal Team')}
          </button>
          <button className={styles.printButton}>
            <FileText size={20} />
            {t('print', 'Print Agreement')}
          </button>
        </section>
      </main>
    </div>
  );
}