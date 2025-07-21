package catalogrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestIt(t *testing.T) {
	ctx := context.TODO()
	repo := New("/Users/marcgrol/src/learnmcp/internal/servicecatalog/service-catalog.sqlite")
	err := repo.Open(ctx)
	if err != nil {
		log.Fatalf("connect error: %s", err)
	}
	defer repo.Close(ctx)

	{
		keyword := "kyc"
		fmt.Printf("List modules with keyword %s:\n", keyword)
		modules, err := repo.ListModules(ctx, keyword)
		if err != nil {
			log.Fatalf("select error: %s", err)
		}
		for _, m := range modules {
			fmt.Printf("%+v\n", m)
		}
	}

	if false {
		{
			modules, err := repo.ListModules(ctx, "")
			if err != nil {
				log.Fatalf("select error: %s", err)
			}
			for _, m := range modules {
				fmt.Printf("%+v\n", m)
			}
		}
		{
			interfaces, err := repo.ListInterfaces(ctx)
			if err != nil {
				log.Fatalf("select error: %s", err)
			}
			for _, m := range interfaces {
				fmt.Printf("%+v\n", m)
			}
		}

		{
			interfaceID := "psp"
			fmt.Printf("Details of module %s:\n", interfaceID)

			module, exists, err := repo.GetModuleOnID(ctx, interfaceID)
			if err != nil {
				log.Fatalf("get error: %s", err)
			}
			if !exists {
				log.Fatalf("module not exists")
			}
			asJson, _ := json.MarshalIndent(module, "", "  ")
			fmt.Printf("%s\n", asJson)
		}

		{
			interfaceID := "com.adyen.services.acm.AcmService"
			fmt.Printf("Details of interface %s:\n", interfaceID)

			module, exists, err := repo.GetInterfaceOnID(ctx, interfaceID)
			if err != nil {
				log.Fatalf("get interface error: %s", err)
			}
			if !exists {
				log.Fatalf("interface not exists")
			}
			asJson, _ := json.MarshalIndent(module, "", "  ")
			fmt.Printf("%s\n", asJson)
		}

		{
			interfaceID := "com.adyen.services.acm.AcmService"
			fmt.Printf("Modules consuming interface %s:\n", interfaceID)

			module, exists, err := repo.ListInterfaceConsumers(ctx, interfaceID)
			if err != nil {
				log.Fatalf("list interface consumers error: %s", err)
			}
			if !exists {
				log.Fatalf("interface not exists")
			}
			asJson, _ := json.MarshalIndent(module, "", "  ")
			fmt.Printf("%s\n", asJson)
		}

		{
			databaseID := "billing"
			fmt.Printf("Modules consuming database %s:\n", databaseID)

			module, exists, err := repo.ListDatabaseConsumers(ctx, databaseID)
			if err != nil {
				log.Fatalf("list database consumers error: %s", err)
			}
			if !exists {
				log.Fatalf("database not exists")
			}
			asJson, _ := json.MarshalIndent(module, "", "  ")
			fmt.Printf("%s\n", asJson)
		}

		{
			teamID := "CustomerArea"
			fmt.Printf("Modules owned by team %s:\n", teamID)

			module, exists, err := repo.ListModulesOfTeam(ctx, teamID)
			if err != nil {
				log.Fatalf("list modules of team error: %s", err)
			}
			if !exists {
				log.Fatalf("team not exists")
			}
			asJson, _ := json.MarshalIndent(module, "", "  ")
			fmt.Printf("%s\n", asJson)
		}
	}
}
