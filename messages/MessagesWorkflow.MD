Title: Comprehensive Messaging Workflow

participant User
participant MessageService
participant NotificationService
participant MailerService
note over MessageService, MailerService: (Broker, e.g., NATS/Kafka, handles pub/sub for integration events)

== 1) Start Conversation ==
User -> MessageService: StartConversationCommand(participants, optional metadata)
note over MessageService:
1. Validate participants (check existence, permissions).
2. Create Conversation entity in DB.
3. Raise domain event: ConversationStarted.
   MessageService -> Broker: Publish ConversationStarted (integration event)

NotificationService: Subscribed to ConversationStarted
NotificationService -> NotificationService: Creates in-app/push notifications for participants
NotificationService: (Optional) Raise NotificationCreated
NotificationService -> Broker: Publish NotificationCreated (if other services need it)

MailerService: Subscribed to NotificationCreated or ConversationStarted (depending on design)
MailerService -> MailerService: Check user email preference, compose "new conversation" email
MailerService -> Broker: Publish EmailCreatedEvent (for auditing or further processing)

== 2) Send Message ==
User -> MessageService: SendMessageCommand(conversationId, content)
note over MessageService:
1. Check user is participant in the conversation and not blocked.
2. Store message in DB (message table or doc store).
3. Raise domain event: MessageSent.
   MessageService -> Broker: Publish MessageSent (integration event)

NotificationService: Subscribed to MessageSent
NotificationService -> NotificationService:
1. Create or update in-app notification for other participants
2. Possibly raise NotificationCreated
   NotificationService -> Broker: Publish NotificationCreated

MailerService: Subscribed to NotificationCreated (or directly to MessageSent)
MailerService -> MailerService:
1. Check if participant is offline or needs email
2. Send email via SMTP/gateway
3. Raise EmailCreatedEvent or EmailSentEvent
   MailerService -> Broker: Publish EmailCreatedEvent

== 3) Read Receipt ==
User -> MessageService: MarkMessageReadCommand(conversationId, messageId)
note over MessageService:
1. Validate user is participant
2. Mark message as read in DB
3. Raise domain event: MessageRead
   MessageService -> Broker: Publish MessageRead (integration event)

NotificationService: Subscribed to MessageRead
NotificationService -> NotificationService:
1. Update notification state (e.g., remove "unread" badge for the message)
2. Possibly raise NotificationUpdated
   NotificationService -> Broker: Publish NotificationUpdated

== 4) Optional: Real-Time Push ==
NotificationService -> WebSocket/Push Service (or direct to client):
1. "User A read the message" or "New message arrived"

== 5) Edit or Delete Message (if allowed) ==
User -> MessageService: EditMessageCommand or DeleteMessageCommand
MessageService -> MessageService:
1. Validate user permission to edit/delete
2. Update or soft-delete in DB
3. Raise domain event: MessageEdited or MessageDeleted
   MessageService -> Broker: Publish MessageEdited/MessageDeleted

NotificationService: Subscribed to MessageEdited/MessageDeleted
NotificationService -> NotificationService: Possibly alert other participants
MailerService: Could also notify or send updates (less common for editing, but possible)

== 6) Archive or Close Conversation (if allowed) ==
User -> MessageService: ArchiveConversationCommand(conversationId)
MessageService:
1. Mark conversation status = archived for this user
2. Raise domain event: ConversationArchived
   MessageService -> Broker: Publish ConversationArchived

NotificationService: Possibly reflect updated status
MailerService: Rare to email about archived conversation, but possible

== 7) Activity & Analytics (Optional) ==
ActivityService: Subscribed to all message & conversation events
ActivityService -> ActivityService:
1. Log or store usage metrics (ConversationStarted, MessageSent, etc.)
2. Possibly trigger alerts for unusual activity or compliance logging