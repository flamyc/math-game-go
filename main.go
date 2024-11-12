package main

import (
	"encoding/json"
	"fmt"
	"math-game/domain"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
)

const (
	totalPoints       int = 100
	pointsPerQuestion int = 100
)

var id uint64 = 1

func main() {
	var userList []domain.User
	userList = getUsers()
	fmt.Println("Вітаємо у грі!")

	for {
		menu()
		choice := ""
		fmt.Scan(&choice)
		switch choice {
		case "1":
			user := play()
			userList = append(userList, user)
			sortAndSave(userList)
		case "2":
			for _, user := range userList {
				fmt.Printf("[%v]\tName: %s\tTime: %v\n",
					user.Id, user.Name, user.TimeSpent)
			}
		case "3":
			return
		default:
			fmt.Println("Оберіть опцію 1, 2 або 3.")
		}
	}
}

func menu() {
	fmt.Println("1. Почати гру")
	fmt.Println("2. Рейтинг")
	fmt.Println("3. Вийти")
}

func play() domain.User {
	for i := 3; i >= 1; i-- {
		fmt.Printf("Гра почнеться через: %v\n", i)
		time.Sleep(time.Second)
	}

	startTime := time.Now()
	myPoints := 0
	for myPoints < totalPoints {
		x, y := rand.Intn(100), rand.Intn(100)

		fmt.Printf("%v + %v = ", x, y)

		ans := ""
		fmt.Scan(&ans)

		ansInt, err := strconv.Atoi(ans)
		if err != nil {
			fmt.Println("Введіть число.")
			continue
		}
		if ansInt == x+y {
			myPoints += pointsPerQuestion
			fmt.Printf("Отримано %v балів. Залишилося набрати %v\n", pointsPerQuestion, totalPoints-myPoints)
		} else {
			fmt.Println("Неправильна відповідь.")
		}
	}

	timeSpent := time.Since(startTime)

	fmt.Println("Вітаю! Остаточний час: ", timeSpent)
	fmt.Println("Введіть ім'я: ")

	name := ""
	fmt.Scan(&name)
	user := domain.User{
		Id:        id,
		Name:      name,
		TimeSpent: timeSpent,
	}
	id++

	return user
}

func sortAndSave(users []domain.User) {
	sort.SliceStable(users, func(i, j int) bool {
		return users[i].TimeSpent < users[j].TimeSpent
	})

	file, err := os.OpenFile("users.json", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		fmt.Printf("sortAndSave -> os.OpenFile: %s", err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("sortAndSave -> file.Close: %s", err)
		}
	}(file)

	encoder := json.NewEncoder(file)
	err = encoder.Encode(users)
	if err != nil {
		fmt.Printf("sortAndSave -> encoder.Encode: %s", err)
	}
}

func getUsers() []domain.User {
	file, err := os.Open("users.json")
	if err != nil {
		fmt.Printf("getUsers -> os.Open: %s", err)
		return nil
	}

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Printf("sortAndSave -> file.Close: %s", err)
		}
	}(file)

	var users []domain.User
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&users)
	if err != nil {
		fmt.Printf("getUsers -> decoder.Decode: %s", err)
		return nil
	}

	return users
}
