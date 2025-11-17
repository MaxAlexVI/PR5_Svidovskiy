package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// .env не обязателен; если файла нет — ошибка игнорируется
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// fallback — прямой DSN в коде (только для учебного стенда!)
		dsn = "postgres://postgres:1234@localhost:5433/todo?sslmode=disable"
	}

	db, err := openDB(dsn)
	if err != nil {
		log.Fatalf("openDB error: %v", err)
	}
	defer db.Close()

	repo := NewRepo(db)

	// 1) Вставим пару задач через массовую вставку
	ctxCreate, cancelCreate := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelCreate()

	titles := []string{"Сделать ПЗ №5", "Купить кофе", "Проверить отчёты", "Изучить Go", "Написать документацию"}
	err = repo.CreateMany(ctxCreate, titles)
	if err != nil {
		log.Fatalf("CreateMany error: %v", err)
	}
	log.Printf("Inserted %d tasks", len(titles))

	// 2) Прочитаем список задач
	ctxList, cancelList := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelList()

	tasks, err := repo.ListTasks(ctxList)
	if err != nil {
		log.Fatalf("ListTasks error: %v", err)
	}

	// 3) Напечатаем
	fmt.Println("=== Tasks ===")
	for _, t := range tasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	// 4) Печатаем DoneList - таски с Done = true
	ctxDone, doneList := context.WithTimeout(context.Background(), 3*time.Second)
	defer doneList()

	doneTasks, err := repo.ListDone(ctxDone, true)
	if err != nil {
		log.Fatalf("ListDone error: %v", err)
	}

	fmt.Println("\n=== Done Tasks List ===")
	for _, t := range doneTasks {
		fmt.Printf("#%d | %-24s | done=%-5v | %s\n",
			t.ID, t.Title, t.Done, t.CreatedAt.Format(time.RFC3339))
	}

	ctxFindId, cancelFindId := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancelFindId()

	if len(tasks) > 0 {
		task, err := repo.FindById(ctxFindId, tasks[0].ID)
		if err != nil {
			log.Fatalf("FindByID error: %v", err)
		}

		fmt.Println("\n=== Task by ID ===")
		fmt.Printf("ID: %d\nTitle: %s\nDone: %v\nCreated: %s\n",
			task.ID, task.Title, task.Done, task.CreatedAt.Format(time.RFC3339))
	}

}
