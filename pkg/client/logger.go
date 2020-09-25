package client

type Logger interface {
	Log(v ...interface{})
	Logf(format string, v ...interface{})
	Info(message string)
	Debug(message string)
}

func (c *Client) info(message string) {

	if c.logger == nil {
		return
	}

	c.logger.Info(message)
}

func (c *Client) debug(message string) {

	if c.logger == nil {
		return
	}

	c.logger.Debug(message)
}
