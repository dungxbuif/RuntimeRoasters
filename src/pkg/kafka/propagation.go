package kafka

import "github.com/segmentio/kafka-go"

type kafkaHeadersCarrier []kafka.Header

func (c *kafkaHeadersCarrier) Get(key string) string {
	if c == nil {
		return ""
	}
	headers := []kafka.Header(*c)
	for i := len(headers) - 1; i >= 0; i-- {
		if headers[i].Key == key {
			return string(headers[i].Value)
		}
	}
	return ""
}

func (c *kafkaHeadersCarrier) Set(key string, value string) {
	if c == nil {
		return
	}
	headers := []kafka.Header(*c)
	for i := range headers {
		if headers[i].Key == key {
			headers[i].Value = []byte(value)
			*c = headers
			return
		}
	}
	*c = append(headers, kafka.Header{Key: key, Value: []byte(value)})
}

func (c *kafkaHeadersCarrier) Keys() []string {
	if c == nil {
		return nil
	}
	headers := []kafka.Header(*c)
	keys := make([]string, 0, len(headers))
	for _, header := range headers {
		keys = append(keys, header.Key)
	}
	return keys
}
