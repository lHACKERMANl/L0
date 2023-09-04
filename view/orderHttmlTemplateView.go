package View

import (
	"fmt"
	"html/template"
	"strings"
)

var orderQueryTemplate = `
<section>
    <h1>Enter Order ID</h1>
    <form id="orderForm" action="/order" method="get">
        <label for="order_id">Order ID:</label>
        <input type="text" id="order_id" name="order_id" required>
        <br>
        <input type="submit" value="Submit" onclick="redirectToOrder()">
    </form>
</section>
`
var orderQueryTemplateFuncs = template.FuncMap{
	"OrderhtmlSafe": func(text string) template.HTML {
		return template.HTML(text)
	},
}

func OrderHTMLView() string {
	tmpl, err := template.New("orderQueryTemplate").Funcs(orderQueryTemplateFuncs).Parse(orderQueryTemplate)
	if err != nil {
		fmt.Sprintf("Error parsing template: %v", err)
	}

	var html strings.Builder
	var emptyData struct{}
	err = tmpl.Execute(&html, emptyData)
	if err != nil {
		fmt.Sprintf("Error executing data: %v", err)
	}
	return html.String()
}
