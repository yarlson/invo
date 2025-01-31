package main

import (
	"log"
	"strconv"
	"strings"
	"time"

	"invo/pkg/config"
	"invo/pkg/invoice"

	"github.com/spf13/cobra"
)

var (
	invoiceFlag string
	year        int
	month       int
	qtyFlag     string
	configFile  string
)

func main() {
	currentYear := time.Now().Year()

	rootCmd := &cobra.Command{
		Use:   "invoice",
		Short: "invo: a minimalist CLI tool for generating professional PDF invoices",
		Run: func(cmd *cobra.Command, args []string) {
			qtyStrs := strings.Split(qtyFlag, ",")
			var qtys []int
			for _, s := range qtyStrs {
				s = strings.TrimSpace(s)
				if s == "" {
					continue
				}
				q, err := strconv.Atoi(s)
				if err != nil {
					log.Fatalf("Invalid quantity value %q: %v", s, err)
				}
				qtys = append(qtys, q)
			}

			cfg, err := config.LoadConfig(configFile)
			if err != nil {
				log.Fatalf("Error loading config: %v", err)
			}

			inv := invoice.NewInvoice(invoiceFlag, year, month, qtys, cfg)
			filename, err := inv.GeneratePDF()
			if err != nil {
				log.Fatalf("Error generating PDF: %v", err)
			}
			log.Printf("PDF generated: %s", filename)
		},
	}

	rootCmd.Flags().IntVarP(&year, "year", "y", currentYear, "Invoice year (defaults to current year)")
	rootCmd.Flags().IntVarP(&month, "month", "m", 1, "Invoice month (1-12)")
	rootCmd.Flags().StringVarP(&qtyFlag, "qty", "q", "1", "Comma separated quantities for each invoice item (e.g. \"2,1\")")
	rootCmd.Flags().StringVarP(&invoiceFlag, "invoice", "i", "01", "Invoice number (e.g. \"01\")")
	rootCmd.Flags().StringVarP(&configFile, "config", "c", "config.yaml", "Path to the configuration file")

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
