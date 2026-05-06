/**
 * DUPLICATE CODE ELIMINATION UTILITY
 * Consolidates repeated patterns found across CreateModals, components, and hooks
 * 
 * FIX 55: Eliminate all duplicate localStorage, debounce, cleanFilters, and validation logic
 */
// ================== CONSOLIDATED DEBOUNCE IMPLEMENTATION ==================
export const createUnifiedDebounce = (func, wait = 300, options = {}) => {
  let timeout;
  let lastArgs;
  let lastThis;
  const { leading = false, trailing = true, maxWait } = options;
  let maxTimeout;
  let lastCallTime;
  let lastInvokeTime = 0;
  function invokeFunc(time) {
    const args = lastArgs;
    const thisArg = lastThis;
    lastArgs = lastThis = undefined;
    lastInvokeTime = time;
    return func.apply(thisArg, args);
  }
  function startTimer(pendingFunc, wait) {
    return setTimeout(pendingFunc, wait);
  }
  function cancelTimer(id) {
    clearTimeout(id);
  }
  function leadingEdge(time) {
    lastInvokeTime = time;
    timeout = startTimer(timerExpired, wait);
    return leading ? invokeFunc(time) : undefined;
  }
  function timerExpired() {
    const time = Date.now();
    if (shouldInvoke(time)) {
      return trailingEdge(time);
    }
    timeout = startTimer(timerExpired, remainingWait(time));
  }
  function trailingEdge(time) {
    timeout = undefined;
    if (trailing && lastArgs) {
      return invokeFunc(time);
    }
    lastArgs = lastThis = undefined;
  }
  function shouldInvoke(time) {
    const timeSinceLastCall = time - lastCallTime;
    const timeSinceLastInvoke = time - lastInvokeTime;
    return (lastCallTime === undefined || 
            timeSinceLastCall >= wait ||
            timeSinceLastCall < 0 ||
            (maxWait !== undefined && timeSinceLastInvoke >= maxWait));
  }
  function remainingWait(time) {
    const timeSinceLastCall = time - lastCallTime;
    const timeSinceLastInvoke = time - lastInvokeTime;
    const timeWaiting = wait - timeSinceLastCall;
    return maxWait === undefined 
      ? timeWaiting 
      : Math.min(timeWaiting, maxWait - timeSinceLastInvoke);
  }
  function debounced(...args) {
    const time = Date.now();
    const isInvoking = shouldInvoke(time);
    lastArgs = args;
    lastThis = this;
    lastCallTime = time;
    if (isInvoking) {
      if (timeout === undefined) {
        return leadingEdge(lastCallTime);
      }
      if (maxWait) {
        timeout = startTimer(timerExpired, wait);
        return invokeFunc(lastCallTime);
      }
    }
    if (timeout === undefined) {
      timeout = startTimer(timerExpired, wait);
    }
    return undefined;
  }
  debounced.cancel = () => {
    if (timeout !== undefined) {
      cancelTimer(timeout);
    }
    if (maxTimeout !== undefined) {
      cancelTimer(maxTimeout);
    }
    lastInvokeTime = 0;
    lastArgs = lastCallTime = lastThis = timeout = maxTimeout = undefined;
  };
  debounced.flush = () => {
    return timeout === undefined ? undefined : trailingEdge(Date.now());
  };
  debounced.pending = () => timeout !== undefined;
  return debounced;
};
// ================== CONSOLIDATED FILTERS CLEANING ==================
export const createUnifiedFilterCleaner = () => {
  const numericFields = new Set([
    'minPrice', 'maxPrice', 'minYear', 'maxYear', 'minArea', 'maxArea', 
    'minStock', 'maxStock', 'price', 'year', 'mileage', 'rooms', 'bathrooms'
  ]);
  const arrayFields = new Set([
    'tags', 'categories', 'features', 'amenities', 'colors', 'sizes'
  ]);
  return {
    cleanFilters: (filters = {}) => {
      const cleaned = {};
      for (const [key, value] of Object.entries(filters)) {
        // Skip empty values
        if (value === null || value === undefined || value === '') continue;
        if (Array.isArray(value) && value.length === 0) continue;
        if (typeof value === 'string' && value.trim() === '') continue;
        // Handle numeric fields
        if (numericFields.has(key)) {
          const numValue = Number(value);
          if (!isNaN(numValue) && numValue >= 0) {
            cleaned[key] = numValue;
          }
        }
        // Handle array fields
        else if (arrayFields.has(key)) {
          if (Array.isArray(value)) {
            const filteredArray = value.filter(item => item !== null && item !== undefined && item !== '');
            if (filteredArray.length > 0) {
              cleaned[key] = filteredArray;
            }
          } else if (typeof value === 'string') {
            const arrayValue = value.split(',').map(item => item.trim()).filter(Boolean);
            if (arrayValue.length > 0) {
              cleaned[key] = arrayValue;
            }
          }
        }
        // Handle regular fields
        else {
          cleaned[key] = value;
        }
      }
      return cleaned;
    },
    cleanRequestBody: (requestBody = {}) => {
      const cleaned = {};
      for (const [key, value] of Object.entries(requestBody)) {
        if (value === '' || value === null || value === undefined) continue;
        if (Array.isArray(value) && value.length === 0) continue;
        if (typeof value === 'object' && value !== null && Object.keys(value).length === 0) continue;
        cleaned[key] = value;
      }
      return cleaned;
    }
  };
};
// ================== CONSOLIDATED LOCALSTORAGE DRAFT SYSTEM ==================
export const createUnifiedDraftStorage = () => {
  const DRAFT_PREFIX = 'unified_draft_';
  const MAX_DRAFTS_PER_USER = 10;
  const DRAFT_EXPIRY_DAYS = 7;
  const getDraftKey = (entityType, userId) => `${DRAFT_PREFIX}${entityType}_${userId}`;
  const isExpired = (timestamp) => {
    const expiryTime = DRAFT_EXPIRY_DAYS * 24 * 60 * 60 * 1000;
    return Date.now() - timestamp > expiryTime;
  };
  return {
    saveDraft: (entityType, userId, draft) => {
      try {
        const draftKey = getDraftKey(entityType, userId);
        const existingDrafts = JSON.parse(localStorage.getItem(draftKey) || '[]');
        // Remove expired drafts
        const validDrafts = existingDrafts.filter(d => !isExpired(d.timestamp));
        // Find existing draft with same ID
        const draftIndex = draft.id 
          ? validDrafts.findIndex(d => d.id === draft.id)
          : -1;
        const draftData = {
          ...draft,
          id: draft.id || `draft_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
          timestamp: Date.now(),
          entityType,
          lastModified: new Date().toISOString()
        };
        if (draftIndex >= 0) {
          validDrafts[draftIndex] = draftData;
        } else {
          validDrafts.unshift(draftData);
          // Keep only MAX_DRAFTS_PER_USER most recent
          if (validDrafts.length > MAX_DRAFTS_PER_USER) {
            validDrafts.splice(MAX_DRAFTS_PER_USER);
          }
        }
        localStorage.setItem(draftKey, JSON.stringify(validDrafts));
        return { success: true, draftId: draftData.id, timestamp: draftData.timestamp };
      } catch (error) {
        return { success: false, error: error.message };
      }
    },
    getDrafts: (entityType, userId) => {
      try {
        const draftKey = getDraftKey(entityType, userId);
        const drafts = JSON.parse(localStorage.getItem(draftKey) || '[]');
        // Filter out expired drafts
        const validDrafts = drafts.filter(d => !isExpired(d.timestamp));
        // Save cleaned drafts back
        if (validDrafts.length !== drafts.length) {
          localStorage.setItem(draftKey, JSON.stringify(validDrafts));
        }
        return validDrafts;
      } catch (error) {
        return [];
      }
    },
    removeDraft: (entityType, userId, draftId) => {
      try {
        const draftKey = getDraftKey(entityType, userId);
        const drafts = JSON.parse(localStorage.getItem(draftKey) || '[]');
        const filteredDrafts = drafts.filter(d => d.id !== draftId);
        localStorage.setItem(draftKey, JSON.stringify(filteredDrafts));
        return { success: true };
      } catch (error) {
        return { success: false, error: error.message };
      }
    },
    clearExpiredDrafts: (entityType, userId) => {
      try {
        const draftKey = getDraftKey(entityType, userId);
        const drafts = JSON.parse(localStorage.getItem(draftKey) || '[]');
        const validDrafts = drafts.filter(d => !isExpired(d.timestamp));
        localStorage.setItem(draftKey, JSON.stringify(validDrafts));
        return { success: true, removed: drafts.length - validDrafts.length };
      } catch (error) {
        return { success: false, error: error.message };
      }
    }
  };
};
// ================== CONSOLIDATED VALIDATION SYSTEM ==================
export const createUnifiedValidator = () => {
  const patterns = {
    email: /^[^\s@]+@[^\s@]+\.[^\s@]+$/,
    phone: /^[\+]?[1-9][\d]{0,15}$/,
    url: /^https?:\/\/.+\..+/,
    price: /^\d+(\.\d{1,2})?$/,
    zipcode: /^[0-9]{2}-[0-9]{3}$|^[0-9]{5}$/,
    alphanumeric: /^[a-zA-Z0-9]+$/,
    slug: /^[a-z0-9-]+$/
  };
  const validators = {
    required: (value, fieldName, message) => {
      if (!value || (typeof value === 'string' && value.trim() === '')) {
        return { [fieldName]: message || `${fieldName} is required` };
      }
      return {};
    },
    minLength: (value, fieldName, minLength, message) => {
      if (value && value.length < minLength) {
        return { [fieldName]: message || `${fieldName} must be at least ${minLength} characters` };
      }
      return {};
    },
    maxLength: (value, fieldName, maxLength, message) => {
      if (value && value.length > maxLength) {
        return { [fieldName]: message || `${fieldName} must be no more than ${maxLength} characters` };
      }
      return {};
    },
    pattern: (value, fieldName, patternName, message) => {
      if (value && !patterns[patternName]?.test(value)) {
        return { [fieldName]: message || `Invalid ${fieldName} format` };
      }
      return {};
    },
    range: (value, fieldName, min, max, message) => {
      const numValue = Number(value);
      if (!isNaN(numValue) && (numValue < min || numValue > max)) {
        return { [fieldName]: message || `${fieldName} must be between ${min} and ${max}` };
      }
      return {};
    },
    custom: (value, fieldName, validatorFn, message) => {
      if (!validatorFn(value)) {
        return { [fieldName]: message || `Invalid ${fieldName}` };
      }
      return {};
    }
  };
  return {
    validateField: (value, fieldName, rules = []) => {
      let errors = {};
      for (const rule of rules) {
        const ruleErrors = validators[rule.type]?.(value, fieldName, ...rule.params, rule.message);
        if (ruleErrors && Object.keys(ruleErrors).length > 0) {
          errors = { ...errors, ...ruleErrors };
        }
      }
      return errors;
    },
    validateForm: (formData, validationRules) => {
      let allErrors = {};
      for (const [fieldName, rules] of Object.entries(validationRules)) {
        const fieldErrors = this.validateField(formData[fieldName], fieldName, rules);
        allErrors = { ...allErrors, ...fieldErrors };
      }
      return allErrors;
    },
    patterns,
    validators
  };
};
// ================== CONSOLIDATED FILE PROCESSING ==================
export const createUnifiedFileProcessor = () => {
  const generateFileHash = async (file) => {
    try {
      const arrayBuffer = await file.arrayBuffer();
      const hashBuffer = await crypto.subtle.digest('SHA-256', arrayBuffer);
      const hashArray = Array.from(new Uint8Array(hashBuffer));
      return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
    } catch (error) {
      return `fallback_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`;
    }
  };
  return {
    processFiles: async (files, options = {}) => {
      const {
        maxFileSize = 10 * 1024 * 1024, // 10MB
        maxFiles = 10,
        allowedTypes = ['image/*', 'video/*', 'application/pdf', 'text/*'],
        generatePreviews = true,
        checkDuplicates = true,
        existingHashes = new Set()
      } = options;
      const processedFiles = [];
      const errors = [];
      if (files.length > maxFiles) {
        errors.push(`Maximum ${maxFiles} files allowed`);
        return { processedFiles: [], errors };
      }
      for (const file of files) {
        try {
          // Validate file type
          const isValidType = allowedTypes.some(type => {
            if (type.endsWith('/*')) {
              return file.type.startsWith(type.slice(0, -1));
            }
            return file.type === type;
          });
          if (!isValidType) {
            errors.push(`Unsupported file type: ${file.name}`);
            continue;
          }
          // Validate file size
          if (file.size > maxFileSize) {
            errors.push(`File too large: ${file.name} (max ${Math.round(maxFileSize / (1024 * 1024))}MB)`);
            continue;
          }
          // Generate hash for duplicate detection
          const hash = checkDuplicates ? await generateFileHash(file) : null;
          if (checkDuplicates && existingHashes.has(hash)) {
            errors.push(`Duplicate file: ${file.name}`);
            continue;
          }
          // Generate preview if needed
          let preview = null;
          if (generatePreviews && file.type.startsWith('image/')) {
            preview = URL.createObjectURL(file);
          }
          const processedFile = {
            file,
            name: file.name,
            type: file.type,
            size: file.size,
            hash,
            preview,
            uploadProgress: 0,
            id: `file_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
          };
          processedFiles.push(processedFile);
          if (checkDuplicates) existingHashes.add(hash);
        } catch (error) {
          errors.push(`Failed to process: ${file.name}`);
        }
      }
      return { processedFiles, errors };
    },
    generateFileHash,
    cleanupPreviews: (files) => {
      files.forEach(file => {
        if (file.preview && file.preview.startsWith('blob:')) {
          URL.revokeObjectURL(file.preview);
        }
      });
    }
  };
};
// ================== EXPORT UNIFIED UTILITIES ==================
export const UnifiedUtils = {
  debounce: createUnifiedDebounce,
  filters: createUnifiedFilterCleaner(),
  drafts: createUnifiedDraftStorage(),
  validator: createUnifiedValidator(),
  fileProcessor: createUnifiedFileProcessor()
};
export default UnifiedUtils; 