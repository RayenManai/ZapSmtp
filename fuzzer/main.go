package main

import (
	"fmt"
	"log"
	"time"
)

func main() {
	// Run different fuzzing scenarios without signature
	fuzzer, err := NewFuzzer("", "", nil)
	if err != nil {
		log.Fatalf("Failed to create fuzzer without signature: %v", err)
	}
	runFuzzingTests(fuzzer)

	// Run different fuzzing scenarios with signature only
	fuzzer, err = NewFuzzer("_test/cert1.pem", "_test/key1.pem", nil)
	if err != nil {
		log.Fatalf("Failed to create fuzzer with signature only: %v", err)
	}
	runFuzzingTests(fuzzer)

	// Run different fuzzing scenarios with encryption only
	fuzzer, err = NewFuzzer("", "", []string{"_test/cert1.pem"})
	if err != nil {
		log.Fatalf("Failed to create fuzzer with encryption only: %v", err)
	}
	runFuzzingTests(fuzzer)

	// Run different fuzzing scenarios with signature and encryption
	fuzzer, err = NewFuzzer("_test/cert1.pem", "_test/key1.pem", []string{"_test/cert1.pem"})
	if err != nil {
		log.Fatalf("Failed to create fuzzer with signature and encryption: %v", err)
	}
	runFuzzingTests(fuzzer)

	defer fuzzer.Close()
}

func runFuzzingTests(fuzzer *Fuzzer) {

	// Test Case 1: normal output message
	fuzzer.FuzzLog("Null Pointer Exception: line 1337")
	time.Sleep(time.Second * 5)

	// Test case 2: different character formats
	fuzzer.FuzzLog(GenerateRandomLog())
	time.Sleep(time.Second * 5)

	// Test case 3: Random long message
	fuzzer.FuzzLog(GenerateRandomLongMessage(10000))
	time.Sleep(time.Second * 5)

	// Test case 4: start with many backspace characters
	fuzzer.FuzzLog(GenerateMalformedLog())
	time.Sleep(time.Second * 5)

	fmt.Println("Fuzzing completed")
}
