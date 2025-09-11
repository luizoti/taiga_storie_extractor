package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"taiga_storie_extractor/internal/api"
	"taiga_storie_extractor/internal/excel"
	"taiga_storie_extractor/internal/structs"
	"taiga_storie_extractor/internal/versioning"
	"time"
)

func main() {
	fmt.Printf("Application: %s\n", versioning.AppName)
	fmt.Printf("Version	  : %s\n", versioning.Version)
	fmt.Println()
	auth := api.GetToken()
	headers := api.GetAuthenticatedHeaders(func() structs.AuthResponse { return auth })
	fmt.Println("Lista de projetos:")

	projects := api.GetAllProjects(headers)

	for _, project := range projects {
		fmt.Printf("\t %d - %s - %s\n", project.ID, project.Slug, project.Name)
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Digite o ID do projeto: ")
	userInput, userInputError := reader.ReadString('\n')

	if userInputError != nil {
		panic(userInputError)
	}
	projectSelected := strings.Trim(userInput, "\n")
	projectSelected = strings.Trim(projectSelected, "\r")
	projectSelected = strings.Trim(projectSelected, "\t")
	projectSelected = strings.TrimSpace(projectSelected)
	selectedProjectID, selectedProjectAtoiError := strconv.Atoi(projectSelected)
	if selectedProjectAtoiError != nil {
		fmt.Println()
		fmt.Println(selectedProjectAtoiError)
		fmt.Printf("ERROR: O valor inserido precisa ser um numeral: %s", userInput)
		consoleHook := bufio.NewReader(os.Stdin)
		fmt.Print("Aperte qualquer teclas para sair: ")
		_, _ = consoleHook.ReadString('\n')
		os.Exit(0)
	}
	_, existProjectInMap := projects[selectedProjectID]
	if !existProjectInMap {
		fmt.Println()
		fmt.Printf("ERROR: O projeto com ID: %d não existe, vamos recomeçar.\n", selectedProjectID)
		main()
	}
	stories := api.GetAllStoriesFromBoard(headers, selectedProjectID)

	// Buscar comentários de cada história
	storyCommentsMap := make(map[int][]structs.StorieDetails)
	customFieldsMap := make(map[int]map[string]string)
	for _, story := range stories {
		customFieldsMap[story.ID] = api.UserStoryCustomAttributes(headers, story.ID)
		storyCommentsMap[story.ID] = api.GetStorieDetailsComment(headers, story.ID)
		fmt.Printf("Card processado: %d - %s\n", story.ID, story.Name)
	}
	now := time.Now()

	selectedProject := projects[selectedProjectID]
	projectSlug := projects[selectedProjectID].Slug
	reader = bufio.NewReader(os.Stdin)
	fmt.Println()
	fmt.Println("Formato de relatório:")
	fmt.Println("\t 1. Comments")
	fmt.Println("\t 2. Historias")
	fmt.Print("Escolha o formato: ")
	userInputType, userInputError := reader.ReadString('\n')

	switch userInputType {
	case "1\n":
		err := excel.ExportMergedComments(
			selectedProject,
			stories,
			storyCommentsMap,
			customFieldsMap,
			filepath.Join(projectSlug, fmt.Sprintf("%s_relatorio_taiga_%s.xlsx", now.Format("2006-01-02_15-04-05"), projectSlug)))

		if err != nil {
			panic(err)
		}
	case "2\n":
		err := excel.ExportStoriesOnly(
			selectedProject,
			stories,
			customFieldsMap,
			filepath.Join(projectSlug, fmt.Sprintf("%s_relatorio_taiga_%s.xlsx", now.Format("2006-01-02_15-04-05"), projectSlug)))

		if err != nil {
			panic(err)
		}
	}

}
