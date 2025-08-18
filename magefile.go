//go:build mage
// +build mage

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"taiga_storie_extractor/internal/versioning"

	_ "github.com/magefile/mage/mg"
)

var (
	buildDir = "build"
	cmdName  = versioning.CmdName
	version  = versioning.Version
)

// All executa o build para o sistema atual e gera zip correspondente
func All() error {
	Clean()
	UpdateDeps()

	fmt.Printf("🌐 Detectado sistema operacional: %s\n", runtime.GOOS)

	switch runtime.GOOS {
	case "windows":
		if err := BuildWindows(); err != nil {
			return err
		}
		if err := Zip("windows"); err != nil {
			return err
		}
	case "linux":
		if err := BuildLinux(); err != nil {
			return err
		}
		if err := Zip("linux"); err != nil {
			return err
		}
	default:
		return fmt.Errorf("Sistema operacional '%s' não suportado", runtime.GOOS)
	}

	return nil
}

func UpdateDeps() error {
	fmt.Println("🔄 Atualizando dependências Go...")
	cmds := [][]string{
		{"go", "get", "-u", "./..."},
		{"go", "mod", "tidy"},
		{"go", "mod", "download"},
	}

	for _, cmd := range cmds {
		fmt.Println("Executando:", cmd)
		out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
		fmt.Println(string(out))
		if err != nil {
			return err
		}
	}

	return nil
}

// BuildLinux compila o projeto para Linux com CGO
func BuildLinux() error {
	return buildTarget("linux", "amd64")
}

// BuildWindows compila o projeto para Windows com CGO
func BuildWindows() error {
	return buildTarget("windows", "amd64")
}

func buildTarget(goos, goarch string) error {
	fmt.Printf("\n🔨 Compilando para %s/%s...\n", goos, goarch)

	if err := checkGCC(); err != nil {
		return err
	}

	dir := filepath.Join(buildDir, goos)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	binaryName := cmdName
	if goos == "windows" {
		binaryName += ".exe"
	}
	outputPath := filepath.Join(dir, binaryName)

	cmd := exec.Command("go", "build", "-buildvcs=false", "-ldflags", "-X main.version="+version, "-o", outputPath, "./cmd/taiga_storie_extractor")
	cmd.Env = append(os.Environ(),
		"GOOS="+goos,
		"GOARCH="+goarch,
		"CGO_ENABLED=1",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	for _, file := range []string{"CHANGELOG.MD", "config.json.exemple"} {
		dst := filepath.Join(dir, filepath.Base(file))
		if err := copyFile(file, dst); err != nil {
			return err
		}
	}
	_ = os.Rename(filepath.Join(dir, "config.json.exemple"), filepath.Join(dir, "config.json"))
	return nil
}

func Zip(target string) error {
	fmt.Println("📦 Gerando zip para:", target)
	dir := filepath.Join(buildDir, target)
	zipName := fmt.Sprintf("%s_%s_%s.zip", cmdName, version, target)

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		if _, err := exec.LookPath("powershell"); err != nil {
			return fmt.Errorf("powershell não encontrado. Requerido para compactar no Windows")
		}
		cmd = exec.Command("powershell", "-Command", fmt.Sprintf("Compress-Archive -Path %s\\* -DestinationPath %s", dir, zipName))
	} else {
		if _, err := exec.LookPath("zip"); err != nil {
			return fmt.Errorf("zip não encontrado. Instale com: sudo apt install zip")
		}
		cmd = exec.Command("zip", "-r", zipName, dir)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Clean remove tudo que foi gerado
func Clean() error {
	fmt.Println("🧹 Limpando arquivos...")
	_ = os.RemoveAll(buildDir)

	files, _ := os.ReadDir(".")
	for _, f := range files {
		if strings.HasPrefix(f.Name(), cmdName+"_") && strings.HasSuffix(f.Name(), ".zip") {
			_ = os.Remove(f.Name())
		}
	}
	return nil
}

// copyFile copia um arquivo de src para dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, input, 0644)
}

// checkGCC garante que o gcc esteja disponível para builds com CGO
func checkGCC() error {
	_, err := exec.LookPath("gcc")
	if err != nil {
		return fmt.Errorf("❌ GCC não encontrado no PATH. Verifique se está instalado e acessível")
	}
	return nil
}
