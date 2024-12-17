package main

import (
	"bufio"
	"log/slog"
	"os"
	"regexp"
	"strconv"
)

func ExtractUncorruptedMuls(corruptedString string) []string {
	re := regexp.MustCompile(`mul\([0-9]{1,3},[0-9]{1,3}\)`)
	uncorruptedMuls := re.FindAllString(corruptedString, -1)
	return uncorruptedMuls
}

func ExtractUncorruptedMulsWithInstructions(corruptedString string) []string {
	re := regexp.MustCompile(`mul\([0-9]{1,3},[0-9]{1,3}\)|do\(\)|don't\(\)`)
	uncorruptedMulsWitInstructions := re.FindAllString(corruptedString, -1)
	return uncorruptedMulsWitInstructions
}

func ProcessMuls(muls []string) int {
	result := 0
	re := regexp.MustCompile(`[0-9]{1,3}`)
	for _, mul := range muls {
		operandsString := re.FindAllString(mul, 2)
		leftOperand, err := strconv.Atoi(operandsString[0])
		rightOperand, err := strconv.Atoi(operandsString[1])
		if err != nil {
			slog.Warn("day 2.1", "err", err, "mul", mul, "operands", operandsString)
		} else {
			result += leftOperand * rightOperand
		}
	}
	return result
}

func ProcessMulsWithInstructions(mulsWithInstructions []string, process bool) (int, bool) {
	result := 0
	operandRe := regexp.MustCompile(`[0-9]{1,3}`)
	dos := 0
	donts := 0
	for _, mulOrInstruction := range mulsWithInstructions {
		if mulOrInstruction == "do()" {
			process = true
			// slog.Debug("day 2.2", "do()", mulOrInstruction)
			dos += 1
			continue
		} else if mulOrInstruction == "don't()" {
			// slog.Debug("day 2.2", "don't()", mulOrInstruction)
			process = false
			donts += 1
			continue
		}
		if process {
			operandsString := operandRe.FindAllString(mulOrInstruction, 2)
			leftOperand, err := strconv.Atoi(operandsString[0])
			rightOperand, err := strconv.Atoi(operandsString[1])
			if err != nil {
				slog.Warn("day 2.2", "err", err, "mul", mulOrInstruction, "operands", operandsString)
			} else {
				result += leftOperand * rightOperand
			}
			// slog.Debug("day 2.2", "process", mulOrInstruction)
		} else {
			// slog.Debug("day 2.2", "skip", mulOrInstruction)
		}
	}
	slog.Debug("day 2.2", "dos", dos, "donts", donts)
	return result, process
}

func ReadCorruptedStrings(path string) ([]string, error) {
	inputFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(inputFile)
	scanner.Split(bufio.ScanLines)

	corruptedStrings := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		corruptedStrings = append(corruptedStrings, line)
	}
	return corruptedStrings, nil
}

func ComputeFirstResult(corruptedStrings []string) int {
	result := 0
	for _, corruptedString := range corruptedStrings {
		muls := ExtractUncorruptedMuls(corruptedString)
		partialResult := ProcessMuls(muls)
		result += partialResult
		slog.Debug("day 2.1", "partialResult", partialResult)
	}
	return result
}

func ComputeSecondResult(corruptedStrings []string) int {
	result := 0
	process := true
	for _, corruptedString := range corruptedStrings {
		muls := ExtractUncorruptedMulsWithInstructions(corruptedString)
		partialResult, dos := ProcessMulsWithInstructions(muls, process)
		process = dos
		result += partialResult
		slog.Debug("day 2.2", "partialResult", partialResult)
	}
	return result
}

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(logger)

	path := "/Users/giacomoponziani/projects/advent-of-code/day3/input_01.txt"

	corruptedStrings, err := ReadCorruptedStrings(path)
	if err != nil {
		slog.Error("day 2", "err", err)
		return
	}

	firstResult := ComputeFirstResult(corruptedStrings)
	secondResult := ComputeSecondResult(corruptedStrings)

	slog.Info("day 2", "first result", firstResult, "second result", secondResult)
}
