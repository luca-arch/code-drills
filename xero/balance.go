/*
Package xero provides a REST API client for xero that sends and receives JSON data.
*/
package xero

import "encoding/json"

type Cell struct {
	Attributes []struct {
		ID    string `json:"ID"`
		Value string `json:"Value"`
	} `description:"Cell attributes map" json:"Attributes,omitempty"`
	Value string `description:"Cell value" json:"Value"`
}

type ReportResponse struct {
	Reports []Report `description:"Reports list" json:"Reports"`
}

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

type Row struct {
	Cells   []Cell `description:"Row cells" json:"Cells"`
	RowType string `description:"Row type (Header, Row, Section, SummaryRow)" json:"RowType"`
	Rows    []Row  `description:"Section children (only if the RowType is Section)" json:"Rows,omitempty"`
	Title   string `description:"Section title (only if the RowType is Section)" json:"Title"`
}
