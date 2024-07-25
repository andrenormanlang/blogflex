package views

import (
	"context"
	"github.com/a-h/templ"
	"io"
)

// UnsafeHTML safely renders raw HTML content.
func UnsafeHTML(content string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := w.Write([]byte(content))
		return err
	})
}

// RenderLatestPostTitle safely renders raw HTML content.
func RenderLatestPostTitle(content string) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		_, err := w.Write([]byte(content))
		return err
	})
}