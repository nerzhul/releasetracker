package utils

import (
	"bytes"
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateHTTPLoggerConfig() middleware.LoggerConfig {
	loggerConfig := middleware.DefaultLoggerConfig
	loggerConfig.Format = `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}",` +
		`"host":"${host}","method":"${method}","uri":"${uri}","user_agent":"${user_agent}",` +
		`"status":${status},"error":"${error}","latency":${latency},"latency_human":"${latency_human}"` +
		`,"bytes_in":${bytes_in},"bytes_out":${bytes_out},"custom":${custom}}` + "\n"
	loggerConfig.CustomTagFunc = func(c echo.Context, buf *bytes.Buffer) (int, error) {
		switch v := c.Get("username").(type) {
		case string:
			b, err := json.Marshal(struct {
				Username string `json:"username"`
			}{
				Username: v,
			})

			if err != nil {
				return 0, err
			}

			buf.Write(b)
		default:
			buf.WriteString(`{}`)
		}
		return 0, nil
	}

	return loggerConfig
}
