# Tool Service Operations Reference

This document lists all available operations for each tool service in the assistants/internal/application/tools directory.

## Activity Tool Service
- `find`, `get_interaction` - Find an interaction by ID
- `create_activity` - Create a new activity for a user
- `find_activity` - Find activity by user ID
- `add_like` - Add a like interaction
- `add_dislike` - Add a dislike interaction
- `update_interaction` - Update an existing interaction
- `remove_interaction` - Remove an interaction
- `get_interactions` - Get all interactions for an activity
- `get_most_liked` - Get most liked items of a type
- `get_most_disliked` - Get most disliked items of a type
- `archive_activity` - Archive an activity
- `restore_activity` - Restore an archived activity

## Basket Tool Service
- `start_basket`, `create` - Create a new basket
- `get_basket`, `find` - Get basket by ID
- `get_current_basket` - Get current basket for user
- `add_item` - Add item to basket
- `remove_item` - Remove item from basket
- `update_quantity` - Update item quantity
- `get_total` - Get basket total
- `checkout` - Checkout basket
- `cancel` - Cancel basket
- `list_baskets` - List all baskets
- `clear_basket` - Clear basket contents
- `get_analytics` - Get basket analytics

## Category Tool Service
- `add`, `create` - Create a new category
- `get`, `find` - Get category by ID
- `get_by_slug` - Get category by slug
- `list`, `get_all` - List all categories
- `get_main` - Get main categories
- `get_sub` - Get subcategories
- `update` - Update category
- `remove`, `delete` - Delete category
- `archive` - Archive category
- `rebrand` - Rebrand category
- `search` - Search categories

## Comment Tool Service
- `add`, `create` - Create a new comment
- `get`, `find` - Get comment by ID
- `get_comments`, `list` - List comments
- `approve` - Approve comment
- `reject` - Reject comment
- `flag` - Flag comment
- `edit` - Edit comment
- `delete` - Delete comment

## Geocoding Tool Service
- `geocode`, `encode` - Convert address to coordinates
- `reverse_geocode`, `decode` - Convert coordinates to address
- `get_coordinates` - Get coordinates for address
- `get_address` - Get address for coordinates
- `search_places` - Search for places
- `get_nearby` - Get nearby locations
- `validate_address` - Validate address
- `normalize_address` - Normalize address format

## Mailer Tool Service
- `send_email` - Send email
- `send_marketing_email` - Send marketing email
- `send_transactional_email` - Send transactional email
- `send_bulk_email` - Send bulk emails
- `get_email_status` - Get email status
- `get_email_stats` - Get email statistics
- `schedule_email` - Schedule email for later
- `cancel_scheduled_email` - Cancel scheduled email
- `get_templates` - Get email templates
- `validate_email` - Validate email address

## Media Tool Service
- `upload_image` - Upload image
- `upload_video` - Upload video
- `get_media`, `find` - Get media by ID
- `list_media`, `list` - List media
- `delete_media` - Delete media
- `update_media` - Update media metadata
- `get_user_media` - Get media by user
- `search_media` - Search media
- `get_media_by_type` - Get media by type
- `generate_thumbnail` - Generate thumbnail

## Message Tool Service
- `send_message`, `create` - Send a message
- `get_message`, `find` - Get message by ID
- `list_messages`, `list` - List messages
- `get_conversation` - Get conversation between users
- `mark_as_read` - Mark message as read
- `delete_message` - Delete message
- `search_messages` - Search messages
- `get_unread_count` - Get unread message count

## Metric Tool Service
- `track_event` - Track an event
- `get_metrics` - Get metrics
- `get_event_count` - Get event count
- `get_unique_users` - Get unique user count
- `get_conversion_rate` - Get conversion rate
- `get_user_metrics` - Get metrics for user
- `get_product_metrics` - Get metrics for product
- `get_revenue_metrics` - Get revenue metrics
- `get_funnel_analysis` - Get funnel analysis
- `export_metrics` - Export metrics data

## Newsletter Tool Service
- `subscribe` - Subscribe to newsletter
- `unsubscribe` - Unsubscribe from newsletter
- `update_preferences` - Update subscription preferences
- `get_subscriber`, `find` - Get subscriber details
- `list_subscribers`, `list` - List subscribers
- `send_newsletter` - Send newsletter
- `schedule_newsletter` - Schedule newsletter
- `get_campaigns` - Get newsletter campaigns
- `get_campaign_stats` - Get campaign statistics
- `export_subscribers` - Export subscriber list

## Notification Tool Service
- `list_alerts` - List all alerts
- `get_alerts_by_type` - Get alerts by type
- `get_alert`, `find` - Get alert by ID
- `get_alerts_by_user`, `search` - Get alerts for user
- `get_unread_alerts` - Get unread alerts
- `get_read_alerts` - Get read alerts
- `search_alerts` - Search alerts
- `count_alerts` - Count total alerts
- `count_unread_alerts` - Count unread alerts

## Offer Tool Service
### General Offers
- `create_offer`, `create` - Create offer
- `activate_offer`, `activate` - Activate offer
- `close_offer`, `close` - Close offer
- `accept_offer`, `accept` - Accept offer
- `get_offer`, `find` - Get offer by ID
- `list_offers`, `get_offers` - List offers
- `get_offers_by_product` - Get offers for product
- `get_offers_by_user` - Get offers by user
- `search_offers`, `search` - Search offers
- `request_offer_negotiation` - Request negotiation
- `accept_offer_negotiation` - Accept negotiation
- `decline_offer_negotiation` - Decline negotiation

### Buy Now
- `create_buynow`, `create_buy_now` - Create buy now
- `confirm_buynow`, `confirm_buy_now` - Confirm buy now
- `cancel_buynow`, `cancel_buy_now` - Cancel buy now
- `request_buynow_negotiation`, `request_buy_now_negotiation` - Request negotiation
- `accept_buynow_negotiation`, `accept_buy_now_negotiation` - Accept negotiation
- `decline_buynow_negotiation`, `decline_buy_now_negotiation` - Decline negotiation

### Lease
- `create_lease` - Create lease
- `start_lease` - Start lease
- `make_lease_payment` - Make lease payment
- `execute_lease_buyout` - Execute lease buyout
- `end_lease` - End lease
- `cancel_lease` - Cancel lease
- `default_lease` - Default on lease
- `get_active_leases` - Get active leases
- `request_lease_negotiation` - Request negotiation
- `accept_lease_negotiation` - Accept negotiation
- `decline_lease_negotiation` - Decline negotiation

### Buyback
- `create_buyback`, `create_buy_back` - Create buyback
- `redeem_buyback`, `redeem_buy_back` - Redeem buyback
- `expire_buyback`, `expire_buy_back` - Expire buyback
- `cancel_buyback`, `cancel_buy_back` - Cancel buyback
- `get_active_buybacks`, `get_active_buy_backs` - Get active buybacks
- `request_buyback_negotiation`, `request_buy_back_negotiation` - Request negotiation
- `accept_buyback_negotiation`, `accept_buy_back_negotiation` - Accept negotiation
- `decline_buyback_negotiation`, `decline_buy_back_negotiation` - Decline negotiation

### Reservation
- `create_reservation` - Create reservation
- `redeem_reservation` - Redeem reservation
- `expire_reservation` - Expire reservation
- `cancel_reservation` - Cancel reservation
- `get_active_reservations` - Get active reservations
- `request_reservation_negotiation` - Request negotiation
- `accept_reservation_negotiation` - Accept negotiation
- `decline_reservation_negotiation` - Decline negotiation

## Order Tool Service
- `create_order`, `create` - Create order
- `get_order`, `find` - Get order by ID
- `list_orders`, `list` - List orders
- `update_order_status` - Update order status
- `cancel_order`, `cancel` - Cancel order
- `get_user_orders` - Get orders by user
- `search_orders`, `search` - Search orders
- `get_order_items` - Get order items
- `add_order_note` - Add note to order
- `process_return` - Process return
- `process_refund` - Process refund
- `get_order_history` - Get order history
- `export_orders` - Export orders

## Payment Tool Service
- `create_payment_intent`, `create_intent` - Create payment intent
- `capture_payment` - Capture payment
- `cancel_payment` - Cancel payment
- `refund_payment`, `refund` - Refund payment
- `get_payment`, `find` - Get payment by ID
- `list_payments`, `list` - List payments
- `get_payment_methods` - Get payment methods
- `add_payment_method` - Add payment method
- `remove_payment_method` - Remove payment method
- `get_transaction_history` - Get transaction history
- `process_payout` - Process payout
- `get_balance` - Get account balance

## Post Tool Service
- `create_post`, `create` - Create post
- `get_post`, `find` - Get post by ID
- `update_post`, `update` - Update post
- `delete_post`, `delete` - Delete post
- `list_posts`, `list` - List posts
- `get_user_posts` - Get posts by user
- `search_posts`, `search` - Search posts
- `publish_post` - Publish post
- `unpublish_post` - Unpublish post
- `get_draft_posts` - Get draft posts
- `get_published_posts` - Get published posts

## Product Tool Service
- `search`, `search_by_term` - Search products
- `find`, `get` - Get product by ID
- `filter`, `search_with_filters` - Search with filters
- `suggest` - Get product suggestions
- `category_search` - Search by category
- `add`, `create` - Create product
- `update` - Update product
- `remove`, `delete` - Delete product

## Review Tool Service
- `create_review`, `create` - Create review
- `get_review`, `find` - Get review by ID
- `update_review`, `update` - Update review
- `delete_review`, `delete` - Delete review
- `list_reviews`, `list` - List reviews
- `get_product_reviews` - Get reviews for product
- `get_user_reviews` - Get reviews by user
- `approve_review` - Approve review
- `reject_review` - Reject review
- `flag_review` - Flag review
- `respond_to_review` - Respond to review
- `get_review_stats` - Get review statistics

## Service Tool Service
- `find`, `get` - Get service by ID
- `search` - Search services
- `suggest` - Get service suggestions
- `get_services`, `list` - List services
- `get_catalog`, `catalog` - Get service catalog
- `get_public_catalog`, `public_catalog` - Get public catalog
- `search_by_category` - Search by category ID
- `search_by_category_slug` - Search by category slug
- `add`, `create` - Create service
- `update` - Update service
- `remove`, `delete` - Delete service
- `archive` - Archive service
- `mark_sold` - Mark service as sold
- `mark_leased` - Mark service as leased
- `filter` - Filter services

Additional filter options:
- `active` - Filter active services
- `inactive` - Filter inactive services
- `business` - Filter business services
- `individual`, `private` - Filter individual/private services

## Shipping Tool Service
- `create_shipping`, `create` - Create shipping
- `track_shipping`, `track`, `get_tracking` - Track shipping
- `update_shipping_status`, `update_status` - Update status
- `get_shipping_rates`, `calculate_rates` - Get rates
- `get_shipping_by_order`, `get_order_shipping` - Get by order
- `cancel_shipping`, `cancel` - Cancel shipping
- `get_shipping_labels`, `generate_label` - Get labels

## Support Tool Service
- `start_support`, `create_support` - Start support session
- `create_ticket`, `create` - Create ticket
- `list_tickets`, `list` - List tickets
- `get_ticket`, `find` - Get ticket by ID
- `update_ticket`, `update` - Update ticket
- `delete_ticket`, `close_ticket`, `close` - Close ticket
- `add_ticket_comment`, `comment` - Add comment
- `get_ticket_comments` - Get ticket comments
- `search_tickets`, `search` - Search tickets
- `get_user_tickets` - Get tickets by user

## User Tool Service
- `find`, `get` - Find user
- `get_base_user` - Get base user info
- `list_users` - List all users
- `list_participating_users` - List participating users
- `create_user`, `create` - Create user
- `update_user`, `update` - Update user
- `login_user`, `login` - Login user
- `search`, `search_users` - Search users

## Variant Tool Service
- `create_variant`, `create` - Create variant
- `get_variant`, `find` - Get variant by ID
- `update_variant`, `update` - Update variant
- `delete_variant`, `delete` - Delete variant
- `list_variants`, `list` - List variants
- `get_product_variants` - Get product variants
- `search_variants` - Search variants
- `update_inventory` - Update inventory
- `get_inventory` - Get inventory

## Vector Tool Service
- `search_similar_entities`, `search_similar` - Search similar entities
- `get_entity_context`, `get_context` - Get entity context
- `get_recommendations`, `recommend` - Get recommendations
- `check_vector_service_health`, `health` - Check service health

## Wishlist Tool Service
- `add_to_wishlist`, `add` - Add to wishlist
- `remove_from_wishlist`, `remove` - Remove from wishlist
- `get_wishlist`, `find` - Get wishlist
- `list_wishlists`, `list` - List wishlists
- `clear_wishlist`, `clear` - Clear wishlist
- `share_wishlist` - Share wishlist
- `get_shared_wishlist` - Get shared wishlist
- `search_wishlists`, `search` - Search wishlists
- `move_to_cart` - Move item to cart
- `check_item_in_wishlist` - Check if item in wishlist

---

Note: Many operations have aliases (alternative names) that can be used interchangeably. The first name listed is typically the primary operation name, followed by common aliases.