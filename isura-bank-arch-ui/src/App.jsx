import { useState, useEffect } from "react";

// ═══════════════════════════════════════════════
// ISURA BANK — System Design Completo
// "Isura" = Tesouro em Yorubá
// ═══════════════════════════════════════════════

// ─── Reset CSS inline ───
const globalStyles = `
  *, *::before, *::after { box-sizing: border-box; margin: 0; padding: 0; }
  html { -webkit-font-smoothing: antialiased; -moz-osx-font-smoothing: grayscale; }
  body { margin: 0; background: #07080C; overflow-x: hidden; }
  ::-webkit-scrollbar { width: 6px; }
  ::-webkit-scrollbar-track { background: #111318; }
  ::-webkit-scrollbar-thumb { background: #23262F; border-radius: 3px; }
  ::-webkit-scrollbar-thumb:hover { background: #334155; }
  @import url('https://fonts.googleapis.com/css2?family=Playfair+Display:wght@400;600;700&family=Source+Sans+3:wght@400;500;600;700&family=JetBrains+Mono:wght@400;600;700&display=swap');
`;

const C = {
  bg: "#07080C",
  surface: "#111318",
  surfaceHover: "#181A20",
  surfaceRaised: "#1C1E26",
  border: "#23262F",
  gold: "#C9A84C",
  goldLight: "#E4C96A",
  goldDim: "#8B7632",
  goldGlow: "rgba(201,168,76,0.08)",
  ember: "#D4845A",
  terra: "#A67C52",
  sage: "#6B9E7A",
  sky: "#5B8FA8",
  coral: "#C75F5F",
  violet: "#8B7BB8",
  text: "#E8E4DC",
  textSoft: "#B8B0A2",
  textMuted: "#7A7468",
  textFaint: "#4A453E",
};

const ff = {
  display: "'Playfair Display', Georgia, serif",
  body: "'Source Sans 3', 'Segoe UI', sans-serif",
  mono: "'JetBrains Mono', 'Fira Code', monospace",
};

// ─── Service Data ───

const SERVICE_CATEGORIES = [
  {
    id: "core", name: "Core Banking", color: C.gold,
    description: "Serviços fundamentais que formam a espinha dorsal do banco. Sem eles, nada funciona.",
    services: [
      {
        id: "account-svc", name: "Account Service", shortName: "account", priority: "P0",
        description: "Gestão do ciclo de vida completo da conta corrente — abertura, encerramento, bloqueio, dados cadastrais.",
        responsibilities: ["Abertura de conta (onboarding flow)", "Ciclo de vida: PENDING → ACTIVE → BLOCKED → CLOSED", "Vinculação com customer (1 customer : N accounts)", "Gestão de limites (Pix diário, TED, saque)", "Agência + número + dígito verificador"],
        ownsData: ["accounts", "account_limits", "account_status_history"],
        dependencies: ["identity-svc (KYC aprovado para ativar)"],
        exposesAPI: ["gRPC (interno)", "REST via gateway (externo)"],
        events: ["account.created", "account.activated", "account.blocked", "account.closed"],
      },
      {
        id: "identity-svc", name: "Identity Service", shortName: "identity", priority: "P0",
        description: "Onboarding do cliente (PF/PJ), KYC, validação de documentos e compliance AML.",
        responsibilities: ["Cadastro de cliente (PF e PJ)", "Orquestração de KYC (documento + selfie + facematch)", "Validação CPF/CNPJ via Serpro", "Checagem PEP, sanções, listas restritivas (AML)", "Background check contínuo (re-KYC)", "Gestão de consentimentos LGPD"],
        ownsData: ["customers", "kyc_results", "documents", "consents_lgpd"],
        dependencies: ["Provedor KYC externo (idwall/similar)", "Serpro (CPF)"],
        exposesAPI: ["gRPC (interno)", "REST (onboarding mobile)"],
        events: ["customer.created", "kyc.approved", "kyc.rejected", "kyc.review_needed"],
      },
      {
        id: "ledger-svc", name: "Ledger Service", shortName: "ledger", priority: "P0",
        description: "Motor contábil double-entry. TODA movimentação financeira passa por aqui. Fonte da verdade de saldos.",
        responsibilities: ["Lançamentos double-entry (débito + crédito = 0)", "Saldo em tempo real (balance cache + entries)", "Plano de contas contábil (chart of accounts)", "Idempotência de transações", "Hold/Release (reservas para cartão e Pix)", "Conciliação com parceiros (D+1)", "Multimoeda (BRL principal, USD/EUR para cartão internacional)"],
        ownsData: ["chart_of_accounts", "ledger_accounts", "transactions", "entries"],
        dependencies: ["Nenhuma — é folha na árvore de dependência"],
        exposesAPI: ["gRPC (APENAS interno — nenhum client externo acessa direto)"],
        events: ["ledger.entry.created", "ledger.balance.updated", "ledger.hold.placed", "ledger.hold.released"],
      },
      {
        id: "auth-svc", name: "Auth Service", shortName: "auth", priority: "P0",
        description: "Autenticação e autorização. JWT, MFA, device fingerprint, sessões.",
        responsibilities: ["Login (email/CPF + senha)", "MFA (TOTP, SMS, biometria)", "Emissão de JWT (access + refresh tokens)", "Device management e fingerprinting", "Revogação de sessões", "Rate limiting de tentativas", "PIN transacional (6 dígitos para operações)"],
        ownsData: ["credentials", "sessions", "devices", "mfa_configs"],
        dependencies: ["Redis (sessões e rate limiting)"],
        exposesAPI: ["REST (login flow)", "gRPC (validação interna de token)"],
        events: ["auth.login.success", "auth.login.failed", "auth.mfa.required", "auth.device.new"],
      },
    ],
  },
  {
    id: "payments", name: "Pagamentos", color: C.sage,
    description: "Serviços que movimentam dinheiro — Pix, transferências, boletos.",
    services: [
      {
        id: "pix-svc", name: "Pix Service", shortName: "pix", priority: "P0",
        description: "Integração completa com SPI e DICT do BACEN. Envio, recebimento, chaves Pix, QR Code, devolução.",
        responsibilities: ["Registro/remoção de chaves no DICT", "Envio de Pix (débito → SPI → confirmação)", "Recebimento de Pix (SPI → crédito no ledger)", "QR Code estático e dinâmico", "Devolução (total e parcial)", "Portabilidade e reivindicação de chaves", "Saga: débito → envio SPI → compensação em falha"],
        ownsData: ["pix_keys", "pix_transfers", "pix_qrcodes", "pix_devolutions"],
        dependencies: ["ledger-svc", "antifraud-svc", "BACEN SPI", "BACEN DICT"],
        exposesAPI: ["REST (app)", "gRPC (interno)", "Webhook (SPI callback)"],
        events: ["pix.sent", "pix.received", "pix.key.registered", "pix.devolution.created"],
      },
      {
        id: "transfer-svc", name: "Transfer Service", shortName: "transfer", priority: "P1",
        description: "TED, TEF (interna) e agendamento de transferências.",
        responsibilities: ["TEF — transferência interna (instantânea)", "TED — transferência interbancária (CIP/STR)", "Agendamento de transferências", "Validação de limites diários", "Comprovantes"],
        ownsData: ["transfers", "transfer_schedules"],
        dependencies: ["ledger-svc", "account-svc (limites)", "CIP (TED)"],
        exposesAPI: ["REST (app)", "gRPC (interno)"],
        events: ["transfer.completed", "transfer.scheduled", "transfer.failed"],
      },
      {
        id: "boleto-svc", name: "Boleto Service", shortName: "boleto", priority: "P1",
        description: "Emissão e liquidação de boletos registrados via CIP.",
        responsibilities: ["Emissão de boleto (código de barras + linha digitável)", "Registro na CIP (CNAB 240 ou API)", "Recebimento de retorno de liquidação", "Pagamento de boletos externos", "Cálculo de multa/juros pós-vencimento", "Geração de PDF do boleto"],
        ownsData: ["boletos", "boleto_payments", "boleto_returns"],
        dependencies: ["ledger-svc", "CIP", "notif-svc"],
        exposesAPI: ["REST (app)", "Webhook (retorno CIP)"],
        events: ["boleto.issued", "boleto.paid", "boleto.expired"],
      },
    ],
  },
  {
    id: "cards", name: "Cartões", color: C.violet,
    description: "Emissão e processamento de cartões de débito e crédito com suporte internacional.",
    services: [
      {
        id: "card-svc", name: "Card Service", shortName: "card", priority: "P0",
        description: "Ciclo de vida do cartão — emissão, bloqueio, cancelamento, virtual. NÃO processa autorizações.",
        responsibilities: ["Emissão de cartão físico e virtual", "Bloqueio/desbloqueio temporário", "Cancelamento definitivo", "Configuração (limite, categorias, contactless)", "Tokenização para wallets (Apple Pay, Google Pay)"],
        ownsData: ["cards", "card_configs", "card_tokens"],
        dependencies: ["account-svc", "Processadora (Dock/Swap)"],
        exposesAPI: ["REST (app)", "gRPC (interno)"],
        events: ["card.issued", "card.activated", "card.blocked", "card.cancelled"],
      },
      {
        id: "card-auth-svc", name: "Card Authorization", shortName: "card-auth", priority: "P0",
        description: "Processa autorizações de compra em TEMPO REAL (<100ms). Decide aprovar ou recusar.",
        responsibilities: ["Receber autorização (ISO 8583 via processadora)", "Validar status do cartão", "Verificar saldo (débito) ou limite (crédito)", "Consultar antifraude real-time", "Converter moeda (internacional → BRL + IOF)", "Hold no ledger", "Responder APPROVED/DECLINED em < 100ms"],
        ownsData: ["authorizations", "auth_decisions"],
        dependencies: ["card-svc", "ledger-svc", "antifraud-svc", "exchange-svc"],
        exposesAPI: ["gRPC (processadora → nós)"],
        events: ["card.auth.approved", "card.auth.declined", "card.auth.reversed"],
      },
      {
        id: "card-billing-svc", name: "Card Billing", shortName: "card-billing", priority: "P1",
        description: "Fatura do cartão de crédito — fechamento mensal, pagamento, parcelamento.",
        responsibilities: ["Fechamento mensal de fatura", "Cálculo: compras + juros + IOF + anuidade", "Pagamento total, mínimo ou parcial", "Parcelamento de fatura", "Geração de PDF da fatura"],
        ownsData: ["invoices", "invoice_items", "invoice_payments"],
        dependencies: ["card-auth-svc", "ledger-svc", "credit-svc (juros)"],
        exposesAPI: ["REST (app)", "Cron (fechamento mensal)"],
        events: ["invoice.closed", "invoice.paid", "invoice.overdue"],
      },
    ],
  },
  {
    id: "credit", name: "Crédito", color: C.ember,
    description: "Motor de crédito — scoring, empréstimos, contratos e cobrança.",
    services: [
      {
        id: "credit-engine-svc", name: "Credit Engine", shortName: "credit-engine", priority: "P1",
        description: "Motor de decisão de crédito. Score, limites, análise de risco. Não empresta — só decide.",
        responsibilities: ["Score interno (0-1000)", "Consulta bureau externo (Serasa/SPC)", "Política de crédito (regras + ML)", "Definição de limite do cartão", "Pré-aprovação de empréstimos", "Re-scoring periódico"],
        ownsData: ["credit_scores", "credit_policies", "bureau_results", "pre_approvals"],
        dependencies: ["account-svc", "ledger-svc (histórico)", "Bureau externo"],
        exposesAPI: ["gRPC (interno)"],
        events: ["credit.score.calculated", "credit.preapproval.granted", "credit.limit.updated"],
      },
      {
        id: "loan-svc", name: "Loan Service", shortName: "loan", priority: "P1",
        description: "Gestão de empréstimos — simulação, contratação, parcelas, cobrança.",
        responsibilities: ["Simulação (Price, SAC)", "Contratação com contrato digital", "Desembolso (crédito via ledger)", "Geração de parcelas", "Antecipação de parcelas", "Régua de cobrança", "CET, IOF, TAC — cálculos regulatórios"],
        ownsData: ["loans", "loan_installments", "loan_contracts"],
        dependencies: ["credit-engine-svc", "ledger-svc", "boleto-svc", "notif-svc"],
        exposesAPI: ["REST (app)", "Cron (vencimentos)"],
        events: ["loan.contracted", "loan.disbursed", "loan.installment.paid", "loan.overdue", "loan.settled"],
      },
    ],
  },
  {
    id: "invest", name: "Investimentos", color: C.sky,
    description: "Plataforma de investimentos — renda fixa, fundos, custódia.",
    services: [
      {
        id: "invest-svc", name: "Investment Service", shortName: "invest", priority: "P2",
        description: "Aplicação e resgate em produtos de investimento. Integra com B3/CETIP.",
        responsibilities: ["Catálogo de produtos (CDB, LCI, LCA, Tesouro, Fundos)", "Simulação de rendimento", "Aplicação e resgate", "Marcação a mercado diária", "Cálculo de IR e IOF", "Come-cotas semestral", "Informe de rendimentos anual"],
        ownsData: ["products", "positions", "yields", "tax_records"],
        dependencies: ["ledger-svc", "Custodiante (B3/CETIP)", "account-svc"],
        exposesAPI: ["REST (app)"],
        events: ["invest.applied", "invest.redeemed", "invest.yield.calculated"],
      },
    ],
  },
  {
    id: "platform", name: "Plataforma", color: C.terra,
    description: "Serviços transversais consumidos por todos os domínios.",
    services: [
      {
        id: "api-gateway", name: "API Gateway", shortName: "gateway", priority: "P0",
        description: "Ponto único de entrada externo. Roteia, autentica, limita.",
        responsibilities: ["Roteamento para serviço correto", "Validação de JWT", "Rate limiting (Redis)", "REST → gRPC transformation", "CORS, security headers", "Circuit breaker", "API versioning", "Request logging + tracing"],
        ownsData: ["Nenhum — stateless"],
        dependencies: ["auth-svc", "Redis", "Todos os serviços"],
        exposesAPI: ["REST público (HTTPS)"],
        events: [],
      },
      {
        id: "antifraud-svc", name: "Antifraud Service", shortName: "antifraud", priority: "P0",
        description: "Análise de fraude em tempo real para todas as operações financeiras.",
        responsibilities: ["Análise real-time (< 50ms)", "Rules engine configurável", "Score de risco (0-100) por transação", "Device fingerprinting", "Block/allow lists", "Análise comportamental", "Feedback loop", "Integração com provedores externos"],
        ownsData: ["fraud_rules", "fraud_scores", "blocklists", "device_fingerprints"],
        dependencies: ["Redis (feature store)", "account-svc"],
        exposesAPI: ["gRPC (APENAS interno)"],
        events: ["fraud.alert.created", "fraud.transaction.blocked"],
      },
      {
        id: "exchange-svc", name: "Exchange Rate Service", shortName: "exchange", priority: "P1",
        description: "Cotações de moedas para compras internacionais no cartão.",
        responsibilities: ["Ingestão de cotações (BACEN + provedores)", "Conversão USD/EUR/etc → BRL", "Spread configurável", "Cache (Redis, TTL curto)", "Histórico de cotações"],
        ownsData: ["exchange_rates", "exchange_spreads"],
        dependencies: ["Feed externo", "Redis (cache)"],
        exposesAPI: ["gRPC (interno)"],
        events: ["exchange.rate.updated"],
      },
      {
        id: "notif-svc", name: "Notification Service", shortName: "notification", priority: "P1",
        description: "Central de notificações — push, SMS, email, in-app.",
        responsibilities: ["Push (Firebase/APNs)", "SMS (Twilio/Zenvia)", "Email transacional (SES/SendGrid)", "In-app (WebSocket/SSE)", "Templates por evento", "Preferências opt-in/opt-out", "Retry com backoff", "Audit log"],
        ownsData: ["notification_templates", "notification_log", "user_preferences"],
        dependencies: ["Kafka/Pulsar (consome eventos)", "Provedores externos"],
        exposesAPI: ["WebSocket (real-time)", "REST (preferências)"],
        events: ["notification.sent", "notification.failed"],
      },
      {
        id: "reconciliation-svc", name: "Reconciliation Service", shortName: "reconciliation", priority: "P1",
        description: "Conciliação automática entre ledger interno e parceiros externos.",
        responsibilities: ["Conciliação D+1 com SPI", "Conciliação com processadora de cartão", "Conciliação com CIP (boletos)", "Detecção de divergências", "Alertas para time financeiro", "Relatórios regulatórios (BACEN)"],
        ownsData: ["reconciliation_runs", "discrepancies"],
        dependencies: ["ledger-svc", "Arquivos de retorno dos parceiros"],
        exposesAPI: ["REST (dashboard interno)", "Cron (batch diário)"],
        events: ["recon.completed", "recon.discrepancy.found"],
      },
      {
        id: "audit-svc", name: "Audit Service", shortName: "audit", priority: "P1",
        description: "Log imutável de todas as ações. Compliance BACEN, LGPD, auditoria.",
        responsibilities: ["Ingestão de eventos de todos os serviços", "Armazenamento imutável (append-only)", "Busca por período, serviço, usuário", "Retenção 5+ anos", "Exportação para auditoria externa", "Integração com SIEM"],
        ownsData: ["audit_events (append-only, particionado)"],
        dependencies: ["Kafka (consome tudo)", "S3 (cold storage)"],
        exposesAPI: ["REST (consulta interna)", "gRPC (ingestão)"],
        events: [],
      },
    ],
  },
];

const ALL_SERVICES = SERVICE_CATEGORIES.flatMap((c) =>
  c.services.map((s) => ({ ...s, categoryColor: c.color, categoryName: c.name }))
);

const PRIORITIES = {
  P0: { label: "P0 — Dia 1", color: C.coral, desc: "Sem isso o banco não abre" },
  P1: { label: "P1 — MVP", color: C.gold, desc: "Necessário para operar" },
  P2: { label: "P2 — Fase 2", color: C.sky, desc: "Pós-lançamento" },
};

// ─── Components ───

function ServiceCard({ service, isExpanded, onToggle }) {
  const pri = PRIORITIES[service.priority];
  return (
    <div
      onClick={onToggle}
      style={{
        background: isExpanded ? C.surfaceRaised : C.surface,
        border: `1px solid ${isExpanded ? service.categoryColor + "55" : C.border}`,
        borderRadius: "10px",
        padding: "18px 20px",
        cursor: "pointer",
        transition: "all 0.2s ease",
        marginBottom: "8px",
      }}
    >
      <div style={{ display: "flex", justifyContent: "space-between", alignItems: "flex-start" }}>
        <div style={{ flex: 1 }}>
          <div style={{ display: "flex", alignItems: "center", gap: "10px", marginBottom: "6px", flexWrap: "wrap" }}>
            <span style={{ fontFamily: ff.mono, fontSize: "14px", fontWeight: 700, color: service.categoryColor }}>
              {service.name}
            </span>
            <span style={{
              padding: "2px 8px", borderRadius: "4px", fontSize: "10px", fontFamily: ff.mono,
              fontWeight: 600, background: pri.color + "18", color: pri.color,
              border: `1px solid ${pri.color}33`, letterSpacing: "0.5px",
            }}>
              {service.priority}
            </span>
          </div>
          <p style={{ fontFamily: ff.body, fontSize: "13px", color: C.textSoft, lineHeight: 1.5, margin: 0 }}>
            {service.description}
          </p>
        </div>
        <span style={{
          color: C.textMuted, fontSize: "18px", marginLeft: "12px",
          transform: isExpanded ? "rotate(45deg)" : "rotate(0deg)",
          transition: "transform 0.2s", flexShrink: 0,
        }}>
          +
        </span>
      </div>

      {isExpanded && (
        <div style={{ marginTop: "18px", borderTop: `1px solid ${C.border}`, paddingTop: "16px" }}
          onClick={(e) => e.stopPropagation()}>
          <DetailSection title="Responsabilidades" items={service.responsibilities} color={service.categoryColor} />
          <DetailSection title="Dados que possui" items={service.ownsData} color={C.gold} />
          <DetailSection title="Dependências" items={service.dependencies} color={C.coral} />
          <DetailSection title="APIs expostas" items={service.exposesAPI} color={C.sage} />
          {service.events.length > 0 && (
            <div style={{ marginTop: "14px" }}>
              <div style={{ fontFamily: ff.mono, fontSize: "11px", color: C.textMuted, marginBottom: "8px", letterSpacing: "0.5px", textTransform: "uppercase" }}>
                Domain Events (Kafka)
              </div>
              <div style={{ display: "flex", flexWrap: "wrap", gap: "6px" }}>
                {service.events.map((e, i) => (
                  <code key={i} style={{
                    fontFamily: ff.mono, fontSize: "11px", padding: "3px 8px", borderRadius: "4px",
                    background: C.ember + "15", color: C.ember, border: `1px solid ${C.ember}25`,
                  }}>
                    {e}
                  </code>
                ))}
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}

function DetailSection({ title, items, color }) {
  return (
    <div style={{ marginTop: "14px" }}>
      <div style={{ fontFamily: ff.mono, fontSize: "11px", color: C.textMuted, marginBottom: "8px", letterSpacing: "0.5px", textTransform: "uppercase" }}>
        {title}
      </div>
      {items.map((item, i) => (
        <div key={i} style={{ display: "flex", alignItems: "flex-start", gap: "8px", padding: "3px 0", fontFamily: ff.body, fontSize: "13px", color: C.textSoft, lineHeight: 1.5 }}>
          <span style={{ color: color, fontSize: "8px", marginTop: "6px", flexShrink: 0 }}>●</span>
          {item}
        </div>
      ))}
    </div>
  );
}

// ─── Dependency Map ───

function DependencyMap() {
  return (
    <div>
      <p style={{ fontFamily: ff.body, fontSize: "14px", color: C.textSoft, lineHeight: 1.7, margin: "0 0 20px 0" }}>
        Grafo de dependências entre serviços. O <strong style={{ color: C.gold }}>ledger-svc</strong> é o
        serviço mais demandado — todos que movimentam dinheiro dependem dele.
      </p>
      <div style={{ overflowX: "auto" }}>
        <pre style={{
          fontFamily: ff.mono, fontSize: "12px", lineHeight: 1.8, color: C.textSoft,
          background: C.surface, border: `1px solid ${C.border}`, borderRadius: "10px", padding: "20px", whiteSpace: "pre",
        }}>
{`                    ┌─────────────┐
                    │ API Gateway │  ← ponto de entrada
                    └──────┬──────┘
                           │ REST → gRPC
        ┌──────────┬───────┼───────┬──────────┬──────────┐
        ▼          ▼       ▼       ▼          ▼          ▼
  ┌──────────┐ ┌──────┐ ┌─────┐ ┌──────┐ ┌───────┐ ┌────────┐
  │ account  │ │ pix  │ │card │ │ loan │ │invest │ │ boleto │
  │   svc    │ │ svc  │ │ svc │ │ svc  │ │  svc  │ │  svc   │
  └────┬─────┘ └──┬───┘ └──┬──┘ └──┬───┘ └──┬────┘ └───┬────┘
       │          │        │       │         │          │
       │          ▼        ▼       │         │          │
       │     ┌─────────┐ ┌────────┴──┐      │          │
       │     │antifraud│ │ card-auth │      │          │
       │     │   svc   │ │    svc    │      │          │
       │     └────┬────┘ └─────┬─────┘      │          │
       │          │            │             │          │
       │    ┌─────┘     ┌──────┘             │          │
       │    │           │    ┌───────────────┘          │
       │    │           │    │        ┌─────────────────┘
       ▼    ▼           ▼    ▼        ▼
  ╔═══════════════════════════════════════╗
  ║         LEDGER SERVICE               ║
  ║   (double-entry, saldos, holds)      ║
  ╚═══════════════════════════════════════╝
       │              │              │
       ▼              ▼              ▼
  ┌──────────┐  ┌──────────┐  ┌──────────┐
  │PostgreSQL│  │  Redis   │  │  Kafka   │
  │(per svc) │  │ (cache)  │  │ (events) │
  └──────────┘  └──────────┘  └──────────┘

  Transversais (consomem eventos de todos):
  ┌──────────┐ ┌───────────────┐ ┌──────────┐
  │  notif   │ │reconciliation │ │  audit   │
  │   svc    │ │     svc       │ │   svc    │
  └──────────┘ └───────────────┘ └──────────┘`}
        </pre>
      </div>

      <div style={{ marginTop: "24px" }}>
        <div style={{ fontFamily: ff.mono, fontSize: "11px", color: C.textMuted, marginBottom: "12px", letterSpacing: "0.5px", textTransform: "uppercase" }}>
          Regras de comunicação
        </div>
        {[
          { rule: "Síncrono (gRPC)", desc: "Apenas quando o caller PRECISA da resposta para continuar. Ex: pix-svc → ledger-svc (débito).", color: C.sage },
          { rule: "Assíncrono (Kafka)", desc: "Para tudo eventual. Ex: ledger → notification (avisar cliente do Pix recebido).", color: C.ember },
          { rule: "Database per Service", desc: "Cada serviço tem seu schema. NENHUM serviço lê o banco de outro.", color: C.gold },
          { rule: "Sem ciclo", desc: "Se A depende de B sincronamente, B NUNCA depende de A. Usa eventos.", color: C.coral },
        ].map((r, i) => (
          <div key={i} style={{ display: "flex", gap: "12px", marginBottom: "12px", alignItems: "flex-start" }}>
            <code style={{ fontFamily: ff.mono, fontSize: "11px", color: r.color, whiteSpace: "nowrap", padding: "3px 8px", background: r.color + "12", borderRadius: "4px", border: `1px solid ${r.color}25` }}>
              {r.rule}
            </code>
            <span style={{ fontFamily: ff.body, fontSize: "13px", color: C.textSoft, lineHeight: 1.5 }}>{r.desc}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

// ─── Summary ───

function SummaryView() {
  const groups = [
    { label: "P0 — Dia 1 (obrigatório)", services: ALL_SERVICES.filter(s => s.priority === "P0"), color: C.coral },
    { label: "P1 — MVP (operação mínima)", services: ALL_SERVICES.filter(s => s.priority === "P1"), color: C.gold },
    { label: "P2 — Fase 2 (pós-lançamento)", services: ALL_SERVICES.filter(s => s.priority === "P2"), color: C.sky },
  ];

  return (
    <div>
      <p style={{ fontFamily: ff.body, fontSize: "14px", color: C.textSoft, lineHeight: 1.7, margin: "0 0 24px 0" }}>
        O Isura Bank precisa de <strong style={{ color: C.gold }}>{ALL_SERVICES.length} serviços</strong> em {SERVICE_CATEGORIES.length} domínios.
      </p>

      {groups.map((group, gi) => (
        <div key={gi} style={{ marginBottom: "24px" }}>
          <div style={{ display: "flex", alignItems: "center", gap: "10px", marginBottom: "12px" }}>
            <span style={{ fontFamily: ff.mono, fontSize: "12px", fontWeight: 700, color: group.color }}>{group.label}</span>
            <span style={{ fontFamily: ff.mono, fontSize: "11px", color: C.textMuted }}>({group.services.length})</span>
          </div>
          <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))", gap: "8px" }}>
            {group.services.map((s, i) => (
              <div key={i} style={{
                background: C.surface, border: `1px solid ${C.border}`, borderLeft: `3px solid ${s.categoryColor}`,
                borderRadius: "6px", padding: "12px 14px",
              }}>
                <div style={{ fontFamily: ff.mono, fontSize: "12px", fontWeight: 600, color: s.categoryColor, marginBottom: "4px" }}>{s.shortName}</div>
                <div style={{ fontFamily: ff.body, fontSize: "11px", color: C.textMuted }}>{s.categoryName}</div>
              </div>
            ))}
          </div>
        </div>
      ))}

      <div style={{ marginTop: "32px", background: C.goldGlow, border: `1px solid ${C.gold}30`, borderRadius: "10px", padding: "20px" }}>
        <div style={{ fontFamily: ff.display, fontSize: "16px", color: C.gold, marginBottom: "12px" }}>Infra compartilhada</div>
        <div style={{ display: "grid", gridTemplateColumns: "repeat(auto-fill, minmax(200px, 1fr))", gap: "12px" }}>
          {[
            { name: "PostgreSQL 16", role: "1 instance por serviço" },
            { name: "Redis 7 Cluster", role: "Cache, rate limiting, idempotência" },
            { name: "Apache Kafka", role: "Event streaming, audit trail, sagas" },
            { name: "Apache Pulsar", role: "Notificações, dead letter" },
            { name: "Kubernetes (EKS)", role: "Orquestração, auto-scaling" },
            { name: "OpenTelemetry", role: "Traces, métricas, logs" },
            { name: "Vault / HSM", role: "Secrets, PAN cartão" },
            { name: "S3", role: "Docs, comprovantes, backups" },
          ].map((infra, i) => (
            <div key={i}>
              <div style={{ fontFamily: ff.mono, fontSize: "12px", color: C.text, marginBottom: "3px" }}>{infra.name}</div>
              <div style={{ fontFamily: ff.body, fontSize: "11px", color: C.textMuted, lineHeight: 1.4 }}>{infra.role}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}

// ─── Main App ───

const TABS = [
  { id: "services", label: "Serviços por Domínio" },
  { id: "deps", label: "Mapa de Dependências" },
  { id: "summary", label: "Resumo & Prioridade" },
];

export default function App() {
  const [activeTab, setActiveTab] = useState("services");
  const [expandedService, setExpandedService] = useState(null);
  const [filterPriority, setFilterPriority] = useState("ALL");
  const [mounted, setMounted] = useState(false);

  useEffect(() => setMounted(true), []);

  const filteredCategories =
    filterPriority === "ALL"
      ? SERVICE_CATEGORIES
      : SERVICE_CATEGORIES
          .map((c) => ({ ...c, services: c.services.filter((s) => s.priority === filterPriority) }))
          .filter((c) => c.services.length > 0);

  return (
    <>
      <style dangerouslySetInnerHTML={{ __html: globalStyles }} />
      <div style={{ minHeight: "100vh", background: C.bg, color: C.text, fontFamily: ff.body }}>

        {/* Header */}
        <header style={{
          padding: "32px 32px 24px", borderBottom: `1px solid ${C.border}`,
          opacity: mounted ? 1 : 0, transform: mounted ? "translateY(0)" : "translateY(-8px)",
          transition: "all 0.5s ease",
        }}>
          <div style={{ display: "flex", alignItems: "baseline", gap: "12px", marginBottom: "6px" }}>
            <h1 style={{ fontFamily: ff.display, fontSize: "28px", fontWeight: 700, color: C.gold, margin: 0 }}>
              Isura Bank
            </h1>
            <span style={{ fontFamily: ff.body, fontSize: "13px", color: C.textMuted, fontStyle: "italic" }}>
              tesouro em Yorubá
            </span>
          </div>
          <p style={{ fontFamily: ff.body, fontSize: "14px", color: C.textSoft, margin: "8px 0 0 0", lineHeight: 1.6 }}>
            Fase 1 — Mapeamento de Serviços
          </p>
        </header>

        {/* Tabs */}
        <nav style={{ display: "flex", padding: "0 32px", borderBottom: `1px solid ${C.border}`, background: C.surface }}>
          {TABS.map((tab) => (
            <button
              key={tab.id}
              onClick={() => setActiveTab(tab.id)}
              style={{
                padding: "14px 20px", background: "transparent", border: "none",
                borderBottom: activeTab === tab.id ? `2px solid ${C.gold}` : "2px solid transparent",
                color: activeTab === tab.id ? C.gold : C.textMuted,
                fontFamily: ff.body, fontSize: "13px", fontWeight: activeTab === tab.id ? 600 : 400,
                cursor: "pointer", transition: "all 0.15s",
              }}
            >
              {tab.label}
            </button>
          ))}
        </nav>

        {/* Content */}
        <main style={{ padding: "24px 32px 60px", maxWidth: "960px", margin: "0 auto" }}>
          {activeTab === "services" && (
            <div>
              <div style={{ display: "flex", gap: "8px", marginBottom: "24px", flexWrap: "wrap" }}>
                {[
                  { key: "ALL", label: "Todos", color: C.textSoft },
                  ...Object.entries(PRIORITIES).map(([k, v]) => ({ key: k, label: v.label, color: v.color })),
                ].map((f) => (
                  <button
                    key={f.key}
                    onClick={() => setFilterPriority(f.key)}
                    style={{
                      padding: "6px 14px", borderRadius: "6px",
                      border: `1px solid ${filterPriority === f.key ? f.color + "66" : C.border}`,
                      background: filterPriority === f.key ? f.color + "15" : "transparent",
                      color: filterPriority === f.key ? f.color : C.textMuted,
                      fontFamily: ff.mono, fontSize: "11px", fontWeight: 600,
                      cursor: "pointer", transition: "all 0.15s",
                    }}
                  >
                    {f.label}
                  </button>
                ))}
              </div>

              {filteredCategories.map((category) => (
                <div key={category.id} style={{ marginBottom: "32px" }}>
                  <div style={{ marginBottom: "14px" }}>
                    <h2 style={{ fontFamily: ff.display, fontSize: "19px", fontWeight: 600, color: category.color, margin: "0 0 4px 0" }}>
                      {category.name}
                    </h2>
                    <p style={{ fontFamily: ff.body, fontSize: "13px", color: C.textMuted, margin: 0 }}>{category.description}</p>
                  </div>
                  {category.services.map((svc) => (
                    <ServiceCard
                      key={svc.id}
                      service={{ ...svc, categoryColor: category.color }}
                      isExpanded={expandedService === svc.id}
                      onToggle={() => setExpandedService(expandedService === svc.id ? null : svc.id)}
                    />
                  ))}
                </div>
              ))}
            </div>
          )}
          {activeTab === "deps" && <DependencyMap />}
          {activeTab === "summary" && <SummaryView />}
        </main>
      </div>
    </>
  );
}
