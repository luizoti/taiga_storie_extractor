package excel

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"taiga_storie_extractor/internal/config"
	"taiga_storie_extractor/internal/structs"

	"github.com/xuri/excelize/v2"
)

func ExportStoriesOnly(
	project structs.Project,
	stories []structs.Storie,
	CustomFieldsMap map[int]map[string]string,
	outputPath string,
) error {
	f := excelize.NewFile()
	sheet := "Historias"
	err := f.SetSheetName("Sheet1", sheet)
	if err != nil {
		return err
	}

	// Cabeçalhos atualizados com todos os campos
	headers := []string{
		"Projeto Criado em",
		"Projeto Modificado em",

		"História Criado em",
		"História Modificado em",
		"História Concluída em",
		"História Data Limite",

		"Projeto ID",
		"Projeto Nome",
		"Projeto Slug",
		"Projeto Descrição",
		"História ID", "História Ref",
		"História Nome",
		"História Motivo da Data Limite",
		"História Status da Data Limite",
		"História Comentário Inicial",
	}

	rowIndex := 2

	for _, story := range stories {
		row := []string{
			string(project.CreatedDate),
			string(project.ModifiedDate),

			string(story.CreatedDate),
			string(story.ModifiedDate),
			string(story.FinishDate),
			string(story.DueDate),

			strconv.Itoa(project.ID),
			project.Name,
			project.Slug,
			project.Description,

			strconv.Itoa(story.ID),
			strconv.Itoa(story.Ref),
			story.Name,
			story.DueDateReason,
			story.DueDateStatus,
			story.Comment,
		}
		for name, value := range CustomFieldsMap[story.ID] {
			row = append(row, value)
			if !slices.Contains(headers, name) {
				headers = append(headers, name)
			}
		}
		cell, _ := excelize.CoordinatesToCellName(1, rowIndex)
		err := f.SetSheetRow(sheet, cell, &row)
		if err != nil {
			return err
		}
		rowIndex++
	}

	err = f.SetSheetRow(sheet, "A1", &headers)
	if err != nil {
		return err
	}
	workDirectory, _ := config.GetWorkDirectory()
	exportsDirectory := filepath.Join(workDirectory, "exports", filepath.Dir(outputPath))
	makeExportDirectoryError := os.MkdirAll(exportsDirectory, 0777)
	if makeExportDirectoryError != nil {
		panic(fmt.Sprintf("Cannot Create %f", makeExportDirectoryError))

	}
	xlsxFilePath := filepath.Join(exportsDirectory, filepath.Base(outputPath))
	fmt.Printf("✅ Salvo em: %s\n", xlsxFilePath)

	if err := f.SaveAs(xlsxFilePath); err != nil {
		return fmt.Errorf("erro ao salvar planilha: %w", err)
	}
	consoleHook := bufio.NewReader(os.Stdin)
	fmt.Print("Aperte qualquer teclas para sair: ")
	_, _ = consoleHook.ReadString('\n')

	return nil
}

func ExportMergedComments(
	project structs.Project,
	stories []structs.Storie,
	storyComments map[int][]structs.StorieDetails,
	CustomFieldsMap map[int]map[string]string,
	outputPath string,
) error {
	f := excelize.NewFile()
	sheet := "Comentários"
	err := f.SetSheetName("Sheet1", sheet)
	if err != nil {
		return err
	}

	// Cabeçalhos atualizados com todos os campos
	headers := []string{
		"Projeto Criado em",
		"Projeto Modificado em",

		"História Criado em",
		"História Modificado em",
		"História Concluída em",
		"História Data Limite",

		"Comentário Criado em",

		"Projeto ID",
		"Projeto Nome",
		"Projeto Slug",
		"Projeto Descrição",
		"História ID", "História Ref",
		"História Nome",
		"História Motivo da Data Limite",
		"História Status da Data Limite",
		"História Comentário Inicial",
		"Comentário ID",
		"Comentário Texto",
		"Comentário HTML",
	}

	rowIndex := 2

	for _, story := range stories {

		comments := storyComments[story.ID]
		for _, comment := range comments {
			row := []string{
				string(project.CreatedDate),
				string(project.ModifiedDate),

				string(story.CreatedDate),
				string(story.ModifiedDate),
				string(story.FinishDate),
				string(story.DueDate),

				string(comment.CreatedDate),

				strconv.Itoa(project.ID),
				project.Name,
				project.Slug,
				project.Description,
				//project.CsvUUID,

				strconv.Itoa(story.ID),
				strconv.Itoa(story.Ref),
				story.Name,
				story.DueDateReason,
				story.DueDateStatus,
				story.Comment,

				comment.UuId,
				comment.Comment,
				comment.CommentHtml,
			}

			for name, value := range CustomFieldsMap[story.ID] {
				row = append(row, value)
				if !slices.Contains(headers, name) {
					headers = append(headers, name)
				}
			}
			cell, _ := excelize.CoordinatesToCellName(1, rowIndex)
			err := f.SetSheetRow(sheet, cell, &row)
			if err != nil {
				return err
			}
			rowIndex++
		}
	}

	err = f.SetSheetRow(sheet, "A1", &headers)
	if err != nil {
		return err
	}
	workDirectory, _ := config.GetWorkDirectory()
	exportsDirectory := filepath.Join(workDirectory, "exports", filepath.Dir(outputPath))
	makeExportDirectoryError := os.MkdirAll(exportsDirectory, 0777)
	if makeExportDirectoryError != nil {
		panic(fmt.Sprintf("Cannot Create %f", makeExportDirectoryError))

	}
	xlsxFilePath := filepath.Join(exportsDirectory, filepath.Base(outputPath))
	fmt.Printf("✅ Salvo em: %s\n", xlsxFilePath)

	if err := f.SaveAs(xlsxFilePath); err != nil {
		return fmt.Errorf("erro ao salvar planilha: %w", err)
	}
	consoleHook := bufio.NewReader(os.Stdin)
	fmt.Print("Aperte qualquer teclas para sair: ")
	_, _ = consoleHook.ReadString('\n')

	return nil
}
