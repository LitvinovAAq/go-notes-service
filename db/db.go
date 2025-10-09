package db

import (
    "database/sql"        
    "fmt"                 
    "os"                  

    _ "github.com/lib/pq" 
)

func getEnv(key, def string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return def
}

func dsnFromEnv() string {
    host := getEnv("PG_HOST", "localhost")       
    port := getEnv("PG_PORT", "5432")      
    user := getEnv("PG_USER", "postgres")  
    pass := getEnv("PG_PASSWORD", "admin") 
    name := getEnv("PG_DB", "GoStud")      
    ssl  := getEnv("PG_SSLMODE", "disable")

    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        host, port, user, pass, name, ssl,
    )
}

func GetDB() (*sql.DB, error) {
    return sql.Open("postgres", dsnFromEnv())
}
