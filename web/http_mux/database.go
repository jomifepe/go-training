package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func configureDB() *sql.DB {
	fmt.Println("DB - Connecting to the database...")
	psqlInfo := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DATABASE_HOST"),
		os.Getenv("DATABASE_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("DB - Successfully connected to the database!")
	migrateDatabase(db)
	return db
}

func migrateDatabase(instance *sql.DB) {
	fmt.Println("DB - Running migrations...")
	driver, err := postgres.WithInstance(instance, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", os.Getenv("POSTGRES_DB"), driver)

	if err != nil {
		log.Fatal(err)
	}

	migErr := m.Up()
	if migErr == nil {
		fmt.Println("DB - Migrations completed successfully!")
		return
	}
	if migErr == migrate.ErrNoChange {
		fmt.Println("DB - Nothing to migrate...")
		return
	}
	log.Fatal(migErr)
}

func dbGetAllTasks() (TaskCollection, error) {
	rows, err := db.Query(`SELECT * FROM "Tasks"`)
	if err != nil {
		fmt.Println("DB - Couldn't get all tasks", err)
		return []Task{}, err
	}

	var tasks []Task
	defer rows.Close()
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Description, &t.Completed, &t.CreatedAt); err != nil {
			fmt.Println("DB - Couldn't read database row:", rows)
			return []Task{}, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func dbGetTask(id int) (Task, error) {
	var t Task
	row := db.QueryRow(`SELECT * FROM "Tasks" WHERE id = $1`, id)
	err := row.Scan(&t.ID, &t.Description, &t.Completed, &t.CreatedAt)
	if err != nil {
		fmt.Printf("DB - Couldn't get task with id %v | %v\n", id, err)
		return Task{}, err
	}

	return t, nil
}

func dbCreateTask(t Task) (int, error) {
	newID := 0
	if err := db.QueryRow(`
		INSERT INTO "Tasks" (description) 
		VALUES ($1)
		RETURNING id`, t.Description).Scan(&newID); err != nil {

		fmt.Printf("DB - Failed to create task: %+v | %v", t, err)
		return -1, err
	}

	t.ID = newID
	fmt.Printf("DB - Created new task: %+v", t)
	return newID, nil
}

func dbUpdateTask(t Task) (Task, error) {
	res, execErr := db.Exec(`
		UPDATE "Tasks"
		SET description = $2, completed = $3
		WHERE id = $1;`, t.ID, t.Description, t.Completed)
	if execErr != nil {
		fmt.Printf("DB - Couldn't update task with id %v | %v\n", t.ID, execErr)
		return t, execErr
	}

	if _, err := res.RowsAffected(); err != nil {
		panic(err)
	}

	fmt.Printf("DB - Updated existing Task: %+v", t)
	return t, nil
}

func dbDeleteTask(id int) error {
	if _, err := db.Exec(`DELETE FROM "Tasks" WHERE id = $1`, id); err != nil {
		fmt.Printf("DB - Couldn't delete task with id %v | %v", id, err)
		return err
	}
	return nil
}
