# AGENTS Guide: frontend

## Business Purpose
Frontend web client workspace and UI build/runtime assets.

## Key Files In This Directory
- `sfx_markt/frontend/.claude/settings.local.json`
- `sfx_markt/frontend/.env`
- `sfx_markt/frontend/.env.development`
- `sfx_markt/frontend/.env.example`
- `sfx_markt/frontend/.env.local`
- `sfx_markt/frontend/.env.production`
- `sfx_markt/frontend/.local`
- `sfx_markt/frontend/ACCESSIBILITY_IMPROVEMENTS.md`
- `sfx_markt/frontend/ATOMIC_FIX_PLAN.md`
- `sfx_markt/frontend/AUTONOMOUS_WORK_PLAN.md`
- `sfx_markt/frontend/DEPLOYMENT_GUIDE.md`
- `sfx_markt/frontend/Dockerfile`
- `sfx_markt/frontend/PRODUCTION_READINESS_ANALYSIS.md`
- `sfx_markt/frontend/PRODUCTION_READINESS_FINAL_REPORT.md`
- `sfx_markt/frontend/PRODUCTION_READINESS_PROGRESS.md`
- `sfx_markt/frontend/PRODUCTION_READINESS_SUMMARY.md`
- `sfx_markt/frontend/SECURE_TOKEN_MIGRATION.md`
- `sfx_markt/frontend/__mocks__/api/mockFeedData.js`
- `sfx_markt/frontend/__mocks__/api/searchApi.jsx`
- `sfx_markt/frontend/__tests__/__mocks__/authContext.js`
- `sfx_markt/frontend/__tests__/__mocks__/authContextMock.jsx`
- `sfx_markt/frontend/__tests__/__mocks__/fileMock.js`
- `sfx_markt/frontend/__tests__/__mocks__/nats.ws.js`
- `sfx_markt/frontend/__tests__/__mocks__/styleMock.js`
- `sfx_markt/frontend/__tests__/api/axiosInstance.test.js`
- `sfx_markt/frontend/__tests__/api/searchApi.test.jsx`
- `sfx_markt/frontend/__tests__/api/userApi.test.js`
- `sfx_markt/frontend/__tests__/app/client-layout.test.js`
- `sfx_markt/frontend/__tests__/app/homepage.test.js`
- `sfx_markt/frontend/__tests__/app/layouts.test.js`
- `sfx_markt/frontend/__tests__/app/providers-real-intl.test.js`
- `sfx_markt/frontend/__tests__/app/providers.test.js`
- `sfx_markt/frontend/__tests__/auth/auth.utils.test.js`
- `sfx_markt/frontend/__tests__/components/ErrorBoundary.test.js`
- `sfx_markt/frontend/__tests__/components/Feed/Feed.test.jsx`
- `sfx_markt/frontend/__tests__/components/Feed/FeedProvider.test.jsx`
- `sfx_markt/frontend/__tests__/components/Feed/LeftsideIntegration.test.jsx`
- `sfx_markt/frontend/__tests__/components/Header.test.js`
- `sfx_markt/frontend/__tests__/components/navigation.test.jsx`
- `sfx_markt/frontend/__tests__/config/next-config.test.js`
- `sfx_markt/frontend/__tests__/context/AuthContext.test.jsx`
- `sfx_markt/frontend/__tests__/context/CategoryContext.test.jsx`
- `sfx_markt/frontend/__tests__/context/NATSContext.test.jsx`
- `sfx_markt/frontend/__tests__/context/NavBarContext.test.jsx`
- `sfx_markt/frontend/__tests__/features/CreateDealModal.test.jsx`
- `sfx_markt/frontend/__tests__/features/CreatePostModal.test.jsx`
- `sfx_markt/frontend/__tests__/features/CreateVehicleModal.test.jsx`
- `sfx_markt/frontend/__tests__/feed/Feed.test.jsx`
- `sfx_markt/frontend/__tests__/feed/FeedProvider.test.jsx`
- `sfx_markt/frontend/__tests__/feed/LeftsideIntegration.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useActivityApi.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useAutoSave.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useBasketActions.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useBodyLock.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useCategoryContext.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useCategoryKeyboardNav.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useChatHistory.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useCommentsActions.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useConversation.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useEntityCache.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useFeedQuery.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useFocusTrap.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useFormValidation.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useHeaderScroll.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useInteractions.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useIsMobile.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useListItems.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useMedia.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useMediaQuery.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useMobileDetection.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/usePerformance.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/usePrefetchOnHover.test.jsx`
- `sfx_markt/frontend/__tests__/hooks/useWishlist.test.jsx`
- `sfx_markt/frontend/__tests__/i18n/comprehensive-i18n.test.js`
- `sfx_markt/frontend/__tests__/i18n/live-routing.test.js`
- `sfx_markt/frontend/__tests__/integration/add-dropdown-integration.test.js`
- `sfx_markt/frontend/__tests__/integration/client-layout-real.test.js`
- `sfx_markt/frontend/__tests__/integration/comprehensive-integration-summary.test.js`
- `sfx_markt/frontend/__tests__/integration/fixed-navigation-url-capture.test.js`
- `sfx_markt/frontend/__tests__/integration/header-navigation-integration.test.js`

## How To Work In This Directory
1. Keep changes aligned with the owning service/business module contracts.
2. Do not edit generated/build/vendor artifacts directly.
3. Validate with targeted build/test/deploy commands relevant to affected modules.

## Relationship To Commands / Queries / Proto / Events
- This directory is support/infrastructure/integration-oriented unless it contains a dedicated service module pattern.
- For adding/removing commands, queries, RPCs, domain events, and SQL read models, edit the owning service directory (for example `sfx_markt/users`, `sfx_markt/products`, `sfx_markt/ordering`, `sfx_markt/payments`, `sfx_markt/search`).
- Regeneration baseline when shared contracts change:
  - `cd <workspace-root>/sfx_markt && go generate ./...`

## SQL / Data Safety
- Prefer additive, forward-only schema/data changes.
- Never rewrite historical migrations or previously applied seed files.
