package utils

func CelciusToKelvin(celcius float64) float64 {
	return celcius + 273
}

func CelciusToFahrenheit(celcius float64) float64 {
	return celcius*9/5 + 32.0
}
