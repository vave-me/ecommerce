"use client";

import React, { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useMutation } from '@tanstack/react-query';
import { 
  ArrowLeft,
  Save,
  AlertTriangle,
  Info,
  Eye,
  EyeOff,
  HelpCircle
} from 'lucide-react';
import { useAuth } from '@/context/AuthContext';
import { useUserRole } from '@/hooks/useUserRole';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import { addConnector, getConnectorTypes } from '@/api/client/admin/erpApi';
import styles from './AddConnectorForm.module.css';

// Auth config field component
const AuthConfigField = ({ 
  label, 
  name, 
  value, 
  onChange, 
  type = 'text', 
  required = false, 
  placeholder = '',
  helpText = ''
}) => {
  const [showPassword, setShowPassword] = useState(false);
  const isPassword = type === 'password';
  
  return (
    <div className={styles.formField}>
      <label htmlFor={name}>
        {label}
        {required && <span className={styles.required}>*</span>}
        {helpText && (
          <span className={styles.helpIcon} title={helpText}>
            <HelpCircle size={14} />
          </span>
        )}
      </label>
      <div className={styles.inputWrapper}>
        <input
          id={name}
          name={name}
          type={isPassword && !showPassword ? 'password' : 'text'}
          value={value}
          onChange={(e) => onChange(name, e.target.value)}
          placeholder={placeholder}
          required={required}
          className={styles.formInput}
        />
        {isPassword && (
          <button
            type="button"
            className={styles.togglePassword}
            onClick={() => setShowPassword(!showPassword)}
          >
            {showPassword ? <EyeOff size={16} /> : <Eye size={16} />}
          </button>
        )}
      </div>
    </div>
  );
};

// Auth config forms for different ERP types
const AuthConfigForms = {
  netsuite: ({ config, onChange }) => (
    <>
      <AuthConfigField
        label="Account ID"
        name="account_id"
        value={config.account_id || ''}
        onChange={onChange}
        required
        placeholder="e.g., 123456"
        helpText="Your NetSuite account ID"
      />
      <AuthConfigField
        label="Consumer Key"
        name="consumer_key"
        value={config.consumer_key || ''}
        onChange={onChange}
        required
        placeholder="OAuth 1.0a consumer key"
      />
      <AuthConfigField
        label="Consumer Secret"
        name="consumer_secret"
        value={config.consumer_secret || ''}
        onChange={onChange}
        type="password"
        required
        placeholder="OAuth 1.0a consumer secret"
      />
      <AuthConfigField
        label="Token ID"
        name="token_id"
        value={config.token_id || ''}
        onChange={onChange}
        required
        placeholder="Access token ID"
      />
      <AuthConfigField
        label="Token Secret"
        name="token_secret"
        value={config.token_secret || ''}
        onChange={onChange}
        type="password"
        required
        placeholder="Access token secret"
      />
    </>
  ),
  
  odoo: ({ config, onChange }) => (
    <>
      <AuthConfigField
        label="Database"
        name="database"
        value={config.database || ''}
        onChange={onChange}
        required
        placeholder="e.g., mycompany"
        helpText="Odoo database name"
      />
      <AuthConfigField
        label="Username"
        name="username"
        value={config.username || ''}
        onChange={onChange}
        required
        placeholder="redacted-email@example.com"
      />
      <AuthConfigField
        label="Password"
        name="password"
        value={config.password || ''}
        onChange={onChange}
        type="password"
        required
        placeholder="Your Odoo password"
      />
      <AuthConfigField
        label="API Key (Optional)"
        name="api_key"
        value={config.api_key || ''}
        onChange={onChange}
        type="password"
        placeholder="API key for enhanced security"
        helpText="Use API key instead of password if available"
      />
    </>
  ),
  
  dynamics365: ({ config, onChange }) => (
    <>
      <AuthConfigField
        label="Tenant ID"
        name="tenant_id"
        value={config.tenant_id || ''}
        onChange={onChange}
        required
        placeholder="e.g., 12345678-1234-1234-1234-123456789012"
        helpText="Azure AD tenant ID"
      />
      <AuthConfigField
        label="Client ID"
        name="client_id"
        value={config.client_id || ''}
        onChange={onChange}
        required
        placeholder="Application (client) ID"
      />
      <AuthConfigField
        label="Client Secret"
        name="client_secret"
        value={config.client_secret || ''}
        onChange={onChange}
        type="password"
        required
        placeholder="Client secret value"
      />
      <AuthConfigField
        label="Resource URL"
        name="resource"
        value={config.resource || ''}
        onChange={onChange}
        required
        placeholder="https://org.crm.dynamics.com"
        helpText="Your Dynamics 365 instance URL"
      />
    </>
  ),
  
  sap: ({ config, onChange }) => (
    <>
      <AuthConfigField
        label="Client"
        name="client"
        value={config.client || ''}
        onChange={onChange}
        required
        placeholder="e.g., 100"
        helpText="SAP client number"
      />
      <AuthConfigField
        label="Username"
        name="username"
        value={config.username || ''}
        onChange={onChange}
        required
        placeholder="SAP username"
      />
      <AuthConfigField
        label="Password"
        name="password"
        value={config.password || ''}
        onChange={onChange}
        type="password"
        required
        placeholder="SAP password"
      />
      <AuthConfigField
        label="Language"
        name="language"
        value={config.language || 'EN'}
        onChange={onChange}
        placeholder="e.g., EN, DE"
        helpText="SAP system language"
      />
    </>
  ),
  
  // Default form for other types
  default: ({ config, onChange }) => (
    <>
      <AuthConfigField
        label="API Key"
        name="api_key"
        value={config.api_key || ''}
        onChange={onChange}
        required
        type="password"
        placeholder="Your API key"
      />
      <AuthConfigField
        label="API Secret (Optional)"
        name="api_secret"
        value={config.api_secret || ''}
        onChange={onChange}
        type="password"
        placeholder="API secret if required"
      />
    </>
  )
};

const AddConnectorForm = () => {
  const t = useTranslations('ERPConnectors');
  const router = useRouter();
  const { user } = useAuth();
  const { isAdmin } = useUserRole();

  const [formData, setFormData] = useState({
    name: '',
    type: '',
    environment: 'production',
    baseUrl: '',
    authType: 'basic',
    authConfig: {},
    webhookEnabled: false,
    syncEnabled: true,
    syncIntervalSeconds: 300,
    batchSize: 100,
    rateLimitRequestsPerSecond: 10,
    rateLimitBurst: 20
  });

  const [errors, setErrors] = useState({});

  // Add connector mutation
  const addMutation = useMutation({
    mutationFn: async (data) => {
      return await addConnector(data);
    },
    onSuccess: (response) => {
      
      router.push('/admin/erp/connectors');
    },
    onError: (error) => {
      // Error: 'Failed to add connector:', error...
      setErrors({
        submit: error.response?.data?.message || error.message || 'Failed to add connector'
      });
    }
  });

  const handleChange = (field, value) => {
    setFormData(prev => ({
      ...prev,
      [field]: value
    }));
    
    // Clear error for this field
    if (errors[field]) {
      setErrors(prev => ({
        ...prev,
        [field]: undefined
      }));
    }
  };

  const handleAuthConfigChange = (field, value) => {
    setFormData(prev => ({
      ...prev,
      authConfig: {
        ...prev.authConfig,
        [field]: value
      }
    }));
  };

  const handleTypeChange = (type) => {
    // Reset auth config when type changes
    setFormData(prev => ({
      ...prev,
      type,
      authConfig: {},
      authType: getDefaultAuthType(type)
    }));
  };

  const getDefaultAuthType = (type) => {
    switch (type) {
      case 'netsuite':
        return 'oauth1';
      case 'dynamics365':
        return 'oauth2';
      case 'odoo':
      case 'sap':
        return 'basic';
      default:
        return 'apikey';
    }
  };

  const getBaseUrlPlaceholder = (type) => {
    switch (type) {
      case 'netsuite':
        return 'https://123456.suitetalk.api.netsuite.com';
      case 'odoo':
        return 'https://mycompany.odoo.com';
      case 'dynamics365':
        return 'https://myorg.api.crm.dynamics.com';
      case 'sap':
        return 'https://mysap.company.com:8000';
      default:
        return 'https://api.example.com';
    }
  };

  const validateForm = () => {
    const newErrors = {};

    if (!formData.name.trim()) {
      newErrors.name = 'Connector name is required';
    }

    if (!formData.type) {
      newErrors.type = 'Connector type is required';
    }

    if (!formData.baseUrl.trim()) {
      newErrors.baseUrl = 'Base URL is required';
    } else {
      try {
        new URL(formData.baseUrl);
      } catch (e) {
        newErrors.baseUrl = 'Invalid URL format';
      }
    }

    // Validate auth config based on type
    const requiredAuthFields = getRequiredAuthFields(formData.type);
    requiredAuthFields.forEach(field => {
      if (!formData.authConfig[field]) {
        newErrors[`auth_${field}`] = `${field.replace(/_/g, ' ')} is required`;
      }
    });

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const getRequiredAuthFields = (type) => {
    switch (type) {
      case 'netsuite':
        return ['account_id', 'consumer_key', 'consumer_secret', 'token_id', 'token_secret'];
      case 'odoo':
        return ['database', 'username', 'password'];
      case 'dynamics365':
        return ['tenant_id', 'client_id', 'client_secret', 'resource'];
      case 'sap':
        return ['client', 'username', 'password'];
      default:
        return ['api_key'];
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    
    if (!validateForm()) {
      return;
    }

    // Submit the form
    addMutation.mutate(formData);
  };

  if (!isAdmin) {
    return (
      <div className={styles.accessDenied}>
        <AlertTriangle size={48} />
        <h2>Access Denied</h2>
        <p>You don't have permission to add ERP connectors.</p>
      </div>
    );
  }

  const AuthForm = AuthConfigForms[formData.type] || AuthConfigForms.default;

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <button 
          className={styles.backButton}
          onClick={() => router.push('/admin/erp/connectors')}
        >
          <ArrowLeft size={20} />
          Back to Connectors
        </button>
      </div>

      <form onSubmit={handleSubmit} className={styles.form}>
        <div className={styles.section}>
          <h2>Basic Information</h2>
          
          <div className={styles.formGroup}>
            <label htmlFor="name">
              Connector Name <span className={styles.required}>*</span>
            </label>
            <input
              id="name"
              type="text"
              value={formData.name}
              onChange={(e) => handleChange('name', e.target.value)}
              placeholder="e.g., Production SAP System"
              className={errors.name ? styles.errorInput : ''}
              required
            />
            {errors.name && <span className={styles.errorText}>{errors.name}</span>}
          </div>

          <div className={styles.formRow}>
            <div className={styles.formGroup}>
              <label htmlFor="type">
                ERP Type <span className={styles.required}>*</span>
              </label>
              <select
                id="type"
                value={formData.type}
                onChange={(e) => handleTypeChange(e.target.value)}
                className={errors.type ? styles.errorInput : ''}
                required
              >
                <option value="">Select ERP type...</option>
                <option value="sap">SAP ERP</option>
                <option value="netsuite">NetSuite</option>
                <option value="odoo">Odoo</option>
                <option value="dynamics365">Microsoft Dynamics 365</option>
                <option value="erpnext">ERPNext</option>
                <option value="frappe">Frappe</option>
              </select>
              {errors.type && <span className={styles.errorText}>{errors.type}</span>}
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="environment">
                Environment
              </label>
              <select
                id="environment"
                value={formData.environment}
                onChange={(e) => handleChange('environment', e.target.value)}
              >
                <option value="production">Production</option>
                <option value="staging">Staging</option>
                <option value="development">Development</option>
                <option value="sandbox">Sandbox</option>
              </select>
            </div>
          </div>

          <div className={styles.formGroup}>
            <label htmlFor="baseUrl">
              Base URL <span className={styles.required}>*</span>
            </label>
            <input
              id="baseUrl"
              type="url"
              value={formData.baseUrl}
              onChange={(e) => handleChange('baseUrl', e.target.value)}
              placeholder={getBaseUrlPlaceholder(formData.type)}
              className={errors.baseUrl ? styles.errorInput : ''}
              required
            />
            {errors.baseUrl && <span className={styles.errorText}>{errors.baseUrl}</span>}
          </div>
        </div>

        {formData.type && (
          <div className={styles.section}>
            <h2>Authentication Configuration</h2>
            <div className={styles.infoBox}>
              <Info size={16} />
              <span>Configure authentication for {formData.type.toUpperCase()}</span>
            </div>
            
            <AuthForm 
              config={formData.authConfig} 
              onChange={handleAuthConfigChange} 
            />
          </div>
        )}

        <div className={styles.section}>
          <h2>Sync Configuration</h2>
          
          <div className={styles.formGroup}>
            <label className={styles.checkboxLabel}>
              <input
                type="checkbox"
                checked={formData.syncEnabled}
                onChange={(e) => handleChange('syncEnabled', e.target.checked)}
              />
              <span>Enable automatic synchronization</span>
            </label>
          </div>

          {formData.syncEnabled && (
            <>
              <div className={styles.formRow}>
                <div className={styles.formGroup}>
                  <label htmlFor="syncIntervalSeconds">
                    Sync Interval (seconds)
                  </label>
                  <input
                    id="syncIntervalSeconds"
                    type="number"
                    min="60"
                    max="86400"
                    value={formData.syncIntervalSeconds}
                    onChange={(e) => handleChange('syncIntervalSeconds', parseInt(e.target.value) || 300)}
                  />
                  <span className={styles.fieldHelp}>How often to sync data (60-86400 seconds)</span>
                </div>

                <div className={styles.formGroup}>
                  <label htmlFor="batchSize">
                    Batch Size
                  </label>
                  <input
                    id="batchSize"
                    type="number"
                    min="1"
                    max="1000"
                    value={formData.batchSize}
                    onChange={(e) => handleChange('batchSize', parseInt(e.target.value) || 100)}
                  />
                  <span className={styles.fieldHelp}>Records per sync batch (1-1000)</span>
                </div>
              </div>
            </>
          )}

          <div className={styles.formGroup}>
            <label className={styles.checkboxLabel}>
              <input
                type="checkbox"
                checked={formData.webhookEnabled}
                onChange={(e) => handleChange('webhookEnabled', e.target.checked)}
              />
              <span>Enable webhooks for real-time updates</span>
            </label>
          </div>
        </div>

        <div className={styles.section}>
          <h2>Rate Limiting</h2>
          
          <div className={styles.formRow}>
            <div className={styles.formGroup}>
              <label htmlFor="rateLimitRequestsPerSecond">
                Requests per Second
              </label>
              <input
                id="rateLimitRequestsPerSecond"
                type="number"
                min="1"
                max="100"
                value={formData.rateLimitRequestsPerSecond}
                onChange={(e) => handleChange('rateLimitRequestsPerSecond', parseInt(e.target.value) || 10)}
              />
              <span className={styles.fieldHelp}>Max requests per second (1-100)</span>
            </div>

            <div className={styles.formGroup}>
              <label htmlFor="rateLimitBurst">
                Burst Size
              </label>
              <input
                id="rateLimitBurst"
                type="number"
                min="1"
                max="200"
                value={formData.rateLimitBurst}
                onChange={(e) => handleChange('rateLimitBurst', parseInt(e.target.value) || 20)}
              />
              <span className={styles.fieldHelp}>Max burst requests (1-200)</span>
            </div>
          </div>
        </div>

        {errors.submit && (
          <div className={styles.errorBox}>
            <AlertTriangle size={20} />
            <span>{errors.submit}</span>
          </div>
        )}

        <div className={styles.formActions}>
          <button
            type="button"
            className={styles.cancelButton}
            onClick={() => router.push('/admin/erp/connectors')}
          >
            Cancel
          </button>
          
          <button
            type="submit"
            className={styles.submitButton}
            disabled={addMutation.isPending}
          >
            {addMutation.isPending ? (
              <>
                <LoadingSpinner size={16} />
                <span>Creating...</span>
              </>
            ) : (
              <>
                <Save size={16} />
                <span>Create Connector</span>
              </>
            )}
          </button>
        </div>
      </form>
    </div>
  );
};

export default AddConnectorForm;