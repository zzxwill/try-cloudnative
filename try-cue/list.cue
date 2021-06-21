p: {
	a?: int
	b:  1
}

c: "aaa": [{
	for _, v in {
		b: p.b
		if p.a != _|_ {
			a: p.a
		}
	} {v}
}]
