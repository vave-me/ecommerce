# ShareComponent Documentation

A comprehensive, reusable sharing component that provides multiple sharing methods across the application.

## Features

- **Web Share API**: Native sharing on supported devices
- **Social Media Sharing**: Facebook, Twitter, LinkedIn integration
- **Copy to Clipboard**: Fallback sharing method
- **Email Sharing**: Direct email composition
- **Analytics Tracking**: Built-in share event tracking
- **Multiple Variants**: Button, dropdown, modal, and inline layouts
- **Authentication Support**: Optional login requirements
- **Responsive Design**: Mobile-first approach
- **Accessibility**: Full ARIA support and keyboard navigation

## Usage Examples

### 1. Basic Button (Default)

```jsx
import ShareComponent from '@/components/Shared/ShareComponent';

<ShareComponent
    url="https://example.com/product/123"
    title="Amazing Product"
    description="Check out this amazing product!"
    contentId="123"
    contentType="product"
/>
```

### 2. Integration with Engagement Bar

```jsx
// Already integrated in src/components/shared/Engagement.jsx
<Engagement
    // ... other props
    shareUrl={productUrl}
    shareTitle={productTitle}
    shareDescription={productDescription}
    contentId={productId}
    contentType="product"
    requireAuth={false}
/>
```

### 3. Dropdown Variant

```jsx
<ShareComponent
    variant="dropdown"
    url={window.location.href}
    title="Share this page"
    platforms={['native', 'copy', 'facebook', 'twitter']}
    onShareSuccess={(platform, data) => {
        console.log('Shared via:', platform);
    }}
/>
```

### 4. Modal Variant

```jsx
<ShareComponent
    variant="modal"
    size="large"
    url={contentUrl}
    title={contentTitle}
    description={contentDescription}
    requireAuth={true}
    onShareError={(platform, error) => {
        if (platform === 'unauthorized') {
            // Show login modal
        }
    }}
/>
```

### 5. Inline Sharing Buttons

```jsx
<ShareComponent
    variant="inline"
    platforms={['facebook', 'twitter', 'linkedin', 'email']}
    url={shareUrl}
    title={shareTitle}
    className="my-custom-share-row"
/>
```

### 6. Icon Only Button

```jsx
<ShareComponent
    iconOnly={true}
    size="small"
    showCount={true}
    count={shareCount}
    className="engagement-share-button"
/>
```

## Props Reference

### Content Props
- `url` (string): URL to share (defaults to current page)
- `title` (string): Title for sharing (defaults to page title)
- `description` (string): Description for sharing
- `image` (string): Image URL for social media sharing

### Configuration Props
- `variant` ('button' | 'dropdown' | 'modal' | 'inline'): Display variant
- `size` ('small' | 'medium' | 'large'): Component size
- `platforms` (array): Available sharing platforms
- `showCount` (boolean): Show share count badge
- `count` (number): Current share count
- `requireAuth` (boolean): Require authentication to share

### Styling Props
- `className` (string): Additional CSS classes
- `buttonText` (string): Custom button text
- `iconOnly` (boolean): Show only icon, no text

### Event Props
- `onShare` (function): Called on any share action
- `onShareSuccess` (function): Called on successful share
- `onShareError` (function): Called on share error

### Analytics Props
- `trackingData` (object): Additional analytics data
- `contentId` (string): Unique content identifier
- `contentType` (string): Type of content being shared

## Platform Support

### Available Platforms
- `native`: Web Share API (mobile/PWA)
- `copy`: Copy link to clipboard
- `facebook`: Facebook sharing
- `twitter`: Twitter sharing
- `linkedin`: LinkedIn sharing
- `email`: Email composition

### Platform Detection
The component automatically detects available platforms:
- Web Share API availability
- Clipboard API support
- Browser compatibility

## Analytics & Tracking

### Automatic Tracking
- Google Analytics integration (if gtag available)
- Custom event tracking via callbacks
- Share count updates
- Platform-specific metrics

### Custom Analytics
```jsx
<ShareComponent
    trackingData={{
        category: 'product',
        label: 'premium-item',
        userId: currentUser.id
    }}
    onShareSuccess={(platform, data) => {
        // Custom analytics
        analytics.track('content_shared', {
            platform,
            contentId: data.contentId,
            userId: data.userId
        });
    }}
/>
```

## Context Provider

### ShareProvider Setup
```jsx
import { ShareProvider } from '@/components/Shared/ShareComponent/ShareProvider';

<ShareProvider defaultConfig={{
    defaultPlatforms: ['native', 'copy', 'facebook', 'twitter'],
    requireAuthForSharing: false,
    enableAnalytics: true,
    trackingEndpoint: '/api/shares'
}}>
    <App />
</ShareProvider>
```

### Using Share Context
```jsx
import { useShareContext } from '@/components/Shared/ShareComponent/ShareProvider';

function MyComponent() {
    const { 
        getShareCount, 
        recordGlobalShare,
        analytics 
    } = useShareContext();
    
    const shareCount = getShareCount('product-123');
    
    return (
        <div>
            <p>Shared {shareCount} times</p>
            <ShareComponent contentId="product-123" />
        </div>
    );
}
```

## Custom Hook

### useShareAPI Hook
```jsx
import { useShareAPI } from '@/hooks/useShareAPI';

function CustomShareButton() {
    const { 
        handleShare, 
        isNativeShareSupported,
        loading,
        error 
    } = useShareAPI();
    
    const handleClick = async () => {
        await handleShare({
            url: window.location.href,
            title: document.title,
            contentId: 'custom-content',
            contentType: 'article',
            platform: 'native'
        });
    };
    
    return (
        <button onClick={handleClick} disabled={loading}>
            {isNativeShareSupported() ? 'Share' : 'Copy Link'}
        </button>
    );
}
```

## Styling & Theming

### CSS Variables
The component uses CSS variables for theming:
```css
:root {
    --bg-primary: #ffffff;
    --bg-secondary: #f8fafc;
    --text-primary: #111827;
    --text-secondary: #6b7280;
    --border-color: #e1e5e9;
    --primary: #3b82f6;
    --success-bg: #dcfce7;
    --success-text: #166534;
}
```

### Dark Mode Support
Automatic dark mode detection via `prefers-color-scheme`:
```css
@media (prefers-color-scheme: dark) {
    /* Dark theme variables */
}
```

### Custom Styling
```jsx
<ShareComponent
    className="my-custom-share"
    style={{
        '--bg-primary': '#1f2937',
        '--text-primary': '#f9fafb'
    }}
/>
```

## Accessibility

### Features
- Full ARIA label support
- Keyboard navigation
- Screen reader compatibility
- High contrast mode support
- Reduced motion support

### ARIA Labels
All buttons include proper ARIA labels:
- Share button: "Share this content"
- Platform buttons: "Share on Facebook", etc.
- Copy button: "Copy link to clipboard"

## Browser Support

### Requirements
- Modern browsers (ES2018+)
- Optional: Web Share API (mobile)
- Optional: Clipboard API (copy functionality)

### Fallbacks
- No Web Share API → Copy to clipboard
- No Clipboard API → Manual URL display
- Graceful degradation for all features

## Performance

### Optimizations
- Lazy loading via dynamic imports
- Memoized components and callbacks
- Efficient event handling
- Minimal re-renders

### Bundle Size
- Core component: ~15KB gzipped
- Dynamic loading reduces initial bundle
- Tree-shakeable platform modules 