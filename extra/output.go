package extra

import (
	"fmt"
	"github.com/gookit/color"
)

func ShowMenu() {

	for i, item := range MainMenu {
		if i < 9 {
			fmt.Printf("\033[36m[\033[33m0%d\033[36m] \033[36m>> \033[37m%s\033[0m\n", i+1, item)
		} else {
			fmt.Printf("\033[36m[\033[33m%d\033[36m] \033[36m>> \033[37m%s\033[0m\n", i+1, item)
		}
	}
}

func ShowLogo() {
	logo := "\n███████╗████████╗ █████╗ ██████╗     ██╗      █████╗ ██████╗ ███████╗\n██╔════╝╚══██╔══╝██╔══██╗██╔══██╗    ██║     ██╔══██╗██╔══██╗██╔════╝\n███████╗   ██║   ███████║██████╔╝    ██║     ███████║██████╔╝███████╗\n╚════██║   ██║   ██╔══██║██╔══██╗    ██║     ██╔══██║██╔══██╗╚════██║\n███████║   ██║   ██║  ██║██║  ██║    ███████╗██║  ██║██████╔╝███████║\n╚══════╝   ╚═╝   ╚═╝  ╚═╝╚═╝  ╚═╝    ╚══════╝╚═╝  ╚═╝╚═════╝ ╚══════╝\n                                                                     "

	color.Cyan.Printf(logo + "\n")

}

func ShowDevInfo() {
	fmt.Println("\033[36mVERSION: \033[33m1.3\033[33m")
	fmt.Println("\033[36mDEV: \033[33mhttps://t.me/StarLabsTech\033[33m")
	fmt.Println("\033[36mGitHub: \033[33mhttps://github.com/0xStarLabs/StarLabs-Twitter\033[33m")
	fmt.Println("\033[36mDONATION EVM ADDRESS: \033[33m0x620ea8b01607efdf3c74994391f86523acf6f9e1\033[0m")
	fmt.Println()
}

var MainMenu = []string{
	"Follow",
	"Retweet",
	"Like",
	"Tweet",
	"Tweet with Picture",
	"Quote Tweet",
	"Comment",
	"Comment with Picture",
	"Change Description",
	"Change Username",
	"Change Name",
	"Change Background",
	"Change Password",
	"Change Birthdate",
	"Change Location",
	"Change profile Picture",
	"Check if account is valid",
	//"Unfreeze Accounts",
	"Mutual Subscription",
}
