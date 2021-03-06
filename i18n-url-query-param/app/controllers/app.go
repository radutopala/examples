// Copyright (c) Jeevanandam M. (https://github.com/jeevatkm)
// aahframework.org/examples source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package controllers

import "aahframe.work"

// AppController struct application controller
type AppController struct {
	*aah.Context
}

// Index method is application home page.
func (c *AppController) Index() {
	c.Reply().Ok()
}
