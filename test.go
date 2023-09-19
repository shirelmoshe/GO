package main

import (
	"fmt"
	"regexp"
	"sync"
)

const (
	conferenceName     = "ShirelCinema"
	totalTickets       = 200
	maxTicketsPerUser  = 10
	ticketsSoldMessage = "All tickets have been sold"
	notEnoughTickets   = "Not enough tickets available"
)

type User struct {
	FirstName string
	LastName  string
	Email     string
	TicketQty uint
}

func main() {
	var remainingTickets uint = totalTickets
	userMap := make(map[string]User)
	var wg sync.WaitGroup

	fmt.Printf("Welcome to %v cinema application\n", conferenceName)
	fmt.Printf("There are %v tickets in stock\n", remainingTickets)
	fmt.Println("Get your tickets here to attend")

	userInputChan := make(chan User, 2)

	for remainingTickets > 0 {
		wg.Add(1)

		go func() {
			defer wg.Done()
			user := getUserInput()
			userInputChan <- user
		}()

		if len(userInputChan) >= 2 {
			user1 := <-userInputChan
			user2 := <-userInputChan

			if user1.TicketQty > maxTicketsPerUser {
				fmt.Printf("Sorry, %v can buy at most %v tickets at once.\n", user1.FirstName, maxTicketsPerUser)
				continue
			}

			if user1.TicketQty > remainingTickets {
				fmt.Printf("%v, %v\n", notEnoughTickets, user1.FirstName)
				continue
			}

			remainingTickets -= user1.TicketQty
			userMap[user1.Email] = user1

			fmt.Printf("%v %v %v %v ticket(s) purchased successfully.\n", user1.FirstName, user1.LastName, user1.Email, user1.TicketQty)
			fmt.Printf("There are %v tickets left.\n", remainingTickets)

			if user2.TicketQty > maxTicketsPerUser {
				fmt.Printf("Sorry, %v can buy at most %v tickets at once.\n", user2.FirstName, maxTicketsPerUser)
				return
			}

			if user2.TicketQty > remainingTickets {
				fmt.Printf("%v, %v\n", notEnoughTickets, user2.FirstName)
				return
			}

			remainingTickets -= user2.TicketQty
			userMap[user2.Email] = user2

			fmt.Printf("%v %v %v %v ticket(s) purchased successfully.\n", user2.FirstName, user2.LastName, user2.Email, user2.TicketQty)
			fmt.Printf("There are %v tickets left.\n", remainingTickets)
		}
	}

	wg.Wait()

	fmt.Println(ticketsSoldMessage)

	displayAttendees(userMap)
}

func getUserInput() User {
	var user User

	for {
		fmt.Println("Please enter your first name")
		_, err := fmt.Scan(&user.FirstName)
		if err != nil || !isValidName(user.FirstName) {
			fmt.Println("Invalid first name. Please enter a valid name.")
			continue
		}
		break
	}

	for {
		fmt.Println("Please enter your last name")
		_, err := fmt.Scan(&user.LastName)
		if err != nil || !isValidName(user.LastName) {
			fmt.Println("Invalid last name. Please enter a valid name.")
			continue
		}
		break
	}

	for {
		fmt.Println("Please enter your email address")
		_, err := fmt.Scan(&user.Email)
		if err != nil || !isValidEmail(user.Email) {
			fmt.Println("Invalid email address. Please enter a valid email.")
			continue
		}
		break
	}

	for {
		fmt.Printf("How many tickets (maximum %v) do you want to buy?\n", maxTicketsPerUser)
		_, err := fmt.Scan(&user.TicketQty)
		if err != nil {
			fmt.Println("Invalid input. Please enter a valid number of tickets.")
			continue
		}

		if user.TicketQty <= maxTicketsPerUser {
			break
		}

		fmt.Printf("Sorry, you can buy at most %v tickets at once.\n", maxTicketsPerUser)
	}

	return user
}

func isValidName(name string) bool {

	return regexp.MustCompile("^[A-Za-z]+$").MatchString(name)
}

func isValidEmail(email string) bool {

	return regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`).MatchString(email)
}

func displayAttendees(userMap map[string]User) {
	fmt.Println("List of attendees:")
	for email, user := range userMap {
		fmt.Printf("Name: %v %v, Email: %v, Tickets: %v\n", user.FirstName, user.LastName, email, user.TicketQty)
	}
}
