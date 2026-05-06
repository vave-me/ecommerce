package domain

type TOTPSetupDetails struct {
	Secret      string // Base32 encoded secret for manual entry
	QRCodeURI   string // otpauth:// URI
	QRCodeImage string // Optional: base64 encoded QR code image data
	// TempSecret    string // Optional: Unencrypted secret to hold temporarily if not storing before verification
}
