package main

import (
	"fmt"
	"log"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func main() {
	dbPath := "./data/drama_generator.db"

	// Check if file exists
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("Error: Không tìm thấy database tại %s. Xin hãy chạy script này ở thư mục gốc của project (huobao-drama).", dbPath)
	}

	fmt.Println("Đang kết nối tới Database...")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("Không thể kết nối Database: %v", err)
	}

	fmt.Println("1. Đang sửa các trạng thái video Upscale hoàn thành...")
	result1 := db.Exec("UPDATE video_generations SET status = 'upscaled' WHERE status = 'completed' AND is_upscaled = 1")
	if result1.Error != nil {
		log.Printf("❌ Lỗi Query 1: %v\n", result1.Error)
	} else {
		fmt.Printf("✅ Đã vá thành công %d bản ghi Upscale hoàn tất (chuyển sang 'upscaled').\n", result1.RowsAffected)
	}

	fmt.Println("2. Đang sửa các trạng thái Upscale thất bại...")
	result2 := db.Exec("UPDATE video_generations SET status = 'upscale_failed' WHERE status = 'failed' AND (video_url IS NOT NULL OR local_path IS NOT NULL)")
	if result2.Error != nil {
		log.Printf("❌ Lỗi Query 2: %v\n", result2.Error)
	} else {
		fmt.Printf("✅ Đã cứu thành công %d bản ghi Upscale thất bại (chuyển sang 'upscale_failed').\n", result2.RowsAffected)
	}

	fmt.Println("\n🎉 Chúc mừng! Database đã được nâng cấp theo hệ thống trạng thái mới.")
	fmt.Println("Bạn có thể F5 lại web để xem thành quả!")
}
