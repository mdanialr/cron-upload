package service

import "fmt"

// BytesToAnyBit convert bytes to any bytes using bit (1024) (Kibibyte, Mebibyte, Gibibyte).
func BytesToAnyBit(b int64, unit string) (string, error) {
	var conversion int64 = 1024

	switch unit {
	case "Kb":
		return fmt.Sprintf("%dKb", b/conversion), nil
	case "Mb":
		conversion = conversion * 1024
		return fmt.Sprintf("%dMb", b/conversion), nil
	case "Gb":
		conversion = conversion * 1024
		conversion = conversion * 1024
		return fmt.Sprintf("%dGb", b/conversion), nil
	}

	return "", fmt.Errorf("unit is not supported")
}
