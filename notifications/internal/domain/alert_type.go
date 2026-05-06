package domain

import "github.com/stackus/errors"

type AlertType string

const (
	UnknownType     AlertType = ""
	BasketType      AlertType = "basket"
	ProductType     AlertType = "product"
	UserType        AlertType = "user"
	WishlistType    AlertType = "wishlist"
	InteractionType AlertType = "interaction"
	MessageType     AlertType = "message"
	CommentType     AlertType = "comment"
	OfferType       AlertType = "offer"
	SupportType     AlertType = "support"
	OrderType       AlertType = "order"
	ReviewType      AlertType = "review"
	PaymentType     AlertType = "payment"
	FollowingType   AlertType = "following"
)

func (t AlertType) String() string {
	switch t {
	case BasketType, ProductType, UserType, WishlistType, InteractionType, MessageType, CommentType, OfferType, SupportType, OrderType, ReviewType, PaymentType, FollowingType:
		return string(t)
	default:
		return ""
	}
}

func ToAlertType(s string) (AlertType, error) {
	switch s {
	case string(BasketType):
		return BasketType, nil
	case string(ProductType):
		return ProductType, nil
	case string(UserType):
		return UserType, nil
	case string(WishlistType):
		return WishlistType, nil
	case string(InteractionType):
		return InteractionType, nil
	case string(MessageType):
		return MessageType, nil
	case string(CommentType):
		return CommentType, nil
	case string(OfferType):
		return OfferType, nil
	case string(SupportType):
		return SupportType, nil
	case string(OrderType):
		return OrderType, nil
	case string(ReviewType):
		return ReviewType, nil
	case string(PaymentType):
		return PaymentType, nil
	case string(FollowingType):
		return FollowingType, nil
	default:
		return UnknownType, errors.Wrap(errors.ErrBadRequest, "invalid alert type: "+s)
	}
}
