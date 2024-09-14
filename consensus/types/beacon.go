package types
type BeaconBlock struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    Bytes32
	StateRoot     Bytes32
	Body          BeaconBlockBody
}
type BeaconBlockBody struct {
	RandaoReveal          SignatureBytes
	Eth1Data              Eth1Data
	Graffiti              Bytes32
	ProposerSlashings     [16]ProposerSlashing
	AttesterSlashings     [2]AttesterSlashing
	Attestations          [128]Attestation
	Deposits              [16]Deposit
	VoluntaryExits        [16]SignedVoluntaryExit
	SyncAggregate         SyncAggregate
	ExecutionPayload      ExecutionPayload
	BlsToExecutionChanges [16]SignedBlsToExecutionChange // For Capella and Deneb
	BlobKzgCommitments    [4096]ByteVector               // For Deneb
}
func (b BeaconBlockBody) Default() BeaconBlockBody {
	return BeaconBlockBody{}
}
type SignedBlsToExecutionChange struct {
	Message   BlsToExecutionChange
	Signature SignatureBytes
}
type BlsToExecutionChange struct {
	ValidatorIndex     uint64
	FromBlsPubkey      BLSPubKey
	ToExecutionAddress Address
}
type ExecutionPayload struct {
	ParentHash    Bytes32
	FeeRecipient  Address
	StateRoot     Bytes32
	ReceiptsRoot  Bytes32
	LogsBloom     LogsBloom
	PrevRandao    Bytes32
	BlockNumber   uint64
	GasLimit      uint64
	GasUsed       uint64
	Timestamp     uint64
	ExtraData     ByteList
	BaseFeePerGas uint64 // Assuming U256 is a uint64
	BlockHash     Bytes32
	Transactions  []ByteList   // List of transactions with a max length
	Withdrawals   []Withdrawal // For Capella and Deneb
	BlobGasUsed   uint64       // For Deneb
	ExcessBlobGas uint64       // For Deneb
}
func (e ExecutionPayload) Default() ExecutionPayload {
	// Implement default initialization if needed
	return ExecutionPayload{}
}
type Withdrawal struct {
	Index          uint64
	ValidatorIndex uint64
	Address        Address
	Amount         uint64
}
type ProposerSlashing struct {
	SignedHeader1 SignedBeaconBlockHeader
	SignedHeader2 SignedBeaconBlockHeader
}
type SignedBeaconBlockHeader struct {
	Message   BeaconBlockHeader
	Signature SignatureBytes
}
type BeaconBlockHeader struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    Bytes32
	StateRoot     Bytes32
	BodyRoot      Bytes32
}
type AttesterSlashing struct {
	Attestation1 IndexedAttestation
	Attestation2 IndexedAttestation
}
type IndexedAttestation struct {
	AttestingIndices [2048]uint64
	Data             AttestationData
	Signature        SignatureBytes
}
type Attestation struct {
	AggregationBits [2048]byte
	Data            AttestationData
	Signature       SignatureBytes
}
type AttestationData struct {
	Slot            uint64
	Index           uint64
	BeaconBlockRoot Bytes32
	Source          Checkpoint
	Target          Checkpoint
}
type Checkpoint struct {
	Epoch uint64
	Root  Bytes32
}
type SignedVoluntaryExit struct {
	Message   VoluntaryExit
	Signature SignatureBytes
}
type VoluntaryExit struct {
	Epoch          uint64
	ValidatorIndex uint64
}
type Deposit struct {
	Proof [33]Bytes32
	Data  DepositData
}
type DepositData struct {
	Pubkey                BLSPubKey
	WithdrawalCredentials Bytes32
	Amount                uint64
	Signature             SignatureBytes
}
type Eth1Data struct {
	DepositRoot  Bytes32
	DepositCount uint64
	BlockHash    Bytes32
}
type Bootstrap struct {
	Header                     Header
	CurrentSyncCommittee       SyncCommittee
	CurrentSyncCommitteeBranch []Bytes32
}
type Update struct {
	AttestedHeader          Header
	NextSyncCommittee       SyncCommittee
	NextSyncCommitteeBranch []Bytes32
	FinalizedHeader         Header
	FinalityBranch          []Bytes32
	SyncAggregate           SyncAggregate
	SignatureSlot           uint64
}
type FinalityUpdate struct {
	AttestedHeader  Header
	FinalizedHeader Header
	FinalityBranch  []Bytes32
	SyncAggregate   SyncAggregate
	SignatureSlot   uint64
}
type OptimisticUpdate struct {
	AttestedHeader Header
	SyncAggregate  SyncAggregate
	SignatureSlot  uint64
}
type Header struct {
	Slot          uint64
	ProposerIndex uint64
	ParentRoot    Bytes32
	StateRoot     Bytes32
	BodyRoot      Bytes32
}
type SyncCommittee struct {
	Pubkeys         [512]BLSPubKey
	AggregatePubkey BLSPubKey
}
type SyncAggregate struct {
	SyncCommitteeBits      []byte 
	SyncCommitteeSignature SignatureBytes
}
type GenericUpdate struct {
	AttestedHeader          Header
	SyncAggregate           SyncAggregate
	SignatureSlot           uint64
	NextSyncCommittee       *SyncCommittee
	NextSyncCommitteeBranch *[]Bytes32
	FinalizedHeader         *Header
	FinalityBranch          *[]Bytes32
}
func (g *GenericUpdate) FromUpdate(update *Update) {
	g.AttestedHeader = update.AttestedHeader
	g.SyncAggregate = update.SyncAggregate
	g.SignatureSlot = update.SignatureSlot
	g.NextSyncCommittee = &update.NextSyncCommittee
	g.NextSyncCommitteeBranch = &update.NextSyncCommitteeBranch
	g.FinalizedHeader = &update.FinalizedHeader
	g.FinalityBranch = &update.FinalityBranch
}
func (g *GenericUpdate) FromFinalityUpdate(update *FinalityUpdate) {
	g.AttestedHeader = update.AttestedHeader
	g.SyncAggregate = update.SyncAggregate
	g.SignatureSlot = update.SignatureSlot
	g.NextSyncCommittee = nil
	g.NextSyncCommitteeBranch = nil
	g.FinalizedHeader = &update.FinalizedHeader
	g.FinalityBranch = &update.FinalityBranch
}
func (g *GenericUpdate) FromOptimisticUpdate(update *OptimisticUpdate) {
	g.AttestedHeader = update.AttestedHeader
	g.SyncAggregate = update.SyncAggregate
	g.SignatureSlot = update.SignatureSlot
	g.NextSyncCommittee = nil
	g.NextSyncCommitteeBranch = nil
	g.FinalizedHeader = nil
	g.FinalityBranch = nil
}
