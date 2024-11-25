package videoxlink

import (
	"context"

	"github.com/influxdata/telegraf"
	client "github.com/lukirs95/goxlinkclient"
)

func (xlink *VideoXLink) Start(_ telegraf.Accumulator) error {
	xlink.wgMu.Lock()
	defer xlink.wgMu.Unlock()

	updates := make(client.UpdateChan, 10)
	stats := make(client.StatsChan, 10)
	ctx, stop := context.WithCancel(context.Background())
	xlink.stop = stop

	xlink.wg.Add(1)
	go func(){
		defer xlink.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				xlink.handleUpdate(update)
			case stat := <- stats:
				xlink.handleStats(stat)
			}
		}
	}()

	for _, system := range xlink.Systems {
		xlink.wg.Add(1)
		go func(system string) {
			defer xlink.wg.Done()
			x := client.NewClient(system)
			x.Connect(ctx, updates, stats)
		}(system)
	}
	return nil
}

func (xlink *VideoXLink) Stop() {
	xlink.wgMu.Lock()
	defer xlink.wgMu.Unlock()

	if xlink.stop != nil {
		xlink.stop()
	}

	xlink.wg.Wait()
}

