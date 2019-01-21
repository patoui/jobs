package jobs

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/spiral/roadrunner/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func viperConfig(cfg string) service.Config {
	v := viper.New()
	v.SetConfigType("json")

	err := v.ReadConfig(bytes.NewBuffer([]byte(cfg)))
	if err != nil {
		panic(err)
	}

	return &configWrapper{v}
}

// configWrapper provides interface bridge between v configs and service.Config.
type configWrapper struct {
	v *viper.Viper
}

// Get nested config section (sub-map), returns nil if section not found.
func (w *configWrapper) Get(key string) service.Config {
	sub := w.v.Sub(key)
	if sub == nil {
		return nil
	}

	return &configWrapper{sub}
}

// Unmarshal unmarshal config data into given struct.
func (w *configWrapper) Unmarshal(out interface{}) error {
	return w.v.Unmarshal(out)
}

func jobs(container service.Container) *Service {
	svc, _ := container.Get("jobs")
	return svc.(*Service)
}

func TestService_Init(t *testing.T) {
	c := service.NewContainer(logrus.New())
	c.Register("jobs", &Service{Brokers: map[string]Broker{"ephemeral": &testBroker{}}})

	assert.NoError(t, c.Init(viperConfig(`{
	"jobs":{
		"workers":{
			"command": "php tests/consumer.php",
			"pool.numWorkers": 1
		},
		"pipelines":{"default":{"broker":"ephemeral"}},
    	"dispatch": {
	    	"spiral-jobs-tests-local-*.pipeline": "default"
    	},
    	"consume": ["default"]
	}
}`)))
}

func TestService_ServeStop(t *testing.T) {
	c := service.NewContainer(logrus.New())
	c.Register("jobs", &Service{Brokers: map[string]Broker{"ephemeral": &testBroker{}}})

	assert.NoError(t, c.Init(viperConfig(`{
	"jobs":{
		"workers":{
			"command": "php tests/consumer.php",
			"pool.numWorkers": 1
		},
		"pipelines":{"default":{"broker":"ephemeral"}},
    	"dispatch": {
	    	"spiral-jobs-tests-local-*.pipeline": "default"
    	},
    	"consume": ["default"]
	}
}`)))

	ready := make(chan interface{})
	jobs(c).AddListener(func(event int, ctx interface{}) {
		if event == EventBrokerReady {
			close(ready)
		}
	})

	go func() { c.Serve() }()
	<-ready
	c.Stop()
}

func TestService_GetPipeline(t *testing.T) {
	c := service.NewContainer(logrus.New())
	c.Register("jobs", &Service{Brokers: map[string]Broker{"ephemeral": &testBroker{}}})

	assert.NoError(t, c.Init(viperConfig(`{
	"jobs":{
		"workers":{
			"command": "php tests/consumer.php",
			"pool.numWorkers": 1
		},
		"pipelines":{"default":{"broker":"ephemeral"}},
    	"dispatch": {
	    	"spiral-jobs-tests-local-*.pipeline": "default"
    	},
    	"consume": ["default"]
	}
}`)))

	ready := make(chan interface{})
	jobs(c).AddListener(func(event int, ctx interface{}) {
		if event == EventBrokerReady {
			close(ready)
		}
	})

	go func() { c.Serve() }()
	defer c.Stop()
	<-ready

	assert.Equal(t, "ephemeral", jobs(c).Pipelines().Get("default").Broker())
}

func TestService_StatPipeline(t *testing.T) {
	c := service.NewContainer(logrus.New())
	c.Register("jobs", &Service{Brokers: map[string]Broker{"ephemeral": &testBroker{}}})

	assert.NoError(t, c.Init(viperConfig(`{
	"jobs":{
		"workers":{
			"command": "php tests/consumer.php",
			"pool.numWorkers": 1
		},
		"pipelines":{"default":{"broker":"ephemeral"}},
    	"dispatch": {
	    	"spiral-jobs-tests-local-*.pipeline": "default"
    	},
    	"consume": ["default"]
	}
}`)))

	ready := make(chan interface{})
	jobs(c).AddListener(func(event int, ctx interface{}) {
		if event == EventBrokerReady {
			close(ready)
		}
	})

	go func() { c.Serve() }()
	defer c.Stop()
	<-ready

	svc := jobs(c)
	pipe := svc.Pipelines().Get("default")

	stat, err := svc.Stat(pipe)
	assert.NoError(t, err)

	assert.Equal(t, int64(0), stat.Queue)
	assert.Equal(t, true, stat.Consuming)
}

func TestService_StatNonConsumingPipeline(t *testing.T) {
	c := service.NewContainer(logrus.New())
	c.Register("jobs", &Service{Brokers: map[string]Broker{"ephemeral": &testBroker{}}})

	assert.NoError(t, c.Init(viperConfig(`{
	"jobs":{
		"workers":{
			"command": "php tests/consumer.php",
			"pool.numWorkers": 1
		},
		"pipelines":{"default":{"broker":"ephemeral"}},
    	"dispatch": {
	    	"spiral-jobs-tests-local-*.pipeline": "default"
    	},
    	"consume": []
	}
}`)))

	ready := make(chan interface{})
	jobs(c).AddListener(func(event int, ctx interface{}) {
		if event == EventBrokerReady {
			close(ready)
		}
	})

	go func() { c.Serve() }()
	defer c.Stop()
	<-ready

	svc := jobs(c)
	pipe := svc.Pipelines().Get("default")

	stat, err := svc.Stat(pipe)
	assert.NoError(t, err)

	assert.Equal(t, int64(0), stat.Queue)
	assert.Equal(t, false, stat.Consuming)
}
