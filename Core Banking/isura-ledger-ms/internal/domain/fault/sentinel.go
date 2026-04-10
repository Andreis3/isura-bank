package fault

// Erros sentinela por domínio.
//
// Erros sentinela são variáveis globais que representam condições de erro
// conhecidas. Permitem comparação com errors.Is sem depender de strings.
//
// Convenção:
//   - Prefixo "Err" para erros de negócio (esperados)
//   - Sem prefixo para erros de infraestrutura (wrappados com Wrap())
//
// Uso:
//
//   if errors.Is(err, ErrCustomerNotFound) {
//       // tratar 404
//   }

/* --- Account --- */
var (
	ErrAccountNotFound      = &DomainError{Code: CodeNotFound, FriendlyMessage: "account not found"}
	ErrAccountAlreadyExists = &DomainError{Code: CodeConflict, FriendlyMessage: "account already exists"}
	ErrInsufficientFunds    = &DomainError{Code: CodeConflict, FriendlyMessage: "insufficient funds"}
	ErrTransactionNotFound  = &DomainError{Code: CodeNotFound, FriendlyMessage: "transaction not found"}
	ErrInvalidAmount        = &DomainError{Code: CodeBadRequest, FriendlyMessage: "invalid amount"}
	ErrDuplicateTransaction = &DomainError{Code: CodeConflict, FriendlyMessage: "transaction already exists"}
)
