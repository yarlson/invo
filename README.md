# invo

A minimalist CLI tool for generating professional PDF invoices. Perfect for freelancers and small businesses who need to generate monthly invoices with consistent formatting.

## Quick Start

1. Install:

```bash
brew install yarlson/invo/invo   # macOS/Linux with Homebrew
```

See [Installation](#installation) for other methods.

2. Create config:

```bash
# ~/.config/invo/config.yaml
sender:
 name: "Your Company Name"
 city: "Your City"
 address: "Your Address"
 reg_nr: "12345"
 phone: "+1234567890"

bill_to:
 name: "Client Name"
 address: ["Client Address"]

project_name: "Project Name"

payment:
 bic: "BANKBIC"
 iban: "XX00BANK0000"
 address: "Bank Address"

items:
 - description: "Service Description"
   unit_price: 100.0
```

3. Generate invoice:

```bash
invo                     # Basic invoice for current month
invo -m 1 -n INV-123    # January invoice with custom number
```

## Common Tasks

Generate invoice for specific month:

```bash
invo -m 3                  # March invoice
invo -y 2024 -m 3         # March 2024 invoice
```

Set custom dates:

```bash
invo -d 2024-03-15        # Custom invoice date
invo -d 2024-03-15 -D 2024-04-15  # Custom invoice and due dates
```

Multiple items with different quantities:

```bash
invo -q 2,1,3             # First item: 2 units
                          # Second item: 1 unit
                          # Third item: 3 units
```

Use different config file:

```bash
invo -c path/to/config.yaml
```

## Command Reference

```bash
invo [flags]
```

Core flags:

- `--config, -c`  
  Config file path. Default: `config.yaml`, checks `~/.config/invo/config.yaml` first
- `--number, -n`  
  Invoice number. Default: `01`

Time flags:

- `--year, -y`  
  Year (YYYY). Default: current year
- `--month, -m`  
  Month (1-12). Default: `1`
- `--date, -d`  
  Invoice date (YYYY-MM-DD). Default: last day of month
- `--due, -D`  
  Due date (YYYY-MM-DD). Default: 10th of next month

Item flags:

- `--quantities, -q`  
  Comma-separated quantities. Default: `1`

## Output Format

Generated PDF includes:

- Invoice number: `<INITIALS>-YYYY-MM-<number>` (e.g., `YC-2024-03-INV123`)
- Sender and recipient details
- Project information
- Item list with quantities and prices
- Payment details
- Professional styling with embedded fonts

## Installation

### Homebrew (macOS/Linux)

```bash
brew install yarlson/invo/invo
```

### Direct Download

Download the latest release from [GitHub Releases](https://github.com/yarlson/invo/releases) for your platform:

- Linux (amd64/arm64)
- macOS (amd64/arm64)
- Windows (amd64/arm64)

Move the binary to your PATH (e.g., `/usr/local/bin`).

### Build from Source

Requires Go 1.23+:

```bash
git clone https://github.com/yarlson/invo
cd invo
go build -o invo ./cmd/invo
```

## Configuration Details

The tool looks for config in this order:

1. Path specified by `--config`
2. `~/.config/invo/config.yaml`
3. `./config.yaml`

See [example config](#quick-start) above for the structure.

## Customization

To modify the invoice template, edit:

- `./cmd/invo/main.go` - CLI interface
- `./pkg/invoice` - PDF generation

## License

[MIT License](./LICENSE)
