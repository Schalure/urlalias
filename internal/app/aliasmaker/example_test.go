package aliasmaker_test

import (
	"context"
	"fmt"

	"github.com/Schalure/urlalias/internal/app/aliaslogger/zaplogger"
	"github.com/Schalure/urlalias/internal/app/aliasmaker"
	"github.com/Schalure/urlalias/internal/app/storage/memstor"
)

func Example() {

	stor, err := memstor.NewStorage()
	if err != nil {
		panic("Can't create storage")
	}

	logger, err := zaplogger.NewZapLogger("")
	if err != nil {
		panic("Can't create logger")
	}

	service, err := aliasmaker.New(stor, logger)
	if err != nil {
		panic("Can't create service")
	}

	userID := uint64(1)

	//	Get short key by original URL
	shortKey, err := service.GetShortKey(context.Background(), userID, "https://example.com")
	if err != nil {
		panic("Can't create shortKey")
	}
	fmt.Println(shortKey)

	//	Get original URL by short key
	originalURL, err := service.GetOriginalURL(context.Background(), shortKey)
	if err != nil {
		panic("Can't return originalURL")
	}
	fmt.Println(originalURL)

	// Output:
	// 000000000
	// https://example.com
}
