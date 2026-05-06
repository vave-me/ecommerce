package constants

// EnhancedNaturalLanguagePatterns provides comprehensive mapping between natural language and tool calls
const EnhancedNaturalLanguagePatterns = `
## NATURAL LANGUAGE TO TOOL MAPPING GUIDE

### 1. ACTIVITY & ANALYTICS
**User says** → **Use tool**
- "track this action" → activity_log
- "log user activity" → activity_log
- "show my activity history" → activity_get_user
- "what have I been doing" → activity_get_user
- "show system activities" → activity_get_system
- "what's happening on the platform" → activity_get_system

### 2. SHOPPING & BASKET
**User says** → **Use tool**
- "add to cart/basket" → basket_add_item
- "put this in my basket" → basket_add_item
- "remove from cart" → basket_remove_item
- "take this out of my basket" → basket_remove_item
- "show my cart" → basket_get + basket_get_items
- "what's in my basket" → basket_get_items
- "empty my cart" → basket_clear
- "clear everything" → basket_clear
- "how much is everything" → basket_get_total
- "calculate total" → basket_get_total

### 3. PRODUCT SEARCH & DISCOVERY
**User says** → **Use tool**
- "find products" → product_search_advanced
- "search for [item]" → product_search_advanced with search_text
- "show me electronics under $200" → product_search_advanced with category_id="electronics", max_price=20000
- "cheapest phones" → product_search_advanced with category_id="phones", sort_by="price", sort_order="asc"
- "newest arrivals" → product_search_advanced with sort_by="created_at", sort_order="desc"
- "products like this" → vector_search_similar with entity_type="product"
- "similar items" → vector_search_similar
- "recommendations for me" → vector_get_recommendations
- "suggest something" → product_get_recommendations

### 4. ORDER MANAGEMENT
**User says** → **Use tool**
- "place an order" → order_create
- "create order" → order_create
- "buy these items" → order_create with items array
- "show order #12345" → order_get with order_id="12345"
- "track my order" → order_get + shipping_track
- "my recent orders" → order_get_user
- "cancel my order" → order_update_status with status="cancelled"
- "mark as shipped" → order_update_status with status="shipped"
- "refund this order" → order_refund

### 5. PAYMENT PROCESSING
**User says** → **Use tool**
- "pay for order" → payment_create
- "process payment" → payment_process
- "charge my card" → payment_create with payment_method_id
- "refund payment" → payment_refund
- "payment history" → payment_get_user
- "add payment method" → payment_add_method
- "remove my card" → payment_remove_method
- "saved payment methods" → payment_get_methods

### 6. USER ACCOUNT
**User says** → **Use tool**
- "find user john" → user_find or user_search with query="john"
- "create account" → user_create
- "sign up" → user_create
- "update profile" → user_update
- "change my name" → user_update with firstName/lastName
- "login/sign in" → user_login
- "logout" → user_logout
- "reset password" → user_forgot_password
- "enable account" → user_enable

### 7. COMMUNICATION
**User says** → **Use tool**
- "send notification" → notification_create
- "notify user about" → notification_create
- "email them" → mailer_send_email
- "send message to" → message_send
- "unread messages" → message_get_unread_count
- "subscribe to newsletter" → newsletter_subscribe
- "unsubscribe me" → newsletter_unsubscribe

### 8. REVIEWS & FEEDBACK
**User says** → **Use tool**
- "leave a review" → review_add
- "rate this product" → review_add with rating
- "write feedback" → review_add or comment_add
- "comment on this" → comment_add
- "reply to review" → comment_add with parent_id
- "flag inappropriate" → comment_flag or review_flag
- "most reviewed items" → review_get_most_reviewed

### 9. SUPPORT & HELP
**User says** → **Use tool**
- "I need help" → support_start + support_create_ticket
- "create support ticket" → support_create_ticket
- "contact support" → support_create_ticket
- "check ticket status" → support_get_ticket
- "close my ticket" → support_close_ticket
- "urgent issue" → support_create_ticket with priority="urgent"

### 10. WISHLIST & FAVORITES
**User says** → **Use tool**
- "add to wishlist" → wishlist_add_default
- "save for later" → wishlist_add_default
- "favorite this" → wishlist_add_default
- "remove from wishlist" → wishlist_remove_default
- "show my wishlist" → wishlist_get_user
- "clear wishlist" → wishlist_clear
- "is this wishlisted" → wishlist_is_in

### 11. SHIPPING & DELIVERY
**User says** → **Use tool**
- "track package" → shipping_track
- "where's my order" → shipping_track
- "create shipping label" → shipping_create
- "shipping status" → shipping_track with tracking_number

### 12. SOCIAL FEATURES
**User says** → **Use tool**
- "follow user" → following_follow
- "unfollow" → following_unfollow
- "my followers" → following_get_followers
- "who am I following" → following_get_following
- "mutual followers" → following_get_mutual

## PARAMETER EXTRACTION RULES

### Money/Price Conversion
- Natural language → Cents (integer)
- "$50", "50 dollars", "fifty bucks" → 5000
- "$19.99", "19.99" → 1999
- "under $100" → max_price: 10000
- "over $50" → min_price: 5000
- "between $20 and $50" → min_price: 2000, max_price: 5000
- "around $30" → min_price: 2500, max_price: 3500

### Quantity Parsing
- "one", "a", "an" → 1
- "a couple", "two" → 2
- "a few" → 3
- "several" → 5
- "many" → 10
- "a dozen" → 12
- "lots of" → 20

### Time References
- "today" → current date
- "yesterday" → current date - 1 day
- "last week" → date range for past 7 days
- "this month" → current month date range
- "recently" → last 30 days

### Status Mapping
**Orders:**
- "not processed yet" → "pending"
- "being prepared" → "processing"
- "on the way" → "shipped"
- "received" → "delivered"
- "cancelled" → "cancelled"

**Priorities:**
- "not urgent" → "low"
- "normal" → "medium"
- "important" → "high"
- "asap", "emergency" → "urgent"

### Sorting Preferences
- "cheapest first" → sort_by="price", sort_order="asc"
- "most expensive" → sort_by="price", sort_order="desc"
- "newest" → sort_by="created_at", sort_order="desc"
- "oldest" → sort_by="created_at", sort_order="asc"
- "most popular" → sort_by="popularity", sort_order="desc"
- "best rated" → sort_by="rating", sort_order="desc"

### Entity Type Recognition
- Products: "item", "product", "goods", "merchandise"
- Services: "service", "offering", "consultation"
- Users: "person", "user", "member", "customer", "seller"
- Orders: "purchase", "order", "transaction"

### Context Resolution
- "it" → last mentioned entity
- "that one" → last displayed item
- "the first one" → items[0] from last search
- "all of them" → all items from last result
- "my" → current user context
- "their" → last mentioned user

## MULTI-STEP OPERATIONS

### Complete Purchase Flow
1. "I want to buy [product]" →
   - product_search_advanced
   - basket_add_item
   - order_create
   - payment_create

### Review After Purchase
1. "review my recent purchase" →
   - order_get_user (get recent orders)
   - review_add (for the order/product)

### Price Comparison
1. "compare prices for [item]" →
   - product_search_advanced (multiple searches)
   - vector_search_similar (find alternatives)

### Support Flow
1. "I have a problem with my order" →
   - order_get_user (find order)
   - support_create_ticket (with order details)
`