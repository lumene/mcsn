package main

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Liza-Developer/mcapi2"
)

var (
	sends       []time.Time
	recvs       []time.Time
	statuscodes []string
	email       []string
)

func sendAuto(option string, delay float64) {
	leng := 0
	for _, name := range names {

		if useAuto {
			delay = AutoOffset(false)
		}

		if bearers.Bearers == nil || len(bearers.Bearers) == 0 {
			fmt.Println("Attempting to reauth accounts..")
			authAccs()
		}

		dropTime := mcapi2.DropTime(name)

		fmt.Printf("    Name: %v\n   Delay: %v\nDroptime: %v\n\n", name, delay, dropTime)

		mcapi2.PreSleep(dropTime)

		payload := bearers.CreatePayloads(name)

		mcapi2.Sleep(dropTime, delay)

		fmt.Println()

		for f, accType := range bearers.AccountType {
			switch accType {
			case "Giftcard":
				leng = 6
			case "Microsoft":
				leng = 2
			}

			for i := 0; i < leng; {
				go func() {
					send, recv, status := payload.SocketSending(int64(f))

					sends = append(sends, send)
					recvs = append(recvs, recv)
					statuscodes = append(statuscodes, status)
				}()
				i++
				time.Sleep(40 * time.Microsecond)
			}

			email = append(email, strings.Split(AccountsVer[f], ":")[0])

		}

		time.Sleep(500 * time.Millisecond)

		sort.Slice(sends, func(i, j int) bool {
			return sends[i].Before(sends[j])
		})

		for i, accountType := range payload.AccountType {

			switch accountType {
			case "Giftcard":
				leng = 6
			case "Microsoft":
				leng = 2
			}

			for f := 0; f < leng; {
				if statuscodes[f] != "200" {
					content += fmt.Sprintf("- [%v] Sent @ %v | Recv @ %v - %v\n", statuscodes[f], formatTime(sends[f]), formatTime(recvs[f]), email[i])
					fmt.Printf("[%v] Sent @ %v | Recv @ %v - %v\n", statuscodes[f], formatTime(sends[f]), formatTime(recvs[f]), email[i])
				} else {
					content += fmt.Sprintf("+ [%v] Succesfully sniped %v | %v\n", statuscodes[f], name, email[i])
					fmt.Printf("[%v] Recv @ %v | %v\n", statuscodes[f], formatTime(recvs[f]), email[i])
					sendInfo.ChangeSkin(nil, bearers.Bearers[i])
					sendInfo.SendWebhook(jsonValue(embeds{Content: nil, Embeds: []embed{{Description: fmt.Sprintf("```diff\n%v\n```", content), Color: nil}}}))
					bearers.Bearers = remove(bearers.Bearers, bearers.Bearers[i])
					bearers.AccountType = remove(bearers.AccountType, bearers.AccountType[i])
					payload.Payload = remove(payload.Payload, payload.Payload[i])
				}
				f++
			}
		}

		content = `
+    __  ______________ _   __
-   /  |/  / ____/ ___// | / /
+  / /|_/ / /    \__ \/  |/ / 
- / /  / / /___ ___/ / /|  /  
+/_/  /_/\____//____/_/ |_/

`

	}
}

func singlesniper(name string, delay float64) {
	var leng int

	dropTime = mcapi2.DropTime(name)

	fmt.Printf(`    Name: %v
   Delay: %v
Droptime: %v

`, name, delay, formatTime(time.Unix(dropTime, 0)))

	mcapi2.PreSleep(dropTime)

	payload := bearers.CreatePayloads(name)

	mcapi2.Sleep(dropTime, delay)

	fmt.Println()

	for f, accountType := range payload.AccountType {
		switch accountType {
		case "Giftcard":
			leng = 6
		case "Microsoft":
			leng = 2
		}

		for i := 0; i < leng; {
			go func() {
				send, recv, status := payload.SocketSending(int64(f))

				sends = append(sends, send)
				recvs = append(recvs, recv)
				statuscodes = append(statuscodes, status)
			}()
			i++
			time.Sleep(40 * time.Microsecond)
		}

		email = append(email, strings.Split(AccountsVer[f], ":")[0])
	}

	time.Sleep(500 * time.Millisecond)

	sort.Slice(sends, func(i, j int) bool {
		return sends[i].Before(sends[j])
	})

	var num int

	for i, accountType := range payload.AccountType {

		switch accountType {
		case "Giftcard":
			leng = 6
		case "Microsoft":
			leng = 2
		}

		for f := 0; f < leng; {
			if statuscodes[f] != "200" {
				content += fmt.Sprintf("- [%v] Sent @ %v | Recv @ %v - %v\n", statuscodes[num], formatTime(sends[num]), formatTime(recvs[num]), email[i])
				fmt.Printf("[%v] Sent @ %v | Recv @ %v - %v\n", statuscodes[num], formatTime(sends[num]), formatTime(recvs[num]), email[i])
				num++
			} else {
				content += fmt.Sprintf("+ [%v] Succesfully sniped %v | %v\n", statuscodes[num], name, email[i])
				fmt.Printf("[%v] Recv @ %v | %v\n", statuscodes[num], formatTime(recvs[num]), email[i])
				sendInfo.ChangeSkin(nil, bearers.Bearers[i])
				sendInfo.SendWebhook(jsonValue(embeds{Content: nil, Embeds: []embed{{Description: fmt.Sprintf("```diff\n%v\n```", content), Color: nil}}}))
				num++
			}
			f++
		}
	}
}
