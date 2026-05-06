package domain

import (
	"github.com/stackus/errors"
	"middleman/internal/ddd"
	"middleman/internal/es"
)

const AlertAggregate = "notifications.Alert"

var (
	ErrAlertAlreadyCreated = errors.Wrap(errors.ErrBadRequest, "the notification cannot be recreated")
	ErrAlertHasNoPayload   = errors.Wrap(errors.ErrBadRequest, "the notification has no payload")
	ErrUserIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the notification ID cannot be blank")
	ErrAlertNotFound       = errors.Wrap(errors.ErrNotFound, "alert not found")
)

type Alert struct {
	es.Aggregate
	UserID    string
	AlertType AlertType
	Message   string
	Payload   map[string]interface{}
	IsRead    bool
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*Alert)(nil)

func NewAlert(id string) *Alert {
	return &Alert{
		Aggregate: es.NewAggregate(id, AlertAggregate),
	}
}

func (n *Alert) AddBasketAlert(id, userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddBasketAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	basketID, _ := payload["BasketID"].(string)
	productID, _ := payload["ProductID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &BasketAlertAdded{
		ID:        n.ID(),
		AlertType: string(alertType),
		UserID:    userID,
		BasketID:  basketID,
		ProductID: productID,
		Payload:   payload,
		Message:   message,
		IsRead:    false,
	}
	n.AddEvent(BasketAlertAddedEvent, event)
	return ddd.NewEvent(BasketAlertAddedEvent, n), nil
}

func (n *Alert) AddWishlistAlert(userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddWishlistAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	wishlistID, _ := payload["WishlistID"].(string)
	productID, _ := payload["ProductID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &WishlistAlertAdded{
		ID:         n.ID(),
		AlertType:  string(alertType),
		UserID:     userID,
		WishlistID: wishlistID,
		ProductID:  productID,
		Payload:    payload,
		Message:    message,
		IsRead:     false,
	}
	n.AddEvent(WishlistAlertAddedEvent, event)
	return ddd.NewEvent(WishlistAlertAddedEvent, n), nil

}
func (n *Alert) AddSupportAlert(userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddSupportAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	ticketID, _ := payload["TicketID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &SupportAlertAdded{
		ID:        n.ID(),
		AlertType: string(alertType),
		UserID:    userID,
		TicketID:  ticketID,
		Payload:   payload,
		Message:   message,
		IsRead:    false,
	}
	n.AddEvent(SupportAlertAddedEvent, event)
	return ddd.NewEvent(SupportAlertAddedEvent, n), nil

}
func (n *Alert) AddOfferAlert(userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddOfferAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	offerID, _ := payload["OfferID"].(string)
	productID, _ := payload["ProductID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &OfferAlertAdded{
		ID:        n.ID(),
		AlertType: string(alertType),
		UserID:    userID,
		OfferID:   offerID,
		ProductID: productID,
		Payload:   payload,
		Message:   message,
		IsRead:    false,
	}
	n.AddEvent(OfferAlertAddedEvent, event)
	return ddd.NewEvent(OfferAlertAddedEvent, n), nil

}
func (n *Alert) AddOrderAlert(userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddOrderAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	orderID, _ := payload["OrderID"].(string)
	productID, _ := payload["ProductID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &OrderAlertAdded{
		ID:        n.ID(),
		AlertType: string(alertType),
		UserID:    userID,
		OrderID:   orderID,
		ProductID: productID,
		Payload:   payload,
		Message:   message,
		IsRead:    false,
	}
	n.AddEvent(OrderAlertAddedEvent, event)
	return ddd.NewEvent(OrderAlertAddedEvent, n), nil

}
func (n *Alert) AddMessageAlert(userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	// Handle notification types

	err := validateAddMessageAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	messageID, _ := payload["MessageID"].(string)
	senderID, _ := payload["MessageSenderID"].(string)
	productID, _ := payload["ProductID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &MessageAlertAdded{
		ID:        n.ID(),
		AlertType: string(alertType),
		MessageID: messageID,
		UserID:    userID, // ID of the notification receiving the notification
		SenderID:  senderID,
		ProductID: productID,
		Payload:   payload,
		Message:   message,
		IsRead:    false,
	}
	n.AddEvent(MessageAlertAddedEvent, event)
	return ddd.NewEvent(MessageAlertAddedEvent, n), nil

}

func (n *Alert) AddCommentAlert(userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddCommentAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	commentID, _ := payload["CommentID"].(string)
	senderID, _ := payload["UserAddedID"].(string)
	productID, _ := payload["ProductID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &CommentAlertAdded{
		ID:        n.ID(),
		AlertType: string(alertType),
		UserID:    userID,
		CommentID: commentID,
		SenderID:  senderID,
		ProductID: productID,
		Payload:   payload,
		Message:   message,
		IsRead:    false,
	}
	n.AddEvent(CommentAlertAddedEvent, event)
	return ddd.NewEvent(CommentAlertAddedEvent, n), nil

}

func (n *Alert) AddProductAlert(userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddProductAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	productID, _ := payload["ProductID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &ProductAlertAdded{
		ID:        n.ID(),
		AlertType: string(alertType),
		UserID:    userID,
		ProductID: productID,
		Payload:   payload,
		Message:   message,
		IsRead:    false,
	}
	n.AddEvent(ProductAlertAddedEvent, event)
	return ddd.NewEvent(ProductAlertAddedEvent, n), nil
}
func (n *Alert) Read(id string) (ddd.Event, error) {

	event := &AlertRead{
		ID:     id,
		IsRead: true,
	}
	n.AddEvent(AlertReadEvent, event)
	// Update the aggregate state
	n.IsRead = true
	return ddd.NewEvent(AlertReadEvent, n), nil
}

func (n *Alert) AddInteractionAlert(id, userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddInteractionAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	productID, _ := payload["ProductID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &InteractionAlertAdded{
		ID:        n.ID(),
		AlertType: string(alertType),
		UserID:    userID,
		ProductID: productID,
		Payload:   payload,
		Message:   message,
		IsRead:    false,
	}
	n.AddEvent(InteractionAlertAddedEvent, event)
	return ddd.NewEvent(InteractionAlertAddedEvent, n), nil
}

// Payload Validation Functions

func validateAddBasketAlertPayload(payload map[string]interface{}) error {
	if basketID, ok := payload["BasketID"].(string); !ok || basketID == "" {
		return errors.Wrap(errors.ErrBadRequest, "BasketID is required and cannot be empty")
	}
	if productID, ok := payload["ProductID"].(string); !ok || productID == "" {
		return errors.Wrap(errors.ErrBadRequest, "ProductID is required and cannot be empty")
	}
	return nil
}

func validateAddWishlistAlertPayload(payload map[string]interface{}) error {
	if wishlistID, ok := payload["WishlistID"].(string); !ok || wishlistID == "" {
		return errors.Wrap(errors.ErrBadRequest, "WishlistID is required and cannot be empty")
	}
	if productID, ok := payload["ProductID"].(string); !ok || productID == "" {
		return errors.Wrap(errors.ErrBadRequest, "ProductID is required and cannot be empty")
	}
	return nil
}
func validateAddSupportAlertPayload(payload map[string]interface{}) error {
	if supportID, ok := payload["SupportID"].(string); !ok || supportID == "" {
		return errors.Wrap(errors.ErrBadRequest, "WishlistID is required and cannot be empty")
	}
	return nil
}
func validateAddOfferAlertPayload(payload map[string]interface{}) error {
	if wishlistID, ok := payload["OfferID"].(string); !ok || wishlistID == "" {
		return errors.Wrap(errors.ErrBadRequest, "WishlistID is required and cannot be empty")
	}
	if productID, ok := payload["ProductID"].(string); !ok || productID == "" {
		return errors.Wrap(errors.ErrBadRequest, "ProductID is required and cannot be empty")
	}
	return nil
}
func validateAddOrderAlertPayload(payload map[string]interface{}) error {
	if wishlistID, ok := payload["OrderID"].(string); !ok || wishlistID == "" {
		return errors.Wrap(errors.ErrBadRequest, "OrderID is required and cannot be empty")
	}
	if productID, ok := payload["ProductID"].(string); !ok || productID == "" {
		return errors.Wrap(errors.ErrBadRequest, "ProductID is required and cannot be empty")
	}
	return nil
}
func validateAddNotificationAlertPayload(payload map[string]interface{}) error {
	if userID, ok := payload["UserID"].(string); !ok || userID == "" {
		return errors.Wrap(errors.ErrBadRequest, "UserID is required and cannot be empty")
	}
	return nil
}
func validateAddProductAlertPayload(payload map[string]interface{}) error {
	if productID, ok := payload["ProductID"].(string); !ok || productID == "" {
		return errors.Wrap(errors.ErrBadRequest, "ProductID is required and cannot be empty")
	}
	return nil
}
func validateAddMessageAlertPayload(payload map[string]interface{}) error {
	if messageID, ok := payload["MessageID"].(string); !ok || messageID == "" {
		return errors.Wrap(errors.ErrBadRequest, "MessageID is required and cannot be empty")
	}
	if message, ok := payload["Message"].(string); !ok || message == "" {
		return errors.Wrap(errors.ErrBadRequest, "Message is required and cannot be empty")
	}
	if messageSenderID, ok := payload["MessageSenderID"].(string); !ok || messageSenderID == "" {
		return errors.Wrap(errors.ErrBadRequest, "MessageSenderID is required and cannot be empty")
	}
	return nil
}

func validateAddInteractionAlertPayload(payload map[string]interface{}) error {
	if interactionID, ok := payload["InteractionID"].(string); !ok || interactionID == "" {
		return errors.Wrap(errors.ErrBadRequest, "InteractionID is required and cannot be empty")
	}
	if productID, ok := payload["ProductID"].(string); !ok || productID == "" {
		return errors.Wrap(errors.ErrBadRequest, "ProductID is required and cannot be empty")
	}
	if message, ok := payload["Message"].(string); !ok || message == "" {
		return errors.Wrap(errors.ErrBadRequest, "Message is required and cannot be empty")
	}
	return nil
}

func validateAddCommentAlertPayload(payload map[string]interface{}) error {
	if commentID, ok := payload["CommentID"].(string); !ok || commentID == "" {
		return errors.Wrap(errors.ErrBadRequest, "CommentID is required and cannot be empty")
	}
	if productID, ok := payload["ProductID"].(string); !ok || productID == "" {
		return errors.Wrap(errors.ErrBadRequest, "ProductID is required and cannot be empty")
	}
	if message, ok := payload["Message"].(string); !ok || message == "" {
		return errors.Wrap(errors.ErrBadRequest, "Message is required and cannot be empty")
	}
	return nil
}

func validateAddReviewAlertPayload(payload map[string]interface{}) error {
	if reviewID, ok := payload["ReviewID"].(string); !ok || reviewID == "" {
		return errors.Wrap(errors.ErrBadRequest, "ReviewID is required and cannot be empty")
	}
	if productID, ok := payload["ProductID"].(string); !ok || productID == "" {
		return errors.Wrap(errors.ErrBadRequest, "ProductID is required and cannot be empty")
	}
	if message, ok := payload["Message"].(string); !ok || message == "" {
		return errors.Wrap(errors.ErrBadRequest, "Message is required and cannot be empty")
	}
	return nil
}

func validateAddPaymentAlertPayload(payload map[string]interface{}) error {
	if paymentID, ok := payload["PaymentID"].(string); !ok || paymentID == "" {
		return errors.Wrap(errors.ErrBadRequest, "PaymentID is required and cannot be empty")
	}
	if orderID, ok := payload["OrderID"].(string); !ok || orderID == "" {
		return errors.Wrap(errors.ErrBadRequest, "OrderID is required and cannot be empty")
	}
	if message, ok := payload["Message"].(string); !ok || message == "" {
		return errors.Wrap(errors.ErrBadRequest, "Message is required and cannot be empty")
	}
	return nil
}

func validateAddFollowingAlertPayload(payload map[string]interface{}) error {
	if followerID, ok := payload["FollowerID"].(string); !ok || followerID == "" {
		return errors.Wrap(errors.ErrBadRequest, "FollowerID is required and cannot be empty")
	}
	if message, ok := payload["Message"].(string); !ok || message == "" {
		return errors.Wrap(errors.ErrBadRequest, "Message is required and cannot be empty")
	}
	return nil
}

func (n *Alert) AddReviewAlert(userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddReviewAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	reviewID, _ := payload["ReviewID"].(string)
	productID, _ := payload["ProductID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &ReviewAlertAdded{
		ID:        n.ID(),
		AlertType: string(alertType),
		UserID:    userID,
		ReviewID:  reviewID,
		ProductID: productID,
		Payload:   payload,
		Message:   message,
		IsRead:    false,
	}
	n.AddEvent(ReviewAlertAddedEvent, event)
	return ddd.NewEvent(ReviewAlertAddedEvent, n), nil
}

func (n *Alert) AddPaymentAlert(userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddPaymentAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	paymentID, _ := payload["PaymentID"].(string)
	orderID, _ := payload["OrderID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	event := &PaymentAlertAdded{
		ID:        n.ID(),
		AlertType: string(alertType),
		UserID:    userID,
		PaymentID: paymentID,
		OrderID:   orderID,
		Payload:   payload,
		Message:   message,
		IsRead:    false,
	}
	n.AddEvent(PaymentAlertAddedEvent, event)
	return ddd.NewEvent(PaymentAlertAddedEvent, n), nil
}

func (n *Alert) AddFollowingAlert(userID string, alertType AlertType, payload map[string]interface{}) (ddd.Event, error) {
	if n.AlertType != UnknownType {
		return nil, ErrAlertAlreadyCreated
	}

	if len(payload) == 0 {
		return nil, ErrAlertHasNoPayload
	}

	if userID == "" {
		return nil, ErrUserIDCannotBeBlank
	}

	// Notification common notification data
	n.UserID = userID
	n.AlertType = alertType
	n.Payload = payload

	err := validateAddFollowingAlertPayload(payload)
	if err != nil {
		return nil, err
	}
	// Extract values safely with defaults
	followerID, _ := payload["FollowerID"].(string)
	message, _ := payload["Message"].(string)
	
	// Set message in alert
	n.Message = message
	
	n.AddEvent(FollowingAlertAddedEvent, &FollowingAlertAdded{
		ID:         n.ID(),
		AlertType:  string(alertType),
		UserID:     userID,
		FollowerID: followerID,
		Payload:    payload,
		Message:    message,
		IsRead:     false,
	})
	return ddd.NewEvent(FollowingAlertAddedEvent, n), nil
}

// Snapshot and EventApplier methods implementation (assuming you're using Event Sourcing)

func (n *Alert) ApplyEvent(event ddd.Event) error {
	switch e := event.Payload().(type) {
	case *ProductAlertAdded:
		n.UserID = e.UserID
		n.AlertType = AlertType(e.AlertType)
		n.Message = e.Message
		n.Payload = e.Payload
	case *BasketAlertAdded:
		n.UserID = e.UserID
		n.AlertType = BasketType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			"BasketID":  e.BasketID,
			"ProductID": e.ProductID,
			"Message":   e.Message,
		}
	case *WishlistAlertAdded:
		n.UserID = e.UserID
		n.AlertType = WishlistType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			"WishlistID": e.WishlistID,
			"ProductID":  e.ProductID,
			"Message":    e.Message,
		}
	case *MessageAlertAdded:
		n.UserID = e.UserID
		n.AlertType = MessageType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			"MessageID":       e.MessageID,
			"Message":         e.Message,
			"MessageSenderID": e.SenderID,
		}
	case *InteractionAlertAdded:
		n.UserID = e.UserID
		n.AlertType = InteractionType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			//TODO FIX THIS
			"InteractionID": e.UserID,
			"ProductID":     e.ProductID,
			"Message":       e.Message,
		}
	case *UserAlertAdded:
		n.UserID = e.UserID
		n.AlertType = AlertType(e.AlertType)
		n.Message = e.Message
		n.Payload = e.Payload
	case *OrderAlertAdded:
		n.UserID = e.UserID
		n.AlertType = OrderType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			"OrderID":   e.OrderID,
			"ProductID": e.ProductID,
			"Message":   e.Message,
		}
	case *OfferAlertAdded:
		n.UserID = e.UserID
		n.AlertType = OfferType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			"OfferID":   e.OfferID,
			"ProductID": e.ProductID,
			"Message":   e.Message,
		}
	case *SupportAlertAdded:
		n.UserID = e.UserID
		n.AlertType = SupportType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			"TicketID": e.TicketID,
			"Message":  e.Message,
		}
	case *CommentAlertAdded:
		n.UserID = e.UserID
		n.AlertType = CommentType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			"CommentID": e.CommentID,
			"ProductID": e.ProductID,
			"Message":   e.Message,
		}
	case *AlertRead:
		n.IsRead = true
	case *ReviewAlertAdded:
		n.UserID = e.UserID
		n.AlertType = ReviewType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			"ReviewID":  e.ReviewID,
			"ProductID": e.ProductID,
			"Message":   e.Message,
		}
	case *PaymentAlertAdded:
		n.UserID = e.UserID
		n.AlertType = PaymentType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			"PaymentID": e.PaymentID,
			"OrderID":   e.OrderID,
			"Message":   e.Message,
		}
	case *FollowingAlertAdded:
		n.UserID = e.UserID
		n.AlertType = FollowingType
		n.Message = e.Message
		n.Payload = map[string]interface{}{
			"FollowerID": e.FollowerID,
			"Message":    e.Message,
		}

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", n, event.EventName(), e)
	}
	return nil
}

func (n *Alert) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *AlertV1:
		n.UserID = ss.UserID
		n.AlertType = ss.AlertType
		n.Payload = ss.Payload
		n.Message = ss.Message
		n.IsRead = ss.IsRead

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", n, snapshot)
	}

	return nil
}

func (n *Alert) ToSnapshot() es.Snapshot {
	return &AlertV1{
		UserID:    n.UserID,
		AlertType: n.AlertType,
		Payload:   n.Payload,
		Message:   n.Message,
		IsRead:    n.IsRead,
	}
}
