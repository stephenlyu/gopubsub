package redispubsub

import (
	"github.com/garyburd/redigo/redis"
	"time"
	"github.com/Sirupsen/logrus"
)

var (
	DEFAULT_REDIS_SERVER = "localhost:6379"
)

type PubSubMessageCallback interface {
	OnMessage(channel string, data []byte)
}

type RedisPubSub struct {
	redis.Pool
}

func NewRedisPubSub(address string, password string) *RedisPubSub {
	if address == "" {
		address = DEFAULT_REDIS_SERVER
	}

	return &RedisPubSub{
		Pool: redis.Pool {
			MaxIdle:     3,
			IdleTimeout: 240 * time.Second,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", address)
				if err != nil {
					return nil, err
				}
				if password != "" {
					if _, err := c.Do("AUTH", password); err != nil {
						c.Close()
						return nil, err
					}
				}
				return c, err
			},
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
		},
	}
}

func (this *RedisPubSub) Subscribe(callback PubSubMessageCallback, channels ...interface{}) error {
	c := this.Pool.Get()
	defer c.Close()

	psc := redis.PubSubConn{c}
	psc.Subscribe(channels...)

	for {
		switch v := psc.Receive().(type) {
		case redis.Subscription:
			logrus.Infof("%s: %s %d\n", v.Channel, v.Kind, v.Count)
		case redis.Message://单个订阅subscribe
			if callback != nil {
				callback.OnMessage(v.Channel, v.Data)
			}
		case error:
			return v
		}
	}

	return nil
}

func (this *RedisPubSub) Publish(channel string, data []byte) error {
	c := this.Pool.Get()
	defer c.Close()
	_, err := c.Do("PUBLISH", channel, data)
	return err
}
