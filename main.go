package main

import (
	"flag"
	"fmt"
	"go-zero2/initialize"
	"go-zero2/system/system_models"
)

type Option struct {
	DB bool
}

func main() {
	var opt Option
	flag.BoolVar(&opt.DB, "db", false, "db")
	flag.Parse()
	if opt.DB {
		db := initialize.InitMysql("root:zdsmE2J2GMxn6tHX@tcp(47.103.52.101:3306)/system?charset=utf8mb4&parseTime=True&loc=Local")
		err := db.AutoMigrate(
			&system_models.SubSystem{},
			&system_models.Authorization{},
			&system_models.Lease{},
			&system_models.Version{})
		if err != nil {
			fmt.Printf("数据库迁移失败 err:%s", err.Error())
		}
		fmt.Println("数据库迁移成功!")
	}
}

/*
CREATE TABLE sub_system (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    server_address VARCHAR(255) NOT NULL,
    database_config TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE authorization (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sub_system_id BIGINT NOT NULL,
    token VARCHAR(255) NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (sub_system_id) REFERENCES sub_system(id)
);

CREATE TABLE lease (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sub_system_id BIGINT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    status VARCHAR(50) NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (sub_system_id) REFERENCES sub_system(id)
);

CREATE TABLE version (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    sub_system_id BIGINT NOT NULL,
    version_number VARCHAR(50) NOT NULL,
    release_notes TEXT,
    release_date DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (sub_system_id) REFERENCES sub_system(id)
);
*/
