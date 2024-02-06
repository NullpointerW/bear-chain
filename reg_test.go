package main

import (
	"fmt"
	"regexp"
	"testing"
	"time"
)

func TestReg(t *testing.T) {
	exp := `(\d{1,2})h(\d{1,2})m(\d{1,2})s`
	reg := regexp.MustCompile(exp)
	match := reg.FindStringSubmatch(`{"msg":"You have exceeded the rate limit. Please wait 7h59m4s before you try again"}`)
	fmt.Println(len(match))
	for _, s := range match {
		fmt.Println(s)
	}
}

func TestEstimatedTime(t *testing.T) {
	//d := estimatedTime(`{"msg":"You have exceeded the rate limit. Please wait 7h59m4s before you try again"}`)
	d, _ := time.ParseDuration("8h")
	fmt.Println(d)
}
