package controllers

import (
	"blog/models"

	"github.com/fpay/gopress"
)

const (
	defaultSortBy = "created_at desc"
	notFoundURL   = "/assets/404.html"
)

func getUser(ctx gopress.Context) *models.User {
	return ctx.Get("user").(*models.User)
}
