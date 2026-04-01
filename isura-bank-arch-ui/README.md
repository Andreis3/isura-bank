# Isura Bank — System Design (UI)

Interface interativa de mapeamento de serviços do Isura Bank.

## Pré-requisitos

- **Node.js** 18+ (recomendado 20 LTS)
- **npm** (vem com o Node)

## Como rodar

```bash
# 1. Instalar dependências
npm install

# 2. Rodar em modo dev (abre no browser automaticamente)
npm run dev
```

O app abre em `http://localhost:3000`.

## Build para produção

```bash
npm run build
npm run preview   # preview do build
```

## Estrutura

```
isura-bank-ui/
├── index.html          # Entry HTML
├── package.json        # Dependências (Vite + React)
├── vite.config.js      # Config do Vite
├── README.md
└── src/
    ├── main.jsx        # Bootstrap React
    └── App.jsx         # App completo (tudo em 1 arquivo)
```
