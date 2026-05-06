# Classified Components

This directory contains the unified classified/product card components that follow the established patterns from DealCard and AutomotiveCard.

## Architecture

The classified components follow a **server/client separation pattern** for optimal performance:

- **`ClassifiedCard.jsx`** - Client component handling interactivity and state management
- **`ClassifiedCard.server.jsx`** - Server component for media pre-fetching and optimization
- **`ClassifiedCard.module.css`** - Unified styling with responsive design

## Components

### ClassifiedCard (Client Component)

The main client-side component that handles:

- **Interactive Features**: Like, dislike, comment, message, favorite
- **Media Navigation**: Image carousel with navigation controls
- **Real-time Updates**: Dynamic status indicators and engagement metrics
- **Responsive Design**: Mobile-first approach with breakpoints
- **Accessibility**: ARIA labels, keyboard navigation, screen reader support

#### Props

```javascript
{
  product: Object,           // Product data object
  preloadedMedia: Array,     // Pre-fetched media from server (optional)
  hasPreloadedMedia: Boolean // Flag indicating if media was pre-loaded
}
```

#### Product Data Structure

```javascript
{
  // Core product data
  id: string,
  name: string,
  description: string,
  basePrice: string,
  condition: string,
  thumbnail: string,
  
  // Product-specific data
  brand: string,
  model: string,
  sku: string,
  stock: number,
  hasVariants: boolean,
  shippingCost: string,
  weight: string,
  dimensions: { height, width, depth },
  
  // Business data
  negotiable: boolean,
  middlemanService: boolean,
  userType: string,
  status: string,
  
  // Location and engagement
  location: string,
  tags: Array,
  metrics: Object,
  
  // User data
  userId: string,
  authorName: string,
  avatar: string
}
```

### ClassifiedCardServer (Server Component)

Server-side wrapper that:

- **Pre-fetches Media**: Optimizes loading by fetching images server-side
- **Error Handling**: Graceful fallback to client-side loading on failures
- **Performance**: Reduces client-side API calls and improves perceived performance
- **SEO Optimization**: Better server-side rendering for search engines

#### Usage

```javascript
// Server component usage
import { ClassifiedCardServer } from '@/components/classified';

export default function ProductsPage({ products }) {
  return (
    <div>
      {products.map(product => (
        <ClassifiedCardServer 
          key={product.id} 
          product={product} 
        />
      ))}
    </div>
  );
}

// Client component usage
import { ClassifiedCard } from '@/components/classified';

export default function ProductsList({ products }) {
  return (
    <div>
      {products.map(product => (
        <ClassifiedCard 
          key={product.id} 
          product={product} 
        />
      ))}
    </div>
  );
}
```

## Features

### 🎨 **Modern Design**
- Card-based layout with hover effects
- Gradient backgrounds and shadows
- Smooth animations and transitions
- Status badges and indicators

### 📱 **Responsive Design**
- Mobile-first approach
- Breakpoints: 768px, 480px
- Optimized touch targets
- Adaptive image sizes

### 🖼️ **Media Handling**
- Server-side media pre-fetching
- Image carousel with navigation
- Lazy loading and optimization
- Fallback to default images
- Error handling and recovery

### 💬 **Interactive Features**
- Like/dislike functionality
- Comment system integration
- Direct messaging
- Favorite/bookmark system
- Social sharing capabilities

### 🏷️ **Product Information**
- Price display with currency formatting
- Shipping cost indicators
- Stock availability
- Condition badges
- Brand and model information
- SKU and specifications
- Tag system

### 🔍 **Status Indicators**
- NEW badge for recent listings
- HOT badge for popular items
- Popularity indicators
- Middleman service badges
- Stock availability
- Negotiable pricing indicators

### ♿ **Accessibility**
- ARIA labels and descriptions
- Keyboard navigation support
- Screen reader compatibility
- High contrast support
- Focus management

## Styling

The component uses CSS Modules with:

- **Modern CSS**: Flexbox, Grid, CSS Variables
- **Animations**: Smooth transitions, hover effects, loading states
- **Responsive**: Mobile-first breakpoints
- **Theming**: Consistent color palette and typography
- **Performance**: Optimized selectors and minimal reflows

### CSS Classes

Key CSS classes available for customization:

```css
.card              /* Main card container */
.imageContainer    /* Product image area */
.productTitle      /* Product name/title */
.price            /* Price display */
.badgeRow         /* Status badges container */
.tagsContainer    /* Product tags */
.commentsWrapper  /* Comments section */
```

## Integration

### Feed System Integration

The ClassifiedCard integrates with the unified feed system:

```javascript
// In ConnectedFeedDisplay.jsx
case 'product':
  return (
    <ClassifiedCard
      key={item.id}
      product={item.product}
      preloadedMedia={item.preloadedMedia}
      hasPreloadedMedia={!!item.preloadedMedia}
    />
  );
```

### API Integration

The component integrates with:

- **Media API**: `getMediaByItem()` for image fetching
- **Activity API**: `useActivityApi()` for engagement actions
- **Comments API**: `CommentsSetup` component
- **Messaging API**: Modal system for direct messages

## Performance Optimizations

1. **Server-side Media Pre-fetching**: Reduces client-side API calls
2. **Memoized Calculations**: Price formatting, time calculations
3. **Lazy Loading**: Images and heavy components
4. **Optimized Re-renders**: React.memo and useCallback
5. **Error Boundaries**: Graceful error handling
6. **Bundle Optimization**: Tree-shaking friendly exports

## Testing

The components support testing with:

- **Unit Tests**: Component behavior and props
- **Integration Tests**: API interactions
- **Visual Tests**: Responsive design and styling
- **Accessibility Tests**: ARIA compliance and keyboard navigation

## Migration from ImprovedClassifiedCard

When migrating from the old `ImprovedClassifiedCard`:

1. **Replace imports**:
   ```javascript
   // Old
   import ImprovedClassifiedCard from './ImprovedClassifiedCard';
   
   // New
   import { ClassifiedCard } from '@/components/classified';
   ```

2. **Update props structure**: The new component expects a more structured product object

3. **Remove custom styling**: The new component has built-in responsive design

4. **Update media handling**: Media is now handled automatically with server pre-fetching

## Future Enhancements

Planned improvements:

- **Variant Support**: Product variants and options
- **Advanced Filtering**: Category and attribute filters
- **Comparison Mode**: Side-by-side product comparison
- **Wishlist Integration**: Enhanced favorite functionality
- **Review System**: Product ratings and reviews
- **Internationalization**: Multi-language support

## Dependencies

- React 18+
- Next.js 15+
- Lucide React (icons)
- dayjs (date formatting)
- React Redux (state management)
- React Toastify (notifications) 