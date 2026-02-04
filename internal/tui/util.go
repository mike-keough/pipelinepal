package tui

import (
	"fmt"
	"strings"
)

func clamp(v, lo, hi int) int {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func ellipsize(s string, n int) string {
	s = strings.TrimSpace(s)
	if n <= 0 || len(s) <= n {
		return s
	}
	if n <= 1 {
		return s[:n]
	}
	return s[:n-1] + "…"
}

func fmtLeadLine(name, leadType, source string) string {
	if source == "" {
		return fmt.Sprintf("%s [%s]", name, leadType)
	}
	return fmt.Sprintf("%s [%s] • %s", name, leadType, source)
}
