import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useQuery } from '@tanstack/react-query';
import { Package, MapPin, Clock, CheckCircle, Truck, AlertCircle, Search } from 'lucide-react';
import { trackShipment, getShipmentHistory } from '@/api/client/shippingApi';
import LoadingSpinner from '@/components/Utils/LoadingSpinner';
import styles from './ShipmentTracking.module.css';

const TrackingTimeline = ({ events }) => {
  const t = useTranslations('ShipmentTracking');
  
  if (!events || events.length === 0) {
    return (
      <div className={styles.noEvents}>
        <AlertCircle size={48} />
        <p>{t('noEventsFound', { defaultValue: 'No tracking events found' })}</p>
      </div>
    );
  }

  return (
    <div className={styles.timeline}>
      {events.map((event, index) => (
        <div key={event.id || index} className={styles.timelineItem}>
          <div className={styles.timelineDot}>
            {index === 0 && <CheckCircle size={16} />}
          </div>
          <div className={styles.timelineContent}>
            <h4 className={styles.eventTitle}>{event.description}</h4>
            <p className={styles.eventLocation}>
              <MapPin size={14} />
              {event.location}
            </p>
            <p className={styles.eventTime}>
              <Clock size={14} />
              {new Date(event.timestamp).toLocaleString()}
            </p>
          </div>
        </div>
      ))}
    </div>
  );
};

const ShipmentStatus = ({ status }) => {
  const statusConfig = {
    pending: { color: '#f59e0b', bg: 'rgba(245, 158, 11, 0.1)', text: 'Pending', icon: Clock },
    processing: { color: '#3b82f6', bg: 'rgba(59, 130, 246, 0.1)', text: 'Processing', icon: Package },
    shipped: { color: '#8b5cf6', bg: 'rgba(139, 92, 246, 0.1)', text: 'Shipped', icon: Truck },
    in_transit: { color: '#06b6d4', bg: 'rgba(6, 182, 212, 0.1)', text: 'In Transit', icon: Truck },
    out_for_delivery: { color: '#f97316', bg: 'rgba(249, 115, 22, 0.1)', text: 'Out for Delivery', icon: MapPin },
    delivered: { color: '#22c55e', bg: 'rgba(34, 197, 94, 0.1)', text: 'Delivered', icon: CheckCircle },
    failed: { color: '#ef4444', bg: 'rgba(239, 68, 68, 0.1)', text: 'Failed', icon: AlertCircle },
    returned: { color: '#64748b', bg: 'rgba(100, 116, 139, 0.1)', text: 'Returned', icon: Package }
  };

  const config = statusConfig[status] || statusConfig.pending;
  const Icon = config.icon;

  return (
    <div 
      className={styles.statusBadge}
      style={{ color: config.color, backgroundColor: config.bg }}
    >
      <Icon size={20} />
      <span>{config.text}</span>
    </div>
  );
};

const ShipmentTracking = ({ initialTrackingNumber = '' }) => {
  const t = useTranslations('ShipmentTracking');
  const [trackingNumber, setTrackingNumber] = useState(initialTrackingNumber);
  const [searchInput, setSearchInput] = useState(initialTrackingNumber);

  const { 
    data: trackingData, 
    isLoading, 
    error, 
    refetch 
  } = useQuery({
    queryKey: ['shipmentTracking', trackingNumber],
    queryFn: () => trackShipment(trackingNumber),
    enabled: !!trackingNumber,
    retry: 1
  });

  const { 
    data: historyData,
    isLoading: historyLoading 
  } = useQuery({
    queryKey: ['shipmentHistory', trackingData?.shipping?.id],
    queryFn: () => getShipmentHistory(trackingData.shipping.id),
    enabled: !!trackingData?.shipping?.id
  });

  const handleSearch = (e) => {
    e.preventDefault();
    if (searchInput.trim()) {
      setTrackingNumber(searchInput.trim());
    }
  };

  const shipment = trackingData?.shipping;
  const events = historyData?.events || [];

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <h2 className={styles.title}>{t('title', { defaultValue: 'Track Your Shipment' })}</h2>
        <p className={styles.subtitle}>
          {t('subtitle', { defaultValue: 'Enter your tracking number to see the latest status' })}
        </p>
      </div>

      <form onSubmit={handleSearch} className={styles.searchForm}>
        <div className={styles.searchInput}>
          <input
            type="text"
            value={searchInput}
            onChange={(e) => setSearchInput(e.target.value)}
            placeholder={t('trackingPlaceholder', { defaultValue: 'Enter tracking number' })}
            className={styles.input}
          />
          <button type="submit" className={styles.searchButton}>
            <Search size={20} />
            {t('track', { defaultValue: 'Track' })}
          </button>
        </div>
      </form>

      {isLoading && (
        <div className={styles.loading}>
          <LoadingSpinner />
          <p>{t('loading', { defaultValue: 'Searching for your shipment...' })}</p>
        </div>
      )}

      {error && (
        <div className={styles.error}>
          <AlertCircle size={24} />
          <p>{t('notFound', { defaultValue: 'Shipment not found. Please check your tracking number.' })}</p>
        </div>
      )}

      {shipment && (
        <div className={styles.trackingInfo}>
          <div className={styles.shipmentHeader}>
            <div className={styles.trackingDetails}>
              <h3>{t('trackingNumber', { defaultValue: 'Tracking Number' })}</h3>
              <p className={styles.trackingNumber}>{shipment.trackingNumber}</p>
            </div>
            <ShipmentStatus status={shipment.status} />
          </div>

          <div className={styles.shipmentDetails}>
            <div className={styles.detailSection}>
              <h4>{t('from', { defaultValue: 'From' })}</h4>
              <p className={styles.name}>{shipment.senderName}</p>
              <p className={styles.address}>{shipment.senderAddress}</p>
            </div>
            <div className={styles.detailSection}>
              <h4>{t('to', { defaultValue: 'To' })}</h4>
              <p className={styles.name}>{shipment.receiverName}</p>
              <p className={styles.address}>{shipment.receiverAddress}</p>
            </div>
          </div>

          <div className={styles.shipmentMeta}>
            <div className={styles.metaItem}>
              <span className={styles.metaLabel}>{t('service', { defaultValue: 'Service' })}</span>
              <span className={styles.metaValue}>{shipment.serviceType}</span>
            </div>
            <div className={styles.metaItem}>
              <span className={styles.metaLabel}>{t('carrier', { defaultValue: 'Carrier' })}</span>
              <span className={styles.metaValue}>{shipment.carrierName || 'N/A'}</span>
            </div>
            <div className={styles.metaItem}>
              <span className={styles.metaLabel}>{t('weight', { defaultValue: 'Weight' })}</span>
              <span className={styles.metaValue}>{shipment.weight}</span>
            </div>
          </div>

          {shipment.deliveredAt && (
            <div className={styles.deliveryInfo}>
              <CheckCircle size={24} />
              <div>
                <h4>{t('delivered', { defaultValue: 'Delivered' })}</h4>
                <p>{new Date(shipment.deliveredAt).toLocaleString()}</p>
              </div>
            </div>
          )}

          <div className={styles.trackingEvents}>
            <h3>{t('trackingHistory', { defaultValue: 'Tracking History' })}</h3>
            {historyLoading ? (
              <LoadingSpinner />
            ) : (
              <TrackingTimeline events={events} />
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default ShipmentTracking;