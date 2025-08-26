package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"tempo/internal/storage"
	"tempo/internal/types"

	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [filename]",
	Short: "Import job configurations",
	Long: `Import job configurations from a JSON file.

Examples:
  tempo import backup.json
  tempo import --format yaml jobs.yaml`,
	Args: cobra.ExactArgs(1),
	RunE: runImport,
}

var importFormat string

func init() {
	importCmd.Flags().StringVarP(&importFormat, "format", "f", "json", "Import format (json, yaml)")
}

func runImport(cmd *cobra.Command, args []string) error {
	filename := args[0]

	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	var jobs []types.Job

	switch importFormat {
	case "json":
		err = json.Unmarshal(data, &jobs)
		if err != nil {
			return fmt.Errorf("failed to parse JSON: %v", err)
		}
	case "yaml":
		return fmt.Errorf("YAML import not implemented yet")
	default:
		return fmt.Errorf("unsupported format: %s", importFormat)
	}

	store, err := storage.NewStorage("")
	if err != nil {
		return fmt.Errorf("failed to initialize storage: %v", err)
	}

	// Import jobs to storage
	for _, job := range jobs {
		if err := store.AddJob(job); err != nil {
			return fmt.Errorf("failed to import job '%s': %v", job.ID, err)
		}
	}

	fmt.Printf("✓ Imported %d jobs from %s\n", len(jobs), filename)
	return nil
}
