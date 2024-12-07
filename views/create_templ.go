// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func CreatePost(userID string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Create Post - BlogFlex</title><link href=\"https://fonts.googleapis.com/css2?family=Inter:wght@400;600;700&amp;display=swap\" rel=\"stylesheet\"><link href=\"https://stackpath.bootstrapcdn.com/bootstrap/4.5.2/css/bootstrap.min.css\" rel=\"stylesheet\"><!-- Place the first <script> tag in your HTML's <head> --><script src=\"https://cdn.tiny.cloud/1/7xbog5vcrh51qgt6v64oxx7cvpqbgjezxtc42rq11wtscrsq/tinymce/7/tinymce.min.js\" referrerpolicy=\"origin\"></script><!-- Place the following <script> and <textarea> tags your HTML's <body> --><script>\r\n  tinymce.init({\r\n    selector: 'textarea',\r\n    plugins: 'anchor autolink charmap codesample emoticons image link lists media searchreplace table visualblocks wordcount linkchecker',\r\n    toolbar: 'undo redo | blocks fontfamily fontsize | bold italic underline strikethrough | link image media table | align lineheight | numlist bullist indent outdent | emoticons charmap | removeformat',\r\n  });\r\n</script><style>\r\n        body {\r\n          font-family: 'Inter', sans-serif;\r\n          background-color: #121212;\r\n          color: #e0e0e0;\r\n        }\r\n        .container {\r\n          max-width: 1000px;\r\n        }\r\n        .card {\r\n          background-color: #1e1e1e;\r\n          border: none;\r\n          box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);\r\n          border-radius: 10px;\r\n        }\r\n        .card-body {\r\n          padding: 20px;\r\n        }\r\n        .card-title {\r\n          color: #bb86fc;\r\n        }\r\n        .form-control {\r\n          background-color: #2c2c2c;\r\n          border: none;\r\n          color: #e0e0e0;\r\n        }\r\n        .form-control:focus {\r\n          background-color: #2c2c2c;\r\n          color: #e0e0e0;\r\n        }\r\n        .btn-primary {\r\n          background-color: #bb86fc;\r\n          border: none;\r\n        }\r\n        .btn-primary:hover {\r\n          background-color: #3700b3;\r\n        }\r\n      </style></head><body><div class=\"container mt-5\"><div class=\"row justify-content-center\"><div class=\"col-md-12\"><div class=\"card shadow-sm\"><div class=\"card-body\"><form id=\"create-post-form\" hx-post=\"/protected/posts\" hx-redirect=\"true\" hx-target=\"#response-message\" hx-swap=\"innerHTML\" method=\"POST\"><div class=\"form-group\"><label for=\"title\">Title</label> <textarea id=\"title\" name=\"title\" rows=\"2\" class=\"form-control\" required></textarea></div><div class=\"form-group\"><label for=\"content\">Content</label> <textarea id=\"content\" name=\"content\" rows=\"10\" class=\"form-control\" required></textarea></div><input type=\"hidden\" name=\"user_id\" value=\"")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(userID)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/create.templ`, Line: 80, Col: 59}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><div class=\"text-center\"><button type=\"submit\" class=\"btn btn-primary\">Create Post</button></div></form><div id=\"response-message\" class=\"mt-4\"></div></div></div></div></div></div><script>\r\n        document.addEventListener(\"DOMContentLoaded\", function() {\r\n          tinymce.init({\r\n            selector: '#title, #content',\r\n            plugins: 'anchor autolink charmap codesample emoticons image link lists media searchreplace table visualblocks wordcount checklist mediaembed casechange export formatpainter pageembed linkchecker a11ychecker tinymcespellchecker permanentpen powerpaste advtable advcode editimage advtemplate ai mentions tinycomments tableofcontents footnotes mergetags autocorrect typography inlinecss markdown',\r\n            toolbar: 'undo redo | blocks fontfamily fontsize | bold italic underline strikethrough | link image media table mergetags | addcomment showcomments | spellcheckdialog a11ycheck typography | align lineheight | checklist bullist numlist indent outdent | emoticons charmap | removeformat',\r\n            tinycomments_mode: 'embedded',\r\n            tinycomments_author: 'Author name',\r\n            mergetags_list: [\r\n              { value: 'First.Name', title: 'First Name' },\r\n              { value: 'Email', title: 'Email' },\r\n            ],\r\n            ai_request: (request, respondWith) => respondWith.string(() => Promise.reject(\"See docs to implement AI Assistant\")),\r\n            setup: function(editor) {\r\n              editor.on('change', function(e) {\r\n                editor.save();\r\n              });\r\n            },\r\n            skin: window.matchMedia(\"(prefers-color-scheme: dark)\").matches ? \"oxide-dark\" : \"\",\r\n            content_css: window.matchMedia(\"(prefers-color-scheme: dark)\").matches ? \"dark\" : \"\"\r\n          });\r\n\r\n          document.getElementById('create-post-form').addEventListener('submit', function(e) {\r\n            if (tinymce.get('content').getContent() === '' || tinymce.get('title').getContent() === '') {\r\n              e.preventDefault();\r\n              alert('Both Title and Content are required.');\r\n            }\r\n          });\r\n        });\r\n      </script><script src=\"https://unpkg.com/htmx.org@2.0.0/dist/htmx.min.js\"></script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
