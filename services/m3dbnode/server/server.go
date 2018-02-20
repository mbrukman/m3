	"github.com/m3db/m3db/environment"
	"github.com/m3db/m3db/x/mmap"
	"github.com/m3db/m3db/x/xio"
	"github.com/m3db/m3x/context"
	"github.com/m3db/m3x/ident"
		logger.Fatalf("could not parse new file mode: %v", err)
		logger.Fatalf("could not parse new directory mode: %v", err)
	mmapCfg := cfg.Filesystem.MmapConfiguration()
	shouldUseHugeTLB := mmapCfg.HugeTLB.Enabled
	if shouldUseHugeTLB {
		// Make sure the host supports HugeTLB before proceeding with it to prevent
		// excessive log spam.
		shouldUseHugeTLB, err = hostSupportsHugeTLB()
		if err != nil {
			logger.Fatalf("could not determine if host supports HugeTLB: %v", err)
		}
		if !shouldUseHugeTLB {
			logger.Warnf("host doesn't support HugeTLB, proceeding without it")
		}
	}
		SetMmapEnableHugeTLB(shouldUseHugeTLB).
		SetMmapHugeTLBThreshold(mmapCfg.HugeTLB.Threshold).
		logger.Fatalf("unknown commit log queue size type: %v",
	// Set the series cache policy
	seriesCachePolicy := cfg.Cache.SeriesConfiguration().Policy
	opts = opts.SetSeriesCachePolicy(seriesCachePolicy)

	// Setup the block retriever
	switch seriesCachePolicy {
	case series.CacheAll:
		// No options needed to be set
	default:
		// All other caching strategies require retrieving series from disk
		// to service a cache miss
		retrieverOpts := fs.NewBlockRetrieverOptions().
			SetBytesPool(opts.BytesPool()).
			SetSegmentReaderPool(opts.SegmentReaderPool()).
			SetIdentifierPool(opts.IdentifierPool())
		if blockRetrieveCfg := cfg.BlockRetrieve; blockRetrieveCfg != nil {
			retrieverOpts = retrieverOpts.
				SetFetchConcurrency(blockRetrieveCfg.FetchConcurrency)
		}
		blockRetrieverMgr := block.NewDatabaseBlockRetrieverManager(
			func(md namespace.Metadata) (block.DatabaseBlockRetriever, error) {
				retriever := fs.NewBlockRetriever(retrieverOpts, fsopts)
				if err := retriever.Open(md); err != nil {
					return nil, err
				}
				return retriever, nil
			})
		opts = opts.SetDatabaseBlockRetrieverManager(blockRetrieverMgr)
	}
	hostID, err := cfg.HostID.Resolve()
	if err != nil {
		logger.Fatalf("could not resolve local host ID: %v", err)
	}

		envCfg environment.ConfigureResults

	case cfg.EnvironmentConfig.Service != nil:
		envCfg, err = cfg.EnvironmentConfig.Configure(environment.ConfigurationParameters{
			InstrumentOpts: iopts,
			HashingSeed:    cfg.Hashing.Seed,
		})
			logger.Fatalf("could not initialize dynamic config: %v", err)
	case cfg.EnvironmentConfig.Static != nil:
		envCfg, err = cfg.EnvironmentConfig.Configure(environment.ConfigurationParameters{
			HostID: hostID,
		})
			logger.Fatalf("could not initialize static config: %v", err)
	opts = opts.SetNamespaceInitializer(envCfg.NamespaceInitializer)
	topo, err := envCfg.TopologyInitializer.Init()
		logger.Fatalf("could not initialize m3db topology: %v", err)
			TopologyInitializer: envCfg.TopologyInitializer,
	bs, err := cfg.Bootstrap.New(opts, m3dbClient)
	kvWatchBootstrappers(envCfg.KVStore, logger, timeout, cfg.Bootstrap.Bootstrappers,
			updated, err := cfg.Bootstrap.New(opts, m3dbClient)
	db, err := cluster.NewDatabase(hostID, envCfg.TopologyInitializer, opts)
		kvWatchNewSeriesLimitPerShard(envCfg.KVStore, logger, topo,
		logger.Infof("bytes pool registering bucket capacity=%d, size=%d, "+
	logger.Infof("bytes pool %s init", policy.Type)
	segmentReaderPool := xio.NewSegmentReaderPool(
	var identifierPool ident.Pool
		identifierPool = ident.NewPool(
		identifierPool = ident.NewNativePool(
func hostSupportsHugeTLB() (bool, error) {
	// Try and determine if the host supports HugeTLB in the first place
	withHugeTLB, err := mmap.Bytes(10, mmap.Options{
		HugeTLB: mmap.HugeTLBOptions{
			Enabled:   true,
			Threshold: 0,
		},
	})
		return false, fmt.Errorf("could not mmap anonymous region: %v", err)
	defer mmap.Munmap(withHugeTLB.Result)
	if withHugeTLB.Warning == nil {
		// If there was no warning, then the host didn't complain about
		// usa of huge TLB
		return true, nil

	// If we got a warning, try mmap'ing without HugeTLB
	withoutHugeTLB, err := mmap.Bytes(10, mmap.Options{})
		return false, fmt.Errorf("could not mmap anonymous region: %v", err)
	defer mmap.Munmap(withoutHugeTLB.Result)
	if withoutHugeTLB.Warning == nil {
		// The machine doesn't support HugeTLB, proceed without it
		return false, nil
	}
	// The warning was probably caused by something else, proceed using HugeTLB
	return true, nil