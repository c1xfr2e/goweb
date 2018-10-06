package middleware

import (
	"bytes"
	"io/ioutil"

	"github.com/bluecover/lm/server/codec"
	"github.com/bluecover/lm/server/render"
	"github.com/gin-gonic/gin"
)

func RequestDecoder(dec codec.Decoder) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := decodeBody(c, dec); err != nil {
			render.BadRequest(c, err)
			c.Abort()
			return
		}

		if err := decodeQueryString(c, dec); err != nil {
			render.BadRequest(c, err)
			c.Abort()
			return
		}
	}
}

func decodeBody(c *gin.Context, dec codec.Decoder) error {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return err
	}

	if len(body) == 0 {
		return nil
	}

	body, err = dec.Decode(body)
	if err != nil {
		return err
	}

	c.Request.Body.Close()
	c.Request.Body = ioutil.NopCloser(bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", gin.MIMEPOSTForm)
	return nil
}

func decodeQueryString(c *gin.Context, dec codec.Decoder) error {
	param := c.Query("param")
	if len(param) == 0 {
		return nil
	}

	data, err := dec.Decode([]byte(param))
	if err != nil {
		return err
	}

	c.Request.URL.RawQuery = string(data)
	return nil
}
