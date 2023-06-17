package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

// User modelini temsil eden bir struct tanımlayalım
type User struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func main() {
	// MariaDB bağlantı bilgilerini ayarlayın
	db, err := sql.Open("mysql", "root:password@tcp(db:3306)/database")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping veritabanına erişimi kontrol eder
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// HTTP endpoint'lerini tanımlayalım
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Tüm kullanıcıları getir
			users, err := getUsers(db)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// JSON formatında kullanıcıları döndür
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users)

		case http.MethodPost:
			// Yeni bir kullanıcı kaydet
			var user User
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = saveUser(db, &user)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Kaydedilen kullanıcıyı döndür
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Web sunucusunu başlat
	fmt.Println("Web server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// Kullanıcıları veritabanından getiren bir yardımcı fonksiyon
func getUsers(db *sql.DB) ([]User, error) {
	rows, err := db.Query("SELECT * FROM User")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

// Kullanıcıyı veritabanına kaydeden bir yardımcı fonksiyon
func saveUser(db *sql.DB, user *User) error {
	_, err := db.Exec("INSERT INTO User (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		return err
	}

	return nil
}
