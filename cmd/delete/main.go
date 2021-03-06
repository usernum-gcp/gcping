package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/ImJasonH/gcping/pkg/util"
	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

var project = flag.String("project", "gcping-1369", "Project to use")

func main() {
	flag.Parse()

	svc, err := compute.NewService(context.Background())
	if err != nil {
		log.Fatalf("NewService: %v", err)
	}

	// Delete all instances in parallel.
	if err := util.ForEachRegion(svc, *project, deleteVM); err != nil {
		log.Fatal(err)
	}
}

func deleteVM(svc *compute.Service, region string) error {
	zone := region + "-b"
	log.Println("Deleting instance:", region)
	start := time.Now()
	op, err := svc.Instances.Delete(*project, zone, region).Do()
	if herr, ok := err.(*googleapi.Error); ok && herr.Code == http.StatusNotFound {
		log.Printf("instances.delete (%s): already didn't exist", region)
		return nil
	}
	if err != nil {
		return err
	}
	if err := util.WaitForZoneOp(svc, *project, zone, op); err != nil {
		return err
	}
	log.Printf("instances.delete (%s): ok, took %s", region, time.Since(start))
	return nil
}
