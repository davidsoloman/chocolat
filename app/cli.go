package app

import (
	"fmt"
	"github.com/angdev/chocolat/config"
	"github.com/angdev/chocolat/model"
	"github.com/codegangsta/cli"
	"github.com/satori/go.uuid"
)

func NewCli() *cli.App {
	app := cli.NewApp()
	app.Name = "chocolat"
	app.Usage = "Self-Hosted Event Analytics Service"
	app.Version = "0.0.1"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "create, c",
			Usage: "Create a new project",
		},
		cli.BoolFlag{
			Name:  "run, r",
			Usage: "Run the server",
		},
		cli.BoolFlag{
			Name:  "list, l",
			Usage: "List projects",
		},
		cli.StringFlag{
			Name:  "project, p",
			Usage: "Inspect a project",
		},
		cli.StringFlag{
			Name:  "delete, d",
			Usage: "Delete a project",
		},
		cli.BoolFlag{
			Name:  "routes",
			Usage: "List routes",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("create") {
			createProject()
		} else if c.Bool("list") {
			listProject()
		} else if c.String("project") != "" {
			inspectProject(c.String("project"))
		} else if c.String("delete") != "" {
			deleteProject(c.String("delete"))
		} else if c.Bool("routes") {
			listRoutes()
		} else if c.Bool("run") {
			Chocolat.Run()
		} else {
			cli.ShowAppHelp(c)
		}
	}

	return app
}

func createProject() {
	db := model.DB()
	project := model.Project{
		UUID: uuid.NewV4().String(),
	}
	db.Create(&project)

	fmt.Printf("Created a new project.\n")
	fmt.Printf("Project UUID - %s\n", project.UUID)
	fmt.Printf("Project Master Key - %s\n", project.MasterKey().Value)
	fmt.Printf("Project Read Key - %s\n", project.ReadKey().Value)
	fmt.Printf("Project Write Key - %s\n", project.WriteKey().Value)
}

func listProject() {
	db := model.DB()
	var projects []model.Project

	if db.Find(&projects).RecordNotFound() {
		fmt.Println("No project found")
		return
	}

	for _, p := range projects {
		fmt.Println(p.UUID)
	}
}

func inspectProject(uuid string) {
	project := model.ProjectByUUID(uuid)

	if project == nil {
		fmt.Println("No project found")
		return
	}

	fmt.Printf("Project UUID - %s\n", project.UUID)
	fmt.Printf("Project Master Key - %s\n", project.MasterKey().Value)
	fmt.Printf("Project Read Key - %s\n", project.ReadKey().Value)
	fmt.Printf("Project Write Key - %s\n", project.WriteKey().Value)
}

func deleteProject(uuid string) {
	project := model.ProjectByUUID(uuid)

	if project == nil {
		fmt.Println("No project found")
		return
	}

	project.Delete()

	fmt.Printf("Delete a project %s\n", project.UUID)
}

func listRoutes() {
	routes := config.Routes
	fmt.Printf("%6s\t%s\n", "Method", "Path")
	for _, route := range routes {
		fmt.Printf("%6s\t%s\n", route.HttpMethod, route.PathExp)
	}
}
