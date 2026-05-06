package constants

// NaturalLanguagePatterns helps the LLM understand common user intents
const NaturalLanguagePatterns = `
## Intent Recognition Patterns

### Search/Find Operations
- "find", "search", "look for", "show me" → search operations
- "what's available", "list all" → list/search operations
- "similar to", "like this" → search with similarity

### Create Operations  
- "add", "create", "make", "new" → create operations
- "I want to sell" → product_create
- "place an order" → order_create

### Update Operations
- "change", "update", "modify", "edit" → update operations
- "set the price to" → product_update with price
- "mark as delivered" → order_update_status

### Delete/Remove Operations
- "remove", "delete", "cancel" → delete operations
- "clear my basket" → basket_clear
- "cancel order" → order_update_status with cancelled

## Parameter Extraction

### Price/Money
- "fifty dollars", "$50", "50 bucks" → 5000 (convert to cents)
- "under $100" → max_price: 10000
- "between $50 and $100" → min_price: 5000, max_price: 10000

### Quantities
- "a couple", "two" → 2
- "a few" → 3-5
- "several" → 5-10

### Status Values
- "delivered", "shipped", "pending" → order status
- "active", "inactive" → user status
- "in stock", "out of stock" → product availability

### Common References
- "my orders" → filter by current user_id
- "that product" → reference to previously mentioned item
- "it" → last discussed entity
`