// Package invoice provides functionality for generating PDF invoices.
package invoice

import (
	"fmt"
	"invo/pkg/config"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"

	"invo/pkg/invoice/fonts"
)

// Invoice represents a complete invoice with all necessary data.
type Invoice struct {
	Year        int
	Month       int
	Items       []config.Item
	InvoiceDate time.Time
	DueDate     time.Time
	Period      string
	InvoiceNum  string
	Subtotal    float64
	Config      *config.Config
}

// NewInvoice creates a new Invoice instance with calculated fields.
func NewInvoice(invoiceNum string, year, month int, qtys []int, cfg *config.Config, invoiceDate, dueDate string) *Invoice {
	items := normalizeItems(cfg.Items, qtys)

	inv := &Invoice{
		Year:   year,
		Month:  month,
		Items:  items,
		Config: cfg,
	}

	if invoiceDate != "" {
		date, err := time.Parse("2006-01-02", invoiceDate)
		if err == nil {
			inv.InvoiceDate = date
		}
	}
	if dueDate != "" {
		date, err := time.Parse("2006-01-02", dueDate)
		if err == nil {
			inv.DueDate = date
		}
	}

	// Only calculate dates if not provided
	if inv.InvoiceDate.IsZero() || inv.DueDate.IsZero() {
		inv.calculateDates()
	}

	inv.Period = fmt.Sprintf("%02d/%04d", month, year)
	inv.InvoiceNum = generateInvoiceNumber(invoiceNum, cfg.Sender.Name, year, month)
	inv.Subtotal = calculateSubtotal(items)

	return inv
}

// GeneratePDF creates and saves a PDF invoice.
func (i *Invoice) GeneratePDF() (string, error) {
	pdf := i.initializePDF()

	i.addHeader(pdf)
	i.addSenderInfo(pdf)
	i.addRecipientInfo(pdf)
	i.addProjectInfo(pdf)
	i.addItemsTable(pdf)
	i.addPaymentDetails(pdf)
	i.addTotal(pdf)

	filename := i.generateFilename()
	if err := pdf.OutputFileAndClose(filename); err != nil {
		return "", fmt.Errorf("saving PDF: %w", err)
	}

	return filename, nil
}

func (i *Invoice) initializePDF() *fpdf.Fpdf {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 20, 15)
	pdf.AddUTF8FontFromBytes("Roboto", "", fonts.RobotoTTF)
	pdf.AddUTF8FontFromBytes("Roboto", "B", fonts.RobotoTTF)
	pdf.AddPage()
	return pdf
}

func (i *Invoice) addHeader(pdf *fpdf.Fpdf) {
	pdf.SetFont("Roboto", "B", 20)
	pdf.SetTextColor(34, 139, 34)
	pdf.CellFormat(100, 10, "INVOICE", "", 0, "L", false, 0, "")

	pdf.SetFont("Roboto", "", 10)
	pdf.SetTextColor(0, 0, 0)
	pdf.SetXY(140, 20)
	pdf.CellFormat(0, 5, fmt.Sprintf("Invoice # %s", i.InvoiceNum), "", 1, "R", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Date: %s", i.InvoiceDate.Format("02/01/2006")), "", 1, "R", false, 0, "")
	pdf.CellFormat(0, 5, fmt.Sprintf("Due date: %s", i.DueDate.Format("02/01/2006")), "", 1, "R", false, 0, "")
	pdf.Ln(10)
}

func (i *Invoice) addSenderInfo(pdf *fpdf.Fpdf) {
	pdf.SetFont("Roboto", "B", 11)
	pdf.CellFormat(0, 6, "From", "", 1, "", false, 0, "")
	pdf.SetFont("Roboto", "", 10)
	pdf.CellFormat(0, 5, i.Config.Sender.Name, "", 1, "", false, 0, "")
	pdf.CellFormat(0, 5, i.Config.Sender.City, "", 1, "", false, 0, "")
	pdf.CellFormat(0, 5, i.Config.Sender.Address, "", 1, "", false, 0, "")
	pdf.CellFormat(0, 5, "Reg Nr: "+i.Config.Sender.RegNr, "", 1, "", false, 0, "")
	pdf.CellFormat(0, 5, "Phone: "+i.Config.Sender.Phone, "", 1, "", false, 0, "")
	pdf.Ln(8)
}

func (i *Invoice) addRecipientInfo(pdf *fpdf.Fpdf) {
	pdf.SetFont("Roboto", "B", 11)
	pdf.CellFormat(0, 6, "Bill To", "", 1, "", false, 0, "")
	pdf.SetFont("Roboto", "", 10)
	pdf.CellFormat(0, 5, i.Config.BillTo.Name, "", 1, "", false, 0, "")
	for _, line := range i.Config.BillTo.Address {
		pdf.CellFormat(0, 5, line, "", 1, "", false, 0, "")
	}
	pdf.Ln(5)
}

func (i *Invoice) addProjectInfo(pdf *fpdf.Fpdf) {
	pdf.SetFont("Roboto", "B", 10)
	pdf.CellFormat(30, 5, "Project:", "", 0, "", false, 0, "")
	pdf.SetFont("Roboto", "", 10)
	pdf.CellFormat(60, 5, i.Config.ProjectName, "", 1, "", false, 0, "")
	pdf.SetFont("Roboto", "B", 10)
	pdf.CellFormat(30, 5, "Period:", "", 0, "", false, 0, "")
	pdf.SetFont("Roboto", "", 10)
	pdf.CellFormat(60, 5, i.Period, "", 1, "", false, 0, "")
	pdf.Ln(5)

	currentY := pdf.GetY()
	pdf.SetDrawColor(200, 200, 200)
	pdf.Line(15, currentY, 195, currentY)
	pdf.Ln(7)
}

func (i *Invoice) addItemsTable(pdf *fpdf.Fpdf) {
	pdf.SetFont("Roboto", "B", 10)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(100, 8, "Description", "1", 0, "L", true, 0, "")
	pdf.CellFormat(20, 8, "Qty", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Unit price", "1", 0, "C", true, 0, "")
	pdf.CellFormat(30, 8, "Total price", "1", 0, "C", true, 0, "")
	pdf.Ln(-1)

	pdf.SetFont("Roboto", "", 10)
	pdf.SetFillColor(255, 255, 255)
	for _, item := range i.Items {
		total := float64(item.Quantity) * item.UnitPrice
		pdf.CellFormat(100, 8, item.Description, "LRB", 0, "L", false, 0, "")
		pdf.CellFormat(20, 8, fmt.Sprintf("%d", item.Quantity), "LRB", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("€%.2f", item.UnitPrice), "LRB", 0, "C", false, 0, "")
		pdf.CellFormat(30, 8, fmt.Sprintf("€%.2f", total), "LRB", 0, "C", false, 0, "")
		pdf.Ln(-1)
	}
	pdf.Ln(5)
}

func (i *Invoice) addPaymentDetails(pdf *fpdf.Fpdf) {
	pdf.SetFont("Roboto", "B", 10)
	pdf.SetFillColor(230, 230, 230)
	pdf.CellFormat(120, 8, "Payment details", "1", 1, "L", true, 0, "")

	pdf.SetFont("Roboto", "", 10)
	pdf.SetFillColor(255, 255, 255)
	paymentDetails := [][]string{
		{"Account holder:", i.Config.Sender.Name},
		{"BIC:", i.Config.Payment.BIC},
		{"IBAN:", i.Config.Payment.IBAN},
		{"Address:", i.Config.Payment.Address},
	}
	for _, row := range paymentDetails {
		pdf.CellFormat(40, 8, row[0], "1", 0, "L", false, 0, "")
		pdf.CellFormat(80, 8, row[1], "1", 1, "L", false, 0, "")
	}
	pdf.Ln(5)
}

func (i *Invoice) addTotal(pdf *fpdf.Fpdf) {
	subtotalStr := fmt.Sprintf("€%.2f", i.Subtotal)

	pdf.SetFont("Roboto", "B", 10)
	pdf.CellFormat(100, 6, "", "", 0, "", false, 0, "")
	pdf.CellFormat(20, 6, "", "", 0, "", false, 0, "")
	pdf.CellFormat(30, 6, "Subtotal:", "", 0, "R", false, 0, "")
	pdf.SetFont("Roboto", "", 10)
	pdf.CellFormat(30, 6, subtotalStr, "", 0, "R", false, 0, "")
	pdf.Ln(10)

	pdf.SetFont("Roboto", "B", 14)
	pdf.SetTextColor(255, 0, 128)
	pdf.CellFormat(0, 8, subtotalStr, "", 0, "R", false, 0, "")
}

func (i *Invoice) generateFilename() string {
	safeName := strings.Replace(i.Config.Sender.Name, " ", "_", -1)
	return fmt.Sprintf("%s_%02d_%04d.pdf", safeName, i.Month, i.Year)
}

func (i *Invoice) calculateDates() {
	i.InvoiceDate = lastDayOfMonth(i.Year, i.Month)
	i.DueDate = calculateDueDate(i.Year, i.Month)
}

// Utility functions
func normalizeItems(configItems []config.Item, qtys []int) []config.Item {
	items := make([]config.Item, len(configItems))
	copy(items, configItems)

	for i := range items {
		if i < len(qtys) {
			items[i].Quantity = qtys[i]
		} else {
			items[i].Quantity = 1
		}
	}

	return items
}

func calculateSubtotal(items []config.Item) float64 {
	var total float64
	for _, item := range items {
		total += float64(item.Quantity) * item.UnitPrice
	}
	return total
}

func calculateDueDate(year, month int) time.Time {
	dueMonth := month + 1
	dueYear := year
	if dueMonth > 12 {
		dueMonth = 1
		dueYear++
	}
	return time.Date(dueYear, time.Month(dueMonth), 10, 0, 0, 0, 0, time.UTC)
}

func lastDayOfMonth(year, month int) time.Time {
	t := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return t.AddDate(0, 1, 0).Add(-24 * time.Hour)
}

func generateInvoiceNumber(invoiceNum, senderName string, year, month int) string {
	words := strings.Fields(senderName)
	initials := ""
	for _, word := range words {
		if len(word) > 0 {
			initials += strings.ToUpper(string(word[0]))
		}
	}

	return fmt.Sprintf("%s-%04d-%02d-%s", initials, year, month, invoiceNum)
}
