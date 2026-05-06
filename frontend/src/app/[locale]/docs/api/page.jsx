"use client";

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import { 
  Code,
  Terminal,
  Book,
  Zap,
  Lock,
  Globe,
  Database,
  Server,
  Shield,
  Key,
  FileJson,
  GitBranch,
  Package,
  ShoppingCart,
  Users,
  MessageSquare,
  BarChart3,
  Settings,
  ChevronRight,
  ChevronDown,
  Copy,
  Check,
  ArrowRight,
  ExternalLink,
  Download,
  PlayCircle,
  Cpu,
  Cloud,
  Layers,
  Webhook
} from 'lucide-react';
import styles from './ApiDocs.module.css';

export default function ApiDocsPage() {
  const t = useTranslations('apiDocs');
  const router = useRouter();
  const [activeSection, setActiveSection] = useState('overview');
  const [expandedEndpoints, setExpandedEndpoints] = useState({});
  const [copiedCode, setCopiedCode] = useState(null);

  const handleCopyCode = (id, code) => {
    navigator.clipboard.writeText(code);
    setCopiedCode(id);
    setTimeout(() => setCopiedCode(null), 2000);
  };

  const toggleEndpoint = (id) => {
    setExpandedEndpoints(prev => ({
      ...prev,
      [id]: !prev[id]
    }));
  };

  const sections = [
    { id: 'overview', icon: Book, label: t('overview', 'Overview') },
    { id: 'authentication', icon: Key, label: t('authentication', 'Authentication') },
    { id: 'products', icon: Package, label: t('products', 'Products') },
    { id: 'orders', icon: ShoppingCart, label: t('orders', 'Orders') },
    { id: 'users', icon: Users, label: t('users', 'Users') },
    { id: 'analytics', icon: BarChart3, label: t('analytics', 'Analytics') },
    { id: 'webhooks', icon: Webhook, label: t('webhooks', 'Webhooks') },
    { id: 'sdks', icon: Layers, label: t('sdks', 'SDKs & Libraries') }
  ];

  const endpoints = {
    products: [
      {
        id: 'list-products',
        method: 'GET',
        path: '/api/v1/products',
        description: t('listProductsDesc', 'Retrieve a list of products'),
        parameters: [
          { name: 'page', type: 'integer', description: t('pageParam', 'Page number (default: 1)') },
          { name: 'limit', type: 'integer', description: t('limitParam', 'Items per page (default: 20, max: 100)') },
          { name: 'sort', type: 'string', description: t('sortParam', 'Sort by: created_at, price, name') },
          { name: 'category', type: 'string', description: t('categoryParam', 'Filter by category ID') }
        ],
        example: `curl -X GET https://api.platform.com/v1/products \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json"`,
        response: `{
  "data": [
    {
      "id": "prod_123",
      "name": "Premium Widget",
      "price": 2999,
      "currency": "EUR",
      "stock": 150,
      "created_at": "2024-01-15T10:30:00Z"
    }
  ],
  "meta": {
    "page": 1,
    "limit": 20,
    "total": 245
  }
}`
      },
      {
        id: 'create-product',
        method: 'POST',
        path: '/api/v1/products',
        description: t('createProductDesc', 'Create a new product'),
        parameters: [
          { name: 'name', type: 'string', required: true, description: t('nameParam', 'Product name') },
          { name: 'description', type: 'string', required: true, description: t('descParam', 'Product description') },
          { name: 'price', type: 'integer', required: true, description: t('priceParam', 'Price in cents') },
          { name: 'stock', type: 'integer', description: t('stockParam', 'Initial stock quantity') },
          { name: 'images', type: 'array', description: t('imagesParam', 'Array of image URLs') }
        ],
        example: `curl -X POST https://api.platform.com/v1/products \\
  -H "Authorization: Bearer YOUR_API_KEY" \\
  -H "Content-Type: application/json" \\
  -d '{
    "name": "Premium Widget",
    "description": "High-quality widget for professionals",
    "price": 2999,
    "stock": 150,
    "images": ["https://example.com/image1.jpg"]
  }'`,
        response: `{
  "data": {
    "id": "prod_123",
    "name": "Premium Widget",
    "description": "High-quality widget for professionals",
    "price": 2999,
    "currency": "EUR",
    "stock": 150,
    "images": ["https://example.com/image1.jpg"],
    "created_at": "2024-01-15T10:30:00Z"
  }
}`
      }
    ],
    orders: [
      {
        id: 'list-orders',
        method: 'GET',
        path: '/api/v1/orders',
        description: t('listOrdersDesc', 'Retrieve a list of orders'),
        parameters: [
          { name: 'status', type: 'string', description: t('statusParam', 'Filter by status: pending, processing, shipped, delivered') },
          { name: 'customer', type: 'string', description: t('customerParam', 'Filter by customer ID') },
          { name: 'date_from', type: 'string', description: t('dateFromParam', 'Filter orders from date (ISO 8601)') },
          { name: 'date_to', type: 'string', description: t('dateToParam', 'Filter orders to date (ISO 8601)') }
        ]
      },
      {
        id: 'update-order',
        method: 'PATCH',
        path: '/api/v1/orders/{id}',
        description: t('updateOrderDesc', 'Update order status'),
        parameters: [
          { name: 'status', type: 'string', required: true, description: t('statusUpdateParam', 'New order status') },
          { name: 'tracking_number', type: 'string', description: t('trackingParam', 'Shipping tracking number') },
          { name: 'notes', type: 'string', description: t('notesParam', 'Internal notes') }
        ]
      }
    ]
  };

  const rateLimits = [
    { endpoint: t('allEndpoints', 'All endpoints'), limit: '1000 requests/hour' },
    { endpoint: t('createEndpoints', 'Create endpoints'), limit: '100 requests/hour' },
    { endpoint: t('analyticsEndpoints', 'Analytics endpoints'), limit: '500 requests/hour' }
  ];

  const sdks = [
    { 
      language: 'JavaScript/TypeScript', 
      icon: FileJson,
      package: '@platform/sdk',
      install: 'npm install @platform/sdk'
    },
    { 
      language: 'Python', 
      icon: Code,
      package: 'platform-sdk',
      install: 'pip install platform-sdk'
    },
    { 
      language: 'PHP', 
      icon: Code,
      package: 'platform/sdk',
      install: 'composer require platform/sdk'
    },
    { 
      language: 'Ruby', 
      icon: Code,
      package: 'platform-sdk',
      install: 'gem install platform-sdk'
    }
  ];

  return (
    <div className={styles.container}>
      {/* Header */}
      <header className={styles.header}>
        <div className={styles.headerContent}>
          <div className={styles.headerInfo}>
            <h1 className={styles.title}>
              {t('pageTitle', 'API Documentation')}
            </h1>
            <p className={styles.subtitle}>
              {t('subtitle', 'Build powerful integrations with our RESTful API')}
            </p>
            
            <div className={styles.headerActions}>
              <button className={styles.primaryButton}>
                <Terminal size={20} />
                {t('apiPlayground', 'API Playground')}
              </button>
              <button className={styles.secondaryButton}>
                <Download size={20} />
                {t('downloadPostman', 'Postman Collection')}
              </button>
            </div>
          </div>

          <div className={styles.headerStats}>
            <div className={styles.statCard}>
              <Zap className={styles.statIcon} />
              <div>
                <div className={styles.statValue}>99.99%</div>
                <div className={styles.statLabel}>{t('uptime', 'API Uptime')}</div>
              </div>
            </div>
            <div className={styles.statCard}>
              <Server className={styles.statIcon} />
              <div>
                <div className={styles.statValue}>&lt;50ms</div>
                <div className={styles.statLabel}>{t('responseTime', 'Avg Response')}</div>
              </div>
            </div>
            <div className={styles.statCard}>
              <Globe className={styles.statIcon} />
              <div>
                <div className={styles.statValue}>v1</div>
                <div className={styles.statLabel}>{t('currentVersion', 'Current Version')}</div>
              </div>
            </div>
          </div>
        </div>
      </header>

      {/* Quick Start */}
      <section className={styles.quickStart}>
        <div className={styles.quickStartContainer}>
          <h2>{t('quickStart', 'Quick Start')}</h2>
          <div className={styles.quickStartSteps}>
            <div className={styles.step}>
              <div className={styles.stepNumber}>1</div>
              <h3>{t('step1', 'Get your API key')}</h3>
              <p>{t('step1Desc', 'Generate an API key from your dashboard')}</p>
            </div>
            <div className={styles.step}>
              <div className={styles.stepNumber}>2</div>
              <h3>{t('step2', 'Make your first request')}</h3>
              <p>{t('step2Desc', 'Use curl or your favorite HTTP client')}</p>
            </div>
            <div className={styles.step}>
              <div className={styles.stepNumber}>3</div>
              <h3>{t('step3', 'Explore endpoints')}</h3>
              <p>{t('step3Desc', 'Check out our comprehensive API reference')}</p>
            </div>
          </div>
        </div>
      </section>

      {/* Main Content */}
      <main className={styles.mainContent}>
        <div className={styles.contentGrid}>
          {/* Sidebar Navigation */}
          <aside className={styles.sidebar}>
            <nav className={styles.navigation}>
              {sections.map(section => (
                <button
                  key={section.id}
                  onClick={() => setActiveSection(section.id)}
                  className={`${styles.navItem} ${activeSection === section.id ? styles.active : ''}`}
                >
                  <section.icon size={18} />
                  <span>{section.label}</span>
                </button>
              ))}
            </nav>

            <div className={styles.resources}>
              <h3>{t('resources', 'Resources')}</h3>
              <a href="#" className={styles.resourceLink}>
                <PlayCircle size={16} />
                {t('videoTutorials', 'Video Tutorials')}
              </a>
              <a href="#" className={styles.resourceLink}>
                <Book size={16} />
                {t('apiGuide', 'API Best Practices')}
              </a>
              <a href="#" className={styles.resourceLink}>
                <GitBranch size={16} />
                {t('changelog', 'API Changelog')}
              </a>
              <a href="#" className={styles.resourceLink}>
                <ExternalLink size={16} />
                {t('status', 'API Status')}
              </a>
            </div>
          </aside>

          {/* Documentation Content */}
          <div className={styles.docContent}>
            {/* Overview Section */}
            {activeSection === 'overview' && (
              <section className={styles.section}>
                <h2>{t('apiOverview', 'API Overview')}</h2>
                <p className={styles.intro}>
                  {t('overviewIntro', 'Our API provides programmatic access to read and write platform data. Build custom integrations, automate workflows, and create powerful applications.')}
                </p>

                <div className={styles.features}>
                  <div className={styles.featureCard}>
                    <Shield className={styles.featureIcon} />
                    <h3>{t('secure', 'Secure by Design')}</h3>
                    <p>{t('secureDesc', 'OAuth 2.0 authentication, encrypted connections, and granular permissions')}</p>
                  </div>
                  <div className={styles.featureCard}>
                    <Cpu className={styles.featureIcon} />
                    <h3>{t('performant', 'High Performance')}</h3>
                    <p>{t('performantDesc', 'Optimized for speed with response times under 50ms')}</p>
                  </div>
                  <div className={styles.featureCard}>
                    <Cloud className={styles.featureIcon} />
                    <h3>{t('scalable', 'Built to Scale')}</h3>
                    <p>{t('scalableDesc', 'Handle millions of requests with our global infrastructure')}</p>
                  </div>
                </div>

                <div className={styles.baseUrl}>
                  <h3>{t('baseUrl', 'Base URL')}</h3>
                  <div className={styles.codeBlock}>
                    <code>https://api.platform.com/v1</code>
                    <button 
                      onClick={() => handleCopyCode('base-url', 'https://api.platform.com/v1')}
                      className={styles.copyButton}
                    >
                      {copiedCode === 'base-url' ? <Check size={16} /> : <Copy size={16} />}
                    </button>
                  </div>
                </div>

                <div className={styles.rateLimits}>
                  <h3>{t('rateLimits', 'Rate Limits')}</h3>
                  <table className={styles.table}>
                    <thead>
                      <tr>
                        <th>{t('endpoint', 'Endpoint')}</th>
                        <th>{t('limit', 'Limit')}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {rateLimits.map((limit, index) => (
                        <tr key={index}>
                          <td>{limit.endpoint}</td>
                          <td>{limit.limit}</td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              </section>
            )}

            {/* Authentication Section */}
            {activeSection === 'authentication' && (
              <section className={styles.section}>
                <h2>{t('authenticationTitle', 'Authentication')}</h2>
                <p className={styles.intro}>
                  {t('authIntro', 'All API requests require authentication. We use API keys that you can generate from your dashboard.')}
                </p>

                <div className={styles.authExample}>
                  <h3>{t('usingApiKey', 'Using your API key')}</h3>
                  <p>{t('apiKeyDesc', 'Include your API key in the Authorization header:')}</p>
                  
                  <div className={styles.codeExample}>
                    <div className={styles.codeHeader}>
                      <span>curl</span>
                      <button 
                        onClick={() => handleCopyCode('auth-example', 'Authorization: Bearer YOUR_API_KEY')}
                        className={styles.copyButton}
                      >
                        {copiedCode === 'auth-example' ? <Check size={16} /> : <Copy size={16} />}
                      </button>
                    </div>
                    <pre className={styles.code}>
{`curl https://api.platform.com/v1/products \\
  -H "Authorization: Bearer YOUR_API_KEY"`}
                    </pre>
                  </div>
                </div>

                <div className={styles.securityNote}>
                  <Lock size={20} />
                  <div>
                    <h4>{t('securityNote', 'Security Note')}</h4>
                    <p>{t('securityNoteDesc', 'Keep your API keys secure. Never expose them in client-side code or public repositories.')}</p>
                  </div>
                </div>
              </section>
            )}

            {/* Products Section */}
            {activeSection === 'products' && (
              <section className={styles.section}>
                <h2>{t('productsApi', 'Products API')}</h2>
                <p className={styles.intro}>
                  {t('productsIntro', 'Manage your product catalog programmatically. Create, update, and sync products across channels.')}
                </p>

                <div className={styles.endpoints}>
                  {endpoints.products.map(endpoint => (
                    <div key={endpoint.id} className={styles.endpoint}>
                      <button
                        onClick={() => toggleEndpoint(endpoint.id)}
                        className={styles.endpointHeader}
                      >
                        <div className={styles.endpointInfo}>
                          <span className={`${styles.method} ${styles[endpoint.method.toLowerCase()]}`}>
                            {endpoint.method}
                          </span>
                          <code className={styles.path}>{endpoint.path}</code>
                        </div>
                        <ChevronDown 
                          className={`${styles.toggleIcon} ${expandedEndpoints[endpoint.id] ? styles.expanded : ''}`}
                        />
                      </button>

                      {expandedEndpoints[endpoint.id] && (
                        <div className={styles.endpointContent}>
                          <p className={styles.endpointDesc}>{endpoint.description}</p>
                          
                          {endpoint.parameters && (
                            <div className={styles.parameters}>
                              <h4>{t('parameters', 'Parameters')}</h4>
                              <table className={styles.paramTable}>
                                <thead>
                                  <tr>
                                    <th>{t('name', 'Name')}</th>
                                    <th>{t('type', 'Type')}</th>
                                    <th>{t('description', 'Description')}</th>
                                  </tr>
                                </thead>
                                <tbody>
                                  {endpoint.parameters.map((param, index) => (
                                    <tr key={index}>
                                      <td>
                                        <code>{param.name}</code>
                                        {param.required && <span className={styles.required}>*</span>}
                                      </td>
                                      <td>{param.type}</td>
                                      <td>{param.description}</td>
                                    </tr>
                                  ))}
                                </tbody>
                              </table>
                            </div>
                          )}

                          {endpoint.example && (
                            <div className={styles.example}>
                              <h4>{t('example', 'Example Request')}</h4>
                              <div className={styles.codeExample}>
                                <div className={styles.codeHeader}>
                                  <span>curl</span>
                                  <button 
                                    onClick={() => handleCopyCode(`${endpoint.id}-example`, endpoint.example)}
                                    className={styles.copyButton}
                                  >
                                    {copiedCode === `${endpoint.id}-example` ? <Check size={16} /> : <Copy size={16} />}
                                  </button>
                                </div>
                                <pre className={styles.code}>{endpoint.example}</pre>
                              </div>
                            </div>
                          )}

                          {endpoint.response && (
                            <div className={styles.response}>
                              <h4>{t('response', 'Example Response')}</h4>
                              <div className={styles.codeExample}>
                                <div className={styles.codeHeader}>
                                  <span>JSON</span>
                                  <button 
                                    onClick={() => handleCopyCode(`${endpoint.id}-response`, endpoint.response)}
                                    className={styles.copyButton}
                                  >
                                    {copiedCode === `${endpoint.id}-response` ? <Check size={16} /> : <Copy size={16} />}
                                  </button>
                                </div>
                                <pre className={styles.code}>{endpoint.response}</pre>
                              </div>
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  ))}
                </div>
              </section>
            )}

            {/* SDKs Section */}
            {activeSection === 'sdks' && (
              <section className={styles.section}>
                <h2>{t('sdksTitle', 'SDKs & Libraries')}</h2>
                <p className={styles.intro}>
                  {t('sdksIntro', 'Use our official SDKs to integrate faster with your favorite programming language.')}
                </p>

                <div className={styles.sdkGrid}>
                  {sdks.map((sdk, index) => (
                    <div key={index} className={styles.sdkCard}>
                      <sdk.icon className={styles.sdkIcon} size={32} />
                      <h3>{sdk.language}</h3>
                      <code className={styles.packageName}>{sdk.package}</code>
                      
                      <div className={styles.installCommand}>
                        <pre>{sdk.install}</pre>
                        <button 
                          onClick={() => handleCopyCode(`sdk-${index}`, sdk.install)}
                          className={styles.copyButton}
                        >
                          {copiedCode === `sdk-${index}` ? <Check size={16} /> : <Copy size={16} />}
                        </button>
                      </div>
                      
                      <button className={styles.sdkLink}>
                        {t('viewDocs', 'View Documentation')}
                        <ExternalLink size={16} />
                      </button>
                    </div>
                  ))}
                </div>
              </section>
            )}
          </div>
        </div>
      </main>

      {/* CTA Section */}
      <section className={styles.ctaSection}>
        <div className={styles.ctaContent}>
          <h2>{t('ctaTitle', 'Ready to Build?')}</h2>
          <p>{t('ctaDesc', 'Get your API key and start building amazing integrations today')}</p>
          <div className={styles.ctaButtons}>
            <button 
              onClick={() => router.push('/dashboard/api-keys')}
              className={styles.ctaPrimary}
            >
              <Key size={20} />
              {t('getApiKey', 'Get API Key')}
            </button>
            <button 
              onClick={() => router.push('/support/developers')}
              className={styles.ctaSecondary}
            >
              <MessageSquare size={20} />
              {t('contactSupport', 'Developer Support')}
            </button>
          </div>
        </div>
      </section>
    </div>
  );
}