package output_service

import (
	"fmt"
	"os"

	log_service "github.com/Rfluid/insta-tools/src/log/service"
	output_flag "github.com/Rfluid/insta-tools/src/output/flag"
	"github.com/pterm/pterm"
)

// WriteOutput writes data to the specified OutputPath if provided
func WriteConditionally(data string) error {
	if output_flag.OutputPath == "" {
		log_service.LogConditionally(
			pterm.DefaultLogger.Info,
			"No output path set. Will not write to file.",
		)
		return nil // Do nothing if no output path is set
	}

	log_service.LogConditionally(
		pterm.DefaultLogger.Info,
		fmt.Sprintf("Writing to output file %s", output_flag.OutputPath),
	)

	file, err := os.Create(output_flag.OutputPath) // Create or overwrite the file
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString(data) // Write the data to the file
	if err != nil {
		return fmt.Errorf("failed to write to output file: %w", err)
	}

	log_service.LogConditionally(
		pterm.DefaultLogger.Info,
		fmt.Sprintf("Wrote to output file %s", output_flag.OutputPath),
	)

	return nil
}

func PrintConditionally(data string) {
	if output_flag.OutputPath != "" {
		return
	}
	log_service.LogConditionally(
		pterm.DefaultLogger.Info,
		"No output path set. Will print results.",
	)
	fmt.Println(data)
}
