# invo: a minimalist CLI tool for generating professional PDF invoices

A simple command-line tool written in Go that generates professional PDF invoices. This tool is particularly useful for freelancers and small businesses who need to generate monthly invoices with consistent formatting.

## Features

- Generates professional PDF invoices
- Customizable invoice period (month/year)
- Automatic calculation of due dates
- Includes payment details and company information
- Supports customizable invoice numbers
- Clean, professional design with embedded fonts

## Prerequisites

- Go 1.23 or higher

## Installation

### Direct Download

Pre-built binaries are available for Linux, macOS, and Windows. You can download the latest release from the [GitHub Releases page](https://github.com/yarlson/invo/releases). For example, if you're using version 2, choose the appropriate archive for your operating system and architecture:

- **Linux (amd64/arm64)**
- **macOS (amd64/arm64)**
- **Windows (amd64/arm64)**

After downloading, extract the archive and move the `invo` binary to a directory included in your `PATH` (e.g. `/usr/local/bin` on Unix-like systems).

### Homebrew

If you are on macOS (or Linux with Homebrew), you can install `invo` via Homebrew. The Homebrew formula is maintained in the [homebrew-invo repository](https://github.com/yarlson/homebrew-invo).

To install `invo`, run:

```bash
brew install yarlson/invo/invo
```

To verify the installation, run:

```bash
invo --help
```

### Building from Source

If you prefer to build from source, follow these steps:

1. **Clone the repository:**

   ```bash
   git clone https://github.com/yarlson/invo
   cd invo
   ```

2. **Build the program:**

   ```bash
   go build -o invo ./cmd/invoice
   ```

   The project is configured to build with `CGO_ENABLED=0` for portability and supports multiple operating systems and architectures. For example, you can cross-compile for Linux, macOS, or Windows.

After building, move the `invo` binary to a directory included in your `PATH`.

## Usage

Run the program with optional flags to customize the invoice. The tool now uses [Cobra](https://github.com/spf13/cobra) for command-line parsing, which provides both long and shorthand flag options.

```bash
./invo [flags]
```

### Available Flags

- `--invoice, -i`  
  Custom invoice number (alphanumeric string) that will be used as the final part of the generated invoice number.  
  **Default:** `"01"`

- `--year, -y`  
  Year in 4-digit format.  
  **Default:** current year (e.g. if running in 2025, default is `2025`)

- `--month, -m`  
  Month as a number (1-12).  
  **Default:** `1`

- `--qty, -q`  
  Comma separated quantities for each invoice item (e.g. `"2,1"`).  
  **Default:** `"1"`

- `--config, -c`  
  Path to the configuration file.  
  **Default:** `config.yaml`  
  **Note:** If not explicitly provided, the program first checks for a configuration file in the XDG configuration directory under the `invo` subfolder. If a file exists there, it is used; otherwise, the default file path (`config.yaml`) is used.

### Examples

Generate an invoice for January of the current year with the default invoice number:

```bash
./invo -m 1
```

Generate an invoice for March 2024 with multiple quantities and a custom invoice number:

```bash
./invo -i INV-123 -y 2024 -m 3 -q "2,1"
```

Generate an invoice using a custom configuration file:

```bash
./invo -c /path/to/custom_config.yaml
```

The program will generate a PDF file (with a filename based on the sender and date) in the current directory.

## Configuration File

The configuration file is written in YAML format and contains details about the sender, billing information, project name, payment details, and the invoice items.

### Default File Lookup

If you do not specify a configuration file using the `--config` flag, the tool checks for the file in the following order:

1. **XDG Configuration Directory:**  
   The tool looks for a configuration file in the XDG configuration directory under an `invo` subfolder.

   - **What is the XDG Configuration Directory?**  
     The XDG Base Directory Specification defines a standard location for user-specific configuration files. By default, if the environment variable `XDG_CONFIG_HOME` is not set, it typically defaults to `~/.config` on Unix-like systems.
   - **Example:**  
     On most Linux systems, the file would be located at:
     ```bash
     ~/.config/invo/config.yaml
     ```

2. **Fallback to Local File:**  
   If no file is found in the XDG directory, the tool falls back to using the local file `config.yaml` in the current directory.

### Example Configuration File (`config.yaml`)

```yaml
sender:
  name: "Your Company Name"
  city: "Your City, Country"
  address: "Your Street Address, Suite/Unit Number"
  reg_nr: "Company Registration Number"
  phone: "+1 234-567-8900"

bill_to:
  name: "Client Company Name"
  address:
    - "Client Street Address, Suite/Unit"
    - "Client City, State/Province"
    - "Client Country"

project_name: "Project or Service Name"

payment:
  bic: "BANKBICXXX"
  iban: "XX00BANK00000000000000"
  address: "Your Company Payment Address"

items:
  - description: "Service or Product Description"
    unit_price: 100.0
```

## Output

The generated invoice includes:

- Invoice number in the format: `<INITIALS>-YYYY-MM-<custom_invoice_number>` (where `<INITIALS>` are derived from the sender's name)
- Invoice date (last day of the specified month)
- Due date (10th of the next month)
- Sender and recipient details
- Item descriptions and pricing for multiple services
- Payment details
- Total amount

## Customization

To modify the invoice template, you can edit the code in `./cmd/invoice/main.go` and the package in `./pkg/invoice`. You can update:

- Company details
- Payment information
- Pricing and service items
- Colors and formatting
- Invoice layout

## License

[MIT License](./LICENSE)
