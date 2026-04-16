# TODO: Infraestructura de releases (binarios pre-compilados)

## Estado actual

Los scripts de instalacion (`install-opsx.sh` / `install-opsx.ps1`) requieren
que el usuario tenga **Go 1.24+** instalado porque compilan desde el codigo fuente.

El repo original ([gentle-ai](https://github.com/Gentleman-Programming/gentle-ai))
no requiere Go porque distribuye **binarios pre-compilados** via GitHub Releases.

## Que falta para eliminar el requisito de Go

### 1. Configurar GoReleaser

Crear `.goreleaser.yaml` en la raiz del proyecto. GoReleaser compila el binario
para todas las plataformas automaticamente en cada release tag.

Referencia: el repo original ya tiene uno que se puede adaptar.

Plataformas a soportar:
- `linux/amd64`, `linux/arm64`
- `darwin/amd64`, `darwin/arm64` (macOS Intel + Apple Silicon)
- `windows/amd64`, `windows/arm64`

### 2. Configurar GitHub Actions

Crear `.github/workflows/release.yml` que:
- Se dispare con cada tag `v*` (ej: `v1.0.0`)
- Ejecute GoReleaser para compilar y subir los binarios a GitHub Releases
- Genere `checksums.txt` para verificacion

Ejemplo minimo:

```yaml
name: Release
on:
  push:
    tags: ['v*']
jobs:
  release:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0
      - uses: actions/setup-go@v5
        with:
          go-version: '1.24'
      - uses: goreleaser/goreleaser-action@v6
        with:
          args: release --clean
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### 3. Actualizar los scripts de instalacion

Una vez que existan releases con binarios:
- `install-opsx.sh`: agregar metodo `binary` que descargue el binario directo
  (sin necesidad de Go). Mantener `go build` como fallback.
- `install-opsx.ps1`: idem para Windows.
- Prioridad: binary > go (igual que el repo original).

### 4. (Opcional) Homebrew tap y Scoop bucket

Para instalacion con gestores de paquetes:
- Crear repo `JuanCruzRobledo/homebrew-tap` con la formula de Homebrew
- Crear repo `JuanCruzRobledo/scoop-bucket` con el manifest de Scoop
- GoReleaser puede generar estas formulas automaticamente

## Beneficio

Los alumnos podrian instalar con un solo comando sin necesidad de Go:

```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/JuanCruzRobledo/jr-stack/main/scripts/install-opsx.sh | bash

# Windows
irm https://raw.githubusercontent.com/JuanCruzRobledo/jr-stack/main/scripts/install-opsx.ps1 | iex
```

Los mismos comandos que hoy, pero sin el pre-requisito de Go.
