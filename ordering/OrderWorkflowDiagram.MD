Title: Place Order Saga (Choreography)

participant User
participant BasketService
participant OrderService
participant InventoryService
participant PaymentService
participant ShippingService
participant NotificationService
note over OrderService: (Broker) Publishes & Subscribes Domain Events
note over InventoryService: (Broker) Publishes & Subscribes Domain Events
note over PaymentService: (Broker) Publishes & Subscribes Domain Events
note over ShippingService: (Broker) Publishes & Subscribes Domain Events
note over NotificationService: (Broker) Subscribes to relevant events

== STEP 1: Basket Checkout ==
User -> BasketService: CheckoutBasketCommand
note over BasketService:
1) Validate basket (items, user info, coupons, etc.)
2) Optionally finalize basket for checkout

BasketService -> OrderService: CreateOrderCommand (sync or async)

== STEP 2: Order Created ==
OrderService -> OrderService:
1) Create Order entity in DB (status = CREATED)
2) Raise domain event: OrderCreatedEvent
OrderService -> (Broker): Publish OrderCreatedEvent

== STEP 3: (Optional) Reserve Inventory ==
InventoryService: Subscribed to OrderCreatedEvent
InventoryService -> InventoryService:
1) Check product stock for each line item
2) Reserve or decrement available quantity
3) If insufficient stock => publish StockReservationFailedEvent
4) Else => raise domain event ProductStockAdjustedEvent
InventoryService -> (Broker):
- If insufficient => Publish StockReservationFailedEvent
- Else => Publish ProductStockAdjustedEvent

OrderService: Subscribed to StockReservationFailedEvent
OrderService -> OrderService: Set order status = CANCELED
OrderService -> (Broker): Publish OrderCanceledEvent

== STEP 4: Payment Authorization ==
note over OrderService: Once inventory is confirmed, or directly if no reservation step
OrderService -> PaymentService: RequestPaymentCommand (sync or broker request)
PaymentService -> PaymentService:
1) Contact external payment gateway
2) If success => raise PaymentAuthorizedEvent
3) Else => raise PaymentFailedEvent
PaymentService -> (Broker):
- Publish PaymentAuthorizedEvent or PaymentFailedEvent

== STEP 5: Order Approved or Canceled ==
OrderService: Subscribed to PaymentAuthorizedEvent
OrderService -> OrderService: Set order status = APPROVED
OrderService: Raise domain event OrderApprovedEvent
OrderService -> (Broker): Publish OrderApprovedEvent

OrderService: Subscribed to PaymentFailedEvent or StockReservationFailedEvent
OrderService -> OrderService: Set order status = CANCELED
OrderService -> (Broker): Publish OrderCanceledEvent

== STEP 6: Shipping Creation ==
ShippingService: Subscribed to OrderApprovedEvent
ShippingService -> ShippingService:
1) Create shipment record
2) Raise ShipmentCreatedEvent
ShippingService -> (Broker): Publish ShipmentCreatedEvent

== STEP 7: Shipment Progress ==
note over ShippingService: As the package moves through statuses (e.g., "In Transit", "Delivered")
ShippingService -> (Broker): Publish ShipmentStatusChangedEvent
OrderService: Subscribed => Update Order status (SHIPPED / DELIVERED, etc.)
NotificationService: Subscribed => Notifies user:
"Your order is now in transit" or "Your order was delivered."

== STEP 8: Notifications ==
NotificationService: Also Subscribed to:
- OrderCreatedEvent => "Order #123 created"
- OrderApprovedEvent => "Order #123 approved"
- OrderCanceledEvent => "Order #123 canceled"
- ShipmentCreatedEvent => "Shipping label created"
- ShipmentStatusChangedEvent => "Order in transit/delivered"
NotificationService -> (Email, SMS, Push, In-app):
Send notifications to user channels
