# 📤 **Taiga Storie Extractor**

O **Taiga Storie Extractor** Extrai cards e comentários de boards do taiga.

---

## ⚙️ **Arquivo de Configuração – `config.json`**

Este arquivo define os parâmetros de tempo e segurança para os envios.

A opção `extractAfter` não foi implementada.

### Exemplo de configuração:
```json
{
  "logLevel": "INFO",
  "username": "seuuser",
  "password": "suasenha",
  "extractAfter": "2025/08/06",
  "apiBaseUrl": "https://taiga.SEU_SERVER.com.br/api/v1"
}
```
---

## 🚀 **Como Usar**
        
### 1. **Configure o `config.json`**
Edite o arquivo com os valores necessários.
        
### 2. **Execute o programa**

#### No **Windows**:
- Dê **duplo clique** em `taigaStorieExtractor.exe`.

#### No **macOS/Linux**:
- Abra o terminal na pasta do programa e execute:
```bash
./taigaStorieExtractor
```
