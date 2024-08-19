package config

var runtime_schedule_table map[string]string

func init() {
	runtime_schedule_table = make(map[string]string)
}

func GetScheduleTable(table_name string) string {
	return runtime_schedule_table[table_name]
}

func SetScheduleTable(table_name string, table_content string) {
	runtime_schedule_table[table_name] = table_content
}
