/*
   Copyright 2019 Septian Wibisono

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.

*/
package main

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/septianw/jas/common"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()
	if *Mode == "production" {
		gin.SetMode("release")
	} else {
		gin.SetMode("debug")
	}

	switch STAGE {
	case "development":
		gin.SetMode("debug")
		break
	case "production":
		gin.SetMode("release")
		break
	case "testing":
		gin.SetMode("debug")
		break
	default:
		gin.SetMode("debug")
	}

	// r.GET("/ping", func(c *gin.Context) {
	// 	c.String(http.StatusOK, "pong")
	// })

	r.NoRoute(func(c *gin.Context) {
		common.SendHttpError(c, common.PAGE_NOT_FOUND_CODE, errors.New("Page not found."))
	})

	r.NoMethod(func(c *gin.Context) {
		common.SendHttpError(c, common.PAGE_NOT_FOUND_CODE, errors.New("Page not found."))
	})

	r.GET("/", func(c *gin.Context) {
		middleware, exist := c.Get("middleware")
		if exist {
			c.String(http.StatusOK, "Hello World"+middleware.(string))
		} else {
			c.String(http.StatusOK, "Hello World")
		}
	})

	return r
}
