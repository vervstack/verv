package deploy

import (
	"context"
	"sync"
	"time"

	"go.vervstack.ru/verv/plugins/deployment/nodeclient"
	"go.vervstack.ru/verv/plugins/deployment/nodes"
	"go.vervstack.ru/verv/plugins/deployment/pipeline"
)

const probeTimeout = 3 * time.Second

func probeNodes(ctx context.Context, entries []nodes.Entry) []pipeline.Status {
	statuses := make([]pipeline.Status, len(entries))

	var wg sync.WaitGroup

	for idx := range entries {
		wg.Add(1)

		go probeEntry(ctx, &wg, entries[idx], &statuses[idx])
	}

	wg.Wait()

	return statuses
}

func probeEntry(ctx context.Context, wg *sync.WaitGroup, entry nodes.Entry, status *pipeline.Status) {
	defer wg.Done()

	probeCtx, cancel := context.WithTimeout(ctx, probeTimeout)
	defer cancel()

	client := nodeclient.New(entry.Url, entry.Secret)

	*status = pipeline.Probe(probeCtx, client)
}
