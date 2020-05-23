package rpc

type plasm struct {
	query
}

func (r *plasm) GetCurrentEra() (int, error) {
	eraIndex, err := ReadStorage(r.c, "PlasmRewards", "CurrentEra", r.hash)
	if err != nil {
		return 0, err
	}
	return eraIndex.ToInt(), nil
}
