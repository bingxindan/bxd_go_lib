package discovery

import (
	"context"
	"encoding/json"
	"github.com/bingxindan/bxd_go_lib/logger"
	"github.com/go-errors/errors"
	client3 "go.etcd.io/etcd/client/v3"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Register for grpc server
type Register struct {
	EtcdAddrs   []string
	DialTimeout int

	closeCh     chan struct{}
	kv          client3.KV
	leasesID    client3.LeaseID
	keepAliveCh <-chan *client3.LeaseKeepAliveResponse

	srvInfo  Server
	srvTTL   int64
	maxRetry int
	cli      *client3.Client
}

// NewRegister create a register base on etcd
func NewRegister(etcdAddrs []string) *Register {
	return &Register{
		EtcdAddrs:   etcdAddrs,
		DialTimeout: 3,
	}
}

// Register a service
func (r *Register) Register(ctx context.Context, srvInfo Server, ttl int64) (chan<- struct{}, error) {
	var err error

	if strings.Split(srvInfo.Addr, ":")[0] == "" {
		return nil, errors.New("invalid ip")
	}

	if r.cli, err = client3.New(client3.Config{
		Endpoints:   r.EtcdAddrs,
		DialTimeout: time.Duration(r.DialTimeout) * time.Second,
	}); err != nil {
		return nil, err
	}

	r.srvInfo = srvInfo
	r.srvTTL = ttl

	if err = r.register(ctx); err != nil {
		return nil, err
	}

	r.closeCh = make(chan struct{})

	go r.keepAlive(ctx)

	return r.closeCh, nil
}

// Stop stop register
func (r *Register) Stop() {
	r.closeCh <- struct{}{}
}

// register 注册节点
func (r *Register) register(ctx context.Context) error {
	if ctx == nil {
		ctx = context.Background()
	}

	leaseCtx, cancel := context.WithTimeout(ctx, time.Duration(r.DialTimeout)*time.Second)
	defer cancel()

	leaseResp, err := r.cli.Grant(leaseCtx, r.srvTTL)
	if err != nil {
		return err
	}

	r.leasesID = leaseResp.ID
	if r.keepAliveCh, err = r.cli.KeepAlive(ctx, leaseResp.ID); err != nil {
		return err
	}

	data, err := json.Marshal(r.srvInfo)
	if err != nil {
		return err
	}
	_, err = r.cli.Put(ctx, BuildRegPath(r.srvInfo), string(data), client3.WithLease(r.leasesID))
	return err
}

// unregister 删除节点
func (r *Register) unregister(ctx context.Context) error {
	_, err := r.cli.Delete(ctx, BuildRegPath(r.srvInfo))
	return err
}

// keepAlive
func (r *Register) keepAlive(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(r.srvTTL) * time.Second)
	for {
		select {
		case <-r.closeCh:
			if err := r.unregister(ctx); err != nil {
				logger.E("etcd_register_keepAlive", "unregister failed err: %+v", err)
			}
			if _, err := r.cli.Revoke(ctx, r.leasesID); err != nil {
				logger.E("etcd_register_keepAlive1", "revoke failed: %+v", err)
			}
			return
		case res := <-r.keepAliveCh:
			if res == nil {
				if err := r.register(ctx); err != nil {
					logger.E("etcd_register_keepAlive2", "register failed: %+v", err)
				}
			}
		case <-ticker.C:
			if r.keepAliveCh == nil {
				if err := r.register(ctx); err != nil {
					logger.E("etcd_register_keepAlive3", "register failed, err: %+v", err)
				}
			}
		}
	}
}

// UpdateHandler return http handler
func (r *Register) UpdateHandler(ctx context.Context) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		wi := req.URL.Query().Get("weight")
		weight, err := strconv.Atoi(wi)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		var update = func() error {
			r.srvInfo.Weight = int64(weight)
			data, err := json.Marshal(r.srvInfo)
			if err != nil {
				return err
			}
			_, err = r.cli.Put(ctx, BuildRegPath(r.srvInfo), string(data), client3.WithLease(r.leasesID))
			return err
		}

		if err := update(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		w.Write([]byte("update server weight success"))
	})
}

func (r *Register) GetServerInfo(ctx context.Context) (Server, error) {
	resp, err := r.cli.Get(ctx, BuildRegPath(r.srvInfo))
	if err != nil {
		return r.srvInfo, err
	}
	info := Server{}
	if resp.Count >= 1 {
		if err := json.Unmarshal(resp.Kvs[0].Value, &info); err != nil {
			return info, err
		}
	}
	return info, nil
}
