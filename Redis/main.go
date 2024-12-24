package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sync"
	"time"
)

var rdb *redis.Client

func ConnectRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "251210", // 密码
		DB:       0,        // 数据库
		PoolSize: 20,       // 连接池大小
	})
}

func doCommand() {

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	//执行命令获取结果
	val, err := rdb.Get(ctx, "key").Result()
	fmt.Println(val, err)

	//先获取到命令对象
	cmd := rdb.Get(ctx, "key")
	fmt.Println(cmd.Val())
	fmt.Println(cmd.Err())

	//直接执行命令获取错误  time.Hour设置的键值对持续时间也就是一个小时后过期
	err = rdb.Set(ctx, "key", 10, time.Hour).Err()

	// 直接执行命令获取值
	value := rdb.Get(ctx, "key").Val()
	fmt.Println(value)

}

// 执行任意命令或自定义命令
func doDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	//直接执行命令获取错误
	err := rdb.Do(ctx, "set", "key", 10, "EX", 3600).Err()
	fmt.Println(err)

	//执行命令获取结果
	val, err := rdb.Do(ctx, "get", "key").Result()
	fmt.Println(val, err)
}

// 操作zset示例
func zsetDemo() {
	//key
	zsetKey := "language_rank"
	//value
	//注意：v8版本使用[]*redis.Z,v9版本使用[]redis.Z
	language := []redis.Z{
		{Score: 90.0, Member: "Goland"},
		{Score: 98.0, Member: "Java"},
		{Score: 95.0, Member: "Python"},
		{Score: 97.0, Member: "JavaScript"},
		{Score: 99.0, Member: "C/C++"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	//ZADD
	err := rdb.ZAdd(ctx, zsetKey, language...).Err()
	if err != nil {
		fmt.Printf("zadd data failed! err: %v\n", err)
		return
	}
	fmt.Println("zadd data success!")

	//把Goland分数加10
	newScore, err := rdb.ZIncrBy(ctx, zsetKey, 10, "Goland").Result()
	if err != nil {
		fmt.Printf("zincrby failed! err: %v\n", err)
		return
	}
	fmt.Printf("Goland's new score is: %f.", newScore)

	//取分数最高的三个
	ret := rdb.ZRevRangeWithScores(ctx, zsetKey, 0, 2).Val()
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}
	//取95~100分的
	op := &redis.ZRangeBy{Min: "95", Max: "100"}
	ret, err = rdb.ZRangeByScoreWithScores(ctx, zsetKey, op).Result()
	if err != nil {
		fmt.Printf("zrangebyscore failed, err:%v\n", err)
		return
	}
	for _, z := range ret {
		fmt.Println(z.Member, z.Score)
	}
}

// PipelineDemo 管道示例
// 使用pipeline一次执行100个Get命令，遍历取出100个命令的执行结果
func PipelineDemo() {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	cmds, err := rdb.Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for i := 0; i < 100; i++ {
			pipe.Get(ctx, fmt.Sprintf("key%d", i))
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, cmd := range cmds {
		fmt.Println(cmd.(*redis.StringCmd).Val())
	}
}

// watchDemo  watch:对key进行监视，如果在监视期间其他用户对key进行修改删除等操作，执行exec会返回一个错误
func watchDemo(key string) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	err := rdb.Watch(ctx, func(tx *redis.Tx) error {
		n, err := tx.Get(ctx, key).Int()
		if err != nil && err != redis.Nil {
			return err
		}
		// 假设操作耗时5秒
		// 5秒内我们通过其他的客户端修改key，当前事务就会失败
		time.Sleep(5 * time.Second)
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, key, n+1, time.Hour)
			return nil
		})
		return err
	}, key)
	if err != nil {
		fmt.Printf("报错：%v\n", err)
	}
}

// INCRDemo 使用 GET 、SET和WATCH命令实现一个 INCR 命令的完整示例
func INCRDemo() {
	const routineCount = 100

	//设置五秒的超时
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//increment是一个自定义对key进行递增（+1）的函数
	//使用GET,SET,WATCH实现
	increment := func(key string) error {
		//定义一个函数，这个函数是watch需要传递的参数  事务函数
		txf := func(tx *redis.Tx) error {
			//获取当前值或零值
			n, err := tx.Get(ctx, key).Int()
			if err != nil && err != redis.Nil {
				fmt.Printf("获取key值操作错误，err%v\n", err)
				return err
			}
			//对k值进行加一
			n++

			_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
				pipe.Set(ctx, key, n, 0)
				return nil
			})
			return err
		}
		//最多重试100次
		for retries := routineCount; retries > 0; retries-- {
			err := rdb.Watch(ctx, txf, key)
			if err != redis.TxFailedErr {
				fmt.Println("watch处的err")
				return err
			}
			// 乐观锁丢失
		}
		return errors.New("increment reached maximum number of retries")
	}
	// 开启100个goroutine并发调用increment
	// 相当于对key执行100次递增
	var wg sync.WaitGroup
	wg.Add(routineCount)
	for i := 0; i < routineCount; i++ {
		go func() {
			defer wg.Done()

			if err := increment("lsl"); err != nil {
				fmt.Println("increment error:", err)
			}
		}()
	}
	wg.Wait()

	n, err := rdb.Get(ctx, "lsl").Int()
	fmt.Println("最终结果：", n, err)
}

func main() {
	ConnectRedis()
	//doCommand()
	//zset操作示例
	//zsetDemo()

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	//rdb.Set(ctx, "lsl", 18, time.Hour)
	//watchDemo("lsl")
	fmt.Printf("值：%v\n", rdb.Get(ctx, "lsl").Val())

	//INCRDemo()
}
