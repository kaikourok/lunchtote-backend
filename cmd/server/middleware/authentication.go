package middleware

import (
	"net/http"
	"reflect"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		if c.Request.Method == "GET" {
			if session.Get("cid") != nil {
				c.Next()
				return
			}
		} else {
			if session.Get("cid") != nil && c.Request.Header.Get("X-Auth-Key") == session.Get("csrf-token") {
				c.Next()
				return
			}
		}

		c.AbortWithStatus(http.StatusUnauthorized)
	}
}

func AuthAdministrator() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)

		administrator := session.Get("administrator")

		if session.Get("administrator") == nil || reflect.TypeOf(administrator).Kind() != reflect.Bool || !(administrator.(bool)) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		if c.Request.Method == "GET" {
			if session.Get("cid") != nil {
				c.Next()
				return
			}
		} else {
			if session.Get("cid") != nil && c.Request.Header.Get("X-Auth-Key") == session.Get("csrf-token") {
				c.Next()
				return
			}
		}

		c.AbortWithStatus(http.StatusNotFound)
	}
}
