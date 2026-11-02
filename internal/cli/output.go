package cli

import "github.com/spf13/cobra"

func jsonOutput(root *cobra.Command) bool {
	flag := root.PersistentFlags().Lookup("json")
	return flag != nil && flag.Value.String() == "true"
}
