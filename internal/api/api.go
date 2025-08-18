package api

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"taiga_storie_extractor/internal/config"
	"taiga_storie_extractor/internal/structs"
)

func GetAuthenticatedHeaders(authResponse structs.HeadersProvider) map[string]string {
	requestHeader := map[string]string{
		"Authorization":        "Bearer " + authResponse().Token,
		"Accept":               "application/json, text/plain, */*",
		"Accept-Language":      "pt-br",
		"Connection":           "keep-alive",
		"Sec-Fetch-Dest":       "empty",
		"Sec-Fetch-Mode":       "cors",
		"Sec-Fetch-Site":       "same-origin",
		"User-Agent":           "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36 OPR/120.0.0.0",
		"sec-ch-ua":            `"Opera";v="120", "Not-A.Brand";v="8", "Chromium";v="135"`,
		"sec-ch-ua-mobile":     "?0",
		"sec-ch-ua-platform":   `"Linux"`,
		"x-disable-pagination": "1",
	}
	return requestHeader
}

func GetToken() structs.AuthResponse {
	username := config.GetConfig().Username
	password := config.GetConfig().Password
	if username == "seu_user_name" || username == "" {
		fmt.Println("ERROR: Usuário invalido, preencha corretamente o arquivo `config.json`")
		consoleHook := bufio.NewReader(os.Stdin)
		fmt.Print("Aperte qualquer teclas para sair: ")
		_, _ = consoleHook.ReadString('\n')
		os.Exit(0)
	}
	if password == "sua_senha" || password == "" {
		fmt.Println("ERROR: Senha invalida, preencha corretamente o arquivo `config.json`")
		consoleHook := bufio.NewReader(os.Stdin)
		fmt.Print("Aperte qualquer teclas para sair: ")
		_, _ = consoleHook.ReadString('\n')
		os.Exit(0)
	}
	userCredentials := structs.UserCredentials{
		Type:     "normal",
		Username: config.GetConfig().Username,
		Password: config.GetConfig().Password,
	}

	jsonValue, UserCredentialsMarshalError := json.Marshal(userCredentials)

	if UserCredentialsMarshalError != nil {
		panic(UserCredentialsMarshalError)
	}

	response, AuthRequestError := http.Post(config.GetConfig().ApiBaseUrl+"/auth", "application/json", bytes.NewBuffer(jsonValue))
	if AuthRequestError != nil {
		panic(AuthRequestError)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	if response.StatusCode != http.StatusOK {
		panic(fmt.Sprintf("Erro na autenticação: %s", response.Status))
	}
	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}

	// Agora tente decodificar o JSON manualmente
	var result structs.AuthResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		panic(err)
	}

	return result
}

func GetAllProjects(authenticatedHeaders map[string]string) map[int]structs.Project {
	request, AllProjectsRequestError := http.NewRequest("GET", config.GetConfig().ApiBaseUrl+"/projects", nil)
	if request == nil {
		panic("Request Cannot be nil!")
	}
	query := request.URL.Query()
	query.Add("member", "20")
	query.Add("order_by", "user_order")
	query.Add("slight", "true")

	request.URL.RawQuery = query.Encode()

	if AllProjectsRequestError != nil {
		fmt.Println(AllProjectsRequestError)
	}

	for headerKey, headerValue := range authenticatedHeaders {
		request.Header.Set(headerKey, headerValue)
	}

	// Cria o client HTTP e faz a requisição
	client := &http.Client{}
	response, err := client.Do(request)

	defer func(Body io.ReadCloser) {
		if closerError := Body.Close(); err != nil {
			log.Printf("Erro ao fechar Body: %v", closerError)
		}
	}(response.Body)

	// Lê a resposta
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Agora tente decodificar o JSON manualmente
	var projects []structs.Project
	if err := json.Unmarshal(body, &projects); err != nil {
		panic(err)
	}
	mappedProjects := make(map[int]structs.Project)
	for _, project := range projects {
		mappedProjects[project.ID] = project

	}
	return mappedProjects
}

func GetAllStoriesFromBoard(authenticatedHeaders map[string]string, projectId int) []structs.Storie {
	request, StorieRequestError := http.NewRequest("GET", config.GetConfig().ApiBaseUrl+"/userstories", nil)
	if request == nil {
		panic("Request Cannot be nil!")
	}
	query := request.URL.Query()

	query.Add("project", fmt.Sprint(projectId))
	query.Add("status__is_archived", "false")

	request.URL.RawQuery = query.Encode()

	if StorieRequestError != nil {
		panic(StorieRequestError)
	}

	for headerKey, headerValue := range authenticatedHeaders {
		request.Header.Set(headerKey, headerValue)
	}

	// Cria o client HTTP e faz a requisição
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(response.Body)

	// Lê a resposta
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	// Agora tente decodificar o JSON manualmente
	var Stories []structs.Storie
	if err := json.Unmarshal(body, &Stories); err != nil {
		panic(err)
	}

	return Stories

}

func GetProjectDetailWithSlug(authenticatedHeaders map[string]string, projectSlug string) structs.Project {
	request, ProjectRequestError := http.NewRequest("GET", config.GetConfig().ApiBaseUrl+"/projects/by_slug", nil)
	if request == nil {
		panic("Request Cannot be nil!")
	}
	query := request.URL.Query()
	query.Add("slug", projectSlug)

	request.URL.RawQuery = query.Encode()

	if ProjectRequestError != nil {
		panic(ProjectRequestError)
	}

	for headerKey, headerValue := range authenticatedHeaders {
		request.Header.Set(headerKey, headerValue)
	}

	// Cria o client HTTP e faz a requisição
	client := &http.Client{}
	response, err := client.Do(request)

	defer func(Body io.ReadCloser) {
		bodyCloserErr := Body.Close()
		if bodyCloserErr != nil {
			panic(bodyCloserErr)
		}
	}(request.Body)
	if response == nil {
		panic("Response cannot be nil, stoping.")
	}
	// Lê a resposta
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Agora tente decodificar o JSON manualmente
	var project structs.Project
	if err := json.Unmarshal(body, &project); err != nil {
		panic(err)
	}

	return project
}

func GetStorieDetailsActivity(authenticatedHeaders map[string]string, storyId int) []structs.CustomAttribute {
	// Monta requisição
	request, err := http.NewRequest("GET", config.GetConfig().ApiBaseUrl+"/history/userstory/"+fmt.Sprint(storyId), nil)
	if err != nil {
		panic(err)
	}

	// Headers
	for headerKey, headerValue := range authenticatedHeaders {
		request.Header.Set(headerKey, headerValue)
	}

	// Faz requisição
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		bodyCloserError := Body.Close()
		if bodyCloserError != nil {
			panic(bodyCloserError)
		}
	}(response.Body)

	// Lê corpo
	body, err := io.ReadAll(response.Body)
	if err != nil {
		panic(err)
	}
	// Decodifica todas as versões
	var allVersions []map[string]interface{}
	if err := json.Unmarshal(body, &allVersions); err != nil {
		panic(err)
	}

	if len(allVersions) == 0 {
		return []structs.CustomAttribute{}
	}

	// Encontra a mais recente pelo created_at
	var CustomAttributeArray []structs.CustomAttribute
	diff := allVersions[0]["diff"].(map[string]interface{})
	customAttrsRaw, ok := diff["custom_attributes"]

	if !ok || customAttrsRaw == nil {
		return []structs.CustomAttribute{}
	}

	customAttrs, ok := customAttrsRaw.([]interface{})

	if !ok {
		panic("custom_attributes não é um array")
	}

	for _, group := range customAttrs {
		for _, attr := range group.([]interface{}) {
			attrMap := attr.(map[string]interface{})
			attId := int(attrMap["id"].(float64))

			CustomAttributeArray = append(CustomAttributeArray, structs.CustomAttribute{
				ID:    attId,
				Name:  attrMap["name"].(string),
				Type:  attrMap["type"].(string),
				Value: attrMap["value"].(string),
			})
		}
	}
	return CustomAttributeArray
}

func GetStorieDetailsComment(authenticatedHeaders map[string]string, storyId int) []structs.StorieDetails {
	request, StorieRequestError := http.NewRequest("GET", config.GetConfig().ApiBaseUrl+"/history/userstory/"+fmt.Sprint(storyId), nil)
	if request == nil {
		panic("request cannot be nil, stoping.")
	}
	query := request.URL.Query()
	query.Add("type", "comment")

	request.URL.RawQuery = query.Encode()

	if StorieRequestError != nil {
		panic(StorieRequestError)
	}

	for headerKey, headerValue := range authenticatedHeaders {
		request.Header.Set(headerKey, headerValue)
	}

	// Cria o client HTTP e faz a requisição
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err)
	}
	defer func(Body io.ReadCloser) {
		bodyCloserError := Body.Close()
		if bodyCloserError != nil {
			panic(bodyCloserError)
		}
	}(response.Body)

	// Lê a resposta
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}
	// Agora tente decodificar o JSON manualmente
	var StorieDetails []structs.StorieDetails

	if err := json.Unmarshal(body, &StorieDetails); err != nil {
		panic(err)
	}

	return StorieDetails
}

func UserStoryCustomAttributes(authenticatedHeaders map[string]string, storyId int) map[string]string {
	request, UserStoryCustomAttributesError := http.NewRequest("GET", fmt.Sprintf("%s/userstories/custom-attributes-values/%d", config.GetConfig().ApiBaseUrl, storyId), nil)

	if UserStoryCustomAttributesError != nil {
		panic(UserStoryCustomAttributesError)
	}

	for headerKey, headerValue := range authenticatedHeaders {
		request.Header.Set(headerKey, headerValue)
	}

	// Cria o client HTTP e faz a requisição
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		closeError := Body.Close()
		if closeError != nil {
			panic(closeError)
		}
	}(response.Body)

	// Lê a resposta
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	// Agora tente decodificar o JSON manualmente
	var CustomAttributeRawData map[string]map[string]string
	if unmarshalError := json.Unmarshal(body, &CustomAttributeRawData); err != nil {
		panic(unmarshalError)
	}

	AttributeMap := make(map[string]string)

	for AttributeID, AttributeValue := range CustomAttributeRawData["attributes_values"] {
		for _, ActivityData := range GetStorieDetailsActivity(authenticatedHeaders, storyId) {
			if AttributeID == strconv.Itoa(ActivityData.ID) && AttributeValue == ActivityData.Value {
				_, exist := AttributeMap[ActivityData.Name]
				if !exist {
					AttributeMap[ActivityData.Name] = ActivityData.Value
				}
			}
		}
	}
	return AttributeMap
}
