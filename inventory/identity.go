package inventory

/**
* IDENTITY FAMILY - For Testing
**/
type IdentityOrder struct{}

func (this *IdentityOrder) NextState(accum State) (State, error) {
	return accum, nil
}

func (this *IdentityOrder) RenderEntry() string {
	return "identity"
}

func NewIdentityOrder() StateEntry {
	return new(IdentityOrder)
}
