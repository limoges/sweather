package weather

type Provider interface {
	Temperature(lat, lon string) (Temperature, error)
}
