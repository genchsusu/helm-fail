package main

import (
	"fmt"
	"log"
	"os"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/storage"
	"helm.sh/helm/v3/pkg/storage/driver"
	"k8s.io/client-go/kubernetes"
)

func main() {
	// Get the release name from command line argument
	releaseName := os.Args[1]

	// Create a new Helm configuration
	cfg := new(action.Configuration)

	// Initialize new Helm CLI settings
	settings := cli.New()

	// Initialize the Helm configuration with settings
	err := cfg.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new Helm status action client
	client := action.NewStatus(cfg)

	// Run the client to get the status of the release
	rel, err := client.Run(releaseName)
	if err != nil {
		log.Fatalf("Failed to get status of release %s: %v", releaseName, err)
	} else {
		// Check if the release status is pending
		if rel.Info.Status.IsPending() {
			// If yes, mark the release status as failed
			rel.SetStatus(release.StatusFailed, "manually marked as failed")

			// Load kubernetes config
			clientConfig, err := settings.RESTClientGetter().ToRESTConfig()
			if err != nil {
				log.Fatalf("Failed to load kubernetes config: %v", err)
			}

			// Create a new kubernetes client
			clientset, err := kubernetes.NewForConfig(clientConfig)
			if err != nil {
				log.Fatalf("Failed to create kubernetes client: %v", err)
			}

			// Create a new Secrets client
			secretsClient := clientset.CoreV1().Secrets(settings.Namespace())

			// Create a new storage driver with the secrets client
			storageDriver := driver.NewSecrets(secretsClient)

			// Create a new storage instance using the driver
			storage := storage.Init(storageDriver)

			// Update the release in the storage
			err = storage.Update(rel)
			if err != nil {
				fmt.Printf("Failed to update release: %v", err)
				return
			}

			fmt.Printf("Release '%s' status changed to 'failed'\n", releaseName)
		} else {
			fmt.Printf("Release '%s' status doesn't need to change\n", releaseName)
		}
	}
}
