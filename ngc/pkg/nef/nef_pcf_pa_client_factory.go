package ngcnef

func initializePcfClient(cfg Config) PcfPolicyAuthorization {

	return NewPCFClient(&cfg)

}
