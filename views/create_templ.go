// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func CreatePost() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>BlogFlex</title><link href=\"https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&amp;display=swap\" rel=\"stylesheet\"><link href=\"/static/css/tailwind.generated.css\" rel=\"stylesheet\"><style>\r\n      body {\r\n        font-family: 'Inter', sans-serif;\r\n      }\r\n    </style></head><body class=\"bg-gray-100\"><div class=\"max-w-3xl mx-auto p-4 bg-white shadow-md rounded-lg mt-10\"><h1 class=\"text-3xl font-bold mb-6 text-center\">Create a New Post</h1><form hx-post=\"/posts/create\" hx-target=\"#post-list\" hx-swap=\"innerHTML\" class=\"space-y-6\"><div><label for=\"title\" class=\"block text-sm font-medium text-gray-700\">Title</label><div class=\"mt-1\"><input type=\"text\" id=\"title\" name=\"title\" class=\"block w-full shadow-sm sm:text-sm border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500\" required></div></div><div><label for=\"content\" class=\"block text-sm font-medium text-gray-700\">Content</label><div class=\"mt-1\"><textarea id=\"content\" name=\"content\" rows=\"5\" class=\"block w-full shadow-sm sm:text-sm border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500\" required></textarea></div></div><div class=\"text-center\"><button type=\"submit\" class=\"inline-flex items-center px-6 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500\">Create Post</button></div></form><div id=\"post-list\" class=\"mt-6\"></div></div><script src=\"https://unpkg.com/htmx.org@2.0.0/dist/htmx.min.js\"></script></body>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
