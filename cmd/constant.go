package cmd

var validStatuses = map[string]bool{
	"Pending":   true,
	"Running":   true,
	"Succeeded": true,
	"Failed":    true,
	"Unknown":   true,
}
