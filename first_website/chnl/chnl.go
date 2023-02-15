package chnl

import (
	"bytes"
	"context"
)

type Collection struct {
	Url         string
	ContentType string
	Content     *bytes.Buffer
}

type Store interface {
	Insert(ctx context.Context, ch1 Collection)
}

// func (s1 collection) MakeChan() {
// 	ch1 := make(chan collection)
// 	start(s1, ch1)
// }

// func start(s1 collection, ch1 chan collection) {
// 	go func() {
// 		for {
// 			coll, ok := <-ch1
// 			if !ok {
// 				log.Fatalf("An Error Occured ")
// 				break
// 			}

// 			resp, err := http.Post(coll.url, coll.contentType, coll.content)
// 			if err != nil {
// 				log.Fatalf("An Error Occured %v", err)
// 				break
// 			}

// 			if resp != nil {
// 				log.Print("resp done")
// 				break
// 			}
// 			defer resp.Body.Close()
// 		}
// 	}()

// 	ch1 <- s1
// }
