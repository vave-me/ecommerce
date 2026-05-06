# Production Deployment Guide

## Overview
This guide covers deploying the production-ready frontend application with all security, performance, and monitoring features enabled.

## Prerequisites

- Node.js 18+ and npm/yarn
- Access to deployment platform (Vercel, AWS, etc.)
- Configured environment variables
- Sentry account for error tracking
- Analytics accounts (Google Analytics, etc.)

## Pre-Deployment Checklist

### 1. Environment Configuration
```bash
# Copy environment template
cp .env.example .env.production

# Configure all required variables:
- NEXT_PUBLIC_APP_URL (production URL)
- NEXT_PUBLIC_API_URL (production API)
- NEXT_PUBLIC_SENTRY_DSN (Sentry project DSN)
- Authentication secrets
- Payment provider keys
- Analytics IDs
```

### 2. Security Verification
- [x] All console.log statements removed
- [x] XSS protection implemented (DOMPurify)
- [x] Secure token storage configured
- [x] Environment variables properly secured
- [x] CSP headers configured
- [x] HTTPS enforced

### 3. Performance Optimization
- [x] Code splitting implemented
- [x] Images optimized with next/image
- [x] React.memo on heavy components
- [x] Bundle size analyzed
- [x] Core Web Vitals tested

### 4. Error Handling
- [x] Global error handler configured
- [x] Error boundaries in place
- [x] Sentry integration tested
- [x] User-friendly error messages
- [x] Error tracking enabled

## Build Process

### 1. Install Dependencies
```bash
npm ci --production
```

### 2. Run Tests
```bash
npm test
npm run test:coverage
```

### 3. Build Application
```bash
# Set build-time variables
export NEXT_PUBLIC_BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
export NEXT_PUBLIC_DEPLOYMENT_ID=$(git rev-parse --short HEAD)

# Production build
npm run build

# Analyze bundle size
ANALYZE=true npm run build
```

### 4. Verify Build
```bash
# Check build output
ls -la .next/

# Test production build locally
npm run start
```

## Deployment Platforms

### Vercel (Recommended)
```bash
# Install Vercel CLI
npm i -g vercel

# Deploy
vercel --prod

# Set environment variables
vercel env add NEXT_PUBLIC_SENTRY_DSN production
```

### AWS Amplify
```bash
# Install Amplify CLI
npm i -g @aws-amplify/cli

# Initialize
amplify init

# Deploy
amplify publish
```

### Docker
```dockerfile
# Dockerfile
FROM node:18-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:18-alpine
WORKDIR /app
COPY --from=builder /app/.next ./.next
COPY --from=builder /app/public ./public
COPY --from=builder /app/package*.json ./
RUN npm ci --production
EXPOSE 3000
CMD ["npm", "start"]
```

## Post-Deployment

### 1. Verify Deployment
- [ ] Application loads correctly
- [ ] API connections work
- [ ] Authentication flows function
- [ ] Payment processing works
- [ ] Real-time features (NATS) connected

### 2. Monitor Performance
```javascript
// Check Core Web Vitals
- LCP < 2.5s
- FID < 100ms
- CLS < 0.1
- TTFB < 800ms
```

### 3. Configure Monitoring

#### Sentry
1. Verify events are being received
2. Set up alerts for critical errors
3. Configure release tracking
4. Enable session replay

#### Performance Monitoring
1. Set up custom dashboards
2. Configure alerting thresholds
3. Monitor API response times
4. Track user interactions

### 4. Set Up Alerts
```javascript
// Critical alerts to configure:
- Error rate > 1%
- Page load time > 3s
- API errors > 5/min
- Memory usage > 80%
- Failed transactions
```

## Security Headers

Add these headers in your deployment platform:

```nginx
# Security headers
Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline' 'unsafe-eval' https://www.google-analytics.com https://www.googletagmanager.com; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; img-src 'self' data: https:; connect-src 'self' https://api.your-domain.com wss://your-domain.com https://sentry.io;
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
Referrer-Policy: strict-origin-when-cross-origin
Permissions-Policy: camera=(), microphone=(), geolocation=()
```

## Rollback Plan

### Quick Rollback
```bash
# Vercel
vercel rollback

# Git-based
git revert HEAD
git push origin main

# Docker
docker pull your-app:previous-version
docker stop current-container
docker run your-app:previous-version
```

### Database Rollback
- Keep migration rollback scripts
- Test rollback procedures
- Have database backups ready

## Monitoring Dashboard Access

1. **Sentry**: https://sentry.io/organizations/your-org/
2. **Google Analytics**: https://analytics.google.com/
3. **Application Monitoring**: /admin/monitoring
4. **Server Logs**: Check your hosting provider

## Troubleshooting

### Common Issues

1. **White screen on load**
   - Check browser console for errors
   - Verify all environment variables
   - Check Sentry for error reports

2. **API connection failures**
   - Verify CORS configuration
   - Check API URL in environment
   - Test API endpoints directly

3. **Performance issues**
   - Check bundle size
   - Verify image optimization
   - Review Core Web Vitals
   - Check for memory leaks

4. **Authentication problems**
   - Verify auth secrets match
   - Check token expiration
   - Review secure storage

## Maintenance

### Regular Tasks
- [ ] Review error reports weekly
- [ ] Check performance metrics
- [ ] Update dependencies monthly
- [ ] Review security alerts
- [ ] Backup configurations

### Performance Reviews
- [ ] Analyze bundle size trends
- [ ] Review Core Web Vitals
- [ ] Check error rates
- [ ] Monitor API performance

## Support

For deployment issues:
1. Check error logs in Sentry
2. Review deployment platform logs
3. Contact platform support
4. Check GitHub issues

## Appendix

### Useful Commands
```bash
# Check production build size
du -sh .next/

# Test production locally
NODE_ENV=production npm start

# Generate security report
npm audit

# Update dependencies
npm update --save
```

### Environment Variables Reference
See `.env.example` for complete list of variables and their descriptions.

### Performance Budgets
- Initial JS: < 100KB
- Total JS: < 300KB
- Image sizes: < 200KB
- Time to Interactive: < 3s
- First Contentful Paint: < 1.5s