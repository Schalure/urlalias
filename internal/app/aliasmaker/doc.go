/*
Package aliasmaker povides the implementation of a link shortening service

	//	Create new service
	service, err := aliasmaker.New(stor, logger)
	...

	//	Create new short URL
	shortKey, err := service.GetShortKey(r.Context(), userID, http://example.com)
	...

	//	Get original URL
	originalURL, err := service.GetOriginalURL(r.Context(), shortKey)
	...
*/
package aliasmaker