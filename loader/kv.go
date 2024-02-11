package loader

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/redis/go-redis/v9"
)

type KeyValueLoader struct {
	client *redis.Client
}

func NewKeyValueLoader(client *redis.Client) *KeyValueLoader {
	return &KeyValueLoader{
		client: client,
	}
}

func (kv *KeyValueLoader) Load(ctx context.Context, source string) error {
	// read file
	file, err := os.Open(source)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		fqdn := strings.Split(scanner.Text(), ".")
		if err := kv.client.Set(
			ctx,
			strings.Join(fqdn[:len(fqdn)-1], "."),
			fqdn[len(fqdn)-1],
			0,
		).Err(); err != nil {
			fmt.Println(err)
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return err
}
