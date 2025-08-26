package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"tempo/internal/storage"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [filename]",
	Short: "Export job configurations",
	Long: `Export all job configurations to a JSON file for backup or migration.

Examples:
  tempo export
  tempo export backup.json
  tempo export --format yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: runExport,
}

var exportFormat string

func init() {
	exportCmd.Flags().StringVarP(&exportFormat, "format", "f", "json", "Export format (json, yaml)")
}

func runExport(cmd *cobra.Command, args []string) error {
	filename := "tempo-jobs.json"
	if len(args) > 0 {
		filename = args[0]
	}

	store, err := storage.NewStorage("")
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %v", err)
	}

	jobs := store.GetAllJobs()

	var data []byte
	var err2 error

	switch exportFormat {
	case "json":
		data, err2 = json.MarshalIndent(jobs, "", "  ")
		if err2 != nil {
			return fmt.Errorf("failed to marshal JSON: %v", err2)
		}
	case "yaml":
		return fmt.Errorf("YAML export not implemented yet")
	default:
		return fmt.Errorf("unsupported format: %s", exportFormat)
	}

	err2 = os.WriteFile(filename, data, 0644)
	if err2 != nil {
		return fmt.Errorf("failed to write file: %v", err2)
	}

	fmt.Printf("✓ Exported %d jobs to %s\n", len(jobs), filename)
	return nil
}
