package metrics

import "github.com/PagerDuty/godspeed"

var Metric Metrics = BlackHole{}

type Metrics interface {
	Count(stat string, count float64, tags []string) error
	Set(stat string, value float64, tags []string) error
}

type BlackHole struct{}

func (b BlackHole) Count(stat string, count float64, tags []string) error {
	return nil
}

func (b BlackHole) Set(stat string, value float64, tags []string) error {
	return nil
}

type GodSpeed struct {
	IP        string
	Port      int
	NameSpace string
}

func (b *GodSpeed) newConn() (*godspeed.Godspeed, error) {
	gs, err := godspeed.New(b.IP, b.Port, false)
	if err != nil {
		return gs, err
	}
	gs.Namespace = b.NameSpace

	return gs, err
}

func (b *GodSpeed) Count(stat string, count float64, tags []string) error {
	c, err := b.newConn()
	if err != nil {
		return err
	}
	defer c.Conn.Close()
	return c.Count(stat, count, tags)
}

func (b *GodSpeed) Set(stat string, value float64, tags []string) error {
	c, err := b.newConn()
	if err != nil {
		return err
	}
	defer c.Conn.Close()
	return c.Set(stat, value, tags)
}
