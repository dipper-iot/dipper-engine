package output_redis_extend_queue

type OptionSession struct {
	RedisAddress  string `json:"redis_address"`
	RedisPassword string `json:"redis_password"`
	RedisDb       int    `json:"redis_db"`
	Queue         string `json:"queue"`
	NextSuccess   string `json:"next_success"`
	NextError     string `json:"next_error"`
}
