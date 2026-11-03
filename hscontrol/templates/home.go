package templates

import (
	"github.com/chasefleming/elem-go"
)

// Home renders the landing page at /: what this service is and how to put a
// device on it. serverURL is the control server's public URL. Linux and macOS
// run our own client (thelo); phones and Windows use the Tailscale app, which
// works with Thelo until we ship our own.
func Home(serverURL string) *elem.Element {
	return page(
		ProductName,
		H1(elem.Text(ProductName+", the "+CompanyName+" private network")),
		P(elem.Text("Sign in with your Simlah account and passkey "+
			"to put a device on your "+NetName+"; access is granted by a "+CompanyName+" administrator.")),
		H2(elem.Text("Join a device")),
		Ol(
			elem.Li(nil,
				elem.Text("On Linux or macOS, install "+ProductName+" and point it at this server:"),
				codeBlockText("thelo login --login-server "+serverURL),
			),
			elem.Li(nil,
				elem.Text("On an iPhone, iPad, Android phone or Windows, use the Tailscale app, "+
					"which works with "+ProductName+": follow the "),
				A("/apple", elem.Text("Apple")),
				elem.Text(" or "),
				A("/windows", elem.Text("Windows")),
				elem.Text(" instructions."),
			),
			elem.Li(nil,
				elem.Text("Open the link it gives you within 15 minutes and sign in with your passkey."),
			),
		),
		P(elem.Text("Devices sign in again every 30 days. If your access is withdrawn, your devices "+
			"leave the network within minutes.")),
	)
}
