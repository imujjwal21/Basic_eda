package chnl

import (
	"context"
	"log"
	"net/http"
)

type inMem struct {
	ch chan Collection
}

func NewInMem(ch chan Collection) Store {
	return &inMem{ch}
}

func (i *inMem) Insert(ctx context.Context, s1 Collection) {
	go func() {
		for {
			coll, ok := <-i.ch
			if !ok {
				log.Fatalf("An Error Occured ")
				break
			}

			resp, err := http.Post(coll.Url, coll.ContentType, coll.Content)
			if err != nil {
				log.Fatalf("An Error Occured %v", err)
				break
			}

			if resp != nil {
				log.Print("resp done")
				break
			}
			defer resp.Body.Close()
		}
	}()

	i.ch <- s1
}
