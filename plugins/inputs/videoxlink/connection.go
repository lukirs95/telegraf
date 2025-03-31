package videoxlink

import (
	"context"
	"time"

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
	go func() {
		defer xlink.wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case update := <-updates:
				xlink.handleUpdate(update)
			case stat := <-stats:
				xlink.handleStats(stat)
			}
		}
	}()

	for _, system := range xlink.Systems {
		xlink.wg.Add(1)
		x := client.NewClient(system)
		go xlink.connect(ctx, x, system, updates, stats)
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

func (xlink *VideoXLink) connect(ctx context.Context, x *client.Client, ip string, updates client.UpdateChan, stats client.StatsChan) {
	defer xlink.wg.Done()
	for {
		xlink.Log.Infof("start xlink connection (%s)", ip)
		if err := x.Connect(ctx, updates, stats); err != nil {
			xlink.Log.Errorf("xlink connection fault (%s) %w", ip, err)
			time.Sleep(time.Second * 5)
		} else {
			xlink.Log.Infof("xlink connection closed (%s)", ip)
			return
		}
	}
}
