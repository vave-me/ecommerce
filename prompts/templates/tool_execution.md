<!-- Version: 1.0.0 -->
<!-- Last Updated: 2025-07-21 -->
<!-- Status: Active -->

# Tool Execution Guidelines

## Execution Limits

### Per Response Limits
- Maximum tools per response: 10
- Maximum schema consultations: 3
- Maximum iterations: 5

### Tool Categories & Priorities

1. **Direct Action Tools** (Priority 1)
   - product_search, user_find, order_track
   - Use for immediate user needs

2. **Data Modification Tools** (Priority 2)
   - product_add, order_create, review_add
   - Require validation before execution

3. **Analytics Tools** (Priority 3)
   - metrics_get, statistics_view
   - Use for insights and reporting

4. **Schema Tools** (Priority 4 - USE SPARINGLY)
   - schema_get_fields, schema_generate_contextual_help
   - Maximum 3 calls per response
   - Only use when truly needed

## Execution Patterns

### Sequential Execution
For dependent operations:
```
1. Validate user exists → 2. Create order → 3. Send confirmation
```

### Parallel Execution  
For independent operations:
```
Execute simultaneously:
- Search products
- Get user profile
- Check availability
```

### Conditional Execution
Based on previous results:
```
IF search returns results:
  → Show products
ELSE:
  → Suggest alternatives
```

## Common Tool Combinations

### Product Discovery Flow
1. product_search → Get initial results
2. category_get → Show related categories
3. product_get_details → Deep dive on selection

### Order Creation Flow
1. product_check_availability → Verify stock
2. user_validate → Confirm user
3. order_create → Create order
4. notification_send → Confirm to user

### User Research Flow
1. user_find → Locate user
2. user_get_ratings → Check reputation
3. product_get_by_seller → View their products

## Error Handling

### Tool Failure Recovery
```
Primary: product_search(name="laptop")
Fallback: category_browse(type="electronics")
Last Resort: schema_generate_contextual_help("finding laptops")
```

### Validation Failures
```
If validation fails:
1. Explain issue clearly
2. Request missing information
3. Provide valid examples
```

## Anti-Pattern Prevention

### ❌ AVOID: Schema Loop
```
BAD:
schema_generate_contextual_help → 
schema_generate_contextual_help →
schema_generate_contextual_help (LOOP!)
```

### ✅ BETTER: Direct Response
```
GOOD:
Recognize help request → Provide capability overview → Ask for specifics
```

### ❌ AVOID: Over-Tooling
```
BAD:
Simple "hi" → 15 tool calls to understand greeting
```

### ✅ BETTER: Appropriate Response
```
GOOD:
Simple "hi" → Friendly greeting + capability overview
```

## Tool Selection Matrix

| User Intent | Primary Tool | Fallback | Avoid |
|------------|--------------|----------|--------|
| "Find X" | product_search | category_browse | schema_tools |
| "What can you do?" | None - direct response | help_template | schema_loops |
| "Track order" | order_track | order_list | multiple schemas |
| "User info" | user_find | user_search | excessive validation |

## Performance Guidelines

### Fast Path (< 1 second)
- Single tool execution
- Direct response
- Cached results

### Normal Path (1-3 seconds)
- 2-3 tool executions
- Simple validation
- Standard operations

### Slow Path (> 3 seconds) - AVOID
- Multiple schema consultations
- Excessive iterations
- Complex nested operations

## Success Metrics

Good tool execution has:
- ✓ Clear purpose
- ✓ Efficient path
- ✓ Graceful failures
- ✓ User value
- ✓ Fast response