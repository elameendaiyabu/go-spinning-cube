package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

const (
	width  = 100
	height = 30
)

var (
	A, B, C             float64
	cubeWidth           float64 = 10
	zBuffer             [width * height]float64
	buffer              [width * height]byte
	backgroundASCIICode byte    = ' '
	distanceFromCam     float64 = 100
	horizontalOffset    float64
	K1                  float64 = 40
	incrementSpeed      float64 = 0.6
	x, y, z             float64
	ooz                 float64
	xp, yp              int
	idx                 int
)

func calculateX(i, j, k float64) float64 {
	return j*math.Sin(A)*math.Sin(B)*math.Cos(C) - k*math.Cos(A)*math.Sin(B)*math.Cos(C) +
		j*math.Cos(A)*math.Sin(C) + k*math.Sin(A)*math.Sin(C) + i*math.Cos(B)*math.Cos(C)
}

func calculateY(i, j, k float64) float64 {
	return j*math.Cos(A)*math.Cos(C) + k*math.Sin(A)*math.Cos(C) -
		j*math.Sin(A)*math.Sin(B)*math.Sin(C) + k*math.Cos(A)*math.Sin(B)*math.Sin(C) -
		i*math.Cos(B)*math.Sin(C)
}

func calculateZ(i, j, k float64) float64 {
	return k*math.Cos(A)*math.Cos(B) - j*math.Sin(A)*math.Cos(B) + i*math.Sin(B)
}

func calculateForSurface(cubeX, cubeY, cubeZ float64, ch byte) {
	x = calculateX(cubeX, cubeY, cubeZ)
	y = calculateY(cubeX, cubeY, cubeZ)
	z = calculateZ(cubeX, cubeY, cubeZ) + distanceFromCam

	ooz = 1 / z

	xp = int(float64(width)/2 + horizontalOffset + K1*ooz*x*2)
	yp = int(float64(height)/2 + K1*ooz*y)

	idx = xp + yp*width
	if idx >= 0 && idx < width*height {
		if ooz > zBuffer[idx] {
			zBuffer[idx] = ooz
			buffer[idx] = ch
		}
	}
}

func main() {
	fmt.Print("\x1b[2J")         // Clear screen
	fmt.Print("\x1b[?25l")       // Hide cursor
	defer fmt.Print("\x1b[?25h") // Show cursor on exit

	for {
		// Reset buffers
		for i := range buffer {
			buffer[i] = backgroundASCIICode
			zBuffer[i] = 0
		}

		horizontalOffset = 20 * cubeWidth
		for cubeX := -cubeWidth; cubeX < cubeWidth; cubeX += incrementSpeed {
			for cubeY := -cubeWidth; cubeY < cubeWidth; cubeY += incrementSpeed {
				calculateForSurface(cubeX, cubeY, -cubeWidth, '@')
				calculateForSurface(cubeWidth, cubeY, cubeX, '$')
				calculateForSurface(-cubeWidth, cubeY, -cubeX, '~')
				calculateForSurface(-cubeX, cubeY, cubeWidth, '#')
				calculateForSurface(cubeX, -cubeWidth, -cubeY, ';')
				calculateForSurface(cubeX, cubeWidth, cubeY, '+')
			}
		}
		// Prepare and print the frame
		var output strings.Builder
		output.WriteString("\x1b[H") // Move cursor to home position
		for k := 0; k < width*height; k++ {
			if k%width == 0 && k > 0 {
				output.WriteByte('\n')
			}
			output.WriteByte(buffer[k])
		}
		fmt.Print(output.String())

		// Update rotation angles
		A += 0.05
		B += 0.05
		C += 0.01

		time.Sleep(16 * time.Millisecond)
	}
}
