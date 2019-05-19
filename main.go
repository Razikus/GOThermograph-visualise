package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/png"
	_ "image/png"
	"math"
	"os"
	"strconv"
	"strings"
)

type scalator struct {
	yMin float64
	yMax float64
	xMin float64
	xMax float64
}

func (scalat *scalator) scalator(x float64) float64 {
	h := scalat.xMax - scalat.xMin
	return scalat.yMin + ((scalat.yMax-scalat.yMin)/h)*(x-scalat.xMin)
}

type minmax struct {
	min float64
	max float64
}

func getMinMax(twodArray *[][]float64) minmax {
	min := (*twodArray)[0][0]
	max := (*twodArray)[0][0]
	for _, val := range *twodArray {
		for _, toCheck := range val {
			if toCheck > max {
				max = toCheck
			}
			if toCheck < min {
				min = toCheck
			}
		}
	}
	return minmax{min: min, max: max}
}

type termoImageProcessor struct {
	values         *[][]float64
	image          *image.NRGBA
	scalatorStruct scalator
}

func readValues(path string) *[][]float64 {
	inFile, _ := os.Open(path)
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	var toReturnValues = make([][]float64, 0)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		values := strings.Fields(scanner.Text())
		constructedValues := make([]float64, 0)
		for _, val := range values {
			floated, _ := strconv.ParseFloat(val, 64)
			constructedValues = append(constructedValues, floated)
		}
		toReturnValues = append(toReturnValues, constructedValues)

	}
	return &toReturnValues
}

func (processor *termoImageProcessor) processImage() bool {
	for y, xArray := range *processor.values {
		for x, value := range xArray {
			l := processor.scalatorStruct.scalator(value)
			pi4 := math.Pi * 4
			pi2 := math.Pi * 2
			threeMult := 3.0 * 255
			rCalculated := ((1 + math.Cos((pi4/threeMult)*l)) / 2) * 255
			gCalculated := ((1 + math.Cos((pi4/threeMult)*l-(pi2/3))) / 2) * 255
			bCalculated := ((1 + math.Cos((pi4/threeMult)*l-(pi4/3))) / 2) * 255
			if rCalculated > 255 {
				rCalculated = 255
			}
			if gCalculated > 255 {
				gCalculated = 255
			}
			if bCalculated > 255 {
				bCalculated = 255
			}
			r := uint8(rCalculated)
			g := uint8(gCalculated)
			b := uint8(bCalculated)
			newPixelData := color.NRGBA{R: r, G: g, B: b, A: 255}
			processor.image.SetNRGBA(x, y, newPixelData)

		}
	}
	return true
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: fileWithValues outImage")
		os.Exit(1)
	}

	var readed = readValues(os.Args[1])
	height := len(*readed)
	width := len((*readed)[0])
	newDataRect := image.Rect(0, 0, width, height)
	newDataImage := image.NewNRGBA(newDataRect)
	rgbMinMax := minmax{0, 255}
	minMaxThermoValues := getMinMax(readed)
	scalatorStruct := scalator{rgbMinMax.min, rgbMinMax.max, minMaxThermoValues.max, minMaxThermoValues.min}
	processor := termoImageProcessor{values: readed, scalatorStruct: scalatorStruct, image: newDataImage}
	processor.processImage()
	outputFile, err := os.Create(os.Args[2])
	if err != nil {
	}
	png.Encode(outputFile, newDataImage)
}
