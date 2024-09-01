// Package xero provides a REST API client for xero that sends and receives JSON data.
package xero

import "encoding/json"

// Attributes is a struct that represents Cell attributes.
// https://developer.xero.com/documentation/api/accounting/reports#balance-sheet
type Attributes struct {
	ID    string `json:"ID"`
	Value string `json:"Value"`
}

// Cell is a struct that represents Row cells.
// https://developer.xero.com/documentation/api/accounting/reports#balance-sheet
type Cell struct {
	Attributes []Attributes `description:"Cell attributes map" json:"Attributes,omitempty"`
	Value      string       `description:"Cell value" json:"Value"`
}

// ReportResponse is a struct that represents Xero's GET BalanceSheet API response.
// https://developer.xero.com/documentation/api/accounting/reports#balance-sheet
type ReportResponse struct {
	Reports []Report `description:"Reports list" json:"Reports"`
}

// Report is a struct that represents a Xero Report.
// https://developer.xero.com/documentation/api/accounting/reports#balance-sheet
type Report struct {
	Fields         []json.RawMessage `description:"Report fields, not used in this assessment" json:"Fields,omitempty"`
	ReportID       string            `description:"Report UUID" json:"ReportID"`
	ReportName     string            `description:"Report human-readable label" json:"ReportName"`
	ReportType     string            `description:"Report type (BalanceSheet, SalesTaxReturn, ProfitAndLoss, ...)" json:"ReportType"`
	ReportTitles   []string          `description:"List of titles for usage with breadcrumbs" json:"ReportTitles"`
	ReportDate     string            `description:"Report human-readable date (25 August 2024)" json:"ReportDate"`
	Rows           []Row             `description:"Report rows" json:"Rows"`
	UpdatedDateUTC DateTimeField     `description:"Report last update timestamp" json:"UpdatedDateUTC,omitempty"`
}

// Row is a struct that represents a Row.
// https://developer.xero.com/documentation/api/accounting/reports#balance-sheet
type Row struct {
	Cells   []Cell `description:"Row cells" json:"Cells"`
	RowType string `description:"Row type (Header, Row, Section, SummaryRow)" json:"RowType"`
	Rows    []Row  `description:"Section children (only if the RowType is Section)" json:"Rows,omitempty"`
	Title   string `description:"Section title (only if the RowType is Section)" json:"Title"`
}
