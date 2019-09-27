package types

import (
	"encoding/binary"
	"github.com/huandu/xstrings"
	"reflect"
	"strconv"
	"subscan-end/utiles"
	"subscan-end/utiles/uint128"
	"unicode/utf8"
)

type Compact struct {
	ScaleType
	CompactLength int    `json:"compact_length"`
	CompactBytes  []byte `json:"compact_bytes"`
}

func (c *Compact) ProcessCompactBytes() []byte {
	compactByte := c.GetNextBytes(1)
	byteMod := compactByte[0] % 4
	if byteMod == 0 {
		c.CompactLength = 1
	} else if byteMod == 1 {
		c.CompactLength = 2
	} else if byteMod == 2 {
		c.CompactLength = 4
	} else {
		c.CompactLength = 5 + ((int(compactByte[0]) - 3) / 4)
	}
	if c.CompactLength == 1 {
		c.CompactBytes = compactByte
	} else if utiles.IntInSlice(c.CompactLength, []int{2, 4}) {
		c.CompactBytes = append(compactByte[:], c.GetNextBytes(c.CompactLength - 1)[:]...)
	} else {
		c.CompactBytes = c.GetNextBytes(c.CompactLength - 1)
	}
	return c.CompactBytes
}

func (c *Compact) Process() {
	c.ProcessCompactBytes()
	if c.SubType != "" {
		s := ScaleDecoder{TypeString: c.SubType, Data: ScaleBytes{Data: c.CompactBytes}}
		byteData := s.ProcessAndUpdateData(c.SubType)
		if reflect.TypeOf(byteData).Kind() == reflect.Int && c.CompactLength <= 4 {
			c.Value = []byte(strconv.Itoa(byteData.(int) / 4))
		} else {
			c.Value = byteData
		}
	} else {
		c.Value = c.CompactBytes
	}
}

type CompactU32 struct {
	Compact
}

func (c *CompactU32) Init(data ScaleBytes, subType string, arg ...interface{}) {
	c.TypeString = "Compact<u32>"
	c.ScaleDecoder.Init(data, "", arg...)
}

func (c *CompactU32) Process() {
	c.ProcessCompactBytes()
	if c.CompactLength <= 4 {
		data := make([]byte, len(c.Data.Data))
		copy(data, c.Data.Data)

		compactBytes := c.CompactBytes
		bs := make([]byte, 4-len(compactBytes))
		compactBytes = append(compactBytes[:], bs...)
		c.Data.Data = data
		c.Value = int(binary.LittleEndian.Uint32(compactBytes)) / 4
	} else {
		c.Value = int(binary.LittleEndian.Uint32(c.CompactBytes))
	}
}

func (c *CompactU32) Encode(value int) ScaleBytes {
	if value <= 63 {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(value<<2))
		c.Data.Data = bs[0:1]
	} else if value <= 16383 {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(value<<2)|1)
		c.Data.Data = bs[0:2]
	} else if value <= 1073741823 {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(value<<2)|2)
		c.Data.Data = bs
	}
	return c.Data
}

type BoxProposal struct {
	ScaleType
	CallIndex          string                   `json:"-"`
	Call               MetadataCalls            `json:"-"`
	CallModule         MetadataModules          `json:"-"`
	Params             []map[string]interface{} `json:"-"`
	MetadataCallsIndex map[string]interface{}   `json:"-"`
}

func (b *BoxProposal) Init(data ScaleBytes, subType string, arg ...interface{}) {
	b.TypeString = "Box<Proposal>"
	if len(arg) > 0 {
		b.MetadataCallsIndex = arg[0].(map[string]interface{})
	}
	b.ScaleDecoder.Init(data, "", arg...)
}

func (b *BoxProposal) Process() {
	b.CallIndex = utiles.BytesToHex(b.GetNextBytes(2))
	call := b.MetadataCallsIndex[b.CallIndex].(map[string]interface{})
	b.Call = call["module"].(MetadataCalls)
	b.CallModule = call["CallModule"].(MetadataModules)
	for _, arg := range b.Call.Args {
		argObj := b.ProcessAndUpdateData(arg["type"].(string))
		b.Params = append(b.Params, map[string]interface{}{
			"name":  arg["name"].(string),
			"type":  arg["type"].(string),
			"value": argObj,
		})
	}
	b.Value = map[string]interface{}{
		"call_index":  b.CallIndex,
		"call_name":   b.Call.Name,
		"call_module": b.CallModule.Name,
		"params":      b.Params,
	}
}

type Option struct {
	ScaleType
}

func (o *Option) Process() {
	optionType := o.GetNextBytes(1)
	if o.SubType != "" && utiles.BytesToHex(optionType) != "00" {
		o.Value = o.ProcessAndUpdateData(o.SubType)
	}
}

type Bytes struct {
	ScaleType
}

func (b *Bytes) Init(data ScaleBytes, subType string, arg ...interface{}) {
	b.TypeString = "Vec<u8>"
	b.ScaleDecoder.Init(data, "", arg...)
}

func (b *Bytes) Process() {
	length := b.ProcessAndUpdateData("Compact<u32>").(int)
	value := b.GetNextBytes(int(length))
	if utf8.Valid(value) {
		b.Value = string(value)
	} else {
		b.Value = utiles.BytesToHex(value)
	}
}

type OptionBytes struct {
	ScaleType
}

func (b *OptionBytes) Init(data ScaleBytes, subType string, arg ...interface{}) {
	b.TypeString = "Option<Vec<u8>>"
	b.ScaleDecoder.Init(data, "", arg...)
}

func (b *OptionBytes) Process() {
	optionByte := b.GetNextBytes(1)
	if utiles.BytesToHex(optionByte) != "00" {
		b.Value = b.ProcessAndUpdateData("Bytes").(string)
	}
}

type HexBytes struct {
	ScaleType
}

func (h *HexBytes) Process() {
	length := h.ProcessAndUpdateData("Compact<u32>").(int)
	h.Value = utiles.AddHex(utiles.BytesToHex(h.GetNextBytes(int(length))))
}

type Text struct {
	ScaleType
}

func (t *Text) Process() {
	length := t.ProcessAndUpdateData("Compact<u32>").(int)
	t.Value = string(t.GetNextBytes(int(length)))
}

type U8 struct {
	ScaleType
}

func (u *U8) Process() {
	u.Value = u.GetNextU8()
}

type U32 struct {
	ScaleType
}

func (u *U32) Process() {
	u.Value = uint32(binary.LittleEndian.Uint32(u.GetNextBytes(4)))
}

type U64 struct {
	ScaleType
}

func (u *U64) Process() {
	u.Value = uint64(binary.LittleEndian.Uint64(u.GetNextBytes(8)))
}

type U128 struct {
	ScaleType
}

func (u *U128) Process() {
	if len(u.Data.Data) < 16 {
		u.Data.Data = utiles.HexToBytes(xstrings.RightJustify(utiles.BytesToHex(u.Data.Data), 32, "0"))
	}
	u.Value = uint128.FromBytes(u.GetNextBytes(16))
}

type H256 struct {
	ScaleType
}

func (h *H256) Process() {
	h.Value = utiles.AddHex(utiles.BytesToHex(h.GetNextBytes(32)))
}

type Era struct {
	ScaleType
}

func (e *Era) Process() {
	optionHex := utiles.BytesToHex(e.GetNextBytes(1))
	if optionHex == "00" {
		e.Value = optionHex
	} else {
		e.Value = optionHex + utiles.BytesToHex(e.GetNextBytes(1))
	}
}

type Bool struct {
	ScaleType
}

func (b *Bool) Process() {
	b.Value = b.getNextBool()
}

type Moment struct {
	CompactU32
}

func (m *Moment) Init(data ScaleBytes, subType string, arg ...interface{}) {
	m.TypeString = "Compact<Moment>"
	m.ScaleDecoder.Init(data, subType, arg...)
}

func (m *Moment) Process() {
	intValue := m.ProcessAndUpdateData("Compact<u32>").(int)
	m.Value = intValue
}

type Struct struct {
	ScaleType
}

func (s *Struct) Init(data ScaleBytes, subType string, arg ...interface{}) {
	s.ScaleDecoder.Init(data, subType, arg...)
}

func (s *Struct) Process() {
	result := make(map[string]interface{})
	for _, dataType := range s.StructOrderField {
		result[dataType] = s.ProcessAndUpdateData(s.TypeMapping[dataType])
	}
	s.Value = result
}

type ValidatorPrefs struct {
	Struct
}

func (v *ValidatorPrefs) Init(data ScaleBytes, subType string, arg ...interface{}) {
	v.TypeString = "(Compact<u32>, Compact<Balance>)"
	v.Struct.Init(data, subType, arg...)
}

type AccountId struct {
	H256
}

type AccountIndex struct {
	U32
}

type ReferendumIndex struct {
	U32
}

type PropIndex struct {
	U32
}

type Vote struct {
	U8
}

type SessionKey struct {
	H256
}

type AttestedCandidate struct {
	H256
}

type Balance struct {
	U128
}

type ParaId struct {
	U32
}

type Key struct {
	Bytes
}

type KeyValue struct {
	Struct
}

func (k *KeyValue) Init(data ScaleBytes, subType string, arg ...interface{}) {
	k.TypeString = "(Vec<u8>, Vec<u8>)"
	k.TypeMapping = map[string]string{
		"key":   "Vec<u8>",
		"value": "Vec<u8>",
	}
	k.StructOrderField = []string{"key", "value"}
	k.Struct.Init(data, subType, arg...)
}

type Signature struct {
	ScaleType
}

func (s *Signature) Process() {
	s.Value = utiles.BytesToHex(s.GetNextBytes(64))
}

type BalanceOf struct {
	Balance
}

type BlockNumber struct {
	U64
}

type NewAccountOutcome struct {
	CompactU32
}

type Vec struct {
	ScaleType
	Elements []interface{} `json:"elements"`
}

func (v *Vec) Init(data ScaleBytes, subType string, arg ...interface{}) {
	v.Elements = []interface{}{}
	v.ScaleDecoder.Init(data, subType, arg...)
}

func (v *Vec) Process() {
	elementCount := v.ProcessAndUpdateData("Compact<u32>").(int)
	var result []interface{}
	for i := 0; i < elementCount; i++ {
		element := v.ProcessAndUpdateData(v.SubType)
		v.Elements = append(v.Elements, element)
		result = append(result, element)
	}
	v.Value = result
}

type Address struct {
	ScaleType
	AccountLength string `json:"account_length"`
	AccountId     string `json:"account_id"`
	AccountIndex  string `json:"account_index"`
	AccountIdx    string `json:"account_idx"`
}

func (a *Address) Process() {
	AccountLength := a.GetNextBytes(1)
	a.AccountLength = utiles.BytesToHex(AccountLength)
	if a.AccountLength == "ff" {
		a.AccountId = utiles.BytesToHex(a.GetNextBytes(32))
	} else {
		var AccountIndex []byte
		if a.AccountLength == "fc" {
			AccountIndex = a.GetNextBytes(2)
		} else if a.AccountLength == "fd" {
			AccountIndex = a.GetNextBytes(4)
		} else if a.AccountLength == "fe" {
			AccountIndex = a.GetNextBytes(8)
		} else {
			AccountIndex = AccountLength
		}
		a.AccountIndex = utiles.BytesToHex(AccountIndex)
		a.AccountIdx = strconv.FormatUint(uint64(binary.LittleEndian.Uint32(AccountIndex)), 10)
	}
	a.Value = map[string]string{"account_length": a.AccountLength, "account_id": a.AccountId, "account_index": a.AccountIndex, "account_idx": a.AccountIdx}
}

type RawAddress struct {
	Address
}

type Enum struct {
	ScaleType
	ValueList []string `json:"value_list"`
	Index     int      `json:"index"`
}

func (e *Enum) Init(data ScaleBytes, subType string, arg ...interface{}) {
	e.Index = 0
	e.ValueList = arg[0].([]string)
	e.ScaleDecoder.Init(data, subType, arg...)
}

func (e *Enum) Process() {
	index := utiles.BytesToHex(e.GetNextBytes(1))
	e.Index = utiles.StringToInt(index)
	if e.ValueList[e.Index] != "" {
		e.Value = e.ValueList[e.Index]
	} else {
		e.Value = ""
	}
}

type RewardDestination struct {
	Enum
}

func (r *RewardDestination) Init(data ScaleBytes, subType string, arg ...interface{}) {
	r.ValueList = []string{"Staked", "Stash", "Controller"}
	r.ScaleDecoder.Init(data, subType, arg...)
}

type VoteThreshold struct {
	Enum
}

func (v *VoteThreshold) Init(data ScaleBytes, subType string, arg ...interface{}) {
	v.ValueList = []string{"SuperMajorityApprove", "SuperMajorityAgainst", "SimpleMajority"}
	v.ScaleDecoder.Init(data, subType, arg...)
}

type Inherent struct {
	Bytes
}

type LockPeriods struct {
	U8
}

type Hash struct {
	H256
}

type VoteIndex struct {
	U32
}

type ProposalIndex struct {
	U32
}

type Permill struct {
	U32
}

type IdentityType struct {
	Bytes
}

type VoteType struct {
	Enum
}

func (v *VoteType) Init(data ScaleBytes, subType string, arg ...interface{}) {
	v.Index = 0
	v.TypeString = "voting::VoteType"
	v.ValueList = []string{"Binary", "MultiOption"}
	v.ScaleDecoder.Init(data, subType, arg...)
}

type VoteOutcome struct {
	ScaleType
}

func (v *VoteOutcome) Process() []byte {
	return v.GetNextBytes(32)
}

type Identity struct {
	Bytes
}

type ProposalTitle struct {
	Bytes
}

type ProposalContents struct {
	Bytes
}

type ProposalStage struct {
	Enum
}

func (p *ProposalStage) Init(data ScaleBytes, subType string, arg ...interface{}) {
	p.ValueList = []string{"PreVoting", "Voting", "Completed"}
	p.ScaleDecoder.Init(data, subType, arg...)
}

type ProposalCategory struct {
	Enum
}

func (p *ProposalCategory) Init(data ScaleBytes, subType string, arg ...interface{}) {
	p.ValueList = []string{"Signaling"}
	p.ScaleDecoder.Init(data, subType, arg...)
}

type VoteStage struct {
	Enum
}

func (p *VoteStage) Init(data ScaleBytes, subType string, arg ...interface{}) {
	p.ValueList = []string{"PreVoting", "Commit", "Voting", "Completed"}
	p.ScaleDecoder.Init(data, subType, arg...)
}

type TallyType struct {
	Enum
}

func (t *TallyType) Init(data ScaleBytes, subType string, arg ...interface{}) {
	t.TypeString = "voting::TallyType"
	t.ValueList = []string{"OnePerson", "OneCoin"}
	t.ScaleDecoder.Init(data, subType, arg...)
}

type Attestation struct {
	Bytes
}

type ContentId struct {
	H256
}

type MemberId struct {
	U64
}

type PaidTermId struct {
	U64
}

type SubscriptionId struct {
	U64
}

type SchemaId struct {
	U64
}

type DownloadSessionId struct {
	U64
}

type UserInfo struct {
	Struct
}

func (u *UserInfo) Init(data ScaleBytes, subType string, arg ...interface{}) {
	u.TypeMapping = map[string]string{
		"handle":     "Option<Vec<u8>>",
		"avatar_uri": "Option<Vec<u8>>",
		"about":      "Option<Vec<u8>>",
	}
	u.StructOrderField = []string{"handle", "avatar_uri", "about"}
	u.ScaleDecoder.Init(data, subType, arg...)
}

type Role struct {
	Enum
}

func (r *Role) Init(data ScaleBytes, subType string, arg ...interface{}) {
	r.ValueList = []string{"Storage"}
	r.ScaleDecoder.Init(data, subType, arg...)
}

type ContentVisibility struct {
	Enum
}

func (c *ContentVisibility) Init(data ScaleBytes, subType string, arg ...interface{}) {
	c.ValueList = []string{"Draft", "Public"}
	c.ScaleDecoder.Init(data, subType, arg...)
}

type ContentMetadata struct {
	Struct
}

func (c *ContentMetadata) Init(data ScaleBytes, subType string, arg ...interface{}) {
	c.TypeMapping = map[string]string{
		"owner":        "AccountId",
		"added_at":     "BlockAndTime",
		"children_ids": "Vec<ContentId>",
		"visibility":   "ContentVisibility",
		"schema":       "SchemaId",
		"json":         "Vec<u8>",
	}
	c.StructOrderField = []string{"owner", "added_at", "children_ids", "visibility", "schema", "json"}
	c.ScaleDecoder.Init(data, subType, arg...)
}

type ContentMetadataUpdate struct {
	Struct
}

func (c *ContentMetadataUpdate) Init(data ScaleBytes, subType string, arg ...interface{}) {
	c.TypeMapping = map[string]string{
		"children_ids": "Option<Vec<ContentId>>",
		"visibility":   "Option<ContentVisibility>",
		"schema":       "Option<SchemaId>",
		"json":         "Option<Vec<u8>>",
	}
	c.StructOrderField = []string{"children_ids", "visibility", "schema", "json"}
	c.ScaleDecoder.Init(data, subType, arg...)
}

type LiaisonJudgement struct {
	Enum
}

func (l *LiaisonJudgement) Init(data ScaleBytes, subType string, arg ...interface{}) {
	l.ValueList = []string{"Pending", "Accepted", "Rejected"}
	l.ScaleDecoder.Init(data, subType, arg...)
}

type BlockAndTime struct {
	Struct
}

func (b *BlockAndTime) Init(data ScaleBytes, subType string, arg ...interface{}) {
	b.TypeMapping = map[string]string{
		"block": "BlockNumber",
		"time":  "Moment",
	}
	b.StructOrderField = []string{"block", "time"}
	b.ScaleDecoder.Init(data, subType, arg...)
}

type DataObjectTypeId struct {
	U64
}

func (d *DataObjectTypeId) Init(data ScaleBytes, subType string, arg ...interface{}) {
	d.TypeString = "<T as DOTRTrait>::DataObjectTypeId"
	d.ScaleDecoder.Init(data, subType, arg...)
}

type DataObject struct {
	Struct
}

func (d *DataObject) Init(data ScaleBytes, subType string, arg ...interface{}) {
	d.TypeMapping = map[string]string{
		"owner":             "AccountId",
		"added_at":          "BlockAndTime",
		"type_id":           "DataObjectTypeId",
		"size":              "u64",
		"liaison":           "AccountId",
		"liaison_judgement": "LiaisonJudgement",
	}
	d.StructOrderField = []string{"owner", "added_at", "type_id", "size", "liaison", "liaison_judgement"}
	d.ScaleDecoder.Init(data, subType, arg...)
}

type DataObjectStorageRelationshipId struct {
	U64
}

type ProposalStatus struct {
	Enum
}

func (p *ProposalStatus) Init(data ScaleBytes, subType string, arg ...interface{}) {
	p.ValueList = []string{"Active", "Cancelled", "Expired", "Approved", "Rejected", "Slashed"}
	p.ScaleDecoder.Init(data, subType, arg...)
}

type VoteKind struct {
	Enum
}

func (v *VoteKind) Init(data ScaleBytes, subType string, arg ...interface{}) {
	v.ValueList = []string{"Abstain", "Approve", "Reject", "Slash"}
	v.ScaleDecoder.Init(data, subType, arg...)
}

type TallyResult struct {
	Struct
}

func (t *TallyResult) Init(data ScaleBytes, subType string, arg ...interface{}) {
	t.TypeString = "TallyResult<BlockNumber>"
	t.TypeMapping = map[string]string{
		"proposal_id":  "u32",
		"abstentions":  "u32",
		"approvals":    "u32",
		"rejections":   "u32",
		"slashes":      "u32",
		"status":       "ProposalStatus",
		"finalized_at": "BlockNumber",
	}
	t.StructOrderField = []string{"proposal_id", "abstentions", "approvals", "rejections", "slashes", "status", "finalized_at"}
	t.ScaleDecoder.Init(data, subType, arg...)
}

type StorageHasher struct {
	Enum
}

func (s *StorageHasher) Init(data ScaleBytes, subType string, arg ...interface{}) {
	valueList := []string{"Blake2_128", "Blake2_256", "Twox128", "Twox256", "Twox128Concat"}
	s.Enum.Init(data, "", valueList)
}

func (s *StorageHasher) isBlake2_128() bool {
	return s.Index == 0
}
func (s *StorageHasher) isBlake2_256() bool {
	return s.Index == 1
}
func (s *StorageHasher) isTwox128() bool {
	return s.Index == 2
}
func (s *StorageHasher) isTwox256() bool {
	return s.Index == 3
}
func (s *StorageHasher) isTwox128Concat() bool {
	return s.Index == 4
}

type Other struct {
	Bytes
}

type TokenBalance struct {
	U128
}

type Currency struct {
	U128
}

type CurrencyOf struct {
	U128
}

type SessionIndex struct {
	U32
}

type Keys struct {
	Struct
}

func (d *Keys) Init(data ScaleBytes, subType string, arg ...interface{}) {
	d.TypeString = "(AccountId, AccountId)"
	d.ScaleDecoder.Init(data, subType, arg...)
}

type ScheduleGas struct {
	Struct
}

func (s *ScheduleGas) Init(data ScaleBytes, subType string, arg ...interface{}) {
	s.TypeMapping = map[string]string{
		"version":               "u32",
		"putCodePerByteCost":    "gas",
		"growMemCost":           "gas",
		"regularOpCost":         "gas",
		"returnDataPerByteCost": "gas",
		"eventDataPerByteCost":  "gas",
		"eventPerTopicCost":     "gas",
		"eventBaseCost":         "gas",
		"sandboxDataReadCost":   "gas",
		"sandboxDataWriteCost":  "gas",
		"maxEventTopics":        "u32",
		"maxStackHeight":        "u32",
		"maxMemoryPages":        "u32",
		"enablePrintln":         "bool",
		"maxSubjectLen":         "u32",
	}
	s.StructOrderField = []string{"version", "putCodePerByteCost", "growMemCost", "regularOpCost", "returnDataPerByteCost",
		"eventDataPerByteCost", "eventDataPerByteCost", "eventPerTopicCost", "eventBaseCost", "sandboxDataReadCost", "sandboxDataWriteCost",
		"maxEventTopics", "maxStackHeight", "maxMemoryPages", "enablePrintln", "maxSubjectLen",
	}
	s.ScaleDecoder.Init(data, subType, arg...)
}

type EraIndex struct {
	U32
}

type StakingLedgers struct {
	Struct
}

func (s *StakingLedgers) Init(data ScaleBytes, subType string, arg ...interface{}) {
	s.TypeString = "StakingLedgers<AccountId, RingBalanceOf, KtonBalanceOf, StakingBalance<RingBalanceOf, KtonBalanceOf>, Moment>"
	s.TypeMapping = map[string]string{
		"stash":               "AccountId",
		"total_ring":          "Compact<RingBalanceOf>",
		"total_regular_ring":  "Compact<RingBalanceOf>",
		"active_ring":         "Compact<RingBalanceOf>",
		"active_regular_ring": "Compact<RingBalanceOf>",
		"total_kton":          "Compact<KtonBalanceOf>",
		"active_kton":         "Compact<KtonBalanceOf>",
		"deposit_items":       "Vec<TimeDepositItem>",
		"unlocking":           "Vec<UnlockChunk>",
	}
	s.StructOrderField = []string{"stash", "total_ring", "total_regular_ring", "active_ring", "active_regular_ring", "total_kton", "active_kton", "deposit_items", "unlocking"}
	s.ScaleDecoder.Init(data, subType, arg...)
}

type TimeDepositItem struct {
	Struct
}

func (t *TimeDepositItem) Init(data ScaleBytes, subType string, arg ...interface{}) {
	t.TypeMapping = map[string]string{
		"value":       "Compact<RingBalanceOf>",
		"start_time":  "Compact<Moment>",
		"expire_time": "Compact<Moment>",
	}
	t.StructOrderField = []string{"value", "start_time", "expire_time"}
	t.ScaleDecoder.Init(data, subType, arg...)
}

type ExtendedBalance struct {
	U128
}

type RingBalanceOf struct {
	U128
}

type KtonBalanceOf struct {
	U128
}

type UnlockChunk struct {
	Struct
}

func (u *UnlockChunk) Init(data ScaleBytes, subType string, arg ...interface{}) {
	u.TypeMapping = map[string]string{
		"value":      "StakingBalance",
		"era":        "Compact<EraIndex>",
		"dt_power":   "ExtendedBalance",
		"is_regular": "bool",
	}
	u.StructOrderField = []string{"value", "era", "dt_power", "is_regular"}
	u.ScaleDecoder.Init(data, subType, arg...)
}

type RegularItem struct {
	Struct
}

func (r *RegularItem) Init(data ScaleBytes, subType string, arg ...interface{}) {
	r.TypeMapping = map[string]string{
		"value":       "Compact<RingBalanceOf>",
		"expire_time": "Compact<Moment>",
	}
	r.StructOrderField = []string{"value", "expire_time"}
	r.ScaleDecoder.Init(data, subType, arg...)
}

type StakingBalance struct {
	Enum
}

func (s *StakingBalance) Init(data ScaleBytes, subType string, arg ...interface{}) {
	s.TypeString = "StakingBalance<RingBalanceOf, KtonBalanceOf>"
	valueList := []string{"RingBalanceOf", "KtonBalanceOf"}
	s.Enum.Init(data, subType, valueList)
}

type Perbill struct {
	U32
}

type Gas struct {
	U64
}

type CodeHash struct {
	Hash
}

type Schedule struct {
	Struct
}

type EventIndex struct {
	U32
}

type DigestItem struct {
	Enum
}

func (s *DigestItem) Init(data ScaleBytes, subType string, arg ...interface{}) {
	valueList := []string{"Other", "AuthoritiesChange", "ChangesTrieRoot", "SealV0", "Consensus", "Seal", "PreRuntime"}
	s.Enum.Init(data, subType, valueList)
}

type Digest struct {
	Struct
}

func (d *Digest) Init(data ScaleBytes, subType string, arg ...interface{}) {
	d.TypeMapping = map[string]string{
		"logs": "Vec<DigestItem>",
	}
	d.StructOrderField = []string{"logs"}
	d.ScaleDecoder.Init(data, subType, arg...)
}

type Exposures struct {
	Struct
}

func (s *Exposures) Init(data ScaleBytes, subType string, arg ...interface{}) {
	s.TypeMapping = map[string]string{
		"total":  "ExtendedBalance",
		"own":    "ExtendedBalance",
		"others": "Vec<IndividualExpo>",
	}
	s.StructOrderField = []string{"total", "owner", "others"}
	s.TypeString = "Exposures<AccountId, ExtendedBalance>"
	s.Struct.Init(data, subType, arg...)
}

type VestingSchedule struct {
	Struct
}

func (v *VestingSchedule) Init(data ScaleBytes, subType string, arg ...interface{}) {
	v.TypeString = "VestingSchedule<Balance>"
	v.TypeMapping = map[string]string{
		"offset":        "Balance",
		"perBlock":      "Balance",
		"startingBlock": "BlockNumber",
	}
	v.StructOrderField = []string{"offset", "perBlock", "startingBlock"}
	v.Struct.Init(data, subType, arg...)
}

type StoredPendingChange struct {
	Struct
}

func (s *StoredPendingChange) Init(data ScaleBytes, subType string, arg ...interface{}) {
	s.TypeString = "StoredPendingChange<BlockNumber>"
	s.TypeMapping = map[string]string{
		"scheduledAt":     "BlockNumber",
		"delay":           "BlockNumber",
		"nextAuthorities": "Vec<NextAuthority>",
	}
	s.StructOrderField = []string{"scheduledAt", "delay", "nextAuthorities"}
	s.Struct.Init(data, subType, arg...)
}

type AuthorityId struct {
	AccountId
}

type NextAuthority struct {
	Struct
}

func (n *NextAuthority) Init(data ScaleBytes, subType string, arg ...interface{}) {
	n.TypeString = "(AuthorityId, u64)"
	n.Struct.Init(data, subType, arg...)
}
