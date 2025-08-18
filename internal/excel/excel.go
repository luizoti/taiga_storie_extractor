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

func ExportMergedComments(
	project structs.Project,
	stories []structs.Storie,
	storyComments map[int][]structs.StorieDetails,
	CustomFieldsMap map[int]map[string]string,
	outputPath string,
) error {
	f := excelize.NewFile()
	sheet := "Comentários"
	f.SetSheetName("Sheet1", sheet)

	// Cabeçalhos atualizados com todos os campos
	headers := []string{
		"Projeto ID",
		"Projeto Nome",
		"Projeto Slug",
		"Projeto Descrição",
		"Projeto Criado em",
		"Projeto Modificado em",
		//"Projeto CSV UUID",
		"História ID", "História Ref",
		"História Nome",
		"História Criado em",
		"História Modificado em",
		"História Concluída em",
		"História Data Limite",
		"História Motivo da Data Limite",
		"História Status da Data Limite",
		"História Comentário Inicial",
		"Comentário ID",
		"Comentário Criado em",
		"Comentário Texto",
		"Comentário HTML",
	}

	rowIndex := 2

	for _, story := range stories {

		comments := storyComments[story.ID]
		for _, comment := range comments {
			row := []string{
				strconv.Itoa(project.ID),
				project.Name,
				project.Slug,
				project.Description,
				project.CreatedDate,
				project.ModifiedDate,
				//project.CsvUUID,

				strconv.Itoa(story.ID),
				strconv.Itoa(story.Ref),
				story.Name,
				story.CreatedDate,
				story.ModifiedDate,
				story.FinishDate,
				story.DueDate,
				story.DueDateReason,
				story.DueDateStatus,
				story.Comment,

				comment.UuId,
				comment.CreatedDate,
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
			f.SetSheetRow(sheet, cell, &row)
			rowIndex++
		}
	}

	f.SetSheetRow(sheet, "A1", &headers)

	exportsDirectory := filepath.Join(config.GetWorkDirectory(), "exports", filepath.Dir(outputPath))
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
