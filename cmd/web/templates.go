package main

import "snippetbox.abdulalsh.com/internal/models"

type TemplateData struct {
	Snippet  models.Snippet
	Snippets []models.Snippet
}
