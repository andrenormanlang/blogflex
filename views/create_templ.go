// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func CreatePost(userID uint) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>BlogFlex</title><link href=\"https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&amp;display=swap\" rel=\"stylesheet\"><link href=\"https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css\" rel=\"stylesheet\"><style>\r\n            body {\r\n                font-family: 'Inter', sans-serif;\r\n            }\r\n        </style></head><body class=\"bg-light\"><div class=\"container mt-5\"><div class=\"row justify-content-center\"><div class=\"col-md-8\"><div class=\"card shadow-sm\"><div class=\"card-body\"><h1 class=\"card-title text-center mb-4\">Create a New Post</h1><form hx-post=\"/protected/posts/create\" hx-redirect=\"/protected/posts\" hx-target=\"#response-message\" hx-swap=\"innerHTML\"><div class=\"form-group\"><label for=\"title\">Title</label> <input type=\"text\" id=\"title\" name=\"title\" class=\"form-control\" required></div><div class=\"form-group\"><label for=\"content\">Content</label> <textarea id=\"content\" name=\"content\" rows=\"5\" class=\"form-control\" required></textarea></div><input type=\"hidden\" name=\"user_id\" value=\"` + strconv.FormatUint(uint64(userID), 10) + `\"><div class=\"text-center\"><button type=\"submit\" class=\"btn btn-primary\">Create Post</button></div></form><div id=\"response-message\" class=\"mt-4\"></div></div></div></div></div></div><script src=\"https://unpkg.com/htmx.org@2.0.0/dist/htmx.min.js\"></script></body>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
