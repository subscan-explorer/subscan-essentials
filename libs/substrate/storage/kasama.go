package storage

type RawBabePreDigest struct {
	Primary   *RawBabePreDigestPrimary   `json:"primary,omitempty"`
	Secondary *RawBabePreDigestSecondary `json:"secondary,omitempty"`
}

type RawBabePreDigestPrimary struct {
	AuthorityIndex uint   `json:"authorityIndex"`
	SlotNumber     uint64 `json:"slotNumber"`
	Weight         uint   `json:"weight"`
	VrfOutput      string `json:"vrfOutput"`
	VrfProof       string `json:"vrfProof"`
}

type RawBabePreDigestSecondary struct {
	AuthorityIndex uint   `json:"authorityIndex"`
	SlotNumber     uint64 `json:"slotNumber"`
	Weight         uint   `json:"weight"`
}
