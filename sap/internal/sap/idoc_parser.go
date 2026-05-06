package sap

import (
	"encoding/xml"
	"fmt"
	"strconv"
	"time"
)

// IDocType represents the type of IDoc
type IDocType string

const (
	IDocTypeMATMAS IDocType = "MATMAS" // Material Master
	IDocTypeCOND_A IDocType = "COND_A" // Pricing Conditions
	IDocTypeINVCON IDocType = "INVCON" // Inventory Control
	IDocTypeORDERS IDocType = "ORDERS" // Purchase Orders
)

// IDoc represents a SAP Intermediate Document
type IDoc struct {
	XMLName xml.Name `xml:"IDOC"`
	Begin   string   `xml:"BEGIN,attr"`
	EDI_DC40 struct {
		IDOCTYP string `xml:"IDOCTYP"` // IDoc type
		MESTYP  string `xml:"MESTYP"`  // Message type
		SNDPOR  string `xml:"SNDPOR"`  // Sender port
		RCVPOR  string `xml:"RCVPOR"`  // Receiver port
		CREDAT  string `xml:"CREDAT"`  // Creation date
		CRETIM  string `xml:"CRETIM"`  // Creation time
		DOCNUM  string `xml:"DOCNUM"`  // IDoc number
	} `xml:"EDI_DC40"`
	Data interface{} `xml:",any"`
}

// MATMAS represents Material Master IDoc
type MATMAS struct {
	E1MARAM struct {
		MATNR      string `xml:"MATNR"`      // Material number
		MAKTX      string `xml:"MAKTX"`      // Material description
		MAKTX_LONG string `xml:"MAKTX_LONG"` // Long description
		MEINS      string `xml:"MEINS"`      // Base unit of measure
		MTART      string `xml:"MTART"`      // Material type
		MATKL      string `xml:"MATKL"`      // Material group
		BRGEW      string `xml:"BRGEW"`      // Gross weight
		NTGEW      string `xml:"NTGEW"`      // Net weight
		GEWEI      string `xml:"GEWEI"`      // Weight unit
		VOLUM      string `xml:"VOLUM"`      // Volume
		VOLEH      string `xml:"VOLEH"`      // Volume unit
		LAENG      string `xml:"LAENG"`      // Length
		BREIT      string `xml:"BREIT"`      // Width
		HOEHE      string `xml:"HOEHE"`      // Height
		MEABM      string `xml:"MEABM"`      // Dimension unit
		SPART      string `xml:"SPART"`      // Division
		AEDAT      string `xml:"AEDAT"`      // Changed on
		AETIM      string `xml:"AETIM"`      // Time of change
	} `xml:"E1MARAM"`
	E1MAKTM []struct {
		SPRAS string `xml:"SPRAS"` // Language key
		MAKTX string `xml:"MAKTX"` // Material description
	} `xml:"E1MAKTM"`
}

// INVCON represents Inventory Control IDoc
type INVCON struct {
	E1INVCO struct {
		MATNR string `xml:"MATNR"` // Material number
		WERKS string `xml:"WERKS"` // Plant
		LGORT string `xml:"LGORT"` // Storage location
		SOBKZ string `xml:"SOBKZ"` // Special stock indicator
		MENGE string `xml:"MENGE"` // Quantity
		MEINS string `xml:"MEINS"` // Unit of measure
		AEDAT string `xml:"AEDAT"` // Changed on
		AETIM string `xml:"AETIM"` // Time of change
	} `xml:"E1INVCO"`
}

// COND_A represents Pricing Condition IDoc
type COND_A struct {
	E1KOMG struct {
		KVEWE string `xml:"KVEWE"` // Usage of condition table
		KOTABNR string `xml:"KOTABNR"` // Condition table
		KAPPL string `xml:"KAPPL"` // Application
		KSCHL string `xml:"KSCHL"` // Condition type
		MATNR string `xml:"MATNR"` // Material number
		VKORG string `xml:"VKORG"` // Sales organization
		VTWEG string `xml:"VTWEG"` // Distribution channel
	} `xml:"E1KOMG"`
	E1KONH struct {
		DATAB string `xml:"DATAB"` // Valid from
		DATBI string `xml:"DATBI"` // Valid to
		KOSRT string `xml:"KOSRT"` // Sort sequence
	} `xml:"E1KONH"`
	E1KONP struct {
		KBETR string `xml:"KBETR"` // Rate
		KONWA string `xml:"KONWA"` // Rate unit
		KPEIN string `xml:"KPEIN"` // Pricing unit
		KMEIN string `xml:"KMEIN"` // Unit of measure
		WAERS string `xml:"WAERS"` // Currency
		MWSK1 string `xml:"MWSK1"` // Tax code
	} `xml:"E1KONP"`
}

// ParseIDoc parses a raw IDoc XML payload
func ParseIDoc(payload []byte) (*IDoc, error) {
	var idoc IDoc
	if err := xml.Unmarshal(payload, &idoc); err != nil {
		return nil, fmt.Errorf("parsing IDoc: %w", err)
	}
	return &idoc, nil
}

// ParseMATMAS parses a MATMAS IDoc
func ParseMATMAS(payload []byte) (*MATMAS, error) {
	var matmas MATMAS
	if err := xml.Unmarshal(payload, &matmas); err != nil {
		return nil, fmt.Errorf("parsing MATMAS: %w", err)
	}
	return &matmas, nil
}

// ParseINVCON parses an INVCON IDoc
func ParseINVCON(payload []byte) (*INVCON, error) {
	var invcon INVCON
	if err := xml.Unmarshal(payload, &invcon); err != nil {
		return nil, fmt.Errorf("parsing INVCON: %w", err)
	}
	return &invcon, nil
}

// ParseCOND_A parses a COND_A IDoc
func ParseCOND_A(payload []byte) (*COND_A, error) {
	var conda COND_A
	if err := xml.Unmarshal(payload, &conda); err != nil {
		return nil, fmt.Errorf("parsing COND_A: %w", err)
	}
	return &conda, nil
}

// Helper functions for parsing SAP data types

// ParseSAPDate parses SAP date format (YYYYMMDD)
func ParseSAPDate(dateStr string) (time.Time, error) {
	if len(dateStr) != 8 {
		return time.Time{}, fmt.Errorf("invalid SAP date format: %s", dateStr)
	}
	return time.Parse("20060102", dateStr)
}

// ParseSAPTime parses SAP time format (HHMMSS)
func ParseSAPTime(timeStr string) (time.Time, error) {
	if len(timeStr) != 6 {
		return time.Time{}, fmt.Errorf("invalid SAP time format: %s", timeStr)
	}
	return time.Parse("150405", timeStr)
}

// ParseSAPDateTime combines SAP date and time
func ParseSAPDateTime(dateStr, timeStr string) (time.Time, error) {
	date, err := ParseSAPDate(dateStr)
	if err != nil {
		return time.Time{}, err
	}
	
	if timeStr == "" {
		return date, nil
	}
	
	t, err := ParseSAPTime(timeStr)
	if err != nil {
		return time.Time{}, err
	}
	
	return time.Date(
		date.Year(), date.Month(), date.Day(),
		t.Hour(), t.Minute(), t.Second(), 0,
		time.UTC,
	), nil
}

// ParseSAPFloat parses SAP decimal format
func ParseSAPFloat(value string) (float64, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.ParseFloat(value, 64)
}

// ParseSAPInt parses SAP integer format
func ParseSAPInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	return strconv.Atoi(value)
}