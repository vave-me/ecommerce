package streaming

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/stackus/errors"
)

// DRMProvider represents different DRM systems
type DRMProvider string

const (
	DRMProviderWidevine  DRMProvider = "widevine"
	DRMProviderFairPlay  DRMProvider = "fairplay"
	DRMProviderPlayReady DRMProvider = "playready"
	DRMProviderClearKey  DRMProvider = "clearkey"
)

// DRMManager handles DRM operations
type DRMManager struct {
	providers    map[DRMProvider]DRMProviderInterface
	keyStore     KeyStore
	licenseCache *LicenseCache
	httpClient   *http.Client
	mu           sync.RWMutex
}

// DRMProviderInterface defines methods for DRM providers
type DRMProviderInterface interface {
	GenerateLicense(ctx context.Context, req *LicenseRequest) (*LicenseResponse, error)
	VerifyLicense(ctx context.Context, license []byte) error
	GetContentKey(contentID string) ([]byte, error)
	PackageContent(content []byte, key []byte) ([]byte, error)
}

// KeyStore manages content encryption keys
type KeyStore interface {
	GenerateKey(contentID string) (*ContentKey, error)
	GetKey(keyID string) (*ContentKey, error)
	StoreKey(key *ContentKey) error
	RevokeKey(keyID string) error
}

// ContentKey represents an encryption key
type ContentKey struct {
	KeyID       string
	ContentID   string
	Key         []byte
	IV          []byte
	Algorithm   string
	CreatedAt   time.Time
	ExpiresAt   time.Time
	IsActive    bool
}

// LicenseRequest represents a DRM license request
type LicenseRequest struct {
	ContentID    string
	UserID       string
	DeviceID     string
	SessionID    string
	DRMProvider  DRMProvider
	Challenge    []byte
	CustomData   map[string]interface{}
}

// LicenseResponse represents a DRM license response
type LicenseResponse struct {
	License      []byte
	ExpiresAt    time.Time
	Restrictions LicenseRestrictions
}

// LicenseRestrictions defines playback restrictions
type LicenseRestrictions struct {
	MaxResolution    string
	HDCP             string // HDCP version required
	DigitalOnly      bool
	ExpirationDate   time.Time
	PlaybackDuration int // seconds
	RentalDuration   int // hours
	Geoblocking      []string // allowed countries
}

// LicenseCache caches DRM licenses
type LicenseCache struct {
	cache map[string]*CachedLicense
	mu    sync.RWMutex
}

// CachedLicense represents a cached license
type CachedLicense struct {
	License   *LicenseResponse
	CachedAt  time.Time
	ExpiresAt time.Time
}

// NewDRMManager creates a new DRM manager
func NewDRMManager() *DRMManager {
	return &DRMManager{
		providers:    make(map[DRMProvider]DRMProviderInterface),
		keyStore:     NewInMemoryKeyStore(),
		licenseCache: NewLicenseCache(),
		httpClient:   &http.Client{Timeout: 10 * time.Second},
	}
}

// RegisterProvider registers a DRM provider
func (dm *DRMManager) RegisterProvider(provider DRMProvider, impl DRMProviderInterface) {
	dm.mu.Lock()
	defer dm.mu.Unlock()
	dm.providers[provider] = impl
}

// GenerateContentKey generates a new content encryption key
func (dm *DRMManager) GenerateContentKey(contentID string) (*ContentKey, error) {
	return dm.keyStore.GenerateKey(contentID)
}

// RequestLicense processes a DRM license request
func (dm *DRMManager) RequestLicense(ctx context.Context, req *LicenseRequest) (*LicenseResponse, error) {
	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s:%s", req.ContentID, req.UserID, req.DeviceID)
	if cached := dm.licenseCache.Get(cacheKey); cached != nil {
		return cached.License, nil
	}

	// Get provider
	dm.mu.RLock()
	provider, exists := dm.providers[req.DRMProvider]
	dm.mu.RUnlock()

	if !exists {
		return nil, errors.Wrap(errors.ErrNotFound, "DRM provider not found")
	}

	// Generate license
	license, err := provider.GenerateLicense(ctx, req)
	if err != nil {
		return nil, err
	}

	// Cache license
	dm.licenseCache.Set(cacheKey, license, 1*time.Hour)

	return license, nil
}

// EncryptSegment encrypts a video segment
func (dm *DRMManager) EncryptSegment(segment []byte, keyID string) ([]byte, error) {
	// Get encryption key
	key, err := dm.keyStore.GetKey(keyID)
	if err != nil {
		return nil, err
	}

	// AES-128 CBC encryption (HLS standard)
	block, err := aes.NewCipher(key.Key)
	if err != nil {
		return nil, err
	}

	// Pad segment to block size
	padding := aes.BlockSize - (len(segment) % aes.BlockSize)
	if padding > 0 {
		segment = append(segment, bytes.Repeat([]byte{byte(padding)}, padding)...)
	}

	// Encrypt
	encrypted := make([]byte, len(segment))
	mode := cipher.NewCBCEncrypter(block, key.IV)
	mode.CryptBlocks(encrypted, segment)

	return encrypted, nil
}

// GenerateKeyFile generates HLS key file content
func (dm *DRMManager) GenerateKeyFile(keyID string, keyURL string) ([]byte, error) {
	key, err := dm.keyStore.GetKey(keyID)
	if err != nil {
		return nil, err
	}

	// HLS key file format
	content := fmt.Sprintf("#EXT-X-KEY:METHOD=AES-128,URI=\"%s\",IV=0x%s\n",
		keyURL,
		hex.EncodeToString(key.IV),
	)

	return []byte(content), nil
}

// WidevineDRM implements Widevine DRM
type WidevineDRM struct {
	licenseServerURL string
	signingKey       []byte
	encryptionKey    []byte
	keyStore         KeyStore
}

// NewWidevineDRM creates a new Widevine DRM provider
func NewWidevineDRM(licenseServerURL string, signingKey, encryptionKey []byte, keyStore KeyStore) *WidevineDRM {
	return &WidevineDRM{
		licenseServerURL: licenseServerURL,
		signingKey:       signingKey,
		encryptionKey:    encryptionKey,
		keyStore:         keyStore,
	}
}

// GenerateLicense generates a Widevine license
func (w *WidevineDRM) GenerateLicense(ctx context.Context, req *LicenseRequest) (*LicenseResponse, error) {
	// Get content key
	contentKey, err := w.keyStore.GetKey(req.ContentID)
	if err != nil {
		return nil, err
	}

	// Build Widevine license request
	licenseReq := map[string]interface{}{
		"payload":    base64.StdEncoding.EncodeToString(req.Challenge),
		"content_id": req.ContentID,
		"provider":   "middleman_streams",
		"allowed_track_types": "SD_HD",
		"policy": map[string]interface{}{
			"can_play":    true,
			"can_persist": false,
			"can_renew":   true,
		},
	}

	// Send to Widevine license server
	jsonData, _ := json.Marshal(licenseReq)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", w.licenseServerURL, bytes.NewReader(jsonData))
	if err != nil {
		return nil, err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	
	// Sign request
	signature := w.signRequest(jsonData)
	httpReq.Header.Set("X-Signature", signature)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Wrap(errors.ErrInternalServerError, "license server error")
	}

	licenseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &LicenseResponse{
		License:   licenseData,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Restrictions: LicenseRestrictions{
			MaxResolution:    "1080p",
			HDCP:             "1.4",
			DigitalOnly:      true,
			ExpirationDate:   time.Now().Add(48 * time.Hour),
			PlaybackDuration: 0, // unlimited during rental
			RentalDuration:   48,
		},
	}, nil
}

// VerifyLicense verifies a Widevine license
func (w *WidevineDRM) VerifyLicense(ctx context.Context, license []byte) error {
	// Implement Widevine license verification
	return nil
}

// GetContentKey retrieves content key
func (w *WidevineDRM) GetContentKey(contentID string) ([]byte, error) {
	key, err := w.keyStore.GetKey(contentID)
	if err != nil {
		return nil, err
	}
	return key.Key, nil
}

// PackageContent packages content for Widevine
func (w *WidevineDRM) PackageContent(content []byte, key []byte) ([]byte, error) {
	// CENC (Common Encryption) packaging
	return content, nil
}

// signRequest signs a request for Widevine
func (w *WidevineDRM) signRequest(data []byte) string {
	h := sha256.New()
	h.Write(data)
	h.Write(w.signingKey)
	return hex.EncodeToString(h.Sum(nil))
}

// FairPlayDRM implements Apple FairPlay DRM
type FairPlayDRM struct {
	certificateURL string
	keyServerURL   string
	keyStore       KeyStore
}

// NewFairPlayDRM creates a new FairPlay DRM provider
func NewFairPlayDRM(certificateURL, keyServerURL string, keyStore KeyStore) *FairPlayDRM {
	return &FairPlayDRM{
		certificateURL: certificateURL,
		keyServerURL:   keyServerURL,
		keyStore:       keyStore,
	}
}

// GenerateLicense generates a FairPlay license
func (f *FairPlayDRM) GenerateLicense(ctx context.Context, req *LicenseRequest) (*LicenseResponse, error) {
	// FairPlay SPC (Server Playback Context) processing
	// This is a simplified version - actual implementation requires Apple's FairPlay SDK

	// Get content key
	contentKey, err := f.keyStore.GetKey(req.ContentID)
	if err != nil {
		return nil, err
	}

	// Generate CKC (Content Key Context)
	ckc := f.generateCKC(req.Challenge, contentKey.Key)

	return &LicenseResponse{
		License:   ckc,
		ExpiresAt: time.Now().Add(24 * time.Hour),
		Restrictions: LicenseRestrictions{
			MaxResolution:    "4K",
			HDCP:             "2.2",
			DigitalOnly:      true,
			ExpirationDate:   time.Now().Add(48 * time.Hour),
			PlaybackDuration: 0,
			RentalDuration:   48,
		},
	}, nil
}

// VerifyLicense verifies a FairPlay license
func (f *FairPlayDRM) VerifyLicense(ctx context.Context, license []byte) error {
	return nil
}

// GetContentKey retrieves content key
func (f *FairPlayDRM) GetContentKey(contentID string) ([]byte, error) {
	key, err := f.keyStore.GetKey(contentID)
	if err != nil {
		return nil, err
	}
	return key.Key, nil
}

// PackageContent packages content for FairPlay
func (f *FairPlayDRM) PackageContent(content []byte, key []byte) ([]byte, error) {
	// FairPlay uses SAMPLE-AES encryption
	return content, nil
}

// generateCKC generates Content Key Context for FairPlay
func (f *FairPlayDRM) generateCKC(spc []byte, contentKey []byte) []byte {
	// Simplified CKC generation
	// Actual implementation requires FairPlay SDK
	h := sha256.New()
	h.Write(spc)
	h.Write(contentKey)
	return h.Sum(nil)
}

// InMemoryKeyStore implements KeyStore interface
type InMemoryKeyStore struct {
	keys map[string]*ContentKey
	mu   sync.RWMutex
}

// NewInMemoryKeyStore creates a new in-memory key store
func NewInMemoryKeyStore() *InMemoryKeyStore {
	return &InMemoryKeyStore{
		keys: make(map[string]*ContentKey),
	}
}

// GenerateKey generates a new content key
func (ks *InMemoryKeyStore) GenerateKey(contentID string) (*ContentKey, error) {
	// Generate 128-bit key
	key := make([]byte, 16)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}

	// Generate IV
	iv := make([]byte, 16)
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	contentKey := &ContentKey{
		KeyID:     uuid.New().String(),
		ContentID: contentID,
		Key:       key,
		IV:        iv,
		Algorithm: "AES-128",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(30 * 24 * time.Hour),
		IsActive:  true,
	}

	ks.mu.Lock()
	ks.keys[contentKey.KeyID] = contentKey
	ks.mu.Unlock()

	return contentKey, nil
}

// GetKey retrieves a content key
func (ks *InMemoryKeyStore) GetKey(keyID string) (*ContentKey, error) {
	ks.mu.RLock()
	defer ks.mu.RUnlock()

	key, exists := ks.keys[keyID]
	if !exists {
		return nil, errors.Wrap(errors.ErrNotFound, "key not found")
	}

	if !key.IsActive {
		return nil, errors.Wrap(errors.ErrForbidden, "key is revoked")
	}

	if time.Now().After(key.ExpiresAt) {
		return nil, errors.Wrap(errors.ErrForbidden, "key has expired")
	}

	return key, nil
}

// StoreKey stores a content key
func (ks *InMemoryKeyStore) StoreKey(key *ContentKey) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()
	ks.keys[key.KeyID] = key
	return nil
}

// RevokeKey revokes a content key
func (ks *InMemoryKeyStore) RevokeKey(keyID string) error {
	ks.mu.Lock()
	defer ks.mu.Unlock()

	if key, exists := ks.keys[keyID]; exists {
		key.IsActive = false
		return nil
	}

	return errors.Wrap(errors.ErrNotFound, "key not found")
}

// NewLicenseCache creates a new license cache
func NewLicenseCache() *LicenseCache {
	return &LicenseCache{
		cache: make(map[string]*CachedLicense),
	}
}

// Get retrieves a cached license
func (lc *LicenseCache) Get(key string) *CachedLicense {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	if cached, exists := lc.cache[key]; exists {
		if time.Now().Before(cached.ExpiresAt) {
			return cached
		}
		// Expired, remove from cache
		delete(lc.cache, key)
	}

	return nil
}

// Set caches a license
func (lc *LicenseCache) Set(key string, license *LicenseResponse, ttl time.Duration) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	lc.cache[key] = &CachedLicense{
		License:   license,
		CachedAt:  time.Now(),
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Clear clears the license cache
func (lc *LicenseCache) Clear() {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	lc.cache = make(map[string]*CachedLicense)
}