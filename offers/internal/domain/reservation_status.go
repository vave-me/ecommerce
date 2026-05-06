package domain // ReservationStatus
type ReservationStatus string

const (
	ReservationStatusUnknown             ReservationStatus = ""
	ReservationStatusActive              ReservationStatus = "ACTIVE"
	ReservationStatusRedeemed            ReservationStatus = "REDEEMED"
	ReservationStatusExpired             ReservationStatus = "EXPIRED"
	ReservationStatusCanceled            ReservationStatus = "CANCELED"
	ReservationStatusNegotiationPending  ReservationStatus = "NEGOTIATION_PENDING"
	ReservationStatusNegotiationAgreed   ReservationStatus = "NEGOTIATION_AGREED"
	ReservationStatusNegotiationDeclined ReservationStatus = "NEGOTIATION_DECLINED"
)

func (s ReservationStatus) String() string {
	switch s {
	case ReservationStatusActive, ReservationStatusRedeemed, ReservationStatusExpired, ReservationStatusCanceled, ReservationStatusNegotiationPending, ReservationStatusNegotiationAgreed, ReservationStatusNegotiationDeclined:
		return string(s)
	default:
		return ""
	}
}
