package register

type exists = struct{}

type Register map[any]exists

func (r Register) Add(dep any) {
	r[dep] = exists{}
}

func (r Register) Exists(dep any) bool {
	_, ok := r[dep]
	return ok
}

func (r Register) Del(dep any) {
	delete(r, dep)
}

func New() Register {
	return Register{}
}
