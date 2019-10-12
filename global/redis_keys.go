package global

import "fmt"

func DataKey(s string) string {
	return fmt.Sprintf("task_data/%s", s)
}
