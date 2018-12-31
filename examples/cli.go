package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	jusibe "github.com/azeezolaniran2016/jusibe-go"
)

var (
	accessToken string
	publicKey   string
)

func main() {
	flag.StringVar(&accessToken, "access_token", "", "Jusibe access_token")
	flag.StringVar(&publicKey, "public_key", "", "Jusibe public_key")
	flag.Parse()

	config := &jusibe.Config{
		PublicKey:   publicKey,
		AccessToken: accessToken,
	}

	jusibe, err := jusibe.New(config)
	if err != nil {
		log.Fatalf(err.Error())
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		printHelp()
		line := readLine(scanner, "Enter Command: ")
		switch line {
		case "get_credits":
			{
				getCredits(jusibe)
			}
		case "send_sms":
			{
				sendSMS(jusibe, scanner)
			}
		case "delivery_status":
			{
				checkDeliveryStatus(jusibe, scanner)
			}
		case "exit":
			{
				os.Exit(0)
			}
		default:
			{
				printHelp()
			}
		}
		fmt.Println()
	}
}

func getCredits(jb *jusibe.Jusibe) {
	fmt.Println("Fetching credits...")

	res, _, err := jb.CheckSMSCredits(context.Background())

	if err != nil {
		log.Printf("Failed to check credits - %s\n", err.Error())
		return
	}

	fmt.Printf("Credits => %s\n", res.SMSCredits)
}

func sendSMS(jb *jusibe.Jusibe, scanner *bufio.Scanner) {
	to := readLine(scanner, "Enter To: ")
	from := readLine(scanner, "Enter From: ")
	message := readLine(scanner, "Enter Message: ")

	fmt.Println("Sending SMS....")

	res, _, err := jb.SendSMS(context.Background(), to, from, message)

	if err != nil {
		log.Printf("Failed to send SMS - %s\n", err.Error())
		return
	}

	fmt.Printf("Response => %+v", res)
}

func checkDeliveryStatus(jb *jusibe.Jusibe, scanner *bufio.Scanner) {
	messageID := readLine(scanner, "Enter MessageID: ")

	fmt.Println("Fetching delivery status...")

	res, _, err := jb.CheckSMSDeliveryStatus(context.Background(), messageID)
	if err != nil {
		log.Printf("Failed to check SMS Delivery status - %s\n", err.Error())
		return
	}

	fmt.Printf("Response => %+v", res)
}

func readLine(scanner *bufio.Scanner, prompt string) (line string) {
	if prompt != "" {
		fmt.Print(prompt)
	}
	scanner.Scan()
	line = scanner.Text()

	return
}

func printHelp() {
	fmt.Println("Enter API Method to execute:")
	fmt.Println("\tget_credits - View remaining credits")
	fmt.Println("\tsend_sms - Send SMS")
	fmt.Println("\tdelivery_status - Check SMS delivery status")
	fmt.Println("Enter exit to Quit")
	fmt.Println()
}
