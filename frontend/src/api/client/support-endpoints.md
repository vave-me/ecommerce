# Support Service API Endpoints

## Base URL
All endpoints are prefixed with `/api/support`

## Support Channel Management
- `POST /api/support/channels` - Create support channel
- `GET /api/support/channels/{id}` - Get support channel
- `GET /api/support/users/{user_id}/channels` - Get user support channels
- `PUT /api/support/channels/{id}/settings` - Update channel settings
- `POST /api/support/channels/{id}/close` - Close support channel
- `POST /api/support/channels/{id}/reactivate` - Reactivate support channel

## Ticket Management
- `POST /api/support/tickets` - Create ticket
- `GET /api/support/tickets/{id}` - Get ticket
- `POST /api/support/tickets/batch` - Get multiple tickets
- `GET /api/support/channels/{channel_id}/tickets` - Get channel tickets
- `GET /api/support/tickets` - List all tickets (admin)
- `PUT /api/support/tickets/{id}` - Update ticket
- `POST /api/support/tickets/{id}/assign` - Assign ticket
- `PUT /api/support/tickets/{id}/priority` - Update ticket priority
- `POST /api/support/tickets/{id}/escalate` - Escalate ticket
- `POST /api/support/tickets/{id}/resolve` - Resolve ticket
- `POST /api/support/tickets/{id}/reopen` - Reopen ticket
- `POST /api/support/tickets/{id}/close` - Close ticket
- `POST /api/support/tickets/merge` - Merge tickets
- `POST /api/support/tickets/{ticket_id}/link` - Link tickets

## Communication
- `POST /api/support/tickets/{ticket_id}/replies` - Add ticket reply
- `POST /api/support/tickets/{ticket_id}/notes` - Add internal note
- `GET /api/support/tickets/{ticket_id}/communications` - Get ticket communications

## AI Integration
- `POST /api/support/channels/{channel_id}/ai/enable` - Enable AI support
- `PUT /api/support/channels/{channel_id}/ai/config` - Configure AI assistant
- `POST /api/support/tickets/{ticket_id}/handoff/human` - Handoff to human
- `POST /api/support/tickets/{ticket_id}/handoff/ai` - Handoff to AI
- `GET /api/support/tickets/{ticket_id}/ai/suggestions` - Get AI suggestions

## Knowledge Base
- `POST /api/support/knowledge` - Create knowledge article
- `GET /api/support/knowledge/{id}` - Get knowledge article
- `POST /api/support/knowledge/search` - Search knowledge base
- `POST /api/support/tickets/{ticket_id}/articles` - Link article to ticket
- `POST /api/support/knowledge/{article_id}/rate` - Rate article

## Analytics
- `GET /api/support/analytics/metrics` - Get support metrics
- `GET /api/support/analytics/agents/{agent_id}` - Get agent performance
- `GET /api/support/analytics/tickets` - Get ticket analytics

## Frontend Integration Status

### Working Features:
1. ✅ Support channel creation and management
2. ✅ Ticket creation with proper validation
3. ✅ Ticket listing and filtering
4. ✅ Real-time ticket updates
5. ✅ Priority and status management
6. ✅ Reply and internal note functionality
7. ✅ Knowledge base search
8. ✅ Support metrics display

### API Connection:
- Frontend uses axios with base URL `/api`
- Support API endpoints use `/support` prefix
- Error handling with fallback to mock data
- React Query for caching and synchronization

### Authentication:
- User ID from AuthContext
- JWT tokens handled by axios interceptor
- Admin endpoints require proper permissions