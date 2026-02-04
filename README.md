# Replacciato

Replaces Catpuccin themes with values of conforming themes if you don't want to manually change them yourself.

### Usage

```bash
npm install
# CHOOSE YOUR ARCH AND SYSTEM
npm run build:amd64
npm run build:windows:amd64
npm run build:macos:amd64
npm run build:arm64
npm run build:windows:arm64
npm run build:macos:arm64
# OR, TO JUST MATCH YOUR SYSTEM ARCH
go build -o replacciato main.go
```

```bash
./replacciato --path <file or dir> --theme <theme>.json --type <latte,frappe,macchiato,mocha>
```

### Currently Supports:

- Hex Codes
- RGB (rgb(r,g,b))
- RGB in arrays [r, g, b]
- HSL (hsl(hdeg,s%,l))