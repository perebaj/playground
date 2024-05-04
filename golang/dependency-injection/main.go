package main

type SportPlace interface {
	Play()
}

type Sport struct {
	sp SportPlace
}

func NewSport(sportPlace SportPlace) *Sport {
	return &Sport{sp: sportPlace}
}

type Futebool struct{}

func (f *Futebool) Play() {
	println("Playing Futbool")
}

type Basquete struct{}

func (b *Basquete) Play() {
	println("Playing Basquete")
}

func main() {
	sports := NewSport(&Futebool{})
	sports.sp.Play()

	sports = NewSport(&Basquete{})
	sports.sp.Play()

	sports2 := NewSport2("futebool")
	sports2.Play()
	sports2 = NewSport2("basquete")
	sports2.Play()
}

// Certanly not a good example of dependency injection ⏬⏬⏬

type Sport2 struct {
	kind string
}

func NewSport2(kind string) *Sport2 {
	return &Sport2{kind: kind}
}

func (s *Sport2) Play() {
	switch s.kind {
	case "futebool":
		println("Playing Futbool")
	case "basquete":
		println("Playing Basquete")
	}
}
