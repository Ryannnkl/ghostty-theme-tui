# Ghostty Theme TUI - Implementation Plan

## Objetivo

Criar uma TUI em Go para navegar, pesquisar e aplicar temas do Ghostty com preview imediato.

A ferramenta deve resolver o problema atual: trocar tema manualmente no arquivo `~/.config/ghostty/config` deixa o fluxo lento e pouco visual. A TUI deve permitir testar temas rapidamente usando teclado e confirmar o tema escolhido.

## Experiencia desejada

Tela principal:

- Input de busca no topo.
- Lista filtrada de temas abaixo.
- Navegacao com setas do teclado.
- Ao mover a selecao, o preview do tema muda imediatamente.
- Ao pressionar `Enter`, o tema selecionado e salvo no config do Ghostty.
- Ao pressionar `Esc` ou `Ctrl+C`, a TUI sai sem salvar e restaura o tema original.
- Ao pressionar `/`, o foco volta para o input de busca.
- Ao pressionar `Ctrl+R`, recarrega a lista de temas.

Comportamento esperado:

- A lista deve carregar temas de `ghostty +list-themes --path`.
- A busca deve filtrar por texto, sem diferenciar maiusculas/minusculas.
- O tema atual deve aparecer destacado.
- O item selecionado deve mostrar alguma indicacao clara.
- Se houver erro ao ler temas ou escrever config, mostrar mensagem no rodape.

## Preview

Existem duas possibilidades.

### Opcao A - Preview dentro da propria TUI

Ler o arquivo do tema selecionado e aplicar as cores no terminal atual usando sequencias OSC:

- `OSC 10` para foreground.
- `OSC 11` para background.
- `OSC 12` para cursor.
- `OSC 4` para palette ANSI.

Vantagens:

- Preview instantaneo.
- Nao depende de recarregar janelas externas do Ghostty.
- Funciona dentro da propria TUI.

Cuidados:

- Ao sair sem salvar, restaurar o tema original.
- Ao trocar de tema muitas vezes, evitar flicker.
- Nem todos os temas podem declarar todas as cores.

### Opcao B - Preview alterando config do Ghostty

Alterar `theme = ...` temporariamente no config e depender de reload.

Problema:

- O Ghostty tem action `reload_config`, mas nao parece expor um comando CLI direto para disparar essa action em uma janela ja aberta.
- O usuario teria que recarregar manualmente com `Ctrl+Shift+,`, o que quebra a ideia de preview ao navegar.

Decisao recomendada:

Usar a Opcao A para preview na TUI e salvar no config apenas ao confirmar com `Enter`.

## Stack recomendada

Usar Bubble Tea:

- `github.com/charmbracelet/bubbletea` para loop da TUI.
- `github.com/charmbracelet/bubbles/textinput` para input de busca.
- `github.com/charmbracelet/lipgloss` para layout e estilos.

Motivo:

- E idiomatico para TUI em Go.
- Lida bem com teclado, resize, estados e renderizacao.
- Evita escrever manualmente todo o controle de terminal.

## Estrutura de projeto

```text
ghostty-theme-tui/
  go.mod
  README.md
  PLAN.md
  cmd/
    ghostty-theme-tui/
      main.go
  internal/
    ghostty/
      config.go
      themes.go
    preview/
      osc.go
    tui/
      model.go
      view.go
      update.go
```

## Modelo de dados

```go
type Theme struct {
    Name   string
    Source string // resources ou user
    Path   string
    Colors ThemeColors
}

type ThemeColors struct {
    Foreground string
    Background string
    Cursor     string
    Palette    map[int]string
}
```

## Leitura dos temas

Comando base:

```bash
ghostty +list-themes --path
```

Exemplo de linha:

```text
Catppuccin Mocha (resources) /usr/share/ghostty/themes/Catppuccin Mocha
```

Parsing:

- Nome: antes de ` (resources)` ou ` (user)`.
- Source: `resources` ou `user`.
- Path: caminho depois do source.

Importante:

- Nomes e caminhos podem conter espacos.
- Nao usar `strings.Fields` de forma ingenua.
- Usar regexp:

```regex
^(.+) \((resources|user)\) (.+)$
```

## Parsing dos arquivos de tema

Arquivos de tema do Ghostty usam sintaxe parecida com config:

```ini
background = #1e1e2e
foreground = #cdd6f4
cursor-color = #f5e0dc
palette = 0=#45475a
palette = 1=#f38ba8
```

Implementar parser simples linha a linha:

- Remover comentarios iniciados por `#` somente quando a linha inteira for comentario.
- Ignorar linhas vazias.
- Separar por primeiro `=`.
- Trim em chave e valor.
- Suportar:
  - `background`
  - `foreground`
  - `cursor-color`
  - `palette`

Para `palette`, separar `index=color`.

## Escrita no config do Ghostty

Arquivo alvo:

```text
${XDG_CONFIG_HOME:-$HOME/.config}/ghostty/config
```

Regra:

- Se existir linha `theme = ...`, substituir a primeira.
- Se nao existir, adicionar `theme = Nome Do Tema` no final.
- Preservar o resto do arquivo.
- Criar diretorio se nao existir.

Antes de escrever:

- Fazer backup simples:

```text
config.bak
```

Depois de escrever:

```bash
ghostty +validate-config
```

Se a validacao falhar:

- Restaurar backup.
- Mostrar erro no rodape da TUI.

## Estado da TUI

Campos principais do model:

```go
type model struct {
    themes        []ghostty.Theme
    filtered      []ghostty.Theme
    selected      int
    input         textinput.Model
    currentTheme  string
    originalTheme string
    message       string
    width         int
    height        int
}
```

Eventos:

- `tea.KeyUp`: decrementa selecao.
- `tea.KeyDown`: incrementa selecao.
- `tea.KeyEnter`: salva tema e sai.
- `tea.KeyEsc`: restaura preview original e sai.
- `ctrl+c`: restaura preview original e sai.
- Input text changed: recalcula filtro e seleciona primeiro item.

## Restauracao ao sair

Ao iniciar:

- Ler tema atual do config.
- Guardar como `originalTheme`.
- Aplicar preview do tema atual, se possivel.

Ao navegar:

- Aplicar preview do tema selecionado.

Ao sair sem confirmar:

- Aplicar preview do `originalTheme`.
- Nao alterar config.

Ao confirmar:

- Escrever config.
- Manter preview do tema selecionado.

## Comandos planejados

```bash
ghostty-theme-tui
```

Opcionalmente:

```bash
ghostty-theme-tui --color dark
ghostty-theme-tui --color light
ghostty-theme-tui --color all
```

`--color` pode filtrar usando:

```bash
ghostty +list-themes --color=dark --path
```

## Instalacao local

Durante desenvolvimento:

```bash
go run ./cmd/ghostty-theme-tui
```

Instalar no usuario:

```bash
go install ./cmd/ghostty-theme-tui
```

Se `~/go/bin` nao estiver no PATH, adicionar no `~/.bashrc`:

```bash
export PATH="$HOME/go/bin:$PATH"
```

## Testes recomendados

Unidade:

- Parser de `ghostty +list-themes --path`.
- Parser de arquivo de tema.
- Substituicao de `theme = ...` no config.
- Filtro por texto.

Manuais:

- Abrir TUI dentro do Ghostty.
- Digitar busca.
- Navegar com setas.
- Confirmar com Enter.
- Sair com Esc e verificar se nao salvou.
- Testar tema com espacos no nome.
- Testar config sem linha `theme =`.

## Primeiras tarefas

1. Criar `go.mod`.
2. Adicionar dependencias Bubble Tea, Bubbles e Lip Gloss.
3. Implementar `internal/ghostty/themes.go`.
4. Implementar `internal/ghostty/config.go`.
5. Implementar `internal/preview/osc.go`.
6. Implementar TUI basica com input e lista.
7. Conectar preview ao movimento da selecao.
8. Conectar `Enter` para salvar.
9. Adicionar README com uso.

## Riscos conhecidos

- Sequencias OSC podem nao restaurar exatamente tudo se o tema tiver poucas cores declaradas.
- Alguns temas podem usar includes ou opcoes mais complexas; comecar suportando os campos principais.
- O preview afeta o terminal atual enquanto a TUI roda; por isso a restauracao no cancelamento e importante.
- Se o programa crashar, o preview visual pode ficar aplicado ate reiniciar/recarregar o terminal.

## Resultado esperado

Ao final, o fluxo ideal deve ser:

```bash
ghostty-theme-tui
```

Depois:

- Digitar parte do nome do tema.
- Navegar com setas.
- Ver o preview mudar na hora.
- Pressionar `Enter`.
- O tema fica salvo em `~/.config/ghostty/config`.
