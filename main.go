package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"

	database "cloud.google.com/go/spanner/admin/database/apiv1"
	instance "cloud.google.com/go/spanner/admin/instance/apiv1"
	adminpb "google.golang.org/genproto/googleapis/spanner/admin/database/v1"
	instancepb "google.golang.org/genproto/googleapis/spanner/admin/instance/v1"
)

func createInstance(projectId, instanceId string) error {
	ctx := context.Background()
	instanceAdmin, err := instance.NewInstanceAdminClient(ctx)
	if err != nil {
		return err
	}
	defer instanceAdmin.Close()

	instanceString := fmt.Sprintf("projects/%s/instances/%s", projectId, instanceId)

	ins, _ := instanceAdmin.GetInstance(ctx, &instancepb.GetInstanceRequest{Name: instanceString})
	if ins != nil && ins.State.String() == "READY" {
		fmt.Printf("Instance already created [%s]\n", instanceString)
		return nil
	}

	op, err := instanceAdmin.CreateInstance(ctx, &instancepb.CreateInstanceRequest{
		Parent:     fmt.Sprintf("projects/%s", projectId),
		InstanceId: instanceId,
		Instance: &instancepb.Instance{
			Config:      fmt.Sprintf("projects/%s/instanceConfigs/%s", projectId, "default"),
			DisplayName: instanceId,
		},
	})
	if err != nil {
		return fmt.Errorf("could not create instance %s: %v", instanceString, err)
	}
	// Wait for the instance creation to finish.
	i, err := op.Wait(ctx)
	if err != nil {
		return fmt.Errorf("waiting for instance creation to finish failed: %v", err)
	}
	// The instance may not be ready to serve yet.
	if i.State != instancepb.Instance_READY {
		fmt.Printf("instance state is not READY yet. Got state %v\n", i.State)
	}
	fmt.Printf("Created instance [%s]\n", instanceId)
	return nil
}

func createDatabase(projectId, instanceId, databaseId string) error {
	ctx := context.Background()
	adminClient, err := database.NewDatabaseAdminClient(ctx)
	if err != nil {
		return err
	}
	defer adminClient.Close()

	dbString := fmt.Sprintf("projects/%s/instances/%s/databases/%s", projectId, instanceId, databaseId)

	db, _ := adminClient.GetDatabase(ctx, &adminpb.GetDatabaseRequest{Name: dbString})
	if db != nil && db.State.String() == "READY" {
		fmt.Printf("Database already created [%s]\n", databaseId)
		return nil
	}

	op, err := adminClient.CreateDatabase(ctx, &adminpb.CreateDatabaseRequest{
		Parent:          fmt.Sprintf("projects/%s/instances/%s", projectId, instanceId),
		CreateStatement: fmt.Sprintf("CREATE DATABASE `%s`", databaseId),
	})
	if err != nil {
		return err
	}
	if _, err := op.Wait(ctx); err != nil {
		return err
	}
	fmt.Printf("Created database [%s]\n", databaseId)
	return nil
}

func getIds(args []string) (string, string, string, error) {
	projectId := os.Getenv("SPANNER_PROJECT_ID")
	instanceId := os.Getenv("SPANNER_INSTANCE_ID")
	databaseId := os.Getenv("SPANNER_DATABASE_ID")
	if len(args) > 2 {
		return projectId, instanceId, databaseId, errors.New("Too many args.")
	} else if len(args) == 2 {
		matches := regexp.MustCompile("^(((((projects/)?([^/]*)/)?instances/)?([^/]*)/)?databases/)?([^/]*)$").FindStringSubmatch(args[1])
		if len(matches) > 0 {
			if len(matches[len(matches)-1]) > 0 {
				databaseId = matches[len(matches)-1]
			}
			if len(matches[len(matches)-2]) > 0 {
				instanceId = matches[len(matches)-2]
			}
			if len(matches[len(matches)-3]) > 0 {
				projectId = matches[len(matches)-3]
			}
		}
	}

	if len(projectId) == 0 {
		return projectId, instanceId, databaseId, errors.New("Could not find ProjectId.")
	}

	if len(instanceId) == 0 {
		return projectId, instanceId, databaseId, errors.New("Could not find InstanceId.")
	}

	if len(databaseId) == 0 {
		return projectId, instanceId, databaseId, errors.New("Could not find DatabaseId.")
	}

	return projectId, instanceId, databaseId, nil
}

func create(args []string) error {
	projectId, instanceId, databaseId, err := getIds(args)

	if err != nil {
		return err
	}

	err = createInstance(projectId, instanceId)

	if err != nil {
		return err
	}

	err = createDatabase(projectId, instanceId, databaseId)

	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := create(os.Args)
	if err != nil {
		fmt.Printf("Error:\n")
		fmt.Printf("  %s\n\n", err.Error())
		fmt.Printf("Usage:\n")
		fmt.Printf("  spanner-createdb {databaseId}\n")
		fmt.Printf("  spanner-createdb databases/{databaseId}\n")
		fmt.Printf("  spanner-createdb {instanceId}/databases/{databaseId}\n")
		fmt.Printf("  spanner-createdb instances/{instanceId}/databases/{databaseId}\n")
		fmt.Printf("  spanner-createdb {projectId}/instances/{instanceId}/databases/{databaseId}\n")
		fmt.Printf("  spanner-createdb projects/{projectId}/instances/{instanceId}/databases/{databaseId}\n")
		fmt.Printf("\n")
		fmt.Printf("You can also pass the ids via environment variables:\n")
		fmt.Printf("  SPANNER_PROJECT_ID\n")
		fmt.Printf("  SPANNER_INSTANCE_ID\n")
		fmt.Printf("  SPANNER_DATABASE_ID\n")
		fmt.Printf("\n")
		os.Exit(1)
	} else {
		os.Exit(0)
	}
}
