/**
 * Secure DOM Utilities
 * Prevents XSS vulnerabilities in dynamic DOM creation
 * 
 * SECURITY: These utilities create DOM elements safely without using innerHTML
 * to prevent script injection attacks.
 */
/**
 * Securely create a notification element with safe DOM manipulation
 * @param {Object} options - Notification options
 * @param {string} options.type - Notification type ('update', 'offline', 'install')
 * @param {string} options.title - Notification title
 * @param {string} options.message - Notification message
 * @param {Function} options.onAction - Action button callback
 * @param {Function} options.onDismiss - Dismiss button callback
 * @returns {HTMLElement} - Safe notification element
 */
export function createSecureNotification({ type, title, message, onAction, onDismiss }) {
  // Create container
  const notification = document.createElement('div');
  notification.className = 'pwa-notification';
  // Set safe ID based on type
  const safeId = `pwa-${type}-notification`;
  notification.id = safeId;
  // Create main container with styles
  const container = document.createElement('div');
  container.style.cssText = getNotificationStyles(type);
  // Create content container
  const contentDiv = document.createElement('div');
  contentDiv.style.cssText = 'flex: 1; margin-right: 12px;';
  // Create title element
  const titleDiv = document.createElement('div');
  titleDiv.style.cssText = 'font-weight: 600; margin-bottom: 4px;';
  titleDiv.textContent = title; // Safe text content
  // Create message element
  const messageDiv = document.createElement('div');
  messageDiv.style.cssText = 'opacity: 0.9;';
  messageDiv.textContent = message; // Safe text content
  // Assemble content
  contentDiv.appendChild(titleDiv);
  contentDiv.appendChild(messageDiv);
  container.appendChild(contentDiv);
  // Create action button if provided
  if (onAction) {
    const actionBtn = document.createElement('button');
    actionBtn.style.cssText = getActionButtonStyles(type);
    actionBtn.textContent = type === 'install' ? 'Install' : 'Update';
    actionBtn.addEventListener('click', onAction);
    container.appendChild(actionBtn);
  }
  // Create dismiss button
  const dismissBtn = document.createElement('button');
  dismissBtn.style.cssText = getDismissButtonStyles();
  dismissBtn.textContent = '×';
  dismissBtn.addEventListener('click', onDismiss || (() => notification.remove()));
  container.appendChild(dismissBtn);
  // Add animation styles
  const style = document.createElement('style');
  style.textContent = `
    @keyframes slideUp {
      from { transform: translateY(100%); opacity: 0; }
      to { transform: translateY(0); opacity: 1; }
    }
  `;
  document.head.appendChild(style);
  notification.appendChild(container);
  return notification;
}
/**
 * Get notification styles based on type
 * @param {string} type - Notification type
 * @returns {string} - CSS styles string
 */
function getNotificationStyles(type) {
  const baseStyles = `
    position: fixed;
    bottom: 20px;
    left: 20px;
    right: 20px;
    color: white;
    padding: 16px;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
    z-index: 10000;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    font-size: 14px;
    animation: slideUp 0.3s ease-out;
  `;
  const backgrounds = {
    update: 'background: #2980b9;',
    offline: 'background: #10b981;',
    install: 'background: linear-gradient(135deg, #2980b9, #7c3aed);'
  };
  return baseStyles + (backgrounds[type] || backgrounds.update);
}
/**
 * Get action button styles
 * @param {string} type - Notification type
 * @returns {string} - CSS styles string
 */
function getActionButtonStyles(type) {
  if (type === 'install') {
    return `
      background: white;
      color: #2980b9;
      border: none;
      padding: 12px 24px;
      border-radius: 8px;
      font-weight: 600;
      font-size: 14px;
      cursor: pointer;
      white-space: nowrap;
    `;
  }
  return `
    background: white;
    color: #2980b9;
    border: none;
    padding: 8px 16px;
    border-radius: 8px;
    font-weight: 600;
    font-size: 14px;
    cursor: pointer;
    white-space: nowrap;
  `;
}
/**
 * Get dismiss button styles
 * @returns {string} - CSS styles string
 */
function getDismissButtonStyles() {
  return `
    background: transparent;
    color: white;
    border: none;
    padding: 8px;
    margin-left: 8px;
    cursor: pointer;
    font-size: 18px;
    line-height: 1;
  `;
}
/**
 * Safely set text content with optional fallback
 * @param {HTMLElement} element - Target element
 * @param {string} text - Text to set
 * @param {string} fallback - Fallback text if original is invalid
 */
export function secureSetText(element, text, fallback = '') {
  if (!element || typeof element.textContent === 'undefined') {
    return;
  }
  // Validate text input
  const safeText = (typeof text === 'string' && text.length > 0) ? text : fallback;
  element.textContent = safeText;
}
/**
 * Create a secure install banner
 * @param {Function} onInstall - Install button callback
 * @param {Function} onDismiss - Dismiss callback
 * @returns {HTMLElement} - Safe install banner element
 */
export function createSecureInstallBanner(onInstall, onDismiss) {
  const banner = document.createElement('div');
  banner.id = 'pwa-install-banner';
  const container = document.createElement('div');
  container.style.cssText = `
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    background: linear-gradient(135deg, #2980b9, #7c3aed);
    color: white;
    padding: 16px 20px;
    z-index: 10000;
    display: flex;
    align-items: center;
    justify-content: space-between;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
    transform: translateY(-100%);
    transition: transform 0.3s ease-out;
  `;
  // Content
  const content = document.createElement('div');
  content.style.cssText = 'flex: 1; margin-right: 16px;';
  const title = document.createElement('div');
  title.style.cssText = 'font-weight: 600; font-size: 16px; margin-bottom: 4px;';
  title.textContent = 'Install sfx markt';
  const message = document.createElement('div');
  message.style.cssText = 'font-size: 14px; opacity: 0.9;';
  message.textContent = 'Get the full app experience';
  content.appendChild(title);
  content.appendChild(message);
  // Install button
  const installBtn = document.createElement('button');
  installBtn.style.cssText = `
    background: white;
    color: #2980b9;
    border: none;
    padding: 12px 24px;
    border-radius: 8px;
    font-weight: 600;
    font-size: 14px;
    cursor: pointer;
    margin-right: 12px;
  `;
  installBtn.textContent = 'Install';
  installBtn.addEventListener('click', onInstall);
  // Dismiss button
  const dismissBtn = document.createElement('button');
  dismissBtn.style.cssText = getDismissButtonStyles();
  dismissBtn.textContent = '×';
  dismissBtn.addEventListener('click', onDismiss);
  container.appendChild(content);
  container.appendChild(installBtn);
  container.appendChild(dismissBtn);
  banner.appendChild(container);
  // Animate in
  setTimeout(() => {
    container.style.transform = 'translateY(0)';
  }, 100);
  return banner;
}
/**
 * Remove all PWA notifications safely
 */
export function removeAllPWANotifications() {
  const selectors = [
    '#pwa-update-notification',
    '#pwa-offline-notification', 
    '#pwa-install-banner',
    '.pwa-notification'
  ];
  selectors.forEach(selector => {
    const elements = document.querySelectorAll(selector);
    elements.forEach(element => element.remove());
  });
} 