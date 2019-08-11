package dockertest

//
//var resources *resourceManager
//
//func init() {
//	r, err := newResourceManager()
//	if err != nil {
//		log.Fatal(err, "failed to init resource manager")
//	}
//	resources = r
//}
//
//type resourceManager struct {
//	services sync.Map
//
//	pm   sync.RWMutex
//	pool *dockertest.Pool
//}
//
//func newResourceManager() (*resourceManager, error) {
//	pool, err := dockertest.NewPool("endpoint")
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to init dockertest pool")
//	}
//
//	pool.MaxWait = 1 * time.Minute
//	return &resourceManager{
//		services: sync.Map{},
//		pool:     pool,
//	}, nil
//}
