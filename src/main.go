/*
Package main .
*/
package main

import (
	"AWSResourcesToJson/resources"
	"encoding/json"
	"fmt"
	"log"
	"os"
)

func main() {
	region := os.Getenv("REGION")

	instances, err := resources.GetEC2Instances(region)
	if err != nil {
		log.Fatal(err)
	}

	loadbalancers, err := resources.GetALBResources(region)
	if err != nil {
		log.Fatal(err)
	}

	awsResources := resources.MergeResources(instances, loadbalancers)

	l, _ := json.Marshal(awsResources)
	fmt.Println(string(l))
}