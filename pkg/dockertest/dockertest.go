package dockertest

import (
	"log"
	"sync"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/ory/dockertest"
)

type component struct {
	r    *dockertest.Resource
	port int
}

type componentIdentifier string

type initComponentFunc func(
	pool *dockertest.Pool,
	modifiers ...OptionModifier,
) (rport int, _ *dockertest.Resource, rerr error)

var components *componentManager

func init() {
	cs, err := newComponentManager()
	if err != nil {
		log.Fatal(err, "failed to init component manager")
	}

	components = cs
}

type componentManager struct {
	m          sync.RWMutex
	initFuncs  map[componentIdentifier]initComponentFunc
	components map[componentIdentifier]*component
	pool       *dockertest.Pool
}

func newComponentManager() (*componentManager, error) {
	const endpoint = ""
	pool, err := dockertest.NewPool(endpoint)

	if err != nil {
		return nil, errors.Wrap(err, "failed to init dockertest pool")
	}

	pool.MaxWait = 1 * time.Minute

	return &componentManager{
		initFuncs:  make(map[componentIdentifier]initComponentFunc),
		components: make(map[componentIdentifier]*component),
		pool:       pool,
	}, nil
}

func (cm *componentManager) registerComponent(ci componentIdentifier, f initComponentFunc) {
	cm.m.Lock()
	defer cm.m.Unlock()

	if _, ok := cm.initFuncs[ci]; ok {
		return
	}

	cm.initFuncs[ci] = f
}

func (cm *componentManager) ensureComponent(ci componentIdentifier, modifiers ...OptionModifier) (rport int, _ error) {
	cm.m.RLock()
	c, ok := cm.components[ci]
	cm.m.RUnlock()

	if ok {
		return c.port, nil
	}

	cm.m.Lock()
	defer cm.m.Unlock()

	if c, ok := cm.components[ci]; ok {
		return c.port, nil
	}

	initFunc, ok := cm.initFuncs[ci]
	if !ok {
		return 0, errors.Newf("no such component identifier registered: %+v", ci)
	}

	port, r, err := initFunc(cm.pool, modifiers...)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	cm.components[ci] = &component{
		r:    r,
		port: port,
	}

	return port, nil
}
