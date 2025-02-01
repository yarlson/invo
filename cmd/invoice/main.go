// Package main provides the entry point for the invoice generation tool.
package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"invo/pkg/config"
	"invo/pkg/invoice"
)

type flags struct {
	invoiceNum string
	year       int
	month      int
	quantities string
	configFile string
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	f := &flags{}
	cmd := newRootCommand(f)
	return cmd.Execute()
}

func newRootCommand(f *flags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "invoice",
		Short: "invo: a minimalist CLI tool for generating professional PDF invoices",
		RunE:  createRunFunc(f),
	}

	initFlags(cmd, f)
	return cmd
}

func createRunFunc(f *flags) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		qtys, err := parseQuantities(f.quantities)
		if err != nil {
			return fmt.Errorf("parsing quantities: %w", err)
		}

		cfg, err := config.LoadConfig(f.configFile)
		if err != nil {
			return fmt.Errorf("loading config: %w", err)
		}

		inv := invoice.NewInvoice(f.invoiceNum, f.year, f.month, qtys, cfg)
		filename, err := inv.GeneratePDF()
		if err != nil {
			return fmt.Errorf("generating PDF: %w", err)
		}

		log.Printf("PDF generated: %s", filename)
		return nil
	}
}

func initFlags(cmd *cobra.Command, f *flags) {
	currentYear := time.Now().Year()

	cmd.Flags().StringVarP(&f.invoiceNum, "invoice", "i", "01", "Invoice number (e.g. \"01\")")
	cmd.Flags().IntVarP(&f.year, "year", "y", currentYear, "Invoice year")
	cmd.Flags().IntVarP(&f.month, "month", "m", 1, "Invoice month (1-12)")
	cmd.Flags().StringVarP(&f.quantities, "qty", "q", "1", "Comma separated quantities (e.g. \"2,1\")")
	cmd.Flags().StringVarP(&f.configFile, "config", "c", "config.yaml", "Path to config file")
}

func parseQuantities(qtyStr string) ([]int, error) {
	if qtyStr == "" {
		return nil, nil
	}

	var qtys []int
	for _, s := range strings.Split(qtyStr, ",") {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}

		q, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid quantity value %q: %w", s, err)
		}
		qtys = append(qtys, q)
	}

	return qtys, nil
}
