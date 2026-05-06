package grpc

import (
	"context"
	"database/sql"

	"google.golang.org/grpc"

	"middleman/erp/erppb"
	"middleman/erp/internal/application"
	"middleman/erp/internal/constants"
	"middleman/internal/di"
)

type serverTx struct {
	c di.Container
	erppb.UnimplementedERPServiceServer
}

var _ erppb.ERPServiceServer = (*serverTx)(nil)

// RegisterServerTx registers the transactional gRPC server implementation
func RegisterServerTx(container di.Container, registrar grpc.ServiceRegistrar) error {
	erppb.RegisterERPServiceServer(registrar, serverTx{
		c: container,
	})
	return nil
}

// Connector Management with Transaction Support

// AddConnector creates a new ERP connector with transaction support
func (s serverTx) AddConnector(ctx context.Context, request *erppb.AddConnectorRequest) (resp *erppb.AddConnectorResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.AddConnector(ctx, request)
}

// UpdateConnector updates an existing ERP connector with transaction support
func (s serverTx) UpdateConnector(ctx context.Context, request *erppb.UpdateConnectorRequest) (resp *erppb.UpdateConnectorResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateConnector(ctx, request)
}

// RemoveConnector removes an ERP connector with transaction support
func (s serverTx) RemoveConnector(ctx context.Context, request *erppb.RemoveConnectorRequest) (resp *erppb.RemoveConnectorResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.RemoveConnector(ctx, request)
}

// ToggleConnector activates or deactivates a connector with transaction support
func (s serverTx) ToggleConnector(ctx context.Context, request *erppb.ToggleConnectorRequest) (resp *erppb.ToggleConnectorResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ToggleConnector(ctx, request)
}

// Invoice Granular Commands

// CreateInvoice creates a new invoice with transaction support
func (s serverTx) CreateInvoice(ctx context.Context, request *erppb.CreateInvoiceRequest) (resp *erppb.CreateInvoiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.CreateInvoice(ctx, request)
}

// ApproveInvoice approves an invoice with transaction support
func (s serverTx) ApproveInvoice(ctx context.Context, request *erppb.ApproveInvoiceRequest) (resp *erppb.ApproveInvoiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ApproveInvoice(ctx, request)
}

// SendInvoice sends an invoice with transaction support
func (s serverTx) SendInvoice(ctx context.Context, request *erppb.SendInvoiceRequest) (resp *erppb.SendInvoiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.SendInvoice(ctx, request)
}

// VoidInvoice voids an invoice with transaction support
func (s serverTx) VoidInvoice(ctx context.Context, request *erppb.VoidInvoiceRequest) (resp *erppb.VoidInvoiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.VoidInvoice(ctx, request)
}

// RecordInvoicePayment records a payment against an invoice with transaction support
func (s serverTx) RecordInvoicePayment(ctx context.Context, request *erppb.RecordInvoicePaymentRequest) (resp *erppb.RecordInvoicePaymentResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.RecordInvoicePayment(ctx, request)
}

// ProcessInvoice processes an invoice in the ERP system with transaction support (DEPRECATED)
func (s serverTx) ProcessInvoice(ctx context.Context, request *erppb.ProcessInvoiceRequest) (resp *erppb.ProcessInvoiceResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ProcessInvoice(ctx, request)
}

// ProcessReturn processes a return in the ERP system with transaction support
func (s serverTx) ProcessReturn(ctx context.Context, request *erppb.ProcessReturnRequest) (resp *erppb.ProcessReturnResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ProcessReturn(ctx, request)
}

// ProcessWebhook processes incoming webhooks from ERP systems with transaction support
func (s serverTx) ProcessWebhook(ctx context.Context, request *erppb.ProcessWebhookRequest) (resp *erppb.ProcessWebhookResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ProcessWebhook(ctx, request)
}

// RegisterConnector registers a new ERP connector with transaction support
func (s serverTx) RegisterConnector(ctx context.Context, request *erppb.RegisterConnectorRequest) (resp *erppb.RegisterConnectorResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.RegisterConnector(ctx, request)
}

// SendOrder sends an order to the ERP system with transaction support
func (s serverTx) SendOrder(ctx context.Context, request *erppb.SendOrderRequest) (resp *erppb.SendOrderResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.SendOrder(ctx, request)
}

// SyncCustomers synchronizes customers from the ERP system with transaction support
func (s serverTx) SyncCustomers(ctx context.Context, request *erppb.SyncCustomersRequest) (resp *erppb.SyncCustomersResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.SyncCustomers(ctx, request)
}

// SyncPrices synchronizes prices from the ERP system with transaction support
func (s serverTx) SyncPrices(ctx context.Context, request *erppb.SyncPricesRequest) (resp *erppb.SyncPricesResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.SyncPrices(ctx, request)
}

// SyncProducts synchronizes products from the ERP system with transaction support
func (s serverTx) SyncProducts(ctx context.Context, request *erppb.SyncProductsRequest) (resp *erppb.SyncProductsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.SyncProducts(ctx, request)
}

// SyncStock synchronizes stock levels from the ERP system with transaction support
func (s serverTx) SyncStock(ctx context.Context, request *erppb.SyncStockRequest) (resp *erppb.SyncStockResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.SyncStock(ctx, request)
}

// UpdateInventoryReservation updates inventory reservations in the ERP system with transaction support
func (s serverTx) UpdateInventoryReservation(ctx context.Context, request *erppb.UpdateInventoryReservationRequest) (resp *erppb.UpdateInventoryReservationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.UpdateInventoryReservation(ctx, request)
}

// GetConnectorStatus retrieves the status of a specific connector
func (s serverTx) GetConnectorStatus(ctx context.Context, request *erppb.GetConnectorStatusRequest) (resp *erppb.GetConnectorStatusResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.GetConnectorStatus(ctx, request)
}

// GetSyncHistory retrieves the synchronization history for a connector
func (s serverTx) GetSyncHistory(ctx context.Context, request *erppb.GetSyncHistoryRequest) (resp *erppb.GetSyncHistoryResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.GetSyncHistory(ctx, request)
}

// ListConnectors lists all registered connectors
func (s serverTx) ListConnectors(ctx context.Context, request *erppb.ListConnectorsRequest) (resp *erppb.ListConnectorsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ListConnectors(ctx, request)
}

// Return Granular Commands

// CreateReturn creates a new return request with transaction support
func (s serverTx) CreateReturn(ctx context.Context, request *erppb.CreateReturnRequest) (resp *erppb.CreateReturnResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.CreateReturn(ctx, request)
}

// ApproveReturn approves a return request with transaction support
func (s serverTx) ApproveReturn(ctx context.Context, request *erppb.ApproveReturnRequest) (resp *erppb.ApproveReturnResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ApproveReturn(ctx, request)
}

// RejectReturn rejects a return request with transaction support
func (s serverTx) RejectReturn(ctx context.Context, request *erppb.RejectReturnRequest) (resp *erppb.RejectReturnResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.RejectReturn(ctx, request)
}

// ProcessReturnStart starts processing a return with transaction support
func (s serverTx) ProcessReturnStart(ctx context.Context, request *erppb.ProcessReturnStartRequest) (resp *erppb.ProcessReturnStartResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ProcessReturnStart(ctx, request)
}

// CompleteReturn completes a return with transaction support
func (s serverTx) CompleteReturn(ctx context.Context, request *erppb.CompleteReturnRequest) (resp *erppb.CompleteReturnResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.CompleteReturn(ctx, request)
}

// RestockReturnItems records restocked items with transaction support
func (s serverTx) RestockReturnItems(ctx context.Context, request *erppb.RestockReturnItemsRequest) (resp *erppb.RestockReturnItemsResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.RestockReturnItems(ctx, request)
}

// SyncReturnToERP syncs a return to ERP with transaction support
func (s serverTx) SyncReturnToERP(ctx context.Context, request *erppb.SyncReturnToERPRequest) (resp *erppb.SyncReturnToERPResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.SyncReturnToERP(ctx, request)
}

// Inventory Reservation Granular Commands

// CreateInventoryReservation creates a new inventory reservation with transaction support
func (s serverTx) CreateInventoryReservation(ctx context.Context, request *erppb.CreateInventoryReservationRequest) (resp *erppb.CreateInventoryReservationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.CreateInventoryReservation(ctx, request)
}

// ReleaseInventoryReservation releases an inventory reservation with transaction support
func (s serverTx) ReleaseInventoryReservation(ctx context.Context, request *erppb.ReleaseInventoryReservationRequest) (resp *erppb.ReleaseInventoryReservationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.ReleaseInventoryReservation(ctx, request)
}

// FulfillInventoryReservation marks a reservation as fulfilled with transaction support
func (s serverTx) FulfillInventoryReservation(ctx context.Context, request *erppb.FulfillInventoryReservationRequest) (resp *erppb.FulfillInventoryReservationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.FulfillInventoryReservation(ctx, request)
}

// TransferInventoryReservation transfers a reservation to another warehouse with transaction support
func (s serverTx) TransferInventoryReservation(ctx context.Context, request *erppb.TransferInventoryReservationRequest) (resp *erppb.TransferInventoryReservationResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.TransferInventoryReservation(ctx, request)
}

// SyncInvoiceToERP syncs an invoice to ERP with transaction support
func (s serverTx) SyncInvoiceToERP(ctx context.Context, request *erppb.SyncInvoiceToERPRequest) (resp *erppb.SyncInvoiceToERPResponse, err error) {
	ctx = s.c.Scoped(ctx)
	defer func(tx *sql.Tx) {
		err = s.closeTx(tx, err)
	}(di.Get(ctx, constants.DatabaseTransactionKey).(*sql.Tx))

	next := server{app: di.Get(ctx, constants.ApplicationKey).(application.App)}
	return next.SyncInvoiceToERP(ctx, request)
}

// closeTx handles transaction completion (commit or rollback)
func (s serverTx) closeTx(tx *sql.Tx, err error) error {
	if p := recover(); p != nil {
		_ = tx.Rollback()
		panic(p)
	} else if err != nil {
		_ = tx.Rollback()
		return err
	} else {
		return tx.Commit()
	}
}