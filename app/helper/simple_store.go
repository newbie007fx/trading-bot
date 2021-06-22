package helper

type simpleTemp struct {
	store map[string]string
}

func (st *simpleTemp) Set(key, value string) {
	st.store[key] = value
}

func (st *simpleTemp) Get(key string) *string {
	value, ok := st.store[key]
	if ok {
		return &value
	}

	return nil
}

var simple_store *simpleTemp

func GetSimpleStore() *simpleTemp {
	if simple_store == nil {
		simple_store = &simpleTemp{
			store: map[string]string{},
		}
	}

	return simple_store
}
