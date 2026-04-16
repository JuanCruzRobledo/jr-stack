# JR Stack — OPSX Edition

**Fork del [Gentle-AI Stack](https://github.com/Gentleman-Programming/gentle-ai) original, completamente rebrandeado y actualizado para usar el workflow OPSX de OpenSpec.**

> Este repositorio reemplaza tanto el branding original como el flujo Legacy (SDD con fases rigidas) por **JR Stack** con el flujo **OPSX** (acciones fluidas e iterativas).

---

## Quick Start

### Paso 1: Instalar

#### Linux / macOS

```bash
curl -fsSL https://raw.githubusercontent.com/JuanCruzRobledo/jr-stack/main/scripts/install-opsx.sh | bash
```

#### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/JuanCruzRobledo/jr-stack/main/scripts/install-opsx.ps1 | iex
```

> **Requisito:** Go 1.24+ ([descargar](https://go.dev/dl/)) y git.
> El script clona el repo, compila el binario `jr-stack`, y crea toda la configuracion desde cero con OPSX.

### Installation Modes

The installer prompts you to choose a mode:

| Mode | What it installs | Best for |
|------|-----------------|----------|
| **Lite** (default) | OPSX orchestrator + Engram memory + Context7 docs | Quick start, students, first-time users |
| **Full** | Lite + Skills + GGA + Persona | Complete experience, power users |
| **Custom** | Binary only (configure via TUI) | Advanced users who want full control |

Press Enter for Lite (recommended), or choose 2/3 for other modes.

### Paso 2: Lanzar la interfaz

```bash
jr-stack
```

Esto abre el TUI interactivo donde podes:
- Seleccionar que agentes configurar (Claude Code, OpenCode, Cursor, etc.)
- Elegir persona y preset
- Configurar modelos por accion OPSX
- Activar Strict TDD mode
- Desinstalar la version original de Gentleman si la tenias

### Paso 3: Verificar

Reinicia tu agente de IA y preguntale:

```
Quien sos y que podes hacer? Explicame tu flujo de trabajo completo.
```

Deberia responder mencionando:
- Que es el **OPSX Orchestrator**
- Que trabaja con el CLI `openspec`
- Que el flujo es: **explore -> propose -> apply -> archive**
- Que los comandos son `/opsx:explore`, `/opsx:propose`, `/opsx:apply`, `/opsx:archive`

---

## Que cambio respecto al original?

| Aspecto | Gentleman Stack (Legacy) | JR Stack (OPSX) |
|---------|------------------------|-------------------|
| **Comando** | `gentle-ai` | `jr-stack` |
| **Config dir** | `~/.gentle-ai/` | `~/.jr-stack/` |
| **Flujo de trabajo** | Fases rigidas y bloqueantes | Fluido e iterativo: cualquier accion en cualquier momento |
| **Comandos** | `/sdd-apply`, `/sdd-archive`, etc. | `/opsx:explore`, `/opsx:propose`, `/opsx:apply`, `/opsx:archive` |
| **Orchestrator** | Coordina sub-agentes SDD con fases | Coordina via skills + CLI `openspec` como fuente de verdad |
| **Self-update** | Apunta a repo de Gentleman | Apunta a este repo (`JuanCruzRobledo/jr-stack`) |
| **Desinstalar original** | Manual | Opcion integrada en el TUI |

---

## Comandos CLI

```bash
jr-stack                     # Lanzar TUI interactivo
jr-stack install [flags]     # Instalar desde CLI (sin TUI)
jr-stack sync                # Sincronizar configs de agentes
jr-stack sync --lite         # Sync essentials only (OPSX + Engram + Context7)
jr-stack update              # Buscar actualizaciones
jr-stack upgrade             # Aplicar actualizaciones
jr-stack restore             # Restaurar un backup
jr-stack version             # Mostrar version
jr-stack help                # Ayuda
```

### Re-sincronizar configs

```bash
jr-stack sync
```

---

## Flujo OPSX

```
/opsx:explore   (opcional: pensar antes de comprometerse)
       |
       v
/opsx:propose   (crear cambio + propuesta + diseno + tareas)
       |
       v
/opsx:apply     (implementar las tareas del cambio)
       |
       v
/opsx:archive   (sincronizar specs + cerrar el cambio)
```

**No hay fases rigidas.** Podes volver atras, saltear pasos, o repetir cualquier accion en cualquier momento.

| Comando | Que hace |
|---------|----------|
| `/opsx:explore [tema]` | Modo exploracion: investigar ideas, aclarar requisitos, pensar. No genera archivos. |
| `/opsx:propose [nombre]` | Crea un cambio con todos los artefactos: `proposal.md`, `design.md`, `tasks.md` |
| `/opsx:apply [nombre]` | Implementa las tareas del cambio, marcandolas como completadas |
| `/opsx:archive [nombre]` | Sincroniza delta specs con los specs principales y archiva el cambio |

### Estructura de un cambio

```
openspec/changes/<nombre-del-cambio>/
  .openspec.yaml    <- metadata del cambio
  proposal.md       <- que y por que
  design.md         <- como (enfoque tecnico)
  tasks.md          <- checklist de implementacion
  specs/            <- delta specs (requisitos que cambian)
```

---

## Migrando desde Gentleman (gentle-ai)

Si ya tenias instalado el stack original:

1. Ejecuta `jr-stack`
2. Selecciona **"Uninstall Gentleman (gentle-ai)"** del menu
3. Esto limpia: binario `gentle-ai`, config dir `~/.gentle-ai/`, env vars del shell profile
4. Run `jr-stack sync` to apply full config, or `jr-stack sync --lite` for essentials only

---

## Troubleshooting

### El agente sigue respondiendo con el flujo Legacy (SDD)

**Causa:** Quedan archivos del stack original.

**Solucion:**
1. Usa la opcion "Uninstall Gentleman" del TUI
2. Corre `jr-stack sync` para reinyectar configs OPSX

### OpenCode no reconoce los comandos `/opsx:*`

- Reinicia OpenCode despues del sync
- Verifica que `~/.config/opencode/commands/` tenga los archivos `opsx-*.md`

### El build falla con errores de Go

- Verifica Go 1.24+ con `go version`
- Asegurate de estar en la raiz del repositorio (donde esta `go.mod`)

---

## Creditos

- **Stack original:** [Gentleman Programming — gentle-ai](https://github.com/Gentleman-Programming/gentle-ai)
- **OpenSpec / OPSX:** [Fission AI — OpenSpec](https://github.com/Fission-AI/OpenSpec)
- **JR Stack:** [JuanCruzRobledo](https://github.com/JuanCruzRobledo/jr-stack)

---

<div align="center">
<a href="LICENSE"><img src="https://img.shields.io/badge/License-MIT-blue.svg" alt="License: MIT"></a>
</div>
