package main

type Simulation struct {
	Time   float64
	Scale  float64
	Points []PointMass
}

func (s *Simulation) Step() {

	// apply newtownian gravity
	for i := 0; i < len(s.Points); i++ {
		for j := i + 1; j < len(s.Points); j++ {
			fg := s.Points[i].ForceV(&s.Points[j])
			s.Points[i].ApplyForce(fg, s.Time)
			s.Points[j].ApplyForce(fg.Mul(-1), s.Time)
		}
	}
	for i := range s.Points {
		s.Points[i].Move(s.Time)
	}
}
