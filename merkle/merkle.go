package merkle

type MerkleSDK struct {
	ApiKey string

	transactions *TransactionStream
	pool         *PrivatePool
	builder      *BuilderSDK
	simulation   *SimulationAPI
	overwatch    *OverwatchAPI
}

func New() *MerkleSDK {
	return &MerkleSDK{}
}

func (m *MerkleSDK) SetApiKey(apiKey string) {
	m.ApiKey = apiKey
}

// get the api key
func (m *MerkleSDK) GetApiKey() string {
	return m.ApiKey
}

func (m *MerkleSDK) Pool() *PrivatePool {
	if m.pool == nil {
		m.pool = NewPrivatePool(m)
	}
	return m.pool
}

func (m *MerkleSDK) Builder() *BuilderSDK {
	if m.builder == nil {
		m.builder = NewBuilderSDK(m)
	}
	return m.builder
}

func (m *MerkleSDK) Transactions() *TransactionStream {
	if m.transactions == nil {
		m.transactions = NewTransactionStream(m)
	}
	return m.transactions
}

func (m *MerkleSDK) Simulation() *SimulationAPI {
	if m.simulation == nil {
		m.simulation = NewSimulationAPI(m)
	}
	return m.simulation
}

func (m *MerkleSDK) Overwatch() *OverwatchAPI {
	if m.overwatch == nil {
		m.overwatch = NewOverwatchAPI(m)
	}
	return m.overwatch
}
