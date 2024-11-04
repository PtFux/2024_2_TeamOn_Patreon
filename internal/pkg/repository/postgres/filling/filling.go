package main

import (
	"context"
	"fmt"
	"github.com/go-park-mail-ru/2024_2_TeamOn_Patreon/internal/pkg/repository/postgres/filling/consts"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
)

type CustomSub struct {
	customName string
	authorName string
	layer      string
}

type Author string
type Layer string
type PostTitle string
type User string

var (
	CustomSubToAuthor = map[string]CustomSub{}
	Subscriptions     = map[string]CustomSub{}

	Posts = make(map[Author]map[Layer][]PostTitle)
	Users = make(map[Author]map[Layer][]User)
)

func main() {
	// Укажите, сколько пользователей нужно создать
	n := consts.COUNT_USER

	// Параметры подключения к PostgreSQL
	dbHost := "postgres"
	dbPort := 5432
	dbUser := "admin"
	dbPassword := "adminpass"
	dbName := "testdb"

	// Формируем строку подключения
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUser, dbPassword, dbHost, dbPort, dbName)

	//dbURL := "postgres://your_user:your_password@localhost:5432/your_dbname?sslmode=disable"
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("Ошибка создания пула подключений: %v", err)
	}
	defer pool.Close()

	if err := createUsers(context.Background(), pool, n); err != nil {
		log.Fatalf("Ошибка при создании пользователей: %v", err)
	}

	fmt.Printf("Создано %d пользователей\n", n)

	if err := createAuthors(context.Background(), pool, n); err != nil {
		log.Fatalf("Ошибка при создании авторов: %v", err)
	}

	fmt.Printf("Создано %d авторов \n", n)

	if err := createPage(context.Background(), pool, n); err != nil {
		log.Fatalf("Ошибка при создании страниц: %v", err)
	}

	fmt.Printf("Создано %d страниц \n", n)

	if err := createCustomSubscriptions(context.Background(), pool, n); err != nil {
		log.Fatalf("Ошибка при создании кастомных подписок: %v", err)
	}

	fmt.Printf("Создано %d кастомных подписок \n", n)

	if err := createSubscriptions(context.Background(), pool, n); err != nil {
		log.Fatalf("Ошибка при создании подписок: %v", err)
	}

	fmt.Printf("Создано %d подписок \n", n*n)

	if err := createPosts(context.Background(), pool, n); err != nil {
		log.Fatalf("Ошибка при создании постов: %v", err)
	}

	fmt.Printf("Создано %d постов \n", n*n)
	index, err := createPostLikes(context.Background(), pool, n)
	if err != nil {
		log.Fatalf("Ошибка при создании лайков на посты: %v", err)
	}

	fmt.Printf("Создано %d лайков \n", index)
}

func createUsers(ctx context.Context, pool *pgxpool.Pool, n int) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("не удалось получить соединение из пула: %v", err)
	}
	defer conn.Release()

	batch := &pgx.Batch{}

	// Подготавливаем данные для пользователей и связанных записей
	for i := 0; i < n; i++ {
		userID := uuid.New()
		username := fmt.Sprintf(consts.USERNAME, i+1)
		email := fmt.Sprintf(consts.EMAIL_DOMAIN_NAME, username)

		// Генерация хеша пароля
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(consts.PASSWORD), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("ошибка хеширования пароля: %v", err)
		}

		// Запрос на добавление пользователя
		batch.Queue(`
            INSERT INTO People (user_id, username, email, role_id, hash_password) 
            VALUES ($1, $2, $3, (SELECT role_id FROM Role WHERE role_default_name = 'Reader'), $4)`,
			userID, username, email, passwordHash)

	}

	// Выполнение батч-запроса
	br := conn.Conn().SendBatch(ctx, batch)
	defer br.Close()

	// Проверка результатов выполнения батча
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("ошибка выполнения батч-запроса на шаге %d: %v", i+1, err)
		}
	}

	return nil
}

func createAuthors(ctx context.Context, pool *pgxpool.Pool, n int) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("не удалось получить соединение из пула: %v", err)
	}
	defer conn.Release()

	batch := &pgx.Batch{}

	// Подготавливаем данные для пользователей и связанных записей
	for i := 0; i < n; i++ {
		userID := uuid.New()
		username := fmt.Sprintf(consts.AUTHOR_NAME, i+1)
		email := fmt.Sprintf(consts.EMAIL_DOMAIN_NAME, username)

		// Генерация хеша пароля
		passwordHash, err := bcrypt.GenerateFromPassword([]byte(consts.PASSWORD), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("ошибка хеширования пароля: %v", err)
		}

		// Запрос на добавление пользователя
		batch.Queue(`
            INSERT INTO People (user_id, username, email, role_id, hash_password) 
            VALUES ($1, $2, $3, (SELECT role_id FROM Role WHERE role_default_name = 'Author'), $4)`,
			userID, username, email, passwordHash)

	}

	// Выполнение батч-запроса
	br := conn.Conn().SendBatch(ctx, batch)
	defer br.Close()

	// Проверка результатов выполнения батча
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("ошибка выполнения батч-запроса на шаге %d: %v", i+1, err)
		}
	}

	return nil
}

func createPage(ctx context.Context, pool *pgxpool.Pool, n int) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("не удалось получить соединение из пула: %v", err)
	}
	defer conn.Release()

	batch := &pgx.Batch{}

	// Подготавливаем данные для пользователей и связанных записей
	for i := 0; i < n; i++ {
		username := fmt.Sprintf(consts.AUTHOR_NAME, i+1)
		about := fmt.Sprintf(consts.PAGE_INFO, username)

		// Запрос на добавление пользователя
		batch.Queue(`
INSERT INTO Page (page_id, user_id, info) VALUES
(gen_random_uuid(), (SELECT user_id FROM People WHERE username = $1), $2)
  `,
			username, about)
	}

	// Выполнение батч-запроса
	br := conn.Conn().SendBatch(ctx, batch)
	defer br.Close()

	// Проверка результатов выполнения батча
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("ошибка выполнения батч-запроса на шаге %d: %v", i+1, err)
		}
	}

	return nil
}

func createCustomSubscriptions(ctx context.Context, pool *pgxpool.Pool, n int) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("не удалось получить соединение из пула: %v", err)
	}
	defer conn.Release()

	batch := &pgx.Batch{}

	// Подготавливаем данные для пользователей и связанных записей
	for i := 0; i < n; i++ {
		authorName := fmt.Sprintf(consts.AUTHOR_NAME, i+1)
		username := fmt.Sprintf(consts.AUTHOR_NAME, i+1)
		customSub := fmt.Sprintf(consts.CUSTOM_NAME, i+1, authorName)
		layer := i % 4

		cost := consts.CUSTOM_COST * layer

		CustomSubToAuthor[customSub] = CustomSub{layer: string(layer), authorName: username, customName: customSub}

		// Запрос на добавление пользователя
		batch.Queue(`
    INSERT INTO Custom_Subscription (custom_subscription_id, author_id, custom_name, cost, info, subscription_layer_id) VALUES
    (gen_random_uuid(), (SELECT user_id FROM People WHERE username = $1), $2, $3, $4, (SELECT subscription_layer_id FROM Subscription_Layer WHERE layer = $5) )
  `,
			username, customSub, cost, consts.CUSTOM_SUB_INFO, layer)

	}

	// Выполнение батч-запроса
	br := conn.Conn().SendBatch(ctx, batch)
	defer br.Close()

	// Проверка результатов выполнения батча
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("ошибка выполнения батч-запроса на шаге %d: %v", i+1, err)
		}
	}

	return nil
}

func createSubscriptions(ctx context.Context, pool *pgxpool.Pool, n int) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("не удалось получить соединение из пула: %v", err)
	}
	defer conn.Release()

	batch := &pgx.Batch{}

	// Подготавливаем данные для пользователей и связанных записей
	for i := 0; i < n; i++ {
		authorName := fmt.Sprintf(consts.AUTHOR_NAME, i+1)
		customSub := fmt.Sprintf(consts.CUSTOM_NAME, i+1, authorName)
		layer := i % 4
		for j := 0; j < n; j++ {
			username := fmt.Sprintf(consts.USERNAME, j+1)

			Subscriptions[username] = CustomSubToAuthor[customSub]

			if Users[Author(authorName)] == nil {
				Users[Author(authorName)] = make(map[Layer][]User)
			}
			Users[Author(authorName)][Layer(string(layer))] = append(Users[Author(authorName)][Layer(string(layer))], User(username))

			// Запрос на добавление пользователя
			batch.Queue(`
  INSERT INTO Subscription (subscription_id, user_id, custom_subscription_id, started_date, finished_date) VALUES
    (gen_random_uuid(), (SELECT user_id FROM People WHERE username = $1), (SELECT custom_subscription_id FROM Custom_Subscription WHERE custom_name = $2), NOW(), NOW() + INTERVAL '30 days')
`,
				username, customSub)
		}
	}

	// Выполнение батч-запроса
	br := conn.Conn().SendBatch(ctx, batch)
	defer br.Close()

	// Проверка результатов выполнения батча
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("ошибка выполнения батч-запроса на шаге %d: %v", i+1, err)
		}
	}

	return nil
}

func createPosts(ctx context.Context, pool *pgxpool.Pool, n int) error {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("не удалось получить соединение из пула: %v", err)
	}
	defer conn.Release()

	batch := &pgx.Batch{}

	// Подготавливаем данные для пользователей и связанных записей
	for i := 0; i < n; i++ {
		authorName := fmt.Sprintf(consts.AUTHOR_NAME, i+1)
		about := fmt.Sprintf(consts.ABOUT, authorName)
		layer := i % 4

		for j := 0; j < n; j++ {
			title := fmt.Sprintf(consts.TITLE, i, j, authorName)

			if Posts[Author(authorName)] == nil {
				Posts[Author(authorName)] = make(map[Layer][]PostTitle)
			}
			Posts[Author(authorName)][Layer(string(layer))] = append(Posts[Author(authorName)][Layer(string(layer))], PostTitle(title))

			// Запрос на добавление пользователя
			batch.Queue(`
INSERT INTO Post (post_id, user_id, title, about, subscription_layer_id) VALUES
    (gen_random_uuid(), (SELECT user_id FROM People WHERE username = $1), $2, $3 , (SELECT subscription_layer_id FROM Subscription_Layer WHERE layer = $4))
`,
				authorName, title, about, layer)
		}
	}

	// Выполнение батч-запроса
	br := conn.Conn().SendBatch(ctx, batch)
	defer br.Close()

	// Проверка результатов выполнения батча
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return fmt.Errorf("ошибка выполнения батч-запроса на шаге %d: %v", i+1, err)
		}
	}

	return nil
}

func createPostLikes(ctx context.Context, pool *pgxpool.Pool, n int) (index int, err error) {
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return 0, fmt.Errorf("не удалось получить соединение из пула: %v", err)
	}
	defer conn.Release()

	batch := &pgx.Batch{}
	index = 0

	for authorName, layerPosts := range Posts {
		for layer, posts := range layerPosts {
			for _, title := range posts {

				userNames := Users[Author(authorName)][Layer(string(layer))]
				for _, username := range userNames {
					index = index + 1

					top := index % 6
					if top == 0 {
						continue
					}
					if rand.Int()%133%7 == 0 {
						continue
					}

					// Запрос на добавление пользователя
					batch.Queue(`
INSERT INTO Like_Post (like_post_id, post_id, user_id, posted_date) VALUES
    (gen_random_uuid(), (SELECT post_id FROM Post WHERE title = $1), (SELECT user_id FROM People WHERE username = $2), NOW())
   `,
						string(title), string(username))
				}
			}
		}
	}

	// Выполнение батч-запроса
	br := conn.Conn().SendBatch(ctx, batch)
	defer br.Close()

	// Проверка результатов выполнения батча
	for i := 0; i < batch.Len(); i++ {
		_, err := br.Exec()
		if err != nil {
			return 0, fmt.Errorf("ошибка выполнения батч-запроса на шаге %d: %v", i+1, err)
		}
	}

	return index, nil
}