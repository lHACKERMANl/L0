package View

import (
	"fmt"
	"html/template"
	"strings"
)

var orderTemplate = `
<section>
    <style>
        table {
            border-collapse: collapse;
            width: 100%;
        }
        th, td {
            border: 1px solid #ddd;
            padding: 8px;
            text-align: left;
        }
        th {
            background-color: #f2f2f2;
        }
    </style>
    <h1>Order Details</h1>
    <table>
        <tr>
            <th>Field</th>
            <th>Value</th>
        </tr>
        <tr>
            <td>Order ID</td>
            <td>{{ .OrderID }}</td>
        </tr>
        <tr>
            <td>Order Details</td>
            <td>
                <table>
                    <tr>
                        <td>Order UID</td>
                        <td>{{ .Order.OrderUID }}</td>
                    </tr>
                    <tr>
                        <td>Track Number</td>
                        <td>{{ .Order.TrackNumber }}</td>
                    </tr>
                    <tr>
                        <td>Entry</td>
                        <td>{{ .Order.Entry }}</td>
                    </tr>
                    <tr>
                        <td>Locale</td>
                        <td>{{ .Order.Locale }}</td>
                    </tr>
                    <tr>
                        <td>InternalSignature</td>
                        <td>{{ .Order.InternalSignature }}</td>
                    </tr>
                    <tr>
                        <td>CustomerId</td>
                        <td>{{ .Order.CustomerId }}</td>
                    </tr>
                    <tr>
                        <td>DeliveryService</td>
                        <td>{{ .Order.DeliveryService }}</td>
                    </tr>
                    <tr>
                        <td>Shardkey</td>
                        <td>{{ .Order.Shardkey }}</td>
                    </tr>
                    <tr>
                        <td>SmID</td>
                        <td>{{ .Order.SmID }}</td>
                    </tr>
                    <tr>
                        <td>DataCreated</td>
                        <td>{{ .Order.DataCreated }}</td>
                    </tr>
                    <tr>
                        <td>OofShard</td>
                        <td>{{ .Order.OofShard }}</td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td>Order Payment</td>
            <td>
                <table>
                    <tr>
                        <td>Transaction</td>
                        <td>{{ .Payment.Transaction }}</td>
                    </tr>
                    <tr>
                        <td>Request ID</td>
                        <td>{{ .Payment.RequestID }}</td>
                    </tr>
                    <tr>
                        <td>Currency</td>
                        <td>{{ .Payment.Currency }}</td>
                    </tr>
                    <tr>
                        <td>Provider</td>
                        <td>{{ .Payment.Provider }}</td>
                    </tr>
                    <tr>
                        <td>amount</td>
                        <td>{{ .Payment.Amount }}</td>
                    </tr>
                    <tr>
                        <td>PaymentDT</td>
                        <td>{{ .Payment.PaymentDT }}</td>
                    </tr>
                    <tr>
                        <td>Bank</td>
                        <td>{{ .Payment.Bank }}</td>
                    </tr>
                    <tr>
                        <td>DeliveryCost</td>
                        <td>{{ .Payment.DeliveryCost }}</td>
                    </tr>
                    <tr>
                        <td>GoodsTotal</td>
                        <td>{{ .Payment.GoodsTotal }}</td>
                    </tr>
                    <tr>
                        <td>CustomFee</td>
                        <td>{{ .Payment.CustomFee }}</td>
                    </tr>
                </table>
            </td>
        </tr>
        <tr>
            <td>Order Items</td>
            <td>
                <table>
					{{ range .Items }}
                    <tr>
                        <td>Chrt ID</td>
                        <td>{{ .ChrtID }}</td>
                    </tr>
                    <tr>
                        <td>Track Number</td>
                        <td>{{ .TrackNumber }}</td>
                    </tr>
                    <tr>
                        <td>Price</td>
                        <td>{{ .Price }}</td>
                    </tr>
                    <tr>
                        <td>Name</td>
                        <td>{{ .Name }}</td>
                    </tr>
                    <tr>
                        <td>Rid</td>
                        <td>{{ .Rid }}</td>
                    </tr>
                    <tr>
                        <td>Sale</td>
                        <td>{{ .Sale }}</td>
                    </tr>
                    <tr>
                        <td>Size</td>
                        <td>{{ .Size }}</td>
                    </tr>
                    <tr>
                        <td>TotalPrice</td>
                        <td>{{ .TotalPrice }}</td>
                    </tr>
                    <tr>
                        <td>NmID</td>
                        <td>{{ .NmID }}</td>
                    </tr>
                    <tr>
                        <td>Brand</td>
                        <td>{{ .Brand }}</td>
                    </tr>
                    <tr>
                        <td>Status</td>
                        <td>{{ .Status }}</td>
                    </tr>
                    <tr>
                        <td>OrderID</td>
                        <td>{{ .OrderID }}</td>
                    </tr>
					{{ end }}
                </table>
            </td>
        </tr>
        <tr>
            <td>Order Delivery</td>
            <td>
                <table>
                    <tr>
                        <td>Name</td>
                        <td>{{ .Delivery.Name }}</td>
                    </tr>
                    <tr>
                        <td>Phone</td>
                        <td>{{ .Delivery.Phone }}</td>
                    </tr>
                    <tr>
                        <td>Zip</td>
                        <td>{{ .Delivery.Zip }}</td>
                    </tr>
                    <tr>
                        <td>City</td>
                        <td>{{ .Delivery.City }}</td>
                    </tr>
                    <tr>
                        <td>Address</td>
                        <td>{{ .Delivery.Address }}</td>
                    </tr>
                    <tr>
                        <td>Region</td>
                        <td>{{ .Delivery.Region }}</td>
                    </tr>
                    <tr>
                        <td>Email</td>
                        <td>{{ .Delivery.Email }}</td>
                    </tr>
                    <tr>
                        <td>OrderID</td>
                        <td>{{ .Delivery.OrderID }}</td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</section>
`

var templateFuncs = template.FuncMap{
	"htmlSafe": func(text string) template.HTML {
		return template.HTML(text)
	},
}

func OrderHTMLView(orderDetails OrderDetails) string {
	tmpl, err := template.New("orderTemplate").Funcs(templateFuncs).Parse(orderTemplate)
	if err != nil {
		return fmt.Sprintf("Error parsing template: %v", err)
	}

	var html strings.Builder
	err = tmpl.Execute(&html, orderDetails)
	if err != nil {
		return fmt.Sprintf("Error executing template: %v", err)
	}
	return html.String()
}
