package main

import (
	"ProtoDemo/pb"
	"fmt"
	"google.golang.org/protobuf/proto"
	"log"
	"os"
)

func main() {
	UnmarshalBook("serialize.txt")
}

func CreateBook() {
	//创建一本书
	book := &pb.Book{}
	book.Title = "你当像鸟飞往你的山"
	price := &pb.Price{MarketPrice: 2, SalePrice: 3}
	book.Price = price
	book.Author = ".."
	info := &pb.Information{
		Introduction: "比尔盖茨推荐",
		Publication:  "人民出版社",
		IsPublic:     false,
		BookTable:    0,
	}
	book.Infos = make(map[string]*pb.Information)
	book.Infos["info"] = info
	book.Tables = append(book.Tables, "literature")
	book.BookTable = pb.BookTable_BOOK_TABLE_STORY

	//将book序列化存入文件中
	var fileName string = "serialize.txt"
	bytes, err := proto.Marshal(book)
	if err != nil {
		log.Fatalf("marshal failed, err: %v\n", err)
		return
	}
	//创建文件
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("create file failed,err: %v\n", err)
		return
	}
	defer file.Close()
	_, err = file.Write(bytes)
	fmt.Printf("book: %v\n", book)
}

func UnmarshalBook(str string) {
	b, err := os.ReadFile(str)
	if err != nil {
		log.Fatalf("read file failed, err: %v\n", err)
		return
	}
	var book pb.Book
	err = proto.Unmarshal(b, &book)
	if err != nil {
		log.Fatalf("unmarshal failed, err:%v\n", err)
		return
	}
	fmt.Println(book)
}
