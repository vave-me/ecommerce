"use client";
import { memo, useMemo } from 'react';
import { cn } from '../../utils/cn';
/**
 * Advanced skeleton loading component system
 * Features:
 * - Multiple animation types
 * - Responsive design
 * - Accessibility support
 * - Customizable shapes and sizes
 * - Performance optimized
 */
// Base skeleton component
const SkeletonBase = memo(function SkeletonBase({
  className,
  animation = 'pulse',
  speed = 'normal',
  children,
  ...props
}) {
  const animationClasses = useMemo(() => {
    const speedMap = {
      slow: 'duration-2000',
      normal: 'duration-1500',
      fast: 'duration-1000',
    };
    const animationMap = {
      pulse: `animate-pulse ${speedMap[speed]}`,
      wave: `animate-shimmer ${speedMap[speed]}`,
      none: '',
    };
    return animationMap[animation] || animationMap.pulse;
  }, [animation, speed]);
  return (
    <div
      className={cn(
        'bg-gray-200 dark:bg-gray-700 rounded',
        animationClasses,
        className
      )}
      role="status"
      aria-label="Loading content"
      aria-live="polite"
      {...props}
    >
      {children}
      <span className="sr-only">Loading...</span>
    </div>
  );
});
// Text skeleton
export const SkeletonText = memo(function SkeletonText({
  lines = 1,
  className,
  lastLineWidth = '75%',
  ...props
}) {
  const textLines = useMemo(() => {
    return Array.from({ length: lines }, (_, index) => (
      <SkeletonBase
        key={index}
        className={cn(
          'h-4 mb-2 last:mb-0',
          index === lines - 1 && lines > 1 ? `w-[${lastLineWidth}]` : 'w-full'
        )}
        {...props}
      />
    ));
  }, [lines, lastLineWidth, props]);
  if (lines === 1) {
    return (
      <SkeletonBase
        className={cn('h-4 w-full', className)}
        {...props}
      />
    );
  }
  return (
    <div className={cn('space-y-2', className)}>
      {textLines}
    </div>
  );
});
// Avatar skeleton
export const SkeletonAvatar = memo(function SkeletonAvatar({
  size = 'md',
  className,
  ...props
}) {
  const sizeClasses = useMemo(() => {
    const sizeMap = {
      xs: 'w-6 h-6',
      sm: 'w-8 h-8',
      md: 'w-10 h-10',
      lg: 'w-12 h-12',
      xl: 'w-16 h-16',
      '2xl': 'w-20 h-20',
    };
    return sizeMap[size] || sizeMap.md;
  }, [size]);
  return (
    <SkeletonBase
      className={cn(
        'rounded-full flex-shrink-0',
        sizeClasses,
        className
      )}
      {...props}
    />
  );
});
// Card skeleton
export const SkeletonCard = memo(function SkeletonCard({
  hasImage = true,
  hasAvatar = false,
  textLines = 3,
  className,
  imageHeight = 'h-48',
  ...props
}) {
  return (
    <div className={cn('border border-gray-200 dark:border-gray-700 rounded-lg p-4 space-y-4', className)}>
      {hasImage && (
        <SkeletonBase
          className={cn('w-full rounded-lg', imageHeight)}
          {...props}
        />
      )}
      <div className="space-y-3">
        {hasAvatar && (
          <div className="flex items-center space-x-3">
            <SkeletonAvatar size="md" {...props} />
            <div className="flex-1">
              <SkeletonText lines={1} className="w-1/3" {...props} />
              <SkeletonText lines={1} className="w-1/4 mt-1" {...props} />
            </div>
          </div>
        )}
        <SkeletonText lines={textLines} {...props} />
        <div className="flex space-x-2">
          <SkeletonBase className="h-8 w-20 rounded-md" {...props} />
          <SkeletonBase className="h-8 w-16 rounded-md" {...props} />
        </div>
      </div>
    </div>
  );
});
// List item skeleton
export const SkeletonListItem = memo(function SkeletonListItem({
  hasAvatar = true,
  hasImage = false,
  textLines = 2,
  className,
  ...props
}) {
  return (
    <div className={cn('flex items-start space-x-3 p-3', className)}>
      {hasAvatar && <SkeletonAvatar size="md" {...props} />}
      <div className="flex-1 space-y-2">
        <SkeletonText lines={textLines} {...props} />
        {hasImage && (
          <SkeletonBase className="w-full h-32 rounded-lg" {...props} />
        )}
      </div>
    </div>
  );
});
// Table skeleton
export const SkeletonTable = memo(function SkeletonTable({
  rows = 5,
  columns = 4,
  hasHeader = true,
  className,
  ...props
}) {
  const tableRows = useMemo(() => {
    return Array.from({ length: rows }, (_, rowIndex) => (
      <tr key={rowIndex} className="border-b border-gray-200 dark:border-gray-700">
        {Array.from({ length: columns }, (_, colIndex) => (
          <td key={colIndex} className="p-3">
            <SkeletonBase
              className={cn(
                'h-4',
                colIndex === 0 ? 'w-3/4' : 'w-full'
              )}
              {...props}
            />
          </td>
        ))}
      </tr>
    ));
  }, [rows, columns, props]);
  return (
    <div className={cn('w-full', className)}>
      <table className="w-full">
        {hasHeader && (
          <thead>
            <tr className="border-b-2 border-gray-200 dark:border-gray-700">
              {Array.from({ length: columns }, (_, index) => (
                <th key={index} className="p-3 text-left">
                  <SkeletonBase className="h-4 w-2/3" {...props} />
                </th>
              ))}
            </tr>
          </thead>
        )}
        <tbody>
          {tableRows}
        </tbody>
      </table>
    </div>
  );
});
// Form skeleton
export const SkeletonForm = memo(function SkeletonForm({
  fields = 4,
  hasSubmitButton = true,
  className,
  ...props
}) {
  const formFields = useMemo(() => {
    return Array.from({ length: fields }, (_, index) => (
      <div key={index} className="space-y-2">
        <SkeletonBase className="h-4 w-1/4" {...props} />
        <SkeletonBase className="h-10 w-full rounded-md" {...props} />
      </div>
    ));
  }, [fields, props]);
  return (
    <div className={cn('space-y-4', className)}>
      {formFields}
      {hasSubmitButton && (
        <SkeletonBase className="h-10 w-32 rounded-md" {...props} />
      )}
    </div>
  );
});
// Feed skeleton (for social media feeds)
export const SkeletonFeed = memo(function SkeletonFeed({
  posts = 3,
  className,
  ...props
}) {
  const feedPosts = useMemo(() => {
    return Array.from({ length: posts }, (_, index) => (
      <div key={index} className="border-b border-gray-200 dark:border-gray-700 pb-6 mb-6 last:border-b-0 last:mb-0">
        <div className="flex items-center space-x-3 mb-4">
          <SkeletonAvatar size="md" {...props} />
          <div className="flex-1">
            <SkeletonText lines={1} className="w-1/3" {...props} />
            <SkeletonText lines={1} className="w-1/4 mt-1" {...props} />
          </div>
        </div>
        <SkeletonText lines={3} className="mb-4" {...props} />
        <SkeletonBase className="w-full h-64 rounded-lg mb-4" {...props} />
        <div className="flex items-center space-x-6">
          <SkeletonBase className="h-8 w-16 rounded-md" {...props} />
          <SkeletonBase className="h-8 w-16 rounded-md" {...props} />
          <SkeletonBase className="h-8 w-16 rounded-md" {...props} />
        </div>
      </div>
    ));
  }, [posts, props]);
  return (
    <div className={cn('space-y-6', className)}>
      {feedPosts}
    </div>
  );
});
// Product grid skeleton
export const SkeletonProductGrid = memo(function SkeletonProductGrid({
  products = 8,
  columns = 4,
  className,
  ...props
}) {
  const productItems = useMemo(() => {
    return Array.from({ length: products }, (_, index) => (
      <div key={index} className="space-y-3">
        <SkeletonBase className="w-full aspect-square rounded-lg" {...props} />
        <SkeletonText lines={2} {...props} />
        <SkeletonBase className="h-6 w-1/3 rounded-md" {...props} />
        <SkeletonBase className="h-8 w-full rounded-md" {...props} />
      </div>
    ));
  }, [products, props]);
  const gridClasses = useMemo(() => {
    const columnMap = {
      1: 'grid-cols-1',
      2: 'grid-cols-1 sm:grid-cols-2',
      3: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3',
      4: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4',
      5: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5',
      6: 'grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6',
    };
    return columnMap[columns] || columnMap[4];
  }, [columns]);
  return (
    <div className={cn('grid gap-6', gridClasses, className)}>
      {productItems}
    </div>
  );
});
// Main skeleton component with variants
const AdvancedSkeleton = memo(function AdvancedSkeleton({
  variant = 'text',
  ...props
}) {
  const variants = {
    text: SkeletonText,
    avatar: SkeletonAvatar,
    card: SkeletonCard,
    listItem: SkeletonListItem,
    table: SkeletonTable,
    form: SkeletonForm,
    feed: SkeletonFeed,
    productGrid: SkeletonProductGrid,
  };
  const SkeletonComponent = variants[variant] || SkeletonText;
  return <SkeletonComponent {...props} />;
});
AdvancedSkeleton.displayName = 'AdvancedSkeleton';
export default AdvancedSkeleton;
// Export individual components
export {
  SkeletonBase,
  SkeletonText,
  SkeletonAvatar,
  SkeletonCard,
  SkeletonListItem,
  SkeletonTable,
  SkeletonForm,
  SkeletonFeed,
  SkeletonProductGrid,
}; 