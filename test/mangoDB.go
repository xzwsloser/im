package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Person struct {
	Name string `bson:"name"`
	Age  int    `bson:"age"`
	City string `bson:"city"`
}

func InsertData(err error, collection *mongo.Collection) {
	var p = Person{
		Name: "Hello",
		Age:  10,
		City: "WuHan",
	}

	iRes, err := collection.InsertOne(context.Background(), &p)

	if err != nil {
		fmt.Println(err)
		return
	}

	id := iRes.InsertedID.(primitive.ObjectID)
	fmt.Println("自增 ID", id.Hex())
}

func SearchData(collection *mongo.Collection) {
	cond := Person{
		Name: "Hello",
	}

	// 1. 获取到 cursor 对象
	cursor, err := collection.Find(context.Background(), cond)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 2. 延时关闭 cursor
	defer func() {
		if err = cursor.Close(context.Background()); err != nil {
			fmt.Println(err)
		}
	}()

	// 3. 遍历 cursor 得到对象
	for cursor.Next(context.Background()) {
		var p Person
		if cursor.Decode(&p) != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(p)
	}

	// 4. 另外一种遍历方式,但是注意 cursor 的位置
	var results []Person
	if err = cursor.All(context.Background(), &results); err != nil {
		fmt.Println(err)
		return
	}

	for _, result := range results {
		fmt.Println(result)
	}

}

func main() {
	// 设置客户端连接配置
	var (
		client     *mongo.Client // mongoDB 连接对象
		err        error
		db         *mongo.Database   // mongoDB 数据库
		collection *mongo.Collection // mongoDB 数据集合
	)

	// 1. 建立连接
	client, err = mongo.Connect(context.Background(),
		options.Client().ApplyURI("mongodb://localhost:27017").SetConnectTimeout(5*time.Second))

	if err != nil {
		panic(err)
	}

	// 2. 选择数据库
	db = client.Database("runoob")
	// 3. 选择表
	collection = db.Collection("data")

	// CRUD 操作
	// 1. 插入一条数据
	//InsertData(err, collection)

	// 2. 查询数据
	SearchData(collection)
}
