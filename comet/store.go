package comet

type clientUs struct {
	m *SMap
}
type ctrlUs struct {
	m *SMap
}

func (c *clientUs) get(id string) *client {
	v := c.m.Get(id)
	if v == nil {
		return nil
	}
	return v.(*client)
}
