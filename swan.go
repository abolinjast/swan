package main

import (
	"bufio"
	"fmt"
	"os"
    "strings"

	"github.com/fatih/color"
)

func getInput(prompt string, r *bufio.Reader) (string, error) {
    fmt.Print(prompt)
    input, err := r.ReadString('\n')

    return strings.TrimSpace(input), err
}

func main() {
    ipsecConfPath := "/usr/local/etc/ipsec.conf"
    ipsecSecretsPath := "/usr/local/etc/ipsec.secrets"

    // Recieving user input
    ipsecConfig := ""
    secretsConfig := ""
    reader := bufio.NewReader(os.Stdin)
    color.Set(color.FgGreen)    

    fmt.Println("Welcome to the swan helper")
    left, _ := getInput("Enter the Local IP Address: ", reader)
    leftsubnet, _ := getInput("Enter the local subnet: ", reader)
    right, _ := getInput("Enter the Remote IP Address: ", reader)
    rightsubnet, _ := getInput("Enter the remote subnet: ", reader)
    preSharedKey, _ := getInput("Enter the Pre Shared Key that you wish to use: ", reader)
    color.Unset()


    // Config Parsing
    ipsecConfig += fmt.Sprintf(`
# /usr/local/etc/ipsec.conf

# ipsec.conf - strongSwan configuration file

config setup
    charondebug="ike 2, knl 2, cfg 2"  # Debugging level, optional

#
# Connection definitions
conn my-ipsec-tunnel
    authby=secret                   # Use pre-shared key authentication
    left=%v                         # Local IP address
    leftsubnet=%v                   # Local protected subnet
    right=%v                        # Remote IP address
    rightsubnet=%v                  # Remote protected subnet
    ike=aes256-sha256-modp2048      # Phase 1 proposal
    esp=aes256-sha256               # Phase 2 proposal
    keyexchange=ikev2               # Use IKEv2
    auto=start                      # Automatically start the connection
    `, left, leftsubnet, right, rightsubnet) 


    secretsConfig += fmt.Sprintf(`
# /usr/local/etc/ipsec.secrets

%v %v : PSK "%v"

    `,left, right, preSharedKey) 

    // Write configuration to files
    if err := writeToFile(ipsecConfPath, ipsecConfig); err != nil {
        fmt.Printf(color.RedString("Error writing to ipsec.conf: %v", err))
        return
    } 
    if err := writeToFile(ipsecSecretsPath, secretsConfig); err != nil {
        fmt.Printf(color.RedString("Error writing to ipsec.secrets: %v", err))
        return
    }

}

func writeToFile(path, config string) error {
    file, err := os.Create(path)
    if err != nil {
        return err
    }
    defer file.Close()

    _, err = file.WriteString(config)
    return err
}
