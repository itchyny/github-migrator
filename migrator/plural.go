package migrator

import "fmt"

func pluralUnit(count int, unit string) string {
	if count == 1 {
		return unit
	}
	return unit + "s"
}

func plural(count int, unit string) string {
	return fmt.Sprint(count) + " " + pluralUnit(count, unit)
}
